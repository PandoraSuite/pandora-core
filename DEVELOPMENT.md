# Pandora Core — Development Guide

Welcome to the developer guide for **Pandora Core**, the backend service of the Pandora open source project. This document is intended for contributors and maintainers who want to understand, run, and contribute to the project locally.

## :wrench: Requirements

Before getting started, make sure you have the following tools installed:

* [Go 1.23.4](https://go.dev/dl/) — ⚠️ required version
* [Air](https://github.com/air-verse/air) — for live-reloading
* [Delve](https://github.com/go-delve/delve) — for debugging
* [PostgreSQL](https://www.postgresql.org/) — as the default database

You can install Air and Delve via:

```bash
go install github.com/air-verse/air@v1.61.7
go install github.com/go-delve/delve/cmd/dlv@v1.24.1
```

## :rocket: Running the Project

### 1. Run modes

Run only the HTTP server:

```bash
go run ./cmd/http/main.go
```

Run only the gRPC server:

```bash
go run ./cmd/grpc/main.go
```

Run both servers:

```bash
go run ./cmd/main.go
```

### 2. Using Air for auto-reloading

**Pandora Core** includes an Air template configuration. To use it:

```bash
cp .air.template.toml .air.toml
```

Then edit the `.air.toml` file. Update the `cmd` value to match the desired target:

```toml
cmd = "go build -gcflags=\"all=-N -l\" -o ./tmp/main ./cmd/main.go"
```

You may change the final path depending on whether you want to run the HTTP server, gRPC server, or both (using `./cmd/main.go`).

To start the app:

```bash
air
```

## :bug: Debugging

**Pandora Core** supports Delve debugging. With Air and Delve installed:

1. First read [Using Air for auto-reloading](#2-using-air-for-auto-reloading)
1. Run the app with:

   ```bash
   air
   ```
2. Attach a debugger of your choice (e.g., VSCode).

### :mag: VSCode Example

If using Visual Studio Code, you can create or edit `.vscode/launch.json`:

**Local Debugging Configuration**
  
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to server [Local]",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "${workspaceFolder}",
      "port": 2345,
      "host": "127.0.0.1",
      "apiVersion": 2,
      "showLog": true,
      "trace": "verbose"
    }
  ]
}
```

**Docker Debugging Configuration**
  
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to server [Docker]",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "/app",
      "port": 2345,
      "host": "127.0.0.1",
      "substitutePath": [
        {
          "from": "${workspaceFolder}",
          "to": "/app"
        }
      ],
      "apiVersion": 2,
      "showLog": true,
      "trace": "verbose"
    }
  ]
}
```

Press `F5` to attach the debugger.

> :bulb: If you're using another IDE like Goland, Vim, or Sublime Text, feel free to document your setup here to help others.

## :gear: Required Environment Variables

**Pandora Core** requires the following environment variables to run properly:

* `PANDORA_DIR` — (optional) (default: `/etc/pandora`)
* `PANDORA_DB_DNS` — (optional) PostgreSQL connection string (default: `host=localhost port=5432 user=postgres password= dbname=pandora sslmode=disable timezone=UTC`)
* `PANDORA_JWT_SECRET` — (optional) Secret key used for signing authentication tokens (default: Randomly generated on startup. Consider setting a fixed value for consistent local development)
* `PANDORA_HTTP_PORT` — (optional) HTTP server port (default: `80`)
* `PANDORA_GRPC_PORT` — (optional) gRPC server port (default: `50051`)
* `PANDORA_EXPOSE_VERSION` — (optional) (default: `true`)

You can export them manually in your shell before starting the application

## :card_file_box: Generated Files and Folders

While running locally, **Pandora Core** may generate:

* `./tmp/` — temporary directory for compiled binaries when using Air
* `./{$PANDORA_DIR}/adminPanel/credentials.json` — root admin credentials file (created on first run)


## :whale: Running with Docker Compose

**Pandora Core** provides a `docker-compose.yml` setup for local development. To run the service in a containerized environment, you'll need **Docker** and **Docker Compose** installed on your system.

```bash
docker compose -f docker/docker-compose.yml up --build
```

**To attach a debugger to the container**: You can see the section [VSCode Example](#-vscode-example).

## :test_tube: Running Tests

You can run all unit tests with:

```bash
go test ./...
```

We encourage writing tests for new features and keeping existing tests passing.

## :file_folder: Project Structure

```bash
pandora-core/
├── cmd/                  # Entrypoints: HTTP, gRPC, or combined
├── db/                   # Init scripts and database Dockerfile
├── docker/               # Docker setup
├── internal/
│   ├── adapters/         # HTTP, gRPC, persistence, security
│   ├── app/              # Use case orchestration per domain
│   ├── config/           # Configuration loading and env support
│   ├── domain/           # DTOs, entities, enums, and errors
│   ├── ports/            # Interfaces
│   ├── utils/            # Helper functions (e.g., time)
│   ├── validator/        # Custom validators
│   └── version/          # Build/versioning info
├── proto/                # Protocol Buffers definitions
```

## :raising_hand: Contributing

We welcome contributions! Please see [`CONTRIBUTING.md`](./CONTRIBUTING.md) for guidelines.

If you encounter bugs, missing features, or have suggestions, open an issue or pull request.

## :lock: Security

If you discover a security vulnerability, **do not open an issue or PR**. Instead, please follow the secure disclosure process described in [`SECURITY.md`](./SECURITY.md).

## :handshake: Thanks

Thank you for helping improve **Pandora Core**. Your feedback and contributions help make this project better for everyone.
