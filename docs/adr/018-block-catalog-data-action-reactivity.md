# ADR-018: Block Catalog Data, Action, and Reactivity Contract

| Field | Value |
|-------|-------|
| **Status** | Accepted (2026-06-30; ADR-017 prerequisite for Ready entries; accepted with matrix 0.1 baseline) |
| **Depends on** | 002, 004, 005, 006, 007, 011, 012, 013, 017 |
| **Amends** | 004 (extends `ReactiveStore` with action execution and portable expression evaluation), 005 (extends `UiAction` beyond button descriptors), 007 (formalizes normalized action results and lifecycle events) |
| **Impact** | Adds `ActionSource`, `ActionAdapter`, `ActionResult`, `ActionError`, `ActionEffect`, `ActionRunner`, action lifecycle events, Portable Expression AST, built-in function registry, and RuntimeContext action wiring to `@bitcode/ui-schema` and Flutter ports |

---

## Context

ADR-017 introduces a Block Catalog where many entries are not static content. Forms submit, pricing tables select plans, product grids add to cart, record pickers quick-create records, import wizards parse files, workflow actions approve or reject records, chatter posts comments, notifications mark items read, and commerce checkout invokes payment adapters. These behaviors need a mutation contract as portable and explicit as ADR-002's `DataSource` read contract.

The current stack has important foundations but not the complete block runtime contract:

- ADR-002 defines 4-layer `DataSource` for reads.
- ADR-004 defines `ReactiveStore` and DATA / VALIDATION / UI / VALUE reactions, but default expressions are JS functions and string expressions are evaluator-dependent.
- ADR-005 defines `UiSchema` and `UiAction`, but `UiAction` is primarily a button/action descriptor and not a general mutation protocol.
- ADR-011 ports the renderer and reactive primitives to Flutter, so any Ready block contract must work in Dart without assuming JavaScript closures.
- ADR-012 makes `@bitcode/ui-schema` the canonical cross-platform type package.
- ADR-013 provides AppProvider/BcAppConfig, hooks, animation, adapt, and Flutter parity infrastructure.

The candidate matrix requires 5 foundational work items before any block can be Ready: block-catalog types, action-source, portable expression AST, block renderer, and validation tooling. ADR-017 owns the catalog and renderer. This ADR owns the data/action/reactivity contract that makes catalog entries executable without hardcoding BitCode engine assumptions.

## Scope

- Define `ActionSource`, `ActionAdapter`, `ActionResult`, `ActionError`, `ActionEffect`, and `ActionRunner` lifecycle.
- Define lifecycle events and Dart notifier parity.
- Define typed safe effects and deterministic effect ordering.
- Define Portable Expression AST for required Ready Catalog reactivity.
- Define a versioned built-in function registry and host function manifest.
- Define `RuntimeContext` wiring for actions, functions, adapters, permissions, i18n, navigation, and services.
- Extend `ReactiveStore` with `runAction` semantics in TypeScript and Dart.
- Define security and privacy gates for action-capable blocks.

## Out of Scope

- Block catalog entry schema and render strategy, owned by ADR-017.
- BitCode engine endpoint implementation and feature wiring, owned by ADR-014.
- Provider SDK implementations such as Stripe, PayPal, Google, Facebook, or SoundCloud. Core may define adapter contracts and mock/local adapters for tests only.
- Full visual builder components. This ADR defines expression/action contracts those builders may edit.

## Decision

### Section 1: ActionSource

Introduce `ActionSource` as the write/mutation mirror of `DataSource`.

```ts
type ActionSource =
  | LocalActionSource
  | HandlerActionSource
  | UrlActionSource
  | ModelActionSource;
```

#### `local`

`local` is a serializable action resolved entirely by the runtime. It is used for actions such as setting values, resetting forms, closing overlays, navigating through local routes, or emitting local events.

```ts
interface LocalActionSource {
  type: 'local';
  effects: ActionEffect[];
}
```

#### `handler`

`handler` is host-injected code. It is allowed in application code, but not in Ready Catalog JSON because it is not serializable or portable.

```ts
interface HandlerActionSource {
  type: 'handler';
  handlerKey: string;
}
```

Rules:

