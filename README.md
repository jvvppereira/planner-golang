# Trip Planner API (Golang)

[![Maintainability](https://qlty.sh/gh/jvvppereira/projects/planner-golang/maintainability.svg)](https://qlty.sh/gh/jvvppereira/projects/planner-golang)

A robust backend service for a trip planning application, built with Go. This project provides a RESTful API to manage trips, participants, and activities, featuring email notifications and automated code generation from OpenAPI specifications.

## 🚀 Technologies and Libraries

This project leverages modern Go libraries and tools to ensure performance, reliability, and developer productivity:

### Core Framework & Routing
- **[Go (Golang)](https://golang.org/):** Version 1.27.0.
- **[Chi (v5)](https://github.com/go-chi/chi):** A lightweight, idiomatic, and composable router for building HTTP services.
- **[Chi Render](https://github.com/go-chi/render):** Managed content-type negotiation and response rendering.

### Database & Persistence
- **[PostgreSQL](https://www.postgresql.org/):** The primary relational database.
- **[pgx (v5)](https://github.com/jackc/pgx):** A high-performance PostgreSQL driver and toolkit for Go.
- **[sqlc](https://sqlc.dev/):** Compile-safe SQL generation for Type-safe Go.
- **[tern](https://github.com/jackc/tern):** A standalone database migration tool for PostgreSQL.
- **[Docker / Docker Compose](https://www.docker.com/):** For local orchestration of the database, Mailpit, and the application.

### API Specification & Tooling
- **[OpenAPI / Swagger](https://www.openapis.org/):** Used to define the API contract.
- **[goapi-gen](https://github.com/discord-gophers/goapi-gen):** OpenAPI 3 code generator for Go, ensuring the implementation stays in sync with the spec.
- **[kin-openapi](https://github.com/getkin/kin-openapi):** Tools for loading and validating OpenAPI specifications.

### Utilities & Quality of Life
- **[Uber-zap](https://github.com/uber-go/zap):** Blazing fast, structured, leveled logging.
- **[Go-playground Validator (v10)](https://github.com/go-playground/validator):** For robust request data validation.
- **[Google UUID](https://github.com/google/uuid):** Unique identifier generation.
- **[Go-mail](https://github.com/wneessen/go-mail):** Versatile email sending library.
- **[gutils](https://github.com/phenpessoa/gutils):** Shared utility functions for common tasks.

---

## 🏗️ Project Structure

- `cmd/planner-golang`: Entry point of the application.
- `internal/api`: API handlers and logic implementation.
  - `api.go`: Core types, interfaces, and constructor.
  - `trip_helpers.go`: Shared helper functions for trip operations.
  - `*_get.go`, `*_create.go`, `*_update.go`: Individual endpoint handlers.
- `internal/api/spec`: OpenAPI 3 specification and generated code.
- `internal/pgstore`: Database interactions and queries (sqlc generated).
- `internal/mailer`: Email delivery logic (configured for Mailpit).

---

## 🛠️ Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Go](https://golang.org/doc/install) (optional, if running locally without Docker)

### Running with Docker

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd planner-golang
   ```

2. **Configure Environment Variables:**
   Create a `.env` file in the root directory and fill in the necessary database credentials:
   ```env
   PLANNER_DATABASE_NAME=planner
   PLANNER_DATABASE_USER=postgres
   PLANNER_DATABASE_PASSWORD=yourpassword
   PLANNER_DATABASE_PORT=5432
   ```

3. **Start the services:**
   ```bash
   docker-compose up -d
   ```
   This will start:
   - The Go application at `http://localhost:8080`
   - PostgreSQL database at `http://localhost:5432`
   - Mailpit (Web Interface) at `http://localhost:8025`
   - pgAdmin at `http://localhost:8081`

---

## 📬 Email Testing

All emails sent by the application are captured by **Mailpit** during development. You can view them by accessing the dashboard at [http://localhost:8025](http://localhost:8025).

## 📜 License

This project is licensed under the MIT License.

---
*Developed with 💜 during Rocketseat's NLW*
