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
| 1 | BcForm | views | P+A* | ✅ `bc-view-form` | ✅ `views/form` | ❌ | ❌ | ❌ | Reactive engine (ADR-004 Level 1) terintegrasi opt-in via `reactions` prop; `*` sub-component picker mungkin Adapt |
| 2 | BcList | views | P | ✅ `bc-view-list` | ❌ | ❌ | ❌ | ❌ | Migrate dari Stencil |
| 3 | BcDataTable | views | **A** | ✅ `bc-datatable` | ✅ `data/datatable` | ❌ | ❌ | ❌ | Upgrade ke spec ADR-003 §4; mobile = MobileList variant |
| 4 | BcDescriptions | data-display | P | ❌ | ✅ `data-display/descriptions` | ✅ | ⚠️ | ⚠️ | grid layout, span per item, em-dash placeholder; source-level tests pass |
| 5 | BcSteps / BcStepper | navigation | P | ❌ | ✅ `navigation/steps` | ✅ | ✅ | ⚠️ | aria-current=step on active marker, clickable gate, token-driven colors |
| 6 | BcAccordion / BcCollapse | layout | P | ❌ | ✅ `layout/accordion` | ✅ | ✅ | ⚠️ | aria-expanded reflects open, single/multi-expand, chevron rotates with transition |
| 7 | BcContextMenu | overlays | **A** | ❌ | ✅ | ❌ | ❌ | ❌ | Right-click web vs long-press actionsheet mobile |
| 8 | BcSplitter | layout | **A** | ❌ | ✅ | ❌ | ❌ | ❌ | Drag-resize web vs stacked-collapse mobile |
| 9 | BcDrawer | overlays | **A** | ❌ | ✅ | ❌ | ❌ | ❌ | Slide panel web vs sheet mobile |
| 10 | BcTransfer | forms | P | ❌ | ✅ `forms/transfer` | ❌ | ❌ | ❌ | Net-new |
| 11 | BcCascader | forms | **A** | ❌ | ✅ | ❌ | ❌ | ❌ | Dropdown web vs picker mobile |
| 12 | BcTreeSelect | forms | **A** | ❌ | ✅ `forms/treeselect` | ❌ | ❌ | ❌ | Dropdown web vs sheet mobile |
| 13 | BcStatistic | data-display | P | ❌ | ✅ `data-display/statistic` | ✅ | ⚠️ | ⚠️ | trend coloring via semantic tokens, tabular-nums, decorative arrows aria-hidden |
| 14 | BcEmpty | feedback | P | ❌ | ✅ `feedback/empty` | ✅ | ✅ | ⚠️ | default SVG illustration + icon override, decorative visuals aria-hidden, action <Slot /> |
| 15 | BcResult | feedback | P | ❌ | ✅ `feedback/result` | ✅ | ✅ | ⚠️ | 7 status variants (success/error/warning/info/404/403/500), semantic token colors, action <Slot /> |
| 16 | BcOtp / BcPin | forms | P | ❌ | ✅ `forms/otp` | ❌ | ❌ | ❌ | Net-new |
| 17 | **BcSlot** | primitives | P | ❌ | ✅ `primitives/slot` | ✅ | – | ⚠️ | Composition primitive; `display:contents` wrapper, source-level tests pass |
| 18 | **BcPortal** | primitives | **A** | ❌ | ✅ `primitives/portal` | ✅ | – | ⚠️ | Adapt: web `createPortal` (SSR-safe + Esc → onRequestClose), native `Modal`; 11/11 contract tests pass |
| 19 | **BcVisuallyHidden** | primitives | P | ❌ | ✅ `primitives/visually-hidden` | ✅ | ✅ | ⚠️ | sr-only utility, focus-restore variant, no `aria-hidden` (kept in a11y tree) |
| 20 | **BcStack family** (Stack/XStack/YStack/Box) | primitives | P | ❌ | ⚠️ `BcStack` only | ✅ `BcStack` | – | ⚠️ | BcStack PoC done (Phase 0); XStack/YStack/Box belum |
| 21 | **BcBackdrop / BcLoadingOverlay** | feedback | P | ❌ | ✅ `feedback/{backdrop,loading-overlay}` | ✅ | ✅ | ⚠️ | Backdrop: dim layer + dismissable + reduced-motion; LoadingOverlay: aria-busy + role=status + sr-only label |
| 22 | **BcToastQueue** (Sonner-style) | overlays | P | ❌ | ❌ | ❌ | ❌ | ❌ | **Pending** — depends on existing `overlays/toast` consolidation |

**P0 Tier breakdown (estimasi):** 14 Pure + 7 Adapt + 1 Hybrid (BcForm).

**P0 status:** 13/22 ada di Stencil (legacy frozen), 19/22 ada di Mitosis (skeleton + 5 fondasi baru), 5/22 contract+a11y verified, 0/22 verified Done (DoD §5 butuh Mitosis build smoke + e2e). **Fondasi 5/6 (#17–#19, #21) selesai source+contract+test.** Sisa fondasi: BcStack family ekstension (#20) + BcToastQueue (#22).

> **Update 2026-06-14 — Fase C batch 1.** Five missing P0 primitives landed source-side: BcSlot, BcPortal (Adapt), BcVisuallyHidden, BcBackdrop, BcLoadingOverlay. 38 new source-level contract tests pass (no Mitosis-build dependency). Each follows the BcStack/BcPopover convention: `BcXxx.contract.ts` + `BcXxx.lite.tsx` (Pure) or `BcXxx.{web,native}.tsx` + `index.ts` (Adapt) + `BcXxx.test.ts`. Tokens drift check passes.

> **Update 2026-06-14 — Fase C batch 2.** Six existing P0 skeletons gained `BcXxx.contract.ts` + source-level `BcXxx.test.ts` without source rewrites: BcDescriptions (#4), BcSteps (#5), BcAccordion (#6), BcStatistic (#13), BcEmpty (#14), BcResult (#15). +46 source-level tests pass — verifying token-only colors, aria attributes, semantic state markers, and contract surface. Total cumulative: **136 pass / 11 fail** (the 11 failing tests remain the pre-existing BcStack output-level checks awaiting `npm run build`).

---

## P1 — High Value

### Forms (P1)
| Component | Group | Stencil | Mitosis | Contract | A11y | Notes |
|---|---|---|---|---|---|---|
| BcSegmented | forms | ❌ | ✅ | ❌ | ❌ | Net-new |
| BcMentions | forms | ❌ | ✅ | ❌ | ❌ | Untuk Chatter |
| BcSlider | forms | ❌ | ✅ `forms/slider` | ❌ | ❌ | Net-new |
| BcRangeSlider | forms | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcCombobox / BcAutocomplete | forms | ❌ | ✅ `forms/combobox` | ❌ | ❌ | Net-new |
| BcMultiSelect | forms | ❌ | ❌ | ❌ | ❌ | Net-new |
| BcColorPicker | forms | ✅ `bc-field-color` | ✅ `forms/colorpicker` | ❌ | ❌ | Migrate + reuse |
| BcRate | forms | ✅ `bc-field-rating` | ✅ `forms/rate` | ❌ | ❌ | — |
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
