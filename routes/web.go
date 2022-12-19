package routes

import "github.com/gofiber/fiber/v2"

func SetupRouteWeb(app *fiber.App) {
	app.Static("/", "./views")
	app.Static("/password", "./views")

	app.Get("/", func(res *fiber.Ctx) error {
		return res.Render("views/login.html", fiber.Map{
			"Title": "Welcome to API",
		})
	})

	app.Get("/register", func(res *fiber.Ctx) error {
		return res.Render("views/register.html", fiber.Map{
			"Title": "Welcome to API",
		})
	})

	app.Get("/password/reset", func(res *fiber.Ctx) error {
		return res.Render("views/reset-password.html", fiber.Map{
			"Title": "Welcome to API",
		})
	})
}
