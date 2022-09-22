package repository

import (
	settings "drone/v2/settings"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db, err = Init()
	if err != nil {
		fmt.Printf("error in setup test database: %v", err)
	}
	os.Exit(m.Run())
}

func Test_droneRepo_Create(t *testing.T) {
	type fields struct {
		client *gorm.DB
	}
	type args struct {
		drone *Drone
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		fixtures   []Drone
		want       int
		wantErr    bool
		wantObject Drone
		wantMsgExp string
	}{
		{
			name: "test create new drone",
			fields: fields{
				client: db,
			},
			args: args{
				drone: &Drone{
					SerialNumber:    "serial 1",
					Weight:          300,
					State:           "idle",
					Model:           "lightweight",
					BatteryCapacity: 100,
					CurrentPayload:  0,
				},
			},
			fixtures: []Drone{},
			want:     1,
			wantErr:  false,
			wantObject: Drone{
				SerialNumber:    "serial 1",
				Weight:          300,
				State:           settings.GetDroneState()["idle"],
				Model:           settings.GetDroneModels()["lightweight"],
				BatteryCapacity: 100,
				CurrentPayload:  0,
			},
			wantMsgExp: "",
		},
		{
			name: "test can not create new drone with dublicate serial number",
			fields: fields{
				client: db,
			},
			args: args{
				drone: &Drone{
					SerialNumber:    "serial 1",
					Weight:          300,
					State:           "idle",
					Model:           "lightweight",
					BatteryCapacity: 100,
					CurrentPayload:  0,
				},
			},
			fixtures: []Drone{
				{
					SerialNumber:    "serial 1",
					Weight:          300,
					State:           "idle",
					Model:           "lightweight",
					BatteryCapacity: 100,
					CurrentPayload:  0,
				},
			},
			want:       1,
			wantErr:    true,
			wantObject: Drone{},
			wantMsgExp: "duplicate key value violates unique constraint .*serial_number",
		},
		{
			name: "test can not create new drone with serial number more than 100 characters",
			fields: fields{
				client: db,
			},
			args: args{
				drone: &Drone{
					SerialNumber:    "6qjUThKzS4mdhvCCXh9QEH2tdhYxwbu3rPJYD8tQMQwS456hn4KyzDBh24VHDiFgbZkkMna49agPiydhN5eXTkvieRd9CXv7QrDnF",
					Weight:          300,
					State:           "idle",
					Model:           "lightweight",
					BatteryCapacity: 100,
					CurrentPayload:  0,
				},
			},
			fixtures:   []Drone{},
			want:       1,
			wantErr:    true,
			wantObject: Drone{},
			wantMsgExp: "value too long for type character varying.*",
		},
		// {
		// 	name: "test can not create new drone with weight more than 500",
		// 	fields: fields{
		// 		client: db,
		// 	},
		// 	args: args{
		// 		drone: &Drone{
		// 			SerialNumber:    "serial 1",
		// 			Weight:          501,
		// 			State:           "idle",
		// 			Model:           "lightweight",
		// 			BatteryCapacity: 100,
		// 			CurrentPayload:  0,
		// 		},
		// 	},
		// 	fixtures:   []Drone{},
		// 	want:       1,
		// 	wantErr:    true,
		// 	wantObject: Drone{},
		// 	wantMsgExp: "can not save drone with weight more that 500",
		// },
		// {
		// 	name: "test can not create new drone with model not exist",
		// 	fields: fields{
		// 		client: db,
		// 	},
		// 	args: args{
		// 		drone: &Drone{
		// 			SerialNumber:    "serial 1",
		// 			Weight:          230,
		// 			State:           "idle",
		// 			Model:           "test model",
		// 			BatteryCapacity: 100,
		// 			CurrentPayload:  0,
		// 		},
		// 	},
		// 	fixtures:   []Drone{},
		// 	want:       1,
		// 	wantErr:    true,
		// 	wantObject: Drone{},
		// 	wantMsgExp: "can not save drone with model not exist",
		// },
		// {
		// 	name: "test can not create new drone with state not exist",
		// 	fields: fields{
		// 		client: db,
		// 	},
		// 	args: args{
		// 		drone: &Drone{
		// 			SerialNumber:    "serial 1",
		// 			Weight:          400,
		// 			State:           "test state",
		// 			Model:           "lightweight",
		// 			BatteryCapacity: 100,
		// 			CurrentPayload:  0,
		// 		},
		// 	},
		// 	fixtures:   []Drone{},
		// 	want:       1,
		// 	wantErr:    true,
		// 	wantObject: Drone{},
		// 	wantMsgExp: "can not save drone with state not exist",
		// },
		// {
		// 	name: "test can not create new drone with negative battery level",
		// 	fields: fields{
		// 		client: db,
		// 	},
		// 	args: args{
		// 		drone: &Drone{
		// 			SerialNumber:    "serial 1",
		// 			Weight:          400,
		// 			State:           settings.GetDroneState()["idle"],
		// 			Model:           "lightweight",
		// 			BatteryCapacity: -10,
		// 			CurrentPayload:  0,
		// 		},
		// 	},
		// 	fixtures:   []Drone{},
		// 	want:       1,
		// 	wantErr:    true,
		// 	wantObject: Drone{},
		// 	wantMsgExp: "can not save drone with negative battery level",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear all record
			db.Where("1 = 1").Delete(&Drone{})
			trx := db.Begin()
			defer trx.Rollback()
			if len(tt.fixtures) > 0 {
				result := trx.Create(&tt.fixtures)
				if result.Error != nil {
					t.Errorf("Can't create fixtures: %v", result.Error)
				}
			}

			d := &droneRepo{
				client: trx,
			}
			_, err := d.Create(tt.args.drone)
			if (err != nil) != tt.wantErr {
				t.Errorf("droneRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r, _ := regexp.Compile(tt.wantMsgExp)
			if err != nil && !r.MatchString(err.Error()) {
				t.Errorf("errorMsg not matching = %v, wantErr %v", err.Error(), tt.wantMsgExp)
				return
			}
			if !tt.wantErr {
				var got Drone
				result := trx.Where("serial_number = ? ", tt.args.drone.SerialNumber).Find(&got)
				if result.Error != nil {
					t.Errorf("Can't retrive drone: %v", result.Error)
				}

				if got.ID == 0 {
					t.Errorf("expected drone has id not equal zero but got %v", got.ID)
				}
				tt.wantObject.ID = got.ID
				if !reflect.DeepEqual(got, tt.wantObject) {
					t.Errorf("expected drone = %v, want %v", tt.wantObject, got)
				}
			}

		})
	}
}

