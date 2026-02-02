package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"fmt"
)

type StudentClassRepository interface {
	Create(studentID, classID int) error
	Delete(studentID, classID int) error
	Exists(studentID, classID int) (bool, error)
	FindByStudent(studentID int) ([]*domain.Class, error)
	FindByClass(classID int) ([]*domain.User, error)
	DeleteByStudent(studentID int) error
	DeleteByClass(classID int) error
}

type studentClassRepository struct {
	db *sql.DB
}

func NewStudentClassRepository(db *sql.DB) StudentClassRepository {
	return &studentClassRepository{db: db}
}

// ! ==================== METHODS PRO ====================
func (r *studentClassRepository) Create(studentID, classID int) error {
	_, err := r.db.Exec(
		`INSERT INTO student_classes (student_id, class_id) VALUES ($1, $2)`,
		studentID, classID,
	)
	if err != nil {
		return fmt.Errorf("create student-class: %w", err)
	}
	return nil
}

func (r *studentClassRepository) Delete(studentID, classID int) error {
	result, err := r.db.Exec(
		`DELETE FROM student_classes WHERE student_id=$1 AND class_id=$2`,
		studentID, classID)
	if err != nil {
		return fmt.Errorf("delete student-class: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrStudentClassNotFound
	}
	return nil
}

func (r *studentClassRepository) Exists(studentID, classID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM student_classes WHERE student_id=$1 AND class_id=$2)`,
		studentID, classID).Scan(&exists)
	return exists, err
}

func (r *studentClassRepository) FindByStudent(studentID int) ([]*domain.Class, error) {
	rows, err := r.db.Query(`
        SELECT c.id, c.school_id, c.name, c.level, c.section, c.capacity, c.academic_year, 
            c.created_at, c.updated_at
        FROM student_classes sc
        JOIN classes c ON sc.class_id = c.id 
        WHERE sc.student_id = $1 
        ORDER BY c.name`, studentID)
	if err != nil {
		return nil, fmt.Errorf("find student classes: %w", err)
	}
	defer rows.Close()

	classRepo := NewClassRepository(r.db)
	var classes []*domain.Class
	for rows.Next() {
		class := &domain.Class{}
		if err := classRepo.ScanClassRow(rows, class); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, rows.Err()
}

func (r *studentClassRepository) FindByClass(classID int) ([]*domain.User, error) {
	rows, err := r.db.Query(`
        SELECT u.id, u.school_id, u.email, u.password_hash, u.first_name, u.last_name, 
            u.phone, u.role, u.avatar_url, u.status, u.created_at, u.updated_at
        FROM student_classes sc
        JOIN users u ON sc.student_id = u.id 
        WHERE sc.class_id = $1 AND u.role = $2 
        ORDER BY u.first_name, u.last_name`, classID, domain.RoleStudent)
	if err != nil {
		return nil, fmt.Errorf("find class students: %w", err)
	}
	defer rows.Close()

	userRepo := NewUserRepository(r.db)
	var students []*domain.User
	for rows.Next() {
		student := &domain.User{}
		if err := userRepo.ScanUserRow(rows, student); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, rows.Err()
}

func (r *studentClassRepository) DeleteByStudent(studentID int) error {
	_, err := r.db.Exec(`DELETE FROM student_classes WHERE student_id = $1`, studentID)
	if err != nil {
		return fmt.Errorf("delete student classes: %w", err)
	}
	return nil
}

func (r *studentClassRepository) DeleteByClass(classID int) error {
	_, err := r.db.Exec(`DELETE FROM student_classes WHERE class_id = $1`, classID)
	if err != nil {
		return fmt.Errorf("delete class students: %w", err)
	}
	return nil
}
