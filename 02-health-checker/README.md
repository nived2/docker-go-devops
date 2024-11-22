# Health Checker Service

A comprehensive health checking service that monitors multiple components of a distributed system. This project demonstrates the use of Docker Compose, multi-stage builds, and environment variables.

## Features

- Basic health check endpoint (`/health`)
- Detailed health status for all services (`/health/detailed`)
- Metrics endpoint (`/metrics`)
- Prometheus integration for monitoring
- Multi-service architecture with Docker Compose
- Environment variable configuration

## Services Monitored

- Database (PostgreSQL)
- Cache (Redis)
- API (Nginx)
- Monitoring (Prometheus)

## Prerequisites

- Docker
- Docker Compose
- Go 1.16+ (for local development)

## Running the Application

1. Start all services using Docker Compose:
```bash
docker-compose up --build
```

2. Access the endpoints:
- Health Check: http://localhost:8080/health
- Detailed Health: http://localhost:8080/health/detailed
- Metrics: http://localhost:8080/metrics
- Prometheus: http://localhost:9090

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection URL
- `CACHE_URL`: Redis connection URL
- `API_URL`: API service URL
- `MONITORING_URL`: Prometheus URL

## Project Structure

```
.
├── main.go              # Main application code
├── Dockerfile           # Multi-stage build configuration
├── docker-compose.yml   # Multi-service orchestration
├── prometheus.yml       # Prometheus configuration
└── README.md           # This file
```

## Learning Objectives

1. Docker Compose
   - Multi-service orchestration
   - Service dependencies
   - Environment variables
   - Volume mounting

2. Go Development
   - HTTP server with multiple endpoints
   - JSON response handling
   - Environment variable configuration
   - Health check implementation

3. Monitoring
   - Basic Prometheus setup
   - Metrics exposure
   - Service health monitoring

4. Best Practices
   - Multi-stage Docker builds
   - Environment variable usage
   - Service discovery
   - Logging and monitoring

## Testing

To test individual endpoints:

```bash
# Basic health check
curl http://localhost:8080/health

# Detailed health status
curl http://localhost:8080/health/detailed

# Metrics
curl http://localhost:8080/metrics
```

## Next Steps

1. Add real service health checks
2. Implement authentication
3. Add more metrics
4. Create Grafana dashboards
5. Add alerting configuration
