package main

import (
	"ArticleBackend/database"
	"ArticleBackend/routes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Check interface {
	HasData() bool
}

type Data struct {
	Exist bool
}

func (data Data) HasData() bool {
	return data.Exist
}

func IsDataExist(check Check) {
	fmt.Println(check.HasData())
}

func main() {
	database.ConnectDatabase()

	app := fiber.New()
	app.Use(cors.New())

	routes.SetupRoute(app)

	app.Listen(":8000")
}
