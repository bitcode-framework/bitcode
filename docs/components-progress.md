# Component Progress Tracker

> **Sumber kebenaran progress per komponen.** Update file ini setiap komponen
> berubah status. Untuk *kenapa* dan *arsitektur*, baca [ADR](./adr/README.md).
> Untuk *bagaimana spec engine*, baca `engine/docs/features/*.md`.

**Last updated:** 2026-06-14 (post-ADR-006)

## Major Architecture Update (ADR-006)

**ADR-006 Universal Component Strategy** mengubah model tier komponen:

- **2 tier baru:** Pure (1 source Mitosis) vs Adapt (web + native impl, kontrak sama).
- **Semua komponen wajib jalan di 4 target:** Stencil (legacy frozen), React, Vue, React Native.
- **Tidak ada lagi "Tier 3 platform-only"** — semua komponen wajib universal di caller side.
- **Phase 0 spike** wajib selesai sebelum P0 ekspansi (lihat section "Phase 0" di bawah).

Lihat [ADR-006](./adr/006-universal-component-strategy.md) untuk detail lengkap.

## Ringkasan

| Bucket | Total Target | Stencil ✅ | Mitosis ✅ | Contract ✅ | A11y ✅ | Done (semua) |
|---|---|---|---|---|---|---|
| **Phase 0 (Spike & Cleanup)** | 6 tasks | – | – | – | – | 0 |
| **P0 (ERP-Critical)** | 22 | 13 | 14 | 0 | 0 | 0 |
| **P1 (High Value)** | ~58 | ~40 | ~25 | 0 | 0 | 0 |
| **P2 (Polish)** | ~29 | 14 | 4 | 0 | 0 | 0 |
| **Total komponen visual** | ~110 | 127* | 59 | 0 | 0 | 0 |

\* 127 Stencil = 128 folder, beberapa overlap target / akan disatukan saat porting. Stencil **frozen** post-ADR-006 (tidak dapat fitur baru).

## Status Legend

- **Tier (Pure/Adapt)** = klasifikasi implementasi per ADR-006 Section 2 (Pure = 1 source `.lite.tsx`, Adapt = `.web.tsx` + `.native.tsx` manual)
- **Stencil** = ada di `packages/components/src/components/` (legacy, frozen post-ADR-006)
- **Mitosis** = ada di `packages/ui/src/components/` sebagai `.lite.tsx` (Pure) atau `.web.tsx`+`.native.tsx` (Adapt)
- **Contract** = `BcXxx.contract.ts` ada + sesuai Section 4 ADR-002 (Field/DataView/Overlay) + ADR-006 Section 5
- **A11y** = aria roles + keyboard nav (web) + screen reader test (web+mobile) + touch target ≥44px
- **Done** = semua kolom ✅ + 4-target build pass + token-based + smoke test pass per ADR-006 Section 5 DoD
- `–` = tidak applicable
- `?` = belum diaudit

## Aturan Update

1. Komponen baru selesai → ubah baris di tabel + update Ringkasan.
2. **Jangan tandai Done sebelum semua sub-cell ✅** sesuai ADR-006 Section 5 DoD.
3. Status Mitosis = "ada filenya" (`.lite.tsx` untuk Pure, `.web.tsx`+`.native.tsx` untuk Adapt).
4. Status Contract/A11y = "sudah diverifikasi sesuai ADR-002 + ADR-006".
5. Kalau komponen di-rename / di-merge → tambah catatan di kolom Notes, jangan hapus baris.

---

## Phase 0 — Spike & Cleanup (PREREQUISITE)

**Wajib selesai sebelum ekspansi ke P0** (per ADR-006 Section 6).

| # | Task | Status | Notes |
|---|---|---|---|
| 1 | Drop Lynx artifacts (`output/lynx/`, `scripts/lynx-*.mjs`, `demos/lynx-app/`) | ✅ | Artifacts already removed; ADR-001 §15 marked superseded by ADR-006 (2026-06-14) |
| 2 | Auto-generate token pipeline (`scripts/generate-tokens.mjs`) | ✅ | `tokens:generate` + `tokens:check` (CI guard); emits `tailwind.preset.cjs`, `tailwind.theme.css`, `nativewind.preset.cjs`, `tokens.ts` from `theme.css`. NativeWind preset now exposes 25 color keys (incl `info`) + dark palette via `__rawPalette.dark`. |
| 3 | Tambah RN target ke Mitosis (`mitosis.config.js` + `package.json` exports) | ✅ | `targets:[stencil,react,vue,reactNative]`, `stylesType:'native-wind'`. Mitosis 0.13.1 NATIVE_WIND_STYLES_PLUGIN had `[object Object]` bug for dynamic class bindings — fixed via `patch-package` (`patches/@builder.io+mitosis+0.13.1.patch`). Verified: 65/65 RN outputs emit proper static `className="..."` or dynamic `className={\`tpl\`}`. Package exports `./native` + `react-native`/`nativewind` peerDeps added. |
| 4 | Build resolver setup (Vite `.web.tsx` + Metro `.native.tsx`) | ✅ | Vite `extensions:['.web.tsx',...]` configured. `demos/rn-app/scripts/resolver-check.mjs` replays both bundlers' algorithms — 3/3 PASS. |
| 5 | PoC Pure component (BcStack atau BcBadge) jalan di 4 target | ✅ | BcStack authored. Mitosis emits 66/66 components × 4 targets. Contract tests 16/16 PASS (react/vue/reactNative/stencil). RN output uses `<View className={...}>`, no `[object Object]` regression. |
| 6 | PoC Adapt component (BcPopover atau BcDrawer) di-resolve correctly | ✅ | `src/components/adapt-poc/popover/` with shared `BcPopover.contract.ts`, web (`<span>` + absolute positioning), native (`Modal` bottom sheet). Resolver smoke 3/3 PASS. Contract tests 8/8 PASS. |

**Acceptance gate Phase 0:** Semua 6 task done + ADR-006 di-update dengan temuan konkret (replace "estimasi" jadi "verified") + ADR-006 di-promote dari "Proposed pending spike" jadi "Accepted".

---

## Next 5 to Build (rekomendasi urutan eksekusi POST-Phase-0)

Setelah Phase 0 done, urutan build P0:

