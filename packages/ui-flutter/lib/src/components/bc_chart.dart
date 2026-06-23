library bc_chart;

import 'dart:convert';
import 'dart:math' as math;

import 'package:flutter/material.dart';
import '../core/types.dart';
import '../theme/bc_theme.dart';

/// Full chart widget supporting bar, line, pie, and doughnut types
/// (ADR-011 Phase 5c — M5 data display).
///
/// Uses [CustomPaint] to render each chart kind natively in Flutter.
/// For bar/line, multi-series data with category axis is supported.
/// For pie/doughnut, single-series data is rendered as proportional slices.
///
/// Controlled: all data is provided externally via props; the widget is
/// stateless with respect to chart data (re-renders on prop change).
class BcChart extends StatelessWidget {
  final BcChartType type;
  final List<BcChartSeries> series;
  final List<String> categories;
  final List<BcChartItem> items;
  final List<Color> colors;
  final String? chartTitle;
  final bool legend;
  final bool tooltipEnabled;
  final bool animate;
  final double height;
  final double? width;
  final bool loading;
  final bool stacked;
  final bool horizontal;
  final String? ariaLabel;
  final bool showAxisLabels;
  final double doughnutInnerRadiusFactor;
  final String? data;

  static const List<Color> _defaultPalette = [
    BcTheme.primary,
    BcTheme.info,
    BcTheme.success,
    BcTheme.warning,
    BcTheme.destructive,
    BcTheme.mutedForeground,
  ];

  const BcChart({
    super.key,
    this.type = BcChartType.bar,
    this.series = const [],
    this.categories = const [],
    this.items = const [],
    this.colors = _defaultPalette,
    this.chartTitle,
    this.legend = true,
    this.tooltipEnabled = true,
    this.animate = true,
    this.height = 320,
    this.width,
    this.loading = false,
    this.stacked = false,
    this.horizontal = false,
    this.ariaLabel,
    this.showAxisLabels = true,
    this.doughnutInnerRadiusFactor = 0.55,
    this.data,
  });

  @override
  Widget build(BuildContext context) {
    final parsedItems = _parseData();
    final effectiveItems = items.isNotEmpty ? items : parsedItems;
    final effectiveSeries = series;

    return Semantics(
      container: true,
      label: ariaLabel ?? chartTitle ?? 'Chart',
      child: Container(
        key: const Key('bc_chart'),
        width: width,
        height: height,
        decoration: BoxDecoration(
          color: BcTheme.background,
          border: Border.all(color: BcTheme.border),
          borderRadius: BorderRadius.circular(BcTheme.radius),
        ),
        child: loading
            ? _buildLoading(context)
            : _buildChartContent(context, effectiveItems, effectiveSeries),
      ),
    );
  }

  List<BcChartItem> _parseData() {
    if (data == null || data!.isEmpty) return const [];
    try {
      final decoded = jsonDecode(data!);
      if (decoded is List) {
        return decoded
            .map((e) => BcChartItem(
                  name: e['name'] as String? ?? '',
                  value: (e['value'] as num?)?.toDouble() ?? 0,
                ))
            .toList();
      }
    } catch (_) {}
    return const [];
  }

