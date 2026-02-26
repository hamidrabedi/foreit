# Agent Goals

## Primary Goal
Make Forge production-ready by closing framework gaps and proving features end-to-end.

## What Big Frameworks Do (Django, Laravel, ASP.NET)
These frameworks are reliable because they focus on:
- **Stability over features** — Don't break existing APIs
- **Comprehensive testing** — Every feature has tests
- **Documentation** — Clear docs for every public API
- **Backward compatibility** — Deprecate gracefully
- **UI/UX polish** — Admin looks good and works intuitively
- **Edge cases** — Handle errors gracefully

## Priorities (In Order)
1. **Stability** — Don't break existing functionality
2. **Testing** — Every change needs tests
3. **Documentation** — Update docs when adding/changing features
4. **UI/UX** — Admin and ecommerce must look good
5. **Example Parity** — Ecommerce proves the framework works
6. **Features** — Only add what NEEDS to be added
7. **Roadmap** — Keep PLAN.md current

## Mandatory Coverage Areas
- **Admin:** CRUD flows, filters, actions, auth, permissions, UI/UX stability
- **ORM:** Query building, relations, aggregation, update/select/prefetch
- **Schema:** Model definitions, traits, hooks, relation integrity
- **Migrations:** Generation, apply, rollback, drift/safety checks
- **API/Server:** Serializers, viewsets, auth, permissions, pagination, error handling
- **Example Parity:** Ecommerce must work end-to-end

## Definition of Done (Rolling)
1. Every new/updated feature has test coverage
2. Feature works in ecommerce example
3. Tests pass (`go test ./...`)
4. UI/UX tested manually (forms, lists, filters work)
5. Documentation updated if needed
6. `STATUS.md`, `HISTORY.md`, `PLAN.md` updated

## Current Priorities
- Don't add features unless needed — focus on stability
- Match Django/Laravel quality in admin UI
- Ensure ecommerce runs and looks good
- Fix real bugs, not nice-to-haves
