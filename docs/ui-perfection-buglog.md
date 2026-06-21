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
| BC-001 | 2026-06-20 | BcInput | P0 | A1.2 | source+plugin | fixed | handleChange crash (value-first pattern) |
| BC-002 | 2026-06-20 | BcInput | P1 | A1.2 | source | fixed | slotProps.label attributes wired via partAttr across 30 components (data-testid + aria-label) |
| BC-003 | 2026-06-20 | BcInput+all | P0 | A1.2 | mitosis-output | fixed | `if (cond) return;` → `if (cond); return;` |
| BC-004 | 2026-06-20 | BcSelect | P1 | A1.3 | a11y | fixed | listbox aria-label + combobox aria-controls + listbox id |
| BC-005 | 2026-06-20 | BcSelect | P1 | A1.3 | source | fixed | clearable prop — clear button now rendered and wired |
| BC-006 | 2026-06-20 | BcCombobox | P1 | A1.3 | source | fixed | controlled value sync via onMount |
| BC-007 | 2026-06-20 | BcCheckbox/Toggle | P1 | A1.3 | a11y | fixed | button aria-label + label htmlFor |
| BC-008 | 2026-06-20 | BcDialog | P0 | A1.3 | source | fixed | default props set `open: false`, breaking uncontrolled mode |
| BC-009 | 2026-06-20 | BcDialog | P1 | A1.3 | source | fixed | closeOnEscape — keydown listener via onMount/onUnMount |
| BC-010 | 2026-06-20 | BcPopover | P1 | A1.3 | mitosis-output | fixed | hover trigger — if-guard replaces ternary IIFE |
| BC-011 | 2026-06-20 | BcPopover | P1 | A1.3 | mitosis-output | fixed | named slot — triggerContent render-prop replaces Slot |
| BC-012 | 2026-06-20 | BcContextMenu | P1 | A1.3 | a11y | fixed | separator role=separator (not menuitem) |
| BC-013 | 2026-06-20 | BcCommand | P0 | A1.3 | source+plugin | fixed | search filtering — value-first handler |
| BC-014 | 2026-06-20 | BcCommand | P1 | A1.3 | a11y | fixed | dialog aria-label, listbox aria-label, combobox aria-controls |
| BC-015 | 2026-06-20 | BcCommand | P1 | A1.3 | source | fixed | stale activeIndex — local var computation |
| BC-016 | 2026-06-20 | BcAccordion | P1 | A1.3 | source | fixed | onToggle — compute nextIds before passing to callback |
| BC-017 | 2026-06-20 | BcSplitter | P2 | A1.3 | test+source | fixed | (a) jsdom pointer-event workaround via MouseEvent constructor, (b) stale-dragging fixed with useRef |
| BC-018 | 2026-06-20 | BcForm | P1 | A1.3 | source+plugin | fixed | select onChange — e.target.value extraction |
| BC-019 | 2026-06-20 | BcForm | P1 | A1.3 | source | fixed | stale-state closure in setField — validate against nextValues |
| BC-020 | 2026-06-20 | BcList | P1 | A1.3 | source | fixed | local mode — local newTotal variable |
| BC-021 | 2026-06-20 | BcRenderer | P1 | A1.3 | mitosis-output | fixed | data-bc-renderer-custom — whitelist in plugin |
| BC-022 | 2026-06-20 | BcDataTable | P0 | A1.3 | mitosis-output | fixed | defaultGetRowId — module scope + can() to useStore + IIFE fixes |
| BC-023 | 2026-06-21 | BcDataTable | P1 | A1.3 | source | fixed | props.loading wired in onMount |
| BC-024 | 2026-06-21 | BcReactiveScope | P1 | A1.3 | source | fixed | (a) test-timing fixed with waitFor poll; (b) cleanup memory leak fixed with useRef + unsubscribe; pattern also applied to BcForm + AppProvider |
| BC-025 | 2026-06-22 | BcBackdrop | P0 | B2 | source | open | `opacityClass` ReferenceError at render — Mitosis codegen drops getter |
| BC-026 | 2026-06-22 | BcFieldDuration | P0 | B2 | source | open | `inputClass` ReferenceError at render — Mitosis codegen drops getter |
| BC-027 | 2026-06-22 | BcFieldDynLink | P0 | B2 | source | open | `modelDisabled` ReferenceError at render — Mitosis codegen drops getter |
| BC-028 | 2026-06-22 | BcFieldMorph | P0 | B2 | source | open | `typeDisabled` ReferenceError at render — Mitosis codegen drops getter |
| BC-029 | 2026-06-22 | BcFieldTableSelect | P0 | B2 | source | open | `tableWrapClass` ReferenceError at render — Mitosis codegen drops getter |
| BC-030 | 2026-06-22 | BcLoadingOverlay | P0 | B2 | source | open | `visible` ReferenceError at render — Mitosis codegen drops getter |
| BC-031 | 2026-06-22 | BcViewGantt | P0 | B2 | source | open | `chartH` ReferenceError at render — Mitosis codegen drops getter |
| BC-032 | 2026-06-22 | BcButton | P1 | B2 | a11y | open | button-name axe violation — missing accessible name |
| BC-033 | 2026-06-22 | BcCode/BcJson | P1 | B2 | a11y | open | label axe violation — missing accessible name |
| BC-034 | 2026-06-22 | BcTextarea/BcNumber/BcMentions | P1 | B2 | a11y | open | label axe violation — missing accessible name |
| BC-035 | 2026-06-22 | BcDateTime | P1 | B2 | a11y | open | label axe violation — missing accessible name |
| BC-036 | 2026-06-22 | BcFile | P1 | B2 | a11y | open | label axe violation — missing accessible name |
| BC-037 | 2026-06-22 | BcSlider | P1 | B2 | a11y | open | label axe violation — missing accessible name |
| BC-038 | 2026-06-22 | BcRichText | P1 | B2 | a11y | open | aria-input-field-name axe violation |
| BC-039 | 2026-06-22 | BcDropZone | P1 | B2 | a11y | open | aria-allowed-role + nested-interactive axe violations |
| BC-040 | 2026-06-22 | BcColorArea/BcColorWheel/BcKnob | P1 | B2 | a11y | open | nested-interactive axe violations |
| BC-041 | 2026-06-22 | BcMeterGroup | P1 | B2 | a11y | open | aria-required-attr axe violation |
| BC-042 | 2026-06-22 | BcSearch | P1 | B2 | a11y | open | aria-allowed-attr axe violation |
| BC-043 | 2026-06-22 | BcSeparator | P1 | B2 | a11y | open | aria-allowed-attr axe violation |
| BC-044 | 2026-06-22 | BcTimelineSocial | P1 | B2 | a11y | open | aria-required-children axe violation |
| BC-045 | 2026-06-22 | Media viewers (6) | P2 | B2 | test | wontfix | axe-core cannot analyze iframes in jsdom — verify in A3 browser |

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
- **Root cause:** Source `BcInput.lite.tsx` imports `partAttr` but never calls it. Only `partClass(props.slotProps, 'label', base)` is used on the label element. Non-class slotProps entries (data-*, aria-*, etc.) are silently dropped. Same gap exists on every part (root, control, input, description, error, etc.).
- **Correction (2026-06-21 grilling session):** The original bug log stated "Mitosis can't spread `{...partAttr()}`" and classified this as an architecture change. This was inaccurate. ADR-009 (Runtime Part Resolution) documents that dynamic object spread (`{...props.slotProps?.root}`) is not supported by Mitosis across all targets (Stencil, Vue, Angular, Flutter). The design solution already exists: `partAttr(slotProps, partName, attrName)` is a single-attribute accessor that avoids spread entirely. BcInput simply has not wired `partAttr` calls for each supported attribute on each part. This is a source-level feature gap (missing wiring), not an architecture limitation.
- **Fix:** Wire `partAttr` calls explicitly per supported attribute per part. Example for label: `data-testid={partAttr(props.slotProps, 'label', 'data-testid')}`, `aria-label={partAttr(props.slotProps, 'label', 'aria-label')}`. The supported attribute categories per ADR-009 are: className (via partClass), data-*, aria-*, style (via partStyle). This is a Phase B audit task across all components, not an architecture change.
- **Status:** fixed (2026-06-21). partAttr wiring applied to 30 components (BcInput pilot + 29 via batch audit). Each `data-part="X"` element now wires `data-testid={partAttr(props.slotProps, 'X', 'data-testid') as string}`. Input/button/label parts also wire `aria-label`. Build + 1558 source + 476 runtime tests pass.

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

