package mocks

import (
	"drone/v2/usecase"
)

type IMedicationMockUsecase interface {
	RegisterMedication(object usecase.MedicationObject) (int, error)
}

type medicationMockUsecase struct {
}

func NewMedicationMockUsecase() IMedicationMockUsecase {
	return &medicationMockUsecase{}
}

func (medication *medicationMockUsecase) RegisterMedication(object usecase.MedicationObject) (int, error) {
	return 0, nil
}
