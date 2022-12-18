package models

import "time"

type Post struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Post string `json:"post"`

	CreatedAt time.Time
}
