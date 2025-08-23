package services

import (
	"database/sql"
	"errors"

	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
)

type UserService struct {
	DB *sql.DB
}

// Construtor
func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

// Retorna o usuário
func (s *UserService) GetUser(userID uint) (*models.User, error) {
	return repository.FindUserByID(s.DB, userID)
}

// Atualiza dados do usuário
func (s *UserService) UpdateUser(
	userID uint,
	input utils.UserUpdateInput,
) (*models.User, error) {
	updatedUser := models.User{
		ID:        userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	if err := repository.UpdateUser(s.DB, &updatedUser); err != nil {
		return nil, err
	}

	return s.GetUser(userID)
}

// Atualiza o saldo
func (s *UserService) UpdateBalance(userID uint, balance float64) error {
	user := models.User{
		ID:      userID,
		Balance: balance,
	}
	return repository.UpdateUserBalance(s.DB, &user)
}

// Deleta usuário
func (s *UserService) DeleteUser(userID uint, password string) error {
	user, err := s.GetUser(userID)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return errors.New("senha incorreta")
	}

	return repository.DeleteUser(s.DB, userID)
}

// Atualiza senha
func (s *UserService) UpdatePassword(
	userID uint,
	currentPassword, newPassword string,
) error {
	user, err := s.GetUser(userID)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(currentPassword, user.Password) {
		return errors.New("senha incorreta")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	updatedUser := models.User{
		ID:       userID,
		Password: hashedPassword,
	}

	return repository.UpdatePassword(s.DB, &updatedUser)
}
