package repository

import (
	"errors"
	"time"

	"github.com/daviolvr/Fintrack/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func FindUserByID(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User

	err := db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // mantém o comportamento de não encontrado
		}
		return nil, err
	}

	return &user, nil
}

// Atualiza os dados do usuário
func UpdateUser(db *gorm.DB, user *models.User) error {
	result := db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Atualiza o saldo da conta do usuário
func UpdateUserBalance(db *gorm.DB, user *models.User) error {
	result := db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("balance", user.Balance)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Deleta o usuário
func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&models.User{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Atualiza a senha do usuário
func UpdatePassword(db *gorm.DB, user *models.User) error {
	result := db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("password_hash", user.Password)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Incrementa as falhas de login e retorna o total atualizado
func IncrementFailedLogin(db *gorm.DB, userID uint) (int64, error) {
	var user models.User

	// Incrementa e retorna a nova contagem de failed_logins
	result := db.Model(&user).
		Where("id = ?", userID).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "failed_logins"}}}).
		UpdateColumn("failed_logins", gorm.Expr("failed_logins + ?", 1))

	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
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
