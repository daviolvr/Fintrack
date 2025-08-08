package repository

import (
	"database/sql"

	"github.com/daviolvr/Fintrack/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (first_name, last_name, email, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`

	return db.QueryRow(
		query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.ID)
}
