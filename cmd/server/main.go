package main

import (
	"log"
	"net/http"

	"github.com/RohanDSkaria/hospital-management-system/api"
	_ "github.com/RohanDSkaria/hospital-management-system/docs"
	"github.com/RohanDSkaria/hospital-management-system/internal/database"
	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"github.com/RohanDSkaria/hospital-management-system/internal/repository"
	"github.com/RohanDSkaria/hospital-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Hospital Management System API
// @version         1.0
// @description     This is the API for a simple hospital management system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		// @Summary      Get user profile
		// @Description  Returns the current user's profile information from the JWT token.
		// @Tags         Profile
		// @Accept       json
		// @Produce      json
		// @Success      200  {object}  map[string]interface{}
		// @Failure      401  {object}  map[string]interface{}
		// @Security     BearerAuth
		// @Router       /profile [get]
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

	// @Summary      Health check
	// @Description  Simple health check endpoint to verify the server is running.
	// @Tags         Health
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  map[string]interface{}
	// @Router       /ping [get]
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
