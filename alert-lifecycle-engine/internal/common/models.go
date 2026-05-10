package common

// DeviceEvent represents an incoming event from a device
type DeviceEvent struct {
	TenantID  string `json:"tenant_id"`
	DeviceID  string `json:"device_id"`
	EventType string `json:"event_type"` // DEVICE_DOWN / DEVICE_UP
	Timestamp int64  `json:"timestamp"`
	EventID   string `json:"event_id"`
}

// Alert represents the lifecycle of a device alert
type Alert struct {
	AlertID    string `json:"alert_id"`
	DeviceID   string `json:"device_id"`
	State      string `json:"state"` // ACTIVE, RESOLVED, WARNING, ATTENTION, CRITICAL
	CreatedAt  int64  `json:"created_at"`
	ResolvedAt int64  `json:"resolved_at,omitempty"`
}

// Log entry stored in Log Service
type Log struct {
	LogID    string      `json:"log_id"`
	DeviceID string      `json:"device_id"`
	Data     interface{} `json:"data"`
}
