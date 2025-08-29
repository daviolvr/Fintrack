package models

import (
	"time"
)

type User struct {
	ID           uint       `gorm:"primaryKey"`
	FirstName    string     `gorm:"not null;size:100" json:"first_name"`
	LastName     string     `gorm:"not null;size:100" json:"last_name"`
	Email        string     `gorm:"unique;not null;size:100" json:"email"`
	Password     string     `gorm:"not null;size:255" json:"-"`
	Balance      float64    `gorm:"default:0" json:"balance"`
	FailedLogins uint       `gorm:"default:0" json:"failed_logins"`
	LockedUntil  *time.Time `json:"locked_until,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Categoria da transação (ex: Alimentação, Transporte)
type Category struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null" json:"user_id"`
	User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Name   string `gorm:"not null;size:50" json:"name"`
}

type Transaction struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	CategoryID  uint      `gorm:"not null" json:"category_id"`
	Category    Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	Type        string    `gorm:"not null;size:20" json:"type"` // "income" ou "expense"
	Amount      float64   `gorm:"not null" json:"amount"`
	Description string    `gorm:"size:255" json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
