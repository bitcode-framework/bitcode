import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:bc_ui_flutter/src/components/bc_stack.dart';

void main() {
  Widget _harness(Widget child) {
    return MaterialApp(
      home: Scaffold(
        body: SizedBox(
          width: 300,
          height: 300,
          child: child,
        ),
      ),
    );
  }

  testWidgets('BcStack positions children using alignment and stackFit (VAL-M1-004)', (tester) async {
    await tester.pumpWidget(_harness(
      BcStack(
        alignment: Alignment.topLeft,
        stackFit: StackFit.expand,
        children: const [
          SizedBox(key: ValueKey('c1'), width: 50, height: 50),
          SizedBox(key: ValueKey('c2'), width: 50, height: 50),
        ],
      ),
    ));

    expect(find.byKey(const ValueKey('c1')), findsOneWidget);
    expect(find.byKey(const ValueKey('c2')), findsOneWidget);
    expect(find.descendant(of: find.byType(BcStack), matching: find.byType(Stack)), findsOneWidget);

    final c1Rect = tester.getRect(find.byKey(const ValueKey('c1')));
    expect(c1Rect.topLeft, equals(const Offset(0, 0)),
        reason: 'topLeft alignment should pin first child to top-left');

    final c2Rect = tester.getRect(find.byKey(const ValueKey('c2')));
    expect(c2Rect.size, equals(const Size(300, 300)),
        reason: 'StackFit.expand should make children fill stack bounds');
  });

  testWidgets('BcStack alignment change repositions children without rebuild errors (VAL-M1-005)', (tester) async {
    await tester.pumpWidget(_harness(
      BcStack(
        key: const ValueKey('stack'),
        alignment: Alignment.topLeft,
        children: const [
          SizedBox(key: ValueKey('c1'), width: 50, height: 50),
        ],
      ),
    ));
    await tester.pumpAndSettle();

    var c1Rect = tester.getRect(find.byKey(const ValueKey('c1')));
    expect(c1Rect.topLeft, equals(const Offset(0, 0)));

    await tester.pumpWidget(_harness(
      BcStack(
        key: const ValueKey('stack'),
        alignment: Alignment.bottomRight,
        children: const [
          SizedBox(key: ValueKey('c1'), width: 50, height: 50),
        ],
      ),
    ));
    await tester.pumpAndSettle();

    c1Rect = tester.getRect(find.byKey(const ValueKey('c1')));
    expect(c1Rect.bottomRight, equals(const Offset(300, 300)),
        reason: 'bottomRight alignment should pin child to bottom-right');
    expect(tester.takeException(), isNull);
  });

  testWidgets('BcStack renders empty without assertion (VAL-M1-006)', (tester) async {
    await tester.pumpWidget(_harness(
      const BcStack(),
    ));
    await tester.pump();
    expect(tester.takeException(), isNull);
    expect(find.descendant(of: find.byType(BcStack), matching: find.byType(Stack)), findsOneWidget);
  });
}
