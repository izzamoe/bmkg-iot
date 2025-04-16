package main

import (
	_ "bmkg/migrations"
	"bmkg/src/handler"
	"bmkg/src/repository"
	"bmkg/src/worker/bmkg"
	"bmkg/src/worker/mqtt"
	"bmkg/src/worker/telegram"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"log"
	"os"
	"strings"
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

	// Configure MQTT for HiveMQ public broker
	mqttconfig := mqtt.Config{
		Broker:   "tcp://broker.emqx.io:1883", // Using TCP port for standard MQTT
		ClientID: "BE",                        // Unique client ID
		QoS:      1,
	}

	// Initialize MQTT client
	mqttClient := mqtt.NewMQTTClient(mqttconfig)
	mqttClient.Connect()
	bot, err := telegram.NewBot(app)

	// new repo
	bmkgRepo := repository.NewBMKGRepository(app)
	iotRepo := repository.NewIotRepository(app)

	bmkgWorker := bmkg.NewBMKGWorker(bmkgRepo, app, mqttClient, bot)
	if err != nil {
		return
	}

	go bot.Start()

	// handler
	iotHandler := handler.NewIotHandler(mqttClient, iotRepo, app)
	bmkgHandler := handler.NewBMKGHandler()

	//

	app.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		mqttClient.Disconnect()
		return nil
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		bmkgWorker.StartWorker()
		bmkgWorker.StopWorker()

		bmkgHandler.AddBMKGHandler(se.Router)
		handler.AddAdminHandler(se.Router)
		iotHandler.AddIotHandler(se.Router)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
