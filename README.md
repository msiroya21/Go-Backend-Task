# Go Backend Task

This is a Go backend REST API project managing users with CRUD operations, using PostgreSQL as the database. The project demonstrates clean architecture, PostgreSQL integration with SQLC, Uber Zap logging, request validation, pagination, and middleware for request tracing.

---

## Features

- User CRUD endpoints (`/users`, `/users/:id`)
- Paginated user listing
- Request validation with structured error responses
- Logging using Uber Zap
- Middleware for request ID injection and request logging
- Unit testing for utility functions like age calculation

---

## Tech Stack

- Go (Golang)
- Fiber web framework
- PostgreSQL database
- SQLC for type-safe DB queries
- Uber Zap for structured logging
- Validator package for input validation

---

## Setup Instructions

### Prerequisites

- Go installed (version X.X or higher)
- PostgreSQL installed and running
- Git CLI installed (optional)

### Environment Variables

Create a `.env` file and define the following:

DATABASE_URL=postgresql://postgres:<your_password>@localhost:5432/go_backend_task?sslmode=disable

text

Replace `<your_password>` with your actual DB password.

### Database Setup

Run the PostgreSQL SQL migration scripts located at `db/migrations` to setup tables.

### Running the Application

go run cmd/server/main.go

text

The server listens on port `3000`.

---

## API Usage

### List users (paginated)

`GET /users?page=1&limit=10`

### Get user by ID

`GET /users/:id`

### Create user

`POST /users`

Body:

{
"name": "Alice",
"dob": "2000-01-01"
}

text

### Update user

`PUT /users/:id`

### Delete user

`DELETE /users/:id`

---

## Testing

Run unit tests with:

go test ./internal/service

text

---

## Logging

Request logs and error logs are outputted to the console using Uber Zap.

---

## Docker (optional)

Dockerfile and docker-compose.yml included for containerized deployment.

To build and run containers:

docker compose up --build

text

Note: Docker usage requires Docker Desktop and a compatible environment.

---

## Known Issues & Notes

- Sensitive data like database password should be set using environment variables, not hardcoded.
- Docker image pull errors may occur on some network setups.
- Pagination is implemented on the users listing endpoint.