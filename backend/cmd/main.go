package main

import (
	"log"
	"os"

	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db, err := repository.ConnectToDB()
	if err != nil {
		log.Fatalf("Erro ao conectar : %v", err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	log.Printf("Servidor rodando em http://localhost:%s\n", port)
	r.Run(":" + port)
}
