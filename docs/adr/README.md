# Architectural Decision Records — BitCode UI

> Index ringkas. Untuk detail keputusan, baca file ADR masing-masing.
> Untuk progress komponen, lihat [`../components-progress.md`](../components-progress.md).

## Status di Repo

ADR-001 hingga ADR-015 berstatus **Accepted** (lihat tabel di bawah untuk
detail tanggal). ADR-001 sebagian **superseded** oleh ADR-002 (lihat catatan
kolom).

> **Catatan paralelisasi (ADR-013/014/015):** Ketiga ADR dirancang untuk
> dapat dikerjakan paralel oleh agent berbeda. Scope boundary per ADR
> (lihat section "File Ownership" di masing-masing ADR) meminimalkan
> bentrok. ADR-015 (Chart) dan ADR-014 (Engine UI) **zero file overlap**
> dengan apapun. ADR-013 (Cross-tier Infra) punya **low overlap** dengan
> ADR-014 pada beberapa file komponen `.lite.tsx` — resolusi via phase
> sequencing (ADR-013 core stabilize dulu, lalu ADR-14 mulai konsumsi).

## Daftar ADR

| # | Judul | Status | Depends on | Dampak utama |
|---|---|---|---|
| [001](./001-mitosis-component-library.md) | Mitosis as Component Compiler | **Accepted** *(2026-06-21 Phase C1; parts superseded by 002, §15 Lynx superseded by 006)* | — | Mitosis dipilih sebagai compiler; theme system shadcn-standard; 4-target output verified |
| [002](./002-cross-platform-architecture.md) | Cross-Platform Architecture | **Accepted** *(extended by 006, 007, 009, 010)* | extends 001 | 3-layer (token, contract, target); Tier A/B/C; Section 9 cross-tier infrastructure |
| [003](./003-component-roadmap.md) | Enterprise Component Roadmap & Taxonomy | **Accepted** *(2026-06-21 Phase C1)* | 002 | 9-library research; final taxonomy 14 grup; P0/P1/P2 prioritization; Stencil + antd/Prime parity matrices complete |
| [004](./004-reactive-engine.md) | Reactive Engine, State & Expression | **Accepted** *(2026-06-21 Phase C1; Subscription Cleanup section added 2026-06-21)* | 002, 003 | Cascading data, conditional UI/validation, repeater, cross-component reactivity; dual storage pattern (useStore + useRef) for cleanup; `unsubscribe(listener)` method |
| [005](./005-metadata-driven-ui.md) | Metadata-Driven UI (BcRenderer) | **Accepted** *(Phases 1-3 implemented 2026-06-15)* | 002, 003, 004 | Engine baca `views/*.json` + `models/*.json` → render UI runtime |
| [006](./006-universal-component-strategy.md) | Universal Component Strategy (Single Caller API, Platform-Aware Implementation) | Accepted *(Phase 0 spike complete 2026-06-14; auto-mapping primitives documented 2026-06-21)* | extends 002, 003 | 2-tier (Pure/Adapt); platform resolver Vite/Metro; token auto-gen pipeline; Mitosis auto-maps div→View, span→Text, button→Button for RN |
| [007](./007-enterprise-universal-ui-api-contract.md) | Enterprise Universal UI API Contract | **Accepted** *(Phase 2 implementation complete 2026-06-18)* | 002, 003, 004, 005, 006 | Standard event/state/method/prop/a11y contracts that complement reactive and metadata-driven UI integration |
| [008](./008-contract-first-api-freeze.md) | Contract-First API Freeze for `packages/ui` | Accepted *(Phase 2 complete 2026-06-18; partAttr wiring completed 2026-06-21)* | 004, 005, 006, 007 | Freezes feature expansion until Core Integration Surface passes contract matrix, semantic events, adapter mapping, and a11y behavior gates |
| [009](./009-runtime-part-resolution.md) | Runtime Part Resolution | Accepted | 007, 008 | Solves Mitosis spread limitation via pure function helpers (`partClass`/`partAttr`); enables cross-platform slotProps |
| [010](./010-web-component-event-bridge.md) | Web Component Event Bridge | Accepted | 007, 008 | Framework-agnostic semantic events via `dispatchBcEvent`; enables Web Component CustomEvent emission |
| [011](./011-flutter-adapter-architecture.md) | Flutter Adapter Architecture | **Accepted** *(2026-06-22; BcRenderer hybrid port complete 2026-06-24)* | 002, 007, 008 | Tier C Flutter: 75% shared / 25% separate; pre-computed hex const token bridge; 3-layer theme (`BcTheme` + `BcApp`); hybrid reactive store; hybrid BcRenderer port (`BcForm` enhancement + standalone field/data loaders + compositor); FPW adaptive + hand-roll Material-only |
| [012](./012-extract-ui-schema-standalone-package.md) | Extract `@bitcode/ui-schema` as Standalone Package | **Accepted** *(2026-06-22)* | 002, 004, 005, 011 | Portable types + ReactiveStore extracted to zero-dependency package; re-export bridges in `@bitcode/ui` preserve 50+ contract files; unblocks Flutter port (Phase 5a) |
| [013](./013-cross-tier-infrastructure-completion.md) | Cross-tier Infrastructure Completion (ADR-002 Section 9 finish line) | **Proposed** *(2026-06-26)* | 002, 006, 007, 011, 012 | Finish AppProvider adoption, P0/P1 hooks rollout, animation & adapt contract implementation across Tier A (web/RN) + Tier C (Flutter); parallel-safe with ADR-014/015 via file-ownership boundary |
| [014](./014-engine-ui-feature-wiring.md) | Engine UI Feature Wiring (SSR → component-backed render) | **Proposed** *(2026-06-26)* | 005, 007, 010, 011 | Wire 9 ERP UI features (visual builders, chatter, activity, export, pivot, email/notification) into BitCode engine; extends ADR-005 BcRenderer with component-backed SSR render mode; parallel-safe with ADR-013/015 |
| [015](./015-flutter-echarts-chart-strategy.md) | Flutter ECharts Chart Strategy (deprecate CustomPaint) | **Accepted** *(2026-06-26)* | 002, 011 | Replace 912-line hobby-grade CustomPaint with `flutter_echarts`; parity 26/26 chart types via shared ECharts option contract; CustomPaint removed entirely (no deprecation period) |
| [016](./016-demo-app-qa-bug-audit.md) | Demo App QA Bug Audit & Improvement Plan | **Accepted** *(2026-06-29 all 4 phases complete; 5 architectural decisions: useRef migration + lint guardrail, post-build portal AST, @source Tailwind consumption contract, contract-alignment recipe, echarts dependency-injection loader; 28 stub components ported across 5 clusters)* | 001, 002, 006, 007, 011 | Full browser QA of React + Flutter demos: 2 critical bugs (Mitosis ref compiler, overlay portal), 6 high, 10 medium, 28 stubs ported; charts fixed via DI loader pattern; 3 pre-existing runtime bugs fixed (getAnimationClasses import, BcTabs API alignment, BcSplitter dragging init) |
| [017](./017-block-catalog-architecture.md) | Block Catalog Architecture (`@bitcode/ui-blocks`) | **Proposed** *(2026-06-30; matrix draft 0.1 ready)* | 002, 003, 005, 006, 007, 008, 009, 010, 011, 012, 013, 014, 016 | Amends ADR-003 (metadata catalog layer, no source taxonomy change) and ADR-005 (new schema kind surface); defines Ready/Candidate registries, render strategies, RuntimeContext, Ready gates, plugin model |
| [018](./018-block-catalog-data-action-reactivity.md) | Block Catalog Data, Action, Reactivity Contract | **Proposed** *(2026-06-30; required by ADR-017 before any Ready entry)* | 002, 004, 005, 006, 007, 011, 012, 013, 017 | Adds ActionSource, ActionResult, ActionEffect, ActionLifecycle; Portable Expression AST and built-in function registry; cross-platform parity discipline |

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
              │                 │                  │
              └────────┬────────┘                  │
                       ▼                           ▼
                   ADR-005 ───────────────────► ADR-007
              (metadata-driven UI)       (enterprise universal
                                            API contract)
                                                   │
                                                   ▼
                                                ADR-008
                                          (contract-first freeze)
                                                   │
                                    ┌──────────────┴──────────────┐
                                    ▼                              ▼
                               ADR-009                       ADR-010
                       (Runtime Part Resolution)    (Web Component Event Bridge)
                                                                   │
                                                                   ▼
                                                               ADR-011
                                                   (Flutter Adapter Architecture)
                                                                   │
                                                                   ▼
                                                               ADR-012
                                          (Extract @bitcode/ui-schema standalone package)
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
2. Nomor urut berikutnya: **ADR-013**.
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
