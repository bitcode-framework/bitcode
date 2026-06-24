# Next Missions -- Parallel Execution Plan

**Created:** 2026-06-24
**Context:** ADR-011 Phase 5c just completed (130 Flutter widgets, 1158 tests, 297/297 assertions). This handoff defines the next two parallel missions.

---

## Mission A: BcRenderer Dart Port (Flutter package only)

### What
Port the schema-driven UI renderer from TypeScript (`packages/ui/src/core/renderer/`) to Dart (`packages/ui-flutter/lib/src/renderer/`). This enables the Go engine to auto-generate Flutter UIs from the same `models/*.json` definitions the web uses.

### Why
Without BcRenderer, Flutter widgets exist but need manual composition. With it, the "JSON is the source code" vision works on mobile/desktop natively.

### Scope
1. Port `UiSchema` document types (already partially done in `schema/ui_schema.dart` + `schema/ui_schema_types.dart`)
2. Port `resolveFieldType()` mapper (input/textarea/select/number/date/toggle/checkbox/email/password/url/tel)
3. Port `resolveField()`, `resolveFields()`, `defaultFieldResolver()` pure functions
4. Port BcRenderer widget with 6 schema kinds: form, detail, list, section, grid, custom
5. Wire ReactiveStore into BcRenderer for reactive forms (Level 1 per ADR-004)
6. Tests: field resolver, per-kind rendering, section/grid recursion, custom sentinel
7. Integration gate: BcRenderer renders a full form from a UiSchema document

### Key Files to Read First
- `packages/ui/src/core/renderer/` -- TypeScript source to port
- `packages/ui-flutter/lib/src/schema/` -- existing Dart UiSchema types
- `packages/ui-flutter/lib/src/core/reactive_store.dart` -- Dart reactive engine
- `docs/adr/005-metadata-driven-ui.md` -- BcRenderer architecture
- `docs/adr/011-flutter-adapter-architecture.md` -- Flutter adapter design

### Estimated Effort
1-2 mission sessions (medium complexity, mostly porting + testing)

### Dependencies
None -- can start immediately. Does NOT depend on Mission B.

---

## Mission B: Phase 2 Core Gaps (Engine only)

### What
Implement 12 high-priority engine features that make the platform usable for real ERP scenarios.

### Why
These are the most requested features blocking real-world usage. Each is independent and can be tackled as separate features within one or more missions.

### Features (from docs/features.md Phase 2)

| # | Feature | Effort | Description |
|---|---------|--------|-------------|
| 5 | Computed Field Evaluator | M | Expression parser + evaluator at query time (e.g., `sum(lines.subtotal)`) |
| 11 | Field-Level Permission | M | Per-field permission config in model JSON + API/form filter |
| 17 | Record Activity Timeline | M | API endpoint + timeline UI component |
| 27 | Assignment Rules | M | Assignment rule JSON definition + evaluator on create/update |
| 28 | Webhook System | M | Webhook definition (URL, events, headers) + event bus dispatcher |
| 33 | Multi-Step Form / Wizard | M | Wire `bc-dialog-wizard` component + wizard JSON definition |
| 35 | Web Form (Public) | M | Public form route (bypass auth) + CAPTCHA + rate limiting |
| 41 | Export Data (CSV/Excel) | M | Server-side export handler per model |
| 57 | 2FA (TOTP) | M | TOTP library + setup flow + verification middleware |
| 59 | Backup & Restore | M | DB dump/restore commands + optional scheduling |
| 63 | Comment & Mention | M | Comment model + API + wire `bc-chatter` component |
| 67 | In-App Notification | M | Notification model + preferences + bell UI + WebSocket delivery |

### Key Files to Read First
- `docs/features.md` -- full feature inventory with status
- `docs/architecture.md` -- system design
- `engine/internal/` -- Go engine source
- `engine/docs/features/` -- per-feature docs

### Estimated Effort
2-3 mission sessions (12 features, each M-sized, can batch 3-4 per mission)

### Dependencies
None -- can start immediately. Does NOT depend on Mission A.

### Suggested Batching
- **Batch 1:** Export Data (#41), Backup & Restore (#59), 2FA TOTP (#57) -- infrastructure features, no UI dependencies
- **Batch 2:** Webhook System (#28), Assignment Rules (#27), Computed Field Evaluator (#5) -- engine logic features
- **Batch 3:** Comment & Mention (#63), In-App Notification (#67), Record Activity Timeline (#17) -- collaboration features (depend on WebSocket)
- **Batch 4:** Field-Level Permission (#11), Multi-Step Form (#33), Web Form Public (#35) -- UI + auth features

---

## After Both Complete: Mission C -- Offline Mode Phase 3

### What
Sync engine: device registration, outbox recording, sync push/pull, envelope grouping.

### Why
Depends on both: Flutter UI layer (Mission A) for client-side rendering, and engine features (Mission B) for webhook/notification infrastructure.

### Key Files
- `docs/plans/2026-04-28-offline-mode-design.md`
- `docs/plans/impl/offline-mode-implementation.md`
- `engine/docs/features/offline-mode.md`
- `packages/tauri/src-tauri/src/main.rs` -- Rust SQLite migrations

---

## Execution Notes

- **Parallelism:** Missions A and B can run in separate Droid sessions simultaneously. They touch different packages with no code overlap.
- **Sequential within each:** Mission B batches should run in order (Batch 1 first, etc.) since later batches may depend on earlier infrastructure.
- **After Mission A:** The Flutter BcRenderer enables a future mission to wire the Go engine's view system to Flutter via UiSchema JSON, completing the full-stack mobile story.
- **After Mission B:** The platform becomes genuinely usable for ERP deployments. Phase 3 (Visual Builders) and Phase 4 (Enterprise Features) follow.
- **Test commands:**
  - Flutter: `cd packages/ui-flutter && flutter test`
  - Engine: `cd engine && rtk go test ./... -count-1`
  - go-json: `cd packages/go-json && go test ./... -v`
