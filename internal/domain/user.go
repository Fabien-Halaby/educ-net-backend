package domain

import (
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//! User représente l'entité métier Utilisateur
type User struct {
	ID           int `json:"id"`
	SchoolID     int `json:"school_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"lat_name"`
	Phone        string `json:"phone"`
	Role         string  `json:"role"`	//! admin, teacher, student, parent
	AvatarURL    string `json:"avatar_url"`
	Status       string  `json:"status"`	//! pending, approved, rejected, inactive
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

//! NewUser crée un nouvel utilisateur avec validation
func NewUser(schoolID int, email, password, firstName, lastName, phone, role string) (*User, error) {
	if email == "" {
		return nil, ErrUserEmailRequired
	}
	if !isValidEmail(email) {
		return nil, ErrUserEmailInvalid
	}

	if password == "" || len(password) < 6 {
		return nil, ErrUserPasswordTooShort
	}

	if firstName == "" {
		return nil, ErrUserNameRequired
	}

	if !isValidRole(role) {
		return nil, ErrUserInvalidRole
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
		Status:       "pending", //! Par défaut pending
	}, nil
}

//! NewAdminUser crée un admin (directement approved)
func NewAdminUser(schoolID int, email, password, firstName, lastName, phone string) (*User, error) {
	user, err := NewUser(schoolID, email, password, firstName, lastName, phone, "admin")
	if err != nil {
		return nil, err
	}
	user.Status = "approved"
	return user, nil
}

//! VerifyPassword vérifie le mot de passe
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

//! Approve approuve l'utilisateur
func (u *User) Approve() {
	u.Status = "approved"
}

//! Reject rejette l'utilisateur
func (u *User) Reject() {
	u.Status = "rejected"
}

//! IsApproved vérifie si l'utilisateur est approuvé
func (u *User) IsApproved() bool {
	return u.Status == "approved"
}

//! IsAdmin vérifie si l'utilisateur est admin
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

//! GetFullName retourne le nom complet
func (u *User) GetFullName() string {
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

//! Private helpers
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
