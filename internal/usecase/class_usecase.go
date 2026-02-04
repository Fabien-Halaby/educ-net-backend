package usecase

import (
	"educnet/internal/handler/dto"
	"educnet/internal/repository"
	"fmt"
)

// ! ClassUseCase interface
type ClassUseCase interface {
	GetAllBySchoolID(schoolID int) ([]*dto.ClassInfo, error)
}

// ! classUseCase implémentation
type classUseCase struct {
	classRepo repository.ClassRepository
}

func NewClassUsecase(classRepo repository.ClassRepository) ClassUseCase {
	return &classUseCase{
		classRepo: classRepo,
	}
}

func (uc *classUseCase) GetAllBySchoolID(schoolID int) ([]*dto.ClassInfo, error) {
	classes, err := uc.classRepo.FindBySchoolID(schoolID)
	if err != nil {
		return nil, fmt.Errorf("[ERROR_GET_CLASSES_BY_SCHOOL_ID]: %w", err)
	}

	var resp []*dto.ClassInfo
	for _, c := range classes {
		resp = append(resp, dto.ClassInfoFromDomain(c))
	}

	return resp, nil
}
