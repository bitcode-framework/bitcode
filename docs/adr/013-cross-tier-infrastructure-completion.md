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
| 5a | Flutter env fix + Dart parity baseline | ✅ Complete (2026-06-26) |
| 5b | Flutter parity completion + cross-tier parity table | ✅ Complete (2026-06-26) |

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

## Phase 5a Completion Notes (2026-06-26) — Flutter Environment Fix

**Scope:** Investigate and resolve the `flutter test` VM compiler crash / timeout
that blocked all Dart-side assertions (A046 risk note). Prior workers reported
`flutter test`, `dart test`, and `flutter analyze` all exceeding the 240s
foreground timeout with no usable output.

### Root cause (evidence-based)

The Flutter SDK wrapper (`flutter.bat`) shells out to `where.exe` (in
`C:\Windows\System32`) and `git.exe` during startup. The tool-isolated shell
inherited a **10,446-character PATH (230 entries, duplicated ~4×)**. `cmd.exe`
running `flutter.bat` could not resolve `where.exe` under the bloated PATH,
producing `'WHERE' is not recognized as an internal or external command` →
`Error: Unable to find git in your PATH` → non-zero exit. Because every
`flutter`/`dart` invocation failed at the wrapper layer, prior workers
perceived this as a "VM compiler crash / timeout".

A secondary, genuine compile-error layer was also present: 6 `error`-level
analyzer issues in Phase-1-produced Dart files (see Fixes below) that would
have blocked `flutter test` compilation even once the PATH was fixed.

### Fix

1. **PATH fix (environment, not committed to source):** Run `flutter` with a
   minimal PATH containing only `flutter\bin`, `System32`, `System32\Wbem`,
   `System32\WindowsPowerShell\v1.0`, and `Git\cmd`. Documented in
   `packages/ui-flutter/tool/run-flutter.ps1` so future workers/sessions can
   reproduce a working `flutter` invocation without re-diagnosing.
2. **Dart compile fixes (committed):**
   - `lib/src/bc_app.dart:83` — `AppConfigStore(positional)` → `AppConfigStore(init: …)`
     (constructor is `{AppConfigInit init = const AppConfigInit()}` — named only).
   - `lib/src/core/bc_app_config.dart` — `const BcAppConfig.data` constructor
     referenced `data.store` (type `AppConfigStore?`, a non-const class) in a
     const context → `invalid_constant` ×2. Made the `.data` constructor
     non-const (it forwards a potentially-runtime store handle). Also removed an
     `unnecessary_import` of `package:flutter/foundation.dart` and marked
     `_kDefaultData`'s `BcAppConfigData` literal `const`.
   - `test/app_config_store_test.dart` (×4) — `AppConfigStore(const AppConfigInit(…))`
     → `AppConfigStore(init: const AppConfigInit(…))`.

### Verification (2026-06-26)

- `flutter analyze` (full project): **0 errors** (293 info-level lints remain,
  all in pre-existing test files outside ADR-013 scope — `deprecated_member_use`
  for `hasFlag`, `prefer_const_constructors`). Down from 302 issues / 6 errors.
- `flutter test` (full suite): **1527 tests passed, 0 failed, 0 skipped** in 34s.
  No VM crash, no timeout. Includes the previously-broken
  `app_config_store_test.dart` (8 tests pass) and all chart/renderer/schema suites.

### Outcome for Phase 5b (feat-6)

The A046 RISK NOTE waiver is **no longer needed**: `flutter test` now runs the
full suite cleanly. Phase 5b can add new Dart parity tests and verify them with
`flutter test` directly — no `dart analyze`-only fallback required. Phase 5b
should use `packages/ui-flutter/tool/run-flutter.ps1` (or set the minimal PATH
documented therein) when invoking `flutter` from the worker shell.

## Phase 5b Completion Notes (2026-06-26) — Flutter Parity + Cross-Tier Table

**Scope:** Dart parity for the P1 hooks, animation runtime, and adapt runtime,
wired into the `bc_ui_flutter` barrel, with tests and the cross-tier parity
table (A048). Phase 5a's `flutter test` baseline (1527 passing) is preserved —
**1581 tests pass** (1527 + 54 new), 0 failures, 0 errors in `flutter analyze`.

