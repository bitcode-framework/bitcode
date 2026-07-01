# ADR-014: Engine ↔ UI Feature Wiring

- **Status:** Accepted (2026-07-02 — all 9 features complete across 4 phases)
- **Date:** 2026-06-26
- **Decision Makers:** Stakeholder, UI Team, Engine Team
- **Supersedes:** —
- **Amends:** —
- **Related:** ADR-005 (Metadata-Driven UI / BcRenderer), ADR-006 (Universal Component Strategy), ADR-007 (Enterprise Universal UI API Contract), ADR-008 (Contract-First API Freeze)

## Context

BitCode engine currently ships 110/110 enterprise-grade components in `@bitcode/ui`
(Mitosis) + 130 widgets in `@bitcode/ui-flutter`, plus 122 legacy Stencil components.
However, the **engine itself is still SSR-based** (Go HTML templates) and many components
that are fully built and tested are **not wired** to engine features. This ADR defines the
contract, scope, and execution plan for wiring 9 high-value engine features to their
existing UI components.

### Problem Statement

Components exist and are production-ready, but the engine cannot use them because:

1. **No API contract** between engine feature (Go) and component (TS/Dart) — components
   define props, but engine does not emit the matching JSON/data shape.
2. **No SPA mode** — engine renders HTML server-side; components that need interactivity
   (form builder, workflow builder, pivot table) cannot function without client-side
   rendering.
3. **Backend gaps** — some features (pivot table, email/notification) have UI but no
   backend endpoint to feed data.
4. **No reference integration** — even where wiring is trivial (BcChatter, BcActivity,
   BcExport), there is no canonical pattern showing how engine JSON → component props.

### Scope — 9 Features

| # | Feature | UI Component | Backend | Effort | Status |
|---|---------|--------------|---------|--------|--------|
| 1 | Visual Form Builder | New (BcFormBuilder) | Partial (view JSON exists) | L | Belum |
| 2 | Visual Workflow Builder | New (BcWorkflowBuilder) | Partial (process JSON exists) | L | Belum |
| 3 | Visual Schema Builder | New (BcSchemaBuilder) | Partial (model JSON exists) | L | Belum |
| 4 | Multi-Step Wizard | `BcWizard` ✅ | — | M | Belum |
| 5 | Comment & Mention | `BcChatter` ✅ | Belum (no comment endpoint) | M | Component ada, belum wire |
| 6 | Activity Feed | `BcActivity` ✅ | Partial (audit_log exists) | M | Component ada, belum wire |
| 7 | Export Data (CSV/Excel) | `BcExport` ✅ | Partial (query builder exists) | M | Component ada, belum wire |
| 8 | Pivot Table | `BcDataTable` ✅ + new `BcPivot` | Belum (no pivot aggregation endpoint) | L | Component ada, backend belum |
| 9 | Email/Notification UI | New (BcNotificationCenter) | Partial (email infra exists, no notification store) | L | Belum |

## Decision

### D1: Three-tier wiring model

Each feature follows a 3-tier contract:

```
┌─────────────────────────────────────────────────────────┐
│ Tier 1: Engine API (Go)                                 │
│   - REST endpoint OR WebSocket action OR template data  │
│   - Emits JSON matching component contract (ADR-008)    │
└──────────────────────┬──────────────────────────────────┘
                       │ JSON / data-fetcher
┌──────────────────────▼──────────────────────────────────┐
│ Tier 2: BcRenderer integration (ADR-005)                │
│   - view.json `part` declaration                        │
│   - data-source binding (ADR-005 Section 3)             │
│   - event bridge (ADR-010) for client actions           │
└──────────────────────┬──────────────────────────────────┘
                       │ props
┌──────────────────────▼──────────────────────────────────┐
│ Tier 3: Component (Mitosis/Flutter)                     │
│   - Consumes props per contract (ADR-007)               │
│   - Emits events per event bridge (ADR-010)             │
└─────────────────────────────────────────────────────────┘
```

### D2: Feature-specific contracts

#### 5. Comment & Mention (BcChatter) — M

- **Engine:** New `comment` built-in model (auto-created per `comments: true` model flag).
  Fields: `id`, `model`, `record_id`, `author_id`, `body`, `mentions` (JSON array),
  `parent_id` (threading), `created_at`. CRUD at `/api/v1/{module}/comments`.
