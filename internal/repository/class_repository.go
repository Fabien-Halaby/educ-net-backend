package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"fmt"
)

type ClassRepository interface {
	Create(class *domain.Class) error
	FindByID(id int) (*domain.Class, error)
	FindBySchoolID(schoolID int) ([]*domain.Class, error)
	FindBySchoolAndYear(schoolID int, academicYear string) ([]*domain.Class, error)
}

type classRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(class *domain.Class) error {
	query := `
		INSERT INTO classes (school_id, name, level, section, capacity, academic_year)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		class.SchoolID,
		class.Name,
		class.Level,
		class.Section,
		class.Capacity,
		class.AcademicYear,
	).Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create class: %w", err)
	}

	return nil
}

func (r *classRepository) FindByID(id int) (*domain.Class, error) {
	query := `
		SELECT id, school_id, name, level, section, capacity, academic_year, created_at, updated_at
		FROM classes
		WHERE id = $1
	`

	class := &domain.Class{}
	var section sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&class.ID,
		&class.SchoolID,
		&class.Name,
		&class.Level,
		&section,
		&class.Capacity,
		&class.AcademicYear,
		&class.CreatedAt,
		&class.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrClassNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find class: %w", err)
	}

	if section.Valid {
		class.Section = section.String
	}

	return class, nil
}

func (r *classRepository) FindBySchoolID(schoolID int) ([]*domain.Class, error) {
	query := `
		SELECT id, school_id, name, level, section, capacity, academic_year, created_at, updated_at
		FROM classes
		WHERE school_id = $1
		ORDER BY level, name
	`

	rows, err := r.db.Query(query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to find classes: %w", err)
	}
	defer rows.Close()

	classes := []*domain.Class{}
	for rows.Next() {
		class := &domain.Class{}
		var section sql.NullString

		err := rows.Scan(
			&class.ID,
			&class.SchoolID,
			&class.Name,
			&class.Level,
			&section,
			&class.Capacity,
			&class.AcademicYear,
			&class.CreatedAt,
			&class.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan class: %w", err)
		}

		if section.Valid {
			class.Section = section.String
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func (r *classRepository) FindBySchoolAndYear(schoolID int, academicYear string) ([]*domain.Class, error) {
	query := `
		SELECT id, school_id, name, level, section, capacity, academic_year, created_at, updated_at
		FROM classes
		WHERE school_id = $1 AND academic_year = $2
		ORDER BY level, name
	`

	rows, err := r.db.Query(query, schoolID, academicYear)
	if err != nil {
		return nil, fmt.Errorf("failed to find classes: %w", err)
	}
	defer rows.Close()

	classes := []*domain.Class{}
	for rows.Next() {
		class := &domain.Class{}
		var section sql.NullString

		err := rows.Scan(
			&class.ID,
			&class.SchoolID,
			&class.Name,
			&class.Level,
			&section,
			&class.Capacity,
			&class.AcademicYear,
			&class.CreatedAt,
			&class.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan class: %w", err)
		}

		if section.Valid {
			class.Section = section.String
		}

		classes = append(classes, class)
	}

	return classes, nil
}
