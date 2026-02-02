package domain

import "testing"

func TestNewSchool(t *testing.T) {
	tests := []struct {
		name        string
		schoolName  string
		slug        string
		address     string
		phone       string
		email       string
		wantErr     bool
		expectedErr error
	}{
		{
			name:       "Valid school",
			schoolName: "Lycée Andohalo",
			slug:       "lycee-andohalo",
			email:      "admin@andohalo.mg",
			address:    "Andohalo, Tana",
			phone:      "+261 34 12 345 67",
			wantErr:    false,
		},
		{
			name:        "Empty name",
			schoolName:  "",
			slug:        "test",
			email:       "test@test.mg",
			wantErr:     true,
			expectedErr: ErrSchoolNameRequired,
		},
		{
			name:        "Empty slug",
			schoolName:  "Test School",
			slug:        "",
			email:       "test@test.mg",
			wantErr:     true,
			expectedErr: ErrSchoolSlugRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			school, err := NewSchool(tt.schoolName, tt.slug, tt.address, tt.email, tt.phone)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewSchool() expected error but got nil")
					return
				}
				if err != tt.expectedErr {
					t.Errorf("NewSchool() error = %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewSchool() unexpected error = %v", err)
				return
			}

			if school.Name != tt.schoolName {
				t.Errorf("NewSchool() Name = %v, want %v", school.Name, tt.schoolName)
			}
			if school.Slug != tt.slug {
				t.Errorf("NewSchool() Slug = %v, want %v", school.Slug, tt.slug)
			}
			if school.Status != "active" {
				t.Errorf("NewSchool() Status = %v, want active", school.Status)
			}
		})
	}
}

func TestSchool_SetAdmin(t *testing.T) {
	school, _ := NewSchool("Test School", "test-school", "Address", "test@test.mg", "123")

	adminID := 42
	school.SetAdmin(adminID)

	if school.AdminUserID == nil {
		t.Error("SetAdmin() AdminUserID is nil")
		return
	}

	if *school.AdminUserID != adminID {
		t.Errorf("SetAdmin() AdminUserID = %v, want %v", *school.AdminUserID, adminID)
	}
}

func TestSchool_IsActive(t *testing.T) {
	school, _ := NewSchool("Test", "test", "", "test@test.mg", "")

	if !school.IsActive() {
		t.Error("IsActive() = false, want true for new school")
	}

	school.Status = SchoolStatusInactive
	if school.IsActive() {
		t.Error("IsActive() = true, want false for inactive school")
	}
}
