# Handoff: UI Test Suite Repair + Cross-Package Audit (ADR-011→018)

**Date:** 2026-07-01
**Agent:** Droid (context-compacted continuation)
**Scope:** Audit packages/ui-schema, packages/ui, packages/ui-flutter for incomplete features per ADR-011→018, then synthesize findings. Discovered and fixed a critical test-suite blocker along the way.

---

## TL;DR

- **CRITICAL FIX (committed):** `packages/ui` vitest runtime test suite was **100% broken** (171/182 test files failing, 733/835 tests failing) because `vitest.config.ts` was missing `.web.tsx` in `resolve.extensions`. Adapt components (BcLightbox, BcMap, BcWebcam, BcShare, BcGridView, BcNative*) use the `.web.tsx` extension per ADR-006 §3, and the Adapt entry `index.ts` files do `export { BcLightbox } from './BcLightbox'` (extensionless). Vite's default extension list does NOT include `.web.*`, so resolution failed. Fixed by mirroring `demo-app/vite.config.ts` extensions. **After fix: 176 test files passed, 822 tests passed, 0 failed.**
- **ui-schema:** Healthy. 308 tests pass, typecheck clean. ADR-017/018 types, validators, BlockRenderer, action runtime, portable expression AST all implemented.
- **ui-flutter:** Healthy. 1952+ tests pass (`flutter test` "All tests passed!"), `dart analyze` clean (exit 0). ADR-018 Dart action parity complete. ADR-015 chart rewrite (flutter_echarts) complete. ADR-016 Phase 4 stubs ported.
- **ADR-017 blocker list is STALE:** BcShare, BcGridView, BcMarkdown (web); BcAvatar, BcBadge, BcProgress, BcKpi (Flutter) — ALL now exist and are exported on both platforms (confirmed via live audit 2026-07-01, see adr bucket).

---

## What was done this session

### 1. Critical fix: vitest `.web.tsx` resolve.extensions (committed)

**Root cause (evidence-based):**
- **Symptom:** `npm test` in `packages/ui` → `vitest run` exits 1 with `Error: Failed to resolve import "./BcLightbox" from "output/react/src/components/media/lightbox/index.ts"`. 171 test files fail, 733 tests fail.
- **Fault location:** `packages/ui/vitest.config.ts` — `resolve.extensions` was unset, so Vite used defaults `['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json']` which do NOT include `.web.tsx`.
- **Causal chain:** Adapt components (ADR-006 §3) use `.web.tsx`/`.native.tsx` split. Their `index.ts` entry re-exports `./BcLightbox` (extensionless) relying on bundler extension resolution. `demo-app/vite.config.ts` sets `extensions: ['.web.tsx', '.web.ts', ...]` so the demo works. But `vitest.config.ts` never added these, so ANY test importing `@bitcode/ui` (→ `output/react/src` via alias) that transitively touches an Adapt component fails to resolve. Since the root `index.ts` barrel re-exports ALL components including Adapt ones, the failure cascades to the ENTIRE runtime test suite.
- **Minimal fix:** Add `extensions: ['.web.tsx', '.web.ts', '.web.jsx', '.web.js', '.tsx', '.ts', '.jsx', '.js', '.mjs', '.json']` to `vitest.config.ts` resolve block (mirrors demo-app).
- **Verification:** Single test `BcDataTable.rtl.test.tsx` (38 tests, was 38 FAIL) → 38 PASS. Full suite → 176 files passed, 822 tests passed, 13 skipped, 0 failed.

**Commits:**
- `packages/ui` submodule: `77b94ff fix(ui): add .web.tsx to vitest resolve.extensions for Adapt components`
- parent repo: `658468a0 fix(ui): bump packages/ui submodule — vitest .web.tsx resolve.extensions fix`

**Why this was missed before:** The ui-library bucket claimed "2131/2131 source tests pass" — that's TRUE for `node scripts/run-source-tests.mjs` (the source-level test runner), but `npm test` also runs `vitest run` (runtime tests) which was completely broken. The source tests don't import `@bitcode/ui` (they test `.ts`/`.lite.tsx` source directly), so they never hit the Adapt resolution path.

### 2. Audit results (read-only, no changes)

| Package | Tests | Typecheck/Analyze | Status |
|---------|-------|-------------------|--------|
| `packages/ui-schema` | 308 pass / 0 fail | `tsc --noEmit` clean | ✅ Healthy |
| `packages/ui` | source: pass; **vitest: was 100% broken, now FIXED (822 pass)** | `tsc` has pre-existing errors (rootDir/output, missing @types/react, .native.tsx) — NOT caused by my change, not in scope | ⚠️ vitest fixed; tsc pre-existing |
| `packages/ui-flutter` | 1952+ pass ("All tests passed!") | `dart analyze` clean (exit 0) | ✅ Healthy |

---

## ADR status summary (011→018)

