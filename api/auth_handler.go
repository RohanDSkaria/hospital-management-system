package api

import (
	"net/http"

	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"github.com/RohanDSkaria/hospital-management-system/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRequest defines the structure for the user registration request body
type RegisterRequest struct {
	FullName string     `json:"full_name" binding:"required"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=8"`
	Role     model.Role `json:"role" binding:"required"`
}

// @Summary      Register a new user
// @Description  Creates a new user account (receptionist or doctor).
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user body RegisterRequest true "User Registration Info"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /register [post]
// RegisterHandler handles the user registration endpoint
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	// 1. Bind and validate the incoming JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Quick validation for role
	if req.Role != model.Doctor && req.Role != model.Receptionist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role specified"})
		return
	}

	// 2. Call the service to perform the business logic
	user, err := h.authService.RegisterUser(req.FullName, req.Email, req.Password, req.Role)
	if err != nil {
		// Check for specific error from the service layer
		if err.Error() == "user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		// For all other errors, return a generic server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	// 3. Format and send the success response
	response := gin.H{
		"id":         user.ID,
		"full_name":  user.FullName,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	}
	c.JSON(http.StatusCreated, response)
}

// LoginRequest defines the structure for the user login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Login user
// @Description  Authenticates a user and returns a JWT token for access to protected endpoints.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user body LoginRequest true "User Login Info"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /login [post]
// LoginHandler handles the user login endpoint
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest

	// Bind and validate the incoming JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to perform login logic
	token, err := h.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		// For any error from the service (e.g., "invalid credentials"), return Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Send the token back in the response
	c.JSON(http.StatusOK, gin.H{"token": token})
}
