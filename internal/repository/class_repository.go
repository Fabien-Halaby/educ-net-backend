package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"fmt"
	"time"
)

type ClassRepository interface {
	Create(class *domain.Class) error
	FindByID(id int) (*domain.Class, error)
	FindBySchoolID(schoolID int) ([]*domain.Class, error)
	FindBySchoolAndYear(schoolID int, academicYear string) ([]*domain.Class, error)

	GetAll(schoolID int) ([]*domain.Class, error)
	Update(class *domain.Class) error
	Delete(id int) error
	ExistsByName(schoolID int, name string, excludeID int) (bool, error)
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

func (r *classRepository) GetAll(schoolID int) ([]*domain.Class, error) {
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

// ! Update met à jour les informations d'une classe existante.
func (r *classRepository) Update(class *domain.Class) error {
	query := `
		UPDATE classes
		SET name = $1, level = $2, section = $3, capacity = $4, 
		    academic_year = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.Exec(
		query,
		class.Name,
		class.Level,
		class.Section,
		class.Capacity,
		class.AcademicYear,
		time.Now(),
		class.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update class: %w", err)
	}

	return nil
}

// ! Delete supprime une classe par son ID.
func (r *classRepository) Delete(id int) error {
	query := `DELETE FROM classes WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrClassNotFound
	}

	return nil
}

func (r *classRepository) ExistsByName(schoolID int, name string, excludeID int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM classes 
			WHERE school_id = $1 AND name = $2 AND id != $3
		)
	`

	err := r.db.QueryRow(query, schoolID, name, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check class name existence: %w", err)
	}

	return exists, nil
}
