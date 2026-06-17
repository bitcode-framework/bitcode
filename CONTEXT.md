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
