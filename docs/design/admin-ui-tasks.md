# Admin UI v3 — Execution Plan

Spec: [`admin-ui-system.md`](./admin-ui-system.md). Implementers read the spec first.

Working dir for all tasks: `forge/admin/ui/web`
Gate after **every** task: `./node_modules/.bin/tsc -b` && `./node_modules/.bin/vitest run` && `./node_modules/.bin/eslint src/` (use local bins — `npx` is intercepted by rtk)

Branch / PR: `feat/admin-ui-redesign` → draft **PR #204**. Current gate: tsc clean, vitest 46 passed / 1 skipped, eslint 0 errors.

Status: `todo` → `dispatched` → `verified` / `rejected`

### Bundle measurement log (gzipped JS in `dist/index.html` modulepreload)

| Date | After | Entry `index` | Eager total | Notes |
|---|---|---|---|---|
| 2026-09-12 | baseline | — | 501 kB (single chunk) | everything in one chunk |
| 2026-09-13 | 7.3 vendor split + lazy recharts | 272 kB | ~445 kB | |
| 2026-09-13 | 7.3b lazy routes | 205 kB | ~430 kB | lucide namespace import still eager |
| 2026-09-13 | 7.3c StatsWidget named icons | **62 kB** | **357 kB** | `charts` (113 kB) still wrongly preloaded; without it ≈244 kB → task 7.3d |

## P1 — Foundation

| # | Task | Files | Status |
|---|---|---|---|
| 1.1 | Token layer: surface ramp (fixes dark `--card == --background`), border scale, semantic state, chart ramp, type scale, density, motion, radii. Delete `glass-lite`/`glass-thick`/`premium-shadow`. | `src/index.css` | verified |
| 1.2 | Expose every new token to Tailwind: colors, fontSize, spacing, borderRadius, boxShadow, transitionTimingFunction | `tailwind.config.js` | verified |
| 1.4 | **Validated chart ramp** (Appendix D) — replaces eyeballed values that failed CVD check | `src/index.css` | verified |
| 1.3 | Self-host Inter + JetBrains Mono via `@fontsource`; wire `--font-sans` / `--font-mono` | `package.json`, `src/main.tsx`, `src/index.css` | verified |

## P2 — Primitives
| # | Task | Files | Status |
|---|---|---|---|
| 2.1 | `PageHeader` | `src/components/ui/page-header.tsx` | verified |
| 2.2 | `StatusBadge` (semantic tones) + purge raw palette colors from badges | `src/components/ui/status-badge.tsx` | verified |
| 2.3 | `EmptyState` + `EmptyValue` (VB-03) | `src/components/ui/empty-state.tsx` | verified |
| 2.4 | `Tooltip` (radix) + replace `title=` on icon buttons | `src/components/ui/tooltip.tsx` | verified |
| 2.5 | `StatTile` replacing `StatsCard` (VB-04 color tokens); old `StatsCard` deleted as dead code | `src/components/widgets/stat-tile.tsx` | verified |
| 2.6 | Retune `button`/`card`/`input`/`table` to new radii + density + no card shadow | `src/components/ui/*` | verified |

## P3 — Shell
| # | Task | Files | Status |
|---|---|---|---|
| 3.1 | **VB-11** hoist `SidebarItem` out of `AdminLayout` render body | `src/components/layout/` | verified |
| 3.2 | **VB-01** sidebar grouping rewrite | `src/components/layout/` | verified |
| 3.3 | Split `AdminLayout` (830→<300) into `Sidebar`/`TopBar`/`AdminLayout` (VB-12) | `src/components/layout/` | verified |
| 3.4a | **VB-07** verbose names in palette + token restyle (Appendix C) | `src/components/layout/GlobalSearch.tsx` | verified |
| 3.4b | `cmdk` swap — all 4 tests + testids preserved; jsdom polyfills moved to `src/test/setup.ts` (3.4c) | `src/components/layout/` | verified |
| 3.5 | `sonner` Toaster replacing `use-toast` (VB-09) | `src/components/ui/` | verified |
| 3.6 | `framer-motion` `layoutId` active-nav indicator — nav only, no route transitions | `src/components/layout/` | verified |

## P4 — Data surfaces
| # | Task | Files | Status |
|---|---|---|---|
| 4.1 | Split `ModelListPage` into `src/components/list/`: ListFilterPanel (4.1a), ListBulkToolbar (4.1b), ListCell (4.1c), ListToolbar + ListPagination (4.1d) — 1343→906 lines; 4.1e (table header, row actions, save-view dialog) to reach <800 | `src/components/list/` | 4.1a–d verified, 4.1e dispatched |
| 4.2 | **VB-02** FK label resolution (list + view) | `src/components/data/`, `ModelViewPage` | verified |
| 4.3a | `Select` primitive (radix wrapper, none existed) | `src/components/ui/select.tsx` | verified |
| 4.3b | **VB-05/VB-13** typed filters: radix Select + numeric min/max | `src/pages/ModelListPage.tsx` | verified |
| 4.3c | Convert remaining 2 raw `<select>` (saved-views, page-size) — out of scope in 4.3b | `src/pages/ModelListPage.tsx` | verified |
| 4.4 | **VB-08** sticky identity column + scroll-shadow | `src/components/data/` | verified |
| 4.5 | **VB-04** shared chart theme + tooltip | `src/components/widgets/` | verified |
| 4.6 | Dashboard recompose to `canvas` density, hero metric, editorial rhythm | `src/pages/DashboardPage.tsx` | verified |

## P5 — Forms
| # | Task | Files | Status |
|---|---|---|---|
| 5.1 | Extract `Field` + widget registry from `ModelUpsertPage` | `src/components/form/` | verified |
| 5.2 | Split `ModelUpsertPage` 1389 → <400 (VB-12) | `src/pages/` | verified |
| 5.3 | Recompose `ModelViewPage` to the new system | `src/pages/` | verified |

## P5 addendum
| # | Task | Files | Status |
|---|---|---|---|
| 5.4 | Backend validation errors were discarded: UI read `data.details`, API sends `data.error.details`. `parseApiError` (`src/api/errors.ts`, 10 tests) + form-level `role=alert` banner for `non_field_errors` | `src/api/errors.ts`, `ModelUpsertPage.tsx` | verified |

## P6 — Backend
| # | Task | Files | Status |
|---|---|---|---|
| 6.1 | **VB-06** emit `read_only: true` for auto timestamps in admin metadata — spec in `backend-handoff.md` BH-1 | `forge/admin/**.go` | todo (backend) |
| 6.2 | **VB-02 phase 2** `<field>__display` relation labels in list serializer — spec in `backend-handoff.md` BH-2 | `forge/admin/**.go` | todo (backend) |

## P7 — Verification
| # | Task | Files | Status |
|---|---|---|---|
| 7.1 | Playwright sweep: routes × light/dark × 375/1280, assert no h-scroll + no console errors + rendered | `e2e/design-sweep.spec.ts` | authored |
| 7.2 | Lint guards: no raw hex, no Tailwind palette colors, no `text-[Npx]` in `src/` | `eslint.config.js` | verified |
| 7.3 | **Bundle budget.** 7.3: vendor code-splitting (rolldown `codeSplitting.groups`) + lazy recharts. 7.3b: lazy non-entry routes (445→430 kB gz initial). 7.3c: StatsWidget `import * as LucideIcons` removed (namespace import defeated tree-shaking). Budget 300 kB gz — see measurement note below. | `vite.config.ts`, `src/routes/`, widgets | 7.3/7.3b/7.3c verified; 7.3d (eager charts preload) queued |
