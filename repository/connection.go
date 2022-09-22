package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=drone port=5432 sslmode=disable TimeZone=Africa/Cairo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

var FixturesDrones []Drone

func LoadFixtures(db *gorm.DB) error {
	FixturesDrones = []Drone{
		{State: "IDLE", BatteryCapacity: 100, Medications: []Medication{}},
		{State: "IDLE", BatteryCapacity: 100, Medications: []Medication{}},
		{State: "IDLE", BatteryCapacity: 100, Medications: []Medication{}},
		{State: "IDLE", BatteryCapacity: 100, Medications: []Medication{}},
		{State: "IDLE", BatteryCapacity: 100, Medications: []Medication{}},
		{State: "IDLE", BatteryCapacity: 100, Medications: []Medication{}},
	}

	if result := db.Create(&FixturesDrones); result.Error != nil {
		return result.Error
	}

	fixturesMedication := []Medication{
		{Name: "medication 1", Code: "123", Weight: 50},
		{Name: "medication 2", Code: "147", Weight: 100},
		{Name: "medication 3", Code: "159", Weight: 20},
	}
	if result := db.Create(&fixturesMedication); result.Error != nil {
		return result.Error
	}
	fmt.Println("Fixtures loaded ...")
	return nil
}
