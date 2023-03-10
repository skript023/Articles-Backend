package controller

import (
	"ArticleBackend/database"
	"ArticleBackend/joaat"
	"ArticleBackend/models"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Category struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"category_name"`
}

func responseCategory(category models.Category) Category {
	return Category{
		ID:           category.ID,
		CategoryName: category.CategoryName,
	}
}

func CreateCategory(res *fiber.Ctx) error {
	type Input struct {
		Category_name string `json:"category_name" validate:"required"`
	}

	var input Input
	var category models.Category

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_CATEGORY_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := validate.Struct(input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("REQUIREMENT_DOES_NOT_MATCH"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	category.CategoryName = input.Category_name

	database.DB.Create(&category)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("CREATE_CATEGORY_SUCCESS"),
		"message": "Category created successfully.",
	})
}

func createCategory(name string) (uint, error) {
	var category models.Category

	category.CategoryName = name

	result := database.DB.Create(&category)

	if result.Error != nil {
		return 0, result.Error
	}

	return category.ID, nil
}

func GetCategories(res *fiber.Ctx) error {
	categories := []models.Category{}

	database.DB.Find(&categories)
	response_categories := []Category{}
	for _, category := range categories {
		response_category := responseCategory(category)
		response_categories = append(response_categories, response_category)
	}

	return res.Status(fiber.StatusOK).JSON(response_categories)
}

func findCategory(id int, category *models.Category) error {
	database.DB.Find(&category, "id = ?", id)
	if category.ID == 0 {
		return errors.New("Category does not exist")
	}

	return nil
}

func GetCategory(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var category models.Category

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findCategory(id, &category); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_CATEGORY_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(category)
}

func getCategoryData(id int) Category {
	var category models.Category

	if err := findCategory(id, &category); err != nil {
		return Category{}
	}

	return responseCategory(category)
}

func UpdateCategory(res *fiber.Ctx) error {

	id, err := res.ParamsInt("id")

	var category models.Category

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findCategory(id, &category); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_CATEGORY_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	var update Category

	if err := res.BodyParser(&update); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	category.ID = update.ID
	category.CategoryName = update.CategoryName

	database.DB.Save(&category)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("UPDATE_CATEGORY_SUCCESS"),
		"message": "Category updated successfully",
	})
}

func DeleteCategory(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var category models.Category

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findCategory(id, &category); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_CATEGORY_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		return res.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
			"status":  joaat.Hash("DELETE_CATEGORY_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("DELETE_CATEGORY_SUCCESS"),
		"message": "Category deleted successfully",
	})
}
