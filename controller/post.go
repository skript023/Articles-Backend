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
	PostTitle  string `json:"post_title"`
	Post       string `json:"post"`
	PostSlug   string `json:"post_slug"`
	PostImage  string `json:"post_image"`
	PostStatus string `json:"post_status"`
	CreatedAt  time.Time
	Owner      User        `json:"owner"`
	Category   Category    `json:"category"`
	Comment    CommentPost `json:"comments"`
}

func responsePost(post models.Post) Post {
	return Post{
		ID:         post.ID,
		PostTitle:  post.PostTitle,
		Post:       post.Post,
		PostSlug:   post.PostSlug,
		PostImage:  post.PostImage,
		PostStatus: post.PostStatus,
		CreatedAt:  post.CreatedAt,
		Owner:      getUserData(int(post.OwnerID)),
		Category:   getCategoryData(int(post.CategoryID)),
		Comment:    getCommentData(post.ID),
	}
}

func CreatePost(res *fiber.Ctx) error {
	type Input struct {
		Title    string `json:"title" validate:"required"`
		Post     string `json:"post" validate:"required"`
		Category string `json:"category" validate:"required"`
		Image    string `json:"image"`
	}

	var input Input
	var post models.Post
	var category models.Category
	var created_category_id uint

	if err := res.BodyParser(&input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	if err := validate.Struct(input); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("REQUIREMENT_DOES_NOT_MATCH"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	user_id := GetIdFromToken(res)

	if user_id == 0 {
		return res.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  joaat.Hash("USER_OWNER_ID_INVALID"),
			"message": "Forbidden, user are invalid",
		})
	}

	result := database.DB.Find(&category, "category_name = ?", input.Category)
	created_category_id = category.ID

	if result.Error != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("CREATE_CATEGORY_FAILED"),
			"message": fmt.Sprintf("Error : %s", result.Error),
		})
	}

	if created_category_id == 0 {
		var err error
		if created_category_id, err = createCategory(input.Category); err != nil {
			return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  joaat.Hash("CREATE_CATEGORY_FAILED"),
				"message": fmt.Sprintf("Error : %s", err.Error()),
			})
		}
	}

	file, err := res.FormFile("image")

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("FILE_DOES_NOT_VALID"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	post.PostTitle = input.Title
	post.Post = input.Post
	post.CategoryID = created_category_id
	post.OwnerID = user_id
	post.PostSlug = slug.Make(input.Title)
	post.PostStatus = "draft"

	var filename string

	if file != nil {
		filename = post.PostSlug
		if err := res.SaveFile(file, fmt.Sprintf("./public/uploads/post/%s.jpg", post.PostSlug)); err != nil {
			return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  joaat.Hash("FILE_UPLOAD_FAILED"),
				"message": fmt.Sprintf("Error : %s", err.Error()),
			})
		}
	}

	post.PostImage = filename

	database.DB.Create(&post)

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("CREATE_POST_SUCCESS"),
		"message": "Post created successfully.",
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

	return res.Status(fiber.StatusOK).JSON(response_posts)
}

func findPost(id int, post *models.Post) error {
	database.DB.Find(&post, "id = ?", id)
	if post.ID == 0 {
		return errors.New("Post does not exist")
	}

	return nil
}

func findTitle(slug string, post *models.Post) error {
	database.DB.Find(&post, "post_slug = ?", slug)
	if post.ID == 0 {
		return errors.New("post does not exist")
	}

	return nil
}

func GetPost(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var post models.Post

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findPost(id, &post); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	return res.Status(fiber.StatusOK).JSON(responsePost(post))
}

func getPostData(id int) Post {
	var post models.Post

	if err := findPost(id, &post); err != nil {
		return Post{}
	}

	return responsePost(post)
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
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
		})
	}

	if err := findPost(id, &post); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
		})
	}

	var update Update

	if err := res.BodyParser(&update); err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("UPDATE_POST_SUCCESS"),
		"message": "Post updated successfully",
	})
}

func DeletePost(res *fiber.Ctx) error {
	id, err := res.ParamsInt("id")

	var post models.Post

	if err != nil {
		return res.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Please, ensure that id is an integer",
			"data":    "NO_DATA_ACQUIRED",
		})
	}

	if err := findPost(id, &post); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
			"data":    "NO_DATA_ACQUIRED",
		})
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		return res.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
			"status":  joaat.Hash("DELETE_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
			"data":    "NO_DATA_ACQUIRED",
		})
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("DELETE_POST_SUCCESS"),
		"message": "Post deleted successfully",
		"data":    "NO_DATA_ACQUIRED",
	})
}

func ReadPost(res *fiber.Ctx) error {
	title := res.Params("title")

	var post models.Post

	if title == "" {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("ENSURE_ID_VALID"),
			"message": "Request cannot be proceed, invalid title",
			"data":    "NO_DATA_ACQUIRED",
		})
	}

	if err := findTitle(title, &post); err != nil {
		return res.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  joaat.Hash("FIND_POST_FAILED"),
			"message": fmt.Sprintf("Error : %s", err.Error()),
			"data":    "NO_DATA_ACQUIRED",
		})
	}

	return res.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  joaat.Hash("READING_POST_SUCCESS"),
		"message": "Reading post successfully",
		"data":    responsePost(post),
	})
}
