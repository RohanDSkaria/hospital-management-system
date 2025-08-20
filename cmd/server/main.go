package main

import (
	"github.com/RohanDSkaria/hospital-management-system/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	database.Connect()
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Println("Server is starting on port 8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
