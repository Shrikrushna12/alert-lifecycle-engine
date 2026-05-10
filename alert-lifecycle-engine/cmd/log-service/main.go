package main

import (
	"encoding/json"
	"log"
	"net/http"

	"alert-lifecycle-engine/internal/common"
	"alert-lifecycle-engine/internal/logs"
)

func main() {
	repo, err := logs.NewPostgresLogRepository("postgres://postgres:postgres@postgres:5432/alertsdb?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		allLogs, err := repo.GetAllLogs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if allLogs == nil {
			allLogs = []common.Log{} // ✅ use common.Log, not logs.Log
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(allLogs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Path[len("/logs/"):]
		deviceLogs, err := repo.GetLogsByDevice(deviceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if deviceLogs == nil {
			deviceLogs = []common.Log{} // ✅ same fix here
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(deviceLogs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("✅ Log Service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