### BC-010 — BcPopover: hover trigger mode crashes on click (Mitosis compiles ternary to broken IIFE)

- **Date:** 2026-06-20
- **Component:** BcPopover
- **Severity:** P1 (hover-trigger mode broken in React output)
- **Phase:** A1.3 (runtime test — BcPopover)
- **Classification:** mitosis-output
- **Reproduction:**
  1. Render `<BcPopover trigger="hover">Body</BcPopover>`.
  2. Click the trigger button.
- **Expected:** Nothing happens (hover mode, click is a no-op).
- **Actual:** `TypeError: (props.trigger === 'click' ? () => toggle() : undefined)()` — calling `undefined()` throws.
- **Root cause:** Source `onClick={props.trigger === 'click' ? () => state.toggle() : undefined}` compiles to `onClick={event => (props.trigger === 'click' ? () => toggle() : undefined)()}`. The ternary is wrapped in an IIFE that invokes the result — when the result is `undefined`, calling it throws.
- **Fix:** Refactor source to avoid ternary in JSX event handler. Use `onClick={props.trigger === 'click' ? () => state.toggle() : undefined}` as a direct assignment (not an IIFE), or use a named handler that checks the trigger mode internally. May require a Mitosis plugin fix.
- **Status:** deferred to Phase B (Mitosis plugin hardening)