- **Mention:** `@username` parsed server-side → creates `notification` record (see #9).
- **Component:** `BcChatter` already has `messages`, `onSubmit`, `onMention` props.
  Wire to `GET/POST /api/v1/{module}/comments?record_id=X`.
- **BcRenderer:** `part: "chatter"` with `data-source: { model: "comment", filter: { record_id: "{{record.id}}" } }`.

#### 6. Activity Feed (BcActivity) — M

- **Engine:** New `/api/v1/activity` endpoint aggregating `audit_log` + domain events.
  Query params: `model`, `record_id`, `user`, `since`, `limit`.
- **Component:** `BcActivity` has `items`, `variant` props. Map audit_log entries →
  `BcActivityItem[]` via server-side transformer.
- **BcRenderer:** `part: "activity"` with `data-source: { endpoint: "/api/v1/activity", params: { record_id: "{{record.id}}" } }`.

#### 7. Export Data (BcExport) — M

- **Engine:** New `POST /api/v1/{module}/{model}/export` endpoint. Accepts:
  `{ format: "csv"|"xlsx"|"json", columns: [...], filter: {...}, sort: {...} }`.
  Uses existing query builder (ADR in features.md). Streams file response.
- **Component:** `BcExport` has `format`, `columns`, `onExport` props. Wire to export endpoint.
- **BcRenderer:** `part: "export"` with `data-source: { model: "{{model}}" }`.

#### 4. Multi-Step Wizard (BcWizard) — M

- **Engine:** New `wizard` view type in `view.json` — `type: "wizard"` with `steps[]`,
  each step is a form fragment. Process engine drives step transitions.
- **Component:** `BcWizard` already complete (16350 bytes, 16k). Wire `steps` from view JSON.
- **BcRenderer:** `part: "wizard"` declarative in view JSON.

#### 8. Pivot Table (BcPivot) — L

- **Engine:** New `POST /api/v1/{module}/{model}/pivot` endpoint. Accepts:
  `{ rows: [...], cols: [...], values: [{ field, agg: "sum"|"count"|"avg" }], filter: {...} }`.
  Returns aggregated 2D matrix.
- **Component:** New `BcPivot` component (extend `BcDataTable` with pivot mode).
- **BcRenderer:** `part: "pivot"` with pivot config in view JSON.

#### 1-3. Visual Builders (Form/Workflow/Schema) — L each

- **Pattern:** Shared `BcBuilder` base — drag-drop canvas + property panel + JSON preview.
- **Engine:** Endpoints to validate & save generated JSON:
  - `POST /api/v1/admin/views/validate` (form builder)
  - `POST /api/v1/admin/processes/validate` (workflow builder)
  - `POST /api/v1/admin/models/validate` (schema builder)
- **Component:** New `BcFormBuilder`, `BcWorkflowBuilder`, `BcSchemaBuilder`.
- **BcRenderer:** Admin-only parts, mounted in `/admin` panel (ADR-005 Section 4).

#### 9. Email/Notification UI (BcNotificationCenter) — L

- **Engine:** New `notification` built-in model. Fields: `id`, `user_id`, `type`,
  `title`, `body`, `data` (JSON), `read_at`, `created_at`. Endpoints:
  `GET /api/v1/notifications`, `POST /api/v1/notifications/{id}/read`,
  `POST /api/v1/notifications/read-all`.
- **WebSocket:** Push new notifications via existing WebSocket hub (already broadcasts
  domain events).
- **Component:** New `BcNotificationCenter` (dropdown panel) + `BcNotificationItem`.
- **BcRenderer:** `part: "notification-center"` in layout navbar.

### D3: SPA mode opt-in

Engine gains an opt-in SPA mode per view: `view.json` `{ "render": "spa" }` flag.
When set, engine returns JSON (view definition + initial data) instead of HTML;
client-side `BcRenderer` (ADR-005) renders the component tree. SSR remains default
for backwards compatibility.

This is **required** for visual builders (#1-3) and pivot table (#8) which need
client-side interactivity. Other features (#4-7, #9) work in both SSR + SPA.

### D4: Phased rollout

| Phase | Features | Rationale |
|-------|----------|-----------|
| **Phase 1: Wire existing** | #5 Comment, #6 Activity, #7 Export, #4 Wizard | Components ready, only API + BcRenderer part needed. Fastest value. |
| **Phase 2: SPA foundation** | SPA mode + BcRenderer client runtime | Prerequisite for visual builders. |
| **Phase 3: New components** | #8 Pivot, #9 Notification | New component + backend. |
| **Phase 4: Visual builders** | #1 Form, #2 Workflow, #3 Schema | Most complex, depends on SPA + BcRenderer maturity. |

## Scope Boundary — Parallel Execution

This ADR touches the following file ownership boundaries to enable parallel work
with ADR-013 (Cross-tier Infrastructure) and ADR-015 (Chart Strategy):

### ADR-014 OWNS (exclusive):

```
engine/internal/presentation/api/         — new endpoints (comments, activity, export, pivot, notifications)
engine/internal/presentation/view/         — new view types (wizard, pivot), SPA mode
engine/internal/presentation/template/     — template wiring for new parts
engine/internal/presentation/admin/        — admin panel for visual builders
engine/internal/domain/                    — new domain models (comment, notification) IF new
engine/internal/runtime/                   — pivot aggregation, export streaming
packages/ui/src/components/views/          — BcWizard wiring (if needed)
packages/ui/src/components/social/         — BcChatter, BcActivity wiring (IF prop changes needed)
packages/ui/src/components/print/          — BcExport wiring (IF prop changes needed)
packages/ui/src/components/overlays/       — BcNotificationCenter (new)
packages/ui-flutter/lib/src/components/    — bc_chatter.dart, bc_activity.dart, bc_export.dart (IF wiring)
```

### ADR-014 MUST NOT TOUCH (owned by others):

```
packages/ui/src/core/                      — ADR-013 owns (hooks, providers, adapt, animation)
packages/ui/src/components/charts/         — ADR-015 owns
packages/ui-flutter/lib/src/core/          — ADR-013 owns
packages/ui-flutter/lib/src/components/bc_chart.dart — ADR-015 owns
packages/components/src/components/charts/ — ADR-015 owns
```

### Shared files (merge coordinator required):

```
packages/ui/src/index.ts                   — barrel exports (coordinate with ADR-013, ADR-015)
docs/components-progress.md                — progress tracker
docs/adr/README.md                         — ADR index
engine/docs/features/                      — feature docs (per-feature, low conflict)
```

### Sequencing with ADR-013

ADR-014 Phase 1 (wire existing) can start **immediately** — it does not depend on
ADR-013. ADR-014 Phase 2 (SPA mode) touches `BcRenderer` which ADR-013 may also
touch for AppProvider adoption. **Rule: ADR-014 Phase 2 starts after ADR-013
core refaktor stabilizes** (estimated Week 2).

## Consequences

### Positive

- 9 high-value ERP features become usable (comments, activity, export, wizards,
  pivot, notifications, visual builders).
- Establishes canonical engine↔UI wiring pattern reusable for future features.
- SPA mode unlocks client-side rendering path (prerequisite for visual builders).
- Reuses 100% of existing components — no new component needed for #4-7.

### Negative

- Visual builders (#1-3) are L effort each — significant investment.
- SPA mode is a paradigm shift from current SSR-only — risk of dual-render bugs.
- Pivot backend is non-trivial (OLAP aggregation in Go).
- New built-in models (comment, notification) add to base module surface.

### Neutral

- BcRenderer (ADR-005) must mature to client-side runtime — currently SSR-only.

## Acceptance Criteria

### Phase 1 (Wire Existing)
- [ ] `comment` model auto-created for models with `comments: true`
- [ ] `GET/POST /api/v1/{module}/comments` working with `BcChatter`
- [ ] `GET /api/v1/activity` working with `BcActivity`
- [ ] `POST /api/v1/{module}/{model}/export` working with `BcExport` (CSV + XLSX)
- [ ] `wizard` view type working with `BcWizard`
- [ ] At least 1 engine test per feature (handler + integration)

### Phase 2 (SPA Foundation)
- [x] `view.json` `{ "render": "spa" }` flag respected — engine parser accepts `Render` field on `ViewDefinition` with `RenderSSR`/`RenderSPA` constants
- [x] BcRenderer client runtime renders component tree from JSON — `packages/ui/src/spa-runtime.ts` (web/Mitosis) + `packages/ui-flutter/lib/src/spa_runtime.dart` (Flutter) fetch envelope, map to `UiSchema`, build `BcRendererProps` + `ReactiveStore`
- [x] SSR mode unchanged (backwards compatible) — `view.IsSPAMode` predicate returns true only for `"spa"`; views with `render:"ssr"` or no `render` continue to return HTML
- [x] Engine SPA envelope = `{ view, data, meta }` (data-parity with SSR enforced by `TestDataParity_SPAEqualsSSR_FormView`)
- [x] Test coverage: 11 engine tests, 13 web tests, 28 Flutter tests, 100% passing

### Phase 3 (New Components)
- [ ] `POST /api/v1/{module}/{model}/pivot` endpoint
- [ ] `BcPivot` component renders aggregated matrix
- [ ] `notification` model + endpoints + WebSocket push
- [ ] `BcNotificationCenter` in navbar

### Phase 4 (Visual Builders)
- [ ] `BcFormBuilder` drag-drop → valid `view.json`
- [ ] `BcWorkflowBuilder` drag-drop → valid `process.json`
- [ ] `BcSchemaBuilder` drag-drop → valid `model.json`
- [ ] All 3 validate via admin endpoints before save

## References

- ADR-005: Metadata-Driven UI (BcRenderer) — `part` declaration, data-source binding
- ADR-006: Universal Component Strategy — Pure/Adapt pattern for new components
- ADR-007: Enterprise Universal UI API Contract — component prop contracts
- ADR-008: Contract-First API Freeze — engine↔component contract freeze
- ADR-010: Web Component Event Bridge — event emission from components
- `docs/features.md` — feature inventory & roadmap
- `engine/docs/features/` — per-feature engine docs
