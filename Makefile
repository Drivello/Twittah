# Variables
SERVICES = auth user tweet gateway
SERVICES_DIR = services

.PHONY: all build test lint tidy run docker-compose all-up

all: build

build: tidy
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $(SERVICES_DIR)/$$service && go build ./... ; \
		cd - > /dev/null ; \
	done

test:
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd $(SERVICES_DIR)/$$service && go test ./... ; \
		cd - > /dev/null ; \
	done

lint:
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		export PATH=$$PATH:$$HOME/go/bin; \
	fi
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd $(SERVICES_DIR)/$$service && golangci-lint run ./... ; \
		cd - > /dev/null ; \
	done

tidy:
	@for service in $(SERVICES); do \
		echo "Running go mod tidy in $$service..."; \
		cd $(SERVICES_DIR)/$$service && go mod tidy ; \
		cd - > /dev/null ; \
	done

envs:
	@for service in $(SERVICES); do \
		if [ ! -f "$(SERVICES_DIR)/$$service/.env" ]; then \
			echo "Creating empty .env for $$service..."; \
				touch "$(SERVICES_DIR)/$$service/.env"; \
		fi \
	done

run:
	docker-compose up --build

docker-compose:
	docker-compose up

all-up: envs tidy build test lint run
	@echo "✅ Project compiled, tested, lint OK and containers running 🚀"
