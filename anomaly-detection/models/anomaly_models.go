package models

import (
    "time"
)

// TrafficData представляет данные о сетевом трафике
type TrafficData struct {
    SourceIP   string    `json:"source_ip"`
    DestIP     string    `json:"dest_ip"`
    Port       int       `json:"port"`
    Protocol   string    `json:"protocol"`
    DataVolume int64     `json:"data_volume"`
    Timestamp  time.Time `json:"timestamp"`
}

// Anomaly представляет информацию об обнаруженной аномалии
type Anomaly struct {
    Detected bool   `json:"detected"`
    Type     string `json:"type"`
    Details  string `json:"details"`
}