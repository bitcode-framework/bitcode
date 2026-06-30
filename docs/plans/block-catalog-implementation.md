# Block Catalog Implementation Plan

Status: Draft implementation appendix, ready for executor planning
Date: 2026-06-30
Architecture baseline: ADR-017 and ADR-018 Accepted
Matrix baseline: `sprints/package-ui-elementor/block-catalog-candidate-matrix.json` version `0.1.0`, status `accepted-baseline`, `baselineFor: ["ADR-017", "ADR-018"]`

This document is the execution appendix for implementing the Block Catalog after ADR-017 and ADR-018 acceptance. It is intentionally more operational than the ADRs. The goal is that an implementator can execute the work in safe, reviewable batches without re-litigating architecture decisions or accidentally promoting incomplete blocks to Ready.

## 1. Non-negotiable execution rules

1. ADR-017 and ADR-018 are accepted architecture. Do not change their contract during implementation unless a new ADR or documented ADR amendment is approved.
2. Matrix `0.1.0` is the accepted candidate baseline. Later candidate changes require a matrix version bump or a documented change record.
3. Runtime reads Ready catalog entries only. Candidate entries may exist for planning, tests, reports, and builder diagnostics, but not as default runtime input.
4. No Ready entry may depend on a stub, placeholder, coming-soon component, unsupported platform message, or hidden fallback.
5. Ready means full web and Flutter support from day one, including source, public export, tests, demos, a11y evidence, adaptive policy, security policy, versioning, and validation.
6. Core remains engine-neutral. BitCode engine integration is one adapter path, not the owner of the catalog contract.
7. Provider SDKs stay out of core. Core may ship generic adapter contracts and mock/local adapters for tests and demos only.
8. ADR-014 owns endpoint, model, permission, realtime, persistence, and engine feature wiring for overlap features. This plan owns catalog/runtime contracts and standalone UI package implementation.
9. Every batch must leave the repository in a coherent state, with validators passing and no partial Ready claims.
10. If a block fails any Ready gate, keep it Candidate and record the blocker in the matrix/report. Do not weaken gates.

## 2. Source artifacts and authority

| Artifact | Role | Authority |
|---|---|---|
| `docs/adr/017-block-catalog-architecture.md` | Catalog architecture, Ready gates, render strategies, plugin boundaries | Accepted architecture |
| `docs/adr/018-block-catalog-data-action-reactivity.md` | ActionSource, ActionRunner, effects, lifecycle, Portable Expression AST, RuntimeContext wiring | Accepted architecture |
| `sprints/package-ui-elementor/block-catalog-candidate-matrix.json` | Machine-readable candidate baseline | Accepted matrix baseline |
| `sprints/package-ui-elementor/block-catalog-parity-matrix.md` | Human-readable mirror of the matrix | Report, must match JSON |
| `sprints/package-ui-elementor/REVIEW-FINDINGS.md` | Review evidence and resolved gap notes | Historical evidence |
| `docs/components-progress.md` | Component progress tracker | Human-facing status summary |

When these artifacts disagree, use this order:

1. Accepted ADR text.
2. Matrix JSON baseline plus documented change records.
3. Generated parity report.
4. Progress docs and handoffs.

## 3. Repository target map

Exact final filenames may evolve to match existing package conventions, but every deliverable must land in the correct package boundary.

### 3.1 `packages/ui-schema`

Canonical, framework-neutral contracts live here.

Expected areas:

- `packages/ui-schema/src/block-catalog/`
  - Block catalog entry types.
  - Candidate entry types.
  - Category, channel, capability, readiness, and render strategy enums.
  - Adaptive, security, cache, offline, migration, i18n, a11y evidence, and readiness report metadata.
  - JSON-schema or runtime validation helpers if this package already owns schema validation.
- `packages/ui-schema/src/actions/`
  - `ActionSource`, `ActionAdapter`, `ActionResult`, `ActionError`, `ActionEffect`, `ActionLifecycleEvent`.
  - Action manifest and adapter manifest types.
- `packages/ui-schema/src/expressions/`
  - Portable Expression AST node types.
  - Built-in function manifest types.
  - Shared golden fixture type definitions.
- `packages/ui-schema/src/runtime/`
  - `RuntimeContext`, `ActionRuntimeContext`, adapter registries, permission-filtered scope shape.
- `packages/ui-schema/src/index.ts`
  - Public exports.

Do not put React, web component, DOM, Flutter, engine, or provider SDK dependencies in `ui-schema`.

### 3.2 `packages/ui`

Web implementation, catalog source, validators, generated indexes, demo integration.

Expected areas:

- `packages/ui/src/block-catalog/`
  - Per-category catalog sources, for example `content.ts`, `data.ts`, `interaction.ts`, `workflow.ts`, `commerce.ts`, `analytics.ts`, `navigation.ts`, `builder.ts`.
  - Ready index, candidate index if needed for tooling, and public block catalog API.
  - `BcBlockRenderer` implementation for schema and component strategy.
  - Web block component registry.
- `packages/ui/src/core/actions/`
  - TypeScript `ActionRunner`.
  - Local/url/model/mock action adapters.
  - Lifecycle event bridge.
  - Effect runner.
- `packages/ui/src/core/expressions/`
  - TypeScript Portable Expression AST evaluator.
  - Built-in functions.
  - Function registry.
- `packages/ui/src/core/runtime/`
  - RuntimeContext validators and permission-filtered scope construction.
- `packages/ui/scripts/`
  - Catalog validation CLI.
  - Matrix/report drift checker.
  - Source/export checker for web and Flutter if implemented from Node.
  - Readiness report generator.
