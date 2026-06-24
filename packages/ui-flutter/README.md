# BitCode UI Flutter

`packages/ui-flutter` is the ADR-011 Tier C Flutter adapter for BitCode UI. It provides contract-driven Dart widgets and a schema renderer that can consume the same UI schema shape used by the other BitCode UI adapters.

## Renderer APIs

The package root exports the renderer and schema APIs:

- `BcRenderer` and `BcRendererProps` for schema-driven composition.
- `UiSchema`, `UiSchemaKind`, `FieldNode`, `FieldType`, `UiAction`, and `ColumnDef` for renderer schemas.
- `buildField` via `BcFieldBuilder` for `FieldNode` to Flutter widget rendering.
- `loadDataSource` via `BcDataSourceLoader` for null, local, fetcher, url, and model data sources.

`BcRenderer` supports `form`, `detail`, `list`, `section`, `grid`, and `custom` schema kinds. `BcForm` also supports schema-driven fields, actions, and layout while preserving legacy child-only usage.

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

- Full suite: 1427 tests passing.
- Renderer scoped tests: 269 tests passing.
- `dart analyze lib/`: clean.
- `build_runner`: generated files in sync.
