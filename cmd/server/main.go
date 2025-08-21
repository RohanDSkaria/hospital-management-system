package main

import (
	"log"
	"net/http"

	"github.com/RohanDSkaria/hospital-management-system/internal/database"
	"github.com/RohanDSkaria/hospital-management-system/internal/repository"
	"github.com/RohanDSkaria/hospital-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	database.Connect()
	db := database.DB

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	_ = authService

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Println("Server is starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
