package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/daviolvr/Fintrack/internal/dto"
	"github.com/daviolvr/Fintrack/internal/services"
	"github.com/daviolvr/Fintrack/internal/utils"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

// @BasePath /api/v1
// @Summary Cria uma transação
// @Description Cria uma transação para o usuário em questão
// @Tags transaction
// @Accept json
// @Produce json
// @Param transaction body dto.TransactionCreateParam true "Request body"
// @Success 201 {object} dto.TransactionCreateResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input dto.TransactionInput
	if !utils.BindJSON(c, &input) {
		return
	}

	tx, err := h.Service.CreateTransaction(userID, input.CategoryID, input.Type, input.Amount, input.Description, input.Date)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp := dto.TransactionCreateResponse{
		CategoryID:  tx.CategoryID,
		Type:        tx.Type,
		Amount:      tx.Amount,
		Description: tx.Description,
		Date:        tx.Date,
	}

	c.JSON(http.StatusCreated, resp)
}

// @BasePath /api/v1
// @Summary Retorna uma transação
// @Description Retorna uma transação do usuário em questão
// @Tags transaction
// @Accept json
// @Produce json
// @Param id path int true "ID da transação"
// @Success 200 {object} dto.TransactionResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /transactions/{id} [get]
func (h *TransactionHandler) Retrieve(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	paramID, err := utils.GetIDParam(c, "id")
	id := uint(paramID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	tx, err := h.Service.RetrieveTransaction(userID, id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, err.Error())
		return
	}

	resp := dto.TransactionResponse{
		CategoryID:  tx.CategoryID,
		Type:        tx.Type,
		Amount:      tx.Amount,
		Description: tx.Description,
		Date:        tx.Date,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// @BasePath /api/v1
// @Summary Lista as transações
// @Description Lista as transações de um usuário
// @Tags transaction
// @Accept json
// @Produce json
// @Success 200 {object} dto.PaginatedTransactionResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /transactions [get]
func (h *TransactionHandler) List(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	// Conversões de query params
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

	var fromDatePtr, toDatePtr *time.Time
	if from := c.Query("from_date"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			fromDatePtr = &t
		}
	}
	if to := c.Query("to_date"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			toDatePtr = &t
		}
	}

	var categoryIDPtr *uint
	if cat := c.Query("category_id"); cat != "" {
		if id, err := strconv.ParseUint(cat, 10, 64); err == nil {
			val := uint(id)
			categoryIDPtr = &val
		}
	}

	var minAmountPtr, maxAmountPtr *float64
	if min := c.Query("min_amount"); min != "" {
		if val, err := strconv.ParseFloat(min, 64); err == nil {
			minAmountPtr = &val
		}
	}
	if max := c.Query("max_amount"); max != "" {
		if val, err := strconv.ParseFloat(max, 64); err == nil {
			maxAmountPtr = &val
		}
	}

	var typePtr *string
	if t := c.Query("type"); t != "" {
		typePtr = &t
	}

	txs, total, err := h.Service.ListTransactions(userID, fromDatePtr, toDatePtr, categoryIDPtr, minAmountPtr, maxAmountPtr, typePtr, page, limit)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	var respTxs []dto.TransactionResponse
	for _, tx := range txs {
		respTxs = append(respTxs, dto.TransactionResponse{
			CategoryID:  tx.CategoryID,
			Type:        tx.Type,
			Amount:      tx.Amount,
			Description: tx.Description,
			Date:        tx.Date,
			CreatedAt:   tx.CreatedAt,
			UpdatedAt:   tx.UpdatedAt,
		})
	}

	resp := dto.PaginatedTransactionResponse{
		Data:       respTxs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + limit - 1) / limit,
	}

	c.JSON(http.StatusOK, resp)
}

// @BasePath /api/v1
// @Summary Atualiza uma transação
// @Description Atualiza uma transação do usuário em questão
// @Tags transaction
// @Accept json
// @Produce json
// @Param transaction body dto.TransactionUpdateParam true "Request body"
// @Success 200 {object} dto.TransactionResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /transactions/{id} [put]
func (h *TransactionHandler) Update(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	// Pega o ID da transação
	paramID, err := utils.GetIDParam(c, "id")
	id := uint(paramID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	var input dto.TransactionInput
	if !utils.BindJSON(c, &input) {
		return
	}

	tx, err := h.Service.UpdateTransaction(userID, id, input.CategoryID, input.Type, input.Amount, input.Description, input.Date)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp := dto.TransactionResponse{
		CategoryID:  tx.CategoryID,
		Type:        tx.Type,
		Amount:      tx.Amount,
		Description: tx.Description,
		Date:        tx.Date,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// @BasePath /api/v1
// @Summary Deleta uma transação
// @Description Deleta uma transação do usuário em questão
// @Tags transaction
// @Accept json
// @Produce json
// @Param id path int true "ID da transação"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) Delete(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	paramID, err := utils.GetIDParam(c, "id")
	id := uint(paramID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	if err := h.Service.DeleteTransaction(userID, id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
