# Bucket: ui-library

> Memory for `packages/ui` Mitosis cross-framework component library work.
> Read this when resuming work on Mitosis components, contracts, ADR-007/008/009/010, token pipeline.
> Keep entries short. Newest on top. Date every entry.

---

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
