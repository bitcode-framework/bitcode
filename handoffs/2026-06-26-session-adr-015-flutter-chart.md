# Session Instruksi: ADR-015 — Flutter Chart Strategy (flutter_echarts & Shared Chart Option Contract)

> **Untuk session baru.** Baca file ini + ADR-015 sampai habis sebelum mulai coding.
> Repo root: `D:\NodeJS\lowcode`. UI-Flutter package: `D:\NodeJS\lowcode\packages\ui-flutter`.

---

## Cara Memulai

```
Kerjakan ADR-015: Flutter Chart Strategy — flutter_echarts & Shared Chart Option Contract.
Baca docs/adr/015-flutter-echarts-chart-strategy.md sampai habis sebelum mulai.

WAJIB berpikir kritis, detail, mateng, lengkap dan jujur.
Setelah selesai per phase, update semua docs terkait lalu commit.

PENTING:
- Ada agent lain yang mungkin sedang mengerjakan ADR-013 (Cross-Tier Infrastructure)
  dan ADR-014 (Engine UI Wiring) secara paralel. Jangan sembarangan remove atau
  merusak kerjaan agent lain. Konfirmasikan jika butuh hal di luar scope kamu.
- Scope file kamu sudah diatur di ADR-015 section "Scope Boundary untuk Paralelisasi".
  PATUHI boundary tersebut.
- ADR-015 ZERO overlap dengan ADR-013 dan ADR-014 (file chart terisolasi).
  Kamu bisa full paralel tanpa khawatir bentrok, asal patuhi boundary.
```

---

## ADR Referensi

- **ADR-015:** `docs/adr/015-flutter-echarts-chart-strategy.md` (Proposed)
- **Supersedes (partial):** ADR-011 §"Chart widget" (M5 CustomPaint approach)
- **Related:** ADR-002 (Cross-Platform), ADR-003 (Component Roadmap §1 Chart Taxonomy),
  ADR-011 (Flutter Adapter Architecture), ADR-013 (Cross-tier Infra), ADR-014 (Engine UI Wiring)

---

## Konteks Singkat

BitCode UI punya 3 layer chart dengan parity tidak seimbang:

| Layer | Chart Types | Engine |
|-------|-------------|--------|
| Stencil (legacy, 26 komponen) | 26/26 | Apache ECharts ✅ |
| Mitosis (`@bitcode/ui`, 1 komponen `BcChart` discriminator `type`) | 26/26 | Apache ECharts (lazy import) ✅ |
| Flutter (`@bitcode/ui-flutter`, 1 komponen `BcChart`) | **4/26** (bar, line, pie, doughnut) | **CustomPaint manual ❌** |

CustomPaint Flutter saat ini (`bc_chart.dart`, 912 lines) adalah **hobby-grade, bukan production-grade**:
- Bug value negatif (barH negatif → crash/render kosong, tidak ada baseline 0)
- Tidak ada tooltip (prop `tooltipEnabled` ada tapi diabaikan)
- Tidak ada animasi (prop `animate` ada tapi render static)
- Tidak ada Y-axis labels (grid lines ada, angka skala tidak ada)
- Tidak ada nice-number axis scaling (maxVal manual)
- Tidak ada legend interactivity (legend statis, tidak bisa hide/show series)
- `shouldRepaint` inefficient (compare List by reference)
- Stacked math redundant (maxVal dihitung 2x)

**Keputusan (Opsi C):** Full `flutter_echarts`, deprecate CustomPaint.
- `flutter_echarts` (by entronad) bungkus Apache ECharts di WebView Flutter
- Langsung dapat 26 type + tooltip + animasi + zoom + toolbox + themes + a11y
- Contract shared: `buildChartOption()` di Mitosis `factory.ts` = single source of truth,
  Flutter terima option yang sama sebagai JSON string → kirim ke WebView
- Zero duplication logika chart antara Mitosis dan Flutter
- WebView overhead acceptable (1-3 chart per dashboard screen)

---

## Scope File (PATUHI untuk Paralelisasi)

### ✅ Owned (boleh modify)

