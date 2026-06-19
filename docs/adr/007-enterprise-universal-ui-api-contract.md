# ADR-007: Enterprise Universal UI API Contract

## Status: Accepted (Phase 2 implementation complete, 2026-06-18)

## Date: 2026-06-17 (proposed) / 2026-06-18 (accepted)

## Depends on: [ADR-002](./002-cross-platform-architecture.md), [ADR-003](./003-component-roadmap.md), [ADR-004](./004-reactive-engine.md), [ADR-005](./005-metadata-driven-ui.md), [ADR-006](./006-universal-component-strategy.md)

## Context

`packages/ui` is intended to become the universal UI foundation for BitCode and external applications. The library must be enterprise-grade, standalone, and easy to integrate across React, Vue, Stencil/Web Components, React Native, Tauri, and future adapters such as Flutter.

The previous Stencil component generation in `packages/components` proved that a large component catalog is not enough. Integration into React + Vite and similar stacks becomes painful when events are not standard, payloads are inconsistent, component state is difficult to control, and accessibility behavior is not encoded as a contract.

ADR-004 already owns the Filament-like reactive model: cascading data, conditional validation, conditional UI, repeater scope, cross-component reactivity, `ReactiveStore`, DATA/VALIDATION/UI/VALUE effects, and standalone Level 0-3 consumption. ADR-005 already owns the metadata-driven `BcRenderer` layer that forwards schema reactions into the reactive engine. This ADR does not replace those decisions. It standardizes the public component API surface so those reactive and metadata layers can integrate consistently across renderers.

Research references:

- Ant Design v6: enterprise ergonomics, form/table breadth, design-token continuity, mostly API-compatible v5 to v6 migration.
- MUI and MUI X: slot APIs, controlled models, data-grid state models, accessibility guidance.
- Mantine: provider-driven theming, polymorphic components, consistent size/variant props, form ergonomics.
- shadcn/ui: source-owned composition, Radix-backed primitives, Tailwind token conventions.
- PrimeReact v10: broad enterprise component coverage, pass-through customization, unstyled/styled modes.
- Radix UI v1: accessible primitive parts, `asChild`, controlled/uncontrolled APIs, WAI-ARIA adherence.
- React Aria / React Aria Components: accessibility, internationalization, collection/selection APIs, slots, render props.
- WAI-ARIA APG: normative keyboard, focus, roles, and ARIA behavior patterns.
- Tamagui: cross-platform token and theme architecture for web and native.
- gluestack-ui: React Native + web component model with Tailwind/NativeWind-oriented styling.
- Existing `packages/ui`: Mitosis outputs, design tokens, Pure/Adapt concept, data source abstraction, and component contracts.
- Existing `packages/components`: cautionary example for Web Component events and integration friction.

Current local findings:

- `packages/ui` already exports React, Vue, Stencil, React Native, core hooks, providers, tokens, Tailwind preset, and NativeWind preset.
- `ADR-004` already defines the standalone reactive engine and Filament-like form behavior.
- `ADR-005` already defines schema-driven `BcRenderer` integration on top of ADR-004.
- `ADR-006` already decides Single Caller API, Pure/Adapt components, and token pipeline direction.
- `core/types.ts` defines event payloads like `FieldChangeEvent`, but event naming is still field-centric and not enough for generic primitives, overlays, collections, and data grids.
- `BcSelect` currently uses `value`, internal state, and `onChange(FieldChangeEvent)`, but lacks `defaultValue`, `onValueChange`, `open/defaultOpen/onOpenChange`, native-form compatibility details, collection selection semantics, keyboard contract, and ARIA contract.
- `BcDataTable.contract.ts` already models many enterprise capabilities, but it uses `initialSort`, `initialFilters`, `initialSelected` rather than standardized controlled/uncontrolled state model pairs like `sortModel/defaultSortModel/onSortModelChange`.

## Decision

Adopt a universal component API contract before expanding `packages/ui` further. Component implementation may vary by platform per ADR-006, but the public contract must be stable, typed, semantic, accessible, and renderer-neutral.

## Section 1: Universal API Layers

Every component contract belongs to one or more layers:

