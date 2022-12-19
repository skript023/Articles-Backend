package models

import "time"

type Contact struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Fullname  string `json:"fullname" gorm:"size:255"`
	Email     string `json:"email" gorm:"size:255"`
	Message   string `json:"message"`
	CreatedAt time.Time
}
