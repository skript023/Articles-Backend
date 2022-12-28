package models

type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Role string `json:"role" gorm:"not null;size:255"`
}
