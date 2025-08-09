package repository

import (
	"database/sql"

	"github.com/daviolvr/Fintrack/internal/models"
)

// Cria categoria para o usuário
func CreateCategory(db *sql.DB, category *models.Category) error {
	query := `INSERT into categories (user_id, name) VALUES ($1, $2) RETURNING id`
	return db.QueryRow(query, category.UserID, category.Name).Scan(&category.ID)
}
