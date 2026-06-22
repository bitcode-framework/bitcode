/// Core types ported from @bitcode/ui-schema (ADR-011 Section 6).
///
/// Enum-based event types (non-generic) for Dart type safety.
/// Mirrors the TS BcChangeMeta/BcEventSource/BcChangeHandler contract
/// without generic overhead.

library core_types;

/// Source of a user-initiated change.
///
/// Mirrors TS `BcEventSource` but trimmed to Dart-relevant values.
/// `touch` is folded into `pointer` (Flutter gesture system unifies them).
enum BcEventSource {
  pointer,
  keyboard,
  programmatic,
  external,
}

/// Reason metadata for change events.
///
/// Mirrors TS string literal reasons as a Dart enum for exhaustiveness
/// checking in switch statements (ADR-011 Section 6).
enum BcChangeReason {
  input,
  reaction,
  programmatic,
  select,
  toggle,
  submit,
}

/// Metadata attached to every change event.
///
/// Non-generic (ADR-011 Section 6): consumer does not annotate
/// `BcChangeMeta<'input'>`. The enum already provides type safety.
class BcChangeMeta {
  final BcChangeReason reason;
  final BcEventSource source;
  final Object? nativeEvent;

  const BcChangeMeta({
    required this.reason,
    this.source = BcEventSource.programmatic,
    this.nativeEvent,
  });

  @override
  String toString() =>
      'BcChangeMeta(reason: $reason, source: $source, nativeEvent: $nativeEvent)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcChangeMeta &&
          reason == other.reason &&
          source == other.source;

  @override
  int get hashCode => Object.hash(reason, source);
}

/// Callback signature for value changes.
///
/// Mirrors TS `BcChangeHandler<TValue, TReason>`.
/// Dart uses generics only for the value type; reason is enum-typed.
typedef BcChangeHandler<TValue> =
    void Function(TValue value, BcChangeMeta meta);

/// Key type for selection and collection identifiers.
typedef BcKey = Object; // String or int

/// Component size variants.
enum BcComponentSize { sm, md, lg }

/// Component visual variants.
enum BcComponentVariant {
  defaultVariant,
  primary,
  secondary,
  destructive,
  outline,
  ghost,
  link,
}
