package handlers

import (
	"database/sql"
	"net/http"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// @BasePath /api/v1
// @Summary Retorna dados do usuário
// @Description Retorna os dados do usuário em questão
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} utils.UserMeResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /users/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := repository.FindUserByID(h.DB, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	resp := utils.UserMeResponse{
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
// @Param user body utils.UserUpdateParam true "Request body"
// @Success 200 {object} utils.UserUpdateResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /users/me [put]
func (h *UserHandler) Update(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input utils.UserUpdateInput

	if !utils.BindJSON(c, &input) {
		return
	}

	updatedUser := models.User{
		ID:        userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	if err := repository.UpdateUser(h.DB, &updatedUser); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	user, err := repository.FindUserByID(h.DB, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	resp := utils.UserUpdateResponse{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
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
// @Param data body utils.BalanceUpdateParam true "Novo saldo"
// @Success 200 {object} utils.UserUpdateBalanceResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /users/me/balance [patch]
func (h *UserHandler) UpdateBalance(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input utils.UserUpdateBalanceInput

	if !utils.BindJSON(c, &input) {
		return
	}

	user := models.User{
		ID:      userID,
		Balance: input.Balance,
	}

	if err := repository.UpdateUserBalance(h.DB, &user); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	resp := utils.UserUpdateBalanceResponse(input)

	c.JSON(http.StatusOK, resp)
}

// @BasePath /api/v1
// @Summary Deleta um usuário
// @Description Deleta o usuário em questão
// @Tags user
// @Accept json
// @Produce json
// @Success 204
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /users/me [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input utils.UserDeleteInput

	if !utils.BindJSON(c, &input) {
		return
	}

	user, err := repository.FindUserByID(h.DB, userID)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	if err := repository.DeleteUser(h.DB, userID); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// @BasePath /api/v1
// @Summary Atualiza a senha do usuário
// @Description Atualiza a senha do usuário em questão
// @Tags user
// @Param user body utils.UserChangePasswordParam true "Request body"
// @Success 200 {object} utils.MessageResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /users/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	var input utils.UserUpdatePasswordInput

	if !utils.BindJSON(c, &input) {
		return
	}

	user, err := repository.FindUserByID(h.DB, userID)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		utils.RespondError(c, http.StatusUnauthorized, utils.ErrUnauthorized.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	updatedUser := models.User{
		ID:       userID,
		Password: hashedPassword,
	}

	if err := repository.UpdatePassword(h.DB, &updatedUser); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	utils.RespondMessage(c, "Usuário atualizado com sucesso")
}
