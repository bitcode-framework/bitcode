# Bucket: ADR (Architecture Decision Records)

> Source of truth untuk keputusan arsitektur BitCode platform.
> ADR = sumber kebenaran. Kode kontradiksi ADR → perbaiki kode atau buat ADR baru.

**Index lengkap**: `docs/adr/README.md`
**Total ADR**: 018 (012-018 Accepted/complete except 013/014 Proposed)

---

## 2026-07-01 — Wizard schema-driven multi-step form (ADR-005/004/008 alignment)

- Implemented Filament-style wizard with schema-driven step content per ADR-005 §2 (compositor pattern)
- Added `WizardStepSchema`, `WizardValidationMode` types to `packages/ui-schema` (ADR-012 canonical types)
- Web `BcWizard` uses `renderStepContent` callback prop (NOT direct BcRenderer import) — keeps wizard standalone, consumer wires BcRenderer
- Per-step validation: client (ADR-004 ValidationRule) + server (generic `validationEndpoint`, consumer-set for standalone)
- `onFinish(mergedValues)` plain callback + optional `finishAction?: ActionSource` (ADR-018) for declarative mode
- Flutter parity: `BcWizardStep` gains `icon`, `schema`, `validationMode`, `validationEndpoint`, `skippable`; `_BcWizardPageState` gains `_validateCurrentStep()`, `_mergedValues`, `_stepValues`, `_stepErrors`, `_validating`
- Mitosis build fix: `reactFixesPlugin` extended with 5b (`)` → `const [` / `function` semicolons) and 5c (`};` → newline before `function`/`export`) to handle multi-useState + interface-on-same-line patterns
- `export interface` blocks removed from `BcWizard.lite.tsx` — Babel's TypeScript parser can't round-trip them when collapsed to a single line by Mitosis

## 2026-07-01 — ADR-017/018 production-grade: Dart action parity + first Ready promotions (Droid)

- **What:** Completed Phase 3 (Dart action runtime parity) + Phase 4 (first Ready catalog promotions). ADR-017/018 now end-to-end production-grade across TS + Dart + demo-app.
- **Dart parity (Phase 3):** 13 source files in `packages/ui-flutter/lib/src/actions/` mirroring TS `actions/` module field-for-field (sealed unions for ActionSource/ActionEffect, ActionResult normalization, EffectSink + 5 adapter interfaces, ActionRunner with dedup/preflight/timeout-race/classifyThrow, redact, lifecycle). 7 test files, 25/25 pass. ReactiveStore.runAction() bridges store mutations to ActionRunner. block_renderer typed (ActionSource map, ActionHandlerPlan, ActionRunner deps). Commit `670ec98`.
- **Ready promotion (Phase 4):** `content.plain-text` + `content.faq` promoted to `readiness: 'ready'` in new `ready-promotions.ts` (distinct from ALL_CANDIDATE_ENTRIES which stays 44 candidates). Full evidence: schemaPreset, a11yEvidence (WCAG AA), test references, state coverage for data-source, dual platform morph. Deep validator fix: UiSchemaKind check scoped to actual UiSchema nodes only (was false-positive flagging FieldNode.type values). PILOT_SUMMARY.readyCount now derived from READY_CATALOG_ENTRIES.length. 9 new tests. Commit `687cab5`.
- **Verification:** ui-schema tsc 0 errors; ui-schema 308/308 tests pass; packages/ui 2131/2131 source tests pass; contract audit 0 violations; demo-app vite build 0 errors; flutter test actions 25/25 pass; dart analyze 0 errors/warnings (13 info lints).
- **Files touched:** `packages/ui-flutter/lib/src/actions/*.dart` (13 new), `packages/ui-flutter/test/actions/*.dart` (7 new), `packages/ui-flutter/lib/{bc_ui_flutter.dart,src/block_renderer/*,src/core/reactive_store.dart}`, `packages/ui-schema/src/block-catalog/{ready-promotions.ts,index.ts,pilot-candidates.ts,validate-catalog-deep.ts}`, `packages/ui-schema/test/block-catalog/{ready-promotions.test.ts,pilot-candidates.test.ts}`
- **Gotcha / note:** `validate-catalog-deep.ts` deepWalk originally checked `obj.type` as UiSchemaKind for ALL objects with `fieldPath.startsWith('schemaPreset')` — this incorrectly flagged FieldNode objects inside `fields[]` whose `type` is a `FieldType` ('textarea', 'input', etc.), not a `UiSchemaKind`. Fixed by scoping the check to `fieldPath === 'schemaPreset'` or matching `/\.sections\[\d+\]$/`. Also: BcAlert.contract.ts has a pre-existing JSDoc `*/` parse bug (line 21 `data-*/aria-*`) NOT introduced by this work.

---

