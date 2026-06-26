# Session Instruksi: ADR-013 — Cross-Tier Infrastructure Completion

> **Untuk session baru.** Baca file ini + ADR-013 sampai habis sebelum mulai coding.
> Repo root: `D:\NodeJS\lowcode`.

---

## Cara Memulai

```
Kerjakan ADR-013: Cross-Tier Infrastructure Completion.
Baca docs/adr/013-cross-tier-infrastructure-completion.md sampai habis sebelum mulai.

WAJIB berpikir kritis, detail, mateng, lengkap dan jujur.
Setelah selesai per phase, update semua docs terkait lalu commit.

PENTING:
- Ada agent lain yang mungkin sedang mengerjakan ADR-014 (Engine UI Wiring) dan
  ADR-015 (Flutter Chart Strategy) secara paralel. Jangan sembarangan remove atau
  merusak kerjaan agent lain. Konfirmasikan jika butuh hal di luar scope kamu.
- Scope file kamu sudah diatur di ADR-013 section "Scope Boundary untuk Paralelisasi".
  PATUHI boundary tersebut.
```

---

## ADR Referensi

- **ADR-013:** `docs/adr/013-cross-tier-infrastructure-completion.md` (Accepted)
- **Amends:** ADR-002 (Section 9 — Cross-tier Infrastructure), ADR-006
- **Related:** ADR-004 (Reactive Engine), ADR-011 (Flutter Adapter)

---

## Konteks Singkat

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

110/110 komponen enterprise "Done" secara contract+a11y+token+tests, tetapi **tidak
semua benar-benar memakai cross-tier infrastructure**. Banyak komponen masih punya
logic lokal yang seharusnya didelegasikan ke hooks/provider.

---

## Scope File (PATUHI untuk Paralelisasi)

### ✅ Boleh di-touch (milik ADR-013)

```
packages/ui/src/core/hooks/                    # P0 adoption + P1 baru
packages/ui/src/core/providers/                # AppProvider adoption
packages/ui/src/core/adapt/                    # Adapt contract runtime
packages/ui/src/core/animation/                # Animation contract impl
packages/ui/src/core/reactive/                 # Reactive store (jika perlu)
packages/ui-flutter/lib/src/core/app_provider.dart      # Flutter AppProvider port
packages/ui-flutter/lib/src/core/hooks/                 # Flutter hooks port
packages/ui-flutter/lib/src/core/adapt/                 # Flutter adapt port
packages/ui-flutter/lib/src/core/animation/             # Flutter animation port
packages/ui-flutter/lib/src/core/reactive_store.dart    # Flutter reactive port
packages/ui-flutter/test/                                # Tests untuk core baru
docs/adr/013-cross-tier-infrastructure-completion.md    # Update status di ADR
```

### ⚠️ Shared files — butuh koordinasi (JANGAN modify tanpa konfirmasi)

```
packages/ui/src/index.ts                       # Barrel export — merge coordinator
docs/components-progress.md                    # Progress tracker — merge coordinator
docs/adr/README.md                             # ADR index — merge coordinator
packages/ui-flutter/lib/src/core/types.dart    # Shared global types — JANGAN tambah
                                                # chart types di sini (ADR-015 punya
                                                # bc_chart_types.dart sendiri)
```

### ❌ DILARANG touch (milik ADR lain)

```
packages/ui/src/components/charts/             # ADR-015 (chart)
packages/ui-flutter/lib/src/components/bc_chart.dart  # ADR-015 (chart)
packages/components/src/components/charts/     # ADR-015 (chart)
engine/internal/presentation/                  # ADR-014 (engine wiring)
engine/internal/runtime/                       # ADR-014 (engine wiring)
```

**Catatan khusus ADR-013 vs ADR-014 overlap:**
Komponen `BcChatter.lite.tsx`, `BcActivity.lite.tsx`, `BcExport.lite.tsx` bisa di-touch
ADR-014 (engine wiring). ADR-013 **hanya** touch komponen ini di phase terakhir (adoption
pass) untuk refaktor pakai hooks/provider. Jika ADR-014 agent sedang aktif:
- Koordinasikan dulu sebelum refaktor komponen tersebut
- Atau lakukan adoption pass setelah ADR-014 selesai wiring

---

## Task Overview

### Phase 1: Core Infrastructure (selesaiin dulu, sebelum paralel full)