- Ready Catalog JSON cannot include function bodies.
- A Ready block may require a named handler only through `requiredAdapters` or `requiredFunctions`, with an explicit manifest and fallback policy.
- `handler` cannot satisfy universal Ready portability unless every required platform declares a conforming handler manifest and fallback policy.
- Missing handler registration is a runtime validation error.

#### `url`

`url` invokes an HTTP endpoint through a registered action adapter.

```ts
interface UrlActionSource {
  type: 'url';
  method: 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  url: string;
  headers?: Record<string, string>;
  body?: PortableExpression | Record<string, unknown>;
  timeoutMs?: number;
  idempotencyKey?: string;
}
```

Rules:

- GET-like reads remain `DataSource`, not `ActionSource`.
- URL actions require URL allowlist validation when used by Ready Catalog entries.
- Secrets, bearer tokens, and API keys must not be embedded in schema. Auth belongs to adapters.
- Abort and timeout behavior must be supported by the action adapter.

#### `model`

`model` delegates to a registered model adapter. It is general-purpose and not BitCode-specific.

```ts
interface ModelActionSource {
  type: 'model';
  model: string;
  operation: 'create' | 'update' | 'delete' | 'upsert' | 'restore' | 'custom';
  id?: PortableExpression;
  payload?: PortableExpression | Record<string, unknown>;
  adapter?: string;
}
```

Rules:

- The model adapter decides whether the backend is BitCode, REST, GraphQL, local SQLite, mock data, or another provider.
- Permission checks are adapter/runtime responsibilities, not hardcoded RBAC in catalog schema.
- `custom` operations must be named, documented, and declared in `requiredAdapters`.

### Section 2: ActionResult and ActionError

All action paths normalize to one result shape.

```ts
interface ActionResult<T = unknown> {
  ok: boolean;
  status: number | 'idle' | 'loading' | 'success' | 'error' | 'validation-error';
  data?: T;
  message?: string;
  errors?: Record<string, string | string[]>;
  errorList?: ActionError[];
  effects?: ActionEffect[];
  meta?: Record<string, unknown>;
}

interface ActionError {
  code: string;
  message: string;
  path?: string;
  field?: string;
  severity?: 'info' | 'warning' | 'error';
  recoverable?: boolean;
  details?: Record<string, unknown>;
}
```

Rules:

- `errors` is a convenience shape for form validation and field errors.
- `errorList` is the structured canonical list for machine handling.
- `ok = false` with validation errors must trigger validation lifecycle, not generic error lifecycle only.
- Adapters must never throw raw network/library errors past `ActionRunner`; they must normalize them into `ActionResult` or `ActionError`.
- Developer diagnostics may include adapter names and action IDs, but never secrets.

### Section 3: ActionLifecycle Events

Actions emit standardized lifecycle events.

Web Component event names:

- `bc-action-start`
- `bc-action-success`
- `bc-action-error`
- `bc-action-validation-error`
- `bc-action-settled`

Callback/notifier equivalents exist for React/Vue/RN and Dart.

Lifecycle payload:

```ts
interface ActionLifecycleEvent {
  actionId: string;
  source: ActionSource;
  phase: 'start' | 'success' | 'error' | 'validation-error' | 'settled';
  payload?: unknown;
  result?: ActionResult;
  error?: ActionError;
  meta: BcChangeMeta;
  timestamp: number;
}
```

Rules:

- `start` fires before adapter invocation.
- `success` fires after an `ok: true` result and before `settled`.
- `validation-error` fires for structured validation failures.
- `error` fires for non-validation failures.
- `settled` always fires exactly once after success, validation-error, error, or cancellation normalization.
- Dart uses typed notifiers/streams instead of DOM `CustomEvent`, but payload semantics must match.

### Section 4: ActionEffect

Effects are typed, safe, and intentionally small. No arbitrary string command execution is allowed.

Allowed effect kinds for v1:

- `setValue`
- `setMeta`
- `setError`
- `refreshDataSource`
- `invalidateKey`
- `toast`
- `navigate`
- `closeOverlay`
- `resetForm`
- `emit`

Conceptual type:

