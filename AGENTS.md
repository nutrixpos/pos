# Agent Guidance for NutrixPOS

## Project Overview
- Go 1.24 monorepo with MongoDB backend (mongodb-driver v1.16)
- Point-of-sale system for restaurants/shops: inventory, sales, products
- Active development - no backward compatibility guarantee

## Build Commands
```bash
go build ./...        # build all packages
go run ./cmd/pos      # run the CLI
```

## Architecture
- `/cmd/` - CLI entrypoints
- `/modules/` - business logic (core, hubsync modules)
- `/common/` - shared utilities (database, config, logger)
- No tests in this repo

## Database
- Use `common.GetDatabaseClient()` singleton - never create new `mongo.Connect()` connections
- Singleton pattern in `common/database.go` ensures single connection

## Common Pitfalls to Avoid

### 1. Database Connection Pattern (CRITICAL)
❌ Wrong:
```go
clientOptions := options.Client().ApplyURI(...)
ctx, cancel := context.WithTimeout(...)
client, err := mongo.Connect(ctx, clientOptions)
```
✅ Correct:
```go
client, err := common.GetDatabaseClient(logger, &config)
if err != nil {
    return err
}
ctx := context.Background()
```

### 2. Imports After Refactoring
When changing mongo.Connect to GetDatabaseClient:
- Remove: `"go.mongodb.org/mongo-driver/mongo"`, `"go.mongodb.org/mongo-driver/mongo/options"`
- Keep: `"go.mongodb.org/mongo-driver/mongo"` only if using `mongo.ErrNoDocuments`
- Add: `"github.com/nutrixpos/pos/common"` if not present

## Dependencies
- `go.mongodb.org/mongo-driver` - MongoDB driver
- `github.com/gorilla/mux` - HTTP router
- `github.com/spf13/cobra` + `viper` - CLI framework
- `github.com/rs/zerolog` - logging (not used everywhere)

## Entities
- `Material`, `Component` and `Inventory Item` are the same entity
- `Product` and `Recipe` are the same entity