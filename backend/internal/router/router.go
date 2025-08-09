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

	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)
	r.GET("/api/v1/me", middlewares.AuthMiddleware(), userHandler.Me)
}
