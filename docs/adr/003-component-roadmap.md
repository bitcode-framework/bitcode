# ADR-003: Enterprise Component Roadmap & Taxonomy

## Status: **Proposed**

## Date: 2026-06-06

## Depends on: [ADR-002](./002-cross-platform-architecture.md)

ADR-002 menetapkan fondasi (tiering, token, contract, konvensi). ADR-003
menetapkan **apa** yang dibangun: taksonomi grup final, gap vs library enterprise
terkenal, daftar komponen, prioritas, dan acceptance criteria.

---

## Context

### Kondisi Saat Ini

| | Jumlah |
|---|---|
| Komponen Stencil lama (`packages/components`) | 127 `.tsx` (128 folder) |
| Target Mitosis (ADR-001) | ~85 |
| Target paska revisi ADR-003 (riset 9 library) | ~100–110 |
| Mitosis sudah dibangun (`packages/ui`) | 65 `.lite.tsx` (59 komponen unik) |

Komponen Mitosis bukan menggantikan Stencil — ini migrasi belum selesai.
Stencil lama tetap utuh sampai engine BitCode konsumsi output Mitosis. Beberapa
primitif baru (BcAvatar, BcSpinner, BcSkeleton, BcLabel, BcPopover, BcTooltip,
BcKpi, BcKbd) tambahan di luar ADR-001 — gaya shadcn/Radix. Sekitar 70 komponen
Stencil (charts 26, fields 36, media 8, sebagian views/social/print/widgets)
belum di-port.

### Riset Pembanding (9 library)

Revisi ADR-003 memperluas riset dari 5 ke 9 library agar gap analysis menutup
domain a11y primitives, layout cross-platform, dan input-only fields yang
sebelumnya tidak tertangkap.

#### Tier 1 — Library "lengkap" (sumber utama taxonomy)

Inventaris resmi dari **Ant Design v6, MUI v9 + MUI X v9, Mantine v9,
shadcn/ui, PrimeReact v10**. Temuan kunci:

- **PrimeReact & Ant Design** paling lengkap out-of-the-box untuk ERP.
- **MUI** lengkap hanya dengan paket komersial (MUI X Pro/Premium).
- **shadcn/ui** sengaja ramping (primitif copy-paste, ~55 item) + Sonner toast
  + Sidebar block + Chart wrapper.
- **Mantine** unggul di hooks library (`@mantine/hooks` ~60+) & Spotlight.
- Komponen enterprise yang sering terlewat tapi penting: Transfer (dual-list),
  Cascader, TreeSelect, Mentions, Segmented, Steps/Stepper, Timeline,
  Descriptions, Statistic, Empty, Result, Splitter, ContextMenu, Command palette,
  OTP/PIN, TagsInput, TreeTable, editable DataGrid, Drawer, Notification (vs Toast).

#### Tier 2 — Headless / a11y primitives

Inventaris resmi dari **Radix UI v1, react-aria / react-aria-components**.
Tier 2 bukan dipakai sebagai source komponen siap-pakai, tapi sebagai sumber
**primitive a11y** dan **input granular** yang belum tertangkap Tier 1:

- **Radix:** `Slot`, `Portal`, `AspectRatio`, `VisuallyHidden`, `Toggle`/
  `ToggleGroup`, `OneTimePassword` (preview). `Slot` adalah dasar pola
  `asChild` — relevan untuk output Mitosis lintas-framework.
- **react-aria:** `DateField`/`TimeField`/`NumberField` (input-only, tanpa
  picker overlay), `ColorArea`/`ColorWheel`/`ColorSlider`/`ColorSwatch`
  (sub-primitives ColorPicker), `DropZone`, `Disclosure`/`DisclosureGroup`
  (Accordion spec-correct), `TagGroup` (read-only ≠ TagsInput),
  `Calendar`/`RangeCalendar` (terpisah dari Picker), I18n ICU built-in.

#### Tier 3 — Universal / cross-platform

