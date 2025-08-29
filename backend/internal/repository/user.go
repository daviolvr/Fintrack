package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/daviolvr/Fintrack/internal/models"
	"gorm.io/gorm"
)

// Cria usuário
func CreateUser(db *gorm.DB, user *models.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// Busca usuário pelo email
func FindUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Busca usuário pelo ID
func FindUserByID(db *sql.DB, id uint) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, balance, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	row := db.QueryRow(query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Atualiza os dados do usuário
func UpdateUser(db *sql.DB, user *models.User) error {
	query := `
		UPDATE users SET first_name = $1, last_name = $2, email = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := db.Exec(query, user.FirstName, user.LastName, user.Email, user.ID)
	return err
}

// Atualiza o saldo da conta do usuário
func UpdateUserBalance(db *sql.DB, user *models.User) error {
	query := `
		UPDATE users SET balance = $1
		WHERE id = $2
	`

	_, err := db.Exec(query, user.Balance, user.ID)
	return err
}

// Deleta o usuário
func DeleteUser(db *sql.DB, id uint) error {
	query := `
		DELETE from users
		WHERE id = $1
	`
	_, err := db.Exec(query, id)
	return err
}

// Atualiza a senha do usuário
func UpdatePassword(db *sql.DB, user *models.User) error {
	query := `
		UPDATE users SET password_hash = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := db.Exec(query, user.Password, user.ID)
	return err
}

// Incrementa as falhas de login e retorna o total atualizado
func IncrementFailedLogin(db *gorm.DB, userID uint) (int64, error) {
	var user models.User

	result := db.Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("failed_logins", gorm.Expr("failed_logins + ?", 1))

	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	// Recupera o valor atualizado
	if err := db.Select("failed_logins").First(&user, userID).Error; err != nil {
		return 0, err
	}

	return int64(user.FailedLogins), nil
}

// Bloqueia o acesso do usuário até "until"
func LockUser(db *gorm.DB, userID uint, until time.Time) error {
	result := db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("locked_until", until)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Reseta o contador e tira o bloqueio (quando login for bem-sucedido)
func ResetFailedLogin(db *gorm.DB, userID uint) error {
	result := db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"failed_logins": 0,
			"locked_until":  nil,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Apenas desbloqueia o usuário (quando bloqueio já expirou, mas login ainda não teve sucesso)
// func UnlockUser(db *sql.DB, userID uint) error {
// 	_, err := db.Exec(`
// 		UPDATE users
// 		SET locked_until = NULL
// 		WHERE id = $1
// 	`, userID)
// 	return err
// }
