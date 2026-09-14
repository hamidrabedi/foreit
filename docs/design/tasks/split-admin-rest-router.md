Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

ROLE: You are the delegated coder. The rule in CLAUDE.md/AGENTS.md that forbids the primary agent from editing code applies to Claude, not to you: edit the files yourself. Do NOT call any other agent. Use US spelling in comments.

# Split forge/admin/api/rest/router.go (1944 lines) by responsibility — pure moves, no behaviour change

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit files in forge/admin/api/rest (non-test files). Other agents edit forge/orm, forge/db and forge/identity at the same time: ignore build errors there.

## Rules (mandatory)
- Pure move: cut functions, types, vars and consts out of router.go into new files in the same package. Do not rename anything, do not change any function body, signature, comment text or order of statements inside a function. The only allowed edits are package clauses and import lists.
- File names follow the repo naming rules: lowercase, underscore-separated, named after the responsibility, no `_helpers`/`_utils`/`_impl` suffixes.
- Target layout (move each declaration to the file whose responsibility it serves; if one does not fit, leave it in router.go):
  - router.go: `Router` struct, `NewRouter`, route registration (`RegisterRoutes`/`registerModelRoutes` and similar), middleware wiring.
  - auth_handlers.go: login/logout/session/me handlers, `adminCredentials`, `secureEqual`, `decodeLoginPayload`, `loginKeys`, `checkLoginRateLimit`, `authenticateAdmin`, and the `adminSessionStore` type with its methods (or put the store in session_store.go if it is larger than ~120 lines).
  - crud_handlers.go: list/retrieve/create/update/delete/bulk handlers and their request parsing (`normalizePathID`, `normalizeBulkID`, `bulkFailure`, validation detail helpers).
  - metadata_handlers.go: model list, metadata, dashboard/environment, display labels.
  - export_handlers.go: export handler and `exportListFields` and CSV/JSON writers.
  - saved_views.go: `savedViewStore`, `savedView`, `newSavedViewStore` and the saved-view handlers.
  - responses.go: `respondJSON`, `respondError` and other response writers.
- Every resulting file must be under 600 lines; if crud_handlers.go is larger, split list/retrieve (crud_read_handlers.go) from create/update/delete/bulk (crud_write_handlers.go).

## Acceptance (from the module root)
    gofmt -l admin                                    # empty
    go build ./admin/... && go vet ./admin/...
    staticcheck ./admin/...                           # no output
    go test -race ./admin/... -count=1
    wc -l admin/api/rest/*.go | sort -n | tail -12    # no non-test file over 600 lines

Also confirm the move is pure: the total number of non-blank, non-import, non-package lines across the package's non-test files must equal the count before your change (count before you start: grep -v -E '^\s*$|^package |^import|^\t"|^\)$' admin/api/rest/*.go excluding _test.go files; report both numbers).

Report (max 20 lines): files created with line counts, line-count check before/after, each command result.
