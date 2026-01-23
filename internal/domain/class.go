package domain

import (
	"errors"
	"time"
)

var (
	ErrClassNotFound = errors.New("class not found")
	ErrInvalidClass  = errors.New("invalid class data")
)

type Class struct {
	ID           int       `json:"id"`
	SchoolID     int       `json:"school_id"`
	Name         string    `json:"name"`
	Level        string    `json:"level"`
	Section      string    `json:"section"`
	Capacity     int       `json:"capacity"`
	AcademicYear string    `json:"academic_year"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewClass(schoolID int, name, level, academicYear string) (*Class, error) {
	if schoolID <= 0 {
		return nil, ErrInvalidClass
	}
	if name == "" {
		return nil, errors.New("class name is required")
	}
	if level == "" {
		return nil, errors.New("class level is required")
	}
	if academicYear == "" {
		return nil, errors.New("academic year is required")
	}

	return &Class{
		SchoolID:     schoolID,
		Name:         name,
		Level:        level,
		AcademicYear: academicYear,
		Capacity:     40, // Default
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
