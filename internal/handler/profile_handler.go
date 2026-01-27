package handler

import (
	"educnet/internal/handler/dto"
	"educnet/internal/middleware"
	"educnet/internal/usecase"
	"educnet/internal/utils"
	"encoding/json"
	"net/http"
)

type ProfileHandler struct {
	profileUC usecase.ProfileUseCase
}

func NewProfileHandler(profileUC usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{profileUC: profileUC}
}

// GET /api/me (déjà existe dans user_handler, mais on peut l'améliorer)
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	profile, err := h.profileUC.GetProfile(claims.UserID)
	if err != nil {
		utils.NotFound(w, err.Error())
		return
	}

	utils.OK(w, "Profile retrieved", profile)
}

// PUT /api/me
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	profile, err := h.profileUC.UpdateProfile(claims.UserID, &req)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "Profile updated successfully", profile)
}

// PUT /api/me/password
func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	if err := h.profileUC.ChangePassword(claims.UserID, &req); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "Password changed successfully", nil)
}

// GET /api/me/subjects (for teachers)
func (h *ProfileHandler) GetMySubjects(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	subjects, err := h.profileUC.GetTeacherSubjects(claims.UserID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "Subjects retrieved", subjects)
}

// GET /api/me/class (for students)
func (h *ProfileHandler) GetMyClass(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	class, err := h.profileUC.GetStudentClass(claims.UserID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "Class retrieved", class)
}