```
packages/ui/src/components/charts/              # Mitosis chart (factory extraction, refactor)
  ├── chart/
  │   ├── BcChart.lite.tsx                       # refaktor: pakai shared option builder
  │   └── options/                               # extract to shared contract
  │       ├── factory.ts                         # → jadi single source of truth
  │       ├── helpers.ts
  │       ├── types.ts
  │       └── *.test.ts
  └── (sparkline, BcKpi, BcProgress — jangan ganggu, milik chart group tapi ADR-015 fokus BcChart)

packages/ui-flutter/                             # Flutter chart rewrite
  ├── pubspec.yaml                               # tambah dependency flutter_echarts
  ├── lib/src/components/
  │   ├── bc_chart.dart                          # REWRITE: CustomPaint → flutter_echarts
  │   └── bc_chart_types.dart                    # NEW: chart type enum + helpers (JANGAN sentuh core/types.dart)
  ├── lib/src/core/
  │   └── chart_option_bridge.dart               # NEW: JSON option → WebView bridge
  └── test/
      ├── bc_chart_test.dart                     # rewrite tests
      └── bc_chart_echarts_test.dart             # NEW: parity tests 26 type

# Shared chart option contract (NEW location — pilih salah satu):
packages/ui/src/contracts/chart-option/          # Opsi A: di @bitcode/ui (TS)
  ├── chart-option.schema.ts                     # JSON schema for ECharts option
  └── chart-option.builder.ts                    # = factory.ts yang di-extract
# ATAU
packages/ui-schema/src/chart/                    # Opsi B: di @bitcode/ui-schema (per ADR-012)
  ├── option-builder.ts
  └── types.ts
```

### ❌ DILARANG touch (milik agent lain / shared coordinator)

```
packages/ui/src/index.ts                         # barrel export — MERGE COORDINATOR only
packages/ui/src/core/                            # ADR-013 territory (hooks/providers/adapt/animation)
packages/ui/src/core/types.ts                    # ADR-013 territory
packages/ui-flutter/lib/src/core/types.dart      # ADR-013 territory (reactive store) — pakai bc_chart_types.dart sendiri
packages/ui/src/components/{forms,overlays,layout,views,...}/  # ADR-013/014 territory
engine/internal/presentation/                    # ADR-014 territory (engine Go)
docs/components-progress.md                      # MERGE COORDINATOR only
docs/adr/README.md                               # MERGE COORDINATOR only
```

### ⚠️ Koordinasi dengan agent lain

- **ADR-013 (Cross-tier Infra):** ZERO file overlap dengan ADR-015. Kamu aman.
- **ADR-014 (Engine UI Wiring):** ZERO file overlap. ADR-014 mungkin konsumsi `BcChart`
  sebagai black box untuk pivot/export, tapi tidak modify chart internals. Kalau ADR-014
  agent butuh chart API change, koordinasi via issue/comment, bukan modify langsung.

---

## MILESTONES & ACCEPTANCE CRITERIA

### Phase 1: Chart Option Contract Extraction

**Goal:** `buildChartOption()` jadi shared contract, terpisah dari Mitosis component.

- [ ] 1.1 Pilih lokasi contract (Opsi A `packages/ui/src/contracts/chart-option/` atau
      Opsi B `packages/ui-schema/src/chart/`). **Saran:** Opsi A dulu (lebih simpel,
      ADR-012 ui-schema belum fully extracted). Kalau ADR-012 sudah done, migrate ke Opsi B.
- [ ] 1.2 Extract `factory.ts` → `chart-option.builder.ts` (pure function, no Mitosis deps)
- [ ] 1.3 Extract `types.ts` → `chart-option.schema.ts` (BcChartType, ChartData, EChartsOption, BuildOptionContext)
- [ ] 1.4 `BcChart.lite.tsx` import dari contract baru (bukan dari `./options/`)
- [ ] 1.5 Tests: semua test chart tetap pass (no regression)
- [ ] 1.6 Contract audit: 0 violations
- [ ] 1.7 Docs: update `packages/components/docs/` chart docs, `packages/ui/docs/api/BcChart.md`

**Verify:**
```powershell
cd D:\NodeJS\lowcode\packages\ui
npm run test -- --filter chart    # atau: npx vitest run src/components/charts
npm run build                      # Mitosis build tetap OK
```

### Phase 2: flutter_echarts Integration

**Goal:** Flutter `BcChart` pakai `flutter_echarts`, deprecate CustomPaint.

- [ ] 2.1 `pubspec.yaml`: tambah `flutter_echarts: ^latest` + `webview_flutter`
- [ ] 2.2 Buat `bc_chart_types.dart` — port `BcChartType` enum (26 type) +
      `BcChartSeries`, `BcChartItem` (jangan sentuh `core/types.dart`)
