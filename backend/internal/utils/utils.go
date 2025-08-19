package utils

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized   = errors.New("não autorizado")
	ErrInvalidID      = errors.New("ID inválido")
	ErrNotFound       = errors.New("registro não encontrado")
	ErrInternalServer = errors.New("erro interno do servidor")
)

// Pega o ID do usuário e retorna
func GetUserID(c *gin.Context) (int64, error) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("usuário não autenticado")
	}

	userID, ok := userIDValue.(int64)
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

// Calcula saldo do usuário após uma transação
func CalculateBalanceAfterTransaction(
	c *gin.Context,
	amount, balance float64,
	txType string,
) (float64, error) {
	if txType == "income" {
		after_balance := balance + amount
		RespondMessage(c, "Transação feita com sucesso")
		return after_balance, nil
	}

	if txType == "expense" {
		after_balance := balance - amount
		if after_balance < 0 {
			RespondError(c, http.StatusBadRequest, "Saldo insuficiente")
			return 0, errors.New("valor da transação ultrapassa o saldo da conta")
		}

		return after_balance, nil
	}

	return 0, errors.New("tipo de transação inválida")
}

// Calcula saldo do usuário após ele deletar uma transação
// func CalculateBalanceAfterDeleteTransaction(
// 	c *gin.Context,
// 	amount, balance float64,
// ) (float64, error) {
// 	afterBalance := balance + amount

// }
