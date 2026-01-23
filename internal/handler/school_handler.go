package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"educnet/internal/domain"
	"educnet/internal/handler/dto"
	"educnet/internal/usecase"
	"educnet/internal/utils"
)

//! SchoolHandler gère les requêtes HTTP liées aux écoles
type SchoolHandler struct {
	schoolUseCase usecase.SchoolUseCase
}

//! NewSchoolHandler crée un nouveau handler
func NewSchoolHandler(schoolUseCase usecase.SchoolUseCase) *SchoolHandler {
	return &SchoolHandler{
		schoolUseCase: schoolUseCase,
	}
}

//! CreateSchool - POST /api/schools/register
func (h *SchoolHandler) CreateSchool(w http.ResponseWriter, r *http.Request) {
	//! 1. Parse request
	var req dto.CreateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	//! 2. Mapper vers use case input
	input := usecase.CreateSchoolInput{
		SchoolName:    req.SchoolName,
		AdminEmail:    req.AdminEmail,
		AdminPassword: req.AdminPassword,
		AdminName:     req.AdminName,
		Phone:         req.Phone,
		Address:       req.Address,
	}

	//! 3. Exécuter use case
	output, err := h.schoolUseCase.CreateSchool(input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	//! 4. Mapper vers response DTO
	response := dto.CreateSchoolResponse{
		School: dto.SchoolDTO{
			ID:          output.School.ID,
			Name:        output.School.Name,
			Slug:        output.School.Slug,
			Address:     output.School.Address,
			Phone:       output.School.Phone,
			Status:      output.School.Status,
			AdminUserID: output.School.AdminUserID,
			CreatedAt:   output.School.CreatedAt,
		},
		Admin: dto.UserDTO{
			ID:        output.Admin.ID,
			SchoolID:  output.Admin.SchoolID,
			Email:     output.Admin.Email,
			FirstName: output.Admin.FirstName,
			LastName:  output.Admin.LastName,
			Phone:     output.Admin.Phone,
			Role:      output.Admin.Role,
			Status:    output.Admin.Status,
			CreatedAt: output.Admin.CreatedAt,
		},
		Token: output.Token,
	}

	//! 5. Répondre avec succès
	utils.Created(w, "School created successfully", response)
}

//! handleError mappe les erreurs domain vers HTTP status codes
func (h *SchoolHandler) handleError(w http.ResponseWriter, err error) {
	log.Println("Error:", err)
	//! Type assertion vers DomainError
	if domainErr, ok := err.(*domain.DomainError); ok {
		switch domainErr {
		case domain.ErrSchoolNameRequired,
			domain.ErrUserEmailRequired,
			domain.ErrUserEmailInvalid,
			domain.ErrUserPasswordTooShort,
			domain.ErrUserNameRequired:
			utils.BadRequest(w, domainErr.Message)
		case domain.ErrSchoolAlreadyExists,
			domain.ErrUserAlreadyExists:
			utils.Conflict(w, domainErr.Message)
		case domain.ErrSchoolNotFound,
			domain.ErrUserNotFound:
			utils.NotFound(w, domainErr.Message)
		default:
			utils.InternalServerError(w, "An error occurred")
		}
		return
	}

	//! Erreur générique
	utils.InternalServerError(w, "An error occurred")
}
