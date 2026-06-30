# Bucket: demo-app

> React demo (`packages/ui/demo-app`, 127+ components / 14+ groups) + Flutter demo (`packages/ui-flutter/example`) parity work.
> Source of truth for status: `docs/adr/016-demo-app-qa-bug-audit.md` + handoffs.

## Active threads

### [2026-07-01] Block Catalog Preview Demos — Agent K
- **Status:** done
- **Files touched:** `packages/ui/demo-app/src/demos/block-catalog/` (9 new files: demoData, demoActionRuntime, DiagnosticPanel, 6 demo components), `packages/ui/demo-app/src/App.tsx` (appended 1 group), `packages/ui-flutter/example/lib/sections/block_catalog_section.dart` (new), `packages/ui-flutter/example/lib/main.dart` (appended 1 nav item)
- **Evidence:** Web tsc clean for all block-catalog files. Flutter analyze 0 issues on my files. Vite build blocked by pre-existing BcChip i18n barrel breakage (not mine, from HEAD of packages/ui submodule).
- **What:** Added "Block Catalog (Preview)" demo group: CatalogOverview (candidate validation), SimpleBlockPreview (CRM hero via BcRenderer), NestedLayoutPreview (Sales Order detail), ExpressionReactivityPreview (employee form reactions + AST preview), ActionBindingPreview (Create Contact via committed ActionRunner), ValidationDiagnosticsPreview (expression/action/result diagnostics). Flutter: BlockCatalogSection (Sales Order detail via Flutter BcRenderer + platform-aware button).
- **Honest labels:** All demos labeled "Preview". Uses production BcRenderer as lower-level renderer P4 BcBlockRenderer will delegate to. ActionRunner from committed P3 (main a8a28a8). Mock adapters are demo-local (demoActionRuntime.ts). No Ready promotion.
- **Gotcha:** `UiSchema`/`FieldNode`/`Reaction` types must be imported from `@bitcode/ui-schema`, NOT `@bitcode/ui` (the demo-app's custom bitcode-ui.d.ts doesn't declare them). Demo-app resolves `@bitcode/ui-schema` transitively via packages/ui.

### [2026-06-29] ADR-016 Demo App QA — ALL 4 PHASES COMPLETE
- **Status**: Done. 5 architectural decisions landed:
  1. `useRef` migration + lint guardrail (no `useStore` for DOM refs)
  2. Post-build portal AST transform (overlay portal fix)
  3. `@source` Tailwind consumption contract
  4. Contract-alignment recipe
  5. ECharts dependency-injection loader pattern
- 28 stub components ported across 5 clusters.
- 3 pre-existing runtime bugs fixed: `getAnimationClasses` import, BcTabs API alignment, BcSplitter dragging init.
- Charts fixed via DI loader (replaces direct echarts import).
- **Next**: none blocking. Re-verify after any component contract change.

### [2026-06-29] Stub Completion + Flutter Review (separate from ADR-014)
- Scope: `packages/ui` + `packages/ui-flutter`, OUTSIDE engine integration.
- 12 mandatory principles in handoff (robust, reliable, future-proof, production-grade, etc.).
- Forbidden: `as any`, `@ts-ignore`, hardcoded colors, stub placeholders, `useStore` for DOM refs.

## Flutter demo parity gap (as of 2026-06-28)

React demo has 127 components / 14 groups. Flutter demo had only ~25 across 6 sections.
Groups still at 0 in Flutter demo: Input(Extended), Primitives, Feedback, Media, Charts, Print, Views, Social.

## Files touched (recent)
- `packages/ui/demo-app/**` — React demo
- `packages/ui-flutter/example/**` — Flutter demo
- `packages/ui/src/**/*.lite.tsx` — component sources (stub ports)
- `scripts/post-build-portal-ast.*` — overlay portal AST transform
