# Block Catalog Ready-Promotion Gate Audit

Audit date: 2026-07-01
Owner: Agent M (readiness-metadata + gate-evidence scope)
ADR baseline: ADR-017 / ADR-018 Accepted; matrix v0.1.0 `accepted-baseline`
Related: `REVIEW-FINDINGS.md`, `block-catalog-candidate-matrix.json`, `docs/plans/block-catalog-implementation.md`

## Verdict

**No Block Catalog component can be promoted to `readiness: 'ready'` in this session.**

Two foundational pieces are missing or unverified, and the ADR-017 §4 Ready Gate cannot be enforced end-to-end without them:

| Blocker | State | Evidence |
|---|---|---|
| **P3 Action runtime** (ADR-018 S8) | UNTRACKED + 2 failing tests | `packages/ui-schema/src/actions/` and `packages/ui-schema/test/actions/` are untracked. `runner.test.ts` A4 (malformed-result normalization: `200 !== 'error'`) and E14 (`resolveFirst is not a function`) fail. Not merged, not verified. |
| **P4 BcBlockRenderer + validator CLI** | DOES NOT EXIST | `rg 'BcBlockRenderer\|catalog:validate\|validate-block-catalog'` across `packages/` returns zero implementation matches. Only the P1 shape validator exists, which explicitly emits `BLOCK_READY_GATE_INCOMPLETE` for every Ready entry: *"cross-package gates land in P4."* |

The P1 shape validator (`packages/ui-schema/src/block-catalog/validate.ts`) confirms this: any Ready entry receives an `info`-severity `BLOCK_READY_GATE_INCOMPLETE` diagnostic stating the full cross-package gate is enforced by the P4 catalog validator. Promoting a block now would emit that reminder and silently rely on gates 2, 3, 7-13 being unenforced.

Per the implementation plan §1 rule #10: *"If a block fails any Ready gate, keep it Candidate and record the blocker in the matrix/report. Do not weaken gates."* This audit records the blockers.

## Smallest safe Ready slice (post-P4)

**`content.hero` + `content.plain-text`** are the first promotable blocks once P4 lands.

Rationale: their matrix capabilities are `static`/`responsive`/`i18n`/`a11y` only. They have **no** `action-source`, **no** `data-source`, **no** `reactive`, **no** `permission-aware`, **no** `large-data`, **no** `payment-adapter`. Therefore they are NOT blocked by the failing P3 action runtime (ADR-018), and they only require the static-content subset of the P4 gates (renderer, web/Flutter export, demo, a11y).

## Gate-evidence table — `content.hero`

Source/export evidence (exact, file:line):

| Platform | Component | Export line | Source file |
|---|---|---|---|
| Web | `BcSection` | `packages/ui/src/index.ts:137` | `components/layout/section/BcSection.lite.tsx` |
| Web | `BcStack` | `packages/ui/src/index.ts:140` | `components/layout/stack/BcStack.lite.tsx` |
| Web | `BcImage` | `packages/ui/src/index.ts:69` | `components/data-display/image/BcImage.lite.tsx` |
| Web | `BcButton` | `packages/ui/src/index.ts:192` | `components/primitives/button/BcButton.lite.tsx` |
| Flutter | `bc_stack` | `bc_ui_flutter.dart` exports `src/components/bc_stack.dart` | `packages/ui-flutter/lib/src/components/bc_stack.dart` |
| Flutter | `bc_image` | `bc_ui_flutter.dart` exports `src/components/bc_image.dart` | `packages/ui-flutter/lib/src/components/bc_image.dart` |
| Flutter | `bc_native_button` | `bc_ui_flutter.dart` exports `src/components/bc_native_button.dart` | `packages/ui-flutter/lib/src/components/bc_native_button.dart` |

ADR-017 §4 gate status:

