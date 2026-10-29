# Bucket: engine

> Go engine + JSON definitions + runtime + modules. Lightweight session memory.
> Source of truth: AGENTS.md, docs/architecture.md, docs/features.md, engine/docs/features/

## Current State (2026-07-13)

## 2026-10-29 — go-json examples, CLI fixes, cross-platform installers
- **What:** Extracted every README snippet into `packages/go-json/examples/` (7 CLI programs + 3 Go embedding programs + `examples_test.go` that executes them), fixed 8 README snippets that could not run as written, fixed two real CLI bugs, and shipped installers.
- **CLI bugs found by running the docs:**
  1. `WithLimits` replaced the whole `Limits` struct, zeroing `MaxSteps` → every program failed with "step limit (0) exceeded". Now merges only non-zero fields.
  2. Flags after the program path were silently dropped (Go's `flag` stops at the first positional). New `cmd/cli/args.go` `reorderArgs`/`parseFlags`, applied to all 10 parse sites.
- **Distribution decision:** two binaries are unavoidable — `go-json-runtimes` depends on `go-json`, and core keeps zero deps, so core cannot import the runtimes without a module cycle. Renamed core CLI to `cmd/go-json-lite`; the full binary from runtimes owns the `go-json` name.
- **`go install` blocker:** Go refuses to install a module whose `go.mod` has a `replace` directive. Removed it from `go-json-runtimes/go.mod`; `packages/go.work` holds the override instead. Note: a module in `use` cannot be replaced at all versions — the replace must be version-specific (`replace ... v0.1.0 => ./go-json`).
- **Installers:** `packages/go-json-runtimes/install.sh` (Linux/macOS/Windows-via-Git-Bash) and `install.ps1` (Windows). Detect OS/arch, download the release asset, verify SHA-256, install to a PATH dir, fall back to `go install`. Both verified end to end including the checksum-mismatch guard.
- **Release:** `.github/workflows/release.yml` builds linux/darwin/windows × amd64/arm64 with `CGO_ENABLED=0` (cross-compiles cleanly), emits `.sha256` sidecars, publishes on tag push.
- **Files touched:** `packages/go-json/{examples,cmd/cli,runtime,README.md,docs/cli-reference.md,AGENTS.md}`, `packages/go-json-runtimes/{go.mod,README.md,install.sh,install.ps1,.github/workflows/release.yml}`, `packages/go.work`, `docs/codebase.md`, root `AGENTS.md`
- **Gotcha / note:** published. `go-json` tagged v0.1.0–v0.1.2, `go-json-runtimes` tagged v0.1.0–v0.1.4; latest release is **v0.1.4** with 12 assets. Use `go install .../cmd/go-json@v0.1.4` (explicit version) — the Go proxy's `@latest` endpoint is a stale cache pinned to the poisoned v0.1.0 and will keep failing the sumdb check until it refreshes. The install scripts are unaffected because they read the GitHub releases API, which correctly reports v0.1.4.
- **Two traps hit while publishing (both cost a release cycle):**
  1. `go.work` masks a missing `go.sum` entry. Bumping the go-json dependency left `go.sum` stale; local builds passed, standalone builds (CI + consumers) failed with "missing go.sum entry". Always `GOWORK=off go mod tidy` before tagging. A CI step now enforces this.
  2. **Never re-point a pushed tag.** sum.golang.org is append-only, so a re-tagged version permanently fails with a checksum mismatch for everyone. A poisoned version can only be abandoned by releasing a higher number.
- **Version reporting:** `Version` is a fallback const; `buildVersion` is set by release ldflags; `ResolveVersion()` prefers ldflags, then `debug.ReadBuildInfo()` module version, then the const.
- **Tests:** go-json `go test ./...` all pass, `go build`/`go vet` clean; runtimes `go test ./...` all pass standalone (`GOWORK=off`).

## 2026-07-13 — Scripts runtime resolution inventory
- **What:** Checked existing script/runtime resolution vs ADR-020 layered render engines.
- **Decision / outcome captured in ADR-021:** scripts are on-demand assets; runtime order is step → module default → extension → family → concrete engine; do not collapse with module layer or UI render resolver.
- **Files checked:** `engine/internal/runtime/executor/steps/script.go`, `engine/internal/runtime/embedded/{registry,script_runner,script_loader}.go`, `engine/internal/runtime/plugin/manager.go`, `engine/docs/features/plugins.md`, `docs/adr/020-pluggable-render-engine-architecture.md`
- **Gotcha / note:** `.js` path currently ambiguous between embedded/default docs and plugin detect defaults; ADR chooses embedded javascript unless explicit `runtime:"node"`.

## 2026-07-13 — Process vs workflow inventory for ADR-021 discussion
- **What:** Checked current process/workflow code and sample JSON for n8n-format decision.
- **Decision / outcome:** Keep BitCode-native format. Existing process JSON supports linear `steps[]` and DAG `nodes{}` + `edges[]`. Workflow JSON is separate state machine (`states` + `transitions`, model-bound). Not n8n-compatible as source of truth; n8n may only be optional import/export later.
- **Files checked:** `engine/internal/compiler/parser/process.go`, `engine/internal/compiler/parser/workflow.go`, `samples/erp/modules/crm/processes/*`, `engine/modules/sales/processes/confirm_order.json`, `engine/docs/features/processes.md`, `engine/docs/features/workflows.md`
- **Gotcha / note:** Naming is easy to confuse: UI "workflow builder" targets process graph authoring, while engine `workflows` means model status state-machine.

## Current State (2026-07-01)

### Test baseline
- 1025 Go engine tests + 82 go-json tests, 0 failures. Build OK.
- Run: `cd engine && rtk go test ./... -count=1` (rtk summarizes). Full output: `| Out-String`.

### 2026-07-01 — Form Builder Engine (ADR-014 Phase 4a)
- **What:** POST /api/v1/admin/views/validate endpoint + ViewValidator service
- **Decision / outcome:** ViewValidator in runtime/builder checks view.json structure + field refs against model registry; BuilderHandler in presentation/api wired in app.go under /api/v1/admin with JWT auth
- **Files touched:** `internal/runtime/builder/view_validator.go`, `internal/runtime/builder/view_validator_test.go`, `internal/presentation/api/builder_handler.go`, `internal/presentation/api/builder_handler_test.go`, `internal/app.go`, `embedded/modules/base/i18n/*.json` (41 builder keys in 11 locales), `engine/docs/features/form-builder.md`
- **Tests:** 16 validator unit tests + 5 handler tests, all passing

### 2026-07-01 — Notification system (ADR-014 Phase 3)
- **What:** Full notification feature — model, service, REST endpoints, WS push, mention integration
- **Decision / outcome:** Dedicated NotificationService (not convention CRUD) for user-scope, unread-first, and WS push via `BroadcastToUser` on hub
- **Files touched:** `embedded/modules/base/models/notification.json`, `internal/runtime/notification/service.go`, `internal/presentation/api/notification_handler.go`, `internal/presentation/websocket/hub.go` (+BroadcastToUser), `internal/runtime/comment/service.go` (mention→notification), `internal/app.go`
- **Tests:** 16 service + 7 handler = 23 new engine tests. Full suite 1004 pass.
- **Gotcha:** `MarkRead` with GORM `Updates` map on filtered row was unreliable on SQLite — switched to single `Update("read_at", now)`. @worker

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
