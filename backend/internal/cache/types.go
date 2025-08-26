package cache

import (
	"github.com/daviolvr/Fintrack/internal/models"
)

type CategoryCacheData struct {
	Categories []models.Category
	Total      int
}

type TransactionRetrieveCacheData struct {
	Transaction models.Transaction
}

type TransactionListCacheData struct {
	Transactions []models.Transaction
	Total        int
}
