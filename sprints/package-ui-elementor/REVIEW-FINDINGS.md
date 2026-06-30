# Block Catalog Review Findings

Review date: 2026-06-30

## Scope

Evidence-first review of `block-catalog-candidate-matrix.json`, ADR-017, ADR-018, Flutter component parity, web root exports, and ADR-014 overlap.

### [HIGH] catalog-counts: Candidate count mismatch
- Lokasi: `docs/adr/017-block-catalog-architecture.md` Context, `docs/components-progress.md` Block Catalog counts
- Masalah: The matrix source of truth contains 44 candidates, but ADR-017 and `components-progress.md` still say 43 total candidates. Status buckets sum to 44: schema-preset-only 8, existing-component-needs-binding 21, component-gap 10, contract-gap 2, adapter-gap 1, duplicate-alias 1, excluded 1.
- Rekomendasi: Update ADR-017 and `components-progress.md` to say 44 candidates and keep bucket counts unchanged.
- Status: RESOLVED IN DOCS

### [HIGH] flutter-parity: Missing Flutter components referenced by candidate mappings
- Lokasi: `packages/ui-flutter/lib/src/components/`, `packages/ui-flutter/lib/bc_ui_flutter.dart`, matrix rows `content.testimonial`, `content.team-member`, `data.account-summary`, `workflow.state-indicator`, `analytics.kpi-grid`, `analytics.chart`
- Masalah: `bc_avatar.dart`, `bc_badge.dart`, `bc_progress.dart`, and `bc_kpi.dart` are absent from Flutter source and absent from the package root export. Rows that map to `BcAvatar`, `BcBadge`, `BcProgress`, or `BcKpi` cannot become Ready unless Flutter parity is added or the presets are remapped to existing Flutter widgets such as `bc_statistic`, `bc_taggroup`, or `bc_box`.
- Rekomendasi: Keep these rows non-Ready. Either add the missing Flutter widgets with tests and barrel exports, or adjust matrix mappings to existing Flutter-ready primitives before promotion.
- Status: OPEN

### [MEDIUM] flutter-verified-existing: Flutter social and markdown mappings verified
- Lokasi: `packages/ui-flutter/lib/src/components/`, `packages/ui-flutter/lib/bc_ui_flutter.dart`
- Masalah: `bc_markdown.dart`, `bc_chatter.dart`, `bc_activity.dart`, and `bc_timeline_social.dart` exist and are exported from the Flutter root barrel. These rows do not need parity-gap treatment for those specific widgets.
- Rekomendasi: Keep these mappings, but still validate each block's data/action contract before Ready promotion.
- Status: OPEN

### [HIGH] web-export-gap: BcShare and BcGridView source exists but root exports are missing
- Lokasi: `packages/ui/src/components/social/share/`, `packages/ui/src/components/views/grid/`, `packages/ui/src/index.ts`
- Masalah: `BcShare` and `BcGridView` source implementations exist, but neither is exported from `packages/ui/src/index.ts`. Matrix rows already flag this for `content.social-links`, `data.record-list`, and `commerce.product-grid`.
- Rekomendasi: Do not mark dependent rows Ready until the generated root export/registry path includes these components. Coordination needed because `packages/ui/src/index.ts` is shared and auto-generated.
- Status: OPEN

### [HIGH] web-markdown-gap: BcMarkdown web source and root export are missing
- Lokasi: `packages/ui/src/components/**`, `packages/ui/src/index.ts`, matrix row `content.rich-text`
- Masalah: `BcMarkdown` is referenced in `mapsTo` for rich text and exists on Flutter, but there is no web source under `packages/ui/src` and no root export in `packages/ui/src/index.ts`.
- Rekomendasi: Either add a web `BcMarkdown` component with tests and export, or remove it from web-facing Ready requirements and rely on `BcRichText` plus sanitize policy until web parity exists.
- Status: OPEN

### [MEDIUM] adr014-overlap: ADR-017/018 must not claim ADR-014 implementation ownership
- Lokasi: `docs/adr/014-engine-ui-feature-wiring.md`, ADR-017/018 scope sections
- Masalah: ADR-014 already owns engine feature wiring for `BcChatter`, `BcActivity`, `BcExport`, `BcWizard`, pivot, notification center, and visual builders. Block Catalog ADRs would conflict if they claim endpoint, engine model, or component wiring ownership for those features.
- Rekomendasi: ADR-017 should define catalog classification, Ready gates, render strategies, and plugin model only. ADR-018 should define generic data/action/reactivity contracts. Engine feature wiring remains ADR-014.
- Status: OPEN

### [MEDIUM] runtime-claim-risk: Current skeleton overclaims component counts without current evidence
- Lokasi: `docs/adr/017-block-catalog-architecture.md` Context
- Masalah: `packages/ui/src/index.ts` says `Components: 187`, and Flutter progress docs say full parity reached, but the matrix review found specific Flutter parity gaps and root export gaps. A count alone is not proof that every mapped component is available for Ready validation.
- Rekomendasi: ADR-017 should treat counts as background only, and make component file plus root export validation mandatory for Ready entries.
- Status: OPEN
