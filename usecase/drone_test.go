package usecase

import (
	"bytes"
	repo "drone/v2/repository"
	repoEnity "drone/v2/repository"
	mosks "drone/v2/repository/mocks"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func Test_droneUsecase_Register(t *testing.T) {
	type args struct {
		object DorneObject
	}
	tests := []struct {
		name        string
		d           *droneUsecase
		args        args
		want        int
		wantErr     bool
		ErrorMsgExp string
	}{
		{
			name: "test cant not register drone without serial number",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					Model:   "Lightweight",
					Weight:  100,
					State:   "IDLE",
					Battery: 100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "Serial Number is not provided",
		},
		{
			name: "test can not register drone with serial number larger than 100 characters",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(101),
					Model:        "Lightweight",
					Weight:       100,
					State:        "IDLE",
					Battery:      100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: `^serial_number:.*does not validate as stringlength.*$`,
		},
		{
			name: "test can not register drone with serial number less than 10 characters",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(5),
					Model:        "Lightweight",
					Weight:       100,
					State:        "IDLE",
					Battery:      100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: `^serial_number:.*does not validate as stringlength.*$`,
		},
		{
			name: "test cant not register drone without model",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(50),
					Weight:       100,
					State:        "IDLE",
					Battery:      100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "Model is not provided",
		},
		{
			name: "test cant not register drone without weight",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(80),
					Model:        "Lightweight",
					State:        "IDLE",
					Battery:      100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "Weight is not provided",
		},
		{
			name: "test cant not register drone without mandatory data",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					State:   "IDLE",
					Battery: 100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "Model is not provided;Serial Number is not provided;Weight is not provided",
		},
		{
			name: "test cant not register drone with model that not exist",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(50),
					Model:        "Model Not exist",
					State:        "IDLE",
					Weight:       250,
					Battery:      100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "model: Model Not exist does not validate as matches.*",
		},
		{
			name: "test cant not register drone with weight larger than 500",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(50),
					Model:        "Lightweight",
					State:        "IDLE",
					Battery:      100,
					Weight:       501,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "weight: [0-9]* does not validate as range.*",
		},
		{
			name: "test cant not register drone with weight less than 10",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(50),
					Model:        "Lightweight",
					State:        "IDLE",
					Battery:      100,
					Weight:       9,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "weight: [0-9]* does not validate as range.*",
		},
		{
			name: "test cant not register drone with state that not exist",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(50),
					Model:        "Lightweight",
					State:        "State not exist",
					Weight:       250,
					Battery:      100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "state: State not exist does not validate as matches.*",
		},
		{
			name: "test cant not register drone with negative battery charge",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(50),
					Model:        "Lightweight",
					State:        "IDLE",
					Weight:       250,
					Battery:      -100,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "battery: -[0-9]* does not validate as range.*",
		},
		{
			name: "test cant not register drone with battery charge more than 100",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(50),
					Model:        "Lightweight",
					State:        "IDLE",
					Weight:       250,
					Battery:      101,
				},
			},
			want:        0,
			wantErr:     true,
			ErrorMsgExp: "battery: [0-9]* does not validate as range.*",
		},
		{
			name: "test register drone successfully",
			d:    &droneUsecase{droneRepo: mosks.NewDroneRepoMock()},
			args: args{
				object: DorneObject{
					SerialNumber: generateRandomSerialNumber(50),
					Model:        "Lightweight",
					State:        "IDLE",
					Weight:       250,
					Battery:      100,
				},
			},
			want:        1,
			wantErr:     false,
			ErrorMsgExp: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.RegisterDrone(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("Error in register operation error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Error in result = %v, want %v", got, tt.want)
			}
			r, _ := regexp.Compile(tt.ErrorMsgExp)
			if err != nil && !r.MatchString(err.Error()) {
				t.Errorf("errorMsg not matching = %v, wantErr %v", err.Error(), tt.ErrorMsgExp)
				return
			}
		})
	}
}

func generateRandomSerialNumber(length int) string {
	rand.Seed(time.Now().UnixNano())
	var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	ran_serial := make([]rune, length)
	for i := range ran_serial {
		ran_serial[i] = charset[rand.Intn(len(charset))]
	}
	return string(ran_serial)
}

