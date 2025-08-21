package main

import (
	"log"
	"net/http"

	"github.com/RohanDSkaria/hospital-management-system/api"
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
	authHandler := api.NewAuthHandler(authService)

	router := gin.Default()

	// Public routes group
	v1Public := router.Group("/api/v1")
	{
		v1Public.POST("/register", authHandler.RegisterHandler)
		v1Public.POST("/login", authHandler.LoginHandler)
	}

	// Protected routes group
	v1Protected := router.Group("/api/v1")
	v1Protected.Use(api.AuthMiddleware()) // Apply the middleware to this group
	{
		// This is a sample protected route for testing
		v1Protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			userRole, _ := c.Get("userRole")

			c.JSON(http.StatusOK, gin.H{
				"message":   "This is a protected route",
				"user_id":   userID,
				"user_role": userRole,
			})
		})

	}

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
