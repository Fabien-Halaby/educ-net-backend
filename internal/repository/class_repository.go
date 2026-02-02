package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"errors"
	"fmt"
)

type ClassRepository interface {
	Create(class *domain.Class) error
	FindByID(id int) (*domain.Class, error)
	FindBySchoolID(schoolID int) ([]*domain.Class, error)
	GetAll(schoolID int) ([]*domain.Class, error)
	FindBySchoolAndYear(schoolID int, academicYear string) ([]*domain.Class, error)
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

// ! ==================== PRO SCANNER ====================
func (r *classRepository) scanClassRow(row domainScanner, classObj *domain.Class) error {
	var section sql.NullString
	err := row.Scan(
		&classObj.ID, &classObj.SchoolID, &classObj.Name, &classObj.Level,
		&section, &classObj.Capacity, &classObj.AcademicYear,
		&classObj.CreatedAt, &classObj.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return fmt.Errorf("scan class row: %w", err)
	}

	classObj.Section = nullString(section)
	return nil
}

// ! ==================== METHODS PRO ====================
func (r *classRepository) Create(class *domain.Class) error {
	err := r.db.QueryRow(
		`INSERT INTO classes (school_id,name,level,section,capacity,academic_year,status)
        VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id,created_at,updated_at`,
		class.SchoolID, class.Name, class.Level, class.Section, class.Capacity,
		class.AcademicYear, class.Status,
	).Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create class: %w", err)
	}
	return nil
}

func (r *classRepository) FindByID(id int) (*domain.Class, error) {
	class := &domain.Class{}
	query := `
        SELECT id, school_id, name, level, section, capacity, academic_year, 
               created_at, updated_at
        FROM classes WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&class.ID,
		&class.SchoolID,
		&class.Name,
		&class.Level,
		&class.Section,
		&class.Capacity,
		&class.AcademicYear,
		&class.CreatedAt,
		&class.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrClassNotFound
		}
		return nil, fmt.Errorf("find class by id %d: %w", id, err)
	}
	return class, nil
}

func (r *classRepository) FindBySchoolID(schoolID int) ([]*domain.Class, error) {
	rows, err := r.db.Query(
		`SELECT id,school_id,name,level,section,capacity,academic_year,created_at,updated_at 
         FROM classes WHERE school_id=$1 ORDER BY level,name`, schoolID)
	if err != nil {
		return nil, fmt.Errorf("find classes by school: %w", err)
	}
	defer rows.Close()

	var classes []*domain.Class
	for rows.Next() {
		classObj := &domain.Class{}
		if err := r.scanClassRow(rows, classObj); err != nil {
			return nil, err
		}
		classes = append(classes, classObj)
	}
	return classes, rows.Err()
}

func (r *classRepository) GetAll(schoolID int) ([]*domain.Class, error) {
	return r.FindBySchoolID(schoolID)
}

func (r *classRepository) FindBySchoolAndYear(schoolID int, academicYear string) ([]*domain.Class, error) {
	rows, err := r.db.Query(
		`SELECT id,school_id,name,level,section,capacity,academic_year,created_at,updated_at 
        FROM classes WHERE school_id=$1 AND academic_year=$2 ORDER BY level,name`,
		schoolID, academicYear)
	if err != nil {
		return nil, fmt.Errorf("find classes by school/year: %w", err)
	}
	defer rows.Close()

	var classes []*domain.Class
	for rows.Next() {
		classObj := &domain.Class{}
		if err := r.scanClassRow(rows, classObj); err != nil {
			return nil, err
		}
		classes = append(classes, classObj)
	}
	return classes, rows.Err()
}

func (r *classRepository) Update(class *domain.Class) error {
	result, err := r.db.Exec(
		`UPDATE classes SET name=$1,level=$2,section=$3,capacity=$4,academic_year=$5,status=$6,updated_at=NOW() 
         WHERE id=$7`,
		class.Name, class.Level, class.Section, class.Capacity, class.AcademicYear, class.Status, class.ID)
	if err != nil {
		return fmt.Errorf("update class: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrClassNotFound
	}
	return nil
}

func (r *classRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM classes WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete class: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrClassNotFound
	}
	return nil
}

func (r *classRepository) ExistsByName(schoolID int, name string, excludeID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM classes WHERE school_id=$1 AND name=$2 AND id!=$3)`,
		schoolID, name, excludeID).Scan(&exists)
	return exists, err
}
