# ADR-004: Reactive Engine, State & Expression

## Status: **Proposed**

## Date: 2026-06-06

## Depends on: [ADR-002](./002-cross-platform-architecture.md), [ADR-003](./003-component-roadmap.md)

## Related: `packages/expr-js` (expr-lang port), `packages/go-json` (expr-lang in Go)

---

## Context

### Kebutuhan

BitCode adalah platform low-code/ERP. Form dan view-nya membutuhkan
**reaktivitas** yang tidak ada di component library biasa (MUI/Mantine/shadcn),
tapi standar di form engine & low-code (Filament, Formily, Ant Design Form,
JsonForms, SurveyJS, Retool/Appsmith). Kasus nyata:

- **Cascading data:** provinsi -> kabupaten -> kecamatan -> desa. Pilih provinsi A,
  kabupaten re-query by A. Clear provinsi -> semua dependant ikut clear (rekursif).
- **Conditional validation:** field `tipe_identitas` = NPWP -> `nomor` wajib 16 digit;
  KTP -> 16 digit; SIM -> pola lain.
- **Conditional UI:** pilih nilai tertentu -> section/tombol/field tampil/sembunyi,
  aktif/nonaktif, wajib/opsional.
- **Repeater:** array field dinamis (baris invoice), add/remove/reorder, dengan
  reaksi se-baris (`qty * harga = subtotal`) dan agregasi ke induk (`sum(subtotal)`).
- **Cross-component (opt-in):** filter bar -> DataTable re-query; pilih baris table ->
  form detail terisi; date range -> chart update; selection -> tombol aktif.

### Masalah yang dihindari

- **Tight coupling:** field tidak boleh saling kenal (BcSelect provinsi tak boleh
  tahu BcSelect kabupaten). Itu mimpi buruk maintenance.
- **Global singleton state:** merusak SSR (lihat temuan H2) + merusak standalone
  (dua form/scope di satu halaman saling ganggu).
- **Platform lock-in:** engine tak boleh React-specific (akan dipakai Vue, Svelte,
  Solid, Angular, RN, Lynx via Mitosis; Flutter via track terpisah).
- **`eval()`:** reaksi bisa berasal dari JSON BitCode (semi-untrusted) -> haram.

### Aset yang sudah ada

- `core/types.ts` -> `FetchParams.dependValues` (temuan M2: dirancang untuk
  cascading, belum di-wire). Akan dipakai engine ini.
- `packages/expr-js` -> port TypeScript penuh dari expr-lang (v1.17.8), sedang
  dikerjakan. Menjadi **default expression evaluator** untuk JS targets.
- `packages/go-json` -> expr-lang di Go (backend). Frontend (expr-js) dan backend
  (go-json) berbagi **bahasa ekspresi yang sama** -> parity frontend<->backend.

---

## Decision Summary

1. **Satu primitif: `ReactiveStore`** (logika murni, framework-agnostic). Semua
   reaktivitas (form, repeater, cross-component) dibangun di atasnya.
2. **Reaksi deklaratif** dengan 4 jenis efek: DATA, VALIDATION, UI, VALUE.
3. **Evaluator interface** dengan **default impl = fungsi JS** (zero-dep, unopinionated).
   `expr-js` adalah **opsional plug-in** untuk kasus preview/offline; bukan default.
4. **BitCode JSON feed**: ekspresi dievaluasi **server-side oleh go-json** (sumber
   kebenaran tunggal). Frontend terima hasil. Tidak ada evaluator duplikat.
5. **Scope resolution** `$self` / `$parent` / `$root` untuk repeater & nested.
6. **Konsumen berlapis & opt-in:** Level 0 controlled -> Level 1 BcForm ->
   Level 2 BcReactiveScope -> Level 3 BitCode JSON. Standalone di tiap level.
7. **Per-instance, bukan singleton** (SSR-safe, multi-instance aman).
8. **Headless** (pola TanStack): engine tak meng-render; komponen tetap dumb.

