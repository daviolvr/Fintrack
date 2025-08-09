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

// Busca categorias do usuário
func FindCategoriesByUser(db *sql.DB, userID int64) ([]models.Category, error) {
	query := `SELECT id, user_id, name FROM categories WHERE user_id = $1 ORDER BY name`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
