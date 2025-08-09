package repository

import (
	"database/sql"

	"github.com/daviolvr/Fintrack/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (first_name, last_name, email, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`

	return db.QueryRow(
		query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.ID)
}

func FindUserByEmail(db *sql.DB, email string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, password_hash, created_at, updated_at
	FROM users WHERE email = $1`

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

func FindUserByID(db *sql.DB, id int64) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, password_hash, created_at, updated_at
	FROM users WHERE id = $1`

	row := db.QueryRow(query, id)

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

func UpdateUser(db *sql.DB, user *models.User) error {
	query := `UPDATE users SET first_name = $1, last_name = $2, email = $3, updated_at = NOW()
	WHERE id = $4`
	_, err := db.Exec(query, user.FirstName, user.LastName, user.Email, user.ID)
	return err
}

// func DeleteUser(db *sql.DB, id int64) error {
// 	query := `DELETE from users WHERE id = $1`
// 	_, err := db.Exec(query, id)
// 	return err
// }
