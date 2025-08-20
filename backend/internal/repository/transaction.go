package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/daviolvr/Fintrack/internal/models"
)

// Cria transação
func CreateTransaction(db *sql.DB, t *models.Transaction) error {
	// Começa a transação
	sqlTx, err := db.Begin()
	if err != nil {
		return err
	}

	// Pega saldo atual do usuário
	var currentBalance float64
	err = sqlTx.QueryRow(`
		SELECT balance
		FROM users
		WHERE id = $1
		FOR UPDATE
	`, t.UserID).Scan(&currentBalance)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	// Se for despesa e não tiver asldo suficiente, bloqueia
	if t.Type == "expense" && currentBalance < t.Amount {
		sqlTx.Rollback()
		return fmt.Errorf("saldo insuficiente")
	}

	// Insere a transação
	query := `
		INSERT INTO transactions (user_id, category_id, type, amount, description, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	err = sqlTx.QueryRow(
		query,
		t.UserID,
		t.CategoryID,
		t.Type,
		t.Amount,
		t.Description,
		t.Date,
	).Scan(&t.ID)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	// Atualiza o saldo
	var balanceChange float64
	if t.Type == "income" {
		balanceChange = t.Amount
	} else {
		balanceChange = -t.Amount
	}

	_, err = sqlTx.Exec(`
		UPDATE users
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2
	`, balanceChange, t.UserID)
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	// Confirma a transação
	return sqlTx.Commit()
}

// Busca uma única transação do usuário
func RetrieveTransactionByIDAndUserID(
	db *sql.DB,
	userID, transactionID int64,
) (*models.Transaction, error) {
	query := `
		SELECT id, user_id, category_id, type, amount, description, date, created_at, updated_at
		FROM transactions
		WHERE id = $1 AND user_id = $2
	`

	row := db.QueryRow(query, transactionID, userID)

	var transaction models.Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.CategoryID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.Description,
		&transaction.Date,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &transaction, nil
}

// Busca transação com filtros opcionais
func FindTransactionsByUser(
	db *sql.DB,
	userID int64,
	fromDate, toDate *time.Time,
	categoryID *int64,
	minAmount, maxAmount *float64,
	txType *string,
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

	if minAmount != nil {
		query += " AND amount >= $" + strconv.Itoa(argIndex)
		args = append(args, *minAmount)
		argIndex++
	}

	if maxAmount != nil {
		query += " AND amount <= $" + strconv.Itoa(argIndex)
		args = append(args, *maxAmount)
		argIndex++
	}

	if txType != nil {
		query += " AND type = $" + strconv.Itoa(argIndex)
		args = append(args, *txType)
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
	sqlTx, err := db.Begin()
	if err != nil {
		return err
	}

	// Pega a transação pra saber o valor
	var amount float64
	var txType string
	err = sqlTx.QueryRow(`
		SELECT amount, type
		FROM transactions
		WHERE id = $1 AND user_id = $2
	`, transactionID, userID).Scan(&amount, &txType)

	if err != nil {
		sqlTx.Rollback()
		return err
	}

	// Deleta a transação
	result, err := sqlTx.Exec(`
		DELETE FROM transactions
		WHERE id = $1 AND user_id = $2
	`, transactionID, userID)

	if err != nil {
		sqlTx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sqlTx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		sqlTx.Rollback()
		return sql.ErrNoRows
	}

	// Atualiza saldo do usuário de acordo com tipo de transação deletada
	if txType == "income" {
		_, err = sqlTx.Exec(`UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, userID)
	} else {
		_, err = sqlTx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, userID)
	}
	if err != nil {
		sqlTx.Rollback()
	}

	return sqlTx.Commit()
}
