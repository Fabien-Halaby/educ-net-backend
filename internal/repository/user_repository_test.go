package repository

import (
	"testing"

	"educnet/internal/domain"
	"educnet/internal/testutil"
)

func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewUserRepository(db)

	//! Seed school
	schoolID := testutil.SeedTestSchool(t, db, "Test School", "test", "test@school.mg")

	user := &domain.User{
		SchoolID:     schoolID,
		Email:        "teacher@test.mg",
		PasswordHash: "hashed_password",
		FirstName:    "John",
		LastName:     "Doe",
		Phone:        "+261 34 00 000 00",
		Role:         "teacher",
		Status:       "pending",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if user.ID == 0 {
		t.Error("Create() ID was not set")
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewUserRepository(db)

	//! Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	testutil.SeedTestUser(t, db, schoolID, "teacher@test.mg", "teacher")

	//! Test
	user, err := repo.FindByEmail("teacher@test.mg")
	if err != nil {
		t.Fatalf("FindByEmail() error = %v", err)
	}

	if user.Email != "teacher@test.mg" {
		t.Errorf("FindByEmail() Email = %v, want teacher@test.mg", user.Email)
	}
	if user.Role != "teacher" {
		t.Errorf("FindByEmail() Role = %v, want teacher", user.Role)
	}
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewUserRepository(db)

	//! Seed
	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	testutil.SeedTestUser(t, db, schoolID, "existing@test.mg", "teacher")

	tests := []struct {
		name   string
		email  string
		exists bool
	}{
		{
			name:   "Existing email",
			email:  "existing@test.mg",
			exists: true,
		},
		{
			name:   "Non-existing email",
			email:  "notfound@test.mg",
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistsByEmail(tt.email)
			if err != nil {
				t.Fatalf("ExistsByEmail() error = %v", err)
			}
			if exists != tt.exists {
				t.Errorf("ExistsByEmail() = %v, want %v", exists, tt.exists)
			}
		})
	}
}
