CREATE TABLE IF NOT EXISTS alerts (
  key TEXT PRIMARY KEY,
  data JSONB
);

CREATE TABLE IF NOT EXISTS logs (
  log_id TEXT PRIMARY KEY,
  device_id TEXT,
  data JSONB
);

CREATE TABLE IF NOT EXISTS processed_events (
  event_id TEXT PRIMARY KEY,
  processed_at TIMESTAMP DEFAULT NOW()
);