- `packages/ui/src/index.ts`
  - Public exports, including fixed `BcShare`, `BcGridView`, catalog APIs, renderer, action/runtime APIs as applicable.
- `packages/ui/demo-app/`
  - Block Catalog demo section for web.

Do not make builder tooling a dependency of simple runtime block rendering.

### 3.3 `packages/ui-flutter`

Dart and Flutter parity implementation.

Expected areas:

- `packages/ui-flutter/lib/src/block_catalog/`
  - Dart catalog types if not generated from schema.
  - Flutter `BcBlockRenderer` equivalent.
  - Flutter block component registry.
- `packages/ui-flutter/lib/src/actions/`
  - Dart `ActionRunner` equivalent.
  - Lifecycle notifier/stream parity.
  - Local/url/model/mock adapters.
  - Effect runner.
- `packages/ui-flutter/lib/src/expressions/`
  - Dart Portable Expression AST evaluator.
  - Built-in function registry.
- `packages/ui-flutter/lib/src/runtime/`
  - RuntimeContext validators and permission-filtered scope construction.
- `packages/ui-flutter/lib/src/components/`
  - Canonical missing widgets: `BcAvatar`, `BcBadge`, `BcProgress`, `BcKpi`.
- `packages/ui-flutter/lib/bc_ui_flutter.dart`
  - Public exports.
- `packages/ui-flutter/example/`
  - Block Catalog demo section for Flutter.
- `packages/ui-flutter/test/`
  - Widget tests, action tests, expression golden tests.

Flutter `BcRenderer` is infrastructure and must be validated as renderer infrastructure, not by looking for a component file named `bc_renderer.dart` in component folders only.

### 3.4 `docs` and `sprints`

- `docs/components-progress.md`: update status after each implementation batch.
- `sprints/package-ui-elementor/block-catalog-candidate-matrix.json`: update only with version bump or documented change record when candidate status changes.
- `sprints/package-ui-elementor/block-catalog-parity-matrix.md`: regenerate or update from JSON after every matrix change.
- `sprints/package-ui-elementor/REVIEW-FINDINGS.md`: only update if a finding status materially changes.

## 4. Core data model contract

Implementation may refine syntax, but it must preserve these responsibilities.

### 4.1 `BlockCatalogEntry`

Required fields for Ready entries:

| Field | Requirement |
|---|---|
| `id` | Namespaced stable slug, for example `content.hero` |
| `version` | Semver-like block version |
| `schemaVersion` | Catalog schema version |
| `labelKey` | i18n key |
| `fallbackLabel` | English fallback string |
| `category` | One of the ADR-017 categories |
| `channels` | Non-empty array of channels |
| `capabilities` | Capability tags, each mapped to validators if behavioral |
| `renderStrategy` | `schema` or `component` |
| `requiredComponents` | Explicit component keys for web and Flutter validation |
| `requiredFunctions` | Explicit built-in or host function requirements |
| `requiredAdapters` | Explicit adapter requirements |
| `adaptivePolicy` | Declared or inherited with proof |
| `securityPolicy` | Declared for all entries, stricter for URL, file, rich text, payment, provider, action |
| `cachePolicy` | Declared, even if `none` |
| `offlinePolicy` | Declared, even if `unsupported` |
| `migration` | Present for versioning and compatibility |
| `readiness` | `ready`, `candidate`, or `deprecated` |
| `a11yEvidence` | Link to catalog metadata, generated readiness report, or demo test reference |
| `tests` | References to schema, web, Flutter, action, a11y tests as applicable |

Optional fields:

- `subcategory`.
- `schemaPreset` when `renderStrategy = schema`.
- `componentKey` when `renderStrategy = component`.
- `aliases` if retained in Ready metadata.
- `deprecatedBy` for deprecated entries.

### 4.2 `CandidateCatalogEntry`

Candidate-only fields may include:

- `sourceReferences`.
- `decision`: `keep`, `merge`, `split`, `alias`, `exclude`.
- `matrixStatus`.
- `gaps`.
- `priorityBatch`.
- `blockedBy`.
- `readyPromotionChecklist`.
- `changeRecord`.

Runtime must not consume Candidate entries by default.

### 4.3 Capabilities and validators

Every behavioral capability needs a validator. Minimum mapping:

| Capability | Required validation |
|---|---|
| `data-source` | Data adapter exists, source shape valid, loading/empty/error states declared |
| `action-source` | ADR-018 ActionSource valid, ActionRunner available, mock adapter tests pass |
| `reactive` | Portable Expression AST or approved built-in function use |
| `permission-aware` | Permission adapter exists, scope filtered before expression evaluation |
| `i18n` | `labelKey` and `fallbackLabel`, locale fallback behavior |
| `a11y` | Contract and evidence link |
| `responsive` | Breakpoints or inherited adaptive policy |
| `adapt` | Platform adaptation policy, including iOS/macOS and Android behavior where applicable |
| `provider-adapter` | Generic adapter contract, no provider SDK in core |
| `payment-adapter` | Generic payment adapter plus mock/local adapter, no secrets in schema |
| `large-data` | Pagination, virtualization, bounded scrolling, debounce/abort policy |
| `offline` | Queue persistence, retry/backoff, idempotency key, replay, conflict, permission revalidation, visible recovery states |
| `sanitize` | Sanitizer defaults and opt-out policy |
| `security` | URL/file/rich text/redaction/provider policy |

Capability alias rule: the matrix uses `action-trigger` for some pre-ADR candidates. Implementation must normalize `action-trigger` to the canonical ADR capability `action-source` before Ready validation.

## 5. Runtime contract implementation details

### 5.1 `BcBlockRenderer`

`BcBlockRenderer` must stay a thin orchestration layer.

Required behavior:

1. Resolve catalog entry by `id` and optional `version`.
2. Reject non-Ready entries unless caller explicitly opts into Candidate tooling mode.
3. Validate `RuntimeContext` against required components, functions, adapters, permission capability, and i18n entries.
4. For `schema` strategy, materialize serializable `UiSchema` and delegate to `BcRenderer`.
5. For `component` strategy, resolve `componentKey` from platform registry.
6. Emit typed diagnostics for missing or incompatible requirements.
7. Never silently downgrade to partial rendering unless an explicit fallback is declared.

Typed diagnostic examples:

- `BLOCK_ENTRY_NOT_FOUND`.
- `BLOCK_ENTRY_NOT_READY`.
- `BLOCK_COMPONENT_MISSING`.
- `BLOCK_COMPONENT_NOT_EXPORTED`.
- `BLOCK_FUNCTION_MISSING`.
- `BLOCK_ADAPTER_MISSING`.
- `BLOCK_PERMISSION_SCOPE_INVALID`.
- `BLOCK_I18N_MISSING`.
- `BLOCK_SCHEMA_INVALID`.
- `BLOCK_RENDER_STRATEGY_INVALID`.
- `BLOCK_PLATFORM_UNSUPPORTED`.

### 5.2 Platform component registries

Registry keys are canonical component keys such as `BcGridView`, not file paths.

Rules:

- Web registry validates source and public export from `packages/ui/src/index.ts` or approved public entrypoint.
- Flutter registry validates source and public export from `packages/ui-flutter/lib/bc_ui_flutter.dart`.
- Registry must support manual registration for application/plugin components.
- Registry must fail loudly on duplicate keys unless replacement is explicitly allowed in test/demo mode.
- Registry must not import all builder surfaces into simple runtime bundles unless tree-shaking or lazy registration is proven.

### 5.3 RuntimeContext and permission-filtered scopes

Construction rules:

1. Host creates a per-instance context.
2. Block runtime derives a permission-filtered expression scope before evaluating sensitive payloads.
3. Unauthorized values, data, services, and permission facts are removed, not just hidden in UI.
4. Action adapters receive only normalized, allowed payloads.
5. Logs and lifecycle events redact sensitive fields.

Minimum scope keys:

- `values`.
- `meta`.
- `data`.
- `services`.
- `permissions`.
- `functions`.

## 6. ADR-018 implementation details

### 6.1 Portable Expression AST

Minimum node support:

- `literal`.
- `path`.
- `op`.
- `call`.
- `if`.

Minimum operators:

- Comparisons: `eq`, `ne`, `gt`, `gte`, `lt`, `lte`.
- Boolean: `and`, `or`, `not`.
- Arithmetic: `add`, `sub`, `mul`, `div`.
- Aggregation: `sum`, `count`, `min`, `max`, `avg`.
- String/list: `concat`, `contains`, `startsWith`, `endsWith`.
- Empty/null: `isEmpty`, `isNotEmpty`, `coalesce`.
- Conditional: `if`.

Golden fixtures must run against both TypeScript and Dart. Include at least:

- Scalar comparisons.
- Nested path resolution.
- Missing path typed error.
- Null propagation.
- Array aggregation.
- String/list containment.
- Date formatting with locale/timezone behavior.
- Built-in function success.
- Missing function typed error.
- Permission-filtered scope denial.
- Invalid node shape.
- Deterministic error path to node or field.

String expressions remain escape hatches and cannot be a Ready foundation unless evaluator parity is proven and declared.

### 6.2 Built-in function registry

Minimum namespaces:

- `math.*`.
- `string.*`.
- `array.*`.
- `date.*`.
- `format.*`.
- `validate.*`.

Every built-in requires a manifest:

- Name.
- Version.
- Supported platforms.
- Pure flag.
- Determinism.
- Input schema if available.
- Output schema if available.
- Fallback behavior.

Ready universal entries may depend only on built-ins that exist on all required platforms.

### 6.3 ActionRunner

Execution order must match ADR-018:

1. Validate action source and required context.
2. Validate permission requirements and construct a permission-filtered expression scope.
3. Evaluate Portable Expression AST payloads against the filtered scope.
4. Emit start lifecycle.
5. Invoke local/url/model adapter or effect runner.
6. Normalize result.
7. Apply validation errors to store meta.
8. Apply result effects in deterministic order.
9. Emit success, validation-error, or error lifecycle.
10. Emit settled exactly once.
11. Return normalized `ActionResult`.

Effect order:

1. `setValue`, `setMeta`, `setError`.
2. `invalidateKey`, `refreshDataSource`.
3. `toast`, `closeOverlay`, `resetForm`.
4. `navigate`.
5. `emit`.

Required edge cases:

- Cancellation and abort.
- Timeout.
- Repeated submission.
- Validation error normalization.
- Adapter throws raw error.
- Adapter returns partial malformed response.
- Missing adapter.
- Missing function.
- Missing permission adapter for permission-aware action.
- Sensitive field redaction.
- Offline-capable replay with permission revalidation.

### 6.4 Offline-capable action blocks

If a block declares offline support for actions, it must prove:

- Durable queue persistence.
- Retry/backoff policy.
- Idempotency key behavior.
- Replay policy.
- Conflict normalization.
- Permission revalidation on replay.
- Visible failure and recovery states.
- Tests for enqueue, replay success, replay conflict, replay permission denial, and user-visible failure.

No action-capable block may claim offline support with only a cache flag.

## 7. Validation tooling contract

Add a catalog validator before any Ready promotion. Suggested command name:

```powershell
Set-Location 'D:\NodeJS\lowcode\packages\ui'; npm run catalog:validate
```

If script naming changes, update this plan and `docs/components-progress.md`.

