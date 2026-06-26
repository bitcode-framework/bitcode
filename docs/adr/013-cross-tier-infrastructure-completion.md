# ADR-013: Cross-Tier Infrastructure Completion

- **Status:** Accepted
- **Date:** 2026-06-26
- **Supersedes:** —
- **Amends:** ADR-002 (Section 9 — Cross-tier Infrastructure), ADR-006
- **Related:** ADR-004 (Reactive Engine), ADR-011 (Flutter Adapter)

## Phase Status

| Phase | Description | Status |
|-------|-------------|--------|
| 1 | AppProvider Consumption Surface (useAppContext + useAppConfigStore hooks, Flutter AppConfigStore + BcAppConfig InheritedWidget parity) | ✅ Complete (2026-06-26) |
| 2 | Hooks P0 Adoption (BcDialog/Popover/Tooltip/Drawer/Dropdown consume P0 hooks) | ✅ Complete (2026-06-26) |
| 3 | Hooks P1 New (8 industry-standard headless hooks) | ✅ Complete (2026-06-26) |
| 4 | Animation + Adapt runtime (resolveAnimation + resolveAdapt, wired into 5 overlay/navigation components) | ✅ Complete (2026-06-26) |
| 5a | Flutter env fix + Dart parity baseline | ⬜ Pending |
| 5b | Flutter parity completion + cross-tier parity table | ⬜ Pending |

## Phase 2 Completion Notes (2026-06-26)

**Scope:** 5 overlay/navigation components refactored to consume P0 hooks
(`resolveControllable`, `DisclosureController`/`useDisclosure`, `createFocusTrap`)
instead of ad-hoc `useStore` boolean patterns. Public API, events, and contracts
preserved.

| Component | ADR-013 Ref | Hooks Adopted | Notes |
|-----------|-------------|---------------|-------|
| BcDialog | A011/A012/A013 | `resolveControllable`, `createFocusTrap` | Controlled+uncontrolled `open`; focus trap activates on modal open, deactivates on close. SSR-safe eager init. |
| BcPopover | A014 | `resolveControllable`, `useDisclosure` | Controlled+uncontrolled `open`; `DisclosureController` drives open/close/toggle with dedup; `applyOpen` always emits `onOpenChange` + WC event. |
| BcTooltip | A015 | (excluded) | Pure CSS hover (`group-hover:visible`/`group-hover:opacity-100`), no JS open-state. Exclusion documented via source-text test. Contract keeps `open`/`defaultOpen`/`onOpenChange` for adapter parity. |
| BcDrawer | A016 | `resolveControllable`, `createFocusTrap` | Controlled+uncontrolled `open`; added `defaultOpen`/`onOpenChange`/`closeOnEscape`; Escape key handler; focus trap on open. |
| BcDropdown | A017 | `resolveControllable`, `useDisclosure` | Controlled+uncontrolled `open`; added `open`/`defaultOpen`/`onOpenChange`; `applyOpen(next, reason)` with semantic reasons (`trigger`/`select`/`escape`/`programmatic`). |

**Reactivity pattern (proven across BcDialog/BcPopover/BcDrawer/BcDropdown):**
- `openCell` (`{value}` plain object in `useStore`) backs `resolveControllable`'s
  read path — stable reference, mutated in place, persists across renders.
- `reactiveOpen` (useState boolean) triggers Mitosis re-renders on flip.
- `DisclosureController.onOpenChange` callback syncs both `openCell.value` and
  `reactiveOpen` (state only, no event emission).
- `applyOpen(next, reason)` seeds the cell (uncontrolled), drives
  `disclosure.set(next)`, then **always** emits `props.onOpenChange` + WC event
  + reason-specific cleanup. Emitting in `applyOpen` (not in the disclosure
  callback) avoids `DisclosureController` dedup swallowing controlled-mode emits.
- `createFocusTrap` is initialized **eagerly** in the store (not lazily) to avoid
  `useState` async-update staleness. SSR-safe (`createFocusTrap()` returns
  NOOP_TRAP on server). Lifecycle driven via `onUpdate([state.reactiveOpen, props.open, ...])`.

