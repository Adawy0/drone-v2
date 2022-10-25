package server

import (
	"drone/v2/usecase"
	"fmt"
	"log"
	"net/http"
)

type LogsAPI interface {
	List(w http.ResponseWriter, r *http.Request)
}

type logsAPI struct {
	logsUC usecase.LogUsecase
}

func NewLogsAPI(uc usecase.LogUsecase) LogsAPI {
	return &logsAPI{
		logsUC: uc,
	}
}

func (api logsAPI) List(w http.ResponseWriter, r *http.Request) {
	response, err := api.logsUC.List()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
