package handler

import (
	"educnet/internal/middleware"
	"educnet/internal/repository"
	"educnet/internal/utils"
	"net/http"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

//! GET /api/me - Get current user info
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	//! Get user from context (set by JWT middleware)
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	//! Get full user info from database
	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		utils.NotFound(w, "User not found")
		return
	}

	utils.OK(w, "User retrieved", map[string]interface{}{
		"id":        user.ID,
		"email":     user.Email,
		"full_name": user.GetFullName(),
		"role":      user.Role,
		"status":    user.Status,
		"school_id": user.SchoolID,
	})
}
