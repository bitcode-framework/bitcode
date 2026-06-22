# ADR-012: Extract `@bitcode/ui-schema` as Standalone Package

## Status: Accepted (2026-06-22)

## Date: 2026-06-22

## Depends on: [ADR-002](./002-cross-platform-architecture.md), [ADR-004](./004-reactive-engine.md), [ADR-005](./005-metadata-driven-ui.md), [ADR-011](./011-flutter-adapter-architecture.md)

## Context

ADR-011 Section 1 specifies that the Flutter adapter (Tier C) will port 4 infrastructure layers from TypeScript to Dart:

1. Token palette (`theme.css` → `bc_theme.dart`)
2. API surface (`*.contract.ts` → Dart contract classes)
3. **Reactive engine** (`store.ts` → `reactive_store.dart`)
4. **Schema renderer** (`UiSchema` types → Dart `UiSchema` class)

Layers 3 and 4 currently live inside `packages/ui/src/core/` — a Mitosis project that compiles `.lite.tsx` into React/Vue/RN output. This creates three problems for the Flutter port:

### Problem 1: Coupling to Mitosis build pipeline

`packages/ui` has a complex build pipeline: `mitosis build` → `apply-adapt-overrides` → `sync-regular-ts` → `reformat-react-output` → `generate-index`. The Flutter token generator (`generate-flutter-tokens.mjs`) already runs outside this pipeline (it reads `theme.css` directly). But the reactive engine and schema types are buried inside `src/core/reactive/` and `src/core/renderer/`, intermixed with Mitosis-specific files (`.lite.tsx`, `scope-context.ts` that imports `@builder.io/mitosis`).

Extracting the portable layers into a standalone package means the Flutter port can depend on `@bitcode/ui-schema` (or port from its clean source) without pulling in Mitosis, Stencil, or any web-specific dependency.

### Problem 2: Type duplication risk

Without a canonical source, the Dart port would be based on a snapshot of the TS types. When the TS types change (new `Effect` kind, new `FieldType`, new `BcEventSource` value), the Dart port drifts silently. A standalone package with a single set of type definitions serves as the **canonical source** — both the TS renderer (`@bitcode/ui`) and the Dart renderer (`@bitcode/ui-flutter`) derive from it.

### Problem 3: Standalone usage blocked

JsonForms-style standalone usage (schema → reactive store, no renderer) currently requires importing from `@bitcode/ui`, which pulls in Mitosis types and Stencil component references. A standalone package enables headless usage: define a `UiSchema`, create a `ReactiveStore`, wire reactions — without any UI framework.

### Prior art

- **JsonForms** — Separates `@jsonforms/core` (schema + reactive engine) from `@jsonforms/react` (React renderer). Same pattern: headless core, framework-specific renderer.
- **Zustand** — Extracted `zustand/vanilla` from `zustand` so non-React consumers can use the store without React bindings.
- **TanStack Table** — `@tanstack/table-core` (headless) + `@tanstack/react-table` (React adapter) + `@tanstack/vue-table` (Vue adapter).

## Decision

Extract portable, framework-agnostic UI schema types and the `ReactiveStore` engine into a new standalone package: **`@bitcode/ui-schema`**.

### Section 1 — Package scope

**`@bitcode/ui-schema` contains:**

| Module | Contents | Source (extracted from) |
|---|---|---|
| `src/types.ts` | Core types: `BcChangeMeta`, `BcEventSource`, `DataSource`, `BcKey`, `BcSelection`, field events, component sizes/variants | `packages/ui/src/core/types.ts` |
| `src/reactive/types.ts` | `Reaction`, `Effect`, `Evaluator`, `Snapshot`, `FieldMeta`, `ValidationRule`, `ReactiveScopeLike` | `packages/ui/src/core/reactive/types.ts` |
| `src/reactive/store.ts` | `ReactiveStore` class — headless reactive engine (pure TS) | `packages/ui/src/core/reactive/store.ts` |
| `src/schema/types.ts` | `UiSchema`, `FieldNode`, `FieldType`, `UiAction`, `BcRendererProps` | `packages/ui/src/core/renderer/types.ts` |
| `src/schema/resolver.ts` | `FieldResolver` interface + `resolveField` / `resolveFieldType` / `resolveFields` pure functions + `defaultFieldResolver` | `packages/ui/src/core/renderer/resolver.ts` |