**Verification (2026-06-26):**
- Source-text tests: 1604 pass, 0 fail (added hook-adoption assertions A011-A018).
- Runtime tests (overlays + dropdown + command): 147 pass, 0 fail across 15 files.
- Build: clean, 0 reformat failures.
- No `as any` / `@ts-ignore` introduced (A050).

## Phase 3 Completion Notes (2026-06-26)

**Scope:** Added the P1 hook set — 8 industry-standard headless SSR-safe hooks
in `packages/ui/src/core/hooks/`. The original 6-hook plan (useHover,
useFocusVisible, useEventListener, useOutsideClick, usePrevious, useInterval)
was superseded by the validation-contract (A020-A027) which specifies 8 hooks
with broader overlay/utility coverage. Each hook is a headless utility class or
factory function (no JSX/DOM in the hook itself); browser APIs are gated behind
`typeof window` / `typeof document` so SSR and React Native imports stay safe.

| Hook | File | API | SSR behavior |
|------|------|-----|--------------|
| usePrevious | `use-previous.ts` | `PreviousValue<T>` class + `usePrevious()` sugar | Pure (no DOM) — works everywhere |
| useIntersection | `use-intersection.ts` | `createIntersectionObserver(opts)` → `{state, observe, unobserve, subscribe, dispose}` | NOOP observer (`isIntersecting:false`, never fires) |
| useInterval | `use-interval.ts` | `createInterval(cb, delay, opts)` → `{active, start, stop, dispose}` | NOOP controller (no `setInterval`) |
| useTimeout | `use-timeout.ts` | `createTimeout(cb, delay, opts)` → `{scheduled, start, clear, dispose}` | NOOP controller (no `setTimeout`) |
| useClickOutside | `use-click-outside.ts` | `createClickOutsideHandler(refs, handler, opts)` → `{setEnabled, dispose}` | NOOP handler (no listeners attached) |
| useKeyboard | `use-keyboard.ts` | `createKeyboardHandler(handlerMap, opts)` → `{attach, detach, dispose}` | NOOP handler (no listeners, attach is no-op) |
| useScrollLock | `use-scroll-lock.ts` | `createScrollLock(opts)` → `{locked, lock, unlock, dispose}` | NOOP lock (ref-counted; shared registry per element) |
| useLocalStorage | `use-local-storage.ts` | `createLocalStorage<T>(key, default)` → `{get, set, subscribe, dispose}` | Returns default, never fires, `set` inert |

**Design patterns (consistent with existing P0 hooks):**
- Headless: factory returns a controller object; no JSX, no render. The Mitosis
  component owns the instance in `useStore` and binds it to elements on mount.
- SSR-safe: each factory checks `typeof window`/`typeof document` (or feature
  presence like `IntersectionObserver`) and returns a frozen NOOP controller so
  server-rendered trees import and construct without crashing.
- Injectable timers: `useInterval`/`useTimeout` accept `setTimer`/`clearTimer`
  options so tests drive a fake clock (mirrors `DebouncedValue`).
- Disposal: every controller exposes `dispose()` (idempotent) to detach
  listeners / cancel timers / clear subscriber sets.
- Ref-counted scroll lock: `createScrollLock` shares a per-element registry via
  `WeakMap` so nested overlays compose — only the last `unlock()` restores the
  original `overflow`.
- Cross-tab storage sync: `useLocalStorage` listens to the `storage` event so
  writes in other tabs propagate to subscribers; malformed JSON is treated as
  absent (returns default, never throws).

**Exports (A028):** All 8 hooks exported from `packages/ui/src/core/hooks/index.ts`
barrel. `packages/ui/src/index.ts` re-exports via `export * from './core/hooks'`
(auto-generated barrel picks up new exports automatically).

**Package.json (A030):** `./core/hooks` export map entry already present and
unchanged — verified no regression.

**Verification (2026-06-26):**
- Source tests: 1652 pass, 0 fail (1604 baseline + 48 new across 8 hook test
  files). Each hook test verifies SSR no-op path + core logic + disposal.
- No `as any` / `@ts-ignore` / `@ts-expect-error` in source or test files (A050).
  Global shimming in tests uses typed casts (`globalThis as { document?: ... }`)
  and `Object.defineProperty` for Node 22's read-only `navigator`.

## Phase 4 Completion Notes (2026-06-26)