### BC-011 — BcPopover: named slot `<Slot name="trigger" />` compiles to `props.trigger` (the string prop)

- **Date:** 2026-06-20
- **Component:** BcPopover
- **Severity:** P1 (trigger content never renders)
- **Phase:** A1.3 (runtime test — BcPopover)
- **Classification:** mitosis-output
- **Reproduction:**
  1. Render `<BcPopover>Body</BcPopover>` with trigger slot content.
  2. Inspect the trigger button.
- **Expected:** The trigger button shows the slot content.
- **Actual:** The trigger button shows the literal string "click" or "hover" (the `trigger` prop value). Mitosis compiles `<Slot name="trigger" />` to `props.trigger` (the prop), not the named slot children.
- **Root cause:** Mitosis does not properly resolve named slots in the React target. `<Slot name="trigger" />` is compiled to `props.trigger` instead of extracting the named slot from `props.children`.
- **Fix:** Either avoid named slots in Mitosis components (use a `trigger` render-prop instead), or fix the Mitosis slot compilation for the React target. This affects any component using `<Slot name="..." />`.
- **Status:** deferred to Phase B (Mitosis slot compilation fix)

### BC-012 — BcContextMenu: separator items wrapped in role="menuitem" (should be role="separator")

- **Date:** 2026-06-20
- **Component:** BcContextMenu
- **Severity:** P1 (a11y)
- **Phase:** A1.3 (runtime test — BcContextMenu)
- **Classification:** a11y (source)
- **Reproduction:**
  1. Render `<BcContextMenu items={[...]} />` with at least one separator item.
  2. Open the context menu.
  3. Run axe-core.