Keputusan rinci di bagian berikut.

---

## Section 1 — ReactiveStore (Core Primitive)

Satu primitif logika murni. Bukan React hooks, bukan DOM. Modul TS biasa yang
bisa dikonsumsi Mitosis (web/RN/Lynx) dan di-port ke Dart (Flutter).

### Tanggung jawab
- Menyimpan **cells** bernama (state addressable): `values`, `meta` per field.
- Mengelola **reactions** (graph dependency).
- Mengevaluasi **expression** via `Evaluator` yang di-inject.
- Meng-emit perubahan ke subscriber (komponen mengikat lewat ini).

### Bentuk (konsep, bukan final)
```
interface ReactiveStore {
  get(path: string): unknown;
  set(path: string, value: unknown): void;     // memicu reaksi
  getMeta(path: string): FieldMeta;             // visible/disabled/required/options/rules/error
  subscribe(listener: (change) => void): () => void;
  unsubscribe(listener: (change) => void): void;  // symmetric removal (BC-024)
  evaluate(expr: string, scope?: Scope): unknown;
  run(): void;                                  // recompute reaksi terdampak
  snapshot(): Record<string, unknown>;
  dispose(): void;                              // teardown (no leak)
}

interface FieldMeta {
  visible: boolean; disabled: boolean; required: boolean;
  options?: unknown[]; rules?: ValidationRule[]; error?: string;
  loading?: boolean;
}
```

### Prinsip terkunci
- **Per-instance**: tiap BcForm/Scope = satu store. Tidak ada singleton global.
  (Pelajaran temuan H2: singleton tidak SSR-safe.)
- **Headless**: store tak meng-render apa pun (pola TanStack). Komponen dumb yang
  membaca `get`/`getMeta` dan memanggil `set`.
- **Sinkron + async**: efek UI/VALUE/VALIDATION sinkron; efek DATA (re-query)
  async dengan state `loading`.
- **Deterministic recompute**: perubahan -> tandai reaksi terdampak -> evaluasi
  topologis -> emit sekali per batch (hindari cascade berulang/`glitch`).

---

## Section 2 — Reaction Model (4 Effects)

Reaksi deklaratif. Field tidak saling kenal; store yang menghubungkan.

```
interface Reaction {
  when: string | string[];   // field pemicu (path)
  target: string;            // field/elemen terdampak
  effect: DataEffect | ValidationEffect | UiEffect | ValueEffect;
  condition?: string;        // expr opsional; reaksi jalan jika true
}
```

### (1) DATA — re-query opsi/sumber
```
{ when:'provinsi', target:'kabupaten',
  data:{ source:{ type:'model', name:'kabupaten',
                  filters:{ provinsi_id:'$provinsi' } },
         clearOnChange:true } }
```
- Saat `provinsi` berubah: re-resolve `source` kabupaten (pakai `dependValues`,
  temuan M2). `clearOnChange` -> kosongkan target + cascade ke dependant-nya.

### (2) VALIDATION — ganti aturan
```
{ when:'tipe_identitas', target:'nomor',
  validation:{ map:{ NPWP:{length:16}, KTP:{length:16}, SIM:{pattern:'...'} } } }
```

### (3) UI — visible/disabled/required
```
{ when:'punya_npwp', target:'section_npwp',
  ui:{ visible:'$punya_npwp == true' } }
```
- Efek UI: `visible`, `disabled`, `required` (ekspresi boolean).

### (4) VALUE — set/clear nilai turunan
```
{ when:['qty','harga'], target:'subtotal',
  value:{ expr:'$qty * $harga' } }
```

### Cascade clear (semantik)
- `clearOnChange`/clear pada source -> store menelusuri graph dependant secara
  rekursif dan mengosongkan semuanya (provinsi -> kabupaten -> kecamatan -> desa).
- Deteksi siklus -> error saat build graph (bukan saat runtime).

---

## Section 3 — Expression & Evaluator

