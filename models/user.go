package models

import "time"

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Fullname  string `json:"fullname" gorm:"not null;size:255"`
	Username  string `json:"username" gorm:"unique;not null;size:255"`
	Email     string `json:"email" gorm:"unique;not null;size:255"`
	Password  string `json:"password" gorm:"not null;size:255"`
	Avatar    string `json:"avatar" gorm:"size:255"`
	Status    string `json:"status" gorm:"default:unverified;size:255"`
	RoleID    uint   `json:"role_id"`
	Role      Role   `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
}
