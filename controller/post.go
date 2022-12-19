package controller

import (
	"ArticleBackend/database"
	"ArticleBackend/joaat"
	"ArticleBackend/models"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
)

type Post struct {
	ID         uint   `json:"id"`
	OwnerID    uint   `json:"owner_id"`
	CategoryID uint   `json:"category_id"`
	PostTitle  string `json:"post_title"`
	Post       string `json:"post"`
	PostSlug   string `json:"post_slug"`
	PostImage  string `json:"post_image"`
	PostStatus string `json:"post_status"`
	CreatedAt  time.Time
}

func responsePost(post models.Post) Post {
	return Post{
		ID:         post.ID,
		OwnerID:    post.OwnerID,
		CategoryID: post.CategoryID,
		PostTitle:  post.PostTitle,
		Post:       post.Post,
		PostSlug:   post.PostSlug,
		PostImage:  post.PostImage,
		PostStatus: post.PostStatus,
		CreatedAt:  post.CreatedAt,
	}
}

func CreatePost(res *fiber.Ctx) error {
	var post models.Post

	if err := res.BodyParser(&post); err != nil {
		return res.Status(400).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	user_id := GetIdFromToken(res)

	if user_id == 0 {
		return res.Status(400).JSON(fiber.Map{
			"status":  joaat.Hash("USER_OWNER_ID_INVALID"),
			"message": "Invalid post owner",
		})
	}
	post.OwnerID = uint(user_id)
	post.PostSlug = slug.Make(post.PostTitle)

	database.DB.Create(&post)

	return res.Status(200).JSON(fiber.Map{
		"status":  joaat.Hash("CREATE_POST_SUCCESS"),
		"message": "Create Post Success.",
	})
}

func GetPosts(res *fiber.Ctx) error {
	posts := []models.Post{}

	database.DB.Find(&posts)
	response_posts := []Post{}
	for _, post := range posts {
		response_post := responsePost(post)
		response_posts = append(response_posts, response_post)
	}

	return res.Status(200).JSON(response_posts)
}

func findPost(id int, post *models.Post) error {
	database.DB.Find(&post, "id = ?", id)
	if post.ID == 0 {
		return errors.New("Post does not exist")
	}

	return nil
}

func GetPost(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var post models.Post

	if err != nil {
		return res.Status(400).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_EXIST"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findPost(id, &post); err != nil {
		return res.Status(404).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(200).JSON(post)
}

func UpdatePost(res *fiber.Ctx) error {
	type Update struct {
		CategoryID uint   `json:"category_id"`
		PostTitle  string `json:"post_title"`
		Post       string `json:"post"`
		PostSlug   string `json:"post_slug"`
		PostImage  string `json:"post_image"`
		PostStatus string `json:"post_status"`
	}

	id, err := res.ParamsInt("id")

	var post models.Post

	if err != nil {
		return res.Status(400).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_EXIST"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findPost(id, &post); err != nil {
		return res.Status(404).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	var update Update

	if err := res.BodyParser(&update); err != nil {
		return res.Status(500).JSON(fiber.Map{
			"status":  joaat.Hash("BODY_PARSING_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	post.CategoryID = update.CategoryID
	post.PostTitle = update.PostTitle
	post.Post = update.Post
	post.PostSlug = update.PostTitle
	post.PostImage = update.PostImage
	post.PostStatus = update.PostStatus

	database.DB.Save(&post)

	return res.Status(200).JSON(fiber.Map{
		"status":  joaat.Hash("UPDATE_POST_SUCCESS"),
		"message": "Post updated successfully",
	})
}

func DeletePost(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var post models.Post

	if err != nil {
		return res.Status(400).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_EXIST"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findPost(id, &post); err != nil {
		return res.Status(404).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		return res.Status(404).JSON(fiber.Map{
			"status":  joaat.Hash("DELETE_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(200).JSON(fiber.Map{
		"status":  joaat.Hash("DELETE_POST_SUCCESS"),
		"message": "Post deleted successfully",
	})
}
