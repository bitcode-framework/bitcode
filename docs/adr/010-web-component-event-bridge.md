# ADR-010: Web Component Event Bridge

## Status: Accepted (2026-06-18)

## Date: 2026-06-18

## Depends on: [ADR-007](./007-enterprise-universal-ui-api-contract.md), [ADR-008](./008-contract-first-api-freeze.md)

## Context

ADR-007 Section 3 defines semantic event contracts for cross-framework compatibility. However, Web Components (Stencil output) require a different event mechanism than React/Vue callbacks:

- **React/Vue**: Use semantic callbacks like `onValueChange(value, meta)`
- **Web Components**: Must dispatch `CustomEvent` on the host element with kebab-case names like `bc-value-change`

The challenge: how do we bridge semantic callbacks to Web Component CustomEvents while maintaining:
1. Framework-agnostic event dispatch
2. SSR/React Native safety (no `document` or `CustomEvent` in those environments)
3. Consistent event detail shape across all 14 semantic events
4. Shadow DOM piercing (events must bubble through shadow boundaries)
5. Type safety for event payloads

## Decision

Use a **Web Component Event Bridge**: a single `dispatchBcEvent()` function that dispatches kebab-case `bc-*` CustomEvents on the host element, with SSR-safe fallbacks.

### Implementation: `core/wc-event-bridge.ts`

```typescript
import type { BcChangeMeta, BcEventSource } from './types';

export interface BcWcEventDetail<T = unknown> {
  value: T;
  meta: BcChangeMeta;
}

export function dispatchBcEvent<T>(
  host: unknown,
  eventName: string,
  value: T,
  meta: BcChangeMeta,
): void {
  if (
    typeof document === 'undefined' ||
    typeof CustomEvent === 'undefined' ||
    !host ||
    typeof (host as Element).dispatchEvent !== 'function'
  ) {
    return;  // SSR/RN safe: silently no-op
  }
  if (typeof eventName !== 'string' || !eventName) {
    return;
  }
  try {
    const detail: BcWcEventDetail<T> = { value, meta };
    (host as Element).dispatchEvent(
      new CustomEvent(eventName, {
        detail,
        bubbles: true,
        composed: true,  // Pierces shadow DOM
      }),
    );
  } catch {
    // SSR-safe: silently swallow errors in restricted environments
  }
}

export function buildMeta<TReason extends string = string>(
  reason: TReason,
  source: BcEventSource = 'programmatic',
  nativeEvent?: unknown,
): BcChangeMeta<TReason> {
  const meta: BcChangeMeta<TReason> = { reason, source };
  if (nativeEvent !== undefined) {
    meta.nativeEvent = nativeEvent;
  }
  return meta;
}

export const BC_EVENTS = {
  VALUE_CHANGE: 'bc-value-change',
  OPEN_CHANGE: 'bc-open-change',
  CHECKED_CHANGE: 'bc-checked-change',
  SELECTION_CHANGE: 'bc-selection-change',
  ACTIVE_KEY_CHANGE: 'bc-active-key-change',
  EXPANDED_CHANGE: 'bc-expanded-change',
  SORT_MODEL_CHANGE: 'bc-sort-model-change',
  FILTER_MODEL_CHANGE: 'bc-filter-model-change',
  SEARCH_CHANGE: 'bc-search-change',
  COLUMN_STATE_CHANGE: 'bc-column-state-change',
  ACTION: 'bc-action',
  SUBMIT: 'bc-submit',
  VALIDATE: 'bc-validate',
  SIZE_CHANGE: 'bc-size-change',
} as const;
```

### Usage in Components

```typescript
import { dispatchBcEvent, buildMeta, BC_EVENTS } from '@/core/wc-event-bridge';

export default function BcSelect(props: BcSelectProps) {
  const handleSelect = (value: string, nativeEvent?: Event) => {
    const meta = buildMeta('select', 'pointer', nativeEvent);
    
    // Semantic callback for React/Vue
    props.onValueChange?.(value, meta);
    
    // Web Component CustomEvent
    dispatchBcEvent(
      props.host,  // Host element reference
      BC_EVENTS.VALUE_CHANGE,
      value,
      meta,
    );
  };
  
  // ...
}
```

## 14 Semantic Events Defined

The bridge supports 14 canonical events per ADR-007 Section 3:

| Event | Use Case | Example Components |
|---|---|---|
| `bc-value-change` | Value changes | BcInput, BcSelect, BcCombobox |
| `bc-open-change` | Overlay open/close | BcDialog, BcDrawer, BcPopover |
| `bc-checked-change` | Checkbox/toggle state | BcCheckbox, BcToggle |
| `bc-selection-change` | Collection selection | BcDataTable, BcList |
| `bc-active-key-change` | Tab/accordion active item | BcTabs, BcAccordion |
| `bc-expanded-change` | Tree/accordion expansion | BcAccordion, BcDisclosure |
| `bc-sort-model-change` | Data table sort | BcDataTable |
| `bc-filter-model-change` | Data table filter | BcDataTable |
| `bc-search-change` | Search input | BcCommand, BcDataTable |
| `bc-column-state-change` | Data table columns | BcDataTable |
| `bc-action` | Generic action | BcDropdown, BcContextMenu |
| `bc-submit` | Form submission | BcForm |
| `bc-validate` | Validation result | BcForm, BcInput |
| `bc-size-change` | Splitter resize | BcSplitter |

## Error Handling

