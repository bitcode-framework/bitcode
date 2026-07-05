# Bucket: ui-library

> Memory for `packages/ui` Mitosis cross-framework component library work.
> Read this when resuming work on Mitosis components, contracts, ADR-007/008/009/010, token pipeline.
> Keep entries short. Newest on top. Date every entry.

---

## 2026-07-05 — Mitosis compiler bug sweep (10 source-level fixes)

- **Status:** done. All 10 bugs found during demo-app browser QA, fixed at source level.
- **Pattern 1 — const-outside-useStore (4 fixes):** Mitosis silently drops `const` declarations outside `useStore`. Fixed in BcChip (variant), BcAlert (severity, showBanner), BcKanbanBoard (selectedCard), BcLookupModal (colSpan). Solution: move to `useStore` getters.
- **Pattern 2 — string style prop (4 fixes):** Mitosis passes `style` attribute value unchanged. React expects object, not CSS string. Fixed in BcWorkflowBuilder, BcTextField, BcPrintPreview, BcDashboard.
- **Pattern 3 — dynamic tag (1 fix):** `const Tag = isInteractive ? 'button' : 'div'; <Tag>` — Mitosis doesn't support dynamic element types. Fixed BcChip with explicit `<Show>` branches.
- **Pattern 4 — adapt-poc index.ts exports (14 fixes):** `export { CompName, default } from './Comp'` fails when target file only has `export default`. Fixed to `export { default as CompName } from './Comp'; export { default } from './Comp'`.
- **Null safety:** BcLightbox.web.tsx `active.type` → `active?.type` (4 access points).
- **Test updates:** BcChip.test.ts, BcLightbox.test.ts, BcShare.test.ts regex patterns updated to match fixed source code.
- **Verification:** vitest 822 pass / 0 fail. tsc EXIT 0. vite build EXIT 0.
- **Docs:** ADR-016 Phase 5 documents all fixes with root causes.

---

## 2026-07-01 — BcNotificationCenter + BcNotificationItem (Phase 3 Notification)

- **BcNotificationCenter.lite.tsx** — `packages/ui/src/components/overlays/notification/`. Bell button + dropdown panel, unread badge (99+ overflow), mark-read per item, mark-all-read, infinite scroll pagination, loading/error/empty states. Declarative `data-source` attribute for API binding. Uses `useStore` pattern.
- **BcNotificationItem.lite.tsx** — Subcomponent. 6 variants (info/warning/error/success/mention/message) with semantic icons + token-driven palette. role=article, unread dot, mark-read button, tabular-nums timestamp.
- **BcNotificationCenter.contract.ts** — Types: `NotificationType`, `NotificationItem`, `BcNotificationCenterProps`. 6 type variants.
- **BcNotificationCenter.test.ts** — 21 source tests: contract, bell button a11y, badge, dropdown dialog, mark-all-read, mark-read, infinite scroll, loading/error/empty states, declarative data source, token palette, 6 variants, BcNotificationItem subcomponent.
- **Barrel exports:** `packages/ui/src/index.ts` updated with BcNotificationCenter + BcNotificationItem.
- **Platform audit:** `packages/ui-schema/src/block-catalog/platform-audit.ts` — added BcNotificationCenter + BcNotificationItem entries.
- **Companion:** Flutter `bc_notification_center.dart` + 10 widget tests (see flutter bucket). @worker

---

## 2026-07-01 — BcPivot Mitosis component (Phase 3 Pivot Table)

- **BcPivot.lite.tsx** — `packages/ui/src/components/data/pivot/`. Renders pivot matrix from declarative data source (endpoint+module+model) or inline `result` prop. Expandable/collapsible rows, row totals, grand total, no-data state. Uses `useStore` pattern (no refs, no hooks).
- **BcPivot.contract.ts** — Types: `PivotValueField`, `PivotResult`, `BcPivotConfig`. 5 aggregation types (sum/count/avg/min/max). declarative endpoint config.
- **BcPivot.test.ts** — 14 source tests: result structure, 5 agg types, aria attributes, loading/error/no-data states, expandable rows, row+grand totals, declarative data source, token palette, tabular-nums, formatCell null/NaN.
- **Block Catalog:** `analytics.pivot` entry in `packages/ui-schema/src/block-catalog/categories/analytics.ts`.
- **Build:** `npm run build` — 198 components. `tsc --noEmit` passes. Source tests 14/14 pass.
- **Companion:** Flutter `bc_pivot.dart` + 12 widget tests (see flutter bucket).

