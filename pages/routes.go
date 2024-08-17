package pages

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicknad/krankentransport/pages/admin"
	"github.com/nicknad/krankentransport/pages/home"
	"github.com/nicknad/krankentransport/pages/login"
)

func RegisterRoutes(app *fiber.App) {
	home.RegisterHome(app)
	login.RegisterLogin(app)
	admin.RegisterAdmin(app)
}
