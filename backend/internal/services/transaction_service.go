package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/daviolvr/Fintrack/internal/cache"
	"github.com/daviolvr/Fintrack/internal/dto"
	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
	"gorm.io/gorm"
)

type TransactionService struct {
	DB    *gorm.DB
	cache *cache.Cache
}

func NewTransactionService(db *gorm.DB, cache *cache.Cache) *TransactionService {
	return &TransactionService{DB: db, cache: cache}
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

	// Invalida cache de transações do usuário
	s.cache.InvalidateUserTransactions(userID)

	return transaction, nil
}

// Recupera uma transação
func (s *TransactionService) RetrieveTransaction(userID, transactionID uint) (*models.Transaction, error) {
	cacheKey := fmt.Sprintf("transactions:user=%d:transaction=%d", userID, transactionID)

	var cached dto.TransactionRetrieveCacheData
	found, err := s.cache.Get(cacheKey, &cached)
	if err == nil && found {
		fmt.Println("Pegando do cache:", cacheKey)
		return &cached.Transaction, nil
	}

	tx, err := repository.RetrieveTransactionByIDAndUserID(s.DB, userID, transactionID)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("transação não encontrada")
	}

	// Salva no cache
	if err := s.cache.Set(cacheKey, dto.TransactionRetrieveCacheData{
		Transaction: *tx,
	}, time.Minute*2); err != nil {
		fmt.Println("Erro ao salvar no cache:", err)
	}

	return tx, nil
}

// Lista transações com filtros e paginação
func (s *TransactionService) ListTransactions(
	userID uint,
	fromDate, toDate *time.Time,
	categoryID *uint,
	minAmount, maxAmount *float64,
	txType *string,
	page, limit int,
) ([]models.Transaction, int, error) {
	// Monta a chave do cache
	cacheKey := fmt.Sprintf(
		"transactions:user=%d:from=%s:to=%s:cat=%s:min=%s:max=%s:type=%s:page=%d:limit=%d",
		userID,
		utils.FormatTime(fromDate),
		utils.FormatTime(toDate),
		utils.FormatUint(categoryID),
		utils.FormatFloat(minAmount),
		utils.FormatFloat(maxAmount),
		utils.FormatString(txType),
		page,
		limit,
	)

	// Verifica se existe no cache
	var cached dto.TransactionListCacheData
	found, err := s.cache.Get(cacheKey, &cached)
	if err == nil && found {
		fmt.Println("Pegando do cache:", cacheKey)
		return cached.Transactions, cached.Total, nil
	}

	transactions, total, err := repository.FindTransactionsByUser(
		s.DB,
		userID,
		fromDate,
		toDate,
		categoryID,
		minAmount,
		maxAmount,
		txType,
		page,
		limit,
	)
	if err != nil {
		return nil, 0, err
	}

	// Salva no cache
	if err := s.cache.Set(cacheKey, dto.TransactionListCacheData{
		Transactions: transactions,
		Total:        total,
	}, time.Minute*2); err != nil {
		fmt.Println("Erro ao salvar no cache:", err)
	}

	return transactions, total, nil
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

	// Invalida cache de transações do usuário
	s.cache.InvalidateUserTransactions(userID)

	return tx, nil
}

// Deleta transação
func (s *TransactionService) DeleteTransaction(userID, transactionID uint) error {
	err := repository.DeleteTransactionByUser(s.DB, userID, transactionID)
	if err != nil {
		return err
	}

	// Invalida cache de transações do usuário
	s.cache.InvalidateUserTransactions(userID)

	return nil
}
