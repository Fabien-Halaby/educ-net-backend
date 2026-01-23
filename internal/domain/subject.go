package domain

import (
	"errors"
	"time"
)

var (
	ErrSubjectNotFound = errors.New("subject not found")
	ErrInvalidSubject  = errors.New("invalid subject data")
)

type Subject struct {
	ID          int       `json:"id"`
	SchoolID    int       `json:"school_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewSubject(schoolID int, name, code string) (*Subject, error) {
	if schoolID <= 0 {
		return nil, ErrInvalidSubject
	}
	if name == "" {
		return nil, errors.New("subject name is required")
	}
	if code == "" {
		return nil, errors.New("subject code is required")
	}

	return &Subject{
		SchoolID:  schoolID,
		Name:      name,
		Code:      code,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
