# Stage 1: compile the Go application in an isolated builder image.
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Buildx injects these args so the binary can be compiled for the target platform.
ARG TARGETOS
ARG TARGETARCH

# Copy dependency files first to take advantage of Docker layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code once dependencies are cached.
COPY . .

# Build a static binary for the requested OS/architecture.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o /movies-api ./api/main.go

# Stage 2: create a minimal runtime image containing only the compiled binary.
FROM alpine:3.20

WORKDIR /app

# Copy the binary from the builder stage to keep the final image small.
COPY --from=builder /movies-api /usr/local/bin/movies-api

EXPOSE 8080

ENTRYPOINT ["movies-api"]