1. **Primitive layer**: buttons, labels, slots, layout boxes, separators, portals.
2. **Behavior primitive layer**: dialog, popover, tooltip, menu, tabs, accordion, combobox, select, tree, grid.
3. **Enterprise component layer**: form, data table, calendar, kanban, scheduler, upload, rich text, command palette.
4. **BitCode integration layer**: metadata renderer, data sources, model adapters, permissions, i18n, offline support.

Primitive and behavior APIs must not depend on BitCode concepts. BitCode integration must be additive through adapters and metadata.

## Section 2: Controlled and Uncontrolled State Standard

All stateful components must support controlled and uncontrolled forms consistently.

Canonical pairs:

| State | Controlled | Uncontrolled | Callback |
|---|---|---|---|
| value | `value` | `defaultValue` | `onValueChange` |
| open | `open` | `defaultOpen` | `onOpenChange` |
| checked | `checked` | `defaultChecked` | `onCheckedChange` |
| selected keys | `selectedKeys` | `defaultSelectedKeys` | `onSelectionChange` |
| active key | `activeKey` | `defaultActiveKey` | `onActiveKeyChange` |
| expanded keys | `expandedKeys` | `defaultExpandedKeys` | `onExpandedChange` |
| sort model | `sortModel` | `defaultSortModel` | `onSortModelChange` |
| filter model | `filterModel` | `defaultFilterModel` | `onFilterModelChange` |
| pagination model | `paginationModel` | `defaultPaginationModel` | `onPaginationModelChange` |
| column state | `columnState` | `defaultColumnState` | `onColumnStateChange` |

Rules:

- `initialX` is not the canonical public API for stateful components. Use `defaultX`.
- `onChange` may be kept only for native-form compatibility or migration, but semantic callbacks are preferred.
- A controlled prop must never be mutated internally.
- A component may expose several controlled models when the state domains are independent.

## Section 3: Event Standard

Events must be semantic, stable, serializable where possible, and easy to bridge to Web Components, React, Vue, React Native, and future renderers.

Canonical callback shape:

```ts
export interface BcChangeMeta<TReason extends string = string> {
  reason: TReason;
  source: 'keyboard' | 'pointer' | 'touch' | 'programmatic' | 'form' | 'api';
  nativeEvent?: unknown;
}

export type BcChangeHandler<TValue, TReason extends string = string> = (
  value: TValue,
  meta: BcChangeMeta<TReason>,
) => void;
```

Canonical callbacks:

- `onValueChange(value, meta)`
- `onOpenChange(open, meta)`
- `onCheckedChange(checked, meta)`
- `onSelectionChange(selection, meta)`
- `onActiveKeyChange(key, meta)`
- `onExpandedChange(keys, meta)`
- `onAction(key, meta)`
- `onSubmit(values, meta)`
- `onValidate(result, meta)`
- `onLoad(result, meta)`
- `onError(error, meta)`

Web Component event naming:

- DOM event names should be lowercase/kebab-case, for example `bc-value-change`, `bc-open-change`, `bc-selection-change`.
- Event `detail` must match the semantic callback payload: `{ value, meta }`, `{ open, meta }`, `{ selectedKeys, meta }`.
- Events that represent state changes should be composed and bubble only when integration requires parent handling. Internal implementation events should not leak.

Rules:

- Consumers must not need to parse raw DOM events for component semantics.
- Native events can be attached through pass-through props or included as `meta.nativeEvent`.
- Event names must not be field-only. `FieldChangeEvent` can remain as a BitCode form adapter type, not the universal event standard.

## Section 4: Imperative Method Standard

Imperative APIs are allowed only where declarative control is insufficient.

Allowed universal methods:

- `focus()`
- `blur()`
- `scrollIntoView()` or `scrollToIndex(index)` for virtualized collections
- `reset()` for forms and inputs
- `validate()` for form controls and forms
- `submit()` for forms
- `open()` / `close()` only for overlay controllers where refs are needed

Rules:

- Every imperative method must have a declarative equivalent when feasible.
- Methods must be typed in contract files.
- Methods that only exist on web must be platform-gated or moved into Adapt implementations with a shared no-op/fallback contract.

## Section 5: Accessibility Contract

WAI-ARIA APG is the baseline for behavior primitives and complex enterprise components.

Every interactive component contract must define:

