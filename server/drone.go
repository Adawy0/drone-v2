package server

import (
	"drone/v2/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	usecaseEntity "drone/v2/usecase"

	"github.com/gorilla/mux"
)

type IDroneAPI interface {
	RegisterDrone(w http.ResponseWriter, r *http.Request)
	RegisterMedication(w http.ResponseWriter, r *http.Request)
	LoadingMedication(w http.ResponseWriter, r *http.Request)
	CheckLoadingMedication(w http.ResponseWriter, r *http.Request)
	CheckAvailableDrones(w http.ResponseWriter, r *http.Request)
	CheckDroneBattery(w http.ResponseWriter, r *http.Request)
}

type droneAPI struct {
	droneUsecase      usecase.IDroneUsecase
	medicationUsecase usecase.IMedicationUsecase
}

func NewDroneAPI(droneUsecase usecase.IDroneUsecase) IDroneAPI {
	return &droneAPI{
		droneUsecase: droneUsecase,
	}
}

func (api *droneAPI) RegisterDrone(w http.ResponseWriter, r *http.Request) {
	var drone DornePayload
	if r.Body == nil {
		http.Error(w, "register drone must have json payload", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&drone)
	if err != nil {
		http.Error(w, "Invaild json payload", http.StatusBadRequest)
		return
	}

	id, err := api.droneUsecase.RegisterDrone(usecaseEntity.DorneObject(drone))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var p RegisterDronePayload
	p.DroneId = id
	payload, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (api *droneAPI) RegisterMedication(w http.ResponseWriter, r *http.Request) {
	var medication MedicationPayload
	if r.Body == nil {
		http.Error(w, "register medication must have json payload", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&medication)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := api.medicationUsecase.RegisterMedication(usecaseEntity.MedicationObject(medication))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload := RegisterMediactionPayload{
		MedicationId: id,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (api *droneAPI) LoadingMedication(w http.ResponseWriter, r *http.Request) {
	args, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Couldn't find id in request URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(args)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Invaild drone id", http.StatusBadRequest)
		return
	}
	var medication MedicationPayload
	if r.Body == nil {
		http.Error(w, "load medication end point must have json payload", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&medication)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = api.droneUsecase.LoadingMedication(id, usecaseEntity.MedicationObject(medication))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

func (api *droneAPI) CheckLoadingMedication(w http.ResponseWriter, r *http.Request) {
	args, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Couldn't find id in request URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(args)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Invaild drone id", http.StatusBadRequest)
		return
	}
	status, err := api.droneUsecase.CheckLoadingMedication(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload := CheckLoadingMedicationPayload{
		DroneId: id,
		Status:  status,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (api *droneAPI) CheckAvailableDrones(w http.ResponseWriter, r *http.Request) {
	drones := api.droneUsecase.CheckAvailableDroneForLoading()
	data, err := json.Marshal(drones)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (api *droneAPI) CheckDroneBattery(w http.ResponseWriter, r *http.Request) {
	args, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Couldn't find id in request URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(args)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Invaild drone id", http.StatusBadRequest)
		return
	}
	battryLevel, err := api.droneUsecase.CheckBatteryLevel(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload := BatteryLevelPayload{
		DroneId:      id,
		BatteryLevel: battryLevel,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