- **What:** Canonical Candidate metadata seed for all 44 accepted-baseline matrix rows as `CandidateCatalogEntry` records. Every entry `readiness: 'candidate'` — NO Ready promotion.
- **Files:** `packages/ui-schema/src/block-catalog/categories/{content,data,interaction,commerce,workflow,builder,analytics,navigation,plugin-excluded}.ts` + `categories/index.ts` barrel/helpers, `platform-audit.ts` (static snapshot w/ provenance), `pilot-candidates.ts`. 7 test files (42 tests). Commit `42cbf292`.
- **Decision / outcome:** requiredComponents use ACTUAL per-component platform availability (not blanket web+flutter). Gaps/blockedBy explicit: `P3-actions-runtime-unmerged`, `P4-block-renderer-validator`, flutter parity gaps (`bc_section`, `bc_timeline` web-only), ADR-015 chart parity verification. Clean pilots: `content.plain-text`, `content.faq`. Conditional: `content.hero` (BcSection flutter gap), `analytics.chart` (ADR-015 parity), `workflow.diagram-viewer` (BcChartGraph vs BcChart visual parity).
- **Files touched:** `packages/ui-schema/src/block-catalog/{categories/,platform-audit.ts,pilot-candidates.ts}`, `packages/ui-schema/test/block-catalog/*.test.ts`
- **Gotcha / note:** Live audit 2026-07-01 found ADR-017 blocker list (BcShare/BcGridView/BcMarkdown web; Flutter BcAvatar/BcBadge/BcProgress/BcKpi) is STALE — all now exist+exported on both platforms. Only `BcSection`/`BcTimeline` are web-only. Action types+runner exist+pass tests in ui-schema but NOT wired into `ReactiveStore.runAction` (P7 merge). matrix-parity test loads frozen baseline via `fs.readFileSync` (not resolveJsonModule). Barrel exports for categories/platform-audit/pilot were added by Agent L's commit 465826e1; my untracked source files make those imports resolve.

---

## 2026-07-01 — Block Catalog policy validators (Agent L, ADR-017/018)

- **What:** Implemented 5 deterministic build-time policy validators for Block Catalog safety & Ready gates. Compose over Agent I's shape validator via same `BlockCatalogDiagnostic` model; separate `BLOCK_POLICY_*` code space; never flip readiness (Amendment 6).
- **Policies:** security (field-scoped secret scan with false-positive guards, unsafe URL schemes, allowlist declaration, sanitization default-on, untrusted HTML, payment adapter, external-link, unsafe call-string, reactive/action runtime declaration), accessibility (evidence presence/shape only — NO link I/O, label presence, WCAG level, focus-management), adaptive (platform morph cupertino+material per ADR-011, state coverage, breakpoints, reduced motion), performance (nesting/dependency limits, large-data perf declaration, expression complexity, debounce/abort — conservative: warnings default, error only when Ready + relevant cap), offline (full queue semantics for read-write action blocks per ADR-018 S6.4, action/data compatibility).
- **Orchestrator:** `validateBlockCatalogPolicies(entry, config?)` runs all 5 in fixed order, merges via `mergeResults`. Severity scales via `isReadyEntry()`: Ready-blocking gates = error on Ready, warning on Candidate; Deprecated skipped.
- **Files:** `packages/ui-schema/src/block-catalog/policies/{types,security-policy,accessibility-policy,adaptive-policy,performance-policy,offline-policy,index}.ts` + 2 barrel exports + 6 test files.
- **Verification:** `npx tsc --noEmit` clean (EXIT=0); 80 new policy tests pass; full ui-schema suite 232 pass / 0 fail.
- **Commit:** `465826e1` — `feat(ui-schema): add block catalog policy validators (ADR-017/018)`.
- **Gotcha / note:** `Reaction` type (committed) uses JS functions, not PortableExpression AST, and `BlockCatalogEntry.schemaPreset` embeds `UiSchema` whose `actions` are `UiAction` button descriptors — NOT `ActionSource`/AST nodes. So the security policy does NOT fabricate an AST walk over schemaPreset; it scans serializable strings + validates capability/adapter declaration consistency. This stays honest to the merged contract. Agent E's `runner.ts:427` parameter-property syntax (pre-existing test blocker at session start) was fixed by Agent E mid-session; full suite now green.

---

## Status Saat Ini (per 2026-06-30)

