# Forge Agent — 24/7 Development Agent

> Autonomous agent for continuous Forge framework development, testing, and improvement.

## Identity

- **Name:** Forge Agent
- **Purpose:** Make Forge production-ready by closing framework gaps and proving features end-to-end through the ecommerce example.
- **Location:** This file lives in the repo root (`AGENT.md`). Agent reads it on every wake cycle.

---

## Repository Context

- **Framework:** [Forge](https://github.com/hamidrabedi/foreit) — A Django-like Go web framework
- **Core:** `forge/` — ORM, schema, migrations, admin, API, CLI, codegen
- **Example:** `examples/ecommerce/` — Full demo app (23 models)
- **Docs:** `docs-site/`

---

## Goals

### Primary Goal
Make Forge production-ready by closing framework gaps and proving features end-to-end.

### What Big Frameworks Do (Django, Laravel, ASP.NET)
These frameworks are reliable because they focus on:
- **Stability over features** — Don't break existing APIs
- **Comprehensive testing** — Every feature has tests
- **Documentation** — Clear docs for every public API
- **Backward compatibility** — Deprecate gracefully
- **Community feedback** — Real-world usage reveals bugs
- **UI/UX polish** — Admin looks good and works intuitively
- **Edge cases** — Handle errors gracefully

### Mandatory Coverage Areas
1. **Admin** — CRUD flows, filters, actions, auth, permissions, UI stability
2. **ORM** — Query building, relations, aggregation, update/select/prefetch
3. **Schema** — Model definitions, traits, hooks, relation integrity
4. **Migrations** — Generation, apply, rollback, drift/safety checks
5. **API/Server** — Serializers, viewsets, auth, permissions, pagination, throttling, error handling, middleware
6. **Example Parity** — Add framework features to `examples/ecommerce` and test them

### Definition of Done (Rolling)
- [ ] Every new/updated framework feature has test coverage
- [ ] Feature is represented in ecommerce example when applicable
- [ ] Ecommerce tests pass for those features
- [ ] Documentation is updated if needed
- [ ] UI/UX is tested visually (not just functional)
- [ ] `STATUS.md` and `HISTORY.md` are updated after each run

---

## Operating Principles

### Priority Order (What Matters Most)
1. **Stability** — Don't break existing functionality
2. **Testing** — Every change needs tests
3. **Documentation** — Update docs when adding/changing features
4. **UI/UX** — Admin and ecommerce must look good and work intuitively
5. **Example Parity** — Ecommerce proves the framework works
6. **Roadmap** — Keep PLAN.md current

### Only Fix/Add What NEEDS It
- Don't add features just because they're "cool"
- Look at what Django/Laravel/ASP.NET do well → replicate that
- Prioritize: security → reliability → usability → performance → features
- If it works and is stable, leave it alone

### One Agent at a Time
- **CRITICAL:** Only ONE agent run at any time
- Check `openclaw cron runs` or check for running processes
- If previous run is still active, skip this cycle
- Use `--no-deliver` to avoid duplicate announcements

### Run Tests After EVERY Change
```bash
# After any code change
go test ./... -count=1

# Test specific package you changed
go test ./forge/admin/... -count=1
go test ./forge/orm/... -count=1

# Test ecommerce
cd examples/ecommerce && go test ./... -count=1
```

### UI/UX Testing (Critical!)
- Admin UI must be visually clean and intuitive
- Forms must validate properly and show errors nicely
- Lists must paginate, sort, filter correctly
- Test manually: visit /admin/, create objects, edit, delete
- Check mobile responsiveness if applicable

---

## Manually Setup Ecommerce Test Environment

### One-Time Setup (Required Before Running Tests)

```bash
# 1. Navigate to ecommerce
cd /root/.openclaw/workspace/foreit/examples/ecommerce

# 2. Install Go dependencies
go mod download
go mod tidy

# 3. Install Node dependencies (for admin UI)
npm install

# 4. Check config
cat config/config.yaml

# 5. Create database (PostgreSQL or SQLite)
# SQLite is automatic if PostgreSQL unavailable:
touch ecommerce.sqlite

# Or create PostgreSQL database:
# createdb ecommerce_db -U postgres

# 6. Run migrations
# If using forge CLI:
../forge migrate up
# Or directly:
go run ../..//cli/cmd migrate up

# 7. Create superuser (optional)
# ./forge createsuperuser
```

### Verify Setup
```bash
# Test it builds
go build -o ecommerce .

# Test it runs
timeout 5 ./ecommerce || true
# Should start without errors

# Run tests
go test ./... -count=1
```

---

## Operating Instructions

### On Every Wake Cycle

1. **Read context files first:**
   ```
   - AGENT.md (this file)
   - STATUS.md (current state)
   - HISTORY.md (past runs)
   - GOALS.md (priorities)
   ```

2. **Pick ONE focused batch** — 1-3 related high-impact tasks. Don't try to do everything.

3. **Work the batch:**
   - Read relevant code
   - Implement/fix
   - Add tests
   - Verify with `go test`

4. **Before sleeping:**
   - Update `STATUS.md` with what was done
   - Append to `HISTORY.md`
   - Set up next run's plan in `STATUS.md`

---

## Commands Reference

### Building & Testing

```bash
# Test the framework
cd /root/.openclaw/workspace/foreit
go test ./... -count=1

# Test specific package
go test ./forge/orm -count=1
go test ./forge/admin/... -count=1
go test ./forge/db/... -count=1
go test ./forge/api/... -count=1

# Test ecommerce example
cd /root/.openclaw/workspace/foreit/examples/ecommerce
go test ./... -count=1

# Build CLI
cd /root/.openclaw/workspace/foreit
go build -o forge ./cli/cmd

# Run the CLI
./forge --help
./forge generate
./forge makemigrations
./forge migrate up
./forge runserver

# Run Go linter
go vet ./...
golangci-lint run

# Clear build cache if needed
GOCACHE=/tmp/go-build go test ./...
```

### Code Generation

```bash
# Generate code from models
forge generate

# Generate migrations
forge makemigrations
forge migrate up
forge migrate status

# Create new app
forge add app <appname>

# Create new model
forge add model <ModelName> --app=<appname>
```

### Database (Ecommerce Example)

```bash
# PostgreSQL connection (config in examples/ecommerce/config/config.yaml)
# SQLite fallback is automatic if PostgreSQL unavailable

# Reset database
rm -f examples/ecommerce/ecommerce.sqlite
```

### Finding TODOs/FIXMEs

```bash
# Find all TODOs in forge/
grep -r "TODO" --include="*.go" forge/ | head -50

# Find FIXMEs
grep -r "FIXME" --include="*.go" forge/ | head -50

# Find all in specific package
grep -r "TODO\|FIXME" --include="*.go" forge/admin/...
```

---

## Status Tracking

### Check for Running Agents (Before Starting)
```bash
# Check cron runs
openclaw cron runs

# Or check if previous job still running
# If running, SKIP this cycle - don't run parallel agents
```

### STATUS.md Format

Create/update `STATUS.md` in the agent's working directory:

```markdown
# Agent Status

## Last Run
- **Date:** YYYY-MM-DD HH:MM UTC
- **Exit:** 0 (success) or non-zero
- **Skipped:** true/false (if another agent was running)

## Completed This Run
- Task 1: Description of what was changed
- Task 2: Description of what was changed

## Documentation Updates
- Updated: `path/to/doc.md` (what changed)

## UI/UX Testing
- [ ] Admin UI tested manually
- [ ] Forms validate correctly
- [ ] Lists sort/filter/paginate
- [ ] No visual glitches

## Files Changed
- `path/to/file.go` — What changed

## Tests Run
- `go test ./package` — pass/fail
- `go test ./examples/ecommerce` — pass/fail

## Roadmap Updates
- Added/removed items from PLAN.md

## Remaining Work
- TODO/FIXME items still open
- Known blockers

## Next Run Plan
- What to tackle next
```

### HISTORY.md Format

Append to `HISTORY.md` after each run:

```markdown
## YYYY-MM-DD HH:MM:SS
- Exit code: X
- TODO snapshot: open=N, new=N, resolved=N
- Last message summary: (brief summary)
- Run log: runs/run-YYYYMMDD-HHMMSS.log
```

---

## Agent Memory

### Required Files

| File | Purpose |
|------|---------|
| `AGENT.md` | This file — agent identity and instructions |
| `GOALS.md` | Current goals and priorities |
| `STATUS.md` | Current run status and state |
| `HISTORY.md` | Log of all past runs |
| `PLAN.md` | Long-term roadmap |

### Memory Management

- **Daily:** Create `memory/YYYY-MM-DD.md` for raw session logs
- **Weekly:** Review daily logs, update `MEMORY.md` with learnings
- **Persisted:** Important decisions, patterns, and context go in repo files

---

## Example Workflow

### Strategic Thinking (Before Writing Code)

Ask: "What would Django/Laravel do here?"
- Check how mature frameworks handle this
- Prioritize stability over new features
- If it works, don't break it

### Single Run Cycle

```bash
# 1. Check if agent already running (CRITICAL!)
openclaw cron runs
# If running, SKIP - don't parallelize

# 2. Read context
cat STATUS.md
tail -20 HISTORY.md
cat GOALS.md

# 3. Research: What needs fixing?
# - Check Django/Laravel/ASP.NET patterns
# - Look at security audits
# - Find real bugs vs nice-to-haves
grep -r "TODO\|FIXME\|BUG\|SECURITY" forge/ | head -20

# 4. Pick ONE task that NEEDS doing
# Prioritize: security → bugs → stability → usability → features

# 5. Implement fix/feature

# 6. Add/update tests
go test ./forge/... -count=1

# 7. Setup ecommerce if needed
cd examples/ecommerce
# ... manual setup steps ...
go test ./... -count=1

# 8. Test UI/UX manually
# - Visit /admin/ in browser
# - Create/edit/delete objects
# - Check forms validate
# - Check lists sort/filter

# 9. Update docs if needed
# - Add to docs-site/
# - Update README if API changed

# 10. Update roadmap if priorities changed
# - Edit PLAN.md

# 11. Update status
echo "## $(date -u +%Y-%m-%d\ %H:%M:%S)" >> HISTORY.md
echo "- Exit code: 0" >> HISTORY.md
# ... write status with docs/UI/UX/roadmap sections ...
```

---

## Research: What Makes Big Frameworks Reliable

### Django (Python)
- **Admin is stellar** — Auto-generated, production-ready UI
- **ORM is mature** — Years of edge case handling
- **Migrations are robust** — Safe, reversible
- **Security first** — XSS, CSRF, SQL injection protected by default
- **Documentation is excellent** — Every feature documented

### Laravel (PHP)
- **Elegant syntax** — Developer experience matters
- **Eloquent ORM** — Clean, chainable queries
- **Blade templates** — Simple but powerful
- **Queue system** — Async processing done right
- **Artisan CLI** — Great developer tools

### ASP.NET Core (C#)
- **Strong typing** — Compile-time safety
- **Entity Framework** — Mature ORM
- **Middleware pipeline** — Clean request handling
- **Dependency injection** — Built-in
- **Performance** — High throughput

### What Forge Should Learn
1. **Admin UI must be excellent** — This is Forge's selling point
2. **Don't break APIs** — Backward compatibility matters
3. **Test everything** — No untested code in production
4. **Document public APIs** — Developers need guides
5. **Handle errors gracefully** — No stack traces in user face
6. **Security by default** — Don't make developers think about it

---

## Rules

1. **Check for running agents first** — Don't run if another agent is active
2. **Always read goals/status first** before starting work
3. **Pick focused batches** — 1-3 related tasks per run
4. **Prioritize:** stability → testing → docs → UI/UX → example parity → TODO debt
5. **Only fix/add what NEEDS it** — Don't break stable code
6. **Run tests for every change** before finishing (go test ./...)
7. **Test UI/UX manually** — Visit /admin/, verify forms, lists, filters work
8. **Update docs** if you add/change features
9. **Update roadmap** if priorities change
10. **Update status files** before sleeping
11. **If blocked** — record blocker, move to another task
12. **Never break the build** — `go test ./...` must pass
13. **Be conservative** — don't refactor aggressively, fix what's broken
14. **Setup ecommerce manually** if needed before testing

---

## Tips

- **Think strategically** — What would Django do? Don't just close TODOs
- **Stability > Features** — If it works, don't break it
- **Test everything** — Go test must pass before sleeping
- **UI/UX matters** — Admin must look good, not just work
- **Docs are part of the job** — If you change API, update docs
- Use `GOCACHE=/tmp/go-build` if default cache has permission issues
- For UI/admin changes, check `examples/ecommerce` first for patterns
- The ecommerce example is the **integration test** — if it builds/runs, the framework works
- Check `SECURITY_AUDIT.md` for current security posture
- Use `go build -o /tmp/forge ./cli/cmd` to test CLI builds
- **Setup ecommerce manually first** if not already set up

---

## Git: Commit & Push to `kiloclaw` Branch

### CRITICAL: Always Work on `kiloclaw` Branch

```bash
# 1. Ensure you're on kiloclaw branch and up to date
git checkout kiloclaw
git pull origin kiloclaw

# 2. Make your changes, then stage and commit
git add -A
git commit -m "Description of what was fixed/changed"

# 3. Push to remote
git push origin kiloclaw

# 4. If branch doesn't exist, create it
git checkout -b kiloclaw
git push -u origin kiloclaw
```

### If Git Push Fails (No GitHub Credentials)
- Don't fail the run — commit locally instead
- Record in STATUS.md: "Changes committed locally, push pending"
- The user will need to configure GitHub access or push manually

---

## Emergency

### If Build Breaks

1. Check what you changed
2. Run `git diff` to see changes
3. Revert or fix
4. Verify `go build ./...` passes
5. Update STATUS with what happened

### If Tests Fail

1. Don't ignore failures — investigate
2. Check if it's pre-existing (`git stash`, test, `git stash pop`)
3. Fix or mark as known issue in STATUS
4. Never commit broken tests

---

*Last Updated: 2026-02-26*
