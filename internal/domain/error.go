package domain

import "fmt"

// ! DomainError représente une erreur métier
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// ! SCHOOL ERRORS
var (
	ErrSchoolNameRequired  = NewError("SCHOOL_NAME_REQUIRED", "School name is required")
	ErrSchoolSlugRequired  = NewError("SCHOOL_SLUG_REQUIRED", "School slug is required")
	ErrSchoolNotFound      = NewError("SCHOOL_NOT_FOUND", "School not found")
	ErrSchoolAlreadyExists = NewError("SCHOOL_ALREADY_EXISTS", "School with this name already exists")
)

// ! COMMON ERRORS
var (
	ErrNameRequired       = NewError("NAME_REQUIRED", "Name is required")
	ErrEmailRequired      = NewError("EMAIL_REQUIRED", "Email is required")
	ErrEmailInvalid       = NewError("EMAIL_INVALID", "Invalid email format")
	ErrEmailAlreadyExists = NewError("EMAIL_ALREADY_EXISTS", "Email already exists")
	ErrPasswordTooShort   = NewError("PASSWORD_TOO_SHORT", "Password must be at least 6 characters")
	ErrInvalidCredentials = NewError("INVALID_CREDENTIALS", "Invalid credentials")
	ErrInvalidRole        = NewError("INVALID_ROLE", "Invalid user role")
	ErrNotFound           = NewError("NOT_FOUND", "Resource not found")
	ErrInvalidPhoneFormat = NewError("INVALID_PHONE_FORMAT", "Phone number format is invalid")
)

// ! CLASS ERRORS
var (
	ErrClassNameRequired  = NewError("CLASS_NAME_REQUIRED", "Class name is required")
	ErrClassLevelRequired = NewError("CLASS_LEVEL_REQUIRED", "Class level is required")
	ErrClassYearRequired  = NewError("CLASS_YEAR_REQUIRED", "Academic year is required")
	ErrClassInvalidID     = NewError("CLASS_INVALID_ID", "Invalid school ID")
	ErrClassNotFound      = NewError("CLASS_NOT_FOUND", "Class not found")
)

// ! SUBJECT ERRORS
var (
	ErrSubjectNameRequired = NewError("SUBJECT_NAME_REQUIRED", "Subject name is required")
	ErrSubjectCodeRequired = NewError("SUBJECT_CODE_REQUIRED", "Subject code is required")
	ErrSubjectInvalidID    = NewError("SUBJECT_INVALID_ID", "Invalid school ID")
	ErrSubjectNotFound     = NewError("SUBJECT_NOT_FOUND", "Subject not found")
)

// ! USER ERRORS
var (
	ErrUserNotFound = NewError("USER_NOT_FOUND", "User not found")
)
