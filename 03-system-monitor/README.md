# System Monitor

A comprehensive system monitoring application built with Go that provides real-time metrics about system resources including CPU, Memory, Disk usage, and running processes.

## Features

- Real-time system metrics monitoring
- CPU usage and core count information
- Memory usage statistics (total, used, free, usage percentage)
- Disk usage information for all mounted partitions
- Process monitoring with top CPU and memory consuming processes
- RESTful API endpoints for accessing metrics
- Docker containerization support

## API Endpoints

- `GET /` - Home page with basic information
- `GET /metrics` - Complete system metrics including CPU, Memory, and Disk usage
- `GET /processes` - List of top 10 processes by CPU usage
- `GET /health` - Health check endpoint

## Building and Running

### Local Development

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

### Using Docker

1. Build the Docker image:
   ```bash
   docker build -t system-monitor .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 system-monitor
   ```

## Environment Variables

- `PORT` - Server port (default: 8080)

## API Response Examples

### System Metrics (/metrics)
```json
{
  "timestamp": "2023-11-01T12:00:00Z",
  "cpu": {
    "usage": 25.5,
    "core_count": 8,
    "load_average": [1.5, 1.2, 1.0]
  },
  "memory": {
    "total": 16000000000,
    "used": 8000000000,
    "free": 8000000000,
    "usage_percentage": 50.0
  },
  "disk": [
    {
      "device": "/dev/sda1",
      "mount_point": "/",
      "total": 250000000000,
      "used": 100000000000,
      "free": 150000000000,
      "usage_percentage": 40.0
    }
  ],
  "process_count": 100,
  "top_processes": [
    {
      "pid": 1234,
      "name": "chrome",
      "cpu_percent": 10.5,
      "memory_usage": 500000000
    }
  ]
}
```

## Dependencies

- github.com/shirou/gopsutil/v3 - System metrics collection
- Standard Go libraries for HTTP server and JSON handling

## Notes

- The application requires appropriate permissions to access system metrics
- When running in Docker, some system metrics might be container-specific
- Process monitoring might require elevated privileges on some systems
