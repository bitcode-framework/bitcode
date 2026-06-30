# ADR-017: Block Catalog Architecture

| Field | Value |
|-------|-------|
| **Status** | Proposed (2026-06-30; full body drafted against matrix draft 0.1) |
| **Depends on** | 002, 003, 005, 006, 007, 008, 009, 010, 011, 012, 013, 014, 016 |
| **Amends** | 003 (adds metadata catalog layer, does not change source taxonomy), 005 (extends schema surface with block-level composition) |
| **Impact** | Adds Block Catalog types to `@bitcode/ui-schema`, Ready/Candidate registries, per-category preset catalogs in `packages/ui`, render strategies, RuntimeContext requirements, validation gates, and plugin/extension model |

---

## Context

BitCode now has a broad standalone UI stack: `packages/ui` root export currently reports 187 web component exports, `packages/ui-flutter` has a mature Dart widget surface, ADR-002 defines the 4-layer `DataSource`, ADR-004 defines `ReactiveStore`, ADR-005 defines `BcRenderer`, ADR-007/008 define the universal component contract, ADR-011 defines Flutter as a first-class Tier C track, and ADR-013 completes cross-tier providers, hooks, animation, and adapt infrastructure.

That component inventory is not enough for a page builder or application builder. A builder needs curated, reusable page-level units such as hero, record list, checkout, chatter, workflow state, and dashboard KPI grid. Those units must be portable across web, desktop, and Flutter, must not depend on BitCode engine internals, and must fail loudly when a host lacks a required component, function, adapter, permission capability, or platform implementation.

The source of truth for this ADR is `sprints/package-ui-elementor/block-catalog-candidate-matrix.json` draft 0.1. It maps Elementor, Website/CMS, EasyERP, portal, commerce, workflow, analytics, and builder references into 44 curated candidates. Current status distribution is: 8 `schema-preset-only`, 21 `existing-component-needs-binding`, 10 `component-gap`, 2 `contract-gap`, 1 `adapter-gap`, 1 `duplicate-alias`, and 1 `excluded`.

Evidence from the 2026-06-30 review matters for readiness discipline:

- `BcShare` and `BcGridView` source files exist in `packages/ui`, but are not exported from `packages/ui/src/index.ts`; this is treated as an export/registry bug, not a component gap.
- `BcMarkdown` exists and is exported on Flutter, but has no `packages/ui` web source or root export; the resolved direction is to implement canonical web `BcMarkdown`.
- `bc_avatar.dart`, `bc_badge.dart`, `bc_progress.dart`, and `bc_kpi.dart` are absent from Flutter source and root exports; the resolved direction is to add canonical Flutter widgets instead of remapping to generic primitives.
- ADR-014 already owns engine feature wiring for chatter, activity, export, wizard, pivot, notification center, and visual builders. This ADR defines catalog architecture only. It must not claim ADR-014 endpoint, model, permission, realtime, persistence, or engine wiring ownership.

Therefore, the Block Catalog must be matrix-first and validation-first. A candidate can be listed without being Ready. A Ready entry is a stronger claim: it is loadable, tested, secure, accessible, adaptive, and fully implemented on web and Flutter.

## Scope

- Define the Block Catalog surface, metadata model, categories, channels, and capability tags.
- Define Ready and Candidate registries, render strategies, RuntimeContext requirements, and validation rules.
- Define how built-in presets live in `packages/ui` while canonical types live in `@bitcode/ui-schema`.
- Define Flutter-first parity discipline and fail-loud validation.
- Define plugin and adapter catalog boundaries for provider-specific blocks.

## Out of Scope

- `ActionSource`, `ActionResult`, `ActionEffect`, lifecycle, and Portable Expression AST internals. These are owned by ADR-018.
- BitCode engine endpoint wiring, built-in models, and SSR to component-backed feature integration. These remain ADR-014 scope.
- Module-specific CRM/HRM/ERP presets in the core Ready Catalog.
- Provider-specific implementations for Stripe, PayPal, Google Maps, Facebook, SoundCloud, and similar services.

## Decision

### Section 1: Catalog Surface

Introduce `BlockCatalogEntry` as a metadata descriptor, not a new source taxonomy and not a new component class by default.

Canonical types live in `@bitcode/ui-schema`:

