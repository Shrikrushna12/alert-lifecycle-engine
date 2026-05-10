package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"alert-lifecycle-engine/internal/common"
)

type LogClient struct {
	BaseURL string
}

func NewLogClient(baseURL string) *LogClient {
	return &LogClient{BaseURL: baseURL}
}

func (lc *LogClient) Log(alert common.Alert) error {
	entry := common.Log{
		LogID:    "log-" + time.Now().Format("20060102150405"),
		DeviceID: alert.DeviceID,
		Data:     alert,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	resp, err := http.Post(lc.BaseURL+"/logs", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