### 7.1 Inputs

- Candidate matrix JSON.
- Per-category catalog source files.
- Generated Ready index.
- Web source tree and public export file.
- Flutter source tree and public export file.
- Function manifests.
- Adapter manifests.
- i18n catalog.
- Test/readiness report metadata.

### 7.2 Outputs

- Exit code `0` when all Ready entries pass.
- Non-zero exit when any Ready gate fails.
- Machine-readable diagnostics, ideally JSON, with:
  - `code`.
  - `severity`.
  - `entryId`.
  - `fieldPath`.
  - `message`.
  - `suggestedFix`.
  - `sourceFile`.
- Human-readable summary by gate.

### 7.3 Required checks

| Check | Required before Ready |
|---|---|
| Matrix status count | Yes |
| Matrix baseline version/change record | Yes |
| Entry schema validation | Yes |
| Candidate vs Ready consistency | Yes |
| Web source exists | Yes |
| Web public export exists | Yes |
| Flutter source exists | Yes |
| Flutter public export exists | Yes |
| Renderer infrastructure export check | Yes, separate from component source check |
| Required functions manifest | Yes |
| Required adapters manifest | Yes |
| Permission-filtered scope rule | Yes |
| Security policy | Yes |
| URL allowlist policy | If URL/provider/media action capability |
| Rich text sanitizer policy | If rich text/HTML/markdown capability |
| File policy | If file/download/upload/import capability |
| Payment policy | If commerce payment capability |
| i18n key and English fallback | Yes |
| A11y contract/evidence link | Yes |
| Adaptive policy | Yes |
| Performance policy | If data-heavy/large-data |
| Cache/offline policy | Yes, even if unsupported |
| Version and migration metadata | Yes |
| No-stub scan | Yes |
| Demo reference | Yes |
| Test reference | Yes |

## 8. Batch execution overview

This plan uses dependency-first ordering. It refines ADR-017 implementation phases by making Foundation and ERP blockers explicit before broad content work.

| Batch | Name | Purpose | Ready promotions allowed? |
|---|---|---|---|
| P0 | Baseline hygiene | Validate accepted matrix/docs and avoid drift | No |
| P1 | Schema and catalog types | Add canonical types in `ui-schema` | No |
| P2 | Expression and function runtime | Add AST evaluators and built-ins in TS/Dart | No, except tooling demos |
| P3 | Action runtime | Add ActionSource, ActionRunner, effects, lifecycle, adapters | No, except demo-only entries |
| P4 | Block renderer and validators | Add `BcBlockRenderer`, registries, validator CLI | Limited pilot Ready entries |
| P5 | Export and parity blockers | Fix `BcShare`, `BcGridView`, Flutter `BcAvatar/BcBadge/BcProgress/BcKpi`, web `BcMarkdown` | Yes, for unblocked entries |
| P6 | ERP data and interaction | Promote ERP data/interaction entries | Yes |
| P7 | Workflow | Promote workflow entries excluding editor until builder batch | Yes |
| P8 | Analytics | Promote KPI/chart, pivot if component exists | Yes |
| P9 | Commerce | Promote generic commerce entries, isolate checkout complexity | Yes |
| P10 | Content and navigation | Promote content/website and navigation entries | Yes |
| P11 | Builder surface | Implement builder components without coupling simple runtime | Yes, only if gates pass |
| P12 | Plugin extensibility | Formalize plugin catalog and provider adapters | Core Ready no, plugin Ready scoped |

## 9. Detailed batch plan

### P0: Baseline hygiene

Purpose: Ensure implementators start from a stable accepted baseline.

Tasks:

1. Confirm ADR-017 and ADR-018 status is Accepted.
2. Confirm matrix JSON status is `accepted-baseline`, schema version `0.1.0`, candidate count `44`.
3. Confirm parity matrix matches JSON status counts and candidate IDs.
4. Confirm no Ready runtime catalog exists yet, or if it exists, mark it experimental until validated.
5. Decide generated file policy:
   - Commit per-category source files.
   - Commit generated Ready index if package convention already commits generated indexes.
   - Commit generated readiness report if it is the a11y/test evidence source.
6. Add an implementation checklist issue or tracking file only if the team wants per-batch tracking. Do not change ADR content for implementation progress.

Exit gate:

- Matrix parse check passes.
- Candidate ID set is stable.
- No code changes beyond optional scripts/tests in P0.

### P1: Schema and catalog types

Purpose: Create the portable contract surface.

Deliverables:

1. `BlockCatalogEntry`.
2. `CandidateCatalogEntry`.
3. `BlockCategory`, `BlockChannel`, `BlockCapability`, `BlockReadiness`, `BlockRenderStrategy`.
4. `BlockAdaptivePolicy`, `BlockSecurityPolicy`, `BlockCachePolicy`, `BlockOfflinePolicy`.
5. `BlockMigrationMetadata`.
6. `BlockA11yEvidence` and readiness report reference types.
7. `RequiredComponent`, `RequiredFunction`, `RequiredAdapter` declarations if simple strings are insufficient.
8. Export surface from `@bitcode/ui-schema`.
9. Minimal type tests or compile tests.

Order:

1. Add types in `ui-schema`.
2. Export from package entrypoints.
3. Add examples as test fixtures, not Ready catalog claims.
4. Run TypeScript verification.

Exit gate:

- `@bitcode/ui-schema` type surface compiles.
- No DOM/web/Flutter/engine/provider dependency introduced.
- Example entries validate structurally.

### P2: Portable Expression AST and built-in functions

Purpose: Establish cross-platform deterministic expression behavior.

Deliverables:

