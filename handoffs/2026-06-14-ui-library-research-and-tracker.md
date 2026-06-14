# Handoff: UI Library Research + ADR Expansion + Tracker Setup

**Tanggal:** 2026-06-14
**Commit terakhir:** `24865e2` — docs: add ADR index + component progress tracker
**Status:** Semua commit berhasil. Working tree masih kotor (bukan scope kita).

---

## Instruksi untuk Sesi Baru

```
Lanjutkan pekerjaan dari handoff file: handoffs/2026-06-14-ui-library-research-and-tracker.md

Konteks singkat: riset 9 UI library sudah selesai, ADR-002 + ADR-003 sudah di-patch
dan committed, ADR index + component tracker sudah dibuat dan committed. Saatnya
mulai implementasi fondasi (token + contract + primitives a11y/layout).

Sebelum mulai, baca file-file berikut:
1. docs/adr/README.md (index 5 ADR)
2. docs/components-progress.md (tracker P0/P1/P2)
3. docs/adr/003-component-roadmap.md Section 3 & Section 5 (priority + build sequence)
4. docs/adr/002-cross-platform-architecture.md Section 9 (cross-tier infrastructure)

Diskusi terakhir: user menanyakan "kalau diselesaikan semua berapa komponen?"
Jawaban: ~110 komponen visual + 1 chart polymorphic + 14 hooks + 2 contracts.
Lihat detail di section "Total Target" di bawah.

User belum memutuskan milestone target (MVP/v1.0/v2.0) — tanyakan dulu sebelum mulai.
User juga belum memutuskan apakah akan mulai implementasi atau lanjut diskusi planning.
```

---

## Apa yang Sudah Dikerjakan (Session Ini)

### 1. Riset 9 UI Library

Dilakukan riset komprehensif terhadap:

| Tier | Library | Coverage |
|---|---|---|
| Tier 1 (lengkap) | Ant Design v6, MUI v9 + MUI X v9, Mantine v9, shadcn/ui, PrimeReact v10 | Enterprise ERP components |
| Tier 2 (a11y) | Radix UI v1, react-aria / react-aria-components | Headless primitives, WAI-ARIA spec |
| Tier 3 (universal) | Tamagui v1, gluestack-ui v2 | Cross-platform (RN + Web) |

**Hasil riset:** 9 library menyumbang total ~110 target komponen + 14 hooks + 2 contracts.
4 library (Radix, react-aria, Tamagui, gluestack) **belum masuk** ADR-003 riset
sebelumnya (cuma 5 library). Setelah patch, semua 9 library sudah diakomodir.

### 2. Patch ADR-002 + ADR-003

**ADR-003 `docs/adr/003-component-roadmap.md`** (274 → 449 baris):
- Section "Kondisi Saat Ini": inventaris updated (127 Stencil .tsx / 128 folder, 65 Mitosis .lite.tsx / 59 unik)
- Section "Riset Pembanding": expanded dari 5 → 9 library (3 tier)
- Section 1 (Taxonomy): + primitives a11y (Slot/Portal/AspectRatio/VisuallyHidden), + layout primitives (Stack family/Box/Center), + input-only fields (DateField/TimeField/NumberField), + ColorPicker sub-primitives, + Backdrop/ToastQueue, + Adapt/Actionsheet, + Sparkline, + Sidebar block, + menu variants
- Section 2C baru: Net-New dari Tier 2/3
- Section 2D baru: Cross-tier infrastructure pointer ke ADR-002 Section 9
- Section 3 (Prioritization): P0 +5 fondasi, P1 +9 items, P2 +7 items
- References: updated

**ADR-002 `docs/adr/002-cross-platform-architecture.md`** (509 → 624 baris):
- Section 9 BARU: Cross-tier Infrastructure
  - 9.1 `core/providers/` — App/ConfigProvider runtime (theme/locale/RTL/messageQueue)
  - 9.2 `core/hooks/` — 8 hooks P0 + 6 hooks P1
  - 9.3 `core/animation/` — AnimationProps contract, reduced-motion, no Framer Motion dep
  - 9.4 `core/adapt/` — Responsive morphing (Popover↔Sheet, Dropdown↔Actionsheet)
  - 9.5 Acceptance Criteria Section 9
- References: updated dengan 9 library + animation/hooks refs

### 3. ADR Index

**`docs/adr/README.md`** (79 baris) — BARU:
- Tabel 5 ADR (nomor, judul, status, depends-on, dampak utama)
- Dependency graph: 001 → 002 → {003, 004} → 005
- Catatan ADR-001 parts superseded by ADR-002
- Aturan tambah ADR baru
- Hubungan dengan dokumen lain

