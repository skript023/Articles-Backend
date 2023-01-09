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
	RoleID   uint   `json:"role_id"`
	Role     Role   `json:"user_role"`
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
		RoleID:   user.RoleID,
		Role:     getRoleData(int(user.RoleID)),
	}
}

func CreateUser(res *fiber.Ctx) error {
	type Input struct {
		Fullname string `json:"fullname" validate:"required"`
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var input Input
	var user models.User

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := validate.Struct(input); err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	hash, err := hashPassword(input.Password)
	if err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	user.RoleID = 1
	user.Fullname = input.Fullname
	user.Username = input.Username
	user.Email = input.Email
	user.Password = hash
	user.Status = "verified"

	result := database.DB.Where("username = ?", user.Username).Or("email = ?", user.Email).FirstOrCreate(&user)

	if result.Error != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", result.Error),
		})
	}

	if result.RowsAffected == 1 {
		return res.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_USER_SUCCESS"),
			"message": "Registeration Success.",
		})
	}

	return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  joaat.Hash("USERNAME_OR_EMAIL_ALREADY_EXIST"),
		"message": "Email or Username already exist",
	})
}

func GetUsers(res *fiber.Ctx) error {
	type Input struct {
		Role uint32 `json:"role" validate:"required"`
	}

	var input = new(Input)
	users := []models.User{}

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("GET_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
			"users":   fiber.Map{},
		})
	}

	if err := validate.Struct(input); err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("GET_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
			"users":   fiber.Map{},
		})
	}

	if input.Role != joaat.Hash("admin") {
		return res.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  joaat.Hash("GET_USER_FAILED"),
			"message": "Restricted Area.",
			"users":   fiber.Map{},
		})
	}

	database.DB.Find(&users)
	response_users := []User{}
	for _, user := range users {
		response_user := responseUser(user)
		response_users = append(response_users, response_user)
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("GET_USER_FAILED"),
		"message": "Users information acquired successfully",
		"users":   response_users,
	})
}

func UsersCount(res *fiber.Ctx) error {
	users := models.User{}
	var counts int64
	database.DB.Model(&users).Count(&counts)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("USERS_COUNT_ACQUIRED"),
		"message": "Users count acquired successfully",
		"users":   counts,
	})
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
			"status":  joaat.Hash("ENSURE_ID_VALID"),
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
			"status":  joaat.Hash("ENSURE_ID_VALID"),
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
			"status":  joaat.Hash("ENSURE_ID_VALID"),
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
