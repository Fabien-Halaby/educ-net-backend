package domain

import "testing"

func TestNewClass(t *testing.T) {
	tests := []struct {
		name        string
		schoolID    int
		className   string
		level       string
		section     string
		year        string
		wantErr     bool
		expectedErr error
	}{
		{
			name:      "Valid class",
			schoolID:  1,
			className: "6ème A",
			level:     "6ème",
			section:   "A",
			year:      "2025-2026",
			wantErr:   false,
		},
		{
			name:        "Invalid school ID",
			schoolID:    0,
			className:   "6ème A",
			level:       "6ème",
			year:        "2025-2026",
			wantErr:     true,
			expectedErr: ErrClassInvalidID,
		},
		{
			name:        "Empty class name",
			schoolID:    1,
			className:   "",
			level:       "6ème",
			year:        "2025-2026",
			wantErr:     true,
			expectedErr: ErrClassNameRequired,
		},
		{
			name:        "Empty level",
			schoolID:    1,
			className:   "6ème A",
			level:       "",
			section:     "A",
			year:        "2025-2026",
			wantErr:     true,
			expectedErr: ErrClassLevelRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cls, err := NewClass(
				tt.schoolID,
				tt.className,
				tt.level,
				tt.section,
				tt.year,
			)

			if tt.wantErr {
				if err == nil || err != tt.expectedErr {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cls.SchoolID != tt.schoolID {
				t.Errorf("SchoolID = %d, want %d", cls.SchoolID, tt.schoolID)
			}

			if cls.Status != ClassStatusActive {
				t.Errorf("Status = %s, want %s", cls.Status, ClassStatusActive)
			}
		})
	}
}

func TestClass_IsActive(t *testing.T) {
	cls, _ := NewClass(1, "6ème A", "6ème", "A", "2025-2026")

	if !cls.IsActive() {
		t.Error("New class should be active")
	}

	cls.Status = ClassStatusArchived
	if cls.IsActive() {
		t.Error("Archived class should not be active")
	}
}
