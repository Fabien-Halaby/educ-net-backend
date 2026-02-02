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
			expectedErr: ErrEmailRequired,
		},
		{
			name:        "Invalid email",
			schoolID:    1,
			email:       "not-an-email",
			password:    "password123",
			firstName:   "John",
			role:        "teacher",
			wantErr:     true,
			expectedErr: ErrEmailInvalid,
		},
		{
			name:        "Password too short",
			schoolID:    1,
			email:       "test@test.mg",
			password:    "123",
			firstName:   "John",
			role:        "teacher",
			wantErr:     true,
			expectedErr: ErrPasswordTooShort,
		},
		{
			name:        "Empty name",
			schoolID:    1,
			email:       "test@test.mg",
			password:    "password123",
			firstName:   "",
			role:        "teacher",
			wantErr:     true,
			expectedErr: ErrNameRequired,
		},
		{
			name:        "Invalid role",
			schoolID:    1,
			email:       "test@test.mg",
			password:    "password123",
			firstName:   "John",
			role:        "invalid",
			wantErr:     true,
			expectedErr: ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// schoolID int, email, password, firstName, lastName, phone, role string
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

func TestUserRoles(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"Valid admin", RoleAdmin, true},
		{"Valid teacher", RoleTeacher, true},
		{"Valid student", RoleStudent, true},
		{"Valid parent", RoleParent, true},
		{"Invalid role", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// schoolID int, email, password, firstName, lastName, phone, role string
			user, err := NewUser(
				1,
				"test@example.com",
				"password123",
				"John",
				"Doe",
				"+261 34",
				tt.role,
			)
			if tt.want && err != nil {
				t.Errorf("Expected valid role, got error: %v", err)
			}
			if !tt.want && err == nil {
				t.Errorf("Expected error for invalid role, got nil")
			}
			if tt.want && user != nil && user.Role != tt.role {
				t.Errorf("Expected role %s, got %s", tt.role, user.Role)
			}
		})
	}
}

func TestUserStatus(t *testing.T) {
	// schoolID int, email, password, firstName, lastName, phone, role string
	user, _ := NewUser(1, "test@example.com", "password123", "John", "Doe", "+261 34", RoleTeacher)

	//! Default status should be pending
	if user.Status != UserStatusPending {
		t.Errorf("Expected default status to be pending, got %s", user.Status)
	}

	if !user.IsPending() {
		t.Error("Expected IsPending() to return true")
	}

	//! Test approval
	user.Approve()
	if user.Status != UserStatusApproved {
		t.Errorf("Expected status to be approved after approval, got %s", user.Status)
	}

	if !user.IsApproved() {
		t.Error("Expected IsApproved() to return true after approval")
	}

	//! Test suspension
	user.Suspend()
	if user.Status != UserStatusSuspended {
		t.Errorf("Expected status to be suspended, got %s", user.Status)
	}
}

func TestUserRoleChecks(t *testing.T) {
	// schoolID int, email, password, firstName, lastName, phone string
	admin, _ := NewAdminUser(1, "admin@test.com", "password123", "Admin", "User", "+261 34")
	teacher, _ := NewUser(1, "teacher@test.com", "password123", "Teacher", "User", "+261 34", RoleTeacher)
	student, _ := NewUser(1, "student@test.com", "password123", "Student", "User", "+261 34", RoleStudent)

	if !admin.IsAdmin() {
		t.Error("Expected admin.IsAdmin() to return true")
	}

	if !teacher.IsTeacher() {
		t.Error("Expected teacher.IsTeacher() to return true")
	}

	if !student.IsStudent() {
		t.Error("Expected student.IsStudent() to return true")
	}

	//! Cross-checks
	if teacher.IsAdmin() {
		t.Error("Expected teacher.IsAdmin() to return false")
	}

	if student.IsTeacher() {
		t.Error("Expected student.IsTeacher() to return false")
	}
}
