package routes

import (
	"ArticleBackend/controller"
	"ArticleBackend/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func SetupRouteApi(app *fiber.App) {
	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	user := api.Group("/user")
	post := api.Group("/post")
	contact := api.Group("/contact")
	sessions := session.New()

	api.Get("/test", middleware.Protected(), func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		//    set value to the session store
		store.Set("name", "King Windrol")

		name := store.Get("name")
		c.Status(200).JSON(fiber.Map{
			"name": name,
		})
		return store.Save()
	})

	api.Get("/test2", middleware.Protected(), controller.SessionData)

	auth.Post("/login", controller.Login)
	auth.Post("/refresh", controller.AuthUser)

	post.Get("/all", controller.GetPosts)
	post.Get("/:id", controller.GetPost)
	post.Post("/create", middleware.Protected(), controller.CreatePost)
	post.Patch("/update/:id", middleware.Protected(), controller.UpdatePost)
	post.Delete("/delete/:id", middleware.Protected(), controller.DeletePost)
	post.Get("/read/:title", controller.ReadPost)

	user.Post("/create", controller.CreateUser)
	user.Post("/all", middleware.Protected(), controller.GetUsers)
	user.Patch("/update/:id", middleware.Protected(), controller.UpdateUser)
	user.Delete("/delete/:id", middleware.Protected(), controller.DeleteUser)
	user.Get("/counts", middleware.Protected(), controller.UsersCount)
	user.Get("/avatar/:name", func(c *fiber.Ctx) error {
		image := c.Params("name")
		return c.SendFile(fmt.Sprintf("./public/users/avatar/%v", image))
	})

	contact.Post("/create", controller.CreateContact)
	contact.Get("/all", middleware.Protected(), controller.GetContacts)
	contact.Get("/counts", middleware.Protected(), controller.ContactCount)
}
