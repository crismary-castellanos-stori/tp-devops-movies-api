#first stage: build the application
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o /movies-api ./api/main.go

#second stage: create the final image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /movies-api /usr/local/bin/movies-api

EXPOSE 8080

ENTRYPOINT ["movies-api"]
