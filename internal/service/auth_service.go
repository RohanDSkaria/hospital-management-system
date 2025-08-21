package service

import (
	"errors"

	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"github.com/RohanDSkaria/hospital-management-system/internal/repository"
	"github.com/RohanDSkaria/hospital-management-system/pkg/utils"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication services
type AuthService interface {
	RegisterUser(fullName, email, password string, role model.Role) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// RegisterUser handles the business logic for creating a new user
func (s *authService) RegisterUser(fullName, email, password string, role model.Role) (*model.User, error) {
	// 1. Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // A real database error occurred
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// 2. Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 3. Create the user model
	newUser := &model.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	// 4. Save the new user to the database
	if err := s.userRepo.SaveUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
