# Agent Guidance for NutrixPOS

## Project Overview
- Go 1.26 monorepo with MongoDB backend (mongodb-driver v1.16)
- Point-of-sale system for restaurants/shops: inventory, sales, products
- Active development - no backward compatibility guarantee

## Build Commands
```bash
go run ./cmd/pos      # run the CLI
npm install           # installs dev tooling (husky) and git hooks
npm run test          # run tests against the embedded FerretDB backend
npm run test-race
npm run test-mongo    # run against external MongoDB (TEST_MONGO_URI, default mongodb://127.0.0.1:27017)
npm run cover         # scoped coverage + enforce >=80% (.husky/coverage-threshold.mjs)
npm run coverage-check # re-check the existing coverage.out against the threshold
npm run vuln          # govulncheck (v1.8.0), scans the active Go toolchain; needs network for the vuln DB
npm run lint          # golangci-lint v2.14.0 (config: .golangci.yml); runs via `go run`, needs network first time
npm run build vet fmt-check cover  # other targets, see package.json scripts
```

The npm scripts use an explicit package list instead of `./...`, because a local
docker-compose Mongo volume at `data/mongo` can be unreadable and makes the
`./...` pattern fail. Add any new top-level Go directory to the package list in
`package.json` (and `.husky/test-mongo.mjs`).

## Git hooks
- Husky (from root `package.json`) installs hooks on `npm install`.
- Hooks are thin Node scripts (no shell logic): `.husky/pre-commit` runs
  `.husky/pre-commit.mjs`, `.husky/pre-push` runs `.husky/pre-push.mjs`.
- `pre-commit` (Go changes only): `gofmt -w` + `git add` (auto-stage), then `npm run vet test`.
  It aborts with an error if a staged Go file also has unstaged changes, to preserve partial staging.
- `pre-push`: `npm run test-race` then `npm run lint` (golangci-lint) then `npm run vuln` (govulncheck).
- Bypass with `git commit -n` / `git push --no-verify`, or `HUSKY=0`.
- Keep the package lists in `package.json` and the hook scripts in sync.

## Architecture
- `/cmd/` - CLI entrypoints
- `/modules/` - business logic (core, hubsync modules)
- `/common/` - shared utilities (database, config, logger)
- `/internal/testutil/` - test-only DB harness (never imported by production code)

## Testing
- Use `internal/testutil.NewTestEnv(t, testutil.BackendFromEnv())` for an isolated database:
  embedded FerretDB with a temp dir by default, or external Mongo when `TEST_MONGO_URI` is set.
- The DB client is a process-wide singleton: never call `t.Parallel()` in tests that use a database,
  and create at most one `TestEnv` per test.
- Integration tests in `modules/core/services` cover materials, purchase orders (including
  rollback via the `testHookAfterPOUpdate` seam), sales and order display ids, on both backends.
- Coverage: the goal is **>=80%** statement coverage over the in-scope packages
  (`common/config`, `common/helpers`, `modules/auth/middlewares`, `modules/core/handlers`,
  `modules/core/middlewares`, `modules/core/services`). Scope is set with `go test -coverpkg=...`, and the hard-to-test
  files `modules/core/services/receipt.go`, `modules/core/services/notification_singleton.go`
  and `modules/core/handlers/notifications.go` are excluded (see `.husky/coverage-threshold.mjs`
  and `codecov.yml`). `npm run cover` enforces the threshold; CI runs the same check.
- CI gates: gofmt, `go vet`, `go build`, `go test -race` (both backends), the 80% coverage
  threshold, golangci-lint, govulncheck.

## Database
- Use `common.GetDatabaseClient()` singleton - never create new `mongo.Connect()` connections
- Singleton pattern in `common/database.go` ensures single connection
- Two backends, selected by `databases[0].type` in config:
  - `type: mongo` - connects to an external/centralized MongoDB server (host/port)
  - `type: ferret` - starts an embedded FerretDB (SQLite, pure Go) in-process; no database service install required. Data is stored in `databases[0].file_path` (default `./data/db`)
- The app code always talks to a `*mongo.Client` regardless of backend
- FerretDB embedded caveat: aggregation `$project` with expression operators (`$size`, `$slice`) is unsupported - fetch the doc with `FindOne` and do the work in Go instead

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

## API http schema
When calling the backend api from the frontend vue app, make sure to include the VITE_APP_BACKEND_HOST and VITE_APP_MODULE_CORE_API_PREFIX env vars in the request path