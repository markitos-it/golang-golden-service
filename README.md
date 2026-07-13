# Golden Service

This project is a "Golden Path" service written in Go. It provides a template for creating new microservices with a focus on best practices, including a clean architecture, gRPC API, and a complete development environment.

## Features

*   **gRPC API:** Exposes a gRPC API with server reflection enabled for easy discovery and testing.
*   **PostgreSQL Integration:** Uses GORM for database integration with a PostgreSQL database.
*   **Configuration Management:** Manages configuration using Viper, loading from environment variables or a local `app.env` file.
*   **Project Scaffolding:** Includes a `clonator` CLI tool to generate new services from this template.
*   **Makefile-driven Workflow:** Simplifies common development tasks like running, testing, and database management.
*   **Dockerized Environment:** Provides a Docker Compose setup for running a local PostgreSQL database.

## Prerequisites

Before you begin, ensure you have the following installed:

*   [Go](https://golang.org/doc/install) (version 1.22 or newer)
*   [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
*   [Make](https://www.gnu.org/software/make/)
*   [`grpcurl`](https://github.com/fullstorydev/grpcurl) for interacting with the gRPC API.
*   Git Hooks Tools: For automated code quality and security checks.

## Getting Started

Follow these steps to get the service up and running.

### 0. Setup Git Hooks (Highly Recommended)

This project uses native Git hooks to automate code quality checks, formatting, and security scanning before you commit and push your code. This ensures consistency and prevents common errors from reaching the repository.

To set them up, please follow the detailed instructions in the guide:

**➡️ [Git Hooks Installation Guide](etc/git-hooks-appsec.md)**


### 1. Clone the Repository

```sh
git clone https://github.com/your-username/markitos-it-svc-golden.git
cd markitos-it-svc-golden
```

### 2. Start the Database

This command uses Docker Compose to start a PostgreSQL container.

```sh
make db-start
make db-create
```

### 3. Run the Application

This command will start the gRPC server.

```sh
make start
```

The server will be listening on the address specified in your configuration (default is `localhost:50051`).

## Usage

### Running the Service

To start the gRPC server:

```sh
make start
```

### Running Tests

Run all unit tests:

```sh
make test
```

Run end-to-end tests using `grpcurl`:

```sh
make test-e2e
```

### Database Management

The `Makefile` provides commands to manage the PostgreSQL database container:

*   **Start the database:** `make db-start`
*   **Stop the database:** `make db-stop`
*   **Create the database:** `make db-create` (Note: The database is created automatically on the first run)
*   **Drop the database:** `make db-drop`

### Protocol Buffers

To regenerate the gRPC code from the `.proto` file:

```sh
make proto
```

## API

The service exposes a gRPC API defined in the Protocol Buffers file located at `internal/infrastructure/proto/app.proto`. With the server running and reflection enabled, you can explore the available services and methods using `grpcurl`.

**List Services:**

```sh
grpcurl -plaintext localhost:50051 list
```

**Describe a Service:**

```sh
grpcurl -plaintext localhost:50051 describe .goldenservice.Goldenservice
```

## Project Structure

The project follows a standard Go project layout:

*   `bin/`: Shell scripts for common tasks (starting, testing, etc.).
*   `cmd/`: Main application entry points.
    *   `app/`: The main gRPC service application.
    *   `clonator/`: A CLI tool for scaffolding new projects from this template.
*   `etc/`: Configuration files, including `docker-compose.yaml`.
*   `internal/`: Private application and library code.
    *   `domain/`: Core business logic, models, and services.
    *   `infrastructure/`: Components that interact with external systems (database, gRPC server).
*   `Makefile`: Defines and automates common development workflows.

## Configuration

The application is configured using a flexible system that loads settings from multiple sources. The order of priority is as follows:

1.  **Environment Variables:** The most prioritized method. Any variable set in the environment (e.g., `export DATABASE_DSN=...`) will override all other settings.
2.  **`.env` File:** You can create a `.env` file in the root directory. Its values will override those from `config.yaml`. This is useful for local development.
3.  **`config.yaml` File:** A `config.yaml` file can be created in the root directory to define base configuration values.

The `main.go` file loads the configuration at startup. Refer to `internal/infrastructure/configuration/config.go` for all available configuration options.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
