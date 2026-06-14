# ADR-006: Universal Component Strategy (Single Caller API, Platform-Aware Implementation)

## Status: **Proposed (pending Phase 0 spike)**

## Date: 2026-06-14

## Depends on: [ADR-002](./002-cross-platform-architecture.md), [ADR-003](./003-component-roadmap.md)

## Extends ADR-002 Section 1 (Tier A) and ADR-002 Section 9.4 (Adapt contract).

---

## Context

### Visi Stakeholder

> "Sekali bikin satu komponen, jadi di web + desktop + mobile. Saya tidak mau pusing di deploy sebagai apa — langsung jalan."

Visi ini **tidak baru** — sudah implisit di ADR-002 Section 1 yang memasukan React Native ke Tier A. Tapi ADR-002 **belum** mendefinisikan **bagaimana caller** (developer yang pakai library) memperoleh pengalaman "satu import, jalan semua platform" tanpa peduli implementasi internal.

ADR-006 mengisi gap itu: **execution strategy untuk komponen lintas-platform** dari sudut pandang caller.

### Diskusi Stakeholder yang Mengarahkan ADR Ini

1. **Tim besar** (bukan solo dev) → constraint waktu/throughput tidak ketat.
2. **Mobile = prioritas tinggi**, sejajar dengan web (bukan after-thought).
3. **Desktop ada user nyata** via Tauri (`packages/tauri/`) — bukan eksperimen.
4. **BitCode engine menyesuaikan UI**, bukan sebaliknya — tidak ada blocker downstream.
5. **Stencil dipertahankan** sebagai legacy target (engine BitCode existing konsumsi Stencil).
6. **Tamagui/gluestack tetap ditunda** (Tier B per ADR-002) — visi "satu Tailwind" tidak boleh dilanggar.
7. **Lynx eksperimen drop**: `output/lynx/`, `scripts/lynx-generator.mjs`, `demos/lynx-app/` dibersihkan.

### Temuan Audit `packages/ui/` (per 2026-06-14)

| Aspek | State sekarang | Konsekuensi untuk visi |
|---|---|
| `mitosis.config.js` targets | `['stencil', 'react', 'vue']` | **Belum** ada `react-native` target → blocker mobile |
| `package.json` peerDeps | `react>=18, vue>=3.4` | **Belum** ada `react-native` → blocker mobile |
| `tokens/tailwind.preset.cjs` | Pakai `hsl(var(--bc-*))` dari `theme.css` | OK, single source = `theme.css` |
| `tokens/nativewind.preset.cjs` | **Hardcoded HSL palette** di file itu sendiri | **Drift risk** — tidak auto-generate dari `theme.css` |
| Mitosis 0.13.1 RN output | Tidak terverifikasi suport NativeWind `className` | **Spike wajib** sebelum lock strategi |
| Build resolver per-platform (`.web.tsx`/`.native.tsx`) | **Tidak ada** | Blocker untuk komponen Adapt |
| Tauri integration | Repo terpisah `packages/tauri/`, belum konsumsi `@bitcode/ui` | OK secara implisit (Tauri = web bundle) |

---

## Decision Summary

1. **Single Caller API** — caller selalu `import { BcXxx } from '@bitcode/ui'`. Tidak ada import path platform-specific. Bundler resolve otomatis.
2. **Two-Implementation Tier**:
   - **Pure**: 1 source `.lite.tsx` Mitosis → compile ke semua target. Implementasi identik antar platform.
   - **Adapt**: 2 source manual (`.web.tsx` + `.native.tsx`), kontrak prop sama, bundler resolve.
   - **Tidak ada** Tier "platform-only". Semua komponen wajib jalan di semua target deploy (web + RN + Tauri). Yang berbeda: berapa file source yang author tulis. Caller tidak peduli.
3. **Platform Resolver di Build Tooling**:
   - Vite (web/Tauri): `.web.tsx` > `.tsx` > `.ts`
   - Metro (RN): `.native.tsx` > `.tsx` > `.ts` (default Metro behavior)
   - Tauri: pakai output Vite (web), tidak ada target terpisah
4. **Token Pipeline = Single Source of Truth**:
   - `theme.css` (Tailwind 4 `@theme` syntax) = source.
   - `nativewind.preset.cjs` **auto-generated** dari `theme.css` (saat ini hardcoded — Phase 0 task).
   - Tidak boleh ada drift antara web dan RN tokens.