```ts
interface BlockCatalogEntry {
  id: string;
  version: string;
  schemaVersion: string;
  labelKey: string;
  fallbackLabel: string;
  category: BlockCategory;
  subcategory?: string;
  channels: BlockChannel[];
  capabilities: BlockCapability[];
  renderStrategy: 'schema' | 'component';
  schemaPreset?: UiSchema;
  componentKey?: string;
  requiredComponents?: string[];
  requiredFunctions?: string[];
  requiredAdapters?: string[];
  adaptivePolicy?: BlockAdaptivePolicy;
  securityPolicy?: BlockSecurityPolicy;
  cachePolicy?: BlockCachePolicy;
  offlinePolicy?: BlockOfflinePolicy;
  migration?: BlockMigrationMetadata;
  readiness: 'candidate' | 'ready' | 'deprecated';
}
```

The exact field names can evolve during implementation, but the responsibilities are fixed:

- `id` is a namespaced slug such as `content.hero` or `commerce.product-grid`.
- Component exports keep `Bc*` naming. Block preset IDs never become component names.
- `version` and `schemaVersion` are mandatory. Catalog migration is explicit, never inferred from shape alone.
- User-facing labels use `labelKey` plus `fallbackLabel`. Ready entries must have an English fallback.
- Built-in presets are stored in per-category files in `packages/ui` with a generated catalog index.
- Engine or module JSON may reference and override catalog entries, but the core catalog remains standalone.

A separate `CandidateCatalogEntry` can carry matrix-only fields such as aliases, source references, status, decision, and gap notes. Runtime must read only the Ready Catalog.

### Section 2: Categories, Channels, and Capabilities

Categories describe capability families. They do not replace ADR-003 component taxonomy.

Allowed top-level categories for catalog v1:

1. `content`
2. `data`
3. `interaction`
4. `commerce`
5. `workflow`
6. `builder`
7. `analytics`
8. `navigation`

`account` is a subcategory or tag, not a top-level category.

Channels describe where a block is commonly used:

- `website`
- `portal`
- `internal`
- `saas`

Channels are metadata and filtering hints. They are not security boundaries.

Capabilities drive validation and runtime requirements. Examples include:

- `static`
- `data-source`
- `action-source`
- `reactive`
- `permission-aware`
- `i18n`
- `a11y`
- `responsive`
- `adapt`
- `provider-adapter`
- `commerce`
- `payment-adapter`
- `workflow`
- `builder`
- `large-data`
- `offline`
- `sanitize`
- `security`

A capability that implies runtime behavior must map to a validator. For example, `action-source` requires ADR-018 action support, `permission-aware` requires a permission/authz adapter, and `large-data` requires an explicit performance policy.

### Section 3: Render Strategies

Every entry declares one explicit render strategy.

#### `schema`

A schema block resolves to `UiSchema` and delegates to `BcRenderer`. This is the default for composed blocks such as hero, FAQ, record list variants, pricing table, workflow state indicator, and product grid.

Schema strategy is valid only if:

- The schema contains only serializable values.
- Required components are declared.
- Required functions and adapters are declared.
- Reactions use ADR-018 Portable Expression AST for Ready behavior.
- Action behavior uses ADR-018 `ActionSource`, not opaque callback strings.

#### `component`

A component block resolves to a registered component key. It is used for behavior-heavy blocks such as import wizard, cart/checkout, workflow diagram editor, query builder, report builder, or notification center when schema composition is not enough.

Component strategy is valid only if:

- The catalog JSON references a string `componentKey`, never a function or class.
- Each platform registers the implementation with a typed registry, for example `registerBlockComponent('BcCheckout', BcCheckout)` on web and an equivalent Dart registry on Flutter.
- Required props are serializable or adapter-resolved.
- The component has web and Flutter implementations before Ready promotion.

`BcBlockRenderer` is a thin orchestration layer over `BcRenderer` and the platform component registry. It must not become a second UI renderer.

Resolution order:

1. Load the catalog entry by ID and version.
2. Validate host RuntimeContext against `requiredComponents`, `requiredFunctions`, and `requiredAdapters`.
3. If `renderStrategy = 'schema'`, materialize `UiSchema` and delegate to `BcRenderer`.
4. If `renderStrategy = 'component'`, resolve `componentKey` from the platform registry.
5. Fail loudly if anything is missing or incompatible.

### Section 4: Ready Gate

Ready means production-grade. It does not mean "planned", "source exists", or "can be mocked later".

A candidate can become Ready only when all gates pass:

