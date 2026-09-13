# Admin UI v3 — Execution Plan

Spec: [`admin-ui-system.md`](./admin-ui-system.md). Implementers read the spec first.

Working dir for all tasks: `forge/admin/ui/web`
Gate after **every** task: `npx tsc -b` && `npx vitest run` && `npx eslint .`

Status: `todo` → `dispatched` → `verified` / `rejected`

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
| 3.4b | `cmdk` swap, deferred — must preserve all 4 tests + testids | `src/components/layout/` | todo |
| 3.5 | `sonner` Toaster replacing `use-toast` (VB-09) | `src/components/ui/` | verified |
| 3.6 | `framer-motion` `layoutId` active-nav indicator — nav only, no route transitions | `src/components/layout/` | verified |

## P4 — Data surfaces
| # | Task | Files | Status |
|---|---|---|---|
| 4.1 | Extract `DataGrid` shell out of `ModelListPage` | `src/components/data/` | todo |
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

## P6 — Backend
| # | Task | Files | Status |
|---|---|---|---|
| 6.1 | **VB-06** emit `read_only: true` for auto timestamps in admin metadata | `forge/admin/**.go` | todo |

## P7 — Verification
| # | Task | Files | Status |
|---|---|---|---|
| 7.1 | Playwright sweep: routes × light/dark × 375/1280, assert no h-scroll + no console errors + rendered | `e2e/design-sweep.spec.ts` | authored |
| 7.2 | Lint guards: no raw hex, no Tailwind palette colors, no `text-[Npx]` in `src/` | `eslint.config.js` | verified |
| 7.3 | **Bundle budget breach.** Task authored (`tasks/t7.3.md`): manualChunks + lazy recharts. Measured 2026-09-12: JS 1,813 kB raw / **501 kB gzip** vs the 300 kB app-page cap in `rules/web/performance.md`. CSS 105 kB / 27.8 kB gzip is fine. Route-level code splitting + dynamic import of recharts. | `vite.config.ts`, routes | todo |
