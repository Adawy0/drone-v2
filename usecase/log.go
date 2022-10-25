package usecase

import (
	repo "drone/v2/repository"
	repoEnity "drone/v2/repository"
	"encoding/json"
)

type LogUsecase interface {
	Create(log repoEnity.Log) error
	List() ([]byte, error)
}

type logUsecase struct {
	logRepo repo.ILogRepository
}

func NewlogUseCase(repo repo.ILogRepository) LogUsecase {
	return logUsecase{
		logRepo: repo,
	}
}

func (l logUsecase) List() ([]byte, error) {

	logs, err := l.logRepo.List()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(logs)
}

func (l logUsecase) Create(log repoEnity.Log) error {
	result := l.logRepo.Create(log)
	return result
}
