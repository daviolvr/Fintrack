package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/daviolvr/Fintrack/internal/auth"
	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
)

type AuthService struct {
	DB *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) RegisterUser(input utils.RegisterInput) error {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	if err := utils.ValidateEmail(input.Email, []string{"gmail.com", "outlook.com"}); err != nil {
		return err
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
	}

	return repository.CreateUser(s.DB, &user)
}

func (s *AuthService) LoginUser(input utils.LoginInput) (string, string, error) {
	user, err := repository.FindUserByEmail(s.DB, input.Email)
	if err != nil {
		return "", "", errors.New("email ou senha incorretos")
	}

	now := time.Now()

	if user.LockedUntil != nil && user.LockedUntil.After(now) {
		return "", "", errors.New("conta bloqueada. Tente mais tarde")
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		failed, _ := repository.IncrementFailedLogin(s.DB, user.ID)
		if failed >= 5 {
			repository.LockUser(s.DB, user.ID, now.Add(10*time.Minute))
			return "", "", errors.New("conta bloqueada. Tente novamente em 10 minutos")
		}
		return "", "", errors.New("email ou senha incorretos")
	}

	repository.ResetFailedLogin(s.DB, user.ID)

	accessToken, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
