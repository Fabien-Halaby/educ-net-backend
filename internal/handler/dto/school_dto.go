package dto

import "time"

//! CreateSchoolRequest DTO de requête HTTP
type CreateSchoolRequest struct {
	SchoolName    string `json:"school_name"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
	AdminName     string `json:"admin_name"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
}

//! CreateSchoolResponse DTO de réponse HTTP
type CreateSchoolResponse struct {
	School SchoolDTO `json:"school"`
	Admin  UserDTO   `json:"admin"`
	Token  string    `json:"token"`
}

//! SchoolDTO représentation HTTP d'une école
type SchoolDTO struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Address     string    `json:"address,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Status      string    `json:"status"`
	AdminUserID *int      `json:"admin_user_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

//! UserDTO représentation HTTP d'un utilisateur
type UserDTO struct {
	ID        int       `json:"id"`
	SchoolID  int       `json:"school_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone,omitempty"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
