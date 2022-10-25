package server

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIs struct {
	DroneAPI IDroneAPI
	LogsAPI  LogsAPI
}

func StartServer(apis APIs) {
	port := flag.String("port", "4000", "Port to listen on")
	flag.Parse()

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	droneSubRouter := r.PathPrefix("/drone").Subrouter()
	droneSubRouter.HandleFunc("/", apis.DroneAPI.RegisterDrone).Methods("POST")
	droneSubRouter.HandleFunc("/{id}/load-medication", apis.DroneAPI.LoadingMedication).Methods("POST")
	droneSubRouter.HandleFunc("/{id}/check-battery", apis.DroneAPI.CheckDroneBattery).Methods("GET")
	droneSubRouter.HandleFunc("/available-drone", apis.DroneAPI.CheckAvailableDrones).Methods("GET")
	droneSubRouter.HandleFunc("/log", apis.LogsAPI.List).Methods("GET")
	start(*port, r)
}

func start(port string, r http.Handler) {
	loggingRouter := loggingHandler(r)
	log.Fatal(http.ListenAndServe(":"+port, loggingRouter))
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