Inventaris resmi dari **Tamagui v1, gluestack-ui v2**. Tier 3 menjadi sumber
pattern lintas-platform yang ADR-002 tunda implementasinya tapi contract-nya
harus ada sejak sekarang:

- **Tamagui:** `Stack`/`XStack`/`YStack`/`ZStack`, `Adapt` (responsive
  morphing: Popover desktop → Sheet mobile), `AnimatePresence` (animation
  contract), `Image` dengan placeholder/blurhash slot, `Theme` runtime
  sub-theming.
- **gluestack-ui:** `Actionsheet` (mobile bottom-sheet ≠ Drawer), `VStack`/
  `HStack`/`Box`/`Center`/`Square`/`Circle`, `Fab`, `PinInput` (=OTP),
  `KeyboardAvoidingView`/`SafeAreaView` integrations.

#### Pola lintas-tier yang muncul

Pola berikut muncul di 2+ library lintas tier, masuk taxonomy ADR-003 (Section 1):

- **Provider runtime:** antd `App`/`ConfigProvider`, Mantine `MantineProvider`,
  Tamagui `TamaguiProvider`. Berbeda dari token statis — tema/i18n/RTL/locale
  swap saat runtime.
- **Toast queue modern:** Sonner (shadcn), Mantine Notifications, react-aria
  ToastRegion. Pola stacked dengan promise-based API.
- **Hooks library portable:** Mantine ~60+, react-aria ~50+, react-use.
  Belum ada strategi `core/hooks/` di ADR.
- **Sparkline / inline chart:** MUI X Sparkline, Tamagui blocks. Tidak ada di
  Stencil 26 chart.
- **Backdrop / LoadingOverlay:** MUI Backdrop, Mantine LoadingOverlay,
  Chakra Skeleton overlay. Beda dari Dialog & Spinner.

### Tujuan

Mencapai **parity enterprise** dengan tetap memenuhi prinsip ADR-002
(simple, flexible, maintainable, powerful, lintas-platform).

---

## Section 1 — Final Taxonomy

Taksonomi sekarang (`charts, data, forms, input, layout, navigation, overlays,
primitives`) punya masalah: `input` terpisah dari `forms`; `BcProgress` ada di
`charts` (harusnya feedback); `BcKpi` di `charts` (harusnya data-display);
tidak ada grup `feedback`. Revisi (selaras antd + Mantine + Radix + Tamagui):

```
primitives/    Button, Badge, Avatar, Card, Separator, Label, Skeleton,
               Spinner, Kbd, ThemeIcon,
               Slot, Portal, AspectRatio, VisuallyHidden,         (Radix/a11y)
               Stack, XStack, YStack, ZStack, Flex, Box,           (Tamagui/gluestack)
               Center, Square, Circle                              (gluestack)
forms/         (data entry) Input, Textarea, Number, Select, MultiSelect,
               Combobox, Autocomplete, Checkbox, Radio, Toggle/Switch,
               ToggleGroup,                                        (Radix)
               Slider, RangeSlider, DateTime, ColorPicker, Rate, OTP/PIN,
               TagsInput, Cascader, TreeSelect, Transfer, Mentions, File,
               RichText, Code, Json, Signature, Barcode, Geo, Segmented,
               DateField, TimeField, NumberField,                  (react-aria, input-only)
               ColorArea, ColorWheel, ColorSlider, ColorSwatch,    (react-aria sub-primitives)
               DropZone                                            (react-aria)
data-display/  Statistic/KPI, Descriptions, Tag, Timeline, Tree, Empty,
               Image, Carousel, Avatar-group, QRCode, Watermark,
               Disclosure, DisclosureGroup,                        (react-aria, ≠Accordion)
               TagGroup,                                           (react-aria, read-only)
               ImageWithPlaceholder                                (Tamagui blurhash slot)
views/         (data views) DataTable, List, Form, Kanban, Calendar,
               Gantt, Map, Report, Tree-table, Activity, Editor
feedback/      Alert, Progress, Result, Toast, Notification, Spinner-state, Tour,
               Backdrop, LoadingOverlay,                           (MUI/Mantine)
               ToastQueue                                          (Sonner-style stacked)
navigation/    Breadcrumb, Pagination, Menu, Dropdown, Steps/Stepper,
               Anchor, Tabs, Menubar, NavigationMenu, BackTop, Affix,
               PanelMenu, TieredMenu, SlideMenu, MegaMenu,         (PrimeReact variants)
               Sidebar                                             (shadcn block)
layout/        Row, Column, Section, Header, Footer, Splitter,
               ButtonGroup, HtmlBlock, Space, Grid, Sheet,
               KeyboardAvoidingView, SafeAreaView                  (gluestack/RN-only contract)
overlays/      Dialog, AlertDialog, Prompt, Wizard, Drawer, Popover,
               Tooltip, ContextMenu, HoverCard, Command (palette),
               Actionsheet,                                        (gluestack, mobile bottom-sheet)
               Adapt                                               (Tamagui responsive morphing)
charts/        BcChart (type prop) + ChartLegend, ChartTooltip helpers,
               Sparkline                                           (MUI X)
media/         ViewerPdf, ViewerImage, ViewerDocument, ViewerVideo,
               ViewerAudio, ViewerYoutube, ViewerInstagram, ViewerTiktok, Galleria
social/        Chatter, Activity, Timeline-social
print/         Export, Print, ReportLink
utility/       Copy, Phone, Email, Url, Priority, Handle, Domain,
               Placeholder, SyncStatus, StatusBar
```