- [ ] 2.3 Buat `chart_option_bridge.dart` — terima JSON option string, kirim ke
      `flutter_echarts` WebView, handle lifecycle (dispose, resize)
- [ ] 2.4 REWRITE `bc_chart.dart`:
  - Hapus semua `_BarChartPainter`, `_LineChartPainter`, `_PieChartView` (912 lines CustomPaint)
  - Ganti dengan `flutter_echarts` widget yang konsumsi option JSON
  - Pertahankan API surface: `type`, `series`, `categories`, `items`, `colors`, `legend`,
    `tooltipEnabled`, `animate`, `height`, `width`, `loading`, `stacked`, `horizontal`,
    `ariaLabel`, `data` (semua props existing tetap ada)
  - Build ECharts option di Dart (mirror `factory.ts` logic) ATAU terima dari caller
- [ ] 2.5 Loading state: pakai ECharts built-in `showLoading` option
- [ ] 2.6 Accessibility: pakai ECharts `aria` option (a11y data table fallback)
- [ ] 2.7 Theme: map `BcTheme` colors → ECharts theme (runtime adaptive dark/light)

**Verify:**
```powershell
cd D:\NodeJS\lowcode\packages\ui-flutter
flutter pub get
flutter analyze                   # 0 errors/warnings
flutter test test/bc_chart_test.dart
```

### Phase 3: 26/26 Parity Tests

**Goal:** Buktikan semua 26 chart type render di Flutter via `flutter_echarts`.

- [ ] 3.1 Buat `bc_chart_echarts_test.dart` — test semua 26 `BcChartType`:
  - area, bar, boxplot, candlestick, custom, funnel, gauge, graph, heatmap,
    kpi, line, mixed, parallel, pictorialbar, pie, pivot, polar, progress,
    radar, sankey, scatter, scorecard, sunburst, themeriver, tree, treemap
  - Plus `doughnut` (Flutter variant)
- [ ] 3.2 Untuk setiap type: assert
  - Widget tree build tanpa error
  - ECharts option JSON ter-generate dengan `type` field yang benar
  - WebView ter-inject dengan option yang benar
  - Loading state berfungsi
  - Empty data state berfungsi
- [ ] 3.3 Golden test (opsional): screenshot parity untuk 3-5 type utama (bar, line, pie, heatmap, radar)
- [ ] 3.4 Compare option JSON dengan output `factory.ts` Mitosis (parity check)

**Verify:**
```powershell
cd D:\NodeJS\lowcode\packages\ui-flutter
flutter test test/bc_chart_echarts_test.dart
# Semua 26 type + doughnut = 27 test cases pass
```

### Phase 4: Migration & Deprecation

**Goal:** CustomPaint dihapus bersih, migration path jelas.

- [ ] 4.1 Hapus semua CustomPaint painter classes (`_BarChartPainter`, `_LineChartPainter`,
      `_PieChartView`, `_BarChartView`, `_LineChartView`)
- [ ] 4.2 Hapus `_defaultPalette` hardcode (pakai ECharts theme)
- [ ] 4.3 Update `CONTRACT_INVENTORY.md` — chart entry reflect new strategy
- [ ] 4.4 Update `packages/ui-flutter/README.md` — chart section: `flutter_echarts` strategy
- [ ] 4.5 Deprecation note di CHANGELOG: "CustomPaint chart deprecated, replaced by flutter_echarts"
- [ ] 4.6 Verify: `flutter analyze` 0 errors, `flutter test` all pass (1158+ tests, 0 fail)

**Verify:**
```powershell
cd D:\NodeJS\lowcode\packages\ui-flutter
flutter analyze
flutter test                       # semua test pass
```

### Phase 5: Docs & Commit

- [ ] 5.1 Update `docs/components-progress.md` — chart row: Flutter ✅ 26/26 via flutter_echarts
      (MERGE COORDINATOR bantu kalau bentrok)
- [ ] 5.2 Update `packages/components/docs/bc-chart.md` (kalau ada) atau `packages/ui/docs/api/BcChart.md`
- [ ] 5.3 Update ADR-015 status: Proposed → Accepted (setelah semua phase done)
- [ ] 5.4 Update `engine/docs/features/` — chart feature doc reflect cross-platform parity
- [ ] 5.5 Commit per phase dengan message:
      - `feat(ui): extract chart option contract (ADR-015 Phase 1)`
      - `feat(ui-flutter): integrate flutter_echarts, deprecate CustomPaint (ADR-015 Phase 2)`
      - `test(ui-flutter): 26/26 chart type parity via flutter_echarts (ADR-015 Phase 3)`
      - `refactor(ui-flutter): remove deprecated CustomPaint chart (ADR-015 Phase 4)`
      - `docs(adr): ADR-015 accepted, chart parity achieved (ADR-015 Phase 5)`

