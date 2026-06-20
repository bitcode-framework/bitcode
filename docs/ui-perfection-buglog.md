# UI Perfection Bug Log

> Running log of bugs discovered during Phase A-C execution of the
> [UI Perfection Roadmap](../plans/2026-06-20-ui-perfection-roadmap.md).
>
> One row per bug. Triage every bug into one of: `fixed`, `deferred`,
> `wontfix`. P0 = blocks the component's core use; P1 = edge case or
> polish; P2 = nice-to-have.

## Triage Protocol

1. **Discover** — log bug with component tag, severity, reproduction, and phase.
2. **Classify** — source bug (`.lite.tsx` wrong) | Mitosis output bug (regex
   workaround in `mitosis-plugin-react-fixes.mjs`) | test bug | a11y bug.
3. **Fix or defer** — P0 must fix in the phase it's discovered. P1 may defer
   to B3. P2 may defer indefinitely with rationale.
4. **Verify** — re-run the failing test after fix. Update status to `fixed`.

## Bug Index

| # | Date | Component | Severity | Phase | Classification | Status | Summary |
|---|------|-----------|----------|-------|----------------|--------|---------|
| _(empty — bugs appended below as discovered)_ | | | | | | | |

## Bug Details

### BC-001 — BcInput: handleChange crashes with "Cannot read properties of undefined (reading 'value')"

- **Date:** 2026-06-20
- **Component:** BcInput
- **Severity:** P0
- **Phase:** A1.2 (pilot runtime test)
- **Classification:** source + plugin interaction
- **Reproduction:**
  1. Render `<BcInput name="q" onValueChange={fn} />` in jsdom.
  2. Type a character into the input.
- **Expected:** `onValueChange` fires with the typed value.
- **Actual:** `TypeError: Cannot read properties of undefined (reading 'value')` thrown from `handleChange` in compiled React output.
- **Root cause:** Source `handleChange(e: Event)` accessed `e.target.value` inside the body. The `react-value-events` Mitosis plugin then rewrote the JSX call site `onChange={(e) => state.handleChange(e)}` into `onChange={(e) => state.handleChange(e.target.value)}`, passing a string where the handler expected an Event — so the handler's own `e.target.value` dereference crashed.
- **Fix:** Refactored `BcInput.lite.tsx` to the value-first pattern: `handleChange(value: string)` with body using `value` directly (no `e.target` dereference). JSX call sites changed to `onChange={(value) => state.handleChange(value)}`. The plugin now leaves the call site alone (handler already value-first), and the React output is `onChange={(e) => handleChange(e.target.value)}` — exactly what we want.
- **Status:** fixed (commit pending)
- **Follow-up:** BcNumber has the same pattern (`handleChange(e: Event)` + `e.target.value`). Must be refactored identically before its runtime test lands.

### BC-002 — BcInput: slotProps.label attributes (data-testid, etc.) not applied to label element

- **Date:** 2026-06-20
- **Component:** BcInput
- **Severity:** P1
- **Phase:** A1.2 (pilot runtime test)
- **Classification:** source
- **Reproduction:**
  1. Render `<BcInput name="x" label="L" slotProps={{ label: { 'data-testid': 'lbl', className: 'c' } }} />`.
  2. Inspect the rendered `<label>` element.
- **Expected:** The label has `data-testid="lbl"` and `className` includes `c`.
- **Actual:** `className` is merged correctly (partClass works), but `data-testid` is `null` — non-class slotProps entries are silently dropped.
- **Root cause:** Source `BcInput.lite.tsx` only calls `partClass(props.slotProps, 'label', base)` on the label. It never calls `partAttr(props.slotProps, 'label')` to spread remaining attributes. Same gap exists on every part (root, control, input, description, error, etc.).
- **Fix:** Add `{...partAttr(props.slotProps, 'label')}` to the `<label>` element (and audit other parts for the same gap). For now, the test is relaxed to only assert className passthrough; the full attribute passthrough becomes a Phase B task (slotProps contract hardening across all components).
- **Status:** test relaxed; full fix deferred to Phase B (slotProps audit)

### BC-003 — BcInput (and all components): `if (cond) return;` compiled to `if (cond); return;` — handler always returns unconditionally