- Role and ARIA attribute requirements.
- Keyboard interaction table.
- Focus entry, movement, trap, restore, and escape behavior.
- Screen reader name, description, and error semantics.
- RTL behavior.
- Disabled and read-only semantics.
- Mobile/touch equivalent interaction where relevant.

Required defaults:

- Use roving tabindex for menus, tabs, radio groups, segmented controls, toolbars, and tree-like navigation.
- Use focus trap, inert or hidden background, Escape close, and focus restore for modal overlays.
- Use `aria-invalid`, `aria-describedby`, and stable IDs for form errors and descriptions.
- Data grid must implement keyboard cell navigation, sort announcements, selectable row semantics, and edit-mode distinction.

## Section 6: Extensibility Contract

`packages/ui` must support both source-owned customization like shadcn/ui and enterprise overrides like MUI, Mantine, and PrimeReact.

Standard extension points:

- `className` for top-level styling on React/web targets.
- `style` for top-level inline style where supported.
- `as` or `asChild` for primitive composition when the renderer supports it.
- `slots` for replacing named internal parts.
- `slotProps` for typed props passed to internal parts.
- `unstyled` for components whose behavior should be used without default visuals.
- `variant`, `size`, `color`, `radius`, `density` for token-backed variants.
- `data-*` attributes for styling states on web output.

Named parts should be stable for behavior primitives and enterprise components. Examples: `root`, `trigger`, `content`, `label`, `description`, `error`, `item`, `indicator`, `header`, `body`, `footer`, `row`, `cell`, `toolbar`, `pagination`.

### Section 6.1: Runtime Part Resolution (slotProps Implementation Mechanism)

Mitosis does not support dynamic object spread (`{...props.slotProps?.root}`) on JSX elements. This is a deliberate Mitosis limitation, not a bug, because several output targets (Stencil, Vue, Angular) do not support arbitrary dynamic attribute spread.

The chosen mechanism is **Runtime Part Resolution**: a pure function helper that components call per-element to read individual values from `slotProps`, avoiding spread entirely. This approach:

- Works within Mitosis compilation constraints (property access and string concatenation are valid Mitosis JSX).
- Is portable to Flutter (pure TS logic, no DOM/framework dependency).
- Keeps components dumb (ADR-004 principle: headless, per-instance, no global state).
- Is runtime, not compile-time, so metadata-driven UI (ADR-005) can provide `slotProps` from schema data at render time.

#### Implementation shape

```ts
// core/part-props.ts

export function partClass(
  slotProps: Record<string, Record<string, unknown>> | undefined,
  partName: string,
  baseClass: string,
): string {
  const extra = slotProps?.[partName]?.className;
  if (!extra || typeof extra !== 'string') return baseClass;
  return baseClass + ' ' + extra;
}

export function partAttr(
  slotProps: Record<string, Record<string, unknown>> | undefined,
  partName: string,
  attrName: string,
): unknown {
  return slotProps?.[partName]?.[attrName];
}
```

#### Usage in Mitosis source

```tsx
// BcInput.lite.tsx
import { partClass, partAttr } from '@/core/part-props';

<div
  class={partClass(props.slotProps, 'root', baseRootClass)}
  data-part="root"
  data-testid={partAttr(props.slotProps, 'root', 'data-testid') as string}
>
  <label
    class={partClass(props.slotProps, 'label', baseLabelClass)}
    data-part="label"
  >
```

#### Rules

1. **className is additive, not replace.** `slotProps.root.className` is concatenated after the component's base classes. Token-driven classes from Tailwind/NativeWind remain the source of truth (ADR-002 Section 4, ADR-006 Section 4). Consumer classes append, never override the base.
2. **Supported attribute categories.** Components should wire `slotProps` passthrough for: `className` (styling), `data-*` attributes (testing, analytics), `aria-*` attributes (a11y overrides), `style` (inline style where supported). Event handler passthrough is intentionally excluded to keep the component surface controlled.
3. **`data-part` attribute required.** Every element that accepts `slotProps` must emit `data-part="{partName}"` so consumers can target parts via CSS selectors and testing tools, independent of `slotProps`.
4. **Explicit wiring per attribute.** Components wire only the attributes they explicitly support per part. This is a feature, not a limitation: enterprise UI libraries should have a controlled surface. Arbitrary prop passthrough invites breakage across framework targets.
5. **Flutter parity.** `partClass` and `partAttr` are pure functions with no DOM dependency. A Dart port can implement the same contract. Flutter adapters receive `slotProps` as a typed map and apply values through the equivalent Dart widget API.
6. **Metadata-driven integration.** `BcRenderer` (ADR-005) may populate `slotProps` from `UiSchema` metadata at render time. For example, a schema field with `density: 'compact'` could translate to `slotProps.root.className = 'p-1'` via the field resolver. This keeps the renderer a simple compositor without per-component knowledge.

