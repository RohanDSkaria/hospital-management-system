package repository

import (
	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PatientRepository interface {
	Create(patient *model.Patient) error
	FindAll() ([]model.Patient, error)
	FindByID(id uuid.UUID) (*model.Patient, error)
	Update(patient *model.Patient) error
	Delete(id uuid.UUID) error
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) Create(patient *model.Patient) error {
	return r.db.Create(patient).Error
}

func (r *patientRepository) FindAll() ([]model.Patient, error) {
	var patients []model.Patient
	err := r.db.Find(&patients).Error
	return patients, err
}

func (r *patientRepository) FindByID(id uuid.UUID) (*model.Patient, error) {
	var patient model.Patient
	err := r.db.Where("id = ?", id).First(&patient).Error
	return &patient, err
}

func (r *patientRepository) Update(patient *model.Patient) error {
	return r.db.Save(patient).Error
}

func (r *patientRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Patient{}, id).Error
}