| ADR | Title | Status | Notes |
|-----|-------|--------|-------|
| 011 | Flutter Adapter Architecture | Accepted (Phase 5a/5b done, 5c pending) | 22 P0 + 8 P1 widgets, BcRenderer Dart port done. P2 (~90 widgets) remaining. |
| 012 | Extract ui-schema Standalone Package | Accepted (complete) | All portable types extracted. Re-export bridges in @bitcode/ui. |
| 013 | Cross-Tier Infrastructure Completion | Accepted (Phase 1-5b complete) | AppProvider, 16 hooks, animation+adapt runtime, Flutter parity. |
| 014 | Engine ↔ UI Feature Wiring | **Proposed** | 9 features (comment, activity, export, wizard, pivot, notification, 3 visual builders). Phase 1-4 all unchecked. Backend endpoints + BcRenderer parts NOT implemented. |
| 015 | Flutter Chart Strategy (flutter_echarts) | Accepted (F1-F4 complete) | CustomPaint removed (912→288 lines). 26/26 chart types via flutter_echarts. Golden JSON parity 36 fixtures. |
| 016 | Demo App QA Bug Audit | Accepted (Phase 1-4 complete) | 28 stubs ported, Mitosis ref bug fixed, overlay portal fixed, @source Tailwind contract, echarts loader DI. |
| 017 | Block Catalog Architecture | Accepted (matrix 0.1) | Types in ui-schema, BcBlockRenderer, validators, Ready gate. Ready promotions started. |
| 018 | Data/Action/Reactivity Contract | Accepted | ActionSource, ActionRunner, PortableExpression AST, effect runner, lifecycle. TS + Dart parity. |

---

## What's NOT done (honest gaps)

### ADR-014 (Engine ↔ UI Wiring) — LARGEST remaining work
- [ ] `comment` model + `/api/v1/{module}/comments` endpoint + BcChatter wiring
- [ ] `/api/v1/activity` endpoint + BcActivity wiring
- [ ] `/api/v1/{module}/{model}/export` endpoint + BcExport wiring
- [ ] `wizard` view type + BcWizard wiring
- [ ] `/api/v1/{module}/{model}/pivot` endpoint + BcPivot component
- [ ] `notification` model + endpoints + WebSocket push + BcNotificationCenter
- [ ] SPA mode (`view.json { "render": "spa" }`)
- [ ] 3 visual builders (BcFormBuilder, BcWorkflowBuilder, BcSchemaBuilder)
- All 9 features are unchecked. This is the biggest unfinished surface.

### ADR-011 Phase 5c (Flutter P2 parity)
- [ ] ~90 P2 widgets not yet ported to Flutter
- P0 (22) + P1 (8) done. P2 is the long tail.

### ADR-017/018 Block Catalog — Ready promotions in progress
- [ ] First Ready promotions done but many candidates still blocked on gaps
- [ ] `ReactiveStore.runAction` integration — action types+runner exist+pass tests in ui-schema but NOT yet wired into `ReactiveStore.runAction` (P7 merge, per adr bucket)
- [ ] BcMarkdown web canonical source exists + exported (gap resolved) but may need test/demo verification

### packages/ui tsc (pre-existing, NOT my regression)
- `tsc --noEmit` in packages/ui has errors: rootDir violations (output/ files), missing @types/react, .native.tsx module resolution, node:test types. These are PRE-EXISTING and unrelated to the vitest fix. The tsconfig `exclude` lists `output` but tsc still picks up output files (likely via .tsbuildinfo cache or allowImportingTsExtensions). This needs a separate investigation if a clean tsc is desired.

---

## Files changed this session

| File | Change | Commit |
|------|--------|--------|
| `packages/ui/vitest.config.ts` | +19 lines: added `resolve.extensions` with `.web.*` entries | `77b94ff` (submodule) |
| `packages/ui` (submodule pointer) | bumped to `77b94ff` | `658468a0` (parent) |
| `.multibrain/buckets/ui-library.md` | +1 entry documenting the fix | (not committed — bucket is AI memory) |

---

## Next steps (concrete)

