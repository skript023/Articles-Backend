package main

import (
	"ArticleBackend/database"
	"ArticleBackend/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

var ConfigDefault = csrf.Config{
	KeyLookup:      "header:X-Csrf-Token",
	CookieName:     "csrf_",
	CookieSameSite: "Strict",
	Expiration:     1 * time.Hour,
	KeyGenerator:   utils.UUID,
}

func main() {
	database.ConnectDatabase()

	app := fiber.New()
	app.Use(cors.New())

	routes.SetupRouteWeb(app)
	routes.SetupRouteApi(app)

	app.Listen(":8000")
}
