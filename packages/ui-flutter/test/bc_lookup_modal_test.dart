import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:bc_ui_flutter/src/components/bc_lookup_modal.dart';
import 'package:bc_ui_flutter/src/components/bc_data_table.dart';
import 'package:bc_ui_flutter/src/core/types.dart';

/// Widget tests for [BcLookupModal] (ADR-011 Phase 5b.2).
///
/// Covers VAL-LOOKUP-001 through VAL-LOOKUP-035 plus VAL-CROSS-005.
void main() {
  Widget harness(Widget child) {
    return MaterialApp(home: Scaffold(body: child));
  }

  List<Map<String, dynamic>> sampleData() => <Map<String, dynamic>>[
        <String, dynamic>{
          'id': '1',
          'name': 'Alice',
          'email': 'alice@example.com'
        },
        <String, dynamic>{
          'id': '2',
          'name': 'Bob',
          'email': 'bob@example.com'
        },
        <String, dynamic>{
          'id': '3',
          'name': 'Charlie',
          'email': 'charlie@example.com'
        },
      ];

  List<BcLookupModalColumn> sampleColumns() => <BcLookupModalColumn>[
        const BcLookupModalColumn(key: 'name', label: 'Name'),
        const BcLookupModalColumn(key: 'email', label: 'Email'),
      ];

  // -----------------------------------------------------------------------
  // Rendering
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal does not render dialog when open is false (VAL-LOOKUP-001)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: false,
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsNothing);
  });

  testWidgets('BcLookupModal renders title when open (VAL-LOOKUP-002)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Pick Contact',
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);
    expect(find.text('Pick Contact'), findsOneWidget);
  });

  testWidgets('BcLookupModal renders column labels in order (VAL-LOOKUP-003)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.text('Name'), findsOneWidget);
    expect(find.text('Email'), findsOneWidget);

    // Verify order: Name should appear before Email (left to right).
    final nameCenter = tester.getCenter(find.text('Name'));
    final emailCenter = tester.getCenter(find.text('Email'));
    expect(nameCenter.dx, lessThan(emailCenter.dx));
  });

  testWidgets('BcLookupModal renders data rows (VAL-LOOKUP-004)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.text('Alice'), findsOneWidget);
    expect(find.text('Bob'), findsOneWidget);
  });

  // -----------------------------------------------------------------------
  // Search
  // -----------------------------------------------------------------------

  testWidgets('BcLookupModal search filters rows (VAL-LOOKUP-005)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        searchable: true,
        searchFields: const ['name', 'email'],
      ),
    ));
    await tester.pumpAndSettle();

    await tester.enterText(
      find.byKey(const ValueKey<String>('searchInput')),
      'ali',
    );
    await tester.pumpAndSettle();

    expect(find.text('Alice'), findsOneWidget);
    expect(find.text('Bob'), findsNothing);
    expect(find.text('Charlie'), findsNothing);
  });

  // -----------------------------------------------------------------------
  // Pagination
  // -----------------------------------------------------------------------

  testWidgets('BcLookupModal prev button decrements page (VAL-LOOKUP-006)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        pageSize: 1,
      ),
    ));
    await tester.pumpAndSettle();

    // Start on page 1 of 3
    expect(find.text('1 / 3'), findsOneWidget);

    // Advance to page 2
    await tester.tap(find.byKey(const ValueKey<String>('nextButton')));
    await tester.pumpAndSettle();
    expect(find.text('2 / 3'), findsOneWidget);

    // Go back to page 1
    await tester.tap(find.byKey(const ValueKey<String>('prevButton')));
    await tester.pumpAndSettle();
    expect(find.text('1 / 3'), findsOneWidget);

    // Prev button should be disabled on page 1 (onPressed null)
    final prevButton = tester.widget<IconButton>(
      find.byKey(const ValueKey<String>('prevButton')),
    );
    expect(prevButton.onPressed, isNull);
  });

  testWidgets('BcLookupModal next button increments page (VAL-LOOKUP-007)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        pageSize: 1,
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.text('1 / 3'), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey<String>('nextButton')));
    await tester.pumpAndSettle();
    expect(find.text('2 / 3'), findsOneWidget);

    // Advance to last page
    await tester.tap(find.byKey(const ValueKey<String>('nextButton')));
    await tester.pumpAndSettle();
    expect(find.text('3 / 3'), findsOneWidget);

    // Next button disabled on last page
    final nextButton = tester.widget<IconButton>(
      find.byKey(const ValueKey<String>('nextButton')),
    );
    expect(nextButton.onPressed, isNull);
  });

  testWidgets('BcLookupModal displays page info (VAL-LOOKUP-008)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        pageSize: 1,
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.text('1 / 3'), findsOneWidget);
  });

  // -----------------------------------------------------------------------
  // Single Select
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal single select emits value on confirm (VAL-LOOKUP-009)',
      (tester) async {
    dynamic capturedValue;
    BcChangeReason? capturedReason;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        onValueChange: (value, meta) {
          capturedValue = value;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Alice'));
    await tester.pumpAndSettle();

    await tester.tap(find.byKey(const ValueKey<String>('confirmButton')));
    await tester.pumpAndSettle();

    expect(capturedValue, '1');
    expect(capturedReason, BcChangeReason.confirm);
  });

  testWidgets(
      'BcLookupModal single select value is String (VAL-LOOKUP-010)',
      (tester) async {
    dynamic capturedValue;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        onValueChange: (value, meta) {
          capturedValue = value;
        },
      ),
    ));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Alice'));
    await tester.pumpAndSettle();

    await tester.tap(find.byKey(const ValueKey<String>('confirmButton')));
    await tester.pumpAndSettle();

    expect(capturedValue, isA<String>());
  });

  // -----------------------------------------------------------------------
  // Multiple Select
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal renders checkbox per row in multiple mode (VAL-LOOKUP-011)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        multiple: true,
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(Checkbox), findsNWidgets(3));
  });

  testWidgets(
      'BcLookupModal multiple select emits List<String> on confirm (VAL-LOOKUP-012)',
      (tester) async {
    dynamic capturedValue;
    BcChangeReason? capturedReason;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        multiple: true,
        onValueChange: (value, meta) {
          capturedValue = value;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    // Check first checkbox
    await tester.tap(find.byType(Checkbox).first);
    await tester.pumpAndSettle();

    // Check second checkbox
    await tester.tap(find.byType(Checkbox).at(1));
    await tester.pumpAndSettle();

    await tester.tap(find.byKey(const ValueKey<String>('confirmButton')));
    await tester.pumpAndSettle();

    expect(capturedValue, isA<List<String>>());
    expect(capturedValue as List<String>, containsAll(<String>['1', '2']));
    expect(capturedReason, BcChangeReason.confirm);
  });

  // -----------------------------------------------------------------------
  // Closeable
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal hides close button when closeable false (VAL-LOOKUP-013)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        closeable: false,
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byKey(const ValueKey<String>('closeButton')), findsNothing);
  });

  testWidgets(
      'BcLookupModal shows close button by default (VAL-LOOKUP-014)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byKey(const ValueKey<String>('closeButton')), findsOneWidget);
  });

  // -----------------------------------------------------------------------
  // closeOnOverlay
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal overlay tap closes with overlayClick reason (VAL-LOOKUP-015)',
      (tester) async {
    bool? capturedOpen;
    BcChangeReason? capturedReason;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        closeOnOverlay: true,
        onOpenChange: (open, meta) {
          capturedOpen = open;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    // Tap on the overlay (barrier) - it's the Container with black54
    final overlay = find.byWidgetPredicate(
      (widget) => widget is Container && widget.color == Colors.black54,
    );
    expect(overlay, findsOneWidget);

    await tester.tapAt(tester.getTopLeft(overlay));
    await tester.pumpAndSettle();

    expect(capturedOpen, false);
    expect(capturedReason, BcChangeReason.overlayClick);
  });

  testWidgets(
      'BcLookupModal overlay tap does not close when closeOnOverlay false (VAL-LOOKUP-016)',
      (tester) async {
    bool openChanged = false;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        closeOnOverlay: false,
        onOpenChange: (open, meta) {
          openChanged = true;
        },
      ),
    ));
    await tester.pumpAndSettle();

    final overlay = find.byWidgetPredicate(
      (widget) => widget is Container && widget.color == Colors.black54,
    );
    expect(overlay, findsOneWidget);

    await tester.tapAt(tester.getTopLeft(overlay));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);
    expect(openChanged, isFalse);
  });

  // -----------------------------------------------------------------------
  // Cancel button
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal cancel closes with cancel reason, no value (VAL-LOOKUP-017)',
      (tester) async {
    bool? capturedOpen;
    BcChangeReason? capturedReason;
    int valueCallCount = 0;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        onValueChange: (value, meta) {
          valueCallCount++;
        },
        onOpenChange: (open, meta) {
          capturedOpen = open;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    await tester.tap(find.byKey(const ValueKey<String>('cancelButton')));
    await tester.pumpAndSettle();

    expect(capturedOpen, false);
    expect(capturedReason, BcChangeReason.cancel);
    expect(valueCallCount, 0);
  });

  // -----------------------------------------------------------------------
  // Confirm button
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal confirm emits value and closes (VAL-LOOKUP-018)',
      (tester) async {
    dynamic capturedValue;
    BcChangeReason? capturedValueReason;
    bool? capturedOpen;
    BcChangeReason? capturedOpenReason;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        onValueChange: (value, meta) {
          capturedValue = value;
          capturedValueReason = meta.reason;
        },
        onOpenChange: (open, meta) {
          capturedOpen = open;
          capturedOpenReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Alice'));
    await tester.pumpAndSettle();

    await tester.tap(find.byKey(const ValueKey<String>('confirmButton')));
    await tester.pumpAndSettle();

    expect(capturedValue, '1');
    expect(capturedValueReason, BcChangeReason.confirm);
    expect(capturedOpen, false);
    expect(capturedOpenReason, BcChangeReason.confirm);
  });

  // -----------------------------------------------------------------------
  // onValueChange reasons
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal row tap fires select reason (VAL-LOOKUP-019)',
      (tester) async {
    dynamic capturedValue;
    BcChangeReason? capturedReason;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        onValueChange: (value, meta) {
          capturedValue = value;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Alice'));
    await tester.pumpAndSettle();

    expect(capturedValue, '1');
    expect(capturedReason, BcChangeReason.select);
  });

  testWidgets(
      'BcLookupModal deselect fires clear reason (VAL-LOOKUP-020)',
      (tester) async {
    dynamic capturedValue;
    BcChangeReason? capturedReason;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        multiple: true,
        onValueChange: (value, meta) {
          capturedValue = value;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    // Check then uncheck the only row to reach empty selection
    await tester.tap(find.byType(Checkbox).first);
    await tester.pumpAndSettle();

    await tester.tap(find.byType(Checkbox).first);
    await tester.pumpAndSettle();

    expect(capturedReason, BcChangeReason.clear);
    expect(capturedValue, <String>[]);
  });

  // -----------------------------------------------------------------------
  // onOpenChange reasons
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal external open=true fires toggle reason (VAL-LOOKUP-021)',
      (tester) async {
    bool openProp = false;
    bool? capturedOpen;
    BcChangeReason? capturedReason;

    late StateSetter setStateRef;

    await tester.pumpWidget(harness(
      StatefulBuilder(
        builder: (context, setState) {
          setStateRef = setState;
          return BcLookupModal(
            open: openProp,
            title: 'Test',
            columns: sampleColumns(),
            data: sampleData(),
            onOpenChange: (open, meta) {
              capturedOpen = open;
              capturedReason = meta.reason;
            },
          );
        },
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsNothing);

    setStateRef(() {
      openProp = true;
    });
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);
    expect(capturedOpen, true);
    expect(capturedReason, BcChangeReason.toggle);
  });

  testWidgets(
      'BcLookupModal escape key closes with escape reason (VAL-LOOKUP-022)',
      (tester) async {
    bool? capturedOpen;
    BcChangeReason? capturedReason;

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        onOpenChange: (open, meta) {
          capturedOpen = open;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    await tester.sendKeyEvent(LogicalKeyboardKey.escape);
    await tester.pumpAndSettle();

    expect(capturedOpen, false);
    expect(capturedReason, BcChangeReason.escape);
  });

  // -----------------------------------------------------------------------
  // Empty state
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal shows empty state when data empty (VAL-LOOKUP-023)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: const <Map<String, dynamic>>[],
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byKey(const ValueKey<String>('empty')), findsOneWidget);
    expect(find.text('No records found'), findsOneWidget);
  });

  // -----------------------------------------------------------------------
  // Uncontrolled mode
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal defaultOpen renders dialog open (VAL-LOOKUP-024)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        defaultOpen: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);
  });

  testWidgets(
      'BcLookupModal uncontrolled close button hides dialog (VAL-LOOKUP-025)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        defaultOpen: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey<String>('closeButton')));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsNothing);
  });

  // -----------------------------------------------------------------------
  // Controlled mode
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal controlled open=false hides dialog (VAL-LOOKUP-026)',
      (tester) async {
    bool openProp = true;

    late StateSetter setStateRef;

    await tester.pumpWidget(harness(
      StatefulBuilder(
        builder: (context, setState) {
          setStateRef = setState;
          return BcLookupModal(
            open: openProp,
            title: 'Test',
            columns: sampleColumns(),
            data: sampleData(),
          );
        },
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);

    setStateRef(() {
      openProp = false;
    });
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsNothing);
  });

  // -----------------------------------------------------------------------
  // Imperative methods
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal reset fires programmatic reason (VAL-LOOKUP-027)',
      (tester) async {
    dynamic capturedValue;
    BcChangeReason? capturedReason;

    final key = GlobalKey<BcLookupModalState>();

    await tester.pumpWidget(harness(
      BcLookupModal(
        key: key,
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        defaultValue: '1',
        onValueChange: (value, meta) {
          capturedValue = value;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    key.currentState!.reset();
    await tester.pumpAndSettle();

    expect(capturedReason, BcChangeReason.programmatic);
    expect(capturedValue, '');
  });

  testWidgets(
      'BcLookupModal imperative open fires programmatic reason (VAL-LOOKUP-028)',
      (tester) async {
    bool? capturedOpen;
    BcChangeReason? capturedReason;

    final key = GlobalKey<BcLookupModalState>();

    await tester.pumpWidget(harness(
      BcLookupModal(
        key: key,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        onOpenChange: (open, meta) {
          capturedOpen = open;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsNothing);

    key.currentState!.open();
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);
    expect(capturedOpen, true);
    expect(capturedReason, BcChangeReason.programmatic);
  });

  testWidgets(
      'BcLookupModal imperative close fires programmatic reason (VAL-LOOKUP-029)',
      (tester) async {
    bool? capturedOpen;
    BcChangeReason? capturedReason;

    final key = GlobalKey<BcLookupModalState>();

    await tester.pumpWidget(harness(
      BcLookupModal(
        key: key,
        defaultOpen: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        onOpenChange: (open, meta) {
          capturedOpen = open;
          capturedReason = meta.reason;
        },
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);

    key.currentState!.close();
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsNothing);
    expect(capturedOpen, false);
    expect(capturedReason, BcChangeReason.programmatic);
  });

  // -----------------------------------------------------------------------
  // Search hidden
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal hides search input when searchable false (VAL-LOOKUP-030)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        searchable: false,
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byKey(const ValueKey<String>('searchInput')), findsNothing);
  });

  // -----------------------------------------------------------------------
  // Pre-selected values
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal pre-selected value highlights row (VAL-LOOKUP-031)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        value: '2',
      ),
    ));
    await tester.pumpAndSettle();

    // The selected row (id=2, Bob) should have a 'row-selected-2' container.
    expect(
      find.byKey(const ValueKey<String>('row-selected-2')),
      findsOneWidget,
    );
  });

  testWidgets(
      'BcLookupModal pre-selected multiple values check boxes (VAL-LOOKUP-032)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
        multiple: true,
        value: const <String>['1', '3'],
      ),
    ));
    await tester.pumpAndSettle();

    // Count checked checkboxes
    final checkedCount = tester
        .widgetList<Checkbox>(find.byType(Checkbox))
        .where((cb) => cb.value == true)
        .length;

    expect(checkedCount, 2);
  });

  // -----------------------------------------------------------------------
  // Default pageSize
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal default pageSize paginates large data (VAL-LOOKUP-033)',
      (tester) async {
    final largeData = List<Map<String, dynamic>>.generate(
      15,
      (i) => <String, dynamic>{
        'id': '${i + 1}',
        'name': 'User${i + 1}',
        'email': 'user${i + 1}@example.com',
      },
    );

    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: largeData,
      ),
    ));
    await tester.pumpAndSettle();

    // Default pageSize=10, so 15 items → 2 pages → pagination visible
    expect(find.byKey(const ValueKey<String>('pageInfo')), findsOneWidget);
    expect(find.text('1 / 2'), findsOneWidget);
  });

  // -----------------------------------------------------------------------
  // Title omitted
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal omitted title renders no title text (VAL-LOOKUP-034)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.byType(AlertDialog), findsOneWidget);
    // The close button should still render
    expect(find.byKey(const ValueKey<String>('closeButton')), findsOneWidget);
  });

  // -----------------------------------------------------------------------
  // Confirm disabled
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal confirm disabled without selection (VAL-LOOKUP-035)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Test',
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    final confirmButton = tester.widget<ElevatedButton>(
      find.byKey(const ValueKey<String>('confirmButton')),
    );
    expect(confirmButton.enabled, isFalse);
  });

  // -----------------------------------------------------------------------
  // VAL-CROSS-005: BcDataTable reuse inside dialog
  // -----------------------------------------------------------------------

  testWidgets(
      'BcLookupModal renders BcDataTable inside dialog when open (VAL-CROSS-005)',
      (tester) async {
    await tester.pumpWidget(harness(
      BcLookupModal(
        open: true,
        title: 'Pick Contact',
        columns: sampleColumns(),
        data: sampleData(),
      ),
    ));
    await tester.pumpAndSettle();

    // Dialog surface must be present
    expect(find.byType(AlertDialog), findsOneWidget);
    // Title must match
    expect(find.text('Pick Contact'), findsOneWidget);
    // Exactly one BcDataTable rendered inside the dialog
    expect(
      find.byWidgetPredicate((widget) => widget is BcDataTable),
      findsOneWidget,
    );
    // Data must be reflected
    expect(find.text('Alice'), findsOneWidget);
    expect(find.text('Bob'), findsOneWidget);
  });
}