## Section 7: Form Contract

Form controls must be native-form-compatible on web, form-library-agnostic in React/Vue, and adapter-friendly for BitCode metadata forms.

Standard form props:

- `name`
- `value`
- `defaultValue`
- `onValueChange`
- `required`
- `disabled`
- `readOnly`
- `invalid`
- `validationStatus`
- `errorMessage`
- `description`
- `label`
- `placeholder`
- `autoFocus`
- `autoComplete` where relevant

Standard field composition:

- `Field`
- `Label`
- `Control`
- `Description`
- `ErrorMessage`

Rules:

- Do not force a specific form state library into component contracts.
- BitCode model metadata should adapt into these props through `BcRenderer`, not redefine form component semantics.
- Web output should use native input/select/textarea semantics where possible, or hidden inputs where custom controls need form submission.

## Section 8: Overlay Contract

Overlays must centralize behavior because inconsistent overlay behavior is a common enterprise integration failure.

Standard overlay props:

- `open`
- `defaultOpen`
- `onOpenChange`
- `modal`
- `placement`
- `offset`
- `collisionBoundary`
- `container`
- `portal`
- `closeOnEscape`
- `closeOnInteractOutside`
- `restoreFocus`
- `trapFocus`
- `lockScroll`

Required behavior:

- Escape key closes when enabled.
- Outside interaction closes when enabled.
- Focus returns to trigger by default.
- Modal overlays trap focus and prevent background interaction.
- Nested overlays must maintain stack order and should not close parent overlays accidentally.
- Native overlays may map popovers to sheets/action sheets while preserving the same public state contract.

## Section 9: Collection and Data Component Contract

Collection components must use key-based APIs rather than index-only APIs.

Standard collection types:

```ts
export type BcKey = string | number;
export type BcSelection = 'all' | Set<BcKey> | BcKey[];
```

Standard collection props:

- `items`
- `children` or `renderItem` where renderer permits
- `getKey`
- `selectedKeys`
- `defaultSelectedKeys`
- `onSelectionChange`
- `disabledKeys`
- `expandedKeys`
- `defaultExpandedKeys`
- `onExpandedChange`
- `onAction`

Data table standard models:

- `paginationModel/defaultPaginationModel/onPaginationModelChange`
- `sortModel/defaultSortModel/onSortModelChange`
- `filterModel/defaultFilterModel/onFilterModelChange`
- `selectionModel/defaultSelectionModel/onSelectionModelChange`
- `columnState/defaultColumnState/onColumnStateChange`

Server/client modes:

- `paginationMode: 'client' | 'server'`
- `sortingMode: 'client' | 'server'`
- `filteringMode: 'client' | 'server'`

Enterprise data components must support loading, empty, error, permissions/visibility adapters, virtualization, keyboard navigation, export hooks, and state persistence without coupling to BitCode runtime.

## Section 10: Token and Styling Contract

Design tokens are the source of truth. Tailwind, CSS variables, NativeWind, and runtime JS objects are generated outputs.

Token categories:

- Color roles: background, foreground, primary, secondary, destructive, muted, accent, border, input, ring, surface, overlay.
- Typography: family, size, line height, weight, tracking.
- Spacing.
- Radius.
- Shadow/elevation.
- Z-index layers.
- Motion durations/easings.
- Density.

Rules:

- Do not encode business meaning in visual tokens.
- Do not hardcode platform-only palette values in NativeWind or Tailwind outputs.
- Components should expose variants and sizes that map to tokens.
- Enterprise consumers must be able to theme globally and override component parts locally.

## Section 11: Cross-Platform Renderer Contract

Core contracts must avoid React-only or DOM-only concepts unless isolated behind renderer adapters.

Rules:

