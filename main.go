package main

import (
	"ArticleBackend/database"
	joaat "ArticleBackend/joaat"
	"ArticleBackend/routes"
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	adder := joaat.Hash("Adder")
	fmt.Println(adder)
	var check_existance Data
	check_existance.Exist = true
	IsDataExist(check_existance)

	app := fiber.New()

	routes.InitializeRoute(app)

	app.Listen(":8000")
}