### Reaksi = fungsi JS (default, zero-dep, unopinionated)
Default cara menyatakan reaksi adalah **fungsi JS biasa**:
```
reactions: [
  { when: (s) => s.tipe_identitas, target: 'nomor',
    validation: (s) => ({ length: s.tipe_identitas === 'NPWP' ? 16 : 16 }) },
  { when: (s) => s.punya_npwp, target: 'section_npwp',
    ui: { visible: (s) => s.punya_npwp === true } },
  { when: (s) => [s.qty, s.harga], target: 'subtotal',
    value: (s) => s.qty * s.harga },
]
```
- Zero-dep, unopinionated, semua developer paham.
- `@bitcode/ui` tidak perlu evaluator eksternal.
- Komponen tetap standalone.

### Evaluator interface (swappable, opsional)
```
interface Evaluator {
  evaluate(expr: string, scope: Record<string, unknown>): unknown;
  compile?(expr: string): CompiledExpr;   // opsional: cache AST
}
```
- **Plug-in opsional**, bukan default. Untuk kasus di mana reaksi harus berupa
  **string** (bukan fungsi): JSON-driven, persistable, parity lintas-runtime.
- Implementasi opsional: `expr-js` (expr-lang subset), evaluator AST minimal
  kustom, atau bawaan sendiri. **No-`eval()`** tetap wajib bila pakai string.

### BitCode JSON feed (single source of truth di backend)
- Ekspresi dari JSON BitCode dievaluasi **server-side oleh `go-json`**
  (expr-lang authoritative di Go).
- Frontend terima hasil; tidak perlu evaluator duplikat di UI.
- Keuntungan: satu evaluator (go-json) = tidak ada risiko drift frontend<->backend.
- **expr-js** dipakai di frontend **hanya** untuk:
  - **Live preview** ekspresi di form builder (tanpa round-trip ke server).
  - **Validasi offline** saat mode offline BitCode (Phase 3-5 ADR).
  - Kedua kasus ini **tidak urgent** di Phase A/B; dapat ditunda.

### Batasan
- Reaksi (fungsi JS atau string) = **pure** (no side-effect di evaluator).
  Side-effect (re-query) ditangani store (efek DATA), bukan ekspresi.
- Evaluator (kalau string-based) tidak boleh akses `window`/`fs`/`process`.

---

## Section 4 — Scope Resolution & Repeater

### Path & scope
Ekspresi merujuk field via path. Dalam konteks bersarang (repeater), resolusi:
- `$field` / `$self.field` -> field pada scope saat ini (mis. baris repeater).
- `$parent.field` -> scope induk (form pembungkus).
- `$root.field` -> scope teratas (form akar).
- `$index` -> indeks baris (dalam repeater).

### Repeater (array field)
Padanan Filament Repeater / Formily ArrayField / RHF useFieldArray.
- Tiap baris = **nested ReactiveStore** dengan scope sendiri + tautan ke induk.
- Operasi: `add`, `remove`, `move/reorder`, `duplicate`.
- Reaksi se-baris: `{ when:['qty','harga'], target:'subtotal',
  value:{ expr:'$qty * $harga' } }`.
- Agregasi ke induk: field induk `{ value:{ expr:'sum($self.lines, .subtotal)' } }`.
- Cascade clear & validation tetap berlaku per-scope.

### Bentuk konsumsi (konsep)
```
<BcRepeater name="lines" reactions={[...]}>
  {(row) => (<>
     <BcNumber name="qty" />
     <BcNumber name="harga" />
     <BcNumber name="subtotal" readonly />
  </>)}
</BcRepeater>
```
- `BcRepeater` membuat nested store per baris; field di dalamnya tetap **dumb**.

---

## Section 5 — Consumer Layers (Standalone Guarantee)

Reaktivitas opinionated, tapi **opt-in & additive**. Komponen tetap unopinionated.
Berhenti di level mana pun -> tetap jalan.

