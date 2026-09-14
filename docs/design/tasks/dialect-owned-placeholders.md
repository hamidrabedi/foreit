Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

# Placeholders are generated once for the dialect instead of regex-rewritten before execution

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit forge/orm (non-generated files), forge/db/db.go, forge/db/dialect, and tests.

## Principles (mandatory)
- Tests FIRST: add the failing literal test below before changing code. gofmt; functions < 40 lines; no new package-level mutable state; naming after behaviour.

## Facts (verified on master)
- `SQLBuilder.AddArg` (forge/orm/sql_builder.go:89) always emits `$N`; sql_builder.go:226 and :253 do the same; forge/orm/query_expr.go has 14 `$%d` formats, manager_helpers.go 5, update_builder.go 3.
- Before execution the SQL is rewritten: `(*db.DB).RebindPlaceholders` (forge/db/db.go:220) → `rebindPostgresToSQLite` (db.go:246) regex-replaces every `$N` and `ILIKE` in the whole statement, including inside string literals, so `WHERE note = 'costs $1'` or `'ILIKE'` in data is corrupted on SQLite (B10). Called from manager_helpers.go:415,495,524,544 and queryset.go:208.
- `RawExpression.ToSQL` (forge/orm/update_builder.go:217) re-numbers its own `$1..$n` with `strings.Replace`, which also hits `$1` inside literals and `$10` when replacing `$1` (B11).
- `BaseQuerySet.newSQLBuilder` (queryset.go ~197) already obtains a dialect with `BuildPlaceholders(n)`.

## Change
1. Give `SQLBuilder` a placeholder function from the dialect: PostgreSQL `$N`, SQLite `?` (SQLite also accepts `?NNN`; use `?` with positional args in order, or `?N` if argument reuse requires it — check whether any builder path reuses an arg index). `newSQLBuilder` passes it; a builder without a dialect keeps `$N`.
2. Every `$%d` format in orm (sql_builder.go, query_expr.go, manager_helpers.go, update_builder.go) goes through the builder's placeholder function. Where a helper builds SQL without a builder (manager_helpers `Build*SQL`), pass the placeholder function in from the manager's connection dialect.
3. `ILIKE` on SQLite: the dialect renders case-insensitive LIKE (`LIKE ... COLLATE NOCASE` or `lower(x) LIKE lower(?)`) at build time instead of the rewrite. Find where ILIKE is emitted in orm (lookup `icontains`/`istartswith`...).
4. `RawExpression`: parse its SQL once, replacing only `$N` tokens outside single-quoted string literals, mapping each `$N` to `builder.AddArg(r.Args[N-1])` (so `$10` is not hit by `$1`).
5. Remove `RebindPlaceholders` calls from orm execution paths. Keep `(*db.DB).RebindPlaceholders`/`Rebind` exported for external callers, but make `rebindPostgresToSQLite` skip single-quoted literals so direct callers are not corrupted either.
6. Tests: on SQLite, insert and filter a row whose text value is `'costs $1 ILIKE'` via manager Create + Filter(F(field).Eq(...)) + a raw expression, and read it back unchanged; `RawExpression` with 10+ args maps `$10` correctly; a PostgreSQL-dialect builder still emits `$1..$n` (string assertion).

## Acceptance (from the module root)
    gofmt -l orm db                                   # empty
    go build ./... && (cd ../examples/ecommerce && go build ./...)
    go vet ./orm/... ./db/... ./admin/... ./api/...
    staticcheck ./orm/... ./db/...                    # no output
    go test -race ./orm/... ./db/... ./admin/... ./api/... -count=1

Report (max 25 lines): failing output before, files changed, each command result.
