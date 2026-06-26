# ADR-013: Cross-Tier Infrastructure Completion

- **Status:** Accepted
- **Date:** 2026-06-26
- **Supersedes:** —
- **Amends:** ADR-002 (Section 9 — Cross-tier Infrastructure), ADR-006
- **Related:** ADR-004 (Reactive Engine), ADR-011 (Flutter Adapter)

## Phase Status

| Phase | Description | Status |
|-------|-------------|--------|
| 1 | AppProvider Consumption Surface (useAppContext + useAppConfigStore hooks, Flutter AppConfigStore + BcAppConfig InheritedWidget parity) | ✅ Complete (2026-06-26) |
| 2 | Hooks P0 Adoption (BcDialog/Popover/Tooltip/Drawer/Dropdown consume P0 hooks) | ⬜ Pending |
| 3 | Hooks P1 New (8 industry-standard headless hooks) | ⬜ Pending |
| 4 | Animation + Adapt contracts (per-tier implementation) | ⬜ Pending |
| 5a | Flutter env fix + Dart parity baseline | ⬜ Pending |
| 5b | Flutter parity completion + cross-tier parity table | ⬜ Pending |

## Context

ADR-002 Section 9 mendefinisikan 5 cross-tier infrastructure yang **wajib** tersedia
sebelum komponen dianggap "truly cross-platform": AppProvider, Hooks, Animation contract,
Adapt contract, dan Reactive store.

Status audit per 2026-06-26:

| Item | Status | Notes |
|------|--------|-------|
| AppProvider adoption | Partial | Unit tests pass, belum diadopsi ke semua komponen |
| Hooks subset P0 (8 hooks) | Partial | Implementasi ada, belum dipakai komponen P0 |
| Hooks subset P1 (6 hooks) | Belum | useHover, useFocusVisible, useEventListener, dll |
| Animation contract | Partial | Props exported, implementasi per-tier belum |
| Adapt contract | Partial | Props exported, runtime Popover->Sheet belum |

Saat ini 110/110 komponen enterprise "Done" secara contract+a11y+token+tests,
tetapi **tidak semua benar-benar memakai cross-tier infrastructure**. Banyak
komponen masih punya logic lokal yang seharusnya didelegasikan ke hooks/provider.

## Decision

Selesaikan 5 cross-tier infrastructure agar ADR-002 Section 9 fully realized:

### 1. AppProvider Adoption (full)

- **Scope:** `packages/ui/src/core/providers/AppProvider.lite.tsx` sudah ada.
  Adoption pass = refactor komponen P0/P1 untuk konsumsi `useAppContext()` untuk
  theme/locale/direction/tenant daripada baca `BcSetup` global langsung.
- **Flutter parity:** Port `AppProvider` ke `packages/ui-flutter/lib/src/core/app_provider.dart`
  via `InheritedWidget`. Komponen Flutter baca via `AppProvider.of(context)`.

### 2. Hooks P0 (8 hooks) — adoption + completion

Daftar 8 hooks P0 (sudah ada 9 file di `packages/ui/src/core/hooks/`):

| Hook | Status file | Adoption needed |
|------|-------------|-----------------|
| useId | ✅ ada | ya |
| useDisclosure | ✅ ada | ya |
| useControllableState | ✅ ada | ya |
| useDebouncedValue | ✅ ada | ya |
| useClipboard | ✅ ada | ya |
| useElementSize | ✅ ada | ya |
| useFocusTrap | ✅ ada | ya |
| useMediaQuery | ✅ ada | ya |

Adoption = komponen P0 (Dialog, Popover, Combobox, Select, Drawer, dll) refactor
untuk pakai hook ini, hapus inline logic.

### 3. Hooks P1 (6 hooks) — baru

| Hook | Tier A impl | Tier B/C |
|------|-------------|----------|
| useHover | mouseenter/leave | TapTarget rewrite |
| useFocusVisible | :focus-visible polyfill | KeyEvent |
| useEventListener | addEventListener | platform event |
| useOutsideClick | document click | GestureDetector |
| usePrevious | ref cache | ref cache |
| useInterval | setInterval | Timer.periodic |

### 4. Animation Contract — per-tier implementation

