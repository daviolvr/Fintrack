package dto

import (
	"time"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserMeResponse struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type UserUpdateResponse struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserUpdateBalanceResponse struct {
	Balance float64 `json:"balance"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PaginatedCategoriesResponse struct {
	Data       []CategoryResponse `json:"data"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

type TransactionCreateResponse struct {
	CategoryID  uint      `json:"category_id"`
	Type        string    `json:"type"` // "income" ou "expense"
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type TransactionResponse struct {
	CategoryID  uint      `json:"category_id"`
	Type        string    `json:"type"` // "income" ou "expense"
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaginatedTransactionResponse struct {
	Data       []TransactionResponse `json:"data"`
	Total      int                   `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"totalPages"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