### Perubahan dari struktur sekarang
- Gabung `input/` ke `forms/`.
- Pindah `BcProgress` -> `feedback/`; `BcKpi` -> `data-display/`.
- Tambah grup: `data-display`, `feedback`, `views`, `media`, `social`, `print`, `utility`.
- `data/` lama (datatable/tree/kanban/calendar) -> `views/` (jantung ERP).
- Tambah primitives a11y dari Radix/react-aria: `Slot`, `Portal`, `AspectRatio`,
  `VisuallyHidden` — wajib untuk komposisi Mitosis lintas-framework.
- Tambah layout primitives lintas-platform dari Tamagui/gluestack: `Stack`
  family (`X/Y/Z`), `Box`, `Center`, `Square`, `Circle` — ADR-002 Tier B/C
  parity sejak awal.
- Tambah input-only fields dari react-aria: `DateField`, `TimeField`,
  `NumberField` — beda dari Picker (yang punya overlay calendar/spinner).
- Pecah ColorPicker monolitik jadi sub-primitives: `ColorArea`, `ColorWheel`,
  `ColorSlider`, `ColorSwatch` (komposisi ala react-aria).
- Tambah `Disclosure`/`DisclosureGroup` (Accordion spec-correct WAI-ARIA).
- Tambah `TagGroup` read-only (beda dari `TagsInput` editable).
- Tambah `Backdrop`/`LoadingOverlay` (≠ Dialog ≠ Spinner standalone).
- Tambah `ToastQueue` Sonner-style (stacked promise-based, ≠ Toast tunggal).
- Tambah `Adapt` & `Actionsheet` untuk mobile/cross-platform.
- Tambah `Sparkline` ke charts.
- Tambah variants menu (`PanelMenu`, `TieredMenu`, `SlideMenu`) + `Sidebar`
  block.

### Hal yang BUKAN komponen tapi taxonomy harus ada slot-nya

Pola berikut bukan komponen visual tunggal, tapi infrastruktur — diakomodir
oleh ADR-002 (lihat bagian "Cross-tier infrastructure" yang dibahas di
ADR-002 revisi):

- `core/providers/` — `App`/`ConfigProvider` runtime (theme/i18n/RTL/locale).
- `core/hooks/` — port-able hooks (use-disclosure, use-clipboard, use-debounced-value, dst).
- `core/animation/` — animation contract (AnimatePresence-style).
- `core/adapt/` — responsive morphing helper (komponen pakai, taxonomy slot di `overlays/`).

