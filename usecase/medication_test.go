package usecase

import (
	"testing"
)

func Test_medicationUsecase_RegisterMedication(t *testing.T) {
	type args struct {
		object MedicationObject
	}
	tests := []struct {
		name       string
		medication *medicationUsecase
		args       args
		want       int
		wantErr    bool
		errorMsg   string
	}{
		{
			name: "test can not register medication without name",
			args: args{
				object: MedicationObject{
					Code:   "code",
					Weight: 10,
				},
			},
			want:     0,
			wantErr:  true,
			errorMsg: "Medication name is not provided",
		},
		{
			name: "test can not register medication without code",
			args: args{
				object: MedicationObject{
					Name:   "test",
					Weight: 10,
				},
			},
			want:     0,
			wantErr:  true,
			errorMsg: "Medication code is not provided",
		},
		{
			name: "test can not register medication without weight",
			args: args{
				object: MedicationObject{
					Name: "test",
					Code: "code",
				},
			},
			want:     0,
			wantErr:  true,
			errorMsg: "Medication weight is not provided",
		},
		{
			name: "test can not register medication without mandatory data",
			args: args{
				object: MedicationObject{},
			},
			want:     0,
			wantErr:  true,
			errorMsg: "Medication code is not provided;Medication name is not provided;Medication weight is not provided",
		},
		{
			name: "test medication image field should be vaild url format",
			args: args{
				object: MedicationObject{
					Name:   "test",
					Code:   "code",
					Weight: 10,
					Image:  "invaild _url_format",
				},
			},
			want:     0,
			wantErr:  true,
			errorMsg: "image: invaild _url_format does not validate as url",
		},
		{
			name: "test register medication successfully",
			args: args{
				object: MedicationObject{
					Name:   "name",
					Code:   "code",
					Weight: 10,
					Image:  "http://test/image",
				},
			},
			want:     0,
			wantErr:  false,
			errorMsg: "",
		},
		{
			name: "test cant not register medication with weight less that 1",
			args: args{
				object: MedicationObject{
					Name:   "name",
					Code:   "code",
					Weight: -1,
					Image:  "http://test/image",
				},
			},
			want:     0,
			wantErr:  true,
			errorMsg: "weight: -1 does not validate as range(1|500)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			medication := &medicationUsecase{}
			got, err := medication.RegisterMedication(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("medicationUsecase.RegisterMedication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("medicationUsecase.RegisterMedication() = %v, want %v", got, tt.want)
			}
		})
	}
}
