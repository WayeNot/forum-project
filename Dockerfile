# syntax=docker/dockerfile:1

FROM golang:1.25.1-bookworm AS builder

WORKDIR /src

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/forum ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /app/data

COPY --from=builder /out/forum ./forum
COPY --from=builder /src/internal/db/migrations ./internal/db/migrations
COPY --from=builder /src/web ./web

ENV DB_PATH=/app/data/forum.db

EXPOSE 5500

CMD ["./forum"]
