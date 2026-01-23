package domain

import "testing"

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		schoolID    int
		email       string
		password    string
		firstName   string
		lastName    string
		phone       string
		role        string
		wantErr     bool
		expectedErr error
	}{
		{
			name:      "Valid user",
			schoolID:  1,
			email:     "test@test.mg",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			phone:     "+261 34 00 000 00",
			role:      "teacher",
			wantErr:   false,
		},
		{
			name:        "Empty email",
			schoolID:    1,
			email:       "",
			password:    "password123",
			firstName:   "John",
			role:        "teacher",
			wantErr:     true,
			expectedErr: ErrUserEmailRequired,
		},
		{
			name:        "Invalid email",
			schoolID:    1,
			email:       "not-an-email",
			password:    "password123",
			firstName:   "John",
			role:        "teacher",
			wantErr:     true,
			expectedErr: ErrUserEmailInvalid,
		},
		{
			name:        "Password too short",
			schoolID:    1,
			email:       "test@test.mg",
			password:    "123",
			firstName:   "John",
			role:        "teacher",
			wantErr:     true,
			expectedErr: ErrUserPasswordTooShort,
		},
		{
			name:        "Empty name",
			schoolID:    1,
			email:       "test@test.mg",
			password:    "password123",
			firstName:   "",
			role:        "teacher",
			wantErr:     true,
			expectedErr: ErrUserNameRequired,
		},
		{
			name:        "Invalid role",
			schoolID:    1,
			email:       "test@test.mg",
			password:    "password123",
			firstName:   "John",
			role:        "invalid",
			wantErr:     true,
			expectedErr: ErrUserInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(
				tt.schoolID,
				tt.email,
				tt.password,
				tt.firstName,
				tt.lastName,
				tt.phone,
				tt.role,
			)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error but got nil")
					return
				}
				if err != tt.expectedErr {
					t.Errorf("NewUser() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}

			if user.Email != tt.email {
				t.Errorf("NewUser() Email = %v, want %v", user.Email, tt.email)
			}
			if user.Role != tt.role {
				t.Errorf("NewUser() Role = %v, want %v", user.Role, tt.role)
			}
			if user.Status != "pending" {
				t.Errorf("NewUser() Status = %v, want pending", user.Status)
			}
		})
	}
}

func TestNewAdminUser(t *testing.T) {
	admin, err := NewAdminUser(1, "admin@test.mg", "password123", "Admin", "User", "+261 34")

	if err != nil {
		t.Fatalf("NewAdminUser() unexpected error = %v", err)
	}

	if admin.Role != "admin" {
		t.Errorf("NewAdminUser() Role = %v, want admin", admin.Role)
	}

	if admin.Status != "approved" {
		t.Errorf("NewAdminUser() Status = %v, want approved", admin.Status)
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	password := "testpassword123"
	user, _ := NewUser(1, "test@test.mg", password, "Test", "User", "", "teacher")

	//! Correct password
	if !user.VerifyPassword(password) {
		t.Error("VerifyPassword() = false, want true for correct password")
	}

	//! Wrong password
	if user.VerifyPassword("wrongpassword") {
		t.Error("VerifyPassword() = true, want false for wrong password")
	}
}

func TestUser_Approve(t *testing.T) {
	user, _ := NewUser(1, "test@test.mg", "password123", "Test", "User", "", "teacher")

	if user.Status != "pending" {
		t.Errorf("Initial Status = %v, want pending", user.Status)
	}

	user.Approve()

	if user.Status != "approved" {
		t.Errorf("After Approve() Status = %v, want approved", user.Status)
	}
}

func TestUser_IsAdmin(t *testing.T) {
	admin, _ := NewAdminUser(1, "admin@test.mg", "password123", "Admin", "User", "")
	teacher, _ := NewUser(1, "teacher@test.mg", "password123", "Teacher", "User", "", "teacher")

	if !admin.IsAdmin() {
		t.Error("IsAdmin() = false, want true for admin user")
	}

	if teacher.IsAdmin() {
		t.Error("IsAdmin() = true, want false for teacher user")
	}
}

func TestUser_GetFullName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		want      string
	}{
		{
			name:      "Full name",
			firstName: "John",
			lastName:  "Doe",
			want:      "John Doe",
		},
		{
			name:      "First name only",
			firstName: "John",
			lastName:  "",
			want:      "John",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, _ := NewUser(1, "test@test.mg", "password123", tt.firstName, tt.lastName, "", "teacher")

			got := user.GetFullName()
			if got != tt.want {
				t.Errorf("GetFullName() = %v, want %v", got, tt.want)
			}
		})
	}
}
