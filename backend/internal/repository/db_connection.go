package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectToDB() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Erro ao abrir conexão com o banco: %v", err)
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	// Testa conexão real com o banco
	if err := db.Ping(); err != nil {
		log.Printf("Erro ao conectar ao banco: %v", err)
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	return db, nil
}
