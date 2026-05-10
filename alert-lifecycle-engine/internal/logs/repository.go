package logs

import (
	"database/sql"
	"log"

	"alert-lifecycle-engine/internal/common"

	_ "github.com/lib/pq"
)

type PostgresLogRepository struct {
	db *sql.DB
}

func NewPostgresLogRepository(conn string) (*PostgresLogRepository, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	// Ensure table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logs (
		log_id TEXT PRIMARY KEY,
		device_id TEXT,
		data JSONB
	)`)
	if err != nil {
		return nil, err
	}

	return &PostgresLogRepository{db: db}, nil
}

func (r *PostgresLogRepository) Save(logEntry common.Log) {
	_, err := r.db.Exec(`
        INSERT INTO logs (log_id, device_id, data)
        VALUES ($1, $2, $3)
        ON CONFLICT (log_id) DO UPDATE SET data = EXCLUDED.data
    `, logEntry.LogID, logEntry.DeviceID, logEntry.Data)
	if err != nil {
		log.Println("Error saving log:", err)
	}
}

func (r *PostgresLogRepository) GetAllLogs() ([]common.Log, error) {
	rows, err := r.db.Query(`SELECT log_id, device_id, data FROM logs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []common.Log
	for rows.Next() {
		var l common.Log
		if err := rows.Scan(&l.LogID, &l.DeviceID, &l.Data); err == nil {
			logs = append(logs, l)
		}
	}
	return logs, nil
}

func (r *PostgresLogRepository) GetLogsByDevice(deviceID string) ([]common.Log, error) {
	rows, err := r.db.Query(`SELECT log_id, device_id, data FROM logs WHERE device_id=$1`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []common.Log
	for rows.Next() {
		var l common.Log
		if err := rows.Scan(&l.LogID, &l.DeviceID, &l.Data); err == nil {
			logs = append(logs, l)
		}
	}
	return logs, nil
}