func getModelsAsstring(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func Test_droneUsecase_LoadingMedication(t *testing.T) {
	type args struct {
		id         int
		medication MedicationObject
	}
	tests := []struct {
		name     string
		d        *droneUsecase
		args     args
		wantErr  bool
		errorMsg string
	}{
		{
			name: "test can not register medication without name",
			d: &droneUsecase{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
				medication: MedicationObject{
					Code:   "code",
					Weight: 10,
				},
			},
			wantErr:  true,
			errorMsg: "Medication name is not provided",
		},
		{
			name: "test can not register medication with invaild name format",
			d: &droneUsecase{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
				medication: MedicationObject{
					Name:   "###test",
					Code:   "code",
					Weight: 10,
				},
			},
			wantErr:  true,
			errorMsg: "invald format for medciation name, can be only combination of (letters, numbers, -, _)",
		},
		{
			name: "test can not register medication without code",
			d: &droneUsecase{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
				medication: MedicationObject{
					Name:   "test",
					Weight: 10,
				},
			},
			wantErr:  true,
			errorMsg: "Medication code is not provided",
		},
		{
			name: "test can not register medication without weight",
			d: &droneUsecase{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
				medication: MedicationObject{
					Name: "test",
					Code: "code",
				},
			},
			wantErr:  true,
			errorMsg: "Medication weight is not provided",
		},
		{
			name: "test can not register medication without mandatory data",
			d: &droneUsecase{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id:         1,
				medication: MedicationObject{},
			},
			wantErr:  true,
			errorMsg: "Medication code is not provided;Medication name is not provided;Medication weight is not provided",
		},
		{
			name: "test medication image field should be vaild url format",
			d: &droneUsecase{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
				medication: MedicationObject{
					Name:   "test",
					Code:   "code",
					Weight: 10,
					Image:  "invaild _url_format",
				},
			},
			wantErr:  true,
			errorMsg: "image: invaild _url_format does not validate as url",
		},
		{
			name: "test register medication successfully",
			d: &droneUsecase{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
				medication: MedicationObject{
					Name:   "name123_jhh-",
					Code:   "code",
					Weight: 10,
					Image:  "http://test/image",
				},
			},
			wantErr:  false,
			errorMsg: "",
		},
		{
			name: "test cant not register medication with weight less that 1",
			d: &droneUsecase{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
				medication: MedicationObject{
					Name:   "name",
					Code:   "code",
					Weight: -1,
					Image:  "http://test/image",
				},
			},
			wantErr:  true,
			errorMsg: "weight: -1 does not validate as range(1|500)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// d := &droneUsecase{}
			if err := tt.d.LoadingMedication(tt.args.id, tt.args.medication); (err != nil) != tt.wantErr {
				t.Errorf("droneUsecase.LoadingMedication() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.errorMsg != err.Error() {
				t.Errorf("droneUsecase.LoadingMedication() error = %v, wantErr %v", err.Error(), tt.errorMsg)
			}
		})
	}
}

func Test_validateDroneForLoadingMedication(t *testing.T) {
	type args struct {
		drone  repoEnity.Drone
		weight float32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantMsg string
	}{
		{
			name: "test can not add mediaction that has weight more that drone weight",
			args: args{
				drone: repoEnity.Drone{
					Weight:          300,
					BatteryCapacity: 100,
					CurrentPayload:  0,
				},
				weight: 400,
			},
			wantErr: true,
			wantMsg: fmt.Sprintf(`drone can not be loaded with %f weight, because current weight is %f and Max weight is %f`, 400.000000, 0.000000, 300.000000),
		},
		{
			name: "test can not add mediaction to drone that bettary level less that 25",
			args: args{
				drone: repoEnity.Drone{
					Weight:          500,
					BatteryCapacity: 24,
					CurrentPayload:  100,
				},
				weight: 50,
			},
			wantErr: true,
			wantMsg: fmt.Sprintf(`drone can not be loaded because battery capacity less that %d`, 24),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDroneForLoadingMedication(tt.args.drone, tt.args.weight); (err != nil) != tt.wantErr {
				t.Errorf("validateDroneForLoadingMedication() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && err.Error() != tt.wantMsg {
				t.Errorf("validateDroneForLoadingMedication() error = %v, wantErr %v", err.Error(), tt.wantMsg)
			}
		})
	}
}

func Test_droneUsecase_CheckLoadingMedication(t *testing.T) {
	type fields struct {
		droneRepo repo.IDroneRepository
	}
	type args struct {
		id int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     string
		errorMsg string
	}{
		{
			name: "test check medication for drone that not exist",
			fields: fields{
				droneRepo: mosks.NewDroneRepoFailMock(),
			},
			args: args{
				id: 1,
			},
			want:     "",
			errorMsg: "can not found drone for this id 1",
		},
		{
			name: "test check medication for drone without before adding mediaction",
			fields: fields{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
			},
			want:     "IDLE",
			errorMsg: "",
		},
		{
			name: "test check medication for drone during adding mediaction",
			fields: fields{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 2,
			},
			want:     "LOADING",
			errorMsg: "",
		},
		{
			name: "test check medication for drone after mediaction loaded",
			fields: fields{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 3,
			},
			want:     "LOADED",
			errorMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &droneUsecase{
				droneRepo: tt.fields.droneRepo,
			}
			if got, err := d.CheckLoadingMedication(tt.args.id); got != tt.want {
				t.Errorf("droneUsecase.CheckLoadingMedication() = %v, want %v", got, tt.want)
			} else if err != nil && err.Error() != tt.errorMsg {
				t.Errorf("droneUsecase.CheckLoadingMedication() = %v, want %v", got, tt.want)

			}
		})
	}
}

func Test_droneUsecase_CheckAvailableDroneForLoading(t *testing.T) {
	type fields struct {
		droneRepo repo.IDroneRepository
	}
	tests := []struct {
		name   string
		fields fields
		want   []repo.Drone
	}{
		{
			name: "test not available drone for loading",
			fields: fields{
				droneRepo: mosks.NewDroneRepoFailMock(),
			},
			want: []repoEnity.Drone{},
		},
		{
			name: "test available drone for loading",
			fields: fields{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			want: []repoEnity.Drone{
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &droneUsecase{
				droneRepo: tt.fields.droneRepo,
			}
			if got := d.CheckAvailableDroneForLoading(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("droneUsecase.CheckAvailableDroneForLoading() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_droneUsecase_CheckBatteryLevel(t *testing.T) {
	type fields struct {
		droneRepo repo.IDroneRepository
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		wantMsg string
	}{
		{
			name: "test check battery for drone that not exist",
			fields: fields{
				droneRepo: mosks.NewDroneRepoFailMock(),
			},
			args: args{
				id: 1,
			},
			want:    "",
			wantErr: true,
			wantMsg: "can not found drone for this id 1",
		},
		{
			name: "test check battery for drone that exist",
			fields: fields{
				droneRepo: mosks.NewDroneRepoMock(),
			},
			args: args{
				id: 1,
			},
			want:    "25%",
			wantErr: false,
			wantMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &droneUsecase{
				droneRepo: tt.fields.droneRepo,
			}
			got, err := d.CheckBatteryLevel(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("droneUsecase.CheckBatteryLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("droneUsecase.CheckBatteryLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