### 4. Component Tracker

**`docs/components-progress.md`** (266 baris) — BARU:
- Ringkasan: P0 0/22 done, P1 ~0/58, P2 0/29, Tier 2/3 net-new 0/~20
- P0 (22 items) dengan kolom: Stencil/Mitosis/Contract/A11y/Done + Notes
- P1 (58 items) dibagi per grup (Forms, Navigation, Feedback/Overlays, Data Display/Layout, Media/Social)
- P2 (29 items)
- Charts section (26 Stencil → 1 polymorphic BcChart)
- Cross-tier infrastructure tracker (5 items)
- Foundation phases tracker (ADR-002 Phase 0–5)
- "Next 5 to Build" rekomendasi urutan eksekusi

### 5. AGENTS.md Updated

`AGENTS.md` baris 13–14: link ke `docs/adr/README.md` + `docs/components-progress.md`.

---

## Total Target Komponen (Jika Semua Selesai)

| Kategori | Jumlah |
|---|---|
| P0 — ERP-Critical | 22 |
| P1 — Forms | 21 |
| P1 — Navigation | 10 |
| P1 — Feedback/Overlays | 8 |
| P1 — Data Display/Layout | 8 |
| P1 — Media & Social | 11 |
| P2 — Polish | 29 |
| BcChart (polymorphic, ~26 type variants) | 1 |
| **Total komponen visual** | **110** |
| Hooks P0 + P1 | 14 |
| Contracts (animation + adapt) | 2 |

**Bandingan:** shadcn ~55, Mantine ~120+60 hooks, antd ~70, PrimeReact ~90+, MUI ~80.

**Estimasi waktu (solo dev, dengan ADR-002 DoD):**
- 110 × 1.5 hari = ~165 hari kerja = ~8 bulan full-time (optimistic)
- Realistis dengan distraksi engine BitCode: 12–18 bulan
- **MVP realistis (P0 22 + hooks + token foundation): 3–4 bulan**

---

## Apa yang Belum Selesai

### Belum dimulai sama sekali
- [ ] Implementasi token layer (ADR-002 Phase 0–2)
- [ ] Implementasi contract (Field/DataView/Overlay)
- [ ] Implementasi `core/providers/AppProvider`
- [ ] Implementasi `core/hooks/` P0 subset (8 hooks)
- [ ] Implementasi `core/animation/` contract
- [ ] Implementasi `core/adapt/` contract
- [ ] Implementasi primitives fondasi: BcSlot, BcPortal, BcVisuallyHidden, BcStack family
- [ ] Migrasi BcForm, BcList dari Stencil (P0 jantung ERP)
- [ ] Upgrade BcDataTable ke enterprise spec (ADR-003 Section 4)
- [ ] Migrasi 70+ komponen Stencil lainnya

### Belum diputuskan (perlu user input)
- [ ] **Milestone target**: user belum jawab apakah mau MVP/v1.0/v2.0 phasing
- [ ] **Prioritas**: user belum konfirmasi "Next 5 to Build" atau mau urutan lain
- [ ] **Token pipeline**: pakai Style Dictionary, atau manual, atau Tailwind @theme saja?
- [ ] **Mitosis output targets prioritas**: React dulu, atau React+Vue bareng, atau semua Tier A?
- [ ] **Chart strategy**: BcChart polymorphic (1 komponen, type prop) atau tetap 26 komponen?

---

## File yang Diubah / Dibuat (Session Ini)

### Committed (`103e75c`)
- `docs/adr/001-mitosis-component-library.md` (596 baris) — baseline, tidak diubah
- `docs/adr/002-cross-platform-architecture.md` (631 baris) — patched (Section 9 baru)
- `docs/adr/003-component-roadmap.md` (455 baris) — patched (9-library research)
- `docs/adr/004-reactive-engine.md` (357 baris) — baseline, tidak diubah
- `docs/adr/005-metadata-driven-ui.md` (285 baris) — baseline, tidak diubah

### Committed (`24865e2`)
- `docs/adr/README.md` (79 baris) — BARU
- `docs/components-progress.md` (266 baris) — BARU
- `AGENTS.md` (+2 baris) — updated dengan link

### Tidak di-commit (working tree kotor, BUKAN scope kita)
Banyak file di `samples/erp/`, `engine/`, `packages/expr-js/`, dll. — itu kerjaan
agent/orang lain. Jangan sentuh.

---

## Konteks Penting untuk Sesi Baru

