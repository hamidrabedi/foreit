Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

ROLE: You are the delegated coder: edit the files yourself; do not call any other agent. Use US spelling in comments.

# Keep exported signatures backward compatible; move new parameters to new functions

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
A read-only copy of master is at /home/hamid/Other/projects/foreit-wt/verify-master/forge; read the old signatures there. Never edit verify-master.

## Rule
The refactor changed exported function signatures, which breaks user code and generated projects. For each function below: restore the **original exported signature** (behaving as it did, using sensible defaults for the new information) and give the new signature a **new exported name**. Update in-repo callers to the new names where they have the new information. Mark nothing deprecated unless stated. gofmt; functions < 40 lines; naming after behaviour.

| Original (restore exactly) | New name for the current signature | Default when called through the original |
|---|---|---|
| `cli/core.NewContext() *Context` | `NewContextWithConfig(cfg *config.Config) *Context` | `config.NewConfig()` |
| `db/migrate.Generate(name, modelsDir, migrationsDir string) error` | `GenerateForDriver(name, modelsDir, migrationsDir string, driver core.Driver) error` | driver from `config.NewConfig().GetDriver()` |
| `db/migrate.NewGenerator(modelsDir, migrationsDir string) (*Generator, error)` | `NewGeneratorForDriver(modelsDir, migrationsDir string, driver core.Driver)` | same |
| `db/migrate/generate.NewMigrationGeneratorWithDefaults(modelsDir, migrationsDir string)` | `NewMigrationGeneratorForDriver(modelsDir, migrationsDir string, driver core.Driver)` | same |
| `orm.BuildInsertSQL(instance interface{}, tableName string)` | `BuildInsertSQLForPK(instance interface{}, tableName, pkColumn string, placeholder ...func(int) string)` | pk column `"id"`, PostgreSQL `$N` placeholders (the original behaviour) |
| `orm.BuildBulkInsertSQL(instances []interface{}, tableName string)` | `BuildBulkInsertSQLForPK(instances []interface{}, tableName, pkColumn string, placeholder ...func(int) string)` | same |
| `orm.ExecuteInsert(ctx, database *db.DB, sql string, args []interface{}) (int64, error)` | `ExecuteInsertTx(ctx, dbtx DBTX, d dialect.Dialect, sql string, args []interface{})` | `dbtx` = the DB's `*sql.DB`, `d` = the DB's dialect |
| `orm.ExecuteBulkInsert(ctx, database *db.DB, ...)` | `ExecuteBulkInsertTx(...)` | same |
| `orm.ExecuteUpdate(ctx, database *db.DB, ...)` | `ExecuteUpdateTx(...)` | same |
| `orm.ExecuteDelete(ctx, database *db.DB, ...)` | `ExecuteDeleteTx(...)` | same |

Leave the variadic additions (`BuildUpdateSQL`, `BuildDeleteSQL`, `QueryExpr.ToSQL` gained `placeholder ...func(int) string`) as they are: existing calls still compile.

The in-repo code that loads config once (the CLI root, the makemigrations command) must keep using the explicit-config/driver functions, so config is still read once there.

## Tests
For each restored original, one test that calls it with the original arguments (compile-level compatibility) and checks the result matches the new function with the documented default (for Execute*, use in-memory SQLite through `db` helpers).

## Acceptance (from the module root)
    gofmt -l .                                        # empty
    go build ./... && (cd ../examples/ecommerce && go build ./...)
    go vet ./...
    staticcheck ./cli/... ./db/... ./orm/...          # no output
    go test -race ./cli/... ./db/... ./orm/... -count=1
    grep -rn 'config.NewConfig()' --include=*.go cli/core/registry.go cli/commands/migrations   # the CLI root still loads it once

Report (max 25 lines): each restored signature and its new counterpart, callers updated, each command result.
