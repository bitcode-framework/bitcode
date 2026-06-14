# Component Progress Tracker

> **Sumber kebenaran progress per komponen.** Update file ini setiap komponen
> berubah status. Untuk *kenapa* dan *arsitektur*, baca [ADR](./adr/README.md).
> Untuk *bagaimana spec engine*, baca `engine/docs/features/*.md`.

**Last updated:** 2026-06-14

## Ringkasan

| Bucket | Total Target | Stencil ✅ | Mitosis ✅ | Contract ✅ | A11y ✅ | Done (semua) |
|---|---|---|---|---|---|---|
| **P0 (ERP-Critical)** | 20 | 13 | 8 | 0 | 0 | 0 |
| **P1 (High Value)** | ~45 | ~40 | ~25 | 0 | 0 | 0 |
| **P2 (Polish)** | ~25 | 14 | 4 | 0 | 0 | 0 |
| **Tier 2/3 net-new** | ~20 | 0 | 0 | 0 | 0 | 0 |
| **Total** | ~110 | 127* | 59 | 0 | 0 | 0 |

\* 127 Stencil = 128 folder, beberapa overlap target / akan disatukan saat porting.

## Status Legend

- **Stencil** = ada di `packages/components/src/components/` (legacy, working)
- **Mitosis** = ada di `packages/ui/src/components/` sebagai `.lite.tsx`
- **Contract** = sesuai Section 4 ADR-002 (Field/DataView/Overlay contract penuh)
- **A11y** = aria roles + keyboard nav + focus management diverifikasi
- **Done** = semua kolom di atas ✅ + Tier A pass + token-based styling
- `–` = tidak applicable
- `?` = belum diaudit

## Aturan Update

1. Komponen baru selesai → ubah baris di tabel + update Ringkasan.
2. **Jangan tandai Done sebelum semua sub-cell ✅.**
3. Status Stencil/Mitosis = "ada filenya". Status Contract/A11y = "sudah diverifikasi sesuai ADR".
4. Kalau komponen di-rename / di-merge → tambah catatan di kolom Notes, jangan hapus baris.

---

## Next 5 to Build (rekomendasi urutan eksekusi)

Berdasarkan ADR-003 Section 5 (Build Sequence) + dependency:

1. **Fondasi token + contract** (ADR-002 Phase 0–2) — bukan komponen, tapi prasyarat.
2. **BcSlot + BcPortal + BcVisuallyHidden** (P0 fondasi a11y) — dipakai semua overlay.
3. **BcStack family (Stack/XStack/YStack/Box)** (P0 fondasi layout) — dipakai banyak komponen lain.
4. **BcForm + BcList** (P0 view inti) — jantung low-code ERP.
5. **BcDataTable upgrade** (P0, ADR-003 Section 4) — gold standard grid.

Setelah 5 ini selesai, baru lanjut P0 standalone (Steps, Accordion, Drawer, dll).

---

## P0 — ERP-Critical (Wajib Enterprise Grade)

