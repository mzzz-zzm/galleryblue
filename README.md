# GalleryBlue

A modern, type-safe web application built with Go (ConnectRPC) and React (Vite + TanStack Query).

## Quick Start

```bash
# Start all services with Docker
make docker-up

# Open in browser
open http://localhost:3000
```

## Prerequisites

- **Docker** and **Docker Compose** (recommended)
- Or: Go v1.21+, Node.js v18+, PostgreSQL 16+

## How to Run

### Option A: Docker (Recommended)

```bash
# Production mode
make docker-up

# Development mode with hot-reload
make docker-dev

# Stop services
make docker-down
```

**Services:**
| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Backend | http://localhost:8080 |
| Database | localhost:5432 |

### Option B: Local Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for local setup with Go SDK.

## Available Commands

```bash
make docker-up        # Start production
make docker-dev       # Start with hot-reload
make docker-down      # Stop services
make docker-rebuild   # Rebuild after changes
make docker-logs      # View logs
make help             # Show all commands
```

## Project Structure

```
├── cmd/server/          # Backend entry point
├── internal/
│   ├── handlers/        # gRPC handler implementations
│   └── db/              # Database connection & queries
├── proto/               # Protocol Buffer definitions
├── gen/                 # Generated code (do not edit)
├── frontend/
│   ├── src/pages/       # React pages
│   ├── src/components/  # Reusable components
│   └── src/gen/         # Generated TypeScript (do not edit)
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Service orchestration
└── Makefile             # Development commands
```

## Adding New Features

See [DEVELOPMENT.md](DEVELOPMENT.md) for step-by-step guide on adding new gRPC APIs.

## Tech Stack

- **Backend**: Go, ConnectRPC, PostgreSQL
- **Frontend**: React, Vite, TanStack Query
- **Infrastructure**: Docker, nginx
