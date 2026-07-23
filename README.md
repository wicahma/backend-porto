# API Server

A lightweight, high-performance modular monolith backend built in Go. Designed for minimal memory footprint and fast execution, this API leverages the Go 1.22+ standard library for routing and avoids heavy frameworks.

## Tech Stack

* **Language:** Go 1.22+
* **Routing:** Standard Library `net/http` (`http.ServeMux`)
* **Database:** PostgreSQL (using `jackc/pgx/v5` and `sqlc`)
* **Caching:** Redis (`redis/go-redis/v9`)
* **Storage:** AWS S3 (`aws-sdk-go-v2`)
* **Image Processing:** Standard `image` + `disintegration/imaging`

## Architecture: Modular Monolith

This project follows a strict modular monolith architecture.
* **Domain Isolation:** Features are separated into vertical slices inside the `internal/` directory (e.g., `internal/users`).
* **No Cross-Domain Joins:** Domains do not share database queries. They communicate strictly through Go interfaces.
* **Composition Root:** All dependency injection, configuration loading, and routing are wired together in `cmd/api/main.go`.

## Prerequisites

* Go 1.22 or higher
* sqlc CLI (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
* PostgreSQL and Redis

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd <project-directory>

```

2. **Set up environment variables:**
Copy the example environment file and update the credentials.
```bash
cp .env.example .env

```


3. **Download dependencies:**
```bash
go mod download

```


4. **Run the server:**
```bash
go run ./cmd/api/main.go

```


The server will start at `http://localhost:8080`.

## Database and sqlc Workflow

We do not use a runtime ORM. Instead, we write raw SQL and compile it into type-safe Go code using `sqlc`.

1. Add your table definitions to `internal/<domain>/schema.sql`.
2. Add your queries to `internal/<domain>/queries.sql`.
3. Run the code generator:
```bash
sqlc generate

```


4. Do not manually edit the generated `*.sql.go` files.

## Deployment

This API is designed to run efficiently in constrained micro-server environments. A standard Dockerfile will compile the application into a minimal binary.

To build the image locally:

```bash
docker build -t api-server .

```

For deployment, point your application service to this repository and expose the target port (default `8080`). Ensure all environment variables (Database DSN, Redis URL, S3 Keys) are mapped in your deployment runtime.
