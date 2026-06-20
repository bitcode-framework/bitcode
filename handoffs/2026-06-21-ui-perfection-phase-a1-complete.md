# Handoff: UI Perfection Roadmap — Phase A1 Runtime Tests (Complete)

**Tanggal:** 2026-06-21
**Session:** Phase A1.3 complete with all fixes applied
**Status:** 474 runtime tests pass, 24 bugs logged (6 fixed, 16 deferred, 1 wontfix, 1 deferred), A1.3 DONE

---

## Yang Sudah Selesai

### A1.3 Core Integration Surface — COMPLETE (27 test files, 474 tests)

| # | Component | Tests | Status |
|---|-----------|-------|--------|
| 1 | BcInput | 23 | ✅ |
| 2 | BcSelect | 20 | ✅ (BC-004, BC-005 deferred) |
| 3 | BcCombobox | 16 | ✅ (BC-006 deferred) |
| 4 | BcCheckbox | 13 | ✅ (BC-007 deferred) |
| 5 | BcToggle | 12 | ✅ (BC-007 deferred) |
| 6 | BcSegmented | 11 | ✅ |
| 7 | BcNumberField | 11 | ✅ |
| 8 | BcDateField | 12 | ✅ |
| 9 | BcTimeField | 11 | ✅ |
| 10 | BcMultiSelect | 12 | ✅ |
| 11 | BcDialog | 14 | ✅ (BC-008 fixed, BC-009 deferred) |
| 12 | BcAlertDialog | 14 | ✅ |
| 13 | BcDrawer | 18 | ✅ |
| 14 | BcPopover | 17 | ✅ (BC-010, BC-011 deferred) |
| 15 | BcTooltip | 13 | ✅ |
| 16 | BcDropdown | 18 | ✅ |
| 17 | BcContextMenu | 19 | ✅ (BC-012 deferred) |
| 18 | BcCommand | 20 | ✅ (BC-013 FIXED, BC-014/BC-015 deferred) |
| 19 | BcTabs | 14 | ✅ |
| 20 | BcAccordion | 15 | ✅ (BC-016 deferred) |
| 21 | BcDisclosure | 16 | ✅ |
| 22 | BcSplitter | 19 | ✅ (BC-017 wontfix) |
| 23 | BcForm | 25 | ✅ (BC-019 FIXED, BC-018 deferred) |
| 24 | ReactiveStore + BcReactiveScope | 31 | ✅ (BC-024 deferred) |
| 25 | BcList | 25 | ✅ (BC-020 FIXED) |
| 26 | BcRenderer + AppProvider | 18 | ✅ (BC-021/BC-022 FIXED, BC-022 was the P0 crash) |
| 27 | BcDataTable | 38 | ✅ (BC-022 FIXED, BC-023 deferred) |

### Bugs Fixed in This Session

| # | Bug | Fix |
|---|-----|-----|
| BC-013 | BcCommand search filtering broken | Refactored handleInput to value-first pattern |
| BC-019 | BcForm stale-state closure in setField | Inlined validation against nextValues instead of stale formValues |
| BC-020 | BcList local mode stale total | Used local newTotal variable instead of state.total |
| BC-022 | BcDataTable defaultGetRowId crash | Moved to module scope + moved can() to useStore + fixed IIFE event handlers |

### Bug Summary (24 total)

| Status | Count | Bugs |
|--------|-------|------|
| fixed | 7 | BC-001, BC-003, BC-008, BC-013, BC-019, BC-020, BC-022 |
| deferred | 16 | BC-002, BC-004, BC-005, BC-006, BC-007, BC-009, BC-010, BC-011, BC-012, BC-014, BC-015, BC-016, BC-018, BC-021, BC-023, BC-024 |
| wontfix | 1 | BC-017 (jsdom limitation) |

---

## Commits

### packages/ui
- `553af86` test(runtime): overlay + navigation + layout runtime tests (182 tests, 8 bugs logged)
- `f63fe13` test(runtime): BcForm + ReactiveStore + BcList + BcRenderer/AppProvider (93 tests, 5 bugs logged)
- `a746018` fix(datatable): BC-022 move defaultGetRowId to module scope + move can() to useStore + fix IIFE event handlers
- `5814114` fix(core): BC-013 BcCommand value-first + BC-019 BcForm stale state + BC-020 BcList stale total + improved BcReactiveScope tests

### root repo
- `e264433` docs(ui): bug log BC-010 through BC-017
- `a40a348` docs(ui): bug log BC-018 through BC-022
- (pending) docs(ui): bug log BC-023/BC-024 + status updates

---

## Yang Belum Selesai (Next Steps)

### Phase A1 — DONE
All 27 Core Integration Surface components have runtime tests. 474 tests pass.

### Remaining Deferred Bugs (Phase B)
1. **Mitosis output bugs (3):** BC-010 (IIFE crash in BcPopover), BC-011 (named slot), BC-021 (data-bc-* strip)
2. **React stale-state closures (3):** BC-015 (BcCommand keydown), BC-016 (BcAccordion onToggle), BC-024 (BcReactiveScope onChange)
3. **A11y bugs (5):** BC-004, BC-007, BC-012, BC-014, BC-018
4. **Feature gaps (3):** BC-005 (clearable), BC-006 (controlled sync), BC-009 (Escape key)
5. **Source bugs (2):** BC-023 (BcDataTable loading prop), BC-002 (slotProps attrs)

### Phase A2 — Stencil Parity Audit
### Phase A3 — Demo App Browser Verification
### Phase A4 — 4-Layer Data Source Uniformity
### Phase B1-B3 — Harden (fix all deferred bugs)
### Phase C1-C3 — Formalize

---

## Instruksi untuk Session Baru

```
Lanjutkan pekerjaan dari handoff file: handoffs/2026-06-21-ui-perfection-phase-a1-complete.md

A1.3 is COMPLETE (474 runtime tests, 24 bugs: 7 fixed, 16 deferred, 1 wontfix).
All P0 bugs fixed (BC-001, BC-003, BC-008, BC-013, BC-022).
Next: Phase A2 (Stencil parity audit) or Phase A3 (demo app browser verification),
OR fix remaining deferred bugs in Phase B.
```
