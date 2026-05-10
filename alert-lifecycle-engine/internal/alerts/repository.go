package alerts

import (
	"database/sql"
	"log"

	"alert-lifecycle-engine/internal/common"

	_ "github.com/lib/pq"
)

type PostgresAlertRepository struct {
	db *sql.DB
}

func NewPostgresAlertRepository(conn string) (*PostgresAlertRepository, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	// Ensure tables exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS alerts (
		key TEXT PRIMARY KEY,
		data JSONB
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS processed_events (
		event_id TEXT PRIMARY KEY,
		processed_at TIMESTAMP DEFAULT NOW()
	)`)
	if err != nil {
		return nil, err
	}

	return &PostgresAlertRepository{db: db}, nil
}

func (r *PostgresAlertRepository) Save(alert common.Alert) {
	_, err := r.db.Exec(`
        INSERT INTO alerts (key, data)
        VALUES ($1, $2)
        ON CONFLICT (key) DO UPDATE SET data = EXCLUDED.data
    `, alert.AlertID, alert)
	if err != nil {
		log.Println("Error saving alert:", err)
	}
}

func (r *PostgresAlertRepository) GetAlertByID(alertID string) (common.Alert, error) {
	var a common.Alert
	err := r.db.QueryRow(`SELECT data FROM alerts WHERE key = $1`, alertID).Scan(&a)
	return a, err
}

func (r *PostgresAlertRepository) GetActiveAlerts() ([]common.Alert, error) {
	rows, err := r.db.Query(`SELECT data FROM alerts WHERE data->>'state' != 'RESOLVED'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []common.Alert
	for rows.Next() {
		var a common.Alert
		if err := rows.Scan(&a); err == nil {
			alerts = append(alerts, a)
		}
	}
	return alerts, nil
}

func (r *PostgresAlertRepository) IsEventProcessed(eventID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM processed_events WHERE event_id = $1`, eventID).Scan(&count)
	return count > 0, err
}

func (r *PostgresAlertRepository) MarkEventProcessed(eventID string) error {
	_, err := r.db.Exec(`INSERT INTO processed_events (event_id) VALUES ($1) ON CONFLICT DO NOTHING`, eventID)
	return err
}
