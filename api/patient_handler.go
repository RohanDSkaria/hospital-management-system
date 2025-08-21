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

// @Summary      Create a new patient
// @Description  Creates a new patient record in the system. Only accessible by receptionists.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        patient body PatientRequest true "Patient Information"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /receptionist/patients [post]
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

// @Summary      Get all patients
// @Description  Retrieves a list of all patients in the system. Accessible by both receptionists and doctors.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /receptionist/patients [get]
// @Router       /doctor/patients [get]
// GetAllPatients handles GET requests to fetch all patients
func (h *PatientHandler) GetAllPatients(c *gin.Context) {
	patients, err := h.patientService.GetAllPatients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch patients"})
		return
	}
	c.JSON(http.StatusOK, patients)
}

// @Summary      Get patient by ID
// @Description  Retrieves a specific patient by their unique ID. Accessible by both receptionists and doctors.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        patient_id path string true "Patient ID" format(uuid)
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /receptionist/patients/{patient_id} [get]
// @Router       /doctor/patients/{patient_id} [get]
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

// @Summary      Update patient
// @Description  Updates an existing patient's information. Accessible by both receptionists and doctors.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        patient_id path string true "Patient ID" format(uuid)
// @Param        patient body PatientRequest true "Updated Patient Information"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /receptionist/patients/{patient_id} [put]
// @Router       /doctor/patients/{patient_id} [put]
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

// @Summary      Delete patient
// @Description  Deletes a patient from the system. Only accessible by receptionists.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        patient_id path string true "Patient ID" format(uuid)
// @Success      204  {string}  string "No Content"
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /receptionist/patients/{patient_id} [delete]
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