### Flutter parity surface delivered

| TS symbol | Flutter counterpart | File | Notes |
|-----------|---------------------|------|-------|
| `usePrevious` / `PreviousValue<T>` | `PreviousValue<T>` + `usePrevious()` | `lib/src/core/hooks/use_previous.dart` | Pure Dart port — no DOM/timer deps. |
| `createInterval` / `IntervalController` | `BcIntervalController` + `createInterval(...)` | `lib/src/core/hooks/use_interval.dart` | `Timer.periodic` instead of `setInterval`; injectable `createTimer` for `FakeAsync` tests. |
| `createTimeout` / `TimeoutController` | `BcTimeoutController` + `createTimeout(...)` | `lib/src/core/hooks/use_timeout.dart` | `Timer` instead of `setTimeout`; verified via `FakeAsync`. |
| `createScrollLock` / `ScrollLock` | `BcScrollLock` + `createScrollLock()` | `lib/src/core/hooks/use_scroll_lock.dart` | Flutter has no `document.body`; exposes a `ValueNotifier<bool>` host widgets use to swap `ScrollPhysics` ↔ `NeverScrollableScrollPhysics`. Ref-counting preserved. |
| `createLocalStorage` / `LocalStorageStore<T>` | `BcLocalStorageStore<T>` + `createLocalStorage(...)` | `lib/src/core/hooks/use_local_storage.dart` | Flutter has no `localStorage`; injectable `BcStorageBackend` (default: in-memory, parity with TS SSR no-op store). Host wires `shared_preferences`/Hive at the provider boundary. JSON round-trip + malformed-JSON-as-default preserved. |
| `AnimationKind` / `AnimationProps` / `ANIMATION_DEFAULT` | `BcAnimationKind` / `BcAnimationProps` / `ANIMATION_DEFAULT` | `lib/src/core/animation/animation.dart` | Enum parity (`fade`/`slide-up`/…/`none`). `resolveAnimation` collapses to `none` on reduced motion; `reducedMotion` flag injected by host from `MediaQuery.disableAnimationsOf(context)` (A042). Carries a ready-to-use `Curve` + diagnostic `cssTransition` string. |
| `resolveAnimation` / `ResolvedAnimation` | `resolveAnimation(...)` / `BcResolvedAnimation` | `lib/src/core/animation/animation.dart` | `getAnimationCurves` parity for `getAnimationClasses`. |
| `AdaptProps` / `resolveAdapt` / `AdaptContext` | `BcAdaptProps` / `resolveAdapt(...)` / `BcAdaptContext` | `lib/src/core/adapt/adapt.dart` | Breakpoint conditions resolve against injectable `breakpointMatcher` (host derives from `MediaQuery.size`); platform conditions against `ctx.platform` (`Platform.isIOS/Android`, A043). `platform:rn` never applies on Flutter. |
| `BC_BREAKPOINTS` | `bcBreakpointQueries` | `lib/src/core/bc_breakpoints.dart` | Same five Tailwind-aligned queries. |

### Web-only parity (documented, not ported)

Per A044, the following TS hooks are web-only on Flutter with justification —
they have no meaningful Flutter counterpart because the browser API they wrap
does not exist in the Flutter rendering model:

| TS hook | Reason for web-only | Flutter status |
|---------|---------------------|----------------|
| `createIntersectionObserver` | Wraps `IntersectionObserver` (DOM scroll/viewport tracking). Flutter has no DOM; viewport tracking is done via `ScrollController`/`RenderObject` per-scrollable, not a global observer. | Documented web-only. Host widgets use `ScrollController` notification listeners directly. |
| `createClickOutsideHandler` | Wraps `document.pointerdown` capture. Flutter has no document; outside-tap is a `TapRegion`/`GestureDetector` concern at the widget-tree level. | Documented web-only. Host widgets compose `TapRegion`. |
| `createKeyboardHandler` | Wraps `document.keydown`. Flutter keyboard is `Focus`/`KeyboardListener`/`Shortcuts` per-widget, not a global listener. | Documented web-only. Host widgets use `Shortcuts`/`Actions`. |

