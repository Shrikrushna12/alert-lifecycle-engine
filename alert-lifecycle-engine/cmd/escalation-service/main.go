package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"alert-lifecycle-engine/internal/alerts"
	"alert-lifecycle-engine/internal/common"
	"alert-lifecycle-engine/internal/escalation"
	"alert-lifecycle-engine/pkg/transport"
)

func main() {
	repo, err := alerts.NewPostgresAlertRepository("postgres://postgres:postgres@postgres:5432/alertsdb?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}

	logClient := transport.NewLogClient("http://log-service:8081")
	esService := alerts.NewEscalationService(repo, logClient)

	// Manual endpoint to trigger escalation
	http.HandleFunc("/escalate", func(w http.ResponseWriter, r *http.Request) {
		activeAlerts, err := repo.GetActiveAlerts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		results := []common.Alert{} // ✅ use common.Alert
		for _, a := range activeAlerts {
			updated := esService.EscalateAlert(a)
			results = append(results, updated)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Background scheduler
	scheduler := escalation.NewScheduler(repo, logClient, 10*time.Second)
	go scheduler.Run()

	log.Println("✅ Escalation Service running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
