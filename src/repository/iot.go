package repository

import (
	"bmkg/src/domain"
	"fmt"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"time"
)

// iot repository buat struct untuk save data
type Iot struct {
	App core.App // Pocketbase app instance
}

// NewIotRepository creates a new Iot repository
func NewIotRepository(app core.App) *Iot {
	return &Iot{
		App: app,
	}
}

// SaveData saves data to the database
func (r *Iot) SaveData(data map[string]interface{}) error {
	collection, err := r.App.FindCollectionByNameOrId("iot")
	if err != nil {
		return err
	}

	// Create a new record
	record := core.NewRecord(collection)

	// Set values from map
	for key, value := range data {
		record.Set(key, value)
	}

	// Save the record
	if err := r.App.Save(record); err != nil {
		return err
	}

	return nil
}

// get all data from iot name iot_device to Domain.Iot
func (r *Iot) GetAllData() ([]domain.IotDevice, error) {
	var devices []domain.IotDevice

	err := r.App.DB().
		NewQuery("SELECT id, device_name, created FROM iot_device ORDER BY created DESC").
		All(&devices)

	if err != nil {
		return nil, err
	}

	return devices, nil
}

// GetBrokenDevice identifies devices that are offline based on:
// 1. No updates in the last 5 minutes, or
// 2. Last status update shows inactive (active=false)
func (r *Iot) GetBrokenDevice() ([]domain.IotDevice, error) {
	var devices []domain.IotDevice

	// Calculate the timestamp from 5 minutes ago
	fiveMinutesAgo := time.Now().UTC().Add(-5 * time.Minute)

	// Query to find broken/offline devices
	query := `
		SELECT d.id, d.device_name 
		FROM iot_device d
		LEFT JOIN (
			SELECT h1.device, h1.active, h1.created
			FROM history_iot h1
			LEFT JOIN history_iot h2 ON h1.device = h2.device AND h1.created < h2.created
			WHERE h2.created IS NULL
		) h ON d.id = h.device
		WHERE h.created < {:threshold} OR h.active = false OR h.created IS NULL
	`

	err := r.App.DB().
		NewQuery(query).
		Bind(dbx.Params{
			"threshold": fiveMinutesAgo,
		}).
		All(&devices)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch broken devices: %w", err)
	}

	return devices, nil
}