  Widget _buildLoading(BuildContext context) {
    return const Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          SizedBox(
            width: 24,
            height: 24,
            child: CircularProgressIndicator(
              strokeWidth: 2,
              color: BcTheme.primary,
            ),
          ),
          SizedBox(height: 8),
          Text(
            'Loading...',
            style: TextStyle(fontSize: 12, color: BcTheme.mutedForeground),
          ),
        ],
      ),
    );
  }

  Widget _buildChartContent(
      BuildContext context, List<BcChartItem> effectiveItems, List<BcChartSeries> effectiveSeries) {
    return Padding(
      padding: const EdgeInsets.all(12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (chartTitle != null)
            Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: Text(
                chartTitle!,
                key: const Key('bc_chart_title'),
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: BcTheme.foreground,
                ),
              ),
            ),
          Expanded(child: _buildChartArea(effectiveItems, effectiveSeries)),
          if (legend) _buildLegend(effectiveItems, effectiveSeries),
        ],
      ),
    );
  }

  Widget _buildChartArea(List<BcChartItem> effectiveItems, List<BcChartSeries> effectiveSeries) {
    final hasData = type == BcChartType.pie || type == BcChartType.doughnut
        ? effectiveItems.isNotEmpty
        : effectiveSeries.isNotEmpty && effectiveSeries.any((s) => s.data.isNotEmpty);

    if (!hasData) {
      return const Center(
        child: Text(
          'No data',
          key: Key('bc_chart_no_data'),
          style: TextStyle(fontSize: 13, color: BcTheme.mutedForeground),
        ),
      );
    }

    switch (type) {
      case BcChartType.bar:
        return _buildBarChart(effectiveSeries);
      case BcChartType.line:
        return _buildLineChart(effectiveSeries);
      case BcChartType.pie:
        return _buildPieChart(effectiveItems);
      case BcChartType.doughnut:
        return _buildDoughnutChart(effectiveItems);
    }
  }

  Widget _buildBarChart(List<BcChartSeries> effectiveSeries) {
    return LayoutBuilder(
      builder: (context, constraints) {
        return _BarChartView(
          series: effectiveSeries,
          categories: categories,
          colors: colors,
          horizontal: horizontal,
          stacked: stacked,
          showAxisLabels: showAxisLabels,
          constraints: constraints,
        );
      },
    );
  }

  Widget _buildLineChart(List<BcChartSeries> effectiveSeries) {
    return LayoutBuilder(
      builder: (context, constraints) {
        return _LineChartView(
          series: effectiveSeries,
          categories: categories,
          colors: colors,
          showAxisLabels: showAxisLabels,
          constraints: constraints,
        );
      },
    );
  }

  Widget _buildPieChart(List<BcChartItem> effectiveItems) {
    return _PieChartView(
      items: effectiveItems,
      colors: colors,
      innerRadiusFactor: 0.0,
    );
  }

  Widget _buildDoughnutChart(List<BcChartItem> effectiveItems) {
    return _PieChartView(
      items: effectiveItems,
      colors: colors,
      innerRadiusFactor: doughnutInnerRadiusFactor,
    );
  }

  Widget _buildLegend(List<BcChartItem> effectiveItems, List<BcChartSeries> effectiveSeries) {
    final legendEntries = <({String label, Color color})>[];

    switch (type) {
      case BcChartType.bar:
      case BcChartType.line:
        for (var i = 0; i < effectiveSeries.length; i++) {
          legendEntries.add((
            label: effectiveSeries[i].name,
            color: colors[i % colors.length],
          ));
        }
        break;
      case BcChartType.pie:
      case BcChartType.doughnut:
        for (var i = 0; i < effectiveItems.length; i++) {
          legendEntries.add((
            label: effectiveItems[i].name,
            color: colors[i % colors.length],
          ));
        }
        break;
    }

    if (legendEntries.isEmpty) return const SizedBox.shrink();

    return Padding(
      key: const Key('bc_chart_legend'),
      padding: const EdgeInsets.only(top: 8),
      child: Wrap(
        spacing: 12,
        runSpacing: 4,
        children: legendEntries.asMap().entries.map((entry) {
          final i = entry.key;
          final e = entry.value;
          return Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                key: Key('bc_chart_legend_swatch_$i'),
                width: 10,
                height: 10,
                decoration: BoxDecoration(
                  color: e.color,
                  shape: BoxShape.circle,
                ),
              ),
              const SizedBox(width: 4),
              Text(
                e.label,
                style: const TextStyle(
                  fontSize: 11,
                  color: BcTheme.mutedForeground,
                ),
              ),
            ],
          );
        }).toList(),
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// Bar chart painter
// ---------------------------------------------------------------------------

