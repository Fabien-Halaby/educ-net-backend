package usecase

import (
	"educnet/internal/domain"
	"educnet/internal/handler/dto"
	"educnet/internal/repository"
	"errors"
)

type ProfileUseCase interface {
	GetProfile(userID int) (*dto.ProfileResponse, error)
	UpdateProfile(userID int, req *dto.UpdateProfileRequest) (*dto.ProfileResponse, error)
	ChangePassword(userID int, req *dto.ChangePasswordRequest) error

	UpdateAvatar(userID int, avatarURL string) error
	GetSchool(userID, schoolID int) (*domain.School, error)
	UpdateSchool(userID, schoolID int, req *dto.UpdateSchoolRequest) (*domain.School, error)
	UpdateSchoolLogo(userID, schoolID int, logoURL string) error

	GetTeacherSubjects(userID int) (*dto.TeacherSubjectsResponse, error)
	GetStudentClasses(userID int) (*dto.StudentClassesResponse, error)
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
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrInternal
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
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrInternal
	}

	//! 2. Validate input
	if req.FirstName == "" {
		return nil, domain.ErrNameRequired
	}
	if req.LastName == "" {
		return nil, domain.ErrNameRequired
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
		return domain.ErrPasswordRequired
	}
	if req.NewPassword == "" {
		return domain.ErrPasswordRequired
	}
	if req.NewPassword != req.ConfirmPassword {
		return domain.ErrPasswordDontMatch
	}
	if len(req.NewPassword) < 8 {
		return domain.ErrPasswordTooShort
	}

	//! 2. Get user
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	//! 3. Verify current password
	if !user.VerifyPassword(req.CurrentPassword) {
		return domain.ErrInvalidCredentials
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

func (uc *profileUseCase) GetSchool(userID, schoolID int) (*domain.School, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.SchoolID != schoolID {
		return nil, domain.ErrUnauthorized
	}
	return uc.schoolRepo.FindByID(schoolID)
}

func (uc *profileUseCase) UpdateSchool(userID, schoolID int, req *dto.UpdateSchoolRequest) (*domain.School, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil || !user.IsAdmin() {
		return nil, domain.ErrForbidden
	}
	if user.SchoolID != schoolID {
		return nil, domain.ErrUnauthorized
	}

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

func (uc *profileUseCase) UpdateSchoolLogo(userID, schoolID int, logoURL string) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil || !user.IsAdmin() {
		return domain.ErrForbidden
	}
	if user.SchoolID != schoolID {
		return domain.ErrForbidden
	}
	return uc.schoolRepo.UpdateLogo(schoolID, logoURL)
}

func (uc *profileUseCase) GetTeacherSubjects(userID int) (*dto.TeacherSubjectsResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil || !user.IsTeacher() {
		return nil, domain.ErrForbidden
	}

	subjects, err := uc.teacherSubjectRepo.FindByTeacher(userID)
	if err != nil {
		return nil, err
	}

	subjectsInfo := make([]dto.SubjectInfo, len(subjects))
	for i, subj := range subjects {
		subjectsInfo[i] = dto.SubjectInfo{
			ID:          subj.ID,
			Name:        subj.Name,
			Code:        subj.Code,
			Description: subj.Description,
		}
	}

	return &dto.TeacherSubjectsResponse{
		Subjects: subjectsInfo,
		Total:    len(subjectsInfo),
	}, nil
}

func (uc *profileUseCase) GetStudentClasses(userID int) (*dto.StudentClassesResponse, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil || !user.IsStudent() {
		return nil, domain.ErrForbidden
	}

	classes, err := uc.studentClassRepo.FindByStudent(userID)
	if err != nil {
		return nil, err
	}

	classesInfo := make([]*dto.ClassInfo, len(classes))
	for i, class := range classes {
		classesInfo[i] = &dto.ClassInfo{
			ID:           class.ID,
			Name:         class.Name,
			Level:        class.Level,
			Section:      class.Section,
			Capacity:     class.Capacity,
			AcademicYear: class.AcademicYear,
		}
	}

	return &dto.StudentClassesResponse{
		Classes: classesInfo,
		Total:   len(classesInfo),
	}, nil
}
