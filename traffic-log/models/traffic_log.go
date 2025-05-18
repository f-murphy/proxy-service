package models

import "time"

type TrafficLog struct {
	ID            int       `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	ClientIP      string    `json:"client_ip"`
	Method        string    `json:"method"`
	URL           string    `json:"url"`
	ResponseCode  int       `json:"response_code"`
	BytesSent     int64     `json:"bytes_sent"`
	BytesReceived int64     `json:"bytes_received"`
}