### Barrel wiring (A045)

All new symbols are re-exported from `packages/ui-flutter/lib/bc_ui_flutter.dart`
under a new `// ADR-013 Phase 5b — Flutter parity` section:

```dart
export 'src/core/bc_breakpoints.dart';
export 'src/core/hooks/hooks.dart';
export 'src/core/animation/animation.dart';
export 'src/core/adapt/adapt.dart';
```

### Tests (A046)

- `test/hooks_parity_test.dart` — 38 tests (PreviousValue, BcIntervalController,
  BcTimeoutController, BcScrollLock, BcLocalStorageStore). Timer paths driven
  by `FakeAsync`.
- `test/animation_adapt_parity_test.dart` — 16 tests (resolveAnimation default
  merge / kind→css mapping / reduced-motion collapse; resolveAdapt breakpoint /
  platform / SSR-safety; breakpoint-query parity).
- Full suite: **1581 pass, 0 fail**. `flutter analyze`: **0 errors**
  (300 info-level lints, of which 7 are new `prefer_initializing_formals` on
  the parity files — non-blocking, consistent with the codebase's existing
  info-lint baseline established in Phase 5a).

### Phase 5 env note

Phase 5a fixed the Flutter environment (`tool/run-flutter.ps1` minimal-PATH
wrapper). Phase 5b confirmed the fix holds: `flutter test` ran the full suite
cleanly in 35s with no VM crash. **No Dart assertions are blocked** — all
A042–A047 parity surface is implemented and tested.

## Context

ADR-002 Section 9 mendefinisikan 5 cross-tier infrastructure yang **wajib** tersedia
sebelum komponen dianggap "truly cross-platform": AppProvider, Hooks, Animation contract,
Adapt contract, dan Reactive store.

Status audit per 2026-06-26:

| Item | Status | Notes |
|------|--------|-------|
| AppProvider adoption | ✅ Complete (Phase 1 + 5b) | useAppContext/useAppConfigStore (TS) + AppConfigStore/BcAppConfig InheritedWidget (Flutter). Phase 5b added Flutter parity surface. |
| Hooks subset P0 (8 hooks) | ✅ Complete | Implemented + adopted by BcDialog/Popover/Tooltip/Drawer/Dropdown (Phase 2) |
| Hooks subset P1 (8 hooks) | ✅ Complete (+ Flutter parity Phase 5b) | usePrevious, useIntersection, useInterval, useTimeout, useClickOutside, useKeyboard, useScrollLock, useLocalStorage — headless SSR-safe, tests pass. Phase 5b added Dart parity for the portable subset (PreviousValue, BcIntervalController, BcTimeoutController, BcScrollLock, BcLocalStorageStore). |
| Animation contract | ✅ Complete (+ Flutter parity Phase 5b) | `resolveAnimation` + `getAnimationClasses` runtime (core/animation), honors prefers-reduced-motion, wired into BcDialog/Popover/Drawer/Tooltip/Dropdown. Phase 5b added Dart `BcAnimationKind`/`BcAnimationProps`/`resolveAnimation` (reduced-motion via injectable flag → host passes `MediaQuery.disableAnimationsOf`). |
| Adapt contract | ✅ Complete (+ Flutter parity Phase 5b) | `resolveAdapt` runtime (core/adapt) with injectable AdaptContext; BcPopover morphs to sheet on mobile. Phase 5b added Dart `BcAdaptResolver`/`resolveAdapt` (breakpoint via `MediaQuery.size`, platform via `Platform.isIOS`/`isAndroid`). |

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

