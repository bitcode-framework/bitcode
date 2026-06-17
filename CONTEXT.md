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
