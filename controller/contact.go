package controller

import "time"

type Contact struct {
	ID        uint   `json:"id"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	Message   string `json:"message"`
	CreatedAt time.Time
}
