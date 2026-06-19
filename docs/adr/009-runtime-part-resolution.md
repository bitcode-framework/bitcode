# ADR-009: Runtime Part Resolution

## Status: Accepted (2026-06-18)

## Date: 2026-06-18

## Depends on: [ADR-007](./007-enterprise-universal-ui-api-contract.md), [ADR-008](./008-contract-first-api-freeze.md)

## Context

ADR-007 Section 6 defines the extensibility contract requiring `slotProps` support for named component parts. However, Mitosis (the multi-target component compiler) cannot compile dynamic object spread syntax (`{...props.slotProps?.root}`) on JSX elements.

This is a deliberate Mitosis limitation because several output targets (Stencil, Vue, Angular) do not support arbitrary dynamic attribute spread. Attempting to use spread results in compilation errors or runtime failures in those targets.

The challenge: how do we implement the slotProps extensibility contract (ADR-007 Section 6) within Mitosis constraints while maintaining:
1. Cross-platform compatibility (React, Vue, Stencil, React Native, Flutter)
2. Type safety
3. Runtime flexibility (metadata-driven UI can provide slotProps at render time)
4. No build-time magic or code generation

## Decision

Use **Runtime Part Resolution**: pure function helpers that components call per-element to read individual values from `slotProps`, avoiding dynamic spread entirely.

### Implementation: `core/part-props.ts`

Three pure functions provide the slotProps mechanism:

```typescript
export type SlotPropsMap = Partial<Record<string, Record<string, unknown>>> | undefined;

export function partClass(
  slotProps: SlotPropsMap,
  partName: string,
  baseClass: string,
): string {
  if (!slotProps || typeof slotProps !== 'object') return baseClass;
  const partMap = slotProps[partName];
  if (!partMap) return baseClass;
  const extra = partMap.className;
  if (extra === undefined || extra === null || extra === '') return baseClass;
  if (typeof extra !== 'string') return baseClass;
  if (!baseClass) return extra;
  return baseClass + ' ' + extra;
}

export function partAttr(
  slotProps: SlotPropsMap,
  partName: string,
  attrName: string,
): unknown {
  if (!slotProps || typeof slotProps !== 'object') return undefined;
  const partMap = slotProps[partName];
  if (!partMap) return undefined;
  return partMap[attrName];
}

export function partStyle(
  slotProps: SlotPropsMap,
  partName: string = 'root',
): unknown {
  if (!slotProps || typeof slotProps !== 'object') return undefined;
  const partMap = slotProps[partName];
  if (!partMap) return undefined;
  return partMap.style;
}
```

### Usage in Mitosis Components

```tsx
import { partClass, partAttr } from '@/core/part-props';

export default function BcInput(props: BcInputProps) {
  return (
    <div
      class={partClass(props.slotProps, 'root', 'flex flex-col gap-1')}
      data-part="root"
      data-testid={partAttr(props.slotProps, 'root', 'data-testid') as string}
    >
      <label
        class={partClass(props.slotProps, 'label', 'text-sm font-medium')}
        data-part="label"
      >
        {props.label}
      </label>
      <input
        class={partClass(props.slotProps, 'input', 'border rounded px-2 py-1')}
        data-part="input"
        aria-label={partAttr(props.slotProps, 'input', 'aria-label') as string}
      />
    </div>
  );
}
```

## Why This Works

### 1. Mitosis Compatibility

Property access and string concatenation are valid Mitosis JSX:
- ✅ `props.slotProps?.root?.className` — property access
- ✅ `baseClass + ' ' + extra` — string concatenation
- ❌ `{...props.slotProps?.root}` — dynamic spread (not supported)

### 2. Cross-Platform Portability

Pure functions have no DOM dependency:
- **React/Vue/Svelte**: Direct function calls
- **React Native**: Same function calls, NativeWind classes
- **Flutter/Dart**: Port the same logic to Dart
- **Stencil/Web Components**: Works in shadow DOM

### 3. Runtime Flexibility

Components call helpers at render time, so:
- Metadata-driven UI (ADR-005) can provide `slotProps` from schema data
- No build-time code generation required
- Dynamic `slotProps` from parent components work naturally

### 4. Type Safety

- `SlotPropsMap` type alias for better IDE support
- Type guards validate inputs before use (never throw)
- Return safe defaults for invalid inputs

## Implementation Rules

### Rule 1: className is Additive, Not Replace

`slotProps.root.className` is concatenated **after** the component's base classes, preserving token-driven styling as the source of truth (ADR-002 Section 3, ADR-006 Section 4).

**Example:**
```typescript
// Component base: "flex items-center bg-background"
// Consumer slotProps: { root: { className: "p-4 shadow-lg" } }
// Result: "flex items-center bg-background p-4 shadow-lg"
```

**Rationale:** Consumer classes append, never override the base. This prevents accidental removal of critical styling (e.g., `flex`, `items-center`).

### Rule 2: Supported Attribute Categories

Components wire `slotProps` passthrough for:
- `className` — styling overrides
- `data-*` attributes — testing, analytics
- `aria-*` attributes — a11y overrides
- `style` — inline style (where supported)

Event handler passthrough is **intentionally excluded** to keep the component surface controlled.

### Rule 3: `data-part` Attribute Required

Every element that accepts `slotProps` must emit `data-part="{partName}"`:

```tsx
<div class={partClass(props.slotProps, 'root', baseClass)} data-part="root">
```

**Rationale:** Consumers can target parts via CSS selectors and testing tools, independent of `slotProps`:

```css
[data-part="root"] { margin-bottom: 1rem; }
[data-part="label"] { font-weight: 600; }
```

### Rule 4: Explicit Wiring Per Attribute

Components wire only the attributes they explicitly support per part. This is a **feature**, not a limitation:

- Enterprise UI libraries should have a controlled surface
- Arbitrary prop passthrough invites breakage across framework targets
- Clear which attributes are supported per part
- Easier to audit and maintain

**Trade-off:** More verbose component code, but safer and more predictable.

### Rule 5: Flutter Parity

`partClass` and `partAttr` are pure functions with no DOM dependency. A Dart port can implement the same contract:

```dart
String partClass(
  Map<String, Map<String, dynamic>>? slotProps,
  String partName,
  String baseClass,
) {
  if (slotProps == null) return baseClass;
  final partMap = slotProps[partName];
  if (partMap == null) return baseClass;
  final extra = partMap['className'];
  if (extra == null || extra.isEmpty) return baseClass;
  if (extra is! String) return baseClass;
  if (baseClass.isEmpty) return extra;
  return '$baseClass $extra';
}
```

Flutter adapters receive `slotProps` as a typed map and apply values through the equivalent Dart widget API.

### Rule 6: Metadata-Driven Integration

`BcRenderer` (ADR-005) may populate `slotProps` from `UiSchema` metadata at render time:

```typescript
// Schema field with density: 'compact'
const slotProps = {
  root: { className: 'p-1' },
  input: { className: 'text-xs' },
};

<BcInput slotProps={slotProps} />
```

This keeps the renderer a simple compositor without per-component knowledge.

## Consequences

### Positive

1. **Works across all targets** — React, Vue, Stencil, React Native, Flutter
2. **Type-safe** — Type guards prevent runtime errors
3. **No build magic** — Pure runtime functions, no code generation
4. **Metadata-driven friendly** — `BcRenderer` can provide `slotProps` from schema
5. **Additive styling** — Token-driven classes remain source of truth
6. **Testable** — Pure functions are easy to unit test (34 tests in `part-props.test.ts`)
7. **Portable** — No DOM dependency, works in SSR/React Native

### Negative

1. **Explicit wiring required** — Each component must call `partClass`/`partAttr` per element
2. **More verbose** — Cannot use generic spread, must wire each supported attribute
3. **Discipline required** — Authors must remember to wire all supported parts

### Trade-off Assessment

For enterprise UI libraries, **explicit wiring is preferable to generic spread** because:
- Controlled surface prevents arbitrary prop passthrough
- Clear which attributes are supported per part
- Easier to audit and maintain
- Works across all framework targets

The verbosity cost is acceptable for the safety and portability benefits.

## Acceptance Criteria

All criteria met as of 2026-06-18:

- [x] **All 25 components wired** — Verified by `runtime-part-rollout.test.ts`
- [x] **`partClass` is additive (never replaces)** — Tested in `part-props.test.ts`
- [x] **`data-part` attributes emitted** — All components emit `data-part="{partName}"`
- [x] **Type-safe `slotProps` interface** — `SlotPropsMap` type with type guards
- [x] **Works in SSR/React Native** — No DOM dependency, guarded `typeof` checks
- [x] **Test coverage** — 34 unit tests + 18 runtime tests, 100% passing
- [x] **Documented in ADR-007 Section 6.1** — Cross-referenced

## Implementation Evidence

### Test Files

- `packages/ui/src/core/part-props.test.ts` — 34 unit tests
- `packages/ui/src/core/part-props.runtime.test.ts` — 18 runtime behavior tests
- `packages/ui/src/components/runtime-part-rollout.test.ts` — Verifies all 25 components wired

### Component Coverage

All 25 Core Integration Surface components use Runtime Part Resolution:

**Form Controls (10):** BcInput, BcSelect, BcCombobox, BcCheckbox, BcToggle, BcSegmented, BcNumberField, BcDateField, BcTimeField, BcMultiSelect

**Forms (2):** BcForm, BcReactiveScope

**Overlays (8):** BcDialog, BcAlertDialog, BcDrawer, BcPopover, BcTooltip, BcDropdown, BcContextMenu, BcCommand

**Navigation/Layout (4):** BcTabs, BcAccordion, BcDisclosure, BcSplitter

**Data (1):** BcDataTable

## Related ADRs

- [ADR-007: Enterprise Universal UI API Contract](./007-enterprise-universal-ui-api-contract.md) — Defines slotProps extensibility contract (Section 6)
- [ADR-008: Contract-First API Freeze](./008-contract-first-api-freeze.md) — Enforces contract matrix including slotProps
- [ADR-010: Web Component Event Bridge](./010-web-component-event-bridge.md) — Complementary mechanism for CustomEvent dispatch

## References

- Implementation: `packages/ui/src/core/part-props.ts`
- Mitosis documentation: JSX constraints, property access, string concatenation
- ADR-007 Section 6.1: Runtime Part Resolution specification
- Test suite: `packages/ui/src/core/part-props.test.ts`, `part-props.runtime.test.ts`