### SSR/React Native Safety

The bridge guards against environments without `document` or `CustomEvent`:

```typescript
if (
  typeof document === 'undefined' ||
  typeof CustomEvent === 'undefined' ||
  !host ||
  typeof (host as Element).dispatchEvent !== 'function'
) {
  return;  // Silently no-op
}
```

**Rationale:** React, Vue, and React Native do not need CustomEvents (they use semantic callbacks). Only Stencil/Web Components dispatch CustomEvents.

### Try/Catch for Robustness

Even if guards pass, `dispatchEvent` can throw in restricted environments:

```typescript
try {
  (host as Element).dispatchEvent(new CustomEvent(eventName, { ... }));
} catch {
  // Silently swallow errors
}
```

**Rationale:** Components should never crash due to event dispatch failures.

### Invalid Input Handling

The bridge validates `eventName` and `host`:

```typescript
if (typeof eventName !== 'string' || !eventName) {
  return;
}
```

**Rationale:** Prevents runtime errors from typos or undefined values.

## Shadow DOM Piercing

CustomEvents use `bubbles: true` and `composed: true`:

```typescript
new CustomEvent(eventName, {
  detail,
  bubbles: true,
  composed: true,  // Critical for shadow DOM
})
```

**Rationale:** Stencil components render in shadow DOM. Events must pierce shadow boundaries to reach parent components.

## Event Detail Shape

All events use a consistent detail shape:

```typescript
interface BcWcEventDetail<T = unknown> {
  value: T;
  meta: BcChangeMeta;
}
```

**Example:**
```typescript
// bc-value-change event
{
  detail: {
    value: 'apple',
    meta: {
      reason: 'select',
      source: 'pointer',
      nativeEvent: MouseEvent
    }
  }
}
```

**Rationale:** Consistent shape makes event handling predictable across all 14 events.

## Consequences

### Positive

1. **Framework-agnostic** — Works in React, Vue, Stencil, React Native, Flutter
2. **Consistent API** — All 14 events use same `dispatchBcEvent()` function
3. **SSR-safe** — No-ops in environments without `document`/`CustomEvent`
4. **Shadow DOM piercing** — `composed: true` allows events to escape shadow boundaries
5. **Type-safe** — `BcChangeMeta<TReason>` generic for reason type safety
6. **Testable** — Pure function is easy to unit test (14 tests in `wc-event-bridge.test.ts`)
7. **Semantic callbacks preserved** — React/Vue can ignore CustomEvents and use `onValueChange` directly

### Negative

1. **Explicit calls required** — Each component must call `dispatchBcEvent()` in addition to semantic callbacks
2. **Dual event path** — Components emit both semantic callbacks and CustomEvents (slight overhead)
3. **Host reference needed** — Components must receive `host` prop to dispatch events

### Trade-off Assessment

For cross-framework UI libraries, **explicit event bridging is preferable to automatic dispatch** because:
- Clear which events are dispatched
- Works in all framework targets
- SSR/React Native safe
- No framework-specific magic

The dual event path overhead is negligible (function call + event construction).

## Acceptance Criteria

All criteria met as of 2026-06-18:

- [x] **14 semantic events defined** — `BC_EVENTS` constant in `wc-event-bridge.ts`
- [x] **SSR/React Native safe** — Guarded `typeof` checks for `document`/`CustomEvent`
- [x] **Shadow DOM piercing** — `composed: true` on all CustomEvents
- [x] **Consistent detail shape** — `{ value, meta }` for all events
- [x] **Type-safe** — `BcChangeMeta<TReason>` generic
- [x] **Test coverage** — 14 unit tests + 8 runtime tests, 100% passing
- [x] **All 25 components wired** — Verified by grep (all components import `dispatchBcEvent`)

## Implementation Evidence

### Test Files

- `packages/ui/src/core/wc-event-bridge.test.ts` — 14 unit tests
- `packages/ui/src/core/wc-event-bridge.runtime.test.ts` — 8 runtime behavior tests

### Component Coverage

All 25 Core Integration Surface components use the Web Component Event Bridge:

**Form Controls (10):** BcInput, BcSelect, BcCombobox, BcCheckbox, BcToggle, BcSegmented, BcNumberField, BcDateField, BcTimeField, BcMultiSelect

**Forms (2):** BcForm, BcReactiveScope

**Overlays (8):** BcDialog, BcAlertDialog, BcDrawer, BcPopover, BcTooltip, BcDropdown, BcContextMenu, BcCommand

**Navigation/Layout (4):** BcTabs, BcAccordion, BcDisclosure, BcSplitter

**Data (1):** BcDataTable

## Related ADRs

- [ADR-007: Enterprise Universal UI API Contract](./007-enterprise-universal-ui-api-contract.md) — Defines semantic event contracts (Section 3)
- [ADR-008: Contract-First API Freeze](./008-contract-first-api-freeze.md) — Enforces event contract matrix
- [ADR-009: Runtime Part Resolution](./009-runtime-part-resolution.md) — Complementary mechanism for slotProps

## References

- Implementation: `packages/ui/src/core/wc-event-bridge.ts`
- ADR-007 Section 3: Event standard specification
- Web Components CustomEvent spec: https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent
- Shadow DOM event propagation: https://developer.mozilla.org/en-US/docs/Web/API/Event/composed
- Test suite: `packages/ui/src/core/wc-event-bridge.test.ts`, `wc-event-bridge.runtime.test.ts`
