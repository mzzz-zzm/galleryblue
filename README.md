# GalleryBlue

A modern, type-safe web application built with Go (ConnectRPC) and React (Vite + Connect-Query).

## Prerequisites

- **Go**: v1.21+ (or use the provided local SDK environment)
- **Node.js**: v18+
- **npm**: v9+

## How to Run

You need to run the **Backend** and **Frontend** in separate terminals.

### 1. Start the Backend

The backend runs on port `8080`.

```bash
# If using the local SDK environment set up by the agent:
export GOROOT=$(pwd)/.go-sdk
export PATH=$GOROOT/bin:$PATH

# Run the server
go run cmd/server/main.go
```

You should see: `Server executing on localhost:8080`

### 2. Start the Frontend

The frontend runs on port `5173` (by default).

```bash
cd frontend
npm install  # Only needed the first time
npm run dev
```

You should see: `Local: http://localhost:5173/`

### 3. Access the App

Open your browser to: **[http://localhost:5173](http://localhost:5173)**

- You should see the User Profile for "Jane Doe".
- **Debug UI**: Pass the `userId` prop or check the debug section to verify connection status.
- **Manual Fetch**: Use the "Test Manual Fetch" button to verify if the browser can reach the backend via the Vite Proxy.

## How to Stop

To stop the application:

1.  Go to the terminal running the **Backend**.
2.  Press `Ctrl + C`.
3.  Go to the terminal running the **Frontend**.
4.  Press `Ctrl + C`.

## Project Structure

- `cmd/server/`: Backend entry point.
- `frontend/`: React application (Vite).
- `proto/`: Protocol Buffer definitions (API Schema).
- `gen/`: Generated Go and TypeScript code (Do not edit manually).