---

## 2026-07-01 — Wizard schema-driven multi-step form (Steps 1-4 complete)

- **Step 1 (ui-schema):** Added `WizardStepSchema`, `WizardValidationMode` types + `'wizard'` to `UiSchemaKind` + `steps?: WizardStepSchema[]` to `UiSchema`. Barrels updated. 308/308 tests pass.
- **Step 2 (web BcWizard):** Rewrote `BcWizard.lite.tsx` — per-step icon, schema-driven content via `renderStepContent` callback (ADR-005 §2 compositor, BcRenderer NOT imported directly), client+server validation (ADR-004 ValidationRule), `onFinish(mergedValues)`, `skippable`, `values`/`onValuesChange`. Contract updated with all new fields.
  - **Mitosis compiler workarounds:** (a) removed `export interface` blocks from `.lite.tsx` — Babel's TS parser can't round-trip `=> any;` in interface function-type properties when collapsed to one line by Mitosis; (b) stripped ALL type annotations from `useStore` methods (Mitosis leaks them into compiled output); (c) used `renderStepContent` callback prop instead of direct BcRenderer import (Mitosis places imports mid-file).
  - **Mitosis plugin fix:** Extended `reactFixesPlugin` step 5 in `mitosis-plugin-react-fixes.mjs` — added `)` → `const [` and `)` → `function` semicolon insertion (for multi-useState + useStore method patterns), added `export` to `}` → keyword regex, added `};` → newline before `function`/`export` (prevents Babel parse failure on collapsed interface+function on same line).
  - **Build:** `mitosis build` 190 components, `reformat-react-output` 0 failures (was 1). Source tests 18/18 pass. RTL tests 2/2 pass (axe-core zero violations). Full suite: 2143/2143 source + 822/822 RTL pass.
- **Step 3 (Flutter):** Updated `bc_wizard.dart` — `BcWizardStep` gains `icon`, `schema`, `validationMode`, `validationEndpoint`, `skippable`. `_BcWizardPageState` gains `_stepValues`, `_validating`, `_stepErrors`, client validation (`_validateField`/`_validateStepClient`), server validation (`_validateStepServer` via http.post), `_validateCurrentStep` cascade, `_handleNext`/`_handleFinish` async with validation, `_buildStepNode` renders icon via `_iconData`, `_buildContent` uses 3-arg `stepContentBuilder(index, step, values)`. 9/9 Flutter tests pass. 1950/1950 full Flutter suite pass.
- **Step 4 (docs):** Removed "Wizard schema convention not yet authored" gap from block catalog. Updated `docs/components-progress.md` wizard row. Updated `CONTRACT_INVENTORY.md` Flutter wizard entry. Created `engine/docs/features/wizard.md`.
- **Also fixed:** BcChatter.lite.tsx — moved `parseMentions` standalone function into `useStore` (same Mitosis limitation). This was another agent's WIP that broke the build. All 18 BcChatter source tests pass.
- **Commits:** `88fd11f` (wizard web), `a003b6f` (chatter Mitosis fix), `c9438e9` (ui pointer bump), `e03327b` (wizard Flutter), `c48fcb7` (ui-schema types).

## 2026-07-01 — CRITICAL FIX: vitest .web.tsx resolve.extensions (Droid)

