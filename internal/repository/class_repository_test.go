package repository

import (
	"testing"

	"educnet/internal/domain"
	"educnet/internal/testutil"
)

func TestClassRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewClassRepository(db)

	//! Seed school
	schoolID := testutil.SeedTestSchool(t, db, "Test School", "test", "test@school.mg")

	class := &domain.Class{
		SchoolID:     schoolID,
		Name:         "6ème A",
		Level:        "6ème",
		Section:      "A",
		Capacity:     40,
		AcademicYear: "2025-2026",
		Status:       domain.ClassStatusActive,
	}

	err := repo.Create(class)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if class.ID == 0 {
		t.Error("Create() ID was not set")
	}
	if class.CreatedAt.IsZero() {
		t.Error("Create() CreatedAt was not set")
	}
}

func TestClassRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewClassRepository(db)

	//! Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	classID := testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")

	//! Test
	class, err := repo.FindByID(classID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if class.ID != classID {
		t.Errorf("FindByID() ID = %v, want %v", class.ID, classID)
	}
	if class.Name != "6ème A" {
		t.Errorf("FindByID() Name = %v, want 6ème A", class.Name)
	}
}

func TestClassRepository_FindByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewClassRepository(db)

	_, err := repo.FindByID(99999)
	if err != domain.ErrClassNotFound {
		t.Errorf("FindByID() error = %v, want ErrClassNotFound", err)
	}
}

func TestClassRepository_FindBySchoolID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewClassRepository(db)

	//! Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")
	testutil.SeedTestClass(t, db, schoolID, "5ème B", "5ème", "B", "2025-2026")

	//! Test
	classes, err := repo.FindBySchoolID(schoolID)
	if err != nil {
		t.Fatalf("FindBySchoolID() error = %v", err)
	}

	if len(classes) != 2 {
		t.Errorf("FindBySchoolID() got %d classes, want 2", len(classes))
	}
}

func TestClassRepository_FindBySchoolAndYear(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewClassRepository(db)

	//! Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")
	testutil.SeedTestClass(t, db, schoolID, "6ème B", "6ème", "B", "2025-2026")
	testutil.SeedTestClass(t, db, schoolID, "5ème A", "5ème", "A", "2024-2025") // autre année

	//! Test
	classes, err := repo.FindBySchoolAndYear(schoolID, "2025-2026")
	if err != nil {
		t.Fatalf("FindBySchoolAndYear() error = %v", err)
	}

	if len(classes) != 2 {
		t.Errorf("FindBySchoolAndYear() got %d classes, want 2", len(classes))
	}
}

func TestClassRepository_ExistsByName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewClassRepository(db)

	//! Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	classID := testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")

	tests := []struct {
		name      string
		schoolID  int
		className string
		excludeID int
		exists    bool
	}{
		{
			name:      "Existing name",
			schoolID:  schoolID,
			className: "6ème A",
			excludeID: 0,
			exists:    true,
		},
		{
			name:      "Non-existing name",
			schoolID:  schoolID,
			className: "4ème Z",
			excludeID: 0,
			exists:    false,
		},
		{
			name:      "Same name but excluded ID",
			schoolID:  schoolID,
			className: "6ème A",
			excludeID: classID,
			exists:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistsByName(tt.schoolID, tt.className, tt.excludeID)
			if err != nil {
				t.Fatalf("ExistsByName() error = %v", err)
			}
			if exists != tt.exists {
				t.Errorf("ExistsByName() = %v, want %v", exists, tt.exists)
			}
		})
	}
}

func TestClassRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewClassRepository(db)

	//! Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	classID := testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")

	//! Get class
	class, err := repo.FindByID(classID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	//! Update
	class.Name = "6ème B"
	class.Capacity = 45
	class.Activate()

	err = repo.Update(class)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	//! Verify
	updated, err := repo.FindByID(classID)
	if err != nil {
		t.Fatalf("FindByID() after update error = %v", err)
	}

	if updated.Name != "6ème B" {
		t.Errorf("Update() Name = %v, want 6ème B", updated.Name)
	}
	if updated.Capacity != 45 {
		t.Errorf("Update() Capacity = %v, want 45", updated.Capacity)
	}
}
