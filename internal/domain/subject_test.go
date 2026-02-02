package domain

import "testing"

func TestNewSubject(t *testing.T) {
	tests := []struct {
		name        string
		schoolID    int
		subjectName string
		code        string
		description string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Valid subject",
			schoolID:    1,
			subjectName: "Mathématiques",
			code:        "MATH",
			description: "Mathématiques générales",
			wantErr:     false,
		},
		{
			name:        "Invalid school ID",
			schoolID:    0,
			subjectName: "Mathématiques",
			code:        "MATH",
			wantErr:     true,
			expectedErr: ErrSubjectInvalidID,
		},
		{
			name:        "Empty name",
			schoolID:    1,
			subjectName: "",
			code:        "MATH",
			wantErr:     true,
			expectedErr: ErrSubjectNameRequired,
		},
		{
			name:        "Empty code",
			schoolID:    1,
			subjectName: "Mathématiques",
			code:        "",
			wantErr:     true,
			expectedErr: ErrSubjectCodeRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := NewSubject(tt.schoolID, tt.subjectName, tt.code, tt.description)

			if tt.wantErr {
				if err == nil || err != tt.expectedErr {
					t.Errorf("expected %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sub.SchoolID != tt.schoolID {
				t.Errorf("SchoolID = %d, want %d", sub.SchoolID, tt.schoolID)
			}
			if sub.Status != SubjectStatusActive {
				t.Errorf("Status = %s, want %s", sub.Status, SubjectStatusActive)
			}
		})
	}
}
