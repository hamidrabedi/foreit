Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

ROLE: You are the delegated coder: edit the files yourself; do not call other agents. Use US spelling in comments.

# Split forge/orm/queryset.go (2007 lines) by responsibility — pure move, no behaviour change

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit non-test files in forge/orm.

## Rules (mandatory)
- Pure move: cut declarations out of queryset.go into new files in package orm. Do not rename anything, change any signature, body, comment text, or the order of statements inside a function. The only allowed edits are import lists.
- Work group by group: after moving each group, run `go build ./orm/` and fix imports before the next group.
- File names: lowercase, underscore-separated, after the responsibility; no `_helpers`/`_utils`/`_impl`.

## Target layout (move each declaration to the file listed; anything not listed stays in queryset.go)
- queryset.go: `QuerySet` interface, `OrderField` type and `GetFieldPath`, `IsAscending`, `Asc`, `Desc`, `NewOrderField`; `BaseQuerySet` struct, `NewQuerySet`, `SetDB`, `getDB`, `getDialect`, `newSQLBuilder`, `clone`; chainable builders `Filter`, `Exclude`, `OrderBy`, `extractOrderFieldPath`, `extractOrderFieldAscending`, `Reverse`, `Limit`, `Offset`, `Distinct`, `Select`, `Only`, `Defer`, `SelectRelated`, `PrefetchRelated`, `Aggregate`, `Annotate`, `Values`, `ValuesList`, `UpdateBuilder`; `Union`, `Intersection`, `Difference`.
- queryset_sql.go: `buildSQL`, `mergeJoins`, `pkSubquery`, `buildSelectClause`, `buildWhereClause`, `buildOrderByClause`, `buildLimitClause`, `buildOffsetClause`, `buildCountOrExistsSQL`, `appendWhereAndJoinParts`, `buildCountSQL`, `BuildExistsSQL`, `buildExistsSQL`.
- queryset_joins.go: `fkColumnFor`, `buildJoinClause`, `warnSkippedSelectRelated`, `appendJoinClause`, `recordJoinAliases`, `reverseFKFor`, `createJoinResolver`, and any package-level variable only they use (e.g. the warned-relations cache).
- queryset_scan.go: `scanField` type, `scanRows`, `buildFieldMap`, `prepareScanArgs`, `setFieldValue`.
- queryset_execute.go: `All`, `Get`, `First`, `Last`, `withDefaultPKOrder`, `Count`, `Exists`.
- queryset_write.go: `columnFor`, `Update`, `BulkUpdate`, `Delete`.
- queryset_values.go: `ValuesQuerySet` and `ValuesListQuerySet` interfaces, `BaseValuesQuerySet`, `BaseValuesListQuerySet` and all their methods.
- Every resulting non-test file must be under 600 lines. `prepareScanArgs` (152 lines) and `createJoinResolver` (104 lines) move unchanged; do not refactor them in this task.

## Acceptance (from the module root)
    gofmt -l orm                                      # empty
    go build ./... && (cd ../examples/ecommerce && go build ./...)
    go vet ./orm/...
    staticcheck ./orm/...                             # no output
    go test -race ./orm/... ./admin/... ./api/... ./filter/... -count=1
    ~/go/bin/golangci-lint run --timeout=5m --config=../.golangci.yml ./orm/...   # 0 issues
    wc -l orm/queryset*.go                            # every file under 600 lines

Report (max 15 lines): files created with line counts, each command result.