**`@bitcode/ui-schema` does NOT contain:**

- Mitosis components (`.lite.tsx`) — stay in `@bitcode/ui`
- Framework-specific types (`BcContext`, `BcAuthConfig`) — stay in `@bitcode/ui` (reference runtime functions)
- Mitosis context (`BcReactiveScopeContext`) — stay in `@bitcode/ui` (imports `@builder.io/mitosis`)
- Stencil web component event bridge — stay in `@bitcode/ui`
- Token system / CSS — stay in `@bitcode/ui` (already separately portable via `theme.css`)

### Section 2 — Zero-dependency, zero-platform-globals

The package has:
- **Zero runtime dependencies** (no `@builder.io/mitosis`, no `react`, no `vue`, no `react-native`)
- **Zero platform globals** (no `window`, `document`, `process`, `fetch` at module scope — `fetch` only appears inside a `DATA` effect's `url` branch, which is consumer-opt-in)
- **Node 22+ ESM compatible** — all internal imports use explicit `.ts` extensions for native ESM resolution

This makes it safe to import from:
- Node.js (SSR, CLI tools, test runners)
- React/Vue/RN via Mitosis (the existing `@bitcode/ui` pipeline)
- Dart (port source — no TS-specific runtime semantics)
- Any standalone consumer (JsonForms-style headless usage)

### Section 3 — Re-export bridges in `@bitcode/ui`

The original files in `packages/ui/src/core/` are replaced with **re-export bridges** — thin files that `export from '@bitcode/ui-schema'`. This preserves backwards compatibility:

- **48 contract + lite files** that import from `'../../../core/types'`, `'../../../core/reactive/types'`, `'../../../core/reactive/store'`, `'../../../core/renderer/types'`, `'../../../core/renderer/resolver'` continue to work unchanged (plus 1 root barrel `src/index.ts`).
- **Mitosis `.lite.tsx` components** that import from `'../../../core/reactive/store'` continue to work — the bridge resolves to `@bitcode/ui-schema` at build time.
- **The `sync-regular-ts.mjs` build step** mirrors bridge files into `output/{react,vue,stencil,react-native}/src/` verbatim, and `@bitcode/ui-schema` resolves from `node_modules` in the output tree.

Bridge files:

| File | Re-exports |
|---|---|
| `src/core/types.ts` | All portable types from `@bitcode/ui-schema` + keeps `BcAuthConfig`/`BcContext` local |
| `src/core/reactive/types.ts` | All reactive types from `@bitcode/ui-schema` |
| `src/core/reactive/store.ts` | `ReactiveStore` class + `Listener` type from `@bitcode/ui-schema` |
| `src/core/reactive/index.ts` | Barrel re-export from `@bitcode/ui-schema` |
| `src/core/renderer/types.ts` | All schema types from `@bitcode/ui-schema` |
| `src/core/renderer/resolver.ts` | `resolveField` + helpers from `@bitcode/ui-schema` |
| `src/core/renderer/index.ts` | Barrel re-export from `@bitcode/ui-schema` |

### Section 4 — Dependency relationship

```
@bitcode/ui-schema  ← types + reactive store (pure TS, zero deps)
       ↑
       ├── @bitcode/ui         ← React/Vue/RN renderer (Mitosis)
       ├── @bitcode/ui-flutter  ← Flutter renderer (Dart port, Phase 5a)
       └── BitCode engine       ← domain → UiSchema mapping (optional)
```

`@bitcode/ui-schema` is a **peer** to `go-json`, not a child. The `Evaluator` interface allows `go-json` (or any expression engine) to be plugged in for string-based reaction expressions, but this is optional — reactions default to plain JS functions.

`@bitcode/ui` declares `@bitcode/ui-schema` as a `file:` dependency in `package.json`:
```json
"dependencies": {
  "@bitcode/ui-schema": "file:../ui-schema"
}
```
On `npm install`, this creates a junction (Windows) / symlink (Unix) from `node_modules/@bitcode/ui-schema` to `packages/ui-schema`, so source edits are immediately visible without reinstall.

### Section 5 — File extension convention

All internal imports in `@bitcode/ui-schema` use explicit `.ts` extensions:
```typescript
export { ReactiveStore } from './store.ts';
export type { Reaction } from './types.ts';
```

This is required because:
1. **Node 22+ ESM** requires explicit extensions for relative imports (no `moduleResolution: Bundler` at runtime)
2. **TypeScript `moduleResolution: Bundler`** (used by `@bitcode/ui`) allows `.ts` extensions — no conflict
3. **Mitosis compiler** preserves the import path as-is in generated output, and `@bitcode/ui-schema` resolves from `node_modules` regardless of extension

The existing `@bitcode/ui` source files continue to use extensionless imports (e.g. `from './store'`) because they run under Node's `--test` flag which has special TS resolution. The bridge files in `@bitcode/ui` import from `@bitcode/ui-schema` (a package import, not relative), so no extension is needed there.

## Acceptance Criteria

- [x] `packages/ui-schema/` exists with `package.json`, `src/`, `README.md`
- [x] All portable types extracted (core types, reactive types, schema types)
- [x] `ReactiveStore` class extracted (identical behavior)
- [x] `resolveField` / `resolveFieldType` / `resolveFields` / `defaultFieldResolver` extracted
- [x] Re-export bridges in `@bitcode/ui` (7 files)
- [x] `npm run build` in `@bitcode/ui` passes (179 components generated)
- [x] `npm test` in `@bitcode/ui` passes (1575 source tests + 797 runtime tests, 0 failures)
- [x] `npm run audit:contract` passes (0 violations)
- [x] No existing import paths changed (48 contract/lite files + 1 barrel untouched)

## Consequences

### Positive

- **Flutter port unblocked.** Phase 5a can port from `@bitcode/ui-schema/src/` (clean, no Mitosis) instead of digging through `packages/ui/src/core/`.
- **Canonical source.** One set of type definitions. TS and Dart both derive from the same source.
- **Standalone usage.** Headless `ReactiveStore` + `UiSchema` available without any UI framework.
- **Cleaner `@bitcode/ui`.** Core types are now clearly separated from framework-specific types (`BcContext`, `BcAuthConfig` stay local).
- **No breaking changes.** All 48 contract/lite files, all Mitosis components, all tests continue to work via re-export bridges.

### Negative

- **One more package.** `packages/ui-schema` adds a directory and `package.json` to maintain.
- **`file:` dependency.** Requires `npm install` in `@bitcode/ui` after cloning (already required for `patch-package` postinstall). The junction/symlink means source edits are live, but a fresh clone needs `npm install` to create the link.
- **`.ts` extension convention.** `@bitcode/ui-schema` uses `.ts` extensions in imports (Node ESM requirement), while `@bitcode/ui` uses extensionless imports (Node `--test` resolution). This asymmetry is documented in Section 5 and is the minimal viable approach for a non-monorepo setup.

### Neutral

- **`go-json` as optional evaluator.** The `Evaluator` interface in `@bitcode/ui-schema` allows `go-json` to be plugged in for string-based reaction expressions. This is already the case in the pre-extraction code; extraction doesn't change the relationship.

## References

- [ADR-002](./002-cross-platform-architecture.md) — Cross-platform architecture (Tier A/B/C)
- [ADR-004](./004-reactive-engine.md) — Reactive engine (ReactiveStore, reactions, effects)
- [ADR-005](./005-metadata-driven-ui.md) — Metadata-driven UI (UiSchema, BcRenderer)
- [ADR-011](./011-flutter-adapter-architecture.md) — Flutter adapter architecture (Tier C, 75% shared / 25% separate)
- JsonForms architecture: `@jsonforms/core` (headless) + `@jsonforms/react` (renderer)
- TanStack Table architecture: `@tanstack/table-core` (headless) + framework adapters