```
Level 0 — Controlled (DEFAULT, paling unopinionated)
  <BcInput value={v} onChange={..} />   pure React/Vue, tanpa engine.
  Konsumen yang cuma mau "field bagus" berhenti di sini.

Level 1 — BcForm (opt-in form reactivity)
  Engine di dalam BcForm: cascading, conditional validation, visibility,
  cascade-clear. Dipakai cara React/Vue normal.

Level 2 — BcReactiveScope (opt-in cross-component)
  Field <-> DataTable <-> Chart. Dibungkus eksplisit.
  Tanpa scope -> komponen tetap controlled biasa.

Level 3 — BitCode JSON (platform feed)
  Engine yang sama, di-drive definisi JSON BitCode.
```

### Jaminan prinsip
- **Standalone:** Level 0 tak menyentuh engine/evaluator sama sekali. expr-js hanya
  ter-load bila Level 1+ dipakai (lazy/opt-in import).
- **Unopinionated (komponen):** komponen tak pernah mewajibkan store/Redux/Zustand.
  Engine hidup di container (BcForm/Scope), bukan di komponen.
- **Easy + integrable:** pola headless (TanStack-like). React/Vite/Vue/RN: import
  komponen, opsional bungkus BcForm/Scope. BitCode: feed JSON.
- **Cross-platform:** store + evaluator = logika murni -> ikut ter-port ke Flutter;
  rendering per-target.

### Cross-component (Level 2) — contoh
```
<BcReactiveScope>
  <BcInput name="q" />                       // set scope.q
  <BcDataTable source={{type:'model', name:'contact',
                        filters:{ name:'$q' }}} />  // re-query saat q berubah
</BcReactiveScope>
```
- Bahaya yang dihindari: TIDAK ada global store. Scope = ReactiveStore per-instance.

---

## Section 6 — Phases & Acceptance

### Phase A — Form-scoped (80% kasus)
- ReactiveStore + reaksi DATA/VALIDATION/UI/VALUE (single scope) dengan
  **fungsi JS sebagai evaluator default**.
- BcForm konsumen: cascading (provinsi->...->desa), conditional validation
  (tipe_identitas), visibility/disabled/required.
- **AC:** demo alamat berjenjang + dokumen kondisional jalan; cascade-clear benar.

### Phase B — Repeater (nested scope)
- BcRepeater + nested store; `$self/$parent/$root/$index`.
- Reaksi se-baris + agregasi ke induk.
- **AC:** invoice lines (qty*harga=subtotal, total=sum) jalan; add/remove/reorder
  mempertahankan reaksi.

### Phase C — Cross-component (opt-in)
- BcReactiveScope; field<->table<->chart.
- **AC:** filter->table re-query; selection->tombol enable; date-range->chart.

### Phase D — BitCode JSON feed
- Mapping definisi JSON BitCode -> reaksi engine.
- Ekspresi dievaluasi **server-side oleh go-json**; frontend terima hasil.
- **AC:** satu form ERP nyata ter-drive penuh dari JSON.

### Phase E — Opsional (bukan blocking)
- Integrasi `expr-js` sebagai evaluator plug-in untuk **live preview** di
  form builder dan **validasi offline**.
- **AC:** form builder live-preview tanpa round-trip; mode offline dapat validasi
  client-side.

### Definition of Done (ADR-level)
- [ ] ReactiveStore per-instance, headless, disposable, SSR-safe.
- [ ] 4 efek + cascade-clear + scope resolution teruji.
- [ ] Level 0 tetap standalone (tanpa evaluator eksternal).
- [ ] BitCode JSON dievaluasi go-json server-side (single source of truth).
- [ ] Jalan di target Tier A; contract Flutter (Tier C) terdefinisi.

---

## Consequences

### Positive
- BitCode jadi platform low-code sungguhan (reaktif), bukan sekadar UI kit.
- Satu primitif untuk form/repeater/cross-component -> konsisten, mudah dirawat.
- **Single source of truth** untuk ekspresi = go-json. Tidak ada evaluator
  duplikat, tidak ada risiko drift frontend<->backend.
