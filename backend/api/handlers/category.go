package handlers

import (
	"net/http"
	"strconv"

	"github.com/daviolvr/Fintrack/internal/services"
	"github.com/daviolvr/Fintrack/internal/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	Service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{Service: service}
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

	var input utils.CategoryInput
	if !utils.BindJSON(c, &input) {
		return
	}

	category, err := h.Service.CreateCategory(userID, input.Name)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := utils.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	c.JSON(http.StatusCreated, resp)
}

// @BasePath /api/v1
// @Summary Lista as categorias
// @Description Lista as categorias do usuário
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} utils.PaginatedCategoriesResponse
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

	categories, total, err := h.Service.ListCategories(userID, search, page, limit)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	var respCategories []utils.CategoryResponse
	for _, cat := range categories {
		respCategories = append(respCategories, utils.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	resp := utils.PaginatedCategoriesResponse{
		Data:       respCategories,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: h.Service.TotalPages(total, limit),
	}

	c.JSON(http.StatusOK, resp)
}

// @BasePath /api/v1
// @Summary Atualiza uma categoria
// @Description Atualiza uma categoria do usuário
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 200 {object} utils.CategoryResponse
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

	paramID, err := utils.GetIDParam(c, "id")
	id := uint(paramID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	var input utils.CategoryInput
	if !utils.BindJSON(c, &input) {
		return
	}

	category, err := h.Service.UpdateCategory(userID, id, input.Name)
	if err != nil {
		if utils.HandleNotFound(c, err, utils.ErrNotFound.Error()) {
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := utils.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	c.JSON(http.StatusOK, resp)
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

	paramID, err := utils.GetIDParam(c, "id")
	id := uint(paramID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, utils.ErrInvalidID.Error())
		return
	}

	if err := h.Service.DeleteCategory(id, userID); err != nil {
		if utils.HandleNotFound(c, err, utils.ErrNotFound.Error()) {
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
