# Components Contract Matrix

> Execution checklist for ADR-008 Contract-First API Freeze. This file tracks whether the Core Integration Surface has the stricter enterprise API contract required before `packages/ui` feature expansion resumes.

## Status Legend

- ✅ Done and verified.
- ⚠️ Partial, implemented under older criteria or missing some ADR-008 requirements.
- ❌ Missing or not audited.
- – Not applicable.

## Gate Definitions

| Gate | Required evidence |
|---|---|
| State | Controlled/uncontrolled pairs such as `value/defaultValue/onValueChange`, `open/defaultOpen/onOpenChange`, or model equivalents. |
| Events | Semantic callbacks plus Web Component `bc-*` event mapping where applicable. DOM/native events only as metadata. |
| Methods | Only approved imperative methods, typed in contract files and mapped per adapter when relevant. |
| Slots/Parts | Named parts, `slots`, `slotProps`, pass-through props, `as`/`asChild` where relevant. |
| A11y | Role, ARIA, keyboard table, focus behavior, labels/errors, RTL, disabled/readOnly, touch/mobile equivalent. |
| Adapter Mapping | React, Vue, React Native, Stencil/Web Component, Tauri, and future Flutter schema/contract parity notes. |
| Tests | Contract, semantic event payload, controlled/uncontrolled sync, adapter bridge, SSR, and a11y behavior tests. |

## P0 Audit Order

1. Form controls: `BcInput`, `BcSelect`, `BcCombobox`, `BcCheckbox`, `BcToggle`, `BcSegmented`, `BcNumberField`, `BcDateField`, `BcTimeField`, `BcMultiSelect`.
2. Forms/reactive: `BcForm`, `BcReactiveScope`, `ReactiveStore`, field wrappers.
3. Overlays: `BcDialog`, `BcAlertDialog`, `BcDrawer`, `BcPopover`, `BcTooltip`, `BcDropdown`, `BcContextMenu`, `BcCommand`.
4. Navigation/layout behavior: `BcTabs`, `BcAccordion`, `BcDisclosure`, `BcSplitter`.
5. Data: `BcDataTable`, `BcList`, collection selection contracts.
6. Metadata/infrastructure: `BcRenderer`, `FieldResolver`, `UiSchema`, `AppProvider`, config/context, tokens/theme, portal/focus utilities, Web Component event bridge.

## Contract Matrix

