package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"fmt"
)

// ! SchoolRepository interface (PARFAIT)
type SchoolRepository interface {
	Create(school *domain.School) error
	Update(school *domain.School) error
	GetAll() ([]*domain.School, error)
	FindByID(id int) (*domain.School, error)
	FindBySlug(slug string) (*domain.School, error)
	ExistsBySlug(slug string) (bool, error)
	UpdateLogo(schoolID int, logoURL string) error
}

type schoolRepository struct {
	db *sql.DB
}

func NewSchoolRepository(db *sql.DB) SchoolRepository {
	return &schoolRepository{db: db}
}

// ! ==================== HELPERS ====================
func (r *schoolRepository) scanSchoolRow(row domainScanner, school *domain.School) error {
	var address, phone, email, logoURL sql.NullString
	var adminUserID sql.NullInt64

	err := row.Scan(
		&school.ID, &school.Name, &school.Slug,
		&address, &phone, &email, &logoURL,
		&adminUserID, &school.Status,
		&school.CreatedAt, &school.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return fmt.Errorf("scan school row: %w", err)
	}

	school.Address = nullString(address)
	school.Phone = nullString(phone)
	school.Email = nullString(email)
	school.LogoURL = nullString(logoURL)
	school.AdminUserID = nullInt(adminUserID)

	return nil
}

// ! ==================== METHODS ====================
func (r *schoolRepository) Create(school *domain.School) error {
	err := r.db.QueryRow(
		`INSERT INTO schools (name, slug, address, phone, email, status) 
        VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`,
		school.Name, school.Slug, school.Address, school.Phone, school.Email, school.Status,
	).Scan(&school.ID, &school.CreatedAt, &school.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create school: %w", err)
	}
	return nil
}

func (r *schoolRepository) FindByID(id int) (*domain.School, error) {
	school := &domain.School{}
	row := r.db.QueryRow(
		`SELECT id,name,slug,address,phone,email,logo_url,admin_user_id,status,created_at,updated_at 
        FROM schools WHERE id = $1`, id)

	if err := r.scanSchoolRow(row, school); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSchoolNotFound
		}
		return nil, err
	}
	return school, nil
}

func (r *schoolRepository) FindBySlug(slug string) (*domain.School, error) {
	school := &domain.School{}
	row := r.db.QueryRow(
		`SELECT id,name,slug,address,phone,email,logo_url,admin_user_id,status,created_at,updated_at 
        FROM schools WHERE slug = $1`, slug)

	if err := r.scanSchoolRow(row, school); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSchoolNotFound
		}
		return nil, err
	}
	return school, nil
}

func (r *schoolRepository) GetAll() ([]*domain.School, error) {
	rows, err := r.db.Query(
		`SELECT id,name,slug,address,phone,email,logo_url,admin_user_id,status,created_at,updated_at 
        FROM schools`)
	if err != nil {
		return nil, fmt.Errorf("get all schools: %w", err)
	}
	defer rows.Close()

	var schools []*domain.School
	for rows.Next() {
		school := &domain.School{}
		if err := r.scanSchoolRow(rows, school); err != nil {
			return nil, err
		}
		schools = append(schools, school)
	}
	return schools, rows.Err()
}

func (r *schoolRepository) ExistsBySlug(slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM schools WHERE slug = $1)`, slug).Scan(&exists)
	return exists, err
}

func (r *schoolRepository) Update(school *domain.School) error {
	result, err := r.db.Exec(
		`UPDATE schools SET name=$1,address=$2,phone=$3,email=$4,admin_user_id=$5,updated_at=NOW() 
         WHERE id=$6`,
		school.Name, school.Address, school.Phone, school.Email, school.AdminUserID, school.ID)
	if err != nil {
		return fmt.Errorf("update school: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrSchoolNotFound
	}
	return nil
}

func (r *schoolRepository) UpdateLogo(schoolID int, logoURL string) error {
	result, err := r.db.Exec(`UPDATE schools SET logo_url=$1,updated_at=NOW() WHERE id=$2`, logoURL, schoolID)
	if err != nil {
		return fmt.Errorf("update school logo: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrSchoolNotFound
	}
	return nil
}
