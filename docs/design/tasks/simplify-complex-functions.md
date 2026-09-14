Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

ROLE: You are the delegated coder. The rule in CLAUDE.md/AGENTS.md that forbids the primary agent from editing code applies to Claude, not to you: edit the files yourself. Do NOT call any other agent. Use US spelling in comments.

# Break up three over-complex functions without changing behaviour (D8)

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit: forge/server/static.go, forge/validate/validator.go, forge/admin/api/rest/metadata_handlers.go, and new or existing _test.go files next to them. Another agent is splitting forge/orm/queryset.go at the same time: do not touch forge/orm; if a build error points into forge/orm, wait a minute and rebuild.

## Rules (mandatory)
- No behaviour change and no exported signature change. Unexported helpers may be added.
- Characterization tests FIRST: for each function, write table-driven tests that pin its current outputs (including edge cases and error paths) and run them against the unchanged code; they must pass before you refactor. Keep them passing after. Do not edit existing test expectations.
- Target: each of the three functions and every helper you extract has cyclomatic complexity <= 15 (`~/go/bin/gocyclo -over 15 <file>` prints nothing for them) and is under 50 lines. Prefer table-driven maps/slices for repetitive switch cases.

## Functions (verified with gocyclo)
1. `server.StaticFS` (forge/server/static.go:112, complexity 43). Characterize: file served with correct content type, directory index handling, not-found, path traversal attempts (`../`), cache headers and any options the function reads. Use `httptest` with an `fstest.MapFS` or a temp dir.
2. `getErrorMessage` (forge/validate/validator.go:197, complexity 42): maps validation tags to messages. Characterize every tag branch (read the switch and list every case, including the default) with its parameter formatting. Replace the switch with a `map[string]func(fieldName, param string) string` (or equivalent) plus the default.
3. `(*Router).attachDisplayLabels` (forge/admin/api/rest/metadata_handlers.go:138, complexity 42, 93 lines). Characterize with the existing admin test fixtures (read list_display_test.go and object_labels_test.go for setup): choice labels, foreign key labels, missing related rows, nil values, and fields without labels. Extract named steps (collect ids per relation, load labels, apply labels to rows).

## Acceptance (from the module root)
    gofmt -l server validate admin                  # empty
    go build ./server/... ./validate/... ./admin/... && go vet ./server/... ./validate/... ./admin/...
    staticcheck ./server/... ./validate/... ./admin/...   # no output
    go test -race ./server/... ./validate/... ./admin/... -count=1
    ~/go/bin/gocyclo -over 15 server/static.go validate/validator.go admin/api/rest/metadata_handlers.go   # nothing for the three functions or their new helpers
    ~/go/bin/golangci-lint run --timeout=5m --config=../.golangci.yml ./server/... ./validate/... ./admin/...   # 0 issues

Report (max 20 lines): characterization tests added (names), new helper names, gocyclo before/after per function, each command result.
