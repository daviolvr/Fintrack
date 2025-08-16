package handlers

import (
	"database/sql"
	"math"
	"net/http"
	"strconv"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	DB *sql.DB
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

// @BasePath /api/v1
// @Summary Cria uma categoria
// @Description Cria uma categoria de transação
// @Tags category
// @Accept json
// @Produce json
// @Param name path string true "Nome da categoria"
// @Success 201
// @Failure 401 {object} utils.ErrorResponse
// @Failure 501 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input struct {
		Name string `json:"name" binding:"required,min=2"`
	}
	if !utils.BindJSON(c, &input) {
		return
	}

	category := models.Category{
		UserID: userID,
		Name:   input.Name,
	}

	if err := repository.CreateCategory(h.DB, &category); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	c.JSON(http.StatusCreated, category)
}

// @BasePath /api/v1
// @Summary Lista as categorias
// @Description Lista as categorias do usuário
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} utils.CategoryListResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	// Query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	categories, total, err := repository.FindCategoriesByUser(h.DB, userID, search, page, limit)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
	})
}

// @BasePath /api/v1
// @Summary Atualiza uma categoria
// @Description Atualiza uma categoria do usuário
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 200 {object} utils.MessageResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	id, err := utils.GetIDParam(c, "id")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	var input struct {
		Name string `json:"name" binding:"required,min=2"`
	}
	if !utils.BindJSON(c, &input) {
		return
	}

	category := models.Category{
		ID:     id,
		UserID: userID,
		Name:   input.Name,
	}

	if err := repository.UpdateCategory(h.DB, &category); err != nil {
		if utils.HandleNotFound(c, err, utils.ErrNotFound.Error()) {
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInvalidID.Error())
		return
	}

	utils.RespondMessage(c, "Categoria atualizada com sucesso")
}

// @BasePath /api/v1
// @Summary Deleta uma categoria
// @Description Deleta a categoria em questão
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	id, err := utils.GetIDParam(c, "id")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	if err := repository.DeleteCategory(h.DB, id, userID); err != nil {
		if utils.HandleNotFound(c, err, utils.ErrNotFound.Error()) {
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
