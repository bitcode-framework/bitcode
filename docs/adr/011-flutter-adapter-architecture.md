# ADR-011: Flutter Adapter Architecture

## Status: Accepted (2026-06-22)

## Date: 2026-06-22

## Depends on: [ADR-002](./002-cross-platform-architecture.md), [ADR-007](./007-enterprise-universal-ui-api-contract.md), [ADR-008](./008-contract-first-api-freeze.md)

## Context

ADR-002 Section 1 menempatkan Flutter sebagai **Tier C (track terpisah)**: "di-generate dari lapis Contract, bukan dari `.lite.tsx`. Dikerjakan setelah Tier A stabil." Saat ini Tier A stabil (177 komponen, 4 target, 0 failures). ADR-011 mengkonkretkan bagaimana Tier C diimplementasikan.

### Masalah yang dihindari

1. **Return type duality.** `ThemeData` (Material) dan `CupertinoThemeData` (Cupertino) adalah tipe berbeda yang tidak saling extends. Dart tidak punya union type. Satu method factory tidak bisa return keduanya.
2. **Consumer if-else hell.** Kalau consumer harus `if (Platform.isIOS) CupertinoApp(...) else MaterialApp(...)`, prinsip "simple by default" gagal.
3. **Flutter SDK lock-in di core.** `ChangeNotifier` dari `package:flutter/foundation.dart` mengunci reactive store ke Flutter. Kontrak `store.contract.ts` menyebut "portable" (line 111), violated jika core terikat Flutter SDK.
4. **Runtime HSL conversion.** `HSLColor.fromAHSL().toColor()` menghasilkan non-const object. `ColorScheme` didesain untuk `const`. 30+ token x 2 theme = 60+ konversi per build = GC pressure.
5. **Contract drift.** Kalau Flutter punya API berbeda dari React/RN, consumer harus belajar 2 API. Prinsip "1 code, 5 target" gagal.
6. **Adaptive di contract level.** Kalau `.contract.ts` punya mapping string per-platform (cupertino: ..., material: ...), contract membengkak dan Mitosis terganggu.

### Prior art

- **flutter_platform_widgets (FPW)** — Library yang menyediakan `PlatformApp`, `PlatformText`, `PlatformTextField`, dll. Auto-pilih Cupertino di iOS/macOS, Material di lainnya. Dipakai oleh ribuan app production.
- **MaterialBasedCupertinoThemeData** — Class bridge Flutter native untuk Material → Cupertino theme conversion.
- **JsonForms / Formily** — Schema-driven UI library (React) yang punya separate renderer per target. Pola: shared schema, separate renderer.
- **Tamagui** — Universal UI library yang punai shared token + separate per-platform implementation.

## Decision

Tier C Flutter diimplementasikan via **5 keputusan konkret** yang saling melengkapi:

---

### Section 1 — Flutter Tier Separation (75% shared, 25% separate)

**Shared across all targets (75%):**

| Layer | Source | Port mechanism |
|---|---|---|
| Token palette | `packages/ui/tokens/theme.css` (HSL source of truth) | `generate-flutter-tokens.mjs` → `bc_theme.dart` (const hex) |
| API surface | `*.contract.ts` (props/events/methods/a11y) | Manual port to Dart, validated by contract test |
| Reactive engine | `store.ts` (pure TS) | Port to `reactive_store.dart` (pure Dart) |
| Schema renderer | `BcRenderer.lite.tsx` + `UiSchema` types | Port to `BcRenderer.dart` + `UiSchema` Dart class |

**Separate per Flutter (25%):**

| Layer | Source | Why separate |
|---|---|---|
| Widget behavior | `bc_xx.dart` (hand-written P0, generator-assisted P1/P2) | Mitosis-Flutter compiler tidak reliable (ADR-002 finding) |

**Rationale:** Tidak ada compiler TSX → Dart yang reliable. Tapi 4 layer infra (token, contract, reactive, schema) portable sebagai pure Dart/const. Hanya widget-level behavior yang perlu implementasi terpisah. Phase 5b generator akan stub 80% komponen sederhana dari contract; 20% kompleks (DataTable, Form, dll) hand-write sekali, pakai selamanya.

