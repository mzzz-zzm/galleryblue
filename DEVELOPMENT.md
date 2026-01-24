# Development Guide

This guide covers how to develop and add new features to GalleryBlue.

## Quick Start

```bash
# Start development environment
make docker-dev

# Or use production mode
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

---

## Adding a New gRPC API

### Step 1: Define the Proto

Edit `proto/users/v1/user.proto`:

```protobuf
// Add new message types
message DeleteUserRequest {
  string id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}

// Add RPC to the service
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);  // NEW
}
```

### Step 2: Regenerate Code

```bash
make docker-rebuild
```

This generates:
- `gen/go/users/v1/` - Go types and Connect handlers
- `frontend/src/gen/users/v1/` - TypeScript types and React hooks

### Step 3: Implement Backend Handler

Edit `internal/handlers/user.go`:

```go
func (s *UserServer) DeleteUser(
    ctx context.Context,
    req *connect.Request[usersv1.DeleteUserRequest],
) (*connect.Response[usersv1.DeleteUserResponse], error) {
    userID := req.Msg.Id
    
    err := db.DeleteUser(ctx, userID)
    if err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }
    
    return connect.NewResponse(&usersv1.DeleteUserResponse{
        Success: true,
    }), nil
}
```

### Step 4: Add Database Helper (if needed)

Edit `internal/db/queries.go`:

```go
func DeleteUser(ctx context.Context, userID string) error {
    _, err := DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
    return err
}
```

### Step 5: Add Frontend Integration

Create `frontend/src/pages/DeletePage.tsx`:

```tsx
import { useMutation } from '@connectrpc/connect-query';
import { deleteUser } from '../gen/users/v1/user-UserService_connectquery';
import { transport } from '../lib/transport';

export const DeletePage = () => {
    const deleteMutation = useMutation(deleteUser, { transport });
    
    const handleDelete = async (userId: string) => {
        await deleteMutation.mutateAsync({ id: userId });
    };
    
    return <button onClick={() => handleDelete("123")}>Delete</button>;
};
```

Add route in `frontend/src/App.tsx`:

```tsx
<Route path="/delete" element={<DeletePage />} />
```

### Step 6: Rebuild and Test

```bash
make docker-rebuild
make docker-logs-backend

# Test API directly
curl -X POST http://localhost:8080/users.v1.UserService/DeleteUser \
  -H "Content-Type: application/json" \
  -d '{"id": "123"}'
```

---

## Project Structure

| Directory | Purpose |
|-----------|---------|
| `proto/` | Protocol Buffer definitions |
| `gen/` | Generated code (do not edit) |
| `cmd/server/` | Backend entry point |
| `internal/handlers/` | gRPC handler implementations |
| `internal/db/` | Database connection and queries |
| `frontend/src/pages/` | React pages |
| `frontend/src/components/` | Reusable React components |
| `frontend/src/gen/` | Generated TypeScript (do not edit) |

---

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make docker-up` | Start production services |
| `make docker-dev` | Start with hot-reload |
| `make docker-down` | Stop all services |
| `make docker-rebuild` | Rebuild after code changes |
| `make docker-logs` | View all logs |
| `make docker-clean` | Remove containers/volumes |
| `make help` | Show all commands |

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| "undefined" types after proto change | Run `make docker-rebuild` |
| Backend not finding new handler | Check method is registered in `main.go` |
| Frontend hook not found | Verify `*_connectquery.ts` was generated |
| Database error | Check `init.sql` has required schema |
| Port already in use | Run `make docker-down` first |
