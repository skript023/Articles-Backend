package controller

import (
	"ArticleBackend/database"
	"ArticleBackend/joaat"
	"ArticleBackend/models"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Contact struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

func responseContact(contact models.Contact) Contact {
	return Contact{
		ID:       contact.ID,
		Fullname: contact.Fullname,
		Email:    contact.Email,
		Message:  contact.Message,
	}
}

func CreateContact(res *fiber.Ctx) error {
	type Input struct {
		Fullname string `json:"fullname" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Message  string `json:"message" validate:"required"`
	}

	var input Input
	var contact models.Contact

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_COMMENT_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := validate.Struct(input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("REQUIREMENT_DOES_NOT_MATCH"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	contact.Fullname = input.Fullname
	contact.Email = input.Email
	contact.Message = input.Message

	database.DB.Create(&contact)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("CREATE_COMMENT_SUCCESS"),
		"message": "Your message sent successfully.",
	})
}

func GetContacts(res *fiber.Ctx) error {
	contacts := []models.Contact{}

	database.DB.Find(&contacts)
	response_contacts := []Contact{}
	for _, contact := range contacts {
		response_contact := responseContact(contact)
		response_contacts = append(response_contacts, response_contact)
	}

	return res.Status(fiber.StatusOK).JSON(response_contacts)
}

func findContact(id int, contact *models.Contact) error {
	database.DB.Find(&contact, "id = ?", id)
	if contact.ID == 0 {
		return errors.New("Contact does not exist")
	}

	return nil
}

func GetContact(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var contact models.Contact

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findContact(id, &contact); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_COMMENT_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(contact)
}

func UpdateContact(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var contact models.Contact

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findContact(id, &contact); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_CONTACT_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := res.BodyParser(&contact); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := validate.Struct(contact); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("REQUIREMENT_DOES_NOT_MATCH"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	database.DB.Save(&contact)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("UPDATE_CONTACT_SUCCESS"),
		"message": "Contact updated successfully",
	})
}

func DeleteContact(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var contact models.Contact

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findContact(id, &contact); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_CONTACT_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := database.DB.Delete(&contact).Error; err != nil {
		return res.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
			"status":  joaat.Hash("DELETE_CONTACT_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("DELETE_CONTACT_SUCCESS"),
		"message": "Contact deleted successfully",
	})
}
