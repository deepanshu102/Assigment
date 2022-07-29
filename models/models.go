package models

import (
	"time"
)

type Users struct {
	User_Id      int       `gorm:"primary_key;auto_increment;not_null"`
	First_Name   string    `json:"first_name"`
	Last_Name    string    `json:"last_name"`
	Password     string    `json:"password"`
	Email        string    `json:"email"`
	User_type    string    `json:"user_type"`
	PhoneNumber  string    `json:"phone"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	DateOfBirth  string    `json:"dob"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
