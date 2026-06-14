# ADR-005: Metadata-Driven UI (BcRenderer)

## Status: **Proposed**

## Date: 2026-06-06

## Depends on: [ADR-002](./002-cross-platform-architecture.md), [ADR-003](./003-component-roadmap.md), [ADR-004](./004-reactive-engine.md)

## Related: `packages/go-json` (ekspresi authoritative di backend)

---

## Context

### Kebutuhan

Platform low-code seperti BitCode tidak membangun UI dengan menulis kode
komponen manual per form. Ia **menerima definisi JSON** (`views/*.json`,
`models/*.json`) dan **menghasilkan UI yang berjalan**. Ini berbeda dengan
component library biasa (MUI/shadcn) yang dipakai via kode.

Kasus nyata di BitCode:
- `models/contact.json` -> daftar field + label + tipe + validasi + reaksi.
- `views/contact_list.json` -> layout list, kolom, filter, aksi.
- `views/contact_form.json` -> layout form, section, tab, reaksi, tombol.
- Engine membaca JSON dan menyajikan UI yang dapat dipakai pengguna akhir.

### Masalah yang dihindari

- **Renderer di dalam go-json:** go-json adalah bahasa/logika (expr-lang), bukan
  UI framework. JSON-nya tidak bisa render komponen (Go backend).
- **Renderer di dalam BitCode core (monolit):** merusak standalone `@bitcode/ui`.
  Konsumen yang cuma mau schema-driven UI harus ikut menarik seluruh BitCode.
- **Ekspresi dievaluasi di frontend (default):** go-json sudah authoritative di
  backend. Duplikat evaluator -> risiko drift.

### Prinsip yang harus dijaga

- **Standalone `@bitcode/ui`** (level 0-2 di ADR-004): konsumen bisa pakai
  schema-driven UI tanpa BitCode (seperti JsonForms / Formily / SurveyJS).
- **Single source of truth ekspresi:** go-json di backend. Frontend terima hasil.
- **Cross-platform:** schema renderer harus portable (web via Mitosis, Flutter
  via track terpisah).

### Prior art
- **JsonForms** — renderer + schema (standalone, React/Angular/Material).
- **Formily** — schema + reactions (heavy, MobX).
- **SurveyJS** — JSON -> survey form (domain-specific).
- **Filament** — schema-driven forms via Livewire (PHP, server-rendered).

---

## Decision Summary

1. **Renderer berdiri sendiri di `@bitcode/ui`** sebagai konsumen Level 3 di
   atas reactive engine (ADR-004). Konsumen bisa pakai tanpa BitCode.
2. **Skema UI neutral** (tidak domain-specific BitCode) — mirip JsonForms.
   BitCode memetakan definisi domain-nya ke skema UI ini.
3. **Ekspresi dievaluasi server-side oleh go-json**; renderer terima hasil.
4. **Renderer tidak meng-render sendiri**: ia mengkomposisi komponen yang sudah
   ada di `@bitcode/ui` + reactive engine. Tidak ada layer rendering baru.
5. **Ekstensibel via resolver**: mapping tipe-skema -> komponen bisa diganti
   (untuk custom fields, tema, atau target non-web).

Keputusan rinci di bagian berikut.

---

## Section 1 — Schema Model

Skema UI adalah dokumen netral (bukan domain BitCode). Konsumen standalone
punya skema sendiri (JsonForms-like); BitCode **memetakan** definisi domain
(`models/*.json`, `views/*.json`) ke skema ini.

### Bentuk (konsep)
```
interface UiSchema {
  type: 'form' | 'list' | 'detail' | 'section' | 'grid' | 'custom';
  fields?: FieldNode[];        // untuk form/list/detail
  layout?: 'vertical' | 'horizontal' | 'grid';
  columns?: number;
  sections?: UiSchema[];       // nested
  reactions?: Reaction[];      // dari ADR-004 (Level 1)
  actions?: Action[];          // tombol/aksi form-level
}

interface FieldNode {
  name: string;
  label?: string;
  type: string;                 // 'input' | 'select' | 'number' | 'date' | 'file' | 'lookup' | ...
  required?: boolean;
  disabled?: boolean;
  visible?: string;             // ekspresi (dievaluasi server-side)
  options?: DataSource;         // 4-layer (ADR-002)
  validation?: ValidationRule[];
  hint?: string;
  reactions?: Reaction[];
}
```

