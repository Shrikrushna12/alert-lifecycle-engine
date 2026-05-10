package logs

import (
	"alert-lifecycle-engine/internal/common"
)

type LogService struct {
	repo *PostgresLogRepository
}

func NewLogService(repo *PostgresLogRepository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) AddLog(log common.Log) {
	s.repo.Save(log)
}

func (s *LogService) GetLogs() ([]common.Log, error) {
	return s.repo.GetAllLogs()
}

func (s *LogService) GetLogsByDevice(deviceID string) ([]common.Log, error) {
	return s.repo.GetLogsByDevice(deviceID)
}
