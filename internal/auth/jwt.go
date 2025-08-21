package auth

import (
	"os"
	"time"

	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims defines the structure of the JWT claims
type CustomClaims struct {
	UserID uuid.UUID  `json:"user_id"`
	Role   model.Role `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT for a given user
func GenerateToken(userID uuid.UUID, role model.Role) (string, error) {
	// Get the secret key from environment variables
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

	// Create the claims
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return it as a string
	return token.SignedString(jwtSecret)
}
