# ADR-008: Contract-First API Freeze for `packages/ui`

## Status: Proposed

## Date: 2026-06-17

## Depends on: [ADR-004](./004-reactive-engine.md), [ADR-005](./005-metadata-driven-ui.md), [ADR-006](./006-universal-component-strategy.md), [ADR-007](./007-enterprise-universal-ui-api-contract.md)

## Context

`packages/ui` has strong architecture direction, broad component coverage, a reactive engine, metadata renderer, multi-target output, and token infrastructure. However, the current progress tracker overstates enterprise readiness because many components are marked Contract/A11y complete under older criteria that do not yet require the stricter API surface defined in ADR-007.

The goal is not to add more components first. The goal is to harden the Core Integration Surface so `packages/ui` becomes a genuinely standalone, enterprise-grade, universal UI library that can integrate cleanly with React, Vue, React Native, Stencil/Web Components, Tauri, BitCode, and a future Flutter adapter.

Decisions already resolved in the grilling session:

1. Use a **Contract-First API Freeze** before further feature expansion.
2. Prioritize the **Core Integration Surface**, not all 149 components equally.
3. Use **Semantic Events** as the primary event contract. DOM/native events are metadata only.
4. Require controlled/uncontrolled state pairs for every stateful component.
5. Treat React, Vue, React Native, Stencil/Web Components, and future Flutter as adapters over neutral contracts.
6. Promise Flutter schema/contract parity, not JSX-style API parity.
7. Require a Contract Matrix for each P0 component.
8. Require slots/slotProps/pass-through APIs for P0 complex components, optional for simple primitives.
9. Restrict imperative methods to capabilities that are awkward to express declaratively.
10. Make A11y Behavior Contract mandatory for P0 interactive components.

## Decision

Freeze feature expansion for `packages/ui` until the Core Integration Surface passes the Contract Matrix gate.

During the freeze, work may continue only if it improves one of these hardening goals:

- Public API contract completeness.
- Semantic event consistency.
- Controlled/uncontrolled state consistency.
- Adapter mapping consistency.
- Accessibility behavior contracts.
- Compatibility tests.
- Documentation alignment in the ADR index and component progress tracker.

New non-critical components and broad P1/P2 expansion should wait until the P0 freeze gate passes.

## Core Integration Surface

The freeze applies first to this P0 surface:

- **Form controls:** `BcInput`, `BcSelect`, `BcCombobox`, `BcCheckbox`, `BcToggle`, `BcSegmented`, `BcNumberField`, `BcDateField`, `BcTimeField`, `BcMultiSelect`.
- **Forms and reactive scope:** `BcForm`, `BcReactiveScope`, `ReactiveStore`, form field wrappers.
- **Overlays:** `BcDialog`, `BcAlertDialog`, `BcDrawer`, `BcPopover`, `BcTooltip`, `BcDropdown`, `BcContextMenu`, `BcCommand`.
- **Navigation/layout behavior:** `BcTabs`, `BcAccordion`, `BcDisclosure`, `BcSplitter`.
- **Data:** `BcDataTable`, `BcList`, collection selection contracts.
- **Metadata rendering:** `BcRenderer`, `FieldResolver`, `UiSchema`.
- **Infrastructure:** `AppProvider`, config/context, tokens/theme, portal/focus utilities, Web Component event bridge.

Simple visual components such as badges, cards, avatars, statistics, and empty states may be hardened later unless they block the Core Integration Surface.

## Contract Matrix Gate

Each P0 item must have a contract matrix that covers:

