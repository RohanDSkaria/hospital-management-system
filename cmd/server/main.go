package main

import (
	"log"
	"net/http"

	"github.com/RohanDSkaria/hospital-management-system/api"
	"github.com/RohanDSkaria/hospital-management-system/internal/database"
	"github.com/RohanDSkaria/hospital-management-system/internal/model"
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

	// --- Repositories ---
	userRepo := repository.NewUserRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	// --- Services ---
	authService := service.NewAuthService(userRepo)
	patientService := service.NewPatientService(patientRepo)

	// --- Handlers ---
	authHandler := api.NewAuthHandler(authService)
	patientHandler := api.NewPatientHandler(patientService)

	// --- Router ---
	router := gin.Default()

	// Public routes group
	v1Public := router.Group("/api/v1")
	{
		v1Public.POST("/register", authHandler.RegisterHandler)
		v1Public.POST("/login", authHandler.LoginHandler)
	}

	// Protected routes group
	v1Protected := router.Group("/api/v1")
	v1Protected.Use(api.AuthMiddleware())
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

		// --- Receptionist Routes ---
		receptionistRoutes := v1Protected.Group("/receptionist")
		receptionistRoutes.Use(api.RoleAuthMiddleware(model.Receptionist))
		{
			receptionistRoutes.POST("/patients", patientHandler.CreatePatient)
			receptionistRoutes.GET("/patients", patientHandler.GetAllPatients)
			receptionistRoutes.GET("/patients/:patient_id", patientHandler.GetPatientByID)
			receptionistRoutes.PUT("/patients/:patient_id", patientHandler.UpdatePatient)
			receptionistRoutes.DELETE("/patients/:patient_id", patientHandler.DeletePatient)
		}

		// --- Doctor Routes ---
		doctorRoutes := v1Protected.Group("/doctor")
		doctorRoutes.Use(api.RoleAuthMiddleware(model.Doctor))
		{
			doctorRoutes.GET("/patients", patientHandler.GetAllPatients)
			doctorRoutes.GET("/patients/:patient_id", patientHandler.GetPatientByID)
			doctorRoutes.PUT("/patients/:patient_id", patientHandler.UpdatePatient)
		}
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
