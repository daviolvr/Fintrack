package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Recebe uma senha em texto puro e retorna o hash gerado ou erro
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compara a senha em texto puro com o hash armazenado
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