class _BarChartView extends StatelessWidget {
  final List<BcChartSeries> series;
  final List<String> categories;
  final List<Color> colors;
  final bool horizontal;
  final bool stacked;
  final bool showAxisLabels;
  final BoxConstraints constraints;

  const _BarChartView({
    required this.series,
    required this.categories,
    required this.colors,
    required this.horizontal,
    required this.stacked,
    required this.showAxisLabels,
    required this.constraints,
  });

  @override
  Widget build(BuildContext context) {
    final hasAxisLabels = showAxisLabels && categories.isNotEmpty;
    final axisSpace = hasAxisLabels ? 36.0 : 8.0;

    return Column(
      children: [
        Expanded(
          child: LayoutBuilder(
            builder: (ctx, inner) {
              return CustomPaint(
                key: const Key('bc_chart_canvas'),
                size: Size(inner.maxWidth, inner.maxHeight),
                painter: _BarChartPainter(
                  series: series,
                  categories: categories,
                  colors: colors,
                  horizontal: horizontal,
                  stacked: stacked,
                  chartWidth: inner.maxWidth,
                  chartHeight: inner.maxHeight,
                  axisSpace: axisSpace,
                ),
              );
            },
          ),
        ),
        if (hasAxisLabels) _buildAxisLabels(),
      ],
    );
  }

  Widget _buildAxisLabels() {
    if (categories.isEmpty) return const SizedBox.shrink();

    return SizedBox(
      key: const Key('bc_chart_xaxis'),
      height: 28,
      child: Row(
        children: [
          const SizedBox(
            width: 28,
            child: Text(
              '',
              style: TextStyle(fontSize: 9, color: BcTheme.mutedForeground),
            ),
          ),
          Expanded(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: categories.map((cat) {
                return Flexible(
                  child: Text(
                    cat,
                    key: Key('bc_chart_axis_$cat'),
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                      fontSize: 9,
                      color: BcTheme.mutedForeground,
                    ),
                  ),
                );
              }).toList(),
            ),
          ),
        ],
      ),
    );
  }
}

class _BarChartPainter extends CustomPainter {
  final List<BcChartSeries> series;
  final List<String> categories;
  final List<Color> colors;
  final bool horizontal;
  final bool stacked;
  final double chartWidth;
  final double chartHeight;
  final double axisSpace;

  _BarChartPainter({
    required this.series,
    required this.categories,
    required this.colors,
    required this.horizontal,
    required this.stacked,
    required this.chartWidth,
    required this.chartHeight,
    required this.axisSpace,
  });

  @override
  void paint(Canvas canvas, Size size) {
    if (series.isEmpty) return;

    final w = size.width;
    final h = size.height - axisSpace;
    final plotW = w - 32;
    final plotH = h;

    double maxVal = 0;
    for (final s in series) {
      if (stacked) {
        for (var i = 0; i < s.data.length; i++) {
          maxVal = math.max(maxVal, s.data[i]);
        }
      } else {
        for (final v in s.data) {
          maxVal = math.max(maxVal, v);
        }
      }
    }

    if (stacked) {
      final catCount = categories.isNotEmpty
          ? categories.length
          : (series.isNotEmpty ? series.first.data.length : 0);
      for (var c = 0; c < catCount; c++) {
        double stackSum = 0;
        for (final s in series) {
          if (c < s.data.length) stackSum += s.data[c];
        }
        maxVal = math.max(maxVal, stackSum);
      }
    }

    if (maxVal <= 0) maxVal = 1;

    final catCount = categories.isNotEmpty
        ? categories.length
        : (series.isNotEmpty ? series.first.data.length : 0);
    if (catCount == 0) return;

    final groupSpace = plotW / catCount;
    final groupPadding = groupSpace * 0.15;
    final groupWidth = groupSpace - groupPadding * 2;

    _drawGridLines(canvas, plotW, plotH, maxVal, 32.0);

    final barCount = stacked ? 1 : series.length;
    final barWidth = groupWidth / barCount;

    for (var catIdx = 0; catIdx < catCount; catIdx++) {
      final groupStart = 32.0 + catIdx * groupSpace + groupPadding;

      if (stacked) {
        double stackY = plotH;
        for (var sIdx = 0; sIdx < series.length; sIdx++) {
          final value = catIdx < series[sIdx].data.length
              ? series[sIdx].data[catIdx]
              : 0.0;
          final barH = (value / maxVal) * plotH;
          final rect = Rect.fromLTWH(
            groupStart,
            stackY - barH,
            groupWidth,
            barH,
          );
          final paint = Paint()
            ..color = colors[sIdx % colors.length]
            ..style = PaintingStyle.fill;
          canvas.drawRRect(
            RRect.fromRectAndRadius(rect, const Radius.circular(2)),
            paint,
          );
          stackY -= barH;
        }
      } else {
        for (var sIdx = 0; sIdx < series.length; sIdx++) {
          final value = catIdx < series[sIdx].data.length
              ? series[sIdx].data[catIdx]
              : 0.0;
          final barH = (value / maxVal) * plotH;
          final barX = groupStart + sIdx * barWidth;
          final rect = Rect.fromLTWH(barX, plotH - barH, barWidth, barH);
          final paint = Paint()
            ..color = colors[sIdx % colors.length]
            ..style = PaintingStyle.fill;
          canvas.drawRRect(
            RRect.fromRectAndRadius(rect, const Radius.circular(2)),
            paint,
          );
        }
      }
    }
  }

