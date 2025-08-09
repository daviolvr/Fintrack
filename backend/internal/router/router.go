package router

import (
	"database/sql"

	"github.com/daviolvr/Fintrack/internal/handlers"
	"github.com/daviolvr/Fintrack/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)

	v1 := r.Group("/api/v1", middlewares.AuthMiddleware())

	// Rotas de user
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	v1.GET("/me", userHandler.Me)
	v1.PUT("/me", userHandler.Update)
	v1.DELETE("/me", userHandler.Delete)
	v1.PUT("/change_password", userHandler.UpdatePassword)

	// Rotas de categories
	v1.POST("/categories", categoryHandler.Create)
	v1.GET("/categories", categoryHandler.List)
}
