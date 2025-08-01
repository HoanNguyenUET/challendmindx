package services

import (
	"encoding/json"
	"strings"

	"mindx/config"
	"mindx/models"

	"gorm.io/gorm"
)

// RiskService handles business logic for risk evaluation
type RiskService struct {
	db     *gorm.DB
	config *config.RiskConfig
}

// NewRiskService creates a new RiskService instance
func NewRiskService(db *gorm.DB) *RiskService {
	return &RiskService{
		db:     db,
		config: &config.LoadConfig().Risk,
	}
}

// EvaluateStudentRisks processes student data and evaluates risk levels
func (s *RiskService) EvaluateStudentRisks(students []models.Student) ([]models.RiskEvaluation, error) {
	var results []models.RiskEvaluation

	// Begin transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Process each student
	for _, student := range students {
		// Store student data
		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Evaluate risk
		evaluation, err := s.evaluateRisk(&student)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Store evaluation
		if err := tx.Create(&evaluation).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		results = append(results, evaluation)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetAllRiskEvaluations retrieves all risk evaluations with student data
func (s *RiskService) GetAllRiskEvaluations() ([]models.RiskEvaluation, error) {
	var evaluations []models.RiskEvaluation
	if err := s.db.Preload("Student").Find(&evaluations).Error; err != nil {
		return nil, err
	}
	return evaluations, nil
}

// evaluateRisk evaluates the risk level for a student
func (s *RiskService) evaluateRisk(student *models.Student) (models.RiskEvaluation, error) {
	var attendanceData []struct {
		Date   string `json:"date"`
		Status string `json:"status"`
	}
	var assignmentData []struct {
		Date      string `json:"date"`
		Name      string `json:"name"`
		Submitted bool   `json:"submitted"`
	}
	var contactData []struct {
		Date   string `json:"date"`
		Status string `json:"status"`
	}

	// Parse JSON data
	if err := json.Unmarshal(student.Attendance, &attendanceData); err != nil {
		return models.RiskEvaluation{}, err
	}
	if err := json.Unmarshal(student.Assignments, &assignmentData); err != nil {
		return models.RiskEvaluation{}, err
	}
	if err := json.Unmarshal(student.Contacts, &contactData); err != nil {
		return models.RiskEvaluation{}, err
	}

	// Calculate risk factors
	var riskFactors []string
	score := 0

	// Attendance risk
	attendanceRate := calculateAttendanceRate(attendanceData)
	if attendanceRate < s.config.AttendanceThreshold {
		riskFactors = append(riskFactors, "attendance")
		score++
	}

	// Assignment risk
	assignmentRate := calculateAssignmentRate(assignmentData)
	if assignmentRate < s.config.AssignmentThreshold {
		riskFactors = append(riskFactors, "assignment")
		score++
	}

	// Contact risk
	contactFailures := countContactFailures(contactData)
	if contactFailures >= s.config.ContactThreshold {
		riskFactors = append(riskFactors, "communication")
		score++
	}

	// Determine risk level based on configurable thresholds
	var riskLevel models.RiskLevel
	switch {
	case score >= s.config.HighRiskThreshold:
		riskLevel = models.RiskLevelHigh
	case score >= s.config.MediumRiskThreshold:
		riskLevel = models.RiskLevelMedium
	default:
		riskLevel = models.RiskLevelLow
	}

	// Create note
	var note string
	if len(riskFactors) > 0 {
		note = strings.Join(riskFactors, ", ") + " risk factors"
	} else {
		note = "No signs of disengagement detected"
	}

	// Create evaluation
	evaluation := models.RiskEvaluation{
		StudentID: student.ID,
		Score:     score,
		RiskLevel: riskLevel,
		Note:      note,
	}

	return evaluation, nil
}

// calculateAttendanceRate calculates the attendance rate as a percentage
func calculateAttendanceRate(attendance []struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}) float64 {
	if len(attendance) == 0 {
		return 100.0
	}

	attended := 0
	for _, a := range attendance {
		if a.Status == "ATTEND" {
			attended++
		}
	}

	return float64(attended) / float64(len(attendance)) * 100.0
}

// calculateAssignmentRate calculates the assignment completion rate as a percentage
func calculateAssignmentRate(assignments []struct {
	Date      string `json:"date"`
	Name      string `json:"name"`
	Submitted bool   `json:"submitted"`
}) float64 {
	if len(assignments) == 0 {
		return 100.0
	}

	submitted := 0
	for _, a := range assignments {
		if a.Submitted {
			submitted++
		}
	}

	return float64(submitted) / float64(len(assignments)) * 100.0
}

// countContactFailures counts the number of failed contact attempts
func countContactFailures(contacts []struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}) int {
	failures := 0
	for _, c := range contacts {
		if c.Status == "FAILED" {
			failures++
		}
	}
	return failures
}