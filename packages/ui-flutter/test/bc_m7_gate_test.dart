// M7 Final Parity Gate Tests
//
// Covers: VAL-M7-GATE-001 through VAL-M7-GATE-008 and
//        VAL-CROSS-001 through VAL-CROSS-005.
//
// These tests verify the final integration gate for ADR-011 Phase 5c.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:bc_ui_flutter/bc_ui_flutter.dart';

void main() {
  // ============================================================
  // VAL-M7-GATE-001: Full barrel export audit
  // ============================================================
  group('VAL-M7-GATE-001: Full barrel export audit', () {
    test('all P0 foundation widgets are exported', () {
      expect(BcSlot, isNotNull);
      expect(BcStack, isNotNull);
      expect(BcXStack, isNotNull);
      expect(BcYStack, isNotNull);
      expect(BcBox, isNotNull);
      expect(BcVisuallyHidden, isNotNull);
      expect(BcBackdrop, isNotNull);
      expect(BcLoadingOverlay, isNotNull);
    });

    test('all P0 display widgets are exported', () {
      expect(BcDescriptions, isNotNull);
      expect(BcSteps, isNotNull);
      expect(BcStatistic, isNotNull);
      expect(BcEmpty, isNotNull);
      expect(BcResult, isNotNull);
    });

    test('all P0 advanced input widgets are exported', () {
      expect(BcOtp, isNotNull);
      expect(BcTransfer, isNotNull);
      expect(BcCascader, isNotNull);
      expect(BcTreeSelect, isNotNull);
    });

    test('all M2 forms text input widgets are exported', () {
      expect(BcRichtext, isNotNull);
      expect(BcCode, isNotNull);
      expect(BcJson, isNotNull);
      expect(BcMarkdown, isNotNull);
    });

    test('all M2 forms picker/specialized widgets are exported', () {
      expect(BcColorPicker, isNotNull);
      expect(BcRate, isNotNull);
      expect(BcSlider, isNotNull);
      expect(BcRangeSlider, isNotNull);
      expect(BcTagsInput, isNotNull);
      expect(BcToggleGroup, isNotNull);
      expect(BcMentions, isNotNull);
      expect(BcKnob, isNotNull);
    });

    test('all P0 adaptive widgets are exported', () {
      expect(BcInput, isNotNull);
      expect(BcSelect, isNotNull);
      expect(BcCheckbox, isNotNull);
      expect(BcToggle, isNotNull);
      expect(BcSegmented, isNotNull);
      expect(BcNumberField, isNotNull);
      expect(BcDateField, isNotNull);
      expect(BcTimeField, isNotNull);
      expect(BcDialog, isNotNull);
      expect(BcAlertDialog, isNotNull);
      expect(BcPortal, isNotNull);
      expect(BcTooltip, isNotNull);
      expect(BcDropdown, isNotNull);
    });

    test('all P0 Material-only widgets are exported', () {
      expect(BcCombobox, isNotNull);
      expect(BcMultiSelect, isNotNull);
      expect(BcForm, isNotNull);
      expect(BcDrawer, isNotNull);
      expect(BcPopover, isNotNull);
      expect(BcContextMenu, isNotNull);
      expect(BcCommand, isNotNull);
      expect(BcTabs, isNotNull);
      expect(BcAccordion, isNotNull);
      expect(BcDisclosure, isNotNull);
      expect(BcSplitter, isNotNull);
      expect(BcDataTable, isNotNull);
      expect(BcList, isNotNull);
      expect(BcFieldDuration, isNotNull);
      expect(BcFieldTableSelect, isNotNull);
      expect(BcFieldMorph, isNotNull);
      expect(BcFieldDynLink, isNotNull);
      expect(BcLookupModal, isNotNull);
      expect(BcViewGantt, isNotNull);
    });

    test('all M3 navigation widgets are exported', () {
      expect(BcBreadcrumb, isNotNull);
      expect(BcPagination, isNotNull);
      expect(BcMenubar, isNotNull);
      expect(BcMegamenu, isNotNull);
      expect(BcNavigationMenu, isNotNull);
      expect(BcAffix, isNotNull);
      expect(BcAnchor, isNotNull);
      expect(BcBacktop, isNotNull);
      expect(BcDock, isNotNull);
      expect(BcPanelMenu, isNotNull);
      expect(BcSlideMenu, isNotNull);
    });

    test('all M4 overlay widgets are exported', () {
      expect(BcPrompt, isNotNull);
      expect(BcWizard, isNotNull);
      expect(BcHoverCard, isNotNull);
      expect(BcActionsheet, isNotNull);
      expect(BcSpeedDial, isNotNull);
      expect(BcNotification, isNotNull);
      expect(BcTour, isNotNull);
      expect(BcInplace, isNotNull);
      expect(BcMeterGroup, isNotNull);
      expect(BcAspectRatio, isNotNull);
      expect(BcFlex, isNotNull);
    });

    test('all M2 upload/capture widgets are exported', () {
      expect(BcFile, isNotNull);
      expect(BcSignature, isNotNull);
      expect(BcBarcode, isNotNull);
      expect(BcGeo, isNotNull);
      expect(BcDropZone, isNotNull);
      expect(BcColorArea, isNotNull);
      expect(BcColorSlider, isNotNull);
    });

    test('all M5 data display widgets are exported', () {
      expect(BcCarousel, isNotNull);
      expect(BcOrgchart, isNotNull);
      expect(BcTagGroup, isNotNull);
      expect(BcSparkline, isNotNull);
      expect(BcChart, isNotNull);
      expect(BcKanbanBoard, isNotNull);
      expect(BcViewMap, isNotNull);
      expect(BcViewReport, isNotNull);
      expect(BcViewEditor, isNotNull);
      expect(BcViewActivity, isNotNull);
      expect(BcVirtualScroller, isNotNull);
      expect(BcTieredMenu, isNotNull);
    });

    test('all M6 media widgets are exported', () {
      expect(BcViewerPdf, isNotNull);
      expect(BcViewerImage, isNotNull);
      expect(BcViewerDocument, isNotNull);
      expect(BcViewerVideo, isNotNull);
      expect(BcViewerAudio, isNotNull);
      expect(BcViewerYoutube, isNotNull);
      expect(BcViewerInstagram, isNotNull);
      expect(BcViewerTiktok, isNotNull);
    });

    test('all M6 social widgets are exported', () {
      expect(BcChatter, isNotNull);
      expect(BcActivity, isNotNull);
      expect(BcTimelineSocial, isNotNull);
      expect(BcGalleria, isNotNull);
      expect(BcImage, isNotNull);
    });

    test('all M6 print widgets are exported', () {
      expect(BcExport, isNotNull);
      expect(BcPrint, isNotNull);
      expect(BcReportLink, isNotNull);
    });

    test('all M7 misc widgets are exported', () {
      expect(BcWatermark, isNotNull);
      expect(BcQrcode, isNotNull);
      expect(BcKbd, isNotNull);
      expect(BcCenter, isNotNull);
      expect(BcCircle, isNotNull);
      expect(BcSquare, isNotNull);
      expect(BcColorSwatch, isNotNull);
      expect(BcColorWheel, isNotNull);
      expect(BcSearch, isNotNull);
      expect(BcSidebar, isNotNull);
      expect(BcFavorites, isNotNull);
    });

    test('M1 overlay widget (BcToastQueue) is exported', () {
      expect(BcToastQueue, isNotNull);
    });

    test('all core types and enums are exported', () {
      expect(BcChangeReason.values.length, 18);
      expect(BcEventSource.values.length, 3);
      expect(BcGanttTask, isNotNull);
      expect(BcGanttLink, isNotNull);
      expect(BcGanttTaskType.values.length, 3);
      expect(BcGanttZoomLevel.values.length, 4);
      expect(BcPortalAnimation.values.length, 3);
      expect(BcBread
