package domain

import (
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ! User représente l'entité métier Utilisateur
type User struct {
	ID           int       `json:"id"`
	SchoolID     int       `json:"school_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"` //! admin, teacher, student, parent
	AvatarURL    string    `json:"avatar_url"`
	Status       string    `json:"status"` //! pending, approved, rejected, inactive
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ! Role constants
const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
	RoleParent  = "parent"
)

// ! User status constants
const (
	UserStatusPending   = "pending"
	UserStatusApproved  = "approved"
	UserStatusRejected  = "rejected"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

// ! NewUser crée un nouvel utilisateur avec validation
func NewUser(schoolID int, email, password, firstName, lastName, phone, role string) (*User, error) {
	if email == "" {
		return nil, ErrEmailRequired
	}
	if !isValidEmail(email) {
		return nil, ErrEmailInvalid
	}

	if password == "" || len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	if firstName == "" {
		return nil, ErrNameRequired
	}

	if !isValidRole(role) {
		return nil, ErrInvalidRole
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		SchoolID:     schoolID,
		Email:        email,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        phone,
		Role:         role,
		Status:       UserStatusPending,
	}, nil
}

// ! NewAdminUser crée un admin (directement approved)
func NewAdminUser(schoolID int, email, password, firstName, lastName, phone string) (*User, error) {
	user, err := NewUser(schoolID, email, password, firstName, lastName, phone, "admin")
	if err != nil {
		return nil, err
	}
	user.Status = UserStatusApproved
	user.UpdatedAt = time.Now()
	return user, nil
}

// ! VerifyPassword vérifie le mot de passe
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// ! Approve approuve l'utilisateur
func (u *User) Approve() {
	u.Status = UserStatusApproved
	u.UpdatedAt = time.Now()
}

// ! Reject rejette l'utilisateur
func (u *User) Reject() {
	u.Status = UserStatusRejected
	u.UpdatedAt = time.Now()
}

// ! Suspend deactivates the user account
func (u *User) Suspend() {
	u.Status = UserStatusSuspended
	u.UpdatedAt = time.Now()
}

// ! Deactivate désactive l'utilisateur
func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
}

// ! IsApproved vérifie si l'utilisateur est approuvé
func (u *User) IsApproved() bool {
	return u.Status == UserStatusApproved
}

// ! IsPending checks if user is waiting for approval
func (u *User) IsPending() bool {
	return u.Status == UserStatusPending
}

// ! IsRejected vérifie si l'utilisateur est rejeté
func (u *User) IsRejected() bool {
	return u.Status == UserStatusRejected
}

// ! IsAdmin vérifie si l'utilisateur est admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// ! IsTeacher checks if user is a teacher
func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

// ! IsStudent checks if user is a student
func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

// ! GetFullName retourne le nom complet
func (u *User) GetFullName() string {
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// ! Private helpers
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidRole(role string) bool {
	validRoles := map[string]bool{
		"admin":   true,
		"teacher": true,
		"student": true,
		"parent":  true,
	}
	return validRoles[role]
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// SetPassword hash et définit un nouveau mot de passe
func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.PasswordHash = string(hashedPassword)
	u.UpdatedAt = time.Now()

	return nil
}
