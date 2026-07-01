# Bucket: flutter

> Flutter adapter, token bridge, demo parity, platform widgets.
> Scope: `packages/ui-flutter/` + Flutter-side of `packages/ui` contracts.

## 2026-07-01 — BcNotificationCenter Flutter widget (Phase 3 Notification)

- **bc_notification_center.dart** — `packages/ui-flutter/lib/src/components/`. StatefulWidget. Adaptive rendering via PlatformWidget: iOS Cupertino dropdown (CupertinoButton bell + CupertinoAlertDialog panel), Android Material dropdown (IconButton bell + Dialog panel). Bell button with badge (99+ overflow), notification items list, mark-all-read, infinite scroll, loading/error/empty states.
- **bc_notification_center_test.dart** — 10 tests: fromJson parses all fields, defaults type to info, isUnread false when read_at set, copyWith marks as read, equality by id, bell renders with badge when unread > 0, no badge when 0, 99+ overflow, onMarkRead callback when item tapped, disabled bell doesn't respond.
- **dart analyze:** clean on both files.
- **flutter test:** 10/10 pass.
- **Gotcha:** Row overflow in dropdown header — wrapped title in `Expanded` + `TextOverflow.ellipsis`. @worker

## 2026-07-01 — BcPivot Flutter widget (Phase 3 Pivot Table)

- **bc_pivot.dart** — `packages/ui-flutter/lib/src/components/`. StatelessWidget with StatefulWidget for expand/collapse. Material `DataTable` on Android, Cupertino-styled container on iOS (per ADR-011). Expandable rows, row totals, grand total, no-data state.
- **bc_pivot_test.dart** — 12 widget tests: renders table with data (VAL-PIVOT-001), no-data states (VAL-PIVOT-002/003), grand total show/hide, title, expandable/collapsible, JSON parsing, serialization, agg func mapping, isEmpty.
- **dart analyze:** 0 errors, 0 warnings on pivot files (info-level `prefer_const_constructors` consistent with rest of codebase).
- **flutter test:** 12/12 pass.

## 2026-07-01 — Flutter block catalog registry + 0-lint cleanup (Droid)

- **What:** Created Flutter equivalent of TS `ready-catalog.ts`. Fixed ALL 13 Dart info-level lints.
- **ready_catalog.dart:** `registerReadyCatalog()`, `isReadyEntryReady()`, `getReadyEntryIds()`, `readyCatalogEntries` const, `readyCount` getter. Exported from `bc_ui_flutter.dart`.
- **ready_catalog_test.dart:** 9 tests (registration, idempotency, lookup, unknown id). 34/34 total Flutter tests pass.
- **Lint fixes:** 13 `curly_braces_in_flow_control_structures` lints fixed in effect.dart (7), runner.dart (5), source.dart (1). `dart analyze` now reports **0 issues** on all touched files.
- **Commit:** `93d9550`

---
## Active Threads

### 2026-06-30 — Flutter demo parity in progress
- **Status:** in progress. React demo has 127 components in 14 groups; Flutter demo has ~25 across 6 sections; library has 120+ ported components but most not in demo.
- **Goal:** 1:1 parity between `packages/ui-flutter/example` and `packages/ui/demo-app`.
- **Missing groups entirely:** Input (Extended), Primitives, Feedback, Media, Charts, Print, Views, Social.
- **Source handoff:** `handoffs/2026-06-28-flutter-demo-parity.md`.
- **Next:** port the 0-existing groups first (Charts blocked on ADR-015 DI loader pattern).

### 2026-06-30 — Stub completion sweep (ADR-016 Phase 4)
- **Status:** mostly complete (28 stubs ported across 5 clusters as of 2026-06-29). Remaining: review of Flutter-side parity for newly ported components.
- **Source handoff:** `handoffs/2026-06-29-ui-stub-completion-and-flutter-review.md`.

## Key Decisions (sticky)

- **Flutter Adaptive (ADR-011):** iOS/macOS → Cupertino via `flutter_platform_widgets` (`Platform.*`); Android/Web/Windows/Linux → Material 3. DataTable/Splitter/Accordion stay Material-only (no Cupertino counterpart). The `flutter:` adapter mapping string in `.contract.ts` is Material-canonical + descriptive, NOT binding codegen target.
- **Token Bridge:** `generate-flutter-tokens.mjs` parses `theme.css` HSL triplets → emits `bc_theme.dart` with `static const Color` hex + 4 factories: `materialLight/Dark`, `cupertinoLight/Dark`.
- **Chart strategy (ADR-015):** `flutter_echarts` replaces 912-line CustomPaint. 26/26 chart types via shared ECharts option contract. CustomPaint removed entirely (no deprecation). Charts need DI loader pattern (see demo-app bucket).
- **BcRenderer hybrid port (ADR-011):** `BcForm` enhancement + standalone field/data loaders + compositor. Hybrid reactive store. Complete as of 2026-06-24.

## Gotchas

- `flutter_platform_widgets` must be used for any widget that has a Cupertino counterpart — never hand-roll `if (Platform.isIOS)`.
- Token regeneration is a build step: changing `theme.css` requires re-running `generate-flutter-tokens.mjs` or colors drift.
- ECharts DI loader: charts won't render in Flutter demo without the dependency-injection loader wired (ADR-016 architectural fix #5).

## Files Touched Recently
- `packages/ui-flutter/example/` — demo app sections
- `packages/ui-flutter/lib/` — library components
- `packages/ui-flutter/tool/generate-flutter-tokens.mjs` — token bridge
