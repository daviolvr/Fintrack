package repository

import (
	"database/sql"
	"time"

	"github.com/daviolvr/Fintrack/internal/models"
)

// Cria usuário
func CreateUser(db *sql.DB, user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`

	return db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.ID)
}

// Busca usuário pelo email
func FindUserByEmail(db *sql.DB, email string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	row := db.QueryRow(query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Busca usuário pelo ID
func FindUserByID(db *sql.DB, id int64) (*models.User, error) {
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
func DeleteUser(db *sql.DB, id int64) error {
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
func IncrementFailedLogin(db *sql.DB, userID int64) (int64, error) {
	var failedLogins int64
	err := db.QueryRow(`
        UPDATE users
        SET failed_logins = failed_logins + 1
        WHERE id = $1
        RETURNING failed_logins
    `, userID).Scan(&failedLogins)
	return failedLogins, err
}

// Bloqueia o acesso do usuário
func LockUser(db *sql.DB, userID int64, until time.Time) error {
	_, err := db.Exec(`
		UPDATE users
		SET locked_until = $1
		WHERE ID = $2
	`, until, userID)
	return err
}

// Reseta o contador e tira o bloqueio
func ResetFailedLogin(db *sql.DB, userID int64) error {
	_, err := db.Exec(`
		UPDATE users
		SET failed_logins = 0, locked_until = NULL
		WHERE id = $1
	`, userID)
	return err
}
