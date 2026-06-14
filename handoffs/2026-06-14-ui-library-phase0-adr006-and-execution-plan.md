# Handoff: UI Library Phase 0 — Spike & Cleanup (ADR-006 Execution)

**Tanggal:** 2026-06-14
**Session sebelumnya:** `handoffs/2026-06-14-ui-library-research-and-tracker.md` (research + ADR expansion + tracker setup)
**Session ini:** Diskusi keputusan blocking + draft ADR-006 + update tracker + ADR index
**Session berikutnya:** **EXECUTE Phase 0** (spike + cleanup per ADR-006 Section 6)

---

## Instruksi untuk Sesi Baru

```
Lanjutkan pekerjaan dari handoff file: handoffs/2026-06-14-ui-library-phase0-adr006-and-execution-plan.md

Konteks singkat: ADR-006 (Universal Component Strategy) sudah di-draft dengan status
"Proposed pending Phase 0 spike". Komponen tracker sudah di-restructure ke 2-tier model
(Pure/Adapt). ADR index sudah update dengan ADR-006.

TUGS UTAMA: Eksekusi Phase 0 (6 task di ADR-006 Section 6). Ini spike, bukan produksi.
Kode yang ditulis di Phase 0 = proof-of-concept, boleh dibuang/refactor setelah validasi.

Sebelum mulai, baca:
1. docs/adr/006-universal-component-strategy.md (FULL — terutama Section 6 Phase 0 tasks)
2. docs/components-progress.md (Phase 0 section + P0 tabel dengan kolom Tier)
3. docs/adr/002-cross-platform-architecture.md Section 3 (token pipeline context) dan Section 9.4 (Adapt)
4. packages/ui/mitosis.config.js (state sekarang, target belum include react-native)
5. packages/ui/tokens/nativewind.preset.cjs (HARDCODED — perlu auto-generate)
6. packages/ui/tokens/theme.css (source of truth)

Urutan eksekusi Phase 0:
1. Drop Lynx artifacts (task cepat, risiko rendah)
2. Token auto-gen pipeline (task kedua — script generate-tokens.mjs)
3. Tambah RN target ke Mitosis config + package.json
4. Build resolver demo (Vite web + Metro RN)
5. PoC Pure component (BcStack) — validates Mitosis RN output
6. PoC Adapt component (BcPopover) — validates dual-impl resolver

Setelah Phase 0 selesai, UPDATE ADR-006:
- Replace "estimasi" jadi "verified" di semua section yang relevan
- Promote status dari "Proposed pending Phase 0" jadi "Accepted"
- Commit + update handoff untuk sesi berikutnya (start P0 #17-22)

WAJIB: berpikir kritis, detail, mateng, lengkap dan jujur.
Kalau Phase 0 spike GAGAL di salah satu titik, STOP dan revisi ADR-006.
Jangan push forward kalau asumsi arsitektur tidak validated.

Test infrastructure: DISIAPKAN tapi tidak blocking dev. Runtime test mobile (iOS/Android)
dijalankan terpisah (CI nightly / pre-release), bukan pre-commit.
```

---

## Keputusan Blocking yang SUDAH DIKUNCI (Session Ini)

### 1. Filosofi Library (LOCKED)
> **"Sekali bikin satu komponen, jadi di web + desktop + mobile. Saya tidak mau pusing di deploy sebagai apa — langsung jalan."**

Implementasi konkret via ADR-006:
- **Single Caller API** — caller selalu `import { BcXxx } from '@bitcode/ui'`, tidak perlu tahu platform
- **Platform resolver otomatis** — Vite pick `.web.tsx`, Metro pick `.native.tsx`
- **Library author bayar kompleksitas** (2 impl per Adapt component), caller tidak

### 2. Mitosis Output Targets (LOCKED)
4 target:
- `react` (web utama)
- `vue` (web alternative)
- `react-native` (mobile, NativeWind styling)
- `stencil` (legacy frozen — engine BitCode existing, tidak dapat fitur baru)

Tauri = pakai output React/Vue (web bundle), tidak ada target terpisah.

