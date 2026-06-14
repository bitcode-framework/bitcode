# Architectural Decision Records — BitCode UI

> Index ringkas. Untuk detail keputusan, baca file ADR masing-masing.
> Untuk progress komponen, lihat [`../components-progress.md`](../components-progress.md).

## Status di Repo

Semua ADR di bawah ini adalah **Proposed**. Tidak ada yang Accepted/Rejected
secara formal, tapi praktiknya semua sudah jadi acuan kerja
(`packages/ui/` dibangun mengikuti ADR-002/003).

ADR-001 sebagian **superseded** oleh ADR-002 (lihat catatan kolom).

## Daftar ADR

| # | Judul | Status | Depends on | Dampak utama |
|---|---|---|---|
| [001](./001-mitosis-component-library.md) | Mitosis as Component Compiler | Proposed *(parts superseded by 002, §15 Lynx superseded by 006)* | — | Mitosis dipilih sebagai compiler; theme system shadcn-standard |
| [002](./002-cross-platform-architecture.md) | Cross-Platform Architecture | Proposed *(extended by 006)* | extends 001 | 3-layer (token, contract, target); Tier A/B/C; Section 9 cross-tier infrastructure |
| [003](./003-component-roadmap.md) | Enterprise Component Roadmap & Taxonomy | Proposed | 002 | 9-library research; final taxonomy 14 grup; P0/P1/P2 prioritization |
| [004](./004-reactive-engine.md) | Reactive Engine, State & Expression | Proposed | 002, 003 | Cascading data, conditional UI/validation, repeater, cross-component reactivity |
| [005](./005-metadata-driven-ui.md) | Metadata-Driven UI (BcRenderer) | Proposed | 002, 003, 004 | Engine baca `views/*.json` + `models/*.json` → render UI runtime |
| [006](./006-universal-component-strategy.md) | Universal Component Strategy (Single Caller API, Platform-Aware Implementation) | Proposed *(pending Phase 0 spike)* | extends 002, 003 | 2-tier (Pure/Adapt); platform resolver Vite/Metro; token auto-gen pipeline; Phase 0 wajib sebelum P0 |

## Dependency Graph

```
                    ADR-001 (Mitosis compiler)
                        │
                        ▼  extends + supersedes parts
                    ADR-002 (cross-platform architecture)
                        │  ◄────── extended by ──────┐
              ┌─────────┴─────────┐                  │
              ▼                   ▼              ADR-006
          ADR-003             ADR-004      (universal component
       (taxonomy +         (reactive engine)    strategy: Pure/Adapt,
        components)             │                resolver, token auto-gen)
              │                 │
              └────────┬────────┘
                       ▼
                   ADR-005
              (metadata-driven UI)
```

## Bagian yang Masih Berlaku dari ADR-001

ADR-001 ditulis sebelum ADR-002 memformalkan 3-layer architecture. Yang masih
relevan dari ADR-001:
- Pemilihan Mitosis sebagai compiler.
- Naming convention (`Bc*` PascalCase).
- Inspirasi shadcn theme system (lihat ADR-002 Section 3 untuk implementasi token).

Yang **superseded** oleh ADR-002:
- Detail styling/Tailwind (digantikan token layer ADR-002 Section 3).
- Strategi target framework (digantikan tiering ADR-002 Section 1).

Yang **superseded** oleh ADR-006:
- Section 15 Lynx target — Lynx dropped per ADR-006 Section 6 cleanup.

## Cara Tambah ADR Baru

1. Salin format dari ADR yang ada (Status, Date, Depends on, Context, Decision/Sections, Acceptance Criteria, Consequences, References).
2. Nomor urut berikutnya: **ADR-007**.
3. Update tabel di file ini setelah ADR baru selesai.
4. Jangan ubah ADR yang sudah accepted — bikin ADR baru yang **supersedes** atau **extends**.

## Hubungan dengan Dokumen Lain

| Dokumen | Isi | Update kapan |
|---|---|---|
| `docs/adr/00x-*.md` | Keputusan strategis (this folder) | Setiap keputusan besar baru |
| `docs/adr/README.md` | Index ADR (this file) | Setiap ADR baru / status berubah |
| `docs/components-progress.md` | Tracker per-komponen (Stencil/Mitosis/Contract/A11y) | Setiap komponen done |
| `engine/docs/features/*.md` | Spec per-fitur engine (Go side) | Setiap fitur engine baru |
| `handoffs/YYYY-MM-DD-*.md` | Konteks session AI (lanjutan kerja) | Per session bila konteks panjang |
| `docs/plans/*.md` | Implementation plan per fase | Per fase besar |
| `sprints/notes.md` | Catatan harian/weekly | Bebas |

## Catatan untuk AI/Agent

- ADR adalah **sumber kebenaran arsitektur**. Kalau kode kontradiksi dengan ADR, perbaiki kode atau bikin ADR baru — jangan diam-diam menyimpang.
- Untuk pertanyaan "udah sampai mana?" untuk komponen, baca `docs/components-progress.md`, **bukan** ADR.
- Untuk pertanyaan "kenapa pilih X?", baca ADR yang relevan.
