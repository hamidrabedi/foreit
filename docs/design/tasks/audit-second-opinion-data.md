You are in the repository root /home/hamid/Other/projects/foreit-wt/refactor-audit. READ-ONLY task: do not edit, create or delete files and do not run git write commands. You may run read-only shell commands (grep, sed -n, go vet, go list).

You are a senior Go reviewer who knows Django's ORM, ent, sqlc, bob, goqu, squirrel, golang-migrate, goose and atlas. Give an independent refactor review of the data layer of this framework:
forge/orm, forge/db (including forge/db/dialect and forge/db/migrate/**), forge/filter (including filters/ and widgets/), forge/schema, forge/codegen, forge/config.

Known context (do not re-report these unless you disagree): Union/Intersection/Difference in forge/orm/queryset.go are no-ops; FK columns are resolved by name guessing (fkColumnFor, reverse lookup, prefetch.go); SQL builder emits $N placeholders then rebinds for SQLite by string rewriting; two migration runners (forge/db/migrations.go and forge/db/migrate/execute/executor.go) both wrap golang-migrate next to a home-grown migrate engine; forge/db/migrate/execute/checksum.go and forge/db/migrate/verify/checksum.go are near-identical; forge/filter/filters (1611 LOC) and forge/filter/widgets and forge/db/migrate/dependencies have no importers.

Report, each bullet as `[CATEGORY][SEVERITY] file:line — what — why it matters — action`, CATEGORY in DUP, DEAD, STUB, BUG, DESIGN, SIMPLIFY, TEST, LIB; SEVERITY CRITICAL (wrong results, data loss, SQL injection), MAJOR, MINOR:

1. Up to 20 BUG findings with a concrete failure scenario (SQL injection through identifiers, order-by or raw fragments; placeholder rewriting breaking literals; transactions not propagated; swallowed errors; races on package-level caches; dialect-specific SQL emitted for the wrong database; migrations that cannot roll back).
2. Up to 15 DUP/DEAD/STUB findings not listed in the known context.
3. Up to 10 DESIGN/SIMPLIFY findings: what the simplest correct architecture would be (for example: explicit relation metadata, dialect-owned placeholders, one migration engine), with the concrete files it replaces.
4. Up to 10 TEST findings: useless tests (no assertions, only getters, duplicate, permanently skipped) and critical untested code.
5. Up to 8 LIB recommendations: which proven library should replace which hand-rolled component (check forge/go.mod first), what it removes, and the migration risk.

Verify every claim by reading the code and cite file:line. No speculation, no preamble. Max 130 lines of markdown.