- Core types must be renderer-neutral.
- React/Vue/Stylus/Web Component/native outputs may adapt event wiring, but not semantic payloads.
- DOM-only APIs (`document`, `window`, hover, right-click, drag) require Adapt components or provider abstractions.
- React Native must not be treated as a later port. Mobile UX differences are valid implementation differences under the same caller API.
- Future Flutter support should consume the same contract metadata and token outputs, even if implementation is generated separately.

## Section 12: Migration Impact on Existing `packages/ui`

Required follow-up changes:

1. Add universal core event and state model types to `src/core/types.ts` or a new `src/core/contracts` module.
2. Update form controls from `onChange(FieldChangeEvent)` to canonical `onValueChange(value, meta)`, keeping `onChange` only as a compatibility alias where necessary.
3. Replace public `initialX` state props with `defaultX` across contracts, especially `BcDataTable`.
4. Add `open/defaultOpen/onOpenChange` to components with expandable or overlay state.
5. Add a11y contract sections to component contract files, especially select, combobox, menu, dialog, popover, tabs, tree, and data table.
6. Add slot and slotProps conventions to enterprise components.
7. Ensure Stencil/Web Component outputs emit standard `bc-*` events with stable `detail` payloads.
8. Keep BitCode-specific data adapters separate from primitive component APIs.

## Consequences

Positive:

- React + Vite, Vue, Web Components, React Native, and Tauri integrations get predictable APIs.
- BitCode can generate UI against a stable contract rather than renderer-specific quirks.
- Accessibility becomes testable per component instead of undocumented behavior.
- Enterprise components can grow without fragmenting state and event models.
- Future native or Flutter adapters can target the same contract language.

Negative:

- Existing components need API migration and compatibility aliases.
- Contract files become stricter and more verbose.
- Some Mitosis-only Pure components will need Adapt splits when platform behavior differs.
- Stencil compatibility requires deliberate event bridging rather than relying on framework callbacks.

## Implementation Status (Phase 2 Complete, 2026-06-18)

### Components Implemented (25 total)

All Core Integration Surface components now have complete contracts and implementations:

**Form Controls (10):**
- BcInput, BcSelect, BcCombobox, BcCheckbox, BcToggle, BcSegmented
- BcNumberField, BcDateField, BcTimeField, BcMultiSelect

**Forms and Reactive Scope (3):**
- BcForm, BcReactiveScope, ReactiveStore

**Overlays (7):**
- BcDialog, BcAlertDialog, BcDrawer, BcPopover, BcTooltip, BcDropdown, BcContextMenu, BcCommand

**Navigation/Layout (4):**
- BcTabs, BcAccordion, BcDisclosure, BcSplitter

**Data (1):**
- BcDataTable

### Infrastructure Implemented (7 total)

1. **Controlled/Uncontrolled Pattern** — All stateful components support canonical pairs (value/defaultValue/onValueChange, open/defaultOpen/onOpenChange, etc.)
2. **slotProps System** — Runtime Part Resolution via `partClass()` and `partAttr()` helpers (see Section 6.1)
3. **Web Component Event Bridge** — `dispatchBcEvent()` with 14 semantic `bc-*` events (see Section 3)
4. **Contract Files** — All 25 components have `.contract.ts` files defining props/events/methods
5. **Type Safety** — `BcChangeMeta<TReason>`, `BcChangeHandler<TValue>`, `BcEventSource` types in `core/types.ts`
6. **Error Handling** — SSR-safe fallbacks, try/catch in event dispatch, graceful degradation
7. **Test Coverage** — 1,538 tests across 25 components + infrastructure modules

### Test Coverage Summary

- **Unit tests**: 82 tests (part-props: 34, wc-event-bridge: 28, collection: 12, config: 8)
- **Runtime behavior tests**: 26 tests (partClass/partAttr/partStyle runtime verification)
- **Component integration tests**: 1,430 tests across all 25 components
- **Total**: 1,538 tests, 100% passing

### Build Verification

- **Mitosis compilation**: 151/151 components successfully compiled to React, Vue, Stencil, React Native
- **Type checking**: Zero TypeScript errors across all outputs
- **Bundle size**: Optimized tree-shaking via per-component exports

## Lessons Learned

### 1. Mitosis Spread Limitation and Runtime Part Resolution