### Domain-agnostic
- Tipe field = nama komponen di `@bitcode/ui` (BcInput, BcSelect, dll).
- Tidak ada "CRM"/"HRM" di skema. BitCode memetakan entitas CRM ke skema umum.
- Ekstensi `custom` untuk field non-standar.

### Reaksi dalam skema
- `reactions` pakai sintaks reaksi ADR-004 (fungsi JS atau string).
- Ekspresi string (`visible: "status in ['draft','open']"`) dievaluasi
  **server-side oleh go-json** saat load skema atau saat recompute.

---

## Section 2 — Renderer Architecture

Renderer **bukan layer rendering baru**. Ia kompositor: menerjemahkan skema jadi
pohon komponen `@bitcode/ui` yang sudah ada, dibungkus ReactiveStore.

### Bentuk (konsep)
```
<BcRenderer
  schema: UiSchema
  values?: Record<string, unknown>
  onChange?: (values) => void
  onSubmit?: (values) => Promise<void>
  resolver?: FieldResolver      // opsional, ekstensibel
/>
```

### Tugas renderer
1. **Parse skema** -> daftar node.
2. **Resolve field** (type -> komponen) via `FieldResolver`.
3. **Setup ReactiveStore** dengan reaksi dari skema.
4. **Subscribe perubahan** -> re-render node yang terdampak reaksi.
5. **Wire events** (onChange, onSubmit, cascade) ke komponen.

### FieldResolver (ekstensibel)
```
interface FieldResolver {
  resolve(type: string, node: FieldNode): ComponentType | null;
}
```
- **Default resolver:** peta tipe-skema -> komponen `@bitcode/ui`
  (`'input'` -> `BcInput`, `'select'` -> `BcSelect`, `'date'` -> `BcDateTime`).
- **Custom resolver:** konsumen bisa override untuk field non-standar atau
  target lain (mis. Flutter track, tema kustom).
- **Nested rendering:** section/tab di-resolve rekursif.

### Ekspresi & evaluasi
- Saat skema load: kirim ekspresi (visible, disabled, rules) ke **go-json** via
  BitCode endpoint; terima hasil (boolean/aturan) -> hydrate skema.
- Saat nilai berubah (reaksi): evaluasi ulang ekspresi terdampak di server,
  update skema ter-hydrate, re-render.
- Untuk **standalone** (tanpa BitCode): reaksi pakai fungsi JS, tanpa server.

### Output rendering
- Web (React/Vue/Svelte/Solid/Angular) via Mitosis.
- React Native via NativeWind + komponen RN dari Mitosis target.
- Flutter (Tier C, track terpisah): BcRenderer Flutter yang mengonsumsi skema
  sama + field resolver Dart.

---

## Section 3 — BitCode Integration (Mapping, Bukan Rumah Renderer)

### Peran BitCode
- **Menyimpan** definisi domain (`models/*.json`, `views/*.json`) — domain concern.
- **Memetakan** definisi domain ke `UiSchema` netral (transformer).
- **Men-feed** skema ke `BcRenderer` di frontend.
- **Mengevaluasi** ekspresi (visible, disabled, rules) via **go-json** di backend.

### Peran `@bitcode/ui` (BcRenderer)
- **Menerima** `UiSchema` + (opsional) endpoint ekspresi.
- **Merender** komponen dari skema via resolver.
- **Tidak tahu** tentang CRM/HRM/model BitCode — domain-agnostic.
- **Tidak mengevaluasi** ekspresi (kecuali mode offline/standalone).

### Flow BitCode JSON -> UI
```
[BitCode JSON: models/contact.json]
        │ (engine mapping)
        ▼
[UiSchema netral]
        │ (API / embed)
        ▼
[BcRenderer(schema)]
        │ (FieldResolver)
        ▼
[BcInput, BcSelect, ...] ← wrapped ReactiveStore (ADR-004)
        │ (saat render, ekspresi dikirim)
        ▼
[go-json server] ─ evaluates ─► [hasil: boolean/rules] ─► [hydrate skema]
```

### Ekspresi: server-side first
- BitCode endpoint `/api/expr/evaluate` (atau batch) menerima `{ expr, scope }`.
- go-json evaluasi, return `{ result }`.
- Renderer hydrate skema + re-render node terdampak.
- **expr-js plug-in opsional** untuk: form builder live-preview, offline mode
  (lihat ADR-004 Phase E).

