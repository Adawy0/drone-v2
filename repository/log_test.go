package repository

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestLogDB_List(t *testing.T) {
	commonFixtures := []Log{
		{
			DroneID:         1,
			BatteryCapacity: 50,
			DroneState:      "LOADING",
		},
		{
			DroneID:         2,
			BatteryCapacity: 100,
			DroneState:      "IDLE",
		},
		{
			DroneID:         3,
			BatteryCapacity: 90,
			DroneState:      "LOADED",
		},
	}
	type fields struct {
		client *gorm.DB
	}
	tests := []struct {
		name     string
		fields   fields
		fixtures []Log
		want     []Log
		wantErr  bool
	}{
		{
			name: "test list all logs",
			fields: fields{
				client: db,
			},
			fixtures: commonFixtures,
			want: []Log{
				{
					ID:              1,
					DroneID:         1,
					BatteryCapacity: 50,
					DroneState:      "LOADING",
				},
				{
					ID:              2,
					DroneID:         2,
					BatteryCapacity: 100,
					DroneState:      "IDLE",
				},
				{
					ID:              3,
					DroneID:         3,
					BatteryCapacity: 90,
					DroneState:      "LOADED",
				},
			},
			wantErr: false,
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
			ldb := LogDB{
				client: trx,
			}
			got, err := ldb.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("LogDB.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, _ := range got {
				tt.want[i].ID = got[i].ID
				tt.want[i].CreatedAt = got[i].CreatedAt
				tt.want[i].UpdatedAt = got[i].UpdatedAt
				tt.want[i].DeletedAt = got[i].DeletedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogDB.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogDB_Create(t *testing.T) {
	type fields struct {
		client *gorm.DB
	}
	type args struct {
		log Log
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		records int
	}{
		{
			name: "create log success",
			fields: fields{
				client: db,
			},
			args: args{
				log: Log{
					DroneID:         2,
					BatteryCapacity: 100,
					DroneState:      "IDLE",
				},
			},
			wantErr: false,
			records: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trx := tt.fields.client.Begin()
			defer trx.Rollback()
			ldb := LogDB{
				client: trx,
			}
			if err := ldb.Create(tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("LogDB.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			var got []Log
			ldb.client.Find(&got)
			if tt.records != len(got) {
				t.Errorf("LogDB.Create() count error want: %v got %v", tt.records, len(got))
			}
		})
	}
}
