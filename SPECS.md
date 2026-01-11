# Project Specifications: GalleryBlue

## 1. Overview
**GalleryBlue** is a modern web application designed with strict type safety and high performance in mind. It leverages **ConnectRPC** to bridge a **Go** backend and a **TypeScript/React** frontend, ensuring a seamless, type-safe development experience from database to UI.

## 2. Technology Stack

### Frontend & Web Link
*   **Language**: TypeScript
    *   *Rationale*: Strict typing matches Go's safety and prevents runtime errors.
*   **Framework**: React
    *   *Rationale*: Vast ecosystem and excellent integration with data-fetching libraries.
*   **Build Tool**: Vite
    *   *Rationale*: Superior performance and developer experience compared to Webpack/CRA.
*   **API Client**: TanStack Query (React Query)
    *   *Rationale*: Robust state management for async server state (caching, loading, retries).
*   **Transport**: ConnectRPC (Web)
    *   *Rationale*: Enables direct gRPC-compatible calls from the browser without requiring a simplified HTTP/JSON translation layer or complex Envoy proxies.

### Backend
*   **Language**: Go (Golang)
    *   *Rationale*: High performance, concurrency support, and native Protobuf handling.
*   **Framework**: Connect-Go
    *   *Rationale*: Simple, reliable, and interoperable gRPC support for Go.

### Interface Definition
*   **IDL**: Protocol Buffers (Protobuf)
*   **Tooling**: Buf
    *   *Rationale*: Modern, fast, and user-friendly Protobuf toolkit for linting, formatting, and generation.

---

## 3. Architecture & Project Structure

The project follows a monorepo-style layout where the API definition ensures consistency between the frontend and backend.

```text
/galleryblue
├── buf.gen.yaml       # Buf configuration for code generation
├── go.mod             # Go module definition
├── proto/             # API Definitions (Source of Truth)
│   └── users/
│       └── v1/
│           └── user.proto
├── gen/               # Auto-generated code (DO NOT EDIT)
│   ├── go/            # Generated Go structs & interfaces
│   └── ts/            # Generated TypeScript types & hooks
└── frontend/          # React Application
```

## 4. Development Workflow

1.  **Define API**: Schemas are written in `proto/`.
2.  **Generate Code**: Run `buf generate` to update `gen/` with Go and TypeScript artifacts.
    *   **Backend Output**: `gen/go` (Module: `github.com/mzzz-zzm/galleryblue/gen/go`)
    *   **Frontend Output**: `gen/ts` (TypeScript types, Connect Client, and TanStack Query hooks).
3.  **Implement Backend**: Fulfill the generated Go interfaces in the backend service using the generated structs.
4.  **Implement Frontend**: Import hooks from `gen/ts` (e.g., `useGetUser`) and build UI components.

## 5. Configuration Details

### `buf.gen.yaml`
The configuration is set to version `v2` and defines plugins for both Go and TypeScript generation.
*   **Managed Mode**: Enabled, with `go_package_prefix` set to `github.com/mzzz-zzm/galleryblue/gen/go`.
*   **Go Plugins**: `protoc-gen-go`, `protoc-gen-connect-go`.
*   **TypeScript Plugins**: `protoc-gen-es`, `protoc-gen-connect-es`, `protoc-gen-connect-query`.

### Frontend Usage Pattern
Components should utilize the generated hooks for data fetching, eliminating manual `fetch` calls or manual type definitions for API responses.

```typescript
// Example Usage
import { useGetUser } from "../gen/ts/users/v1/user-UserService_connectquery";

export const UserProfile = ({ userId }: { userId: string }) => {
  const { data, isLoading } = useGetUser({ id: userId });
  if (isLoading) return <Spinner />;
  return <h1>{data?.name}</h1>;
};
```

## 6. Next Steps (Setup)
1.  Initialize Go module: `go mod init github.com/mzzz-zzm/galleryblue`.
2.  Setup `frontend` directory (e.g., `npm create vite@latest frontend`).
3.  Create `proto` directory and initial definitions.
4.  Run `buf generate`.
