package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RiskLevel represents the dropout risk level of a student
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "LOW"
	RiskLevelMedium RiskLevel = "MEDIUM"
	RiskLevelHigh   RiskLevel = "HIGH"
)

// Student represents a student in the database
type Student struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudentID       string         `gorm:"uniqueIndex" json:"student_id"`
	StudentName     string         `json:"student_name"`
	Attendance      JSONB          `gorm:"type:jsonb" json:"attendance"`
	Assignments     JSONB          `gorm:"type:jsonb" json:"assignments"`
	Contacts        JSONB          `gorm:"type:jsonb" json:"contacts"`
	DropoutScore    *int           `json:"dropout_score"`
	DropoutRiskLevel *string       `json:"dropout_risk_level"`
	DropoutNote     *string        `json:"dropout_note"`
	CreatedAt       int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       int64          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// RiskEvaluation represents a risk evaluation in the database
type RiskEvaluation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StudentID uuid.UUID      `gorm:"type:uuid;index" json:"student_id"`
	Score     int            `json:"score"`
	RiskLevel RiskLevel      `json:"risk_level"`
	Note      string         `json:"note"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// JSONB is a wrapper for handling JSON data in GORM
type JSONB []byte

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

// MarshalJSON returns the JSON encoding of JSONB
func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON sets JSONB to a copy of data
func (j *JSONB) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null JSONB")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// AttendanceRecord represents an attendance record
type AttendanceRecord struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

// AssignmentRecord represents an assignment record
type AssignmentRecord struct {
	Date      string `json:"date"`
	Name      string `json:"name"`
	Submitted bool   `json:"submitted"`
}

// ContactRecord represents a contact record
type ContactRecord struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

// GetAttendanceRecords parses the attendance JSONB data
func (s *Student) GetAttendanceRecords() ([]AttendanceRecord, error) {
	var records []AttendanceRecord
	if len(s.Attendance) > 0 {
		if err := json.Unmarshal(s.Attendance, &records); err != nil {
			return nil, err
		}
	}
	return records, nil
}

// GetAssignmentRecords parses the assignments JSONB data
func (s *Student) GetAssignmentRecords() ([]AssignmentRecord, error) {
	var records []AssignmentRecord
	if len(s.Assignments) > 0 {
		if err := json.Unmarshal(s.Assignments, &records); err != nil {
			return nil, err
		}
	}
	return records, nil
}

// GetContactRecords parses the contacts JSONB data
func (s *Student) GetContactRecords() ([]ContactRecord, error) {
	var records []ContactRecord
	if len(s.Contacts) > 0 {
		if err := json.Unmarshal(s.Contacts, &records); err != nil {
			return nil, err
		}
	}
	return records, nil
}