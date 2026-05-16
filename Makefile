APP_NAME := movies-api
IMAGE_NAME := movies-api
IMAGE_TAG := latest
PORT ?= 8080
BASE_URL ?= http://localhost:8080
RENDER_BASE_URL ?= https://movies-api-latest.onrender.com
GOCACHE := $(CURDIR)/.gocache

.PHONY: help test run build docker-build docker-run compose-up compose-down traffic traffic-render clean

help:
	@echo "Available targets:"
	@echo "  make test         Run unit tests"
	@echo "  make run          Run the API locally with go run"
	@echo "  make build        Build the local binary"
	@echo "  make docker-build Build the Docker image"
	@echo "  make docker-run   Run the Docker image using .env"
	@echo "  make compose-up   Start the app with docker compose"
	@echo "  make compose-down Stop the docker compose service"
	@echo "  make traffic      Generate sample traffic for monitoring"
	@echo "  make traffic-render Generate sample traffic against the Render deployment"
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

compose-up:
	docker compose up --build

compose-down:
	docker compose down

traffic:
	BASE_URL=$(BASE_URL) sh ./scripts/generate_traffic.sh

traffic-render:
	BASE_URL=$(RENDER_BASE_URL) sh ./scripts/generate_traffic.sh

clean:
	rm -f $(APP_NAME)
