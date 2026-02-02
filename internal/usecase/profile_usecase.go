package usecase

import (
	"educnet/internal/domain"
	"educnet/internal/handler/dto"
	"educnet/internal/repository"
	"errors"
	"fmt"
)

type ProfileUseCase interface {
	GetProfile(userID int) (*dto.ProfileResponse, error)
	UpdateProfile(userID int, req *dto.UpdateProfileRequest) (*dto.ProfileResponse, error)
	ChangePassword(userID int, req *dto.ChangePasswordRequest) error

	UpdateAvatar(userID int, avatarURL string) error
	GetSchool(schoolID int) (*domain.School, error)
	UpdateSchool(schoolID int, req *dto.UpdateSchoolRequest) (*domain.School, error)
	UpdateSchoolLogo(schoolID int, logoURL string) error

	GetTeacherSubjects(userID int) (*dto.TeacherSubjectsResponse, error)
	GetStudentClass(userID int) (*dto.StudentClassResponse, error)
}

type profileUseCase struct {
	userRepo           repository.UserRepository
	subjectRepo        repository.SubjectRepository
	classRepo          repository.ClassRepository
	teacherSubjectRepo repository.TeacherSubjectRepository
	studentClassRepo   repository.StudentClassRepository
	schoolRepo         repository.SchoolRepository
}

func NewProfileUseCase(
	userRepo repository.UserRepository,
	subjectRepo repository.SubjectRepository,
	classRepo repository.ClassRepository,
	teacherSubjectRepo repository.TeacherSubjectRepository,
	studentClassRepo repository.StudentClassRepository,
	schoolRepo repository.SchoolRepository,
) ProfileUseCase {
	return &profileUseCase{
		userRepo:           userRepo,
		subjectRepo:        subjectRepo,
		classRepo:          classRepo,
		teacherSubjectRepo: teacherSubjectRepo,
		studentClassRepo:   studentClassRepo,
		schoolRepo:         schoolRepo,
	}
}

func (uc *profileUseCase) GetProfile(userID int) (*dto.ProfileResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.GetFullName(),
		Phone:     user.Phone,
		Role:      user.Role,
		Status:    user.Status,
		SchoolID:  user.SchoolID,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (uc *profileUseCase) UpdateProfile(userID int, req *dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	//! 1. Get user
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	//! 2. Validate input
	if req.FirstName == "" {
		return nil, errors.New("first name is required")
	}
	if req.LastName == "" {
		return nil, errors.New("last name is required")
	}

	//! 3. Update fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	//! 4. Save
	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	//! 5. Return updated profile
	return uc.GetProfile(userID)
}

func (uc *profileUseCase) ChangePassword(userID int, req *dto.ChangePasswordRequest) error {
	//! 1. Validate input
	if req.CurrentPassword == "" {
		return errors.New("current password is required")
	}
	if req.NewPassword == "" {
		return errors.New("new password is required")
	}
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	if len(req.NewPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	//! 2. Get user
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	//! 3. Verify current password
	if !user.VerifyPassword(req.CurrentPassword) {
		return errors.New("current password is incorrect")
	}

	//! 4. Hash new password
	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	//! 5. Save
	return uc.userRepo.Update(user)
}

func (uc *profileUseCase) UpdateAvatar(userID int, avatarURL string) error {
	return uc.userRepo.UpdateAvatar(userID, avatarURL)
}

func (uc *profileUseCase) GetSchool(schoolID int) (*domain.School, error) {
	return uc.schoolRepo.FindByID(schoolID)
}

func (uc *profileUseCase) UpdateSchool(schoolID int, req *dto.UpdateSchoolRequest) (*domain.School, error) {
	school, err := uc.schoolRepo.FindByID(schoolID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		school.Name = req.Name
	}
	if req.Address != "" {
		school.Address = req.Address
	}
	if req.Phone != "" {
		school.Phone = req.Phone
	}
	if req.Email != "" {
		school.Email = req.Email
	}

	if err := uc.schoolRepo.Update(school); err != nil {
		return nil, err
	}

	return school, nil
}

func (uc *profileUseCase) UpdateSchoolLogo(schoolID int, logoURL string) error {
	return uc.schoolRepo.UpdateLogo(schoolID, logoURL)
}

func (uc *profileUseCase) GetTeacherSubjects(userID int) (*dto.TeacherSubjectsResponse, error) {
	// 1. Get user and verify it's a teacher
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsTeacher() {
		return nil, errors.New("user is not a teacher")
	}

	// 2. Get teacher's subject IDs
	subjectIDs, err := uc.teacherSubjectRepo.GetTeacherSubjects(userID)
	if err != nil {
		return nil, err
	}

	// 3. Get subject details
	subjects := []dto.SubjectInfo{}
	for _, subjectID := range subjectIDs {
		subject, err := uc.subjectRepo.FindByID(subjectID)
		if err != nil {
			continue // Skip if not found
		}

		subjects = append(subjects, dto.SubjectInfo{
			ID:          subject.ID,
			Name:        subject.Name,
			Code:        subject.Code,
			Description: subject.Description,
		})
	}

	return &dto.TeacherSubjectsResponse{
		Subjects: subjects,
		Total:    len(subjects),
	}, nil
}

func (uc *profileUseCase) GetStudentClass(userID int) (*dto.StudentClassResponse, error) {
	// 1. Get user and verify it's a student
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsStudent() {
		return nil, errors.New("user is not a student")
	}

	// 2. Get student's class
	classInfo, err := uc.studentClassRepo.GetStudentClass(userID)
	if err != nil {
		return nil, err
	}

	if classInfo == nil {
		return nil, fmt.Errorf("student not assigned to any class")
	}

	// 3. Get class details
	class, err := uc.classRepo.FindByID(classInfo.ClassID)
	if err != nil {
		return nil, err
	}

	return &dto.StudentClassResponse{
		Class: &dto.ClassInfo{
			ID:           class.ID,
			Name:         class.Name,
			Level:        class.Level,
			Section:      class.Section,
			Capacity:     class.Capacity,
			AcademicYear: class.AcademicYear,
		},
	}, nil
}
