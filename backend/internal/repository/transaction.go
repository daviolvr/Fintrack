package repository

import (
	"database/sql"
	"strconv"
	"time"

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

// Busca transação com filtros opcionais
func FindTransactionsByUser(
	db *sql.DB,
	userID int64,
	fromDate, toDate *time.Time,
	categoryID *int64,
) ([]models.Transaction, error) {

	query := `SELECT id, user_id, category_id, type, amount, description, date, created_at, updated_at
	FROM transactions
	WHERE user_id = $1`

	args := []any{userID}
	argIndex := 2 // $1 já está ocupado pelo user_id

	// Filtro por data inicial
	if fromDate != nil {
		query += ` AND date >= $` + strconv.Itoa(argIndex)
		args = append(args, *fromDate)
		argIndex++
	}

	// Filtro por data final
	if toDate != nil {
		query += ` AND date <= $` + strconv.Itoa(argIndex)
		args = append(args, *toDate)
		argIndex++
	}

	// Filtro por categoria
	if categoryID != nil {
		query += ` AND category_id = $` + strconv.Itoa(argIndex)
		args = append(args, *categoryID)
		argIndex++
	}

	// Ordena mais recentes primeiro
	query += ` ORDER BY date DESC`

	// Executa a query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Monta slices de transações
	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.CategoryID,
			&t.Type,
			&t.Amount,
			&t.Description,
			&t.Date,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