**Consumer DX impact: ZERO.** Consumer API 100% identik across React/RN/Vue/Flutter. Yang membayar 25% ekstra adalah library developer, bukan consumer.

---

### Section 2 — Token Bridge: Pre-computed Hex Const + Native Alpha (B+)

`generate-flutter-tokens.mjs` parse `theme.css` HSL channel triplets, konversi ke hex sekali saat build-time, emit `bc_theme.dart` dengan `static const Color` values.

```dart
// Generated bc_theme.dart (excerpt):
class BcTheme {
  static const Color primary = Color(0xFF2A4A8C);
  static const Color primaryForeground = Color(0xFFFFFFFF);
  // ... 30+ const colors
}
```

**Tailwind-style alpha parity** via native `Color.withValues()` (Flutter 3.27+):

```dart
Container(color: BcTheme.primary.withValues(alpha: 0.5));
```

Setara dengan web side `hsl(var(--bc-primary) / 0.5)`. Tidak perlu helper custom, tidak perlu runtime HSL conversion.

**Rationale:**

| Alternatif | Kenapa ditolak |
|---|---|
| (A) Runtime `HSLColor.fromAHSL()` | Non-const, GC pressure, 60+ konversi per build. Anti-pattern di Flutter. |
| (C) Hex const + custom alpha helper | `Color.withValues()` sudah native sejak Flutter 3.27. Custom helper redundant. |

`theme.css` tetap single source of truth. Hex di `bc_theme.dart` derived, bukan authored. Sesuai ADR-002 Section 3.

---

### Section 3 — 3-Layer Theme: Const Colors → Explicit Platform → Platform-Aware App

Karena `ThemeData` dan `CupertinoThemeData` tidak bisa di-union (Dart tidak punya union type), "platform-aware" hidup di level **app widget**, bukan di level theme factory.

**Layer 1 — `BcTheme`: const colors + 4 explicit factory methods (pure data):**

```dart
class BcTheme {
  // Const colors (compile-time, zero allocation)
  static const Color primary = Color(0xFF2A4A8C);
  static const Color success = Color(0xFF16A34A);
  static const Color warning = Color(0xFFEAB308);
  // ... 30+ const colors

  // Explicit Material (for Ancodebuddy/Web/Windows/Linux or forced Material)
  static ThemeData materialLight() => ThemeData(
    colorScheme: ColorScheme.light(primary: primary, /* ... */),
    useMaterial3: true,
  );
  static ThemeData materialDark() => ThemeData(
    colorScheme: ColorScheme.dark(primary: primary, /* ... */),
    useMaterial3: true,
  );

  // Explicit Cupertino (for iOS/macOS or forced Cupertino)
  static CupertinoThemeData cupertinoLight() => CupertinoThemeData(
    brightness: Brightness.light,
    primaryColor: primary,
    /* ... */
  );
  static CupertinoThemeData cupertinoDark() => CupertinoThemeData(
    brightness: Brightness.dark,
    primaryColor: primary,
    /* ... */
  );
}
```

**Layer 2 — `BcApp`: platform-aware entry point via FPW `PlatformApp`:**

```dart
class BcApp extends StatelessWidget {
  const BcApp({super.key, this.home, this.routes, this.title});
  final Widget? home;

  @override
  Widget build(BuildContext context) {
    return PlatformApp(
      // FPW auto-picks: iOS/macOS → Cupertino, else → Material
      material: (_, __) => BcTheme.materialLight(),
      materialDarkTheme: (_, __) => BcTheme.materialDark(),
      cupertino: (_, __) => BcTheme.cupertinoLight(),
      home: home,
    );
  }
}
```

**Consumer experience:**

```dart
// SIMPLE (platform-aware, zero config):
void main() => runApp(BcApp(home: MyScreen()));
// → iOS: Cupertino auto, Ancodebuddy: Material auto, Web: Material auto

// POWERFUL (explicit override):
PlatformApp(
  material: (_, __) => BcTheme.materialDark(),
  cupertino: (_, __) => BcTheme.cupertinoLight(),
  home: MyScreen(),
);

// RAW (token only, no theme):
Container(color: BcTheme.primary.withValues(alpha: 0.5));
```

