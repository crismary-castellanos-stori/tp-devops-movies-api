APP_NAME := movies-api
IMAGE_NAME := movies-api
IMAGE_TAG := latest
PORT ?= 8080
GOCACHE := $(CURDIR)/.gocache

.PHONY: help test run build docker-build docker-run clean

help:
	@echo "Available targets:"
	@echo "  make test         Run unit tests"
	@echo "  make run          Run the API locally with go run"
	@echo "  make build        Build the local binary"
	@echo "  make docker-build Build the Docker image"
	@echo "  make docker-run   Run the Docker image using .env"
	@echo "  make clean        Remove generated artifacts"

test:
	GOCACHE=$(GOCACHE) go test ./...

run:
	GOCACHE=$(GOCACHE) go run ./api/main.go

build:
	GOCACHE=$(GOCACHE) go build -o $(APP_NAME) ./api/main.go

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run:
	docker run --rm --name movies-api -p $(PORT):8080 --env-file .env $(IMAGE_NAME):$(IMAGE_TAG)

clean:
	rm -f $(APP_NAME)
