package usecase

import (
	"educnet/internal/domain"
	"educnet/internal/handler/dto"
	"educnet/internal/repository"

	"errors"
	"fmt"
)

type AdminUseCase interface {
	GetPendingUsers(adminUserID int) (*dto.PendingUsersResponse, error)
	ApproveUser(adminUserID, targetUserID int) error
	RejectUser(adminUserID, targetUserID int, reason string) error
	GetAllUsers(adminUserID int, filters map[string]string) (*dto.UserListResponse, error)

	GetAllSubjects(schoolID int) ([]dto.SubjectResponse, error)
	CreateSubject(adminUserID int, req *dto.CreateSubjectRequest) (*dto.SubjectResponse, error)
	UpdateSubject(adminUserID, subjectID int, req *dto.UpdateSubjectRequest) (*dto.SubjectResponse, error)
	DeleteSubject(adminUserID, subjectID int) error

	GetAll(schoolID int) ([]dto.ClassResponse, error)
	CreateClass(adminUserID int, req *dto.CreateClassRequest) (*dto.ClassResponse, error)
	UpdateClass(adminUserID, classID int, req *dto.UpdateClassRequest) (*dto.ClassResponse, error)
	DeleteClass(adminUserID, classID int) error

	GetDashboard(adminUserID int) (*dto.DashboardResponse, error)
}

type adminUseCase struct {
	userRepo           repository.UserRepository
	teacherSubjectRepo repository.TeacherSubjectRepository
	studentClassRepo   repository.StudentClassRepository
	subjectRepo        repository.SubjectRepository
	classRepo          repository.ClassRepository
}

func NewAdminUseCase(
	userRepo repository.UserRepository,
	teacherSubjectRepo repository.TeacherSubjectRepository,
	studentClassRepo repository.StudentClassRepository,
	subjectRepo repository.SubjectRepository,
	classRepo repository.ClassRepository,
) AdminUseCase {
	return &adminUseCase{
		userRepo:           userRepo,
		teacherSubjectRepo: teacherSubjectRepo,
		studentClassRepo:   studentClassRepo,
		subjectRepo:        subjectRepo,
		classRepo:          classRepo,
	}
}

func (uc *adminUseCase) GetPendingUsers(adminUserID int) (*dto.PendingUsersResponse, error) {
	//! 1. Get admin user to verify permissions and get school_id
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, err
	}

	if !admin.IsAdmin() {
		return nil, errors.New("unauthorized: admin role required")
	}

	//! 2. Get pending users from same school
	users, err := uc.userRepo.FindPendingBySchool(admin.SchoolID)
	if err != nil {
		return nil, err
	}

	//! 3. Build response with additional info
	pendingUsers := []dto.PendingUserInfo{}
	for _, user := range users {
		userInfo := dto.PendingUserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.GetFullName(),
			Role:      user.Role,
			Phone:     user.Phone,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		//! Add subjects for teachers
		if user.IsTeacher() {
			subjectIDs, _ := uc.teacherSubjectRepo.GetTeacherSubjects(user.ID)
			subjects := []string{}
			for _, subjectID := range subjectIDs {
				subject, err := uc.subjectRepo.FindByID(subjectID)
				if err == nil {
					subjects = append(subjects, subject.Name)
				}
			}
			userInfo.Subjects = subjects
		}

		//! Add class for students
		if user.IsStudent() {
			classInfo, _ := uc.studentClassRepo.GetStudentClass(user.ID)
			if classInfo != nil {
				userInfo.ClassName = classInfo.ClassName
			}
		}

		pendingUsers = append(pendingUsers, userInfo)
	}

	return &dto.PendingUsersResponse{
		Users: pendingUsers,
		Total: len(pendingUsers),
	}, nil
}

func (uc *adminUseCase) ApproveUser(adminUserID, targetUserID int) error {
	//! 1. Verify admin permissions
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("unauthorized: admin role required")
	}

	//! 2. Get target user
	targetUser, err := uc.userRepo.FindByID(targetUserID)
	if err != nil {
		return err
	}

	//! 3. Verify same school
	if targetUser.SchoolID != admin.SchoolID {
		return errors.New("unauthorized: cannot approve users from other schools")
	}

	//! 4. Verify user is pending
	if !targetUser.IsPending() {
		return fmt.Errorf("user is not pending (status: %s)", targetUser.Status)
	}

	//! 5. Approve user
	targetUser.Approve()

	//! 6. Save
	return uc.userRepo.Update(targetUser)
}

func (uc *adminUseCase) RejectUser(adminUserID, targetUserID int, reason string) error {
	//! 1. Verify admin permissions
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("unauthorized: admin role required")
	}

	//! 2. Get target user
	targetUser, err := uc.userRepo.FindByID(targetUserID)
	if err != nil {
		return err
	}

	//! 3. Verify same school
	if targetUser.SchoolID != admin.SchoolID {
		return errors.New("unauthorized: cannot reject users from other schools")
	}

	//! 4. Verify user is pending
	if !targetUser.IsPending() {
		return fmt.Errorf("user is not pending (status: %s)", targetUser.Status)
	}

	//! 5. Reject user
	targetUser.Reject()

	//! 6. Save (TODO: store rejection reason in a separate table if needed)
	return uc.userRepo.Update(targetUser)
}

