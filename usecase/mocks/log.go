package mocks

import (
	repoEnity "drone/v2/repository"
)

type LogMockUsecase interface {
	Create(log repoEnity.Log) error
	List() ([]byte, error)
}

type logMockUsecase struct {
}

func NewlogMockUseCase() LogMockUsecase {
	return &logMockUsecase{}
}

func (l logMockUsecase) List() ([]byte, error) {
	return nil, nil
}

func (l logMockUsecase) Create(log repoEnity.Log) error {
	return nil
}