| Gate | Requirement |
|---|---|
| Matrix gate | Candidate row has a clear `keep`, `merge`, `split`, or `alias` decision with gaps resolved or explicitly deferred outside Ready scope. |
| Web component gate | Every required web component exists and is exported from `packages/ui/src/index.ts` or the approved public entrypoint. |
| Flutter component gate | Every required Flutter widget source exists and is exported from `packages/ui-flutter/lib/bc_ui_flutter.dart`. |
| Schema gate | Schema is valid `@bitcode/ui-schema`, serializable, versioned, and migration-aware. |
| Action gate | Mutations use ADR-018 `ActionSource`; no arbitrary callback names in Ready JSON. |
| Expression gate | Required reactivity uses ADR-018 Portable Expression AST or built-in functions. String expressions are advanced escape hatches, not Ready foundations. |
| RuntimeContext gate | Required functions, services, adapters, permissions, and i18n entries are declared and validated. |
| Security/privacy gate | No secrets in schema; URL, rich text, file, payment, permission, and provider policies pass. Sanitization defaults on. |
| Adaptive gate | Block declares or inherits adaptive policy. Complex blocks explicitly cover breakpoints, platform morphing, reduced motion, loading, empty, and error states. |
| Performance gate | Data-heavy blocks declare pagination, virtualization, debounce, abort, and bounded Flutter scrolling where applicable. |
| A11y gate | Block-level accessibility contract is documented and tested. |
| Test gate | Schema validation, web render, Flutter widget render, action mock, and a11y tests pass. |
| Demo gate | A demo exists in `packages/ui/demo-app` and `packages/ui-flutter/example`; engine integration demo exists only when engine-bound. |

Current known gaps from matrix review mean several candidates must remain non-Ready until fixed:

- `content.social-links`, `data.record-list`, and `commerce.product-grid` are blocked by web root export gaps for `BcShare` or `BcGridView`. Source exists, so the fix is public export/registry generation, not a new component.
- `content.rich-text` requires canonical web `BcMarkdown` source, tests, and public export before markdown-backed rich text can be Ready.
- Rows depending on `BcAvatar`, `BcBadge`, `BcProgress`, or `BcKpi` require canonical Flutter widgets with source, tests, and barrel exports before Ready. `BcKpi` may wrap `BcStatistic` internally, but the public widget key must exist.

### Section 5: RuntimeContext

Block rendering uses an explicit named registry. It must not use positional arrays, global singletons, or implicit BitCode engine state.

`RuntimeContext` extends ADR-013 AppProvider/BcAppConfig concepts with block-capability registries:

```ts
interface RuntimeContext {
  values: Record<string, unknown>;
  services: Record<string, unknown>;
  functions: FunctionRegistry;
  adapters: {
    data?: DataAdapterRegistry;
    action?: ActionAdapterRegistry;
    permission?: PermissionAdapter;
    i18n?: I18nAdapter;
    navigation?: NavigationAdapter;
    media?: MediaAdapterRegistry;
  };
  permissions?: PermissionSnapshot;
  i18n?: I18nRuntime;
}
```

Rules:

- JSON references functions and adapters by stable names.
- Host/project functions require explicit manifest entries, platform support, and fallback policy.
- Missing required context is a validation error, not silent fallback.
- Permission/authz is a typed adapter capability. Core catalog never hardcodes BitCode RBAC.
- Engine is one possible provider adapter, not the owner of the catalog contract.

### Section 6: Plugin and Adapter Catalog

Provider-specific integrations belong in plugin or adapter catalogs, not the core Ready Catalog.

Examples:

- Stripe and PayPal payment buttons become payment adapter plugins over generic `commerce.payment-action`.
- Google Maps becomes a map provider adapter over generic `data.map-view`.
- Facebook, SoundCloud, Instagram, TikTok, and similar embeds become provider adapters over generic media/comment/chatter/embed primitives.

Core Ready Catalog may include generic equivalents, adapter contracts, and mock/local adapters for tests, but it must not contain provider SDK implementations, provider secrets, provider-specific API keys, or schema fields that imply a single commercial vendor.

A plugin catalog entry must declare:

- Provider ID and version.
- Required adapters and capabilities.
- Security policy and URL allowlist rules.
- Supported platforms.
- Fallback behavior when provider SDK is unavailable.
- Tests or conformance fixtures against the generic contract.

### Section 7: Flutter-First Parity Discipline

Flutter is not a downstream afterthought. A Ready block must fully render on Flutter from day one.

