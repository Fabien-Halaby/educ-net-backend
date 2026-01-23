package usecase

import (
	"database/sql"
	"educnet/internal/domain"
	"educnet/internal/handler/dto"
	"educnet/internal/repository"
	"fmt"
)

type TeacherUseCase interface {
	RegisterTeacher(req *dto.TeacherRegistrationRequest) (*dto.TeacherRegistrationResponse, error)
}

type teacherUseCase struct {
	db                *sql.DB
	userRepo          repository.UserRepository
	schoolRepo        repository.SchoolRepository
	subjectRepo       repository.SubjectRepository
	teacherSubjectRepo repository.TeacherSubjectRepository
}

func NewTeacherUseCase(
	db *sql.DB,
	userRepo repository.UserRepository,
	schoolRepo repository.SchoolRepository,
	subjectRepo repository.SubjectRepository,
	teacherSubjectRepo repository.TeacherSubjectRepository,
) TeacherUseCase {
	return &teacherUseCase{
		db:                 db,
		userRepo:           userRepo,
		schoolRepo:         schoolRepo,
		subjectRepo:        subjectRepo,
		teacherSubjectRepo: teacherSubjectRepo,
	}
}

func (uc *teacherUseCase) RegisterTeacher(req *dto.TeacherRegistrationRequest) (*dto.TeacherRegistrationResponse, error) {
	//! 1. Validate school exists
	school, err := uc.schoolRepo.FindBySlug(req.SchoolSlug)
	if err != nil {
		return nil, fmt.Errorf("school not found: %w", err)
	}

	//! 2. Check if email already exists
	exists, err := uc.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	//! 3. Validate subjects exist and belong to school
	if len(req.SubjectIDs) == 0 {
		return nil, fmt.Errorf("at least one subject is required")
	}

	subjectNames := []string{}
	for _, subjectID := range req.SubjectIDs {
		subject, err := uc.subjectRepo.FindByID(subjectID)
		if err != nil {
			return nil, fmt.Errorf("subject %d not found: %w", subjectID, err)
		}
		if subject.SchoolID != school.ID {
			return nil, fmt.Errorf("subject %d does not belong to this school", subjectID)
		}
		subjectNames = append(subjectNames, subject.Name)
	}

	//! 4. Start transaction
	tx, err := uc.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	//! 5. Create teacher user (status = PENDING)
	// schoolID int, email, password, firstName, lastName, phone, role string
	user, err := domain.NewUser(
		school.ID,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		req.Phone,
		domain.RoleTeacher,
	)

	if err != nil {
		return nil, err
	}
	user.Phone = req.Phone
	user.Status = domain.UserStatusPending

	if err := uc.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create teacher: %w", err)
	}

	//! 6. Assign subjects
	if err := uc.teacherSubjectRepo.AssignSubjects(user.ID, req.SubjectIDs); err != nil {
		return nil, fmt.Errorf("failed to assign subjects: %w", err)
	}

	//! 7. Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	//! 8. Return response
	return &dto.TeacherRegistrationResponse{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: user.GetFullName(),
		SchoolID: school.ID,
		Status:   string(user.Status),
		Subjects: subjectNames,
		Message:  "Teacher registration successful. Your account is pending approval by the school admin.",
	}, nil
}
