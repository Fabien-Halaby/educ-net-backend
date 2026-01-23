package domain

import "time"

//! School représente l'entité métier École
type School struct {
	ID          int `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	LogoURL     string `json:"logo_url"`
	AdminUserID *int `json:"admin_user_id"`
	Status      string `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//! NewSchool crée une nouvelle école avec validation
func NewSchool(name, slug, address, email, phone string) (*School, error) {
	if name == "" {
		return nil, ErrSchoolNameRequired
	}
	if slug == "" {
		return nil, ErrSchoolSlugRequired
	}

	return &School{
		Name:    name,
		Slug:    slug,
		Address: address,
		Phone:   phone,
		Email:   email,
		Status:  "active",
	}, nil
}

//! SetAdmin définit l'admin de l'école
func (s *School) SetAdmin(adminID int) {
	s.AdminUserID = &adminID
}

//! IsActive vérifie si l'école est active
func (s *School) IsActive() bool {
	return s.Status == "active"
}
