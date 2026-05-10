package alerts

import (
	"log"

	"alert-lifecycle-engine/internal/common"
	"alert-lifecycle-engine/pkg/transport"
)

type AlertService struct {
	repo      *PostgresAlertRepository
	logClient *transport.LogClient
}

func NewAlertService(repo *PostgresAlertRepository, logClient *transport.LogClient) *AlertService {
	return &AlertService{repo: repo, logClient: logClient}
}

func (s *AlertService) ProcessEvent(ev common.DeviceEvent) *common.Alert {
	// Check for duplicate events
	processed, err := s.repo.IsEventProcessed(ev.EventID)
	if err != nil {
		log.Println("Error checking event:", err)
		return nil
	}
	if processed {
		return nil // Skip duplicate
	}

	// Mark as processed
	s.repo.MarkEventProcessed(ev.EventID)

	alertID := "alert-" + ev.DeviceID
	existing, err := s.repo.GetAlertByID(alertID)
	if err != nil {
		// If not exists, create new
		alert := &common.Alert{
			AlertID:   alertID,
			DeviceID:  ev.DeviceID,
			CreatedAt: ev.Timestamp,
		}
		if ev.EventType == "DEVICE_DOWN" {
			alert.State = "ACTIVE"
		} else if ev.EventType == "DEVICE_UP" {
			alert.State = "RESOLVED"
			alert.ResolvedAt = ev.Timestamp
		}
		s.repo.Save(*alert)
		s.logClient.Log(*alert)
		return alert
	}

	// Update existing alert based on event
	if ev.EventType == "DEVICE_DOWN" && existing.State == "RESOLVED" {
		existing.State = "ACTIVE"
		existing.CreatedAt = ev.Timestamp
		existing.ResolvedAt = 0
	} else if ev.EventType == "DEVICE_UP" && existing.State != "RESOLVED" {
		existing.State = "RESOLVED"
		existing.ResolvedAt = ev.Timestamp
	}
	s.repo.Save(existing)
	s.logClient.Log(existing)
	return &existing
}