5. **Mitosis Targets Final**:
   - Tier A inti: `react`, `vue`, `react-native` (NativeWind)
   - Legacy preserved: `stencil` (engine BitCode existing — frozen, tidak dapat fitur baru)
   - Total: 4 target Mitosis output
6. **Definition of Done diperluas** — komponen wajib pass kontrak test + smoke compile di **semua 4 target** sebelum ditandai Done. Runtime test mobile (iOS/Android device atau simulator) **tidak** blocking dev, dijalankan di CI nightly atau pre-release.

---

## Section 1 — Single Caller API: Yang Dilihat Developer

Caller **selalu** menulis kode yang sama persis, di semua platform:

```tsx
// di Vite project (web atau Tauri)
import { BcDataTable, BcForm, BcSplitter } from '@bitcode/ui';

<BcSplitter panels={[
  { id: 'left', content: <BcForm ... /> },
  { id: 'right', content: <BcDataTable ... /> },
]} />
```

```tsx
// di React Native (Expo / RN CLI)
import { BcDataTable, BcForm, BcSplitter } from '@bitcode/ui';
// ↑ exact same import path — bundler (Metro) resolve otomatis ke .native.tsx

<BcSplitter panels={[
  { id: 'left', content: <BcForm ... /> },
  { id: 'right', content: <BcDataTable ... /> },
]} />
// ↑ exact same JSX — caller tidak peduli BcSplitter di-render sebagai
//   drag-resize panels (web) atau stacked collapsible (mobile).
```

**Yang berubah otomatis (bundler-resolved, tidak ditulis caller):**
- BcSplitter web: drag handle + continuous resize
- BcSplitter native: stacked panels + collapse toggle (orientation prop di-honor sebisa mungkin)

**Yang TIDAK berubah:** API komponen, prop, event, type signature.

---

## Section 2 — Two-Implementation Tier (Pure vs Adapt)

### 2.1 Pure Components

**Definisi:** Komponen yang implementasi-nya **identik** antar platform — beda hanya output Mitosis (React vs Vue vs RN syntax).

**Contoh kandidat (estimasi):**
- Primitif layout: BcStack, BcXStack, BcYStack, BcBox, BcCenter
- Primitif a11y: BcSlot, BcPortal, BcVisuallyHidden, BcAspectRatio
- Display sederhana: BcBadge, BcAvatar, BcSeparator, BcLabel, BcSkeleton, BcSpinner
- Forms simple: BcInput, BcCheckbox, BcRadio, BcSwitch, BcTextarea
- Feedback: BcEmpty, BcStatistic, BcResult

**Source file:** `BcXxx.lite.tsx` tunggal di `src/components/<group>/<name>/`.

**Mitosis akan emit:**
- `output/react/.../BcXxx.tsx`
- `output/vue/.../BcXxx.vue`
- `output/native/.../BcXxx.tsx` (RN dengan NativeWind classes)
- `output/stencil/.../bc-xxx/bc-xxx.tsx` (frozen — kalau komponen baru, opt-out via flag `mitosisTargets: ['react','vue','native']`)

**Konstrain:** Tidak boleh pakai web-only API (`document`, `window`, `localStorage`, dll) di `.lite.tsx`. Pakai abstraksi `core/hooks/` atau `core/providers/` (lihat ADR-002 Section 9).

### 2.2 Adapt Components

**Definisi:** Komponen yang **harus berbeda** implementasi web vs native untuk memenuhi UX expectation per platform, tapi punya **kontrak prop yang identik**.

**Contoh kandidat (estimasi):**
- Layout interaktif: BcSplitter (drag vs collapse), BcDrawer (slide vs sheet)
- Overlay: BcPopover (floating vs sheet), BcContextMenu (right-click vs long-press actionsheet), BcTooltip (hover vs nonaktif)
- Picker: BcDateField (custom vs native picker), BcTimeField, BcSelect (native vs custom)
- Heavy interaction: BcDataTable (full grid vs MobileList variant)
- Mobile-affine: BcBottomTabs, BcSwipeableRow, BcPullToRefresh — ada implementasi web fallback (Sidebar tabs, hover actions, refresh button)

**Source files:**
```
src/components/overlays/popover/
  BcPopover.contract.ts     ← KONTRAK PROP (TS interface, source of truth)
  BcPopover.web.tsx         ← React/Vue impl (Mitosis)
  BcPopover.native.tsx      ← RN impl (manual, NativeWind)
  BcPopover.test.ts         ← contract test (jalan di semua target)
  index.ts                  ← re-export type + impl, bundler resolve
```

