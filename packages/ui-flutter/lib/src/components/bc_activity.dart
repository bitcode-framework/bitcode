library bc_activity;

import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_platform_widgets/flutter_platform_widgets.dart';
import '../core/types.dart';
import '../theme/bc_theme.dart';

/// Social activity timeline that renders entries in chronological order with
/// timestamps and actor labels.
///
/// Mirrors TS `BcActivity` contract and validation contract VAL-M6-019/020.
///
/// VAL-M6-019: renders N entries in the provided order, each showing actor
/// label and formatted timestamp. Empty items renders emptyState (or default).
///
/// VAL-M6-020: tapping an item emits onItemClick with item id and a
/// BcChangeMeta (pointer source). When disabled, no callback fires. Items
/// without an id do not crash the handler.
class BcActivity extends StatefulWidget {
  final List<BcActivityEntry> items;
  final bool disabled;
  final String emptyText;
  final Widget? emptyState;
  final BcChangeHandler<String>? onItemClick;

  const BcActivity({
    super.key,
    this.items = const [],
    this.disabled = false,
    this.emptyText = 'No activity',
    this.emptyState,
    this.onItemClick,
  });

  @override
  State<BcActivity> createState() => _BcActivityState();
}

class _BcActivityState extends State<BcActivity> {
  Color _variantColor(BcActivityVariant variant) {
    switch (variant) {
      case BcActivityVariant.success:
        return BcTheme.success;
      case BcActivityVariant.warning:
        return BcTheme.warning;
      case BcActivityVariant.destructive:
        return BcTheme.destructive;
      case BcActivityVariant.info:
        return BcTheme.info;
      case BcActivityVariant.defaultVariant:
        return BcTheme.primary;
    }
  }

  IconData _variantIcon(BcActivityVariant variant) {
    switch (variant) {
      case BcActivityVariant.success:
        return Icons.check_circle_outline;
      case BcActivityVariant.warning:
        return Icons.warning_amber_outlined;
      case BcActivityVariant.destructive:
        return Icons.error_outline;
      case BcActivityVariant.info:
        return Icons.info_outline;
      case BcActivityVariant.defaultVariant:
        return Icons.radio_button_unchecked;
    }
  }

  void _onTap(BcActivityEntry entry) {
    if (widget.disabled) return;
    widget.onItemClick?.call(
      entry.id,
      const BcChangeMeta(
        reason: BcChangeReason.select,
        source: BcEventSource.pointer,
      ),
    );
  }

  Widget _buildAvatar(BcActivityEntry entry) {
    final initial = entry.actor.isEmpty
        ? '?'
        : entry.actor.substring(0, 1).toUpperCase();
    return Container(
      width: 32,
      height: 32,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: _variantColor(entry.variant),
        shape: BoxShape.circle,
      ),
      child: Text(
        initial,
        style: const TextStyle(
          fontSize: 14,
          fontWeight: FontWeight.w600,
          color: BcTheme.primaryForeground,
        ),
      ),
    );
  }

  Widget _buildItem(BcActivityEntry entry, int index) {
    final content = Container(
      key: ValueKey<String>('bc_activity_item_${entry.id}_$index'),
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(
        color: Colors.transparent,
        border: Border.all(color: BcTheme.border),
        borderRadius: const BorderRadius.all(Radius.circular(BcTheme.radius)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildAvatar(entry),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Flexible(
                      child: Text(
                        entry.actor,
                        key: ValueKey<String>(
                            'bc_activity_actor_${entry.id}_$index'),
                        style: const TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w600,
                          color: BcTheme.foreground,
                        ),
                      ),
                    ),
                    const SizedBox(width: 6),
                    Flexible(
                      child: Text(
                        entry.action,
                        key: ValueKey<String>(
                            'bc_activity_action_${entry.id}_$index'),
                        style: const TextStyle(
                          fontSize: 13,
                          color: BcTheme.mutedForeground,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 2),
                Row(
                  children: [
                    Icon(
                      _variantIcon(entry.variant),
                      size: 12,
                      color: _variantColor(entry.variant),
                    ),
                    const SizedBox(width: 4),
                    Flexible(
                      child: Text(
                        entry.timestamp,
                        key: ValueKey<String>(
                            'bc_activity_timestamp_${entry.id}_$index'),
                        style: const TextStyle(
                          fontSize: 11,
                          color: BcTheme.mutedForeground,
                        ),
                      ),
                    ),
                  ],
                ),
                if (entry.target != null && entry.target!.isNotEmpty) ...[
                  const SizedBox(height: 2),
                  Text(
                    entry.target!,
                    key: ValueKey<String>(
                        'bc_activity_target_${entry.id}_$index'),
                    style: const TextStyle(
                      fontSize: 12,
                      color: BcTheme.mutedForeground,
                    ),
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );

    if (widget.disabled) {
      return Opacity(opacity: 0.6, child: content);
    }

    return PlatformWidget(
      material: (ctx, _) => Material(
        color: Colors.transparent,
        child: InkWell(
        onTap: () => _onTap(entry),
        borderRadius: const BorderRadius.all(Radius.circular(BcTheme.radius)),
        child: content,
        ),
      ),
      cupertino: (ctx, _) => CupertinoButton(
        onPressed: () => _onTap(entry),
        padding: EdgeInsets.zero,
        minimumSize: Size.zero,
        child: content,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (widget.items.isEmpty) {
      return widget.emptyState ??
          Center(
            key: const ValueKey<String>('bc_activity_empty'),
            child: Text(
              widget.emptyText,
              style: const TextStyle(
                fontSize: 14,
                color: BcTheme.mutedForeground,
              ),
            ),
          );
    }

    return Semantics(
      container: true,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          for (var i = 0; i < widget.items.length; i++)
            Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: _buildItem(widget.items[i], i),
            ),
        ],
      ),
    );
  }
}
