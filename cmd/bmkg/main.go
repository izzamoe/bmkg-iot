package main

import (
	"bmkg/src/repository"
	"bmkg/src/worker/bmkg"
	"log"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

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
		se.Router.GET("/hello/{name}", func(e *core.RequestEvent) error {
			name := e.Request.PathValue("name")

			return e.String(http.StatusOK, "Hello "+name)
		})

		se.Router.GET("/gempa", func(e *core.RequestEvent) error {
			// Get today's date
			now := time.Now()
			todayDateStr := now.Format("2006-01-02")

			// Struktur untuk menerima hasil COUNT
			type CountResult struct {
				Count int `db:"count"`
			}

			var result CountResult

			// Query optimasi: langsung gunakan date() function di SQLite
			query := "SELECT COUNT(*) as count FROM view_gempa WHERE date(created) = date('now')"

			err := app.DB().NewQuery(query).One(&result)

			if err != nil {
				return e.String(http.StatusInternalServerError, "Error querying gempa data: "+err.Error())
			}

			return e.JSON(http.StatusOK, map[string]interface{}{
				"count": result.Count,
				"date":  todayDateStr,
			})

		})

		// Endpoint untuk menghitung total data IOT
		se.Router.GET("/iot", func(e *core.RequestEvent) error {
			// Struktur untuk menerima hasil COUNT
			type CountResult struct {
				Count int `db:"count"`
			}

			var result CountResult

			// Query untuk menghitung seluruh data di tabel iot
			query := "SELECT COUNT(*) as count FROM iot"

			err := app.DB().NewQuery(query).One(&result)

			if err != nil {
				return e.String(http.StatusInternalServerError, "Error querying IOT data: "+err.Error())
			}

			return e.JSON(http.StatusOK, map[string]interface{}{
				"total_count": result.Count,
			})
		})

		// list device broken
		// select from history

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
