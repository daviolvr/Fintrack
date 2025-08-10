package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/services"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	DB *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

// Criar transaction
func (h *TransactionHandler) Create(c *gin.Context) {
	userID, err := services.GetUserID(c)
	if err != nil {
		services.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input struct {
		CategoryID  int64   `json:"category_id" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=income expense"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
		Date        string  `json:"date" binding:"required,datetime=2006-01-02"`
	}
	if !services.BindJSON(c, &input) {
		return
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		services.RespondError(c, http.StatusBadRequest, "Data inválida")
		return
	}

	transaction := models.Transaction{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Type:        input.Type,
		Amount:      input.Amount,
		Description: input.Description,
		Date:        parsedDate,
	}

	if err := repository.CreateTransaction(h.DB, &transaction); err != nil {
		services.RespondError(c, http.StatusInternalServerError, "Erro ao criar transação")
		return
	}

	c.JSON(http.StatusCreated, transaction)
}
