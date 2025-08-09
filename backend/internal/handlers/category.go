package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
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
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID := userIDValue.(int64)

	var input struct {
		Name string `json:"name" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome de categoria inválido"})
		return
	}

	category := models.Category{
		UserID: userID,
		Name:   input.Name,
	}

	if err := repository.CreateCategory(h.DB, &category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar categoria"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// Listar categorias do usuário
func (h *CategoryHandler) List(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID := userIDValue.(int64)

	categories, err := repository.FindCategoriesByUser(h.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar categorias"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Atualizar categoria
func (h *CategoryHandler) Update(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID := userIDValue.(int64)

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome da categoria inválido"})
		return
	}

	category := models.Category{
		ID:     id,
		UserID: userID,
		Name:   input.Name,
	}

	if err := repository.UpdateCategory(h.DB, &category); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Categoria não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar categoria"})
		return
	}

	c.JSON(http.StatusOK, category)
}
