# Forge refactor audit

Baseline: `origin/master` @ 90e6bef plus the gofmt pass in #214 (2026-09-13).
Scope: the Go module `forge/` and the test module `tests/`. The admin UI (`forge/admin/ui/web`) is out of scope.

This document has two jobs. **Part A** is the playbook: how to run the refactor and what to look for in every PR.
**Part B** is the findings list, ranked by severity, with one planned action per finding.

## How the findings were produced

| Source | What it did |
|---|---|
| Tools (Claude ran them) | `deadcode -test ./...` (Go 1.26.6), `deadcode` rooted at `examples/ecommerce`, `dupl -t 80`, `gocyclo`, `staticcheck 2026.2.1`, importer counts per package |
| Survey agents (Sonnet) | one pass over the data layer, one over web/admin/auth |
| Second opinions | agy (Gemini) on web/admin/auth and Nemotron 3 Ultra on the data layer, asked to validate claims and to report what was missed (Part B §9) |
| Claude | read the code behind every CRITICAL/MAJOR claim before listing it |

Legend for every finding:

- Severity: **CRITICAL** = wrong results, data loss or a security hole. **MAJOR** = a bug or design flaw users will hit. **MINOR** = a smell.
- Status: **V** = verified by reading the code. **R** = reported by an agent and spot-checked only; confirm it before acting on it.

Raw tool output lives next to this file in `audit/`: `deadcode-forge-tests.txt` (901 functions unreachable even from tests),
`deadcode-from-ecommerce.txt`, `dupl-clones.txt` (34 clone groups) and `gocyclo-over-25.txt`.

---

# Part A: playbook

## A1. Ground rules

1. **One concern per PR, one package per PR where possible.** "Delete dead package X", "merge limiter A into B", "fix bug Y". Never mix a behaviour change with a move or a rename.
2. **Delete before you refactor.** Removing dead code first shrinks every later diff and removes false duplicates from the picture.
3. **Pin behaviour before you change it.** Before merging duplicates or rewriting a component, add a test that captures current behaviour, including the SQL string where the ORM is involved. The refactor PR must keep that test green (or change it deliberately, with the reason in the PR body).
4. **Incomplete features are flagged, not finished.** For every STUB finding (Part B §3), the refactor PR makes the stub fail loudly: return `errors.ErrNotImplemented` (or a typed error) and add a `// NOT IMPLEMENTED:` comment plus a row in `docs/design/feature-gaps.md`. No silent `return nil`, no clone-and-ignore.
5. **Public API breaks are allowed, but listed.** The framework has no external users yet; still, every removed or renamed exported identifier goes in the PR body under "Breaking".
6. **Gates for every PR**: `gofmt -l` empty, `go vet ./...`, `staticcheck ./...` (2026.2.1), `go test -race ./...` in `forge`, `go build ./...` in `examples/ecommerce`, and the relevant `tests/` integration package. Re-run `deadcode -test ./...` and paste the before/after count.
7. **Delegation**: each item in Part B's plan (§10) is sized to be one agy/codex task prompt under `docs/design/tasks/`, verified by Claude, one PR each, 2-3 delegates at a time at most.

## A2. What to look for (reviewer checklist)

Use this on the code you touch, and on every refactor PR.

**Dead code**
- Package with no importers in `forge/`, `tests/`, `examples/` (check `cli/` templates that generate imports before deleting).
- Function listed in `audit/deadcode-forge-tests.txt` and not reachable from `examples/ecommerce`.
- Interface with exactly one implementation and no test double: inline it.
- "Facade" packages that only re-export another package (`identity/password.go`, `db/migrate/migrate_impl.go`).

**Duplication**
- Two implementations of one concept (limiters, caches, error types, viewsets, migration runners). Pick the one with more tests and callers; turn the other into a thin adapter first, then delete it.
- Copy-pasted blocks in `audit/dupl-clones.txt`, especially the `identity/repository` scan/query helpers.
- Near-identical files (diff shows only a few changed lines): `db/migrate/{execute,verify}/checksum.go`, `tests/{helpers,testhelpers}/*`.

**Silent failure**
- `return nil` or `return clone` from a method that did not do its job.
- `_ = err`, `x, _ :=` on anything that can fail at startup (`var XObjects, _ = orm.NewManager...` in generated code).
- `continue` on an error inside a loop that builds SQL or config.
- No-op writers or exporters selected by configuration.

**Guessing instead of metadata**
- Column, table or FK names derived by string surgery (`+ "_id"`, `TrimSuffix(table, "s")`, `Name+"ID"`).
- Case-insensitive field matching that later emits the caller's spelling into SQL.

**Global state**
- Package-level `var` that is mutated after init (registries, `DefaultSite`, caches).
- `config.NewConfig()` called anywhere below `main`/CLI entry points.

**SQL safety**
- `fmt.Sprintf` that puts anything other than an escaped identifier or a placeholder into SQL.
- Text rewriting of finished SQL (`strings.Replace` on placeholders, regex rebinds).

**Security**
- Request data flowing into writes without an allow-list of editable fields.
- Auth flows whose response or timing differs between "unknown user", "wrong password" and "inactive/locked".
- Hard-coded permissive CORS (`*` or echo-all).

