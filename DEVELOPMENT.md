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

## Database Management

The database schema is defined in `init.sql`. It runs automatically when the PostgreSQL container starts for the first time.

### When to Reset the Database

Reset the database when you:
- Modify `init.sql` (add/change tables)
- Want to clear all data and start fresh
- See errors like `relation "xxx" does not exist`

### How to Reset

```bash
# Reset database (clears all data!)
make docker-db-reset
```

> **Warning**: This deletes all data including users and images.

### Adding New Tables

1. Edit `init.sql` to add your table:
   ```sql
   CREATE TABLE IF NOT EXISTS my_table (
       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
       name VARCHAR(255) NOT NULL,
       created_at TIMESTAMP DEFAULT NOW()
   );
   ```

2. Reset the database:
   ```bash
   make docker-db-reset
   ```

3. Register a new user (old data is cleared)

---

## Adding a New gRPC API

### Step 1: Define the Proto

Edit `proto/users/v1/user.proto`:

```protobuf
message DeleteUserRequest {
  string id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}

service UserService {
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}
```

### Step 2: Regenerate Code

```bash
make docker-rebuild
```

### Step 3: Implement Backend Handler

Edit `internal/handlers/user.go`:

```go
func (s *UserServer) DeleteUser(
    ctx context.Context,
    req *connect.Request[usersv1.DeleteUserRequest],
) (*connect.Response[usersv1.DeleteUserResponse], error) {
    err := db.DeleteUser(ctx, req.Msg.Id)
    if err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }
    return connect.NewResponse(&usersv1.DeleteUserResponse{Success: true}), nil
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

### Step 5: Add Frontend Page

Create `frontend/src/pages/DeletePage.tsx` and add route in `App.tsx`.

### Step 6: Rebuild and Test

```bash
make docker-rebuild
make docker-logs-backend
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
| `frontend/src/components/` | Reusable components |
| `init.sql` | Database schema |

---

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make docker-up` | Start production services |
| `make docker-dev` | Start with hot-reload |
| `make docker-down` | Stop all services |
| `make docker-rebuild` | Rebuild after code changes |
| `make docker-db-reset` | Reset database (clears data!) |
| `make docker-logs` | View all logs |
| `make docker-clean` | Remove containers/volumes/images |
| `make help` | Show all commands |

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| "relation does not exist" | Run `make docker-db-reset` |
| "undefined" types after proto change | Run `make docker-rebuild` |
| Backend not finding handler | Check `main.go` registration |
| Port already in use | Run `make docker-down` first |
| HTTP 413 (file too large) | Check nginx `client_max_body_size` |