1. **BcSlot + BcPortal + BcVisuallyHidden** (P0 fondasi a11y, Pure) — dipakai semua overlay.
2. **BcStack family (Stack/XStack/YStack/Box)** (P0 fondasi layout, Pure) — dipakai banyak komponen lain.
3. **BcBackdrop + BcLoadingOverlay + BcToastQueue** (P0 fondasi feedback, mostly Pure).
4. **BcForm + BcList** (P0 view inti, Pure with adapt sub-component) — jantung low-code ERP.
5. **BcDataTable upgrade** (P0, ADR-003 Section 4, **Adapt** — web grid vs MobileList) — gold standard grid.

Setelah 5 ini selesai, baru lanjut P0 standalone (Steps, Accordion, Drawer, ContextMenu, dll — campuran Pure/Adapt).

---

## P0 — ERP-Critical (Wajib Enterprise Grade)

Kolom **Tier** = klasifikasi ADR-006 Section 2 (P=Pure, A=Adapt). Estimasi awal — boleh berubah saat implementasi kalau ternyata Pure tidak feasible.

| # | Component | Group | Tier | Stencil | Mitosis | Contract | A11y | Done | Notes |
|---|---|---|---|---|---|---|---|---|---|
| 1 | BcForm | views | P+A* | ✅ `bc-view-form` | ✅ `views/form` | ✅ | ⚠️ | ⚠️ | Reactive engine ADR-004 Level 1 opt-in via `reactions` prop. Validation: required + min/max + email + pattern (silent invalid-pattern fallback). Render gates on meta.visible/disabled/required. 14 source-level tests pass. |
| 2 | BcList | views | P | ✅ `bc-view-list` | ✅ `views/list` | ✅ | ✅ | ⚠️ | 4-layer DataSource integration with idle/loading/error/empty/ready states. Custom renderItem slot. Pagination with disabled boundary + a11y labels. 14 source-level tests pass. |
| 3 | BcDataTable | views | **A** | ✅ `bc-datatable` | ✅ `data/datatable` | ✅ | ✅ | ✅ | Enterprise A1+A2+A3+A4 **fully render-integrated**: multi-sort+filter+search+a11y; selection+column ops; tree-flatten+group+aggregate+virtualize (chevron+indent, group header/body/footer, virtualized scroll container with absolute-index reindexing); inline edit (cell-mode dblclick → input with Enter/Escape, row-mode Edit/Save/Cancel action column with `totalColumnCount` accounting); CSV export; state persistence; visual context menu overlay with default Copy row action. 10 pure engines: `engine/{sort,filter,columns,selection,tree,group,virtualize,edit,export,persistence}.ts`. All 16 ADR-003 §4 must-haves have engine + state + render integration. 193 tests pass (109 engine + 84 source-level + contract). |
| 4 | BcDescriptions | data-display | P | ❌ | ✅ `data-display/descriptions` | ✅ | ⚠️ | ⚠️ | grid layout, span per item, em-dash placeholder; source-level tests pass |
| 5 | BcSteps / BcStepper | navigation | P | ❌ | ✅ `navigation/steps` | ✅ | ✅ | ⚠️ | aria-current=step on active marker, clickable gate, token-driven colors |
| 6 | BcAccordion / BcCollapse | layout | P | ❌ | ✅ `layout/accordion` | ✅ | ✅ | ⚠️ | aria-expanded reflects open, single/multi-expand, chevron rotates with transition |
| 7 | BcContextMenu | overlays | **A** | ❌ | ✅ `overlays/contextmenu` | ✅ | ✅ | ✅ | Adapt split landed: web mitosis (right-click anchored menu) + `BcContextMenu.native.tsx` (300ms long-press → bottom action sheet via `Modal+Pressable`). role=menu/menuitem, backdrop dismiss, separators+danger token. |
| 8 | BcSplitter | layout | **A** | ❌ | ✅ `layout/splitter` | ✅ | ✅ | ⚠️ | Currently Pure shape; Adapt split (mobile stacked-collapse) deferred. PointerEvents drag + ArrowKey nudge, aria-label on handle, token-driven palette. 11 tests pass. |
| 9 | BcDrawer | overlays | **A** | ❌ | ✅ `overlays/drawer` | ✅ | ✅ | ✅ | Adapt split landed: web mitosis (4 placements side panel) + `BcDrawer.native.tsx` (idiomatic bottom sheet regardless of placement, optional close button, backdrop honors `closeOnOverlay`). |
| 10 | BcTransfer | forms | P | ❌ | ✅ `forms/transfer` | ✅ | ✅ | ⚠️ | Dual-list, aria-labeled move buttons, disabled gates triple-checked, token surface |
| 11 | BcCascader | forms | **A** | ❌ | ✅ `forms/cascader` | ✅ | ✅ | ✅ | Adapt split landed: web mitosis (anchored side-by-side columns popover) + `BcCascader.native.tsx` (full-screen sheet stepping one column at a time with `‹ Back` button). onChange fires only on leaf, popover tokens, FieldChangeEvent. |
| 12 | BcTreeSelect | forms | **A** | ❌ | ✅ `forms/treeselect` | ✅ | ✅ | ✅ | Adapt split landed: web mitosis (anchored dropdown with flattened tree) + `BcTreeSelect.native.tsx` (bottom sheet wrapping a flattened ScrollView with same expand/collapse + indent semantics). FieldChangeEvent fires for any node (leaf or branch). |
| 13 | BcStatistic | data-display | P | ❌ | ✅ `data-display/statistic` | ✅ | ⚠️ | ⚠️ | trend coloring via semantic tokens, tabular-nums, decorative arrows aria-hidden |
| 14 | BcEmpty | feedback | P | ❌ | ✅ `feedback/empty` | ✅ | ✅ | ⚠️ | default SVG illustration + icon override, decorative visuals aria-hidden, action <Slot /> |
| 15 | BcResult | feedback | P | ❌ | ✅ `feedback/result` | ✅ | ✅ | ⚠️ | 7 status variants (success/error/warning/info/404/403/500), semantic token colors, action <Slot /> |
| 16 | BcOtp / BcPin | forms | P | ❌ | ✅ `forms/otp` | ✅ | ✅ | ⚠️ | semantic fieldset/legend, per-box aria-label, inputmode=numeric, numeric guard, FieldChangeEvent + onComplete |
| 17 | **BcSlot** | primitives | P | ❌ | ✅ `primitives/slot` | ✅ | – | ⚠️ | Composition primitive; `display:contents` wrapper, source-level tests pass |
| 18 | **BcPortal** | primitives | **A** | ❌ | ✅ `primitives/portal` | ✅ | – | ⚠️ | Adapt: web `createPortal` (SSR-safe + Esc → onRequestClose), native `Modal`; 11/11 contract tests pass |
| 19 | **BcVisuallyHidden** | primitives | P | ❌ | ✅ `primitives/visually-hidden` | ✅ | ✅ | ⚠️ | sr-only utility, focus-restore variant, no `aria-hidden` (kept in a11y tree) |
| 20 | **BcStack family** (Stack/XStack/YStack/Box) | primitives | P | ❌ | ✅ all four | ✅ all four | – | ⚠️ | BcStack (Phase 0) + new BcXStack/BcYStack/BcBox primitives. Each enforces fixed direction (XStack=row, YStack=col) or layout-participating box. 27 tests pass across the family. |
| 21 | **BcBackdrop / BcLoadingOverlay** | feedback | P | ❌ | ✅ `feedback/{backdrop,loading-overlay}` | ✅ | ✅ | ⚠️ | Backdrop: dim layer + dismissable + reduced-motion; LoadingOverlay: aria-busy + role=status + sr-only label |
| 22 | **BcToastQueue** (Sonner-style) | overlays | P | ❌ | ✅ `overlays/toast-queue` | ✅ | ✅ | ⚠️ | Stacked toast host + `BcToastQueueController` headless helper (push/dismiss/clear, severity sugar, dedupe by id, maxRetained drop). role=region + aria-live=polite, pointer-events-none container, motion-reduce respect. 21 tests pass (12 source + 9 controller). |