  void _drawGridLines(
    Canvas canvas,
    double plotW,
    double plotH,
    double maxVal,
    double leftPad,
  ) {
    final gridPaint = Paint()
      ..color = BcTheme.border.withValues(alpha: 0.3)
      ..strokeWidth = 0.5
      ..style = PaintingStyle.stroke;

    const ticks = 4;
    for (var i = 0; i <= ticks; i++) {
      final y = (plotH / ticks) * i;
      canvas.drawLine(
        Offset(leftPad, y),
        Offset(leftPad + plotW, y),
        gridPaint,
      );
    }

    final axisPaint = Paint()
      ..color = BcTheme.border
      ..strokeWidth = 1;
    canvas.drawLine(
      Offset(leftPad, 0),
      Offset(leftPad, plotH),
      axisPaint,
    );
    canvas.drawLine(
      Offset(leftPad, plotH),
      Offset(leftPad + plotW, plotH),
      axisPaint,
    );
  }

  @override
  bool shouldRepaint(covariant _BarChartPainter old) =>
      old.series != series ||
      old.categories != categories ||
      old.colors != colors ||
      old.stacked != stacked ||
      old.horizontal != horizontal ||
      old.chartWidth != chartWidth ||
      old.chartHeight != chartHeight;
}

// ---------------------------------------------------------------------------
// Line chart painter
// ---------------------------------------------------------------------------

class _LineChartView extends StatelessWidget {
  final List<BcChartSeries> series;
  final List<String> categories;
  final List<Color> colors;
  final bool showAxisLabels;
  final BoxConstraints constraints;

  const _LineChartView({
    required this.series,
    required this.categories,
    required this.colors,
    required this.showAxisLabels,
    required this.constraints,
  });

  @override
  Widget build(BuildContext context) {
    final hasAxisLabels = showAxisLabels && categories.isNotEmpty;

    return Column(
      children: [
        Expanded(
          child: LayoutBuilder(
            builder: (ctx, inner) {
              return CustomPaint(
                key: const Key('bc_chart_canvas'),
                size: Size(inner.maxWidth, inner.maxHeight),
                painter: _LineChartPainter(
                  series: series,
                  categories: categories,
                  colors: colors,
                  chartWidth: inner.maxWidth,
                  chartHeight: inner.maxHeight,
                  axisSpace: hasAxisLabels ? 36.0 : 8.0,
                ),
              );
            },
          ),
        ),
        if (hasAxisLabels) _buildAxisLabels(),
      ],
    );
  }

