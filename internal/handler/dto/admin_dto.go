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
