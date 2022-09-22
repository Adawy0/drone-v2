package mocks

import (
	repo "drone/v2/repository"
	"errors"
	"fmt"
)

type droneRepoMock struct {
}

func NewDroneRepoMock() repo.IDroneRepository {
	return &droneRepoMock{}
}

func (d *droneRepoMock) Create(drone *repo.Drone) (int, error) {
	return 1, nil
}

func (d *droneRepoMock) Get(id int) (repo.Drone, error) {
	return repo.Drone{}, nil
}

func (d *droneRepoMock) AddMedication(id int, medication *repo.Medication) error {
	return nil
}

var state = []string{"IDLE", "LOADING", "LOADED"}

func (d *droneRepoMock) CheckLoadingMedication(id int) (string, error) {
	return state[id-1], nil
}

func (d *droneRepoMock) AvailableDroneForLoading() []repo.Drone {
	return []repo.Drone{
		{
			ID:           1,
			SerialNumber: "test serial 1",
			Weight:       120,
		},
		{
			ID:           2,
			SerialNumber: "test serial 2",
			Weight:       120,
		},
	}
}
func (d *droneRepoMock) CheckBatteryLevel(id int) (int, error) {
	return 25, nil
}
func (d *droneRepoMock) ReduceBatteries() {
	// reduce all drones batteries
}

// func (d *droneRepoMock) chnageDroneStatus(id int, state string) error {
// 	return nil
// }

type droneRepoFailMock struct {
}

func NewDroneRepoFailMock() repo.IDroneRepository {
	return &droneRepoFailMock{}
}

func (d *droneRepoFailMock) Create(drone *repo.Drone) (int, error) {
	return 1, nil
}

func (d *droneRepoFailMock) Get(id int) (repo.Drone, error) {
	return repo.Drone{}, nil
}

func (d *droneRepoFailMock) AddMedication(id int, medication *repo.Medication) error {
	return nil
}

func (d *droneRepoFailMock) CheckLoadingMedication(id int) (string, error) {
	return "", errors.New(fmt.Sprintf("can not found drone for this id %d", id))
}

func (d *droneRepoFailMock) AvailableDroneForLoading() []repo.Drone {
	return []repo.Drone{}
}

func (d *droneRepoFailMock) CheckBatteryLevel(id int) (int, error) {
	return 0.0, errors.New(fmt.Sprintf("can not found drone for this id %d", id))
}

func (d *droneRepoFailMock) ReduceBatteries() {
	// reduce all drones batteries
}

// func (d *droneRepoFailMock) chnageDroneStatus(id int, state string) error {
// 	return nil
// }
