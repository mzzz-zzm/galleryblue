# Base image with all development tools
FROM ubuntu:22.04 AS base

# Prevent interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# Copy setup script first
COPY scripts/setup.sh /tmp/setup.sh
RUN chmod +x /tmp/setup.sh && /tmp/setup.sh

# Set up Go environment
ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin
ENV GOPATH=/root/go

# -----------------------------------------------------------
# Production build stage
# -----------------------------------------------------------
FROM base AS production

# Copy frontend package.json first for npm caching
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install

# Copy ALL source files
COPY . .

# Generate protobuf code (after COPY so it's not overwritten)
RUN export PATH=$PATH:/usr/local/go/bin:/root/go/bin:./frontend/node_modules/.bin && buf generate

# Update go dependencies (generated code now exists)
RUN go mod tidy && go mod download

# Build the backend
RUN go build -o /app/bin/server ./cmd/server/main.go

# Expose ports
EXPOSE 8080

# Default command
CMD ["./bin/server"]

# -----------------------------------------------------------
# Development build stage (with hot-reload)
# -----------------------------------------------------------
FROM base AS development

# Install air for Go hot-reload
RUN go install github.com/cosmtrek/air@latest

# Copy go.mod first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy frontend package.json for npm caching
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install

# Copy the rest of the source code
COPY . .

# Generate protobuf code
RUN export PATH=$PATH:/usr/local/go/bin:/root/go/bin:./frontend/node_modules/.bin && buf generate

# Expose ports (8080 for backend, 5173 for Vite dev server)
EXPOSE 8080 5173

# Start script for development
COPY scripts/dev-start.sh /usr/local/bin/dev-start.sh
RUN chmod +x /usr/local/bin/dev-start.sh

CMD ["/usr/local/bin/dev-start.sh"]
