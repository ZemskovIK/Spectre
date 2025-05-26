FROM golang:1.24-bullseye AS builder

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN apt update && apt install -y gcc libc6-dev

WORKDIR /spectre

COPY ./src/ /spectre

RUN go mod download

RUN go build -o migrator ./cmd/migrator/main.go

RUN go build -o spectre ./cmd/spectre/main.go

FROM debian:bullseye-slim

WORKDIR /spectre

RUN mkdir -p /spectre/db/migrations && mkdir -p /spectre/bin

COPY ./db/migrations/ /spectre/db/migrations

COPY --from=builder /spectre/migrator /spectre/bin

COPY --from=builder /spectre/spectre /spectre/bin