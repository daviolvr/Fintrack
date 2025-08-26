package handlers

import (
	"net/http"

	"github.com/daviolvr/Fintrack/internal/dto"
	"github.com/daviolvr/Fintrack/internal/services"
	"github.com/daviolvr/Fintrack/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// @BasePath /api/v1
// @Summary Retorna dados do usuário
// @Description Retorna os dados do usuário em questão
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserMeResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	user, err := h.Service.GetUser(userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	resp := dto.UserMeResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// @BasePath /api/v1
// @Summary Atualiza dados do usuário
// @Description Atualiza dados do usuário em questão
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.UserUpdateParam true "Request body"
// @Success 200 {object} dto.UserUpdateResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/me [put]
func (h *UserHandler) Update(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input dto.UserUpdateInput
	if !utils.BindJSON(c, &input) {
		return
	}

	user, err := h.Service.UpdateUser(userID, input)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.UserUpdateResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// @BasePath /api/v1
// @Summary Atualiza o saldo de um usuário
// @Description Atualiza o saldo do usuário em questão
// @Tags user
// @Accept json
// @Produce json
// @Param data body dto.BalanceUpdateParam true "Novo saldo"
// @Success 200 {object} dto.UserUpdateBalanceResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/me/balance [patch]
func (h *UserHandler) UpdateBalance(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input dto.UserUpdateBalanceInput
	if !utils.BindJSON(c, &input) {
		return
	}

	if err := h.Service.UpdateBalance(userID, input.Balance); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.UserUpdateBalanceResponse(input))
}

// @BasePath /api/v1
// @Summary Deleta um usuário
// @Description Deleta o usuário em questão
// @Tags user
// @Accept json
// @Produce json
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/me [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input dto.UserDeleteInput
	if !utils.BindJSON(c, &input) {
		return
	}

	if err := h.Service.DeleteUser(userID, input.Password); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// @BasePath /api/v1
// @Summary Atualiza a senha do usuário
// @Description Atualiza a senha do usuário em questão
// @Tags user
// @Param user body dto.UserChangePasswordParam true "Request body"
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input dto.UserUpdatePasswordInput
	if !utils.BindJSON(c, &input) {
		return
	}

	if err := h.Service.UpdatePassword(userID, input.Password, input.NewPassword); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondMessage(c, "Usuário atualizado com sucesso")
}