**Scope:** Implemented the runtime for `AnimationProps` (CSS-transition driver for
web Tier A) and `AdaptProps` (breakpoint/platform morph resolver), then wired
`animation?` + `adapt?` props into the 5 P0 overlay/navigation components.
Previously these were contract-only `types.ts` files with no runtime behavior.

### Animation runtime (`packages/ui/src/core/animation/runtime.ts`)

| Export | Purpose |
|--------|---------|
| `resolveAnimation(props?, options?)` | Merges `AnimationProps` with `ANIMATION_DEFAULT`; collapses to `'none'` when `prefers-reduced-motion` is active or `kind === 'none'`. Returns `ResolvedAnimation`. |
| `ResolvedAnimation` | `{ kind, duration, easing, cssTransition }` — `cssTransition` is a ready-to-use CSS value (e.g. `opacity 200ms ease-in-out`) or `none`. |
| `getAnimationClasses(propsOrKind)` | Convenience: returns `{ enter, exit }` Tailwind utility class strings per kind. Empty strings for `'none'`/reduced-motion. |
| `prefersReducedMotion()` / `detectPrefersReducedMotion()` | SSR-safe detector (returns `false` on server). |
| `ResolveAnimationOptions` | `{ reducedMotion?: boolean }` — injectable override for tests (avoids DOM stubbing). |

Kind → CSS property mapping: `fade`→`opacity`, `scale`+`slide-*`→`transform`, `none`→`none`.

### Adapt runtime (`packages/ui/src/core/adapt/runtime.ts`)

| Export | Purpose |
|--------|---------|
| `resolveAdapt(props?, ctx?)` | Evaluates the `when` condition (breakpoint via `createMediaQueryWatcher` or `platform:` tag) and returns `ResolvedAdapt`. |
| `ResolvedAdapt` | `{ applies: boolean, target: AdaptTarget | null }`. |
| `AdaptContext` | `{ platform?: AdaptPlatform; breakpointMatcher?: (q: string) => boolean }` — injectable for tests + SSR-safe (no window access unless matcher provided). |

Defaults: `when: 'md'`, `to: 'sheet'`. Breakpoint conditions use
`BC_BREAKPOINTS` from `core/hooks/use-media-query`. Platform conditions
(`platform:mobile` / `platform:rn`) match against `ctx.platform`.

### Component wiring (A039/A040)

| Component | `animation?` | `adapt?` | Notes |
|-----------|:---:|:---:|-------|
| BcDialog | ✅ | ✅ | `resolveAnimation` drives content transition class. |
| BcPopover | ✅ | ✅ | `resolveAdapt` morphs to a sheet on mobile (`applies` → bottom-sheet shell). **A040 proof.** |
| BcDrawer | ✅ | ✅ | `resolveAnimation` drives panel transition class. |
| BcTooltip | ✅ | ✅ | `resolveAnimation` drives opacity transition (replaces hardcoded `transition-opacity`). |
| BcDropdown | ✅ | ✅ | `resolveAnimation` drives menu transition class. |

All 5 components import `resolveAnimation` from `@/core/animation`; BcPopover
additionally imports `resolveAdapt` from `@/core/adapt` and branches render on
`state.adaptResult.applies`.

### Barrel / exports

- `core/animation/index.ts` re-exports `resolveAnimation`, `getAnimationClasses`,
  `prefersReducedMotion`, `detectPrefersReducedMotion`, types.
- `core/adapt/index.ts` re-exports `resolveAdapt`, `ResolvedAdapt`, `AdaptContext`.
- `src/index.ts` (auto-generated) now includes `export * from './core/animation'`
  and `export * from './core/adapt'` (added to `generate-index.mjs` CORE_EXPORTS).
- `package.json` export map already had `./core/animation` + `./core/adapt` entries.

**Verification (2026-06-26):**
- Source tests: 1698 pass, 0 fail (1652 baseline + 46 new: 21 animation, 13 adapt,
  12 component wiring assertions across 5 components).
- Build: `npm run build` clean — 179 components generated for react/vue/reactNative.
- No `as any` / `@ts-ignore` / `@ts-expect-error` introduced (A050).
- SSR-safe: both runtimes guard on `typeof window` / `typeof matchMedia`.

## Context

ADR-002 Section 9 mendefinisikan 5 cross-tier infrastructure yang **wajib** tersedia
sebelum komponen dianggap "truly cross-platform": AppProvider, Hooks, Animation contract,
Adapt contract, dan Reactive store.