**P0 Tier breakdown (estimasi):** 14 Pure + 7 Adapt + 1 Hybrid (BcForm).

**P0 status:** 13/22 ada di Stencil (legacy frozen), 19/22 ada di Mitosis (skeleton + 5 fondasi baru), 5/22 contract+a11y verified, 0/22 verified Done (DoD §5 butuh Mitosis build smoke + e2e). **Fondasi 5/6 (#17–#19, #21) selesai source+contract+test.** Sisa fondasi: BcStack family ekstension (#20) + BcToastQueue (#22).

> **Update 2026-06-14 — Fase C batch 1.** Five missing P0 primitives landed source-side: BcSlot, BcPortal (Adapt), BcVisuallyHidden, BcBackdrop, BcLoadingOverlay. 38 new source-level contract tests pass (no Mitosis-build dependency). Each follows the BcStack/BcPopover convention: `BcXxx.contract.ts` + `BcXxx.lite.tsx` (Pure) or `BcXxx.{web,native}.tsx` + `index.ts` (Adapt) + `BcXxx.test.ts`. Tokens drift check passes.

> **Update 2026-06-14 — Fase C batch 2.** Six existing P0 skeletons gained `BcXxx.contract.ts` + source-level `BcXxx.test.ts` without source rewrites: BcDescriptions (#4), BcSteps (#5), BcAccordion (#6), BcStatistic (#13), BcEmpty (#14), BcResult (#15). +46 source-level tests pass — verifying token-only colors, aria attributes, semantic state markers, and contract surface. Total cumulative: **136 pass / 11 fail** (the 11 failing tests remain the pre-existing BcStack output-level checks awaiting `npm run build`).

> **Update 2026-06-14 — Fase C batch 3.** Six more P0 skeletons gained contract+tests: BcDrawer (#9, source fix `bg-black/50` → `bg-foreground/50` for token compliance), BcTransfer (#10), BcCascader (#11), BcTreSelect (#12), BcOtp (#16), BcContextMenu (#7). +57 source-level tests pass. **Adapt-marked components** (BcContextMenu, BcDrawer, BcCascader, BcTreSelect) currently live as Pure `.lite.tsx`; the Adapt split (web/native) is deferred until `BcPortal` is wired and per-component native UX is designed. Total cumulative: **193 pass / 11 fail**.

> **Update 2026-06-14 — Fase C batch khusus (BcForm + BcList).** The two view "jantung low-code" components gained contracts+tests. **BcForm (#1)**: BcFormContract pulls Reaction from core/reactive — Level 0 controlled validation + Level 1 reactive engine wiring both verified at source level (14 tests). Render gates on meta.visible/disabled/required composability. **BcList (#2)**: BcListContract pulls DataSource from core/types — 4-layer source/data fallback, full lifecycle status union (idle/loading/error/empty/ready), pagination clamping, custom renderItem slot (14 tests). Total cumulative: **220 pass / 11 fail**. **15/22 P0 components now have contract + source-level tests; 6 partial Adapt-target conversions deferred + BcDataTable enterprise upgrade + BcSplitter & BcStack-family + BcToastQueue remain.**

> **Update 2026-06-14 — Fase C subset 1 (BcSplitter + Stack family + BcToastQueue).** Closes the remaining P0 skeleton + family backlog. **BcSplitter (#8)**: contract + 11 tests. **BcStack family (#20)**: BcXStack/BcYStack/BcBox added as new primitives (3 source + 3 contract + 3 test). XStack locks `flex-row`, YStack locks `flex-col` (Tamagui XStack/YStack convention), BcBox layout-participating wrapper distinct from BcSlot's `display:contents`. **BcToastQueue (#22)**: brand-new Sonner-style stacked toast host with companion `BcToastQueueController` headless helper. The controller owns push/dismiss/clear + severity sugar + id-based dedupe + maxRetained drop. Component renders the controlled list at one of 6 viewport corners, schedules per-entry timers in `onMount`/`onUpdate`/`onUnMount`, dispatches `onDismiss(id)` so the controller drops the entry. **+52 tests** (11 splitter + 20 stack family + 21 toast queue). Total cumulative: **272 pass / 11 fail**. **18/22 P0 components now have contract + source-level tests.** Only BcDataTable enterprise upgrade (#3, ADR-003 §4) and the four deferred Adapt splits remain in P0.

> **Update 2026-06-14 — Fase C subset 2 / BcDataTable A1 (multi-sort + filter + search + a11y).** First commit of the BcDataTable enterprise upgrade per ADR-003 §4. **Engine extraction**: pure headless modules `engine/sort.ts` (multi-rule comparator, 3-state cycle helper, locale-aware string + null-aware compare) and `engine/filter.ts` (12 operators: eq/neq/gt/gte/lt/lte/like/startsWith/endsWith/in/notIn/between/isNull/isNotNull, plus `applyGlobalSearch`). Both reusable for BcList / BcTree. **Component upgrade**: source rewrite of `BcDataTable.lite.tsx` adds multi-sort with Shift+click and 1-based priority badges (gated when ≥2 rules), per-column filter row gated on at least one `filterable` column, global search input above the table, a11y `aria-sort` on every sortable header + `tabIndex` + Enter/Space keyboard toggle. Server mode forwards sort/filter/search via `FetchParams`. Side-effect-in-getter bug from earlier source fixed by splitting `localProcessed` (search→filter→sort) from `localView` (paginated slice). **+49 tests** (15 sort engine + 16 filter engine + 18 source-level). Total cumulative: **321 pass / 11 fail**. **A2/A3/A4 (column ops + virtualization + inline edit) follow in next commits.**

> **Update 2026-06-14 — Fase C subset 2 / BcDataTable A2 (column ops + selection).** Second commit of the enterprise upgrade. **Engine extraction**: `engine/columns.ts` (ColumnState + show-hide + reorder + resize + pin with auto re-derived order so left-pinned sort first, right-pinned last) and `engine/selection.ts` (single/multi/checkbox modes; ctrl-toggle, shift-range from anchor index; header `toggleAllVisible` + `isAllVisibleSelected` + `isSomeVisibleSelected` for tri-state checkbox). **Component upgrade**: header gains an optional leading checkbox column with tri-state `aria-checked` (true/false/mixed); rows expose `aria-selected` and a visual `bg-accent/50` cue when selected; clicks dispatch through `handleRowSelectClick` honouring ctrl/meta + shift modifiers. Header drag-drop reorder gated on `reorderable` (HTML5 drag with `text/bc-col` mime), PointerEvent column resize handle (min 40px floor) gated on `resizable`, column visibility menu gated on `showColumnToggle` with `aria-haspopup`/`aria-expanded`. `defaultGetRowId` falls back from `row.id` to the row index when no accessor is supplied. **+43 tests** (13 columns engine + 16 selection engine + 14 source-level). Total cumulative: **364 pass / 11 fail**. **A3 (virtualization + tree-table + grouping) and A4 (inline edit + export + state persistence + context menu) follow.**

> **Update 2026-06-14 — Fase C subset 2 / BcDataTable A3 (engines: virtualization + tree + grouping).** Third commit of the enterprise upgrade. **Three new pure engines**: `engine/tree.ts` (depth-first flatten with cycle break at 64 levels, `expandAll` skips leaves, `toggleExpand` with stable Set), `engine/group.ts` (stable bucket order from first-occurrence with rows preserving source order; aggregates: sum/avg/count/min/max with numeric/Date/string-coercion, NaN min/max for empty numeric set), `engine/virtualize.ts` (`computeWindow` math: top/middle/bottom with overscan, scrollTop clamping for overshoot+negative, paddingTop+visible*rowHeight+paddingBottom invariant). **Component wiring**: `treeRows` getter calls flattenTree with current `expandedTreeKeys`, `toggleTreeRow` emits `onExpandedChange`; `groupedRows` getter calls groupRows with `groupBy + aggregates`; `virtualWindow` + `virtualRows` + `handleVirtualScroll` recompute on scroll; defaults applied (rowHeight=36, viewportHeight=400, overscan=5). Contract re-exports new engine types. Render integration (virtualized scroll container, tree chevron + indent, group header/footer rows) deferred to a dedicated A3-render sub-commit so this batch stays focused on engine correctness + component state wiring. **+34 tests** (10 tree + 8 group + 8 virtualize + 11 source-level). Total cumulative: **398 pass / 11 fail**. **A4 follows.**

> **Update 2026-06-14 — Fase C subset 2 / BcDataTable A4 (inline edit + export + persistence + context menu).** Fourth and final engine commit of the enterprise upgrade. **Three new pure engines**: `engine/edit.ts` (cell-mode + row-mode with per-row staged-changes Map; commitCell/commitRow return diffs the caller applies; cancelCell/cancelRow drop staging cleanly; isCellEditing/isRowEditing pinpoint matchers), `engine/export.ts` (RFC 4180 CSV serializer with UTF-8 BOM for Excel double-click, custom delimiter/newline/header opts, per-column format() callback, Date→ISO + object→JSON coercion, no DOM/filesystem touches), `engine/persistence.ts` (versioned PersistedState snapshot, serialize/deserialize that returns null on bad input rather than throwing, injectable PersistenceStorage interface so localStorage/sessionStorage/memory/server-side all work, createMemoryStorage helper for tests + SSR fallback). **Component wiring**: edit handlers gate on editMode prop, commitCellEdit fires onCellChange only when there's a diff, commitRowEdit fires onRowChange only when buffer non-empty; triggerCsvDownload guards on `typeof document` + `URL.createObjectURL` (SSR/RN-safe no-op) and only includes visible columns via columnStates filter; restorePersistedState runs in onMount **before** loadData so the initial fetch picks up restored sort/filter/search; persistNow + restorePersistedState both early-return when storage/key are missing; row right-click wires through onContextMenu when showRowContextMenu is true. Export CSV button gated on `exportable` with aria-label. Contract re-exports EditMode, CellRef, ExportColumn, ExportOptions, PersistenceStorage, PersistedState. **+49 tests** (13 edit + 11 export + 11 persistence + 14 source-level). Total cumulative: **447 pass / 11 fail**. **All 16 ADR-003 §4 must-haves now have engine + state wiring; A3/A4 render integration (virtualized scroll, tree chevron, group rows, cell input) follow as separate render sub-commits.**

> **Update 2026-06-15 — BcDataTable A3-render (virtualized scroll + tree chevron + group rows).** Render integration of the A3 engines into `BcDataTable.lite.tsx`. **`bodyMode` discriminator** (`'group' | 'tree' | 'virtual' | 'plain'`) routes the tbody to one of four render branches per ADR-003 §4 row 7-9; precedence is group > tree > virtual > plain so a grouped tree-table gracefully degrades to grouped rows (a future render sub-commit may layer them). **Tree branch** renders `state.treeRows` with `aria-level={depth+1}` + `aria-expanded` on `<tr>`, a chevron `<button>` in the first column gated on `hasChildren` (split into two `<Show>` branches because Mitosis 0.13.1's Stencil generator chokes on ternary-with-unicode-string literals), 16px-per-depth indent on column 0, and `stopPropagation` on the chevron click so toggling expand never fires row select. **Group branch** renders three nested `<For>` blocks: header rows (`data-bc-group-header`, full-width `colSpan` showing field+value+count), body rows per group (`data-bc-group-row`), and an optional footer row gated on `showGroupFooter` whose cells call the new `aggregateValueFor(group, key)` helper. **Aggregate formatter** centralised on state (`formatAggregate(value, kind)`): NaN sentinel renders as em-dash, sum/avg use `toFixed(2)` for non-integers, count is plain integer, min/max plain. **Virtual branch** wraps the existing `<table>` in a scroll container capped to `viewportHeight` with `onScroll → handleVirtualScroll`, emits `<tr aria-hidden>` spacers above and below the rendered window with explicit `height` so the native scrollbar reflects the full dataset, and reindexes every row by the absolute `startIndex + rowIdx` so checkbox / aria-selected stay correct after scroll. **`totalColumnCount` getter** centralises base-cols-plus-selection-slot math used by every full-width row (loading / error / empty / group header / virtual spacer); legacy inline calc replaced. **+9 source-level tests** asserting bodyMode discriminator, tree chevron + indent, group header/body/footer wiring, virtualized container + spacers + absolute row index. Total cumulative: **467 pass / 0 fail.** Mitosis full-glob build green: 76/76 components compile across react/vue/reactNative/stencil. **A4-render (cell input + row edit UI + ctxmenu visual) follows as the next sub-commit.**

> **Update 2026-06-15 — BcDataTable A4-render (inline edit + row actions + visual context menu).** Closes the BcDataTable enterprise upgrade. **`cellAsInput(row, idx, field)` discriminator** unifies cell-edit (`editMode === 'cell'` + active cell match) and row-edit (`editMode === 'row'` + active row match) into one boolean so every body branch (plain / tree / virtual / group body) toggles between `<span>` and `<input>` identically. **`editingCellValue`** prefers the staged-changes buffer over the row source so typing survives across renders. **Cell-mode** triggers via `onDblClick` on the `<td>` (gated on `editMode === 'cell'`); the input wires `onInput → stageEdit`, `onBlur → commitCellEdit`, and `onKeyDown → handleCellEditKey` (Enter commits, Escape cancels). **Row-mode** adds an action column at the rightmost edge of every body branch: an "Edit" button when not editing, a "Save / Cancel" pair when editing. The header gains an `<th>Actions</th>` slot and the per-column-filter row gains a matching empty `<th>` so widths line up; `totalColumnCount` now adds the action column too so loading / error / empty / group-header / group-footer / virtualized spacers all span correctly. Group-footer rows pad the action slot with an empty `<td>` so aggregate values keep aligning with their data column. **Visual context menu** upgrades the existing `handleRowContextMenu` from callback-only to also opening a positioned overlay: `{contextMenuOpen, contextMenuX/Y, contextMenuRow, contextMenuRowIndex}` state drives a fixed-position `<ul role="menu">` with a "Copy row" action that writes tab-separated visible-column values to `navigator.clipboard` (SSR / RN safe via typeof guard). A backdrop catches clicks outside the menu via `closeContextMenu()`. **+9 source-level tests** asserting `cellAsInput` / `cellEditActive` / `rowEditActive` discriminators, dblclick wiring, Edit / Save / Cancel button labels, totalColumnCount action-column accounting, virtualized + tree branches indexing, and context-menu open / close / copy-row clipboard guard. Total cumulative: **476 pass / 0 fail.** Mitosis full-glob build green: 76/76. **All 16 ADR-003 §4 must-haves now have engines + state + render integration; BcDataTable upgrade complete.**

> **Update 2026-06-15 — Adapt splits for BcContextMenu / BcDrawer / BcCascader / BcTreeSelect (ADR-006 §2.2).** Closes the four "Adapt" components left in P0. Each component now ships a hand-written `BcXxx.native.tsx` next to its `.lite.tsx` web source; both impls conform to the shared `BcXxx.contract.ts`. The web target keeps the mitosis-built `output/react/.../BcXxx.tsx` (right-click menu, side panel, anchored columns popover, anchored tree dropdown). The RN target gets the hand-tuned mobile UX via a new post-build script `scripts/apply-adapt-overrides.mjs` that copies each `.native.tsx` over the mitosis-generated `output/react-native/.../BcXxx.tsx`. Vue + Stencil targets stay on the web mitosis output (Adapt is a web-vs-mobile split, not a third dimension). **BcContextMenu.native.tsx**: Modal+Pressable bottom sheet that opens on a 300ms long-press (mobile equivalent of right-click); items support icon/separator/disabled/danger same as the web menu. **BcDrawer.native.tsx**: Modal-anchored bottom sheet regardless of `placement` prop (idiomatic on iOS + Android), close button optional, backdrop dismiss honors `closeOnOverlay`. **BcCascader.native.tsx**: full-screen sheet that walks one column at a time with a `‹ Back` step-up button — matches the dominant mobile picker pattern; `onChange` still emits a `FieldChangeEvent` only on leaf selection. **BcTreeSelect.native.tsx**: bottom sheet wrapping a flattened `ScrollView` with the same expand/collapse + indent semantics as the web dropdown. **+22 source-level tests** (5+5+6+6) asserting per-component file presence, contract import, RN-only primitives, mobile-pattern marker (long-press / bottom-sheet / Back button / flatten helper), and `data-bc-platform="native"` parity marker. **`scripts/apply-adapt-overrides.mjs`** is idempotent, fails loud when a `.native.tsx` is missing, prepends an "AUTO-COPIED" banner to the generated mirror; wired into `npm run build` between `fix-react-output` and `generate-index`. Total cumulative: **498 pass / 0 fail.** Mitosis full-glob build still green: 76/76. **All 22 P0 components (5 fondasi + 6 batch 2 + 6 batch 3 + BcForm + BcList + BcSplitter + BcStack family + BcToastQueue + BcDataTable) now have full contract + source-level tests + (where applicable) Adapt split. P0 surface complete.**

> **Update 2026-06-15 — Fase D / BcRenderer (ADR-005 Phases 1-3 landed).** Schema-driven UI compositor lives at `core/renderer/`. `UiSchema` document is **domain-agnostic** (no BitCode coupling) so the renderer is genuinely standalone — JsonForms-style usage works out of the box. Six schema kinds: **`'form'`** wires `<BcForm>` with fields + reactions; **`'detail'`** is the same form rendered `disabled + hideSubmit`; **`'list'`** wires `<BcDataTable>` with `columns_def` / `data` / `source` / `pageSize`; **`'section'`** + **`'grid'`** are grouping wrappers that recurse via `<BcRenderer>` (CSS `grid-template-columns: repeat(N, minmax(0, 1fr))`); **`'custom'`** emits a sentinel `<div data-bc-renderer-custom data-payload={...}>` for host-portal replacement. **`FieldType`** discriminator (input / textarea / select / number / date / toggle / checkbox / email / password / url / tel) maps to `<BcForm>`'s native input string via `resolveFieldType`. Resolver helpers (`resolveField`, `resolveFields`, `defaultFieldResolver`) are pure TS so consumers can pre-process schemas before handing them to the renderer (this is how custom field types like `lookup` / `signature` plug in — Mitosis cannot do runtime dynamic dispatch, so pre-processing is the deliberate extension story). **Reactivity** rides through unchanged: `schema.reactions` forwarded into `<BcForm>`'s existing reactive engine (Level 1 per ADR-004), so all four effect kinds (DATA / VALIDATION / UI / VALUE) work without renderer changes. **`scripts/generate-index.mjs`** updated to also scan `src/core/renderer/` and to mirror non-`components/` source paths into the per-target output barrels. **+22 tests** (file presence + `resolveFieldType` mappings + `resolveField` field-by-field semantics + `resolveFields` purity + barrel re-exports + per-kind source-level branching + section/grid recursion + grid-template-columns wiring + custom-sentinel data attributes). Total cumulative: **520 pass / 0 fail.** Mitosis full-glob build green: **77/77** components. **`@bitcode/ui` is now a complete schema-driven UI library: standalone usage = JsonForms parity; BitCode integration = engine maps `models/*.json` -> `UiSchema` and feeds it to `<BcRenderer>`.** Phase 4 (BitCode JSON transformer) lives in the BitCode engine, outside `@bitcode/ui` scope.

> **Update 2026-06-15 — Batch 1A: P1 Forms skeleton-existing → contract+tests landed.** Six P1 Forms components that already had `.lite.tsx` skeletons gained matching `BcXxx.contract.ts` + source-level tests at the same DoD bar as P0: **BcSegmented** (10 tests; role=tablist + tab + aria-selected, sm/md/lg sizes, FieldChangeEvent, controlled/uncontrolled), **BcMentions** (12 tests; @-prefix token detector with whitespace-close, case-insensitive `label||value` filter, pick replaces token with `prefix+value+space`, emits both onChange + onSelect), **BcSlider** (13 tests; native `<input type="range">` for free keyboard a11y, gradient-driven track via `--bc-primary`/`--bc-muted` tokens, Number coercion guard, aria-label fallback, tabular-nums value display), **BcCombobox** (12 tests; role=combobox + aria-expanded, case-insensitive trim filter with empty-query short-circuit, focus reopens popover, emptyText empty state, label inserted but stable value emitted on pick), **BcColorPicker** (13 tests; native `<input type="color">` for system color pickers, preset swatches with `ring-2 ring-ring` active halo, button type="button" guard so swatches do not submit forms, uppercase tabular-nums hex display), **BcRate** (13 tests; role=radiogroup + per-star aria-pressed + numeric aria-label, antd-style click-to-clear, sm/md/lg icon sizing, `text-warning` token for filled stars (no hardcoded yellow), decorative SVG aria-hidden, hover affordance gated on `!disabled && !readonly`). All six wire `FieldChangeEvent` consistently with name/value/oldValue/timestamp shape. **+73 source-level tests** total. Cumulative: **593 pass / 0 fail.** Mitosis full-glob build still green: 77/77. **6/22 P1 Forms now contract+test verified** (existing skeletons promoted; net-new + Stencil migrate cohorts follow).

---

## P1 — High Value

### Forms (P1)
| Component | Group | Stencil | Mitosis | Contract | A11y | Notes |
|---|---|---|---|---|---|---|
| BcSegmented | forms | ❌ | ✅ | ✅ | ✅ | Net-new. role=tablist + role=tab + aria-selected; controlled/uncontrolled; sm/md/lg sizes; FieldChangeEvent. 10 source-level tests pass. |
| BcMentions | forms | ❌ | ✅ | ✅ | ✅ | Untuk Chatter. @-prefix token detector + filtered popover; whitespace closes; pick replaces token with `prefix+value+space`; emits onChange + onSelect. 12 source-level tests pass. |
| BcSlider | forms | ❌ | ✅ `forms/slider` | ✅ | ✅ | Net-new. Native `<input type="range">` for free a11y; %-driven gradient track via `--bc-primary`/`--bc-muted`; aria-label fallback; tabular-nums value display. 13 source-level tests pass. |
| BcRangeSlider | forms | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcCombobox / BcAutocomplete | forms | ❌ | ✅ `forms/combobox` | ✅ | ✅ | Net-new. role=combobox + aria-expanded; case-insensitive trim filter; focus reopens; emptyText state; FieldChangeEvent on pick (label inserted, value emitted). 12 source-level tests pass. |
| BcMultiSelect | forms | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcColorPicker | forms | ✅ `bc-field-color` | ✅ `forms/colorpicker` | ✅ | ✅ | Migrate + reuse. Native `<input type="color">` for system pickers; preset swatches with active ring; tabular-nums uppercase hex; FieldChangeEvent. 13 source-level tests pass. |
| BcRate | forms | ✅ `bc-field-rating` | ✅ `forms/rate` | ✅ | ✅ | role=radiogroup + aria-pressed per star; click-to-clear (antd idiom); sm/md/lg sizes; warning token (no hardcoded yellow); decorative SVG aria-hidden. 13 source-level tests pass. |
| BcTagsInput | forms | ✅ `bc-field-tags` | ❌ | ❌ | ❌ | Migrate |
| BcRichText | forms | ✅ `bc-field-richtext` | ❌ | ❌ | ❌ | Migrate |
| BcCode | forms | ✅ `bc-field-code` | ❌ | ❌ | ❌ | Migrate |
| BcFile | forms | ✅ `bc-field-file` | ❌ | ❌ | ❌ | Migrate |
| BcJson | forms | ✅ `bc-field-json` | ❌ | ❌ | ❌ | Migrate |
| BcSignature | forms | ✅ `bc-field-signature` | ❌ | ❌ | ❌ | Migrate |
| BcBarcode | forms | ✅ `bc-field-barcode` | ❌ | ❌ | ❌ | Migrate |
| BcGeo | forms | ✅ `bc-field-geo` | ❌ | ❌ | ❌ | Migrate |
| **BcDateField** (input-only) | forms | ❌ | ❌ | ❌ | ❌ | **NEW Tier 2** |
| **BcTimeField** (input-only) | forms | ❌ | ❌ | ❌ | ❌ | **NEW Tier 2** |
| **BcNumberField** (locale-aware) | forms | ❌ | ❌ | ❌ | ❌ | **NEW Tier 2** |
| **BcDropZone** | forms | ❌ | ❌ | ❌ | ❌ | **NEW Tier 2** |
| **BcToggleGroup** | forms | ❌ | ❌ | ❌ | ❌ | **NEW Tier 2** |

### Navigation (P1)
| Component | Group | Stencil | Mitosis | Contract | A11y | Notes |
|---|---|---|---|---|---|---|
| BcCommand (palette) | overlays | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcMenubar | navigation | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcMegaMenu | navigation | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcNavigationMenu | navigation | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcAffix | navigation | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcAnchor | navigation | ❌ | ✅ | ❌ | ❌ | Net-new |
| BcBackTop | navigation | ❌ | ✅ `navigation/backtop` | ❌ | ❌ | Net-new |
| BcBreadcrumb | navigation | ❌ | ✅ | ❌ | ❌ | — |
| BcPagination | navigation | ❌ | ✅ | ❌ | ❌ | — |
| BcDropdown | navigation | ❌ | ✅ | ❌ | ❌ | — |

### Feedback / Overlays (P1)
| Component | Group | Stencil | Mitosis | Contract | A11y | Notes |
|---|---|---|---|---|---|---|
| BcNotification | feedback | ❌ | ❌ | ❌ | ❌ | Beda dari Toast |
| BcTour | feedback | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcAlertDialog | overlays | ✅ `bc-dialog-alert` | ❌ | ❌ | ❌ | Migrate |
| BcPrompt | overlays | ✅ `bc-dialog-prompt` | ❌ | ❌ | ❌ | Migrate |
| BcWizard | overlays | ✅ `bc-dialog-wizard` | ❌ | ❌ | ❌ | Migrate |
| BcHoverCard | overlays | ❌ | ❌ | ❌ | ❌ | Net-new |
| **BcActionsheet** | overlays | ❌ | ❌ | ❌ | ❌ | **NEW Tier 3** mobile bottom-sheet |
| **BcAdapt** | overlays | ❌ | ❌ | – | – | **NEW Tier 3** responsive contract |

### Data Display / Layout (P1)
| Component | Group | Stencil | Mitosis | Contract | A11y | Notes |
|---|---|---|---|---|---|---|
| BcCarousel | data-display | ❌ | ✅ | ❌ | ❌ | Net-new |
| BcOrgChart | data-display | ❌ | ❌ | ❌ | ❌ | HRM critical |
| BcInplace | data-display | ❌ | ❌ | ❌ | ❌ | Inline edit |
| **BcAspectRatio** | primitives | ❌ | ❌ | – | ✅* | **NEW Tier 2** |
| **BcDisclosure / BcDisclosureGroup** | data-display | ❌ | ❌ | ❌ | ❌ | **NEW Tier 2** spec-correct |
| **BcTagGroup** (read-only) | data-display | ❌ | ❌ | ❌ | ❌ | **NEW Tier 2** |
| **BcFlex** | layout | ❌ | ❌ | – | – | **NEW Tier 2** |
| **BcSparkline** | charts | ❌ | ❌ | ❌ | ❌ | **NEW Tier 1 (MUI X)** |

### Media & Social (P1) — semua di Stencil, perlu migrate
| Component | Stencil | Mitosis | Notes |
|---|---|---|
| BcViewerPdf | ✅ | ❌ | Migrate |
| BcViewerImage | ✅ | ❌ | Migrate |
| BcViewerDocument | ✅ | ❌ | Migrate |
| BcViewerVideo | ✅ | ❌ | Migrate |
| BcViewerAudio | ✅ | ❌ | Migrate |
| BcViewerYoutube | ✅ | ❌ | Migrate |
| BcViewerInstagram | ✅ | ❌ | Migrate |
| BcViewerTiktok | ✅ | ❌ | Migrate |
| BcChatter | ✅ | ❌ | Migrate |
| BcActivity | ✅ | ❌ | Migrate |
| BcTimeline (social) | ✅ | ❌ | Migrate |

---

## P2 — Polish

| Component | Group | Stencil | Mitosis | Notes |
|---|---|---|---|
| BcWatermark | data-display | ❌ | ✅ | Net-new |
| BcQRCode | data-display | ❌ | ✅ | Net-new |
| BcSpeedDial | overlays | ❌ | ✅ | Net-new |
| BcKbd | primitives | ❌ | ✅ | Net-new |
| BcVirtualScroller | views | ❌ | ❌ | Net-new |
| BcKnob | forms | ❌ | ❌ | Net-new |
| BcMeterGroup | feedback | ❌ | ❌ | Net-new |
| BcDock | navigation | ❌ | ❌ | Net-new |
| BcGalleria | media | ❌ | ❌ | Net-new |
| BcViewGantt | views | ✅ | ❌ | Migrate |
| BcViewMap | views | ✅ | ❌ | Migrate |
| BcViewReport | views | ✅ | ❌ | Migrate |
| BcViewEditor | views | ✅ | ❌ | Migrate |
| BcViewActivity | views | ✅ | ❌ | Migrate |
| BcExport | print | ✅ | ❌ | Migrate |
| BcPrint | print | ✅ | ❌ | Migrate |
| BcReportLink | print | ✅ | ❌ | Migrate |
| BcSearch | navigation | ✅ | ❌ | Migrate |
| BcFavorites | navigation | ✅ | ❌ | Migrate |
| **BcColorArea** (sub-primitive) | forms | ❌ | ❌ | **NEW Tier 2** |
| **BcColorWheel** | forms | ❌ | ❌ | **NEW Tier 2** |
| **BcColorSlider** | forms | ❌ | ❌ | **NEW Tier 2** |
| **BcColorSwatch** | forms | ❌ | ❌ | **NEW Tier 2** |
| **BcImage** (placeholder/blurhash) | data-display | ❌ | ❌ | **NEW Tier 3** |
| **BcCenter / BcSquare / BcCircle** | primitives | ❌ | ❌ | **NEW Tier 3** |
| **BcPanelMenu** | navigation | ❌ | ❌ | **NEW Tier 1 (Prime)** |
| **BcTieredMenu** | navigation | ❌ | ❌ | **NEW Tier 1 (Prime)** |
| **BcSlideMenu** | navigation | ❌ | ❌ | **NEW Tier 1 (Prime)** |
| **BcSidebar** (nav block) | layout | ❌ | ❌ | **NEW Tier 1 (shadcn)** |

---

## Charts (Stencil — separate migration plan)

26 chart variants di Stencil: area, bar, boxplot, candlestick, custom, funnel,
gauge, graph, heatmap, kpi, line, mixed, parallel, pictorialbar, pie, pivot,
polar, progress, radar, sankey, scatter, scorecard, sunburst, themeriver, tree,
treemap.

**Strategi:** rebuild jadi `BcChart` tunggal dengan prop `type`, sesuai
ADR-003 Section 1 taxonomy. **Belum dimulai di Mitosis.**

| Status | Count |
|---|---|
| Stencil ada | 26 |
| Mitosis ada | 0 |
| Verified Done | 0 |

---

## Cross-tier Infrastructure (ADR-002 Section 9)

Bukan komponen visual. Tracker ini hanya untuk awareness.

| Item | Status | Notes |
|---|---|---|
| `core/providers/AppProvider` | ⚠️ | `AppConfigStore` + `BcAppContext` + `attachSystemThemeWatcher` + Mitosis `.lite.tsx`. 9/9 unit test pass. Belum diadopsi komponen P0. |
| `core/hooks/` subset P0 (8 hooks) | ⚠️ | Implementasi headless: `useId`/`useDisclosure`/`resolveControllable`/`DebouncedValue`/`copyToClipboard`/`createMediaQueryWatcher`/`createFocusTrap`/`createElementSizeObserver`. 13/13 unit test pass. Belum dipakai komponen P0. |
| `core/hooks/` subset P1 (6 hooks) | ❌ | useHover, useFocusVisible, useEventListener, useThrottledCallback, useLocalStorage, useIsomorphicLayoutEffect |
| `core/animation/` contract | ⚠️ | `AnimationProps` + `ANIMATION_DEFAULT` exported. Implementasi per-tier (CSS transition / RN Animated) belum. |
| `core/adapt/` contract | ⚠️ | `AdaptProps` + `isPlatformCondition` + canonical breakpoints. Implementasi runtime (Popover→Sheet) belum. |

> **Update 2026-06-14 — Fase A done.** Cross-tier infrastructure (ADR-002 §9) skeleton landed in `packages/ui/src/core/{providers,hooks,animation,adapt}/` with unit-test coverage. Adoption into 22 P0 components is the next gate.

> **Update 2026-06-14 — Fase B done (partial).** Reactive engine (ADR-004) gained scope resolution (`$self`/`$parent`/`$root`/`$index`), `Evaluator` plug-in interface (Phase E slot), `BcReactiveScope` Level 2 wrapper, and **opt-in integration into `BcForm`** (`reactions` + `onReactiveChange` props). 10 new tests pass on top of the existing 6. Repeater (Phase B per ADR-004) and BitCode JSON feed (Phase D) remain pending.

---

## Foundation Phases (ADR-002 Section 8 + ADR-006 Section 6)

Bukan komponen, tapi prasyarat. Tracker ringkas:

| Phase | Status | Notes |
|---|---|---|
| **Phase 0 (ADR-006)** — Spike & Cleanup | ❌ | **Active — see "Phase 0" section above** |
| Phase 1 — Fondasi (token auto-gen, contract template, exports map) | ❌ | Sebagian setup ada (`tokens/*`), butuh auto-gen pipeline |
| Phase 2 — Perbaiki komponen rusak (BcDataTable, BcInput refactor, BcSelect a11y) | ❌ | Per ADR-002 Section 6 Findings Register |
| Phase 3 — P0 Fondasi a11y/layout (#17-22 di tabel P0) | ❌ | Pasca Phase 0 spike sukses |
| Phase 4 — P0 Views (BcForm, BcList) + DataTable upgrade | ❌ | Jantung ERP |
| Phase 5 — P0 sisanya (Steps, Drawer, ContextMenu, Cascader, dll) | ❌ | — |
| Phase 6 — P1 ekspansi | ❌ | Setelah P0 stabil |
| Phase 7 — P2 polish | ❌ | Opsional, after P1 |
| Phase 8 — Tier B/C (Tamagui/gluestack/Flutter codegen) | ❌ | Per ADR-002, tetap ditunda |

---

## Catatan Penting

1. **Mitosis 59 komponen ≠ 59 komponen siap pakai.** Sebagian besar masih
   skeleton (`.lite.tsx` dengan logic minimal). Status ✅ di kolom "Mitosis"
   hanya berarti file ada — **bukan** berarti contract/a11y/token sudah lengkap.

2. **Stencil 127 komponen masih production-ready** dan dipakai engine BitCode
   sekarang. Migrasi ke Mitosis tidak menghapus Stencil sampai engine konsumsi
   output Mitosis (sesuai ADR-002 dan arahan stakeholder).

3. **Yang membuat komponen Done bukan jumlah file**, tapi: contract penuh +
   a11y verified + token-based + Tier A pass + tests + docs. Prinsip ini dari
   ADR-002 Section 5.3 Definition of Done.

4. **Tracker ini stale kalau tidak diupdate.** Update bersamaan dengan setiap
   commit yang menyentuh `packages/ui/`.
