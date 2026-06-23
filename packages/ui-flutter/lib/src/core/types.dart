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
///
/// P0 reasons (indices 0-5) are preserved for backward compatibility.
/// P1 reasons (indices 6-17) extend the enum for the 8 P1 widgets.
/// TS kebab-case names map to Dart camelCase:
///   type-change -> typeChange
///   value-change -> valueChange
///   model-change -> modelChange
///   overlay-click -> overlayClick
enum BcChangeReason {
  input,
  reaction,
  programmatic,
  select,
  toggle,
  submit,
  escape,
  back,
  deselect,
  clear,
  typeChange,
  valueChange,
  modelChange,
  confirm,
  cancel,
  overlayClick,
  increment,
  decrement,
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

// ---------------------------------------------------------------------------
// P1 typed wrapper classes (ADR-011 Phase 5b.2)
//
// These mirror the TS contract interfaces for composite value types used by
// the 8 P1 widgets. They are plain immutable classes with const constructors
// and value-based equality so widget tests can compare instances directly.
// ---------------------------------------------------------------------------

/// Column descriptor for [BcFieldTableSelect].
///
/// Mirrors TS `BcFieldTableSelectColumn` ({ key, label }).
class BcFieldTableSelectColumn {
  final String key;
  final String label;

  const BcFieldTableSelectColumn({required this.key, required this.label});

  @override
  String toString() =>
      'BcFieldTableSelectColumn(key: $key, label: $label)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcFieldTableSelectColumn &&
          key == other.key &&
          label == other.label;

  @override
  int get hashCode => Object.hash(key, label);
}

/// Single selectable option for a morph `select` widget kind.
///
/// Mirrors TS `BcFieldMorphOption` ({ label, value }).
class BcFieldMorphOption {
  final String label;
  final String value;

  const BcFieldMorphOption({required this.label, required this.value});

  @override
  String toString() => 'BcFieldMorphOption(label: $label, value: $value)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcFieldMorphOption &&
          label == other.label &&
          value == other.value;

  @override
  int get hashCode => Object.hash(label, value);
}

/// Type descriptor for [BcFieldMorph].
///
/// Mirrors TS `BcFieldMorphType` ({ value, label, widget, options? }).
/// The [widget] field is a string tag ('text', 'number', 'date', 'select',
/// 'boolean') that the morph widget maps to a concrete Flutter input.
class BcFieldMorphType {
  final String value;
  final String label;
  final String widget;
  final List<BcFieldMorphOption>? options;

  const BcFieldMorphType({
    required this.value,
    required this.label,
    required this.widget,
    this.options,
  });

  @override
  String toString() =>
      'BcFieldMorphType(value: $value, label: $label, widget: $widget, options: $options)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcFieldMorphType &&
          value == other.value &&
          label == other.label &&
          widget == other.widget &&
          _listEquals(options, other.options);

  @override
  int get hashCode => Object.hash(value, label, widget, options);

  static bool _listEquals(List<BcFieldMorphOption>? a, List<BcFieldMorphOption>? b) {
    if (identical(a, b)) return true;
    if (a == null || b == null) return a == b;
    if (a.length != b.length) return false;
    for (var i = 0; i < a.length; i++) {
      if (a[i] != b[i]) return false;
    }
    return true;
  }
}

/// Composite value for [BcFieldMorph].
///
/// Mirrors TS `BcFieldMorphValue` ({ type, value }).
/// The [value] is `Object?` because the concrete type varies by widget
/// kind (String for text, num for number, bool for boolean, etc.).
class BcFieldMorphValue {
  final String type;
  final Object? value;

  const BcFieldMorphValue({required this.type, this.value});

  @override
  String toString() => 'BcFieldMorphValue(type: $type, value: $value)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcFieldMorphValue && type == other.type && value == other.value;

  @override
  int get hashCode => Object.hash(type, value);
}

/// Model descriptor for [BcFieldDynLink].
///
/// Mirrors TS `BcDynLinkModel` ({ model, label }).
class BcDynLinkModel {
  final String model;
  final String label;

  const BcDynLinkModel({required this.model, required this.label});

  @override
  String toString() => 'BcDynLinkModel(model: $model, label: $label)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcDynLinkModel &&
          model == other.model &&
          label == other.label;

  @override
  int get hashCode => Object.hash(model, label);
}

/// Composite value for [BcFieldDynLink].
///
/// Mirrors TS `BcDynLinkValue` ({ model, id, display }).
class BcDynLinkValue {
  final String model;
  final String id;
  final String display;

  const BcDynLinkValue({
    required this.model,
    required this.id,
    required this.display,
  });

  @override
  String toString() =>
      'BcDynLinkValue(model: $model, id: $id, display: $display)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcDynLinkValue &&
          model == other.model &&
          id == other.id &&
          display == other.display;

  @override
  int get hashCode => Object.hash(model, id, display);
}

/// Record returned by a [BcFieldDynLink] data fetcher.
///
/// Mirrors TS `BcDynLinkRecord` ({ id, label }).
class BcDynLinkRecord {
  final String id;
  final String label;

  const BcDynLinkRecord({required this.id, required this.label});

  @override
  String toString() => 'BcDynLinkRecord(id: $id, label: $label)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcDynLinkRecord && id == other.id && label == other.label;

  @override
  int get hashCode => Object.hash(id, label);
}

/// Column descriptor for [BcLookupModal].
///
/// Mirrors TS `BcLookupModalColumn` ({ key, label }).
/// Intentionally parallels [BcFieldTableSelectColumn] shape.
class BcLookupModalColumn {
  final String key;
  final String label;