**Tests worth deleting or rewriting**
- Tests that only assert a type implements an interface (the compiler already checks this with `var _ I = (*T)(nil)`).
- Tests with no assertion at all that are not race-detector concurrency tests.
- Permanently skipped tests (`t.Skip` without an env-var guard). Fix or delete; never keep.
- Two test helper packages doing the same job.

**Complexity**
- Functions with cyclomatic complexity > 25 (`audit/gocyclo-over-25.txt`) or files > 800 lines. Split by responsibility, not by line count.

---

# Part B: findings

## 1. Critical and major bugs (fix first, behaviour PRs)

| # | Sev | St | Where | Problem | Action |
|---|---|---|---|---|---|
| B1 | CRITICAL | V | `forge/api/errors/idempotency_stores.go:187-231` | `DatabaseStore` builds SQL and discards it; `Set` returns `nil`, `Get` always "not found". Idempotent endpoints backed by it silently re-execute side effects (e.g. payments) | Make `NewDatabaseStore` return `ErrNotImplemented` now (flag); implement with real queries later |
| B2 | CRITICAL | V | `forge/orm/queryset.go` `Union` / `Intersection` / `Difference` | Ignore `other` and return a clone of the receiver: wrong results, no error, no tests | Return an error queryset (`qs.err = ErrNotImplemented`) now; implement UNION/INTERSECT/EXCEPT later |
| B3 | MAJOR | V | `forge/admin/core/admin.go:595-605` → `orm.Manager.UpdateFields` | Admin PATCH copies every request key into the update. `UpdateFields` only checks that the column exists. `Config.ReadOnlyFields` / `GetReadOnlyFields` (`admin/core/config.go:55,64`) are never enforced on writes, so read-only, auto-managed and excluded fields (`created_at`, `password`, `is_superuser` on a users admin) can be overwritten via the API | Filter create/update payloads through an allow-list: form fields minus read-only (config + metadata `read_only`); reject unknown keys with 400 |
| B4 | MAJOR | V | `forge/identity/backends/password.go:55-85` | Unknown user returns before any bcrypt compare (timing enumeration), and `ErrUserInactive` / `ErrUserLocked` are returned **before** the password is checked, so account status is disclosed without the password | Always run bcrypt (dummy hash when user is missing); check the password first, then status; return one generic error to the client |
| B5 | MAJOR | V | `forge/orm/queryset.go:1425`, `orm/update_builder.go:174` | `Update` / `Increment` / `Decrement` emit the caller's key as the column. `GetField` matches Go names and case-insensitively, so `"Price"` passes validation and then produces `SET "Price" = ...`, which fails on Postgres | Resolve every key to `FieldInfo.DBColumn` in one helper used by all write paths |
| B6 | MAJOR | V | `forge/orm/manager.go:406-441` | `Create` sets the generated id back only on an int64 field named `ID`/`Id`/`id`, ignoring `schema.PrimaryKey`; the result is discarded. `Manager.Update` treats id `0` as missing, so non-int primary keys are unsupported | Use schema PK metadata; return an error when the id cannot be set; flag non-integer PKs as unsupported |
| B7 | MAJOR | V | `forge/orm/prefetch.go:179-180` | Many-to-many prefetch guesses through-table columns with `TrimSuffix(table, "s") + "_id"` (`categories` → `categorie_id`) | Use the relation's `Through` metadata and explicit columns (see D1) |
| B8 | MAJOR | V | `forge/filter/persistence.go:34-50` | `InMemoryFilterStorage` map has no mutex; concurrent use from handlers is a data race | Add a `sync.RWMutex`, or delete it with the rest of the unused filter persistence (see §4) |
| B9 | MAJOR | V | `forge/log/logger.go:136-142` | Remote log output returns a no-op writer: configured audit-log shipping silently drops everything | Fail configuration validation with "remote output not implemented" (flag) |
| B10 | MAJOR | R | `forge/db/db.go:238-247` | `rebindPostgresToSQLite` regex-rewrites `$N` and `ILIKE` across the whole statement, including string literals (`DEFAULT '$1 off'`) | Short term: literal-aware rewrite. Target: dialect-owned placeholders (D3) |
| B11 | MAJOR | R | `forge/orm/update_builder.go:192-202` | `RawExpression.ToSQL` rebinds by `strings.Replace` on its own text | Build raw expressions through `SQLBuilder.AddArg` at construction |
| B12 | MINOR | V | `forge/admin/api/rest/router.go:81-91` | CORS echoes any origin with credentials. Admin auth is `Authorization: Bearer` only (no cookies), so this is not a session-theft hole today, but it becomes one the moment cookie auth is added | Allow-list origins from config; default same-origin |
| B13 | MINOR | V | `forge/cli/commands/admin/createsuperuser.go:120-126` | Superuser password is read with terminal echo | `golang.org/x/term.ReadPassword` |
| B14 | MINOR | R | `forge/filter/dialect.go:34,58`, `forge/filter/widgets/autosuggest.go:28-38` | JSON path interpolated into SQL; widget HTML built without escaping (XSS). Both are dead code today | Delete with the dead packages (§4) |
| B15 | MINOR | V | `forge/orm/schema_registry.go` `GetModelSchemaByType` | Only served cached schemas, so reverse-relation lookups depended on test order | Fixed in #215 |
| B16 | MAJOR | V | `forge/api/parsers/multipart.go:24` | `multipart.NewReader(r, "")`: an empty boundary makes every multipart/file-upload parse fail | Take the boundary from `Content-Type` (`mime.ParseMediaType`), or use `r.ParseMultipartForm` |
| B17 | MAJOR | V | `forge/api/viewset.go` `setFieldValue` (default branch, ~860) | `reflect.ValueOf(nil).Type()` panics when JSON sends `null` for a float, pointer, slice or map field | Guard `!valueValue.IsValid()`; long term decode through `encoding/json` / `mapstructure` (already a dependency) |
| B18 | MAJOR | V | `forge/db/migrate/generate/squash.go:37-80` | Squash writes a new, higher-version migration with the combined SQL but keeps the originals ("should be archived" comment), so `migrate up` runs every statement twice | Flag `forge migrations squash` as not implemented until it replaces the range atomically |
| B19 | MAJOR | V | `forge/api/permissions/is_owner.go:22-27,95-99` | Defaults `"user_id"` / `"id"` are looked up with `reflect.FieldByName`, which only matches Go field names (`UserID`, `ID`), so owner checks always deny | Resolve through schema metadata (Go name or column) |
| B20 | MAJOR | V | `forge/api/permissions/is_admin.go:44-62` | `IsAdminUser` looks only for `IsAdmin`; `identity/models.User` has `IsStaff` / `IsSuperuser`, so real admins are always denied | Check `IsStaff`/`IsSuperuser` (or a small `AdminUser` interface) |
| B21 | MAJOR | V | `forge/api/filters/ordering.go:70-73` | Calls variadic `OrderBy(...any)` via `reflect.Call` with one `[]string` argument, so the ORM receives a slice as a single field and builds invalid `ORDER BY` | Pass each field as its own argument, or call the typed interface instead of reflection |
| B22 | MAJOR | V | `forge/api/filters/search.go:38-48` | Calls `Search` via reflection; `orm.QuerySet` has no `Search`, so `?search=` is silently ignored | Build `orm.Or(F(field).IContains(q)...)` like the admin does; share one helper |
| B23 | MAJOR | V | `forge/api/pagination.go:99` | `baseURL` from `r.URL.Scheme` + `r.URL.Host`, which are empty on server requests, so `next`/`previous` become `://path?page=2` | Build from `r.Host` + TLS/`X-Forwarded-Proto` (trusted proxy aware), or return relative links |
| B24 | MAJOR | V | `forge/admin/api/rest/router.go:360-372` | Admin login accepts only the `FORGE_ADMIN_USERNAME`/`FORGE_ADMIN_PASSWORD` env pair; users created by `forge createsuperuser` cannot log in to the admin | Authenticate through `identity` (staff/superuser), keep env credentials as an explicit bootstrap option |
| B25 | MAJOR | V | `forge/api/authentication/token.go:51-53` | A request that presents an unknown token falls through to anonymous (`nil, nil`) instead of 401; DRF raises `AuthenticationFailed` here | Return an invalid-token error when credentials were supplied but did not match |
| B26 | MAJOR | V | `forge/admin/history_manager.go:16-21` | `getMem()` lazily initialises `m.mem` without synchronisation: data race on concurrent requests, lost history entries | Initialise in the constructor; history is in-memory only, flag persistence as a gap |
| B27 | MAJOR | V | `forge/server/security.go:260-307` `SanitizeHTML` (+ SQL keyword blacklist ~120-181) | Regex "sanitizer" misses unquoted handlers (`<img onerror=alert(1) src=x>`); the SQL blacklist is not a defence. No callers today | Delete both; if HTML sanitising is needed, `bluemonday` |
| B28 | MINOR | R | `forge/api/content_negotiation.go:30,58` | `Renderers[0]` / `Parsers[0]` indexed without a length check | Validate configuration in the constructor |
| B29 | MINOR | R | `forge/identity/middleware/auth.go:104-121` | Cookie session auth does not itself require CSRF; safe only if `server` CSRF middleware is mounted on the same routes | Document the requirement or enforce CSRF when auth came from a cookie |
| B30 | MINOR | R | `forge/admin/api/rest/login_limiter.go:66-83` | Cleanup only runs above 10k entries and holds the lock while iterating; username rotation can grow memory | Fold into the shared limiter (§5) with TTL eviction |
| B31 | MAJOR | V | `forge/db/migrate/sql/sqlite.go` `BuildAddForeignKey` → `base.go` | SQLite delegates to the base builder, which emits a PostgreSQL `DO $$ ... pg_constraint ... ::regclass` block, so any SQLite migration that adds a foreign key fails | SQLite: define FKs inside `CREATE TABLE`, or use the table-rebuild recipe; never share PG-only DDL through the base |
| B32 | MAJOR | V | `forge/db/migrate/sql/postgres.go:137`, `sqlite.go:139` (down SQL for `DropColumn`/`DropTable`) | Down SQL returns an error when the original definition is unknown, so generating any migration that drops a column or table fails outright | Carry the previous column/table definition from state into the change, or emit an explicit irreversible marker |
| B33 | MAJOR | V | `forge/db/migrate/sql/sqlite.go:106,138` and `DropForeignKey` case | SQLite down steps for `AddColumn`, `DropConstraint`, `DropForeignKey` are SQL comments; rollback "succeeds" and records the version as reverted while the schema is unchanged | Generate the SQLite table-rebuild sequence, or refuse to generate a no-op down |
| B34 | MAJOR | V | `forge/db/migrate/execute/status.go` `mergeAppliedVersions` | Marks every integer from 1 to the current version as applied. With timestamp versions (`20240101120000`) this loops ~2×10¹³ times and exhausts memory; with gaps it reports versions that never existed | Merge only versions that exist as files or rows |
| B35 | MAJOR | V | `forge/orm/queryset.go` `Last` / `Reverse` | `Reverse` only flips existing `orderBy`; on an unordered queryset `Last()` returns the first row | Default to primary key when no ordering (Django behaviour) |
| B36 | MAJOR | V | `forge/config/settings.go:109` | `GetInt("database.conn_max_lifetime")` on the default `"5m"` returns 0, disabling the connection lifetime; `config.go:171` reads the same key correctly with `GetDuration`, i.e. two settings loaders disagree. `settings_test.go:124-128` asserts the 0 | One loader using `GetDuration`; fix the test |
| B37 | MAJOR | V | `forge/orm/update_builder.go:53-56` | `reflect.TypeOf(nil)` is nil, so `UpdateBuilder.Set(field, nil)` panics on `.AssignableTo` when clearing a nullable column | Handle `nil` explicitly (allowed only for nullable fields) |
| B38 | MAJOR | V | `forge/db/migrate/generate/generator.go:349-367` vs `execute/recover.go:50,79` | The first generated migration creates `schema_migrations(id, name, checksum, applied_at)` while golang-migrate (used to apply) owns `schema_migrations(version, dirty)`; `CREATE TABLE IF NOT EXISTS` silently no-ops, so the checksum column the verifier expects never exists | Let golang-migrate own its table; store checksums in a separate `forge_migration_checksums` table or drop the feature |
| B39 | MINOR | V | `forge/orm/manager_helpers.go:192-198,385-399`, `orm/query_expr.go:183-272` | INSERT/UPDATE/DELETE helpers and `QueryExpr` interpolate table and column names unquoted, so reserved names (`order`, `user`) break; INSERT hardcodes `RETURNING id`. Not an injection path: names come from schema/codegen, and request filters go through `orm.F` → `FieldRef`, which is resolved against the schema and escaped | Quote through the dialect; use the schema PK in `RETURNING` |
| B40 | MINOR | R | `forge/orm/queryset.go` `Exists`, `Update` | `Exists` runs a full `COUNT(*)` instead of `SELECT 1 ... LIMIT 1`; `Update` iterates a map so SET order and placeholder order change between calls | `LIMIT 1`; sort keys |

