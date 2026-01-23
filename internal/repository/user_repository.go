package repository

import (
	"database/sql"
	"fmt"

	"educnet/internal/domain"
)

//! UserRepository interface (contrat)
type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id int) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	ExistsByEmail(email string) (bool, error)
}

//! userRepository implémentation
type userRepository struct {
	db *sql.DB
}

//! NewUserRepository crée un nouveau repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (school_id, email, password_hash, first_name, last_name, phone, role, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.SchoolID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.Status,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) FindByID(id int) (*domain.User, error) {
	query := `
		SELECT id, school_id, email, password_hash, first_name, last_name, phone, role, avatar_url, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	var phone, avatarURL sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.SchoolID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&phone,
		&user.Role,
		&avatarURL,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if phone.Valid {
		user.Phone = phone.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, school_id, email, password_hash, first_name, last_name, phone, role, avatar_url, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	var phone, avatarURL sql.NullString

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.SchoolID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&phone,
		&user.Role,
		&avatarURL,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if phone.Valid {
		user.Phone = phone.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	return user, nil
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
