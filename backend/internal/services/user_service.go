package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/daviolvr/Fintrack/internal/cache"
	"github.com/daviolvr/Fintrack/internal/dto"
	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
	"gorm.io/gorm"
)

type UserService struct {
	DB    *gorm.DB
	cache *cache.Cache
}

// Construtor
func NewUserService(db *gorm.DB, cache *cache.Cache) *UserService {
	return &UserService{DB: db, cache: cache}
}

// Retorna o usuário
func (s *UserService) GetUser(userID uint) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%d:data", userID)

	var user models.User
	found, err := s.cache.Get(cacheKey, &user)
	if err == nil && found {
		fmt.Println("Pegando do cache:", cacheKey)
		return &user, nil
	}

	userPtr, err := repository.FindUserByID(s.DB, userID)
	if err != nil {
		return nil, err
	}

	if userPtr != nil {
		if err = s.cache.Set(cacheKey, userPtr, time.Minute*10); err != nil {
			fmt.Println("Erro ao salvar no cache:", err)
		}
	}

	return userPtr, nil
}

// Atualiza dados do usuário
func (s *UserService) UpdateUser(
	userID uint,
	input dto.UserUpdateInput,
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

	// Invalid o cache do user
	_ = s.cache.InvalidateUserData(userID)
	return s.GetUser(userID)
}

// Atualiza o saldo
func (s *UserService) UpdateBalance(userID uint, balance float64) error {
	user := models.User{
		ID:      userID,
		Balance: balance,
	}
	if err := repository.UpdateUserBalance(s.DB, &user); err != nil {
		return err
	}
	return s.cache.InvalidateUserData(userID)
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

	if err := repository.DeleteUser(s.DB, userID); err != nil {
		return err
	}

	return s.cache.InvalidateUserData(userID)
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

	if err := repository.UpdatePassword(s.DB, &updatedUser); err != nil {
		return err
	}

	return s.cache.InvalidateUserData(userID)
}
