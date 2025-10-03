package handlers

import (
	"net/http"

	"github.com/daviolvr/Fintrack/internal/dto"
	"github.com/daviolvr/Fintrack/internal/services"
	"github.com/daviolvr/Fintrack/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// @BasePath /api/v1
// @Summary Registra um usuário
// @Description Registra um usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterInput true "Request Body with User data"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput
	if !utils.BindJSON(c, &input) {
		return
	}

	if err := h.Service.RegisterUser(input); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso"})
}

// @BasePath /api/v1
// @Summary Login de usuários
// @Description Login de usuários no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginInput true "Request body"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if !utils.BindJSON(c, &input) {
		return
	}

	user, message, err := h.Service.LoginUser(c, input)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.FirstName + " " + user.LastName,
		},
	})
}

// @BasePath /api/v1
// @Summary Atualiza token de acesso
// @Description Atualiza token de acesso do usuário
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body dto.RefreshTokenInput true "Refresh Token"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input dto.RefreshTokenInput
	if !utils.BindJSON(c, &input) {
		return
	}

	newAccessToken, err := h.Service.RefreshToken(input.RefreshToken)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
