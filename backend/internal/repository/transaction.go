package repository

import (
	"database/sql"

	"github.com/daviolvr/Fintrack/internal/models"
)

// Cria transação
func CreateTransaction(db *sql.DB, t *models.Transaction) error {
	query := `INSERT INTO transactions (user_id, category_id, type, amount, description, date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	RETURNING id`

	return db.QueryRow(
		query,
		t.UserID,
		t.CategoryID,
		t.Type,
		t.Amount,
		t.Description,
		t.Date,
	).Scan(&t.ID)
}
