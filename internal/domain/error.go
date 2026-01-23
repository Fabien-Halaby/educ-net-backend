package domain

import "fmt"

//! DomainError représente une erreur métier
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"mesage"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

//! NewError crée une nouvelle erreur domain
func NewError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

//! School errors
var (
	ErrSchoolNameRequired     = NewError("SCHOOL_NAME_REQUIRED", "School name is required")
	ErrSchoolSlugRequired     = NewError("SCHOOL_SLUG_REQUIRED", "School slug is required")
	ErrSchoolNotFound         = NewError("SCHOOL_NOT_FOUND", "School not found")
	ErrSchoolAlreadyExists    = NewError("SCHOOL_ALREADY_EXISTS", "School with this name already exists")
)

//! User errors
var (
	ErrUserEmailRequired      = NewError("USER_EMAIL_REQUIRED", "Email is required")
	ErrUserEmailInvalid       = NewError("USER_EMAIL_INVALID", "Email format is invalid")
	ErrUserPasswordTooShort   = NewError("USER_PASSWORD_TOO_SHORT", "Password must be at least 6 characters")
	ErrUserNameRequired       = NewError("USER_NAME_REQUIRED", "Name is required")
	ErrUserInvalidRole        = NewError("USER_INVALID_ROLE", "Invalid role")
	ErrUserNotFound           = NewError("USER_NOT_FOUND", "User not found")
	ErrUserAlreadyExists      = NewError("USER_ALREADY_EXISTS", "User with this email already exists")
	ErrUserInvalidPassword    = NewError("USER_INVALID_PASSWORD", "Invalid password")
)
