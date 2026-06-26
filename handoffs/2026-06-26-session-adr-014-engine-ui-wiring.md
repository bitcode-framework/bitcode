# Session Instruksi: ADR-014 — Engine ↔ UI Feature Wiring

> **Untuk session baru.** Baca file ini + ADR-014 sampai habis sebelum mulai coding.
> Repo root: `D:\NodeJS\lowcode`.

---

## Cara Memulai

```
Kerjakan ADR-014: Engine ↔ UI Feature Wiring.
Baca docs/adr/014-engine-ui-feature-wiring.md sampai habis sebelum mulai.

WAJIB berpikir kritis, detail, mateng, lengkap dan jujur.
Setelah selesai per phase, update semua docs terkait lalu commit.

PENTING:
- Ada agent lain yang mungkin sedang mengerjakan ADR-013 (Cross-Tier Infrastructure)
  dan ADR-015 (Flutter Chart Strategy) secara paralel. Jangan sembarangan remove atau
  merusak kerjaan agent lain. Konfirmasikan jika butuh hal di luar scope kamu.
- Scope file kamu sudah diatur di ADR-014 section "Scope Boundary — Parallel Execution".
  PATUHI boundary tersebut.
- JANGAN mulai Phase 1 sebelum konfirmasi stakeholder soal SPA vs Progressive Enhancement
  (decision point di D1).
```

---

## ADR Referensi

- **ADR-014:** `docs/adr/014-engine-ui-feature-wiring.md` (Proposed)
- **Related:** ADR-005 (Metadata-Driven UI / BcRenderer), ADR-006 (Universal Component Strategy),
  ADR-007 (Enterprise Universal UI API Contract), ADR-008 (Contract-First API Freeze),
  ADR-010 (Web Component Event Bridge)

---

## Konteks Singkat

Engine BitCode saat ini **SSR-based** (Go HTML templates). Komponen UI sudah complete
(110/110 Mitosis + 130 Flutter + 122 Stencil), tapi **banyak komponen belum di-wire**
ke engine features. ADR-014 mendefinisikan kontrak, scope, dan execution plan untuk
wiring 9 high-value engine features ke komponen UI yang sudah ada / baru.

**Prinsip kunci (D1 — Three-tier wiring model):**
```
Tier 1: Engine API (Go) → emits JSON matching component contract (ADR-008)
Tier 2: BcRenderer (ADR-005) → view.json `part` declaration + data-source binding
Tier 3: Component (Mitosis/Flutter) → consumes props, emits events (ADR-010)
```

**⚠️ DECISION POINT sebelum mulai:** ADR-014 D1 masih Proposed. Pilih salah satu:
- **Opsi A: SPA mode** — engine serve JSON, client (React/Vue/RN/Flutter) render full UI
- **Opsi B: Progressive Enhancement** — SSR tetap, komponen di-hydrate setelah page load
Konfirmasi ke stakeholder sebelum mulai Phase 1.

---

## State Saat Ini (2026-06-26)

- Repo: `D:\NodeJS\lowcode`, branch `main`
- `packages/ui`: 110/110 enterprise components Done, 2366 tests, 0 failures
- `packages/ui-flutter`: 130 widgets, 1158 tests, Phase 5c complete
- `packages/components`: 122 Stencil components (legacy, frozen)
- Engine: Go, SSR-based, 570 tests, 0 failures
- BcRenderer (ADR-005): Phase 1-3 complete, `part` declaration + data-source binding ada
- Contract files (ADR-008): 139 `.contract.ts` files, 0 violations

### Komponen yang sudah ada (siap wire):
- `BcChatter` ✅ — `packages/ui/src/components/social/chatter/`
- `BcActivity` ✅ — `packages/ui/src/components/social/activity/`
- `BcExport` ✅ — `packages/ui/src/components/print/export/`
- `BcWizard` ✅ — `packages/ui/src/components/` (16k, complete)
- `BcDataTable` ✅ — untuk pivot table base

