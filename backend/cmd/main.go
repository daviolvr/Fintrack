package main

import (
	"fmt"
	"log"

	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db, err := repository.ConnectToDB()
	if err != nil {
		log.Fatalf("Erro ao conectar : %v", err)
	}
	defer db.Close()

	fmt.Println("Conexão com o banco de dados bem-sucedida!")
}
