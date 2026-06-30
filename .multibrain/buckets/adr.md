# Bucket: ADR (Architecture Decision Records)

> Source of truth untuk keputusan arsitektur BitCode platform.
> ADR = sumber kebenaran. Kode kontradiksi ADR → perbaiki kode atau buat ADR baru.

**Index lengkap**: `docs/adr/README.md`
**Total ADR**: 018 (012-018 Accepted/complete except 013/014 Proposed)

---

## Status Saat Ini (per 2026-06-30)

| ADR | Judul | Status | Catatan |
|-----|-------|--------|---------|
| 001 | Mitosis Component Library | Accepted (Phase C1) | parts superseded by 002, §15 Lynx superseded by 006 |
| 002 | Cross-Platform Architecture | Accepted | extended by 006/007/009/010 |
| 003 | Component Roadmap & Taxonomy | Accepted (Phase C1) | 14 grup, P0/P1/P2 |
| 004 | Reactive Engine | Accepted | dual storage (useStore+useRef) cleanup |
| 005 | Metadata-Driven UI (BcRenderer) | Accepted (Phases 1-3) | ~5% done per .opencode-context.md |
| 006 | Universal Component Strategy | Accepted (Phase 0 spike) | Pure/Adapt 2-tier |
| 007 | Enterprise Universal UI API Contract | Accepted (Phase 2) | event/state/method/prop/a11y |
| 008 | Contract-First API Freeze | Accepted (Phase 2 + partAttr wiring) | Core Integration Surface complete |
| 009 | Runtime Part Resolution | Accepted | partClass/partAttr pure fns |
| 010 | Web Component Event Bridge | Accepted | dispatchBcEvent |
| 011 | Flutter Adapter Architecture | Accepted | 75% shared/25% separate, BcRenderer hybrid port done |
| 012 | Extract @bitcode/ui-schema | Accepted | zero-dep package, unblocks Flutter port |
| 013 | Cross-tier Infrastructure Completion | **Proposed** (2026-06-26) | AppProvider adoption, P0/P1 hooks, animation+adapt contract |
| 014 | Engine UI Feature Wiring | **Proposed** (2026-06-26) | 9 ERP UI features → component-backed SSR render |
| 015 | Flutter ECharts Chart Strategy | Accepted | deprecate CustomPaint, 26/26 parity |
| 016 | Demo App QA Bug Audit | Accepted (all 4 phases, 2026-06-29) | 28 stubs ported, 5 arch decisions |
| 017 | Block Catalog Architecture | **Accepted** (2026-06-30, matrix 0.1 baseline) | amends ADR-003 + ADR-005 |
| 018 | Block Catalog Data/Action/Reactivity | **Accepted** (2026-06-30) | required by ADR-017 before any Ready entry |

---

## Aturan ADR (rekap)

- ADR yang sudah Accepted **jangan diubah** — buat ADR baru yang supersedes/extends.
- Nomor urut berikutnya: **ADR-019**.
- Update tabel di `docs/adr/README.md` setiap ADR baru/status berubah.
- ADR-013/014/015 dirancang **paralel-safe** (lihat README paralelisasi note).
- Untuk "udah sampai mana komponen?" → baca `docs/components-progress.md`, BUKAN ADR.

---

## Log

<!-- Append-only: newest first. Format: YYYY-MM-DD | summary | entry ref -->
- 2026-07-01 | P2 Portable Expression AST landed in @bitcode/ui-schema (12 TS source files: ast, operators, scope, errors, permission-scope, evaluator iface, function-registry, builtin-manifest 70+ fns, golden-types, validate, serialize, index barrel) + Dart mirror in ui-flutter (10 files: ast manual sealed class, operators, errors, scope, permission-scope, function-registry, builtin-manifest, validate, evaluator, expressions barrel). 9 golden JSON fixtures (122 cases, BOM-stripped) as cross-platform source of truth. 32 TS tests pass (26 unit + 6 golden-fixture structural). 136 Dart tests pass (22 validate + 114 golden-parity). Security: FORBIDDEN_FUNCTION_NAMES denylist + path forbidden-segments + reserved-roots enforcement. Next: P3 ActionSource/ActionAdapter + runtime evaluator. | packages/ui-schema/src/expressions/, packages/ui-flutter/lib/src/expressions/
- 2026-06-30 | P1 Block Catalog foundation types landed in @bitcode/ui-schema: BlockCatalogEntry/CandidateCatalogEntry, enums (category/channel/capability/readiness/renderStrategy), Required* unions, 4 policies, migration/a11y/test metadata, capability normalization (action-trigger -> action-source), minimal shape validator (lenient Candidate/strict Ready + BLOCK_READY_GATE_INCOMPLETE), 2 Candidate fixtures (content.hero schema, interaction.notification-center component). tsc --noEmit passes. No Ready promotion. Next: P2 Portable Expression AST. | packages/ui-schema/src/block-catalog/, commit a2d4857
- 2026-06-30 | ADR-017 (Block Catalog) dan ADR-018 (Data/Action/Reactivity) Accepted setelah final read-through, matrix 0.1 baseline | docs/adr/017, 018
- 2026-06-30 | ADR-017 (Block Catalog) dan ADR-018 (Data/Action/Reactivity) diajukan, matrix draft 0.1 siap | docs/adr/017, 018
- 2026-06-29 | ADR-016 Demo App QA complete — 4 phases, 28 stubs ported, 5 arch decisions | docs/adr/016
- 2026-06-26 | ADR-013, 014, 015 diajukan (cross-tier infra, engine UI wiring, flutter echarts) | docs/adr/013-015
- 2026-06-22 | ADR-011, 012 Accepted (Flutter adapter + ui-schema standalone) | docs/adr/011, 012
