# Handoff: BcDataTable A1-A4 Enterprise Engines Complete — Render Integration + Build Fix Remaining

**Tanggal:** 2026-06-15
**Sesi sebelumnya:** Opsi A (4 subset commit) — BcDataTable enterprise upgrade per ADR-003 §4
**Submodule HEAD:** `packages/ui` @ `da9d297`
**Parent HEAD:** `brainplusplus/bitcode` @ `1616611`

---

## Yang Sudah Selesai (verified)

### Fase A -> Fase C + BcDataTable A1-A4

Total **447 tests pass / 11 fail** (11 = pre-existing BcStack, lihat blok "Blocker" di bawah).

| Fase | Commit submodule | Konten |
|---|---|---|
| Fase A | `34adcdd` | Cross-tier infra: providers/hooks/animation/adapt contracts |
| Fase B | `1cfbe0c` | ReactiveStore scope resolution + Evaluator + BcReactiveScope + BcForm Level 1 integration |
| Fase C batch 1 | `8407c94` | 5 P0 fondasi primitives: BcSlot, BcPortal (Adapt), BcVisuallyHidden, BcBackdrop, BcLoadingOverlay |
| Fase C batch 2 | `65018fa` | 6 skeleton + contract+test: BcDescriptions, BcSteps, BcAccordion, BcStatistic, BcEmpty, BcResult |
| Fase C batch 3 | `09e75e5` | 6 skeleton + contract+test: BcContextMenu, BcDrawer (source fix bg-black to bg-foreground), BcTransfer, BcCascader, BcTreeSelect, BcOtp |
| Fase C khusus | `a6863c8` | BcForm + BcList contracts+tests |
| Fase C subset 1 | `6c85fc1` | BcSplitter + BcXStack/BcYStack/BcBox + BcToastQueue + BcToastQueueController |
| DT A1 | `085b117` | BcDataTable multi-sort + filter + global search + a11y aria-sort |
| DT A2 | `895637d` | BcDataTable column ops + selection (checkbox/single/multi + shift-range) |
| DT A3 | `03cb549` | BcDataTable engines: tree + group + virtualize (state wired, render deferred) |
| DT A4 | `da9d297` | BcDataTable engines: edit + export + persistence + ctx menu (state wired, render deferred) |

### BcDataTable 10 Pure Engines

All in `packages/ui/src/components/data/datatable/engine/`:

1. `sort.ts` — multi-rule comparator, 3-state cycle, locale + null + Date + NaN aware (15 tests)
2. `filter.ts` — 12 operators + global search (16 tests)
3. `columns.ts` — ColumnState: show-hide, reorder, resize, pin (13 tests)
4. `selection.ts` — single/multi/checkbox modes, ctrl-toggle, shift-range from anchor, header tri-state (16 tests)
5. `tree.ts` — depth-first flatten with cycle break at 64 levels, expandAll/collapseAll (10 tests)
6. `group.ts` — stable bucket by first-occurrence, aggregates sum/avg/count/min/max with Date+number coercion (8 tests)
7. `virtualize.ts` — computeWindow math: overscan, scrollTop clamping, paddingTop+visible+paddingBottom invariant (8 tests)
8. `edit.ts` — cell + row mode, staged per-row buffer, commit/cancel returns diffs (13 tests)
9. `export.ts` — RFC 4180 CSV with UTF-8 BOM, custom delimiter/newline/header, per-column format() (11 tests)
10. `persistence.ts` — versioned snapshot, injectable Storage, deserializeState null-safe (11 tests)

### P0 Status: 22/22 components have engine + state + contract + source-level tests

All 16 ADR-003 Section 4 must-have features now have engine + state wiring.
Render integration for A3/A4 (virtualized scroll, tree chevron, group rows, cell input) is deferred.

---

## BLOCKER: 11 Pre-Existing BcStack Test Failures

### Gejala

```
not ok 286 - BcStack react: source exists
not ok 287 - BcStack react: container element is <div
...
# tests 458
# pass 447
# fail 11
```

### Akar Masalah

`BcStack.test.ts` (existing since commit `ffdbfb2`, Phase 0 ADR-006) tests against **Mitosis build output** in `packages/ui/output/`. That directory does not exist in the working tree because `npm run build` has never been run.

Root cause: `npm run build` fails with:
```
npm error could not determine executable to run
(running `npx mitosis build`)
```

`npx mitosis build` cannot find the `mitosis` binary in `node_modules/.bin/`. Possible causes:
1. `@builder.io/mitosis-cli` dependency not properly installed.
2. `node_modules/` stale or incomplete.
3. `npx` wrapper cannot resolve the bin entry point.

### Cara Fix (3 options, top priority first)

**Option 1 (recommended): Fix install + build**
```powershell
cd packages/ui
npm install
npx mitosis --version
npm run build
npm test
```

**Option 2: Tag BcStack test as skip-when-no-output**
Modify `BcStack.test.ts` to skip (not fail) when output files are missing. Reduces safety net.

