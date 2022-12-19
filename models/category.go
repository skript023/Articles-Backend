package models

type Category struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	CategoryName string `json:"category_name" gorm:"size:255"`
}
