# qualifire-home-assignment

LLM Gateway service responsible for routing requests to LLM providers and performing extensible behaviors, such as:

- Virtual API-key management
- Quota and rate enforcement
- Metrics and logging
- Request validation and error normalization

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
    - [Top-level](#top-level)
    - [`cmd/`](#cmd)
    - [`internal/`](#internal)
    - [`tests/`](#tests)
    - [`tmp/`](#tmp)
- [Setup Instructions](#setup-instructions)
    - [Prerequisites](#prerequisites)
    - [Local Setup](#local-setup)
    - [Running with Docker](#running-with-docker)
- [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [`keys.json` Virtual Key Configuration](#keysjson-virtual-key-configuration)
- [Example Client Code](#example-client-code)
    - [Go HTTP Client](#go-http-client)
    - [Example Request / Response](#example-request--response)
- [Implementation Notes](#implementation-notes)
    - [Technologies and Design Choices](#technologies-and-design-choices)
    - [Concurrency Handling](#concurrency-handling)
    - [Testing Strategy](#testing-strategy)

---

## Overview

This service exposes a unified HTTP API for chat completions and related LLM capabilities. Clients call this gateway using **virtual API keys**. The gateway authenticates the virtual key, enforces quotas, collects metrics, and forwards the request to an underlying LLM provider.

The main goals are:

- Provide a stable, provider-agnostic interface.
- Manage cross-cutting concerns (auth, quota, logging, metrics) in one place.
- Make it easy to add or change LLM providers without affecting clients.

---

## Features

- HTTP API for chat completions.
- Integration with multiple LLM providers behind a single interface.
- Virtual API key management and quota enforcement.
- Metrics middleware for observability (latency, counts, etc.).
- Structured error types with consistent JSON responses.
- Comprehensive unit tests around core functionality (auth, quota, metrics, providers, etc.).

---

## Tech Stack

- **Language:** Go 1.24
- **Frameworks / Libraries:**
    - Standard `net/http` for HTTP server.
    - Gin for routing and HTTP middleware in the HTTP layer.
    - `testify` for assertions in tests.
- **Config / Environment:** `.env` files + Go config loader.
- **Containerization:** Docker + `docker-compose`.
- **Modules & Packaging:** Go modules (`go.mod`), layered internal packages.

---

## Project Structure

High-level structure:

### Top-level

- **`go.mod`** – Go module file with dependencies.
- **`.env` / `.env.example`** – Environment configuration for local development.
- **`Dockerfile`** – Docker image build instructions for the API service.
- **`docker-compose.yaml`** – Orchestrates the service and any supporting services.
- **`.air.toml`** – Hot-reload / live-reload configuration for local development.
- **`tmp/`** – Temporary or scratch files, not used in core logic.

### `cmd/`

- **`cmd/api/main.go`**  
  Application entrypoint. Typical responsibilities:
    - Load configuration.
    - Initialize logger, metrics, quota services, and providers.
    - Set up HTTP router, middleware, and routes.
    - Start and gracefully stop the HTTP server.

### `internal/`

Core application logic, organized by domain / responsibility.

#### `internal/configs/`

- **`config.go`** – Configuration loading and parsing (env variables, config paths).
- **`keys.json`** – Virtual API key definitions (used by quota/auth layers).

#### `internal/http/`

HTTP-facing layer, including controllers, errors, middleware, routes, and validation.

- **`controllers/`**
    - `chat_completions.go` – HTTP handlers for chat completion endpoints.
    - `controller.go` – Shared controller helpers (e.g., response helpers, base types).
    - `metrics.go` – Controllers or endpoints related to metrics exposure.

- **`errors/`**
    - `api_provider.go` – Error helpers for provider-level failures (LLM calls).
    - `error.go` – Core error type(s) and factory functions for consistent JSON errors.
    - `validation.go` – Validation-specific errors and helpers.

- **`middleware/`**
    - `metrics.go` – Metrics middleware, capturing per-request data (counts, latency, etc.), aggregating usage, and coordinating with services.
    - `quota.go` – Quota / virtual key middleware, enforcing usage limits before hitting providers.

- **`routes/`**
    - Route registration, wiring URL paths and HTTP methods to controllers and middleware.

- **`validators/`**
    - `chat_comlpletion.go` – Validation logic for chat completion requests (required fields, formats).
    - `validator.go` – Shared validation helpers / interfaces.

#### `internal/loggers/`

- **`logger.go`** – Logging abstraction and initialization (e.g., log levels, formats, request-scoped fields).

#### `internal/models/`

Domain and transport models shared across layers.

- **`logger.go`** – Types related to structured logging context.
- **`model.go`** – Core domain models for chat requests/responses and generic entities.
- **`proxy_request.go`** – Models for outgoing proxy requests to providers (provider-specific representations).

#### `internal/providers/`

LLM provider abstraction and implementations.

- **`provider.go`** – Provider interface(s) used by the rest of the system.
- **`opan_ai.go`** – Implementation for OpenAI-like provider (request/response mapping, error handling).
- **`anthropic.go`** – Implementation for Anthropic-like provider.

#### `internal/services/`

Domain services encapsulating logic not tied directly to HTTP.

- **`metrics.go`** – Metrics aggregation service (counters, timings, derived stats like averages).
- **`quota.go`** – Quota management service:
    - Tracks request counts and token usage per virtual key.
    - Handles time-window logic and limits.
    - Exposes a simple API to check and increment usage.

#### `internal/transports/`

- **`logging.go`** – Transport wrappers (e.g., logging around HTTP client calls to providers, or logging middleware for outbound traffic).

#### `internal/utils/`

Cross-cutting utilities.

- **`auth.go`** – Helpers for extracting and validating virtual API keys from headers.
- **`tokens.go`** – Token-related utilities (e.g., counting or tracking token usage).

### `tests/`

Dedicated tests for core behaviors.

- **`auth_test.go`** – Tests for auth utilities (virtual key extraction and formats).
- **`chat_comlpletion_test.go` / `chat_completions_test.go`** – Tests around chat completion controller/validator behavior.
- **`config_test.go`** – Tests for configuration loading and validation.
- **`error_test.go`** – Tests for error types and their mapping to API responses.
- **`logger_test.go`** – Tests for logger initialization and behavior.
- **`metrics_test.go`** – Tests around metrics middleware and service interactions, including duration and usage recording.
- **`provider_test.go`** – Tests for provider implementations and interface behavior.
- **`quota_test.go`** – Tests for quota limits, time windows, and usage tracking.
- **`routes_test.go`** – Tests to ensure routes are registered correctly and resolve as expected.
- **`tokens_test.go`** – Tests for token utilities.

---

## Setup Instructions

### Prerequisites

- Go **1.24+**
- Docker & Docker Compose (optional but recommended)
- A valid set of real provider API keys (for the backing LLM providers)
- A configured `.env` file and `keys.json` virtual key configuration

Clone the repository:

### Local Setup

1. **Create `.env` file**

   ```bash
   cp .env.example .env
   ```

   Adjust it according to your environment (see [Environment Variables](#environment-variables)).

2. **Configure virtual keys**

   Edit `internal/configs/keys.json` with definitions of virtual keys and their associated providers/limits.

3. **Install dependencies**

   ```bash
   go mod tidy
   ```

4. **Run the API**

   ```bash
   go run ./cmd/api
   ```

   The service will start on the configured port.

5. **Run tests**

   ```bash
   go test ./...
   ```

### Running with Docker

1. Build and start the service:

   ```bash
   docker-compose up --build
   ```

2. The API will be available on the host port defined in `docker-compose.yaml` (for example, `http://localhost:8080`).

3. To stop:

   ```bash
   docker-compose down
   ```

---

## Configuration

### Environment Variables

Typical environment variables (names are illustrative; align with your `config.go`):

- `APP_PORT` – Port on which the HTTP server listens (e.g., `8080`).
- `APP_ENV` – Environment (`local`, `dev`, `prod`, etc.).
- `OPENAI_API_KEY` – Real provider API key for OpenAI-backed routes.
- `ANTHROPIC_API_KEY` – Real provider API key for Anthropic-backed routes.
- `KEYS_CONFIG_PATH` – Path to the `keys.json` file with virtual key definitions.
- `LOG_LEVEL` – Log level (`debug`, `info`, `warn`, `error`).
- `METRICS_ENABLED` – Enable or disable metrics tracking.

Example `.env`:

bash APP_PORT=8080 APP_ENV=local
OPENAI_API_KEY=sk-openai-... ANTHROPIC_API_KEY=sk-anthropic-...
KEYS_CONFIG_PATH=internal/configs/keys.json
LOG_LEVEL=debug METRICS_ENABLED=true

### `keys.json` Virtual Key Configuration

`keys.json` describes which virtual API keys your clients can use and how they map to providers/quotas.

A typical entry includes:

- A virtual key ID (the value sent by the client).
- A provider identifier (e.g., `openai`, `anthropic`).
- Optional quotas (max requests, max tokens, etc.).
- Optional metadata for auditing or routing.

---

## Implementation Notes

### Technologies and Design Choices

- **Go** is used for its strong concurrency model, efficient HTTP server, and rich standard library.
- **Gin** provides ergonomic routing and middleware composition for the HTTP layer.
- **Modular internal packages** keep concerns separated:
    - HTTP layer vs. domain services vs. provider integrations.
    - Config, logging, and utilities are reusable and testable.

The overall design aims to keep the provider-specific logic isolated, while cross-cutting concerns (quota, metrics, auth) are implemented once and reused across providers.

### Concurrency Handling

Go’s concurrency features are used in several places:

- **HTTP Server**  
  Each incoming HTTP request is handled in its own goroutine, as provided by the Go standard library and Gin. This allows the service to process many concurrent requests efficiently.

- **Quota Service**  
  Quota tracking uses concurrency-safe structures (e.g., `sync.Map` for per-key entries and `sync.Mutex` within entries) to handle updates from multiple concurrent requests. This ensures consistent request and token counts per virtual key, even under high load.

- **Metrics Service**  
  Metrics aggregation is designed for concurrent use, so recording stats from different handlers and goroutines remains safe.

- **Context-based Cancellation**  
  Provider calls and long-running operations can leverage `context.Context` for:
    - Request-level deadlines.
    - Canceling downstream operations if the client disconnects or times out.

- **Graceful Shutdown (entrypoint level)**  
  The main process can use signals and context timeouts to:
    - Stop accepting new connections.
    - Allow in-flight requests to finish.
    - Cleanly release resources.

### Testing Strategy

The project has a strong focus on testability:

- **Unit Tests**
    - Utilities such as auth header parsing are covered by isolated tests.
    - Quota and metrics services are tested for correctness, edge cases, and concurrency-related behavior (e.g., resetting windows, updating usage).
    - Error helpers are validated to ensure they produce consistent and predictable API responses.

- **HTTP / Middleware Tests**
    - Middleware such as metrics and quota are tested using in-memory Gin routers and `httptest` recorders.
    - Tests assert that metrics are tracked or skipped in the right scenarios (e.g., specific endpoints, missing keys).
    - Status codes and response formats are validated end-to-end through the HTTP stack.

- **Provider and Routing Tests**
    - Provider abstractions are tested to ensure uniform behavior and proper mapping between internal models and external API calls.
    - Routes tests confirm that endpoints are wired correctly.

This combination helps ensure that core behaviors—auth, quota, metrics, error handling, and routing—remain stable and regressions are caught early.

---