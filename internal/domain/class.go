package domain

import (
	"time"
)

const (
	ClassStatusActive   = "active"
	ClassStatusInactive = "inactive"
	ClassStatusArchived = "archived"
)

type Class struct {
	ID           int       `json:"id"`
	SchoolID     int       `json:"school_id"`
	Name         string    `json:"name"`
	Level        string    `json:"level"`
	Section      string    `json:"section"`
	Capacity     int       `json:"capacity"`
	AcademicYear string    `json:"academic_year"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewClass(schoolID int, name, level, section, academicYear string) (*Class, error) {
	if schoolID <= 0 {
		return nil, ErrClassInvalidID
	}
	if name == "" {
		return nil, ErrClassNameRequired
	}
	if level == "" {
		return nil, ErrClassLevelRequired
	}
	if academicYear == "" {
		return nil, ErrClassYearRequired
	}

	return &Class{
		SchoolID:     schoolID,
		Name:         name,
		Level:        level,
		Section:      section,
		AcademicYear: academicYear,
		Capacity:     40,
		Status:       ClassStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (c *Class) Activate() {
	c.Status = ClassStatusActive
	c.UpdatedAt = time.Now()
}

func (c *Class) Archive() {
	c.Status = ClassStatusArchived
	c.UpdatedAt = time.Now()
}

func (c *Class) IsActive() bool {
	return c.Status == ClassStatusActive
}
