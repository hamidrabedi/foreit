Your workspace is /home/hamid/Other/projects/foreit-wt/refactor-audit. Always use absolute paths. Never search outside it. READ-ONLY task: do not edit, create or delete any file, and do not run git.

You are a senior Go reviewer giving a second opinion on a refactor audit of the framework in
/home/hamid/Other/projects/foreit-wt/refactor-audit/forge (packages: api, admin (Go only, skip admin/ui/web), identity, server, netutil, log, validate, registry, cli, errors).

## Part 1: validate these claims

For EACH claim answer `CONFIRMED`, `WRONG` or `PARTLY`, with file:line evidence and one sentence. Re-rate severity if you disagree.

1. api/errors/idempotency_stores.go ~187-231: DatabaseStore Get/Set/Delete build SQL but never execute it; Set returns nil (CRITICAL: silent loss of idempotency).
2. identity/backends/password.go ~60-85: unknown user returns before bcrypt compare (timing enumeration), and inactive/locked status is returned before the password is checked (status disclosure without password). MAJOR.
3. admin/api/rest/router.go ~81-91: CORS echoes any origin with AllowCredentials true. Claimed MINOR because the admin API authenticates only via `Authorization: Bearer` (no cookies). Check whether any cookie/session auth reaches this router.
4. api/viewset.go, viewset_enhanced.go, viewset_enhanced_integrated.go, viewset_config.go: parallel viewset base implementations; BaseViewSet only used by CLI scaffolds.
5. Three rate limiters: server/ratelimit.go (x/time/rate), api/throttling (hand-rolled fixed window), admin/api/rest/login_limiter.go.
6. Four error systems: api/errors, api/exceptions, forge/errors, server/errors.go.
7. admin/core/notifications.go NotificationHub/SSEHandler never routed; sets Access-Control-Allow-Origin: *.
8. log/logger.go ~136-142: remote output silently returns a no-op writer.
9. registry/plugin.go ~233-241: applyAdminExtensions/applyAPIExtensions are no-op placeholders.
10. api/viewset.go getManagerFromModel ~756 always returns reflect.Value{}.
11. cli/commands/admin/createsuperuser.go ~120: password read with echo.
12. These packages have zero importers in forge/, tests/, examples/: admin/codegen, admin/utils, api/caching, api/serializers, api/versioning, cli/internal, log/hooks. Check each (including reflection/registration or templates in cli that generate imports of them).
13. identity/password.go is only a re-export shim of identity/utils/password.go.
14. api/authentication/* and identity/backends/* are two unrelated authentication abstractions.

## Part 2: what the audit missed

List up to 15 additional CRITICAL/MAJOR findings in scope (bugs with a concrete failure scenario, security holes, duplicated logic, dead code, bad design, over-engineering, useless tests), each as:
`[CATEGORY][SEVERITY] file:line — what — why — action`.

## Part 3: libraries

Up to 8 places where a well-known Go library (or stdlib, Go 1.26) should replace hand-rolled code. Check forge/go.mod first. Format: `hand-rolled file(s) → library — what it removes — risk`.

Output plain markdown, max 120 lines. No preamble.
