package dto

type TeacherRegistrationRequest struct {
	SchoolSlug  string   `json:"school_slug"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Phone       string   `json:"phone"`
	SubjectIDs  []int    `json:"subject_ids"`
}

type TeacherRegistrationResponse struct {
	UserID    int      `json:"user_id"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	SchoolID  int      `json:"school_id"`
	Status    string   `json:"status"`
	Subjects  []string `json:"subjects"`
	Message   string   `json:"message"`
}
