package cache

import (
	"github.com/daviolvr/Fintrack/internal/models"
)

type CategoryCacheData struct {
	Categories []models.Category
	Total      int
}

type TransactionCacheData struct {
	Transactions []models.Transaction
	Total        int
}
