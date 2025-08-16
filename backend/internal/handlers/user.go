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
// @Success 200 {object} utils.UserResponse
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

	c.JSON(http.StatusOK, gin.H{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

// @BasePath /api/v1
// @Summary Atualiza dados do usuário
// @Description Atualiza dados do usuário em questão
// @Tags user
// @Accept json
// @Produce json
// @Param user body utils.UserUpdateInput true "Request body"
// @Success 200 {object} utils.MessageResponse
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

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if !utils.BindJSON(c, &input) {
		return
	}

	user := models.User{
		ID:        userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	if err := repository.UpdateUser(h.DB, &user); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, utils.ErrInternalServer.Error())
		return
	}

	utils.RespondMessage(c, "Usuário atualizado com sucesso")
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

	var input struct {
		Password string `json:"password"`
	}

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
// @Param user body utils.UserChangePassword true "Request body"
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

	var input struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}
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
