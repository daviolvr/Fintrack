package router

import (
	"database/sql"

	"github.com/daviolvr/Fintrack/internal/handlers"
	"github.com/daviolvr/Fintrack/internal/middlewares"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)

	v1 := r.Group("/api/v1", middlewares.AuthMiddleware())

	// Rotas de auth
	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)
	r.POST("api/v1/refresh", handlers.RefreshToken)

	// Rotas de user
	v1.GET("/users/me", userHandler.Me)
	v1.PUT("/users/me", userHandler.Update)
	v1.PATCH("/users/me/balance", userHandler.UpdateBalance)
	v1.DELETE("/users/me", userHandler.Delete)
	v1.PUT("/users/password", userHandler.UpdatePassword)

	// Rotas de categories
	v1.POST("/categories", categoryHandler.Create)
	v1.GET("/categories", categoryHandler.List)
	v1.PUT("/categories/:id", categoryHandler.Update)
	v1.DELETE("/categories/:id", categoryHandler.Delete)

	// Rotas de transactions
	v1.POST("/transactions", transactionHandler.Create)
	v1.GET("/transactions", transactionHandler.List)
	v1.GET("/transactions/:id", transactionHandler.Retrieve)
	v1.PUT("/transactions/:id", transactionHandler.Update)
	v1.DELETE("/transactions/:id", transactionHandler.Delete)

	// Inicializa Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
