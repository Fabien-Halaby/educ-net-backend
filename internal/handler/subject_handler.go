package handler

import (
	"educnet/internal/usecase"
	"educnet/internal/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ! SubjectHandler gère les requêtes HTTP liées aux écoles
type SubjectHandler struct {
	subjectUseCase usecase.SubjectUseCase
}

// ! NewSubjectHandler crée un nouveau handler
func NewSubjectHandler(subjectUseCase usecase.SubjectUseCase) *SubjectHandler {
	return &SubjectHandler{
		subjectUseCase: subjectUseCase,
	}
}

func (h *SubjectHandler) GetSubjectsBySchoolID(w http.ResponseWriter, r *http.Request) {
	idS := mux.Vars(r)["schoolId"]
	id, err := strconv.Atoi(idS)
	if err != nil {
		utils.BadRequest(w, "Invalid school ID")
		return
	}

	out, err := h.subjectUseCase.GetAllBySchoolID(id)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Subjectes retrieved successfully", out)
}
