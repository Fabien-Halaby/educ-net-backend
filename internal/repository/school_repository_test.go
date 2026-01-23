package repository

import (
	"testing"

	"educnet/internal/domain"
	"educnet/internal/testutil"
)

func TestSchoolRepository_Create(t *testing.T) {
	//! Skip si pas de DB_TEST
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSchoolRepository(db)

	school := &domain.School{
		Name:    "Test School",
		Slug:    "test-school",
		Email:   "test@school.mg",
		Address: "Test Address",
		Phone:   "+261 34 12 345 67",
		Status:  "active",
	}

	err := repo.Create(school)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	//! Vérifier que l'ID a été assigné
	if school.ID == 0 {
		t.Error("Create() ID was not set")
	}

	//! Vérifier que CreatedAt a été set
	if school.CreatedAt.IsZero() {
		t.Error("Create() CreatedAt was not set")
	}
}

func TestSchoolRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSchoolRepository(db)

	//! Seed data
	schoolID := testutil.SeedTestSchool(t, db, "Test School", "test-school", "test@school.mg")

	//! Test
	school, err := repo.FindByID(schoolID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if school.ID != schoolID {
		t.Errorf("FindByID() ID = %v, want %v", school.ID, schoolID)
	}
	if school.Name != "Test School" {
		t.Errorf("FindByID() Name = %v, want Test School", school.Name)
	}
}

func TestSchoolRepository_FindByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSchoolRepository(db)

	_, err := repo.FindByID(99999)
	if err != domain.ErrSchoolNotFound {
		t.Errorf("FindByID() error = %v, want ErrSchoolNotFound", err)
	}
}

func TestSchoolRepository_FindBySlug(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSchoolRepository(db)

	//! Seed
	testutil.SeedTestSchool(t, db, "Test School", "test-school", "test@school.mg")

	//! Test
	school, err := repo.FindBySlug("test-school")
	if err != nil {
		t.Fatalf("FindBySlug() error = %v", err)
	}

	if school.Slug != "test-school" {
		t.Errorf("FindBySlug() Slug = %v, want test-school", school.Slug)
	}
}

func TestSchoolRepository_ExistsBySlug(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSchoolRepository(db)

	//! Seed
	testutil.SeedTestSchool(t, db, "Test School", "test-school", "test@school.mg")

	tests := []struct {
		name   string
		slug   string
		exists bool
	}{
		{
			name:   "Existing slug",
			slug:   "test-school",
			exists: true,
		},
		{
			name:   "Non-existing slug",
			slug:   "non-existing",
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistsBySlug(tt.slug)
			if err != nil {
				t.Fatalf("ExistsBySlug() error = %v", err)
			}
			if exists != tt.exists {
				t.Errorf("ExistsBySlug() = %v, want %v", exists, tt.exists)
			}
		})
	}
}

func TestSchoolRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewSchoolRepository(db)

	//! Seed
	schoolID := testutil.SeedTestSchool(t, db, "Original Name", "original", "original@school.mg")
	adminID := testutil.SeedTestUser(t, db, schoolID, "admin@school.mg", "admin")

	//! Get school
	school, err := repo.FindByID(schoolID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	//! Update
	school.Name = "Updated Name"
	school.Address = "Updated Address"
	school.SetAdmin(adminID)

	err = repo.Update(school)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	//! Verify
	updated, err := repo.FindByID(schoolID)
	if err != nil {
		t.Fatalf("FindByID() after update error = %v", err)
	}
	
	if updated.Name != "Updated Name" {
		t.Errorf("Update() Name = %v, want Updated Name", updated.Name)
	}
	if updated.AdminUserID == nil || *updated.AdminUserID != adminID {
		t.Errorf("Update() AdminUserID = %v, want %v", updated.AdminUserID, adminID)
	}
}

