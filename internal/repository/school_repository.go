package repository

import (
	"database/sql"
	"fmt"

	"educnet/internal/domain"
)

//! SchoolRepository interface (contrat)
type SchoolRepository interface {
	Create(school *domain.School) error
	Update(school *domain.School) error
	FindByID(id int) (*domain.School, error)
	FindBySlug(slug string) (*domain.School, error)
	ExistsBySlug(slug string) (bool, error)
}

//! schoolRepository implémentation
type schoolRepository struct {
	db *sql.DB
}

//! NewSchoolRepository crée un nouveau repository
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
		SET name = $1, address = $2, phone = $3, admin_user_id = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(
		query,
		school.Name,
		school.Address,
		school.Phone,
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
