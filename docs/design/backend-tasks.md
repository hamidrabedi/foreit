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
| A1 | fix/db-transaction-driver-integrity | `WithTx` swallows commit errors; PG→SQLite silent fallback; migration driver from global config | **PR #202** (agy; 1370 tests pass) |
| A2 | same | migration version overflow; non-atomic/colliding migration file generation; semicolon splitting in recovery | todo (spec after A1) |

## Phase 2: concurrency & lifecycle

| # | Slice / branch | Task | Status |
|---|---|---|---|
| B1 | fix/cache-concurrency-lifecycle | 3 unsafe caches (no lock / delete under RLock); unstoppable cleanup goroutines | dispatched |
| B2 | same | atomic throttle increment (anon + user); trusted-proxy client IP (XFF/X-Real-IP spoofing) in `server/ratelimit.go` + `api/throttling/anon_rate.go` | todo (spec after B1) |
| B3 | same | rate-limit store: goroutine leak on `stop`, expiry-based eviction instead of "keep a random half" | todo |
| B4 | same | unsynchronised plugin registries + global API settings | todo (needs audit re-read) |

## Phase 3: filtering & ORM correctness

| # | Slice / branch | Task | Status |
|---|---|---|---|
| C1 | fix/filter-expression-correctness | `AndGroup`/`OrGroup` build no groups; `OrFilter` is a no-op alias | todo |
| C2 | same | `IN`/range with `[]T`, empty `IN` → invalid SQL, HTTP `field__in=` → nil, `isnull=false`, SQLite `EXTRACT`, case-insensitive prefix/suffix, unsigned field path | todo (spec after C1) |
| C3 | same | relation-path joins; dialect-aware SQL generation | todo (larger — design first) |

## Phase 4: server/auth edge cases

| # | Slice / branch | Task | Status |
|---|---|---|---|
| D1 | fix/server-security-edge-cases | CSRF exemption raw-prefix match; `Redirect` panics without request; session auth accepts inactive/locked users | todo |
| D2 | same | query-string API keys written to access logs; admin login brute-force; expired sessions never purged; password trimming | todo |

## Phase 5: admin integration (frontend-touching, after the UI redesign lands)

Custom admin mount prefix end to end; surface backend validation `details` in forms;
inert notification wiring; React Query v5 mutation-callback context argument.

## Known pre-existing failures (not ours)

- `./orm` `TestPrefetchRelated_Integration`: failed on the clean-master baseline 2026-09-13, passed on the A1 run. Treat it as flaky, not ours.
