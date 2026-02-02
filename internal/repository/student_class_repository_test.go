package repository

import (
	"educnet/internal/domain"
	"educnet/internal/testutil"
	"testing"
)

func TestStudentClassRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewStudentClassRepository(db)

	//! Seed prerequisites
	schoolID := testutil.SeedTestSchool(t, db, "Test School", "test", "test@school.mg")
	studentID := testutil.SeedTestUser(t, db, schoolID, "student@test.mg", domain.RoleStudent)
	classID := testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")

	//! Test Create
	err := repo.Create(studentID, classID)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	//! Verify exists
	exists, err := repo.Exists(studentID, classID)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Create() should create association")
	}
}

func TestStudentClassRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewStudentClassRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	studentID := testutil.SeedTestUser(t, db, schoolID, "student@test.mg", domain.RoleStudent)
	classID := testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")

	//! Create first
	repo.Create(studentID, classID)

	//! Test Delete
	err := repo.Delete(studentID, classID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	//! Verify deleted
	exists, _ := repo.Exists(studentID, classID)
	if exists {
		t.Error("Delete() should remove association")
	}
}

func TestStudentClassRepository_FindByStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewStudentClassRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	studentID := testutil.SeedTestUser(t, db, schoolID, "student@test.mg", domain.RoleStudent)
	class1ID := testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")
	class2ID := testutil.SeedTestClass(t, db, schoolID, "5ème B", "5ème", "B", "2025-2026")

	//! Enroll student in 2 classes
	repo.Create(studentID, class1ID)
	repo.Create(studentID, class2ID)

	//! Test FindByStudent
	classes, err := repo.FindByStudent(studentID)
	if err != nil {
		t.Fatalf("FindByStudent() error = %v", err)
	}
	if len(classes) != 2 {
		t.Errorf("FindByStudent() got %d classes, want 2", len(classes))
	}
}

func TestStudentClassRepository_FindByClass(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewStudentClassRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	student1ID := testutil.SeedTestUser(t, db, schoolID, "student1@test.mg", domain.RoleStudent)
	student2ID := testutil.SeedTestUser(t, db, schoolID, "student2@test.mg", domain.RoleStudent)
	classID := testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")

	//! Enroll 2 students in class
	repo.Create(student1ID, classID)
	repo.Create(student2ID, classID)

	//! Test FindByClass
	students, err := repo.FindByClass(classID)
	if err != nil {
		t.Fatalf("FindByClass() error = %v", err)
	}
	if len(students) != 2 {
		t.Errorf("FindByClass() got %d students, want 2", len(students))
	}
}

func TestStudentClassRepository_DeleteByStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewStudentClassRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	studentID := testutil.SeedTestUser(t, db, schoolID, "student@test.mg", domain.RoleStudent)
	classID := testutil.SeedTestClass(t, db, schoolID, "6ème A", "6ème", "A", "2025-2026")

	repo.Create(studentID, classID)
	err := repo.DeleteByStudent(studentID)
	if err != nil {
		t.Fatalf("DeleteByStudent() error = %v", err)
	}

	classes, _ := repo.FindByStudent(studentID)
	if len(classes) != 0 {
		t.Error("DeleteByStudent() should remove all student classes")
	}
}