### 3. Token Pipeline (LOCKED)
- **Source of truth:** `tokens/theme.css` (Tailwind 4 `@theme` syntax)
- **Auto-generate** ke `tailwind.preset.cjs`, `tailwind.theme.css`, `nativewind.preset.cjs`
- **Saat ini:** `nativewind.preset.cjs` HARDCODED → akan diganti auto-gen di Phase 0
- **CI guard:** `npm run tokens:check` untuk detect drift

### 4. Component Tier Model (LOCKED)
2 tier (bukan 3):
- **Pure** (~14 dari 22 P0): 1 source `.lite.tsx` Mitosis → compile ke 4 target
- **Adapt** (~7 dari 22 P0): `.web.tsx` + `.native.tsx` manual, kontrak prop sama
- **Tidak ada** "platform-only" tier — semua komponen wajib jalan di semua target

### 5. Testing Strategy (LOCKED)
- **Pre-commit:** typecheck + lint + Mitosis build sukses (cepat, <30s)
- **Test infrastructure:** disiapkan (Vitest, Playwright, Maestro, Tauri webdriver)
- **Runtime test mobile:** dijalankan **terpisah** (CI nightly / pre-release), tidak blocking dev
- **Definition of Done:** wajib pass kontrak test + smoke di 4 target sebelum komponen ditandai Done

### 6. Stencil & Lynx (LOCKED)
- **Stencil:** dipertahankan sebagai legacy target (frozen, tidak dapat fitur baru)
- **Lynx:** **DROP** — hapus `output/lynx/`, `scripts/lynx-*.mjs`, `demos/lynx-app/`, update ADR-001 §15

### 7. Resource (LOCKED)
- **Tim besar** (bukan solo dev) — constraint throughput tidak ketat
- **Mac tersedia** untuk iOS dev/test
- **Cloud CI tersedia** untuk iOS/Android test
- **BitCode engine menyesuaikan** UI library — tidak ada downstream blocker, sample/erp tidak perlu audit

---

## Apa yang Sudah Dikerjakan (Session Ini)

### 1. Diskusi Keputusan Blocking (5 putaran)
Saya dan user diskusi mendalam selama 5 putaran untuk resolve:
- Milestone target → MVP 22 P0 dengan fondasi yang baik
- Token pipeline → Tailwind 4 `@theme` + NativeWind (shared source)
- Mitosis output → React + Vue + RN + Stencil (legacy)
- Visi "sekali bikin jadi semua" → valid via Single Caller API + Platform-Aware Implementation
- Tier model → 2-tier (Pure/Adapt), bukan 3-tier
- Testing → infra disiapkan, runtime test tidak blocking dev

### 2. ADR-006 Drafted
**File baru:** `docs/adr/006-universal-component-strategy.md` (425 baris)

Status: **Proposed pending Phase 0 spike** — boleh direvisi setelah spike sukses/gagal.

9 section:
1. Single Caller API (yang dilihat developer)
2. Two-Implementation Tier (Pure vs Adapt)
3. Build Resolver per Platform (Vite/Metro/Tauri)
4. Token Pipeline (single source of truth)
5. Definition of Done (diperluas dari ADR-002 §5.3)
6. Phase 0 Cleanup Tasks (6 task prerequisite)
7. Acceptance Criteria (Phase 0, Accepted, MVP)
8. Konsekuensi (positive/negative/risks/mitigations)
9. Hubungan dengan ADR lain

### 3. Component Tracker Updated
**File diupdate:** `docs/components-progress.md`

- Header updated dengan link ke ADR-006
- Ringkasan updated (Phase 0 row added, P0 22 bukan 20)
- Status legend updated dengan kolom Tier (Pure/Adapt)
- **Phase 0 section BARU** di atas "Next 5 to Build" — 6 task prerequisite
- Tabel P0 di-restructure: kolom Tier (P/A) added, estimasi tier per komponen
- Foundation Phases di-restructure (Phase 0 ADR-006 jadi phase aktif)

### 4. ADR Index Updated
**File diupdate:** `docs/adr/README.md`

- Tabel daftar ADR: ADR-006 row added
- Dependency graph: ADR-006 extend ADR-002
- "Bagian superseded dari ADR-001": Lynx §15 superseded oleh ADR-006
- Aturan tambah ADR: nomor urut berikutnya = ADR-007

---

## Apa yang Belum Selesai

