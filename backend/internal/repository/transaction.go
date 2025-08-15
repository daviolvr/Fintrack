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
	page, limit int,
) ([]models.Transaction, int, error) {

	args := []any{userID}
	argIndex := 2

	query := `
		SELECT id, user_id, category_id, type, amount, description, date, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
	`

	if fromDate != nil {
		query += ` AND date >= $` + strconv.Itoa(argIndex)
		args = append(args, *fromDate)
		argIndex++
	}

	if toDate != nil {
		query += ` AND date <= $` + strconv.Itoa(argIndex)
		args = append(args, *toDate)
		argIndex++
	}

	if categoryID != nil {
		query += ` AND category_id = $` + strconv.Itoa(argIndex)
		args = append(args, *categoryID)
		argIndex++
	}

	// Ordena por mais recentes
	query += ` ORDER BY date desc`

	// Contagem total
	var total int
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_sub"
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Paginação
	offset := (page - 1) * limit
	query += ` LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

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
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}

	return transactions, total, nil
}

// Atualiza uma transação pertencente a um usuário
func UpdateTransaction(db *sql.DB, t *models.Transaction) error {
	query := `
	UPDATE transactions
	SET category_id = $1,
		type = $2,
		amount = $3,
		description = $4,
		date = $5,
		updated_at = NOW()
	WHERE id = $6 AND user_id = $7
	`

	result, err := db.Exec(query,
		t.CategoryID,
		t.Type,
		t.Amount,
		t.Description,
		t.Date,
		t.ID,
		t.UserID,
	)
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

// Deleta uma transação pelo ID e pelo userID
func DeleteTransactionByUser(db *sql.DB, userID int64, transactionID int64) error {
	query := `
	DELETE FROM transactions
	WHERE id = $1 AND user_id = $2
	`

	result, err := db.Exec(query, transactionID, userID)
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