- **Expected:** Zero violations.
- **Actual:** `aria-command-name` violation — the `<div role="menuitem">` wrapper around the separator has no accessible name. Separators should use `role="separator"` and should not be wrapped in `role="menuitem"`.
- **Root cause:** Source applies `role="menuitem"` to the outer `<div>` for ALL items (including separators) via `<div key={item.value} role="menuitem">`. Separators get the wrong role.
- **Fix:** Conditionally apply `role="menuitem"` only to non-separator items, and use `role="separator"` for separator items. Or restructure so separators are not wrapped in a menuitem div.
- **Status:** deferred to Phase B2 (a11y runtime hardening)

### BC-013 — BcCommand: search input filtering broken (handler expects Event, receives string)

- **Date:** 2026-06-20
- **Component:** BcCommand
- **Severity:** P0 (core feature broken — search/filter does not work)
- **Phase:** A1.3 (runtime test — BcCommand)
- **Classification:** source + plugin interaction (same pattern as BC-001)
- **Reproduction:**
  1. Render `<BcCommand open items={[...]} />`.
  2. Type into the search input.
- **Expected:** The item list filters to match the query.
- **Actual:** The query state never updates. All items remain visible regardless of what is typed. Typing "zzz" does not show the empty state.
- **Root cause:** Source `handleInput(e: any)` reads `e?.target?.value`. The `react-value-events` Mitosis plugin rewrites `onInput={(e) => state.handleInput(e)}` to `onChange={e => handleInput(e.target.value)}`, passing a string. The handler then reads `.target.value` on the string (undefined), setting `query` to `''` always.
- **Fix:** Refactor `handleInput` to the value-first pattern: `handleInput(value: string)`. Same fix as BC-001 for BcInput.
- **Fix applied:** Refactored `handleInput(e: any)` to `handleInput(value: string)` in BcCommand.lite.tsx. The Mitosis value-events plugin now correctly extracts `e.target.value` and passes it as a string.
- **Status:** fixed (commit 5814114)

### BC-014 — BcCommand: three a11y violations (dialog name, listbox name, combobox aria-controls)

- **Date:** 2026-06-20
- **Component:** BcCommand
- **Severity:** P1 (a11y)
- **Phase:** A1.3 (runtime test — BcCommand)
- **Classification:** a11y (source)
- **Reproduction:**
  1. Render `<BcCommand open items={[...]} />`.
  2. Run axe-core.
- **Expected:** Zero violations.
- **Actual:** Three violations:
  - `aria-dialog-name` — the dialog root has no `aria-label`/`aria-labelledby`.
  - `aria-input-field-name` — the listbox (`role="listbox"`) has no accessible name.
  - `aria-required-attr` — the combobox input is missing `aria-controls` (pointing to the listbox id).
- **Fix:** Add `aria-label` or `aria-labelledby` to the dialog root. Add `aria-label` to the listbox. Add `aria-controls={listboxId}` to the combobox input and `id={listboxId}` to the listbox.
- **Status:** deferred to Phase B2 (a11y runtime hardening)

### BC-015 — BcCommand: stale-state closure in keydown handler (rapid ArrowDown/Up uses captured activeIndex)

- **Date:** 2026-06-20
- **Component:** BcCommand
- **Severity:** P1 (keyboard nav broken for rapid key presses)
- **Phase:** A1.3 (runtime test — BcCommand)
- **Classification:** source (React state pattern)
- **Reproduction:**
  1. Render `<BcCommand open items={[...]} />`.
  2. Rapidly press ArrowDown multiple times using `fireEvent.keyDown` (synchronous, no re-render between events).
