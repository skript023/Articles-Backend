package controller

import (
	"ArticleBackend/database"
	"ArticleBackend/joaat"
	"ArticleBackend/models"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Role struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

func responseRole(category models.Role) Role {
	return Role{
		ID:   category.ID,
		Role: category.Role,
	}
}

func CreateRole(res *fiber.Ctx) error {
	type Input struct {
		Role string `json:"role" validate:"required"`
	}

	var input Input
	var role models.Role

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_ROLE_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := validate.Struct(input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("REQUIREMENT_DOES_NOT_MATCH"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	role.Role = input.Role

	database.DB.Create(&role)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("CREATE_ROLE_SUCCESS"),
		"message": "Role created successfully.",
	})
}

func createRole(name string) (uint, error) {
	var role models.Role

	role.Role = name

	result := database.DB.Create(&role)

	if result.Error != nil {
		return 0, result.Error
	}

	return role.ID, nil
}

func GetRoles(res *fiber.Ctx) error {
	roles := []models.Role{}

	database.DB.Find(&roles)
	response_roles := []Role{}
	for _, role := range roles {
		response_role := responseRole(role)
		response_roles = append(response_roles, response_role)
	}

	return res.Status(fiber.StatusOK).JSON(response_roles)
}

func findRole(id int, role *models.Role) error {
	database.DB.Find(&role, "id = ?", id)
	if role.ID == 0 {
		return errors.New("Role does not exist")
	}

	return nil
}

func GetRole(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var role models.Role

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findRole(id, &role); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_ROLE_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(role)
}

func getRoleData(id int) Role {
	var role models.Role

	if err := findRole(id, &role); err != nil {
		return Role{}
	}

	return responseRole(role)
}

func UpdateRole(res *fiber.Ctx) error {

	id, err := res.ParamsInt("id")

	var role models.Role

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findRole(id, &role); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_ROLE_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	var update Role

	if err := res.BodyParser(&update); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	role.ID = update.ID
	role.Role = update.Role

	database.DB.Save(&role)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("UPDATE_ROLE_SUCCESS"),
		"message": "Role updated successfully",
	})
}

func DeleteRole(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var role models.Role

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findRole(id, &role); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_ROLE_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		return res.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
			"status":  joaat.Hash("DELETE_ROLE_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("DELETE_ROLE_SUCCESS"),
		"message": "Role deleted successfully",
	})
}
