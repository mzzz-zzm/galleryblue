#!/bin/bash
# Development start script - runs both backend and frontend

set -e

echo "🚀 Starting GalleryBlue Development Environment"

# Wait for database
echo "⏳ Waiting for database..."
until pg_isready -h db -U galleryblue; do
    sleep 1
done
echo "✅ Database is ready"

# Regenerate protobuf code (in case of changes)
echo "📦 Generating protobuf code..."
export PATH=$PATH:/usr/local/go/bin:/root/go/bin:/app/frontend/node_modules/.bin
buf generate

# Start frontend dev server in background
echo "🎨 Starting frontend dev server..."
cd /app/frontend
npm run dev -- --host 0.0.0.0 &
FRONTEND_PID=$!

# Start backend with hot-reload
echo "🔧 Starting backend with hot-reload..."
cd /app
air -c .air.toml &
BACKEND_PID=$!

# Handle shutdown
trap "kill $FRONTEND_PID $BACKEND_PID 2>/dev/null" EXIT

echo "✅ Development environment is running!"
echo "   Frontend: http://localhost:5173"
echo "   Backend:  http://localhost:8080"

# Keep container running
wait