- **Expected:** `activeIndex` increments by 1 per keypress, clamped at the last item.
- **Actual:** `activeIndex` only increments once. The handler reads `activeIndex` from the closure (which is 0 on the first render), so all rapid keydowns set it to `Math.min(0 + 1, list.length - 1) = 1`.
- **Root cause:** The compiled handler `setActiveIndex(Math.min(activeIndex + 1, list.length - 1))` captures `activeIndex` from the render scope. Rapid `fireEvent` calls don't trigger re-renders between events, so the closure value is stale. This is a known React pattern issue — the fix is to use the functional updater form: `setActiveIndex(prev => Math.min(prev + 1, list.length - 1))`.
- **Fix:** Use functional state updates in the keydown handler: `setActiveIndex(prev => Math.min(prev + 1, list.length - 1))`. This works correctly regardless of batching. With `userEvent` (which awaits re-renders), the current code works, but `fireEvent` exposes the stale-closure issue.
- **Status:** deferred to Phase B (keyboard nav hardening — B2)

### BC-016 — BcAccordion: onToggle callback receives stale openIds (pre-toggle state)

- **Date:** 2026-06-20
- **Component:** BcAccordion
- **Severity:** P1 (callback contract broken)
- **Phase:** A1.3 (runtime test — BcAccordion)
- **Classification:** source (React state pattern)
- **Reproduction:**
  1. Render `<BcAccordion items={[...]} onToggle={fn} />`.
  2. Click an item to open it.
- **Expected:** `onToggle` is called with the new open IDs array (e.g., `['a']`).
- **Actual:** `onToggle` is called with the stale pre-toggle array (e.g., `[]`).
- **Root cause:** Source calls `props.onToggle?.(state.openIds)` AFTER `setOpenIds(...)` but BEFORE React re-renders. The `state.openIds` still holds the old value. Same issue applies to `dispatchBcEvent` which also receives the stale array.
- **Fix:** Compute the next array first, call `setOpenIds(next)`, then pass `next` to `onToggle` and `dispatchBcEvent`. This avoids relying on state updates being synchronous.
- **Status:** deferred to Phase B (state callback contract hardening)

### BC-017 — BcSplitter: pointer drag unreliable in jsdom (split into BC-017a wontfix + BC-017b fixed)

- **Date:** 2026-06-20 (updated 2026-06-21)
- **Component:** BcSplitter
- **Severity:** P2 (testing limitation)
- **Phase:** A1.3 (runtime test — BcSplitter)
- **Classification:** test (jsdom limitation) + source (stale-state closure)
- **Reproduction:**
  1. Render `<BcSplitter start="A" end="B" onResize={fn} />`.
  2. `fireEvent.pointerDown(handle)` then `fireEvent.pointerMove(handle, { clientX: 700 })`.
- **Expected:** `onResize` called with the computed percentage (70).
- **Actual:** `onResize` called with NaN. Two compounding issues identified and split:

**BC-017a (jsdom limitation, FIXED via workaround):** jsdom's `fireEvent.pointerMove` does not deliver `clientX`/`clientY` to the event handler (they arrive as `undefined`), so the percentage calculation produces NaN. Additionally, jsdom does not implement layout, so `getBoundingClientRect` returns all zeros. **Fix:** the test constructs a `MouseEvent` with `type: 'pointermove'` and `clientX/clientY` properties (jsdom's MouseEvent constructor honors these, unlike fireEvent.pointerMove), and mocks `getBoundingClientRect` on the root element. The test now asserts the exact percentage value (70). Production browsers deliver real PointerEvents and rects without workaround.

**BC-017b (stale-state closure, FIXED):** The `onMove` handler read `dragging` from the React state closure, which could be stale (`false`) when `pointerDown` and `pointerMove` fired without an intervening re-render. This is the same pattern as BC-015 (BcCommand) and BC-019 (BcForm) which were already fixed.

- **Fix (BC-017b):** Added `draggingRef = useRef<boolean>(false)` that mirrors `state.dragging`. `startDrag`/`stopDrag`/`onMove` read from `draggingRef` (always current) instead of `state.dragging` (may be stale in closure). `state.dragging` remains the source of truth for render (`data-dragging` attribute + overlay `Show`). The test now verifies that `onMove` ran (onResize was called) — without the fix, onMove would early-return on the stale `dragging=false` and onResize would never fire.
- **Status:** BC-017a fixed (MouseEvent constructor + getBoundingClientRect mock workaround). BC-017b fixed (useRef for dragging flag).

