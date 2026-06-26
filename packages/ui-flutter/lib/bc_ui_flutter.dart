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

// Chart option contract (ADR-015 — Dart port of the TS chart-option factory).
// Pure functions/enums, no widget deps — testable in isolation.
export 'src/core/chart/chart_contract.dart';

// Reactive engine
export 'src/core/reactive_store.dart';
export 'src/core/reactive_store_notifier.dart';

// AppProvider consumption surface (ADR-013 Phase 1) — cross-tier parity with
// @bitcode/ui core/providers + core/hooks/use-app-context.
export 'src/core/app_config_snapshot.dart';
export 'src/core/app_config_store.dart';
export 'src/core/bc_app_config.dart';

// Hooks P1 parity (ADR-013 Phase 5b) — Dart ports of the TS P1 headless hooks
// (packages/ui/src/core/hooks/*). Headless utility classes / factories.
export 'src/core/hooks/hooks.dart';

// Animation runtime parity (ADR-013 Phase 5b) — Dart port of the TS
// `core/animation` runtime (resolveAnimation, BcAnimationKind, reduced-motion).
export 'src/core/animation/animation.dart';

// Adapt runtime parity (ADR-013 Phase 5b) — Dart port of the TS `core/adapt`
// runtime (resolveAdapt, BcAdaptResolver, BcAdaptContext). Shared breakpoint
// queries (BC_BREAKPOINTS parity) live in bc_breakpoints.dart.
export 'src/core/adapt/adapt.dart';
export 'src/core/bc_breakpoints.dart';

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

// M2 forms — text inputs
export 'src/components/bc_richtext.dart';
export 'src/components/bc_code.dart';
export 'src/components/bc_json.dart';
export 'src/components/bc_markdown.dart';

// M2 forms — pickers and specialized inputs
export 'src/components/bc_color_picker.dart';
export 'src/components/bc_rate.dart';
export 'src/components/bc_slider.dart';
export 'src/components/bc_range_slider.dart';
export 'src/components/bc_tags_input.dart';
export 'src/components/bc_toggle_group.dart';
export 'src/components/bc_mentions.dart';
export 'src/components/bc_knob.dart';

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

// M3 navigation
export 'src/components/bc_breadcrumb.dart';
export 'src/components/bc_pagination.dart';
export 'src/components/bc_menubar.dart';
export 'src/components/bc_megamenu.dart';
export 'src/components/bc_navigation_menu.dart';
export 'src/components/bc_affix.dart';
export 'src/components/bc_anchor.dart';
export 'src/components/bc_backtop.dart';
export 'src/components/bc_dock.dart';
export 'src/components/bc_panel_menu.dart';
export 'src/components/bc_slide_menu.dart';

// M4 overlays — modal
export 'src/components/bc_prompt.dart';
export 'src/components/bc_wizard.dart';

// M4 overlays — reveal
export 'src/components/bc_hover_card.dart';
export 'src/components/bc_actionsheet.dart';
export 'src/components/bc_speed_dial.dart';

// M4 feedback
export 'src/components/bc_notification.dart';
export 'src/components/bc_tour.dart';
export 'src/components/bc_inplace.dart';

// M4 layout utility
export 'src/components/bc_meter_group.dart';
export 'src/components/bc_aspect_ratio.dart';
export 'src/components/bc_flex.dart';

// M2 forms — upload and capture
export 'src/components/bc_file.dart';
export 'src/components/bc_signature.dart';
export 'src/components/bc_barcode.dart';
export 'src/components/bc_geo.dart';
export 'src/components/bc_drop_zone.dart';
export 'src/components/bc_color_area.dart';
export 'src/components/bc_color_slider.dart';

// M5 data display
export 'src/components/bc_carousel.dart';
export 'src/components/bc_orgchart.dart';
export 'src/components/bc_taggroup.dart';
export 'src/components/bc_sparkline.dart';
export 'src/components/bc_chart.dart';
export 'src/components/bc_kanban_board.dart';
export 'src/components/bc_view_map.dart';
export 'src/components/bc_view_report.dart';
export 'src/components/bc_view_editor.dart';
export 'src/components/bc_view_activity.dart';
export 'src/components/bc_virtual_scroller.dart';
export 'src/components/bc_tiered_menu.dart';

// M6 media viewers and players
export 'src/components/bc_viewer_pdf.dart';
export 'src/components/bc_viewer_image.dart';
export 'src/components/bc_viewer_document.dart';
export 'src/components/bc_viewer_video.dart';
export 'src/components/bc_viewer_audio.dart';
export 'src/components/bc_viewer_youtube.dart';
export 'src/components/bc_viewer_instagram.dart';
export 'src/components/bc_viewer_tiktok.dart';

// M6 social widgets
export 'src/components/bc_chatter.dart';
export 'src/components/bc_activity.dart';
export 'src/components/bc_timeline_social.dart';
export 'src/components/bc_galleria.dart';
export 'src/components/bc_image.dart';

// M6 print widgets
export 'src/components/bc_export.dart';
export 'src/components/bc_print.dart';
export 'src/components/bc_report_link.dart';

// M7 misc widgets
export 'src/components/bc_watermark.dart';
export 'src/components/bc_qrcode.dart';
export 'src/components/bc_kbd.dart';
export 'src/components/bc_center.dart';
export 'src/components/bc_circle.dart';
export 'src/components/bc_square.dart';
export 'src/components/bc_color_swatch.dart';
export 'src/components/bc_color_wheel.dart';
export 'src/components/bc_search.dart';
export 'src/components/bc_sidebar.dart';
export 'src/components/bc_favorites.dart';

// Renderer
export 'src/renderer/bc_field_builder.dart';
export 'src/renderer/bc_data_source_loader.dart';
export 'src/renderer/bc_renderer.dart';
