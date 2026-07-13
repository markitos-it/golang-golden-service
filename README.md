# Golden Service

Base (golden path) service in Go with clean architecture, a gRPC API, and PostgreSQL persistence using GORM.

## Features

- gRPC API with reflection enabled for discovery and testing with grpcurl.
- PostgreSQL persistence via GORM.
- Configuration with Viper from config.yaml, .env, and environment variables.
- Centralized gRPC request logging through a unary interceptor.
- Database logging at warn level, ignoring noisy record not found entries.
- Development workflow with Makefile and scripts in bin/.
- Clonator CLI to generate new services from this template.

## Requirements

- Go (use the version defined in go.mod).
- Docker and Docker Compose.
- Make.
- grpcurl (recommended for manual gRPC testing).

## Quick Start

1. Clone the repository:

```sh
git clone https://github.com/your-username/markitos-it-svc-golden.git
cd markitos-it-svc-golden
```

2. Start the database and create the schema:

```sh
make db-start
make db-create
```

3. Start the service:

```sh
make start
```

## Main Commands

- Start app: `make start`
- Run unit/integration tests: `make test`
- Run gRPC e2e tests with grpcurl: `make test-e2e`
- Generate protobuf/gRPC code: `make proto`
- Start DB: `make db-start`
- Stop DB: `make db-stop`
- Create DB: `make db-create`
- Drop DB: `make db-drop`
- Refresh module dependencies: `make tidy`

You can also see the full command list with:

```sh
make help
```

## Configuration

The service loads configuration with this priority order:

1. Environment variables.
2. .env file (if present).
3. config.yaml file (if present).

Relevant variables:

- `DATABASE_DSN`: PostgreSQL connection string.
- `GRPC_SERVER_ADDRESS`: gRPC server address (example: `:50051`).
- `GOLDEN_UPLOADS_BASEDIR`: base path for poster files.

If `GOLDEN_UPLOADS_BASEDIR` is not defined, the fallback is `/tmp/goldens`.

## Logging

The project has two main logging layers:

1. gRPC logging (unary interceptor)
   - Logs every unary call: method, gRPC code, duration, and error.
   - Example format:
   ```text
   [grpc] method=/golden.Goldenservice/GetGolden code=NotFound duration=1.2ms error=rpc error: code = NotFound desc = Resource not found
   ```

2. Logging GORM
   - Configured at warn level.
   - `record not found` is ignored to reduce expected SQL noise on single-record lookups.

Additionally, when a record is missing in specific operations (for example `One` and `Delete`), a custom warning with the id is emitted for better functional traceability.

## API gRPC

Protobuf contract: `internal/infrastructure/proto/app.proto`

With the server running, you can explore services and methods:

```sh
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 describe golden.Goldenservice
```

Example lookup by id:

```sh
grpcurl -plaintext -d '{"id":"fdde6b8c-a3fa-47e9-a25e-7805759bd30f"}' localhost:50051 golden.Goldenservice/GetGolden
```

## Hooks and AppSec

To install and use security/quality hooks, check:

- `etc/hooks/git-hooks-appsec.md`

## Project Structure

- `bin/`: operational scripts (start, test, db, appsec, support).
- `cmd/app/`: gRPC service entrypoint.
- `cmd/clonator/`: scaffolding CLI.
- `etc/`: docker-compose, hooks, and deployment resources.
- `internal/domain/`: business model, services, and types.
- `internal/infrastructure/`: gRPC, configuration, and database.
- `testsuite/`: domain, infrastructure, and e2e tests.

## License

MIT. See LICENSE.
