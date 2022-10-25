package usecase

import (
	repo "drone/v2/repository"
	repoEnity "drone/v2/repository"
	"drone/v2/utils"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/asaskevich/govalidator"
)

type IDroneUsecase interface {
	RegisterDrone(object DorneObject) (int, error)
	LoadingMedication(id int, medication MedicationObject) error
	CheckLoadingMedication(id int) (string, error)
	CheckAvailableDroneForLoading() []repo.Drone
	CheckBatteryLevel(id int) (string, error)
	CheckDronesBatteries()
}

type droneUsecase struct {
	droneRepo repo.IDroneRepository
}

func NewDroneUsecase(d repo.IDroneRepository) IDroneUsecase {
	return &droneUsecase{
		droneRepo: d,
	}
}

func (d *droneUsecase) RegisterDrone(object DorneObject) (int, error) {
	droneValidate, err := govalidator.ValidateStruct(object)
	if err != nil || !droneValidate {
		return 0, err
	}
	data, err := utils.TypeConverter[repo.Drone](&object)
	if err != nil {
		log.Println(err.Error())
	}
	return d.droneRepo.Create(data)
}

func (d *droneUsecase) LoadingMedication(id int, medication MedicationObject) error {
	drone, err := d.droneRepo.Get(id)
	if err != nil {
		return err
	}
	err = utils.ValidateMedicationName(medication.Name)
	if err != nil {
		return err
	}
	medicationValidate, err := govalidator.ValidateStruct(medication)
	if err != nil || !medicationValidate {
		return err
	}
	if err := validateDroneForLoadingMedication(drone, medication.Weight); err == nil {
		data, err := utils.TypeConverter[repo.Medication](&medication)
		if err != nil {
			log.Println(err.Error())
		}
		// simulation Loading item time based on medication weight into drone
		// Let now be Fixed time
		time.Sleep(5 * time.Second)
		return d.droneRepo.AddMedication(id, data)
	}
	return err
}

func (d *droneUsecase) CheckLoadingMedication(id int) (string, error) {
	return d.droneRepo.CheckLoadingMedication(id)
}

func (d *droneUsecase) CheckAvailableDroneForLoading() []repo.Drone {
	return d.droneRepo.AvailableDroneForLoading()
}

func (d *droneUsecase) CheckBatteryLevel(id int) (string, error) {
	batteryLevel, err := d.droneRepo.CheckBatteryLevel(id)
	if err != nil {
		return "", err
	}
	batteryFormat := fmt.Sprintf("%d%s", batteryLevel, "%")
	return batteryFormat, nil
}

func validateDroneForLoadingMedication(drone repoEnity.Drone, weight float32) error {
	if drone.BatteryCapacity < 25 {
		errorMsg := fmt.Sprintf(`drone can not be loaded because battery capacity less that %d`, drone.BatteryCapacity)
		return errors.New(errorMsg)
	}
	if drone.CurrentPayload+weight > drone.Weight {
		errorMsg := fmt.Sprintf(`drone can not be loaded with %f weight, because current weight is %f and Max weight is %f`, weight, drone.CurrentPayload, drone.Weight)
		return errors.New(errorMsg)
	}
	return nil
}

func (d *droneUsecase) CheckDronesBatteries() {
	d.droneRepo.ReduceBatteries()

}
