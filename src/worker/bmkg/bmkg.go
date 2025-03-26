package bmkg

import (
	"bmkg/src/domain"
	"bmkg/src/repository"
	"bmkg/src/utils"
	"context"
	"encoding/json"
	"log"
	"time"
)

// BMKGWorker handles earthquake data processing from BMKG
type BMKGWorker struct {
	repo       *repository.BMKG
	cancelFunc context.CancelFunc
	ctx        context.Context
}

// NewBMKGWorker creates a new instance of BMKGWorker
func NewBMKGWorker(repo *repository.BMKG) *BMKGWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &BMKGWorker{
		repo:       repo,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// StartWorker initiates a worker that periodically fetches earthquake data from the BMKG API.
func (w *BMKGWorker) StartWorker() {
	// Run as a goroutine to prevent blocking
	go func() {
		log.Println("Starting BMKG worker...")

		// Use ticker for periodic execution
		ticker := time.NewTicker(timing * time.Second)
		defer ticker.Stop()

		// Do initial fetch immediately
		w.fetchAndProcessData()

		for {
			select {
			case <-ticker.C:
				w.fetchAndProcessData()
			case <-w.ctx.Done():
				log.Println("BMKG worker stopped")
				return
			}
		}
	}()
}

// fetchAndProcessData handles the API call and data processing
func (w *BMKGWorker) fetchAndProcessData() {
	response, err := hitAPI()
	if err != nil {
		log.Printf("Error fetching data from the BMKG API: %v", err)
		return
	}

	// Process and save the response
	err = w.processAndSaveEarthquake(response)
	if err != nil {
		log.Printf("Error processing earthquake data: %v", err)
	} else {
		log.Println("Successfully processed and saved earthquake data")
	}
}

// StopWorker gracefully stops the BMKG worker
func (w *BMKGWorker) StopWorker() {
	log.Println("Stopping BMKG worker...")
	if w.cancelFunc != nil {
		w.cancelFunc()
	}
}

// processAndSaveEarthquake converts BMKG API response to a map and saves it to the database
func (w *BMKGWorker) processAndSaveEarthquake(response domain.ResponseBmkgAPI) error {
	// Convert struct to map for flexible storage
	gempaData := utils.StructToMap(response.Infogempa.Gempa)

	// Add timestamp for when this data was fetched
	gempaData["fetchedAt"] = time.Now().Format(time.RFC3339)

	// Save to repository
	return w.repo.SaveGempa(gempaData)
}

// hitAPI fetches the earthquake data from the BMKG API.
//
// This function makes an HTTP GET request to the BMKG API endpoint.
// It then parses the JSON response into a domain.ResponseBmkgAPI struct.
//
// Returns:
//   - domain.ResponseBmkgAPI: A struct containing the parsed response from the BMKG API.
//   - error: An error if the request fails or if the response cannot be parsed.
//
// Example:
//
//	response, err := hitAPI()
//	if err != nil {
//	    log.Fatalf("Failed to fetch earthquake data: %v", err)
//	}
//	fmt.Printf("Latest earthquake: %+v\n", response.Infogempa.Gempa)
func hitAPI() (domain.ResponseBmkgAPI, error) {
	// hit api bmkg
	respons, err := utils.NewHTTPClient().Get(bmkgAPIURL)

	if err != nil {
		return domain.ResponseBmkgAPI{}, err
	}

	// decode response
	var response domain.ResponseBmkgAPI
	if err := json.Unmarshal(respons, &response); err != nil {
		return domain.ResponseBmkgAPI{}, err
	}
	return response, nil
}
