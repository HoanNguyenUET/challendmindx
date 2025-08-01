package services

import (
	"fmt"
	"strings"

	"mindx/config"
	"mindx/models"

	"gorm.io/gorm"
)

// StudentService handles business logic for student data
type StudentService struct {
	db     *gorm.DB
	config *config.RiskConfig
}

// NewStudentService creates a new StudentService instance
func NewStudentService(db *gorm.DB) *StudentService {
	return &StudentService{
		db:     db,
		config: &config.LoadConfig().Risk,
	}
}

// ProcessAndEvaluateStudents processes student data from JSON file, evaluates risk, and stores results
func (s *StudentService) ProcessAndEvaluateStudents(students []models.Student) ([]models.Student, error) {
	// Begin transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var updatedStudents []models.Student

	// Process each student
	for i := range students {
		// Check if student already exists
		var existingStudent models.Student
		result := tx.Where("student_id = ?", students[i].StudentID).First(&existingStudent)
		
		var student models.Student
		if result.Error == nil {
			// Student exists, update record
			if err := tx.Model(&existingStudent).Updates(map[string]interface{}{
				"student_name": students[i].StudentName,
				"attendance":   students[i].Attendance,
				"assignments":  students[i].Assignments,
				"contacts":     students[i].Contacts,
			}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			student = existingStudent
		} else if result.Error == gorm.ErrRecordNotFound {
			// Student doesn't exist, create new record
			if err := tx.Create(&students[i]).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			student = students[i]
		} else {
			// Other error
			tx.Rollback()
			return nil, result.Error
		}

		// Evaluate risk
		score, riskLevel, note, err := s.evaluateRisk(&student)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Update student record with risk evaluation
		if err := tx.Model(&student).Updates(map[string]interface{}{
			"dropout_score":      score,
			"dropout_risk_level": riskLevel,
			"dropout_note":       note,
		}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Get updated student record
		if err := tx.First(&student, student.ID).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		
		updatedStudents = append(updatedStudents, student)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return updatedStudents, nil
}

// GetAllStudents retrieves all students with their risk evaluations
func (s *StudentService) GetAllStudents() ([]models.Student, error) {
	var students []models.Student
	if err := s.db.Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

// GetStudentsWithFilters retrieves students with filtering and sorting options
func (s *StudentService) GetStudentsWithFilters(riskLevel, sortBy string) ([]models.Student, error) {
	var students []models.Student
	query := s.db
	
	// Apply risk level filter if provided
	if riskLevel != "" {
		query = query.Where("dropout_risk_level = ?", riskLevel)
	}
	
	// Apply sorting if provided
	if sortBy == "risk_level" {
		// Custom sorting order for risk levels: HIGH, MEDIUM, LOW
		query = query.Order("CASE dropout_risk_level " +
			"WHEN 'HIGH' THEN 1 " +
			"WHEN 'MEDIUM' THEN 2 " +
			"WHEN 'LOW' THEN 3 " +
			"ELSE 4 END")
	} else if sortBy == "risk_level_asc" {
		// Custom sorting order for risk levels: LOW, MEDIUM, HIGH
		query = query.Order("CASE dropout_risk_level " +
			"WHEN 'LOW' THEN 1 " +
			"WHEN 'MEDIUM' THEN 2 " +
			"WHEN 'HIGH' THEN 3 " +
			"ELSE 4 END")
	} else if sortBy == "score" {
		query = query.Order("dropout_score DESC")
	} else if sortBy == "score_asc" {
		query = query.Order("dropout_score ASC")
	} else {
		// Default sorting by student_id
		query = query.Order("student_id")
	}
	
	if err := query.Find(&students).Error; err != nil {
		return nil, err
	}
	
	return students, nil
}

// evaluateRisk evaluates the dropout risk for a student
func (s *StudentService) evaluateRisk(student *models.Student) (int, string, string, error) {
	attendanceRecords, err := student.GetAttendanceRecords()
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to parse attendance data: %w", err)
	}

	assignmentRecords, err := student.GetAssignmentRecords()
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to parse assignment data: %w", err)
	}

	contactRecords, err := student.GetContactRecords()
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to parse contact data: %w", err)
	}

	// Calculate risk factors
	var riskFactors []string
	score := 0

	// Attendance risk
	attendanceRate := s.calculateAttendanceRate(attendanceRecords)
	if attendanceRate < s.config.AttendanceThreshold {
		riskFactors = append(riskFactors, "attendance")
		score++
	}

	// Assignment risk
	assignmentRate := s.calculateAssignmentRate(assignmentRecords)
	if assignmentRate < s.config.AssignmentThreshold {
		riskFactors = append(riskFactors, "assignment")
		score++
	}

	// Contact risk
	contactFailures := s.countContactFailures(contactRecords)
	if contactFailures >= s.config.ContactThreshold {
		riskFactors = append(riskFactors, "communication")
		score++
	}

	// Determine risk level
	var riskLevel string
	switch {
	case score >= s.config.HighRiskThreshold:
		riskLevel = string(models.RiskLevelHigh)
	case score >= s.config.MediumRiskThreshold:
		riskLevel = string(models.RiskLevelMedium)
	default:
		riskLevel = string(models.RiskLevelLow)
	}

	// Create note
	var note string
	if len(riskFactors) > 0 {
		note = strings.Join(riskFactors, ", ") + " risk factors"
	} else {
		note = "No signs of disengagement detected"
	}

	return score, riskLevel, note, nil
}

// calculateAttendanceRate calculates the attendance rate as a percentage
func (s *StudentService) calculateAttendanceRate(attendance []models.AttendanceRecord) float64 {
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
func (s *StudentService) calculateAssignmentRate(assignments []models.AssignmentRecord) float64 {
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
func (s *StudentService) countContactFailures(contacts []models.ContactRecord) int {
	failures := 0
	for _, c := range contacts {
		if c.Status == "FAILED" {
			failures++
		}
	}
	return failures
}