## 2. Security review items to re-check during the refactor

- Admin writes: B3 (field allow-list) and create path, same check.
- Auth: B4; confirm `api/authentication/*` session and API-key paths have the same generic-error behaviour.
- `admin/core/notifications.go` sets `Access-Control-Allow-Origin: *` on an SSE stream. It is unrouted now (§3); if it is ever mounted it must sit behind admin auth with an origin allow-list.

## 3. Incomplete features to flag (do not implement now)

Each becomes a loud failure plus a row in `docs/design/feature-gaps.md`.

| Where | What is missing | Flag as |
|---|---|---|
| `orm/queryset.go` `Union` / `Intersection` / `Difference` | set operations | error queryset (B2) |
| `api/errors/idempotency_stores.go` `DatabaseStore` | DB-backed idempotency | constructor error (B1) |
| `api/viewset.go:756` `getManagerFromModel` | always `reflect.Value{}` | delete with the old viewset (§5) or error |
| `log/logger.go:136` remote exporter | remote log shipping | config error (B9) |
| `registry/plugin.go:233-241` `applyAdminExtensions` / `applyAPIExtensions` | admin/API plugin extensions | registration returns `ErrNotImplemented` |
| `admin/core/notifications.go` `NotificationHub` / `SSEHandler` | real-time notifications (planned, see backend-tasks.md) | leave unrouted, add gap row |
| `filter/filters/model_choice.go:88-92` `GetOptions` | choice options | deleted with `filter/filters` (§4) |
| `filter/filters/lookup.go:29-33` | `lookup:value` parsing | deleted with `filter/filters` |
| `orm/queryset.go:625` `buildJoinClause` | `SelectRelated` supports one level only (multi-hop exists for filters) | error on multi-hop `SelectRelated` |
| `cli/commands/development/test.go:48` | `forge test` only prints a hint | remove the command or document it |
| Non-integer primary keys (`Manager.Update`, `Create` id set-back) | UUID/string PKs | explicit error at schema build (B6) |
| MySQL | `db/dialect/dialect.go` docs and `filter/dialect.go` `MySQLAdapter` imply support; no driver exists | remove MySQL mentions; state "Postgres and SQLite only" in docs |
| `orm/queryset.go:437` `Aggregate` | aggregates are appended to a slice that is never rendered into SQL | error until implemented |
| Transactions in the ORM | `Manager` / `QuerySet` hold `*db.DB`; `getDB` never looks at a transaction, so ORM calls cannot join a `db.Tx` | document as unsupported now; D12 fixes it |
| `forge migrations squash` | see B18 | command returns "not implemented" |
| SQLite rollback of `AddColumn` / `DropConstraint` / `DropForeignKey` | down SQL is a comment, see B33 | generator refuses to emit a no-op down step |

