# BcUiDemoApp — Flutter Demo

A single-page Flutter showcase app that renders **all 32 `bc_ui_flutter` widgets** in one navigable page. Useful for manual QA, screenshot verification, and contract validation across Material and Cupertino platforms.

## What it demonstrates

- **6 demo sections** — fields, overlays, layout, data display, views, and misc widgets.
- **Material ↔ Cupertino adaptation** via [`flutter_platform_widgets`](https://pub.dev/packages/flutter_platform_widgets) (FPW), matching the adaptive behavior of the Stencil/React targets.
- **Theme tokens** consumed from `bc_theme.dart` (light + dark palettes generated from `theme.css` at build time).
- **Event tracing** — widget interactions `debugPrint` to the console for inspection.
- **Drawer + section list navigation** — every widget is reachable without scrolling through one giant list.

## Run it

```bash
cd packages/ui-flutter/example

# Web (fastest for quick QA)
flutter run -d chrome

# Mobile/desktop (requires an attached device or emulator)
flutter run

# Production web build
flutter build web
```

> On machines with a bloated PATH, invoke Flutter through `../tool/run-flutter.ps1` instead of the bare `flutter` command.

## Verify

```bash
flutter analyze    # expect 0 issues
flutter test       # smoke test: BcUiDemoApp renders, drawer + 6 sections present
flutter build web  # expect success
```

## Structure

```
example/
├── lib/
│   ├── main.dart              # App entry, theme, drawer + section list
│   ├── sections/              # 6 demo section files
│   │   ├── fields_section.dart
│   │   ├── overlays_section.dart
│   │   ├── layout_section.dart
│   │   ├── data_display_section.dart
│   │   ├── views_section.dart
│   │   └── misc_section.dart
│   └── shared/                # Shared demo primitives
│       ├── demo_section.dart  # Section card wrapper
│       ├── demo_heading.dart
│       └── demo_state.dart
├── test/
│   └── widget_test.dart       # Smoke test for BcUiDemoApp
└── web/                       # Flutter web bootstrap
```

The app is a path-dependency consumer of `bc_ui_flutter` — there is no published package to install. Edit widgets in `packages/ui-flutter/lib/` and the demo picks them up on the next hot reload.
