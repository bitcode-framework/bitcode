# Handoff: UI Perfection — Full Bug Fix Sweep (BC-002/017/024 + Pre-existing Failures)

**Tanggal:** 2026-06-21
**Session:** Lanjutan dari grilling session. User minta "fix semua" — fix BC-017a, BC-024a, BcForm/AppProvider latent bugs, build warning, dan 9 pre-existing test failures.
**Status:** 475 runtime tests pass, 1554 source tests pass (10 failures fixed), 0 build warnings.

---

## Apa yang dikerjakan

Session ini adalah eksekusi penuh dari janji "no tech debt" user. Setiap item yang sebelumnya saya tandai sebagai "belum selesai" difix.

### 1. BC-017a (drag value verification) — FIXED

**Problem:** `fireEvent.pointerMove(handle, { clientX: 700 })` di jsdom deliver `clientX=undefined` (PointerEvent tidak ada di jsdom, fireEvent salah deliver).

**Fix:** Test pakai `new MouseEvent('pointermove', { clientX: 700, bubbles: true })` lalu `handle.dispatchEvent(evt)`. jsdom menghormati MouseEvent constructor properties. Plus mock `getBoundingClientRect` pada root element. Test sekarang asert `onResize(70)` eksak, bukan hanya "called".

### 2. BC-024a (test-timing) — FIXED

**Problem:** `store.set()` bisa fire sebelum `useEffect` subscribe selesai di jsdom.

**Fix:** Test pakai `waitFor` dengan poll probe value (unik per iterasi untuk hindari Object.is guard). Setelah `onChange` pernah terpangil (subscription aktif), clear mock lalu asert `onChange(42)` terjadi pada `store.set('x', 42)`.

### 3. BcForm + AppProvider latent bugs — FIXED

**Problem yang lebih parah dari yang saya kira sebelumnya:**
- `BcForm`: Selain `reactiveOff` (stale closure), `reactiveStore` juga stale closure di cleanup. `setReactiveStore(newStore)` dipangil di mount tapi cleanup `if (reactiveStore) reactiveStore.dispose()` membaca `reactiveStore` dari render pertama (null). **`dispose()` tidak pernah dipangil** → listener BENAR-BENAR LEAK, bukan "aman secara accidental" seperti yang saya klaim sebelumnya.
- `AppProvider`: `state.disposers` di-mutate via `.push()` (bukan replace), jadi secara accidental aman. Tapi pattern fragile — kalau Mitosis compile berubah ke seter pattern, leak muncul.

**Fix (dual storage pattern):**
- `BcForm`: tambah `reactiveListenerRef = useRef<Listener | null>(null)` dan `reactiveStoreRef = useRef<ReactiveStore | null>(null)`. Mount set kedua ref. Cleanup baca dari ref. `state.reactiveStore` tetap di useStore untuk render path (`state.reactiveStore.set(...)`).
- `AppProvider`: pindah `disposers` ke `disposersRef = useRef<Array<() => void>>(mkDisposers())`. Mount push ke `disposersRef`. Cleanup iterate `disposersRef`.

### 4. Build warning: use-debounced-value.ts babel parse — FIXED

**Root cause:** Output `use-debounced-value.ts` di `output/react/` adalah **stale file dari run Mitosis lama** (timestamp 19/06). Source pakai `private readonly onChange?` (TypeScript class field), plugin react-fixes regex booleanRegex mengubah `readonly` → `readOnly` (camelCase) yang bukan valid TypeScript keyword.

Saya sebelumnya fix regex negative lookahead di `mitosis-plugin-react-fixes.mjs` (`(?!\s*[\w-]+\??\s*[?:])` untuk handle optional `?`). Tapi root cause sebenarnya: file `.ts` regular tidak di-copy Mitosis (sumber: `files: 'src/**/**.lite.tsx'`). Stale file hapus, rebuild bersih. Babel parse: 0 failures.

**Perbaikan regex tetap dipertahankan** sebagai defensive guard jika di masa depan ada pattern serupa.

### 5. 10 source test failures — FIXED

1 failure baru (BcForm lifecycle) + 9 pre-existing failures. Semua adalah **stale assertions** yang tidak match pattern setelah refactor BC-022 (can() jadi `state.can()`, ternary event handler jadi inline guard).

- `BcAccordion.test.ts` — pattern `state.openIds = [id]` → `nextIds = [id]` + `state.openIds = nextIds` (BC-016 fix).
- `BcContextMenu.test.ts` — `role="menuitem"` literal → `role={item.separator ? 'separator' : 'menuitem'}` expression (BC-012 fix).
- `BcDataTable.test.ts` (7 tests) — `can('select')` → `state.can('select')` (BC-022 fix), `onKeyDown={col.sortable ? ...}` → `onKeyDown={(e) => { if (col.sortable) ... }}` (BC-022 IIFE fix), `onContextMenu={props.showRowContextMenu}` → `onContextMenu={(e) => { if (props.showRowContextMenu) ... }}` (BC-022), `onDblClick={... ternary}` → inline guard (BC-022), `onScroll={... ternary}` → inline guard (BC-022).
- `BcForm.test.ts` — `state.reactiveStore.dispose()` → `reactiveStoreRef.dispose()` (BC-024 fix), `store.subscribe((snap)` → `const listener: Listener = (snap)` + `store.subscribe(listener)` (BC-024 fix).

