package handler

import (
	"fmt"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"net/http"
	"strings"
	"sync"
	"time"
)

// AddAdminHandler add handler admin
func AddAdminHandler(router *router.Router[*core.RequestEvent]) {
	group := router.Group("/admin")
	group.GET("/stat", stats)
}

func stats(e *core.RequestEvent) error {
	// Get today's date
	now := time.Now()
	todayDateStr := now.Format("2006-01-02")

	// Define combined response struct
	type StatsResponse struct {
		GempaCount int    `json:"gempa_count"`
		IotCount   int    `json:"iot_count"`
		Date       string `json:"date"`
	}

	// Struktur untuk menerima hasil COUNT
	type CountResult struct {
		Count int `db:"count"`
	}

	// Initialize response
	response := StatsResponse{
		Date: todayDateStr,
	}

	// Use a mutex to protect concurrent writes
	var mu sync.Mutex
	var errors []string

	// Use a WaitGroup to wait for both goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine for gempa query
	go func() {
		defer wg.Done()

		var gempaResult CountResult
		gempaQuery := "SELECT COUNT(*) as count FROM view_gempa WHERE date(created) = date('now')"

		err := e.App.DB().NewQuery(gempaQuery).One(&gempaResult)
		if err != nil {
			mu.Lock()
			errors = append(errors, fmt.Sprintf("Error querying gempa data: %s", err.Error()))
			mu.Unlock()
			return
		}

		mu.Lock()
		response.GempaCount = gempaResult.Count
		mu.Unlock()
	}()

	// Goroutine for IOT query
	go func() {
		defer wg.Done()

		var iotResult CountResult
		iotQuery := "SELECT COUNT(*) as count FROM iot_device"

		err := e.App.DB().NewQuery(iotQuery).One(&iotResult)
		if err != nil {
			mu.Lock()
			errors = append(errors, fmt.Sprintf("Error querying IOT data: %s", err.Error()))
			mu.Unlock()
			return
		}

		mu.Lock()
		response.IotCount = iotResult.Count
		mu.Unlock()
	}()

	// Wait for both queries to complete
	wg.Wait()

	// Check for errors
	if len(errors) > 0 {
		return e.String(http.StatusInternalServerError, strings.Join(errors, "; "))
	}

	return e.JSON(http.StatusOK, response)
}
