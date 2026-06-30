# ADR-018: Block Catalog Data, Action, and Reactivity Contract

| Field | Value |
|-------|-------|
| **Status** | Proposed (2026-06-30; required by ADR-017 before any Ready entry; ADR body to be drafted) |
| **Depends on** | 002, 004, 005, 006, 007, 011, 012, 013, 017 |
| **Amends** | 004 (extends ReactiveStore contract with `ReactiveStore.runAction`), 005 (extends `UiAction` to describe `ActionSource` and typed effects), 007 (formalises ActionResult/ActionError shapes) |
| **Impact** | Adds `ActionSource`, `ActionAdapter`, `ActionResult`, `ActionError`, `ActionEffect`, `ActionRunner` lifecycle to `@bitcode/ui-schema`; introduces Portable Expression AST and built-in function registry |

---

## Context

The Block Catalog (ADR-017) requires a portable mutation contract (`ActionSource`) and a portable reactive contract that can be evaluated on web, Flutter, and any future Dart-based runtime (Swift, KMP) without depending on string-expression evaluators that are not available on those runtimes. ADR-004 defines a reactive engine with `ReactiveStore` and four effect kinds (DATA / VALIDATION / UI / VALUE) but treats expressions primarily as plain JS functions, with an optional `Evaluator` interface for string expressions. ADR-005 defines `UiAction` as a button descriptor only.

A candidate parity matrix (0.1) at `sprints/package-ui-elementor/block-catalog-candidate-matrix.json` lists 5 foundational work items required before any Ready Block can be marked ready: block-catalog types, action-source, portable expression AST, block renderer, validation tooling. This ADR proposes the action and reactivity foundations that ADR-017 will compose.

## Scope

- Define `ActionSource`, `ActionAdapter`, `ActionResult`, `ActionError`, `ActionEffect`, and `ActionRunner` lifecycle.
- Define a Portable Expression AST portable to TS, Dart, Swift, and Kotlin (without depending on string-evaluator parity across runtimes).
- Define the built-in function registry manifest, host injection pattern, and fail-loud runtime contract.
- Wire `ActionSource` to `ReactiveStore.runAction` and to `UiAction.effects`.

## Out of scope

- Block catalog registry schema (ADR-017).
- Block renderer implementation (ADR-017).
- Third-party plugin catalog.

## Decision Summary (Draft)

1. **ActionSource mirrors DataSource**: `local`, `handler`, `url` (`POST|PUT|PATCH|DELETE`), `model` (operation + adapter). Ready Catalog presets only allow serializable variants (`local`, `url`, `model`); `handler` is host-injected at runtime.
2. **ActionResult is normalized**: `{ ok, status, data, message, errors, errorList, effects }` so web/Flutter/server validation maps to one shape.
3. **ActionLifecycle events are standardized**: `bc-action-start`, `bc-action-success`, `bc-action-error`, `bc-action-validation-error`, `bc-action-settled` (Web Component) and corresponding Dart notifiers.
4. **ActionEffects are typed and safe**: `setValue`, `setMeta`, `setError`, `refreshDataSource`, `invalidateKey`, `toast`, `navigate`, `closeOverlay`, `resetForm`, `emit`. No arbitrary string commands. Sanitization defaults to on and may only be disabled via explicit `security` policy metadata.
5. **Portable Expression AST is the foundation for Ready Catalog reactivity**: structured descriptors with a built-in whitelist of portable functions (`math.*`, `string.*`, `array.*`, `date.*`, `format.*`, `validate.*`, plus a strict ternary/conditional). String expressions remain an advanced escape hatch that requires evaluator parity.
6. **Built-in function registry is versioned and validated** at build-time. Host functions must register with explicit manifest (name, version, platforms, pure, deterministic).

## Acceptance Criteria (Draft)

- `ActionSource`, `ActionAdapter`, `ActionResult`, `ActionError`, `ActionEffect`, and `ActionRunner` types exported from `@bitcode/ui-schema` with Dart ports.
- `ReactiveStore.runAction(action, payload, ctx)` exists in TS and Dart, with deterministic effect ordering and lifecycle emissions.
- Portable Expression AST evaluators exist in TS and Dart with parity test suite.
- Built-in function registry manifest validates host function declarations and fails loud on missing or impure functions.

## Status Notes

This ADR is currently a skeleton. Full Decision, Consequences, and Implementation Phases sections will be drafted after ADR-017 lands and the action/reactivity contract surface is finalised with ADR-004 and ADR-011 owners.
