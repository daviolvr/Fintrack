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

// Atualiza categoria pelo ID e user_id
func UpdateCategory(db *sql.DB, category *models.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2 AND user_id = $3`
	result, err := db.Exec(query, category.Name, category.ID, category.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Deleta categoria pelo ID e user_id
func DeleteCategory(db *sql.DB, id int64, userID int64) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	result, err := db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