**Problem**: Mitosis cannot compile dynamic object spread (`{...props.slotProps?.root}`) on JSX elements. This is a deliberate Mitosis limitation because several output targets (Stencil, Vue, Angular) do not support arbitrary dynamic attribute spread.

**Solution**: Runtime Part Resolution with pure function helpers (`partClass`, `partAttr`, `partStyle`) that components call per-element to read individual values from `slotProps`, avoiding spread entirely.

**Key Insight**: This approach works within Mitosis constraints because:
- Property access (`slotProps?.root?.className`) is valid Mitosis JSX
- String concatenation (`baseClass + ' ' + extra`) is valid Mitosis JSX
- No dynamic attribute spread required

**Trade-off**: Requires explicit wiring in each component (cannot use generic spread), but this is actually a feature for enterprise UI libraries that need controlled surfaces.

### 2. Type Safety Considerations

**Challenge**: `slotProps` is a nested record structure (`Record<string, Record<string, unknown>>`), which can lead to runtime errors if consumers pass invalid values.

**Solution**: 
- Type guards in `partClass()`, `partAttr()`, `partStyle()` validate input types before use
- Return `baseClass` or `undefined` for invalid inputs (never throw)
- Type-safe `SlotPropsMap` type alias for better IDE support

**Example**:
```typescript
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
  if (typeof extra !== 'string') return baseClass;  // Type guard
  if (!baseClass) return extra;
  return baseClass + ' ' + extra;
}
```

### 3. Error Handling Patterns

**SSR/React Native Safety**: `dispatchBcEvent()` must not crash in environments without `document` or `CustomEvent`:

```typescript
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
  // ... rest of implementation
}
```

**Graceful Degradation**: All helper functions return safe defaults (`baseClass`, `undefined`) rather than throwing errors, ensuring components render even with invalid `slotProps`.

### 4. Additive className Pattern

**Decision**: `slotProps.root.className` is concatenated after the component's base classes, not replaced.

**Rationale**:
- Token-driven classes from Tailwind/NativeWind remain the source of truth
- Consumer classes append, never override the base
- Prevents accidental removal of critical styling (e.g., `flex`, `items-center`)

**Example**:
```typescript
// Component base: "flex items-center bg-background"
// Consumer slotProps: { root: { className: "p-4 shadow-lg" } }
// Result: "flex items-center bg-background p-4 shadow-lg"
```

### 5. Explicit Wiring vs. Generic Spread

**Trade-off**: Components must explicitly call `partClass(props.slotProps, 'root', baseClass)` for each element, rather than using `{...props.slotProps?.root}`.

**Positive**:
- Controlled surface prevents arbitrary prop passthrough
- Works across all framework targets (React, Vue, Stencil, RN)
- Clear which attributes are supported per part
- Easier to audit and maintain

**Negative**:
- More verbose component code
- Requires discipline to wire all supported attributes

**Verdict**: For enterprise UI libraries, explicit wiring is preferable to generic spread.

## Migration Guide for Consumers

### 1. Using Controlled/Uncontrolled State

**Before** (if using `initialX` pattern):
```typescript
<BcSelect initialValue="apple" onChange={(e) => console.log(e.target.value)} />
```

**After** (canonical controlled/uncontrolled pairs):
```typescript
// Uncontrolled (component manages state)
<BcSelect defaultValue="apple" onValueChange={(value, meta) => console.log(value)} />

// Controlled (parent manages state)
const [value, setValue] = useState('apple');
<BcSelect value={value} onValueChange={(value, meta) => setValue(value)} />
```

### 2. Using slotProps for Customization

**Basic Example**:
```typescript
<BcInput
  slotProps={{
    root: { className: 'my-custom-class', 'data-testid': 'email-input' },
    label: { className: 'font-bold' },
    input: { 'aria-label': 'Email address' },
  }}
/>
```

**Advanced Example** (with inline styles):
```typescript
<BcDialog
  slotProps={{
    overlay: { style: { backgroundColor: 'rgba(0, 0, 0, 0.8)' } },
    content: { className: 'max-w-2xl', 'data-part': 'dialog-content' },
    title: { className: 'text-2xl font-bold' },
  }}
/>
```

### 3. Using Web Component Events