**Rationale:**

| Alternatif | Kenapa ditolak |
|---|---|
| Single method `static ??? adaptive()` | Return type union tidak ada di Dart. Tidak compile. |
| Sealed class `Either<ThemeData, CupertinoThemeData>` | Consumer tetap if-else (`if (theme is ThemeData) ...`). Hilang "satu method" benefit. |
| Contract-level Cupertino mapping string | Membengkakkan contract, mengganggu Mitosis, freeze ADR-008 violated. |

SRP terjaga: `BcTheme` = pure data (testable tanpa widget tree), `BcApp` = widget (boleh grow dengan locale/router nanti).

---

### Section 4 — Hybrid Reactive Store: Pure Dart Core + Thin Flutter Adapter

**Layer 1 — `reactive_store.dart` (pure Dart, no Flutter import):**

```dart
// Port of store.ts. Pure Dart. Testable via dart test (no Flutter binding).
class ReactiveStore {
  final Map<String, dynamic> _values = {};
  final Set<VoidCallback> _listeners = {};

  void subscribe(VoidCallback listener) { _listeners.add(listener); }
  void unsubscribe(VoidCallback listener) { _listeners.remove(listener); }
  dynamic get(String path) { /* ... */ }
  void set(String path, dynamic value, BcChangeMeta meta) {
    _values[path] = value;
    for (final listener in _listeners) {
      listener();
    }
  }
}
```

**Layer 2 — `reactive_store_notifier.dart` (thin Flutter adapter, ~15 lines):**

```dart
import 'package:flutter/foundation.dart';

class ReactiveStoreNotifier extends ChangeNotifier {
  final ReactiveStore _store;
  late final VoidCallback _listener;

  ReactiveStoreNotifier(this._store) {
    _listener = () => notifyListeners();
    _store.subscribe(_listener);
  }

  ReactiveStore get store => _store;

  dynamic get(String path) => _store.get(path);
  void set(String path, dynamic value, BcChangeMeta meta) =>
      _store.set(path, value, meta);

  @override
  void dispose() {
    _store.unsubscribe(_listener);
    super.dispose();
  }
}
```

**Rationale:**

| Alternatif | Kenapa ditolak |
|---|---|
| (A) `ChangeNotifier` di core | Lock core ke Flutter SDK. Violates kontrak "portable". Core tidak reusable di CLI/server/future Dart targets. |
| (B) Pure Dart tanpa adapter | Setiap widget reinvent wiring = drift. 176 widget = 176 tempat rawan lupa cleanup. Memory leak. |

Mirror arsitektur existing TS: `store.ts` (pure, no DOM) + `wc-event-bridge.ts` (web adapter). Konsistensi lintas bahasa = maintainable.

---

### Section 5 — Adaptive Implementation via FPW + Hand-Roll

Widget dengan native Cupertino counterpart (~12 widget) pakai **flutter_platform_widgets** (`Platform.*`):

```dart
// BcInput adaptive via FPW:
class BcInput extends StatelessWidget {
  const BcInput({
    super.key,
    this.label,
    this.value,
    this.onValueChange,
  });
  final String? label;
  final String? value;
  final void Function(String, BcChangeMeta)? onValueChange;

  @override
  Widget build(BuildContext context) {
    return PlatformTextField(
      material: (_, __) => MaterialTextFieldData(
        decoration: InputDecoration(labelText: label),
        onChanged: (v) => onValueChange?.call(v, BcChangeMeta(reason: BcChangeReason.input, source: BcEventSource.pointer)),
      ),
      cupertino: (_, __) => CupertinoTextFieldData(
        placeholder: label,
        onChanged: (v) => onValueChange?.call(v, BcChangeMeta(reason: BcChangeReason.input, source: BcEventSource.pointer)),
      ),
    );
  }
}
```

Widget tanpa Cupertino counterpart (~10 widget: DataTable, Splitter, Accordion, dll) **stay Material-only** di semua platform:

```dart
// BcDataTable: Material-only, no adaptive:
class BcDataTable extends StatelessWidget {
  // ... pure Material 3 DataGrid
}
```

