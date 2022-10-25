package repository

import (
	"gorm.io/gorm"
)

type IDroneRepository interface {
	Create(drone *Drone) (int, error)
	Get(id int) (Drone, error)
	AddMedication(id int, medication *Medication) error
	CheckLoadingMedication(id int) (string, error)
	AvailableDroneForLoading() []Drone
	CheckBatteryLevel(id int) (int, error)
	ReduceBatteries()
	// chnageDroneStatus(id int, state string) error
}

type droneRepo struct {
	client  *gorm.DB
	logRepo ILogRepository
}

func NewDroneRepo(client *gorm.DB, logRepo ILogRepository) IDroneRepository {
	return &droneRepo{
		client:  client,
		logRepo: logRepo,
	}
}

func (d *droneRepo) Create(drone *Drone) (int, error) {

	// Not Add this validation here becuase this layer responsible for saving data only
	// if drone.Weight > 500 {
	// 	return 0, errors.New("can not save drone with weight more that 500")
	// }
	// if drone.BatteryCapacity < 0 {
	// 	return 0, errors.New("can not save drone with negative battery level")
	// }
	// if _, found := settings.GetDroneModels()[drone.Model]; !found {
	// 	return 0, errors.New("can not save drone with model not exist")
	// }
	// if _, found := settings.GetDroneState()[drone.State]; !found {
	// 	return 0, errors.New("can not save drone with state not exist")
	// }
	// drone.State = settings.GetDroneState()[drone.State]
	// drone.Model = settings.GetDroneModels()[drone.Model]
	result := d.client.Save(&drone)
	if result.Error != nil {
		return 0, result.Error
	}
	return drone.ID, nil
}

func (d *droneRepo) Get(id int) (Drone, error) {
	var drone Drone
	if result := d.client.Preload("Medications").First(&drone, id); result.Error != nil {
		return Drone{}, result.Error
	}
	return drone, nil
}

func (d *droneRepo) AddMedication(id int, medication *Medication) error {
	drone, err := d.Get(id)
	if err != nil {
		return err
	}
	medication.DroneID = drone.ID
	drone.Medications = append(drone.Medications, *medication)
	drone.State = "Loading"
	if result := d.client.Save(&drone); result.Error != nil {
		return result.Error
	}
	// d.chnageDroneStatus(id, settings.GetDroneState()["loading"])
	return nil
}

func (d *droneRepo) CheckLoadingMedication(id int) (string, error) {
	drone, err := d.Get(id)
	if err != nil {
		return "", err
	}
	return drone.State, nil
}

func (d *droneRepo) AvailableDroneForLoading() []Drone {
	var availableDrone []Drone
	result := d.client.Where("state = ?", "IDLE").Preload("Medications").Find(&availableDrone)
	if result.Error != nil {
		return []Drone{}
	}
	return availableDrone
}

func (d *droneRepo) CheckBatteryLevel(id int) (int, error) {
	var drone Drone
	result := d.client.Where("id = ?", id).Find(&drone)
	if result.Error != nil {
		return 0, result.Error
	}
	return drone.BatteryCapacity, nil
}

// func (d *droneRepo) chnageDroneStatus(id int, state string) error {
// 	result := d.client.Where("id = ?", id).Update("state", state)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

func (d *droneRepo) ReduceBatteries() {
	//TODO: refactor can do this logic using ORM
	var drones []Drone
	d.client.Find(&drones)
	drones = d.reduceBatteries(drones)
	d.client.Save(&drones)
}

func (d *droneRepo) reduceBatteries(drones []Drone) []Drone {
	var update []Drone
	for _, o := range drones {
		if o.BatteryCapacity > 1 {
			o.BatteryCapacity = o.BatteryCapacity - 1
			update = append(update, o)
			d.logRepo.Create(Log{
				DroneID:         o.ID,
				DroneState:      o.State,
				BatteryCapacity: o.BatteryCapacity,
			})
		}
	}
	return update
}
