package main

import (
	"bmkg/src/repository"
	"bmkg/src/worker/bmkg"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// new repo
	bmkgRepo := repository.NewBMKGRepository(app)
	bmkgWorker := bmkg.NewBMKGWorker(bmkgRepo)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		bmkgWorker.StartWorker()
		bmkgWorker.StopWorker()

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
