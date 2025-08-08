package handlers

import (
	"database/sql"
	"net/http"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	DB *sql.DB
}

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Hasheia a senha
	hashedPassword, err := services.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao hashear a senha"})
		return
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
	}

	if err := repository.CreateUser(h.DB, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuário criado com sucesso",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Busca usuário pelo email
	user, err := repository.FindUserByEmail(h.DB, input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha incorretos"})
	}

	// Verifica a senha
	if !services.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha incorretos"})
		return
	}

	// Gera token
	token, err := services.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	// Retorna token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
