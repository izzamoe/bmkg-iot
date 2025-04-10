package repository

import (
	"bmkg/src/domain"
	"fmt"
	"github.com/pocketbase/dbx"

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

	//Add additional metadata if needed
	// record.Set("processed", true)

	// Save the record
	if err := r.App.Save(record); err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}

	return nil
}

// GetAllGempa retrieves all earthquake data from the database
func (r *BMKG) GetAllGempa() ([]domain.Gempa, error) {

	var earthquakes []domain.Gempa

	err := r.App.DB().
		NewQuery("SELECT id, tanggal, jam, magnitude, kedalaman, wilayah, coordinates, created FROM earthquake ORDER BY created DESC").
		All(&earthquakes)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch earthquake data: %w", err)
	}

	return earthquakes, nil
}

// DeleteGempa deletes earthquake data from the database by ID
func (r *BMKG) DeleteGempa(id string) error {
	query := "DELETE FROM earthquake WHERE id = {:id}"

	_, err := r.App.DB().
		NewQuery(query).
		Bind(dbx.Params{
			"id": id,
		}).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete earthquake record: %w", err)
	}

	return nil
}

// GetLastGempa retrieves the most recent earthquake data from the database
func (r *BMKG) GetLastGempa() (domain.Gempa, error) {
	var gempa domain.Gempa

	err := r.App.DB().
		NewQuery("SELECT id, tanggal, jam, magnitude, kedalaman, wilayah, coordinates, created FROM earthquake ORDER BY created DESC LIMIT 1").
		One(&gempa)

	if err != nil {
		return domain.Gempa{}, fmt.Errorf("failed to fetch latest earthquake data: %w", err)
	}

	return gempa, nil
}

// get gempa cari yang unik group by Shakemap biar unik gempa hari ini , intinya cari gempa hari ini
func (r *BMKG) GetGempaHariIni() ([]domain.Gempa, error) {
	var gempa []domain.Gempa

	err := r.App.DB().
		NewQuery("SELECT id, tanggal, jam, magnitude, kedalaman, wilayah, coordinates, created FROM earthquake WHERE date(created) = date('now') GROUP BY shakemap ORDER BY created DESC").
		All(&gempa)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's earthquake data: %w", err)
	}

	return gempa, nil
}
