# Variables
GO := $(shell pwd)/.go-sdk/bin/go
export GOROOT := $(shell pwd)/.go-sdk
BIN_DIR := $(shell pwd)/bin
FRONTEND_DIR := frontend
SERVER_CMD := cmd/server/main.go
PID_FILE := .server.pid
FRONTEND_PID_FILE := .frontend.pid

# Check if Go SDK exists
ifeq (,$(wildcard $(GO)))
    $(error "Go SDK not found at $(GO). Please check your setup.")
endif

.PHONY: all
all: build

# Install dependencies
.PHONY: install-deps
install-deps:
	@echo "Installing backend tools..."
	@mkdir -p $(BIN_DIR)
	@export GOBIN=$(BIN_DIR) && $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@export GOBIN=$(BIN_DIR) && $(GO) install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	@echo "Installing frontend dependencies..."
	@cd $(FRONTEND_DIR) && npm install

# Generate code
.PHONY: generate
generate:
	@echo "Generating code..."
	@# We need both $(BIN_DIR) (for Go plugins) and frontend/node_modules/.bin (for TS plugins) in PATH
	@export PATH=$(BIN_DIR):$(shell pwd)/$(FRONTEND_DIR)/node_modules/.bin:$$PATH && buf generate

# Build everything
.PHONY: build
build: generate
	@echo "Building backend..."
	@$(GO) build -o $(BIN_DIR)/server $(SERVER_CMD)
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && npm run build

# Start services
.PHONY: start
start: build
	@echo "Starting backend server..."
	@# Check if port 8080 is free
	@if lsof -i :8080 >/dev/null; then \
		echo "Error: Port 8080 is already in use. Run 'make stop' first."; \
		exit 1; \
	fi
	@# Run the binary directly
	@export GOROOT=$(shell pwd)/.go-sdk && ./bin/server > server.log 2>&1 & echo $$! > $(PID_FILE)
	@echo "Server started with PID $$(cat $(PID_FILE))"
	@echo "Starting frontend dev server..."
	@cd $(FRONTEND_DIR) && npm run dev -- --host 0.0.0.0 > ../frontend.log 2>&1 & echo $$! > $(FRONTEND_PID_FILE)
	@echo "Frontend started with PID $$(cat $(FRONTEND_PID_FILE))"

# Stop services
.PHONY: stop
stop:
	@# Try to kill server from PID file
	@if [ -f $(PID_FILE) ]; then \
		echo "Stopping server (PID $$(cat $(PID_FILE)))..."; \
		kill $$(cat $(PID_FILE)) 2>/dev/null || echo "PID $$(cat $(PID_FILE)) not found or already stopped."; \
		rm $(PID_FILE); \
	fi
	@# Fallback: Check if port 8080 is still in use and kill usage
	@PID_PORT=$$(lsof -t -i :8080 2>/dev/null); \
	if [ -n "$$PID_PORT" ]; then \
		echo "Found process $$PID_PORT on port 8080. Killing..."; \
		kill -9 $$PID_PORT 2>/dev/null || true; \
	fi
	@# Try to kill frontend from PID file
	@if [ -f $(FRONTEND_PID_FILE) ]; then \
		echo "Stopping frontend (PID $$(cat $(FRONTEND_PID_FILE)))..."; \
		kill $$(cat $(FRONTEND_PID_FILE)) 2>/dev/null || echo "PID $$(cat $(FRONTEND_PID_FILE)) not found or already stopped."; \
		rm $(FRONTEND_PID_FILE); \
	fi
	@# Fallback: Check if port 5173 is still in use (Vite default)
	@PID_FRONTEND=$$(lsof -t -i :5173 2>/dev/null); \
	if [ -n "$$PID_FRONTEND" ]; then \
		echo "Found process $$PID_FRONTEND on port 5173. Killing..."; \
		kill -9 $$PID_FRONTEND 2>/dev/null || true; \
	fi
	@echo "Services stopped."

# Clean artifacts
.PHONY: clean
clean: stop
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)/server
	@rm -rf $(FRONTEND_DIR)/dist
	@rm -rf server.log frontend.log
	@# Optional: clean generated code if desired, but often it's better to keep unless doing a full reset
	@# rm -rf gen/go frontend/src/gen
	@echo "Clean complete."
