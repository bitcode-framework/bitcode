# Bucket: demo-app

> React demo (`packages/ui/demo-app`, 127+ components / 14+ groups) + Flutter demo (`packages/ui-flutter/example`) parity work.
> Source of truth for status: `docs/adr/016-demo-app-qa-bug-audit.md` + handoffs.

### [2026-10-29] Demo startup unblocked: patch-package header + .lite resolution

- **Symptom 1:** `npm install` in `demo-app` failed with `Failed to apply patch for package @builder.io/mitosis ... patch file could not be parsed`.
- **Root cause 1:** `patches/@builder.io+mitosis+0.13.1.patch` hunk 1 header said `@@ -237,7 +237,9 @@` but the hunk body has 4 added lines (`+237,11`). patch-package's `verifyHunkIntegrity` throws, then `read.js` mislabels it as a parse error. Fixed the header to `+237,11`. Verified against a pristine `npm pack @builder.io/mitosis@0.13.1` tarball in a temp sandbox: applies clean, idempotent on rerun.
- **Symptom 2 (the real `npm run dev` blocker):** `Cannot assign to "reactiveStoreRef" because it is a constant` during Vite dep scan, pointing at `src/components/views/form/BcForm.lite.tsx`.
- **Root cause 2:** `demo-app/vite.config.ts` `redirectRelativeCoreImports()` rewrites `@/core` → `src/core`, so `src/core/block-renderer/ready-catalog.ts` and `src/core/renderer/BcRenderer.lite.tsx` load, and *their* relative `.lite` imports (`BcLabel.lite`, `BcForm.lite`, `BcDataTable.lite`, …) resolved to raw Mitosis source instead of compiled output. esbuild cannot transpile `.lite.tsx`. NOT a stale-build problem — a full `npm run build` does not fix it.
- **Fix:** added an early branch in `redirectRelativeCoreImports()`: any `.lite` import whose importer lives under `packages/ui/src/` now resolves to the compiled `output/react/src/<rel>.tsx` when that file exists.
- **Also note:** the earlier esbuild warnings were masked because `npm run build` in `packages/ui` was never run after the last `src/` edit; running `mitosis build` alone leaves `output/` unformatted (`BcDataTable.tsx` fails to parse). Always use `npm run build` (runs reformat + sync + index post-steps), never bare `mitosis build`.
- **Env gotcha:** `cmd`-spawned npm scripts in this shell see an empty `PATH`, so `patch-package` resolves as "not recognized". Workaround: prefix `$env:PATH="C:\Program Files\nodejs;C:\Windows\System32;C:\Windows;..."` before `npm install`.
- **Verification:** `npm install` EXIT 0; `npm run build` (ui) EXIT 0 / 204 components; demo-app `vite build` EXIT 0 (1142 modules); `vite` dev server serves both `/` (352k chars body) and `/block-renderer-test.html` with 0 page errors; vitest 822 pass / 0 fail.
- **Files:** `packages/ui/patches/@builder.io+mitosis+0.13.1.patch`, `packages/ui/demo-app/vite.config.ts`, `packages/ui/demo-app/README.md`.
- **Cross-ref:** tooling/build details → `ui-library` bucket.

### [2026-07-05] Full Demo Coverage Milestone (204/204) — Droid

- **Status:** done. All fatal bugs fixed, all 77 missing components demoed.
- **Web:** 11 new demo files (FormsExtended, OverlaysExtended, NavigationExtended, DataDisplayExtended, FeedbackExtended, DataExtended, ChartsExtended, MediaExtended, SocialExtended, LayoutExtended, Builders). 77 new component demos + 11 group entries added to App.tsx. Total: 204 components across 25 groups.
- **Flutter:** 17 components added across 9 sections + 1 new section (builders_section.dart). Missing `bc_pivot.dart` export added to barrel.
- **Bug fixes (10 source-level):** 4 string-style bugs (BcWorkflowBuilder, BcTextField, BcPrintPreview, BcDashboard), 1 dynamic-tag bug (BcChip), 4 const-outside-useStore bugs (BcChip, BcAlert, BcKanbanBoard, BcLookupModal), 1 null-safety fix (BcLightbox.web.tsx), 14 adapt-poc index.ts export fixes.
- **Mitosis patterns learned:** (a) `const` outside `useStore` is dropped by compiler — use getters; (b) `style` prop must be object not string for React target; (c) dynamic-tag `<Tag>` unsupported — use `<Show>` branches.
- **Verification:** tsc EXIT 0, vite build EXIT 0, vitest 822 pass / 0 fail, flutter test 2106 pass / 0 fail, Playwright QA all groups render.
- **Docs:** ADR-016 Phase 5 section appended. See `docs/adr/016-demo-app-qa-bug-audit.md`.

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