### Standalone use-case (tanpa BitCode)
```
<BcRenderer schema={mySchema} resolver={customResolver} />
```
- Skema didefinisikan manual (bukan dari BitCode).
- Reaksi = fungsi JS, tanpa server round-trip.
- Ini sama dengan JsonForms-style UX, zero BitCode dependency.

---

## Section 4 — Phases & Acceptance

### Phase 1 — Form renderer dasar
- `UiSchema` minimal (form + field node) + default `FieldResolver` untuk
  komponen `@bitcode/ui` yang sudah ada.
- `BcRenderer` sebagai kompositor.
- **AC:** skema JSON sederhana -> form berfungsi; field dasar (input/select/
  number/date) jalan; submit handler kerja.

### Phase 2 — Reaktivitas terintegrasi
- Wire `ReactiveStore` (ADR-004) dalam renderer.
- Reaksi dari skema (DATA/VALIDATION/UI/VALUE) aktif.
- **AC:** skema dengan reaksi cascading/conditional validasi/visibility jalan.

### Phase 3 — Layout & nested
- Section, tab, grid, nested schema.
- Repeater (array field) di-render dari skema.
- **AC:** form kompleks (ERP-style) ter-render dari skema tunggal.

### Phase 4 — BitCode JSON mapping
- Transformer BitCode JSON -> `UiSchema`.
- Endpoint ekspresi go-json + hydrate.
- **AC:** satu entitas ERP (contact/lead/invoice) ter-render dari definisi
  BitCode asli, dengan reaksi lintas-field bekerja.

### Phase 5 — Ekstensi & target lain
- Custom resolver (field non-standar, tema).
- Flutter track (BcRenderer Flutter yang mengonsumsi skema sama).
- **AC:** field kustom ter-render via resolver custom; Flutter renderer
  mengonsumsi skema identik.

### Definition of Done (ADR-level)
- [ ] Renderer standalone (tanpa BitCode, pakai skema manual) — AC Phase 1.
- [ ] 4 efek reaksi ADR-004 jalan dalam renderer — AC Phase 2.
- [ ] BitCode JSON -> UiSchema -> UI end-to-end — AC Phase 4.
- [ ] Ekspresi dievaluasi go-json server-side (single source of truth).
- [ ] Custom resolver bekerja; target non-web (Tier C) kontrak terdefinisi.

---

## Consequences

### Positive
- **BitCode jadi low-code sungguhan**: definisi JSON -> UI yang berjalan.
- **`@bitcode/ui` tetap standalone**: renderer bisa dipakai tanpa BitCode,
  mirip JsonForms/Formily — menambah pasar library.
- **Single source of truth ekspresi** (go-json): tidak ada evaluator duplikat,
  tidak ada drift frontend<->backend.
- **Satu skema, banyak target**: web (Mitosis) + Flutter (track terpisah) +
  potensi Lynx/RN — konsisten.
- **Ekstensibel** via resolver: field kustom, tema, atau target baru = plug-in,
  bukan rewrite.

### Negative / Biaya
- Skema netral + mapping BitCode = dua lapisan dokumen (domain + UI). Perlu
  tooling mapping yang teruji.
- Round-trip ekspresi ke server -> butuh latency management (cache, debounce,
  batch).
- Nested schema + repeater = kompleksitas rendering (rekursi, state).

### Risks
- Skema netral kurang ekspresif untuk kasus BitCode tertentu -> mitigasi:
  `custom` type + resolver; jangan memaksakan semua kasus ke skema standar.
- Over-abstraction renderer (Formily-level) -> mitigasi: renderer = kompositor
  sederhana, bukan framework penuh; komponen tetap dumb.
- Flutter track butuh renderer Dart terpisah + resolver Dart -> mitigasi:
  kontrak skema + FieldResolver interface harus konsisten lintas bahasa.

---

## References
- ADR-002 (3-layer, SSR/H2, DataSource), ADR-003 (taksonomi, BcForm P0),
  ADR-004 (ReactiveStore, reaksi, Evaluator interface).
- `packages/go-json` (expr-lang authoritative di Go).
- Prior art: JsonForms, Formily, SurveyJS, Filament (PHP), Retool (low-code).