func (uc *adminUseCase) GetAllUsers(adminUserID int, filters map[string]string) (*dto.UserListResponse, error) {
	//! 1. Verify admin permissions
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, err
	}

	if !admin.IsAdmin() {
		return nil, errors.New("unauthorized: admin role required")
	}

	//! 2. Get users from same school
	users, err := uc.userRepo.FindBySchool(admin.SchoolID, filters)
	if err != nil {
		return nil, err
	}

	//! 3. Build response
	userList := []dto.UserListInfo{}
	for _, user := range users {
		userList = append(userList, dto.UserListInfo{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.GetFullName(),
			Role:      user.Role,
			Status:    user.Status,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.UserListResponse{
		Users: userList,
		Total: len(userList),
	}, nil
}

// ! ========== SUBJECTS ==========
func (uc *adminUseCase) CreateSubject(adminUserID int, req *dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	//! 1. Verify admin
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, err
	}

	if !admin.IsAdmin() {
		return nil, errors.New("unauthorized: admin role required")
	}

	//! 2. Validate input
	if req.Name == "" {
		return nil, errors.New("subject name is required")
	}
	if req.Code == "" {
		return nil, errors.New("subject code is required")
	}

	//! 3. Check if code already exists
	exists, err := uc.subjectRepo.ExistsByCode(admin.SchoolID, req.Code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("subject code '%s' already exists", req.Code)
	}

	//! 4. Create subject
	subject, err := domain.NewSubject(admin.SchoolID, req.Name, req.Code)
	if err != nil {
		return nil, err
	}
	subject.Description = req.Description

	if err := uc.subjectRepo.Create(subject); err != nil {
		return nil, err
	}

	//! 5. Return response
	return &dto.SubjectResponse{
		ID:          subject.ID,
		Name:        subject.Name,
		Code:        subject.Code,
		Description: subject.Description,
		SchoolID:    subject.SchoolID,
	}, nil
}

func (uc *adminUseCase) UpdateSubject(adminUserID, subjectID int, req *dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	//! 1. Verify admin
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, err
	}

	if !admin.IsAdmin() {
		return nil, errors.New("unauthorized: admin role required")
	}

	//! 2. Get subject
	subject, err := uc.subjectRepo.FindByID(subjectID)
	if err != nil {
		return nil, err
	}

	//! 3. Verify same school
	if subject.SchoolID != admin.SchoolID {
		return nil, errors.New("unauthorized: cannot modify subjects from other schools")
	}

	//! 4. Validate input
	if req.Name == "" {
		return nil, errors.New("subject name is required")
	}
	if req.Code == "" {
		return nil, errors.New("subject code is required")
	}

	//! 5. Check if new code conflicts
	if req.Code != subject.Code {
		exists, err := uc.subjectRepo.ExistsByCode(admin.SchoolID, req.Code, subjectID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("subject code '%s' already exists", req.Code)
		}
	}

	//! 6. Update subject
	subject.Name = req.Name
	subject.Code = req.Code
	subject.Description = req.Description

	if err := uc.subjectRepo.Update(subject); err != nil {
		return nil, err
	}

	//! 7. Return response
	return &dto.SubjectResponse{
		ID:          subject.ID,
		Name:        subject.Name,
		Code:        subject.Code,
		Description: subject.Description,
		SchoolID:    subject.SchoolID,
	}, nil
}

func (uc *adminUseCase) DeleteSubject(adminUserID, subjectID int) error {
	//! 1. Verify admin
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("unauthorized: admin role required")
	}

	//! 2. Get subject
	subject, err := uc.subjectRepo.FindByID(subjectID)
	if err != nil {
		return err
	}

	//! 3. Verify same school
	if subject.SchoolID != admin.SchoolID {
		return errors.New("unauthorized: cannot delete subjects from other schools")
	}

	//! 4. Delete subject
	return uc.subjectRepo.Delete(subjectID)
}

func (uc *adminUseCase) GetAllSubjects(schoolID int) ([]dto.SubjectResponse, error) {
	subjects, err := uc.subjectRepo.FindBySchoolID(schoolID)
	if err != nil {
		return nil, err
	}

	return dto.SubjectResponsesFromDomain(subjects), nil
}

//! ========== CLASSES ==========

func (uc *adminUseCase) GetAll(schoolID int) ([]dto.ClassResponse, error) {
	classes, err := uc.classRepo.GetAll(schoolID)
	if err != nil {
		return nil, err
	}

	return dto.ClassResponsesFromDomain(classes), nil
}