---

## Section 2 — Gap Analysis

Dua kategori gap: (A) sudah ada di Stencil lama -> tinggal migrasi;
(B) net-new (belum ada di mana pun) -> ini yang membawa ke enterprise.

### A. Tinggal Migrasi (ada di Stencil, belum di Mitosis)

- **Views (kritis):** List, Form, Gantt, Map, Report, Activity, Editor.
- **Fields:** RichText, Code, File, Lookup, Barcode, Color, Geo, Rating, Json,
  Signature (+ varian: currency, percent, duration, password, radio, tags, time...).
- **Media (belum ada satu pun):** ViewerPdf, ViewerImage, ViewerDocument,
  ViewerVideo, ViewerAudio, ViewerYoutube, ViewerInstagram, ViewerTiktok.
- **Dialogs:** AlertDialog, Prompt, Wizard, QuickEntry.
- **Layout:** ButtonGroup, HtmlBlock, Sidebar, Sheet.
- **Social:** Chatter, Activity, Timeline.
- **Print:** Export, Print, ReportLink.
- **DataTable companion:** Filter (bar/builder/panel), LookupModal, ChildTable.
- **Search:** Search, Favorites.
- **Utility:** Copy, Phone, Email, Url, Priority, Handle, Domain, Placeholder,
  SyncStatus, StatusBar.

### B. Net-New (belum ada di Stencil maupun Mitosis)

Muncul di 3+ library besar tapi tidak ada di BitCode:

- **Forms:** Slider, RangeSlider, OTP/PIN, TagsInput, Cascader, TreeSelect,
  Transfer (dual-list), Mentions, Segmented, Combobox/Autocomplete (standalone),
  MultiSelect, ColorPicker (proper, standalone).
- **Data-display:** Descriptions, Statistic (generik), Empty, Timeline (generik),
  QRCode, Watermark, Carousel.
- **Feedback:** Result, Notification (beda dari Toast), Tour/Onboarding.
- **Navigation:** Steps/Stepper (standalone), Anchor, Menubar, NavigationMenu,
  BackTop, Affix, MegaMenu.
- **Overlays:** Drawer (proper), ContextMenu, HoverCard, Command palette.
- **Layout:** Splitter/Resizable.
- **Specialized:** OrgChart (untuk HRM), Inplace (inline edit), Galleria,
  VirtualScroller, MeterGroup, Knob, Dock.

### C. Net-New dari Tier 2/3 (Radix, react-aria, Tamagui, gluestack)

Pattern di bawah ini muncul di 2+ library Tier 2/3 tapi belum tertangkap di
Tier 1 maupun di Stencil/Mitosis BitCode. Mereka penting untuk a11y,
komposisi Mitosis, dan parity cross-platform:

- **A11y primitives (Radix, react-aria):** `Slot` (composition pattern,
  `asChild`), `Portal`, `AspectRatio`, `VisuallyHidden`, `Toggle`/
  `ToggleGroup`, `Disclosure`/`DisclosureGroup` (≠Accordion), `TagGroup`
  (read-only ≠ TagsInput).
- **Input-only fields (react-aria, Mantine):** `DateField`, `TimeField`,
  `NumberField` — locale-aware, tanpa picker overlay. Penting untuk form
  enterprise yang butuh input cepat tanpa modal.
- **ColorPicker sub-primitives (react-aria):** `ColorArea`, `ColorWheel`,
  `ColorSlider`, `ColorSwatch`, `ColorField` — komposisi power-user; default
  `BcColorPicker` tetap ada sebagai shorthand.
- **Layout primitives (Tamagui, gluestack, MUI, Chakra):** `Stack`/`XStack`/
  `YStack`/`ZStack`, `Flex`, `Box`, `Center`, `Square`, `Circle`. Belum ada di
  taxonomy ADR sebelumnya yang hanya punya `Row`/`Column`.
