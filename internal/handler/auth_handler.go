package handler

import (
	"educnet/internal/handler/dto"
	"educnet/internal/usecase"
	"educnet/internal/utils"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// ! POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	//! Validate
	if req.Email == "" {
		utils.BadRequest(w, "email is required")
		return
	}
	if req.Password == "" {
		utils.BadRequest(w, "password is required")
		return
	}

	//! Login
	resp, err := h.authUC.Login(&req)
	if err != nil {
		utils.Unauthorized(w, err.Error())
		return
	}

	utils.OK(w, "Login successful", resp)
}
