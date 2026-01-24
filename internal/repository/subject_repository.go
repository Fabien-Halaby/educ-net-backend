package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"fmt"
	"time"
)

type SubjectRepository interface {
	Create(subject *domain.Subject) error
	FindByID(id int) (*domain.Subject, error)
	FindBySchoolID(schoolID int) ([]*domain.Subject, error)
	FindBySchoolAndCode(schoolID int, code string) (*domain.Subject, error)

	Update(subject *domain.Subject) error
	Delete(id int) error        
	ExistsByCode(schoolID int, code string, excludeID int) (bool, error)
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



//! Update met à jour un sujet existant
func (r *subjectRepository) Update(subject *domain.Subject) error {
	query := `
		UPDATE subjects
		SET name = $1, code = $2, description = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(
		query,
		subject.Name,
		subject.Code,
		subject.Description,
		time.Now(),
		subject.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update subject: %w", err)
	}

	return nil
}

//! Delete supprime un sujet par son ID
func (r *subjectRepository) Delete(id int) error {
	query := `DELETE FROM subjects WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subject: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrSubjectNotFound
	}

	return nil
}

//! ExistsByCode vérifie si un code de sujet existe déjà dans une école, en excluant un ID spécifique
func (r *subjectRepository) ExistsByCode(schoolID int, code string, excludeID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM subjects 
			WHERE school_id = $1 AND code = $2 AND id != $3
		)
	`

	var exists bool
	err := r.db.QueryRow(query, schoolID, code, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check subject code: %w", err)
	}

	return exists, nil
}