Status audit per 2026-06-26:

| Item | Status | Notes |
|------|--------|-------|
| AppProvider adoption | Partial | Unit tests pass, belum diadopsi ke semua komponen |
| Hooks subset P0 (8 hooks) | Partial | Implementasi ada, belum dipakai komponen P0 |
| Hooks subset P1 (8 hooks) | ✅ Complete | usePrevious, useIntersection, useInterval, useTimeout, useClickOutside, useKeyboard, useScrollLock, useLocalStorage — headless SSR-safe, tests pass |
| Animation contract | ✅ Complete | `resolveAnimation` + `getAnimationClasses` runtime (core/animation), honors prefers-reduced-motion, wired into BcDialog/Popover/Drawer/Tooltip/Dropdown |
| Adapt contract | ✅ Complete | `resolveAdapt` runtime (core/adapt) with injectable AdaptContext; BcPopover morphs to sheet on mobile |

Saat ini 110/110 komponen enterprise "Done" secara contract+a11y+token+tests,
tetapi **tidak semua benar-benar memakai cross-tier infrastructure**. Banyak
komponen masih punya logic lokal yang seharusnya didelegasikan ke hooks/provider.

## Decision

Selesaikan 5 cross-tier infrastructure agar ADR-002 Section 9 fully realized:

### 1. AppProvider Adoption (full)

- **Scope:** `packages/ui/src/core/providers/AppProvider.lite.tsx` sudah ada.
  Adoption pass = refactor komponen P0/P1 untuk konsumsi `useAppContext()` untuk
  theme/locale/direction/tenant daripada baca `BcSetup` global langsung.
- **Flutter parity:** Port `AppProvider` ke `packages/ui-flutter/lib/src/core/app_provider.dart`
  via `InheritedWidget`. Komponen Flutter baca via `AppProvider.of(context)`.

### 2. Hooks P0 (8 hooks) — adoption + completion

Daftar 8 hooks P0 (sudah ada 9 file di `packages/ui/src/core/hooks/`):

| Hook | Status file | Adoption needed |
|------|-------------|-----------------|
| useId | ✅ ada | ya |
| useDisclosure | ✅ ada | ya |
| useControllableState | ✅ ada | ya |
| useDebouncedValue | ✅ ada | ya |
| useClipboard | ✅ ada | ya |
| useElementSize | ✅ ada | ya |
| useFocusTrap | ✅ ada | ya |
| useMediaQuery | ✅ ada | ya |

Adoption = komponen P0 (Dialog, Popover, Combobox, Select, Drawer, dll) refactor
untuk pakai hook ini, hapus inline logic.

### 3. Hooks P1 (8 hooks) — baru

> **Note:** The original plan listed 6 P1 hooks (useHover, useFocusVisible,
> useEventListener, useOutsideClick, usePrevious, useInterval). The
> validation-contract (A020-A027) supersedes this with an 8-hook industry-
> standard set. The implemented set below is the source of truth.

| Hook | Tier A impl | Tier B/C | Status |
|------|-------------|----------|--------|
| usePrevious | ref cache (PreviousValue class) | ref cache | ✅ |
| useIntersection | IntersectionObserver factory | web-only (NOOP) | ✅ |
| useInterval | setInterval controller (start/stop/dispose) | Timer.periodic | ✅ |
| useTimeout | setTimeout controller (start/clear/dispose) | Timer | ✅ |
| useClickOutside | document pointerdown listener | GestureDetector | ✅ |
| useKeyboard | keydown handler-map (mod/ctrl/shift/alt) | KeyEvent | ✅ |
| useScrollLock | body overflow ref-counted lock | web-only (NOOP) | ✅ |
| useLocalStorage | localStorage + cross-tab storage event | async-storage | ✅ |

All 8 are headless SSR-safe utility classes/factories (DOM gated behind
`typeof window`/`typeof document`), exported from `core/hooks/index.ts` and
re-exported by the auto-generated `src/index.ts`. Each has a co-located
`*.test.ts` (48 new tests; SSR no-op paths, core logic, disposal verified).

### 4. Animation Contract — per-tier implementation

