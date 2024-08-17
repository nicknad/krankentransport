package admin

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	auth "github.com/nicknad/krankentransport/auth/session"
	"github.com/nicknad/krankentransport/dataaccess"
	"github.com/nicknad/krankentransport/utils"
)

func RegisterAdmin(app *fiber.App) {
	app.Get("/admin", auth.AssertAuthenticatedMiddleware, auth.AssertAdminMiddleWare, func(c *fiber.Ctx) error {
		krankenfahrten, err := dataaccess.GetKrankenfahrten()

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		users, err := dataaccess.GetUsers()

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return utils.Render(c, AdminLayout(auth.IsAuthenticated(c), krankenfahrten, users))
	})

	app.Post("/admin/fahrt/create", auth.AssertAuthenticatedMiddleware, auth.AssertAdminMiddleWare, func(c *fiber.Ctx) error {
		desc := c.FormValue("description")

		if desc == "" {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		fahrt, err := dataaccess.CreateKrankenfahrt(desc)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return utils.Render(c, FahrtRow(fahrt))

	})

	app.Post("/admin/fahrt/reopen/:id", auth.AssertAuthenticatedMiddleware, auth.AssertAdminMiddleWare, func(c *fiber.Ctx) error {
		s := c.Params("id")
		if s == "" {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		id, err := strconv.Atoi(s)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		fahrt, err := dataaccess.GetKrankenfahrt(id)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		fahrt.AcceptedAt = nil
		fahrt.AcceptedByLogin = nil

		err = dataaccess.UndoAcceptKrankenfahrt(fahrt)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return utils.Render(c, FahrtRow(fahrt))

	})

	app.Post("/admin/user/create", auth.AssertAuthenticatedMiddleware, auth.AssertAdminMiddleWare, func(c *fiber.Ctx) error {
		password := c.FormValue("password")
		login := c.FormValue("user")
		admin := c.FormValue("admincheck")

		user, err := dataaccess.CreateUser(login, password, admin == "true")

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return utils.Render(c, UserRow(user))
	})

	app.Delete("/admin/fahrt/delete/:id", auth.AssertAuthenticatedMiddleware, auth.AssertAdminMiddleWare, func(c *fiber.Ctx) error {
		s := c.Params("id")

		if s == "" {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		i, err := strconv.Atoi(s)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		err = dataaccess.DeleteKrankenfahrt(i)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Send([]byte(""))
	})

	app.Delete("/admin/user/delete/:id", auth.AssertAuthenticatedMiddleware, auth.AssertAdminMiddleWare, func(c *fiber.Ctx) error {
		s := c.Params("id")

		if s == "" {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		i, err := strconv.Atoi(s)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		usr, err := dataaccess.GetUserById(i)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		dataaccess.DeleteUser(usr.Login)

		return c.Send([]byte(""))
	})
}
