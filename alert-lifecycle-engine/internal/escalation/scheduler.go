package escalation

import (
	"log"
	"time"

	"alert-lifecycle-engine/internal/alerts"
	"alert-lifecycle-engine/pkg/transport"
)

// Scheduler periodically escalates unresolved alerts
type Scheduler struct {
	repo      *alerts.PostgresAlertRepository
	logClient *transport.LogClient
	service   *alerts.EscalationService
	interval  time.Duration
}

// NewScheduler creates a new escalation scheduler
func NewScheduler(repo *alerts.PostgresAlertRepository, logClient *transport.LogClient, interval time.Duration) *Scheduler {
	return &Scheduler{
		repo:      repo,
		logClient: logClient,
		service:   alerts.NewEscalationService(repo, logClient),
		interval:  interval,
	}
}

// Run starts the escalation loop
func (s *Scheduler) Run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndEscalate()
		}
	}
}

// checkAndEscalate fetches active alerts and escalates them based on time
func (s *Scheduler) checkAndEscalate() {
	activeAlerts, err := s.repo.GetActiveAlerts()
	if err != nil {
		log.Println("Error fetching active alerts:", err)
		return
	}

	policy := DefaultPolicy()
	now := time.Now().Unix()
	for _, a := range activeAlerts {
		duration := now - a.CreatedAt
		var newState string
		switch {
		case duration >= int64(policy.Thresholds[2]):
			newState = "CRITICAL"
		case duration >= int64(policy.Thresholds[1]):
			newState = "ATTENTION"
		case duration >= int64(policy.Thresholds[0]):
			newState = "WARNING"
		default:
			continue // Not yet time to escalate
		}

		if a.State != newState {
			a.State = newState
			s.repo.Save(a)
			log.Printf("Alert %s escalated to %s\n", a.AlertID, a.State)
			s.logClient.Log(a)
		}
	}
}