| Surface | Component / Module | State | Events | Methods | Slots/Parts | A11y | Adapter Mapping | Tests | Notes |
|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| Form control | BcInput | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Contract + source fully wired per ADR-007 §6.1 Runtime Part Resolution. `partClass` + `data-part` on root, label, control, input, description, error parts. slotProps passthrough verified in build (151 components, 4 targets). |
| Form control | BcSelect | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `partClass` + `data-part="root"` wired per ADR-007 §6.1. Controlled pairs, semantic events, WC bridge, methods, a11y, adapter mapping verified. |
| Form control | BcCombobox | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `partClass` + `data-part="root"` wired per ADR-007 §6.1. Controlled pairs, semantic events, WC bridge, methods, a11y, adapter mapping verified. |
| Form control | BcCheckbox | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `partClass` + `data-part` on root, control parts per ADR-007 §6.1. Checked/defaultChecked/onCheckedChange, semantic events, methods, a11y, adapter mapping verified. |
| Form control | BcToggle | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `partClass` + `data-part="root"` wired per ADR-007 §6.1. Checked/defaultChecked/onCheckedChange, role=switch, semantic events, methods, a11y, adapter mapping verified. |
| Form control | BcSegmented | ✅ | ✅ | – | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onValueChange, Web Component events, a11y, adapter mapping. |
| Form control | BcNumberField | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onValueChange, methods, Web Component events, a11y, adapter mapping. |
| Form control | BcDateField | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onValueChange, methods, Web Component events, a11y, adapter mapping. |
| Form control | BcTimeField | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onValueChange, methods, Web Component events, a11y, adapter mapping. |
| Form control | BcMultiSelect | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with value/open events, methods, Web Component events, a11y, adapter mapping. |
| Forms/reactive | BcForm | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with submit/validate lifecycle, form role a11y, Web Component events, methods, adapter mapping. |
| Forms/reactive | BcReactiveScope | ⚠️ | ⚠️ | ❌ | – | – | ⚠️ | ⚠️ | ADR-004 Level 2 scope; contract needs formal freeze but internal API stable. |
| Forms/reactive | ReactiveStore | ⚠️ | ⚠️ | ⚠️ | – | – | ⚠️ | ⚠️ | Headless core stable per ADR-004; public API freeze deferred to consumer-level freeze. |
| Forms/reactive | Field wrappers | ⚠️ | – | – | ⚠️ | ⚠️ | ⚠️ | ⚠️ | Composition defined in contracts (label/description/error parts per field component). Full standalone wrapper deferred. |
| Overlay | BcDialog | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `partClass` + `data-part` on overlay, content parts per ADR-007 §6.1. open/defaultOpen/onOpenChange, role=dialog, aria-modal, reason meta, WC events, adapter mapping verified. |
| Overlay | BcAlertDialog | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with defaultOpen/onOpenChange, alertdialog a11y, Web Component events, adapter mapping. |
| Overlay | BcDrawer | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with defaultOpen/onOpenChange, drawer a11y, Web Component events, adapter mapping. |
| Overlay | BcPopover | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `partClass` + `data-part="root"` wired per ADR-007 §6.1. open/defaultOpen/onOpenChange, semantic reason meta, adapter mapping verified. |
| Overlay | BcTooltip | ✅ | ✅ | – | ⚠️ | ✅ | ✅ | ✅ | Contract added with open/defaultOpen/onOpenChange, tooltip role, adapter mapping. |
| Overlay | BcDropdown | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with open/defaultOpen/onOpenChange, onAction, Web Component events, adapter mapping. |
| Overlay | BcContextMenu | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with open/action semantic events, long-press mapping, Web Component events, adapter mapping. |
| Overlay | BcCommand | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with open/action semantic events, command palette a11y, Web Component events, adapter mapping. |
| Navigation/layout | BcTabs | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract added with activeKey/defaultActiveKey/onActiveKeyChange, tablist/tab/tabpanel roles, keyboard, adapter mapping. |
| Navigation/layout | BcAccordion | ✅ | ✅ | – | ⚠️ | ✅ | ✅ | ✅ | Contract updated with expandedKeys/defaultExpandedKeys/onExpandedChange, Web Component events, adapter mapping. |
| Navigation/layout | BcDisclosure | ✅ | ✅ | – | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onOpenChange semantic events, Web Component events, adapter mapping. |
| Navigation/layout | BcSplitter | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with size/defaultSize/onSizeChange, Web Component events, separator a11y, stacked mobile mapping. |
| Data | BcDataTable | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract and source migrated to ADR-007 controlled/uncontrolled pairs: `sortModel/defaultSortModel/onSortModelChange`, `filterModel/defaultFilterModel/onFilterModelChange`, `selectedKeys/defaultSelectedKeys`, `expandedKeys/defaultExpandedKeys`, `onColumnStateChange`. Legacy `initialX` and `onSortChange/onFilterChange/onColumnsChange` kept as deprecated aliases. WC events renamed to `bc-sort-model-change`, `bc-filter-model-change`, `bc-column-state-change`. |
| Data | BcList | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with action events, list a11y, parts, methods, adapter mapping. |
| Data | Collection selection contracts | ✅ | ⚠️ | – | – | ⚠️ | ⚠️ | ⚠️ | `BcKey` and `BcSelection` types added to `core/types.ts` per ADR-007 §9. Consumer-level integration pending. |
| Metadata | BcRenderer | ⚠️ | ⚠️ | ❌ | – | – | ⚠️ | ⚠️ | Schema-driven compositor per ADR-005; public API stable; formal ADR-008 contract deferred. |
| Metadata | FieldResolver | ⚠️ | – | – | – | – | ⚠️ | ⚠️ | Stable extension contract per ADR-005; formal ADR-008 freeze deferred. |
| Metadata | UiSchema | ⚠️ | ⚠️ | – | – | ⚠️ | ⚠️ | ⚠️ | Schema types stable per ADR-005; adapter hints deferred. |
| Infrastructure | AppProvider | ⚠️ | ⚠️ | ❌ | – | – | ⚠️ | ⚠️ | Provider contract for theme/locale/rtl/portal exists in core/context; formal ADR-008 freeze deferred. |
| Infrastructure | Config/context | ⚠️ | ⚠️ | ❌ | – | – | ⚠️ | ⚠️ | BcContext interface stable in core/types; SSR rules documented in ADR-006. |
| Infrastructure | Tokens/theme | ⚠️ | – | – | – | – | ⚠️ | ⚠️ | Pipeline exists per ADR-006 Phase 0; Flutter token output deferred. |
| Infrastructure | Portal/focus utilities | ⚠️ | ⚠️ | ⚠️ | – | ⚠️ | ⚠️ | ⚠️ | BcPortal Adapt PoC exists; centralized focus trap deferred. |
| Infrastructure | Web Component event bridge | ✅ | ✅ | – | – | – | ✅ | ✅ | `core/wc-event-bridge.ts` implemented with `dispatchBcEvent()`, `buildMeta()`, and `BC_EVENTS` constants. SSR/RN safe. Covers all ADR-007 semantic events. Per-component wiring into Mitosis source is incremental (components call `dispatchBcEvent` alongside semantic callbacks). |
| Infrastructure | Runtime Part Resolution | ✅ | – | – | ✅ | – | ✅ | ✅ | `core/part-props.ts` implemented per ADR-007 §6.1 with `partClass`, `partAttr`, `partStyle`. 7 source components wired: BcInput, BcSelect, BcCombobox, BcCheckbox, BcToggle, BcDialog, BcPopover. className additive, data-part emitted, Mitosis 4-target build verified. |

## Audit Notes

- This matrix intentionally starts strict. Existing implementations may already satisfy some cells, but they should not be marked ✅ until verified against ADR-007 and ADR-008.
- Prefer migrating one component family at a time and updating this file immediately after verification.
- Compatibility aliases are allowed during migration, but canonical APIs should follow ADR-007.
