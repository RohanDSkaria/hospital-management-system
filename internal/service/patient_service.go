package service

import (
	"github.com/RohanDSkaria/hospital-management-system/internal/model"
	"github.com/RohanDSkaria/hospital-management-system/internal/repository"
	"github.com/google/uuid"
	"time"
)

type PatientService interface {
	CreatePatient(fullName, address, contact string, dob time.Time, history string, registeredByID uuid.UUID) (*model.Patient, error)
	GetAllPatients() ([]model.Patient, error)
	GetPatientByID(id uuid.UUID) (*model.Patient, error)
	UpdatePatient(id uuid.UUID, fullName, address, contact string, dob time.Time, history string) (*model.Patient, error)
	DeletePatient(id uuid.UUID) error
}

type patientService struct {
	patientRepo repository.PatientRepository
}

func NewPatientService(repo repository.PatientRepository) PatientService {
	return &patientService{patientRepo: repo}
}

func (s *patientService) CreatePatient(fullName, address, contact string, dob time.Time, history string, registeredByID uuid.UUID) (*model.Patient, error) {
	patient := &model.Patient{
		FullName:       fullName,
		Address:        address,
		ContactNumber:  contact,
		DateOfBirth:    dob,
		MedicalHistory: history,
		RegisteredByID: registeredByID,
	}
	err := s.patientRepo.Create(patient)
	return patient, err
}

func (s *patientService) GetAllPatients() ([]model.Patient, error) {
	return s.patientRepo.FindAll()
}

func (s *patientService) GetPatientByID(id uuid.UUID) (*model.Patient, error) {
	return s.patientRepo.FindByID(id)
}

func (s *patientService) UpdatePatient(id uuid.UUID, fullName, address, contact string, dob time.Time, history string) (*model.Patient, error) {
	patient, err := s.patientRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	// Update fields
	patient.FullName = fullName
	patient.Address = address
	patient.ContactNumber = contact
	patient.DateOfBirth = dob
	patient.MedicalHistory = history

	err = s.patientRepo.Update(patient)
	return patient, err
}

func (s *patientService) DeletePatient(id uuid.UUID) error {
	return s.patientRepo.Delete(id)
}
