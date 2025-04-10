package domain

import "time"

type IotDevice struct {
}

// history device
// make history_iot
// HistoryIot represents the history of an IoT device's activity
type HistoryIot struct {
	// Device is the identifier of the IoT device
	Device string `json:"device"`
	// Active indicates whether the device is currently active
	Active bool `json:"active"`
	// Created is the timestamp when this history record was created
	Created time.Time `json:"created"`
}