---

## Hal Kritis yang Harus Diperhatikan

### 1. flutter_echarts package health

Cek dulu sebelum mulai:
```powershell
cd D:\NodeJS\lowcode\packages\ui-flutter
flutter pub add flutter_echarts --dry-run   # cek latest version & compatibility
```

Kalau package unmaintained / broken, **STOP dan konfirmasi stakeholder** sebelum lanjut.
Alternatif: fork package, atau raw WebView + ECharts CDN (ADR-015 Section "Alternatives D").

### 2. ECharts option JSON serialization

`buildChartOption()` di TS menghasilkan `Record<string, unknown>`. Flutter terima sebagai
JSON string. **Pastikan:**
- Tidak ada function di option (ECharts option pure data, function tidak JSON-serializable)
- Formatter (tooltip, axis label) harus pakai ECharts template string, bukan JS function
- Kalau perlu custom formatter, dokumentasikan sebagai "web-only" feature

### 3. WebView lifecycle management

Setiap `BcChart` widget punya WebView sendiri. **WAJIB:**
- `dispose()` dipanggil ketika widget unmount → WebView release memory
- `ResizeObserver` equivalent di Flutter: `LayoutBuilder` + debounce rebuild
- Jangan create WebView baru setiap rebuild (cache instance)

### 4. Dark/Light theme runtime

CustomPaint lama hardcode `BcTheme.primary` dll. ECharts theme harus **runtime adaptive**:
- Dengarkan theme change (via `AppProvider` kalau ADR-013 sudah port, atau `BcTheme` listener)
- Rebuild ECharts option dengan color palette baru
- Jangan hardcode palette di Dart (pakai contract dari `factory.ts`)

### 5. Jangan ganggu BcKpi, BcProgress, BcSparkline

`BcKpi`, `BcProgress`, `BcSparkline` ada di `packages/ui/src/components/charts/` tapi
**TIDAK termasuk scope ADR-015** (mereka bukan ECharts-based). Jangan modify file mereka.
ADR-015 fokus hanya `BcChart` (discriminator `type`, 26 type).

### 6. Backwards compatibility

API surface `BcChart` Flutter harus tetap sama:
- Semua props existing (`type`, `series`, `categories`, `items`, `colors`, `legend`,
  `tooltipEnabled`, `animate`, `height`, `width`, `loading`, `stacked`, `horizontal`,
  `ariaLabel`, `data`) tetap ada dengan signature yang sama
- Consumer yang sudah pakai `BcChart(type: BcChartType.bar, ...)` tidak perlu ubah kode
- Yang berubah hanya internals: CustomPaint → flutter_echarts

---

## Testing Strategy

| Layer | Test | Target |
|-------|------|--------|
| Contract | `chart-option.builder.test.ts` (TS) | factory.ts extraction no regression |
| Flutter unit | `bc_chart_test.dart` | widget build, props, loading/empty state |
| Flutter parity | `bc_chart_echarts_test.dart` | 26 type + doughnut = 27 cases |
| Flutter golden | `bc_chart_golden_test.dart` (opsional) | screenshot parity 3-5 type |
| Integration | manual: demo app render chart | visual confirm di device/emulator |

---

## Referensi Cepat

| Apa | Dimana |
|-----|--------|
| ADR-015 full | `docs/adr/015-flutter-echarts-chart-strategy.md` |
| CustomPaint yang akan dihapus | `packages/ui-flutter/lib/src/components/bc_chart.dart` (912 lines) |
| Factory TS (single source of truth) | `packages/ui/src/components/charts/options/factory.ts` |
| Types TS | `packages/ui/src/components/charts/options/types.ts` |
| BcChartType union (26 type) | `types.ts` → `BcChartType` |
| Stencil chart reference (26 komponen) | `packages/components/src/components/charts/bc-chart-*/` |
| flutter_echarts package | https://pub.dev/packages/flutter_echarts |
| ECharts docs | https://echarts.apache.org/ |
| Flutter testing | `packages/ui-flutter/test/` (1158 tests existing) |
| ADR-011 (Flutter adapter, superceded partial) | `docs/adr/011-flutter-adapter-architecture.md` |

