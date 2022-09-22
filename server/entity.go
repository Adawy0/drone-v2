package server

type DornePayload struct {
	SerialNumber string  `json:"serial_number" valid:"required~Serial Number is not provided,stringlength(10|100)"`
	Model        string  `json:"model" valid:"required~Model is not provided,matches(Lightweight|Middleweight|Cruiserweight|Heavyweight)"`
	Weight       float32 `json:"weight" valid:"required~Weight is not provided,range(10|500)"`
	Battery      int     `json:"battery" valid:"optional, range(10|100)"`
	State        string  `json:"state" valid:"optional,matches(IDLE|LOADING|LOADED|DELIVERING|DELIVERED|RETURNING)"`
}

type MedicationPayload struct {
	Name   string  `json:"name" valid:"required~Medication name is not provided"`
	Code   string  `json:"code" valid:"required~Medication code is not provided"`
	Weight float32 `json:"weight" valid:"required~Medication weight is not provided,range(1|500)"`
	Image  string  `json:"image" valid:"optional,url"`
}

type RegisterDronePayload struct {
	DroneId int `json:"drone_id"`
}

type RegisterMediactionPayload struct {
	MedicationId int `json:"medication_id"`
}

type CheckLoadingMedicationPayload struct {
	DroneId int    `json:"drone_id"`
	Status  string `json:"status"`
}

type BatteryLevelPayload struct {
	DroneId      int    `json:"drone_id"`
	BatteryLevel string `json:"battery_level"`
}
