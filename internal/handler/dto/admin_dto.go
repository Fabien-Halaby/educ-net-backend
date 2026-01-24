package dto

type PendingUsersResponse struct {
	Users []PendingUserInfo `json:"users"`
	Total int               `json:"total"`
}

type PendingUserInfo struct {
	ID        int      `json:"id"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	Role      string   `json:"role"`
	Phone     string   `json:"phone,omitempty"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
	Subjects  []string `json:"subjects,omitempty"`  //! For teachers
	ClassName string   `json:"class_name,omitempty"` //! For students
}

type ApproveUserRequest struct {
	UserID int `json:"user_id"`
}

type RejectUserRequest struct {
	UserID int    `json:"user_id"`
	Reason string `json:"reason,omitempty"`
}

type UserListResponse struct {
	Users []UserListInfo `json:"users"`
	Total int            `json:"total"`
}

type UserListInfo struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	Phone     string `json:"phone,omitempty"`
	CreatedAt string `json:"created_at"`
}


//! ========== SUBJECTS ==========

type CreateSubjectRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
}

type UpdateSubjectRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
}

type SubjectResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	SchoolID    int    `json:"school_id"`
}

//! ========== CLASSES ==========

type CreateClassRequest struct {
	Name         string `json:"name"`
	Level        string `json:"level"`
	Section      string `json:"section,omitempty"`
	Capacity     int    `json:"capacity"`
	AcademicYear string `json:"academic_year"`
}

type UpdateClassRequest struct {
	Name         string `json:"name"`
	Level        string `json:"level"`
	Section      string `json:"section,omitempty"`
	Capacity     int    `json:"capacity"`
	AcademicYear string `json:"academic_year"`
}

type ClassResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Level        string `json:"level"`
	Section      string `json:"section,omitempty"`
	Capacity     int    `json:"capacity"`
	AcademicYear string `json:"academic_year"`
	SchoolID     int    `json:"school_id"`
}
