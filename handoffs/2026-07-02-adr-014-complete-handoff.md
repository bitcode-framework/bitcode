# ADR-014 Complete — Final Handoff

**Date:** 2026-07-02
**Phase:** Phase 4 (Final) — Documentation + Handoff
**Status:** ✅ ALL 9 FEATURES COMPLETE

---

## What Was Done

ADR-014 (Engine ↔ UI Feature Wiring) is fully implemented across all 4 phases. All 9 features have engine APIs, web components (Mitosis), Flutter parity, i18n (11 locales), permission enforcement, and tests.

### Feature Summary

| # | Feature | Phase | Engine | Web | Flutter | Tests |
|---|---------|-------|--------|-----|---------|-------|
| 5 | Comment & Mention | 1 | ✅ comment.json + CommentService + @mention parser + WS push | ✅ BcChatter.lite.tsx (upgraded) | ✅ bc_chatter.dart | ✅ 13 engine + 5 handler |
| 6 | Activity Feed | 1 | ✅ ActivityService + GET /api/v1/activity | ✅ BcActivity.lite.tsx (upgraded) | ✅ bc_activity.dart | ✅ 9 engine + 10 handler |
| 7 | Export Data | 1 | ✅ ExportService + POST /export (CSV/XLSX/JSON streaming) | ✅ BcExport.lite.tsx (upgraded) | ✅ bc_export.dart | ✅ 14 engine + 12 handler |
| 4 | Multi-Step Wizard | 1 | ✅ wizard view type + process transitions | ✅ BcWizard.lite.tsx (upgraded) | ✅ bc_wizard.dart | ✅ engine + component |
| 2 | SPA Foundation | 2 | ✅ render:"spa" + JSON envelope | ✅ spa-runtime.ts | ✅ spa_runtime.dart | ✅ engine + 13 source + 28 Flutter |
| 8 | Pivot Table | 3 | ✅ PivotService + POST /pivot (OLAP aggregation) | ✅ BcPivot.lite.tsx (new) | ✅ bc_pivot.dart | ✅ 20 engine + 9 handler |
| 9 | Notification | 3 | ✅ notification.json + NotificationService + WS push | ✅ BcNotificationCenter.lite.tsx (new) | ✅ bc_notification_center.dart | ✅ 8 engine + 7 handler + 21 component |
| 1 | Form Builder | 4 | ✅ ViewValidator + POST /admin/views/validate | ✅ BcFormBuilder.lite.tsx (new) | ✅ bc_form_builder.dart | ✅ 16 engine + 5 handler + component |
| 2 | Workflow Builder | 4 | ✅ ProcessValidator + POST /admin/processes/validate | ✅ BcWorkflowBuilder.lite.tsx (new) | ✅ bc_workflow_builder.dart | ✅ engine + component |
| 3 | Schema Builder | 4 | ✅ ModelValidator + POST /admin/models/validate | ✅ BcSchemaBuilder.lite.tsx (new) | ✅ bc_schema_builder.dart | ✅ engine + component |

### Phase 4 Docs + Handoff (This Feature)

Updated all documentation per AGENTS.md mandatory doc-update rules:

- **docs/codebase.md** — Added all new ADR-014 files: runtime services (comment, activity, export, notification, builder validators), API handlers (comment, activity, export, pivot, notification, builder), SPA renderer, new UI components (BcPivot, builders), Flutter components (bc_chatter through bc_schema_builder), block catalog entries.
- **docs/architecture.md** — Added new "ADR-014: Engine ↔ UI Feature Wiring" section with 3-tier architecture diagram, SPA mode explanation, new JSON models (comment, notification), API endpoint table, new view types, new UI component table.
- **docs/features.md** — Marked all 9 ADR-014 features ✅: #1 Schema Builder, #22 Workflow Builder, #26 Email/Notification, #30 Form Builder, #33 Wizard, #41 Export, #42 Pivot, #63 Comment, #64 Activity, #67 In-App Notification. Updated roadmap checklist items.
- **docs/components-progress.md** — Added "ADR-014 Engine ↔ UI Feature Wiring" section with complete component status table (all 9 features, all 3 tiers).
- **docs/adr/README.md** — Updated ADR-014 status from Proposed → **Accepted** (2026-07-02).
- **docs/adr/014-engine-ui-feature-wiring.md** — Updated status from Proposed → **Accepted**.

### Per-feature docs in engine/docs/features/ (all verified present):
- comments.md, activity.md, export.md, wizard.md, spa-mode.md, pivot.md, notification.md, form-builder.md, workflow-builder.md, schema-builder.md

---

## Verification

### Test Suites Run

| Suite | Result |
|-------|--------|
| Engine Go tests (`rtk go test ./... -count=1`) | ✅ 1077 passed, 51 packages, 0 failures |
| ui-schema tests (`node scripts/run-tests.mjs`) | ✅ 0 failures |
| packages/ui source + vitest (`npm test`) | ✅ Passing |
| Flutter tests (`flutter test`) | ✅ Passing |
| Engine build (`rtk go build ./...`) | ✅ OK |

### Assertion Compliance

- **P1-I18N** (A-S6-001): ✅ All 11 locale files (en, id, ar, de, es, fr, ja, ko, pt-BR, ru, zh-CN) in `engine/embedded/modules/base/i18n/` contain ADR-014 keys (export.*, pivot.*, notification.*, comment.*, activity.*, builder.*, wizard.*).
- **P1-PERM** (A-S6-002): ✅ Export handler checks `CanExport` group permission (403 on denial). Pivot handler checks `CanRead`/`CanSelect`. Comments/notifications/activity honor group-based permissions. Builders are admin-only (validate endpoints under `/api/v1/admin/`).

---

## Key Decisions

1. **comment & notification = JSON models** (not Go DDL) — Gets auto-CRUD, permissions, i18n, record rules, WebSocket events for free. Go service layer handles specialized logic (mention parsing, WS push, activity aggregation).
2. **SPA mode = JSON envelope** — Engine returns `{view, data, meta}` when `render:"spa"`. Existing portable BcRenderer renders client-side. SSR stays default.
3. **New components = Pure/Adapt + Block Catalog** — BcPivot, BcNotificationCenter, 3 builders follow ADR-006 2-tier pattern and register as Block Catalog candidates.
4. **Builder validators = admin-only endpoints** — `POST /api/v1/admin/{views,processes,models}/validate` validate JSON definitions server-side.

---

## Files Changed (This Feature — Phase 4 Docs)

- `docs/codebase.md` — Added all ADR-014 files to file map
- `docs/architecture.md` — Added ADR-014 section (3-tier architecture, SPA mode, new models, API table)
- `docs/features.md` — Updated 9 feature statuses to ✅, fixed roadmap checklist
- `docs/components-progress.md` — Added ADR-014 component status table
- `docs/adr/README.md` — ADR-014 status → Accepted
- `docs/adr/014-engine-ui-feature-wiring.md` — Status → Accepted
- `handoffs/2026-07-02-adr-014-complete-handoff.md` — This file

---

## Next Steps

ADR-014 is **complete**. All 9 features are production-ready with engine APIs, web + Flutter UI, i18n, permissions, and tests. Potential follow-up work:

- **Block Catalog Ready promotion** — All new components are registered as `candidate` entries. Ready-gate promotion (ADR-017/018) is separate scope.
- **SPA mode adoption** — Visual builders (Phase 4) require SPA mode. More views can opt into SPA as needed.
- **Email fan-out** — NotificationService creates notification records. Actual email delivery via SMTP (pkg/email) could be wired as a notification channel.
