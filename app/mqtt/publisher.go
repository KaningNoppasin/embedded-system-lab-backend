package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const defaultTransactionStatusTopic = "esp32/transaction/status"

type Publisher struct {
	client                 mqtt.Client
	transactionStatusTopic string
}

type TransactionStatusPayload struct {
	Status      string    `json:"status"`
	PublishedAt time.Time `json:"published_at"`
}

func NewPublisher() (*Publisher, error) {
	brokerURL := getEnv("MQTT_BROKER_URL", defaultBrokerURL)
	clientID := getEnv("MQTT_PUBLISHER_CLIENT_ID", fmt.Sprintf("embedded-lab-api-publisher-%d", time.Now().UnixNano()))
	transactionStatusTopic := getEnv("MQTT_TRANSACTION_STATUS_TOPIC", defaultTransactionStatusTopic)

	options := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second)

	options.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("mqtt publisher connection lost: %v", err)
	})

	client := mqtt.NewClient(options)
	token := client.Connect()
	if !token.WaitTimeout(connectTimeout) {
		return nil, fmt.Errorf("mqtt publisher connect timeout after %s", connectTimeout)
	}
	if err := token.Error(); err != nil {
		return nil, err
	}

	log.Printf("mqtt publisher connected to broker %s", brokerURL)

	return &Publisher{
		client:                 client,
		transactionStatusTopic: transactionStatusTopic,
	}, nil
}

func (p *Publisher) PublishTransactionStatus(status string) (*TransactionStatusPayload, error) {
	payload := &TransactionStatusPayload{
		Status:      status,
		PublishedAt: time.Now(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	token := p.client.Publish(p.transactionStatusTopic, 0, false, body)
	if !token.WaitTimeout(connectTimeout) {
		return nil, fmt.Errorf("mqtt publish timeout after %s", connectTimeout)
	}
	if err := token.Error(); err != nil {
		return nil, err
	}

	log.Printf("mqtt transaction status published: topic=%s status=%s", p.transactionStatusTopic, status)
	return payload, nil
}

func (p *Publisher) TransactionStatusTopic() string {
	if p == nil || p.transactionStatusTopic == "" {
		return defaultTransactionStatusTopic
	}

	return p.transactionStatusTopic
}

func (p *Publisher) Close() {
	if p == nil || p.client == nil || !p.client.IsConnected() {
		return
	}

	p.client.Disconnect(250)
}
