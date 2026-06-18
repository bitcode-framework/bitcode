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
| Form control | BcSegmented | ✅ | ✅ | – | ✅ | ✅ | ✅ | ✅ | Source implementation created. partClass + data-part on root, option. |
| Form control | BcNumberField | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Source implementation created. partClass + data-part on root, label, input, suffix. |
| Form control | BcDateField | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Source implementation created. partClass + data-part on root, label, input. |
| Form control | BcTimeField | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Source implementation created. partClass + data-part on root, label, input. |
| Form control | BcMultiSelect | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Source implementation created. partClass + data-part on root, label, trigger, chip, content, option. |
| Forms/reactive | BcForm | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Contract + source wired per ADR-007 §6.1. `partClass` + `data-part` on root, field, label, control, description, error, actions, submitButton parts. Submit/validate lifecycle, form a11y, WC events, methods, adapter mapping verified. |
| Forms/reactive | BcReactiveScope | ✅ | ✅ | ✅ | – | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: Parts (root), Methods (dispose), A11y (group role), Adapter mapping for all 6 targets. Internal API stable per ADR-004 Level 2. |
| Forms/reactive | ReactiveStore | ✅ | ✅ | ✅ | – | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: Methods, a11y (headless), adapter mapping (headless-pure-ts-portable). |
| Forms/reactive | Field wrappers | ✅ | – | – | ✅ | ✅ | ✅ | ✅ | Composition wired through BcForm field/label/control/description/error parts. All parts support slotProps passthrough. |
| Overlay | BcDialog | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `partClass` + `data-part` on overlay, content parts per ADR-007 §6.1. open/defaultOpen/onOpenChange, role=dialog, aria-modal, reason meta, WC events, adapter mapping verified. |
| Overlay | BcAlertDialog | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, overlay, content, title, description, confirmButton, cancelButton. defaultOpen/onOpenChange, alertdialog a11y, WC events, adapter mapping verified. |
| Overlay | BcDrawer | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, overlay, content, header, title, description, body, closeButton. defaultOpen/onOpenChange, drawer a11y, WC events, adapter mapping verified. |
| Overlay | BcPopover | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `partClass` + `data-part="root"` wired per ADR-007 §6.1. open/defaultOpen/onOpenChange, semantic reason meta, adapter mapping verified. |
| Overlay | BcTooltip | ✅ | ✅ | – | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root/content. open/defaultOpen/onOpenChange, tooltip role, adapter mapping verified. |
| Overlay | BcDropdown | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, trigger, content, item, separator. open/defaultOpen/onOpenChange, onAction, WC events, adapter mapping verified. |
| Overlay | BcContextMenu | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, backdrop, content, item, separator. open/action semantic events, long-press mapping, WC events, adapter mapping verified. |
| Overlay | BcCommand | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, input, list, group, item, empty. open/action semantic events, command palette a11y, WC events, adapter mapping verified. |
| Navigation/layout | BcTabs | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, tablist, tab. activeKey/defaultActiveKey/onActiveKeyChange, tablist/tab roles, keyboard, adapter mapping verified. |
| Navigation/layout | BcAccordion | ✅ | ✅ | – | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, item, trigger, content. expandedKeys/defaultExpandedKeys/onExpandedChange, WC events, adapter mapping verified. |
| Navigation/layout | BcDisclosure | ✅ | ✅ | – | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, trigger, content. onOpenChange semantic events, WC events, adapter mapping verified. |
| Navigation/layout | BcSplitter | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Runtime Part Resolution wired on root, startPane, handle, endPane. size/defaultSize/onSizeChange, WC events, separator a11y, stacked mobile mapping verified. |
| Data | BcDataTable | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Contract/source migrated to ADR-007 controlled pairs and wired per ADR-007 §6.1. `partClass` + `data-part` on root, toolbar, search, header, headerCell, body, row, cell, loading, error, empty, pagination, contextMenu. Legacy aliases kept deprecated. WC events use `bc-sort-model-change`, `bc-filter-model-change`, `bc-column-state-change`. |
| Data | BcList | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Contract + source wired per ADR-007 §6.1. `partClass` + `data-part` on root, item, pagination, empty, loading, error parts. Action events, list a11y, methods, adapter mapping verified. |
| Data | Collection selection contracts | ✅ | ✅ | ✅ | – | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: `BcKey` and `BcSelection` types in `core/types.ts` per ADR-007 §9. Methods (scrollToKey), WC events (bc-selection-change, bc-expanded-change, bc-action), A11y (listbox/grid/tree/menu roles), Adapter mapping for all 6 targets. |
| Metadata | BcRenderer | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: Parts (root + 6 schema kinds), adapter mapping, a11y contract. Phases 1-3 implemented. |
| Metadata | FieldResolver | ✅ | – | ✅ | – | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: Methods (resolve), A11y (pure function resolver - no rendering), Adapter mapping for all 6 targets. Stable extension contract per ADR-005. |
| Metadata | UiSchema | ✅ | ✅ | – | ✅ | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: Parts (root, field, section, grid), WC events (bc-renderer-change), A11y (form/presentation role depending on schema type), Adapter mapping for all 6 targets. Schema types stable per ADR-005. |
| Infrastructure | AppProvider | ✅ | ✅ | ❌ | – | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: Parts, a11y (htmlDir/htmlLang), adapter mapping. |
| Infrastructure | Config/context | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: Parts (root), Methods (setLocale, setTheme, setRtl), WC events (bc-config-change), A11y (htmlLangAttribute, htmlDirAttribute for screen readers), Adapter mapping for all 6 targets. BcContext interface stable in core/types; SSR rules documented in ADR-006. |
| Infrastructure | Tokens/theme | ✅ | – | – | – | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: A11y (WCAG 2.1 AA minimum contrast compliance, RTL support via CSS logical properties), Adapter mapping for all 6 targets (Flutter pending). Pipeline exists per ADR-006 Phase 0; CI drift detection via `npm run tokens:check`. |
| Infrastructure | Portal/focus utilities | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Formal ADR-008 contract frozen: Parts (root, content), a11y (focus-trap/restore), adapter mapping. |
| Infrastructure | Web Component event bridge | ✅ | ✅ | – | – | – | ✅ | ✅ | `core/wc-event-bridge.ts` wired into 19 source components (BcTooltip skipped: CSS-only hover). SSR/RN safe. |
| Infrastructure | Runtime Part Resolution | ✅ | – | – | ✅ | – | ✅ | ✅ | `core/part-props.ts` implemented per ADR-007 §6.1 with `partClass`, `partAttr`, `partStyle`. 25 Core Integration Surface components wired across forms, overlays, navigation/layout, and data/view surfaces. className additive, data-part emitted, Mitosis 4-target build verified. |

## Audit Notes

- This matrix intentionally starts strict. Existing implementations may already satisfy some cells, but they should not be marked ✅ until verified against ADR-007 and ADR-008.
- Prefer migrating one component family at a time and updating this file immediately after verification.
- Compatibility aliases are allowed during migration, but canonical APIs should follow ADR-007.