1. TypeScript evaluator.
2. Dart evaluator.
3. Shared JSON golden fixtures.
4. Built-in function registry and manifests.
5. Permission-filtered scope helper contract.
6. Typed expression errors.

Order:

1. Add schema/types from P1 if not complete.
2. Implement TypeScript evaluator.
3. Implement Dart evaluator.
4. Add shared fixtures.
5. Add manifest validation.
6. Add permission-filtered scope tests.

Exit gate:

- Golden fixtures pass in TypeScript and Dart.
- Missing paths and missing functions produce typed errors.
- Unauthorized fields are absent before evaluation.
- No JavaScript function bodies are required for Ready behavior.

### P3: ActionSource, ActionRunner, effects, lifecycle

Purpose: Implement mutation contract symmetrically to DataSource.

Deliverables:

1. `ActionSource` variants: `local`, `handler`, `url`, `model`.
2. `ActionResult`, `ActionError`, `ActionEffect`, `ActionLifecycleEvent`.
3. TypeScript `ActionRunner`.
4. Dart `ActionRunner` equivalent.
5. Local/url/model mock adapters.
6. Deterministic effect runner.
7. Web lifecycle events and Dart notifier/stream parity.
8. Redaction utilities.
9. Timeout/abort/cancellation behavior.
10. Offline-capable action policy tests for any offline-supporting action demo.

Order:

1. Add action types to `ui-schema`.
2. Implement local effects first.
3. Implement mock/model/url adapters.
4. Implement lifecycle emission.
5. Integrate with `ReactiveStore.runAction` in TypeScript.
6. Integrate equivalent Dart API.
7. Add parity tests.

Exit gate:

- Effect ordering test passes.
- Lifecycle order test passes.
- Validation error test passes.
- Timeout/abort test passes.
- Missing adapter/function diagnostics are typed.
- Sensitive payload fields are redacted in logs/events.

### P4: Block renderer, registries, and validators

Purpose: Connect catalog entries to rendering and enforce Ready gates.

Deliverables:

1. `BcBlockRenderer` web implementation.
2. `BcBlockRenderer` Flutter implementation or equivalent.
3. Web component registry.
4. Flutter component registry.
5. Build-time catalog validator CLI.
6. Runtime host validator.
7. Generated readiness report format.
8. At least one schema-strategy demo entry and one component-strategy demo entry in Candidate or pilot Ready mode.

Order:

1. Build registry APIs.
2. Add renderer strategy switch.
3. Add runtime validation.
4. Add build-time validator.
5. Add pilot entries.
6. Add demo routes/sections.

Exit gate:

- Missing component, adapter, function, permission, and i18n diagnostics are tested.
- Schema strategy delegates to `BcRenderer`.
- Component strategy resolves through typed registry.
- Validator blocks incomplete Ready entries.

### P5: Export and parity blockers

Purpose: Resolve known source/export and platform parity blockers before broad promotions.

Deliverables:

1. Public web export/registry for `BcShare`.
2. Public web export/registry for `BcGridView`.
3. Canonical Flutter `BcAvatar` source, tests, and barrel export.
4. Canonical Flutter `BcBadge` source, tests, and barrel export.
5. Canonical Flutter `BcProgress` source, tests, and barrel export.
6. Canonical Flutter `BcKpi` source, tests, and barrel export. It may wrap or reuse `BcStatistic` internally, but the public widget key must exist.
7. Canonical web `BcMarkdown` source, tests, and public export.
8. Validator checks for all above.

Order:

1. Fix web exports that already have source: `BcShare`, `BcGridView`.
2. Add Flutter widgets in dependency order: `BcAvatar`, `BcBadge`, `BcProgress`, `BcKpi`.
3. Add web `BcMarkdown`.
4. Add tests and demos for each.
5. Re-run catalog validator.

Exit gate:

- Public exports pass for web and Flutter.
- No candidate remains blocked only by these known export/parity issues.
- Components have tests and demo coverage.

### P6: ERP data and interaction

Candidate scope:

| Candidate | Initial matrix status | Promotion target | Critical blockers |
|---|---|---|---|
| `interaction.form` | existing-component-needs-binding, foundation batch | Ready | ActionSource, AST, Flutter renderer export validation |
| `data.record-list` | existing-component-needs-binding | Ready | DataSource conventions, BcGridView export, performance budget |
| `data.record-detail` | existing-component-needs-binding | Ready | Permission adapter, detail schema preset |
| `data.account-summary` | schema-preset-only | Ready | Flutter BcAvatar |
| `data.download-center` | existing-component-needs-binding | Ready | Safe download action, file validation policy |
| `interaction.multi-step-wizard` | existing-component-needs-binding | Ready | Wizard schema/action contract, ADR-014 engine wiring deferred |
| `interaction.record-picker` | contract-gap | Candidate or Ready after contract | Unified picker contract, quick-create action, snapshot convention |
| `interaction.import-wizard` | component-gap | Candidate until component exists | Parser adapter, preview grid, file policy |
| `interaction.notification-center` | existing-component-needs-binding | Ready only for generic contract | Realtime/action contract, endpoint/model wiring deferred to ADR-014 |

Execution notes:

- `interaction.form` belongs to matrix `foundation` but production promotion is validated in this ERP batch after P2/P3/P4.
- `interaction.import-wizard` must not be forgotten. It remains Candidate until `BcImportWizard` or equivalent component exists.
- Notification Center can be Ready for generic catalog contract only. Engine realtime wiring is ADR-014.

Exit gate:

- At least one full ERP-like demo renders in web and Flutter with mock data/action adapters.
- Data-heavy entries declare pagination, virtualization, bounded scrolling, debounce/abort behavior.
- Permission-aware entries validate filtered scopes.

