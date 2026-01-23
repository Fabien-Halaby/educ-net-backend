package utils

import (
	"testing"
)

func TestCreateSlug(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Simple name",
			input: "Lycée Andohalo",
			want:  "lycee-andohalo",
		},
		{
			name:  "Name with special characters",
			input: "École Primaire #1!",
			want:  "ecole-primaire-1",
		},
		{
			name:  "Name with multiple spaces",
			input: "  Collège   Saint   Paul  ",
			want:  "college-saint-paul",
		},
		{
			name:  "Name with underscores and hyphens",
			input: "Institut_Technique-Avancé",
			want:  "institut-technique-avance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateSlug(tt.input)
			if got != tt.want {
				t.Errorf("CreateSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSplitFullName(t *testing.T) {
	tests := []struct {
		name          string
		fullName      string
		wantFirstName string
		wantLastName  string
	}{
		{
			name:          "First and last name",
			fullName:      "John Doe",
			wantFirstName: "John",
			wantLastName:  "Doe",
		},
		{
			name:          "Single name",
			fullName:      "Madonna",
			wantFirstName: "Madonna",
			wantLastName:  "",
		},
		{
			name:          "Multiple spaces",
			fullName:      "  Alice   Smith  ",
			wantFirstName: "Alice",
			wantLastName:  "Smith",
		},
		{
			name:          "Middle name included",
			fullName:      "Bob A. Johnson",
			wantFirstName: "Bob",
			wantLastName:  "A. Johnson",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFirstName, gotLastName := SplitFullName(tt.fullName)
			if gotFirstName != tt.wantFirstName || gotLastName != tt.wantLastName {
				t.Errorf("SplitFullName(%q) = (%q, %q), want (%q, %q)",
					tt.fullName, gotFirstName, gotLastName, tt.wantFirstName, tt.wantLastName)
			}
		})
	}
}