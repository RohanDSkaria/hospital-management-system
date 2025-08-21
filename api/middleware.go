package api

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/RohanDSkaria/hospital-management-system/internal/auth"
	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware creates a gin middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// The header should be in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

		// Parse and validate the token
		claims := &auth.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// If the token is valid, set the user info in the context for later use
		c.Set("userID", claims.UserID.String())
		c.Set("userRole", claims.Role)

		// Continue to the next handler
		c.Next()
	}
}

// RoleAuthMiddleware checks if the user role from the JWT matches the required role
func RoleAuthMiddleware(requiredRole model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// We get the userRole from the context, which was set by AuthMiddleware
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user role not found in token"})
			return
		}

		// Check if the user's role is the one we require for this endpoint
		if userRole.(model.Role) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you are not authorized to perform this action"})
			return
		}

		// Role is correct, proceed to the handler
		c.Next()
	}
}
