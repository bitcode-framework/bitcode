# Bucket: flutter

> Flutter adapter, token bridge, demo parity, platform widgets.
> Scope: `packages/ui-flutter/` + Flutter-side of `packages/ui` contracts.

## 2026-10-29 — Fix: Media + Views sections blank (unbounded height)

- **Symptom:** Clicking **Media** or **Views** in the Flutter demo left the content panel completely blank. Browser console: 53 errors (Media) / 35 errors (Views).
- **Root cause:** `Expanded` requires a bounded height, but the demo wraps every section in `SingleChildScrollView` (unbounded height). Four widgets violated this:
  - `bc_pivot.dart` — `Column > Expanded` in both `_buildMaterial` and `_buildCupertino`.
  - `bc_viewer_document.dart` — `Column > Expanded` in `build`.
  - `bc_viewer_instagram.dart` / `bc_viewer_tiktok.dart` — `Container(height: double.infinity)`, which becomes an infinite constraint under an unbounded parent.
- **Contract divergence:** the web adapters never had this problem. `BcPivot` is `w-full overflow-auto` (natural height, no `Expanded`); `BcViewerDocument`, `BcViewerInstagram`, and `BcViewerTiktok` all declare `width: '100%'` / `height: '600px'` defaults. The Flutter ports dropped those props and replaced them with `Expanded` / `double.infinity`.
- **Fix:** each widget now uses `LayoutBuilder` and branches on `constraints.hasBoundedHeight` — fill the host when bounded, fall back to the web contract's 600px viewport when not. `bc_viewer_youtube` already used `AspectRatio` and needed no change.
- **Tests:** `test/unbounded_height_test.dart` (13) — pumps each widget inside `SingleChildScrollView` and asserts no layout exception, plus a bounded-box case to prove fill behaviour is preserved. All 13 failed before the fix.
- **Verification:** `flutter test` 2144 pass / 0 fail. Live browser: Media 0 errors, Views 0 errors, `BcPivot` renders the full matrix with Grand Total 2480.
- **Still failing (separate, pre-existing):** `Display` overflows horizontally by 119px at narrow widths (demo layout, not a crash); `Native` fails only on `pumpAndSettle timed out` (endless animation) while rendering cleanly in the browser with 0 console errors. Both are demo/test-harness issues, not library regressions.

## 2026-10-29 — Fix: Block Catalog red blocks + overflow (JSON coercion + grid layout)

- **Symptom:** Block Catalog (Preview) demo section showed a red error block (`TypeError: "2026-06-18": type 'String' is not a subtype of type 'DateTime?'`) and yellow/black `RenderFlex overflowed by 62 px` bars.
- **Root cause A — hard casts at the renderer boundary.** `bc_field_builder.dart` used unchecked `as DateTime?` / `as num?` / `as bool?` / `as String?` on `store.get()`, which returns `dynamic`. Demo values come from JSON (`'orderDate': '2026-06-18'`), so the cast throws and tears down the subtree. Six call sites affected (input, number, date, select, toggle, checkbox). The web adapter's `BcDateField.lite.tsx` declares `value?: string` (ISO), so the Flutter adapter diverged from the shared contract.
- **Root cause B — fixed-height grid cells.** `_buildGrid` used `GridView` + `childAspectRatio: 3.0`; a `detail` with 3-4 fields is taller than the allotted cell. The 100062 px overflow seen in one screenshot had the *same* cause (same grid, different available height) — not a separate bug.
- **Fix A:** added `coerceString` / `coerceNum` / `coerceDateTime` / `coerceBool` at the top of `bc_field_builder.dart` and routed all six builders through them. Accepts native types, ISO strings, epoch millis, `1`/`0`, `'true'`/`'false'`/`'yes'`/`'on'`. Unconvertible input degrades to `null` instead of throwing. Store value is never mutated.
- **Fix B:** replaced `GridView` with `_buildGridRows` — `Row` + `Expanded` cells, `CrossAxisAlignment.start`, 12 px spacing, trailing `SizedBox.shrink()` placeholders to keep column widths aligned. Cells now size to content.
- **Tests added:** `test/renderer/bc_field_coercion_test.dart` (18), `test/renderer/bc_grid_layout_test.dart` (7), `example/test/block_catalog_section_test.dart` (4 — the section had zero coverage before). All three fail on unfixed code (verified by stashing the lib changes) and pass after.
- **Updated:** 5 existing assertions in `bc_renderer_test.dart` / `bc_renderer_integration_test.dart` that asserted `GridView` internals; they now assert the observable contract (column count via widget positions).
- **Verification:** `flutter test` 2131 pass / 0 fail (was 2106 + 25 new). `flutter analyze` 0 errors; 3 pre-existing warnings in `bc_chatter.dart` untouched. Live demo in Chrome: Block Catalog renders clean, 0 console errors, both bugs visually gone.
- **Not fixed (pre-existing, out of scope):** 5 example-app section tests fail identically with and without this change (`Display`, `Media`, `Views`, `Native`, and the section sweep) — unrelated overflow/unbounded-height issues in those sections.

## 2026-07-05 — Full demo parity sweep (17 components)

- **What:** 17 components added to demo parity: BcContextMenu, BcPortal, BcNotificationCenter, BcDrawer, BcStack, BcForm, BcTimeField, BcTagGroup, BcFieldDynLink, BcList, BcViewerDocument, BcViewerInstagram, BcViewerTiktok, BcPivot, BcFormBuilder, BcSchemaBuilder, BcWorkflowBuilder.
- **Files:** 9 existing sections modified (forms, input_extended, overlays, navigation, primitives, data, media, views, main.dart) + 1 new section (`builders_section.dart`).
- **Barrel fix:** Added missing `export 'src/components/bc_pivot.dart'` to `bc_ui_flutter.dart`.
- **Verification:** flutter test 2106 pass / 0 fail; flutter analyze 0 issues on example app.

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