  Widget _buildAxisLabels() {
    if (categories.isEmpty) return const SizedBox.shrink();

    return SizedBox(
      height: 28,
      child: Row(
        children: [
          const SizedBox(width: 28),
          Expanded(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: categories.map((cat) {
                return Flexible(
                  child: Text(
                    cat,
                    key: Key('bc_chart_axis_$cat'),
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                      fontSize: 9,
                      color: BcTheme.mutedForeground,
                    ),
                  ),
                );
              }).toList(),
            ),
          ),
        ],
      ),
    );
  }
}

class _LineChartPainter extends CustomPainter {
  final List<BcChartSeries> series;
  final List<String> categories;
  final List<Color> colors;
  final double chartWidth;
  final double chartHeight;
  final double axisSpace;

  _LineChartPainter({
    required this.series,
    required this.categories,
    required this.colors,
    required this.chartWidth,
    required this.chartHeight,
    required this.axisSpace,
  });

  @override
  void paint(Canvas canvas, Size size) {
    if (series.isEmpty) return;

    final w = size.width;
    final h = size.height - axisSpace;
    const leftPad = 32.0;
    final plotW = w - leftPad - 8;
    final plotH = h;

    double maxVal = 0;
    for (final s in series) {
      for (final v in s.data) {
        if (v > maxVal) maxVal = v;
      }
    }
    if (maxVal <= 0) maxVal = 1;

    _drawGrid(canvas, plotW, plotH, maxVal, leftPad);

    final catCount = categories.isNotEmpty
        ? categories.length
        : (series.isNotEmpty ? series.first.data.length : 0);
    if (catCount == 0) return;

    final stepX = catCount > 1 ? plotW / (catCount - 1) : 0.0;

    for (var sIdx = 0; sIdx < series.length; sIdx++) {
      final s = series[sIdx];
      final color = colors[sIdx % colors.length];
      final points = <Offset>[];

      for (var i = 0; i < s.data.length && i < catCount; i++) {
        final x = leftPad + i * stepX;
        final y = plotH - (s.data[i] / maxVal) * plotH;
        points.add(Offset(x, y));
      }

      if (points.isEmpty) continue;

      final linePaint = Paint()
        ..color = color
        ..strokeWidth = 2.0
        ..style = PaintingStyle.stroke
        ..strokeCap = StrokeCap.round
        ..strokeJoin = StrokeJoin.round;

      final path = Path();
      path.moveTo(points[0].dx, points[0].dy);
      for (var i = 1; i < points.length; i++) {
        path.lineTo(points[i].dx, points[i].dy);
      }
      canvas.drawPath(path, linePaint);

      final dotPaint = Paint()
        ..color = color
        ..style = PaintingStyle.fill;

      for (final p in points) {
        canvas.drawCircle(p, 3.0, dotPaint);
      }
    }
  }

  void _drawGrid(
    Canvas canvas,
    double plotW,
    double plotH,
    double maxVal,
    double leftPad,
  ) {
    final gridPaint = Paint()
      ..color = BcTheme.border.withValues(alpha: 0.3)
      ..strokeWidth = 0.5
      ..style = PaintingStyle.stroke;

    const ticks = 4;
    for (var i = 0; i <= ticks; i++) {
      final y = (plotH / ticks) * i;
      canvas.drawLine(
        Offset(leftPad, y),
        Offset(leftPad + plotW, y),
        gridPaint,
      );
    }

    final axisPaint = Paint()
      ..color = BcTheme.border
      ..strokeWidth = 1;
    canvas.drawLine(Offset(leftPad, 0), Offset(leftPad, plotH), axisPaint);
    canvas.drawLine(
      Offset(leftPad, plotH),
      Offset(leftPad + plotW, plotH),
      axisPaint,
    );
  }

  @override
  bool shouldRepaint(covariant _LineChartPainter old) =>
      old.series != series ||
      old.categories != categories ||
      old.colors != colors ||
      old.chartWidth != chartWidth ||
      old.chartHeight != chartHeight;
}