- `@bitcode/ui` benar-benar unopinionated: reaksi = fungsi JS, zero-dep.

### Negative / Biaya
- Engine + scope resolution + graph recompute = infrastruktur baru.
- Repeater nested scope adalah bagian tersulit (path resolution).
- Round-trip ke server untuk ekspresi BitCode JSON -> butuh koneksi, tapi
  mitigasi: Phase E (expr-js offline) bila offline mode dibutuhkan.

### Risks
- Latensi round-trip ekspresi BitCode JSON -> mitigasi: caching hasil, debounce,
  batch evaluation.
- Over-engineering (meniru Formily/MobX penuh) -> mitigasi: batasi ke dependency
  map + recompute + fungsi JS; bukan observable graph penuh.
- Cross-component (Level 2/3) beririsan dengan rombak engine BitCode -> jaga
  primitif tetap netral agar engine bisa men-drive via JSON nanti.

---

## Subscription Cleanup (BC-024 fix, 2026-06-21)

### Problem

Mitosis `useStore` compiles to React `useState`. When a component stores the
disposer returned by `store.subscribe()` in `useStore` state and reads it in
`onUnMount` (which compiles to `useEffect` cleanup with `deps: []`), the
cleanup closure captures the **initial** value from the first render (null),
not the updated value after `setScopeOff(off)` in the mount effect. The
`unsubscribe` is never invoked, leaking the listener — especially dangerous
for externally-owned stores where `dispose()` is not called on unmount.

### Solution: Dual Storage Pattern

Components that subscribe to a `ReactiveStore` must use **two storage slots**:

1. **`useRef` for the cleanup path** — `useRef.current` is a mutable container
   that is always current regardless of closure capture. The mount effect
   writes the listener (and store, if internally created) to refs; the
   unmount cleanup reads from refs.
2. **`useStore` for the render path** — state that the render function reads
   (e.g. `state.reactiveStore` for `handleSet` to call `store.set()`) stays
   in `useStore` so React re-renders when it changes.

### Symmetric `unsubscribe(listener)` method

`ReactiveStore` now exposes `unsubscribe(listener)` alongside the existing
`subscribe(listener) => () => void` disposer API. This mirrors the DOM
`EventTarget` add/removeEventListener pattern and lets components that store
the **listener function** (stable reference) clean up deterministically,
without relying on the disposer closure being current.

Both APIs are supported. The disposer return-value pattern remains for
backward compatibility and one-shot subscriptions; the symmetric
`unsubscribe(listener)` is preferred for component lifecycle cleanup.

### Affected components

- `BcReactiveScope` — `changeListenerRef` (useRef) + `store.unsubscribe`.
- `BcForm` — `reactiveListenerRef` + `reactiveStoreRef` (both useRef) +
  `store.unsubscribe` + `store.dispose`. (Before the fix, `state.reactiveStore`
  was also stale in cleanup, so `dispose()` never ran — the listener genuinely
  leaked, not just theoretically.)
- `AppProvider` — `disposersRef` (useRef array) holds all disposers;
  cleanup iterates the ref.

### Build reproducibility

Regular `.ts` files (stores, types, utilities) under `src/` are mirrored to
each Mitosis output target's `src/` directory by `scripts/sync-regular-ts.mjs`
(run as a build step). Without this, compiled components that import
`ReactiveStore` would resolve to stale artifacts in `output/` and silently
use old source — the `unsubscribe` method added to `store.ts` was not visible
to `BcForm.tsx` until this sync step was added.

---

## References
- ADR-002 (3-layer, SSR/H2), ADR-003 (BcForm P0, taksonomi).
- `packages/expr-js` (port TS expr-lang, plug-in opsional), `packages/go-json`
  (expr-lang di Go, authoritative).
- Prior art: Filament, Formily, Ant Design Form, JsonForms, SurveyJS, react-hook-form,
  Retool/Appsmith (reactive graph), TanStack (headless pattern).
