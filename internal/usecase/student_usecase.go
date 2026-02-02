package usecase

import (
	"database/sql"
	"educnet/internal/domain"
	"educnet/internal/handler/dto"
	"educnet/internal/repository"
	"errors"
	"fmt"
)

type StudentUseCase interface {
	RegisterStudent(req *dto.StudentRegistrationRequest) (*dto.StudentRegistrationResponse, error)
}

type studentUseCase struct {
	db               *sql.DB
	userRepo         repository.UserRepository
	schoolRepo       repository.SchoolRepository
	classRepo        repository.ClassRepository
	studentClassRepo repository.StudentClassRepository
}

func NewStudentUseCase(
	db *sql.DB,
	userRepo repository.UserRepository,
	schoolRepo repository.SchoolRepository,
	classRepo repository.ClassRepository,
	studentClassRepo repository.StudentClassRepository,
) StudentUseCase {
	return &studentUseCase{
		db:               db,
		userRepo:         userRepo,
		schoolRepo:       schoolRepo,
		classRepo:        classRepo,
		studentClassRepo: studentClassRepo,
	}
}

func (uc *studentUseCase) RegisterStudent(req *dto.StudentRegistrationRequest) (*dto.StudentRegistrationResponse, error) {
	//! 1. Validate school exists
	school, err := uc.schoolRepo.FindBySlug(req.SchoolSlug)
	if errors.Is(err, domain.ErrSchoolNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, domain.ErrInternal
	}

	//! 2. Check if email already exists
	exists, err := uc.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, domain.ErrInternal
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	//! 3. Validate class exists and belongs to school
	class, err := uc.classRepo.FindByID(req.ClassID)
	if errors.Is(err, domain.ErrClassNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, domain.ErrInternal
	}
	if class.SchoolID != school.ID {
		return nil, domain.ErrForbidden
	}

	//! 4. Start transaction
	tx, err := uc.db.Begin()
	if err != nil {
		return nil, domain.ErrInternal
	}
	defer tx.Rollback()

	//! 5. Create student user (status = PENDING)
	user, err := domain.NewUser(
		school.ID, req.Email, req.Password, req.FirstName, req.LastName, req.Phone, domain.RoleStudent,
	)
	if err != nil {
		return nil, err
	}
	user.Status = domain.UserStatusPending

	if err := uc.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create student: %w", err)
	}

	//! 6. Enroll in class
	if err := uc.studentClassRepo.Create(user.ID, class.ID); err != nil {
		return nil, fmt.Errorf("failed to enroll student: %w", err)
	}

	//! 7. Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	//! 8. Return response
	return &dto.StudentRegistrationResponse{
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.GetFullName(),
		SchoolID:  school.ID,
		Status:    user.Status,
		ClassName: class.Name,
		Message:   "Student registration successful. Pending admin approval.",
	}, nil
}
