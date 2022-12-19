package controller

import (
	"ArticleBackend/database"
	"ArticleBackend/joaat"
	"ArticleBackend/models"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func responseUser(user models.User) User {
	return User{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}
}

func CreateUser(res *fiber.Ctx) error {
	var user models.User

	if err := res.BodyParser(&user); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	user.Password = hash
	database.DB.Create(&user)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("CREATE_USER_SUCCESS"),
		"message": "Registeration Success.",
	})
}

func GetUsers(res *fiber.Ctx) error {
	users := []models.User{}

	database.DB.Find(&users)
	response_users := []User{}
	for _, user := range users {
		response_user := responseUser(user)
		response_users = append(response_users, response_user)
	}

	return res.Status(fiber.StatusOK).JSON(response_users)
}

func findUser(id int, user *models.User) error {
	database.DB.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

func GetUser(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var user models.User

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_EXIST"),
			"message": "Please, ensure that id is an integer",
			"data":    fiber.Map{},
		})
	}

	if err := findUser(id, &user); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
			"data":    fiber.Map{},
		})
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("USER_RETRIEVED"),
		"message": "User information retrieved successfully",
		"data":    responseUser(user),
	})
}

func getUserData(id int) User {
	var user models.User

	if err := findUser(id, &user); err != nil {
		return User{}
	}

	return responseUser(user)
}

func UpdateUser(res *fiber.Ctx) error {
	type Update struct {
		Username string `json:"username"`
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}

	id, err := res.ParamsInt("id")

	var user models.User

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_EXIST"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findUser(id, &user); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	var update Update

	if err := res.BodyParser(&update); err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	user.Username = update.Username
	user.Fullname = update.Fullname
	user.Email = update.Email
	user.Avatar = update.Avatar

	database.DB.Save(&user)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("UPDATE_USER_SUCCESS"),
		"message": "User information updated successfully",
	})
}

func DeleteUser(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var user models.User

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_EXIST"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findUser(id, &user); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("DELETE_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("DELETE_USER_SUCCESS"),
		"message": "User deleted successfully",
	})
}
