# Bucket: engine

> Go engine + JSON definitions + runtime + modules. Lightweight session memory.
> Source of truth: AGENTS.md, docs/architecture.md, docs/features.md, engine/docs/features/

## Current State (2026-06-30)

### Test baseline
- 570 Go engine tests + 82 go-json tests, 0 failures. Build OK.
- Run: `cd engine && rtk go test ./... -count=1` (rtk summarizes). Full output: `| Out-String`.

### Recent completions (engine-side)
- Data Migration System (Laravel-style, 29 tests, MongoDB parity via MigrationStore)
- Group-Based Permission System (Odoo-style 12 perms, bi-directional JSON↔DB sync)
- Convention-Driven CRUD (`"api": true` auto-REST, `bc-datatable`, `bitcode publish:crud`)
- Auto Swagger/OpenAPI + GraphQL API + WebSocket CRUD (per-model `api.protocols.*`)
- Model lifecycle events (16 types) + Field validation (40+ validators) + Field sanitization (14)
- Client-side onchange API (`POST /api/{model}/onchange`)
- Model process registry (dynamic finders, OQL, module-qualified names)

### Open / Remaining (engine)
- [ ] Redis cache wiring — wire into permission checker + query result cache
- [ ] Marketplace — community module sharing
- [ ] NATS event bus — replace in-process bus for distributed deployments
- [ ] ADR-014 (Proposed) — Engine UI Feature Wiring: 9 ERP UI features → component-backed SSR render. Extends ADR-005 BcRenderer. Parallel-safe with ADR-013/015.

### Conventions (engine)
- All PK/FK = UUID (TEXT in SQLite, UUID in Postgres, CHAR(36) in MySQL)
- DDD internally, flat JSON externally. Convention over configuration.
- `engine/internal/compiler/parser/` one file per definition type.
- Infrastructure tables (audit_log, data_revisions, ir_migration, _sync_*) are Go-generated DDL, NOT JSON models — do not appear in admin panel, no auto-CRUD.

### ADR-014 Phase 1 progress (2026-07-01)
- [x] **#6 Activity Feed** — `GET /api/v1/activity` endpoint. `runtime/activity.Service` transforms `audit_logs` → `BcActivityEntry[]` (verb/variant/icon/actor). Permission-aware (record-scoped: read perm + record rule; global timeline: auth-only). Pagination (page/page_size, max 200). `ViewActivity` wired in `view/renderer.go` + `component_compiler.CompileActivity`. i18n keys (11 locales). Engine tests: 12 service + 10 handler + 1 compiler. BcActivity.lite.tsx contract upgraded (dataSource, variant, loading, onLoadMore, onItemAction) + 15 source tests. Flutter bc_activity.dart parity (variant keys, loading skeleton) + 14 widget tests.
- [x] **#5 Comment & Mention** (sibling, separate worker) — comment.json model + endpoint + mention→notification.
- [x] **#7 Export Data** — `POST /api/v1/{module}/{model}/export` (+ GET variant). `runtime/export.Service` renders CSV (encoding/csv + UTF-8 BOM), XLSX (excelize streaming writer), JSON (pretty array) in batches (default 1000). **Field-level permission**: columns sanitized per field `groups` (dropped if user not in group), masked fields masked when user lacks `can_mask` (reuses `field_filter.go` `maskValue`). Export permission checked per group (12-perm `export`). Filter+sort via `persistence.Query` + OQL. Streaming response (CSV/JSON via `SetBodyStreamWriter`, XLSX buffered). Engine tests: 21 service + 12 handler (csv, xlsx, json, permission denial, filter, columns, field mask, group sanitization, superuser bypass, GET variant). BcExport.lite.tsx contract upgraded (endpoint, columns, loading, progress, onExportComplete/onProgress, declarative fetch) + 13 source tests + 2 vitest. Flutter bc_export.dart parity (endpoint, columns, loading indicator, onExportComplete/onProgress) + 8 widget tests. i18n keys (11 locales). Wired in `app.go`. List views emit `<bc-export>` via component_compiler.
- [x] **#42 Pivot Table** (Phase 3) — `POST /api/v1/{module}/{model}/pivot` endpoint. `runtime/pivot.Service` pure Go OLAP aggregation (GROUP BY rows+cols, sum/count/avg/min/max), multiple value fields, multiple row/col dimensions, null handling, grand totals. Permission-aware (read perm + field-level sanitization + RLS filter). `ViewPivot` wired in `view/renderer.go` + `component_compiler.CompilePivot`. i18n keys (11 locales). Engine tests: 20 service + 9 handler. BcPivot.lite.tsx + 14 source tests. Flutter bc_pivot.dart + 12 widget tests. Block Catalog `analytics.pivot`.

### Notes
- Someone may be working on `sprints/bitcode-flexible-module` and `sprints/package-ui-elementor` — confirm before touching sprint-scoped files.
