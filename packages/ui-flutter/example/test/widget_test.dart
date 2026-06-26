// This is a basic Flutter widget test.
//
// To perform an interaction with a widget in your test, use the WidgetTester
// utility in the flutter_test package. For example, you can send tap and scroll
// gestures. You can also use WidgetTester to find child widgets in the widget
// tree, read text, and verify that the values of widget properties are correct.

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
}