- **What:** The entire `packages/ui` vitest runtime test suite was BROKEN (171/182 test files failed, 733 tests failed) because `vitest.config.ts` was missing `.web.tsx`/`.web.ts`/`.web.jsx`/`.web.js` in `resolve.extensions`. Adapt components (BcLightbox, BcMap, BcWebcam, BcShare, BcGridView, BcNative*) use the `.web.tsx` extension per ADR-006 §3, and their entry `index.ts` does `export { BcLightbox } from './BcLightbox'` (extensionless). Vite's default extensions list does NOT include `.web.*`, so resolution failed for ANY test importing `@bitcode/ui` (which aliases to `output/react/src`).
- **Root cause:** The Adapt entry `index.ts` comments say "Vite sees `extensions: ['.web.tsx', ...]`" — but that only holds in `demo-app/vite.config.ts`, NOT in `vitest.config.ts`. The vitest config was never updated when Adapt components were introduced (ADR-016 Phase 4, commit `9fbd97c`).
- **Fix:** Added `extensions: ['.web.tsx', '.web.ts', '.web.jsx', '.web.js', '.tsx', '.ts', '.jsx', '.js', '.mjs', '.json']` to `vitest.config.ts` `resolve` block, mirroring `demo-app/vite.config.ts`.
- **Verification:** Before = 733 fail / 89 pass. After = **822 pass / 0 fail / 13 skipped** (176 test files passed, 6 skipped). Single-test confirm: BcDataTable.rtl.test.tsx 38/38 pass.
- **Commits:** `packages/ui` `77b94ff` (config fix) + root `658468a0` (submodule bump).
- **Gotcha:** The bucket's prior claim "2131/2131 source tests pass" was TRUE but MISLEADING — `npm test` runs `node scripts/run-source-tests.mjs && vitest run`. The source tests (node --test) passed, but `vitest run` was completely broken. Always run BOTH halves of `npm test` and check vitest exit code separately.
- **Note:** `npx tsc --noEmit` in packages/ui has PRE-EXISTING errors (rootDir/output inclusion, missing @types/react, .native.tsx node:test). These are unrelated to this fix and not introduced by it. The vitest.config.ts change is not even in the tsc `include` glob.

## 2026-07-01 — Agent cleanup and verification sweep

**Status:** done
**Files touched:** `packages/ui-schema/src/actions/validate.ts`, `packages/ui-schema/src/block-catalog/{validate,validate-catalog-deep}.ts`, `packages/ui-schema/src/block-catalog/policies/performance-policy.ts`, `packages/ui-schema/src/block-renderer/{renderer,validate-runtime}.ts`, `.multibrain/session.md`, `.multibrain/buckets/ui-library.md`
**Evidence:** `npx tsc --noEmit` PASS; `node scripts/run-tests.mjs` 299 pass / 0 fail; `node scripts/run-source-tests.mjs` 2131 pass / 0 fail; `node scripts/contract-audit.mjs` PASS 0 violations; `npm run build` in demo-app still blocked only by pre-existing demo/ui export/type gaps after ui-schema unused-import errors were cleared.

Cleared strict TypeScript unused-import diagnostics left across Agent G/H/I/L files. Agent I output is committed at `6c1ebaf8` and independently green; no Ready promotion was made.

## 2026-07-01 — Deep Block Catalog validator (Agent I)

**Status:** done
**Files touched:** `packages/ui-schema/src/block-catalog/{catalog-diagnostic,validate-catalog-deep}.ts`, `packages/ui-schema/src/{block-catalog/index,index}.ts`, `packages/ui-schema/test/block-catalog/validate-catalog-deep.test.ts`
**Evidence:** `npx tsc --noEmit` EXIT 0; `node scripts/run-tests.mjs` 299 pass / 0 fail.

Added `validateBlockCatalogEntryDeep()` without Ready promotion. It composes P1 shape validation + Agent F expression validator (`71885aa`) + Agent G action validator (`a8a28a8`), with path-scoped guards for embedded expressions/actions, DataSource structural checks, serializability warning/error, nesting/cycle detection, minimal action/permission adapter gates, and per-platform component registry probes. Policy readiness gates stay Agent L scope.

## 2026-07-01 — BcBlockRenderer portable engine + web binding (P4, ADR-017 S3, Agent H)

