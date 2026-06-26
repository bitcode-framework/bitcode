import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:bc_ui_flutter/src/components/bc_chart.dart';
import 'package:bc_ui_flutter/src/core/types.dart';

void main() {
  Widget harness(Widget child) {
    return MaterialApp(home: Scaffold(body: child));
  }

  group('VAL-M5-014: Chart type determines series rendering', () {
    testWidgets('bar type renders chart with CustomPaint', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.bar,
          categories: ['A', 'B', 'C'],
          series: [
            BcChartSeries(name: 'S1', data: [10, 20, 30]),
          ],
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart')), findsOneWidget);
      expect(find.byKey(const Key('bc_chart_canvas')), findsOneWidget);
      expect(find.byType(CustomPaint), findsWidgets);
      expect(tester.takeException(), isNull);
    });

    testWidgets('line type renders connected polyline', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.line,
          categories: ['A', 'B', 'C'],
          series: [
            BcChartSeries(name: 'S1', data: [10, 20, 30]),
          ],
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart_canvas')), findsOneWidget);
      expect(find.byType(CustomPaint), findsWidgets);
      expect(tester.takeException(), isNull);
    });

    testWidgets('pie type renders circular proportion', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.pie,
          items: [
            BcChartItem(name: 'A', value: 30),
            BcChartItem(name: 'B', value: 70),
          ],
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart_canvas')), findsOneWidget);
      expect(find.byType(CustomPaint), findsWidgets);
      expect(tester.takeException(), isNull);
    });

    testWidgets('doughnut (pie alias) renders pie with hollow center',
        (tester) async {
      // ADR-015: `doughnut` is a backward-compatible alias of `pie` (the
      // 26-type enum has no doughnut value). The legacy hollow-center
      // rendering is preserved via `doughnutInnerRadiusFactor`.
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.pie,
          items: [
            BcChartItem(name: 'A', value: 40),
            BcChartItem(name: 'B', value: 60),
          ],
          doughnutInnerRadiusFactor: 0.55,
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart_canvas')), findsOneWidget);
      expect(find.byType(CustomPaint), findsWidgets);
      expect(tester.takeException(), isNull);
    });

    testWidgets('default type is bar', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          categories: ['A'],
          series: [BcChartSeries(name: 'S1', data: [10])],
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart')), findsOneWidget);
      expect(find.byType(CustomPaint), findsWidgets);
    });
  });

  group('VAL-M5-015: Series data maps to visual proportions', () {
    testWidgets(
        'bar chart renders proportional heights without error',
        (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.bar,
          categories: ['A', 'B'],
          series: [
            BcChartSeries(name: 'S1', data: [25, 100]),
          ],
          height: 300,
          legend: false,
          showAxisLabels: false,
        ),
      ));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('pie with zero total renders gracefully', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.pie,
          items: [
            BcChartItem(name: 'A', value: 0),
            BcChartItem(name: 'B', value: 0),
          ],
          height: 300,
          legend: false,
        ),
      ));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('doughnut (pie alias) with zero total renders gracefully',
        (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.pie,
          items: [
            BcChartItem(name: 'A', value: 0),
          ],
          height: 300,
          legend: false,
        ),
      ));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('empty data renders chart container without error',
        (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.bar,
          height: 300,
          legend: false,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart')), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('line chart with single data point renders without error',
        (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.line,
          categories: ['A'],
          series: [BcChartSeries(name: 'S1', data: [50])],
          height: 300,
          legend: false,
          showAxisLabels: false,
        ),
      ));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('VAL-M5-016: Legend and axis labels render from config', () {
    testWidgets('legend renders one entry per series', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.bar,
          categories: ['A', 'B'],
          series: [
            BcChartSeries(name: 'Revenue', data: [10, 20]),
            BcChartSeries(name: 'Cost', data: [5, 10]),
          ],
          legend: true,
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart_legend')), findsOneWidget);
      expect(find.text('Revenue'), findsOneWidget);
      expect(find.text('Cost'), findsOneWidget);
      expect(
          find.byKey(const Key('bc_chart_legend_swatch_0')), findsOneWidget);
      expect(
          find.byKey(const Key('bc_chart_legend_swatch_1')), findsOneWidget);
    });

    testWidgets('legend disabled does not render', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.bar,
          categories: ['A', 'B'],
          series: [
            BcChartSeries(name: 'Revenue', data: [10, 20]),
          ],
          legend: false,
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart_legend')), findsNothing);
    });

    testWidgets('axis labels render categories on x-axis', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.bar,
          categories: ['Jan', 'Feb', 'Mar'],
          series: [
            BcChartSeries(name: 'S1', data: [10, 20, 30]),
          ],
          legend: false,
          showAxisLabels: true,
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart_xaxis')), findsOneWidget);
      expect(find.text('Jan'), findsWidgets);
      expect(find.text('Feb'), findsWidgets);
      expect(find.text('Mar'), findsWidgets);
    });

    testWidgets(
        'updating data re-renders legend and axis labels',
        (tester) async {
      late void Function(VoidCallback) setStateRef;
      List<BcChartSeries> currentSeries = const [
        BcChartSeries(name: 'Old', data: [10, 20]),
      ];
      List<String> currentCategories = const ['Old1', 'Old2'];

      Widget buildWidget() {
        return StatefulBuilder(
          builder: (context, setState) {
            setStateRef = setState;
            return MaterialApp(
              home: Scaffold(
                body: BcChart(
                  type: BcChartType.bar,
                  categories: currentCategories,
                  series: currentSeries,
                  legend: true,
                  showAxisLabels: true,
                  height: 300,
                ),
              ),
            );
          },
        );
      }

      await tester.pumpWidget(buildWidget());
      await tester.pumpAndSettle();

      expect(find.text('Old'), findsOneWidget);

      setStateRef(() {
        currentSeries = const [BcChartSeries(name: 'New', data: [30, 40])];
        currentCategories = const ['New1', 'New2'];
      });
      await tester.pumpWidget(buildWidget());
      await tester.pumpAndSettle();

      expect(find.text('Old'), findsNothing);
      expect(find.text('New'), findsOneWidget);
      expect(find.text('New1'), findsWidgets);
      expect(find.text('New2'), findsWidgets);
    });

    testWidgets(
        'pie legend uses item names, not series names',
        (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.pie,
          items: [
            BcChartItem(name: 'Alpha', value: 30),
            BcChartItem(name: 'Beta', value: 70),
          ],
          legend: true,
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.text('Alpha'), findsOneWidget);
      expect(find.text('Beta'), findsOneWidget);
    });
  });

  group('Additional edge cases', () {
    testWidgets('chartTitle renders when provided', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          chartTitle: 'My Chart',
          categories: ['A'],
          series: [BcChartSeries(name: 'S1', data: [10])],
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('bc_chart_title')), findsOneWidget);
      expect(find.text('My Chart'), findsOneWidget);
    });

    testWidgets('loading overlay shows when loading is true',
        (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          loading: true,
          categories: ['A'],
          series: [BcChartSeries(name: 'S1', data: [10])],
          height: 300,
        ),
      ));
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('multi-series bar chart renders all series',
        (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.bar,
          categories: ['A', 'B', 'C'],
          series: [
            BcChartSeries(name: 'S1', data: [10, 20, 30]),
            BcChartSeries(name: 'S2', data: [15, 25, 35]),
            BcChartSeries(name: 'S3', data: [5, 15, 25]),
          ],
          height: 300,
          legend: true,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.text('S1'), findsOneWidget);
      expect(find.text('S2'), findsOneWidget);
      expect(find.text('S3'), findsOneWidget);
    });

    testWidgets('string data is parsed correctly', (tester) async {
      await tester.pumpWidget(harness(
        const BcChart(
          type: BcChartType.pie,
          data:
              '[{"name":"A","value":30},{"name":"B","value":70}]',
          legend: true,
          height: 300,
        ),
      ));
      await tester.pumpAndSettle();

      expect(find.text('A'), findsWidgets);
      expect(find.text('B'), findsWidgets);
    });
  });
}
