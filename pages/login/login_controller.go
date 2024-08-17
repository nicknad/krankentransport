package login

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	auth "github.com/nicknad/krankentransport/auth/session"
	"github.com/nicknad/krankentransport/utils"
)

func RegisterLogin(app *fiber.App) {
	app.Get("/login", func(c *fiber.Ctx) error {
		return utils.Render(c, loginTempl(auth.IsAuthenticated(c)))
	})

	app.Post("/action/login", func(c *fiber.Ctx) error {
		password := c.FormValue("password")
		login := c.FormValue("user")
		isValid := auth.IsValidPassword(login, password)
		if !isValid {
			return utils.Render(c, InvalidPasswordError(), templ.WithStatus(http.StatusUnprocessableEntity))
		}

		auth.SetSession(c, login)
		c.Response().Header.Set("HX-Redirect", "/")

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/action/logout", auth.AssertAuthenticatedMiddleware, func(c *fiber.Ctx) error {

		_, ok := c.Locals("l_login").(string)

		if !ok {
			return utils.Render(c, UnexpectedError(), templ.WithStatus(http.StatusInternalServerError))
		}

		auth.ClearSession(c)
		c.Response().Header.Set("HX-Redirect", "/login")
		c.Set("HX-Redirect", "/login")
		return c.Redirect("/login")
	})

}
