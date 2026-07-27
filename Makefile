GO ?= go
DOCKER_COMPOSE ?= docker compose
E2E_PROJECT ?= miniflux-mcp-e2e
E2E_COMPOSE_FILE := .github/e2e/compose.yml
E2E_DOCKER_ARCH = $(shell docker version --format '{{.Server.Arch}}')
E2E_PLATFORM = $(if $(filter arm64 aarch64,$(E2E_DOCKER_ARCH)),linux/arm64,linux/amd64)
E2E_COMPOSE = env DOCKER_DEFAULT_PLATFORM=$(E2E_PLATFORM) $(DOCKER_COMPOSE) -p $(E2E_PROJECT) -f $(E2E_COMPOSE_FILE)
E2E_VERSION ?= e2e
VERSION ?= $(shell git describe --tags --always --dirty)
REVISION ?= $(shell git rev-parse --short HEAD)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS = -X main.Version=$(VERSION) -X main.Revision=$(REVISION) -X main.BuildDate=$(BUILD_DATE)

.PHONY: build test lint e2e

build:
	$(GO) build -ldflags "$(LDFLAGS)" ./...

test:
	$(GO) test ./...

lint:
	golangci-lint run

e2e:
	@set -eu; \
	test_dir=$$(mktemp -d); \
	cleanup() { \
		$(E2E_COMPOSE) down --volumes || true; \
		rm -rf "$$test_dir"; \
	}; \
	trap cleanup EXIT INT TERM; \
	if ! $(E2E_COMPOSE) up -d --wait; then \
		$(E2E_COMPOSE) logs; \
		exit 1; \
	fi; \
	$(GO) build -ldflags "-X main.Version=$(E2E_VERSION)" -o "$$test_dir/miniflux-mcp" .; \
	if ! MCP_SERVER_PATH="$$test_dir/miniflux-mcp" \
		EXPECTED_SERVER_VERSION=$(E2E_VERSION) \
		MINIFLUX_URL=http://localhost:8080 \
		MINIFLUX_USERNAME=admin \
		MINIFLUX_PASSWORD=test123 \
		$(GO) test -tags=e2e -v -timeout=2m ./e2e; then \
		$(E2E_COMPOSE) logs; \
		exit 1; \
	fi
