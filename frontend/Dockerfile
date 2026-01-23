# Stage 1: Generate protobuf code (TypeScript only)
FROM ubuntu:22.04 AS generator

ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /app

# Install minimal dependencies for code generation
RUN apt-get update && \
    apt-get install -y curl wget ca-certificates && \
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Install Buf CLI
RUN curl -sSL "https://github.com/bufbuild/buf/releases/download/v1.28.1/buf-Linux-x86_64" -o /usr/local/bin/buf && \
    chmod +x /usr/local/bin/buf

# Copy proto files and buf config
COPY proto/ ./proto/
COPY buf.gen.frontend.yaml ./buf.gen.yaml
COPY buf.work.yaml ./

# Copy frontend package.json for protoc plugins
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install

# Generate TypeScript code only
RUN export PATH=$PATH:./frontend/node_modules/.bin && buf generate

# Stage 2: Build frontend
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY frontend/package*.json ./
RUN npm install

# Copy source code
COPY frontend/ .

# Copy generated TypeScript files from generator stage
COPY --from=generator /app/frontend/src/gen ./src/gen

# Build the application
RUN npm run build

# Stage 3: Production with nginx
FROM nginx:alpine

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
