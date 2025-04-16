package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient defines the interface for MQTT operations
type MQTTClient interface {
	Connect() error
	Disconnect()
	Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error
	Publish(topic string, qos byte, retained bool, payload interface{}) error
	IsConnected() bool
}

// Config holds MQTT client configuration
type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
	QoS      byte
}

// mqttClient implements MQTTClient interface
type mqttClient struct {
	client mqtt.Client
	config Config
}

// NewMQTTClient creates a new MQTT client instance
func NewMQTTClient(config Config) MQTTClient {
	opts := mqtt.NewClientOptions().
		AddBroker(config.Broker).
		SetClientID(config.ClientID).
		SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	return &mqttClient{
		client: client,
		config: config,
	}
}

func (m *mqttClient) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *mqttClient) Disconnect() {
	m.client.Disconnect(250)
}

func (m *mqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if token := m.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *mqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if token := m.client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *mqttClient) IsConnected() bool {
	return m.client.IsConnected()
}
