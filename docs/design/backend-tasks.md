# Backend fix slices (taken over from the codex agent, 2026-09-13)

The codex agent ran five read-only audits, then started four fix slices in `/tmp`
worktrees and hit its usage limit. A reboot wiped `/tmp`, so **all four slices' uncommitted
work was lost**. Only PR #201 survived. Every slice below is re-implemented from codex's
verified audit findings, which were re-checked line by line against `origin/master` @ ca30f66.

Workflow per task: Claude writes the prompt (`tasks/b-*.md`) → an opencode free model
implements it in the slice's worktree → Claude verifies the diff plus
`gofmt`/`vet`/`go test -race` → Claude commits, pushes, and opens or updates the PR.
Worktrees: `/home/hamid/Other/projects/foreit-wt/<slice>` (never `/tmp`).

Status: `todo` → `dispatched` → `verified` → `pushed`

## Phase 1: critical security & data integrity

| # | Slice / branch | Task | Status |
|---|---|---|---|
| P1 | fix/jwt-authentication-verification | JWT signature + HS256 allowlist; propagate invalid-credential errors | **PR #201 green, mergeable** |
| A1 | fix/db-transaction-driver-integrity | `WithTx` swallows commit errors; PG→SQLite silent fallback; migration driver from global config | **PR #202** (A1+A2) |
| A2 | same | migration version overflow; non-atomic/colliding migration file generation; semicolon splitting in recovery | **on PR #202** |

## Phase 2: concurrency & lifecycle

| # | Slice / branch | Task | Status |
|---|---|---|---|
| B1 | fix/cache-concurrency-lifecycle | throttling + API caches (no lock / delete under RLock); unstoppable cleanup goroutine | **PR #203** (filter cache removed, see Planned) |
| B2 | same | atomic throttle increment (anon + user); trusted-proxy client IP (XFF/X-Real-IP spoofing) in `server/ratelimit.go` + `api/throttling/anon_rate.go` | **on PR #203** |
| B3 | same | rate-limit store: goroutine leak on `stop`, expiry-based eviction instead of "keep a random half" | **on PR #203** |
| B4 | fix/registry-settings-concurrency | unsynchronised plugin registry + global API settings | **PR #207** |

## Phase 3: filtering & ORM correctness

| # | Slice / branch | Task | Status |
|---|---|---|---|
| C1 | fix/filter-expression-correctness | `AndGroup`/`OrGroup` build no groups; `OrFilter` is a no-op alias | **PR #205** |
| C2 | same | `IN`/range with `[]T`, empty `IN` → invalid SQL, HTTP `field__in=` → nil, `isnull=false`, SQLite `EXTRACT`, case-insensitive prefix/suffix, unsigned field path | **on PR #205** (EXTRACT + unsigned deferred to C3) |
| C3a | same | year/month/day filters emit invalid SQL on every DB (ComparisonExpression has no date-part case); dialect-aware SQLBuilder (EXTRACT vs SQLite strftime) | spec ready (tasks/b-c3a-date-parts-dialect.md) |
| C3b | same | relation-path joins | todo (design first) |

## Phase 4: server/auth edge cases

| # | Slice / branch | Task | Status |
|---|---|---|---|
| D1 | fix/server-security-edge-cases | CSRF exemption raw-prefix match; `Redirect` panics without request; session auth accepts inactive/locked users | **PR #206** |
| D2 | same | query-string API keys written to access logs; admin login brute-force; expired sessions never purged; password trimming | **on PR #206** |

## Admin UI redesign

Draft **PR #204** (`feat/admin-ui-redesign`). Extracted ListFilterPanel, ListBulkToolbar, ListCell (ModelListPage 1343→1035, pushed). 7.3b route lazy-loading pushed (initial JS 445→430 kB gz, still over 300). 7.3c lucide namespace import removed (committed, size to be measured). In flight: 5.4 backend validation errors in forms. Next: 4.1d toolbar + pagination.

## Phase 5: admin integration (frontend-touching, after the UI redesign lands)

Custom admin mount prefix end to end; surface backend validation `details` in forms;
inert notification wiring; React Query v5 mutation-callback context argument.

## Planned / deferred (not now)

### Filter result caching (`forge/filter/cache.go`)

The owner has decided against caching filters for now: "we will implement that later". The
fixes were removed from PR #203. `FilterCache` has no callers today. When filter caching is
built, fix these known defects in the existing type first:

- `GetParsedTree` / `GetCompiledSQL` / `GetMetadata` `delete` expired entries while holding
  only `RLock`, a concurrent map write under a read lock. Re-check and delete under the write lock.
- `NewFilterCache` starts a cleanup goroutine that can never stop. Add an idempotent `Close()`.
- Cached values (`*FilterNode`, metadata maps) are returned by reference. Define a copy
  contract so callers cannot mutate shared cached trees.
- Decide on invalidation (schema change, per-model TTL) before wiring it into request paths.

## Known pre-existing failures (not ours)

- `./orm` `TestPrefetchRelated_Integration`: failed on the clean-master baseline 2026-09-13, passed on the A1 run. Treat it as flaky, not ours.
