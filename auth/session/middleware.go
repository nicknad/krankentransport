package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nicknad/krankentransport/dataaccess"
	"golang.org/x/crypto/bcrypt"
)

const sessionIdLength = 32

var sessionManager = NewSessionManager()

func ClearExpiredSessions() {
	sessionManager.ClearExpiredSessions()
}

func ApiAuthMiddleware(c *fiber.Ctx) error {
	if !IsAuthenticated(c) {
		ClearSession(c)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	SetRole(c)

	return c.Next()
}

func AssertAdminMiddleWare(c *fiber.Ctx) error {
	SetRole(c)
	isAdmin, ok := c.Locals("l_isAdmin").(bool)

	if ok && isAdmin {
		return c.Next()
	}

	return c.SendStatus(fiber.StatusForbidden)
}

func SetRole(c *fiber.Ctx) {
	login, ok := c.Locals("l_login").(string)

	if ok {
		user, err := dataaccess.GetUser(login)
		if err != nil {
			c.Locals("l_isAdmin", false)
			return
		}

		c.Locals("l_isAdmin", user.Admin)
		return
	}

	c.Locals("l_isAdmin", false)
}

func AssertAuthenticatedMiddleware(c *fiber.Ctx) error {
	if !IsAuthenticated(c) {
		ClearSession(c)
		c.Set("HX-Redirect", "/login")
		return c.Redirect("/login")
	}
	return c.Next()
}

func IsAuthenticated(c *fiber.Ctx) bool {
	authCookie := c.Cookies("fahrten_auth_token")
	if authCookie != "" {

		checkedSession, exists := sessionManager.CheckSession(authCookie)
		if !exists {
			return false
		}

		c.Locals("l_login", checkedSession.Login)

		return true
	}

	// if c.Path() != "/login" {
	// 	return false
	// }

	authHeader := c.Get("Authorization")
	if authHeader != "" {
		authAttempt := strings.TrimPrefix(authHeader, "Basic ")
		s := strings.SplitN(authAttempt, " ", 2)
		if len(s) != 2 {
			return false
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			return false
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			return false
		}

		valid := IsValidPassword(pair[0], pair[1])
		if valid {
			SetSession(c, pair[0])
		}

		return valid
	}

	return false
}

func GetUserSessionId(c *fiber.Ctx) string {
	return c.Cookies("fahrten_auth_token")
}

func SetSession(c *fiber.Ctx, login string) string {
	newSessionId := generateSessionId()
	sessionManager.CreateSession(login, newSessionId)
	c.Cookie(&fiber.Cookie{
		Name:     "fahrten_auth_token",
		Value:    newSessionId,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "strict",
		MaxAge:   3600,
	})
	return newSessionId
}

func ClearSession(c *fiber.Ctx) {
	token := c.Cookies("fahrten_auth_token")
	sessionManager.ClearSession(token)
	c.Cookie(&fiber.Cookie{
		Name:     "fahrten_auth_token",
		Value:    "",
		HTTPOnly: true,
		// Secure:   true,
		SameSite: "strict",
		MaxAge:   3600,
	})
}

func generateSessionId() string {
	bytes := make([]byte, sessionIdLength)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func IsValidPassword(login string, password string) bool {
	user, err := dataaccess.GetUser(login)

	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