### Phase 0 — Spike & Cleanup (PREREQUISITE, belum mulai)

6 task dari ADR-006 Section 6:

1. **Drop Lynx artifacts** ❌
   - Hapus `packages/ui/output/lynx/`
   - Hapus `packages/ui/scripts/lynx-generator.mjs`, `scripts/generate-lynx.mjs`
   - Hapus `packages/ui/demos/lynx-app/`
   - Update `mitosis.config.js`, `package.json` scripts, `src/index.ts` exports
   - Update ADR-001 Section 15 → tandai superseded oleh ADR-006

2. **Auto-generate token pipeline** ❌
   - Tulis `packages/ui/scripts/generate-tokens.mjs` (BARU)
   - Parse `tokens/theme.css` (Tailwind 4 `@theme`)
   - Emit `tailwind.preset.cjs`, `tailwind.theme.css`, `nativewind.preset.cjs`, `tokens.ts`
   - Tambah `npm run tokens:check` untuk CI drift guard

3. **Tambah RN target ke Mitosis** ❌
   - Update `mitosis.config.js`: `targets: ['stencil', 'react', 'vue', 'react-native']`
   - Verify Mitosis 0.13.1 emit `className` (NativeWind compatible), bukan StyleSheet inline
   - Kalau Mitosis emit StyleSheet → tulis post-processing plugin atau upgrade Mitosis
   - Update `package.json` exports: tambah `./native` entry
   - Update `peerDependencies`: tambah `react-native` optional

4. **Build resolver demo** ❌
   - Vite demo: tambah `resolve.extensions` config (`.web.tsx` priority)
   - Metro demo: buat RN demo app baru di `packages/ui/demos/rn-app/` (atau reuse Tauri structure)
   - Verify `.native.tsx` auto-resolve di Metro

