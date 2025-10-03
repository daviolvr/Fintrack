package services

import (
	"errors"
	"os"
	"time"

	"github.com/daviolvr/Fintrack/internal/auth"
	"github.com/daviolvr/Fintrack/internal/dto"
	"github.com/daviolvr/Fintrack/internal/models"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) RegisterUser(input dto.RegisterInput) error {
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

func (s *AuthService) LoginUser(c *gin.Context, input dto.LoginInput) (*models.User, string, error) {
	user, err := repository.FindUserByEmail(s.DB, input.Email)
	if err != nil {
		return nil, "", errors.New("email ou senha incorretos")
	}

	now := time.Now()

	if user.LockedUntil != nil && user.LockedUntil.After(now) {
		return nil, "", errors.New("conta bloqueada. Tente mais tarde")
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		failed, _ := repository.IncrementFailedLogin(s.DB, user.ID)
		if failed >= 5 {
			repository.LockUser(s.DB, user.ID, now.Add(10*time.Minute))
			return nil, "", errors.New("conta bloqueada. Tente novamente em 10 minutos")
		}
		return nil, "", errors.New("email ou senha incorretos")
	}

	repository.ResetFailedLogin(s.DB, user.ID)

	accessToken, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	// Configura cookies HTTP-only
	c.SetCookie("access_token", accessToken, 3600, "/", "localhost", true, true)     // 1 hora
	c.SetCookie("refresh_token", refreshToken, 259200, "/", "localhost", true, true) // 3 dias

	return user, "Login realizado com sucesso", nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	jwtRefreshSecret := []byte(os.Getenv("JWT_REFRESH_SECRET"))

	// Valida o refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		return jwtRefreshSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return "", errors.New("refresh token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("refresh token inválido")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("refresh token malformado")
	}

	// Gera novo access token
	newAccessToken, err := auth.GenerateJWT(uint(userIDFloat))
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
