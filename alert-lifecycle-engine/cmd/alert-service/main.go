package main

import (
	"encoding/json"
	"log"
	"net/http"

	"alert-lifecycle-engine/internal/alerts"
	"alert-lifecycle-engine/internal/common"
	"alert-lifecycle-engine/pkg/transport"
)

func main() {
	repo, err := alerts.NewPostgresAlertRepository("postgres://postgres:postgres@postgres:5432/alertsdb?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}

	logClient := transport.NewLogClient("http://log-service:8081")
	service := alerts.NewAlertService(repo, logClient)

	http.HandleFunc("/alerts/events", func(w http.ResponseWriter, r *http.Request) {
		var ev common.DeviceEvent
		if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		alert := service.ProcessEvent(ev)
		if alert == nil {
			http.Error(w, "no alert generated", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(alert)
	})

	log.Println("✅ Alert Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
