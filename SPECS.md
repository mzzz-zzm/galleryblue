# Project Specifications: GalleryBlue

## 1. Overview

**GalleryBlue** is a modern web application for sharing and managing images. Built with strict type safety using **ConnectRPC** to bridge a **Go** backend and a **TypeScript/React** frontend.

## 2. Technology Stack

### Infrastructure
- **Docker Compose**: Orchestrates frontend, backend, and database containers
- **nginx**: Serves frontend and proxies API requests

### Frontend
- **TypeScript** + **React** + **Vite**
- **TanStack Query**: Server state management
- **ConnectRPC Web**: Type-safe gRPC calls from browser

### Backend
- **Go** + **Connect-Go**: High-performance gRPC server
- **PostgreSQL**: Primary database
- **bcrypt**: Password hashing

### API Definition
- **Protocol Buffers** (Protobuf) + **Buf** tooling

---

## 3. Database Schema

### `users` Table
| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | Primary Key, auto-generated |
| email | VARCHAR | Unique, Not Null |
| password_hash | VARCHAR | Not Null |
| display_name | VARCHAR | Unique |
| created_at | TIMESTAMP | Default NOW() |
| updated_at | TIMESTAMP | Default NOW() |

### `images` Table (NEW)
| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | Primary Key, auto-generated |
| owner_id | UUID | Foreign Key → users.id, Not Null |
| filename | VARCHAR | Not Null |
| content_type | VARCHAR | Not Null (e.g., "image/jpeg") |
| data | BYTEA | Not Null (binary image data) |
| title | VARCHAR | Optional |
| description | TEXT | Optional |
| created_at | TIMESTAMP | Default NOW() |
| updated_at | TIMESTAMP | Default NOW() |

---

## 4. Features

### Authentication
1. **Register** (`/register`): Create account with email/password
2. **Login** (`/login`): Authenticate and receive session token
3. **Update Profile** (`/update`): Modify name/email/password (requires current password)

### Image Gallery (NEW)
1. **Upload Image**: Authenticated user uploads JPEG file
   - Max file size: 5MB
   - Supported formats: JPEG only (for now)
   - Image stored in database with owner reference

2. **View Own Images** (`/my-images`): User sees their uploaded images
   - Can view, edit, delete their images

3. **View All Images** (`/gallery`): Browse all public images
   - Shows image with uploader info
   - Users can only view others' images (no edit/delete)

4. **Edit Image** (owner only): Update title/description
5. **Delete Image** (owner only): Remove from database

---

## 5. API Services

### AuthService
```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
}
```

### UserService
```protobuf
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
}
```

### ImageService (NEW)
```protobuf
service ImageService {
  // Upload a new image (owner = authenticated user)
  rpc UploadImage(UploadImageRequest) returns (UploadImageResponse);
  
  // Get single image by ID (anyone can view)
  rpc GetImage(GetImageRequest) returns (GetImageResponse);
  
  // List all images (public gallery)
  rpc ListImages(ListImagesRequest) returns (ListImagesResponse);
  
  // List images owned by current user
  rpc ListMyImages(ListMyImagesRequest) returns (ListMyImagesResponse);
  
  // Update image metadata (owner only)
  rpc UpdateImage(UpdateImageRequest) returns (UpdateImageResponse);
  
  // Delete image (owner only)
  rpc DeleteImage(DeleteImageRequest) returns (DeleteImageResponse);
}
```

---

## 6. Authorization Rules

| Action | Who Can Perform |
|--------|-----------------|
| Upload image | Authenticated user |
| View image | Anyone (public) |
| View uploader info | Anyone (public) |
| Edit image metadata | Owner only |
| Delete image | Owner only |

---

## 7. Project Structure

```
/galleryblue
├── docker-compose.yml    # Container orchestration
├── Dockerfile            # Multi-stage build (production/development)
├── Makefile              # Development commands
├── proto/
│   └── users/v1/
│       └── user.proto    # User & Image API definitions
├── gen/                  # Generated code (DO NOT EDIT)
│   └── go/
├── internal/
│   ├── handlers/         # gRPC implementations
│   └── db/               # Database queries
├── frontend/
│   ├── src/pages/        # React pages
│   ├── src/components/   # Reusable components
│   └── src/gen/          # Generated TypeScript
└── init.sql              # Database schema
```

---

## 8. Development Workflow

1. **Start services**: `make docker-up`
2. **Define API**: Edit `proto/users/v1/user.proto`
3. **Generate code**: `make docker-rebuild`
4. **Implement backend**: Add handlers in `internal/handlers/`
5. **Implement frontend**: Create pages/components using generated hooks
6. **Test**: `make docker-logs` to debug

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed guide.
