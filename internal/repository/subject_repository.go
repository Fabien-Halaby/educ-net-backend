package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"errors"
	"fmt"
)

type SubjectRepository interface {
	Create(subject *domain.Subject) error
	FindByID(id int) (*domain.Subject, error)
	FindBySchoolID(schoolID int) ([]*domain.Subject, error) // GetAll(schoolID)
	FindBySchoolAndCode(schoolID int, code string) (*domain.Subject, error)
	Update(subject *domain.Subject) error
	Delete(id int) error
	ExistsByCode(schoolID int, code string, excludeID int) (bool, error)

	//! HELPER
	ScanSubjectRow(row domainScanner, subjectObj *domain.Subject) error
}

type subjectRepository struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

// ! ==================== PRO SCANNER ====================
func (r *subjectRepository) ScanSubjectRow(row domainScanner, subjectObj *domain.Subject) error {
	var description sql.NullString
	err := row.Scan(
		&subjectObj.ID,
		&subjectObj.SchoolID,
		&subjectObj.Name,
		&subjectObj.Code,
		&description,
		&subjectObj.CreatedAt,
		&subjectObj.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return fmt.Errorf("scan subject row: %w", err)
	}

	subjectObj.Description = nullString(description)
	return nil
}

// ! ==================== METHODS PRO ====================
func (r *subjectRepository) Create(subject *domain.Subject) error {
	err := r.db.QueryRow(
		`INSERT INTO subjects (school_id, name, code, description)
         VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at`,
		subject.SchoolID, subject.Name, subject.Code, subject.Description,
	).Scan(&subject.ID, &subject.CreatedAt, &subject.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create subject: %w", err)
	}
	return nil
}

func (r *subjectRepository) FindByID(id int) (*domain.Subject, error) {
	subject := &domain.Subject{}
	query := `
        SELECT id, school_id, name, code, description, created_at, updated_at
        FROM subjects WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&subject.ID,
		&subject.SchoolID,
		&subject.Name,
		&subject.Code,
		new(sql.NullString),
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSubjectNotFound
		}
		return nil, fmt.Errorf("find subject by id %d: %w", id, err)
	}
	return subject, nil
}

func (r *subjectRepository) FindBySchoolID(schoolID int) ([]*domain.Subject, error) {
	rows, err := r.db.Query(
		`SELECT id,school_id,name,code,description,created_at,updated_at 
         FROM subjects WHERE school_id=$1 ORDER BY name`, schoolID)
	if err != nil {
		return nil, fmt.Errorf("find subjects by school: %w", err)
	}
	defer rows.Close()

	var subjects []*domain.Subject
	for rows.Next() {
		subjectObj := &domain.Subject{}
		if err := r.ScanSubjectRow(rows, subjectObj); err != nil {
			return nil, err
		}
		subjects = append(subjects, subjectObj)
	}
	return subjects, rows.Err()
}

func (r *subjectRepository) FindBySchoolAndCode(schoolID int, code string) (*domain.Subject, error) {
	subject := &domain.Subject{}
	query := `
        SELECT id, school_id, name, code, description, created_at, updated_at
        FROM subjects WHERE school_id=$1 AND code=$2`

	err := r.db.QueryRow(query, schoolID, code).Scan(
		&subject.ID,
		&subject.SchoolID,
		&subject.Name,
		&subject.Code,
		new(sql.NullString),
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSubjectNotFound
		}
		return nil, fmt.Errorf("find subject by school/code %s: %w", code, err)
	}
	return subject, nil
}

func (r *subjectRepository) Update(subject *domain.Subject) error {
	result, err := r.db.Exec(
		`UPDATE subjects SET name=$1,code=$2,description=$3,updated_at=NOW() 
         WHERE id=$4`,
		subject.Name, subject.Code, subject.Description, subject.ID)
	if err != nil {
		return fmt.Errorf("update subject: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrSubjectNotFound
	}
	return nil
}

func (r *subjectRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM subjects WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete subject: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrSubjectNotFound
	}
	return nil
}

func (r *subjectRepository) ExistsByCode(schoolID int, code string, excludeID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM subjects WHERE school_id=$1 AND code=$2 AND id!=$3)`,
		schoolID, code, excludeID).Scan(&exists)
	return exists, err
}
