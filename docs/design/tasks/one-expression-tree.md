Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

ROLE: You are the delegated coder: edit the files yourself; do not call any other agent. Use US spelling in comments.

# One expression tree in the ORM internals; the deprecated FieldExpr/QueryExpr API stays for users

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit forge/orm (non-generated files) and tests.

## Principles (mandatory)
- No exported identifier is removed or has its signature changed. `FieldExpr`, `NewFieldExpr`, `QueryExpr` and their methods are marked "Deprecated ... will be removed in v2.0": they stay until v2.0.
- Tests FIRST where behaviour is touched. gofmt; functions < 40 lines; naming after behaviour.

## Facts (verified)
- `orm.Expression` (expression.go) is schema-resolved, escapes identifiers and uses the dialect's placeholders; `QuerySet.Filter/Exclude` accept it.
- `QueryExpr` (query_expr.go) is built only by the deprecated `FieldExpr[T]`. Inside the ORM it is still used by `AnnotationExpr.Expr` / `NewAnnotation` (annotations.go), `SQLBuilder.BuildWhere(conditions, excludes []QueryExpr)` (sql_builder.go) and annotation SQL building in queryset.go.

## Change
1. Annotations accept any `Expression`: add `func NewExpressionAnnotation(name string, expr Expression) AnnotationExpr` and an `Expression` field on `AnnotationExpr` used when set. Keep `NewAnnotation(name string, expr QueryExpr)` working: convert the `QueryExpr` to an equivalent `Expression` (write an unexported adapter type implementing `Expression` that renders the same SQL through the builder's placeholder function), so there is one rendering path in queryset.go.
2. `SQLBuilder.BuildWhere` keeps its signature; implement it through the same adapter so the ORM renders one tree.
3. Test first: an annotation built with `NewExpressionAnnotation` and one built with `NewAnnotation` from a `FieldExpr` produce the same SELECT fragment for an equivalent comparison on SQLite and PostgreSQL builders; existing annotation tests keep passing.

## Acceptance (from the module root)
    gofmt -l orm                                      # empty
    go build ./... && (cd ../examples/ecommerce && go build ./...)
    go vet ./orm/... ./api/... ./admin/... ./filter/...
    staticcheck ./orm/...                             # no output (deprecation warnings inside orm itself are expected; report any others)
    go test -race ./orm/... ./api/... ./admin/... ./filter/... -count=1

Report (max 20 lines): API added, adapter, tests, each command result.