- **Status:** done
- **Files touched:** packages/ui-schema/src/block-renderer/ (5 files), packages/ui-schema/test/block-renderer/renderer.test.ts (29 tests), packages/ui-schema/src/index.ts (export), packages/ui/src/core/block-renderer/ (4 files + test), packages/ui/scripts/generate-index.mjs, packages/ui/src/index.ts (regenerated), packages/ui-flutter/lib/src/block_renderer/ (4 files), packages/ui-flutter/lib/bc_ui_flutter.dart
- **Evidence:** ui-schema tsc 0 errors, npm test 283 pass/0 fail; packages/ui audit:contract 0 violations, npm test 2131 pass/0 fail; flutter dart analyze 0 issues.
- **Commits:** root 3628e4b (engine + Dart mirror), packages/ui submodule 58ffd5b (web binding), root 32cc976 (pointer bump).
- **Gotcha:** G (ActionRunner) not yet ported to Dart — Dart contract mirror types action-bearing surfaces as Object? with deferred TODO. F (expressions) is ported to Dart; renderer engine uses F's buildPermissionFilteredScope verbatim. Coordination: Agent L owns block-catalog/ (categories/policies/pilot); Agent H owns block-renderer/ exclusively. index.ts is shared — both agents' exports preserved verbatim.

Landed P4 BcBlockRenderer foundation as a render-plan two-tier engine per the 10 review amendments. Portable engine emits RenderPlan union (Schema/Component/Hidden/Diagnostic) — NEVER framework output. Registry stores ComponentMetadata ONLY. Actions deferred (runner NOT invoked at render). Web binding translates plans via BcRenderer + WebBlockComponentRegistry. No entries promoted to Ready.

---

## 2026-07-01 — Block Catalog Ready-promotion gate audit (Agent M)

**Status:** done (no promotion; groundwork + blocker report)
**Files touched:** `sprints/package-ui-elementor/ready-promotion-audit.md`, `packages/ui-schema/test/block-catalog/ready-gate.test.ts`
**Evidence:** `npx tsc --noEmit` EXIT 0; ready-gate test 8/8 pass; `rg readiness.*'ready'` in catalog source = 0 claims.

Verdict: NO Ready promotion this session. P3 ActionRunner is untracked + 2 failing tests (A4 malformed-result, E14 resolveFirst); P4 `BcBlockRenderer`/`catalog:validate` CLI do not exist. P1 shape validator emits `BLOCK_READY_GATE_INCOMPLETE` for any Ready entry. Smallest safe Ready slice identified: `content.hero` + `content.plain-text` (static-only capabilities, not blocked by P3 action runtime). Web+Flutter source/export gates verified PASS for both with exact file:line evidence. Added future-proof gate test (contract: "every Ready entry passes strict shape validation"; NOT "Ready==0 forever"). Rejected-candidate inventory + unblock map in audit artifact. Do NOT create `packages/ui/src/block-catalog/` registry until Agent H/J define ownership.

---

## 2026-07-01 - Portable Expression AST foundation (P2, ADR-018 S5-S6)

**Status:** done
**Files touched:** packages/ui-schema/src/expressions/ (12 TS files), packages/ui-flutter/lib/src/expressions/ (10 Dart files), packages/ui-schema/test/ (3 test files + 9 fixture JSONs), packages/ui-flutter/test/expressions/ (2 test files)
**Evidence:** 32 TS tests pass, 136 Dart tests pass, dart analyze 0 issues, tsc --noEmit passes.

Landed P2 Portable Expression AST: 5-node-kind discriminated union (Literal/Path/Op/Call/If), 24 operators with arity bounds, 70+ built-in functions manifest, static validator with FORBIDDEN_FUNCTION_NAMES denylist + path reserved-roots/forbidden-segments, permission-filtered scope builder, versioned JSON serialization, 9 golden JSON fixtures (122 cases) as cross-platform source of truth. Dart mirror uses manual sealed class for wire-shape parity (no freezed). Next: P3 ActionSource + runtime evaluator.

## 2026-06-30 - Agent B/D registry integration fix

**Status:** done
**Files touched:** `packages/ui/scripts/generate-registry.mjs`, `packages/ui/demos/playground/src/registry.ts`, root `packages/ui` submodule pointer
**Evidence:** `packages/ui` commit `e6094ce`; root commit `b55aeec`; registry has 196 components incl. `BcShare`, `BcGridView`, `BcMarkdown`; source tests 2126 pass; contract audit 0 violations.