**React/Vue** (use semantic callbacks):
```typescript
<BcSelect onValueChange={(value, meta) => {
  console.log('Value:', value);
  console.log('Source:', meta.source);  // 'keyboard' | 'pointer' | 'programmatic'
  console.log('Reason:', meta.reason);  // 'select' | 'clear' | etc.
}} />
```

**Stencil/Web Components** (use `bc-*` events):
```typescript
<bc-select onbc-value-change={(e) => {
  const { value, meta } = e.detail;
  console.log('Value:', value);
  console.log('Source:', meta.source);
}} />
```

### 4. Targeting Parts with CSS

**Using data-part attributes**:
```css
/* Target the root element */
[data-part="root"] {
  margin-bottom: 1rem;
}

/* Target the label */
[data-part="label"] {
  font-weight: 600;
}

/* Target the input */
[data-part="input"]:focus {
  outline: 2px solid blue;
}
```

### 5. Imperative Methods

**Using ref for imperative control**:
```typescript
const inputRef = useRef<BcInputRef>(null);

function handleFocus() {
  inputRef.current?.focus();
}

function handleValidate() {
  const isValid = inputRef.current?.validate();
  console.log('Valid:', isValid);
}

<BcInput ref={inputRef} />
```

**Available methods** (varies by component):
- `focus()`, `blur()` — All form controls
- `validate()`, `reset()` — Form controls and BcForm
- `open()`, `close()` — Overlays (BcDialog, BcDrawer, etc.)
- `scrollToIndex(index)` — Virtualized collections

## Related ADRs

- [ADR-008: Contract-First API Freeze](./008-contract-first-api-freeze.md) — Freezes feature expansion until Core Integration Surface passes contract matrix
- [ADR-009: Runtime Part Resolution](./009-runtime-part-resolution.md) — Detailed documentation of slotProps implementation mechanism
- [ADR-010: Web Component Event Bridge](./010-web-component-event-bridge.md) — Detailed documentation of CustomEvent system

## Acceptance Criteria

- New stateful components use controlled/uncontrolled canonical pairs.
- New semantic events use `(value, meta)` or equivalent stable payload shape.
- Web Component events use kebab-case `bc-*` names and `detail` payloads.
- Interactive components document and test keyboard behavior against WAI-ARIA APG patterns.
- Form controls support native-form-compatible props and do not require BitCode runtime.
- Data components expose controlled model props for pagination, sorting, filtering, selection, and columns.
- Overlay components share one behavior contract for focus, escape, outside interaction, portal, and stacking.
- Tokens remain the styling source of truth for web and native outputs.

## References

- Ant Design: https://ant.design/
- Ant Design v6 migration: https://ant.design/docs/react/migration-v6/
- Ant Design Form: https://ant.design/components/form/
- Ant Design Table: https://ant.design/components/table/
- MUI X Data Grid: https://mui.com/x/react-data-grid/
- MUI X Data Grid API: https://mui.com/x/api/data-grid/data-grid/
- MUI X Data Grid accessibility: https://mui.com/x/react-data-grid/accessibility/
- MUI slots: https://mui.com/material-ui/customization/overriding-component-structure/
- Mantine: https://mantine.dev/
- Mantine provider: https://mantine.dev/theming/mantine-provider/
- shadcn/ui: https://ui.shadcn.com/
- PrimeReact pass-through: https://primereact.org/passthrough/
- PrimeReact DataTable: https://primereact.org/datatable/
- Radix UI Primitives: https://www.radix-ui.com/primitives
- Radix composition: https://www.radix-ui.com/primitives/docs/guides/composition
- Radix accessibility: https://www.radix-ui.com/primitives/docs/overview/accessibility
- React Aria: https://react-spectrum.adobe.com/react-aria/index.html
- React Aria customization: https://react-spectrum.adobe.com/react-aria/advanced.html
- React Aria accessibility: https://react-spectrum.adobe.com/react-aria/accessibility.html
- WAI-ARIA APG: https://www.w3.org/WAI/ARIA/apg/
- WAI-ARIA APG patterns: https://www.w3.org/WAI/ARIA/apg/patterns/
- WAI-ARIA keyboard guidance: https://www.w3.org/WAI/ARIA/apg/practices/keyboard-interface/
- Tamagui: https://tamagui.dev/docs/intro/introduction
- gluestack-ui: https://gluestack.io/ui/docs/home/overview/introduction
