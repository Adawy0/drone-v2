package main

import (
	"drone/v2/repository"
	db "drone/v2/repository"
	server "drone/v2/server"
	"drone/v2/usecase"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func runCornJob(d usecase.IDroneUsecase) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Minutes().Do(func() {
		d.CheckDronesBatteries()
	})

	s.StartAsync()
}
func main() {
	fmt.Println("Dorne start")

	DB, err := db.Init()
	if err != nil {
		log.Println("cant connect to database")
		return
	}
	droneRepo := repository.NewDroneRepo(DB)
	droneUseCase := usecase.NewDroneUsecase(droneRepo)
	droneAPI := server.NewDroneAPI(droneUseCase)
	apis := server.APIs{
		DroneAPI: droneAPI,
	}

	go runCornJob(droneUseCase)

	server.StartServer(apis)
}
