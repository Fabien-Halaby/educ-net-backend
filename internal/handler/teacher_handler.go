// package handler

// import (
// 	"educnet/internal/handler/dto"
// 	"educnet/internal/usecase"
// 	"educnet/internal/utils"
// 	"encoding/json"
// 	"net/http"
// )

// type TeacherHandler struct {
// 	teacherUC usecase.TeacherUseCase
// }

// func NewTeacherHandler(teacherUC usecase.TeacherUseCase) *TeacherHandler {
// 	return &TeacherHandler{teacherUC: teacherUC}
// }

// //! POST /api/teachers/register
// func (h *TeacherHandler) Register(w http.ResponseWriter, r *http.Request) {
// 	var req dto.TeacherRegistrationRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.BadRequest(w, "Invalid request body")
// 		return
// 	}

// 	//! Validate required fields
// 	if req.SchoolSlug == "" {
// 		utils.BadRequest(w, "school_slug is required")
// 		return
// 	}
// 	if req.Email == "" {
// 		utils.BadRequest(w, "email is required")
// 		return
// 	}
// 	if req.Password == "" {
// 		utils.BadRequest(w, "password is required")
// 		return
// 	}
// 	if req.FirstName == "" {
// 		utils.BadRequest(w, "first_name is required")
// 		return
// 	}
// 	if len(req.SubjectIDs) == 0 {
// 		utils.BadRequest(w, "at least one subject is required")
// 		return
// 	}

// 	//! Register teacher
// 	resp, err := h.teacherUC.RegisterTeacher(&req)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.Created(w, "Teacher created successfully", resp)
// }

package handler

import (
	"encoding/json"
	"net/http"

	"educnet/internal/handler/dto"
	"educnet/internal/middleware"
	"educnet/internal/usecase"
	"educnet/internal/utils"
)

type TeacherHandler struct {
	teacherUC usecase.TeacherUseCase
}

func NewTeacherHandler(teacherUC usecase.TeacherUseCase) *TeacherHandler {
	return &TeacherHandler{teacherUC: teacherUC}
}

// POST /api/teachers/register
func (h *TeacherHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.TeacherRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	resp, err := h.teacherUC.RegisterTeacher(&req)
	if err != nil {
		utils.HandleUseCaseError(w, err) // 400/409/422/500 auto
		return
	}

	utils.Created(w, "Teacher registered successfully", resp)
}

func (h *TeacherHandler) GetMySubjects(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	subjects, err := h.teacherUC.GetTeacherSubjects(claims.UserID)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Subjects retrieved", subjects)
}