## 4. Dead code to remove

Packages with **no importers** anywhere in `forge/`, `tests/`, `examples/` (V, importer grep) and, where listed, 100% unreachable in `deadcode -test`:

| Package | Go LOC (non-test) | Notes |
|---|---|---|
| `filter/filters` | 1611 | 107/107 funcs unreachable; also carries B14 and the `GetOptions`/`Parse` stubs |
| `api/serializers` | 445 | 25/25 unreachable. Keep `api/serializers/fields` only if something uses it |
| `admin/utils` | 434 | 0 importers |
| `cli/internal` | 224 | 14/14 unreachable |
| `db/migrate/dependencies` | 206 | 6/6 unreachable |
| `log/hooks` | 202 | 0 importers |
| `api/caching` | 184 | one of four caches (§5) |
| `api/versioning` | 151 | 0 importers |
| `admin/codegen` | 125 | 0 importers |
| `filter/widgets` | 86 | 6/6 unreachable; XSS (B14) |

Also unreachable, but only used inside their own package or by the identity router; decide per item:

- `identity/handlers` and `identity/serializers` (100% unreachable from tests; only `identity/router.go` imports handlers, and nothing mounts that router). Either wire identity routes into the example and test them, or delete the router, handlers and serializers together.
- `db/migrate/execute/executor.go` `Executor` (only its own integration test uses it) — duplicate of `db.MigrationRunner` (§5).
- `filter/dialect.go`, `filter/optimizer.go` (`QueryOptimizer` built but `Optimize()` never called), `filter/persistence.go` (if nothing saves filters).
- `db/migrate/migrate_impl.go` facade (`NewDetector` and friends only re-export `generate`/`sql`/`state`).
- `identity/password.go` re-export shim of `identity/utils/password.go`.
- Files whose exported types are referenced only in their own file (V, grep): `filter/cache.go` (`FilterCache`, deferred feature, see backend-tasks.md), `filter/metrics.go` (`Metrics`, `AlertChecker`), `filter/persistence.go` (`SavedFilter`), `filter/relations.go` (placeholder `JOIN %s ON ...` strings), `filter/typed_filter.go` (its `Build` keeps only the last filter), `schema/registry.go` (`RegisterFieldType`), `schema/field_traits.go` (`FieldTrait`), `orm/safe_accessor.go` (`SafeAccessor`, used only by its test). Not dead: `schema/helpers.go` (`IndexOn` has 10 users).
- `orm/update_builder_helpers.go` `NewUpdateBuilderFromQuerySet` duplicates `NewUpdateBuilder` (R).
- High dead ratio, prune function by function from `audit/deadcode-forge-tests.txt`: `schema` 78/114, `registry` 35/52, `validate` 46/70, `filter` 110/173, `api/errors` 60/119, `api/authentication` 13/23, `errors` 10/19.

