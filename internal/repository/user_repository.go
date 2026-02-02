package repository

import (
	"database/sql"
	"fmt"

	"educnet/internal/domain"
)

// ! UserRepository interface (contrat)
type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id int) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	ExistsByEmail(email string) (bool, error)

	Update(user *domain.User) error
	UpdateAvatar(userID int, avatarURL string) error
	FindPendingBySchool(schoolID int) ([]*domain.User, error)
	FindBySchool(schoolID int, filters map[string]string) ([]*domain.User, error)
}

// ! userRepository implémentation
type userRepository struct {
	db *sql.DB
}

// ! NewUserRepository crée un nouveau repository
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

func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET first_name = $1, 
		    last_name = $2, 
		    phone = $3, 
		    avatar_url = $4, 
		    password_hash = $5,
		    status = $6, 
		    updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.AvatarURL,
		user.PasswordHash,
		user.Status,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) UpdateAvatar(userID int, avatarURL string) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, avatarURL, userID)
	return err
}

func (r *userRepository) FindPendingBySchool(schoolID int) ([]*domain.User, error) {
	query := `
		SELECT id, school_id, email, password_hash, first_name, last_name, 
		    phone, role, avatar_url, status, created_at, updated_at
		FROM users
		WHERE school_id = $1 AND status = 'pending'
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending users: %w", err)
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		var phone, avatarURL sql.NullString

		err := rows.Scan(
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
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if phone.Valid {
			user.Phone = phone.String
		}
		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) FindBySchool(schoolID int, filters map[string]string) ([]*domain.User, error) {
	query := `
		SELECT id, school_id, email, password_hash, first_name, last_name, 
		    phone, role, avatar_url, status, created_at, updated_at
		FROM users
		WHERE school_id = $1
	`

	args := []interface{}{schoolID}
	argCount := 1

	//! Applliquer la filtre
	if role, ok := filters["role"]; ok && role != "" {
		argCount++
		query += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, role)
	}

	if status, ok := filters["status"]; ok && status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		var phone, avatarURL sql.NullString

		err := rows.Scan(
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
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if phone.Valid {
			user.Phone = phone.String
		}
		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		}

		users = append(users, user)
	}

	return users, nil
}