### P7: Workflow

Candidate scope:

| Candidate | Initial status | Promotion target | Notes |
|---|---|---|---|
| `data.audit-log` | existing-component-needs-binding | Ready for generic activity contract | Endpoint/model wiring stays ADR-014 |
| `workflow.approval-actions` | component-gap | Candidate until component exists, then Ready | Requires ActionSource lifecycle |
| `workflow.state-indicator` | contract-gap | Ready after Flutter BcBadge/BcProgress | Dedicated state component only if semantics require it |
| `workflow.diagram-viewer` | schema-preset-only | Ready | Read-only process graph preset |
| `workflow.diagram-editor` | component-gap, builder-surface | Candidate until P11 | Authoring component, persistence, undo/redo |
| `workflow.chatter` | existing-component-needs-binding | Ready for generic comment/message contract | Endpoint/model/permission/realtime/persistence stays ADR-014 |

Exit gate:

- Viewer side is Ready separately from editor.
- Approval/chatter actions use ADR-018 lifecycle and effects.
- ADR-014-owned concerns are not implemented or claimed as part of core catalog readiness.

### P8: Analytics

Candidate scope:

| Candidate | Initial status | Promotion target | Notes |
|---|---|---|---|
| `analytics.kpi-grid` | schema-preset-only | Ready | Depends on canonical Flutter BcKpi |
| `analytics.chart` | existing-component-needs-binding | Ready | Uses ADR-015 ECharts option contract and Flutter BcProgress |
| `data.pivot-table` | component-gap | Candidate until pivot component/contract exists | Aggregation endpoint/engine wiring remains ADR-014 |

Exit gate:

- KPI and chart demos exist in web and Flutter.
- Dashboard performance budget declared.
- Pivot remains explicitly Candidate if no dedicated component/contract is delivered.

### P9: Commerce

Candidate scope:

| Candidate | Initial status | Promotion target | Notes |
|---|---|---|---|
| `commerce.product-grid` | existing-component-needs-binding | Ready | Commerce variant of `data.record-list`, BcGridView export required |
| `commerce.product-detail` | existing-component-needs-binding | Ready | Commerce variant of `data.record-detail`, cart ActionSource |
| `commerce.cart-checkout` | component-gap | Candidate until dedicated flow exists | Payment adapter, security review, likely own batch if complex |
| `commerce.payment-action` | adapter-gap | Ready for generic contract plus mock/local adapter only | Provider SDKs are plugins, not core |
| `commerce.review` | existing-component-needs-binding | Ready | Generic datasource/action shape |

Exit gate:

- No provider secret or provider SDK in core.
- Payment actions use generic payment adapter contract plus mock/local adapter.
- Cart/checkout is not promoted until dedicated UX, security, and platform tests pass.

### P10: Content and navigation

Content scope:

| Candidate | Initial status | Promotion target | Critical blockers |
|---|---|---|---|
| `content.hero` | schema-preset-only | Ready | Schema preset and Flutter demo |
| `content.cta` | existing-component-needs-binding | Ready | ActionSource trigger support |
| `content.rich-text` | existing-component-needs-binding | Ready | Web BcMarkdown, sanitize policy |
| `content.plain-text` | schema-preset-only | Ready | Text primitive/schema support |
| `content.media` | existing-component-needs-binding | Ready | Media URL safety, provider metadata |
| `content.gallery` | existing-component-needs-binding | Ready | Datasource shape |
| `content.faq` | schema-preset-only | Ready | Heading hierarchy and i18n |
| `content.testimonial` | schema-preset-only | Ready | Flutter BcAvatar |
| `content.team-member` | schema-preset-only | Ready | Flutter BcAvatar |
| `content.pricing-table` | existing-component-needs-binding | Ready | Plan-selection ActionSource |
| `content.countdown` | component-gap | Candidate until `BcCountdown` or portable timer block exists | Time/reactive component |
| `content.social-links` | existing-component-needs-binding | Ready | BcShare export, URL safety, provider icon strategy |
| `data.map-view` | existing-component-needs-binding | Ready for generic map adapter | Google provider remains plugin |

Navigation scope:

| Candidate | Initial status | Promotion target | Notes |
|---|---|---|---|
| `navigation.menu` | existing-component-needs-binding | Ready | Responsive/adapt policies per menu variant |
| `navigation.search-toc-sitemap` | existing-component-needs-binding | Ready | Search datasource and CMS navigation schemas |

Exit gate:

- Static website block demo complete in web and Flutter.
- External links and provider URLs pass safety policies.
- Rich text sanitizer defaults on.

### P11: Builder surface

Candidate scope:

| Candidate | Initial status | Promotion target | Notes |
|---|---|---|---|
| `builder.expression-builder` | component-gap, foundation batch | Candidate/Ready only after visual builder exists | P2 provides AST foundations, P11 provides visual editor |
| `builder.query-builder` | component-gap | Candidate/Ready after dedicated component | Builder surface, not runtime dependency |
| `builder.page-builder` | component-gap | Candidate/Ready after dedicated component | Consumes Block Catalog runtime |
| `builder.report-builder` | component-gap | Candidate/Ready after dedicated component | Export/report engine wiring remains ADR-014 |
| `workflow.diagram-editor` | component-gap | Candidate/Ready after dedicated authoring component | Reopened here from workflow |

Execution notes:

- Split `builder.expression-builder` into two deliverables:
  - Foundation AST/evaluator work in P2.
  - Visual expression builder component in P11.
- Builder packages may depend on Block Catalog runtime. Block Catalog runtime must not depend on builder packages.

Exit gate:

- Builder surfaces have undo/versioning, validation, persistence strategy, and security model.
- Simple schema/component rendering works without importing builder surfaces.