| # | Gate | Status | Blocker |
|---|---|---|---|
| 1 | Matrix gate | PASS | decision `keep`, status `schema-preset-only`, no unresolved alias |
| 2 | Web component gate | PASS | all 4 web components exported (see table above) |
| 3 | Flutter component gate | PARTIAL | `bc_section` is not a distinct Flutter widget; matrix maps Flutter to `bc_stack` (verified exported). Verify preset composes only exported Flutter widgets. |
| 4 | Schema gate | PASS (shape) | `candidateContentHero` fixture shape-validates; full serializable preset lands at promotion |
| 5 | Action gate | N/A | capabilities are static — no action-source |
| 6 | Expression gate | N/A | no reactive capability |
| 7 | RuntimeContext gate | BLOCKED | no `BcBlockRenderer` (P4) to validate required functions/adapters |
| 8 | Security/privacy gate | BLOCKED | no P4 validator to enforce URL/sanitize policy checks |
| 9 | Adaptive gate | BLOCKED | no P4 validator; adaptivePolicy must be declared at promotion |
| 10 | Performance gate | N/A | static content — no large-data |
| 11 | A11y gate | BLOCKED | no P4 validator + no a11y test artifact linked |
| 12 | Test gate | BLOCKED | no P4 validator; web/Flutter render tests not linked |
| 13 | Demo gate | BLOCKED | no demo route in `packages/ui/demo-app` / `packages/ui-flutter/example` |

**Ready today?** NO — gates 3, 7, 8, 9, 11, 12, 13 blocked.

## Gate-evidence table — `content.plain-text`

Source/export evidence (exact, file:line):

| Platform | Component | Export line | Source file |
|---|---|---|---|
| Web | `BcLabel` | `packages/ui/src/index.ts:197` | `components/primitives/label/BcLabel.lite.tsx` |
| Web | `BcSection` | `packages/ui/src/index.ts:137` | `components/layout/section/BcSection.lite.tsx` |
| Flutter | `bc_box` | `bc_ui_flutter.dart` exports `src/components/bc_box.dart` | `packages/ui-flutter/lib/src/components/bc_box.dart` |

ADR-017 §4 gate status: identical blocking pattern to `content.hero` (gates 7, 8, 9, 11, 12, 13 blocked by missing P4). Additionally, matrix gap notes *"Need text node support or canonical text primitive in BlockCatalog schema"* — the schema preset must be finalized at promotion.

**Ready today?** NO.

## Components rejected (blocked), grouped by blocker

### Blocked by P3 action runtime (ADR-018 S8) — capabilities include action-source/action-trigger
Every candidate whose capabilities include `action-source` (matrix alias `action-trigger`) or `reactive` is blocked until P3 ActionRunner is merged and verified:
- `content.cta`, `content.pricing-table`
- `interaction.form`, `interaction.multi-step-wizard`, `interaction.record-picker`, `interaction.import-wizard`, `interaction.notification-center`
- `workflow.approval-actions`, `workflow.chatter`
- `commerce.product-detail`, `commerce.cart-checkout`, `commerce.payment-action`, `commerce.review`

### Blocked by P4 BcBlockRenderer / validator (every Ready block)
All 42 non-excluded candidates are blocked by the absence of P4 (gates 7, 8, 9, 11, 12, 13 cannot be enforced).

### Blocked by Flutter parity widgets
Resolved partially by commit `f6b2287` (`bc_avatar`, `bc_badge`, `bc_kpi`, `bc_progress` + tests + barrel exports added). Remaining validation: P4 export-check must run against these before the dependent candidates (`content.testimonial`, `content.team-member`, `data.account-summary`, `workflow.state-indicator`, `analytics.kpi-grid`, `analytics.chart`) can be considered.

### Blocked by web `BcMarkdown`
`content.rich-text` requires canonical web `BcMarkdown`. Web source now exists (Agent D, `packages/ui/src/components/content/markdown/`), exported from `packages/ui/src/index.ts`. Sanitize policy + P4 validator still required.

### Blocked by web export/registry bugs
- `content.social-links` — `BcShare` source exists; web export verified at `packages/ui/src/index.ts` (`export { default as BcShare }`). Registry generation (commit `b55aeec`) included it. **Export gate unblocked;** remaining blockers: URL safety policy, provider icon strategy, P4 validator.
- `data.record-list`, `commerce.product-grid` — `BcGridView` web export verified at `packages/ui/src/index.ts` (`export { default as BcGridView }`). **Export gate unblocked;** remaining: datasource conventions, performance budget, P4 validator.

