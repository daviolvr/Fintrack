package models

import (
	"time"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserUpdateInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserChangePassword struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type CategoryCreateUpdateResponse struct {
	Name string `json:"name"`
}

type CategoryListResponse struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Name   int64 `json:"name"`
}