**Contract impact: ZERO.** `.contract.ts` tetap Material-canonical dengan satu `flutter:` mapping string (descriptive, bukan binding codegen). ADR-008 freeze tidak violated. Tidak ada dual mapping string per-platform.

**Rationale:**

| Alternatif | Kenapa ditolak |
|---|---|
| Hand-roll semua 22 widget | 12 widget yang ada Cupertino equivalent akan terlihat non-native di iOS. |
| FPW untuk semua 22 widget | FPW tidak punya Cupertino counterpart untuk DataTable/Splitter/Accordion. Pakai Material via FPW tetap Material. |
| Contract-level adaptive mapping | Membengkakkan contract, mengganggu Mitosis, ADR-008 freeze violated. |

---

### Section 6 — Enum-Based Event Types (Non-Generic BcChangeMeta)

```dart
enum BcEventSource { pointer, keyboard, programmatic, external }

enum BcChangeReason { input, reaction, programmatic, select, toggle, submit }

class BcChangeMeta {
  final BcChangeReason reason;
  final BcEventSource source;
  final Object? nativeEvent;

  const BcChangeMeta({
    required this.reason,
    this.source = BcEventSource.programmatic,
    this.nativeEvent,
  });
}
```

**Rationale:**

| Alternatif | Kenapa ditolak |
|---|---|
| (A) Generic `BcChangeMeta<TReason extends String>` | Dart enum lebih type-safe dari string literal union. Generic menambah complexity di consumer. |
| (B) String-based reason seperti TS | Dart enum = exhaustiveness check di switch. Lebih robust. |

Non-generic karena: (1) consumer tidak perlu annotate `BcChangeMeta<'input'>`, (2) enum sudah type-safe, (3) mirror TS contract tanpa generic overhead.

---

## Consequences

### Positive

1. **Consumer DX = 100% identik** across React/RN/Vue/Flutter. 1 API, 5 target.
2. **Core portable** — `reactive_store.dart` reusable di CLI/server/future Dart targets.
3. **Const-first** — `static const Color` = compile-time, zero allocation, idiomatic Flutter.
4. **Platform-aware by default** — `BcApp(home: ...)` satu line, auto Cupertino/Material.
5. **Contract untouched** — ADR-008 freeze preserved. Flutter mapping string descriptive only.
6. **Token parity** — `Color.withValues(alpha:)` = Tailwind `hsl(var(--x) / 0.5)`.
7. **Adaptive where it matters** — iOS user gets CupertinoTextField, not Material TextField.
8. **Future generator-ready** — Phase 5b generator stub 80% simple widget dari contract.

### Negative

1. **25% library-dev overhead** — `bc_xx.dart` written separately from `.lite.tsx`. Tidak ada compiler TSX→Dart yang reliable.
2. **FPW dependency** — `flutter_platform_widgets` added to `pubspec.yaml`. Lock-in ke library eksternal.
3. **Return type duality** — `BcTheme` expose 4 factory methods (2 ThemeData + 2 CupertinoThemeData). Consumer yang pakai `CupertinoApp` explicit harus import `cupertinoLight()`.
4. **No runtime HSL mutation** — Hex const tidak bisa di-mutate runtime untuk dynamic theming. Token change = re-run generator. Acceptable: BitCode tidak punya use case runtime hue-shift.
5. **Hand-write P0 first** — 22 widget hand-written sebelum generator (Phase 5a). Effort upfront, tapi validate pattern dengan kode nyata.

### Trade-off Assessment

Trade-off utama: **25% library-dev overhead** vs **100% consumer DX parity**. Prinsip project: "DX user prioritized over library developer convenience." Keputusan ini selaras.

FPW lock-in mitigated: FPW thin wrapper over Flutter native widgets. Kalau FPW deprecated, migration path = replace `Platform.*` dengan native widget + manual `Platform.isIOS` check. Bukan rewrite.

---

## Acceptance Criteria

Phase 5a (current scope):

