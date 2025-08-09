package handlers

import (
	"database/sql"
	"net/http"

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