func (uc *adminUseCase) CreateClass(adminUserID int, req *dto.CreateClassRequest) (*dto.ClassResponse, error) {
	//! 1. Verify admin
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, err
	}

	if !admin.IsAdmin() {
		return nil, errors.New("unauthorized: admin role required")
	}

	//! 2. Validate input
	if req.Name == "" {
		return nil, errors.New("class name is required")
	}
	if req.Level == "" {
		return nil, errors.New("class level is required")
	}
	if req.AcademicYear == "" {
		return nil, errors.New("academic year is required")
	}
	if req.Capacity <= 0 {
		req.Capacity = 40 //! Default
	}

	//! 3. Create class
	class, err := domain.NewClass(admin.SchoolID, req.Name, req.Level, req.Section, req.AcademicYear)
	if err != nil {
		return nil, err
	}
	class.Section = req.Section
	class.Capacity = req.Capacity

	if err := uc.classRepo.Create(class); err != nil {
		return nil, err
	}

	//! 4. Return response
	return &dto.ClassResponse{
		ID:           class.ID,
		Name:         class.Name,
		Level:        class.Level,
		Section:      class.Section,
		Capacity:     class.Capacity,
		AcademicYear: class.AcademicYear,
		SchoolID:     class.SchoolID,
	}, nil
}

func (uc *adminUseCase) UpdateClass(adminUserID, classID int, req *dto.UpdateClassRequest) (*dto.ClassResponse, error) {
	//! 1. Verify admin
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, err
	}

	if !admin.IsAdmin() {
		return nil, errors.New("unauthorized: admin role required")
	}

	//! 2. Get class
	class, err := uc.classRepo.FindByID(classID)
	if err != nil {
		return nil, err
	}

	//! 3. Verify same school
	if class.SchoolID != admin.SchoolID {
		return nil, errors.New("unauthorized: cannot modify classes from other schools")
	}

	//! 4. Validate input
	if req.Name == "" {
		return nil, errors.New("class name is required")
	}
	if req.Level == "" {
		return nil, errors.New("class level is required")
	}
	if req.AcademicYear == "" {
		return nil, errors.New("academic year is required")
	}

	//! 5. Update class
	class.Name = req.Name
	class.Level = req.Level
	class.Section = req.Section
	class.Capacity = req.Capacity
	class.AcademicYear = req.AcademicYear

	if err := uc.classRepo.Update(class); err != nil {
		return nil, err
	}

	//! 6. Return response
	return &dto.ClassResponse{
		ID:           class.ID,
		Name:         class.Name,
		Level:        class.Level,
		Section:      class.Section,
		Capacity:     class.Capacity,
		AcademicYear: class.AcademicYear,
		SchoolID:     class.SchoolID,
	}, nil
}

func (uc *adminUseCase) DeleteClass(adminUserID, classID int) error {
	//! 1. Verify admin
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return errors.New("unauthorized: admin role required")
	}

	//! 2. Get class
	class, err := uc.classRepo.FindByID(classID)
	if err != nil {
		return err
	}

	//! 3. Verify same school
	if class.SchoolID != admin.SchoolID {
		return errors.New("unauthorized: cannot delete classes from other schools")
	}

	//! 4. Delete class
	return uc.classRepo.Delete(classID)
}

func (uc *adminUseCase) GetDashboard(adminUserID int) (*dto.DashboardResponse, error) {
	// 1. Verify admin
	admin, err := uc.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, err
	}

	if !admin.IsAdmin() {
		return nil, errors.New("unauthorized: admin role required")
	}

	// 2. Get all users from school
	allUsers, err := uc.userRepo.FindBySchool(admin.SchoolID, map[string]string{})
	if err != nil {
		return nil, err
	}

	// 3. Calculate stats
	stats := dto.DashboardStats{}
	stats.TotalUsers = len(allUsers)

	for _, user := range allUsers {
		switch user.Role {
		case "teacher":
			stats.TotalTeachers++
		case "student":
			stats.TotalStudents++
		case "admin":
			stats.TotalAdmins++
		}

		switch user.Status {
		case "pending":
			stats.PendingUsers++
		case "approved":
			stats.ApprovedUsers++
		case "rejected":
			stats.RejectedUsers++
		}
	}

	// 4. Get subjects and classes count
	subjects, _ := uc.subjectRepo.FindBySchoolID(admin.SchoolID)
	stats.TotalSubjects = len(subjects)

	classes, _ := uc.classRepo.FindBySchoolID(admin.SchoolID)
	stats.TotalClasses = len(classes)

	// 5. Get pending users (limit 5 for dashboard)
	pendingUsers, err := uc.userRepo.FindPendingBySchool(admin.SchoolID)
	if err != nil {
		return nil, err
	}

	pendingList := []dto.PendingUserInfo{}
	limit := 5
	if len(pendingUsers) < limit {
		limit = len(pendingUsers)
	}

	for i := 0; i < limit; i++ {
		user := pendingUsers[i]
		pendingList = append(pendingList, dto.PendingUserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.GetFullName(),
			Role:      user.Role,
			Phone:     user.Phone,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.DashboardResponse{
		Stats:        stats,
		PendingUsers: pendingList,
	}, nil
}
