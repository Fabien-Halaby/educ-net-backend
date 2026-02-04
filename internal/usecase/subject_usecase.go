package usecase

import (
	"educnet/internal/handler/dto"
	"educnet/internal/repository"
	"fmt"
)

// ! SubjectUseCase interface
type SubjectUseCase interface {
	GetAllBySchoolID(schoolID int) ([]*dto.SubjectInfo, error)
}

// ! SubjectUseCase implémentation
type subjectUseCase struct {
	subjectRepo repository.SubjectRepository
}

func NewSubjectUsecase(subjectRepo repository.SubjectRepository) SubjectUseCase {
	return &subjectUseCase{
		subjectRepo: subjectRepo,
	}
}

func (uc *subjectUseCase) GetAllBySchoolID(schoolID int) ([]*dto.SubjectInfo, error) {
	subjects, err := uc.subjectRepo.FindBySchoolID(schoolID)
	if err != nil {
		return nil, fmt.Errorf("[ERROR_GET_Subjects_BY_SCHOOL_ID]: %w", err)
	}

	var resp []*dto.SubjectInfo
	for _, c := range subjects {
		resp = append(resp, dto.SubjectInfoFromDomain(c))
	}

	return resp, nil
}
