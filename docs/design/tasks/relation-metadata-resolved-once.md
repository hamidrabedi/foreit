Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

ROLE: You are the delegated coder: edit the files yourself; do not call any other agent. Use US spelling in comments.

# Relation columns are resolved once when the schema is built, not guessed at query time (D1)

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit forge/orm (non-generated files) and tests. Do not add options to forge/schema, and do not change or remove any exported function signature (the refactor must stay backward compatible; recent commits restored compatible signatures in forge/orm — keep them).

## Principles (mandatory)
- Tests FIRST. gofmt; functions < 40 lines; errors wrapped with %w; naming after behaviour.
- **No working behaviour may break.** Every query that returns rows today must return the same rows after this change. Where current behaviour is silently wrong or silently skipped, make it visible (a warning), not a new error.

## Facts (verified on master)
- `RelationInfo` (forge/orm/schema.go) stores Name, Type, TargetModel, TargetField, OnDelete, OnUpdate, RelatedName, FieldName, Through. It has no foreign key column.
- `fkColumnFor` (forge/orm/queryset.go) tries eight naming guesses each time a join is built (DBColumn or Name equal to the relation name, struct field `Name`, `Name+"ID"`, `name_id`, any `*ID` field with a matching prefix, lowercase guesses, `<target>_id`). `buildJoinClause` silently `continue`s when nothing matches, so `SelectRelated` quietly does nothing.
- `reverseFKFor` guesses `<table>_id` / `<type>_id` on the target.
- `prefetchForeignKey` (forge/orm/prefetch.go) guesses the struct field `rel.Name + "ID"`.
- `throughColumns` (prefetch.go) derives `<source model>_id` / `<target model>_id` (Django's default for auto-created through tables).
- `RelationInfo` values are created in `BuildModelSchema` (forge/orm/schema.go) from `schema.Relation`.

## Change
1. Add to `RelationInfo` (additive, exported fields are fine): `FKColumn string`, `FKStructField string`, `ThroughSourceColumn string`, `ThroughTargetColumn string`, each documented.
2. In `BuildModelSchema`, resolve them once:
   - FK/O2O: run the **same matching rules `fkColumnFor` uses today, in the same order**, once, and store the result (column and Go struct field). This is what keeps existing models working: any model that joins today must resolve to the same column. Move the rules into one unexported resolver function called at schema build.
   - M2M: store the columns `throughColumns` computes today.
3. Use the stored values:
   - `fkColumnFor(schema, rel)` returns `rel.FKColumn` when set; when empty (relations created without `BuildModelSchema`, e.g. hand-built schemas in tests), fall back to the resolver so behaviour is identical.
   - `reverseFKFor`: prefer the `FKColumn` of the relation on `target` that points back to `current`; keep the existing guesses as the fallback.
   - `prefetchForeignKey` uses `rel.FKStructField` when set, otherwise the existing `rel.Name + "ID"`.
   - `prefetchManyToMany` uses the stored through columns when set, otherwise `throughColumns`.
4. Make silent skips visible without breaking queries: when `buildJoinClause` cannot resolve a selected relation, keep skipping the join (current behaviour, rows still returned) but log one warning per model+relation using `log/slog` (`slog.Warn("select_related skipped: no foreign key column", "model", ..., "relation", ..., "expected_column", snake(relation)+"_id")`), guarded by a `sync.Map` of already-warned keys (a package-level cache is acceptable here; document why).

## Tests (new file relation_metadata_test.go, in-memory SQLite like select_related_test.go)
- `BuildModelSchema` fills FKColumn/FKStructField for a `ForeignKeyField("Category", "Category")` with a `category_id` field, and for a model whose Go field is `AuthorID` with DB column `author_id`.
- Equivalence: for every model used by the existing tests in forge/orm (select_related_test.go, relation_path_join_test.go, prefetch_related_test.go, m2m_prefetch_test.go and the to-many filter tests), the column stored at build time equals what the old `fkColumnFor` returned (write the comparison before deleting nothing; keep the resolver as the single source of truth).
- `SelectRelated` on a relation with no matching FK field still returns the rows (no error) and emits exactly one warning for two queries (capture slog with a buffer handler via `slog.SetDefault` and restore it with `t.Cleanup`).
- All existing tests pass unchanged; do not edit existing test expectations.

## Acceptance (from the module root)
    gofmt -l orm                                      # empty
    go build ./... && (cd ../examples/ecommerce && go build ./...)
    go vet ./orm/... ./admin/... ./api/... ./filter/...
    staticcheck ./orm/...                             # no output
    go test -race ./orm/... ./admin/... ./api/... ./filter/... -count=1
    ~/go/bin/golangci-lint run --timeout=5m --config=../.golangci.yml ./orm/...   # 0 issues

Report (max 25 lines): fields added, resolver location, equivalence test result, warning test, confirmation that no existing test expectation changed, each command result.