### BC-018 — BcForm: select onChange passes Event instead of value (Mitosis value-events plugin mismatch)

- **Date:** 2026-06-20
- **Component:** BcForm
- **Severity:** P1 (select field values not stored correctly)
- **Phase:** A1.3 (runtime test — BcForm)
- **Classification:** source + plugin interaction
- **Reproduction:**
  1. Render `<BcForm fields={[{ name: 'role', type: 'select', options: [...] }]} />`.
  2. Change the select value.
- **Expected:** `onChange` fires with the form values including the selected option value.
- **Actual:** `onChange` fires but the stored value for the select field is the Event object, not the selected value string.
- **Root cause:** The compiled `<select>` uses `onChange={value => setField(field, value)}` which passes the React Event object as `value`. The Mitosis value-events plugin rewrote `onInput` to `onChange` for text/textarea (correctly extracting `e.target.value`), but the source `<select>` used `onChange` directly (not `onInput`), so the plugin left it as-is and the handler receives the Event instead of the string value.
- **Fix:** Change the source `<select>` to use `onInput` (like text/textarea) so the plugin applies the same value-first transform. Alternatively, refactor `setField` to accept the Event for select specifically.
- **Status:** deferred to Phase B (value-first refactor)

### BC-019 — BcForm: stale-state closure in validateAll and setField (errors may not update correctly)

- **Date:** 2026-06-20
- **Component:** BcForm
- **Severity:** P1 (validation errors may not render on first submit or clear on correction)
- **Phase:** A1.3 (runtime test — BcForm)
- **Classification:** source (React state pattern)
- **Reproduction:**
  1. Render `<BcForm fields={[{ name: 'email', type: 'email', required: true }]} values={{ email: 'bad' }} />`.
  2. Click Submit.
- **Expected:** The email validation error "Invalid email address" appears.
- **Actual:** The error may not render. `validateAll` reads `formValues` from the closure, which is seeded from `props.values` at mount. The validation runs against the correct values, but the `setErrors` state update may not flush before the DOM is checked. Similarly, `setField` uses stale `formValues` when computing the next error state after a user correction.
- **Root cause:** React state batching means `setErrors(nextErrors)` may not be visible immediately after `validateAll()` returns. The closure-captured `formValues` in `setField` is stale after `setFormValues` is called.
- **Fix:** Use functional state updates: `setFormValues(prev => ({ ...prev, [field.name]: value }))` and `setErrors(prev => { ... })`. Also compute validation results synchronously and pass them directly to callbacks.
- **Fix applied:** Refactored `setField` to validate against the NEW value (`nextValues[field.name]`) instead of the stale `formValues` closure. Inlined the validation logic in `setField` to avoid the stale-state read. Error clearing on field correction now works.
- **Status:** fixed (commit 5814114)

### BC-020 — BcList: local mode status set to 'empty' instead of 'ready' (stale total closure)

- **Date:** 2026-06-20
- **Component:** BcList
- **Severity:** P1 (pagination never renders for local data)
- **Phase:** A1.3 (runtime test — BcList)
- **Classification:** source (React state pattern)
- **Reproduction:**
  1. Render `<BcList data={[...5 items...]} pageSize={2} />`.
  2. Check for pagination controls.