### Komponen yang belum ada (perlu buat baru):
- `BcFormBuilder` — drag-drop → valid `view.json`
- `BcWorkflowBuilder` — drag-drop → valid `process.json`
- `BcSchemaBuilder` — drag-drop → valid `model.json`
- `BcPivot` — aggregated matrix view
- `BcNotificationCenter` — navbar notification dropdown

---

## Scope — 9 Features

| # | Feature | UI Component | Backend | Effort | Status |
|---|---------|--------------|---------|--------|--------|
| 1 | Visual Form Builder | New (BcFormBuilder) | Partial (view JSON) | L | Belum |
| 2 | Visual Workflow Builder | New (BcWorkflowBuilder) | Partial (process JSON) | L | Belum |
| 3 | Visual Schema Builder | New (BcSchemaBuilder) | Partial (model JSON) | L | Belum |
| 4 | Multi-Step Wizard | `BcWizard` ✅ | — | M | Belum |
| 5 | Comment & Mention | `BcChatter` ✅ | Belum (no endpoint) | M | Component ada |
| 6 | Activity Feed | `BcActivity` ✅ | Partial (audit_log) | M | Component ada |
| 7 | Export Data | `BcExport` ✅ | Partial (query builder) | M | Component ada |
| 8 | Pivot Table | `BcDataTable` + new `BcPivot` | Belum (no pivot endpoint) | L | Component ada |
| 9 | Email/Notification UI | New (BcNotificationCenter) | Partial (email infra) | L | Belum |

---

## Scope Boundary untuk Paralelisasi (PATUHI STRICTLY)

### File ownership (HANYA ADR-014 agent boleh modify):

```
engine/internal/presentation/          ← ADR-014 ONLY (semua subfolder)
  ├── admin/                           ← admin panel wiring
  ├── api/                             ← REST endpoints baru
  ├── graphql/                         ← GraphQL resolvers baru
  ├── template/                        ← Go HTML templates (jika SSR)
  ├── view/                            ← view renderer, part declarations
  └── websocket/                       ← WebSocket actions baru

engine/internal/domain/                ← ADR-014 ONLY (model baru: comment, notification)
  └── (comment, notification entities)

engine/internal/runtime/               ← ADR-014 ONLY (pivot aggregation, export streaming)

engine/modules/base/                   ← ADR-014 ONLY (comment model auto-create)

packages/ui/src/components/
  ├── form-builder/                    ← ADR-014 ONLY (NEW, BcFormBuilder)
  ├── workflow-builder/                ← ADR-014 ONLY (NEW, BcWorkflowBuilder)
  ├── schema-builder/                  ← ADR-014 ONLY (NEW, BcSchemaBuilder)
  ├── pivot/                           ← ADR-014 ONLY (NEW, BcPivot)
  └── notification/                    ← ADR-014 ONLY (NEW, BcNotificationCenter)

packages/ui-flutter/lib/src/components/
  ├── bc_form_builder.dart             ← ADR-014 ONLY (NEW)
  ├── bc_workflow_builder.dart         ← ADR-014 ONLY (NEW)
  ├── bc_schema_builder.dart           ← ADR-014 ONLY (NEW)
  ├── bc_pivot.dart                    ← ADR-014 ONLY (NEW)
  └── bc_notification_center.dart      ← ADR-014 ONLY (NEW)

engine/docs/features/                  ← ADR-014 ONLY (per-feature docs baru)
```

### Files yang TIDAK boleh di-touch (milik ADR lain):

```
packages/ui/src/core/                  ← ADR-013 territory (hooks, providers, adapt, animation)
packages/ui/src/components/charts/     ← ADR-015 territory
packages/ui-flutter/lib/src/core/      ← ADR-013 territory
packages/ui-flutter/lib/src/components/bc_chart.dart  ← ADR-015 territory
packages/ui-flutter/lib/src/theme/     ← ADR-013 territory
```

