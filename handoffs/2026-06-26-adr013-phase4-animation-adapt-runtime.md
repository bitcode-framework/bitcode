# ADR-013 Phase 4 — Animation + Adapt Runtime

**Date:** 2026-06-26
**Feature:** `feat-4` (ADR-013 Phase 4)
**Status:** ✅ Complete

## What Was Done

### Animation Runtime (`packages/ui/src/core/animation/`)
- `runtime.ts` — `resolveAnimation(props?, options?)` returns `ResolvedAnimation`
  `{ kind, duration, easing, cssTransition }`. Merges with `ANIMATION_DEFAULT`,
  collapses to `'none'` on reduced-motion. `ResolveAnimationOptions` allows
  injecting `reducedMotion` for test-friendly SSR-safe detection.
- `getAnimationClasses(propsOrKind)` returns `{ enter, exit }` Tailwind utility
  class strings. Kind→CSS mapping: fade→opacity, scale+slide-*→transform,
  none→none.
- `prefersReducedMotion()` / `detectPrefersReducedMotion()` SSR-safe detectors
  (guarded by `typeof window` / `typeof document`).
- `index.ts` barrel + `types.ts` (`AnimationProps`, `ResolvedAnimation`,
  `AnimationKind`, `ResolveAnimationOptions`).

### Adapt Runtime (`packages/ui/src/core/adapt/`)
- `runtime.ts` — `resolveAdapt(props?, ctx?)` returns `ResolvedAdapt`
  `{ applies, target, when }`. Uses `BC_BREAKPOINTS` from hooks. Defaults:
  `when: 'md'`, `to: 'sheet'`. Platform conditions (`platform:mobile` /
  `platform:rn`) match against `ctx.platform`.
- `AdaptContext` type `{ platform?, breakpointMatcher? }` for injectable /
  test-friendly detection.
- `index.ts` barrel + `types.ts` (`AdaptProps`, `ResolvedAdapt`, `AdaptContext`,
  `AdaptTarget`).

### Component Wiring (5 overlays — A039/A040)
Added `animation?: AnimationProps` + `adapt?: AdaptProps` props to contracts and
lite sources of: **BcDialog, BcPopover, BcDrawer, BcTooltip, BcDropdown**.
- All import `resolveAnimation` from `@/core/animation` and apply animation
  classes.
- BcPopover consumes `resolveAdapt` to morph to a bottom sheet on mobile (A040
  proof — `return resolveAdapt(props.adapt)` in output).
- Component tests add wiring assertions (12 assertions across 5 test files).

### Barrel / Exports
- `generate-index.mjs` `CORE_EXPORTS` updated to include
  `export * from './core/animation'` and `export * from './core/adapt'`.
- Regenerated `src/index.ts`.

## Verification (2026-06-26)
- `npm run test:source`: **1698 pass, 0 fail** (1652 baseline + 46 new — 21
  animation + 13 adapt runtime + 12 component wiring assertions). (A049)
- Build (`mitosis build`): **179 components generated** for react/vue/reactNative,
  exit 0, clean. (A049)
- No `as any` / `@ts-ignore` / `@ts-expect-error` introduced in source or tests
  (A050). Global shimming uses typed casts (`globalThis as { document?: ... }`).
- No hardcoded user-facing strings (A051).
- Phase commit per A052 (single phase commit).

## Validation Contract Coverage
- A032 `resolveAnimation` default merge — ✅
- A033 reduced-motion collapse — ✅
- A034 `getAnimationClasses` enter/exit utilities — ✅
- A035 kind→CSS mapping (fade/scale/slide-*/none) — ✅
- A036 `AnimationProps` on 5 overlay contracts — ✅
- A037 SSR-safe `prefersReducedMotion` — ✅
- A038 `resolveAdapt` returns applies/target/when — ✅
- A039 `AdaptProps` on 5 overlay contracts — ✅
- A040 BcPopover→sheet morph proof in output — ✅
- A041 `AdaptContext` injectable for tests — ✅
- A049 tests pass + build clean — ✅
- A050 no type suppression — ✅
- A051 no hardcoded user-facing strings — ✅
- A052 commit per phase — ✅

## Files Changed
- `packages/ui/src/core/animation/` — runtime.ts, types.ts, index.ts (new)
- `packages/ui/src/core/adapt/` — runtime.ts, types.ts, index.ts (new)
- `packages/ui/src/core/animation/runtime.test.ts` (new, 21 tests)
- `packages/ui/src/core/adapt/runtime.test.ts` (new, 13 tests)
- `packages/ui/src/index.ts` — regenerated barrel
- `packages/ui/scripts/generate-index.mjs` — CORE_EXPORTS update
- 5 component contract.ts + lite.tsx (BcDialog/Popover/Drawer/Tooltip/Dropdown)
- 5 component test files — wiring assertions
- `docs/adr/013-cross-tier-infrastructure-completion.md` — Phase 4 status ✅
- `docs/codebase.md` — core/ description updated

## What's Left (NOT Phase 4)
Phase 4 is fully complete. Next ADR-013 phases (separate features):
- **Phase 5a** — Flutter env fix + Dart parity baseline
- **Phase 5b** — Flutter parity completion + cross-tier parity table

## Important Context / Gotchas
1. **`tsc --noEmit` is NOT the project's verification method.** The ui package
   has many pre-existing TS errors (missing `@types/react`, `node:test` module
   declarations, echarts types, `toBeAccessible` matcher). These are pre-existing
   and unrelated to ADR-013. The actual verification is `npm run test:source`
   (Node native test runner) + `npm run build` (mitosis).
2. **Animation classes use Tailwind utilities** (e.g.
   `transition-opacity duration-200 ease-in-out`). These exist in the Tailwind
   preset; no CSS additions needed.
3. **Reduced-motion detection is lazy** — `prefersReducedMotion()` reads
   `matchMedia` only when called (not at module load), so SSR is safe.
4. **BcPopover adapt morph is a runtime hint**, not a full visual morph — the
   `resolveAdapt` result is computed in the component and exposed via output,
   proving the wiring path. A full visual morph (rendering BcSheet) is a future
   enhancement, out of Phase 4 scope (contract + resolver only).
