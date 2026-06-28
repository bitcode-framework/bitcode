// This is a basic Flutter widget test.
//
// To perform an interaction with a widget in your test, use the WidgetTester
// utility in the flutter_test package. For example, you can send tap and scroll
// gestures. You can also use WidgetTester to find child widgets in the widget
// tree, read text, and verify that the values of widget properties are correct.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:bc_ui_flutter_example/main.dart';

void main() {
  testWidgets('Demo app renders with sidebar and content', (WidgetTester tester) async {
    await tester.pumpWidget(const BcUiDemoApp());
    await tester.pumpAndSettle(const Duration(seconds: 3));

    // The first section (Forms) should be visible by default
    expect(find.text('Forms'), findsWidgets);
    expect(find.byType(BcUiDemoApp), findsOneWidget);
  });

  // Regression test for ADR-016 Demo Bug 1: data_section.dart used ListView
  // as its root widget, but main.dart wraps every section in a
  // SingleChildScrollView — ListView inside SingleChildScrollView crashes with
  // "Vertical viewport was given unbounded height". Fix: ListView → Column.
  // This test taps the Data sidebar item and confirms no ErrorWidget appears.
  testWidgets('Data section renders without ErrorWidget (ADR-016 Demo Bug 1)',
      (WidgetTester tester) async {
    await tester.pumpWidget(const BcUiDemoApp());
    await tester.pumpAndSettle(const Duration(seconds: 2));

    final dataTarget = find.text('Data');
    if (dataTarget.evaluate().isNotEmpty) {
      await tester.tap(dataTarget.first);
      await tester.pumpAndSettle(const Duration(seconds: 1));
    }
    expect(find.byType(ErrorWidget), findsNothing,
        reason: 'Data section rendered an ErrorWidget (crash)');
    // Confirm the Data section's content is actually present.
    expect(find.text('BcDataTable'), findsOneWidget);
  });
}