### Blocked by component-gap (new component needed)
- `content.countdown` (needs `BcCountdown`)
- `interaction.import-wizard` (needs import block + file parsing adapter)
- `workflow.diagram-editor` (needs authoring component)
- `builder.query-builder`, `builder.expression-builder`, `builder.page-builder`, `builder.report-builder`
- `commerce.cart-checkout` (needs dedicated cart/checkout components)
- `data.pivot-table` (needs pivot component/contract)

### Blocked by contract-gap
- `interaction.record-picker` (unified picker contract)
- `workflow.state-indicator` (resolved: Flutter BcBadge/BcProgress added; distinct state component only if semantics require)

### Blocked by ADR-014 engine wiring (out of catalog scope)
Endpoint/model/permission/realtime/persistence wiring for: `workflow.chatter`, `data.audit-log`, `interaction.notification-center`, `data.pivot-table`, `builder.report-builder`. Catalog owns generic contract only.

### Excluded / alias
- `plugin.embed-social` — `duplicate-alias`, provider plugin not core Ready.
- `excluded.wordpress-specific` — `excluded`.

## Unblock map (blocker → resolving batch)

| Blocker | Batch | Exit signal |
|---|---|---|
| P3 ActionRunner failing/untracked | P3 | `test/actions/runner.test.ts` A4 + E14 fixed; actions dir committed; `node scripts/run-tests.mjs` green |
| P4 BcBlockRenderer + validator CLI | P4 | `catalog:validate` CLI exists; `BcBlockRenderer` schema+component strategies implemented; registries exist |
| Demo routes | P4/P5+ | `packages/ui/demo-app` + `packages/ui-flutter/example` block catalog sections |
| A11y test artifacts | P5+ | a11y test file linked in entry `tests.a11y` / `a11yEvidence.testReference` |
| Flutter widget render tests | P5+ | widget test per Ready block under `packages/ui-flutter/test/` |
| ADR-014 engine wiring | ADR-014 | endpoint/model/permission wiring merged (out of catalog plan scope) |
| Component-gap components | P6-P11 | dedicated `Bc*` component exists on web + Flutter |

## Future-proof gate enforcement (this session)

Added `packages/ui-schema/test/block-catalog/ready-gate.test.ts` (8 tests, all pass). Contract enforced: *"Every Ready catalog entry must pass strict shape validation."* The test is NOT "Ready count must be 0 forever" — when a real Ready promotion lands, the test must keep passing, which is the regression value. Today it holds vacuously (zero Ready fixtures) and also proves the shape validator has teeth (rejects Ready entries missing required fields: `requiredComponents`, `a11yEvidence`, `tests`, `schemaPreset`, `componentKey`).

When P4 ships the cross-package validator, this test extends to also require the P4 validator result as evidence — but the contract ("Ready ⇒ complete evidence") stays the same.

## Verification (this session)

| Command | Result |
|---|---|
| `cd packages/ui-schema; npx tsc --noEmit` | EXIT 0 (new test type-checks) |
| `cd packages/ui-schema; node --test test/block-catalog/ready-gate.test.ts` | 8 pass, 0 fail |
| `rg "readiness: 'ready'" packages/ui-schema/src packages/ui/src/block-catalog` | 0 matches (no premature Ready claim) |

Pre-existing failures (out of Agent M scope, reported not fixed): `packages/ui-schema/test/actions/runner.test.ts` A4 + E14 — untracked Agent E/F action-runtime work.

## Files changed (this session)

- `sprints/package-ui-elementor/ready-promotion-audit.md` (this file) — audit artifact
- `packages/ui-schema/test/block-catalog/ready-gate.test.ts` — future-proof gate test
- `.multibrain/buckets/ui-library.md` — memory entry
- `.multibrain/buckets/adr.md` — memory entry
