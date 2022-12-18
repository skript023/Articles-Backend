package models

import "time"

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Fullname  string `json:"fullname" gorm:"not null"`
	Username  string `json:"username" gorm:"unique_index;not null"`
	Email     string `json:"email" gorm:"unique_index;not null"`
	Password  string `json:"password" gorm:"not null"`
	Avatar    string `json:"avatar" gorm:"not null"`
	Status    string `json:"status" gorm:"default:unverified"`
	CreatedAt time.Time
}
