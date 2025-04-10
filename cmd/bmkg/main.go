package main

import (
	"bmkg/src/repository"
	"bmkg/src/worker/bmkg"
	"fmt"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/router"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "bmkg/migrations"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	// new repo
	bmkgRepo := repository.NewBMKGRepository(app)
	bmkgWorker := bmkg.NewBMKGWorker(bmkgRepo)

	// on boostrap
	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		bmkgWorker.StartWorker()
		bmkgWorker.StopWorker()

		return e.Next()
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		funcName(se)

		se.Router.GET("/stats", func(e *core.RequestEvent) error {
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

				err := app.DB().NewQuery(gempaQuery).One(&gempaResult)
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
				iotQuery := "SELECT COUNT(*) as count FROM iot"

				err := app.DB().NewQuery(iotQuery).One(&iotResult)
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
		})

		// list device broken
		// select from history

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func funcName(se *core.ServeEvent) *router.Route[*core.RequestEvent] {
	return se.Router.GET("/hello/{name}", func(e *core.RequestEvent) error {
		name := e.Request.PathValue("name")

		return e.String(http.StatusOK, "Hello "+name)
	})
}
