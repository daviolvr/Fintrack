package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	DB *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

// @BasePath /api/v1
// @Summary Cria uma transação
// @Description Cria uma transação para o usuário em questão
// @Tags transaction
// @Accept json
// @Produce json
// @Param transaction body utils.TransactionCreateResponse true "Request body"
// @Success 201
// @Failure 401 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input struct {
		CategoryID  int64   `json:"category_id" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=income expense"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
		Date        string  `json:"date" binding:"required,datetime=2006-01-02"`
	}
	if !utils.BindJSON(c, &input) {
		return
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Data inválida")
		return
	}

	// Busca o saldo atual do usuário
	user, err := repository.FindUserByID(h.DB, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Erro ao buscar usuário")
		return
	}

	// Calcula o novo saldo
	newBalance, err := utils.CalculateBalanceAfterTransaction(c, input.Amount, user.Balance, input.Type)
	if err != nil {
		return
	}

	// Atualiza saldo no banco
	user.Balance = newBalance
	if err := repository.UpdateUserBalance(h.DB, user); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Erro ao atualizar saldo")
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
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) Retrieve(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	transactionID, err := utils.GetIDParam(c, "id")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	transaction, err := repository.RetrieveTransactionByIDAndUserID(h.DB, userID, transactionID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	if transaction == nil {
		utils.RespondError(c, http.StatusNotFound, "Transação não encontrada")
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// @BasePath /api/v1
// @Summary Lista as transações
// @Description Lista as transações de um usuário
// @Tags transaction
// @Accept json
// @Produce json
// @Success 200 {object} utils.TransactionListResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /transactions [get]
func (h *TransactionHandler) List(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	// Filtros de data
	var fromDatePtr, toDatePtr *time.Time
	if from := c.Query("from_date"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			fromDatePtr = &t
		} else {
			utils.RespondError(c, http.StatusBadRequest, "Formato de from_date inválido")
			return
		}
	}

	if to := c.Query("to_date"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			toDatePtr = &t
		} else {
			utils.RespondError(c, http.StatusBadRequest, "Formato de to_date inválido")
			return
		}
	}

	// Filtro por categoria
	var categoryIDPtr *int64
	if cat := c.Query("category_id"); cat != "" {
		if id, err := strconv.ParseInt(cat, 10, 64); err == nil {
			categoryIDPtr = &id
		} else {
			utils.RespondError(c, http.StatusBadRequest, "category_id inválido")
			return
		}
	}

	// Filtro por valor mínimo
	var minAmountPtr *float64
	if min := c.Query("min_amount"); min != "" {
		if val, err := strconv.ParseFloat(min, 64); err == nil {
			minAmountPtr = &val
		} else {
			utils.RespondError(c, http.StatusBadRequest, "min_amount inválido")
			return
		}
	}

	// Filtro por valor máximo
	var maxAmountPtr *float64
	if max := c.Query("max_amount"); max != "" {
		if val, err := strconv.ParseFloat(max, 64); err == nil {
			maxAmountPtr = &val
		} else {
			utils.RespondError(c, http.StatusBadRequest, "max_mount inválido")
			return
		}
	}

	// Filtro por tipo de transação (income ou expense)
	var typePtr *string
	if t := c.Query("type"); t != "" {
		if t == "income" || t == "expense" {
			typePtr = &t
		} else {
			utils.RespondError(c, http.StatusBadRequest, "type deve ser 'income' ou 'expense'")
			return
		}
	}

	// Parâmetros de paginação
	page := 1
	limit := 10
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	transactions, total, err := repository.FindTransactionsByUser(
		h.DB,
		userID,
		fromDatePtr,
		toDatePtr,
		categoryIDPtr,
		minAmountPtr,
		maxAmountPtr,
		typePtr,
		page,
		limit,
	)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       transactions,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
	})
}

// @BasePath /api/v1
// @Summary Atualiza uma transação
// @Description Atualiza uma transação do usuário em questão
// @Tags transaction
// @Accept json
// @Produce json
// @Param transaction body utils.TransactionUpdateResponse true "Request body"
// @Success 200 {object} utils.MessageResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /transactions/{id} [put]
func (h *TransactionHandler) Update(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	// ID da transação
	transactionID, err := utils.GetIDParam(c, "id")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	var input struct {
		CategoryID  int64   `json:"category_id" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=income expense"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
		Date        string  `json:"date" binding:"required,datetime=2006-01-02"`
	}
	if !utils.BindJSON(c, &input) {
		return
	}

	// Converte data
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Data inválida")
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
	if utils.HandleNotFound(c, err, utils.ErrNotFound.Error()) {
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	utils.RespondMessage(c, "Transação atualizada com sucesso")
}

// @BasePath /api/v1
// @Summary Deleta uma transação
// @Description Deleta uma transação do usuário em questão
// @Tags transaction
// @Accept json
// @Produce json
// @Param id path int true "ID da transação"
// @Success 200 {object} utils.MessageResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) Delete(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	transactionID, err := utils.GetIDParam(c, "id")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	user, err := repository.FindUserByID(h.DB, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Erro ao buscar usuário")
		return
	}

	// Busca a transação a ser deletada
	transaction, err := repository.RetrieveTransactionByIDAndUserID(h.DB, userID, transactionID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Erro ao buscar transação")
		return
	}

	// Atualiza o saldo do usuário
	user.Balance += transaction.Amount

	// Salva o novo saldo do usuário no banco
	err = repository.UpdateUserBalance(h.DB, user)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	// Deleta a transação
	err = repository.DeleteTransactionByUser(h.DB, userID, transactionID)
	if err != nil {
		if utils.HandleNotFound(c, err, utils.ErrNotFound.Error()) {
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	utils.RespondMessage(c, "Transação deletada com sucesso")
}
