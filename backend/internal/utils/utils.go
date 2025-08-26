package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized   = errors.New("não autorizado")
	ErrInvalidID      = errors.New("ID inválido")
	ErrNotFound       = errors.New("registro não encontrado")
	ErrInternalServer = errors.New("erro interno do servidor")
)

// Pega o ID do usuário e retorna
func GetUserID(c *gin.Context) (uint, error) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("usuário não autenticado")
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, errors.New("ID do usuário inválido")
	}

	return userID, nil
}

// Retorna um erro padronizado em JSON
func RespondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

// Faz parse de um parâmetro numérico da URL
func GetIDParam(c *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}

// Faz bind do JSON e retorna false se inválido
func BindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		RespondError(c, http.StatusBadRequest, "Dados inválidos")
		return false
	}
	return true
}

// Checa se é sql.ErrNoRows e responde NotFound
func HandleNotFound(c *gin.Context, err error, msg string) bool {
	if err == sql.ErrNoRows {
		RespondError(c, http.StatusNotFound, msg)
		return true
	}
	return false
}

// Resposta de sucesso com mensagem
func RespondMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// Valida o formato do email e os domínios
func ValidateEmail(email string, allowedDomains []string) error {
	// Regex pra validar formato básico de email
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		return errors.New("email inválido")
	}

	if len(allowedDomains) > 0 {
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return errors.New("email inválido")
		}
		domain := parts[1]
		allowed := false
		for _, d := range allowedDomains {
			if strings.EqualFold(domain, d) {
				allowed = true
				break
			}
		}
		if !allowed {
			return errors.New("domínio de email não permitido")
		}
	}

	return nil
}

func FormatTime(t *time.Time) string {
	if t == nil {
		return "nil"
	}
	return t.Format("2006-01-02")
}

func FormatUint(u *uint) string {
	if u == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *u)
}

func FormatFloat(f *float64) string {
	if f == nil {
		return "nil"
	}
	return fmt.Sprintf("%.2f", *f)
}

func FormatString(s *string) string {
	if s == nil {
		return "nil"
	}
	return *s
}
