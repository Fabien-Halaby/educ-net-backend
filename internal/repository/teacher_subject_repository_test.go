package repository

import (
	"educnet/internal/domain"
	"educnet/internal/testutil"
	"testing"
)

func TestTeacherSubjectRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewTeacherSubjectRepository(db)

	//! Seed prerequisites
	schoolID := testutil.SeedTestSchool(t, db, "Test School", "test", "test@school.mg")
	teacherID := testutil.SeedTestUser(t, db, schoolID, "teacher@test.mg", domain.RoleTeacher)
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Mathématiques", "MATH", "Test maths subject")

	//! Test Create
	err := repo.Create(teacherID, subjectID)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	//! Verify exists
	exists, err := repo.Exists(teacherID, subjectID)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Create() should create association")
	}
}

func TestTeacherSubjectRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewTeacherSubjectRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	teacherID := testutil.SeedTestUser(t, db, schoolID, "teacher@test.mg", domain.RoleTeacher)
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Math", "MATH", "Test maths subject")

	//! Create first
	repo.Create(teacherID, subjectID)

	//! Test Delete
	err := repo.Delete(teacherID, subjectID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	//! Verify deleted
	exists, _ := repo.Exists(teacherID, subjectID)
	if exists {
		t.Error("Delete() should remove association")
	}
}

func TestTeacherSubjectRepository_FindByTeacher(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewTeacherSubjectRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	teacherID := testutil.SeedTestUser(t, db, schoolID, "teacher@test.mg", domain.RoleTeacher)
	subject1ID := testutil.SeedTestSubject(t, db, schoolID, "Mathématiques", "MATH", "Test maths subject")
	subject2ID := testutil.SeedTestSubject(t, db, schoolID, "Français", "FR", "Test frs subject")

	//! Assign 2 subjects to teacher
	repo.Create(teacherID, subject1ID)
	repo.Create(teacherID, subject2ID)

	//! Test FindByTeacher
	subjects, err := repo.FindByTeacher(teacherID)
	if err != nil {
		t.Fatalf("FindByTeacher() error = %v", err)
	}
	if len(subjects) != 2 {
		t.Errorf("FindByTeacher() got %d subjects, want 2", len(subjects))
	}
}

func TestTeacherSubjectRepository_FindBySubject(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewTeacherSubjectRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	teacher1ID := testutil.SeedTestUser(t, db, schoolID, "teacher1@test.mg", domain.RoleTeacher)
	teacher2ID := testutil.SeedTestUser(t, db, schoolID, "teacher2@test.mg", domain.RoleTeacher)
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Mathématiques", "MATH", "Test maths subject")

	// Assign 2 teachers to subject
	repo.Create(teacher1ID, subjectID)
	repo.Create(teacher2ID, subjectID)

	// Test FindBySubject
	teachers, err := repo.FindBySubject(subjectID)
	if err != nil {
		t.Fatalf("FindBySubject() error = %v", err)
	}
	if len(teachers) != 2 {
		t.Errorf("FindBySubject() got %d teachers, want 2", len(teachers))
	}
}

func TestTeacherSubjectRepository_DeleteByTeacher(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewTeacherSubjectRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	teacherID := testutil.SeedTestUser(t, db, schoolID, "teacher@test.mg", domain.RoleTeacher)
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Math", "MATH", "Test maths subject")

	repo.Create(teacherID, subjectID)
	err := repo.DeleteByTeacher(teacherID)
	if err != nil {
		t.Fatalf("DeleteByTeacher() error = %v", err)
	}

	subjects, _ := repo.FindByTeacher(teacherID)
	if len(subjects) != 0 {
		t.Error("DeleteByTeacher() should remove all teacher subjects")
	}
}

func TestTeacherSubjectRepository_Exists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	repo := NewTeacherSubjectRepository(db)

	schoolID := testutil.SeedTestSchool(t, db, "Test", "test", "test@school.mg")
	teacherID := testutil.SeedTestUser(t, db, schoolID, "teacher@test.mg", domain.RoleTeacher)
	subjectID := testutil.SeedTestSubject(t, db, schoolID, "Math", "MATH", "Test maths subject")

	repo.Create(teacherID, subjectID)

	exists, err := repo.Exists(teacherID, subjectID)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() = false, want true")
	}
}