### P12: Plugin extensibility

Scope:

| Candidate | Matrix status | Handling |
|---|---|---|
| `plugin.embed-social` | duplicate-alias | Alias to generic provider/embed capability, not core Ready primitive |
| Provider maps | plugin/provider | Host-injected plugin entries over generic map adapter |
| Provider payments | plugin/provider | Host-injected plugin entries over generic payment adapter |
| `excluded.wordpress-specific` | excluded | Remains excluded with reason |

Rules:

- Plugin catalog entries have plugin-scoped readiness, not core Ready status.
- Core runtime ignores plugin entries unless host injects them.
- Plugin entries must declare provider ID, supported platforms, adapters, security policy, URL allowlist, fallback behavior, and conformance fixtures.
- Provider secrets stay in host/plugin configuration, never schema.

Exit gate:

- Generic plugin registration mechanism works.
- Provider examples can be demonstrated with mock/local adapters or safe public URLs.
- Core Ready Catalog remains provider-neutral.

## 10. Candidate promotion checklist

A candidate can be promoted to Ready only if all items below are checked.

### Metadata

- [ ] Candidate decision is `keep`, `merge`, or approved `split` with target entry.
- [ ] Alias/exclusion is documented if not promoted.
- [ ] `id`, `version`, and `schemaVersion` are present.
- [ ] `labelKey` and `fallbackLabel` are present.
- [ ] Category, subcategory if needed, channels, capabilities are present.
- [ ] `renderStrategy` is explicit.
- [ ] `requiredComponents`, `requiredFunctions`, and `requiredAdapters` are complete.
- [ ] `migration` metadata exists.

### Web

- [ ] Every required web component source exists.
- [ ] Every required web component is publicly exported.
- [ ] Web registry has correct key.
- [ ] Web render test passes.
- [ ] Web demo is present.

### Flutter

- [ ] Every required Flutter widget source exists.
- [ ] Every required Flutter widget is exported from `bc_ui_flutter.dart`.
- [ ] Flutter registry has correct key.
- [ ] Flutter widget test passes.
- [ ] Flutter demo is present.
- [ ] Platform-aware behavior follows ADR-011.

### Data/action/reactivity

- [ ] DataSource shape is declared when relevant.
- [ ] ActionSource shape is declared when relevant.
- [ ] ActionRunner tests cover success, validation error, generic error, cancellation/timeout if relevant.
- [ ] Effects are deterministic.
- [ ] Lifecycle parity exists on web and Flutter.
- [ ] Portable Expression AST is used for required reactivity.
- [ ] Built-in or host functions have manifests.
- [ ] Permission-filtered scope is built before sensitive expression evaluation.

### Security, privacy, and reliability

- [ ] No secrets in schema or catalog metadata.
- [ ] URL allowlist policy exists when needed.
- [ ] Rich text sanitizer policy exists when needed.
- [ ] File policy exists when needed.
- [ ] Payment policy exists when needed.
- [ ] Logs/events redact sensitive payload fields.
- [ ] Permission/authz adapter requirement exists for permission-aware blocks.
- [ ] Cache policy exists.
- [ ] Offline policy exists, even if unsupported.
- [ ] Offline-capable actions meet ADR-018 queue/replay/conflict requirements.

### UX quality

- [ ] Loading, empty, error, and disabled states exist where relevant.
- [ ] A11y contract exists.
- [ ] A11y evidence is linked.
- [ ] Adaptive policy exists.
- [ ] Performance policy exists for data-heavy blocks.
- [ ] i18n fallback behavior is tested.
- [ ] Demo includes realistic edge cases, not only happy path.

## 11. Verification commands

Use package-local commands because the repo root may not have a workspace `package.json`.

### 11.1 Existing commands

```powershell
Set-Location 'D:\NodeJS\lowcode\packages\ui'; npm run tokens:check
Set-Location 'D:\NodeJS\lowcode\packages\ui'; npm run registry
Set-Location 'D:\NodeJS\lowcode\packages\ui'; npm run index
Set-Location 'D:\NodeJS\lowcode\packages\ui'; npm run audit:contract
Set-Location 'D:\NodeJS\lowcode\packages\ui'; npm run test
Set-Location 'D:\NodeJS\lowcode\packages\ui'; npm run build
```

```powershell
Set-Location 'D:\NodeJS\lowcode\packages\ui-flutter'; flutter analyze
Set-Location 'D:\NodeJS\lowcode\packages\ui-flutter'; flutter test
```

For `packages/ui-schema`, add a verification script if absent. Until then, use local TypeScript directly if dependencies are installed:

```powershell
Set-Location 'D:\NodeJS\lowcode\packages\ui-schema'; npx tsc --noEmit
```

### 11.2 New commands to add

Recommended scripts:

```json
{
  "catalog:validate": "node scripts/validate-block-catalog.mjs",
  "catalog:report": "node scripts/generate-block-catalog-report.mjs",
  "catalog:check": "npm run catalog:validate && npm run catalog:report -- --check",
  "test:expressions": "vitest run src/core/expressions",
  "test:actions": "vitest run src/core/actions"
}
```

Recommended Flutter targets:

```powershell
Set-Location 'D:\NodeJS\lowcode\packages\ui-flutter'; flutter test test/expressions test/actions test/block_catalog
```

### 11.3 Per-batch minimum verification