| Area | Required content |
|---|---|
| State props | Controlled and uncontrolled pairs such as `value/defaultValue/onValueChange`, `open/defaultOpen/onOpenChange`, or model equivalents. |
| Semantic events | Stable renderer-neutral callbacks and Web Component `bc-*` event mapping. |
| Imperative methods | Only allowed methods, such as `focus`, `blur`, `validate`, `reset`, `submit`, `scrollToIndex`, `scrollToKey`, `open`, `close`, `measure`. |
| Slots and parts | Named parts, `slots`, `slotProps`, pass-through props, `as`/`asChild` where relevant. |
| A11y behavior | Roles, ARIA attributes, keyboard table, focus behavior, labels/descriptions/errors, RTL, disabled/readOnly, touch/mobile equivalent. |
| Form compatibility | Native form behavior on web, hidden inputs where needed, name/required/invalid/error semantics. |
| Adapter mapping | React, Vue, React Native, Stencil/Web Component, Tauri, and future Flutter schema/contract mapping. |
| SSR safety | No unguarded `window`, `document`, layout, storage, or browser-only APIs. |
| Tests | Contract tests, source-level behavior tests, adapter smoke tests, and event payload tests. |

A component is not enterprise-ready until its matrix is complete and verified.

## Event Policy

Semantic Events are the primary contract.

Examples:

```ts
onValueChange(value, meta)
onOpenChange(open, meta)
onSelectionChange(selection, meta)
onSortModelChange(sortModel, meta)
```

Web Components must emit corresponding kebab-case CustomEvents:

```ts
bc-value-change: { detail: { value, meta } }
bc-open-change: { detail: { open, meta } }
bc-selection-change: { detail: { selection, meta } }
```

DOM or native events may be included as `meta.nativeEvent`, but consumers must not need to parse `event.target.value` or platform-specific payloads to use the component.

## Adapter Policy

Neutral TypeScript contracts are the source of truth. Framework outputs are adapters.

- React and Vue map contracts to idiomatic props/callbacks.
- React Native maps contracts to native primitives and mobile UX where needed.
- Stencil/Web Components map contracts to attributes/properties and `bc-*` CustomEvents.
- Tauri uses the web adapter.
- Future Flutter maps the same schema, tokens, semantic events, and reactive model into Dart idioms. It does not need JSX prop parity.

## A11y Policy

A11y Behavior Contract is required for P0 interactive components.

Each P0 interactive component must define:

- Role and ARIA attribute model.
- Keyboard interaction table.
- Focus entry, movement, trap, restore, and escape behavior.
- Label, description, error, and live-region behavior.
- RTL behavior.
- Disabled and read-only semantics.
- Touch/mobile equivalent.

WAI-ARIA APG remains the baseline for web behavior. React Native and Flutter adapters must preserve the semantic behavior even when platform APIs differ.

## Progress Tracker Alignment

`docs/components-progress.md` must stop treating old `Contract ✅` and `A11y ✅` as final enterprise readiness for the Core Integration Surface. `docs/components-contract-matrix.md` is the execution checklist for this freeze. During this freeze, the tracker must include a dedicated Contract-First API Freeze gate with status columns for:

- Matrix.
- Semantic Events.
- Controlled State.
- Slots/Parts.
- Methods.
- A11y Behavior.
- Adapter Mapping.
- Tests.

Existing component completion status remains historical implementation progress, not the final enterprise-ready claim under ADR-008.

## Consequences

Positive:

- Prevents repeating the `packages/components` problem where broad component count hid integration friction.
- Makes React/Vite, React Native, Web Component, BitCode, and future Flutter integration predictable.
- Aligns reactive engine and metadata renderer with a stable public component surface.
- Gives future contributors an explicit reason why feature expansion is paused.

Negative:

- Slows visible feature expansion.
- Requires revisiting components previously marked complete.
- May introduce compatibility aliases for older APIs while moving to semantic events.
- Requires additional contract and adapter tests before components can be called enterprise-ready.

## Acceptance Criteria

- `docs/components-progress.md` includes the Contract-First API Freeze gate.
- All Core Integration Surface components have a completed Contract Matrix.
- Semantic event payloads are standardized and tested.
- Controlled/uncontrolled state pairs are standardized and tested.
- P0 complex components define slots/parts/pass-through rules.
- P0 interactive components define A11y Behavior Contracts.
- Web Component event bridge emits `bc-*` events with semantic `detail` payloads.
- React Native and future Flutter mappings are documented at the contract level.
- No new non-critical P1/P2 expansion is treated as priority until this gate passes.
