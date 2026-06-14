# Handoff: UI Library — Phase 0 Done, ADR-006 Accepted, Start P0 Foundation

**Tanggal:** 2026-06-14
**Session sebelumnya:** `handoffs/2026-06-14-ui-library-phase0-adr006-and-execution-plan.md` (Phase 0 spike execution)
**Session ini:** Completed Phase 0 spike (6/6 tasks). ADR-006 promoted to **Accepted**.
**Session berikutnya:** **Start P0 #17-22 foundation** (a11y + layout primitives).

---

## Instruksi untuk Sesi Baru

```
Lanjutkan pekerjaan dari handoff file: handoffs/2026-06-14-ui-library-phase0-done-start-p0.md

Konteks singkat: Phase 0 spike ADR-006 COMPLETE. Semua 6 acceptance criteria
terpenuhi. Mitosis RN target verified dengan patch-package fix. ADR-006 status
sekarang "Accepted" (bukan lagi "Proposed pending Phase 0").

TUGAS UTAMA: Mulai P0 foundation components (#17-22 per docs/components-progress.md).
Mulai dengan #20 BcStack family — sebagian sudah ada (BcStack dari Phase 0 PoC),
tinggal tambah BcXStack, BcYStack, BcBox. Lalu BcSlot/BcPortal/BcVisuallyHidden.

Sebelum mulai, baca:
1. docs/adr/006-universal-component-strategy.md (FULL — status sekarang Accepted,
   semua Phase 0 task ada verdict konkret di Section 6 dan 7)
2. docs/components-progress.md (Phase 0 row semua ✅, lanjut ke P0 section)
3. packages/ui/src/components/layout/stack/BcStack.lite.tsx (PATTERN — Pure
   component authoring style)
4. packages/ui/src/components/adapt-poc/popover/ (PATTERN — Adapt component
   authoring style untuk komponen yang beda web vs native)
5. packages/ui/scripts/generate-tokens.mjs (token pipeline source)
6. packages/ui/patches/@builder.io+mitosis+0.13.1.patch (Mitosis RN bug fix)

Verifikasi pre-work sebelum mulai (run, jangan asumsi):
- cd packages/ui && npm run tokens:check   → harus "in sync"
- cd packages/ui && npx mitosis build      → harus 66/66 × 4 targets
- node packages/ui/demos/rn-app/scripts/resolver-check.mjs → harus 3/3 PASS
- node --test packages/ui/src/components/layout/stack/BcStack.test.ts → 16/16 PASS

Urutan eksekusi P0 foundation:
1. Extend BcStack → BcXStack (row alias), BcYStack (column alias), BcBox (generic div)
2. BcSlot (Pure) — used by semua overlay untuk asChild pattern
3. BcPortal (Adapt — web: createPortal, native: Modal)
4. BcVisuallyHidden (Pure) — a11y utility
5. BcBackdrop (Pure) — overlay background
6. BcLoadingOverlay (Adapt — web: overlay, native: ActivityIndicator)
7. BcToastQueue (Pure core, Adapt renderer) — queue logic Pure, rendering Adapt

Setelah foundation selesai (P0 #17-22 done), lanjut ke:
- P0 #1-2: BcForm, BcList (jantung ERP)
- P0 #3: BcDataTable upgrade (complex Adapt)

WAJIB: berpikir kritis, detail, mateng, lengkap dan jujur.
Setiap komponen baru wajib:
  a) Punya .contract.ts (Pure: opsional tapi direkomendasi; Adapt: wajib)
  b) Lulus Mitosis build ke 4 target (Pure) atau di-resolve web+native (Adapt)
  c) Lulus contract test (minimal: smoke source + container element + class
     presence + RN [object Object] regression)
  d) Update docs/components-progress.md (Done = ✅)
  e) Commit dengan format: feat(ui): add BcXxx (Pure|Adapt) per ADR-006

Jangan sentuh working tree kotor (samples/erp, engine, expr-js).
```

---

## Apa yang Sudah Selesai (Session Ini — Phase 0 Spike)

### 7 commits clean (all on `main`)

