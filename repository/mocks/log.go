package mocks

import (
	repo "drone/v2/repository"
)

type LogRepoMock struct {
}

func NewLogRepoMock() repo.ILogRepository {
	return &LogRepoMock{}
}

func (d *LogRepoMock) Create(log repo.Log) error {
	return nil
}

func (d *LogRepoMock) List() ([]repo.Log, error) {
	return []repo.Log{}, nil
}
