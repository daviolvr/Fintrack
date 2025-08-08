package router

import (
	"database/sql"

	"github.com/daviolvr/Fintrack/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	authHandler := handlers.NewAuthHandler(db)

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
}