Rule: before deleting an exported function that `deadcode` lists, grep `forge/cli` templates and `docs-site/` for it; generated user code may call it.

## 5. Duplicates to merge

| Concept | Copies | Keep | Merge plan |
|---|---|---|---|
| Viewset base | `api/viewset.go` (864), `viewset_enhanced.go` (787), `viewset_enhanced_integrated.go` (108), `viewset_config.go` (80) | the enhanced stack, collapsed into one type | Point CLI scaffolds (`cli/commands/project/add_api.go:126`, `auth.go:156`) at it, delete `BaseViewSet`, fold "integrated" into the main type |
| Rate limiting | `server/ratelimit.go` (x/time/rate), `api/throttling/*` (fixed window), `admin/api/rest/login_limiter.go` | one limiter store on `golang.org/x/time/rate` | Throttle classes and login limiter become policies over the shared store |
| In-memory cache | `api/caching/memory.go`, `api/viewset_cache.go`, `api/throttling/cache.go`, `filter/cache.go` (deferred feature) | one internal TTL cache (or `otter`/`ristretto`) | Delete `api/caching` (dead); others use the shared cache |
| API errors | `api/errors` (RFC 7807, sanitizer), `api/exceptions`, `forge/errors` (`server/errors.go` is only an adapter over `api/errors`) | `api/errors` (most complete, has the sanitizer) | `api/exceptions` becomes constructors for `api/errors` problems, then is removed |
| Authentication | `api/authentication/*` vs `identity/backends/*` (+ `identity/middleware`, `identity/service/auth.go`) | `identity` owns users and credentials; `api/authentication` adapts it to HTTP | One `Authenticator` interface; delete the parallel abstraction |
| Migration runner | `db/migrations.go` `MigrationRunner` vs `db/migrate/execute/executor.go` `Executor` (80-line blocks copied, same Windows URL fallback chain) | `MigrationRunner` (the one callers use) | Delete `Executor`; extract `buildMigrationsURL` helper |
| Migration checksum | `db/migrate/execute/checksum.go` vs `db/migrate/verify/checksum.go` (91 lines, 1 line differs) | one in `db/migrate/core` | Move and delete the copy |
| FK / column resolution | `fkColumnFor`, `reverseFKFor`, `prefetch.go:43-54`, `prefetch.go:179` | one schema-driven resolver (D1) | All join/prefetch paths call it |
| Dialects | `db/dialect/postgres.go` vs `sqlite.go` (~95 lines similar), `db/migrate/sql/postgres.go` vs `sqlite.go` | shared base + per-dialect overrides | Embed a `baseDialect` |
| Test helpers | `tests/helpers` vs `tests/testhelpers` (`postgres_features.go`, `sql_assertions.go` differ by 1 line), `forge/internal/testutils` vs `forge/identity/testutils` | `tests/testhelpers` and `forge/internal/testutils` | Move callers, delete the copies |
| Repository boilerplate | `identity/repository/{user,permission,session,token}.go` (8 clone groups) | generic scan/query helpers | Extract `scanOne`/`queryList` helpers, or move repositories onto the ORM |
| Password helpers | `identity/password.go` shim vs `identity/utils/password.go` | `identity/utils` | Delete the shim |
| Admin CRUD helpers | `admin/core/admin.go:884-948` (three copies) | one helper | Extract |