- **Expected:** Pagination renders (5 items / 2 per page = 3 pages).
- **Actual:** No pagination. `loadData()` for local mode calls `setTotal(props.data.length)` then checks `total === 0` to set status. But `total` is the stale closure value (0 from initial state), so status is set to 'empty' instead of 'ready'. Pagination only renders when `status === 'ready'`.
- **Root cause:** `setStatus(total === 0 ? 'empty' : 'ready')` uses the stale `total` from the closure (initial value 0), not the updated value after `setTotal(props.data.length)`.
- **Fix:** Compute the new total first: `const newTotal = props.data ? props.data.length : 0; setTotal(newTotal); setStatus(newTotal === 0 ? 'empty' : 'ready');`
- **Fix applied:** Used a local `newTotal` variable in `loadData()` instead of reading `state.total` after `state.total = ...`. Pagination now renders correctly for local data.
- **Status:** fixed (commit 5814114)

### BC-021 — BcRenderer: custom sentinel `data-bc-renderer-custom` attribute stripped by Mitosis plugin

- **Date:** 2026-06-20
- **Component:** BcRenderer
- **Severity:** P1 (custom sentinel unidentifiable in DOM)
- **Phase:** A1.3 (runtime test — BcRenderer)
- **Classification:** mitosis-output (plugin strips data-bc-* attributes)
- **Reproduction:**
  1. Render `<BcRenderer schema={{ type: 'custom', name: 'my-widget' }} />`.
  2. Query for `[data-bc-renderer-custom="my-widget"]`.
- **Expected:** The sentinel div has `data-bc-renderer-custom="my-widget"`.
- **Actual:** The attribute is missing. The `mitosis-plugin-react-fixes.mjs` step strips all `data-bc-*` attributes from the compiled output, including the custom sentinel identifier. Only `data-payload` survives.
- **Root cause:** The plugin's `data-bc-*` stripping rule is too aggressive — it strips ALL `data-bc-*` attributes, including those that are intentionally part of the component's API (like `data-bc-renderer-custom`, `data-bc-dialog`, `data-bc-drawer`, etc.).
- **Fix:** Either (1) change the stripping rule to only strip `data-bc-*` attributes that are NOT in a whitelist, or (2) use a different attribute naming convention for sentinel markers (e.g., `data-renderer-custom` without the `bc-` prefix), or (3) remove the stripping rule entirely and handle the underlying issue differently.
- **Status:** deferred to Phase B (Mitosis plugin hardening)

### BC-022 — BcDataTable: `defaultGetRowId is not defined` (Mitosis compilation issue)

- **Date:** 2026-06-20
- **Component:** BcDataTable
- **Severity:** P0 (component crashes on render)
- **Phase:** A1.3 (runtime test — BcRenderer list schema)
- **Classification:** mitosis-output
- **Reproduction:**
  1. Render `<BcRenderer schema={{ type: 'list', columns_def: [...], data: [...] }} />`.
  2. Or render `<BcDataTable columns={[...]} data={[...]} />` directly.
- **Expected:** The table renders with rows.
- **Actual:** `ReferenceError: defaultGetRowId is not defined` — the compiled output references a function that Mitosis failed to emit.
- **Root cause:** The BcDataTable source likely defines or imports a `defaultGetRowId` function that the Mitosis compiler does not include in the React output. This could be a scope issue where the function is defined inside a closure that Mitosis doesn't traverse.
- **Fix:** Investigate the BcDataTable source for `defaultGetRowId`. If it's a local function, move it to module scope. If it's imported, ensure the import is preserved in the compiled output. This blocks all BcDataTable runtime tests.
- **Fix applied:** Moved `defaultGetRowId` from function body to module scope. Also moved `can()` permission helper from function-body `const` into `useStore` (same Mitosis drop pattern). Also fixed all ternary-in-JSX event handlers (onClick, onKeyDown, onContextMenu, onDblClick, onScroll, onDragStart) that compiled to broken IIFEs (BC-010 pattern).
- **Status:** fixed (commit a746018)

### BC-023 — BcDataTable: props.loading not wired to status state

- **Date:** 2026-06-21
- **Component:** BcDataTable
- **Severity:** P1
- **Phase:** A1.3 (runtime test — BcDataTable)
- **Classification:** source (feature gap)
- **Reproduction:**
  1. Render `<BcDataTable columns={[...]} data={[]} loading />`.
  2. Check for loading state indicator.
