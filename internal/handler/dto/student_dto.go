package dto

type StudentRegistrationRequest struct {
	SchoolSlug string `json:"school_slug"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	ClassID    int    `json:"class_id"`
}

type StudentRegistrationResponse struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	SchoolID  int    `json:"school_id"`
	Status    string `json:"status"`
	ClassName string `json:"class_name"`
	Message   string `json:"message"`
}
