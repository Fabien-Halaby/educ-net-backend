package handler

import (
	"educnet/internal/usecase"
	"educnet/internal/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ! ClassHandler gère les requêtes HTTP liées aux écoles
type ClassHandler struct {
	classUseCase usecase.ClassUseCase
}

// ! NewClassHandler crée un nouveau handler
func NewClassHandler(ClassUseCase usecase.ClassUseCase) *ClassHandler {
	return &ClassHandler{
		classUseCase: ClassUseCase,
	}
}

func (h *ClassHandler) GetClassesBySchoolID(w http.ResponseWriter, r *http.Request) {
	idS := mux.Vars(r)["schoolId"]
	id, err := strconv.Atoi(idS)
	if err != nil {
		utils.BadRequest(w, "Invalid school ID")
		return
	}

	out, err := h.classUseCase.GetAllBySchoolID(id)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}

	utils.OK(w, "Classes retrieved successfully", out)
}
