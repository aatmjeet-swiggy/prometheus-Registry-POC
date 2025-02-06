# prometheus-grafana-stack/prometheus-grafana-stack/README.md

# Prometheus and Grafana Stack

This project sets up a dummy Prometheus and Grafana stack for testing purposes using Docker. It includes the necessary configuration files and Dockerfiles to build and run the services.

## Project Structure

```
prometheus-grafana-stack
├── prometheus
│   ├── prometheus.yml
│   └── Dockerfile
├── grafana
│   ├── grafana.ini
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Building and Running the Stack

1. Clone the repository:

   ```
   git clone <repository-url>
   cd prometheus-grafana-stack
   ```

2. Build and start the services using Docker Compose:

   ```
   docker-compose up --build
   ```

3. Access Prometheus at `http://localhost:9090` and Grafana at `http://localhost:3000`.

### Configuration

- **Prometheus**: The configuration file is located at `prometheus/prometheus.yml`. You can modify the scrape configurations and other settings as needed.
- **Grafana**: The configuration file is located at `grafana/grafana.ini`. Adjust the settings for the Grafana server and database as required.

### Stopping the Stack

To stop the services, run:

```
docker-compose down
```

## License

This project is licensed under the MIT License.