Validation must check both source and public export:

- Web: component exists and is exported from `packages/ui/src/index.ts` or an approved public entrypoint.
- Flutter: `bc_*.dart` source exists and is exported from `packages/ui-flutter/lib/bc_ui_flutter.dart`.
- Renderer infrastructure such as Flutter `BcRenderer` is validated as infrastructure, not as a component file.
- Platform-aware behavior follows ADR-011: iOS/macOS should use Cupertino patterns where appropriate, Android/web/desktop use Material patterns where appropriate, and Material-only components are acceptable only when no native Cupertino equivalent exists.

No Ready block may rely on a "coming soon", placeholder, unsupported platform message, or hidden stub.

### Section 8: Validation

Validation has two layers.

#### Build-time validation

Build-time validators run in package CI and catalog generation:

- JSON schema validation for Ready and Candidate catalogs.
- Status and gap consistency checks against the candidate matrix.
- Required component source and root export checks for web and Flutter.
- Required function and adapter manifest checks.
- Security policy checks for URL, rich text, file, payment, and provider capabilities.
- i18n key plus English fallback checks.
- A11y and adaptive policy presence checks.
- No-stub checks for Ready entries.

#### Runtime validation

Runtime validators protect host integration:

- `BcBlockRenderer` checks that required functions, adapters, components, permissions, and i18n entries are registered before rendering.
- Missing requirements produce typed errors suitable for developer tooling and builder diagnostics.
- Runtime must not silently degrade a Ready block into a partial rendering unless the entry declares an explicit fallback.

## Consequences

### Positive

- Builders get a curated, production-grade block surface without changing ADR-003 component taxonomy.
- Ready status becomes auditable and enforceable, not aspirational.
- Flutter parity is protected by validation before entries are promoted.
- Core catalog remains standalone and engine-neutral.
- Provider-specific complexity is isolated behind plugins/adapters.
- Schema strategy reuses `BcRenderer` instead of creating another renderer stack.

### Negative

- Ready promotion is slower because every entry needs web, Flutter, tests, demos, a11y, security, and runtime validation.
- Catalog metadata adds a new maintenance surface in `@bitcode/ui-schema` and `packages/ui`.
- Some entries that look visually simple remain blocked until ADR-018 ActionSource or Portable Expression AST exists.
- Export and registry hygiene becomes a hard blocker, so generated barrels must be reliable.

### Trade-offs

- The catalog chooses fewer Ready entries with strict guarantees over many optimistic entries with stubs.
- Schema presets are preferred for maintainability, but component strategy remains available for behavior-heavy blocks.
- Provider-specific widgets are excluded from core even if they are popular, because generic adapter contracts are safer and more future-proof.

## Implementation Phases

Phases follow the matrix `priorityBatch` values. Each phase can promote only entries that pass the Ready gate.

### Phase 0: Review and validation foundation

- Keep `block-catalog-candidate-matrix.json` as source of truth.
- Generate the markdown parity report from the JSON matrix.
- Add catalog schema validation and status-count checks.
- Create `REVIEW-FINDINGS.md` for factual gaps.
- Fix documentation count drift.

### Phase 1: Foundational catalog infrastructure

- Add BlockCatalogEntry, CandidateCatalogEntry, category/channel/capability enums, adaptive policy, security policy, i18n metadata, offline/cache policy, versioning, and migration metadata to `@bitcode/ui-schema`.
- Add per-category catalog files in `packages/ui` and a generated Ready index.
- Add `BcBlockRenderer` with `schema` and `component` strategies.
- Add platform component registries for web and Flutter.
- Add build-time and runtime validators.

### Phase 2: ERP data and interaction

- Promote record list/detail, account summary, download center, form, wizard, record picker, import wizard, and notification center only after ADR-018 action/runtime requirements are available.
- Resolve web export gaps for `BcGridView` before record-list and product-grid promotion.
- Validate mock demos and engine-bound demos where applicable.

### Phase 3: Workflow

- Promote audit log, approval actions, state indicator, diagram viewer, and chatter according to matrix status.
- Keep workflow diagram editor as component-gap until dedicated authoring component exists.
- Defer ADR-014 endpoint, model, permission, realtime, persistence, and engine feature wiring to ADR-014 for chatter, activity, export, wizard, pivot, notification center, and visual builders.

### Phase 4: Builder surface