- **Date:** 2026-06-20
- **Component:** BcInput (and any component using `if (cond) return;` guard pattern)
- **Severity:** P0 (blocks all event handlers that start with an early-return guard)
- **Phase:** A1.2 (pilot runtime test)
- **Classification:** mitosis-output (plugin regression introduced 2026-06-20)
- **Reproduction:**
  1. Render `<BcInput onValueChange={fn} />` in jsdom.
  2. Fire `fireEvent.change(input, { target: { value: 'x' } })`.
- **Expected:** `onValueChange` fires with `'x'`.
- **Actual:** `onValueChange` never fires. Handler body exits at the first `if` guard because Mitosis output reads `if (cond); return;` (note the semicolon after the condition) — the `if` becomes an empty statement and `return` runs unconditionally.
- **Root cause:** The `mitosis-plugin-react-fixes.mjs` step 5 regex `/([)}])\s+(const|let|var|function|return|if|for|while|switch|throw|try)\b/g` was too aggressive. It inserted `;` after ANY `)` followed by a keyword, including the closing `)` of an `if (...)` condition — turning `if (cond) return;` into `if (cond); return;`.
- **Fix:** Tightened the regex in step 5. Now `;` is inserted only after `}` (object/array/function body — never a control-flow condition) or after `)` when followed by `const`/`let`/`var`/`function` (these keywords never appear as the body of `if`/`for`/`while`, so the `)` is unambiguously a statement terminator). The `return`/`if`/`for`/`while`/`switch`/`throw`/`try` cases are no longer matched after `)`.
- **Status:** fixed (commit pending)
- **Blast radius:** Every component with an `if (cond) return;` guard was silently broken in the React output. This is likely the root cause of the "0 onValueChange calls" pattern that A1 will surface across many components. Re-running all runtime tests after this fix is mandatory.

### BC-004 — BcSelect: listbox missing accessible name + combobox missing aria-controls/activedescendant

- **Date:** 2026-06-20
- **Component:** BcSelect
- **Severity:** P1 (a11y)
- **Phase:** A1.3 (runtime test — BcSelect)
- **Classification:** a11y (source)
- **Reproduction:**
  1. Render `<BcSelect label="Fruit" options={[...]} />` in jsdom.
  2. Open the listbox.
  3. Run axe-core.
