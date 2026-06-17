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
| Form control | BcInput | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract added with semantic `onValueChange`, legacy `onChange` compatibility, approved methods, a11y contract, adapter mapping, and source-level ADR-008 tests. Slots/parts are typed but not fully wired in implementation yet. |
| Form control | BcSelect | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract added with value/open controlled pairs, semantic value/open events, Web Component event bridge, methods, a11y contract, adapter mapping, and source tests. Slots/parts typed but not fully wired. |
| Form control | BcCombobox | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with value/open controlled pairs, semantic value/open events, Web Component event bridge, methods, a11y contract, adapter mapping, and source tests. Slots/parts typed but not fully wired. |
| Form control | BcCheckbox | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract added with checked/defaultChecked/onCheckedChange, semantic events, methods, a11y, adapter mapping. Slots typed not fully wired. |
| Form control | BcToggle | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract added with checked/defaultChecked/onCheckedChange, role=switch, semantic events, methods, a11y, adapter mapping. Slots typed not fully wired. |
| Form control | BcSegmented | ✅ | ✅ | – | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onValueChange, Web Component events, a11y, adapter mapping. |
| Form control | BcNumberField | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onValueChange, methods, Web Component events, a11y, adapter mapping. |
| Form control | BcDateField | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onValueChange, methods, Web Component events, a11y, adapter mapping. |
| Form control | BcTimeField | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onValueChange, methods, Web Component events, a11y, adapter mapping. |
| Form control | BcMultiSelect | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with value/open events, methods, Web Component events, a11y, adapter mapping. |
| Forms/reactive | BcForm | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with submit/validate lifecycle, form role a11y, Web Component events, methods, adapter mapping. |
| Forms/reactive | BcReactiveScope | ⚠️ | ⚠️ | ❌ | – | – | ⚠️ | ⚠️ | ADR-004 Level 2 scope; contract needs formal freeze but internal API stable. |
| Forms/reactive | ReactiveStore | ⚠️ | ⚠️ | ⚠️ | – | – | ⚠️ | ⚠️ | Headless core stable per ADR-004; public API freeze deferred to consumer-level freeze. |
| Forms/reactive | Field wrappers | ⚠️ | – | – | ⚠️ | ⚠️ | ⚠️ | ⚠️ | Composition defined in contracts (label/description/error parts per field component). Full standalone wrapper deferred. |
| Overlay | BcDialog | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract added with open/defaultOpen/onOpenChange, role=dialog, aria-modal, aria-labelledby/describedby, reason meta, Web Component events, adapter mapping. |
| Overlay | BcAlertDialog | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with defaultOpen/onOpenChange, alertdialog a11y, Web Component events, adapter mapping. |
| Overlay | BcDrawer | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with defaultOpen/onOpenChange, drawer a11y, Web Component events, adapter mapping. |
| Overlay | BcPopover | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract added with open/defaultOpen/onOpenChange, semantic reason meta, adapter mapping. |
| Overlay | BcTooltip | ✅ | ✅ | – | ⚠️ | ✅ | ✅ | ✅ | Contract added with open/defaultOpen/onOpenChange, tooltip role, adapter mapping. |
| Overlay | BcDropdown | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with open/defaultOpen/onOpenChange, onAction, Web Component events, adapter mapping. |
| Overlay | BcContextMenu | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with open/action semantic events, long-press mapping, Web Component events, adapter mapping. |
| Overlay | BcCommand | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with open/action semantic events, command palette a11y, Web Component events, adapter mapping. |
| Navigation/layout | BcTabs | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract added with activeKey/defaultActiveKey/onActiveKeyChange, tablist/tab/tabpanel roles, keyboard, adapter mapping. |
| Navigation/layout | BcAccordion | ✅ | ✅ | – | ⚠️ | ✅ | ✅ | ✅ | Contract updated with expandedKeys/defaultExpandedKeys/onExpandedChange, Web Component events, adapter mapping. |
| Navigation/layout | BcDisclosure | ✅ | ✅ | – | ⚠️ | ✅ | ✅ | ✅ | Contract updated with onOpenChange semantic events, Web Component events, adapter mapping. |
| Navigation/layout | BcSplitter | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with size/defaultSize/onSizeChange, Web Component events, separator a11y, stacked mobile mapping. |
| Data | BcDataTable | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with ADR-008 parts, methods, Web Component events, grid a11y, adapter mapping. `initialX` to `defaultX` migration still pending in source. |
| Data | BcList | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ | Contract updated with action events, list a11y, parts, methods, adapter mapping. |
| Data | Collection selection contracts | ⚠️ | ⚠️ | – | – | ⚠️ | ⚠️ | ⚠️ | BcKey/BcSelection semantics specified in ADR-007; formal types deferred to core/types. |
| Metadata | BcRenderer | ⚠️ | ⚠️ | ❌ | – | – | ⚠️ | ⚠️ | Schema-driven compositor per ADR-005; public API stable; formal ADR-008 contract deferred. |
| Metadata | FieldResolver | ⚠️ | – | – | – | – | ⚠️ | ⚠️ | Stable extension contract per ADR-005; formal ADR-008 freeze deferred. |
| Metadata | UiSchema | ⚠️ | ⚠️ | – | – | ⚠️ | ⚠️ | ⚠️ | Schema types stable per ADR-005; adapter hints deferred. |
| Infrastructure | AppProvider | ⚠️ | ⚠️ | ❌ | – | – | ⚠️ | ⚠️ | Provider contract for theme/locale/rtl/portal exists in core/context; formal ADR-008 freeze deferred. |
| Infrastructure | Config/context | ⚠️ | ⚠️ | ❌ | – | – | ⚠️ | ⚠️ | BcContext interface stable in core/types; SSR rules documented in ADR-006. |
| Infrastructure | Tokens/theme | ⚠️ | – | – | – | – | ⚠️ | ⚠️ | Pipeline exists per ADR-006 Phase 0; Flutter token output deferred. |
| Infrastructure | Portal/focus utilities | ⚠️ | ⚠️ | ⚠️ | – | ⚠️ | ⚠️ | ⚠️ | BcPortal Adapt PoC exists; centralized focus trap deferred. |
| Infrastructure | Web Component event bridge | ⚠️ | ⚠️ | – | – | – | ⚠️ | ⚠️ | Event naming convention (`bc-*` + detail payload) defined in ADR-007; per-component events typed in contracts. Implementation deferred to Stencil output post-processing. |

## Audit Notes

- This matrix intentionally starts strict. Existing implementations may already satisfy some cells, but they should not be marked ✅ until verified against ADR-007 and ADR-008.
- Prefer migrating one component family at a time and updating this file immediately after verification.
- Compatibility aliases are allowed during migration, but canonical APIs should follow ADR-007.