### 1. User sedang dalam fase **planning**, bukan implementation
User merasa "pusing" dengan konteks yang tersebar. Session ini berhasil:
- Riset 9 library selesai
- ADR di-patch untuk akomodasi 9 library
- Tracker dibuat supaya "udah sampai mana?" punya jawaban konkret
- 2 commit masuk git

Tapi user **belum mulai implementasi**. Kemungkinan next session user mau:
- (a) Lanjut diskusi planning (milestone, prioritas, strategi)
- (b) Mulai implementasi fondasi (token + contract)
- (c) Mulai implementasi komponen P0 tertentu
- **TANYAKAN DULU sebelum mulai.**

### 2. Mitosis 59 komponen ≠ production-ready
Banyak file `.lite.tsx` di `packages/ui/` masih skeleton (logic minimal).
Status ✅ di tracker "Mitosis" hanya berarti file ada — bukan Done.
Definisi Done ada di ADR-002 Section 5.3.

### 3. Stencil 127 komponen masih production
Engine BitCode sekarang masih pakai output Stencil. Migrasi ke Mitosis tidak
menghapus Stencil sampai engine konsumsi output Mitosis.

### 4. ADR-001 status
Statusnya "Proposed" tapi praktiknya **sebagian superseded oleh ADR-002**:
- Mitosis sebagai compiler — tetap berlaku
- Naming convention (Bc*) — tetap berlaku
- Detail styling/Tailwind — superseded oleh ADR-002 Section 3 token layer
- Target framework strategy — superseded oleh ADR-002 Section 1 tiering

### 5. 3 lapis dokumentasi yang sudah dibuat
| Lapis | Dokumen | Pertanyaan |
|---|---|---|
| Strategis | `docs/adr/` (5 ADR + README index) | "kenapa pilih X?" |
| Taktis | `docs/components-progress.md` | "udah sampai mana?" |
| Operasional | TodoWrite + handoffs/ | "hari ini kerjain apa?" |

Jangan campur. Jangan bikin dokumen master raksasa.

### 6. Token pipeline belum diputuskan
ADR-002 Section 3 menyebut token layer (3-layer architecture), tapi implementasi
konkret belum ada. Pilihan: Style Dictionary, manual, atau Tailwind @theme saja.
Ini blocking decision untuk Phase 0.

### 7. Packages/ui struktur sekarang
```
packages/ui/
├── src/
│   ├── components/     59 komponen unik di 9 grup
│   ├── core/           (kosong — belum ada providers/hooks/animation/adapt)
│   ├── index.ts
│   └── themes/
├── tokens/
│   ├── theme.css
│   ├── tailwind.preset.cjs
│   ├── tailwind.theme.css
│   └── nativewind.preset.cjs
├── output/
│   ├── react/          (Mitosis output)
│   ├── vue/
│   └── stencil/
├── mitosis.config.js
├── package.json        (@bitcode/ui, peerDeps React 18+ / Vue 3.4+)
└── scripts/            (dev, build, fix-react-output, generate-index, registry)
```

---

## Next Steps (Rekomendasi Urutan)

Saran saya, **sebelum mulai implementasi**, tanyakan ke user:

1. **"Mau mulai implementasi atau lanjut planning?"**
   - Kalau planning → bahas milestone target (MVP 22 / v1.0 80 / v2.0 110)
   - Kalau implementasi → lanjut ke pertanyaan #2

2. **"Token pipeline mau pakai apa?"** (Style Dictionary / manual / Tailwind @theme)
   - Ini blocking decision untuk semua kerja implementasi

3. **"Target Mitosis output mana dulu?"** (React only / React+Vue / all Tier A)
   - Menentukan scope testing dan CI

4. Kalau semua dijawab, mulai dari **"Next 5 to Build"** di tracker:
   - Token + contract foundation (ADR-002 Phase 0–2)
   - BcSlot + BcPortal + BcVisuallyHidden (a11y fondasi)
   - BcStack family (layout fondasi)
   - BcForm + BcList (jantung ERP)
   - BcDataTable upgrade

---

## Verifikasi

```
git log --oneline -3
24865e2 docs: add ADR index + component progress tracker
103e75c docs(adr): add ADR-001..005 baseline + expand 002/003 to 9-library research
b565896 chore: bump expr-js submodule (patch.ts doc)
```

- 2 commit sukses, 8 file (5 ADR baseline + 3 new docs + 1 updated AGENTS.md)
- Tidak ada file di luar scope yang ter-commit
- Working tree masih kotor (samples/erp dll — bukan scope kita)
- Tidak ada kode `packages/ui/` yang diubah (session ini ADR + docs only)