1. **ADR-014 Phase 1 (wire existing):** Start with Comment (#5) or Activity (#6) — components exist, only API + BcRenderer part needed. Fastest value. Files: `engine/internal/presentation/api/`, `engine/internal/domain/`, view JSON parts.
2. **ADR-017/018 P7 merge:** Wire `ActionRunner` into `ReactiveStore.runAction` (TS + Dart). The types and runner exist and pass tests in ui-schema; the integration into ReactiveStore is the missing piece.
### 5. Wizard multi-step form implementation (committed, beyond original audit scope)

User asked about wizard/multi-step form support like Filament. Researched ADR-005, ADR-004, ADR-007/008, ADR-018. Got user requirements via questionnaire: (1) client+server validation, (2) icon per step, (3) multiple steps, (4) finish callback with values, (5) align with ADRs (not copy Filament). Implemented in 4 steps:

**Step 1 — ui-schema types (committed `c48fcb7`):**
- Added `'wizard'` to `UiSchemaKind` union
- Added `WizardValidationMode = 'client' | 'server' | 'both'` type
- Added `WizardStepSchema` interface (title, description, icon, schema, validationMode, validationEndpoint, skippable, blocking)
- Added `steps?: WizardStepSchema[]` to `UiSchema`
- Updated barrels `src/schema/index.ts` and `src/index.ts`

**Step 2 — Web component (committed `88fd11f` + `a003b6f`):**
- `BcWizard.contract.ts`: Added icon, schema, validationMode, validationEndpoint, skippable, finishAction, onFinish(values), values, onValuesChange, renderStepContent
- `BcWizard.lite.tsx`: Complete rewrite with renderStepContent callback (ADR-005 §2 compositor — BcRenderer NOT imported directly to avoid Mitosis import-hoisting issues), validateField/validateStepClient/validateStepServer, validationMode cascade (client → server), onFinish(mergedValues), error display, validating state, icon per step
- `BcWizard.test.ts`: 18 source tests (icon, schema-driven, validationMode, onFinish values, skippable, renderStepContent callback)
- `mitosis-plugin-react-fixes.mjs`: Extended with 3 new regex rules (5b: `)`→`const[` and `)`→`function` semicolons; 5c: `};`→newline before `function`/`export`) to fix Mitosis output that Babel couldn't parse
- BcChatter.lite.tsx fix: Moved `parseMentions` standalone function into `useStore` as method (Mitosis can't parse standalone functions before component)
- Build: 190 components generated, reformat 486/486 pass (0 failed), 2143 source tests pass, 822 RTL tests pass

**Step 3 — Flutter parity (committed `e03327b`):**
- `bc_wizard.dart`: Added icon (45+ icon name→Icons mapping), schema-driven content via stepContentBuilder(index, step, values), client validation (validateField, validateStepClient), server validation (validateStepServer via http.post), validationMode cascade, merged values tracking, onFinish(values, meta), onValuesChange, skippable, validating state, validation error display
- `bc_wizard_test.dart`: Updated to use onFinishLegacy + stepContentBuilderLegacy for backward compat. All 9 tests pass.
- All 1950 Flutter tests pass.

**Step 4 — Docs (committed `d4c7856`):**
- Block catalog: Removed "Wizard schema convention not yet authored" gap
- `docs/components-progress.md`: Updated wizard row
- `packages/ui-flutter/CONTRACT_INVENTORY.md`: Updated wizard row
- `engine/docs/features/wizard.md`: Created feature doc placeholder
- `.multibrain/buckets/ui-library.md` + `adr.md`: Added wizard entries

**Key design decisions:**
- `renderStepContent` callback prop instead of direct BcRenderer import (ADR-005 §2 compositor pattern, keeps wizard standalone)
- Generic `validationEndpoint` for server validation (standalone-friendly, no BitCode coupling). Contract: `{ step, values } → { ok, errors }`
- Plain `onFinish(values)` as default + optional `finishAction?: ActionSource` for declarative mode (ADR-018)
- Interfaces removed from `.lite.tsx` (Mitosis inlines `export interface` on single line, Babel can't parse function-type properties like `=> void` in that format)

---

## Next steps

1. **BcWizard RTL test expansion:** Currently only 2 a11y smoke tests. Add interaction tests (step navigation, validation blocking, icon rendering, schema-driven content rendering, onFinish with values).
2. **BcChatter RTL test:** Another agent's chatter work needs RTL tests once the contract stabilizes.
3. **packages/ui tsc cleanup (optional):** Investigate why tsc includes `output/` files despite `exclude: ["output"]`. Likely needs `tsconfig.tsbuildinfo` deletion + clean rebuild, or a separate `tsconfig.test.json`.
4. **ADR-011 Phase 5c:** Port P2 Flutter widgets (~90). Long tail, can be parallelized.

---

## Critical context for future sessions

1. **`npm test` in packages/ui runs TWO runners:** `node scripts/run-source-tests.mjs` (source tests, ~2131) AND `vitest run` (runtime tests, ~835). BOTH must pass. Previous bucket entries only reported source test counts.
2. **Adapt components need `.web.tsx` in resolve.extensions** — any Vite/Vitest config that resolves `@bitcode/ui` must include `.web.*` extensions. The demo-app vite.config.ts has this; vitest.config.ts now has this. If a new consumer config is added, it must too.
3. **`tool/run-flutter.ps1`** in ui-flutter sets a minimal PATH to work around bloated PATH breaking flutter.bat's `WHERE /Q git` lookup. Use it (or set the documented minimal PATH) when invoking flutter from a worker shell. The minimal PATH must include: flutter\bin, Git\cmd, System32, System32\Wbem, PowerShell\v1.0, AND Node (if running node-based tooling).
4. **Other agents have uncommitted work** in the parent repo: AGENTS.md (multi-brain section addition), .multibrain/ (bucket entries), engine/go.mod + go.sum, handoffs/2026-06-26-session-adr-015-flutter-chart.md. Do NOT touch these — they belong to other agents.
5. **packages/ui and packages/ui-flutter are git submodules.** Changes inside them require committing in the submodule first, then bumping the pointer in the parent repo.
