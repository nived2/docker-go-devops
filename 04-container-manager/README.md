# Container Manager

A Docker container management tool built with Go that provides a simple interface for managing Docker containers, images, and networks.

## Features

- Container lifecycle management (create, start, stop, remove)
- Image management (pull, list, remove)
- Network management (create, list, connect, disconnect)
- Container monitoring and stats
- Custom Docker commands execution
- RESTful API interface

## Prerequisites

- Go 1.19 or later
- Docker Engine
- Docker API access

## Installation

1. Clone the repository
2. Build the application:
   ```bash
   go build -o container-manager
   ```
3. Run the application:
   ```bash
   ./container-manager
   ```

## API Endpoints

### Containers
- `GET /containers` - List all containers
- `POST /containers` - Create a new container
- `GET /containers/{id}` - Get container details
- `POST /containers/{id}/start` - Start a container
- `POST /containers/{id}/stop` - Stop a container
- `DELETE /containers/{id}` - Remove a container
- `GET /containers/{id}/stats` - Get container statistics

### Images
- `GET /images` - List all images
- `POST /images` - Pull a new image
- `DELETE /images/{id}` - Remove an image

### Networks
- `GET /networks` - List all networks
- `POST /networks` - Create a new network
- `GET /networks/{id}` - Get network details
- `DELETE /networks/{id}` - Remove a network
- `POST /networks/{id}/connect` - Connect a container to a network
- `POST /networks/{id}/disconnect` - Disconnect a container from a network

## Configuration

The application can be configured using environment variables:

- `DOCKER_HOST` - Docker daemon socket (default: "unix:///var/run/docker.sock")
- `API_PORT` - Port for the HTTP API (default: 8080)
- `LOG_LEVEL` - Logging level (default: "info")

## Development

This project uses:
- [Go Docker SDK](https://pkg.go.dev/github.com/docker/docker/client)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Logrus](https://github.com/sirupsen/logrus) for logging

## License

MIT License
