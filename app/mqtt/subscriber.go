package mqtt

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	defaultBrokerURL = "tcp://localhost:1883"
	defaultTopic     = "esp32/temperature"
	connectTimeout   = 10 * time.Second
)

type Subscriber struct {
	client mqtt.Client
}

type TemperatureStore interface {
	StoreTemperature(topic string, value float64) error
}

func NewTemperatureSubscriber(store TemperatureStore) (*Subscriber, error) {
	brokerURL := getEnv("MQTT_BROKER_URL", defaultBrokerURL)
	topic := getEnv("MQTT_TOPIC", defaultTopic)
	clientID := getEnv("MQTT_CLIENT_ID", fmt.Sprintf("embedded-lab-api-%d", time.Now().UnixNano()))

	options := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second)

	options.OnConnect = func(client mqtt.Client) {
		token := client.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
			payload := strings.TrimSpace(string(msg.Payload()))
			temperature, err := strconv.ParseFloat(payload, 64)
			if err != nil {
				log.Printf("mqtt message parse failed for topic %s: payload=%q err=%v", msg.Topic(), payload, err)
				return
			}

			if err := store.StoreTemperature(msg.Topic(), temperature); err != nil {
				log.Printf("influxdb write failed for topic %s value=%f: %v", msg.Topic(), temperature, err)
				return
			}

			log.Printf("mqtt temperature received: topic=%s value=%f", msg.Topic(), temperature)
		})

		if token.WaitTimeout(connectTimeout) && token.Error() != nil {
			log.Printf("mqtt subscribe failed for topic %s: %v", topic, token.Error())
			return
		}

		log.Printf("mqtt subscribed to topic %s", topic)
	}

	options.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("mqtt connection lost: %v", err)
	})

	client := mqtt.NewClient(options)
	token := client.Connect()
	if !token.WaitTimeout(connectTimeout) {
		return nil, fmt.Errorf("mqtt connect timeout after %s", connectTimeout)
	}
	if err := token.Error(); err != nil {
		return nil, err
	}

	log.Printf("mqtt connected to broker %s", brokerURL)
	return &Subscriber{client: client}, nil
}

func (s *Subscriber) Close() {
	if s == nil || s.client == nil || !s.client.IsConnected() {
		return
	}

	s.client.Disconnect(250)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
