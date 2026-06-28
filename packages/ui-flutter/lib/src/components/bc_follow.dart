library bc_follow;

import 'package:flutter/material.dart';

import '../core/types.dart';
import '../theme/bc_theme.dart';

/// Reason the follow state changed.
enum BcFollowChangeReason { follow, unfollow, programmatic }

/// Follow/unfollow toggle with optional follower count (ADR-007 parity).
class BcFollow extends StatefulWidget {
  final bool? following;
  final bool defaultFollowing;
  final int? count;
  final bool disabled;
  final String followLabel;
  final String followingLabel;
  final BcChangeHandler<bool>? onFollowingChange;

  const BcFollow({
    super.key,
    this.following,
    this.defaultFollowing = false,
    this.count,
    this.disabled = false,
    this.followLabel = 'Follow',
    this.followingLabel = 'Following',
    this.onFollowingChange,
  });

  @override
  State<BcFollow> createState() => _BcFollowState();
}

class _BcFollowState extends State<BcFollow> {
  late bool _internalFollowing;
  bool get _isControlled => widget.following != null;
  bool get _isFollowing => _isControlled ? widget.following! : _internalFollowing;

  @override
  void initState() {
    super.initState();
    _internalFollowing = widget.defaultFollowing;
  }

  void _toggle() {
    if (widget.disabled) return;
    final next = !_isFollowing;
    if (!_isControlled) {
      setState(() => _internalFollowing = next);
    }
    widget.onFollowingChange?.call(
      next,
      BcChangeMeta(
        reason: next ? BcChangeReason.submit : BcChangeReason.cancel,
        source: BcEventSource.pointer,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final following = _isFollowing;
    final bg = following ? BcTheme.secondary : BcTheme.primary;
    final fg = following ? BcTheme.secondaryForeground : BcTheme.primaryForeground;
    return MergeSemantics(
      child: Semantics(
        label: following ? widget.followingLabel : widget.followLabel,
        button: true,
        toggled: following,
        enabled: !widget.disabled,
        child: Opacity(
          opacity: widget.disabled ? 0.5 : 1.0,
          child: Material(
            color: Colors.transparent,
            child: InkWell(
              onTap: widget.disabled ? null : _toggle,
              borderRadius: BorderRadius.circular(BcTheme.radius),
              child: Container(
                key: const Key('bc-follow-root'),
                constraints: const BoxConstraints(minHeight: 36),
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                decoration: BoxDecoration(
                  color: bg,
                  borderRadius: BorderRadius.circular(BcTheme.radius),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      following ? widget.followingLabel : widget.followLabel,
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w500,
                        color: fg,
                      ),
                    ),
                    if (widget.count != null) ...[
                      const SizedBox(width: 8),
                      Text(
                        '${widget.count}',
                        style: TextStyle(
                          fontSize: 12,
                          color: fg.withValues(alpha: 0.8),
                          fontFeatures: [const FontFeature.tabularFigures()],
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
