package repository

import (
	"database/sql"
	"fmt"

	"github.com/daviolvr/Fintrack/internal/models"
)

// Cria categoria para o usuário
func CreateCategory(db *sql.DB, category *models.Category) error {
	query := `
		INSERT into categories (user_id, name) VALUES ($1, $2) 
		RETURNING id
	`
	return db.QueryRow(query, category.UserID, category.Name).Scan(&category.ID)
}

// Busca categorias do usuário
func FindCategoriesByUser(
	db *sql.DB,
	userID uint,
	search string,
	page, limit int,
) ([]models.Category, int, error) {

	query := `
		SELECT id, user_id, name
		FROM categories
		WHERE user_id = $1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM categories
		WHERE user_id = $1
	`

	// Filtros
	args := []any{userID}
	if search != "" {
		query += ` AND name ILIKE $2`
		countQuery += ` AND name ILIKE $2`
		args = append(args, "%"+search+"%")
	}

	// Contagem total
	var total int
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Ordenação + paginação
	offset := (page - 1) * limit
	query += ` ORDER by name LIMIT $` + fmt.Sprint(len(args)+1) + ` OFFSET $` + fmt.Sprint(len(args)+2)
	args = append(args, limit, offset)

	// Execução
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name); err != nil {
			return nil, 0, err
		}
		categories = append(categories, c)
	}

	return categories, total, nil
}

// Atualiza categoria pelo ID e user_id
func UpdateCategory(db *sql.DB, category *models.Category) error {
	query := `
		UPDATE categories SET name = $1 
		WHERE id = $2 AND user_id = $3
	`
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
func DeleteCategory(db *sql.DB, id, userID uint) error {
	query := `
		DELETE FROM categories 
		WHERE id = $1 AND user_id = $2
	`
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
