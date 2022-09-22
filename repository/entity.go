package repository

type Medication struct {
	Name    string `json:"name" gorm:"uniqueIndex"`
	Code    string `json:"code" gorm:"primaryKey"`
	Weight  int    `json:"weight"`
	Image   []byte `json:"image"`
	DroneID int    `gorm:"foreignKey:DroneID"`
}

func (Medication) TableName() string {
	return `"drone"."medications"`
}

type Drone struct {
	ID              int     `json:"id" gorm:"primaryKey"`
	SerialNumber    string  `json:"serial_number" gorm:"type:varchar(100);uniqueIndex"`
	Weight          float32 `json:"weight"`
	State           string  `json:"state" gorm:"default:IDLE"`
	Model           string  `json:"model"`
	BatteryCapacity int     `json:"battery_capactiy" gorm:"default:100"`
	Medications     []Medication
	CurrentPayload  float32 `json:"current_payload" gorm:"default:0"`
}

func (Drone) TableName() string {
	return `"drone"."drones"`
}