- **Expected:** The table shows a loading indicator.
- **Actual:** The table shows the empty state. `props.loading` is declared in the interface and defaults to false, but the component never reads it to set `status` to `'loading'`. The loading state only appears when `source` is a fetcher (async).
- **Fix:** In `onMount` or a `useEffect`, check `props.loading` and set `status = 'loading'` when true.
- **Status:** deferred to Phase B

### BC-024 — BcReactiveScope: onChange subscription timing + cleanup memory leak (split into BC-024a deferred + BC-024b fixed)

- **Date:** 2026-06-21 (updated 2026-06-21)
- **Component:** BcReactiveScope
- **Severity:** P1
- **Phase:** A1.3 (runtime test — BcReactiveScope)
- **Classification:** source (React effect timing + cleanup closure)
- **Reproduction:**
  1. Render `<BcReactiveScope store={externalStore} onChange={fn}>...</BcReactiveScope>`.
  2. Call `externalStore.set('x', 1)`.
- **Expected:** `onChange` fires with the updated values.
- **Actual:** Two distinct issues identified and split:

**BC-024a (test-timing, FIXED with waitFor):** The `useEffect` that subscribes `onChange` runs on mount, but in jsdom + React test env, `store.set()` may fire before the subscription is established (effect not yet flushed). **Fix:** the test uses `waitFor` to poll until `onChange` has fired at least once on a probe value, then mutates the real value and asserts the callback receives it. This replaces the previous fixed `setTimeout(50)` which was unreliable.

**BC-024b (cleanup memory leak, FIXED):** The compiled `useEffect` cleanup (deps `[]`) captured `scopeOff` from the first render (initial value `null` from `mkScopeOff()`). After `setScopeOff(off)` in the mount effect triggered a re-render, the cleanup closure still held the stale `null` — so `unsubscribe` was never invoked on unmount. With an externally-owned store (`props.store` provided), the store is NOT disposed on unmount, so the listener leaked: `store.set()` after unmount would still fire the stale `onChange` callback.

- **Fix (BC-024b):** Two changes:
  1. Added `unsubscribe(listener)` method to `ReactiveStore` (symmetric with `subscribe`, mirrors DOM EventTarget add/removeEventListener pattern).
  2. Refactored `BcReactiveScope.lite.tsx` to store the listener function in a `useRef` (`changeListenerRef`) instead of `useStore` state. `useRef.current` is a mutable container that is always current regardless of closure capture, so `onUnMount` reads the real listener and calls `store.unsubscribe(changeListenerRef.current)`. Removed `ScopeOff` type and `mkScopeOff` factory (no longer needed).
- **Test:** Added "unsubscribes onChange listener on unmount (BC-024b fix)" test in `reactive.rtl.test.tsx` — renders with external store + onChange, unmounts, then asserts `onChange` is NOT called after `store.set()` post-unmount.
- **Note (extended fix):** During the 2026-06-21 follow-up session, the same stale-disposer pattern was found in `BcForm.lite.tsx` and `AppProvider.lite.tsx`. Deeper investigation revealed that **`BcForm` had a worse latent bug**: `reactiveStore` was also read from the cleanup closure (always null from the first render), so `dispose()` was never invoked — meaning the listener actually DID leak, contrary to the initial assumption that disposal saves it. **Fix:** both components now use the dual-storage pattern: `useStore` for render-reactive state and `useRef` for cleanup-stable references. `BcForm` adds `reactiveStoreRef` + `reactiveListenerRef`; `AppProvider` moves `disposers` to `disposersRef`. All three components now have deterministic cleanup regardless of closure capture.
- **Status:** BC-024a fixed (waitFor poll workaround for test-timing). BC-024b fixed (useRef + unsubscribe method). Same pattern applied to BcForm + AppProvider.