```ts
type ActionEffect =
  | { type: 'setValue'; path: string; value: PortableExpression | unknown }
  | { type: 'setMeta'; path: string; meta: Record<string, unknown> }
  | { type: 'setError'; path: string; error: string | ActionError }
  | { type: 'refreshDataSource'; key: string }
  | { type: 'invalidateKey'; key: string }
  | { type: 'toast'; level: 'info' | 'success' | 'warning' | 'error'; messageKey?: string; fallbackMessage: string }
  | { type: 'navigate'; to: string; replace?: boolean; params?: Record<string, unknown> }
  | { type: 'closeOverlay'; overlayId?: string }
  | { type: 'resetForm'; formId?: string; scope?: string }
  | { type: 'emit'; event: string; payload?: unknown };
```

Rules:

- Effects are declarative. They do not embed executable code.
- Effects can be produced by static action config, adapters, or `ActionResult.effects`.
- Effects run in deterministic order.
- Effects that require adapters, such as `navigate` or `toast`, fail loudly when the adapter is missing unless an explicit fallback exists.
- `emit` is limited to typed application events. It is not a shell for arbitrary commands.

### Section 5: Portable Expression AST

Ready Catalog required behavior must use a portable structured expression descriptor, not JavaScript functions or unbounded string expressions.

```ts
type PortableExpression =
  | { kind: 'literal'; value: unknown }
  | { kind: 'path'; path: string }
  | { kind: 'op'; op: PortableOperator; args: PortableExpression[] }
  | { kind: 'call'; fn: string; args: PortableExpression[] }
  | { kind: 'if'; test: PortableExpression; then: PortableExpression; else: PortableExpression };
```

Small required operator set:

- Comparisons: `eq`, `ne`, `gt`, `gte`, `lt`, `lte`
- Boolean: `and`, `or`, `not`
- Arithmetic: `add`, `sub`, `mul`, `div`
- Aggregation: `sum`, `count`, `min`, `max`, `avg`
- String/list: `concat`, `contains`, `startsWith`, `endsWith`
- Null/empty: `isEmpty`, `isNotEmpty`, `coalesce`
- Conditional: `if`

Rules:

- Evaluation scope is explicit and named: `values`, `meta`, `data`, `services`, `permissions`, and function arguments.
- The scope must be permission-filtered before expression evaluation; unauthorized fields, services, and permission facts must not be exposed to the evaluator.
- Path resolution must be deterministic and must not access global objects.
- Expressions are pure. Side effects belong to `ActionEffect` or `DataSource` refresh, not expression evaluation.
- Evaluation errors are typed and point to the expression node or path.
- String expressions remain an advanced/server-evaluated escape hatch. They are not the foundation for Ready Catalog behavior unless platform evaluator parity is proven and declared.

### Section 6: Built-in Function Registry

Portable expressions may call whitelisted pure built-ins.

Built-in namespaces:

- `math.*`
- `string.*`
- `array.*`
- `date.*`
- `format.*`
- `validate.*`

Function manifest shape:

```ts
interface FunctionManifestEntry {
  name: string;
  version: string;
  platforms: Array<'web' | 'flutter' | 'react-native' | 'server'>;
  pure: true;
  deterministic: boolean;
  inputSchema?: unknown;
  outputSchema?: unknown;
  fallback?: 'error' | 'literal' | 'empty' | 'host';
}
```

Rules:

- Ready universal entries may depend only on built-ins that exist on all required platforms.
- Host/project functions require explicit registration and a manifest.
- Runtime fails loudly if a required host function is not registered.
- Non-pure functions are adapters, not expression functions.
- Date/time functions must declare timezone behavior.
- Formatting functions must declare locale behavior and fallback.

### Section 7: RuntimeContext Wiring

ADR-017 defines `RuntimeContext` at the catalog level. This ADR defines action/runtime expectations inside it.

Required registries for actions:

```ts
interface ActionRuntimeContext {
  values: Record<string, unknown>;
  meta: Record<string, unknown>;
  data: Record<string, unknown>;
  functions: FunctionRegistry;
  adapters: {
    action: ActionAdapterRegistry;
    data?: DataAdapterRegistry;
    permission?: PermissionAdapter;
    navigation?: NavigationAdapter;
    toast?: ToastAdapter;
    i18n?: I18nAdapter;
  };
  services?: Record<string, unknown>;
}
```

Rules:

- Context is explicit and per-instance. No mutable singleton registry may be required.
- `ActionRunner` reads only from provided context and adapter registries.
- Permission and context validation happen before sensitive payload expression evaluation and before adapter invocation when an action declares permission requirements.
- Expression scopes must expose only authorized values, data, services, and permission facts.
- Adapters own auth and secrets.
- Host services must be registered by name and never inferred from global variables.

