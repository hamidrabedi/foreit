Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

ROLE: You are the delegated coder. The rule in CLAUDE.md/AGENTS.md that forbids the primary agent from editing code applies to Claude, not to you: edit the files yourself. Do NOT call any other agent. Use US spelling in comments.

# Split forge/admin/core/admin.go (1209 lines) by responsibility; break up ListObjects

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit files in forge/admin/core (and add tests there). Other agents edit other packages concurrently: ignore build errors outside forge/admin/core.

## Step 1 — pure move (no behaviour change)
Move declarations out of admin.go into new files in package core. Do not rename anything or change any body; only imports change.
- admin.go: `Admin` type, `NewAdmin`, `SetDB`, the small getters (`ModelName`, `ModelType`, `Schema`, `Manager`, `ModelSchema`, `Config`, `ManagerInterface`, `ConfigInterface`, `PageType`), `GetMetadata`, `GetQueryset`, `applyConfigDefaults`.
- admin_list.go: `ListObjects`, `orderableFields`, `sanitizeOrdering`, `isSafeOrderField`, `parseLookup`, `Autocomplete`, `ObjectLabels`.
- admin_crud.go: `GetObject`, `validateData`, `CreateObject`, `UpdateObject`, `DeleteObject`, `SaveModel`, `DeleteModel`, `safeGetObjectByID`, `decodeData`, `stringToDateTimeHook`, `getObjectID`, `getObjectLabel`, `toInt64`.
- admin_actions.go: `ExecuteAction`.
- admin_history.go: `GetHistory`, `LogAction`, `resolveUserID`.
- admin_permissions.go: `isSuperuser`, `HasAddPermission`, `HasChangePermission`, `HasDeletePermission`, `HasViewPermission`, `HasModulePermission`, `getPermissionsMetadata`.
Run acceptance; everything must pass before step 2.

## Step 2 — ListObjects under 50 lines (D8)
`ListObjects` is 164 lines. Extract, inside admin_list.go, named steps with the same behaviour: search filter, field lookup filters (`field__lookup` query params), ordering, count + pagination, row serialization. `ListObjects` calls them in order. No change to what it returns for any input.
Before refactoring, add a table-driven test `TestListObjects_FiltersOrderingPagination` (in-memory SQLite, see existing admin/core tests for setup, e.g. write_fields_test.go) covering: search across search fields, an `__icontains` lookup, an `__in` lookup, ordering asc/desc with an unsafe field ignored, page 2 of size 2 with the total count. Run it against the old code (must pass), refactor, keep it passing.
The two identical permission checks `HasChangePermission`/`HasDeletePermission`/`HasViewPermission` (22 lines each, differing only in the config hook and codename): extract one unexported helper they all call, same behaviour.

## Acceptance (from the module root)
    gofmt -l admin                                    # empty
    go build ./admin/... && go vet ./admin/...
    staticcheck ./admin/...                           # no output
    go test -race ./admin/... -count=1
    wc -l admin/core/admin*.go                        # every file under 600 lines

Report (max 20 lines): files and line counts, new test, each command result.
