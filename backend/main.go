package main

import (
	"log"
	"os"
	"time"

	"github.com/daviolvr/Fintrack/docs"
	"github.com/daviolvr/Fintrack/internal/repository"
	"github.com/daviolvr/Fintrack/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Carrega o .env
	_ = godotenv.Load()
	
	frontEndPort := os.Getenv("FRONTEND_PORT")

	// Conecta ao banco
	db, err := repository.ConnectToDB()
	if err != nil {
		log.Fatalf("Erro ao conectar : %v", err)
	}
	defer db.Close()

	// Utiliza o engine do Gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:" + frontEndPort},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	// Seta as rotas
	router.SetupRoutes(r, db)

	// Configuraçẽos do Swagger
	docs.SwaggerInfo.Title = "Fintrack API"
	docs.SwaggerInfo.Description = "API para controle financeiro pessoal"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

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
