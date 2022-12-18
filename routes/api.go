package routes

import (
	"ArticleBackend/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoute(app *fiber.App) {
	app.Get("/api", func(res *fiber.Ctx) error {
		return res.SendString("Welcome to my API")
	})

	app.Post("/api/v1/user/create", controller.CreateUser)
	app.Delete("/api/v1/user/delete/:id", controller.DeleteUser)
}
