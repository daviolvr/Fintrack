package repository

import (
	"database/sql"

	"github.com/daviolvr/Fintrack/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (name, email, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`

	return db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)
}
