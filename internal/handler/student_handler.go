// package handler

// import (
// 	"educnet/internal/handler/dto"
// 	"educnet/internal/usecase"
// 	"educnet/internal/utils"
// 	"encoding/json"
// 	"net/http"
// )

// type StudentHandler struct {
// 	studentUC usecase.StudentUseCase
// }

// func NewStudentHandler(studentUC usecase.StudentUseCase) *StudentHandler {
// 	return &StudentHandler{studentUC: studentUC}
// }

// //! POST /api/students/register
// func (h *StudentHandler) Register(w http.ResponseWriter, r *http.Request) {
// 	var req dto.StudentRegistrationRequest
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
// 	if req.ClassID == 0 {
// 		utils.BadRequest(w, "class_id is required")
// 		return
// 	}

// 	//! Register student
// 	resp, err := h.studentUC.RegisterStudent(&req)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.Created(w, "Student created successfully", resp)
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

type StudentHandler struct {
	studentUC usecase.StudentUseCase
}

func NewStudentHandler(studentUC usecase.StudentUseCase) *StudentHandler {
	return &StudentHandler{studentUC: studentUC}
}

// POST /api/students/register
func (h *StudentHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.StudentRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	resp, err := h.studentUC.RegisterStudent(&req)
	if err != nil {
		utils.HandleUseCaseError(w, err) // 400/409/422/500 auto
		return
	}

	utils.Created(w, "Student registered successfully", resp)
}

func (h *StudentHandler) GetMyClasses(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	classes, err := h.studentUC.GetStudentClasses(claims.UserID)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Classes retrieved", classes)
}
