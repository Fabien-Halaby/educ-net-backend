package domain

import "time"

type Message struct {
	ID          int64     `json:"id"`
	Content     string    `json:"content"`
	UserID      int       `json:"user_id"`
	ClassID     int       `json:"class_id"`
	MessageType string    `json:"message_type"`
	FileURL     *string   `json:"file_url,omitempty"`
	IsPinned    bool      `json:"is_pinned"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        UserInfo  `json:"user"`
	ClassName   string    `json:"class_name"`
}

type UserInfo struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	FullName  string  `json:"full_name"`
	Role      string  `json:"role"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type CreateMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type MessageResponse struct {
	Message Message `json:"message"`
}
