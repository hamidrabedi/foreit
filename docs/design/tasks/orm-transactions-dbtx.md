Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

# ORM managers and querysets can run inside a transaction

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit forge/orm (non-generated files) and forge/db/transaction.go, plus tests.

## Principles (mandatory)
- Tests FIRST (a test that fails today). gofmt; functions < 40 lines; errors wrapped with %w; no new package-level mutable state; naming after behaviour.
- Existing behaviour and exported signatures keep working; this is additive.

## Facts (verified on master)
- `Manager[T]` holds `db *db.DB` (forge/orm/manager.go:17); `ExecuteInsert/ExecuteBulkInsert/ExecuteUpdate/ExecuteDelete` take `*db.DB` (forge/orm/manager_helpers.go ~408-538).
- `BaseQuerySet[T].db` is `interface{}` and `getDB` (forge/orm/queryset.go:179) returns `*sql.DB` through `GetSQLDB` (forge/orm/types.go:26), which only accepts `*db.DB` or `*sql.DB`. `getDialect` works the same way.
- `db.Tx` (forge/db/transaction.go) wraps `*sql.Tx` and keeps its `*DB`; `(*DB).WithTx(ctx, fn func(*Tx) error)` exists.
- So ORM calls always use the pool; nothing done through a manager can join a `db.Tx`, and Django-style `transaction.atomic` is impossible.

## Change
1. orm: define
   ```go
   // DBTX is satisfied by *sql.DB and *sql.Tx.
   type DBTX interface {
       ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
       QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
       QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
   }
   ```
   Replace `getDB` / `GetSQLDB` results used for execution with a `DBTX` resolved from the stored connection: `*db.DB` → its `*sql.DB`; `*sql.DB` → itself; `*db.Tx` → its `*sql.Tx`; `*sql.Tx` → itself. The dialect for a `*db.Tx` comes from its parent `*DB` (add an unexported accessor or an exported `(*Tx).DB()` in db/transaction.go).
   Make the `Execute*` helpers take `DBTX` (plus whatever they need from the dialect) instead of `*db.DB`.
2. `func (m *Manager[T]) WithTx(tx *db.Tx) *Manager[T]` returns a copy bound to the transaction; querysets created from that copy use the transaction too. `Manager.db` becomes a small internal connection value that holds either a `*db.DB` or a `*db.Tx`; `SetDB(*db.DB)` keeps working.
3. Tests on in-memory SQLite (see existing orm tests that open SQLite): inside `database.WithTx`, create a row with `manager.WithTx(tx)`, read it back through the same tx-bound manager (visible), return an error so the tx rolls back, then the pool-bound manager's `Count` is 0. A second test commits and sees the row. A queryset from the tx-bound manager (`Filter(...).All`) sees uncommitted rows.

## Acceptance (from the module root)
    gofmt -l orm db                                   # empty
    go build ./... && (cd ../examples/ecommerce && go build ./...)
    go vet ./orm/... ./db/... ./admin/... ./api/...
    staticcheck ./orm/... ./db/...                    # no output
    go test -race ./orm/... ./db/... ./admin/... ./api/... -count=1

Report (max 25 lines): the failing test output before, signatures changed, each command result.
