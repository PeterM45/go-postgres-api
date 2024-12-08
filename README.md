# Go PostgreSQL Auth API

Simple, lightweight authentication system built with Go and PostgreSQL.

## Prerequisites

- Go 1.19+
- Docker and Docker Compose

## Quick Start

1. Clone and setup

```bash
git clone https://github.com/PeterM45/go-postgres-api
cd go-postgres-api
```

2. Create `.env` file

```env
DB_USER=admin
DB_PASS=password
DB_NAME=myapp
DB_HOST=localhost
DB_PORT=5433
PORT=8080
JWT_SECRET=your-secret-key
```

3. Run

```bash
# Start PostgreSQL
docker compose up -d

# Run API
go run cmd/main.go
```

## API Routes

### Public

- `POST /api/users` - Create user
- `POST /api/auth/login` - Login

### Protected (Requires JWT)

- `GET /api/users` - List users
- `GET /api/users/{id}` - Get user
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

## Test with curl

```bash
# Create user
curl -X POST http://localhost:8080/api/users \
-H "Content-Type: application/json" \
-d '{"username": "test", "email": "test@example.com", "password": "password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
-H "Content-Type: application/json" \
-d '{"email": "test@example.com", "password": "password123"}'
```
