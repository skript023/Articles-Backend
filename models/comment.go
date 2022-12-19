package models

import "time"

type Comment struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	PostID    uint   `json:"post_id"`
	Post      Post   `gorm:"foreignKey:PostID"`
	Fullname  string `json:"fullname" gorm:"size:255"`
	Email     string `json:"email" gorm:"size:255"`
	Comment   string `json:"comment"`
	Status    string `json:"status" gorm:"size:255"`
	CreatedAt time.Time
}