- Implement builder components only as component strategy entries.
- Keep builder runtime separate from Ready block runtime so page/application builders do not become dependencies of simple runtime rendering.
- Add validation, undo/versioning, persistence, and authoring-specific security gates.

### Phase 5: Analytics

- Promote KPI grid and chart blocks after canonical Flutter `BcKpi` and `BcProgress` widgets exist with tests and barrel exports.
- Reuse ADR-015 shared ECharts option contract for chart variants.
- Enforce large-data and dashboard performance budgets.

### Phase 6: Commerce

- Promote product grid/detail as commerce variants of generic record list/detail.
- Implement cart/checkout and payment action only with generic payment adapters.
- Keep provider-specific buttons in plugin catalogs.

### Phase 7: Content and navigation/CMS

- Promote hero, CTA, rich text, media, gallery, FAQ, testimonial, team member, pricing table, countdown, social links, menu, and search/TOC/sitemap according to resolved gaps.
- Add resolved missing Flutter/web components before Ready: canonical Flutter `BcAvatar`, `BcBadge`, `BcProgress`, `BcKpi`; canonical web `BcMarkdown`; public exports/registry for `BcShare` and `BcGridView`.
- Enforce URL allowlists, rich text sanitization, external-link a11y, and responsive/adaptive policies.

## Acceptance Criteria

### ADR-level

- [ ] `BlockCatalogEntry` and `CandidateCatalogEntry` types live in `@bitcode/ui-schema`.
- [ ] Category, channel, capability, readiness, and render strategy enums are typed and validated.
- [ ] Built-in catalog source is split by category in `packages/ui` with a generated index.
- [ ] Runtime reads only Ready entries by default.
- [ ] `BcBlockRenderer` delegates schema entries to `BcRenderer` and component entries to a platform registry.
- [ ] Build-time validators reject Ready entries with missing web or Flutter source/export requirements.
- [ ] Runtime validators fail loudly on missing components, functions, adapters, permission capability, or i18n entries.
- [ ] Ready entries have web and Flutter demos.
- [ ] Ready entries with data/action behavior have mock demos; engine integration demos are required only when engine-bound.
- [ ] Third-party-specific providers are represented by plugin/adapter catalog entries, not core Ready entries.

### Per Ready block

- [ ] Web component source and root export verified.
- [ ] Flutter component source and root export verified.
- [ ] Required components/functions/adapters declared.
- [ ] Security/privacy policy declared and validated.
- [ ] i18n key and English fallback declared.
- [ ] A11y contract declared and tested.
- [ ] Adaptive policy declared or inherited.
- [ ] Performance policy declared for data-heavy blocks.
- [ ] Schema/action/expression contracts are serializable.
- [ ] Tests pass for schema validation, web render, Flutter render, action mock where applicable, and a11y.

## Related ADRs

- ADR-002: Cross-platform architecture, tokens, data source, infrastructure slots.
- ADR-003: Component taxonomy and roadmap. This ADR adds metadata above it without changing it.
- ADR-004: Reactive engine. ADR-018 extends it for action execution and portable catalog expressions.
- ADR-005: Metadata-driven UI. This ADR composes `BcRenderer` through `BcBlockRenderer`.
- ADR-006: Single Caller API and Pure/Adapt implementation strategy.
- ADR-007: Universal component API contract.
- ADR-008: Contract-first API freeze and Ready quality bar.
- ADR-009: Runtime Part Resolution for customization.
- ADR-010: Web Component event bridge.
- ADR-011: Flutter adapter architecture and platform-aware behavior.
- ADR-012: `@bitcode/ui-schema` as canonical portable type package.
- ADR-013: AppProvider, hooks, animation, adapt, and Flutter parity infrastructure.
- ADR-014: Engine UI feature wiring. Related but separate implementation ownership.
- ADR-016: Demo app QA and no-stub discipline.
- ADR-018: Data, action, and reactivity contract required by block runtime.

## References

- `sprints/package-ui-elementor/block-catalog-candidate-matrix.json`
- `sprints/package-ui-elementor/block-catalog-parity-matrix.md`
- `sprints/package-ui-elementor/REVIEW-FINDINGS.md`
- `CONTEXT.md` glossary for Website Block, ERP Interaction Component, Block Catalog, Ready Block
- `packages/ui-schema/src/`
- `packages/ui/src/core/renderer/`
- `packages/ui/src/index.ts`
- `packages/ui-flutter/lib/bc_ui_flutter.dart`
