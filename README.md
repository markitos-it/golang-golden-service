<h1 align="center">
  <br>
  <img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" alt="Go Logo" width="150">
  <br>
  🟡 markitos-it-service-golden
  <br>
</h1>

<h4 align="center">A production-ready, highly opinionated Go 1.26 gRPC microservice template for Markitos IT.</h4>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/gRPC-v1.60+-244c5a?style=flat-square&logo=google" alt="gRPC Version">
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker" alt="Docker Ready">
</p>

<p align="center">
  <a href="#-about-the-project">About</a> •
  <a href="#-key-features">Features</a> •
  <a href="#-getting-started">Getting Started</a> •
  <a href="#-environment-variables">Environment Variables</a> •
  <a href="#-docker">Docker</a>
</p>

---

## 📖 About The Project

**markitos-it-service-golden** is the foundational template ("Golden Service") for all new gRPC microservices. It allows developers to bypass the boilerplate setup and immediately start writing business logic. 

Built on **Go 1.26**, this service implements industry best practices, including a highly optimized multi-stage Docker build utilizing Alpine and Distroless for a minimal footprint and maximum security.

## ✨ Key Features

*   **Go 1.26**: Leveraging the latest standard library features.
*   **gRPC & Protobuf**: High-performance, strongly-typed remote procedure calls.
*   **Self-Sufficient Binary**: Compiled with `CGO_ENABLED=0` for maximum portability without external dependencies.
*   **Distroless Docker Image**: Uses `gcr.io/distroless/static-debian12` for an ultra-secure, minimal production environment.

## 🚀 Create a New Service (Clone)

You can quickly generate a new microservice based on this template (*Golden Service*) using our interactive Go CLI tool, **Clonator**.

```bash
make clonator
```

This will start a local server. Open your browser and navigate to http://localhost:8080.
There you will find an intuitive form where you can fill in:
- **Entity Singular:** (e.g., `user`) - Lowercase letters only.
- **Entity Plural:** (e.g., `users`) - Lowercase letters only.
- **Service Name:** (e.g., `markitos-it-service-users`) - Kebab-case format.

After filling in the details, you will see a confirmation screen. Upon acceptance, the tool will automatically clone, clean, and configure the new project for you. 😎

### Option 2: Interactive CLI 💻

If you prefer to perform the process directly from your terminal, you can run the interactive wizard:

```bash
make clone
```

The Bash script will guide you by asking for the required values step by step. It also supports parameter injection via environment variables (`MDK_CLONE_ENTITY_SINGULAR`, `MDK_CLONE_ENTITY_PLURAL`, `MDK_CLONE_TARGET_SERVICE_NAME`) for use in automated CI/CD processes.

## 🚀 Getting Started

### Prerequisites

*   Go >= 1.26
*   Docker

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/markitos-it/markitos-it-service-golden.git
   cd markitos-it-service-golden
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Run the service:
   ```bash
   go run cmd/app/main.go
   ```

## ⚙️ Environment Variables

You can configure the application using the following environment variables. The variables below are the default values, which can be overridden at runtime:

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `DATABASE_DSN` | `host=localhost user=admin password=admin dbname=markitos-it-svc-golden sslmode=disable` | Database connection string |
| `GRPC_SERVER_ADDRESS` | `:30000` | Address and port where the gRPC server will listen |

## 🐳 Docker

The project includes a highly optimized, multi-stage Dockerfile ready for production.

### Build the image
```bash
docker build -t markitos-it-svc-golden:1.0.0 .
```

### Run the container
```bash
docker run -p 3000:3000 -p 30000:30000 \
  -e DATABASE_DSN="host=db user=admin password=admin dbname=markitos-it-svc-golden sslmode=disable" \
  -e GRPC_SERVER_ADDRESS=":30000" \
  markitos-it-svc-golden:1.0.0
```

*gRPC API will be available at `localhost:30000`*