### ⚠️ POTENTIAL OVERLAP dengan ADR-013 (resolusi konflik):

Komponen yang SUDAH ada dan akan di-wire oleh ADR-014:
- `BcChatter.lite.tsx`, `BcActivity.lite.tsx`, `BcExport.lite.tsx`

ADR-013 mungkin juga touch file ini untuk **adoption pass** (refaktor pakai hooks/AppProvider).

**ATURAN RESOLUSI:**
1. ADR-014 agent boleh **read** file ini untuk understand contract.
2. ADR-014 agent boleh **modify** hanya untuk wiring (tambah data-source binding, event handlers).
3. ADR-014 agent **TIDAK boleh** refaktor ke hooks/AppProvider (itu ADR-013 job).
4. Kalau ADR-013 belum selesai adoption pass saat ADR-014 mulai wiring, kerjakan wiring
   di **template layer** (`engine/internal/presentation/view/`) dulu, bukan di komponen langsung.
5. Kalau bentrok tidak terhindarkan: **stale agent rebase**, jangan force push.

### Shared files (merge coordinator manages — JANGAN modify langsung, PR/coordinate):

```
packages/ui/src/index.ts               ← barrel exports (coordinate dengan agent lain)
packages/ui-flutter/lib/bc_ui.dart     ← barrel exports (coordinate)
docs/components-progress.md            ← update progress (coordinate)
docs/features.md                       ← update feature status (coordinate)
docs/codebase.md                       ← update file map (coordinate)
docs/adr/README.md                     ← update ADR status (coordinate)
engine/internal/app.go                 ← wire new endpoints (coordinate — shared wiring file)
```

---

## Execution Plan (4 Phases)

### Phase 1: Wiring Patterns & Contract (blocks semua)
- [ ] Pilih SPA vs Progressive Enhancement (decision point D1) — konfirmasi stakeholder
- [ ] Define JSON contract per feature (engine → component), match ADR-008
- [ ] BcRenderer `part` declaration pattern untuk 9 features
- [ ] Reference integration: wire `BcChatter` ke dummy endpoint sebagai contoh pattern
- [ ] Event bridge (ADR-010) wiring pattern untuk client actions

### Phase 2 (Easy Wins — M effort, component sudah ada)
- [ ] `comment` model auto-create + `/api/v1/{module}/comments` CRUD endpoint
- [ ] `BcChatter` wire ke comment endpoint + `@mention` parsing → notification
- [ ] `/api/v1/activity` endpoint (aggregate audit_log + domain events)
- [ ] `BcActivity` wire ke activity endpoint
- [ ] `POST /api/v1/{module}/{model}/export` endpoint (CSV/XLSX/JSON streaming)
- [ ] `BcExport` wire ke export endpoint
- [ ] `wizard` view type di `view.json` + `BcWizard` wire steps from view JSON

### Phase 3 (New Components)
- [ ] `POST /api/v1/{module}/{model}/pivot` endpoint (aggregation)
- [ ] `BcPivot` component (aggregated matrix render, base on BcDataTable)
- [ ] `notification` model + endpoints + WebSocket push
- [ ] `BcNotificationCenter` in navbar

### Phase 4 (Visual Builders — L effort)
- [ ] `BcFormBuilder` drag-drop → valid `view.json` (validate via admin endpoint)
- [ ] `BcWorkflowBuilder` drag-drop → valid `process.json`
- [ ] `BcSchemaBuilder` drag-drop → valid `model.json`
- [ ] Semua 3 builder validate via admin endpoints sebelum save

---

## Testing Requirements

### Engine (Go)
```bash
cd engine
rtk go test ./internal/presentation/... -count=1
rtk go test ./internal/domain/... -count=1
rtk go test ./internal/runtime/... -count=1
```

### UI (Mitosis)
```bash
cd packages/ui
npm test
npm run build
```

