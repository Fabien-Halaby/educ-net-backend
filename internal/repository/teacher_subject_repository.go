package repository

import (
	"database/sql"
	"fmt"
)

type TeacherSubjectRepository interface {
	AssignSubjects(teacherID int, subjectIDs []int) error
	GetTeacherSubjects(teacherID int) ([]int, error)
	RemoveAllSubjects(teacherID int) error
}

type teacherSubjectRepository struct {
	db *sql.DB
}

func NewTeacherSubjectRepository(db *sql.DB) TeacherSubjectRepository {
	return &teacherSubjectRepository{db: db}
}

func (r *teacherSubjectRepository) AssignSubjects(teacherID int, subjectIDs []int) error {
	if len(subjectIDs) == 0 {
		return nil
	}

	//! Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	//! Insert each subject
	stmt, err := tx.Prepare(`
		INSERT INTO teacher_subjects (teacher_id, subject_id)
		VALUES ($1, $2)
		ON CONFLICT (teacher_id, subject_id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, subjectID := range subjectIDs {
		_, err := stmt.Exec(teacherID, subjectID)
		if err != nil {
			return fmt.Errorf("failed to assign subject %d: %w", subjectID, err)
		}
	}

	return tx.Commit()
}

func (r *teacherSubjectRepository) GetTeacherSubjects(teacherID int) ([]int, error) {
	query := `
		SELECT subject_id
		FROM teacher_subjects
		WHERE teacher_id = $1
	`

	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher subjects: %w", err)
	}
	defer rows.Close()

	subjectIDs := []int{}
	for rows.Next() {
		var subjectID int
		if err := rows.Scan(&subjectID); err != nil {
			return nil, fmt.Errorf("failed to scan subject ID: %w", err)
		}
		subjectIDs = append(subjectIDs, subjectID)
	}

	return subjectIDs, nil
}

func (r *teacherSubjectRepository) RemoveAllSubjects(teacherID int) error {
	query := `DELETE FROM teacher_subjects WHERE teacher_id = $1`
	_, err := r.db.Exec(query, teacherID)
	if err != nil {
		return fmt.Errorf("failed to remove subjects: %w", err)
	}
	return nil
}
