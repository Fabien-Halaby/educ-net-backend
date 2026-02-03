// package handler

// import (
// 	"educnet/internal/handler/dto"
// 	"educnet/internal/middleware"
// 	"educnet/internal/usecase"
// 	"educnet/internal/utils"
// 	"encoding/json"
// 	"net/http"
// 	"strconv"

// 	"github.com/gorilla/mux"
// )

// type AdminHandler struct {
// 	adminUC usecase.AdminUseCase
// }

// func NewAdminHandler(adminUC usecase.AdminUseCase) *AdminHandler {
// 	return &AdminHandler{adminUC: adminUC}
// }

// // ! GET /api/admin/users/pending
// func (h *AdminHandler) GetPendingUsers(w http.ResponseWriter, r *http.Request) {
// 	//! Get admin from context
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	//! Get pending users
// 	resp, err := h.adminUC.GetPendingUsers(claims.UserID)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Pending users retrieved", resp)
// }

// // ! POST /api/admin/users/{id}/approve
// func (h *AdminHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
// 	//! Get admin from context
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	//! Get user ID from URL
// 	vars := mux.Vars(r)
// 	userID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		utils.BadRequest(w, "Invalid user ID")
// 		return
// 	}

// 	//! Approve user
// 	if err := h.adminUC.ApproveUser(claims.UserID, userID); err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "User approved successfully", nil)
// }

// // ! POST /api/admin/users/{id}/reject
// func (h *AdminHandler) RejectUser(w http.ResponseWriter, r *http.Request) {
// 	//! Get admin from context
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	//! Get user ID from URL
// 	vars := mux.Vars(r)
// 	userID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		utils.BadRequest(w, "Invalid user ID")
// 		return
// 	}

// 	//! Get optional reason
// 	var req struct {
// 		Reason string `json:"reason"`
// 	}
// 	json.NewDecoder(r.Body).Decode(&req)

// 	//! Reject user
// 	if err := h.adminUC.RejectUser(claims.UserID, userID, req.Reason); err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "User rejected successfully", nil)
// }

// // ! GET /api/admin/users
// func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
// 	//! Get admin from context
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	//! Get query params
// 	filters := map[string]string{
// 		"role":   r.URL.Query().Get("role"),
// 		"status": r.URL.Query().Get("status"),
// 	}

// 	//! Get users
// 	resp, err := h.adminUC.GetAllUsers(claims.UserID, filters)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Users retrieved", resp)
// }

// // ========== SUBJECTS ==========

// // POST /api/admin/subjects
// func (h *AdminHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	var req dto.CreateSubjectRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.BadRequest(w, "Invalid request body")
// 		return
// 	}

// 	resp, err := h.adminUC.CreateSubject(claims.UserID, &req)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.Created(w, "Subject created successfully", resp)
// }

// // PUT /api/admin/subjects/{id}
// func (h *AdminHandler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	subjectID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		utils.BadRequest(w, "Invalid subject ID")
// 		return
// 	}

// 	var req dto.UpdateSubjectRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.BadRequest(w, "Invalid request body")
// 		return
// 	}

// 	resp, err := h.adminUC.UpdateSubject(claims.UserID, subjectID, &req)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Subject updated successfully", resp)
// }

// // DELETE /api/admin/subjects/{id}
// func (h *AdminHandler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	subjectID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		utils.BadRequest(w, "Invalid subject ID")
// 		return
// 	}

// 	if err := h.adminUC.DeleteSubject(claims.UserID, subjectID); err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Subject deleted successfully", nil)
// }

// func (h *AdminHandler) GetAllSubjects(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	subjects, err := h.adminUC.GetAllSubjects(claims.SchoolID)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Subjects retrieved", subjects)
// }

// // ! ========== CLASSES ==========
// func (h *AdminHandler) GetAllClasses(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	classes, err := h.adminUC.GetAllClasses(claims.SchoolID)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Classes retrieved", classes)
// }

// // ! POST /api/admin/classes
// func (h *AdminHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	var req dto.CreateClassRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.BadRequest(w, "Invalid request body")
// 		return
// 	}

// 	resp, err := h.adminUC.CreateClass(claims.UserID, &req)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.Created(w, "Class created successfully", resp)
// }

// // ! PUT /api/admin/classes/{id}
// func (h *AdminHandler) UpdateClass(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	classID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		utils.BadRequest(w, "Invalid class ID")
// 		return
// 	}

// 	var req dto.UpdateClassRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.BadRequest(w, "Invalid request body")
// 		return
// 	}

// 	resp, err := h.adminUC.UpdateClass(claims.UserID, classID, &req)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Class updated successfully", resp)
// }

