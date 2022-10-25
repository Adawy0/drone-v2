package repository

import (
	"gorm.io/gorm"
)

type LogDB struct {
	client *gorm.DB
}

type ILogRepository interface {
	Create(log Log) error
	List() ([]Log, error)
}

func NewLogRepository(client *gorm.DB) ILogRepository {
	return &LogDB{client: client}
}

func (ldb LogDB) List() ([]Log, error) {

	logs := []Log{}
	result := ldb.client.Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}
	return logs, nil
}

func (ldb LogDB) Create(log Log) error {
	result := ldb.client.Create(&log)
	return result.Error
}
