package handlers

import (
	"database/sql"
	"net/http"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/services"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	DB *sql.DB
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

// Criar categoria
func (h *CategoryHandler) Create(c *gin.Context) {
	userID, err := services.GetUserID(c)
	if err != nil {
		services.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input struct {
		Name string `json:"name" binding:"required,min=2"`
	}
	if !services.BindJSON(c, &input) {
		return
	}

	category := models.Category{
		UserID: userID,
		Name:   input.Name,
	}

	if err := repository.CreateCategory(h.DB, &category); err != nil {
		services.RespondError(c, http.StatusInternalServerError, "Erro ao criar categoria")
		return
	}

	c.JSON(http.StatusCreated, category)
}

// Listar categorias do usuário
func (h *CategoryHandler) List(c *gin.Context) {
	userID, err := services.GetUserID(c)
	if err != nil {
		services.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	categories, err := repository.FindCategoriesByUser(h.DB, userID)
	if err != nil {
		services.RespondError(c, http.StatusInternalServerError, "Erro ao listar categorias")
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Atualizar categoria
func (h *CategoryHandler) Update(c *gin.Context) {
	userID, err := services.GetUserID(c)
	if err != nil {
		services.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := services.GetIDParam(c, "id")
	if err != nil {
		services.RespondError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input struct {
		Name string `json:"name" binding:"required,min=2"`
	}
	if !services.BindJSON(c, &input) {
		return
	}

	category := models.Category{
		ID:     id,
		UserID: userID,
		Name:   input.Name,
	}

	if err := repository.UpdateCategory(h.DB, &category); err != nil {
		if services.HandleNotFound(c, err, "Categoria não encontrada") {
			return
		}
		services.RespondError(c, http.StatusInternalServerError, "Erro ao atualizar categoria")
		return
	}

	services.RespondMessage(c, "Categoria atualizada com sucesso")
}

// Deletar categoria
func (h *CategoryHandler) Delete(c *gin.Context) {
	userID, err := services.GetUserID(c)
	if err != nil {
		services.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := services.GetIDParam(c, "id")
	if err != nil {
		services.RespondError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := repository.DeleteCategory(h.DB, id, userID); err != nil {
		if services.HandleNotFound(c, err, "Categoria não encontrada") {
			return
		}
		services.RespondError(c, http.StatusInternalServerError, "Erro ao deletar categoria")
		return
	}

	c.Status(http.StatusNoContent)
}