// // ! DELETE /api/admin/classes/{id}
// func (h *AdminHandler) DeleteClass(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	classID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		utils.BadRequest(w, "Invalid class ID")
// 		return
// 	}

// 	if err := h.adminUC.DeleteClass(claims.UserID, classID); err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Class deleted successfully", nil)
// }

// // ! GET /api/admin/dashboard
// func (h *AdminHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := middleware.GetUserFromContext(r.Context())
// 	if !ok {
// 		utils.Unauthorized(w, "Unauthorized")
// 		return
// 	}

// 	dashboard, err := h.adminUC.GetDashboard(claims.UserID)
// 	if err != nil {
// 		utils.BadRequest(w, err.Error())
// 		return
// 	}

// 	utils.OK(w, "Dashboard retrieved", dashboard)
// }

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"educnet/internal/handler/dto"
	"educnet/internal/middleware"
	"educnet/internal/usecase"
	"educnet/internal/utils"

	"github.com/gorilla/mux"
)

/*
	====== ALL HANDLERS ======

	func (h *AdminHandler) ApproveUser(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) CreateClass(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) CreateSubject(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) DeleteClass(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) DeleteSubject(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) GetAllClasses(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) GetAllSubjects(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) GetDashboard(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) GetPendingUsers(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) RejectUser(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) UpdateClass(w http.ResponseWriter, r *http.Request)
	func (h *AdminHandler) UpdateSubject(w http.ResponseWriter, r *http.Request)

*/

type AdminHandler struct {
	adminUC usecase.AdminUseCase
}

func NewAdminHandler(adminUC usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUC: adminUC}
}

// GET /api/admin/users/pending
func (h *AdminHandler) GetPendingUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	resp, err := h.adminUC.GetPendingUsers(claims.UserID)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Pending users retrieved", resp)
}

// POST /api/admin/users/{id}/approve
func (h *AdminHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid user ID") // 400
		return
	}

	if err := h.adminUC.ApproveUser(claims.UserID, userID); err != nil {
		utils.HandleUseCaseError(w, err) // 403/404/500
		return
	}

	utils.OK(w, "User approved successfully", nil)
}

// POST /api/admin/users/{id}/reject
func (h *AdminHandler) RejectUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid user ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	if err := h.adminUC.RejectUser(claims.UserID, userID, req.Reason); err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "User rejected successfully", nil)
}

// GET /api/admin/users?role=student&status=pending
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	filters := map[string]string{
		"role":   r.URL.Query().Get("role"),
		"status": r.URL.Query().Get("status"),
	}

	resp, err := h.adminUC.GetAllUsers(claims.UserID, filters)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Users retrieved", resp)
}

// ========== SUBJECTS ==========

// POST /api/admin/subjects
func (h *AdminHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	var req dto.CreateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	resp, err := h.adminUC.CreateSubject(claims.UserID, &req)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.Created(w, "Subject created successfully", resp)
}

// PUT /api/admin/subjects/{id}
func (h *AdminHandler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	subjectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid subject ID")
		return
	}

	var req dto.UpdateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	resp, err := h.adminUC.UpdateSubject(claims.UserID, subjectID, &req)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Subject updated successfully", resp)
}

// DELETE /api/admin/subjects/{id}
func (h *AdminHandler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	subjectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid subject ID")
		return
	}

	if err := h.adminUC.DeleteSubject(claims.UserID, subjectID); err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Subject deleted successfully", nil)
}

// GET /api/admin/subjects
func (h *AdminHandler) GetAllSubjects(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	subjects, err := h.adminUC.GetAllSubjects(claims.SchoolID)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Subjects retrieved", subjects)
}

// ========== CLASSES ==========

// GET /api/admin/classes
func (h *AdminHandler) GetAllClasses(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	classes, err := h.adminUC.GetAllClasses(claims.SchoolID)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Classes retrieved", classes)
}

// POST /api/admin/classes
func (h *AdminHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	var req dto.CreateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	resp, err := h.adminUC.CreateClass(claims.UserID, &req)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.Created(w, "Class created successfully", resp)
}

// PUT /api/admin/classes/{id}
func (h *AdminHandler) UpdateClass(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	classID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid class ID")
		return
	}

	var req dto.UpdateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	resp, err := h.adminUC.UpdateClass(claims.UserID, classID, &req)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Class updated successfully", resp)
}

// DELETE /api/admin/classes/{id}
func (h *AdminHandler) DeleteClass(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	classID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid class ID")
		return
	}

	if err := h.adminUC.DeleteClass(claims.UserID, classID); err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Class deleted successfully", nil)
}

// GET /api/admin/dashboard
func (h *AdminHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	dashboard, err := h.adminUC.GetDashboard(claims.UserID)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Dashboard retrieved", dashboard)
}
