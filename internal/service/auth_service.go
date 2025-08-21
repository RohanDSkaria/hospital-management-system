package service

import (
	"errors"

	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"github.com/RohanDSkaria/hospital-management-system/internal/repository"
	"github.com/RohanDSkaria/hospital-management-system/internal/auth"
	"github.com/RohanDSkaria/hospital-management-system/pkg/utils"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication services
type AuthService interface {
	RegisterUser(fullName, email, password string, role model.Role) (*model.User, error)
	LoginUser(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// LoginUser handles the business logic for user login
func (s *authService) LoginUser(email, password string) (string, error) {
	// 1. Find the user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// If the record is not found, or any other DB error, return invalid credentials
		return "", errors.New("invalid credentials")
	}

	// 2. Compare the provided password with the stored hash
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	// 3. Generate a JWT
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

// RegisterUser handles the business logic for creating a new user
func (s *authService) RegisterUser(fullName, email, password string, role model.Role) (*model.User, error) {
	// 1. Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
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
