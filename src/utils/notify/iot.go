package notify

import "bmkg/src/worker/mqtt"

// send to mqtt
// function to send notification to mqtt
func SendMQTTNotification(mqttclient mqtt.MQTTClient, topic, message string) error {
	// Implement the logic to send a notification to MQTT
	// This is a placeholder implementation
	return nil
}
