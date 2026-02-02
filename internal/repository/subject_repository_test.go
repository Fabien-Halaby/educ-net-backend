package repository

import (
	"errors"
	"testing"

	"educnet/internal/domain"
	"educnet/internal/testutil"
)

func TestSubjectRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)

	// Seed school
	schoolID := testutil.SeedTestSchool(t, db, "Test School", "test", "test@school.mg")

	subject := &domain.Subject{
		SchoolID:    schoolID,
		Name:        "Mathématiques",
		Code:        "MATH",
		Description: "Cours de mathématiques",
	}

	err := repo.Create(subject)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if subject.ID == 0 {
		t.Error("Create() ID was not set")
	}
}

func TestSubjectRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)

	// Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Maths", "MATH", "Mathématiques")

	// Test
	subject, err := repo.FindByID(subjectID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if subject.ID != subjectID {
		t.Errorf("FindByID() ID = %v, want %v", subject.ID, subjectID)
	}
	if subject.Code != "MATH" {
		t.Errorf("FindByID() Code = %v, want MATH", subject.Code)
	}
}

func TestSubjectRepository_FindByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)

	_, err := repo.FindByID(99999)
	if !errors.Is(err, domain.ErrSubjectNotFound) {
		t.Errorf("FindByID() error = %v, want ErrSubjectNotFound", err)
	}
}

func TestSubjectRepository_FindBySchoolID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)

	// Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	testutil.SeedTestSubject(t, db, schoolID, "Maths", "MATH", "Mathématiques")
	testutil.SeedTestSubject(t, db, schoolID, "Français", "FR", "Français")

	// Test
	subjects, err := repo.FindBySchoolID(schoolID)
	if err != nil {
		t.Fatalf("FindBySchoolID() error = %v", err)
	}

	if len(subjects) != 2 {
		t.Errorf("FindBySchoolID() got %d subjects, want 2", len(subjects))
	}
}

func TestSubjectRepository_FindBySchoolAndCode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)

	// Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	testutil.SeedTestSubject(t, db, schoolID, "Maths", "MATH", "Mathématiques")

	// Test
	subject, err := repo.FindBySchoolAndCode(schoolID, "MATH")
	if err != nil {
		t.Fatalf("FindBySchoolAndCode() error = %v", err)
	}

	if subject.Code != "MATH" {
		t.Errorf("FindBySchoolAndCode() Code = %v, want MATH", subject.Code)
	}
}

func TestSubjectRepository_FindBySchoolAndCode_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")

	_, err := repo.FindBySchoolAndCode(schoolID, "PHYS")
	if !errors.Is(err, domain.ErrSubjectNotFound) {
		t.Errorf("FindBySchoolAndCode() error = %v, want ErrSubjectNotFound", err)
	}
}

func TestSubjectRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)

	// Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Maths", "MATH", "Mathématiques")

	// Get subject
	subject, err := repo.FindByID(subjectID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	// Update
	subject.Name = "Algèbre"
	subject.Code = "ALG"

	err = repo.Update(subject)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify
	updated, err := repo.FindByID(subjectID)
	if err != nil {
		t.Fatalf("FindByID() after update error = %v", err)
	}

	if updated.Name != "Algèbre" {
		t.Errorf("Update() Name = %v, want Algèbre", updated.Name)
	}
	if updated.Code != "ALG" {
		t.Errorf("Update() Code = %v, want ALG", updated.Code)
	}
}

func TestSubjectRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)

	// Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Maths", "MATH", "Mathématiques")

	err := repo.Delete(subjectID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.FindByID(subjectID)
	if !errors.Is(err, domain.ErrSubjectNotFound) {
		t.Errorf("Delete() subject still exists, want ErrSubjectNotFound")
	}
}

// REMPLACER test "Same code but excluded ID"
func TestSubjectRepository_ExistsByCode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSubjectRepository(db)

	// Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Maths", "MATH", "Mathématiques")

	tests := []struct {
		name      string
		schoolID  int
		code      string
		excludeID int
		exists    bool
	}{
		{
			name:      "Existing code",
			schoolID:  schoolID,
			code:      "MATH",
			excludeID: 0,
			exists:    true,
		},
		{
			name:      "Non-existing code",
			schoolID:  schoolID,
			code:      "PHYS",
			excludeID: 0,
			exists:    false,
		},
		{
			name:      "Same code but excluded ID",
			schoolID:  schoolID,
			code:      "MATH",
			excludeID: subjectID,
			exists:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistsByCode(tt.schoolID, tt.code, tt.excludeID)
			if err != nil {
				t.Fatalf("ExistsByCode() error = %v", err)
			}
			if exists != tt.exists {
				t.Errorf("ExistsByCode() = %v, want %v", exists, tt.exists)
			}
		})
	}
}
