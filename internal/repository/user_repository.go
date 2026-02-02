package repository

import (
	"database/sql"
	"educnet/internal/domain"
	"fmt"
)

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

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// ! ==================== PRO HELPERS ====================
func (r *userRepository) scanUserRow(row domainScanner, user *domain.User) error {
	var phone, avatarURL sql.NullString
	err := row.Scan(
		&user.ID, &user.SchoolID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &phone, &user.Role,
		&avatarURL, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return err // ← Passer directement
	}
	if err != nil {
		return fmt.Errorf("scan user row: %w", err)
	}

	user.Phone = nullString(phone)
	user.AvatarURL = nullString(avatarURL)
	return nil
}

// ! ==================== METHODS PRO ====================
func (r *userRepository) Create(user *domain.User) error {
	err := r.db.QueryRow(
		`INSERT INTO users (school_id,email,password_hash,first_name,last_name,phone,role,status)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id,created_at,updated_at`,
		user.SchoolID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.Role, user.Status,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(id int) (*domain.User, error) {
	user := &domain.User{}
	row := r.db.QueryRow(
		`SELECT id,school_id,email,password_hash,first_name,last_name,phone,role,avatar_url,status,created_at,updated_at 
         FROM users WHERE id=$1`, id)

	if err := r.scanUserRow(row, user); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	row := r.db.QueryRow(
		`SELECT id,school_id,email,password_hash,first_name,last_name,phone,role,avatar_url,status,created_at,updated_at 
         FROM users WHERE email=$1`, email)

	if err := r.scanUserRow(row, user); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}

func (r *userRepository) Update(user *domain.User) error {
	result, err := r.db.Exec(
		`UPDATE users SET first_name=$1,last_name=$2,phone=$3,avatar_url=$4,password_hash=$5,status=$6,updated_at=NOW() 
         WHERE id=$7`,
		user.FirstName, user.LastName, user.Phone, user.AvatarURL, user.PasswordHash, user.Status, user.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdateAvatar(userID int, avatarURL string) error {
	result, err := r.db.Exec(`UPDATE users SET avatar_url=$1,updated_at=NOW() WHERE id=$2`, avatarURL, userID)
	if err != nil {
		return fmt.Errorf("update user avatar: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) FindPendingBySchool(schoolID int) ([]*domain.User, error) {
	rows, err := r.db.Query(
		`SELECT id,school_id,email,password_hash,first_name,last_name,phone,role,avatar_url,status,created_at,updated_at 
         FROM users WHERE school_id=$1 AND status='pending' ORDER BY created_at DESC`, schoolID)
	if err != nil {
		return nil, fmt.Errorf("find pending users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := r.scanUserRow(rows, user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *userRepository) FindBySchool(schoolID int, filters map[string]string) ([]*domain.User, error) {
	query := `SELECT id,school_id,email,password_hash,first_name,last_name,phone,role,avatar_url,status,created_at,updated_at 
              FROM users WHERE school_id=$1`
	args := []interface{}{schoolID}

	if role, ok := filters["role"]; ok && role != "" {
		query += ` AND role=$2`
		args = append(args, role)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query += ` AND status=$` + fmt.Sprint(len(args)+1)
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("find users by school: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := r.scanUserRow(rows, user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