---

## Estimasi Effort

| Phase | Estimasi | Dependency |
|-------|----------|------------|
| Phase 1: Contract extraction | 1-2 hari | None (fully isolated) |
| Phase 2: flutter_echarts integration | 3-5 hari | Phase 1 done |
| Phase 3: 26/26 parity tests | 2-3 hari | Phase 2 done |
| Phase 4: Migration & deprecation | 1-2 hari | Phase 3 done |
| Phase 5: Docs & commit | 0.5 hari | All phases done |
| **Total** | **7-12 hari** | |

**Bisa full paralel dengan ADR-013 dan ADR-014** (zero file overlap).

---

## ✅ COMPLETION SUMMARY (2026-06-26)

ADR-015 is **COMPLETE** — all 4 features done. 26/26 chart type parity achieved across Stencil, Mitosis, and Flutter layers.

### Commits
| Phase | Commit | Description |
|-------|--------|-------------|
| F1 | `b2f4c84` (+ submodule `7c2859e`) | Extract chart-option contract to `packages/ui/src/contracts/chart-option/` |
| F2 | `87e511f` | Dart port of chart-option contract with golden-fixture parity (48 tests, 0 diff) |
| F3 | `e4739c7` | Rewrite BcChart with flutter_echarts wrapper, 26-type parity (CustomPaint removed) |
| F4 | (this commit) | Migration cleanup, docs, ADR-015 → Accepted |

### What was delivered
- **Contract extraction (F1):** TS `factory.ts` → `packages/ui/src/contracts/chart-option/` (pure functions, no Mitosis deps)
- **Dart contract port (F2):** `packages/ui-flutter/lib/src/core/chart/` — `chart_types.dart` (26-value enum), `chart_helpers.dart`, `chart_factory.dart` (buildChartOption), `chart_contract.dart` (barrel). 48 golden-fixture parity tests + 34 helper tests = 82 tests, all pass with 0 diff vs TS output.
- **flutter_echarts rewrite (F3):** `bc_chart.dart` 912 lines → 288 lines. All CustomPaint code removed. 29 widget-contract tests (WC-001..012 + RG-001/002). Fake WebViewPlatform test helper for widget tests.
- **Migration cleanup (F4):** Migration guide, docs updated, ADR-015 → Accepted.

### Test results
- Chart tests: 111 pass (48 factory + 34 helpers + 29 bc_chart), 0 fail
- Full suite: 1237 pass, 7 fail (all pre-existing `bc_app_config.dart` errors from ADR-013, NOT our changes)
- `flutter analyze`: 0 issues on all chart files

### Key decisions
1. **Fake WebViewPlatform** instead of `@visibleForTesting` — the standard, recommended approach per webview_flutter's own error message. Created `test/helpers/fake_webview_platform.dart` with no-op stubs for `WebViewPlatform`, `PlatformWebViewController`, `PlatformWebViewWidget`, `PlatformNavigationDelegate`, `PlatformWebViewCookieManager`.
2. **CustomPaint removed entirely** (not deprecated behind a flag) — cleaner than the ADR's original v1.0/v1.1.0 plan. Zero legacy rendering code remains.
3. **XS-001 honored**: widget tests do NOT assert on rendered pixels (WebView can't render in test env). Tests verify the widget *contract* surface — container structure, option JSON, loading/empty states, semantics, data normalization, backward-compat.

---

## Selesai Jika

- [x] `packages/ui/src/contracts/chart-option/` berisi shared contract
- [x] `BcChart.lite.tsx` konsumsi contract baru, semua test pass
- [x] `bc_chart.dart` pakai `flutter_echarts`, CustomPaint dihapus (912 lines → 288 lines)
- [x] 26/26 chart type + doughnut parity test pass di Flutter (48 golden-fixture + 29 widget tests)
- [x] `flutter analyze` 0 errors/warnings on all chart files
- [x] `flutter test` all chart tests pass (111 tests, 0 fail); 7 pre-existing bc_app_config failures are ADR-013's
- [x] API surface `BcChart` Flutter backwards compatible (all legacy props preserved)
- [x] ADR-015 status: Accepted
- [x] Docs updated: components-progress, chart-migration guide, CONTRACT_INVENTORY, codebase.md, ADR-015
- [x] Commit per phase: `b2f4c84`, `87e511f`, `e4739c7`, F4 (this commit)