func Test_droneRepo_Get(t *testing.T) {
	type fields struct {
		client *gorm.DB
	}
	type args struct {
		id int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     Drone
		fixtures []Drone
		wantErr  bool
		wantMsg  string
	}{
		{
			name: "test get drone that not exist",
			fields: fields{
				client: db,
			},
			args: args{
				id: 1,
			},
			fixtures: []Drone{},
			want:     Drone{},
			wantErr:  true,
			wantMsg:  gorm.ErrRecordNotFound.Error(),
		},
		{
			name: "test get drone that exist",
			fields: fields{
				client: db,
			},
			args: args{
				id: 1,
			},
			fixtures: []Drone{
				{
					SerialNumber: "get serial 1",
					State:        "IDLE",
					Model:        "Lightweight",
				},
			},
			want: Drone{
				SerialNumber:    "get serial 1",
				State:           "IDLE",
				Model:           "Lightweight",
				BatteryCapacity: 100,
			},
			wantErr: false,
			wantMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear all record
			db.Where("1 = 1").Delete(&Drone{})
			trx := db.Begin()
			defer trx.Rollback()
			var createDrone Drone
			if len(tt.fixtures) > 0 {
				result := trx.Create(&tt.fixtures)
				if result.Error != nil {
					t.Errorf("Can't create fixtures: %v", result.Error)
				}
				trx.Where("serial_number = ?", tt.fixtures[0].SerialNumber).Find(&createDrone)
			}

			d := &droneRepo{
				client: trx,
			}

			got, err := d.Get(createDrone.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("droneRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want.ID = got.ID
			tt.want.Medications = got.Medications
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("droneRepo.Get() = %v, want %v", got, tt.want)
			}
			if err != nil && tt.wantMsg != err.Error() {
				t.Errorf("droneRepo.Get() error message = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_droneRepo_AddMedication(t *testing.T) {
	type fields struct {
		client *gorm.DB
	}
	type args struct {
		id         int
		medication *Medication
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantErr            bool
		fixtures           []Drone
		wantMedication     Medication
		medicatiionNUmbers int
	}{
		{
			name: "test add medication for drone that has empty mediactions",
			fields: fields{
				client: db,
			},
			args: args{
				id: 1,
				medication: &Medication{
					Name:   "medication",
					Code:   "code 1",
					Weight: 10,
				},
			},
			fixtures: []Drone{
				{
					SerialNumber: "ser 1",
					State:        "IDLE",
					Model:        "Lightweight",
				},
			},
			wantErr: false,
			wantMedication: Medication{
				Name:   "medication",
				Code:   "code 1",
				Weight: 10,
			},
			medicatiionNUmbers: 1,
		},
		{
			name: "test add medication for drone that has mediactions",
			fields: fields{
				client: db,
			},
			args: args{
				id: 1,
				medication: &Medication{
					Name:   "medication3",
					Code:   "code 3",
					Weight: 10,
				},
			},
			fixtures: []Drone{
				{
					SerialNumber: "ser 1",
					State:        "IDLE",
					Model:        "Lightweight",
					Medications: []Medication{
						{
							Name:   "medication1",
							Code:   "code 1",
							Weight: 10,
						},
						{
							Name:   "medication2",
							Code:   "code 2",
							Weight: 10,
						},
					},
				},
			},
			wantErr: false,
			wantMedication: Medication{
				Name:   "medication3",
				Code:   "code 3",
				Weight: 10,
			},
			medicatiionNUmbers: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear all record
			db.Where("1 = 1").Delete(&Drone{})
			trx := db.Begin()
			defer trx.Rollback()
			var createDrone Drone
			if len(tt.fixtures) > 0 {
				result := trx.Create(&tt.fixtures)
				if result.Error != nil {
					t.Errorf("Can't create fixtures: %v", result.Error)
				}
				trx.Where("serial_number = ?", tt.fixtures[0].SerialNumber).Find(&createDrone)
			}

			d := &droneRepo{
				client: trx,
			}

			if err := d.AddMedication(createDrone.ID, tt.args.medication); (err != nil) != tt.wantErr {
				t.Errorf("droneRepo.AddMedication() error = %v, wantErr %v", err, tt.wantErr)
			}

			trx.Where("serial_number = ?", tt.fixtures[0].SerialNumber).Preload("Medications").Find(&createDrone)
			if len(createDrone.Medications) != tt.medicatiionNUmbers {
				t.Errorf("droneRepo.AddMedication() lenght error = %v, wantErr %v", len(createDrone.Medications), tt.medicatiionNUmbers)
			}
			if !foundMedication(tt.wantMedication, createDrone.Medications) {
				t.Errorf("expected to found new medication %v in %v but not exist", tt.wantMedication, createDrone.Medications)

			}
		})
	}
}

func foundMedication(m Medication, medications []Medication) bool {
	for _, o := range medications {
		if o.Name == m.Name {
			return true
		}
	}
	return false
}

func Test_droneRepo_CheckLoadingMedication(t *testing.T) {
	type fields struct {
		client *gorm.DB
	}
	type args struct {
		id int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		fixtures []Drone
		want     string
		wantErr  bool
	}{
		{
			name: "test check loading medication",
			fields: fields{
				client: db,
			},
			args: args{
				id: 1,
			},
			fixtures: []Drone{
				{
					SerialNumber: "ser 1",
					State:        "IDLE",
					Model:        "Lightweight",
				},
			},
			want:    settings.GetDroneState()["idle"],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear all record
			db.Where("1 = 1").Delete(&Drone{})
			trx := db.Begin()
			defer trx.Rollback()
			var createDrone Drone
			if len(tt.fixtures) > 0 {
				result := trx.Create(&tt.fixtures)
				if result.Error != nil {
					t.Errorf("Can't create fixtures: %v", result.Error)
				}
				trx.Where("serial_number = ?", tt.fixtures[0].SerialNumber).Find(&createDrone)
			}

			d := &droneRepo{
				client: trx,
			}
			got, err := d.CheckLoadingMedication(createDrone.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("droneRepo.CheckLoadingMedication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("droneRepo.CheckLoadingMedication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_droneRepo_AvailableDroneForLoading(t *testing.T) {
	type fields struct {
		client *gorm.DB
	}
	tests := []struct {
		name     string
		fields   fields
		fixtures []Drone
		want     []Drone
	}{
		{
			name: "test get available drone for loading",
			fields: fields{
				client: db,
			},
			fixtures: []Drone{
				{
					SerialNumber: "ser 1",
					State:        settings.GetDroneState()["idle"],
					Model:        settings.GetDroneModels()["lightweight"],
				},
				{
					SerialNumber: "ser 2",
					State:        settings.GetDroneState()["loading"],
					Model:        settings.GetDroneModels()["lightweight"],
				},
				{
					SerialNumber: "ser 3",
					State:        settings.GetDroneState()["loaded"],
					Model:        settings.GetDroneModels()["lightweight"],
				},
				{
					SerialNumber: "ser 4",
					State:        settings.GetDroneState()["delivering"],
					Model:        settings.GetDroneModels()["lightweight"],
				},
				{
					SerialNumber: "ser 5",
					State:        settings.GetDroneState()["delivered"],
					Model:        settings.GetDroneModels()["lightweight"],
				},
				{
					SerialNumber: "ser 7",
					State:        settings.GetDroneState()["idle"],
					Model:        settings.GetDroneModels()["lightweight"],
				},
			},
			want: []Drone{
				{
					SerialNumber:    "ser 1",
					State:           settings.GetDroneState()["idle"],
					Model:           settings.GetDroneModels()["lightweight"],
					BatteryCapacity: 100,
				},
				{
					SerialNumber:    "ser 7",
					State:           settings.GetDroneState()["idle"],
					Model:           settings.GetDroneModels()["lightweight"],
					BatteryCapacity: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear all record
			db.Where("1 = 1").Delete(&Drone{})
			trx := db.Begin()
			defer trx.Rollback()
			var createDrone Drone
			if len(tt.fixtures) > 0 {
				result := trx.Create(&tt.fixtures)
				if result.Error != nil {
					t.Errorf("Can't create fixtures: %v", result.Error)
				}
				trx.Where("serial_number = ?", tt.fixtures[0].SerialNumber).Find(&createDrone)
			}

			d := &droneRepo{
				client: trx,
			}
			got := d.AvailableDroneForLoading()
			for i, _ := range got {
				tt.want[i].ID = got[i].ID
				tt.want[i].Medications = got[i].Medications
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected drone = %v, want %v", tt.want, got)
			}
		})
	}
}

func Test_droneRepo_CheckBatteryLevel(t *testing.T) {
	type fields struct {
		client *gorm.DB
	}
	type args struct {
		id int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		fixtures []Drone
		want     int
		wantErr  bool
	}{
		{
			name: "test check battery level for given new drone",
			fields: fields{
				client: db,
			},
			args: args{
				id: 1,
			},
			fixtures: []Drone{
				{
					SerialNumber: "ser 1",
					State:        settings.GetDroneState()["idle"],
					Model:        settings.GetDroneModels()["lightweight"],
				},
			},
			want:    100,
			wantErr: false,
		},
		{
			name: "test check battery level for given exist drone",
			fields: fields{
				client: db,
			},
			args: args{
				id: 1,
			},
			fixtures: []Drone{
				{
					SerialNumber:    "ser 1",
					State:           settings.GetDroneState()["idle"],
					Model:           settings.GetDroneModels()["lightweight"],
					BatteryCapacity: 50,
				},
			},
			want:    50,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear all record
			db.Where("1 = 1").Delete(&Drone{})
			trx := db.Begin()
			defer trx.Rollback()
			var createDrone Drone
			if len(tt.fixtures) > 0 {
				result := trx.Create(&tt.fixtures)
				if result.Error != nil {
					t.Errorf("Can't create fixtures: %v", result.Error)
				}
				trx.Where("serial_number = ?", tt.fixtures[0].SerialNumber).Find(&createDrone)
			}

			d := &droneRepo{
				client: trx,
			}
			got, err := d.CheckBatteryLevel(createDrone.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("droneRepo.CheckBatteryLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("droneRepo.CheckBatteryLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_droneRepo_ReduceBatteries(t *testing.T) {
	type fields struct {
		client *gorm.DB
	}
	tests := []struct {
		name          string
		fields        fields
		fixtures      []Drone
		wantBatteries []int
	}{
		{
			name: "test reduce batteries",
			fields: fields{
				client: db,
			},
			fixtures: []Drone{
				{
					SerialNumber: "ser 1",
					State:        settings.GetDroneState()["idle"],
					Model:        settings.GetDroneModels()["lightweight"],
				},
				{
					SerialNumber:    "ser 2",
					State:           settings.GetDroneState()["idle"],
					Model:           settings.GetDroneModels()["lightweight"],
					BatteryCapacity: 40,
				},
			},
			wantBatteries: []int{99, 39},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear all record
			db.Where("1 = 1").Delete(&Drone{})
			trx := db.Begin()
			defer trx.Rollback()
			if len(tt.fixtures) > 0 {
				result := trx.Create(&tt.fixtures)
				if result.Error != nil {
					t.Errorf("Can't create fixtures: %v", result.Error)
				}
			}

			d := &droneRepo{
				client: trx,
			}
			d.ReduceBatteries()
			var createDrones []Drone
			trx.Find(&createDrones)
			for i, _ := range createDrones {
				if tt.wantBatteries[i] != createDrones[i].BatteryCapacity {
					t.Errorf("ReduceBatteries = %v, want %v", tt.wantBatteries[i], createDrones[i].BatteryCapacity)
				}
			}
		})
	}
}