`AnimationProps` (AnimationKind, AnimationEasing) sudah di-export. Yang belum:
- Tier A (web): map ke CSS transitions / Web Animations API
- Tier A (RN): map ke `Animated` API
- Tier C (Flutter): map ke `AnimationController` + `Tween`

### 5. Adapt Contract — runtime Popover->Sheet

`AdaptProps` (breakpoint/platform condition) sudah di-export. Yang belum:
runtime resolver yang mengubah `<BcPopover>` menjadi `<BcSheet>` di mobile viewport
(< `adaptBreakpoint`, default `md` = 768px).

## File Ownership Boundary (untuk paralelisasi agent)

**WAJIB dipatuhi agar tidak bentrok dengan agent ADR-014/ADR-015:**

| Path | Owner ADR |
|------|-----------|
| `packages/ui/src/core/hooks/**` | ADR-013 (exclusive) |
| `packages/ui/src/core/providers/**` | ADR-013 (exclusive) |
| `packages/ui/src/core/adapt/**` | ADR-013 (exclusive) |
| `packages/ui/src/core/animation/**` | ADR-013 (exclusive) |
| `packages/ui/src/core/reactive/**` | ADR-013 (exclusive) |
| `packages/ui-flutter/lib/src/core/app_provider.dart` | ADR-013 (exclusive) |
| `packages/ui-flutter/lib/src/core/hooks_*.dart` | ADR-013 (exclusive, file baru) |
| `packages/ui-flutter/lib/src/core/types.dart` | **SHARED** — ADR-013 hanya add reactive types, ADR-015 hanya add chart types. Bagian terpisah, jangan overlap. |
| Komponen `.lite.tsx` (P0/P1) | ADR-013 boleh touch UNTUK ADOPTION (ganti inline logic → hook). ADR-014 boleh touch UNTUK ENGINE WIRING (tambah data-binding). Koordinasi via rebase. |
| `packages/ui/src/index.ts` (barrel) | **SHARED** — semua ADR add export. Lihat "Merge Coordinator" di bawah. |

## Consequences

**Positif:**
- ADR-002 Section 9 fully realized — komponen benar-benar cross-platform, bukan
  cuma "file ada di 4 target"
- Logic duplikat berkurang (inline → shared hooks)
- Flutter parity untuk infrastructure, bukan hanya widget

**Negatif:**
- Refactor 110 komponen berisiko regression — perlu test gating ketat
- Adoption pass butuh waktu (estimasi: 3-5 hari untuk P0, 5-7 hari P1)
- `types.dart` Flutter jadi shared file → butuh disiplin

**Mitigasi:**
- Adoption pass bertahap per group (forms → overlays → navigation → dst)
- Test gating: `npm test` + `flutter test` harus 0 failure sebelum commit
- `types.dart` Flutter: tambah section comment `// === ADR-013 REACTIVE ===`
  dan `// === ADR-015 CHART ===` untuk memisahkan

## Merge Coordinator (paralelisasi)

File berikut di-touch oleh multiple ADR. Owner merge = human/coordinator:

- `packages/ui/src/index.ts` — barrel export
- `packages/ui-flutter/lib/src/core/types.dart` — shared types
- `docs/components-progress.md` — progress tracker
- `docs/adr/README.md` — ADR index

**Aturan:** agent commit ke branch masing-masing, coordinator rebase + resolve
conflict di file-file ini. Agent TIDAK boleh resolve conflict sendiri.

## Open Questions

1. Reactive store port: apakah perlu ADR terpisah, atau cukup section di ADR-013?
   (Saran: cukup section, karena sudah dicakup ADR-004 scope)
2. Adapt default breakpoint: `md` (768px) atau `sm` (640px)? (Saran: `md`, konsisten
   dengan Tailwind default)

## References

- ADR-002 Section 9 — Cross-tier Infrastructure (definisi awal)
- ADR-004 — Reactive Engine
- ADR-006 — Universal Component Strategy (Pure/Adapt)
- ADR-011 — Flutter Adapter Architecture (parity target)
- `packages/ui/src/core/hooks/` — 9 hook files existing
- `packages/ui/src/core/providers/AppProvider.lite.tsx` — AppProvider existing
- `packages/ui/src/core/adapt/types.ts` — AdaptProps existing
- `packages/ui/src/core/animation/types.ts` — AnimationProps existing
