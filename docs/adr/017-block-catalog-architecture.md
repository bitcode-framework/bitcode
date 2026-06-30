# ADR-017: Block Catalog Architecture

| Field | Value |
|-------|-------|
| **Status** | Proposed (2026-06-30; matrix draft 0.1 ready in `sprints/package-ui-elementor/block-catalog-candidate-matrix.json`; ADR body to be drafted after matrix approval) |
| **Depends on** | 002, 003, 005, 006, 007, 008, 009, 010, 011, 012, 013, 016 |
| **Amends** | 003 (adds metadata catalog layer, does not change source taxonomy), 005 (extends schema kind surface with `block`) |
| **Impact** | Adds `@bitcode/ui-blocks` package proposal, Ready/Candidate registries, capability categories, channels, render strategies, RuntimeContext, plugin/extension model |

---

## Context

The BitCode standalone UI library (`packages/ui` + `packages/ui-flutter`) already ships 187 web components and 160+ Flutter widgets, a 4-layer `DataSource` contract (ADR-002), a reactive engine (ADR-004), a metadata-driven renderer (ADR-005), and a contract-first API freeze (ADR-008). However, there is no curated block catalog that lets builders compose pages from reusable page-level units in a portable, no-stub, Flutter-first way.

A candidate parity matrix has been drafted (0.1) at `sprints/package-ui-elementor/block-catalog-candidate-matrix.json`, mapping Elementor/EasyERP references and a curated BitCode/ADR-014 set into 43 candidate blocks. After reclassification the candidates distribute as: 8 schema-preset-only, 21 existing-component-needs-binding, 10 component-gap, 2 contract-gap, 1 adapter-gap, 1 duplicate-alias, 1 excluded.

This ADR proposes the architecture that turns the candidate matrix into a production-grade Block Catalog.

## Scope

- Define a Block Catalog architecture that is matrix-first, runtime-first, and Flutter-first.
- Establish a Ready Block gate that forbids stubs and partial Flutter parity.
- Define runtime registries and validation that make Ready entries actually loadable.

## Out of scope

- Implementing the 5 foundational work items (handled by ADR-018).
- Module-specific presets for BitCode/ERP/CRM/HRM.
- Third-party plugin implementations (pay, social embed, maps) beyond plugin model.

## Decision Summary (Draft)

1. **Catalog surface**: `BlockCatalogEntry` is a metadata-only descriptor. Source taxonomy (ADR-003) and component code are unchanged. Categories are capability descriptors (`content`, `data`, `interaction`, `commerce`, `workflow`, `builder`, `analytics`, `navigation`). `account` is a subcategory/tag.
2. **Channel tags**: `website`, `portal`, `internal`, `saas` are channels, not categories. A single entry can declare multiple channels.
3. **Capability tags**: `static`, `data-source`, `action-source`, `reactive`, `mutation`, `permission-aware`, `i18n`, `a11y`, `responsive`, `performance-budget`, etc. Capabilities gate runtime requirements.
4. **Render strategies**: `schema` (resolves to `UiSchema` and delegates to `BcRenderer`) or `component` (resolves to a registered `Bc*` component via platform registry). `BcBlockRenderer` is a thin layer over `BcRenderer`.
5. **Ready gate**: a candidate becomes Ready only if `webComponents`, `flutterComponents`, `requiredFunctions`, `requiredAdapters`, schema/action/expression/security/i18n/adaptive/test gates, and per-block tests all pass. Runtime validators fail loud if any requirement is unmet.
6. **RuntimeContext**: extends existing AppProvider/BcAppConfig capability injection (ADR-013) with named functions and services. JSON only references `name`; the host registers typed implementations.
7. **Plugin model**: third-party embeds/payments/maps live in a separate plugin catalog. Core Ready Catalog never references provider-specific keys.
8. **Flutter-first parity discipline**: Flutter component file existence and root export are validated at build time. Entries referencing missing or unexported components cannot be Ready.

## Acceptance Criteria (Draft)

- BlockCatalogEntry, Ready/Candidate types, and enums (category, channel, capability, renderStrategy) live in `@bitcode/ui-schema`.
- `BcBlockRenderer` resolves an entry into either `UiSchema` (delegated to `BcRenderer`) or a registered component.
- Build-time and runtime validators enforce the Ready gate; CI fails on violation.
- Per-block tests cover schema validation, web render, Flutter widget render, action lifecycle, and a11y keyboard.

## Status Notes

This ADR is currently a skeleton. The matrix in `sprints/package-ui-elementor/block-catalog-candidate-matrix.json` is the authoritative source of evidence. Full Decision, Consequences, and Implementation Phases sections will be drafted after matrix approval.
