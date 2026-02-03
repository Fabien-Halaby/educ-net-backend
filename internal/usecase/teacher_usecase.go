package usecase

import (
	"database/sql"
	"educnet/internal/domain"
	"educnet/internal/handler/dto"
	"educnet/internal/repository"
	"errors"
	"fmt"
)

type TeacherUseCase interface {
	RegisterTeacher(req *dto.TeacherRegistrationRequest) (*dto.TeacherRegistrationResponse, error)

	GetTeacherSubjects(teacherID int) ([]dto.SubjectResponse, error)
}

type teacherUseCase struct {
	db                 *sql.DB
	userRepo           repository.UserRepository
	schoolRepo         repository.SchoolRepository
	subjectRepo        repository.SubjectRepository
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

	//! 3. Validate subjects exist and belong to school
	if len(req.SubjectIDs) == 0 {
		return nil, domain.ErrValidation
	}

	subjectNames := []string{}
	for _, subjectID := range req.SubjectIDs {
		subject, err := uc.subjectRepo.FindByID(subjectID)
		if errors.Is(err, domain.ErrSubjectNotFound) {
			return nil, domain.ErrNotFound
		}
		if err != nil {
			return nil, domain.ErrInternal
		}
		if subject.SchoolID != school.ID {
			return nil, domain.ErrForbidden
		}
		subjectNames = append(subjectNames, subject.Name)
	}

	//! 4. Start transaction
	tx, err := uc.db.Begin()
	if err != nil {
		return nil, domain.ErrInternal
	}
	defer tx.Rollback()

	//! 5. Create teacher user (status = PENDING)
	user, err := domain.NewUser(
		school.ID, req.Email, req.Password, req.FirstName, req.LastName, req.Phone, domain.RoleTeacher,
	)
	if err != nil {
		return nil, err
	}
	user.Status = domain.UserStatusPending

	if err := uc.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create teacher: %w", err)
	}

	//! 6. Assign subjects (utilise repo injecté)
	for _, subjectID := range req.SubjectIDs {
		if err := uc.teacherSubjectRepo.Create(user.ID, subjectID); err != nil {
			return nil, fmt.Errorf("failed to assign subject %d: %w", subjectID, err)
		}
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
		Status:   user.Status,
		Subjects: subjectNames,
		Message:  "Teacher registration successful. Pending admin approval.",
	}, nil
}

func (uc *teacherUseCase) GetTeacherSubjects(teacherID int) ([]dto.SubjectResponse, error) {
	user, err := uc.userRepo.FindByID(teacherID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	if !user.IsTeacher() {
		return nil, domain.ErrForbidden
	}

	subjects, err := uc.teacherSubjectRepo.FindByTeacher(teacherID)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return dto.SubjectResponsesFromDomain(subjects), nil
}
