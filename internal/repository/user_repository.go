package repository

import (
	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	SaveUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

// userRepository is the implementation of UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// SaveUser persists a new user to the database
func (r *userRepository) SaveUser(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByEmail finds a user by their email address
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
