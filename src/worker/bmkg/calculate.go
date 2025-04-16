package bmkg

import (
	"bmkg/src/db"
	"bmkg/src/domain"
	"bmkg/src/utils"
	"bmkg/src/utils/ngitung"
	"bmkg/src/utils/notify"
	"bmkg/src/worker/mqtt"
	"bmkg/src/worker/telegram"
	"fmt"
	"github.com/pocketbase/pocketbase/core"
	"strconv"
)

// CalculateAndNotify processes earthquake data and sends notifications to affected users and devices
func CalculateAndNotify(app core.App, mqtt mqtt.MQTTClient, response domain.ResponseBmkgAPI, telegrambot *telegram.Bot) error {
	// Extract earthquake data
	magnitude, err := strconv.ParseFloat(response.Infogempa.Gempa.Magnitude, 64)
	if err != nil {
		return fmt.Errorf("failed to parse magnitude: %w", err)
	}

	earthquakeLat, err := strconv.ParseFloat(utils.ExtractNumber(response.Infogempa.Gempa.Lintang), 64)
	if err != nil {
		return fmt.Errorf("failed to parse earthquake latitude: %w", err)
	}

	earthquakeLon, err := strconv.ParseFloat(utils.ExtractNumber(response.Infogempa.Gempa.Bujur), 64)
	if err != nil {
		return fmt.Errorf("failed to parse earthquake longitude: %w", err)
	}

	earthquakeLocation := ngitung.Location{
		Lat: earthquakeLat,
		Lon: earthquakeLon,
	}

	// Build notification message with earthquake details
	notificationMsg := buildEarthquakeNotificationMessage(response)

	// Notify IoT devices in affected areas
	if err := notifyAffectedDevices(app, mqtt, earthquakeLocation, magnitude, notificationMsg); err != nil {
		return fmt.Errorf("error notifying devices: %w", err)
	}

	// Notify users in affected areas
	if err := notifyAffectedUsers(app, earthquakeLocation, magnitude, notificationMsg, telegrambot); err != nil {
		return fmt.Errorf("error notifying users: %w", err)
	}

	return nil
}

// NotifyAffectedDevices sends notifications to IoT devices that would feel the earthquake
func notifyAffectedDevices(app core.App, mqtt mqtt.MQTTClient, earthquakeLocation ngitung.Location, magnitude float64, message string) error {
	// Get all IoT devices
	iotDevices, err := app.FindAllRecords("iot_device")
	if err != nil {
		return fmt.Errorf("failed to retrieve IoT devices: %w", err)
	}

	// Loop through each device and check if it would feel the earthquake
	for _, device := range iotDevices {
		var deviceInfo db.IotDevice
		deviceInfo.SetProxyRecord(device)

		// Get device location
		deviceLat, err := strconv.ParseFloat(deviceInfo.Lintang(), 64)
		if err != nil {
			continue // Skip devices with invalid latitude
		}

		deviceLon, err := strconv.ParseFloat(deviceInfo.Bujur(), 64)
		if err != nil {
			continue // Skip devices with invalid longitude
		}

		deviceLocation := ngitung.Location{
			Lat: deviceLat,
			Lon: deviceLon,
		}

		// Check if the device would feel the earthquake
		isFelt, _, _ := ngitung.IsWithinFeltRadius(earthquakeLocation, deviceLocation, magnitude)
		if !isFelt {
			continue
		}

		// Send notification to device
		if err := notify.SendMQTTNotification(mqtt, "device/"+deviceInfo.Id, message); err != nil {
			// Log error but continue processing other devices
			fmt.Printf("Failed to send notification to device %s: %v\n", deviceInfo.Id, err)
		}
	}

	return nil
}

// NotifyAffectedUsers sends notifications to users who would feel the earthquake
// NotifyAffectedUsers sends notifications to users who would feel the earthquake
func notifyAffectedUsers(app core.App, earthquakeLocation ngitung.Location, magnitude float64, message string, telegrambot *telegram.Bot) error {
	// Get all users with notification preferences
	usersToNotify, err := app.FindAllRecords("user_notify")
	if err != nil {
		return fmt.Errorf("failed to retrieve users: %w", err)
	}

	// Loop through each user and send notification if they would feel the earthquake
	for _, user := range usersToNotify {
		var userInfo db.UserNotify
		userInfo.SetProxyRecord(user)

		// Get user location
		userLat, err := strconv.ParseFloat(userInfo.Lintang(), 64)
		if err != nil {
			continue // Skip users with invalid latitude
		}

		userLon, err := strconv.ParseFloat(userInfo.Bujur(), 64)
		if err != nil {
			continue // Skip users with invalid longitude
		}

		userLocation := ngitung.Location{
			Lat: userLat,
			Lon: userLon,
		}

		// Check if the user would feel the earthquake
		isFelt, _, _ := ngitung.IsWithinFeltRadius(earthquakeLocation, userLocation, magnitude)
		if !isFelt {
			continue
		}

		// Send notification based on user preference
		if err := sendNotificationByType(userInfo, message, *telegrambot); err != nil {
			return fmt.Errorf("failed to send notification to user: %w", err)
		}
	}

	return nil
}

// BuildEarthquakeNotificationMessage creates a formatted notification message with earthquake details
func buildEarthquakeNotificationMessage(response domain.ResponseBmkgAPI) string {
	gempa := response.Infogempa.Gempa
	return "Earthquake Alert! Magnitude: " + gempa.Magnitude +
		"\nLocation: " + gempa.Lintang + ", " + gempa.Bujur +
		"\nDepth: " + gempa.Kedalaman +
		"\nTime: " + gempa.Tanggal + " " + gempa.Jam +
		"\nRegion: " + gempa.Wilayah
}

// SendNotificationByType sends a notification using the user's preferred notification method
// SendNotificationByType sends a notification using the user's preferred notification method
func sendNotificationByType(userInfo db.UserNotify, message string, telegrambot telegram.Bot) error {
	identifier := userInfo.Identifier()

	switch userInfo.Type() {
	case db.Telegram:
		chatID, err := strconv.ParseInt(identifier, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse telegram chat ID: %w", err)
		}
		return telegrambot.SendMessage(chatID, message)
	case db.Wa:
		return notify.SendWhatsAppNotification(identifier, message)
	default:
		return nil
	}
}
