package main

import (
	"log"
	"os"

	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega o .env
	_ = godotenv.Load()

	// Conecta ao banco
	db, err := repository.ConnectToDB()
	if err != nil {
		log.Fatalf("Erro ao conectar : %v", err)
	}
	defer db.Close()

	// Utiliza o engine do Gin
	r := gin.Default()

	// Seta as rotas
	router.SetupRoutes(r, db)

	// Tenta usar a porta do .env
	// Caso não tenha porta no .env, usa default 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	// Inicializa o servidor
	log.Printf("Servidor rodando em http://localhost:%s\n", port)
	r.Run(":" + port)
}
