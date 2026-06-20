# UI Perfection Roadmap — Runtime Tests, Parity Audit, and Definition-of-Done

**Date:** 2026-06-20
**Status:** Proposed (awaiting review)
**Scope:** `packages/ui` (Mitosis component library) + parity with `packages/components` (Stencil)
**Related ADRs:** 001, 003, 004, 006, 008
**Out of scope (explicit):** Flutter target (deferred until Phase A-C complete)

---

## 1. Context — Why This Plan Exists

### 1.1 Honest Gap Analysis (2026-06-20)

Current state of `packages/ui`:

- **109 components** across P0/P1/P2 tiers have `.lite.tsx` source + `.contract.ts` + source-level tests.
- **1554 source-level tests pass**, but these are static-analysis tests (regex match against source string content). They assert patterns like `assert.match(src, /onInput=\{...\}/)` — they do NOT render the component or fire events.
- **Mitosis build green** for React/Vue/React Native/Stencil (151 components, 4 targets).
- **Demo app builds successfully** (375 modules, 1.54s) as of 2026-06-20 after fixing BcDataTable syntax errors, BcViewEditor literal-newline-in-string bug, BcViewGantt `<For>` callback bug, and adding a react-native web stub.

Critical gaps that the "Done" status in `docs/components-progress.md` does NOT reflect:

1. **No runtime tests.** Zero components have a test that mounts them in a DOM, fires an event, and asserts a state change. The 1554 tests are string-match assertions.
2. **No Stencil parity audit.** `packages/components` has 122 Stencil components with `@Prop`/`@Event`/`@Method` decorators that have been battle-tested in the BitCode engine. No side-by-side matrix exists to confirm the Mitosis reimplementations have feature parity.
3. **No E2E tests.** ADR-006 §5 DoD #6 requires Playwright (web) + Maestro/Detox (mobile) + Tauri webdriver smoke tests. None exist.
4. **No demo app browser verification.** The demo app builds but has never been click-tested in a browser. Bugs like BcDataTable syntax errors only surfaced when someone tried to build the demo app — implying more bugs remain hidden.
5. **4-layer data source inconsistent.** `local | fetcher | url | model` is declared on BcDataTable/BcList, but error/retry/abort/race paths are untested, and other data-bearing components (Combobox async, MultiSelect async, Cascader lazy, TreeSelect lazy, LookupModal, FieldDynLink, FieldTableSelect, Chart, Calendar, Kanban, Tree, OrgChart) have inconsistent coverage.
6. **ADR-001, ADR-003, ADR-004 still "Proposed"** — they were never promoted to Accepted because their acceptance criteria were never formally verified.

### 1.2 Definition of "Perfect" (per user, 2026-06-20)

This plan targets a concrete, measurable definition of "perfect" — the intersection of three criteria:

| # | Criterion | Source |
|---|-----------|--------|
| C1 | **Stencil parity 100%** — every `@Prop`, `@Event`, `@Method` in `packages/components` has a matching prop/event/method in the Mitosis contract, with documented behavior equivalence. | User request |
| C2 | **antd/Prime parity** — all enterprise features present in antd and PrimeReact for equivalent component categories are present in `packages/ui`. | User request |
| C3 | **ADR-006 §5 DoD fully satisfied** — contract + token + a11y + runtime smoke (React/Vue/RN) + E2E (Playwright web + Maestro mobile + Tauri webdriver). | ADR-006 |
| C4 (implicit) | **Runtime test coverage ≥80%** for Core Integration Surface (25 components). | ADR-008 spirit |
| C5 (implicit) | **0 P0 bugs** in issue tracker for Core Integration Surface. | Pragmatic gate |

### 1.3 Non-Goals

- **Flutter target.** Deferred. A separate ADR (ADR-011) will be written when Phase A-C complete. Effort estimate 3-6 months standalone.
- **New components beyond the 109.** Feature freeze until C1-C5 satisfied for existing 109.
- **Engine-side work** (Redis cache, NATS, Marketplace, Offline Mode Phase 3-5). Out of scope — this plan is UI-only.
- **BitCode JSON → UiSchema transformer** (ADR-005 Phase 4). Tracked separately in engine roadmap.

---

## 2. Phase Structure (Single-Track, Sequential)

