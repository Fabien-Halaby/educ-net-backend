package dto

import "educnet/internal/domain"

// UpdateProfileRequest pour modifier le profil
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// ChangePasswordRequest pour changer le mot de passe
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// ProfileResponse réponse complète du profil
type ProfileResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	SchoolID  int    `json:"school_id"`
	AvatarURL string `json:"avatar_url,omitempty"`
	CreatedAt string `json:"created_at"`
}

// TeacherSubjectsResponse liste des matières d'un enseignant
type TeacherSubjectsResponse struct {
	Subjects []SubjectInfo `json:"subjects"`
	Total    int           `json:"total"`
}

type SubjectInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
}

func SubjectInfoFromDomain(subject *domain.Subject) *SubjectInfo {
	return &SubjectInfo{
		ID:          subject.ID,
		Name:        subject.Name,
		Code:        subject.Code,
		Description: subject.Description,
	}
}

// StudentClassResponse classe d'un étudiant
type StudentClassesResponse struct {
	Classes []*ClassInfo `json:"classes"`
	Total   int          `json:"total"`
}

type ClassInfo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Level        string `json:"level"`
	Section      string `json:"section,omitempty"`
	Capacity     int    `json:"capacity"`
	AcademicYear string `json:"academic_year"`
}

func ClassInfoFromDomain(class *domain.Class) *ClassInfo {
	return &ClassInfo{
		ID:           class.ID,
		Name:         class.Name,
		Level:        class.Level,
		Section:      class.Section,
		Capacity:     class.Capacity,
		AcademicYear: class.AcademicYear,
	}
}

// ! School and Upload DTOs
type UpdateSchoolRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type UploadResponse struct {
	URL string `json:"url"`
}
