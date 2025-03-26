package repository

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

// BMKG repository for handling earthquake data
type BMKG struct {
	App core.App
}

func NewBMKGRepository(app core.App) *BMKG {
	return &BMKG{
		App: app,
	}
}

// SaveGempa saves earthquake data to the database
func (r *BMKG) SaveGempa(data map[string]interface{}) error {
	collection, err := r.App.FindCollectionByNameOrId("earthquake")
	if err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	// Create a new record
	record := core.NewRecord(collection)

	// Set values from map
	for key, value := range data {
		record.Set(key, value)
	}

	// Add additional metadata if needed
	// record.Set("processed", true)

	// Save the record
	if err := r.App.Save(record); err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}

	return nil
}
