package routes

import (
	"ArticleBackend/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoute(app *fiber.App) {
	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	user := api.Group("/user")

	api.Get("/welcome", func(res *fiber.Ctx) error {
		return res.SendString("Welcome to my API")
	})

	auth.Post("/login", controller.Login)

	user.Post("/create", controller.CreateUser)
	user.Patch("/update/:id", controller.UpdateUser)
	user.Delete("/delete/:id", controller.DeleteUser)
}
