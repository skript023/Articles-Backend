package controller

import (
	"ArticleBackend/database"
	"ArticleBackend/joaat"
	"ArticleBackend/models"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Comment struct {
	ID        uint   `json:"id"`
	PostID    uint   `json:"post_id"`
	Post      Post   `json:"post"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	Comment   string `json:"comment"`
	Status    string `json:"status"`
	CreatedAt time.Time
}

type CommentPost struct {
	ID       uint   `json:"id"`
	PostID   uint   `json:"post_id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Comment  string `json:"comment"`
	Status   string `json:"status"`
}

func responseComment(comment models.Comment) Comment {
	return Comment{
		ID:       comment.ID,
		PostID:   comment.PostID,
		Post:     getPostData(int(comment.PostID)),
		Fullname: comment.Fullname,
		Email:    comment.Email,
		Comment:  comment.Comment,
		Status:   comment.Status,
	}
}

func responsePostComment(comment models.Comment) CommentPost {
	return CommentPost{
		ID:       comment.ID,
		PostID:   comment.PostID,
		Fullname: comment.Fullname,
		Email:    comment.Email,
		Comment:  comment.Comment,
		Status:   comment.Status,
	}
}

func CreateComment(res *fiber.Ctx) error {
	var comment models.Comment

	if err := res.BodyParser(&comment); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	database.DB.Create(&comment)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("CREATE_USER_SUCCESS"),
		"message": "Registeration Success.",
	})
}

func GetComments(res *fiber.Ctx) error {
	comments := []models.Comment{}

	database.DB.Find(&comments)
	response_comments := []Comment{}
	for _, comment := range comments {
		response_comment := responseComment(comment)
		response_comments = append(response_comments, response_comment)
	}

	return res.Status(fiber.StatusOK).JSON(response_comments)
}

func findComment(id int, comment *models.Comment) error {
	database.DB.Find(&comment, "id = ?", id)
	if comment.ID == 0 {
		return errors.New("Post does not exist")
	}

	return nil
}

func GetComment(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var comment models.Comment

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findComment(id, &comment); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_COMMENT_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(comment)
}

func findCommentPost(id uint, comment *models.Comment) error {
	database.DB.Find(&comment, "post_id = ?", id)
	if comment.ID == 0 {
		return errors.New("Post does not exist")
	}

	return nil
}

func getCommentData(id uint) CommentPost {
	var comments models.Comment

	if err := findCommentPost(id, &comments); err != nil {
		return CommentPost{}
	}

	return responsePostComment(comments)
}

func UpdateComment(res *fiber.Ctx) error {
	type Update struct {
		Comment string `json:"comment"`
	}

	id, err := res.ParamsInt("id")

	var user models.Comment

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findComment(id, &user); err != nil {
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

	user.Comment = update.Comment

	database.DB.Save(&user)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("UPDATE_USER_SUCCESS"),
		"message": "User information updated successfully",
	})
}

func DeleteComment(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var comment models.Comment

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findComment(id, &comment); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_USER_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
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
