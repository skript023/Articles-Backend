package controller

import (
	"ArticleBackend/config"
	"ArticleBackend/database"
	"ArticleBackend/joaat"
	"ArticleBackend/models"
	"errors"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Email: e}).Find(&user).Error; err != nil {
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
	if err := db.Where(&models.User{Username: u}).Find(&user).Error; err != nil {
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

func Login(res *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(LoginInput)
	var ud UserData

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": "Error on login request",
			"data":    err,
		})
	}

	identity := input.Identity
	pass := input.Password
	user, email, err := new(models.User), new(models.User), *new(error)

	if valid(identity) {
		email, err = getUserByEmail(identity)
		if err != nil {
			return res.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  joaat.Hash("EMAIL_INVALID"),
				"message": "Credential does not match",
				"data":    err,
			})
		}
	} else {
		user, err = getUserByUsername(identity)
		if err != nil {
			return res.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  joaat.Hash("USERNAME_INVALID"),
				"message": "Credential does not match",
				"data":    err,
			})
		}
	}

	if email == nil && user == nil {
		return res.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  joaat.Hash("USER_NOT_FOUND"),
			"message": "Credential does not match",
			"data":    err,
		})
	}

	if email != nil {
		ud = UserData{
			ID:       email.ID,
			Username: email.Username,
			Email:    email.Email,
			Password: email.Password,
		}

	}
	if user != nil {
		ud = UserData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
		}
	}

	if !CheckPasswordHash(pass, ud.Password) {
		return res.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  joaat.Hash("PASSWORD_INVALID"),
			"message": "Credential does not match",
			"data":    nil,
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = ud.Username
	claims["user_id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Env("SECRET")))
	if err != nil {
		return res.SendStatus(fiber.StatusInternalServerError)
	}

	return res.JSON(fiber.Map{
		"status":  joaat.Hash("AUTH_SUCCESS"),
		"message": "Success login",
		"data":    t,
	})
}
