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

type TransactionCreateResponse struct {
	CategoryID  int64   `json:"category_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type TransactionListResponse struct {
	CategoryID  int64     `json:"category_id"`
	Type        string    `json:"type"` // "income" ou "expense"
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TransactionUpdateResponse struct {
	CategoryID  int64   `json:"category_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type BalanceUpdateInput struct {
	Balance float64 `json:"balance"`
}