| Commit | Type | Phase 0 Step |
|---|---|---|
| `8651b73` | docs(adr) | Checkpoint: ADR-006 + tracker + handoff dari sesi sebelumnya |
| `aa0a4de` | feat(ui) | Initial @bitcode/ui baseline (164 files, 65 lite components + tokens + scripts + react demo) |
| `4761539` | chore(adr-001) | **Task 1**: Lynx artifacts drop — ADR-001 §15 marked superseded |
| `6c8b730` | feat(ui) | **Task 2**: Auto-gen token pipeline (`scripts/generate-tokens.mjs`) + `tokens:check` CI guard |
| `09c3897` | feat(ui) | **Task 3**: RN Mitosis target — `stylesType: 'native-wind'` + patch-package fix for `[object Object]` bug |
| `f09c3c7` | feat(ui) | **Tasks 4-6 + Acceptance**: Vite/Metro resolver demos + BcStack Pure PoC + BcPopover Adapt PoC + ADR-006 ACCEPTED |

### Critical Finding: Mitosis 0.13.1 RN Bug + Fix

**Bug:** `NATIVE_WIND_STYLES_PLUGIN` di Mitosis 0.13.1 join `node.bindings.class` (object)
dengan `Array.join(' ')` → menghasilkan literal string `"[object Object]"` di output.
**Affects:** Semua komponen yang pakai dynamic class expression (template literal,
ternary, identifier ref) — ~50% dari 65 existing components.
**Fix:** `patch-package` patch (`patches/@builder.io+mitosis+0.13.1.patch`, 43 lines
added/23 removed). Patch membaca `binding.code` dan emit proper template literal
binding. Postinstall hook auto-applies via `patch-package`.
**Verified:** 66/66 RN outputs clean (0 `[object Object]`). Test di BcStack dan
BcSeparator/BcBadge/BcKpi semua PASS.

### Phase 0 Acceptance Matrix (all green)

| Task | Verification |
|---|---|
| 1. Drop Lynx | grep 0 hits di packages/ui/, ADR-001 §15 banner updated |
| 2. Token auto-gen | `tokens:check` pass, `tokens:generate` idempotent, drift detection works |
| 3. Mitosis RN target | `mitosis build` emits 66/66 × 4 targets; 0 `[object Object]` |
| 4. Build resolver | `resolver-check.mjs` 3/3 PASS (Vite `.web.tsx`, Metro `.native.tsx`, index.ts barrel) |
| 5. PoC Pure | BcStack: 16/16 contract tests PASS across react/vue/reactNative/stencil |
| 6. PoC Adapt | BcPopover: 8/8 contract tests PASS (web/native parity + resolver) |

---

## Apa yang Belum Selesai (Next Session Scope)

### P0 Foundation (#17-22 per docs/components-progress.md)

Sesuai urutan prioritas:

1. **BcStack family extension** (Pure, P0 #20) — BcStack sudah ada, tinggal:
   - `BcXStack.lite.tsx` (row alias untuk BcStack direction="row")
   - `BcYStack.lite.tsx` (column alias)
   - `BcBox.lite.tsx` (generic `<div>` container)
2. **BcSlot** (Pure, P0 #17) — `asChild` pattern, slotted rendering
3. **BcPortal** (Adapt, P0 #18) — web: `react-dom createPortal`, native: `<Modal>`
4. **BcVisuallyHidden** (Pure, P0 #19) — sr-only utility
5. **BcBackdrop** (Pure, P0 #21) — overlay background (fixed inset, click handler)
6. **BcLoadingOverlay** (Adapt, P0 #22) — web: absolute overlay, native: `<ActivityIndicator>`
7. **BcToastQueue** (Pure core + Adapt renderer) — queue logic Pure, rendering Adapt

Total 6-7 komponen foundation. Setelah ini P0 #1-16 (BcForm, BcList, BcDataTable, dll).

### Catatan Tambahan

- **BcPopover di `adapt-poc/`** adalah PoC, bukan production component. P0 #11 (popover
  yang dipakai BitCode) akan rebuild `overlays/popover/BcPopover.lite.tsx` yang existing.
  PoC boleh stay sebagai reference atau di-move.

- **`overlays/popover/BcPopover.lite.tsx` existing** belum migrated ke adapt-poc pattern.
  Itu backlog P0 — sat itu dikerjakan, BcPopover akan jadi Adapt component proper.

- **Tauri integration**: belum di-test di Phase 0 (opsional per ADR-006). Kalau diperlukan
  saat P0 (BcForm/BcList akan dipakai engine BitCode), Tauri test bisa dilakukan ad-hoc
  — Tauri = web bundle jadi pakai Vite resolver sama seperti react-app.

- **NativeWind dark mode**: belum wired. NativeWind preset meng-export `__rawPalette.dark`
  (via `nativewind.preset.cjs`) tapi caller harus wire runtime switch. Backlog post-MVP.

---

## Konteks Penting untuk Sesi Baru

### 1. Pure vs Adapt — Pattern Authoring

**Pure (satu file `.lite.tsx`, Mitosis compile):**
```
src/components/layout/stack/
  BcStack.lite.tsx          ← source of truth
  BcStack.contract.ts       ← optional type contract + helpers
  BcStack.test.ts           ← contract test (run via Node test runner)
```

**Adapt (web + native files, shared contract):**
```
src/components/overlays/popover/  (atau adapt-poc/popover/)
  BcPopover.contract.ts     ← wajib — shared type surface
  BcPopover.web.tsx         ← Vite/Tauri picks via .web.tsx extension priority
  BcPopover.native.tsx      ← Metro picks via .native.tsx default priority
  BcPopover.test.ts         ← cross-impl parity assertions
  index.ts                  ← barrel re-exports (import from './BcPopover')
```

**Key:** Pure pakai `class="..."` (Mitosis canonical). Adapt bebas pakai React
hooks + platform API tapi **wajib** conform ke shared contract.

### 2. Mitosis Build Commands

```bash
cd packages/ui

# Verify all changes are clean
npm run tokens:check            # prebuild hook — fail if drift

# Build (regenerate output/)
npx mitosis build               # 4 targets, ~5-10s

# Post-build fixes (already wired in npm run build)
npm run build                   # mitosis build + fix-react-output + generate-index

# Test PoC components (Node native test runner)
node --test src/components/layout/stack/BcStack.test.ts
node --test src/components/adapt-poc/popover/BcPopover.test.ts
```

### 3. patch-package — Jangan Lupa

Mitosis 0.13.1 RN+NativeWind plugin punya upstream bug. Patch ada di
`packages/ui/patches/@builder.io+mitosis+0.13.1.patch`. **postinstall** hook
auto-applies patch sat `npm install`. Jangan delete patch file.

Kalau `npm install` dari clean state → pastikan patch applied:
```bash
cd packages/ui
npm run postinstall            # explicit apply
# or
npx patch-package              # apply all patches
```

### 4. Token Pipeline — Jangan Edit Artifact

File-file ini **auto-generated**, jangan hand-edit:
- `tokens/tailwind.preset.cjs`
- `tokens/tailwind.theme.css`
- `tokens/nativewind.preset.cjs`
- `tokens/tokens.ts`

Source of truth: `tokens/theme.css`. Edit sana, lalu `npm run tokens:generate`.

CI guard: `npm run tokens:check` (fails sat drift). Sudah wired sebagai
`prebuild` hook.

### 5. Test Pattern (Node native test runner)

```ts
import { test } from 'node:test';
import { strict as assert } from 'node:assert';
// Node 22 supports .ts imports natively (no build step needed).
```

Untuk komponen Pure: assert Mitosis output di `output/{react,vue,react-native,stencil}/`
berisi class yang expected + container element (`<div>` web, `<View>` RN) + no
`[object Object]` regression.

Untuk komponen Adapt: assert contract imported by kedua impls, no platform leakage,
export shape konsisten.

### 6. Commit Convention

```
feat(ui): add BcXxx (Pure|Adapt) per ADR-006

[Description of what the component does, which ADR-006 section it implements,
which P0 number it closes, verification commands run]
```

---

## File yang Berubah/Dibuat di Session Ini

### Updated (existing)
- `docs/adr/006-universal-component-strategy.md` — Status: Proposed → **Accepted**; Section 6 + 7 verdicts
- `docs/adr/001-mitosis-component-library.md` — §15 marked superseded (Task 1)
- `docs/components-progress.md` — Phase 0 row all ✅; verdicts added
- `docs/codebase.md` — Lynx references removed, packages/ui/ tree updated
- `packages/ui/mitosis.config.js` — targets += reactNative, options.reactNative.stylesType=native-wind
- `packages/ui/package.json` — exports ./native, peerDeps react-native+nativewind, postinstall=patch-package, scripts tokens:generate/check
- `packages/ui/tokens/*` — all 4 artifacts now auto-generated; README updated
- `packages/ui/demos/react-app/vite.config.ts` — resolve.extensions prepends .web.tsx

### New (created)
- `packages/ui/scripts/generate-tokens.mjs` — token pipeline
- `packages/ui/tokens/tokens.ts` — typed palette (auto-generated)
- `packages/ui/patches/@builder.io+mitosis+0.13.1.patch` — RN bug fix
- `packages/ui/src/components/layout/stack/BcStack.lite.tsx` — Pure PoC
- `packages/ui/src/components/layout/stack/BcStack.contract.ts`
- `packages/ui/src/components/layout/stack/BcStack.test.ts`
- `packages/ui/src/components/adapt-poc/popover/BcPopover.contract.ts` — Adapt PoC
- `packages/ui/src/components/adapt-poc/popover/BcPopover.web.tsx`
- `packages/ui/src/components/adapt-poc/popover/BcPopover.native.tsx`
- `packages/ui/src/components/adapt-poc/popover/index.ts`
- `packages/ui/src/components/adapt-poc/popover/BcPopover.test.ts`
- `packages/ui/demos/rn-app/scripts/resolver-check.mjs` — bundler resolver smoke test
- `packages/ui/.gitignore` — local (output/, node_modules/, dist/)

### Generated artifacts (gitignored but tracked if regenerated)
- `packages/ui/output/react/` — Mitosis react output (66 components)
- `packages/ui/output/vue/` — Mitosis vue output (66 components)
- `packages/ui/output/react-native/` — Mitosis reactNative output (66 components)
- `packages/ui/output/stencil/` — Mitosis stencil output (66 components)

---

## Verifikasi Akhir (Run Sebelum Mulai Pekerjaan Baru)

```bash
cd packages/ui

# 1. Token pipeline in sync
npm run tokens:check
# Expected: "tokens: ok (in sync with theme.css)"

# 2. Mitosis build clean
npx mitosis build
# Expected: "Mitosis: generation completed." with 66 components × 4 targets

# 3. Resolver smoke
node demos/rn-app/scripts/resolver-check.mjs
# Expected: "resolver smoke: 3 pass, 0 fail, 0 skipped"

# 4. Pure PoC tests
node --test src/components/layout/stack/BcStack.test.ts
# Expected: "# tests 16 / # pass 16 / # fail 0"

# 5. Adapt PoC tests
node --test src/components/adapt-poc/popover/BcPopover.test.ts
# Expected: "# tests 8 / # pass 8 / # fail 0"

# 6. RN output sanity
Select-String -Path "output/react-native/**/*.tsx" -Pattern "\[object Object\]" | Measure-Object | Select-Object -ExpandProperty Count
# Expected: 0
```

Kalau salah satu fail, JANGAN mulai P0. Investigasi dulu — Phase 0 spike ini
sudah verified green. Kalau ada regression, itu sinyal ada yang broken di setup
(petunjuk: cek patch-package applied, cek node_modules fresh, cek theme.css
tidak ter-edit manual).

---

## Total File Count (Session Ini)

| Category | Files | Lines |
|---|---|---|
| New source (lite.tsx, .web.tsx, .native.tsx, contract.ts, test.ts) | 8 | ~600 |
| New scripts (generate-tokens, resolver-check) | 2 | ~560 |
| New docs/ADR updates | 4 | ~250 |
| New patch (patch-package) | 1 | 96 |
| New tokens (auto-gen artifacts) | 4 | ~200 |
| New config (mitosis.config, package.json, vite.config, .gitignore, package.json) | 5 | ~50 |
| **Total new/changed in this session** | **~24 files** | **~1,750 lines** |
| UI library baseline (pre-existing committed at start) | 164 files | 21,278 lines |

**Total `packages/ui/` tracked:** ~190 files, ~23k lines code + tests + tokens + configs.

---

## Next Steps Summary

1. **(next session)** P0 #20 BcStack family extension (XStack, YStack, Box)
2. P0 #17 BcSlot (Pure) — asChild pattern
3. P0 #18 BcPortal (Adapt)
4. P0 #19 BcVisuallyHidden (Pure)
5. P0 #21 BcBackdrop (Pure)
6. P0 #22 BcLoadingOverlay (Adapt)
7. BcToastQueue (Pure core + Adapt renderer)

Setelah 6-7 foundation components selesai, lanjut ke:
- P0 #1-2: BcForm, BcList
- P0 #3: BcDataTable upgrade (Adapt complex)
- P0 #4-16: campuran Pure/Adapt standalone components

Target MVP: 22 P0 components Done per ADR-006 Section 5 DoD.