**Option 3: Commit build output**
Add `output/` to git. Not recommended (repo bloat).

---

## Yang Belum Selesai (Prioritas Kerja)

### Prioritas 1: Fix Build Blocker

Follow Option 1 in BLOCKER section above. This is the only thing causing 11 test failures.

### Prioritas 2: BcDataTable A3-Render (Virtualized Scroll + Tree + Group Rows)

Engine complete (A3 commit `03cb549`), render not yet integrated:

1. **Virtualized scroll container** — render `<div onScroll>` with paddingTop + paddingBottom from `state.virtualWindow`, render only `state.virtualRows`.

2. **Tree-table rendering** — per `TreeRow` in `state.treeRows`:
   - Indent based on `row.depth * 16px`.
   - Chevron expand/collapse: collapsed shows triangle-right, expanded shows triangle-down, gated on `hasChildren`.
   - `onClick` calls `state.toggleTreeRow(row.key)`.

3. **Group header + footer rows** — per `GroupedRow` in `state.groupedRows`:
   - Header row with colSpan showing group key.
   - Footer row (when `showGroupFooter`): aggregate values per column.

Estimated effort: ~400-600 lines source rewrite in body section. Engine contracts unchanged.

### Prioritas 3: BcDataTable A4-Render (Cell Input + Row Edit UI)

Engine complete (A4 commit `da9d297`), render not yet integrated:

1. **Cell edit input** — when `editMode === 'cell'` and `state.cellIsEditing(rowId, field)`:
   - Render `<input>` instead of `<span>`.
   - `onInput` calls `state.stageEdit`.
   - `onBlur` calls `state.commitCellEdit`.
   - `onKeyDown Escape` calls `state.cancelCellEdit`.

2. **Row edit mode** — when `editMode === 'row'` and `state.rowIsEditing(rowId)`:
   - All cells render as inputs.
   - Confirm/Cancel buttons.
   - Confirm calls `state.commitRowEdit`. Cancel calls `state.cancelRowEdit`.

3. **Context menu visual** — `showRowContextMenu` already wires `onContextMenu` on `<tr>`, but no menu UI renders. Integrate with `BcContextMenu` component or native browser menu.

Estimated effort: ~300-500 lines source + ~100 lines test.

### Prioritas 4: Adapt Splits for 4 Components

Adapt-marked components still Pure `.lite.tsx` (need `.web.tsx` + `.native.tsx`):
- `#7 BcContextMenu` — right-click web vs long-press actionsheet mobile
- `#9 BcDrawer` — slide panel web vs bottom sheet mobile
- `#11 BcCascader` — dropdown columns web vs picker sheet mobile
- `#12 BcTreeSelect` — dropdown web vs sheet mobile

Pattern exists in `adapt-poc/popover/` — follow same: `BcXxx.contract.ts` + `BcXxx.web.tsx` + `BcXxx.native.tsx` + `index.ts` + `BcXxx.test.ts`.

### Prioritas 5: Fase D (ADR-005 BcRenderer + UiSchema)

Separate ADR. BcRenderer = metadata-driven UI that consumes `views/*.json` + `models/*.json` from BitCode engine to render UI at runtime. Needs:
- `UiSchema` type + `FieldNode` types
- Schema to component tree compositor
- BitCode JSON to UiSchema transformer
- Reactive engine integration (ADR-004 already ready from Fase B)

---

## File yang Berubah di Sesi Ini

### Submodule `packages/ui` (bitcode-framework/ui)

**New files (44+ total)** across Fase A through A4. See git log for full list.

**Modified files:**
- `package.json` — exports map + ./native fix + tsconfig changes
- `tsconfig.json` — allowImportingTsExtensions + noEmit
- `src/components/overlays/drawer/BcDrawer.lite.tsx` — bg-black/50 to bg-foreground/50 (token regression fix)
- `src/components/views/form/BcForm.lite.tsx` — reactive engine integration (Level 1 opt-in)
- `src/core/reactive/{types,store,index}.ts` — scope resolution + Evaluator interface
- `src/components/data/datatable/BcDataTable.lite.tsx` — major rewrite from 245 to 992 lines

### Parent repo

- `docs/components-progress.md` — tracker updated extensively
- `packages/ui` submodule pointer bumped 11 times

---

## Cara Melanjutkan di Session Baru

```
Lanjutkan pekerjaan dari handoff file: handoffs/2026-06-15-bcdatatable-A1-A4-done-render-and-build-fix.md

PRIORITY 1: Fix 11 BcStack test failures (npm run build blocker)
PRIORITY 2: BcDataTable A3-render (virtualized scroll + tree chevron + group rows)
PRIORITY 3: BcDataTable A4-render (cell input + row edit UI)
```

### Verifikasi cepat di awal session baru:

```powershell
cd packages/ui
npm test                    # harus tetap 447 pass / 11 fail
npm run tokens:check        # harus "tokens: ok"
git log --oneline -3        # HEAD harus da9d297
```
