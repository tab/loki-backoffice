# Loki-Backoffice

Loki-backoffice backend application

## Key Features

- Create and update roles, permissions, scopes, users and tokens
- Comprehensive logging and telemetry support (OpenTelemetry) for easier monitoring and tracing
- Easily integrate into a microservices architecture

## Prerequisites

Before starting this application, you must have the loki-infrastructure running:

```sh
git clone git@github.com/tab/loki-infrastructure.git
cd loki-infrastructure
```

```sh
docker-compose up
```

## Setup and Configuration

**Environment Variables**:

Use `.env` files (e.g., `.env.development`) or provide environment variables for:

- `DATABASE_DSN` for PostgreSQL
- `TELEMETRY_URI` for OpenTelemetry

**Database Migrations**:

Run the following command to apply database migrations:

```sh
GO_ENV=development make db:drop db:create db:migrate
```

**Run the Services**:

```sh
docker-compose build
docker-compose up
```

**Check health status**:

```sh
curl -X GET http://localhost:8081/health
```

## Architecture

The application follows a layered architecture and clean code principles:

- **cmd/backoffice**: Application entry point
- **internal/app**: Core application logic, including services, controllers, repositories, and DTOs
- **internal/config**: Configuration loading and setup, server startup, middleware, router initialization, and telemetry configuration
- **pkg**: Reusable utilities

## License

Distributed under the MIT License. See `LICENSE` for more information.
