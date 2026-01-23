package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"fmt"
)

type SubjectRepository interface {
	Create(subject *domain.Subject) error
	FindByID(id int) (*domain.Subject, error)
	FindBySchoolID(schoolID int) ([]*domain.Subject, error)
	FindBySchoolAndCode(schoolID int, code string) (*domain.Subject, error)
}

type subjectRepository struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(subject *domain.Subject) error {
	query := `
		INSERT INTO subjects (school_id, name, code, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		subject.SchoolID,
		subject.Name,
		subject.Code,
		subject.Description,
	).Scan(&subject.ID, &subject.CreatedAt, &subject.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create subject: %w", err)
	}

	return nil
}

func (r *subjectRepository) FindByID(id int) (*domain.Subject, error) {
	query := `
		SELECT id, school_id, name, code, description, created_at, updated_at
		FROM subjects
		WHERE id = $1
	`

	subject := &domain.Subject{}
	var description sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&subject.ID,
		&subject.SchoolID,
		&subject.Name,
		&subject.Code,
		&description,
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrSubjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find subject: %w", err)
	}

	if description.Valid {
		subject.Description = description.String
	}

	return subject, nil
}

func (r *subjectRepository) FindBySchoolID(schoolID int) ([]*domain.Subject, error) {
	query := `
		SELECT id, school_id, name, code, description, created_at, updated_at
		FROM subjects
		WHERE school_id = $1
		ORDER BY name
	`

	rows, err := r.db.Query(query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to find subjects: %w", err)
	}
	defer rows.Close()

	subjects := []*domain.Subject{}
	for rows.Next() {
		subject := &domain.Subject{}
		var description sql.NullString

		err := rows.Scan(
			&subject.ID,
			&subject.SchoolID,
			&subject.Name,
			&subject.Code,
			&description,
			&subject.CreatedAt,
			&subject.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subject: %w", err)
		}

		if description.Valid {
			subject.Description = description.String
		}

		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (r *subjectRepository) FindBySchoolAndCode(schoolID int, code string) (*domain.Subject, error) {
	query := `
		SELECT id, school_id, name, code, description, created_at, updated_at
		FROM subjects
		WHERE school_id = $1 AND code = $2
	`

	subject := &domain.Subject{}
	var description sql.NullString

	err := r.db.QueryRow(query, schoolID, code).Scan(
		&subject.ID,
		&subject.SchoolID,
		&subject.Name,
		&subject.Code,
		&description,
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrSubjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find subject: %w", err)
	}

	if description.Valid {
		subject.Description = description.String
	}

	return subject, nil
}
