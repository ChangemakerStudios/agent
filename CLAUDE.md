# Swarmpit Agent

Docker Swarm monitoring agent that collects host/container stats and Docker events, forwarding them to a Swarmpit application server.

## Build & Run

```bash
# Build locally
CGO_ENABLED=0 go build -o agent .

# Build Docker image (multi-arch supported: amd64, arm64, armv7, armv6)
docker build --build-arg GOARCH=amd64 -t swarmpit/agent:latest .

# Run
docker run -d --name agent -v /var/run/docker.sock:/var/run/docker.sock swarmpit/agent:latest
```

## Architecture

- `main.go` — Entry point: Docker client init, health check, starts event/stats goroutines, HTTP server on `:8080`
- `router.go` — gorilla/mux router with logging middleware; gzip middleware available but disabled
- `handler.go` — `GET /` returns agent config; `GET /logs/{container}?since=` returns container logs
- `setup/args.go` — Configuration from environment variables
- `swarmpit/client.go` — HTTP client posting events to Swarmpit app + startup health check loop
- `swarmpit/task/event.go` — Streams Docker events and forwards to Swarmpit
- `swarmpit/task/stats.go` — Periodic collection of CPU, memory, disk, and per-container stats

## Configuration (env vars)

| Variable | Default | Description |
|----------|---------|-------------|
| `STATS_FREQUENCY` | `30` | Stats collection interval in seconds |
| `EVENT_ENDPOINT` | `http://app:8080/events` | Swarmpit event ingestion URL |
| `HEALTH_CHECK_ENDPOINT` | `http://app:8080/version` | Swarmpit health check URL |
| `DEBUG_EVENT` | `false` | Log Docker events to stdout |
| `DEBUG_STATS` | `false` | Log host stats to stdout |

## Tech Stack

- Go 1.12, Docker Engine SDK 19.03, gorilla/mux, gopsutil
- Multi-stage Dockerfile producing a `scratch`-based image
- CI/CD via Travis CI with multi-arch manifest pushes to Docker Hub

## Conventions

- Package `setup` owns all config/env parsing
- Package `swarmpit` owns the HTTP client for Swarmpit communication
- Package `swarmpit/task` owns background collector goroutines
- Logging uses `log.Printf` with level prefixes: `INFO:`, `ERROR:`, `DEBUG:`