### Flutter
```bash
cd packages/ui-flutter
flutter test
```

### Integration (per feature)
- [ ] BcChatter: POST comment → appear in feed, @mention → notification record
- [ ] BcActivity: audit_log entry → appear in feed
- [ ] BcExport: export CSV → file download, correct columns
- [ ] BcPivot: aggregation result correct (sum/count/avg)
- [ ] Visual builders: output JSON valid (pass `bitcode validate`)

---

## Docs Update (MANDATORY setelah setiap phase)

| Yang berubah | Update docs |
|--------------|-------------|
| Endpoint baru | `README.md` (config/CLI), `docs/features.md`, `engine/docs/features/*.md` |
| Model baru (comment, notification) | `docs/codebase.md`, `docs/features.md` |
| Komponen baru | `docs/components-progress.md`, `packages/components/docs/` (jika Stencil), `packages/ui/docs/api/` |
| BcRenderer part baru | `engine/docs/features/metadata-driven-ui.md` (ADR-005 doc) |
| Feature status berubah | `docs/features.md` — update ✅/⚠️/❌ |
| File map berubah | `docs/codebase.md` — add new files/dirs |

---

## Acceptance Criteria (per ADR-014)

- [ ] 9 features wired (3 easy wins + pivot + notification + 3 builders + wizard)
- [ ] JSON contract per feature match ADR-008 (contract-first)
- [ ] BcRenderer `part` declaration untuk semua 9 features
- [ ] Event bridge (ADR-010) berfungsi untuk client actions
- [ ] SPA mode ATAU Progressive Enhancement dipilih & diimplementasi konsisten
- [ ] SSR mode tetap backwards compatible (jika Progressive Enhancement)
- [ ] Engine tests pass (Go)
- [ ] UI tests pass (Mitosis + Flutter)
- [ ] Integration tests pass per feature
- [ ] Semua docs updated

---

## Catatan Penting

1. **Jangan mulai Phase 1 sebelum SPA vs Progressive Enhancement dikonfirmasi.**
   Ini decision point yang menentukan arsitektur semua 9 features.

2. **Visual builders (Phase 4) paling kompleks.** Pertimbangkan pakai library drag-drop
   yang framework-agnostic (mis. dnd-kit untuk React, port ke Flutter via reorderable).
   Tapi cek dulu apakah ada yang sudah ada di codebase.

3. **Comment & Mention butuh `notification` model** — ada dependency antara feature #5
   dan #9. Kerjakan `notification` model dulu di Phase 3, baru mention parsing di Phase 2.

4. **`engine/internal/app.go` adalah shared wiring file.** Kalau agent lain juga tambah
   endpoint, coordinate via merge coordinator. Jangan langsung modify tanuh check `git pull`.

5. **BcChatter/BcActivity/BcExport MUNGKIN di-touch ADR-013 untuk adoption pass.**
   Lihat aturan resolusi di section "POTENTIAL OVERLAP" di atas.

---

## Referensi Tambahan

- `docs/architecture.md` — system design, data flow
- `docs/codebase.md` — file map
- `docs/features.md` — feature inventory (67 features)
- `engine/docs/features/` — per-feature deep docs
- `packages/ui/src/core/renderer/` — BcRenderer implementation (ADR-005)
- `packages/ui/src/core/data-source.ts` — data-source binding
- `packages/ui/src/core/wc-event-bridge.ts` — event bridge (ADR-010)

---

## Quick Start Command

```bash
# 1. Verify state
cd D:\NodeJS\lowcode
git status
git pull

# 2. Baca ADR
cat docs/adr/014-engine-ui-feature-wiring.md

# 3. Cek engine structure
ls engine/internal/presentation/

# 4. Cek BcRenderer (ADR-005)
ls packages/ui/src/core/renderer/

# 5. Mulai Phase 1 (setelah SPA decision confirmed)
```
