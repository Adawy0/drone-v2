package main

import (
	entity "drone/v2/repository"
	"time"

	"gorm.io/gorm"
)

// Up is executed when this migration is applied
func Up_20220808232129(txn *gorm.DB) {
	type Drone struct {
		ID              int     `json:"id" gorm:"primaryKey"`
		SerialNumber    string  `json:"serial_number" gorm:"type:varchar(100);uniqueIndex"`
		Weight          float32 `json:"weight" gorm:"size:500"`
		State           string  `json:"state" gorm:"default:IDLE"`
		Model           string  `json:"model"`
		BatteryCapacity int     `json:"battery_capactiy" gorm:"default:100"`
		Medications     []entity.Medication
		CurrentPayload  float32 `json:"current_payload" gorm:"default:0"`
	}
	txn.AutoMigrate(&Drone{})

	type Medication struct {
		Name    string `json:"name" gorm:"uniqueIndex"`
		Code    string `json:"code" gorm:"primaryKey"`
		Weight  int    `json:"weight"`
		DroneID int    `gorm:"foreignKey:DroneID"`
		Image   []byte `json:"image"`
	}
	txn.AutoMigrate(&Medication{})

	type Log struct {
		ID              int            `json:"-" gorm:"primaryKey"`
		CreatedAt       time.Time      `json:"date"`
		UpdatedAt       time.Time      `json:"-"`
		DeletedAt       gorm.DeletedAt `json:"-"`
		DroneID         int
		BatteryCapacity int
		DroneState      string
	}
	txn.AutoMigrate(&Log{})
}

// Down is executed when this migration is rolled back
func Down_20220808232129(txn *gorm.DB) {
	txn.Migrator().DropTable("drone")
	txn.Migrator().DropTable("medication")
	txn.Migrator().DropTable("logs")

}
