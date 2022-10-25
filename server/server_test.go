package server

import (
	"drone/v2/usecase"
	mockUsecase "drone/v2/usecase/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func Test_droneAPI_RegisterDrone(t *testing.T) {
	type fields struct {
		droneUsecase      usecase.IDroneUsecase
		medicationUsecase usecase.IMedicationUsecase
	}
	type args struct {
		payload io.Reader
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       string
	}{
		{
			name: "Test can not register drone with invaild payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				payload: strings.NewReader(`{,}`),
			},
			wantStatus: http.StatusBadRequest,
			want:       "Invaild json payload\n",
		},
		{
			name: "Test can not register drone without payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				payload: nil,
			},
			wantStatus: http.StatusBadRequest,
			want:       "register drone must have json payload\n",
		},
		{
			name: "Test register drone accept json payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				payload: strings.NewReader(`{}`),
			},
			wantStatus: http.StatusCreated,
			want:       "",
		},
		{
			name: "Test can not register drone with empty body",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				payload: strings.NewReader(``),
			},
			wantStatus: http.StatusBadRequest,
			want:       "Invaild json payload\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &droneAPI{
				droneUsecase:      tt.fields.droneUsecase,
				medicationUsecase: tt.fields.medicationUsecase,
			}
			request, _ := http.NewRequest(http.MethodPost, "/api/drone", tt.args.payload)
			response := httptest.NewRecorder()
			api.RegisterDrone(response, request)
			if status := response.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusCreated)
			}
			if tt.wantStatus == http.StatusBadRequest {
				got := response.Body.String()
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			}

		})
	}
}

func Test_droneAPI_RegisterMedication(t *testing.T) {
	type fields struct {
		droneUsecase      usecase.IDroneUsecase
		medicationUsecase usecase.IMedicationUsecase
	}
	type args struct {
		payload io.Reader
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       string
	}{
		{
			name: "Test can not register medication with invaild payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				payload: strings.NewReader(`{,}`),
			},
			wantStatus: http.StatusBadRequest,
			want:       "invalid character ',' looking for beginning of object key string\n",
		},
		{
			name: "Test can not register medication without payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				payload: nil,
			},
			wantStatus: http.StatusBadRequest,
			want:       "load medication end point must have json payload\n",
		},
		{
			name: "Test register medication accept json payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				payload: strings.NewReader(`{}`),
			},
			wantStatus: http.StatusCreated,
			want:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &droneAPI{
				droneUsecase:      tt.fields.droneUsecase,
				medicationUsecase: tt.fields.medicationUsecase,
			}
			request, _ := http.NewRequest(http.MethodPost, "/api/medication", tt.args.payload)
			response := httptest.NewRecorder()
			api.RegisterMedication(response, request)
		})
	}
}

func Test_droneAPI_LoadingMedication(t *testing.T) {
	type fields struct {
		droneUsecase      usecase.IDroneUsecase
		medicationUsecase usecase.IMedicationUsecase
	}
	type args struct {
		id      string
		payload io.Reader
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       string
	}{
		{
			name: "Test can not load medication without drone id",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				id:      "",
				payload: strings.NewReader(`{}`),
			},
			wantStatus: http.StatusBadRequest,
			want:       "Couldn't find id in request URL\n",
		},
		{
			name: "Test can not load medication with drone id invaild",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				id:      "a",
				payload: strings.NewReader(`{}`),
			},
			wantStatus: http.StatusBadRequest,
			want:       "Invaild drone id\n",
		},
		{
			name: "Test can not loading medication with invaild payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				id:      "4",
				payload: strings.NewReader(`{,}`),
			},
			wantStatus: http.StatusBadRequest,
			want:       "invalid character ',' looking for beginning of object key string\n",
		},
		{
			name: "Test can not load medication without payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				id:      "4",
				payload: nil,
			},
			wantStatus: http.StatusBadRequest,
			want:       "load medication end point must have json payload\n",
		},
		{
			name: "Test load medication accept json payload",
			fields: fields{
				droneUsecase:      mockUsecase.NewDroneMockUsecase(),
				medicationUsecase: mockUsecase.NewMedicationMockUsecase(),
			},
			args: args{
				id:      "4",
				payload: strings.NewReader(`{}`),
			},
			wantStatus: http.StatusCreated,
			want:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &droneAPI{
				droneUsecase:      tt.fields.droneUsecase,
				medicationUsecase: tt.fields.medicationUsecase,
			}
			request, _ := http.NewRequest(http.MethodPost, "/api/drone//load", tt.args.payload)
			response := httptest.NewRecorder()
			if tt.args.id != "" {
				req := mux.SetURLVars(request, map[string]string{
					"id": tt.args.id,
				})
				api.LoadingMedication(response, req)
			} else {
				api.LoadingMedication(response, request)
			}

			if status := response.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusCreated)
			}
			if tt.wantStatus == http.StatusBadRequest {
				got := response.Body.String()
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			}

		})
	}
}

func Test_logsAPI_List(t *testing.T) {
	type fields struct {
		logsUC usecase.LogUsecase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		body       io.Reader
		fields     fields
		wantStatus int
		wantBody   string
	}{
		{
			name: "test get all logs",
			body: http.NoBody,
			fields: fields{
				logsUC: mockUsecase.NewlogMockUseCase(),
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api := logsAPI{
				logsUC: tt.fields.logsUC,
			}
			request, _ := http.NewRequest(http.MethodGet, "/api/drone/log", nil)
			response := httptest.NewRecorder()
			api.List(response, request)
			if status := response.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		})
	}
}
