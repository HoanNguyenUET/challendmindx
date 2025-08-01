package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"mindx/models"
	"mindx/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	db      *gorm.DB
	service *services.StudentService
}

// NewHandler creates a new Handler instance
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db:      db,
		service: services.NewStudentService(db),
	}
}

// EvaluateRisk handles the POST /evaluate endpoint
// It parses the JSON file, evaluates dropout risk, and stores results in the database
func (h *Handler) EvaluateRisk(c echo.Context) error {
	// Read JSON file
	jsonFile, err := os.ReadFile("data.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read JSON file: " + err.Error(),
		})
	}

	// Parse raw JSON data
	var rawStudents []map[string]interface{}
	if err := json.Unmarshal(jsonFile, &rawStudents); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format: " + err.Error(),
		})
	}
	
	// Convert to proper Student models
	var students []models.Student
	for _, rawStudent := range rawStudents {
		student := models.Student{
			StudentID:   rawStudent["student_id"].(string),
			StudentName: rawStudent["student_name"].(string),
		}
		
		// Convert attendance to JSONB
		if attendance, ok := rawStudent["attendance"]; ok {
			attendanceBytes, err := json.Marshal(attendance)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid attendance data",
				})
			}
			student.Attendance = attendanceBytes
		}
		
		// Convert assignments to JSONB
		if assignments, ok := rawStudent["assignments"]; ok {
			assignmentsBytes, err := json.Marshal(assignments)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid assignments data",
				})
			}
			student.Assignments = assignmentsBytes
		}
		
		// Convert contacts to JSONB
		if contacts, ok := rawStudent["contacts"]; ok {
			contactsBytes, err := json.Marshal(contacts)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid contacts data",
				})
			}
			student.Contacts = contactsBytes
		}
		
		students = append(students, student)
	}

	// Process students and evaluate risk
	results, err := h.service.ProcessAndEvaluateStudents(students)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, results)
}

// ListStudents handles the GET /students endpoint
// It lists all students with evaluated risks
// Supports filtering by risk level and sorting
func (h *Handler) ListStudents(c echo.Context) error {
	// Get query parameters
	riskLevel := c.QueryParam("risk_level")
	sortBy := c.QueryParam("sort_by")
	
	// Get students with filters
	students, err := h.service.GetStudentsWithFilters(riskLevel, sortBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, students)
}