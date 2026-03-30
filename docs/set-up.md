# Project Setup

This project contains:

- A Go API built with Fiber
- MongoDB for user data
- Mosquitto as the MQTT broker
- InfluxDB for temperature time-series data
- Grafana for dashboards
- mongo-express for MongoDB inspection

## Prerequisites

Install these first:

- Go `1.25.6` or compatible
- Docker
- Docker Compose
- `make`

## Services and Ports

Default ports used by the project:

- API: `8080`
- MongoDB: `27017`
- mongo-express: `8081`
- Mosquitto MQTT broker: `1883`
- InfluxDB: `8086`
- Grafana: `3000`

## Environment Variables

You can run the API with these environment variables:

```env
APP_PORT=8080
MONGO_URI=mongodb://admin:admin@localhost:27017
MONGO_DATABASE=embedded_lab
MONGO_USER_COLLECTION=users
MQTT_BROKER_URL=tcp://localhost:1883
MQTT_TOPIC=esp32/temperature
INFLUXDB_URL=http://localhost:8086
INFLUXDB_TOKEN=embedded-lab-token
INFLUXDB_ORG=embedded-lab
INFLUXDB_BUCKET=telemetry
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
```

You can create a local env file from `.env.example`, then add the MQTT and InfluxDB values if needed.

## Development Setup

Start the infrastructure:

```bash
make dev-up
```

This starts:

- Mosquitto
- InfluxDB
- Grafana
- MongoDB
- mongo-express

Run the API locally:

```bash
go run ./cmd
```

If port `8080` is already in use, run the API on another port:

```bash
APP_PORT=8082 go run ./cmd
```

Stop development services:

```bash
make dev-down
```

View development logs:

```bash
make dev-logs
```

## Production-Style Setup

Run the full stack, including the API in Docker:

```bash
make prod-up
```

Stop it:

```bash
make prod-down
```

View logs:

```bash
make prod-logs
```

## End-to-End Verification

1. Start development services:

```bash
make dev-up
```

2. Run the API locally:

```bash
go run ./cmd
```

3. Publish a temperature message to MQTT:

```bash
docker exec -it embedded-lab-mosquitto mosquitto_pub -h localhost -t esp32/temperature -m 27.5
```

4. Confirm the API logs show a received temperature message.

5. Open InfluxDB UI in a browser:

```text
http://localhost:8086
```

Use these initial credentials:

- Username: `admin`
- Password: `adminadmin123`
- Token: `embedded-lab-token`
- Organization: `embedded-lab`
- Bucket: `telemetry`

6. Open Grafana:

```text
http://localhost:3000
```

Use these credentials:

- Username: `admin`
- Password: `adminadmin123`

Grafana is pre-provisioned with the InfluxDB datasource and the `Embedded Lab Overview` dashboard.

## MongoDB Access

Open mongo-express:

```text
http://localhost:8081
```

Open a Mongo shell inside the container:

```bash
make mongo-shell
```

## Build the API Image

Build the Docker image manually:

```bash
make build-api
```

Rebuild without cache:

```bash
make rebuild-api
```

## Notes

- The API subscribes to MQTT topic `esp32/temperature`.
- The MQTT payload must be a float value such as `25.4`.
- Each received temperature is written to InfluxDB measurement `temperature`.
- INA219 telemetry is written to InfluxDB measurement `ina219`.
- Grafana connects to InfluxDB automatically through Docker Compose.
- MongoDB is still used for user-related data.
- `POST /notifications/discord` sends the request body message to `DISCORD_WEBHOOK_URL`.
