# ADR-007: Enterprise Universal UI API Contract

## Status: Proposed

## Date: 2026-06-17

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