```
Phase A (Stabilize)  ──>  Phase B (Harden)  ──>  Phase C (Formalize)  ──>  Phase D (Flutter, deferred)
   A1 Runtime tests         B1 Reactivity E2E       C1 Promote ADRs             D1 ADR-011
   A2 Stencil parity        B2 A11y runtime         C2 Parity matrix gate       D2 Token pipeline
   A3 Demo verification     B3 Edge cases           C3 DoD sign-off             D3 Generator spike
   A4 4-layer uniformity                                                        D4 Adapt decisions
```

**Single-track rationale:** Parallelizing would force context-switching between test infrastructure, audit tooling, and component fixes. User explicitly chose single-track (A2 after A1, not parallel) on 2026-06-20.

### 2.1 Effort Estimates (Honest, P50)

| Phase | Effort | Calendar | Notes |
|-------|--------|----------|-------|
| A1 | 4-6 weeks | 1.5 months | Infrastructure 1 week + 25 components × 1-2 days each |
| A2 | 3-5 weeks | 1 month | Script 1 week + matrix + gap fill |
| A3 | 1-2 weeks | 0.5 month | Manual click-through + bug triage |
| A4 | 2-3 weeks | 0.75 month | Audit + test error paths |
| B1 | 2-3 weeks | 0.75 month | Reactive scope wiring + tests |
| B2 | 2-3 weeks | 0.75 month | Keyboard + focus + axe-core |
| B3 | 4-8 weeks | 1.5-2 months | Variable, depends on bug count from A3 |
| C1-C3 | 1 week | 0.25 month | Documentation + ADR promotion |
| **Total A-C** | **19-31 weeks** | **5-8 months** | Realistic for 109 components |
| D (Flutter) | 3-6 months | 3-6 months | Separate ADR-011, deferred |

**Honest caveat:** These are P50 estimates. UI polish work historically runs 1.5-2x over estimate because hidden bugs compound. Plan assumes 1 engineer full-time; parallel agents reduce wall-clock but not total effort.

---

## 3. Phase A — Stabilize (Blocking Everything)

### Phase A1 — Runtime Test Foundation

**Goal:** Every Core Integration Surface component (25) has runtime tests that mount it in jsdom, fire real events, and assert real DOM state. Coverage ≥80% for these 25.