## 6. Design problems and simplifications

| # | Problem | Evidence | Direction |
|---|---|---|---|
| D1 | **Relations are resolved by guessing names** at query time | `fkColumnFor` tries five naming patterns, reverse lookup, prefetch, M2M `TrimSuffix` | Store FK column, target PK and through-table columns on `RelationInfo` when the schema is built (codegen already parses relations). Delete every guess. Biggest single correctness win |
| D2 | **Global mutable state** | `admin.DefaultSite`, `registry.global*` (4), `admin/core` registries (2), `cli/core` registry, `api.globalCache`, `api/exceptions.globalHandler`, schema cache | An explicit `forge.App` holding site, registries, DB and config; `Default*` helpers stay as thin optional wrappers |
| D3 | **SQL placeholders rewritten after the fact** | builder emits `$N`, `db.RebindPlaceholders` rewrites for SQLite; B10, B11 | Dialect object passed to `SQLBuilder`; placeholders generated once |
| D4 | **Config re-read deep in the stack** | `config.NewConfig()` in `db/migrate/generate/generator.go:35,51`, `execute/executor.go:78`, `cli/core/context.go:21` | Load config once in `main`/CLI root and pass the driver/values down |
| D5 | **Generated managers swallow errors** | `var XObjects, _ = orm.NewManager[...]` in generated `gen.go` | Generate a `MustNewManager` that panics at init, or a constructor used by `App` |
| D6 | **Reflection where codegen exists** | ~150 `reflect.` uses in `orm`; viewset dispatch via `reflect.Value.Call` although `ManagerInterface` is declared | Generate scan/field-access code; dispatch viewsets through the declared interfaces |
| D7 | **God files** | `admin/api/rest/router.go` 1944, `orm/queryset.go` 1831, `codegen/ast_parser.go` 1501, `admin/core/admin.go` 1199 | Split by responsibility: router → auth / crud / meta handlers; queryset → build / joins / scan / write; ast_parser → fields / relations / meta |
| D8 | **Very complex functions** | `ASTParser.extractOptionFromMethod` (74), `server.StaticFS` (43), `Admin.ListObjects` (43), `validate.getErrorMessage` (42), `Router.attachDisplayLabels` (42) | Table-driven option parsing; lookup→expression map in `ListObjects` (the admin already duplicates `filter`'s lookup parsing) |
| D9 | **Two filtering stacks** | `admin/core/admin.go` `ListObjects` hand-parses `field__lookup`, while `forge/filter` has its own parser and converter | Admin list uses `forge/filter` |
| D10 | **Migration generation is a home-grown schema differ plus a DDL lexer** (`db/migrate/{generate,parse,state}`) on top of golang-migrate | ~3k LOC | Evaluate Atlas for diffing (see §8); keep golang-migrate for applying files |
| D11 | **Dialect interface promises more than exists** | 16-method `Dialect` with MySQL docs, two implementations | Shrink to what Postgres and SQLite need |
| D12 | **ORM cannot run inside a transaction** | `Manager.db *db.DB`, `BaseQuerySet.getDB` returns the pool only | A small `DBTX` interface (`ExecContext`, `QueryContext`, `QueryRowContext`) accepted by managers/querysets, plus `db.WithTx` passing a tx-bound manager |
| D13 | **Two expression trees** | `orm/query_expr.go` `QueryExpr` (typed `FieldExpr` codegen API, unquoted, PG-only `EXTRACT`) vs `orm/expression.go` `Expression`/`FieldRef` (schema-resolved, escaped, dialect-aware) | Make typed field helpers return `Expression`; delete `QueryExpr` |
| D14 | **`forge/filter` is mostly unconnected machinery** | cache, metrics, persistence, optimizer, relations, typed filter, dialect adapters: unreferenced; the admin parses lookups itself | Keep a thin "query params → `orm.Expression`" parser used by admin and API; delete the rest |

## 7. Tests to delete or fix

| Where | Problem | Action |
|---|---|---|
| `db/dialect/dialect_test.go` `TestPostgreSQLDialect_ImplementsDialect`, `TestSQLiteDialect_ImplementsDialect`; `db/pool_test.go` `TestOptionType`; `admin/site_test.go` `TestTypeAliases` | assert nothing beyond what the compiler checks | replace with `var _ Dialect = (*PostgresDialect)(nil)`; delete tests |
| `orm/expression_test.go` `TestCombinedExpression_WithValues`, `api/serializer_test.go` `TestBaseSerializer_Validate_Invalid`, `api/parsers/parser_test.go` `TestXMLParser_Parse`, `log/logger_test.go` `TestLoggerTrace` | no assertions | add the assertion the name promises or delete |
| 29 `t.Skip` calls: `orm` 14, `db` 13, `identity/service` 2 | permanently skipped (e.g. `orm/schema_test.go:32,65,79,89` "schema not registered", `orm/date_parts_test.go:224` "no sqlite helper", `identity/service/password_test.go:255` expiry) | fix with existing helpers (sqlite helpers exist in `db`, schema registration now works after #215) or delete; env-guarded Postgres skips may stay |
| Tests for dead packages (`api/caching`, `filter/filters`, `admin/utils`...) | test code that only keeps dead code alive | delete with the package |
| `config/settings_test.go:124-128` | asserts `ConnMaxLifetime == 0` and `ConnMaxIdleTime == 0`, i.e. it locks in bug B36 | fix the loader, then assert 5m |
| `orm/update_builder_test.go` (7 tests), `orm/schema_test.go:32,65,79,89` | skipped because the test model's schema is not registered | register the schema (works after #215) and un-skip |
| `orm/schema_test.go` `TestGetModelSchema`, `TestNewFieldAccessor` | assertions only run `if err == nil`, so the test passes when setup fails | `require.NoError` first (R) |
| `db/pool_test.go` | 12 tests skip without CGO because of `mattn/go-sqlite3` | consider `modernc.org/sqlite` for tests (§8) |
| Concurrency tests without assertions (`TestRegistry_Concurrency`, `TestSettings_Concurrency`, `TestMemoryCache_Concurrency` x2) | fine: they exist for `-race` | keep; make sure CI runs `-race` |
| **Missing**: `schema` package has 0 tests; `validate` has 71 test lines for 1271; set operations, admin write allow-list, M2M prefetch | critical untested code | add with the fixes above |

## 8. Libraries instead of hand-rolled code

Already in `go.mod`: chi, cors, scs (sessions), gorilla/csrf, golang-jwt, golang-migrate, validator/v10, viper, zap, cobra, survey, x/crypto, x/time, lib/pq (tests only), go-sqlite3, testify, uuid, strcase.

| Hand-rolled | Use instead | Removes | Risk |
|---|---|---|---|
| `api/throttling/*`, `admin/api/rest/login_limiter.go` | `golang.org/x/time/rate` (already used by `server/ratelimit.go`) or `github.com/go-chi/httprate` | ~500 LOC and two algorithms | low |
| `api/caching`, `api/viewset_cache.go`, `api/throttling/cache.go` | `github.com/maypok86/otter` (or one small internal TTL map) | three caches | low |
| `forge/log` on zap + `log/hooks`, `log/exporters` | stdlib `log/slog` (handlers for JSON/text; OTel bridge when remote export is needed) | zap dependency, hooks package, no-op exporter | medium (log API change across packages) |
| `db.RebindPlaceholders` regex | dialect-generated placeholders (D3); if a rewrite must stay, `sqlx.Rebind` | fragile regex | low |
| `orm/sql_builder.go` string assembly | keep the builder, but model it on `github.com/stephenafamo/bob` or `github.com/doug-martin/goqu/v9` dialects; a full swap is optional | placeholder/quoting code | high if swapped wholesale; do D3 first |
| `db/migrate/generate` + `parse` (schema diff + DDL lexer) | `ariga.io/atlas` schema diffing, golang-migrate stays for apply | ~2-3k LOC | high; spike first |
| `cli/.../createsuperuser.go` password input | `golang.org/x/term` | echo bug (B13) | low |
| `config` typed getters over viper, re-instantiated per call | load once; optionally `github.com/knadh/koanf/v2` | viper boilerplate | medium |
| `validate/field_validator.go`, `schema_validator.go`, `typed_validator.go` | express rules as `validator/v10` tags / custom validators (already a dependency) | second validation model | medium |
| Postgres test driver `lib/pq` (maintenance mode) | `github.com/jackc/pgx/v5/stdlib` | deprecated driver | low |
| SQLite driver `mattn/go-sqlite3` (CGO) | `modernc.org/sqlite` (pure Go) | CGO requirement; 12 `db/pool_test.go` skips; easier cross-compiles | low-medium (driver name and some pragmas differ) |
| CSRF middleware | keep `gorilla/csrf` or move to Go's `http.CrossOriginProtection` (Go 1.25+) | one dependency | low |
| `identity/repository/*` hand-written SQL | the framework's own ORM (dogfooding) or `sqlc` | ~8 clone groups | medium |

## 9. Second-opinion validation

### agy (Gemini), web/admin/auth

Asked to validate 14 claims and report what was missed.

- **Confirmed 13, partly 1.** Correction taken: `server/errors.go` is an adapter over `api/errors`, not a fourth error system (§5 updated).
- **Severity disagreements:** agy rated `getManagerFromModel` CRITICAL (it breaks `GetManager()` when a queryset has no `Create`), and the unimported packages MAJOR. Kept here as a STUB flag and dead-code removal respectively, since neither returns wrong data.
- **New findings accepted after Claude read the code:** B16-B27 (multipart boundary, `null` panic, squash double-apply, owner/admin permission lookups, ordering/search filters, pagination links, admin login env-only, unknown token falls through, history race, regex sanitizer).
- **New findings kept as R (not re-checked):** B28-B30.
- **Library suggestions accepted:** `bluemonday` (only if sanitising is kept), `golang.org/x/term`, `mapstructure` for payload decoding, `x/time/rate` consolidation. Rejected: `elnormous/contenttype` (the negotiator is small; fix it in place), replacing admin history with a DB table (that is a feature, tracked as a gap).

### Nemotron 3 Ultra, data layer

Did not finish: the run stalled for over 40 minutes while still reading files and produced no findings, so it was stopped. muse-spark 1.3 (free tier hangs; paid tier needs a payment method) and codex (usage limit) were also unavailable, so the same prompt went to agy.

### agy (Gemini), data layer

Prompt: `tasks/audit-second-opinion-data.md`. 20 bugs, 15 dead/dup/stub, 8 design, 8 test and 5 library items reported.

- **Accepted after Claude read the code:** B31-B38 (SQLite FK DDL is PostgreSQL, drop-column migrations cannot be generated, SQLite no-op rollbacks, applied-version loop, `Last()` on unordered sets, duration parsed as int, `Set(nil)` panic, conflicting `schema_migrations` tables), the aggregate stub, missing ORM transactions (D12), the second expression tree (D13), the filter package (D14), the dead-file list in §4, and the test items in §7.
- **Downgraded:** "SQL injection through unquoted identifiers" (CRITICAL → B39 MINOR). Names come from schema/codegen, and request-driven filters go through `orm.F`, which resolves against the schema and escapes. Reserved-word breakage is real.
- **Rejected:** "unfiltered `Update`/`Delete` wipe tables" as a bug. `qs.Delete()` on an unfiltered queryset is intended (Django allows `Model.objects.all().delete()`); an opt-in guard is a possible API choice, not a defect. "`schema/helpers.go` is dead": `IndexOn` has 10 users. "Replace the migration engine with golang-migrate": golang-migrate cannot generate migrations from models, which is the feature; see D10/Atlas instead. `gorilla/schema` for filter params: not a fit for `field__lookup` keys.
- **Kept as R (not re-checked):** annotations discarded on scan, `path_cache` keyed without the model, `NewUpdateBuilderFromQuerySet` duplicate, vacuous `schema_test.go` assertions, B40.

## 10. Plan

Sized for one delegate task and one PR each. Waves can run 2-3 PRs in parallel if they touch disjoint packages.

**Wave 0: stop the bleeding (behaviour fixes, small)**
1. B1 + B2 + B9 + §3 flags: every stub fails loudly, `docs/design/feature-gaps.md` created.
2. B3: admin create/update field allow-list (+ tests for read-only and unknown keys).
3. B4: auth backend constant-time and status-after-password (+ tests).
4. B5 + B6: write paths resolve DB columns; PK set-back uses schema metadata.
5. API layer correctness: B16, B17, B21, B22, B23 (each small, one PR for parsers/decoding, one for filters/pagination).
6. Permissions and auth: B19, B20, B24, B25.
7. B18 (flag squash), B26 (history race), B27 (delete regex sanitizer and SQL blacklist).
7a. Migrations: B31-B34 and B38 (SQLite DDL, down generation, applied-version merge, bookkeeping table). One PR per dialect builder, one for status/bookkeeping.
7b. ORM/config: B35, B36 (+ its test), B37.

**Wave 1: delete dead code (large diffs, no behaviour change)**
8. Delete the ten unimported packages in §4 and their tests.
9. Delete `Executor`, `migrate_impl.go` facade, `identity/password.go` shim, `filter/{dialect,optimizer,widgets}` leftovers, dead test helpers.
10. Prune dead functions in `schema`, `registry`, `validate`, `api/errors` from the deadcode list.
11. Tests: delete or fix everything in §7.

**Wave 2: merge duplicates**
12. One migration runner + one checksum.
13. One rate limiter + one cache.
14. One API error system.
15. One viewset.
16. One test-helper package per module; repository helpers.

**Wave 3: design**
17. D1 explicit relation metadata (codegen emits it; joins/prefetch/M2M use it; B7 goes away).
18. D3 dialect-owned placeholders (B10, B11 go away).
19. D4/D5 config loaded once, generated managers fail at init.
20. D7/D8 split god files and complex functions (pure moves, separate PRs).
21. D2 `forge.App` instead of globals.

**Wave 4: libraries**
22. `log/slog` migration.
23. Atlas spike for migration diffing (decision doc before code).

---

## Appendix: reproduce the tool runs

```bash
# deadcode needs a Go 1.26 toolchain; system go is 1.25
GOTOOLCHAIN=go1.26.6 GOPROXY=file://$HOME/go/pkg/mod/cache/download GOSUMDB=off \
  go install golang.org/x/tools/cmd/deadcode@v0.50.0
cd forge && deadcode -test ./...                                   # unreachable even from tests
cd examples/ecommerce && deadcode -filter github.com/forgego/forge ./...
cd forge && find . -name '*.go' ! -name '*_test.go' ! -path './admin/ui/*' | dupl -files -t 80
cd forge && gocyclo -over 25 -ignore '_test|admin/ui' .
```
