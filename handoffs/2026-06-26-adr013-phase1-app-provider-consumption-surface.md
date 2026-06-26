# Handoff — ADR-013 Phase 1: AppProvider Consumption Surface (feat-1)

**Date:** 2026-06-26
**Feature:** feat-1 — AppProvider Consumption Surface
**Assertions claimed:** A001–A009, A049, A050, A052, A053
**Status:** ✅ Implementation complete, committed.

---

## What was implemented

### TypeScript tier (`packages/ui` submodule)
- **`src/core/hooks/use-app-context.ts`** (NEW) — `useAppContext()` (reads `BcAppContext`, SSR-safe fallback to `DEFAULT_APP_CONFIG_SNAPSHOT`), `useAppConfigStore(store?)` (explicit store or lazy module-default), `DEFAULT_APP_CONFIG_SNAPSHOT` (frozen), `__resetDefaultAppConfigStore()` (test helper).
- **`src/core/hooks/use-app-context.test.ts`** (NEW) — 10 tests: SSR default snapshot, store imperative updates, RTL direction resolution, frozen defaults. All pass.
- **`src/core/internal/mitosis-runtime-shim.ts`** (NEW) — runtime-safe `createContext`/`useContext` fallbacks so context-bearing modules import cleanly under Node (Mitosis primitives are compile-time-only / `undefined` at runtime). `useContext` returns the context default = SSR/no-provider behaviour.
- **`src/core/providers/context.ts`** (MODIFIED) — imports shim instead of `@builder.io/mitosis` directly; now Node-importable.
- **`src/core/hooks/index.ts`** (MODIFIED) — barrel exports new symbols.
- **`scripts/generate-index.mjs`** (MODIFIED) — added `core/providers` + `core/hooks` to `CORE_EXPORTS`.
- **`src/index.ts`** (MODIFIED) — regenerated root barrel (now exports `useAppContext` etc. — satisfies A002).

### Flutter tier (`packages/ui-flutter`)
- **`lib/src/core/app_config_snapshot.dart`** (NEW, A007) — `AppConfigSnapshot` + enums (`AppThemeMode`, `ResolvedAppTheme`, `AppDirection`, `ResolvedAppDirection`, `AppSize`) + `AppMessage`/`AppMessageQueue`/`AppConfigInit` + `rtlLocales` + `resolveDirection()`. Pure Dart, matches TS field set.
- **`lib/src/core/app_config_store.dart`** (NEW, A005) — pure-Dart parity store: `getSnapshot()`, `setTheme/setLocale/setDirection/setSize/setDateLocale/setMessageQueue`, `applyResolvedTheme`, `subscribe()` (disposer), `dispose()`, `isDisposed()`.
- **`lib/src/core/bc_app_config.dart`** (NEW, A006) — `InheritedWidget` exposing `snapshot`+`store` via `BcAppConfig.of(context)` with SSR-style default fallback. Primary constructor takes `snapshot`/`store`/`child` directly.
- **`lib/src/bc_app.dart`** (MODIFIED, A009) — `StatelessWidget`→`StatefulWidget`; hosts `AppConfigStore`, subscribes to snapshot changes, wraps children in `BcAppConfig`. Backward-compatible (optional `appConfigStore`/`appConfigInit` params).
- **`test/app_config_store_test.dart`** (NEW, A008) — 6 tests covering defaults, setTheme resolution, subscribe/dispose, direction auto-resolution, dateLocale fallback.
- **`lib/bc_ui_flutter.dart`** (MODIFIED) — barrel exports the 3 new core files.

### Docs
- **`docs/adr/013-cross-tier-infrastructure-completion.md`** — Phase Status table marks Phase 1 ✅ Complete.
- **`docs/components-progress.md`** — AppProvider row updated with Phase 1 consumption-surface notes.

---

## Verification

| Assertion | Status | Evidence |
|-----------|--------|----------|
| A001 useAppContext hook exists at `core/hooks/use-app-context.ts` | ✅ | file present, exported |
| A002 exported from hooks + root barrels | ✅ | `hooks/index.ts` + regenerated `src/index.ts` |
| A003 `useAppConfigStore` exists | ✅ | same file |
| A004 SSR-safe default snapshot returned | ✅ | 10 runtime tests pass (Node), shim returns context default |
| A005 Flutter `AppConfigStore` parity | ✅ | `app_config_store.dart`, `dart test` exit 0 |
| A006 Flutter `BcAppConfig` InheritedWidget | ✅ | `bc_app_config.dart`, `of(context)` + default fallback |
| A007 `AppConfigSnapshot` parity | ✅ | `app_config_snapshot.dart` |
| A008 Flutter store tests | ✅ | `app_config_store_test.dart`, 6 tests |
| A009 `BcApp` wires provider | ✅ | `bc_app.dart` StatefulWidget + subscribe |
| A049 full source suite green | ✅ | **1585 pass / 0 fail** |
| A050 no type suppression | ✅ | grep for `as any`/`@ts-ignore`/`dynamic` = none in changed files |
| A052 docs updated | ✅ | ADR-013 + components-progress |
| A053 package.json export map | ✅ | `./core/hooks` entry already present |

**TS source suite:** 1585 pass / 0 fail (re-confirmed after all edits).
**Flutter/Dart:** `dart test` returned exit code 0 three times (non-verbose). Verbose runs time out (>240s) due to slow dart tooling on this machine — a known pre-existing environment risk documented in the mission plan (Phase 5a investigates). The `AppConfigStore` is pure Dart (no Flutter imports), so the exit-0 result is authoritative.

---

## Key design decisions / gotchas

1. **Mitosis `createContext`/`useContext` are compile-time-only.** At runtime in Node they are `undefined`, so `context.ts` crashed when imported by tests. Fix: `mitosis-runtime-shim.ts` provides runtime fallbacks (`useContext` returns the context's default value = SSR/no-provider behaviour). The Mitosis compiler still recognises the call sites by identifier and compiles them per-target. **This shim is the load-bearing piece** — without it, no context-bearing module can be imported under Node.

2. **`BcAppConfig` constructor signature.** Primary constructor takes `snapshot`/`store`/`child` directly (NOT a `data` param) to match the `BcApp` call site and the documented usage. There is a `BcAppConfig.value(...)`-style path too, but the primary is the direct-params one.

3. **Path discrepancy resolved:** validation contract A001 binds `packages/ui/src/core/hooks/use-app-context.ts` (chose this over the architecture doc's `providers/` path — A001 is the binding assertion).

4. **`packages/ui` is a git submodule** — changes there are committed inside the submodule, then the parent repo records the new submodule pointer.

---

## Parallel agents (DO NOT TOUCH)

Another agent is actively working in `packages/ui-flutter` on **charts** (`bc_chart.dart`, `core/types.dart`, `pubspec.yaml`, `bc_chart_test.dart`, `chart_factory_test.dart`, `chart_helpers_test.dart`, `core/chart/`, `test/fixtures/`). They also left `dart_test_out.txt` / `dart_test_err.txt` artifacts and a `CONTRACT_INVENTORY.md`. **None of these are mine.** I did not modify or delete any of them.

---

## Next steps (Phase 2+)

- Phase 2: Adopt `useAppContext()` / `BcAppConfig.of(context)` into P0/P1 components (TS + Flutter) so they consume the provider instead of ad-hoc theme/locale props.
- Phase 5a: Investigate/fix the slow `flutter test` / `dart test` tooling environment so verbose test runs don't time out.
- Consider adding a `packages/ui` section to `docs/codebase.md` (the package isn't documented there at all yet — out of scope for this feature).
