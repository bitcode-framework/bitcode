/// BitCode UI Flutter adapter.
///
/// Contract-driven Dart widgets ported from @bitcode/ui (ADR-011).
/// Consumer API is identical across React/Vue/RN/Flutter — only the
/// implementation language differs.
///
/// Simple usage (platform-aware, zero config):
/// ```dart
/// void main() => runApp(BcApp(home: MyScreen()));
/// ```
///
/// Raw token access:
/// ```dart
/// Container(color: BcTheme.primary.withValues(alpha: 0.5));
/// ```

library bc_ui_flutter;

// Core types
export 'src/core/types.dart';

// Reactive engine
export 'src/core/reactive_store.dart';
export 'src/core/reactive_store_notifier.dart';

// UiSchema types (ADR-011 Phase 5b.1)
// BcKey and DataSource are defined in core/types.dart and core/reactive_store.dart;
// hide them here to avoid ambiguous_export conflicts.
export 'src/schema/ui_schema_types.dart' hide BcKey, DataSource;
export 'src/schema/ui_schema_events.dart';
export 'src/schema/ui_schema.dart';
export 'src/schema/field_resolver.dart';

// Theme
export 'src/theme/bc_theme.dart';

// App
export 'src/bc_app.dart';

// M1 layout primitives
export 'src/components/bc_slot.dart';
export 'src/components/bc_stack.dart';
export 'src/components/bc_xstack.dart';
export 'src/components/bc_ystack.dart';
export 'src/components/bc_box.dart';
export 'src/components/bc_visually_hidden.dart';
export 'src/components/bc_backdrop.dart';
export 'src/components/bc_loading_overlay.dart';

// M1 display widgets
export 'src/components/bc_descriptions.dart';
export 'src/components/bc_steps.dart';
export 'src/components/bc_statistic.dart';
export 'src/components/bc_empty.dart';
export 'src/components/bc_result.dart';

// M1 advanced input widgets
export 'src/components/bc_otp.dart';
export 'src/components/bc_transfer.dart';
export 'src/components/bc_cascader.dart';
export 'src/components/bc_tree_select.dart';

// M1 overlays
export 'src/components/bc_toast_queue.dart';

// P0 widgets — adaptive (FPW)
export 'src/components/bc_input.dart';
export 'src/components/bc_select.dart';
export 'src/components/bc_checkbox.dart';
export 'src/components/bc_toggle.dart';
export 'src/components/bc_segmented.dart';
export 'src/components/bc_number_field.dart';
export 'src/components/bc_date_field.dart';
export 'src/components/bc_time_field.dart';
export 'src/components/bc_dialog.dart';
export 'src/components/bc_alert_dialog.dart';
export 'src/components/bc_portal.dart';
export 'src/components/bc_tooltip.dart';
export 'src/components/bc_dropdown.dart';

// P0 widgets — Material-only
export 'src/components/bc_combobox.dart';
export 'src/components/bc_multi_select.dart';
export 'src/components/bc_form.dart';
export 'src/components/bc_drawer.dart';
export 'src/components/bc_popover.dart';
export 'src/components/bc_context_menu.dart';
export 'src/components/bc_command.dart';
export 'src/components/bc_tabs.dart';
export 'src/components/bc_accordion.dart';
export 'src/components/bc_disclosure.dart';
export 'src/components/bc_splitter.dart';
export 'src/components/bc_data_table.dart';
export 'src/components/bc_list.dart';
export 'src/components/bc_field_duration.dart';
export 'src/components/bc_field_table_select.dart';
export 'src/components/bc_field_morph.dart';
export 'src/components/bc_field_dyn_link.dart';
export 'src/components/bc_lookup_modal.dart';
export 'src/components/bc_view_gantt.dart';
