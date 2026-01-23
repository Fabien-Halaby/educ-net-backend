package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type StudentClassRepository interface {
	EnrollStudent(studentID, classID int) error
	GetStudentClass(studentID int) (*StudentClassInfo, error)
	RemoveStudent(studentID, classID int) error
}

type StudentClassInfo struct {
	StudentID      int
	ClassID        int
	ClassName      string
	EnrollmentDate time.Time
	IsActive       bool
}

type studentClassRepository struct {
	db *sql.DB
}

func NewStudentClassRepository(db *sql.DB) StudentClassRepository {
	return &studentClassRepository{db: db}
}

func (r *studentClassRepository) EnrollStudent(studentID, classID int) error {
	query := `
		INSERT INTO student_classes (student_id, class_id, enrollment_date, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (student_id, class_id) DO UPDATE
		SET is_active = TRUE, enrollment_date = $3
	`

	_, err := r.db.Exec(query, studentID, classID, time.Now(), true)
	if err != nil {
		return fmt.Errorf("failed to enroll student: %w", err)
	}

	return nil
}

func (r *studentClassRepository) GetStudentClass(studentID int) (*StudentClassInfo, error) {
	query := `
		SELECT sc.student_id, sc.class_id, c.name, sc.enrollment_date, sc.is_active
		FROM student_classes sc
		JOIN classes c ON sc.class_id = c.id
		WHERE sc.student_id = $1 AND sc.is_active = TRUE
		LIMIT 1
	`

	info := &StudentClassInfo{}
	err := r.db.QueryRow(query, studentID).Scan(
		&info.StudentID,
		&info.ClassID,
		&info.ClassName,
		&info.EnrollmentDate,
		&info.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No active enrollment
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student class: %w", err)
	}

	return info, nil
}

func (r *studentClassRepository) RemoveStudent(studentID, classID int) error {
	query := `
		UPDATE student_classes
		SET is_active = FALSE
		WHERE student_id = $1 AND class_id = $2
	`

	_, err := r.db.Exec(query, studentID, classID)
	if err != nil {
		return fmt.Errorf("failed to remove student: %w", err)
	}

	return nil
}
