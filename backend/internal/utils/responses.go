package utils

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
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PaginatedCategoriesResponse struct {
	Data       []CategoryResponse `json:"data"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

type CategoryListResponse struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Name   int64 `json:"name"`
}

type TransactionGetResponse struct {
	CategoryID  int64     `json:"category_id"`
	Type        string    `json:"type"` // "income" ou "expense"
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
