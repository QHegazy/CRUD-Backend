# Task Management Application - Backend (`task-backend`)

A simple task management application backend written in Go, providing a RESTful API to manage tasks (create, read, update, delete) using in-memory storage.

---

## 📂 Project Structure

```
.
├── .github
│   └── workflows
│       ├── cd.yml               # Continuous deployment pipeline
│       └── ci.yml               # Continuous integration pipeline
├── cmd
│   └── main.go
├── internal
│   ├── dto
│   │   └── task_dto.go         # Data Transfer Objects for API requests/responses
│   ├── handlers
│   │   ├── task_handler.go     # HTTP handlers for task endpoints
│   │   └── task_handler_test.go# Unit tests for handlers
│   ├── middlewares
│   │   └── auth_middlewares.go # Authentication middleware (JWT-based)
│   ├── models
│   │   └── task_model.go       # Task domain model
│   ├── res
│   │   └── res.go              # Standard response formatting
│   ├── router
│   │   └── routes.go           # API route definitions
│   ├── server
│   │   └── server.go           # Server initialization and startup logic
│   ├── services
│   │   ├── task_service.go     # Business logic for task management
│   │   └── task_service_test.go# Unit tests for services
│   └── storage
│       ├── memory.go           # In-memory storage implementation
│       └── task_memory_test.go # Unit tests for storage
└── utils
|    ├── jwt.go                  # JWT utility functions
|    └── obfuscate.go            # Helpers for data obfuscation
├── docker-compose.yml
├── Dockerfile
├── Makefile                    # Build, run, test, and clean commands
├── README.md                   # This documentation file
```

---

## 🚀 Prerequisites

* Go >= 1.24
* Docker & Docker Compose (optional, for containerized setup)
* `make` (to use provided Makefile commands)

---

## 🔧 Setup & Run (Local)

1. **Clone the repository**:

   ```bash
   git clone https://github.com/QHegazy/CRUD-Backend.git
   cd task-backend
   ```
2. **Create a `.env` file**:

   Create a file named `.env` in the root directory with the following content:

   ```env
   PORT={port}
   APP_ENV=local
   CORS={url}
   SECRET_KEY={a_very_secret_key_that_no_one_can_guess}
   ```
3. **Build the application**:

   ```bash
   make build
   ```

4. **Start the server**:

   ```bash
   make run
   ```

5. **Access the API**: By default, the server listens on `:8080`.

   * List tasks:   `GET http://localhost:8080/tasks`
   * Get task by ID: `GET http://localhost:8080/tasks/{id}`
   * Create task:  `POST http://localhost:8080/tasks`
   * Update task:  `PUT http://localhost:8080/tasks/{id}`
   * Delete task:  `DELETE http://localhost:8080/tasks/{id}`

---

## 🐳 Docker & Docker Compose (Optional)

Build and run the service using Docker:

```bash
# Build the Docker image
docker build -t task-backend:latest .

# Run with Docker Compose
docker-compose up -d
```

API will be available at `http://localhost:8080`.

---

## 🧪 Testing

Run the full test suite (handlers, services, storage):

```bash
make test
```

Tests are written using Go's standard `testing` package and cover unit tests for each layer.

---

## 📝 Makefile Commands

| Command      | Description                             |
| ------------ | --------------------------------------- |
| `make all`   | Build binary and run tests              |
| `make build` | Compile the Go application              |
| `make run`   | Execute the compiled binary             |
| `make watch` | Live-reload on file changes (via `air`) |
| `make test`  | Run all unit tests                      |
| `make clean` | Remove build artifacts (`bin/`, `tmp/`) |

---

## 🏗️ Architecture & Design Decisions

1. **Layered Structure:** The codebase is organized in layers (`handlers`, `services`, `storage`) to enforce separation of concerns:

   * **Handlers**: Deal with HTTP specifics, input validation, and response formatting.
   * **Services**: Contain business logic and orchestrate calls between handlers and storage.
   * **Storage**: Abstract storage layer (in-memory implementation), easy to swap with a persistent store later.

2. **DTOs & Models**:

   * **Models** define the domain entities (`Task`).
   * **DTOs** map HTTP request/response payloads to internal models, ensuring loose coupling.

3. **Middleware**:

   * JWT-based authentication middleware to protect endpoints.
   * Easily extendable to include logging, rate-limiting, or CORS.

4. **Router & Server**:

   * Centralized route definitions in `router/routes.go` for clarity.
   * Server initialization in `server/server.go` to encapsulate configuration (ports, middlewares).

5. **Testing Strategy**:

   * Unit tests for each layer in parallel.
   * Use of in-memory store simplifies test setup and teardown.

6. **Dockerization**:

   * Multi-stage Dockerfile for small, secure images (build with CGO disabled, distroless runtime).
   * `docker-compose.yml` for multi-service orchestration 