- **Cross-platform overlays (Tamagui, gluestack):** `Adapt` (Popover desktop
  → Sheet mobile), `Actionsheet` (mobile bottom-sheet ≠ Drawer).
- **Modern feedback:** `Backdrop`/`LoadingOverlay` (MUI/Mantine), Sonner-style
  `ToastQueue` (shadcn).
- **Charts:** `Sparkline` inline (MUI X) — Stencil 26 chart tidak ada.
- **Image:** `Image` dengan `placeholder`/`blurhash` slot (Tamagui, Next/Image).
- **Menu variants (PrimeReact):** `PanelMenu`, `TieredMenu`, `SlideMenu`
  (di luar `MegaMenu` yang sudah di P1).
- **DropZone (react-aria):** wilayah upload terpisah dari `BcFieldFile`.

### D. Cross-tier infrastructure (BUKAN komponen, slot wajib di ADR-002)

Diputuskan masuk ADR-002 (cross-platform architecture) karena lintas grup:

- `core/providers/` — `App`/`ConfigProvider` runtime (theme/i18n/RTL/locale
  swap saat runtime, beda dari token statis CSS).
- `core/hooks/` — port-able hooks lintas-framework (`useDisclosure`,
  `useClipboard`, `useDebouncedValue`, `useElementSize`, `useMediaQuery`,
  `useFocusTrap`, `useId`, `useControllableState`, dst). Kandidat: subset
  Mantine + react-aria yang murni logic (no DOM).
- `core/animation/` — animation contract (`AnimatePresence`/`Transition`).
  Tier A pakai Framer/CSS; Tier B Tamagui native; Tier C Flutter implicit.
- `core/adapt/` — responsive morphing helper (Popover ↔ Sheet, Dropdown ↔
  Actionsheet) berdasarkan media query / platform.

---

## Section 3 — Prioritization

Prioritas berdasarkan: nilai untuk ERP, kemunculan di library terkenal,
ketergantungan modul BitCode (HRM, CRM, dll).

### P0 — ERP-Critical (wajib enterprise grade)

| Komponen | Alasan | Ada di |
|---|---|---|
| BcForm | Jantung low-code; render + validasi dari definisi | semua |
| BcList | View list utama | semua |
| BcDataTable (upgrade) | Grid data-heavy enterprise (lihat Section 4) | semua |
| BcDescriptions | Tampilan key-value detail record | antd, Prime |
| BcSteps/BcStepper | Approval/multi-step standalone | antd, MUI, Mantine, Prime |
| BcAccordion/BcCollapse | Standar mutlak, belum ada | semua |
| BcContextMenu | Klik-kanan grid/tree | Mantine, shadcn, Prime, antd |
| BcSplitter | Layout master-detail | antd, Mantine, shadcn, Prime |
| BcDrawer | Panel samping (proper) | semua |
| BcTransfer | Assign role/permission (dual-list) | antd, MUI, Prime |
| BcCascader | Select hierarki | antd, Prime |
| BcTreeSelect | Org unit, chart of accounts | antd, Mantine, Prime |
| BcStatistic | Stat card dashboard | antd |
| BcEmpty + BcResult | Empty state & halaman hasil | antd, shadcn |
| BcOtp/BcPin | UI 2FA (engine sudah punya 2FA) | antd, Mantine, shadcn, Prime |
| **BcSlot, BcPortal** | Composition primitive (`asChild`) — fondasi Mitosis | Radix |
| **BcVisuallyHidden** | A11y dasar wajib | Radix, react-aria |
| **BcStack family** (Stack/XStack/YStack/Box) | Layout primitive lintas-platform | Tamagui, gluestack, MUI |
| **BcBackdrop/BcLoadingOverlay** | Loading state proper (≠ Spinner standalone) | MUI, Mantine |
| **BcToastQueue** (Sonner-style) | Stacked promise-based notifications | shadcn |

### P1 — High Value

