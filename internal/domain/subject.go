package domain

import (
	"time"
)

const (
	SubjectStatusActive   = "active"
	SubjectStatusInactive = "inactive"
	SubjectStatusArchived = "archived"
)

type Subject struct {
	ID          int       `json:"id"`
	SchoolID    int       `json:"school_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewSubject(schoolID int, name, code, description string) (*Subject, error) {
	if schoolID <= 0 {
		return nil, ErrSubjectInvalidID
	}
	if name == "" {
		return nil, ErrSubjectNameRequired
	}
	if code == "" {
		return nil, ErrSubjectCodeRequired
	}

	return &Subject{
		SchoolID:    schoolID,
		Name:        name,
		Code:        code,
		Description: description,
		Status:      SubjectStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s *Subject) Activate() {
	s.Status = SubjectStatusActive
	s.UpdatedAt = time.Now()
}

func (s *Subject) Archive() {
	s.Status = SubjectStatusArchived
	s.UpdatedAt = time.Now()
}

func (s *Subject) IsActive() bool {
	return s.Status == SubjectStatusActive
}
