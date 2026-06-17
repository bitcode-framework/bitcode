# Handoff: ADR-008 Contract-First API Freeze

**Date:** 2026-06-18
**Topic:** Contract-First API Freeze execution for `packages/ui` Core Integration Surface

## Apa yang sudah dikerjakan

### Research & Planning
- Riset mendalam Ant Design v6, MUI/MUI X, Mantine, shadcn/ui, PrimeReact, Radix UI, React Aria, WAI-ARIA APG, Tamagui, gluestack-ui.
- Grilling session (13 questions) untuk kunci keputusan: Contract-First API Freeze, Core Integration Surface, Semantic Events, controlled/uncontrolled, adapter policy, Flutter parity, Contract Matrix, slots/parts, imperative methods, a11y behavior contract.

### ADRs Created
- **ADR-007** (`docs/adr/007-enterprise-universal-ui-api-contract.md`): Standard event/state/method/prop/a11y contracts.
- **ADR-008** (`docs/adr/008-contract-first-api-freeze.md`): Freeze gate process, Core Integration Surface definition, Contract Matrix requirements, acceptance criteria.

### Documentation
- `CONTEXT.md` (root): glossary with Contract-First API Freeze, Core Integration Surface, Semantic Event terms.
- `docs/components-contract-matrix.md`: execution checklist with per-component gate status (30+ rows).
- `docs/components-progress.md`: ADR-008 freeze section, status legend updated.
- `docs/adr/README.md`: ADR-007/008 in index, dependency graph updated.

### Contract Hardening (30+ components)
**Source updated (7 components):**
- `BcInput`, `BcSelect`, `BcCombobox`, `BcCheckbox`, `BcToggle`, `BcDialog`, `BcPopover`

**Contract updated with ADR-008 surface (23+ components):**
- Overlays: `BcAlertDialog`, `BcDrawer`, `BcContextMenu`, `BcCommand`, `BcTooltip`, `BcDropdown`
- Forms: `BcSegmented`, `BcNumberField`, `BcDateField`, `BcTimeField`, `BcMultiSelect`
- Navigation/layout: `BcTabs`, `BcAccordion`, `BcDisclosure`, `BcSplitter`
- Data: `BcDataTable`, `BcList`
- Views: `BcForm`
- Core: `core/types.ts` (BcChangeMeta, BcChangeHandler, BcEventSource)

Each contract now includes: semantic event types (reason + source), Web Component event bridge (`bc-*`), parts/slots interfaces, imperative methods, a11y behavior contract, adapter mapping (React/Vue/RN/Stencil/Tauri/Flutter).

### Tests
- 1388/1388 source-level tests pass (up from 1354).
- Build green: 151 components, 149 exports, 4 targets (React/Vue/Stencil/RN).

## Apa yang belum selesai

### ⚠️ Slots/Parts wiring (all P0)
- Parts interfaces typed in contracts but `slotProps` passthrough not yet rendered in Mitosis output. Mitosis does not have native slot composition like React.

### ⚠️ BcDataTable source migration
- Contract updated but source still uses `initialSort/initialFilters/initialSelected` instead of `defaultSortModel/defaultFilterModel/defaultSelectedKeys`. Migration deferred.

### ⚠️ ReactiveStore / BcReactiveScope formal freeze
- ADR-004 owns the reactive model. Public API stable but formal contract freeze deferred to consumer-level.

### ⚠️ BcRenderer formal freeze
- ADR-005 owns the schema renderer. Public API stable but formal contract freeze deferred.

### ⚠️ Collection selection types
- `BcKey`, `BcSelection` semantics specified in ADR-007 but not yet added to `core/types.ts`.

### ⚠️ Web Component event bridge implementation
- Event naming convention (`bc-*` + detail payload) defined in contracts. Actual CustomEvent dispatch in Stencil output not yet implemented.

### ⚠️ P1/P2 components
- Not part of Core Integration Surface. Contract hardening for ~80 remaining components deferred.

## File yang diubah

### Main repo (committed `3f20b54`)
- `CONTEXT.md` (new)
- `docs/adr/007-enterprise-universal-ui-api-contract.md` (new)
- `docs/adr/008-contract-first-api-freeze.md` (new)
- `docs/adr/README.md`
- `docs/components-contract-matrix.md` (new)
- `docs/components-progress.md`
- `packages/ui` (submodule ref updated)

### UI submodule (committed `08377a0`)
- `src/core/types.ts`
- 7 source files updated (BcInput, BcSelect, BcCombobox, BcCheckbox, BcToggle, BcDialog, BcPopover)
- 23+ contract files updated
- 7 new test files
- 7 new contract files

## Next steps

1. **Slots/Parts wiring**: Investigate Mitosis slot composition patterns. May need `renderProp` or conditional rendering based on `slotProps`.
2. **BcDataTable source migration**: Replace `initialX` with `defaultX` controlled pairs in source.
3. **Collection types**: Add `BcKey`, `BcSelection` to `core/types.ts`.
4. **WC event bridge**: Add CustomEvent dispatch post-processing in Stencil build step.
5. **ReactiveStore/BcRenderer freeze**: Formal contract when consumer patterns stabilize.
6. **P1/P2 contracts**: Apply same ADR-008 pattern to remaining components.

## Konteks penting

- `packages/ui` is a git submodule: `https://github.com/bitcode-framework/ui.git`
- UI repo must be committed first, then main repo updates submodule ref.
- Tests run via `node --test 'src/**/*.test.ts'` (npm test uses `node` which may not be on PATH in PowerShell).
- Build via `cmd.exe /c "set PATH=C:\Program Files\nodejs;%PATH%&& npm run build"`.
- ADR-004/005 already own reactivity and schema rendering. ADR-007/008 own component API surface only.
