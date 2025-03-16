# Installation

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
