package home

import (
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	auth "github.com/nicknad/krankentransport/auth/session"
	"github.com/nicknad/krankentransport/dataaccess"
	"github.com/nicknad/krankentransport/utils"
)

// Helper function to get a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// Helper function to get a pointer to a time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}

func RegisterHome(app *fiber.App) {
	app.Get("/", auth.AssertAuthenticatedMiddleware, func(c *fiber.Ctx) error {
		krankenfahrten, err := dataaccess.GetKrankenfahrten()

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return utils.Render(c, HomeLayout(auth.IsAuthenticated(c), krankenfahrten))
	})

	app.Post("/action/fahrt/:id", auth.AssertAuthenticatedMiddleware, func(c *fiber.Ctx) error {
		loginname, ok := c.Locals("l_login").(string)

		if !ok {
			return utils.Render(c, UnexpectedError(), templ.WithStatus(http.StatusInternalServerError))
		}

		usr, err := dataaccess.GetUser(loginname)

		if err != nil {
			return utils.Render(c, UnexpectedError(), templ.WithStatus(http.StatusInternalServerError))
		}

		s := c.Params("id")
		krankenfahrtId, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		fahrt, err := dataaccess.GetKrankenfahrt(krankenfahrtId)
		if err != nil {
			return err
		}

		currTime := time.Now()
		fahrt.AcceptedByLogin = &usr.Login
		fahrt.AcceptedAt = &currTime

		err = dataaccess.UpdateKrankenfahrt(fahrt)
		if err != nil {
			return err
		}

		c.Response().SetStatusCode(fiber.StatusOK)

		return utils.Render(c, FahrtCell(fahrt))
	})
}
