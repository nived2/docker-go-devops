# Simple Web Server (Project 1)

A basic Go web server containerized with Docker. This project introduces fundamental concepts of both Go and Docker.

## Features
- Basic HTTP server with multiple endpoints
- Request logging
- Health check endpoint
- Environment variable configuration
- Multi-stage Docker build

## Prerequisites
- Go 1.16 or later
- Docker

## Building and Running

### Local Development
```bash
# Run the server locally
go run main.go
```

### Docker Build and Run
```bash
# Build the Docker image
docker build -t simple-webserver .

# Run the container
docker run -p 8080:8080 simple-webserver
```

## Testing the Application
Once running, you can test the following endpoints:

1. Main endpoint:
```bash
curl http://localhost:8080
```

2. Health check:
```bash
curl http://localhost:8080/health
```

## Learning Objectives
- Basic Go HTTP server implementation
- Docker multi-stage builds
- Container port mapping
- Environment variables in containers
- Docker image optimization

## Docker Commands Learned
- `docker build`: Build an image from a Dockerfile
- `docker run`: Create and start a container
- `docker ps`: List running containers
- `docker logs`: View container logs
- `docker stop`: Stop a running container

## Next Steps
After completing this project, you should understand:
- Basic Go web server structure
- Docker image building process
- Container runtime basics
- Port mapping in Docker
- Basic Docker commands
