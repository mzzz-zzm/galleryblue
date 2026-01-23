# Development container for GalleryBlue
# This container has all tools needed for development

FROM ubuntu:22.04

# Prevent interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# Copy setup script first
COPY scripts/setup.sh /tmp/setup.sh
RUN chmod +x /tmp/setup.sh && /tmp/setup.sh

# Set up Go environment
ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin
ENV GOPATH=/root/go

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
