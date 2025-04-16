package handler

import (
	"bmkg/src/db"
	"bmkg/src/repository"
	"bmkg/src/utils"
	hehe "bmkg/src/worker/mqtt"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"net/http"
	"strings"
)

// IotHandler handles IoT device operations
type IotHandler struct {
	mqtt    hehe.MQTTClient
	iotrepo *repository.Iot
	app     core.App
}

// NewIotHandler creates a new instance of IotHandler
func NewIotHandler(mqttClient hehe.MQTTClient, iotrepos *repository.Iot, app core.App) *IotHandler {
	return &IotHandler{
		mqtt:    mqttClient,
		iotrepo: iotrepos,
		app:     app,
	}
}

// AddIotHandler registers IoT route handlers to the router
func (h *IotHandler) AddIotHandler(router *router.Router[*core.RequestEvent]) {
	group := router.Group("/iot")
	group.GET("/on/{id}", h.soundDevice)
	// register device
	group.GET("/create", h.createIotDevice)

	//	 register subscribe
	h.mqtt.Subscribe("device/register/iot", 1, h.mqttHandler)
}

// soundDevice activates the sound on a specified IoT device
func (h *IotHandler) soundDevice(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")

	// Verify device exists in database
	_, err := e.App.FindRecordById("iot_device", id)
	if err != nil {
		return e.String(http.StatusNotFound, fmt.Sprintf("Device %s not found", id))
	}

	// Publish message to device topic
	err = h.mqtt.Publish("device/"+id, 0, false, "on")
	if err != nil {
		return e.String(http.StatusServiceUnavailable, fmt.Sprintf("Failed to communicate with device %s", id))
	}

	return e.String(http.StatusOK, fmt.Sprintf("Device %s is sounding", id))
}

// create new iot device
func (h *IotHandler) createIotDevice(e *core.RequestEvent) error {
	// Parse request body

	return e.String(http.StatusCreated, "Device created successfully")
}

// mqttHandler make handler mqtt
func (h *IotHandler) mqttHandler(client mqtt.Client, msg mqtt.Message) {

	// Create a new IoT device proxy
	deviceProxy, _ := db.NewProxy[db.IotDevice](h.app)

	// handle the message
	// Split payload by comma to parse device registration data
	str := string(msg.Payload())
	parts := strings.Split(str, ",")

	lintang := parts[0]
	bujur := parts[1]

	// Set device properties
	deviceProxy.SetName("device-id-" + utils.GenerateRandomID())
	deviceProxy.SetLintang(lintang)
	deviceProxy.SetBujur(bujur)

	// Save the record to database
	//if err := deviceProxy.; err != nil {
	h.app.Save(deviceProxy)

	data := deviceProxy.Id

	//// generate device name
	//data, err := h.iotrepo.SaveData(datas)
	//if err != nil {
	//	return
	//}

	fmt.Println("Device registered successfully with ID:", data)

	// publish id
	_ = h.mqtt.Publish("device/berhasil", 0, false, data)

	//fmt.Printf("Received message on topic %s: %s\n", msg.Topic())
}