| # | Component | Group | Stencil | Mitosis | Contract | A11y | Done | Notes |
|---|---|---|---|---|---|---|---|---|
| 1 | BcForm | views | ✅ `bc-view-form` | ❌ | ❌ | ❌ | ❌ | Jantung low-code |
| 2 | BcList | views | ✅ `bc-view-list` | ❌ | ❌ | ❌ | ❌ | Migrate dari Stencil |
| 3 | BcDataTable | views | ✅ `bc-datatable` | ✅ `data/datatable` | ❌ | ❌ | ❌ | Upgrade ke spec ADR-003 §4 |
| 4 | BcDescriptions | data-display | ❌ | ✅ | ❌ | ❌ | ❌ | Net-new |
| 5 | BcSteps / BcStepper | navigation | ❌ | ✅ `navigation/steps` | ❌ | ❌ | ❌ | Net-new |
| 6 | BcAccordion / BcCollapse | layout | ❌ | ✅ `layout/accordion` | ❌ | ❌ | ❌ | Net-new |
| 7 | BcContextMenu | overlays | ❌ | ✅ | ❌ | ❌ | ❌ | Net-new |
| 8 | BcSplitter | layout | ❌ | ✅ | ❌ | ❌ | ❌ | Net-new |
| 9 | BcDrawer | overlays | ❌ | ✅ | ❌ | ❌ | ❌ | Net-new |
| 10 | BcTransfer | forms | ❌ | ✅ `forms/transfer` | ❌ | ❌ | ❌ | Net-new |
| 11 | BcCascader | forms | ❌ | ✅ | ❌ | ❌ | ❌ | Net-new |
| 12 | BcTreeSelect | forms | ❌ | ✅ `forms/treeselect` | ❌ | ❌ | ❌ | Net-new |
| 13 | BcStatistic | data-display | ❌ | ✅ | ❌ | ❌ | ❌ | Net-new |
| 14 | BcEmpty | feedback | ❌ | ✅ | ❌ | ❌ | ❌ | Net-new |
| 15 | BcResult | feedback | ❌ | ✅ `feedback/result` | ❌ | ❌ | ❌ | Net-new |
| 16 | BcOtp / BcPin | forms | ❌ | ✅ `forms/otp` | ❌ | ❌ | ❌ | Net-new |
| 17 | **BcSlot** | primitives | ❌ | ❌ | ❌ | – | ❌ | **NEW from Tier 2/3 — fondasi** |
| 18 | **BcPortal** | primitives | ❌ | ❌ | ❌ | – | ❌ | **NEW** |
| 19 | **BcVisuallyHidden** | primitives | ❌ | ❌ | ❌ | ✅* | ❌ | **NEW** *(by definition)* |
| 20 | **BcStack family** (Stack/XStack/YStack/Box) | primitives | ❌ | ❌ | ❌ | – | ❌ | **NEW** |
| 21 | **BcBackdrop / BcLoadingOverlay** | feedback | ❌ | ❌ | ❌ | ❌ | ❌ | **NEW** |
| 22 | **BcToastQueue** (Sonner-style) | overlays | ❌ | ❌ | ❌ | ❌ | ❌ | **NEW** — beda dari `overlays/toast` tunggal |

**P0 status:** 13/22 ada di Stencil, 14/22 ada di Mitosis (sebagai `.lite.tsx` skeleton),
0/22 verified Done. **Fondasi (#17–#22) belum dimulai.**

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
| `core/providers/AppProvider` | ❌ | Theme/locale/RTL/messageQueue runtime |
| `core/hooks/` subset P0 (8 hooks) | ❌ | useDisclosure, useControllableState, useId, useDebouncedValue, useClipboard, useMediaQuery, useFocusTrap, useElementSize |
| `core/hooks/` subset P1 (6 hooks) | ❌ | useHover, useFocusVisible, useEventListener, useThrottledCallback, useLocalStorage, useIsomorphicLayoutEffect |
| `core/animation/` contract | ❌ | AnimationProps, reduced-motion handling |
| `core/adapt/` contract | ❌ | Popover↔Sheet morphing |

---

## Foundation Phases (ADR-002 Section 8)

Bukan komponen, tapi prasyarat. Tracker ringkas:

| Phase | Status | Notes |
|---|---|---|
| Phase 0 — Spike Token + Contract | ❌ | Belum ada bukti spike |
| Phase 1 — Fondasi (token, contract, exports) | ❌ | — |
| Phase 2 — Perbaiki komponen rusak | ❌ | — |
| Phase 3 — Migrasi Views (jantung ERP) | ❌ | — |
| Phase 4 — Komponen enterprise baru (ADR-003) | ❌ | Sebagian Mitosis sudah ada skeleton |
| Phase 5 — Tier B/C (Tamagui/gluestack/Flutter) | ❌ | Tunda |

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
