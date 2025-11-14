
# qualifire-home-assignment

LLM Gateway service responsible for routing requests to multiple LLM providers (e.g., OpenAI, Anthropic) and performing extensible behaviors such as:

- Virtual API-key management
- Quota and rate enforcement
- Metrics and logging
- Request validation and error normalization

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Setup Instructions](#setup-instructions)
  - [Prerequisites](#prerequisites)
  - [Local setup (without Docker)](#local-setup-without-docker)
  - [Running with Docker](#running-with-docker)
- [Configuration](#configuration)
  - [Environment variables](#environment-variables)
  - [`keys.json` virtual key configuration](#keysjson-virtual-key-configuration)
  - [Sample configuration flow](#sample-configuration-flow)
- [Example Client Code](#example-client-code)
  - [Go HTTP client](#go-http-client)
  - [Example request/response formats](#example-requestresponse-formats)
- [Project Structure](#project-structure)
  - [Top-level](#top-level)
  - [`cmd/`](#cmd)
  - [`internal/`](#internal)
  - [`tests/`](#tests)
  - [`tmp/`](#tmp)
- [Implementation Notes](#implementation-notes)
  - [Concurrency handling](#concurrency-handling)
  - [Testing strategy](#testing-strategy)

---

## Overview

This service exposes a unified HTTP API for chat completions and related LLM capabilities. Incoming requests are:

1. Authenticated and quota-checked using *virtual API keys*.
2. Validated and normalized into a common internal model.
3. Routed to the appropriate provider implementation (OpenAI, Anthropic, etc.).
4. Logged and instrumented with metrics for observability.

The design favours:

- Clear separation between HTTP layer, business logic, and provider integrations.
- Extensibility (adding new providers and behaviors with minimal changes).
- Robust error handling and testability.

---

## Features

- HTTP API for chat completions (and related operations).
- Multiple provider integrations via a common provider interface.
- Virtual API key management with configurable quotas.
- Metrics and logging middleware.
- Centralized error handling with normalized responses.
- Unit and integration tests for core components.

---

## Tech Stack

- **Language:** Go 1.24
- **Dependency management:** Go modules (`go.mod`)
- **HTTP server & routing:** Go standard library (`net/http`), with internal routing abstraction
- **Configuration:** Environment variables + JSON key file (`keys.json`)
- **Logging:** Custom logger abstraction
- **Metrics:** Internal metrics services & middleware (exporter-agnostic)
- **Containerization:** Docker + `docker-compose`

---

## Setup Instructions

### Prerequisites

- Go **1.24+**
- Docker & Docker Compose (optional but recommended)
- Make (optional, if you want to add Makefile targets)

Clone the repository:

- **`go.mod`** – Go module definition and dependencies.
- **`.env` / `.env.example`** – Environment configuration for local development.
- **`Dockerfile`** – Container image build definition.
- **`docker-compose.yaml`** – Compose configuration for running the gateway and related services.
- **`README.md`** – This documentation.
- **`tmp/`** – Scratch or temporary files (ignored at runtime).

### `cmd/`

- **`cmd/api/main.go`**  
  Application entrypoint. Typical responsibilities:
    - Parse configuration.
    - Initialize logger, metrics, quota services.
    - Set up HTTP router and middleware.
    - Start the HTTP server (with graceful shutdown handling).

### `internal/`

This directory contains the main application code, grouped by responsibility.

#### `internal/configs/`

- **`config.go`**  
  Central configuration loading and parsing logic:
    - Reads env variables and config files.
    - Builds a strongly typed configuration struct used across the app.
- **`keys.json`**  
  Default/example virtual key configuration, mapping client keys to providers, quotas, and metadata.

#### `internal/http/`

High-level HTTP layer: routing, controllers, validation, and errors.

- **`controllers/`**
    - `chat_completions.go` – Handler functions for chat completion endpoints:
        - Parse and validate incoming payloads.
        - Call underlying services/providers.
        - Shape responses.
    - `controller.go` – Base helpers or shared controller utilities (e.g., common response helpers).
    - `metrics.go` – HTTP handlers that expose metrics endpoints or controller-side metric helpers.

- **`errors/`**
    - `api_provider.go` – Error mapping specific to provider-related failures (e.g., OpenAI/Anthropic).
    - `error.go` – Central error types and helpers (standardized error structure for responses).
    - `validation.go` – Validation error helpers and conversion to HTTP responses.

- **`middleware/`**
    - `metrics.go` – HTTP middleware that records metrics (latency, status codes, etc.).
    - `quota.go` – HTTP middleware enforcing quotas and checking virtual API keys.

- **`routes/`**
    - Router registration: mapping paths (`/v1/chat/completions`, `/metrics`, etc.) to handlers and middleware chains.

- **`validators/`**
    - `chat_comlpletion.go` – Validation logic for chat completion requests (e.g., required fields, length limits).
    - `validator.go` – Shared validation helpers or interfaces.

#### `internal/loggers/`

- **`logger.go`**  
  Logging abstraction and initialization:
    - Wraps the underlying logger implementation.
    - Exposes helper functions for structured logging (e.g., fields for request ID, key ID, provider).

#### `internal/models/`

- **`logger.go`**  
  Model(s) related to logging context or structured log payloads.
- **`model.go`**  
  Core domain models used across layers:
    - Internal representation of chat requests/responses.
    - Provider-agnostic types (roles, messages, usage).
- **`proxy_request.go`**  
  Models for proxying/forwarding requests:
    - Structures that represent how the request is sent to providers.

#### `internal/providers/`

Abstractions and concrete implementations for LLM providers.

- **`provider.go`**  
  Provider interface definition:
    - Methods like `ChatCompletion(ctx, request)` used by the rest of the system.
- **`open_ai.go`**  
  Implementation for OpenAI:
    - Translates internal requests to OpenAI API.
    - Handles OpenAI-specific errors and responses.
- **`anthropic.go`**  
  Implementation for Anthropic:
    - Translates internal requests to Anthropic API.
    - Handles provider-specific response mapping.

#### `internal/services/`

Domain services that encapsulate business logic decoupled from HTTP.

- **`metrics.go`**  
  Metrics service responsible for:
    - Recording counters, histograms, and gauges.
    - Providing helpers used by middleware and controllers.
- **`quota.go`**  
  Quota management and virtual key service:
    - Loads key definitions from `keys.json`.
    - Tracks usage and enforces limits.
    - Provides lookup and validation for virtual keys.

#### `internal/transports/`

- **`logging.go`**  
  Transport-level logging wrappers:
    - Possibly wraps HTTP client calls to providers.
    - May log request/response metadata, latency, and error details.

#### `internal/utils/`

Utility helpers and cross-cutting concerns.

- **`auth.go`**  
  Helpers for extracting and validating API keys from HTTP requests (e.g., `Authorization` header).
- **`tokens.go`**  
  Token-related utilities:
    - Parsing tokens.
    - Possibly token accounting or tokenization helpers shared across services.

### `tests/`

Dedicated test package(s) for focused unit and behavior tests.

- **`auth_test.go`** – Tests for authentication helpers and middleware behavior.
- **`chat_comlpletion_test.go`, `chat_completions_test.go`** – Tests for chat completion controllers, validators, and services.
- **`config_test.go`** – Tests for configuration loading and validation.
- **`error_test.go`** – Tests for error types and error-to-HTTP mapping.
- **`logger_test.go`** – Tests for logger initialization and structured logging behavior.
- **`metrics_test.go`** – Tests for metrics service and middleware.
- **`provider_test.go`** – Tests for provider interface and provider implementations (OpenAI, Anthropic).
- **`quota_test.go`** – Tests for quota logic and virtual key enforcement.
- **`routes_test.go`** – Tests for routing setup (ensuring endpoints are correctly wired).
- **`tokens_test.go`** – Tests for token utilities.

---

## Implementation Notes

### Concurrency handling

The gateway relies heavily on Go’s concurrency primitives while keeping them encapsulated:

- **HTTP handlers**: Each incoming request is handled in its own goroutine by the Go HTTP server.
- **Provider calls**: Provider integrations use `context.Context` for cancellation and timeouts, allowing you to:
    - Enforce per-request deadlines.
    - Cancel in-flight provider calls if the client disconnects.
- **Quota and metrics**:
    - Designed to be concurrency-safe (e.g., using atomic counters or mutexes internally where needed).
    - Where stateful components exist (e.g. in-memory usage counters), they are protected via proper synchronization primitives.
- **Graceful shutdown**:
    - The main process typically sets up signal handling (`os.Signal`), then uses context-aware shutdown to stop accepting new requests while letting in-flight requests complete.

Overall, concurrency is primarily handled via:

- Goroutines (per request and per provider call where needed).
- Channels or context for cancellation and signalling.
- Mutexes or atomic operations for shared mutable state (quota, metrics, caches).

### Testing strategy

The project emphasizes testability across layers:

- **Unit tests**:
    - Target individual services (quota, metrics, providers) in isolation.
    - Use interfaces and small abstractions to inject test doubles/mocks.
- **HTTP-level tests**:
    - Use `net/http/httptest` to test controllers, routes, and middleware end-to-end in-memory.
    - Verify status codes, response payload shapes, and error handling.
- **Configuration tests**:
    - Ensure that configuration loading fails fast on invalid input and correctly processes valid combinations of environment variables and files.
- **Provider tests**:
    - Use mock HTTP servers or stub transport layers to verify provider request/response handling without hitting real external APIs.

This combination provides good coverage for both correctness and regression safety, while keeping tests fast and deterministic.

---

If you’d like, I can also add a short “Quick Start” section with ready-to-run `curl` examples for the main endpoints.