// ---------------------------------------------------------------------------
// Pie / Doughnut chart painter
// ---------------------------------------------------------------------------

class _PieChartView extends StatelessWidget {
  final List<BcChartItem> items;
  final List<Color> colors;
  final double innerRadiusFactor;

  const _PieChartView({
    required this.items,
    required this.colors,
    required this.innerRadiusFactor,
  });

  @override
  Widget build(BuildContext context) {
    return CustomPaint(
      key: const Key('bc_chart_canvas'),
      painter: _PieChartPainter(
        items: items,
        colors: colors,
        innerRadiusFactor: innerRadiusFactor,
      ),
      child: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (items.isEmpty)
              const Text(
                'No data',
                style: TextStyle(
                  fontSize: 12,
                  color: BcTheme.mutedForeground,
                ),
              ),
          ],
        ),
      ),
    );
  }
}

class _PieChartPainter extends CustomPainter {
  final List<BcChartItem> items;
  final List<Color> colors;
  final double innerRadiusFactor;

  _PieChartPainter({
    required this.items,
    required this.colors,
    required this.innerRadiusFactor,
  });

  @override
  void paint(Canvas canvas, Size size) {
    if (items.isEmpty) return;

    final total = items.fold<double>(0, (sum, item) => sum + item.value);
    if (total <= 0) return;

    final cx = size.width / 2;
    final cy = size.height / 2;
    final radius = math.min(cx, cy) * 0.85;
    final innerRadius = radius * innerRadiusFactor.clamp(0.0, 0.95);

    double startAngle = -math.pi / 2;

    for (var i = 0; i < items.length; i++) {
      final sweepAngle = (items[i].value / total) * 2 * math.pi;
      final color = colors[i % colors.length];

      final paint = Paint()
        ..color = color
        ..style = PaintingStyle.fill
        ..strokeWidth = 0;

      final rect = Rect.fromCircle(
        center: Offset(cx, cy),
        radius: radius,
      );

      if (innerRadius > 0) {
        final arcPath = Path.combine(
          PathOperation.difference,
          _createArcPath(cx, cy, radius, innerRadius, startAngle, sweepAngle),
          Path(),
        );
        canvas.drawPath(arcPath, paint);

        final strokePaint = Paint()
          ..color = BcTheme.background
          ..strokeWidth = 2
          ..style = PaintingStyle.stroke;
        canvas.drawArc(
          rect,
          startAngle,
          sweepAngle,
          true,
          strokePaint,
        );
      } else {
        canvas.drawArc(rect, startAngle, sweepAngle, true, paint);

        final strokePaint = Paint()
          ..color = BcTheme.background
          ..strokeWidth = 2
          ..style = PaintingStyle.stroke;
        canvas.drawArc(rect, startAngle, sweepAngle, true, strokePaint);
      }

      startAngle += sweepAngle;
    }
  }

  Path _createArcPath(
    double cx,
    double cy,
    double outerR,
    double innerR,
    double startAngle,
    double sweepAngle,
  ) {
    final path = Path();
    final outerRect = Rect.fromCircle(
      center: Offset(cx, cy),
      radius: outerR,
    );
    final innerRect = Rect.fromCircle(
      center: Offset(cx, cy),
      radius: innerR,
    );

    path.moveTo(
      cx + outerR * math.cos(startAngle),
      cy + outerR * math.sin(startAngle),
    );
    path.arcTo(outerRect, startAngle, sweepAngle, false);
    path.lineTo(
      cx + innerR * math.cos(startAngle + sweepAngle),
      cy + innerR * math.sin(startAngle + sweepAngle),
    );
    path.arcTo(innerRect, startAngle + sweepAngle, -sweepAngle, false);
    path.close();
    return path;
  }

  @override
  bool shouldRepaint(covariant _PieChartPainter old) =>
      old.items != items ||
      old.colors != colors ||
      old.innerRadiusFactor != innerRadiusFactor;
}