### Section 8: ReactiveStore.runAction Integration

Extend `ReactiveStore` with action execution as an orchestrated operation.

Conceptual API:

```ts
interface ReactiveStore {
  runAction(action: UiAction | ActionSource, payload: unknown, ctx: ActionRuntimeContext): Promise<ActionResult>;
}
```

Execution order:

1. Validate action source and required context.
2. Validate permission requirements and construct a permission-filtered expression scope.
3. Evaluate any Portable Expression AST payloads against the filtered scope.
4. Emit `bc-action-start` / Dart start notifier.
5. Invoke the selected action adapter or local effect runner.
6. Normalize result.
7. Apply validation errors to store meta when present.
8. Apply result effects in deterministic order.
9. Emit success, validation-error, or error lifecycle.
10. Emit settled lifecycle exactly once.
11. Return normalized `ActionResult`.

Effect ordering:

1. Field/meta mutation effects: `setValue`, `setMeta`, `setError`.
2. Data invalidation effects: `invalidateKey`, `refreshDataSource`.
3. UI effects: `toast`, `closeOverlay`, `resetForm`.
4. Navigation effects: `navigate`.
5. Event effects: `emit`.

Rationale:

- State updates happen before data refresh.
- Data refresh happens before navigation when possible.
- UI notifications can reflect final state.
- Navigation is late to avoid losing local state before effects complete.
- `emit` is last so observers see the final store snapshot.

### Section 9: Security and Privacy Gate

Action-capable blocks have a stricter security gate.

Rules:

- No secrets in schema, catalog entries, Portable Expression AST, or action config.
- Adapters own authentication and token injection.
- URL actions require allowlists and method restrictions.
- Rich text sanitization defaults on. Opt-out must be explicit, policy-gated, and auditable.
- File actions require type, size, extension, and scanning policy metadata.
- Payment actions require a payment adapter and must never place provider secret keys in schema. Core ships only the generic contract and mock/local adapter for tests; provider SDKs belong in plugin packages.
- Permission/authz adapter must gate permission-aware actions.
- Validation errors must not echo sensitive server details.
- Logs and lifecycle events must redact sensitive payload fields.
- Offline action queues must declare queue persistence, retry/backoff, idempotency key behavior, replay policy, conflict normalization, permission revalidation on replay, and visible failure/recovery states.

## Consequences

### Positive

- Data reads and mutations become symmetric: `DataSource` for reads, `ActionSource` for writes.
- Ready Catalog JSON remains serializable and portable across web and Flutter.
- Action lifecycle is observable and testable across Web Components and Dart.
- Normalized errors make form validation, toast messages, and adapter failures consistent.
- Portable Expression AST prevents JavaScript-only logic from becoming a hidden Ready dependency.
- Security-sensitive blocks have explicit gates before promotion.

### Negative

- Adds a non-trivial runtime layer that must be implemented and tested in both TypeScript and Dart.
- Portable Expression AST is less ergonomic than writing JS functions by hand.
- Some advanced expressions will need host functions or server evaluation, increasing manifest discipline.
- `ActionRunner` can become over-broad if effect kinds expand without strict governance.

### Trade-offs

- The contract prioritizes portability and auditability over maximum expression flexibility.
- `handler` remains available for application code, but Ready Catalog entries are restricted to serializable actions.
- Built-in function registry is intentionally small to keep parity feasible.

## Implementation Phases

### Phase 1: Types and schema contracts

- Add `ActionSource`, `ActionAdapter`, `ActionResult`, `ActionError`, `ActionEffect`, lifecycle payload, Portable Expression AST, and function manifest types to `@bitcode/ui-schema`.
- Add TypeScript exports through `@bitcode/ui` bridges.
- Port equivalent Dart types to `packages/ui-flutter`.
- Add JSON schema or runtime validators for all serializable shapes.

### Phase 2: Portable Expression AST evaluators

- Implement TypeScript evaluator.
- Implement Dart evaluator.
- Add cross-language golden fixtures covering operators, paths, functions, errors, nulls, arrays, aggregation, dates, and locale-sensitive formatting.
- Add manifest validation for required built-ins.