**`index.ts` pattern:**
```ts
// Default re-export delegates to platform-specific via bundler resolution
export type { BcPopoverProps } from './BcPopover.contract';
export { BcPopover } from './BcPopover.web';  // default to .web for type/web
// Metro auto-resolve will use './BcPopover.native' instead of './BcPopover.web'
// when bundling for RN (via package.json `react-native` field or Metro config)
```

**Strategi resolver alternatif (kalau ternyata tidak bekerja):**
- Pakai **conditional exports** di `package.json`:
  ```json
  "exports": {
    "./components/overlays/popover": {
      "react-native": "./src/components/overlays/popover/BcPopover.native.tsx",
      "default": "./src/components/overlays/popover/BcPopover.web.tsx"
    }
  }
  ```
- Atau pakai **`*.native.tsx`** convention yang Metro auto-resolve secara default

**Phase 0 spike** akan validasi mana yang reliable.

### 2.3 Aturan Pembagian Pure vs Adapt

Komponen di-klasifikasi **Adapt** kalau salah satu kondisi terpenuhi:
- Implementasi web pakai DOM API yang tidak ada padanan langsung di RN (mouse events, hover, drag).
- UX expectation per platform fundamental berbeda (right-click vs long-press; calendar grid vs wheel picker).
- Performa beda asimetris yang butuh strategi optimasi platform-spesifik (virtualization).

Selain itu, **default ke Pure**. Ini meminimalkan maintenance burden.

---

## Section 3 — Build Resolver per Platform

### 3.1 Web (Vite)

`vite.config.ts` di consumer atau di internal demo:
```ts
export default {
  resolve: {
    extensions: ['.web.tsx', '.web.ts', '.tsx', '.ts', '.jsx', '.js'],
  },
};
```

Vite auto-pick `.web.tsx` > `.tsx` saat resolve. Adapt component otomatis dapat web impl.

### 3.2 React Native (Metro)

Metro **default** resolve `.native.tsx` > `.tsx`. Tidak butuh konfigurasi tambahan.

Untuk komponen Adapt, file `BcPopover.native.tsx` akan dipilih otomatis saat Metro bundling.

### 3.3 Tauri

Tauri pakai web bundle (Vite output). Resolver = Vite resolver di atas. Tidak ada target terpisah.

### 3.4 Stencil (legacy, frozen)

Stencil tidak punya konsep platform resolver. Komponen Pure di-emit ke Stencil seperti biasa. Komponen Adapt **tidak** di-emit ke Stencil (stencil target di-skip via `mitosisTargets` flag per komponen). Stencil legacy tidak dapat fitur baru — hanya bug fix kritis.

---

## Section 4 — Token Pipeline (Single Source of Truth)

### 4.1 Source: `tokens/theme.css`

Ditulis pakai **Tailwind 4 `@theme` syntax**:
```css
@theme {
  --color-bc-primary: hsl(222 47% 40%);
  --color-bc-background: hsl(0 0% 100%);
  /* ... */
}

@theme dark {
  --color-bc-background: hsl(222 47% 11%);
  /* ... */
}
```

### 4.2 Auto-Generated Outputs (Phase 0 task)

Build script `scripts/generate-tokens.mjs` (BARU, belum ada) akan:
1. Parse `tokens/theme.css`
2. Emit:
   - `tokens/tailwind.preset.cjs` (Tailwind v3 fallback) — sudah ada manual, di-overwrite oleh script
   - `tokens/tailwind.theme.css` (Tailwind v4 `@theme` blok terpisah) — sudah ada manual, di-overwrite
   - `tokens/nativewind.preset.cjs` (RN palette) — **saat ini hardcoded, akan diganti**
   - `tokens/tokens.ts` (object untuk runtime use, kalau diperlukan)
3. CI guard: `npm run tokens:check` — fail kalau output drift dari source.

### 4.3 Konsekuensi

Setelah Phase 0:
- `theme.css` = satu-satunya tempat ubah token.
- `nativewind.preset.cjs` jadi **artifact** (boleh di-gitignore atau di-commit untuk audit). Saya rekomendasi **commit** supaya CI bisa verifikasi tidak drift.
- Web Tailwind 3 user dan v4 user dapat preset/theme yang konsisten dari source yang sama.

---

## Section 5 — Definition of Done (Diperluas dari ADR-002 Section 5.3)

