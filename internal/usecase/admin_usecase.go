package usecase

import (
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
