package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
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

// Lista as transações por usuário
func (h *TransactionHandler) List(c *gin.Context) {
	userID, err := services.GetUserID(c)
	if err != nil {
		services.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Lê parâmetros de filtro da query string
	var fromDatePtr *time.Time
	if from := c.Query("from_date"); from != "" {
		if parsed, err := time.Parse("2006-01-02", from); err == nil {
			fromDatePtr = &parsed
		} else {
			services.RespondError(c, http.StatusBadRequest, "Formato de from_date inválido")
			return
		}
	}

	var toDatePtr *time.Time
	if to := c.Query("to_date"); to != "" {
		if parsed, err := time.Parse("2006-01-02", to); err == nil {
			toDatePtr = &parsed
		} else {
			services.RespondError(c, http.StatusBadRequest, "Formato de to_date inválido")
			return
		}
	}

	var categoryIDPtr *int64
	if cat := c.Query("category_id"); cat != "" {
		if parsed, err := strconv.ParseInt(cat, 10, 64); err == nil {
			categoryIDPtr = &parsed
		} else {
			services.RespondError(c, http.StatusBadRequest, "category_id inválido")
			return
		}
	}

	// Busca no banco
	transactions, err := repository.FindTransactionsByUser(h.DB, userID, fromDatePtr, toDatePtr, categoryIDPtr)
	if err != nil {
		services.RespondError(c, http.StatusInternalServerError, "Erro ao buscar transações")
		return
	}

	// Retorna JSON
	c.JSON(http.StatusOK, transactions)
}

// Atualiza transação do usuário
func (h *TransactionHandler) Update(c *gin.Context) {
	userID, err := services.GetUserID(c)
	if err != nil {
		services.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	// ID da transação
	transactionID, err := services.GetIDParam(c, "id")
	if err != nil {
		services.RespondError(c, http.StatusBadRequest, "ID inválido")
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

	// Converte data
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		services.RespondError(c, http.StatusBadRequest, "Data inválida")
		return
	}

	// Monta transação para atualização
	transaction := models.Transaction{
		ID:          transactionID,
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Type:        input.Type,
		Amount:      input.Amount,
		Description: input.Description,
		Date:        parsedDate,
	}

	// Atualiza no banco
	err = repository.UpdateTransaction(h.DB, &transaction)
	if services.HandleNotFound(c, err, "Transação não encontrada") {
		return
	}
	if err != nil {
		services.RespondError(c, http.StatusInternalServerError, "Erro ao atualizar transação")
		return
	}

	services.RespondMessage(c, "Transação atualizada com sucesso")
}