Komponen ditandai **Done** hanya jika:

1. ✅ Kontrak prop di `BcXxx.contract.ts` (TS interface) — source of truth untuk semua impl.
2. ✅ Implementasi:
   - **Pure:** `BcXxx.lite.tsx` — Mitosis emit sukses untuk react, vue, react-native, stencil (kecuali stencil di-skip).
   - **Adapt:** `BcXxx.web.tsx` + `BcXxx.native.tsx` — keduanya conform ke kontrak.
3. ✅ Token-based styling — no hardcoded color/size; dari `theme.css` atau Tailwind utility yang resolve ke token.
4. ✅ A11y verified:
   - **Web:** keyboard nav (Tab, Enter, Esc, Arrow), screen reader (VoiceOver/NVDA), touch target ≥44px.
   - **Mobile:** TalkBack (Android) / VoiceOver (iOS), touch target ≥44pt.
5. ✅ Contract test (`BcXxx.test.ts`) **pass di semua target**:
   - Vitest jalan di Node (logic test).
   - Smoke render di React (Vitest + Testing Library).
   - Smoke render di Vue (Vitest + Vue Test Utils).
   - Smoke render di RN (Vitest + RN Testing Library, tidak butuh device).
6. ✅ E2E smoke (manual atau automated) di setidak-tidaknya:
   - Web (Playwright)
   - Mobile (Maestro / Detox — tidak blocking dev, jalan di CI nightly)
   - Tauri (Playwright + Tauri webdriver — sat available)
7. ✅ Docs entry: contract description + sample usage + live playground (kalau ada).

**Yang TIDAK blocking dev (per kebijakan stakeholder):**
- Test runtime dijalankan **terpisah**, bukan pre-commit.
- Pre-commit hanya: typecheck + lint + Mitosis build sukses.
- Runtime mobile test = pre-release gate, bukan per-commit gate.

---

## Section 6 — Cleanup Tasks (Phase 0)

Sebelum mulai komponen baru di bawah arsitektur ini, **wajib cleanup**:

1. **Drop Lynx artifacts:**
   - Hapus `packages/ui/output/lynx/`
   - Hapus `packages/ui/scripts/lynx-generator.mjs`, `scripts/generate-lynx.mjs`
   - Hapus `packages/ui/demos/lynx-app/`
   - Hapus referensi Lynx di `mitosis.config.js`, `package.json` scripts, `src/index.ts` exports
   - Update ADR-001 Section 15 (Lynx) → tandai sebagai **superseded oleh ADR-006**
2. **Auto-generate token pipeline:**
   - Tulis `scripts/generate-tokens.mjs`
   - Replace `tokens/nativewind.preset.cjs` hardcoded jadi auto-generated
   - Tambah `npm run tokens:check` di CI
3. **Tambah RN target ke Mitosis:**
   - Update `mitosis.config.js` targets: `['stencil', 'react', 'vue', 'react-native']`
   - Add NativeWind plugin/config supaya Mitosis emit `className` (bukan StyleSheet) di RN output
   - Update `package.json` exports: tambah `./native` entry
   - Update `peerDependencies`: tambah `react-native` opsional
4. **Build resolver setup:**
   - Demo Vite `.web.tsx` resolver
   - Demo Metro setup (RN demo app baru)
5. **Sample Adapt component:**
   - Pilih 1 komponen sederhana (mis. BcPopover atau BcDrawer) sebagai PoC Adapt
   - Validasi `web.tsx` + `native.tsx` resolve correctly di kedua bundler
   - Validasi kontrak test jalan di semua target
6. **Sample Pure component:**
   - Pilih 1 komponen sederhana (mis. BcStack atau BcBadge) sebagai PoC Pure
   - Validasi Mitosis emit ke 4 target sukses
   - Validasi RN output pakai NativeWind `className`, bukan StyleSheet inline