BcCommand (palette), BcNotification, BcSegmented, BcMentions (lengkapi Chatter),
BcSlider/BcRangeSlider, BcAutocomplete/BcCombobox, BcMenubar/BcMegaMenu/
BcNavigationMenu, BcOrgChart (HRM), BcTour, BcAffix+BcAnchor+BcBackTop,
BcColorPicker, BcCarousel, BcInplace, BcAlertDialog/BcPrompt/BcWizard,
sisa Fields (RichText, Code, File, Json, Signature, Barcode, Geo, Rating),
Media viewers (8), Social (Chatter, Activity, Timeline).

**P1 tambahan dari Tier 2/3:**
- **BcDateField/BcTimeField/BcNumberField** — input-only, locale-aware
  (react-aria). Pelengkap field picker yang sudah ada.
- **BcAspectRatio** — wrapper utility a11y/layout (Radix, shadcn).
- **BcDisclosure/BcDisclosureGroup** — Accordion spec-correct WAI-ARIA
  (react-aria). Boleh hidup berdampingan dengan `BcAccordion` legacy.
- **BcActionsheet** — mobile bottom-sheet (gluestack); paired dengan
  **BcAdapt** (Tamagui responsive morphing).
- **BcFlex** — primitive layout generik (antd v5+, Chakra) di luar
  Row/Column.
- **BcDropZone** — region upload terpisah dari `BcFieldFile` (react-aria).
- **BcTagGroup** — read-only tag list dengan dismiss (react-aria).
- **BcToggleGroup** — segmented-like multi-toggle (Radix), beda dari
  `BcSegmented` single-select.
- **Sparkline chart** — masuk family `BcChart` via `type="sparkline"` atau
  komponen mandiri `BcSparkline` (MUI X).

### P2 — Polish

BcWatermark, BcQRCode, BcSpeedDial/BcFloatButton, BcKbd, BcVirtualScroller,
BcKnob, BcMeterGroup, BcDock, BcGalleria, sisa Views (Gantt, Map, Report, Editor).

**P2 tambahan dari Tier 2/3:**
- **BcColorArea/BcColorWheel/BcColorSlider/BcColorSwatch** — sub-primitives
  ColorPicker untuk power-user (react-aria).
- **BcImage** dengan `placeholder`/`blurhash` slot (Tamagui).
- **BcCenter/BcSquare/BcCircle** — layout shorthand (gluestack).
- **BcPanelMenu/BcTieredMenu/BcSlideMenu** — variants menu (PrimeReact).
- **BcSidebar block** — composable navigation frame (shadcn).

---

## Section 4 — BcDataTable Enterprise Spec

Gap terbesar bukan komponen baru, tapi fitur yang membedakan grid biasa vs
enterprise (referensi: MUI X Premium, PrimeReact DataTable, antd ProTable).
PrimeReact jadi acuan utama karena paling lengkap tanpa lisensi berbayar.

### Fitur Wajib (enterprise)

| Fitur | Catatan |
|---|---|
| 4-layer data | local/fetcher/url/model + state idle/loading/error/empty/ready |
| Pagination | client + server (lazy), pageSize selectable |
| Sort | single + multi-sort |
| Filter | per-kolom, operator (eq/like/gt/between/in...), global search |
| Inline edit | cell + row mode |
| Selection | row (single/multi/checkbox) + cell range |
| Column ops | resize, reorder, show/hide, pin/freeze (left/right) |
| Row ops | reorder (drag), expandable rows |
| Grouping | row grouping + footer aggregate (sum/avg/count) |
| Tree-table | data hierarkis |
| Virtualization | 10k+ baris tanpa lag |
| Export | CSV (built-in), Excel/PDF (adapter opsional) |
| State persistence | simpan sort/filter/kolom (local/session/server) |
| Context menu | integrasi BcContextMenu |
| A11y | header role, keyboard nav, aria-sort |
| RTL + i18n | label/format lokal |

### Pendekatan Implementasi

- Bangun di atas **DataView Contract** (ADR-002 Section 4.3).
- Pisahkan engine (sort/filter/paginate/group) dari render -> reusable untuk
  BcList & BcTree-table.
