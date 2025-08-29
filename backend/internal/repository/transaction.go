package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/daviolvr/Fintrack/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Cria transação
func CreateTransaction(db *gorm.DB, t *models.Transaction) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var user models.User

		// Bloqueia a linha do usuário para atualizar saldo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&user, t.UserID).Error; err != nil {
			return err
		}

		// Checa saldo se for despesa
		if t.Type == "expense" && user.Balance < t.Amount {
			return fmt.Errorf("saldo insuficiente")
		}

		// Insere a transação
		if err := tx.Create(t).Error; err != nil {
			return err
		}

		// Atualiza saldo
		balanceChange := t.Amount
		if t.Type == "expense" {
			balanceChange = -t.Amount
		}

		if err := tx.Model(&user).
			Update("balance", gorm.Expr("balance + ?", balanceChange)).Error; err != nil {
			return err
		}

		return nil
	})
}

// Busca uma única transação do usuário
func RetrieveTransactionByIDAndUserID(
	db *gorm.DB,
	userID, transactionID uint,
) (*models.Transaction, error) {
	var transaction models.Transaction

	err := db.Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &transaction, nil
}

// Busca transação com filtros opcionais
func FindTransactionsByUser(
	db *gorm.DB,
	userID uint,
	fromDate, toDate *time.Time,
	categoryID *uint,
	minAmount, maxAmount *float64,
	txType *string,
	page, limit int,
) ([]models.Transaction, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var transactions []models.Transaction
	var total int64

	// Monta a query base
	query := db.Model(&models.Transaction{}).Where("user_id = ?", userID)

	// Filtros opcionais
	if fromDate != nil {
		query = query.Where("date >= ?", *fromDate)
	}
	if toDate != nil {
		query = query.Where("date <= ?", *toDate)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if minAmount != nil {
		query = query.Where("amount >= ?", *minAmount)
	}
	if maxAmount != nil {
		query = query.Where("amount >= ?", *minAmount)
	}
	if txType != nil {
		query = query.Where("type = ?", *txType)
	}

	// Contagem total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginação e ordenação
	offset := (page - 1) * limit
	if err := query.Order("date desc").Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, int(total), nil
}

// Atualiza uma transação pertencente a um usuário
func UpdateTransaction(db *gorm.DB, t *models.Transaction) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var oldTx models.Transaction
		var user models.User

		// Bloqueia a transação antiga para obter amount e type
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", t.ID, t.UserID).
			First(&oldTx).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound
			}
			return err
		}

		// Bloqueia a linha do usuário para atualizar saldo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&user, t.UserID).Error; err != nil {
			return err
		}

		// Remove efeito antigo da transação
		if oldTx.Type == "income" {
			user.Balance -= oldTx.Amount
		} else {
			user.Balance += oldTx.Amount
		}

		// Aplica novo valor da transação
		if t.Type == "income" {
			user.Balance += t.Amount
		} else {
			user.Balance -= t.Amount
		}

		// Checa saldo negativo
		if user.Balance < 0 {
			return fmt.Errorf("saldo insuficiente")
		}

		// Atualiza a transação
		if err := tx.Model(&oldTx).Updates(models.Transaction{
			CategoryID:  t.CategoryID,
			Type:        t.Type,
			Amount:      t.Amount,
			Description: t.Description,
			Date:        t.Date,
		}).Error; err != nil {
			return err
		}

		// Atualiza saldo do usuário
		if err := tx.Model(&user).Update("balance", user.Balance).Error; err != nil {
			return err
		}

		return nil
	})
}

// Deleta uma transação pelo ID e pelo userID
func DeleteTransactionByUser(db *gorm.DB, userID uint, transactionID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var transaction models.Transaction
		var user models.User

		// Bloqueia a transação para pegar amount e type
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", transactionID, userID).
			First(&transaction).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound
			}
			return err
		}

		// Bloqueia linha do usuário
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&user, userID).Error; err != nil {
			return err
		}

		// Remove efeito da transação do saldo
		if transaction.Type == "income" {
			user.Balance -= transaction.Amount
		} else {
			user.Balance += transaction.Amount
		}

		// Checa saldo negativo (opcional)
		if user.Balance < 0 {
			return fmt.Errorf("saldo insuficiente")
		}

		// Deleta a transação
		if err := tx.Delete(&transaction).Error; err != nil {
			return err
		}

		// Atualiza saldo do usuário
		if err := tx.Model(&user).Update("balance", user.Balance).Error; err != nil {
			return err
		}

		return nil
	})
}