- [ ] `packages/ui-flutter/` scaffolded with `pubspec.yaml`, barrel export, directory structure
- [ ] `generate-flutter-tokens.mjs` parses `theme.css` → emits `bc_theme.dart` with 30+ `const Color`
- [ ] `bc_theme.dart` has `BcTheme` class with const colors + 4 factory methods (materialLight/Dark, cupertinoLight/Dark)
- [ ] `bc_app.dart` has `BcApp` widget wrapping `PlatformApp` with BcTheme defaults
- [ ] `types.dart` has `BcEventSource` enum, `BcChangeReason` enum, `BcChangeMeta` class (non-generic)
- [ ] `reactive_store.dart` (pure Dart) passes unit tests without Flutter binding
- [ ] `reactive_store_notifier.dart` (Flutter adapter) passes widget tests
- [ ] 22 P0 widgets implemented (12 FPW adaptive + 10 Material-only)
- [ ] Consumer API matches `.contract.ts` for all 22 widgets
- [ ] `flutter test` passes for all widget tests
- [ ] `dart analyze` clean (zero warnings)
- [ ] ADR-011 referenced in `docs/adr/README.md` index

Phase 5b (deferred):

- [ ] Generator stubs 80% simple widget dari contract
- [ ] P1 (~60 widget) implemented
- [ ] P2 (~90 widget) implemented

---

## Implementation Phases

### Phase 5a — P0 Hand-Write (current)

1. Scaffold `packages/ui-flutter/` — `pubspec.yaml` + barrel + directory structure
2. `packages/ui/scripts/generate-flutter-tokens.mjs` — parse `theme.css` → `bc_theme.dart`
3. Core types — `types.dart` (BcChangeMeta, BcEventSource, BcChangeReason)
4. Reactive store — `reactive_store.dart` (pure Dart) + `reactive_store_notifier.dart` (Flutter adapter)
5. Token bridge — `bc_theme.dart` (3-layer: colors → explicit → app)
6. `bc_app.dart` — `BcApp` widget (PlatformApp wrapper)
7. 22 P0 widgets with adaptive FPW + hand-roll
8. Tests
9. Docs update (codebase.md, components-progress.md)
10. Commit

### Phase 5b — Generator + P1

1. Contract → Dart stub generator
2. P1 (~60 widget) implemented via generator + hand-write complex
3. ADR-011 acceptance criteria Phase 5b checked

### Phase 5c — P2 + Sisa

1. P2 (~90 widget) implemented
2. Full parity with web side
3. ADR-011 acceptance criteria Phase 5c checked

---

## Related ADRs

- [ADR-002: Cross-Platform Architecture](./002-cross-platform-architecture.md) — Defines Tier C (Flutter) as separate track. ADR-011 implements Tier C.
- [ADR-007: Enterprise Universal UI API Contract](./007-enterprise-universal-ui-api-contract.md) — Defines semantic event/state/method contracts. ADR-011 ports these to Dart.
- [ADR-008: Contract-First API Freeze](./008-contract-first-api-freeze.md) — Freezes contract. ADR-011 adds Flutter mapping string without modifying frozen surfaces.
- [ADR-004: Reactive Engine](./004-reactive-engine.md) — Defines reactive store semantics. ADR-011 ports to pure Dart.
- [ADR-005: Metadata-Driven UI](./005-metadata-driven-ui.md) — Defines UiSchema renderer. ADR-011 ports to Dart (Phase 5a scope: types only, renderer Phase 5b).

## References

- `packages/ui/tokens/theme.css` — HSL token source of truth
- `packages/ui/src/core/store.ts` — Reactive store (TS source for port)
- `packages/ui/src/core/types.ts` — BcChangeMeta, BcEventSource types (TS source for port)
- `packages/ui/src/core/renderer/BcRenderer.contract.ts` — Schema renderer contract
- [flutter_platform_widgets](https://pub.dev/packages/flutter_platform_widgets) — FPW library
- [Flutter Color.withValues()](https://api.flutter.dev/flutter/dart-ui/Color/withValues.html) — Native alpha manipulation
- [MaterialBasedCupertinoThemeData](https://api.flutter.dev/flutter/cupertino/MaterialBasedCupertinoThemeData-class.html) — Material → Cupertino theme bridge