| Batch | Required verification |
|---|---|
| P0 | Matrix JSON parse, status counts, diff review |
| P1 | `ui-schema` typecheck, type export tests |
| P2 | TS expression tests, Dart expression tests, shared fixture parity |
| P3 | TS action tests, Dart action tests, lifecycle/effect order tests |
| P4 | Catalog validator, BcBlockRenderer tests, missing context diagnostics |
| P5 | Web build/test/audit, Flutter analyze/test, export checks |
| P6-P10 | Catalog validator, web test/build, Flutter analyze/test, demos |
| P11 | Builder component tests plus proof runtime does not import builder surfaces |
| P12 | Plugin conformance fixtures, provider-neutral core validator |

## 12. Rollback and compatibility rules

1. Catalog schema changes require `schemaVersion` handling.
2. Block behavior changes require `version` bump and migration metadata.
3. Removing a block requires deprecation metadata and optional `deprecatedBy`.
4. Alias changes require a change record.
5. A previously Ready block can be demoted only with documented reason, affected capabilities, and migration/fallback guidance.
6. Generated indexes must be reproducible. If generation is not deterministic, fix the generator before committing results.
7. Runtime errors must be typed and stable enough for builder diagnostics.

## 13. Known blockers from accepted baseline

| Blocker | Affected candidates | Resolution batch |
|---|---|---|
| `BcShare` source exists but lacks public export/registry | `content.social-links` | P5 |
| `BcGridView` source exists but lacks public export/registry | `data.record-list`, `commerce.product-grid` | P5 |
| Web `BcMarkdown` missing | `content.rich-text` | P5 |
| Flutter `BcAvatar` missing | `content.testimonial`, `content.team-member`, `data.account-summary` | P5 |
| Flutter `BcBadge` missing | `workflow.state-indicator` | P5 |
| Flutter `BcProgress` missing | `workflow.state-indicator`, `analytics.chart` | P5 |
| Flutter `BcKpi` missing | `analytics.kpi-grid` | P5 |
| ActionSource/ActionRunner missing | form, CTA, pricing, commerce, workflow actions, chatter, wizard, picker | P3 |
| Portable Expression AST missing | reactive and builder expression entries | P2 |
| BcBlockRenderer/registries missing | every Ready block | P4 |
| Validator missing | every Ready block | P4 |
| ADR-014 engine wiring pending | chatter, audit/activity, notification, wizard/import, pivot, report/export, page builder persistence | Outside this plan, do not claim in core Ready |

## 14. Risk register

| Risk | Severity | Mitigation |
|---|---:|---|
| Ready entries accidentally ship with stubs | High | Validator no-stub check, demo/test references mandatory |
| Builder surface couples into runtime bundle | High | Separate registries, lazy registration, dependency audit |
| Provider SDK leaks into core | High | Adapter boundary tests, dependency review, secret scan |
| Permission filtering occurs after expression evaluation | High | P2/P3 tests, validator rule for permission-aware actions |
| Offline support overclaimed | High | Offline policy shape plus required queue/replay/conflict tests |
| Flutter parity lags web | High | P5 before promotions, Flutter analyze/test in every promotion batch |
| Matrix/report drift | Medium | `catalog:report --check`, status count checks |
| Public export generation misses source components | Medium | `npm run index`, registry/export validator |
| ActionRunner becomes too broad | Medium | Only ADR-018 effect kinds in v1, new effect requires review |
| Performance regressions in large-data blocks | Medium | Pagination/virtualization/debounce/bounded scroll policies and tests |
| i18n/a11y evidence becomes paperwork only | Medium | Link evidence to test/readiness report, not prose-only claims |

## 15. Definition of Done

### 15.1 Batch Done

A batch is done only when:

- All planned files are implemented or explicitly deferred.
- Existing package tests pass.
- New tests cover new contract behavior.
- Catalog validator passes for any Ready entries touched.
- Matrix/report changes, if any, are synchronized.
- `docs/components-progress.md` is updated with real status.
- No unrelated files are staged.
- Diff is reviewed for secrets and provider credentials.

### 15.2 Ready Catalog Complete

The Ready Catalog is complete only when all 44 baseline candidates are one of:

1. Promoted to Ready with all gates passing.
2. Merged or aliased into another Ready entry with documented reason.
3. Excluded with documented `excludeReason`.
4. Deliberately left Candidate with accepted blocker and no runtime Ready claim.

For final completion, the parity report must show zero unaccounted candidates and no stale unresolved gaps.

## 16. First recommended implementation sequence

For the next coding session, start with the smallest foundation slice that creates enforceable structure without promoting any block:

1. Add `BlockCatalogEntry` and supporting enums/policies in `packages/ui-schema`.
2. Add minimal examples/tests for one `schema` entry and one `component` entry, both Candidate.
3. Add `catalog:validate` skeleton that validates schema shape, matrix status, and Ready no-stub rule.
4. Add fail-loud diagnostics shape.
5. Run `npx tsc --noEmit` in `packages/ui-schema` and `npm run test` in `packages/ui` if touched.
6. Do not mark any block Ready in this first slice.

This keeps the first implementation reversible and gives later batches a real gate.

## 17. References

- [ADR-017 Block Catalog Architecture](../adr/017-block-catalog-architecture.md)
- [ADR-018 Block Catalog Data, Action, and Reactivity Contract](../adr/018-block-catalog-data-action-reactivity.md)
- [ADR-014 Engine UI Feature Wiring](../adr/014-engine-ui-feature-wiring.md)
- [ADR-015 Flutter ECharts Chart Strategy](../adr/015-flutter-echarts-chart-strategy.md)
- [ADR-016 Demo App QA Bug Audit](../adr/016-demo-app-qa-bug-audit.md)
- Candidate matrix: `sprints/package-ui-elementor/block-catalog-candidate-matrix.json`
- Parity report: `sprints/package-ui-elementor/block-catalog-parity-matrix.md`
- Review findings: `sprints/package-ui-elementor/REVIEW-FINDINGS.md`