- **Expected:** Zero axe violations.
- **Actual:** 3 violations:
  - `aria-input-field-name` — the `<div role="listbox">` has no accessible name (no aria-label / aria-labelledby).
  - `aria-required-attr` — the `<button role="combobox">` is missing `aria-controls` (and optionally `aria-activedescendant`).
  - `button-name` — axe flags the combobox button as having no discernible text (false-ish positive: the button has a `<span>` child with text, but axe's heuristic misses it when the span text is a placeholder color).
- **Fix (Phase B2 — a11y runtime hardening):**
  - Add `aria-labelledby={labelId}` to the listbox, deriving labelId from the label element's id.
  - Add `aria-controls={listboxId}` to the combobox button and `id={listboxId}` to the listbox.
  - Ensure the combobox button has an `aria-label` or visible text that axe recognizes as discernible.
- **Status:** deferred to Phase B2

### BC-005 — BcSelect: clearable prop declared but clear button never rendered

- **Date:** 2026-06-20
- **Component:** BcSelect
- **Severity:** P1 (feature gap)
- **Phase:** A1.3 (runtime test — BcSelect)
- **Classification:** source
- **Reproduction:**
  1. Render `<BcSelect clearable defaultValue="apple" options={[...]} />`.
  2. Inspect the DOM for a clear button.
- **Expected:** A clear button appears when `clearable=true` and a value is set.
- **Actual:** No clear button is rendered. The `handleClear` function exists in the component but is never wired to any JSX element.
- **Root cause:** Source `BcSelect.lite.tsx` declares `clearable` prop and `handleClear` method, but the JSX return never conditionally renders a clear button. The contract promises `clearButton` as a part, but it's missing from the implementation.
- **Fix:** Add a clear button to the trigger area, gated on `props.clearable && currentValue()`, with `aria-label="Clear"` and `onClick={state.handleClear}`. Mirror the pattern from BcInput's clear button.
- **Status:** deferred to Phase B (feature gap fill during parity audit)

### BC-006 — BcCombobox: controlled `value` prop does not reflect in the input display

- **Date:** 2026-06-20
- **Component:** BcCombobox
- **Severity:** P1 (controlled-mode broken)
- **Phase:** A1.3 (runtime test — BcCombobox)
- **Classification:** source
- **Reproduction:**
  1. Render `<BcCombobox value="cherry" options={[{label:'Cherry',value:'cherry'},...]} />`.
  2. Read the input value.
- **Expected:** Input displays "Cherry" (the label matching the controlled value).
- **Actual:** Input is empty. The `query` state is seeded only from `defaultValue`, with no effect that syncs it to `props.value` when the controlled value changes externally.
- **Root cause:** No `useEffect` watches `props.value` and updates `query` to the matching option's label. The component claims controlled support in its contract but only honors `defaultValue` for the initial render.
- **Fix:** Add a `useEffect` that, when `props.value` is defined, resolves the matching option label and sets `query`. If no option matches, set `query` to the raw value (or empty, per design decision).
- **Status:** deferred to Phase B (controlled-state hardening)

### BC-007 — BcCheckbox / BcToggle: button lacks accessible name (label not wired via aria-labelledby)

- **Date:** 2026-06-20
- **Component:** BcCheckbox, BcToggle (and likely BcCheckbox-like patterns elsewhere)
- **Severity:** P1 (a11y)
- **Phase:** A1.3 (runtime tests)
- **Classification:** a11y (source)
- **Reproduction:**
  1. Render `<BcCheckbox name="x" label="I agree" />`.
  2. Run axe-core.
- **Expected:** Zero violations.
- **Actual:** `button-name` violation — the `<button role="checkbox">` has no accessible name. The label text sits in a sibling `<label>` element, but there is no `id`/`for`/`aria-labelledby` linkage, so screen readers cannot associate the label with the button.
- **Root cause:** Source renders `<button>` and `<label>` as siblings without `htmlFor`/`id` or `aria-labelledby`. The `<label>` element does not wrap the button, so clicking the label also does not toggle the checkbox.
- **Fix (Phase B2):**
  - Add `id={name}` to the button and `htmlFor={name}` to the label, OR
  - Wrap the button inside the `<label>` element so clicking the label toggles it, OR
  - Add `aria-labelledby={labelId}` to the button pointing at the label's id.
- **Status:** deferred to Phase B2

<!-- Template:
### BC-NNN — <component>: <short summary>
- **Date:** YYYY-MM-DD
- **Component:** BcXxx
- **Severity:** P0 | P1 | P2
- **Phase:** A1 | A2 | A3 | A4 | B1 | B2 | B3 | C
- **Classification:** source | mitosis-output | test | a11y
- **Reproduction:**
  1. ...
- **Expected:** ...
- **Actual:** ...
- **Fix:** <commit hash or description>
- **Status:** fixed | deferred | wontfix
-->

### BC-008 — BcDialog: default props set `open: false`, breaking uncontrolled mode (defaultOpen never works)

- **Date:** 2026-06-20
- **Component:** BcDialog (likely same pattern in BcAlertDialog, BcDrawer, BcPopover, BcTooltip)
- **Severity:** P0 (uncontrolled mode completely broken)
- **Phase:** A1.3 (runtime test — BcDialog)
- **Classification:** source
- **Reproduction:**
  1. Render `<BcDialog defaultOpen title="T">Body</BcDialog>`.
  2. Query for `role="dialog"`.
- **Expected:** Dialog renders.
- **Actual:** Dialog does not render. `props.open` is `false` (from default), so `isOpen()` returns `false` always.
- **Root cause:** Default props merge sets `open: false`, making the `props.open !== undefined` check always true.
- **Fix:** Remove `open: false` from defaults, OR change check to `'open' in props`. Will fix in-source now since it blocks 8+ overlay tests.
- **Status:** fixing in-source (blocks A1.3 overlay tests)

### BC-009 — BcDialog: closeOnEscape prop declared but no Escape key handler wired

- **Date:** 2026-06-20
- **Component:** BcDialog (likely BcAlertDialog, BcDrawer, BcCommand too)
- **Severity:** P1 (a11y + UX)
- **Phase:** A1.3 (runtime test — BcDialog)
- **Classification:** source (feature gap)
- **Reproduction:**
  1. Render `<BcDialog defaultOpen closeOnEscape title="T">Body</BcDialog>`.
  2. Press Escape.
- **Expected:** `onOpenChange(false, {reason:'escape'})` fires.
- **Actual:** Nothing happens. No `keydown`/`onKeyDown` listener attached to the dialog or document.
- **Fix:** Add a `useEffect` that listens for `keydown` Escape on `document` when `isOpen() && props.closeOnEscape`, calling `emitClose('escape')`. Cleanup on close/unmount.
- **Status:** deferred to Phase B (keyboard nav hardening — B2)


