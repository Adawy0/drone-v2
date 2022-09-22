package usecase

import (
	"drone/v2/utils"

	"github.com/asaskevich/govalidator"
)

type IMedicationUsecase interface {
	RegisterMedication(object MedicationObject) (int, error)
}

type medicationUsecase struct {
}

func NewMedicationUsecase() IMedicationUsecase {
	return &medicationUsecase{}
}

func (medication *medicationUsecase) RegisterMedication(object MedicationObject) (int, error) {
	err := utils.ValidateMedicationName(object.Name)
	if err != nil {
		return 0, err
	}
	droneValidate, err := govalidator.ValidateStruct(object)
	if err != nil || !droneValidate {
		return 0, err
	}
	// create medication
	return 0, nil
}
