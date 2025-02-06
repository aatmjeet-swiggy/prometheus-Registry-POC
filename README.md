# Prometheus and gRPC Prometheus POC

This project is a proof of concept (POC) to demonstrate the working of Prometheus registry and `grpc_prometheus` with a gRPC server.

## Project Structure
```
prometheus-grafana-stack 
├── prometheus 
│ ├── prometheus.yml 
│ └── Dockerfile 
├── grafana 
│ ├── grafana.ini 
│ └── Dockerfile 
├── docker-compose.yml 
├── proto 
│ ├── hello.proto 
│ ├── hello.pb.go 
│ └── hello_grpc.pb.go 
├── metrics.go 
└── README.md
```

## Prerequisites

- Docker
- Docker Compose
- Go

## Setting Up Docker Dependencies

To set up the Docker dependencies, run:

```sh
docker-compose up -d
```
This will start Prometheus on `http://localhost:9090` and Grafana on `http://localhost:3000`

## Running the POC
To run the POC, execute the following command:

```sh
go run metrics.go
```
This will start the gRPC server and expose metrics on `http://localhost:9091/metrics`.

## Accessing Metrics
- Prometheus: Access Prometheus at `http://localhost:9090`.
- Grafana: Access Grafana at `http://localhost:3000` (default credentials: `admin/admin`).