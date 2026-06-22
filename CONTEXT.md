# BitCode Platform Context

This context defines the shared language for BitCode platform architecture and the universal UI library.

## Language

**Contract-First API Freeze**:
A hardening phase where each `packages/ui` component must finalize its public contract before further feature expansion. The contract includes props, events, methods, slots, accessibility behavior, controlled and uncontrolled state, and cross-renderer integration semantics.
_Avoid_: polish phase, cleanup, refactor

**Core Integration Surface**:
The P0 subset of `packages/ui` components and infrastructure that defines how host applications integrate with the library. It includes primitives for form input, overlays, data grids, metadata rendering, reactive scope, providers, theming, and Web Component event bridging.
_Avoid_: top components, most used components, first batch

**Semantic Event**:
A renderer-neutral component event whose payload carries the UI meaning directly, such as a changed value, open state, or selection. DOM and native events may be included as metadata, but they are not the primary contract.
_Avoid_: raw event, DOM-first event, framework event payload

**Web Component Event Bridge**:
The `core/wc-event-bridge.ts` utility that dispatches `bc-*` CustomEvents on the host element for Web Component consumers. Each semantic callback (e.g. `onValueChange`) has a corresponding `bc-*` event (e.g. `bc-value-change`). SSR-safe and RN-safe (no-ops when `document` is unavailable).
_Avoid_: event emitter, event dispatcher, Stencil events

**BcKey / BcSelection**:
Standard collection types defined in `core/types.ts` per ADR-007 Section 9. `BcKey = string | number`. `BcSelection = 'all' | Set<BcKey> | BcKey[]`. Used by data components (BcDataTable, BcList) for key-based selection APIs instead of index-only APIs.
_Avoid_: row id, selection set, selected indexes

**Runtime Part Resolution**:
The mechanism for `slotProps` passthrough defined in ADR-007 Section 6.1. Pure function helpers (`partClass`, `partAttr` in `core/part-props.ts`) that components call per-element to read individual values from `slotProps`, avoiding dynamic object spread which Mitosis cannot compile. className is additive (concat after base classes), not replace. Portable to Flutter (no DOM dependency).
_Avoid_: slot injection, compile-time spread, Mitosis plugin hack

**Flutter Adaptive**:
The Flutter implementation strategy where widgets with a native Cupertino counterpart render via `flutter_platform_widgets` (`Platform.*`) and select Cupertino on iOS/macOS, Material 3 on Ancodebuddy/Web/Windows/Linux. Widgets without a Cupertino counterpart (DataTable, Splitter, Accordion, etc.) stay Material-only on every platform. This is a Dart implementation detail; the TypeScript `flutter:` adapter mapping string in `.contract.ts` stays Material-canonical and is descriptive, not a binding codegen target.
_Avoid_: contract-level Cupertino mapping, dual flutter mapping strings, platform branch in props

**Token Bridge**:
The build step (`generate-flutter-tokens.mjs`) that parses `theme.css` HSL channel triplets and emits `bc_theme.dart` with `static const Color` hex values plus four `const` factories: `materialLight()`, `materialDark()`, `cupertinoLight()`, `cupertinoDark()`. Tailwind-style alpha parity is achieved via native `Color.withValues(alpha:)` (Flutter 3.27+), not runtime HSL conversion. `theme.css` remains the single source of truth; hex values in Dart are derived, not authored.
_Avoid_: runtime HSLColor conversion, mutable color tokens, custom alpha helper

**Flutter Tier Separation**:
The accepted split between shared and platform-specific work for the Flutter adapter. Shared across all targets: the `.contract.ts` API surface, `theme.css` tokens, `ReactiveStore` semantics, and `BcRenderer` UiSchema. Separate per Flutter: the Dart behavior implementation (hand-written for P0 in Phase 5a, generator-assisted for P1/P2 in Phase 5b/5c). This is not a flaw but an engineering decision following the ADR-002 finding that the Mitosis-Flutter compiler is not reliable.
_Avoid_: one-source-to-rule-them-all, Mitosis-Flutter codegen, pretending Flutter parity is free