### Phase 3: ActionRunner and ReactiveStore integration

- Implement TypeScript `ActionRunner` and `ReactiveStore.runAction`.
- Implement Dart equivalent.
- Add deterministic effect application.
- Add lifecycle event bridge integration on web and notifiers on Dart.
- Add cancellation, timeout, and adapter error normalization.

### Phase 4: RuntimeContext and adapters

- Add action, permission, navigation, toast, i18n, and model adapter registries.
- Implement local/url/model action adapters for mock/demo contexts.
- Leave BitCode engine adapter as optional integration, not a hard dependency.
- Add fail-loud host validation for missing adapters/functions.

### Phase 5: Block Catalog integration

- Wire ADR-017 `BcBlockRenderer` to validate action requirements before rendering.
- Promote only blocks whose data/action/reactivity contracts pass.
- Add action mock demos to `packages/ui/demo-app` and `packages/ui-flutter/example`.

### Phase 6: Security, offline, and performance hardening

- Add URL allowlist, sanitization, file policy, payment adapter, redaction, and permission validation tests.
- Add offline queue policy for action-capable blocks that declare offline support, including durable persistence, retry/backoff, idempotency keys, replay rules, conflict normalization, permission revalidation on replay, and visible failure/recovery states.
- Add performance tests for debounced actions, aborts, and repeated submissions.

## Acceptance Criteria

- [ ] `ActionSource`, `ActionAdapter`, `ActionResult`, `ActionError`, `ActionEffect`, lifecycle payload, Portable Expression AST, and function manifest types are exported from `@bitcode/ui-schema`.
- [ ] Dart equivalents exist and are exported from `packages/ui-flutter`.
- [ ] Ready Catalog JSON supports only serializable action variants unless a named host handler is declared through conforming manifests and fallback policies for every required platform.
- [ ] `ReactiveStore.runAction(action, payload, ctx)` exists in TypeScript and Dart.
- [ ] Action lifecycle emits start, success, error, validation-error, and settled with parity payloads on web and Dart.
- [ ] Effects run in deterministic order and are covered by tests.
- [ ] Portable Expression AST evaluators exist in TypeScript and Dart with shared golden fixtures, and permission-filtered scopes are validated before sensitive payload expression evaluation.
- [ ] Built-in function registry validates name, version, platform support, purity, determinism, and fallback.
- [ ] Missing functions/adapters fail loudly with typed diagnostics.
- [ ] Security/privacy validators reject secrets in schema, unsafe URLs, unsanitized rich text opt-out without policy, unsafe file actions, and payment actions without a payment adapter.
- [ ] Offline-capable action blocks prove durable queue persistence, retry/backoff, idempotency key behavior, replay rules, conflict normalization, permission revalidation on replay, and visible failure/recovery states.
- [ ] At least one schema-strategy block and one component-strategy block use `ActionSource` successfully in web and Flutter demos.

## Related ADRs

- ADR-002: `DataSource`, tokens, cross-platform architecture, and adapter boundaries.
- ADR-004: `ReactiveStore`, reactions, and evaluator interface. This ADR extends it with actions and Portable Expression AST.
- ADR-005: `UiSchema`, `UiAction`, and `BcRenderer`. This ADR extends actions beyond button descriptors.
- ADR-006: Single Caller API and Pure/Adapt platform strategy.
- ADR-007: Semantic event and normalized contract discipline.
- ADR-011: Flutter adapter architecture and Dart parity requirements.
- ADR-012: `@bitcode/ui-schema` as canonical portable type package.
- ADR-013: Runtime providers, hooks, animation, adapt, and Flutter infrastructure.
- ADR-017: Block Catalog architecture that consumes this contract.

## References

- `sprints/package-ui-elementor/block-catalog-candidate-matrix.json`
- `sprints/package-ui-elementor/block-catalog-parity-matrix.md`
- `sprints/package-ui-elementor/REVIEW-FINDINGS.md`
- `packages/ui-schema/src/types.ts`
- `packages/ui-schema/src/reactive/`
- `packages/ui-schema/src/schema/`
- `packages/ui/src/core/data-source.ts`
- `packages/ui/src/core/renderer/`
- `packages/ui-flutter/lib/src/core/reactive_store.dart`
- `packages/ui-flutter/lib/src/schema/`