Fixed Adapt registry scan bug by reading child folder entries, added output-dir creation before registry write, generated playground registry, then bumped root submodule pointer. Targeted BcMarkdown RTL tests pass with known jsdom canvas/key warnings.

## 2026-06-30 — Canonical web BcMarkdown (Agent D)

- **What:** Implemented canonical web `BcMarkdown` in `packages/ui` per ADR-017/018 P5 batch.
- **Outcome:** Edit/preview modes, toolbar (bold/italic/heading/list/link/code), controlled/uncontrolled value+view, sanitize default on. Pure dependency-free markdown engine (types/sanitize/parser/barrel). 21 parser tests + 12 contract tests + 14 runtime tests all pass. Contract audit 0 violations.
- **Files:** `src/core/markdown/{types,sanitize,parser,index,parser.test}.ts`, `src/components/content/markdown/{BcMarkdown.lite.tsx,BcMarkdown.contract.ts,BcMarkdown.test.ts,BcMarkdown.rtl.test.tsx}`, `scripts/{generate-index,generate-registry}.mjs`.
- **Gotcha:** Do NOT promote `content.rich-text` Ready — gaps remain (Flutter BcMarkdown export hardening, contracts/decisions still tracked in matrix). Pre-existing `BcLightbox` barrel breakage blocks `@bitcode/ui` barrel imports in unrelated RTL tests — not owned by this work.
- **Mitosis quirks discovered:** state methods returning JSX with `<For>` get extracted as standalone functions losing For reference (use `.map()` instead); `onInput` transforms to `onChange={e => fn(e.target.value)}` passing string not event (guard both shapes); `onUpdate` runs after every render (track prev value to avoid resetting uncontrolled state).

## 2026-06-29 — Stub completion + Flutter review session (in progress)

**Source:** `handoffs/2026-06-29-ui-stub-completion-and-flutter-review.md`

- Session scope: `packages/ui` + `packages/ui-flutter` only (NOT engine / ADR-014).
- 28 stub components across 5 clusters being ported to real implementations (ADR-016 follow-up).
- **Mandatory code principles (12 points)** documented in the handoff — read before any component edit. Key ones: single source + platform-aware, mobile platform-aware (iOS Cupertino / Android Material), ADR-008 contract freeze respected.
- **Forbidden anti-patterns:** `as any`/`@ts-ignore`, hardcoded color/size, stub `return null`/placeholder in prod path, `console.log` in prod, `useStore({xxxRef})` for DOM ref (use `useRef` + lint guardrail from ADR-016).

## 2026-06-21 — ADR-008 partAttr wiring complete

- `partAttr` runtime part resolution (ADR-009) fully wired across components.
- Phase C1 of ADR-001 complete.

## 2026-06-18 — ADR-008 Contract-First API Freeze: Core Integration Surface COMPLETE

**Source:** `handoffs/2026-06-18-adr-008-final-session.md`

- All 25 components + 7 infrastructure items have complete ADR-008 contract surfaces (Parts, Methods, WebComponentEvents, A11yContract, AdapterMapping).
- 1538 tests passing, 0 failing. Build verified 151/151 components × 4 targets. 25 runtime validation tests added.
- **Remaining:** Phase 3-5 (Accessibility/Performance testing, Documentation, Flutter/BitCode integration) est. 11-18 hours.

## Standing facts

- `packages/ui` = `@bitcode/ui`, Mitosis 0.13.1, 4-target output (React/Vue/Stencil/RN).
- Test runner: `node --test "src/**/*.test.ts"` via `npm test`. 42 pass, 11 pre-existing BcStack fails (need `npm run build` first — NOT a regression).
- tsconfig: `allowImportingTsExtensions: true` + `noEmit: true`; runtime imports use `.ts` ext.
- LSP errors on `@builder.io/mitosis` = TS server unaware of node_modules; not real.
- Token pipeline: `tokens:check` ok; generator in `scripts/generate-tokens.mjs`.
- ADR-009 = Runtime Part Resolution (`partClass`/`partAttr` pure fns, no dynamic spread — Mitosis can't compile it).
- ADR-010 = Web Component Event Bridge (`dispatchBcEvent`, `bc-*` semantic events, SSR/RN safe no-op).
