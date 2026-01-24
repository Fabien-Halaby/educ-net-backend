package handler

import (
	"educnet/internal/middleware"
	"educnet/internal/usecase"
	"educnet/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AdminHandler struct {
	adminUC usecase.AdminUseCase
}

func NewAdminHandler(adminUC usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUC: adminUC}
}

//! GET /api/admin/users/pending
func (h *AdminHandler) GetPendingUsers(w http.ResponseWriter, r *http.Request) {
	//! Get admin from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	//! Get pending users
	resp, err := h.adminUC.GetPendingUsers(claims.UserID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "Pending users retrieved", resp)
}

//! POST /api/admin/users/{id}/approve
func (h *AdminHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	//! Get admin from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	//! Get user ID from URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid user ID")
		return
	}

	//! Approve user
	if err := h.adminUC.ApproveUser(claims.UserID, userID); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "User approved successfully", nil)
}

//! POST /api/admin/users/{id}/reject
func (h *AdminHandler) RejectUser(w http.ResponseWriter, r *http.Request) {
	//! Get admin from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	//! Get user ID from URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid user ID")
		return
	}

	//! Get optional reason
	var req struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	//! Reject user
	if err := h.adminUC.RejectUser(claims.UserID, userID, req.Reason); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "User rejected successfully", nil)
}

//! GET /api/admin/users
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	//! Get admin from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	//! Get query params
	filters := map[string]string{
		"role":   r.URL.Query().Get("role"),
		"status": r.URL.Query().Get("status"),
	}

	//! Get users
	resp, err := h.adminUC.GetAllUsers(claims.UserID, filters)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "Users retrieved", resp)
}
