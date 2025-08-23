package models

import (
	"time"
)

type User struct {
	ID           uint       `json:"id" db:"id"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"-" db:"password_hash"`
	Balance      float64    `json:"balance" db:"balance"`
	FailedLogins uint       `json:"failed_logins" db:"failed_logins"`
	LockedUntil  *time.Time `json:"locked_until,omitempty" db:"locked_until"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// Categoria da transação (ex: Alimentação, Transporte)
type Category struct {
	ID     uint   `json:"id" db:"id"`
	UserID uint   `json:"user_id" db:"user_id"` // Categoria pode ser custom do Usuário
	Name   string `json:"name" db:"name"`
}

type Transaction struct {
	ID          uint      `json:"id" db:"id"`
	UserID      uint      `json:"user_id" db:"user_id"`
	CategoryID  uint      `json:"category_id" db:"category_id"`
	Type        string    `json:"type" db:"type"` // "income" ou "expense"
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description,omitempty" db:"description"`
	Date        time.Time `json:"date" db:"date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