- Virtualization & berat-fitur via `useTarget` (web pakai virtualizer; RN pakai
  FlatList; Flutter pakai ListView.builder) -> tetap satu contract.
- Fitur berat opt-in (prop flag) supaya kasus sederhana tetap ringan (prinsip simple).

### Acceptance Criteria
- [ ] 4-layer data + semua state ter-surface (termasuk error).
- [ ] Sort/filter/paginate server-side & client-side.
- [ ] Inline edit, selection, column resize/reorder/pin.
- [ ] Tree-table & grouping + aggregate.
- [ ] Virtualization terverifikasi pada 10k baris.
- [ ] Export CSV; a11y keyboard; jalan di Tier A.

---

## Section 5 — Build Sequence

Mengikuti fase ADR-002. Komponen dibangun **setelah** fondasi (token+contract) siap.

### Urutan disarankan

1. **Fondasi dulu** (ADR-002 Phase 0-2): token, contract, export, BcInput refactor,
   BcDataTable rewrite, BcSelect a11y. Tanpa ini, komponen baru menumpuk utang.
2. **Views inti** (P0): BcForm, BcList. Jantung low-code ERP.
3. **P0 net-new ringan**: Descriptions, Steps, Accordion, Empty, Result, Splitter,
   Drawer, ContextMenu, OTP/PIN.
4. **P0 forms kompleks**: Cascader, TreeSelect, Transfer, Statistic.
5. **DataTable enterprise upgrade** (Section 4) — paralel dengan #3/#4 bila kapasitas ada.
6. **P1** sesuai kebutuhan modul (Mentions+Chatter, OrgChart+HRM, Command, Media viewers).
7. **P2** polish.
8. **Tier B/C** (ADR-002 Phase 5): generator Tamagui/gluestack + Flutter codegen.

### Prinsip eksekusi
- Tiap komponen patuh Definition of Done (ADR-002 Section 5.3).
- Net-new dibuat dengan contract penuh sejak awal (jangan retrofit).
- Migrasi dari Stencil: port perilaku, bukan salin mentah; sesuaikan ke contract+token.
- Jangan hapus Stencil lama sampai engine BitCode memakai output Mitosis (dirombak
  setelah UI matang, sesuai arahan stakeholder).

---

## Acceptance Criteria (ADR-level)

- [ ] Taksonomi final diterapkan pada struktur folder.
- [ ] Semua P0 selesai dengan contract penuh + a11y + jalan di Tier A.
- [ ] BcDataTable memenuhi Enterprise Spec (Section 4).
- [ ] Parity enterprise terdokumentasi (checklist vs antd/Prime).
- [ ] Registry & docs ter-update tiap komponen baru.

---

## Consequences

### Positive
- Peta jelas menuju parity enterprise; prioritas berbasis nilai ERP nyata.
- Net-new mengisi gap yang membedakan dari "UI kit biasa".

### Negative
- Volume besar (~85+ komponen target + upgrade grid). Butuh disiplin fase.
- Migrasi + net-new bersamaan -> risiko konteks; mitigasi: P0 dulu, AC ketat.

### Risks
- Scope creep -> kunci ke P0/P1/P2; P2 hanya setelah P0 solid.
- Konflik dengan agent lain (UI dikerjakan orang lain) -> koordinasi sebelum eksekusi.

---

## References
- ADR-001 (Mitosis compiler), ADR-002 (cross-platform architecture).
- **Tier 1 (lengkap):** Ant Design v6, MUI v9 + MUI X v9, Mantine v9,
  shadcn/ui (Sonner, Sidebar block, Chart wrapper), PrimeReact v10.
- **Tier 2 (a11y primitives):** Radix UI v1, react-aria / react-aria-components.
- **Tier 3 (universal/cross-platform):** Tamagui v1, gluestack-ui v2.
- Inventaris faktual BitCode (per revisi): 127 `.tsx` Stencil (128 folder),
  65 `.lite.tsx` Mitosis (59 komponen unik).
