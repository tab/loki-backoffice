# Loki-Backoffice

Administration backend application for the Loki SSO service ecosystem.

## Key Features

- Create and update roles, permissions, scopes, users and tokens
- Comprehensive logging and telemetry support (OpenTelemetry) for easier monitoring and tracing
- Easily integrate into a microservices architecture

## Prerequisites

Before starting this application, you must have the loki-infrastructure running:

```sh
git clone git@github.com/tab/loki-infrastructure.git
cd loki-infrastructure

docker-compose up
```

## Setup and Configuration

### Environment Variables

Use `.env` files (e.g., `.env.development`) or provide environment variables for:

- `DATABASE_DSN` for PostgreSQL
- `TELEMETRY_URI` for OpenTelemetry
- `GRPC_ADDRESS` for communication with the main Loki service

### Generate mTLS Client Certificates

#### JWT Signing Keys

```sh
mkdir -p certs/jwt

# Copy public key from Loki service
cp ../loki/certs/jwt/public.key ./certs/jwt/
```

#### mTLS Certificates

For secure communication with the Loki service, you need to generate client certificates for mTLS:

```sh
# Create directory
mkdir -p certs

# Copy CA from Loki service
cp ../loki/certs/ca.pem ./certs/

# Generate Client Certificate
openssl genrsa -out certs/client.key 4096
openssl req -new -key certs/client.key -out certs/client.csr -config <(
cat <<-EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
distinguished_name = dn

[dn]
CN = loki-backoffice
EOF
)

openssl x509 -req -in certs/client.csr -CA certs/ca.pem -CAkey certs/ca.key -CAcreateserial -out certs/client.pem -days 825 -sha256
```

For more detailed information on certificates, see [Documentation](https://tab.github.io/loki).

### Database Migrations

Run the following command to apply database migrations:

```sh
GO_ENV=development make db:drop db:create db:migrate
```

### Run application

```sh
docker-compose build
docker-compose up
```

### Check health status

```sh
curl -X GET http://localhost:8081/live
```

```sh
curl -X GET http://localhost:8081/ready
```

## Documentation

[Documentation](https://tab.github.io/loki)

## API Documentation

Swagger file is available at [api/swagger.yaml](https://github.com/tab/loki-backoffice/blob/master/api/swagger.yaml)

## Related Repositories

The Loki ecosystem consists of the following repositories:

- [Loki](https://github.com/tab/loki) - Loki SSO & RBAC application
- [Loki Infrastructure](https://github.com/tab/loki-infrastructure) - Infrastructure setup for the Loki ecosystem
- [Loki Proto](https://github.com/tab/loki-proto) - Protocol buffer definitions
- [Loki Frontend](https://github.com/tab/loki-frontend) - Frontend application
- [Smart-ID Client](https://github.com/tab/smartid) - Smart-ID client used for authentication
- [Mobile-ID Client](https://github.com/tab/mobileid) - Mobile-ID client used for authentication

## Architecture

The application follows a layered architecture and clean code principles:

- **cmd/backoffice**: Application entry point
- **internal/app**: Core application logic, including services, controllers, repositories, and DTOs
- **internal/config**: Configuration loading and setup, server startup, middleware, router initialization, and telemetry configuration
- **pkg**: Reusable utilities

## License

Distributed under the MIT License. See `LICENSE` for more information.