  const BcLookupModalColumn({required this.key, required this.label});

  @override
  String toString() => 'BcLookupModalColumn(key: $key, label: $label)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcLookupModalColumn &&
          key == other.key &&
          label == other.label;

  @override
  int get hashCode => Object.hash(key, label);
}

// ---------------------------------------------------------------------------
// Gantt domain types (ADR-011 Phase 5b.2 — BcViewGantt)
// ---------------------------------------------------------------------------

/// Rendering shape for a Gantt task bar.
///
/// Mirrors TS `BcGanttTaskType` ('task' | 'milestone' | 'summary').
enum BcGanttTaskType {
  task,
  milestone,
  summary,
}

/// Zoom level for the Gantt timeline.
///
/// Mirrors TS `BcGanttZoomLevel` ('day' | 'week' | 'month' | 'quarter').
enum BcGanttZoomLevel {
  day,
  week,
  month,
  quarter,
}

/// Single Gantt task.
///
/// Mirrors TS `BcGanttTask`. [start] and [duration] are in time units
/// relative to the timeline start. [progress] is 0-100.
class BcGanttTask {
  final String id;
  final String label;
  final double start;
  final double duration;
  final double? progress;
  final String? color;
  final String? parentId;
  final bool? collapsed;
  final BcGanttTaskType? type;
  final List<String>? assignees;

  const BcGanttTask({
    required this.id,
    required this.label,
    required this.start,
    required this.duration,
    this.progress,
    this.color,
    this.parentId,
    this.collapsed,
    this.type,
    this.assignees,
  });

  @override
  String toString() => 'BcGanttTask(id: $id, label: $label, start: $start, '
      'duration: $duration, progress: $progress, color: $color, '
      'parentId: $parentId, collapsed: $collapsed, type: $type, '
      'assignees: $assignees)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcGanttTask &&
          id == other.id &&
          label == other.label &&
          start == other.start &&
          duration == other.duration &&
          progress == other.progress &&
          color == other.color &&
          parentId == other.parentId &&
          collapsed == other.collapsed &&
          type == other.type &&
          _listEquals(assignees, other.assignees);

  @override
  int get hashCode => Object.hash(
        id,
        label,
        start,
        duration,
        progress,
        color,
        parentId,
        collapsed,
        type,
        Object.hashAll(assignees ?? const []),
      );

  static bool _listEquals(List<String>? a, List<String>? b) {
    if (identical(a, b)) return true;
    if (a == null || b == null) return a == b;
    if (a.length != b.length) return false;
    for (var i = 0; i < a.length; i++) {
      if (a[i] != b[i]) return false;
    }
    return true;
  }
}

/// Dependency link between two Gantt tasks.
///
/// Mirrors TS `BcGanttLink`. The [type] field accepts 'FS', 'FF', 'SS',
/// 'SF' (validated at usage layer, not constructor).
class BcGanttLink {
  final String id;
  final String from;
  final String to;
  final String? type;

  const BcGanttLink({
    required this.id,
    required this.from,
    required this.to,
    this.type,
  });

  @override
  String toString() =>
      'BcGanttLink(id: $id, from: $from, to: $to, type: $type)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcGanttLink &&
          id == other.id &&
          from == other.from &&
          to == other.to &&
          type == other.type;

  @override
  int get hashCode => Object.hash(id, from, to, type);
}

// ---------------------------------------------------------------------------
// M2 upload/capture typed wrappers (ADR-011 Phase 5c — M2 forms-upload)
// ---------------------------------------------------------------------------

/// File reference for [BcFile] and [BcDropZone].
///
/// Mirrors a browser `File` object's salient fields used by the filtering
/// and display logic. [path] is set on native (desktop/mobile) platforms.
class BcFileRef {
  final String name;
  final int size;
  final String? type;
  final String? path;

  const BcFileRef({
    required this.name,
    required this.size,
    this.type,
    this.path,
  });

  @override
  String toString() =>
      'BcFileRef(name: $name, size: $size, type: $type, path: $path)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcFileRef &&
          name == other.name &&
          size == other.size &&
          type == other.type &&
          path == other.path;

  @override
  int get hashCode => Object.hash(name, size, type, path);
}

/// Geographic coordinate pair for [BcGeo].
///
/// Mirrors TS `{ lat: number; lng: number }`.
class BcGeoCoord {
  final double lat;
  final double lng;

  const BcGeoCoord({required this.lat, required this.lng});

  @override
  String toString() => 'BcGeoCoord(lat: $lat, lng: $lng)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcGeoCoord && lat == other.lat && lng == other.lng;

  @override
  int get hashCode => Object.hash(lat, lng);
}

/// Saturation/value (brightness) pair for [BcColorArea].
///
/// Maps the 2D drag area's horizontal axis (saturation, 0–100) and
/// vertical axis (value, 0–100).
class BcSvValue {
  final double saturation;
  final double value;

  const BcSvValue({required this.saturation, required this.value});

  @override
  String toString() =>
      'BcSvValue(saturation: $saturation, value: $value)';

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is BcSvValue &&
          saturation == other.saturation &&
          value == other.value;

  @override
  int get hashCode => Object.hash(saturation, value);
}

// ---------------------------------------------------------------------------
// M5 data-display typed wrappers (ADR-011 Phase 5c — M5)
// ---------------------------------------------------------------------------

/// Rendering style for [BcSparkline].
enum BcSparklineType {
  line,
  bar,
  area,
}

/// Visual variant for a tag in [BcTagGroup].
enum BcTagVariant {
  defaultVariant,
  success,
  warning,
  destructive,
  info,
}

