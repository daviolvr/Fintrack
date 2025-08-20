package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/daviolvr/Fintrack/internal/auth"
	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// @BasePath /api/v1
// @Summary Registra um usuário
// @Description Registra um usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param user body utils.RegisterInput true "Request Body with User data"
// @Success 201 {object} utils.MessageResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input utils.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Hasheia a senha
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao hashear a senha"})
		return
	}

	// Valida o email
	if err := utils.ValidateEmail(input.Email, []string{"gmail.com", "outlook.com"}); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
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

// @BasePath /api/v1
// @Summary Login de usuários
// @Description Login de usuários no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param user body utils.LoginInput true "Request body"
// @Success 200 {object} utils.MessageResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input utils.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Busca usuário pelo email
	user, err := repository.FindUserByEmail(h.DB, input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha incorretos"})
		return
	}

	now := time.Now()

	// 🔒 Primeiro: checa se o usuário está bloqueado
	if user.LockedUntil != nil {
		if user.LockedUntil.After(now) {
			// Ainda dentro do bloqueio
			c.JSON(http.StatusForbidden, gin.H{"error": "Conta bloqueada. Tente mais tarde."})
			return
		} else {
			// Tempo de bloqueio expirou → libera
			repository.ResetFailedLogin(h.DB, user.ID)
			user.LockedUntil = nil
			user.FailedLogins = 0
		}
	}

	// Verifica a senha
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		// Incrementa contador de tentativas
		newFailedLogins, err := repository.IncrementFailedLogin(h.DB, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao incrementar falhas de login"})
			return
		}

		// Se atingiu o limite, bloqueia por 10 minutos
		if newFailedLogins >= 5 {
			repository.LockUser(h.DB, user.ID, now.Add(10*time.Minute))
			c.JSON(http.StatusForbidden, gin.H{"error": "Conta bloqueada. Tente novamente em 10 minutos."})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha incorretos"})
		return
	}

	// Se senha correta, zera contador
	if user.FailedLogins > 0 {
		repository.ResetFailedLogin(h.DB, user.ID)
	}

	// Gera access token
	accessToken, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	// Gera refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar refresh token"})
		return
	}

	// Retorna tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @BasePath /api/v1
// @Summary Atualiza token de acesso
// @Description Atualiza token de acesso do usuário
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body utils.RefreshTokenInput true "Refresh Token"
// @Success 200 {object} utils.MessageResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input utils.RefreshTokenInput

	// Faz o bind do JSON
	if !utils.BindJSON(c, &input) {
		return
	}

	// Pega a chave do refresh token
	jwtRefreshSecret := []byte(os.Getenv("JWT_REFRESH_SECRET"))

	// Valida o refresh token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (any, error) {
		return jwtRefreshSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		utils.RespondError(c, http.StatusUnauthorized, "Refresh token inválido")
		return
	}

	// Extrai claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "Refresh token inválido")
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "Refresh token malformado")
		return
	}

	// Gera novo access token
	newToken, err := auth.GenerateJWT(int64(userIDFloat))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Erro ao gerar token")
		return
	}

	// Retorna o novo access token
	c.JSON(http.StatusOK, gin.H{
		"access_token": newToken,
	})
}