`AnimationProps` (AnimationKind, AnimationEasing) sudah di-export. Yang belum:
- Tier A (web): map ke CSS transitions / Web Animations API
- Tier A (RN): map ke `Animated` API
- Tier C (Flutter): map ke `AnimationController` + `Tween`

### 5. Adapt Contract — runtime Popover->Sheet

`AdaptProps` (breakpoint/platform condition) sudah di-export. Yang belum:
runtime resolver yang mengubah `<BcPopover>` menjadi `<BcSheet>` di mobile viewport
(< `adaptBreakpoint`, default `md` = 768px).

## File Ownership Boundary (untuk paralelisasi agent)

**WAJIB dipatuhi agar tidak bentrok dengan agent ADR-014/ADR-015:**

| Path | Owner ADR |
|------|-----------|
| `packages/ui/src/core/hooks/**` | ADR-013 (exclusive) |
| `packages/ui/src/core/providers/**` | ADR-013 (exclusive) |
| `packages/ui/src/core/adapt/**` | ADR-013 (exclusive) |
| `packages/ui/src/core/animation/**` | ADR-013 (exclusive) |
| `packages/ui/src/core/reactive/**` | ADR-013 (exclusive) |
| `packages/ui-flutter/lib/src/core/app_provider.dart` | ADR-013 (exclusive) |
| `packages/ui-flutter/lib/src/core/hooks_*.dart` | ADR-013 (exclusive, file baru) |
| `packages/ui-flutter/lib/src/core/types.dart` | **SHARED** — ADR-013 hanya add reactive types, ADR-015 hanya add chart types. Bagian terpisah, jangan overlap. |
| Komponen `.lite.tsx` (P0/P1) | ADR-013 boleh touch UNTUK ADOPTION (ganti inline logic → hook). ADR-014 boleh touch UNTUK ENGINE WIRING (tambah data-binding). Koordinasi via rebase. |
| `packages/ui/src/index.ts` (barrel) | **SHARED** — semua ADR add export. Lihat "Merge Coordinator" di bawah. |

## Consequences

**Positif:**
- ADR-002 Section 9 fully realized — komponen benar-benar cross-platform, bukan
  cuma "file ada di 4 target"
- Logic duplikat berkurang (inline → shared hooks)
- Flutter parity untuk infrastructure, bukan hanya widget

**Negatif:**
- Refactor 110 komponen berisiko regression — perlu test gating ketat
- Adoption pass butuh waktu (estimasi: 3-5 hari untuk P0, 5-7 hari P1)
- `types.dart` Flutter jadi shared file → butuh disiplin

**Mitigasi:**
- Adoption pass bertahap per group (forms → overlays → navigation → dst)
- Test gating: `npm test` + `flutter test` harus 0 failure sebelum commit
- `types.dart` Flutter: tambah section comment `// === ADR-013 REACTIVE ===`
  dan `// === ADR-015 CHART ===` untuk memisahkan

## Merge Coordinator (paralelisasi)

File berikut di-touch oleh multiple ADR. Owner merge = human/coordinator:

- `packages/ui/src/index.ts` — barrel export
- `packages/ui-flutter/lib/src/core/types.dart` — shared types
- `docs/components-progress.md` — progress tracker
- `docs/adr/README.md` — ADR index

**Aturan:** agent commit ke branch masing-masing, coordinator rebase + resolve
conflict di file-file ini. Agent TIDAK boleh resolve conflict sendiri.

## Open Questions

1. Reactive store port: apakah perlu ADR terpisah, atau cukup section di ADR-013?
   (Saran: cukup section, karena sudah dicakup ADR-004 scope)
2. Adapt default breakpoint: `md` (768px) atau `sm` (640px)? (Saran: `md`, konsisten
   dengan Tailwind default)

## References

- ADR-002 Section 9 — Cross-tier Infrastructure (definisi awal)
- ADR-004 — Reactive Engine
- ADR-006 — Universal Component Strategy (Pure/Adapt)
- ADR-011 — Flutter Adapter Architecture (parity target)
- `packages/ui/src/core/hooks/` — 9 hook files existing
- `packages/ui/src/core/providers/AppProvider.lite.tsx` — AppProvider existing
- `packages/ui/src/core/adapt/types.ts` — AdaptProps existing
- `packages/ui/src/core/animation/types.ts` — AnimationProps existing
