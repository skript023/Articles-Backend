package controller

import (
	"ArticleBackend/config"
	"ArticleBackend/database"
	"ArticleBackend/joaat"
	"ArticleBackend/models"
	"encoding/gob"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var validate = validator.New()
var store_session = session.New()

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Find(&user, "email = ?", e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Find(&user, "username = ?", u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func authAttempt(credentials interface{}) (*models.User, error) {
	data := credentials.(fiber.Map)
	user, email, err := new(models.User), new(models.User), *new(error)

	if valid(data["identity"].(string)) {
		email, err = getUserByEmail(data["identity"].(string))
		if err != nil {
			return nil, err
		}
	} else {
		user, err = getUserByUsername(data["identity"].(string))
		if err != nil {
			return nil, err
		}
	}

	if email != nil && CheckPasswordHash(data["password"].(string), email.Password) {
		database.DB.Find(&email.Role, "role = ?", email.RoleID)
		return email, nil
	}

	if user != nil && CheckPasswordHash(data["password"].(string), user.Password) {
		database.DB.Find(&user.Role, "role = ?", user.RoleID)
		return user, nil
	}

	return nil, errors.New("credentials not valid")
}

func Login(res *fiber.Ctx) error {
	gob.Register(fiber.Map{})
	gob.Register(User{})
	type LoginInput struct {
		Identity string `json:"identity" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	input := new(LoginInput)

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": "Error on login request",
			"data":    "NO_DATA_AQUIRED",
			"user":    "NO_ACCESS_ACQUIRED",
		})
	}

	if err := validate.Struct(input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": "Credential does not match",
			"data":    "NO_DATA_AQUIRED",
			"user":    "NO_ACCESS_ACQUIRED",
		})
	}

	credential := fiber.Map{
		"identity": input.Identity,
		"password": input.Password,
	}

	user, err := authAttempt(credential)

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("AUTH_FAILED"),
			"message": "Credential does not match",
			"data":    "NO_DATA_AQUIRED",
			"user":    "NO_ACCESS_ACQUIRED",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(config.Env("SECRET")))
	if err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("GENERATE_TOKEN_FAILED"),
			"message": fmt.Sprintf("Error : %v", err),
			"token":   "NO_ACCESS_ACQUIRED",
		})
	}

	avatar := res.BaseURL() + "/api/v1/user/avatar/" + user.Avatar
	user.Avatar = avatar

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash(user.Username + "AUTH_SUCCESS" + user.Email),
		"message": "Success login",
		"token":   t,
		"user":    responseUser(*user),
	})
}

func SessionData(res *fiber.Ctx) error {
	current, _ := store_session.Get(res)
	user := current.Get("user")
	if user == nil {
		return res.Status(500).JSON(fiber.Map{
			"status":  joaat.Hash("USER_INVALID"),
			"message": "Unable get data from server",
		})
	}
	return res.Status(200).JSON(user.(fiber.Map))
}

func AuthUser(res *fiber.Ctx) error {
	type Input struct {
		Identity string `json:"identity" validate:"required"`
		Code     uint32 `json:"code" validate:"required"`
	}

	var input Input

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": "Error on login request",
			"data":    "NO_DATA_AQUIRED",
			"user":    "NO_ACCESS_ACQUIRED",
		})
	}

	if err := validate.Struct(input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("REQUIREMENT_DOES_NOT_MATCH"),
			"message": "Refresh Failed",
			"data":    "NO_DATA_AQUIRED",
			"user":    "NO_ACCESS_ACQUIRED",
		})
	}

	user, err := getUserByUsername(input.Identity)

	if err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("GENERATE_TOKEN_FAILED"),
			"message": fmt.Sprintf("Error : %v", err),
			"token":   "NO_ACCESS_ACQUIRED",
			"user":    "NO_ACCESS_ACQUIRED",
		})
	}

	code := joaat.Hash(user.Username + "AUTH_SUCCESS" + user.Email)

	if input.Code != code {
		return res.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  joaat.Hash("GENERATE_TOKEN_FAILED"),
			"message": "User not authentication",
			"token":   "NO_ACCESS_ACQUIRED",
			"user":    "NO_ACCESS_ACQUIRED",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(config.Env("SECRET")))
	if err != nil {
		return res.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  joaat.Hash("GENERATE_TOKEN_FAILED"),
			"message": fmt.Sprintf("Error : %v", err),
			"token":   "NO_ACCESS_ACQUIRED",
			"user":    "NO_ACCESS_ACQUIRED",
		})
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("REFRESH_TOKEN_SUCCESS"),
		"message": "Refresh success",
		"token":   t,
		"user":    responseUser(*user),
	})
}

func GetIdFromToken(res *fiber.Ctx) uint {
	header := res.Request().Header.Peek("Authorization")
	split := strings.Split(string(header), "Bearer ")
	token := split[1]

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(tokens *jwt.Token) (interface{}, error) {
		return []byte(config.Env("SECRET")), nil
	})

	result := claims["user_id"].(float64)

	return uint(result)
}
