package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"educnet/internal/usecase"
	"educnet/internal/utils"
)

// ! SchoolHandler gère les requêtes HTTP liées aux écoles
type SchoolHandler struct {
	schoolUseCase usecase.SchoolUseCase
}

// ! NewSchoolHandler crée un nouveau handler
func NewSchoolHandler(schoolUseCase usecase.SchoolUseCase) *SchoolHandler {
	return &SchoolHandler{
		schoolUseCase: schoolUseCase,
	}
}

func (h *SchoolHandler) CreateSchool(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateSchoolInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	fmt.Println("SCHOOL: ", input)

	output, err := h.schoolUseCase.CreateSchool(input)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.Created(w, "School created successfully", output)
}

func (h *SchoolHandler) GetAllSchool(w http.ResponseWriter, r *http.Request) {
	out, err := h.schoolUseCase.GetAllSchool()
	if err != nil {
		utils.HandleUseCaseError(w, err)
	}

	utils.OK(w, "Schools retrieved successfully", out)
}

// ! handleError mappe les erreurs domain vers HTTP status codes
func (h *SchoolHandler) handleError(w http.ResponseWriter, err error) {
	utils.HandleUseCaseError(w, err)
}
