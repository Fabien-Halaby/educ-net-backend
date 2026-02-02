package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"fmt"
)

type TeacherSubjectRepository interface {
	Create(teacherID, subjectID int) error
	Delete(teacherID, subjectID int) error
	Exists(teacherID, subjectID int) (bool, error)
	FindByTeacher(teacherID int) ([]*domain.Subject, error)
	FindBySubject(subjectID int) ([]*domain.User, error)
	DeleteByTeacher(teacherID int) error
	DeleteBySubject(subjectID int) error
}

type teacherSubjectRepository struct {
	db *sql.DB
}

func NewTeacherSubjectRepository(db *sql.DB) TeacherSubjectRepository {
	return &teacherSubjectRepository{db: db}
}

func scanTeacherSubjectRow(rows *sql.Rows) (*domain.Subject, error) {
	var subject domain.Subject
	err := rows.Scan(
		&subject.ID,
		&subject.Name,
		&subject.Code,
		&subject.Description,
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)

	subject.SchoolID = 1
	return &subject, err
}

func (r *teacherSubjectRepository) Create(teacherID, subjectID int) error {
	_, err := r.db.Exec(
		`INSERT INTO teacher_subjects (teacher_id, subject_id) VALUES ($1, $2)`,
		teacherID, subjectID,
	)
	if err != nil {
		return fmt.Errorf("create teacher-subject: %w", err)
	}
	return nil
}

func (r *teacherSubjectRepository) Delete(teacherID, subjectID int) error {
	result, err := r.db.Exec(
		`DELETE FROM teacher_subjects WHERE teacher_id=$1 AND subject_id=$2`,
		teacherID, subjectID)
	if err != nil {
		return fmt.Errorf("delete teacher-subject: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrTeacherSubjectNotFound
	}
	return nil
}

func (r *teacherSubjectRepository) Exists(teacherID, subjectID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM teacher_subjects WHERE teacher_id=$1 AND subject_id=$2)`,
		teacherID, subjectID).Scan(&exists)
	return exists, err
}

func (r *teacherSubjectRepository) FindByTeacher(teacherID int) ([]*domain.Subject, error) {
	rows, err := r.db.Query(`
        SELECT s.id, s.name, s.code, s.description, s.created_at, s.updated_at
		FROM teacher_subjects ts
		JOIN subjects s ON ts.subject_id = s.id
		WHERE ts.teacher_id = $1
		ORDER BY s.name`, teacherID)
	if err != nil {
		return nil, fmt.Errorf("find teacher subjects: %w", err)
	}
	defer rows.Close()

	var subjects []*domain.Subject
	for rows.Next() {
		subject := &domain.Subject{}
		subject, err := scanTeacherSubjectRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan subject row: %w", err)
		}
		subjects = append(subjects, subject)
	}
	return subjects, rows.Err()
}

func (r *teacherSubjectRepository) FindBySubject(subjectID int) ([]*domain.User, error) {
	rows, err := r.db.Query(`
        SELECT u.id, u.school_id, u.email, u.password_hash, u.first_name, u.last_name,
            u.phone, u.role, u.avatar_url, u.status, u.created_at, u.updated_at
        FROM teacher_subjects ts
        JOIN users u ON ts.teacher_id = u.id
        WHERE ts.subject_id = $1 AND u.role = $2
        ORDER BY u.first_name, u.last_name`, subjectID, domain.RoleTeacher)
	if err != nil {
		return nil, fmt.Errorf("find subject teachers: %w", err)
	}
	defer rows.Close()

	userRepo := NewUserRepository(r.db)
	var teachers []*domain.User
	for rows.Next() {
		teacher := &domain.User{}
		if err := userRepo.ScanUserRow(rows, teacher); err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	return teachers, rows.Err()
}

func (r *teacherSubjectRepository) DeleteByTeacher(teacherID int) error {
	_, err := r.db.Exec(`DELETE FROM teacher_subjects WHERE teacher_id = $1`, teacherID)
	if err != nil {
		return fmt.Errorf("delete teacher subjects: %w", err)
	}
	return nil
}

func (r *teacherSubjectRepository) DeleteBySubject(subjectID int) error {
	_, err := r.db.Exec(`DELETE FROM teacher_subjects WHERE subject_id = $1`, subjectID)
	if err != nil {
		return fmt.Errorf("delete subject teachers: %w", err)
	}
	return nil
}