**Why blocking:** Without runtime tests, every later claim ("parity achieved", "a11y verified", "bug fixed") is unverifiable. This is the foundation for C3 (ADR-006 §5 DoD #5).

#### A1.1 Infrastructure Setup (Week 1)

| Task | Detail |
|------|--------|
| Add dev deps | `vitest`, `@testing-library/react`, `@testing-library/user-event`, `@testing-library/jest-dom`, `jsdom`, `@vitest/coverage-v8` |
| Config file | `packages/ui/vitest.config.ts` — environment `jsdom`, setup file `test-setup.ts`, coverage thresholds |
| Test script | `package.json`: `"test": "vitest run"`, `"test:watch": "vitest"`, `"test:coverage": "vitest run --coverage"` |
| Keep node --test | Existing source-level tests stay as `*.source.test.ts`; new runtime tests are `*.runtime.test.tsx`. Both run under vitest. |
| React alias | vitest config aliases `@bitcode/ui` → `output/react/src` (same as demo app) so tests hit compiled Mitosis output, not source |
| React-native stub | Reuse `demo-app/src/react-native-web-stub.tsx` — copy to `packages/ui/test/react-native-web-stub.tsx`, alias in vitest config |
| CI integration | Add `npm run test:coverage` to CI (when CI exists). For now, local pre-commit hook. |

**Acceptance A1.1:** `npx vitest run` executes (0 tests pass, 0 fail). Coverage report generates. `npx vitest run --coverage` reports 0% baseline.

#### A1.2 Smoke Test Template (Week 1, end)

Standard smoke test per component, in `packages/ui/src/components/<group>/<name>/<Name>.runtime.test.tsx`:

```tsx
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BcXxx } from '@bitcode/ui';

describe('BcXxx (runtime smoke)', () => {
  it('mounts and renders label', () => {
    render(<BcXxx label="My Field" />);
    expect(screen.getByText('My Field')).toBeInTheDocument();
  });

  it('fires onValueChange when user types', async () => {
    const user = userEvent.setup();
    const handleChange = vi.fn();
    render(<BcXxx onValueChange={handleChange} />);
    await user.type(screen.getByRole('textbox'), 'hello');
    expect(handleChange).toHaveBeenCalledWith('hello', expect.any(Object));
  });

  it('is accessible (axe)', async () => {
    const { container } = render(<BcXxx label="My Field" />);
    // axe-core check (setup in test-setup.ts)
    expect(await axe(container)).toHaveNoViolations();
  });
});
```

**Acceptance A1.2:** Template exists. One pilot component (BcInput) implements it and passes.

#### A1.3 Core Integration Surface Coverage (Weeks 2-6)

Implement runtime tests for all 25 Core Integration Surface components, in ADR-008 audit order:

| # | Component | Group | Critical interactions to test |
|---|-----------|-------|-------------------------------|
| 1 | BcInput | forms | type, clear, maxLength, disabled, readOnly, aria-invalid |
| 2 | BcSelect | forms | open, select, keyboard nav, controlled value |
| 3 | BcCombobox | forms | type-filter, pick, clear, async fetcher |
| 4 | BcCheckbox | forms | toggle, indeterminate, disabled |
| 5 | BcToggle | forms | toggle, role=switch, disabled |
| 6 | BcSegmented | forms | select, keyboard, sm/md/lg |
| 7 | BcNumberField | forms | type, ArrowUp/Down step, min/max clamp, locale parse |
| 8 | BcDateField | forms | type, ISO validation, aria-invalid |
| 9 | BcTimeField | forms | type, HH:mm validation, aria-invalid |
| 10 | BcMultiSelect | forms | add chip, remove chip, maxSelected, async |
| 11 | BcForm | views | submit, validate, reset, reactions, field visibility |
| 12 | BcReactiveScope | core | dispose, effect re-run, scope isolation |
| 13 | ReactiveStore | core | get/set, reactions, batch, transaction |
| 14 | Field wrappers | core | composition via BcForm parts |
| 15 | BcDialog | overlays | open/close, Esc dismiss, backdrop click, focus trap |
| 16 | BcAlertDialog | overlays | confirm/cancel, role=alertdialog, focus trap |
| 17 | BcDrawer | overlays | open/close, 4 placements, Esc, focus trap |
| 18 | BcPopover | overlays | open/close, trigger hover/click, placement |
| 19 | BcTooltip | overlays | hover/focus open, delay, placement |
| 20 | BcDropdown | overlays | trigger, item select, separator, disabled |
| 21 | BcContextMenu | overlays | right-click open, item select, long-press (RN) |
| 22 | BcCommand | overlays | type-filter, arrow nav, enter select, esc |
| 23 | BcTabs | navigation | tab switch, keyboard arrows, controlled activeKey |
| 24 | BcAccordion | layout | expand/collapse, single/multi, controlled |
| 25 | BcDisclosure | layout | expand/collapse, aria-expanded |
| 26 | BcSplitter | layout | drag handle, keyboard nudge, min sizes |
| 27 | BcDataTable | data | sort, filter, search, select, pagination, virtualize, tree, group, edit, export |
| 28 | BcList | data | 4-layer source, pagination, empty/error/loading |
| 29 | BcRenderer | core | 6 schema kinds, recursion, custom sentinel |
| 30 | AppProvider/ReactiveScope/FieldResolver/UiSchema | core | infrastructure contracts |

**Per-component minimum test set (DoD per component):**
- Mount + renders children/label
- Controlled value sync (`value` prop reflects in DOM)
- Uncontrolled default (`defaultValue` seeds state)
- Primary event fires with correct payload (`onValueChange`/`onOpenChange`/etc.)
- Disabled/readOnly blocks interaction
- Keyboard interaction (at least 1 key path: Enter/Escape/Arrow)
- A11y: role + aria-* attributes present; axe-core zero violations
- Error state (if data-bearing): error rendering, retry callback

**Acceptance A1.3:** All 25 components have runtime tests. Coverage ≥80% for these 25. `npm run test:coverage` passes. Zero failing tests.

#### A1.4 Bug Triage Protocol

Every test failure during A1.3 is triaged:
- **Source bug** (logic wrong in `.lite.tsx`) → fix source, verify both source-level + runtime test pass.
- **Mitosis output bug** (compiler generates broken code) → add regex workaround to `mitosis-plugin-react-fixes.mjs`, document in ADR-001 §14.
- **Test bug** (assertion wrong) → fix test.
- **A11y bug** (missing aria, wrong role) → fix source, add to contract.

All bugs logged in `docs/ui-perfection-buglog.md` (created during A1).

**Acceptance A1.4:** Bug log exists. Every discovered bug has either a fix-commit or a documented deferral with rationale.

---

### Phase A2 — Stencil Parity Audit

**Goal:** Produce `docs/stencil-parity-matrix.md` listing every `@Prop`/`@Event`/`@Method` in `packages/components/src/components/*` mapped to its Mitosis equivalent. Fill every gap.

**Why after A1:** Need runtime tests to verify parity claims — string-matching tests can't tell if behavior is equivalent.

#### A2.1 Audit Script (Week 1)

Write `packages/ui/scripts/audit-stencil-parity.mjs`:
- Scan `packages/components/src/components/**/*.tsx` for `@Prop()`, `@Event()`, `@Method()`, `@State`.
- Parse Stencil component name from class declaration.
- Map Stencil kebab-case tag (`bc-input`) → Mitosis component (`BcInput`).
- For each Stencil member, check if Mitosis `*.contract.ts` has a matching prop/event/method (camelCase equivalent).
- Output `docs/stencil-parity-matrix.md` with columns: Stencil member | Stencil type | Stencil default | Mitosis contract member | Status (✅/⚠️/❌) | Notes.

**Acceptance A2.1:** Script runs. Matrix generated. Manual review confirms correctness for 5 sample components.

#### A2.2 Gap Analysis (Week 2)

Review matrix, classify gaps:
- **Missing prop** — Stencil has `@Prop() customColor`, Mitosis contract doesn't.
- **Missing event** — Stencil has `@Event() bcCustomChange`, Mitosis doesn't.
- **Missing method** — Stencil has `@Method() doSomething()`, Mitosis doesn't.
- **Type mismatch** — Stencil `@Prop() value: string`, Mitosis `value?: string | number`.
- **Behavior mismatch** — Same prop but different default, different event payload, etc.

**Acceptance A2.2:** Every ❌ and ⚠️ in matrix has a triage decision: fill / defer with rationale / mark as intentional divergence.

#### A2.3 Gap Fill (Weeks 3-5)

For each "fill" decision:
1. Update `BcXxx.contract.ts` with missing prop/event/method.
2. Update `BcXxx.lite.tsx` to implement.
3. Add runtime test (from A1 template) for the new behavior.
4. Update `BcXxx.test.ts` source-level test if pattern changed.
5. Rebuild Mitosis, verify all 4 targets still compile.
6. Update matrix row to ✅.

**Acceptance A2.3:** Matrix has 0 ❌ for Core Integration Surface. ⚠️ rows all have documented rationale.

---

### Phase A3 — Demo App Browser Verification

**Goal:** Run demo app in browser, click through every route, log every bug, fix every P0 bug.

#### A3.1 Manual Test Pass (Week 1)

- `cd packages/ui/demo-app && npm run dev`
- Open `http://localhost:3002` in Chrome + Firefox + Safari.
- For each of the ~109 demo routes:
  - Click into the route.
  - Interact with the component (type, click, drag, etc.).
  - Verify no console errors.
  - Verify no visual regressions (compare to Stencil equivalent if exists).
  - Log bugs in `docs/ui-perfection-buglog.md`.

#### A3.2 Bug Fix (Week 2)

Fix every P0 bug (crash, broken interaction, wrong rendering). P1 bugs (visual polish, edge cases) logged for B3.

**Acceptance A3:** Bug log shows 0 open P0 bugs for demo app. P1 bugs triaged into B3.

---

### Phase A4 — 4-Layer Data Source Uniformity

**Goal:** Every data-bearing component supports `local | fetcher | url | model` with consistent error/retry/abort semantics.

#### A4.1 Audit (Week 1)

List components that should support data sources:
- BcDataTable, BcList, BcCombobox, BcMultiSelect, BcCascader, BcTreeSelect, BcLookupModal, BcFieldDynLink, BcFieldTableSelect, BcChart, BcCalendar, BcKanban, BcTree, BcOrgChart, BcTimeline, BcGantt.

For each, check:
- Does it accept `source?: DataSource`?
- Does `local` work? `fetcher`? `url`? `model`?
- Error state UI?
- Loading state UI?
- Empty state UI?
- Retry callback?
- Abort on unmount?
- Race condition handling (fast consecutive fetches)?

#### A4.2 Uniformity Fix (Weeks 2-3)

- Extract shared `useDataSource` hook (if not exists) in `src/core/data-source/`.
- Wire all data-bearing components to use it.
- Add runtime tests for error/retry/abort/race per component.

**Acceptance A4:** All data-bearing components pass 4-layer source tests. Error/retry/abort/race covered.

---

## 4. Phase B — Harden

### Phase B1 — Reactivity End-to-End

**Goal:** Cross-component reactive flows verified.

- BcDataTable selection → BcChart filter
- BcForm value → BcDataTable filter
- Conditional UI (`when`, `disabled_when`, `required_if`) in BcRenderer
- Cross-form `parent_field` / `set_value` patterns

**Acceptance B1:** 5+ end-to-end reactive scenario tests pass.

### Phase B2 — A11y Runtime Verification

**Goal:** Every interactive component passes axe-core + manual keyboard test.

- Keyboard nav: Tab/Shift+Tab/Arrow/Enter/Escape per component
- Focus trap in Dialog/Drawer/Popover/Command/AlertDialog
- Focus restore after close
- Screen reader landmark/live region behavior
- Touch target ≥44px on mobile (verify via CSS inspection)

**Acceptance B2:** axe-core zero violations for all 109 components. Manual keyboard test checklist complete.

### Phase B3 — Edge Case Hardening

**Goal:** Address every P1 bug from A3 + known edge cases.

Edge case categories:
- IME composition (CJK input) in Input/Combobox/TagsInput
- RTL ordering in DataTable columns, Kanban board, Breadcrumb separator
- Long content truncation (ellipsis/wrap/scroll)
- Empty/loading/error states for all data-bearing components
- Touch target sizes on mobile
- Reduced motion preference (`prefers-reduced-motion`)
- High contrast mode (`prefers-contrast: high`)
- SSR safety (no `window`/`document` access at module top-level)

**Acceptance B3:** Bug log shows 0 open P1 bugs for Core Integration Surface. Edge case test suite passes.

---

## 5. Phase C — Formalize

### Phase C1 — ADR Promotion

- **ADR-001** Proposed → Accepted. Acceptance criteria: Mitosis 0.13.1 producing 4-target output verified, plugin system documented, Stencil/Vue/RN targets green.
- **ADR-003** Proposed → Accepted. Acceptance criteria: P0/P1/P2 taxonomy applied to folder structure (verified), Stencil parity matrix complete (from A2), antd/Prime parity checklist complete.
- **ADR-004** Proposed → Accepted. Acceptance criteria: ReactiveStore runtime tests pass (from B1), cross-component reactive scenarios verified.

### Phase C2 — Parity Matrix Gate

- `docs/stencil-parity-matrix.md` (from A2) frozen as formal gate.
- `docs/antd-prime-parity-matrix.md` created — checklist of antd/PrimeReact enterprise features per component category, mapped to `packages/ui` status.

### Phase C3 — DoD Sign-Off

- Update `docs/components-progress.md` — replace "Done" historical claims with "Done (Phase C verified)" only for components that pass C1-C3.
- Update `docs/components-contract-matrix.md` — Tests gate now requires runtime test coverage ≥80%, not just source-level.
- Commit ceremony: tag `packages/ui@v1.0.0-rc1`.

---

## 6. Phase D — Flutter (Deferred, Separate ADR)

Triggered only after Phase A-C complete. Will produce:

- **ADR-011 Flutter Target Strategy** — decision on Mitosis custom generator vs hand-port, token pipeline extension, per-component Adapt decisions (3-way: web/native/flutter).
- **Spike:** 1 component (BcButton or BcStack) end-to-end in Flutter to validate generator approach.
- **Token pipeline:** `generate-tokens.mjs` emits Flutter `ThemeData`/`ColorScheme` Dart file.
- **Per-component:** Flutter widget mapping (View→Container, Pressable→InkWell, Text→Text, etc.).

**Effort:** 3-6 months. Not started until A-C done.

---

## 7. Tracking Artifacts

| Artifact | Purpose | Created In |
|----------|---------|-------------|
| `docs/plans/2026-06-20-ui-perfection-roadmap.md` | This document | Now |
| `docs/ui-perfection-buglog.md` | Running bug log | A1.1 |
| `packages/ui/vitest.config.ts` | Test config | A1.1 |
| `packages/ui/test-setup.ts` | Test setup (axe, jest-dom matchers) | A1.1 |
| `packages/ui/test/react-native-web-stub.tsx` | RN stub for jsdom | A1.1 |
| `docs/stencil-parity-matrix.md` | Stencil vs Mitosis parity | A2.1 |
| `docs/antd-prime-parity-matrix.md` | antd/Prime parity checklist | C2 |
| `packages/ui/scripts/audit-stencil-parity.mjs` | Audit script | A2.1 |
| Updated `docs/components-progress.md` | Honest status | C3 |
| Updated `docs/components-contract-matrix.md` | Tests gate tightened | C3 |

---

## 8. Commit & Review Cadence

- **Per-component (A1.3, A2.3):** One commit per component — `test(runtime): BcXxx runtime tests` or `fix(parity): BcXxx missing X prop from Stencil`.
- **Per-phase milestone:** Update `docs/components-progress.md` + this roadmap with status. Commit `docs: Phase A1 complete`.
- **Phase boundary:** Create handoff file in `handoffs/YYYY-MM-DD-ui-perfection-phase-X.md` summarizing what's done, what's next, gotchas.

---

## 9. Risk Register

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Mitosis 0.13.1 has more output bugs discovered during A1 runtime tests | High | Medium | Each bug → regex workaround in `mitosis-plugin-react-fixes.mjs` + document in ADR-001 §14. Upgrade Mitosis version if bugs exceed 20. |
| Stencil parity gaps reveal deep architectural mismatches (e.g., Stencil uses Shadow DOM, Mitosis doesn't) | Medium | High | Document as "intentional divergence" in parity matrix. ADR-006 already establishes Mitosis as primary; Stencil is legacy frozen. |
| Runtime tests reveal ADR-008 "Done" claims were false | High | High | Honest acknowledgment in `docs/components-progress.md`. Re-baseline status. This is expected — the whole point of A1. |
| Effort exceeds estimate (19-31 weeks → 40+ weeks) | Medium | High | Weekly progress check against estimate. If >1.5x over, pause and re-scope with user. |
| Demo app reveals bugs in components not in Core Integration Surface (P1/P2) | High | Medium | Log all, fix P0 in A3, defer P1 to B3. Don't expand A1 scope beyond 25 components. |
| User wants Flutter earlier than A-C complete | Low | High | Defer decision. This plan explicitly sequences Flutter after A-C. Re-open conversation if user pushes. |

---

## 10. Open Questions for Review

1. **Coverage threshold:** 80% for Core Integration Surface — acceptable, or push to 90%? (80% is industry standard for UI libs; 90% adds ~30% effort for diminishing returns.)
2. **antd vs PrimeReact as primary parity reference?** antd has more enterprise features; PrimeReact is closer to BitCode's existing Stencil design. Recommend: antd as reference for feature checklist, PrimeReact as reference for API ergonomics.
3. **E2E framework choice:** Playwright (web) is decided. Mobile: Maestro (newer, simpler) vs Detox (older, RN-native). Tauri: Playwright + tauri-driver. Recommend Maestro for lower setup cost.
4. **Should the 5 P0 components not in Core Integration Surface (BcSlot, BcPortal, BcVisuallyHidden, BcBackdrop, BcLoadingOverlay, BcToastQueue, BcStack family) get A1 runtime tests too?** Recommend yes — they're foundational primitives used by everything else.
5. **Bug log granularity:** One file `docs/ui-perfection-buglog.md` with all bugs, or per-component? Recommend single file with component tag column for filterability.

---

## 11. Next Action (Pending Approval)

Upon user approval of this plan:

1. Create `docs/ui-perfection-buglog.md` (empty template).
2. Install vitest dev deps in `packages/ui`.
3. Create `packages/ui/vitest.config.ts` + `test-setup.ts` + `test/react-native-web-stub.tsx`.
4. Update `package.json` test scripts.
5. Implement pilot runtime test for BcInput.
6. Verify `npx vitest run` works and BcInput test passes.
7. Commit `test(runtime): Phase A1 infrastructure + BcInput pilot`.
8. Proceed to next component (BcSelect) per A1.3 order.

**Pause point after step 7** — user reviews pilot before continuing with remaining 24 components.
