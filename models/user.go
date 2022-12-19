package models

import "time"

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Fullname  string `json:"fullname" gorm:"not null;size:255"`
	Username  string `json:"username" gorm:"unique_index;not null;size:255"`
	Email     string `json:"email" gorm:"unique_index;not null;size:255"`
	Password  string `json:"password" gorm:"not null;size:255"`
	Avatar    string `json:"avatar" gorm:"not null;size:255"`
	Status    string `json:"status" gorm:"default:unverified;size:255"`
	CreatedAt time.Time
}
