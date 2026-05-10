package alerts

import (
	"time"

	"alert-lifecycle-engine/internal/common"
	"alert-lifecycle-engine/pkg/transport"
)

type EscalationService struct {
	repo      *PostgresAlertRepository
	logClient *transport.LogClient
}

func NewEscalationService(repo *PostgresAlertRepository, logClient *transport.LogClient) *EscalationService {
	return &EscalationService{repo: repo, logClient: logClient}
}

func (s *EscalationService) EscalateAlert(alert common.Alert) common.Alert {
	switch alert.State {
	case "ACTIVE":
		alert.State = "WARNING"
	case "WARNING":
		alert.State = "ATTENTION"
	case "ATTENTION":
		alert.State = "CRITICAL"
	}

	alert.ResolvedAt = time.Now().Unix()
	s.repo.Save(alert)
	s.logClient.Log(alert)
	return alert
}