| ADR | Judul | Status | Catatan |
|-----|-------|--------|---------|
| 001 | Mitosis Component Library | Accepted (Phase C1) | parts superseded by 002, §15 Lynx superseded by 006 |
| 002 | Cross-Platform Architecture | Accepted | extended by 006/007/009/010 |
| 003 | Component Roadmap & Taxonomy | Accepted (Phase C1) | 14 grup, P0/P1/P2 |
| 004 | Reactive Engine | Accepted | dual storage (useStore+useRef) cleanup |
| 005 | Metadata-Driven UI (BcRenderer) | Accepted (Phases 1-3) | ~5% done per .opencode-context.md |
| 006 | Universal Component Strategy | Accepted (Phase 0 spike) | Pure/Adapt 2-tier |
| 007 | Enterprise Universal UI API Contract | Accepted (Phase 2) | event/state/method/prop/a11y |
| 008 | Contract-First API Freeze | Accepted (Phase 2 + partAttr wiring) | Core Integration Surface complete |
| 009 | Runtime Part Resolution | Accepted | partClass/partAttr pure fns |
| 010 | Web Component Event Bridge | Accepted | dispatchBcEvent |
| 011 | Flutter Adapter Architecture | Accepted | 75% shared/25% separate, BcRenderer hybrid port done |
| 012 | Extract @bitcode/ui-schema | Accepted | zero-dep package, unblocks Flutter port |
| 013 | Cross-tier Infrastructure Completion | **Proposed** (2026-06-26) | AppProvider adoption, P0/P1 hooks, animation+adapt contract |
| 014 | Engine UI Feature Wiring | **Proposed** (2026-06-26) | 9 ERP UI features → component-backed SSR render |
| 015 | Flutter ECharts Chart Strategy | Accepted | deprecate CustomPaint, 26/26 parity |
| 016 | Demo App QA Bug Audit | Accepted (all 4 phases, 2026-06-29) | 28 stubs ported, 5 arch decisions |
| 017 | Block Catalog Architecture | **Accepted** (2026-06-30, matrix 0.1 baseline) | amends ADR-003 + ADR-005 |
| 018 | Block Catalog Data/Action/Reactivity | **Accepted** (2026-06-30) | required by ADR-017 before any Ready entry |

---

## Aturan ADR (rekap)

- ADR yang sudah Accepted **jangan diubah** — buat ADR baru yang supersedes/extends.
- Nomor urut berikutnya: **ADR-019**.
- Update tabel di `docs/adr/README.md` setiap ADR baru/status berubah.
- ADR-013/014/015 dirancang **paralel-safe** (lihat README paralelisasi note).
- Untuk "udah sampai mana komponen?" → baca `docs/components-progress.md`, BUKAN ADR.

---

## Log

<!-- Append-only: newest first. Format: YYYY-MM-DD | summary | entry ref -->
- 2026-07-01 | Agent L Block Catalog policy validators landed: 5 deterministic build-time policies (security, a11y, adaptive, performance, offline) over Agent I shape validator. Separate BLOCK_POLICY_* codes, no Ready promotion (Amendment 6). 80 policy tests pass, full suite 232/0, tsc clean. Commit 465826e1. | packages/ui-schema/src/block-catalog/policies/
- 2026-07-01 | Agent M Block Catalog Ready-promotion gate audit: NO promotion. P3 action runtime untracked+2 failing tests; P4 BcBlockRenderer/validator CLI missing. Future-proof gate test added (`packages/ui-schema/test/block-catalog/ready-gate.test.ts`, 8/8 pass). Smallest safe Ready slice = content.hero + content.plain-text (static-only, post-P4). Full evidence + unblock map in `sprints/package-ui-elementor/ready-promotion-audit.md`. | ready-promotion-audit.md
- 2026-07-01 | P2 Portable Expression AST landed in @bitcode/ui-schema (12 TS source files: ast, operators, scope, errors, permission-scope, evaluator iface, function-registry, builtin-manifest 70+ fns, golden-types, validate, serialize, index barrel) + Dart mirror in ui-flutter (10 files: ast manual sealed class, operators, errors, scope, permission-scope, function-registry, builtin-manifest, validate, evaluator, expressions barrel). 9 golden JSON fixtures (122 cases, BOM-stripped) as cross-platform source of truth. 32 TS tests pass (26 unit + 6 golden-fixture structural). 136 Dart tests pass (22 validate + 114 golden-parity). Security: FORBIDDEN_FUNCTION_NAMES denylist + path forbidden-segments + reserved-roots enforcement. Next: P3 ActionSource/ActionAdapter + runtime evaluator. | packages/ui-schema/src/expressions/, packages/ui-flutter/lib/src/expressions/
- 2026-06-30 | P1 Block Catalog foundation types landed in @bitcode/ui-schema: BlockCatalogEntry/CandidateCatalogEntry, enums (category/channel/capability/readiness/renderStrategy), Required* unions, 4 policies, migration/a11y/test metadata, capability normalization (action-trigger -> action-source), minimal shape validator (lenient Candidate/strict Ready + BLOCK_READY_GATE_INCOMPLETE), 2 Candidate fixtures (content.hero schema, interaction.notification-center component). tsc --noEmit passes. No Ready promotion. Next: P2 Portable Expression AST. | packages/ui-schema/src/block-catalog/, commit a2d4857
- 2026-06-30 | ADR-017 (Block Catalog) dan ADR-018 (Data/Action/Reactivity) Accepted setelah final read-through, matrix 0.1 baseline | docs/adr/017, 018
- 2026-06-30 | ADR-017 (Block Catalog) dan ADR-018 (Data/Action/Reactivity) diajukan, matrix draft 0.1 siap | docs/adr/017, 018
- 2026-06-29 | ADR-016 Demo App QA complete — 4 phases, 28 stubs ported, 5 arch decisions | docs/adr/016
- 2026-06-26 | ADR-013, 014, 015 diajukan (cross-tier infra, engine UI wiring, flutter echarts) | docs/adr/013-015
- 2026-06-22 | ADR-011, 012 Accepted (Flutter adapter + ui-schema standalone) | docs/adr/011, 012
