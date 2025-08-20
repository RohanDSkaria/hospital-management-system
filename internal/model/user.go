package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Role is a custom type for user roles
type Role string

const (
	Receptionist Role = "receptionist"
	Doctor       Role = "doctor"
)

// User represents a user in the system (receptionist or doctor)
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	FullName     string    `gorm:"size:255;not null"`
	Email        string    `gorm:"size:255;not null;unique"`
	PasswordHash string    `gorm:"not null"`
	Role         Role      `gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BeforeCreate is a GORM hook that runs before a new record is created
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}
