#builder: compile stage
FROM golang:latest as builder
WORKDIR /qualifire/app

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy sources
COPY . ./
RUN mv .env.example .env

# Set GOBIN to a safe place for the binary
ENV GOBIN=/qualifire/app/bin
RUN mkdir -p $GOBIN

# Install the binary
RUN go install ./cmd/api

# ---------------------------------------
# Development stage (with air)
FROM golang:latest as dev-builder
WORKDIR /qualifire/app
COPY --from=builder /qualifire/app/.env .env
COPY . ./
RUN go install github.com/air-verse/air@latest
ENV AIR_CONFIG=/qualifire/app/.air.toml
ENTRYPOINT ["air"]


# ---------------------------------------
# Production stage
FROM debian:13-slim as prod-builder
WORKDIR /qualifire/app
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
# Copy installed binary from $GOBIN
COPY --from=builder /qualifire/app/bin/api ./api

# Copy runtime files
COPY --from=builder /qualifire/app/.env .env
RUN mkdir -p /internal/configs
COPY --from=builder /qualifire/app/internal/configs/keys.json ./internal/configs/keys.json

CMD ["./api"]