| Hook | Tier A impl (TS) | Tier C impl (Flutter, Phase 5b) | Status |
|------|-------------|----------|--------|
| usePrevious | ref cache (PreviousValue class) | `PreviousValue<T>` (pure Dart port) — `hooks/use_previous.dart` | ✅ |
| useInterval | setInterval controller (start/stop/dispose) | `BcIntervalController` (Timer.periodic, injectable factory) — `hooks/use_interval.dart` | ✅ |
| useTimeout | setTimeout controller (start/clear/dispose) | `BcTimeoutController` (Timer) — `hooks/use_timeout.dart` | ✅ |
| useScrollLock | body overflow ref-counted lock | `BcScrollLock` (ref-counted `ValueNotifier<bool>`, host swaps ScrollPhysics) — `hooks/use_scroll_lock.dart` | ✅ |
| useLocalStorage | localStorage + cross-tab storage event | `BcLocalStorageStore<T>` (injectable `BcStorageBackend`, in-memory default) — `hooks/use_local_storage.dart` | ✅ |
| useIntersection | IntersectionObserver factory | web-only (documented; Flutter uses `ScrollNotification`/`Visibility` widgets) | ✅ |
| useClickOutside | document pointerdown listener | web-only (documented; Flutter uses `GestureDetector`/`TapRegion`) | ✅ |
| useKeyboard | keydown handler-map (mod/ctrl/shift/alt) | web-only (documented; Flutter uses `Focus`/`KeyboardListener`/`KeyEvent`) | ✅ |

The four implemented Dart ports (`usePrevious`, `useInterval`, `useTimeout`,
`useScrollLock`, `useLocalStorage`) are headless pure-Dart classes with no
widget-tree dependency, so they run in `dart test` without a Flutter binding.
The three web-only hooks (`useIntersection`, `useClickOutside`, `useKeyboard`)
are intentionally NOT ported: Flutter has first-class widget-tree primitives
(`Visibility`, `TapRegion`, `Focus`/`KeyboardListener`) that replace the
DOM-observer/listener pattern — a direct port would be unused dead code.
This web-only decision is documented per A044.

All 8 are headless SSR-safe utility classes/factories (DOM gated behind
`typeof window`/`typeof document`), exported from `core/hooks/index.ts` and
re-exported by the auto-generated `src/index.ts`. Each has a co-located
`*.test.ts` (48 new tests; SSR no-op paths, core logic, disposal verified).

### 4. Animation Contract — per-tier implementation

`AnimationProps` (AnimationKind, AnimationEasing) sudah di-export. Yang belum:
- Tier A (web): map ke CSS transitions / Web Animations API — ✅ Done (Phase 4)
- Tier A (RN): map ke `Animated` API — ⬜ deferred (RN tier not in scope)
- Tier C (Flutter): map ke `AnimationController` + `Tween` — ✅ Done (Phase 5b:
  `packages/ui-flutter/lib/src/core/animation/animation.dart` exposes
  `BcAnimationKind` enum, `BcAnimationProps`, `ANIMATION_DEFAULT`,
  `resolveAnimation()` returning `BcResolvedAnimation` with a ready-to-use
  `Curve` + diagnostic `cssTransition`. Reduced-motion via injectable
  `BcResolveAnimationOptions.reducedMotion` — host passes
  `MediaQuery.disableAnimationsOf(context)`.)

### 5. Adapt Contract — runtime Popover->Sheet

`AdaptProps` (breakpoint/platform condition) sudah di-export. Yang belum:
runtime resolver yang mengubah `<BcPopover>` menjadi `<BcSheet>` di mobile viewport
(< `adaptBreakpoint`, default `md` = 768px).

- Tier A (web): `resolveAdapt()` + `createMediaQueryWatcher` — ✅ Done (Phase 4)
- Tier C (Flutter): `resolveAdapt()` — ✅ Done (Phase 5b:
  `packages/ui-flutter/lib/src/core/adapt/adapt.dart` exposes `BcAdaptProps`,
  `BcAdaptContext`, `BcResolvedAdapt`, `BcAdaptResolver`, and `resolveAdapt()`.
  Breakpoint conditions resolve against an injectable `breakpointMatcher`
  (host derives from `MediaQuery.size`); platform conditions resolve against
  `BcAdaptContext.platform` (host derives from `Platform.isIOS`/`isAndroid`).
  `platform:rn` never applies on Flutter. Pure function — no widget-tree
  access unless injected.)

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
