package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Patient represents a patient record
type Patient struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;"`
	FullName       string    `gorm:"size:255;not null"`
	DateOfBirth    time.Time
	Address        string
	ContactNumber  string    `gorm:"size:20"`
	MedicalHistory string    `gorm:"type:text"`
	RegisteredByID uuid.UUID // Foreign Key
	RegisteredBy   User      `gorm:"foreignKey:RegisteredByID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// BeforeCreate is a GORM hook for the Patient model
func (patient *Patient) BeforeCreate(tx *gorm.DB) (err error) {
	patient.ID = uuid.New()
	return
}
