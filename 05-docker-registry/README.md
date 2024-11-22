# Docker Registry Manager

A Go-based application for managing a private Docker registry with advanced features for image management, access control, and monitoring.

## Features

### Core Features
- Private Docker registry setup and management
- Image push/pull operations
- Tag management and versioning
- Access control and authentication
- Storage backend integration
- Registry metrics and monitoring
- REST API for registry operations

### API Endpoints

#### Registry Management
- `GET /api/v1/registry/health` - Check registry health
- `GET /api/v1/registry/info` - Get registry information
- `GET /api/v1/registry/metrics` - Get registry metrics

#### Image Management
- `GET /api/v1/images` - List all images
- `GET /api/v1/images/{name}` - Get image details
- `GET /api/v1/images/{name}/tags` - List image tags
- `DELETE /api/v1/images/{name}` - Delete an image
- `DELETE /api/v1/images/{name}/tags/{tag}` - Delete specific image tag

#### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/token` - Get authentication token
- `GET /api/v1/auth/verify` - Verify authentication token

#### User Management
- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/{username}` - Update user
- `DELETE /api/v1/users/{username}` - Delete user

## Technology Stack
- Go 1.19
- Docker Registry API v2
- JWT for authentication
- Prometheus for metrics
- Redis for caching
- PostgreSQL for user management

## Getting Started

### Prerequisites
- Go 1.19 or higher
- Docker
- Docker Compose
- PostgreSQL
- Redis

### Installation
1. Clone the repository
2. Install dependencies: `go mod download`
3. Configure environment variables
4. Run the application: `go run main.go`

### Docker Compose Setup
```bash
docker-compose up -d
```

## Configuration
Environment variables:
- `REGISTRY_URL`: Docker registry URL
- `REGISTRY_PORT`: Registry port (default: 5000)
- `API_PORT`: API server port (default: 8080)
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_NAME`: Database name
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `REDIS_HOST`: Redis host
- `REDIS_PORT`: Redis port
- `JWT_SECRET`: JWT signing secret

## API Documentation
Detailed API documentation is available at `/docs` endpoint when running the server.

## Security
- TLS encryption for all communications
- JWT-based authentication
- Role-based access control
- Secure password hashing
- Rate limiting

## Contributing
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License
This project is licensed under the MIT License.
