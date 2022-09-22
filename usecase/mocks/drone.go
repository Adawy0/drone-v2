package mocks

import (
	repo "drone/v2/repository"
	"drone/v2/usecase"
	"errors"
)

type IDroneMockUsecase interface {
	RegisterDrone(object usecase.DorneObject) (int, error)
	LoadingMedication(id int, medication usecase.MedicationObject) error
	CheckLoadingMedication(id int) (string, error)
	CheckAvailableDroneForLoading() []repo.Drone
	CheckBatteryLevel(id int) (string, error)
	CheckDronesBatteries()
}

type droneMockUsecase struct {
}

func NewDroneMockUsecase() IDroneMockUsecase {
	return &droneMockUsecase{}
}

func (u droneMockUsecase) RegisterDrone(object usecase.DorneObject) (int, error) {
	return 1, nil
}

func (u droneMockUsecase) LoadingMedication(id int, medication usecase.MedicationObject) error {
	return nil
}

func (u droneMockUsecase) CheckLoadingMedication(id int) (string, error) {
	return "", errors.New("")
}

func (u droneMockUsecase) CheckAvailableDroneForLoading() []repo.Drone {
	return []repo.Drone{}
}

func (u droneMockUsecase) CheckBatteryLevel(id int) (string, error) {
	return "", errors.New("")
}

func (u droneMockUsecase) CheckDronesBatteries() {

}
