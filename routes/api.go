package routes

import (
	"ArticleBackend/config"
	"ArticleBackend/controller"
	"ArticleBackend/middleware"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func SetupRouteApi(app *fiber.App) {
	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	user := api.Group("/user")
	post := api.Group("/post")

	api.Get("/test", middleware.Protected(), func(c *fiber.Ctx) error {
		header := c.Request().Header.Peek("Authorization")
		split := strings.Split(string(header), "Bearer ")
		tokens := split[1]

		claims := jwt.MapClaims{}
		jwt.ParseWithClaims(tokens, claims, func(tokens *jwt.Token) (interface{}, error) {
			return []byte(config.Env("SECRET")), nil
		})

		fmt.Printf("value: %v", claims["user_id"])
		return c.SendStatus(200)
	})
	auth.Post("/login", controller.Login)

	post.Post("/create", middleware.Protected(), controller.CreatePost)
	post.Get("/all", middleware.Protected(), controller.GetPosts)
	post.Patch("/update/:id", middleware.Protected(), controller.UpdatePost)
	post.Delete("/delete/:id", middleware.Protected(), controller.DeletePost)

	user.Post("/create", controller.CreateUser)
	user.Patch("/update/:id", middleware.Protected(), controller.UpdateUser)
	user.Delete("/delete/:id", middleware.Protected(), controller.DeleteUser)
}
