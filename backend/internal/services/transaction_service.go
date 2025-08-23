package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
)

type TransactionService struct {
	DB *sql.DB
}

func NewTransactionService(db *sql.DB) *TransactionService {
	return &TransactionService{DB: db}
}

// Cria uma transação
func (s *TransactionService) CreateTransaction(
	userID, categoryID uint,
	txType string,
	amount float64,
	description, dateStr string,
) (*models.Transaction, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("data inválida")
	}

	transaction := &models.Transaction{
		UserID:      userID,
		CategoryID:  categoryID,
		Type:        txType,
		Amount:      amount,
		Description: description,
		Date:        parsedDate,
	}

	if err := repository.CreateTransaction(s.DB, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// Recupera uma transação
func (s *TransactionService) RetrieveTransaction(userID, transactionID uint) (*models.Transaction, error) {
	tx, err := repository.RetrieveTransactionByIDAndUserID(s.DB, userID, transactionID)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("transação não encontrada")
	}
	return tx, nil
}

// Lista transações com filtros e paginação
func (s *TransactionService) ListTransactions(
	userID uint,
	fromDate, toDate *time.Time,
	categoryID *uint,
	minAmoumt, maxAmount *float64,
	txType *string,
	page, limit int,
) ([]models.Transaction, int, error) {
	return repository.FindTransactionsByUser(
		s.DB,
		userID,
		fromDate,
		toDate,
		categoryID,
		minAmoumt,
		maxAmount,
		txType,
		page,
		limit,
	)
}

// Atualiza transação
func (s *TransactionService) UpdateTransaction(
	userID, transactionID, categoryID uint,
	txType string,
	amount float64,
	description, dateStr string,
) (*models.Transaction, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("data inválida")
	}

	tx := &models.Transaction{
		ID:          transactionID,
		UserID:      userID,
		CategoryID:  categoryID,
		Type:        txType,
		Amount:      amount,
		Description: description,
		Date:        parsedDate,
	}

	if err := repository.UpdateTransaction(s.DB, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// Deleta transação
func (s *TransactionService) DeleteTransaction(userID, transactionID uint) error {
	return repository.DeleteTransactionByUser(s.DB, userID, transactionID)
}
