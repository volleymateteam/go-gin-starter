# Go-Gin-Starter

A RESTful API backend built with Go and Gin framework providing user management, authentication, and sports management features.

## Features

- JWT-based authentication
- Role-based access control (RBAC)
- User management with permissions
- Teams, seasons, and matches management
- Audit logging for admin actions
- API versioning (v1)
- Swagger documentation
- Health check endpoints
- Rate limiting
- Structured logging
- CORS support
- Docker ready

## Prerequisites

- Go 1.19+
- PostgreSQL 12+
- Docker & Docker Compose (optional)

## Quick Start

### Local Development

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/go-gin-starter.git
   cd go-gin-starter
   ```

2. Set up environment variables:

   ```bash
   cp .env.example .env
   # Edit .env with your own values
   ```

3. Start PostgreSQL database

4. Run the application:

   ```bash
   go run main.go
   ```

### Docker Development

```bash
docker-compose up
```

## Environment Variables

Key environment variables:

- `ENV` - Environment (development/production)
- `PORT` - Server port
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Database connection
- `JWT_SECRET` - Secret for JWT generation
- `ALLOWED_ORIGINS` - CORS origins (comma separated)

## API Documentation

Access Swagger documentation at:

```bash
http://localhost:8080/api/v1/swagger/index.html
```

## Health Checks

- Basic health: `GET /health`
- Detailed health: `GET /health/details`
- Kubernetes readiness: `GET /readiness`
- Kubernetes liveness: `GET /liveness`

## AWS Deployment

### Preparation

1. Configure AWS CLI and authenticate
2. Create ECR repository for container
3. Create RDS PostgreSQL instance

### Deployment Options

#### ECS Fargate

1. Push Docker image to ECR:

   ```bash
   aws ecr get-login-password | docker login --username AWS --password-stdin [AWS-ACCOUNT-ID].dkr.ecr.[REGION].amazonaws.com
   docker build -t go-gin-starter .
   docker tag go-gin-starter:latest [AWS-ACCOUNT-ID].dkr.ecr.[REGION].amazonaws.com/go-gin-starter:latest
   docker push [AWS-ACCOUNT-ID].dkr.ecr.[REGION].amazonaws.com/go-gin-starter:latest
   ```

2. Create ECS cluster, task definition and service

#### EC2 with Docker

1. Launch EC2 instance with Docker installed
2. Pull Docker image and run with environment variables

## Project Structure

```bash
├── config/         # Configuration management
├── controllers/    # Request handlers
├── database/       # Database connection
├── dto/            # Data transfer objects
├── middleware/     # HTTP middlewares
├── models/         # Data models
├── pkg/            # Shared packages
├── repositories/   # Database interactions
├── routes/         # API routes
├── services/       # Business logic
└── docs/           # Swagger docs
```

## Development

- Run tests: `go test ./...`
- Format code: `go fmt ./...`
- Lint code: `golangci-lint run`
- Generate swagger: `swag init -g main.go`
