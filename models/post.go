package models

import "time"

type Post struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	OwnerID    uint     `json:"owner_id"`
	Owner      User     `gorm:"foreignKey:OwnerID"`
	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID"`
	PostTitle  string   `json:"post_title" gorm:"size:255"`
	Post       string   `json:"post"`
	PostSlug   string   `json:"post_slug" gorm:"size:255"`
	PostImage  string   `json:"post_image" gorm:"size:255"`
	PostStatus string   `json:"post_status" gorm:"size:128"`
	CreatedAt  time.Time
}