5. **PoC Pure component** ❌
   - Pilih BcStack (P0 #20) atau BcBadge (di primitives group)
   - Tulis `BcStack.lite.tsx` (Pure, no platform-specific code)
   - Build ke 4 target (stencil, react, vue, react-native)
   - Verify RN output pakai NativeWind className
   - Deploy ke Vite demo + RN demo + (opsional) Tauri demo
   - Run contract test di semua target

6. **PoC Adapt component** ❌
   - Pilih BcPopover (P0 candidate) atau BcDrawer (P0 #9)
   - Tulis `BcPopover.contract.ts` (props + events interface)
   - Tulis `BcPopover.web.tsx` (React/Vue Mitosis, web popover)
   - Tulis `BcPopover.native.tsx` (RN, sheet/modal based)
   - Tulis `index.ts` dengan conditional exports
   - Verify Vite resolve `.web.tsx`, Metro resolve `.native.tsx`
   - Run contract test di semua target

### Setelah Phase 0 selesai

- Update ADR-006 dengan temuan konkret dari spike
- Promote ADR-006 dari "Proposed pending spike" jadi "Accepted"
- Commit semua Phase 0 work
- Update `components-progress.md` (Phase 0 row → done)
- Tulis handoff baru untuk sesi berikutnya: start P0 #17-22 (fondasi a11y/layout)

### Setelah P0 #17-22 selesai

- P0 #1-2 (BcForm, BcList) — jantung ERP
- P0 #3 (BcDataTable upgrade) — Adapt, complex
- P0 sisanya (#4-16) — campuran Pure/Adapt
- Total MVP: 22 P0 komponen Done per ADR-006 Section 5 DoD

---

## Konteks Penting untuk Sesi Baru

### 1. Token Drift Issue (DITEMUKAN SAAT AUDIT)

**`packages/ui/tokens/nativewind.preset.cjs`** saat ini **HARDCODED** HSL palette:
```js
const light = {
  background: '0 0% 100%',
  foreground: '222 47% 11%',
  // ... 20+ entries hardcoded manual
};
```

Sedangkan `tailwind.preset.cjs` (web) pakai `var(--bc-*)` dari `theme.css`.

**Konsekuensi:** Setiap perubahan token di `theme.css` **tidak otomatis** propagate ke RN.
Ini **bug arsitektur** yang harus diperbaiki di Phase 0 task #2.

Solusi (per ADR-006 Section 4): `scripts/generate-tokens.mjs` parse `theme.css`, emit semua preset termasuk nativewind.

### 2. Mitosis 0.13.1 RN Output — UNVERIFIED

Mitosis config sekarang: `targets: ['stencil', 'react', 'vue']`. Belum ada `react-native`.

**Risk:** Mitosis 0.13.1 mungkin emit StyleSheet inline di RN output (default behavior lama), bukan `className` (yang NativeWind butuhkan).

**Verification needed di Phase 0 task #3:**
- Tambah `react-native` ke targets
- Build BcInput (komponen existing, paling matang)
- Inspect `output/react-native/.../BcInput.tsx`
- Kalau pakai StyleSheet → butuh post-process plugin atau Mitosis upgrade
- Kalau pakai className → sukses, lanjut task #5 (PoC Pure)

Kalau Mitosis tidak bisa emit className untuk RN, ini **blocker arsitektur**. Backup strategi:
- Semua komponen jadi Adapt (manual `.web.tsx` + `.native.tsx`)
- Atau pakai Mitosis plugin custom untuk transformStyleSheetToClassName
- Atau upgrade Mitosis ke versi yang suport NativeWind

### 3. Conditional Exports vs Filename Convention

ADR-006 Section 2.2 menyebut 2 strategi resolver untuk Adapt component:
- **A. Conditional exports di package.json** (`"react-native": "./BcXxx.native.tsx"`)
- **B. Filename convention** (`.native.tsx` auto-resolve oleh Metro)

Strategi B lebih simple (Metro default behavior). Strategi A lebih explicit tapi punya gotcha dengan tooling lama.

**Saran:** Coba B dulu di Phase 0 task #4. Kalau ada masalah, fallback ke A.

### 4. Stencil Frozen — Apa Artinya?

- Komponen baru di Mitosis **tidak** emit ke Stencil (kecuali eksplisit opt-in via `mitosisTargets` flag per komponen)
- Komponen lama yang sudah ada di Stencil (`packages/components/`) tetap utuh
- Bug fix kritis untuk komponen Stencil yang dipakai engine BitCode: lakukan di `packages/components/`, bukan di Mitosis
- Engine BitCode existing tetap konsumsi `packages/components/` (Stencil) sampai migrasi ke `@bitcode/ui` output selesai (post-MVP)

### 5. Tauri Integration State

`packages/tauri/` ada dan terkonfigurasi (lihat `tauri.conf.json` di git status modified). Tapi belum konsumsi `@bitcode/ui`.

**Phase 0 task #4 (opsional):** Validasi Tauri demo bisa import `@bitcode/ui` react output. Tidak perlu setup lengkap — sekadar verify build pipeline.

### 6. Lynx Artifacts — yang harus dihapus

Berdasarkan context history dan glob results:
- `packages/ui/output/lynx/` (folder, generated output)
- `packages/ui/scripts/lynx-generator.mjs` (custom generator)
- `packages/ui/scripts/generate-lynx.mjs` (CLI wrapper)
- `packages/ui/demos/lynx-app/` (folder, demo app dengan Rspeedy)
- Referensi di `mitosis.config.js` (kalau ada, check dulu — saat ini targets tidak include lynx)
- Referensi di `package.json` scripts (check `scripts` field)
- Referensi di `src/index.ts` exports (check kalau ada export Lynx-specific)
- ADR-001 Section 15 → tandai superseded oleh ADR-006 (tidak hapus ADR-001, hanya update status section)

### 7. Mac + Cloud CI Tersedia

User konfirmasi punya Mac untuk iOS dev/test. Cloud CI tersedia untuk iOS/Android test.

**Konsekuensi:** Runtime mobile test (Maestro/Detox) bisa dijalankan locally di Mac. CI cloud untuk automation. Tapi per ADR-006 Section 5, ini **tidak blocking dev** — test dijalankan terpisah, bukan pre-commit.

### 8. Tim Besar — Constraint Throughput Longgar

User konfirmasi "tim besar" (bukan solo dev). Estimasi waktu yang saya kasih di diskusi sebelumnya (12-15 bulan untuk ~55 komponen) **terlalu konservatif** untuk tim.

**Re-estimasi untuk tim besar (3-5 dev):**
- Phase 0: 2-3 minggu (parallelizable sebagian)
- MVP 22 P0: 3-4 bulan (parallelizable per komponen setelah fondasi siap)
- Total ke ~55-60 komponen: 6-9 bulan

Tapi **jangan over-promise**. ADR-006 Phase 0 spike adalah bottleneck — kalau Mitosis RN output ternyata bermasalah, semua timeline geser.

---

## File yang Diubah / Dibuat (Session Ini)

### BARU (untracked)
- `docs/adr/006-universal-component-strategy.md` (425 baris) — ADR-006 draft
- `handoffs/2026-06-14-ui-library-phase0-adr006-and-execution-plan.md` — file ini

### UPDATED (modified, not committed)
- `docs/adr/README.md` — ADR-006 row added, dependency graph updated, ADR-001 §15 superseded note
- `docs/components-progress.md` — Phase 0 section added, P0 tabel restructured (Tier column), Foundation Phases restructured

### TIDAK DI-COMMIT (semua)
User minta: "saat mau coding, saya ingin melakukan di sesi baru". Jadi semua file di atas **uncommitted**. Session baru akan:
1. Review file-file ini
2. Commit sebagai checkpoint sebelum mulai Phase 0
3. Mulai eksekusi Phase 0

---

## Verifikasi

```
Files created/modified this session:
  [NEW]   docs/adr/006-universal-component-strategy.md
  [MOD]   docs/adr/README.md
  [MOD]   docs/components-progress.md
  [NEW]   handoffs/2026-06-14-ui-library-phase0-adr006-and-execution-plan.md

Git status: all untracked/modified (NOT committed — per user request)
Working tree: masih kotor (samples/erp, engine, packages/expr-js — BUKAN scope kita)
```

---

## Next Steps (Rekomendasi Urutan Eksekusi di Sesi Baru)

### Step 1: Review & Commit Checkpoint
1. Baca semua file yang diubah/bikin session ini
2. Kalau setuju → commit dengan message: `docs(adr): add ADR-006 universal component strategy + update tracker & index`
3. Stage HANYA file ini:
   - `docs/adr/006-universal-component-strategy.md`
   - `docs/adr/README.md`
   - `docs/components-progress.md`
   - `handoffs/2026-06-14-ui-library-phase0-adr006-and-execution-plan.md`
4. **Jangan stage** file dari working tree kotor (samples/erp, engine, expr-js)

### Step 2: Phase 0 Task #1 — Drop Lynx ( cepat, risiko rendah)
1. `git rm -r packages/ui/output/lynx/`
2. `git rm packages/ui/scripts/lynx-generator.mjs packages/ui/scripts/generate-lynx.mjs`
3. `git rm -r packages/ui/demos/lynx-app/`
4. Check & clean references di `mitosis.config.js`, `package.json`, `src/index.ts`
5. Update ADR-001 Section 15 → tandai superseded
6. Commit: `chore(ui): drop Lynx artifacts per ADR-006`

### Step 3: Phase 0 Task #2 — Token Auto-Gen Pipeline
1. Tulis `packages/ui/scripts/generate-tokens.mjs`
2. Parser untuk Tailwind 4 `@theme` syntax (bisa pakai `postcss` atau regex simple)
3. Emit: `tailwind.preset.cjs`, `tailwind.theme.css`, `nativewind.preset.cjs`, `tokens.ts`
4. Tambah `npm run tokens:generate` dan `npm run tokens:check` ke `package.json`
5. Run `npm run tokens:generate`, verify output == manual files (kecuali nativewind yang berubah)
6. Commit: `feat(ui): auto-generate token pipeline from theme.css`

### Step 4: Phase 0 Task #3 — Mitosis RN Target
1. Update `mitosis.config.js`: tambah `'react-native'` ke targets
2. Update `package.json` exports & peerDependencies
3. Run `npm run build`, inspect `output/react-native/`
4. **KRITIS:** Cek apakah output pakai `className` (NativeWind compatible) atau `StyleSheet`
5. Kalau StyleSheet → STOP, butuh plugin custom atau Mitosis upgrade
6. Kalau className → lanjut Step 5
7. Commit: `feat(ui): add react-native Mitosis target`

### Step 5: Phase 0 Task #4 — Build Resolver Demo
1. Update Vite demo config: `resolve.extensions: ['.web.tsx', '.web.ts', '.tsx', ...]`
2. Bikin RN demo app baru (`packages/ui/demos/rn-app/`) dengan Metro
3. Verify Metro auto-resolve `.native.tsx`
4. Commit: `chore(ui): add platform resolver demos (Vite + Metro)`

### Step 6: Phase 0 Task #5 — PoC Pure Component
1. Pilih BcStack (P0 #20, layout primitive, simple)
2. Tulis `src/components/primitives/stack/BcStack.lite.tsx`
3. Tulis `BcStack.contract.ts`
4. Build ke 4 target, verify sukses
5. Deploy ke Vite demo + RN demo
6. Tulis contract test, run di React + Vue + RN
7. Commit: `feat(ui): PoC Pure component BcStack (4-target verified)`

### Step 7: Phase 0 Task #6 — PoC Adapt Component
1. Pilih BcPopover atau BcDrawer (P0, overlay)
2. Tulis `BcPopover.contract.ts`
3. Tulis `BcPopover.web.tsx` (React/Vue via Mitosis, atau manual React)
4. Tulis `BcPopover.native.tsx` (RN, Modal-based atau react-native-popover)
5. Tulis `index.ts` dengan conditional exports
6. Verify Vite resolve web, Metro resolve native
7. Contract test di semua target
8. Commit: `feat(ui): PoC Adapt component BcPopover (dual-impl verified)`

### Step 8: Phase 0 Acceptance Review
1. Review semua 6 task — semua centang?
2. Kalau ada yang gagal → revisi ADR-006 section yang relevant
3. Update ADR-006: promote status dari "Proposed pending Phase 0" → "Accepted"
4. Replace "estimasi" di ADR-006 dengan "verified by Phase 0 spike"
5. Update `components-progress.md`: Phase 0 row → done
6. Commit: `docs(adr): promote ADR-006 to Accepted after Phase 0 spike`
7. Tulis handoff baru untuk sesi berikutnya: start P0 #17-22

---

## Catatan untuk AI/Agent Sesi Baru

1. **Phase 0 = SPIKE, bukan produksi.** Kode yang ditulis bokal dibuang/refactor setelah validasi. Jangan over-engineer.

2. **Mitosis RN output = titik kritis.** Kalau task #3 gagal (Mitosis emit StyleSheet bukan className), STOP dan diskusi dengan user sebelum lanjut. Jangan paksa dengan post-process hack kalau root cause-nya Mitosis limitation.

3. **Token drift = bug yang harus diperbaiki.** Saat ini `nativewind.preset.cjs` hardcoded. Task #2 mengganti dengan auto-gen. Jangan commit perubahan token manual setelah auto-gen aktif.

4. **Stencil frozen.** Jangan tambah fitur ke Stencil target. Jangan fix bug di `packages/components/` kecuali kritis dan user explicit request.

5. **Lynx drop = bulk delete.** Tapi cek dulu tidak ada dependensi lain yang import dari `output/lynx/` atau `scripts/lynx-*.mjs`.

6. **Test infra disiapkan tapi tidak blocking.** Pre-commit hanya typecheck + lint + Mitosis build. Runtime test mobile jalan di nightly/pre-release. Jangan setup detox/maestro sebagai pre-commit hook.

7. **Tim besar = parallelizable.** Setelah Phase 0 done, P0 komponen bisa dikerjakan paralel per dev (asal tidak ada shared dependency). Tapi fondasi (#17-22) harus selesai dulu sebelum yang lain.

8. **Jangan sentuh working tree kotor** (samples/erp, engine, packages/expr-js). Itu bukan scope kita.

---

## Total File Count (Session Ini)

| File | Status | Lines |
|---|---|---|
| `docs/adr/006-universal-component-strategy.md` | NEW | 425 |
| `docs/adr/README.md` | MODIFIED | +6 rows, updated graph |
| `docs/components-progress.md` | MODIFIED | +Phase 0 section, restructured P0 tabel |
| `handoffs/2026-06-14-ui-library-phase0-adr006-and-execution-plan.md` | NEW | ~350 |

Semua uncommitted. Session baru akan review dan commit.
