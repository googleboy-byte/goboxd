# goboxd Makefile

BINARY=goboxd
IMAGE=goboxd:latest
SERVER_URL=http://localhost:8080
VERSION=0.1.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT)"

.PHONY: build run test integration load lint clean

build:
	@echo "Building $(BINARY) $(VERSION) ($(COMMIT))..."
	go build $(LDFLAGS) -o $(BINARY) ./cmd/goboxd/main.go

run:
	@echo "Bringing up Docker container $(IMAGE)..."
	-docker kill $(BINARY) 2>/dev/null || true
	-docker rm $(BINARY) 2>/dev/null || true
	docker build -t $(IMAGE) .
	docker run -d --privileged --name $(BINARY) -p 8080:8080 $(IMAGE)
	@echo "Server is starting at $(SERVER_URL)"

test:
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

integration:
	@echo "Running integration tests..."
	@curl -s -o /dev/null --connect-timeout 2 $(SERVER_URL)/healthz || (echo "Error: Server is not running at $(SERVER_URL). Run 'make run' first." && exit 1)
	bash tests/integration/run_all.sh $(SERVER_URL)

load:
	@curl -s -o /dev/null --connect-timeout 2 $(SERVER_URL)/healthz || (echo "Error: Server is not running. Run 'make run' first." && exit 1)
	@if ! command -v hey >/dev/null 2>&1; then \
		echo "hey not found. Install with: go install github.com/rakyll/hey@latest"; \
		exit 1; \
	fi
	bash tests/load/load.sh $(SERVER_URL)

lint:
	@echo "Linting code..."
	go vet ./...
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not found, skipping (go vet passed)"; \
	fi

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY)
	-docker kill $(BINARY) 2>/dev/null || true
	-docker rm $(BINARY) 2>/dev/null || true
