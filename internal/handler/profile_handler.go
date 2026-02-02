package handler

import (
	"educnet/internal/handler/dto"
	"educnet/internal/middleware"
	"educnet/internal/usecase"
	"educnet/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProfileHandler struct {
	profileUC usecase.ProfileUseCase
}

func NewProfileHandler(profileUC usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{profileUC: profileUC}
}

// ! GET /api/me (déjà existe dans user_handler, mais on peut l'améliorer)
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	profile, err := h.profileUC.GetProfile(claims.UserID)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Profile retrieved", profile)
}

// ! PUT /api/me
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

// ! PUT /api/me/password
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
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Password changed successfully", nil)
}

// ! POST /api/me/avatar
func (h *ProfileHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	//! Parse multipart form (max 5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.BadRequest(w, "File too large (max 5MB)")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		utils.BadRequest(w, "No file uploaded")
		return
	}
	defer file.Close()

	//! Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		utils.BadRequest(w, "Only JPG, PNG, WEBP allowed")
		return
	}

	//! Generate unique filename
	filename := fmt.Sprintf("avatar_%d_%d%s", claims.UserID, time.Now().Unix(), ext)
	uploadDir := "./uploads/avatars"

	//! Create directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.InternalServerError(w, "Failed to create upload directory")
		return
	}

	//! Save file
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerError(w, "Failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.InternalServerError(w, "Failed to save file")
		return
	}

	//! Update database
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := h.profileUC.UpdateAvatar(claims.UserID, avatarURL); err != nil {
		utils.InternalServerError(w, "Failed to update avatar")
		return
	}

	utils.OK(w, "Avatar uploaded successfully", dto.UploadResponse{URL: avatarURL})
}

// ! GET /api/me/school
func (h *ProfileHandler) GetSchool(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	//! Get user to get school_id
	user, err := h.profileUC.GetProfile(claims.UserID)
	if err != nil {
		utils.NotFound(w, "User not found")
		return
	}

	school, err := h.profileUC.GetSchool(user.ID, user.SchoolID)
	if err != nil {
		utils.NotFound(w, "School not found")
		return
	}

	utils.OK(w, "School retrieved", school)
}

// PUT /api/me/school (Admin only)
func (h *ProfileHandler) UpdateSchool(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	// Check if admin
	if claims.Role != "admin" {
		utils.Forbidden(w, "Only admin can update school")
		return
	}

	//! Get user to get school_id
	user, err := h.profileUC.GetProfile(claims.UserID)
	if err != nil {
		utils.NotFound(w, "User not found")
		return
	}

	var req dto.UpdateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	school, err := h.profileUC.UpdateSchool(user.ID, user.SchoolID, &req)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "School updated successfully", school)
}

// POST /api/me/school/logo (Admin only)
func (h *ProfileHandler) UploadSchoolLogo(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	// Check if admin
	if claims.Role != "admin" {
		utils.Forbidden(w, "Only admin can update school logo")
		return
	}

	// Get user to get school_id
	user, err := h.profileUC.GetProfile(claims.UserID)
	if err != nil {
		utils.NotFound(w, "User not found")
		return
	}

	// Parse multipart form (max 5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.BadRequest(w, "File too large (max 5MB)")
		return
	}

	file, header, err := r.FormFile("logo")
	if err != nil {
		utils.BadRequest(w, "No file uploaded")
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".svg" {
		utils.BadRequest(w, "Only JPG, PNG, WEBP, SVG allowed")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("logo_%d_%d%s", user.SchoolID, time.Now().Unix(), ext)
	uploadDir := "./uploads/logos"

	// Create directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.InternalServerError(w, "Failed to create upload directory")
		return
	}

	// Save file
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerError(w, "Failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.InternalServerError(w, "Failed to save file")
		return
	}

	// Update database
	logoURL := fmt.Sprintf("/uploads/logos/%s", filename)
	if err := h.profileUC.UpdateSchoolLogo(user.ID, user.SchoolID, logoURL); err != nil {
		utils.InternalServerError(w, "Failed to update logo")
		return
	}

	utils.OK(w, "Logo uploaded successfully", dto.UploadResponse{URL: logoURL})
}

// ! GET /api/me/subjects (for teachers)
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

// ! GET /api/me/class (for students)
func (h *ProfileHandler) GetMyClass(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized")
		return
	}

	class, err := h.profileUC.GetStudentClasses(claims.UserID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.OK(w, "Class retrieved", class)
}
