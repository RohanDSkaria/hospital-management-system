package api

import (
	"net/http"
	"time"

	"github.com/RohanDSkaria/hospital-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PatientHandler struct {
	patientService service.PatientService
}

func NewPatientHandler(s service.PatientService) *PatientHandler {
	return &PatientHandler{patientService: s}
}

type PatientRequest struct {
	FullName       string    `json:"full_name" binding:"required"`
	DateOfBirth    time.Time `json:"date_of_birth" binding:"required"`
	Address        string    `json:"address"`
	ContactNumber  string    `json:"contact_number"`
	MedicalHistory string    `json:"medical_history"`
}

// CreatePatient handles POST requests to create a patient
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req PatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from the context (set by AuthMiddleware)
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	patient, err := h.patientService.CreatePatient(req.FullName, req.Address, req.ContactNumber, req.DateOfBirth, req.MedicalHistory, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create patient"})
		return
	}
	c.JSON(http.StatusCreated, patient)
}

// GetAllPatients handles GET requests to fetch all patients
func (h *PatientHandler) GetAllPatients(c *gin.Context) {
	patients, err := h.patientService.GetAllPatients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch patients"})
		return
	}
	c.JSON(http.StatusOK, patients)
}

// GetPatientByID handles GET requests for a single patient
func (h *PatientHandler) GetPatientByID(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("patient_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}
	patient, err := h.patientService.GetPatientByID(patientID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch patient"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

// UpdatePatient handles PUT requests to update a patient
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("patient_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}
	var req PatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	patient, err := h.patientService.UpdatePatient(patientID, req.FullName, req.Address, req.ContactNumber, req.DateOfBirth, req.MedicalHistory)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update patient"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

// DeletePatient handles DELETE requests to remove a patient
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("patient_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}
	if err := h.patientService.DeletePatient(patientID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete patient"})
		return
	}
	c.Status(http.StatusNoContent)
}
