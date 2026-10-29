# BitCode UI Flutter

`packages/ui-flutter` is the ADR-011 Tier C Flutter adapter for BitCode UI. It provides contract-driven Dart widgets and a schema renderer that can consume the same UI schema shape used by the other BitCode UI adapters.

## Renderer APIs

The package root exports the renderer and schema APIs:

- `BcRenderer` and `BcRendererProps` for schema-driven composition.
- `UiSchema`, `UiSchemaKind`, `FieldNode`, `FieldType`, `UiAction`, and `ColumnDef` for renderer schemas.
- `buildField` via `BcFieldBuilder` for `FieldNode` to Flutter widget rendering.
- `loadDataSource` via `BcDataSourceLoader` for null, local, fetcher, url, and model data sources.

`BcRenderer` supports `form`, `detail`, `list`, `section`, `grid`, and `custom` schema kinds. `BcForm` also supports schema-driven fields, actions, and layout while preserving legacy child-only usage.

### JSON value coercion

`ReactiveStore.get()` returns `dynamic`, and renderer values normally originate
from JSON documents rather than Dart natives. `buildField` therefore coerces
before handing a value to a widget — it never hard-casts. `coerceString`,
`coerceNum`, `coerceDateTime`, and `coerceBool` accept the native type as well
as the JSON shape (ISO 8601 date strings, numeric strings, epoch millis,
`1`/`0`, `'true'`/`'false'`). Unconvertible input degrades to `null` instead of
throwing, and the store value is never mutated.

### Grid layout

`grid` sizes each cell from its content (`Row` + `Expanded`, top-aligned), so a
`detail` with many fields grows its row instead of overflowing a fixed cell
height.

### Height-adaptive widgets

`BcPivot`, `BcViewerDocument`, `BcViewerInstagram`, and `BcViewerTiktok` branch
on the incoming constraints via `LayoutBuilder`. They fill the box when the host
provides a bounded height, and fall back to the web adapter's `600px` default
viewport when it does not — for example inside a `SingleChildScrollView`. This
mirrors the web contracts (`BcPivot` is `overflow-auto` with natural height; the
viewers default to `height: '600px'`) and keeps them safe in any parent.

## Cross-tier infrastructure parity (ADR-013)

The package ships Flutter parity for the `@bitcode/ui` cross-tier infrastructure
so the same hook/animation/adapt contracts work on both tiers:

- **AppProvider** — `AppConfigStore` + `BcAppConfig` `InheritedWidget` expose an
  `AppConfigSnapshot` to descendants (`BcAppConfig.of(context)`), mirroring the
  TS `useAppContext`/`BcAppContext` (Phase 1).
- **Hooks (P1 portable subset)** — `PreviousValue` (usePrevious),
  `BcIntervalController` (useInterval, `Timer.periodic`),
  `BcTimeoutController` (useTimeout, `Timer`),
  `BcScrollLock` (useScrollLock, ref-counted `ValueNotifier<bool>` for
  `ScrollPhysics` swapping), and `BcLocalStorageStore<T>` (useLocalStorage,
  injectable `BcStorageBackend`). The DOM-only hooks
  (`useIntersection`/`useClickOutside`/`useKeyboard`) are intentionally not
  ported — Flutter's widget-tree primitives (`Visibility`/`TapRegion`/`Focus`)
  replace them (Phase 5b).
- **Animation** — `BcAnimationKind` enum, `BcAnimationProps`, `ANIMATION_DEFAULT`,
  and `resolveAnimation()` returning a `BcResolvedAnimation` with a ready-to-use
  `Curve` + diagnostic `cssTransition`. Reduced-motion via an injectable flag
  the host derives from `MediaQuery.disableAnimationsOf(context)` (Phase 5b).
- **Adapt** — `BcAdaptProps`, `resolveAdapt()`, `BcAdaptResolver`, and
  `BcAdaptContext` (breakpoint via an injectable matcher the host derives from
  `MediaQuery.size`; platform via `Platform.isIOS`/`isAndroid`). `platform:rn`
  never applies on Flutter (Phase 5b).

All symbols are exported from the package barrel (`package:bc_ui_flutter/bc_ui_flutter.dart`).
See `CONTRACT_INVENTORY.md` → "ADR-013 Cross-tier Infrastructure Parity" and
`docs/adr/013-cross-tier-infrastructure-completion.md` (Phase 5b) for the full
cross-tier parity table.

## Setup and tests

```bash
cd packages/ui-flutter
flutter pub get
flutter test
dart analyze lib/
dart run build_runner build --delete-conflicting-outputs
```

Renderer scoped tests:

```bash
cd packages/ui-flutter
flutter test test/renderer
```

## Usage

```dart
import 'package:flutter/material.dart';
import 'package:bc_ui_flutter/bc_ui_flutter.dart';

class CustomerForm extends StatelessWidget {
  const CustomerForm({super.key});

  @override
  Widget build(BuildContext context) {
    const schema = UiSchema(
      type: UiSchemaKind.form,
      title: 'Customer',
      fields: [
        FieldNode(
          name: 'name',
          label: 'Name',
          type: FieldType.input,
          required: true,
        ),
        FieldNode(
          name: 'active',
          label: 'Active',
          type: FieldType.checkbox,
        ),
      ],
      actions: [
        UiAction(name: 'save', label: 'Save', submit: true),
      ],
    );

    return BcRenderer(
      props: BcRendererProps(
        schema: schema,
        values: const {'active': true},
        onSubmit: (values) {
          debugPrint('submit: $values');
        },
      ),
    );
  }
}
```

## Validation status

- Full suite: **2144 tests passing**, 0 failing.
- Chart tests (ADR-015): 111 passing — 48 golden-fixture parity + 34 helpers + 29 widget-contract.
- Renderer scoped tests: 294 tests passing, including `bc_field_coercion_test.dart`
  (18 — JSON value coercion) and `bc_grid_layout_test.dart` (7 — content-sized
  grid cells).
- `unbounded_height_test.dart`: 13 passing — `BcPivot` and the viewer widgets
  render inside a vertical scroll view without a layout exception.
- `flutter analyze`: 0 errors (3 pre-existing warnings in `bc_chatter.dart`,
  unrelated).
- `build_runner`: generated files in sync.