1. **Hooks P1 (6 hooks baru)** — implement di `packages/ui/src/core/hooks/`:
   - `useHover`, `useFocusVisible`, `useEventListener`, `useOnClickOutside`,
     `usePrevious`, `useUnmount`
   - Tiap hook: unit test + contract test

2. **Animation contract implementation** — di `packages/ui/src/core/animation/`:
   - Props sudah exported (`AnimationProps`, `AnimationKind`, `AnimationEasing`)
   - Implement runtime: Tier A (CSS transition), Tier B (RN Animated), Tier C (Flutter)
   - Per-tier fallback strategy

3. **Adapt contract runtime** — di `packages/ui/src/core/adapt/`:
   - Props sudah exported (`AdaptProps`, `AdaptBreakpoint`, dll)
   - Implement runtime: Popover → Sheet di mobile, Dialog → Drawer, dll
   - Breakpoint resolution engine

### Phase 2: Flutter Parity Port

4. Port 5 infrastructure ke `packages/ui-flutter/lib/src/core/`:
   - `app_provider.dart` (InheritedWidget)
   - `hooks/` (Flutter equivalent: StatefulWidget + State pattern)
   - `adapt/` (MediaQuery-based breakpoint resolution)
   - `animation/` (AnimationController wrapper)
   - `reactive_store.dart` (ChangeNotifier port, tanpa Flutter SDK dependency)

### Phase 3: Adoption Pass (setelah Phase 1+2 stabilize)

5. Refactor komponen P0 untuk konsumsi infrastructure:
   - Pakai `useAppContext()` bukan baca `BcSetup` global
   - Pakai hooks P0 (useDisclosure, useControllableState, useFocusTrap, dll)
   - Pakai animation contract props
   - Pakai adapt contract untuk responsive behavior
6. Refactor komponen P1 (gradual, per group)

---

## Aturan Main

1. **`git status` + `git pull` sebelum mulai sesi** — pastikan dapat latest dari agent lain
2. **1 phase = 1 commit** (atau lebih jika perlu), dengan message jelas
3. **Update docs** setelah phase selesai:
   - `docs/components-progress.md` (status cross-tier infra)
   - ADR-013 (update status table di Context section)
4. **Test WAJIB pass** sebelum commit:
   - `cd packages/ui && npm test`
   - `cd packages/ui-flutter && flutter test`
5. **Jangan force push** — branch `main` shared
6. **Bentrok file** → rebase, jangan both modify lalu conflict

---

## Build & Test Commands

```powershell
# UI (Mitosis) tests
cd packages/ui
npm test

# Flutter tests
cd packages/ui-flutter
flutter test

# Contract audit (blocking mode)
cd packages/ui
npm run contract:audit
```

---

## Definition of Done (per ADR-013)

- [ ] 6 hooks P1 baru dengan unit test + contract test
- [ ] Animation contract implemented di Tier A (web) + Tier C (Flutter)
- [ ] Adapt contract runtime implemented (Popover→Sheet, Dialog→Drawer)
- [ ] AppProvider adopted oleh minimal 10 komponen P0
- [ ] Flutter parity: 5 infrastructure ported ke Dart
- [ ] Semua tests pass: `npm test` (ui) + `flutter test` (ui-flutter)
- [ ] ADR-013 status table updated
- [ ] `docs/components-progress.md` updated
- [ ] Commit dengan message: `feat(ui): ADR-013 cross-tier infra — <phase>`

---

## Pertanyaan yang Mungkin Muncul

- **"Adapt runtime pakai library atau manual?"** → Manual, sesuai ADR-006 Pure/Adapt pattern. MediaQuery-based.
- **"Flutter hooks gimana? Dart ga ada hooks."** → Pakai StatefulWidget + State pattern + mixin. Bukan literal hooks.
- **"Reactive store port ke Flutter, tapi tanpa Flutter SDK dep?"** → Ya, pure Dart class dengan ChangeNotifier-like pattern. Lihat ADR-004.
- **"Kalau bentrok sama ADR-014 di BcChatter/BcActivity/BcExport?"** → ADR-014 wiring di engine side (Go templates + API). ADR-013 refaktor di komponen internals. Kalau overlap, ADR-013 lakukan adoption pass di akhir setelah ADR-014 selesai wiring.