Semua test assertions di-update untuk match pattern post-refactor. Tidak ada source code yang diubah untuk test failures ini (source sudah benar, test stale).

---

## Files Changed

| File | Change |
|------|--------|
| `packages/ui/src/core/reactive/store.ts` | Exported `Listener` type. Added `unsubscribe(listener)` method. |
| `packages/ui/src/core/reactive/BcReactiveScope.lite.tsx` | `useRef` for `changeListenerRef`. Use `store.subscribe`/`store.unsubscribe`. |
| `packages/ui/src/core/reactive/reactive-scope-types.ts` | Removed `ScopeOff` type and `mkScopeOff` factory. |
| `packages/ui/src/core/reactive/reactive.rtl.test.tsx` | BC-024a: waitFor poll. BC-024b: unsubscribe-on-unmount test. |
| `packages/ui/src/components/layout/splitter/BcSplitter.lite.tsx` | `draggingRef = useRef<boolean>(false)` mirror for stale closure fix. |
| `packages/ui/src/components/layout/splitter/BcSplitter.test.ts` | Updated source test for `draggingRef` pattern. |
| `packages/ui/src/components/layout/splitter/BcSplitter.rtl.test.tsx` | BC-017a: MouseEvent constructor + getBoundingClientRect mock. |
| `packages/ui/src/components/views/form/BcForm.lite.tsx` | Dual storage: `reactiveListenerRef` + `reactiveStoreRef` useRef. Removed `reactiveOff` from useStore. |
| `packages/ui/src/components/views/form/BcForm.test.ts` | Updated assertions for new patterns (subscribe(listener), reactiveStoreRef.dispose). |
| `packages/ui/src/core/providers/AppProvider.lite.tsx` | `disposersRef` useRef instead of useStore. Removed `disposers` from useStore. |
| `packages/ui/scripts/mitosis-plugin-react-fixes.mjs` | Updated `booleanRegex` negative lookahead to handle optional `?` (defensive). |
| `packages/ui/src/components/layout/accordion/BcAccordion.test.ts` | Updated stale assertion (BC-016 nextIds pattern). |
| `packages/ui/src/components/overlays/contextmenu/BcContextMenu.test.ts` | Updated stale assertion (BC-012 separator ternary). |
| `packages/ui/src/components/data/datatable/BcDataTable.test.ts` | Updated 7 stale assertions (BC-022 state.can() + inline guard). |
| `docs/ui-perfection-buglog.md` | BC-017, BC-024, BcForm/AppProvider notes — semua marked fixed. |

## Test Results

- **475 runtime tests pass**, 0 failures, 27 test files
- **1554 source tests pass**, 0 failures (was 10 failures, all fixed)
- **Build clean**: 0 reformat failures, 179 components generated per target
- axe-core a11y violations: 0 across all tested components

## Konteks Penting

1. **Dual storage pattern** (useStore untuk render + useRef untuk cleanup) adalah pattern standar untuk React state yang perlu diakses di `useEffect` cleanup dengan `deps: []`. Mitosis 0.13.1 compile `useStore` ke `useState`, jadi closure capture selalu render pertama. `useRef.current` selalu terbaru.

2. **`unsubscribe(listener)` method di ReactiveStore** adalah tambahan API baru. Symmetric dengan `subscribe`, mirrors DOM EventTarget pattern. Old pattern (`subscribe` return disposer) tetap di-support untuk backward compatibility.

3. **`MouseEvent` constructor workaround** di BcSplitter test: jsdom tidak punya `PointerEvent` constructor, dan `fireEvent.pointerMove` tidak deliver `clientX/clientY`. Workaround: pakai `MouseEvent` dengan type `'pointermove'` — React tetap memanggil `onPointerMove` handler. **Production tidak perlu workaround** — browser asli punya PointerEvent.

4. **BC-002 tetap deferred ke Phase B** — ini feature gap (partAttr wiring belum dikerjakan), bukan bug. Solusi desain ada di ADR-009. Aman defer.

## Remaining Issues

Tidak ada tech debt tersisa dari scope session ini. BC-002 (partAttr wiring) tetap di Phase B (feature gap, bukan bug).

## Next Steps

- Phase A2: Stencil parity audit
- Phase A3: Demo app browser verification (optional — jsdom tests sekarang verify behavior lengkap, termasuk drag value 70 dan onChange firing)
- Phase A4: 4-layer data source uniformity
- Phase B: partAttr wiring audit across all components (BC-002)
