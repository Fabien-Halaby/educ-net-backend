package repository

import (
	"database/sql"
	"fmt"

	"educnet/internal/domain"
)

// ! SchoolRepository interface (contrat)
type SchoolRepository interface {
	Create(school *domain.School) error
	Update(school *domain.School) error
	GetAll() ([]*domain.School, error)
	FindByID(id int) (*domain.School, error)
	FindBySlug(slug string) (*domain.School, error)
	ExistsBySlug(slug string) (bool, error)
	UpdateLogo(schoolID int, logoURL string) error
}

// ! schoolRepository implémentation
type schoolRepository struct {
	db *sql.DB
}

// ! NewSchoolRepository crée un nouveau repository
func NewSchoolRepository(db *sql.DB) SchoolRepository {
	return &schoolRepository{db: db}
}

func (r *schoolRepository) Create(school *domain.School) error {
	query := `
		INSERT INTO schools (name, slug, address, phone, email, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		school.Name,
		school.Slug,
		school.Address,
		school.Phone,
		school.Email,
		school.Status,
	).Scan(&school.ID, &school.CreatedAt, &school.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create school: %w", err)
	}

	return nil
}

func (r *schoolRepository) Update(school *domain.School) error {
	query := `
		UPDATE schools 
		SET name = $1, address = $2, phone = $3, email = $4, admin_user_id = $5, updated_at = NOW()
		WHERE id = $6
	`

	_, err := r.db.Exec(
		query,
		school.Name,
		school.Address,
		school.Phone,
		school.Email,
		school.AdminUserID,
		school.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update school: %w", err)
	}

	return nil
}

func (r *schoolRepository) FindByID(id int) (*domain.School, error) {
	query := `
		SELECT id, name, slug, address, phone, email, logo_url, admin_user_id, status, created_at, updated_at
		FROM schools
		WHERE id = $1
	`

	school := &domain.School{}
	var address, phone, email, logoURL sql.NullString
	var adminUserID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&school.ID,
		&school.Name,
		&school.Slug,
		&address,
		&phone,
		&email,
		&logoURL,
		&adminUserID,
		&school.Status,
		&school.CreatedAt,
		&school.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrSchoolNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find school: %w", err)
	}

	//! Convertir les NULL values
	if address.Valid {
		school.Address = address.String
	}
	if phone.Valid {
		school.Phone = phone.String
	}
	if email.Valid {
		school.Email = email.String
	}
	if logoURL.Valid {
		school.LogoURL = logoURL.String
	}
	if adminUserID.Valid {
		id := int(adminUserID.Int64)
		school.AdminUserID = &id
	}

	return school, nil
}

func (r *schoolRepository) FindBySlug(slug string) (*domain.School, error) {
	query := `
		SELECT id, name, slug, address, phone, email, logo_url, admin_user_id, status, created_at, updated_at
		FROM schools
		WHERE slug = $1
	`

	school := &domain.School{}
	var address, phone, email, logoURL sql.NullString
	var adminUserID sql.NullInt64

	err := r.db.QueryRow(query, slug).Scan(
		&school.ID,
		&school.Name,
		&school.Slug,
		&address,
		&phone,
		&email,
		&logoURL,
		&adminUserID,
		&school.Status,
		&school.CreatedAt,
		&school.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrSchoolNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find school: %w", err)
	}

	// Convertir les NULL values
	if address.Valid {
		school.Address = address.String
	}
	if phone.Valid {
		school.Phone = phone.String
	}
	if email.Valid {
		school.Email = email.String
	}
	if logoURL.Valid {
		school.LogoURL = logoURL.String
	}
	if adminUserID.Valid {
		id := int(adminUserID.Int64)
		school.AdminUserID = &id
	}

	return school, nil
}

func (r *schoolRepository) ExistsBySlug(slug string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM schools WHERE slug = $1)"

	err := r.db.QueryRow(query, slug).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check slug existence: %w", err)
	}

	return exists, nil
}

// ! UpdateLogo met à jour l'URL du logo de l'école
func (r *schoolRepository) UpdateLogo(schoolID int, logoURL string) error {
	query := `UPDATE schools SET logo_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, logoURL, schoolID)
	return err
}

// ! GetAll returns all schools in the database
func (r *schoolRepository) GetAll() ([]*domain.School, error) {
	query := `
		SELECT id, name, slug, address, phone, email, logo_url, admin_user_id, status, created_at, updated_at
		FROM schools
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all schools: %w", err)
	}
	defer rows.Close()

	var schools []*domain.School

	for rows.Next() {
		school := &domain.School{}
		var address, phone, email, logoURL sql.NullString
		var adminUserID sql.NullInt64

		err := rows.Scan(
			&school.ID,
			&school.Name,
			&school.Slug,
			&address,
			&phone,
			&email,
			&logoURL,
			&adminUserID,
			&school.Status,
			&school.CreatedAt,
			&school.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan school: %w", err)
		}

		//! Convertir les NULL values
		if address.Valid {
			school.Address = address.String
		}
		if phone.Valid {
			school.Phone = phone.String
		}
		if email.Valid {
			school.Email = email.String
		}
		if logoURL.Valid {
			school.LogoURL = logoURL.String
		}
		if adminUserID.Valid {
			id := int(adminUserID.Int64)
			school.AdminUserID = &id
		}

		schools = append(schools, school)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return schools, nil
}