**Spike Phase 0 = blocker** untuk lock ADR-006 sebagai Accepted. Kalau spike gagal di salah satu titik (terutama #3 dan #6), ADR-006 di-revisi sebelum ekspansi ke 22 P0.

---

## Section 7 — Acceptance Criteria

### Phase 0 (Spike & Cleanup)
- [ ] Lynx artifacts dihapus, ADR-001 Section 15 ditandai superseded.
- [ ] `scripts/generate-tokens.mjs` jalan dan generate `nativewind.preset.cjs` dari `theme.css`.
- [ ] `mitosis.config.js` targets `['stencil', 'react', 'vue', 'react-native']`.
- [ ] Mitosis emit RN dengan NativeWind classes (verified by inspection of output).
- [ ] PoC Pure component (BcStack atau BcBadge) jalan di Vite demo + RN demo + (opsional) Tauri demo.
- [ ] PoC Adapt component (BcPopover atau BcDrawer) — `.web.tsx` di-resolve Vite, `.native.tsx` di-resolve Metro.
- [ ] Contract test jalan di semua 4 target untuk PoC.

### ADR-006 Accepted (setelah Phase 0)
- [ ] Phase 0 AC semua centang.
- [ ] ADR-006 di-update dengan **temuan konkret** dari spike (replace "estimasi" jadi "verified").
- [ ] `components-progress.md` di-restructure jadi 2-tier (Pure / Adapt).
- [ ] Definition of Done di Section 5 di-validasi pada PoC component.

### MVP (22 P0 selesai)
- [ ] 22 P0 komponen Done sesuai DoD Section 5.
- [ ] Tidak ada komponen yang break di salah satu dari 4 target.
- [ ] Token pipeline auto-sync (CI guard pass).
- [ ] Sample app demo (web + mobile + desktop Tauri) konsumsi semua 22 P0.

---

## Section 8 — Konsekuensi

### Positive
- **Caller experience benar-benar universal** — sesuai visi stakeholder.
- **Token pipeline single-source** mencegah drift web/RN.
- **Library author beban kompleksitas dipindah ke build tooling**, bukan ke caller.
- **Fondasi kuat untuk metadata-driven UI** (ADR-005) — semua komponen punya kontrak prop serializable.
- **Pintu Tier B/C tetap terbuka** karena kontrak komponen sudah platform-neutral.

### Negative / Biaya
- **Maintenance per-Adapt-component 2x lipat** — author tulis web + native impl terpisah.
- **CI lebih berat** — 4 target Mitosis + RN smoke test + Tauri smoke test.
- **Build tooling kompleks** — resolver, token generator, contract test runner.
- **Phase 0 spike sebagai prerequisite** — tidak boleh skip atau ekspansi cepat ke P0 sebelum spike sukses.

### Risks
- **Mitosis 0.13.1 RN output mungkin tidak optimal** — kalau pakai StyleSheet inline (bukan NativeWind className), perlu plugin custom atau upgrade Mitosis.
- **Metro vs Vite resolver behavior** bisa subtle berbeda — perlu integration test di kedua bundler.
- **Conditional exports di package.json** punya gotcha (tooling lama mungkin tidak suport) — backup strategi: pakai `*.native.tsx` filename convention saja.

### Mitigations
- Phase 0 spike validasi semua asumsi sebelum lock.
- ADR-006 status **Proposed pending Phase 0** — boleh direvisi setelah spike.
- Backup strategi tertulis di Section 3 untuk setiap titik risk.

---

## Section 9 — Hubungan dengan ADR Lain

| ADR | Relationship |
|---|---|
| ADR-001 | Section 15 (Lynx) **superseded** oleh ADR-006 Section 6 (drop Lynx). |
| ADR-002 | **Extended.** Section 1 Tier A diperjelas (RN sekarang inti, bukan opsional). Section 9.4 Adapt contract diperluas dengan strategi build resolver. Section 5.3 DoD diperluas jadi Section 5 ADR ini. |
| ADR-003 | **Tidak berubah.** Taxonomy dan prioritization tetap berlaku. Section 5 build sequence diadjust: Phase 0 baru sebelum Phase 1. |
| ADR-004 | **Tidak berubah.** Reactive engine konsumsi kontrak komponen ADR-006. |
| ADR-005 | **Diperkuat.** BcRenderer butuh komponen punya kontrak serializable — ini sudah jadi syarat DoD ADR-006 Section 5. |

---

## References

- ADR-002 Section 1 (target tiering), Section 4 (contract layer), Section 9 (cross-tier infrastructure).
- ADR-003 Section 1 (taxonomy), Section 3 (prioritization P0).
- Mitosis docs: targets, useTarget, plugin API.
- NativeWind v4 docs: preset config, Metro setup.
- Tailwind 4 `@theme` directive documentation.
- React Native Metro resolver: `*.native.tsx` convention.
- Vite resolver: `extensions` config untuk `.web.tsx` priority.
- Tauri 2: webview deployment patterns.
- Conditional exports spec (Node.js, package.json `exports` field).
