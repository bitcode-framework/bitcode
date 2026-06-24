library bc_chatter;

import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_platform_widgets/flutter_platform_widgets.dart';
import '../core/types.dart';
import '../theme/bc_theme.dart';

/// Chat / messaging widget with message list, compose input and send.
///
/// Mirrors TS `BcChatter` contract and validation contract VAL-M6-017/018.
///
/// VAL-M6-017: renders N message bubbles in order. Typing into the input and
/// pressing Enter or the send button clears the input and delivers the
/// composed text to [onSend].
///
/// VAL-M6-018: sending emits [onSend] with the trimmed text and a
/// BcChangeMeta whose source is pointer. Empty / whitespace-only messages
/// are rejected (never fire [onSend]). When [disabled] is true the input is
/// non-interactive, the send button is disabled, and no [onSend] fires.
class BcChatter extends StatefulWidget {
  final List<BcChatterMessage> messages;
  final String? placeholder;
  final bool disabled;
  final String? currentUser;
  final String emptyText;
  final Widget? emptyState;
  final BcChangeHandler<String>? onSend;

  const BcChatter({
    super.key,
    this.messages = const [],
    this.placeholder,
    this.disabled = false,
    this.currentUser,
    this.emptyText = 'No messages',
    this.emptyState,
    this.onSend,
  });

  @override
  State<BcChatter> createState() => _BcChatterState();
}

class _BcChatterState extends State<BcChatter> {
  final TextEditingController _inputController = TextEditingController();
  final FocusNode _focusNode = FocusNode();
  final ScrollController _scrollController = ScrollController();

  @override
  void dispose() {
    _inputController.dispose();
    _focusNode.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  bool _isMineMessage(BcChatterMessage msg) {
    return widget.currentUser != null && msg.author == widget.currentUser;
  }

  void _send({required BcEventSource source}) {
    if (widget.disabled) return;
    final text = _inputController.text.trim();
    if (text.isEmpty) return;

    widget.onSend?.call(
      text,
      BcChangeMeta(
        reason: BcChangeReason.submit,
        source: source,
      ),
    );
    _inputController.clear();

    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted && _scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 200),
          curve: Curves.easeOut,
        );
      }
    });
  }

  Widget _buildAvatar(BcChatterMessage msg, bool mine) {
    if (msg.avatar != null && msg.avatar!.isNotEmpty) {
      return CircleAvatar(
        radius: 14,
        backgroundImage: NetworkImage(msg.avatar!),
        backgroundColor: BcTheme.muted,
      );
    }
    final initial = msg.author.isEmpty
        ? '?'
        : msg.author.substring(0, 1).toUpperCase();
    return Container(
      width: 28,
      height: 28,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: mine ? BcTheme.primary : BcTheme.muted,
        shape: BoxShape.circle,
      ),
      child: Text(
        initial,
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: mine ? BcTheme.primaryForeground : BcTheme.foreground,
        ),
      ),
    );
  }

  Widget _buildMessage(BcChatterMessage msg, int index) {
    final mine = _isMineMessage(msg);
    return Container(
      key: ValueKey<String>('bc_chatter_msg_${msg.id}_$index'),
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      child: Row(
        mainAxisAlignment:
            mine ? MainAxisAlignment.end : MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (!mine) ...[_buildAvatar(msg, mine), const SizedBox(width: 6)],
          Flexible(
            child: Container(
              constraints: BoxConstraints(
                maxWidth: MediaQuery.of(context).size.width * 0.7,
              ),
              padding: const EdgeInsets.symmetric(
                  horizontal: 10, vertical: 6),
              decoration: BoxDecoration(
                color: mine ? BcTheme.primary : BcTheme.card,
                borderRadius: BorderRadius.only(
                  topLeft: const Radius.circular(BcTheme.radius),
                  topRight: const Radius.circular(BcTheme.radius),
                  bottomLeft: mine
                      ? const Radius.circular(BcTheme.radius)
                      : Radius.zero,
                  bottomRight: mine
                      ? Radius.zero
                      : const Radius.circular(BcTheme.radius),
                ),
                border: mine ? null : Border.all(color: BcTheme.border),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (!mine)
                    Text(
                      msg.author,
                      key: ValueKey<String>(
                          'bc_chatter_author_${msg.id}_$index'),
                      style: const TextStyle(
                        fontSize: 11,
                        fontWeight: FontWeight.w600,
                        color: BcTheme.mutedForeground,
                      ),
                    ),
                  Text(
                    msg.content,
                    key: ValueKey<String>(
                        'bc_chatter_content_${msg.id}_$index'),
                    style: TextStyle(
                      fontSize: 13,
                      color: mine ? BcTheme.primaryForeground : BcTheme.foreground,
                    ),
                  ),
                  if (msg.timestamp != null)
                    Text(
                      msg.timestamp!,
                      key: ValueKey<String>(
                          'bc_chatter_timestamp_${msg.id}_$index'),
                      style: TextStyle(
                        fontSize: 10,
                        color: mine
                            ? BcTheme.primaryForeground.withValues(alpha: 0.7)
                            : BcTheme.mutedForeground,
                      ),
                    ),
                ],
              ),
            ),
          ),
          if (mine) ...[const SizedBox(width: 6), _buildAvatar(msg, mine)],
        ],
      ),
    );
  }

  Widget _buildInput() {
    final platform = Theme.of(context).platform;
    final isIOS = platform == TargetPlatform.iOS ||
        platform == TargetPlatform.macOS;

    final input = PlatformTextField(
      key: const ValueKey<String>('bc_chatter_input'),
      controller: _inputController,
      focusNode: _focusNode,
      enabled: !widget.disabled,
      keyboardType: TextInputType.multiline,
      textInputAction: TextInputAction.send,
      onSubmitted: (_) => _send(source: BcEventSource.pointer),
      material: (_, __) => MaterialTextFieldData(
        decoration: InputDecoration(
          isDense: true,
          hintText: widget.placeholder ?? 'Type a message...',
          hintStyle: const TextStyle(
              fontSize: 13, color: BcTheme.mutedForeground),
          contentPadding: const EdgeInsets.symmetric(
              horizontal: 10, vertical: 8),
          enabledBorder: const OutlineInputBorder(
            borderRadius:
                BorderRadius.all(Radius.circular(BcTheme.radius)),
            borderSide: BorderSide(color: BcTheme.border),
          ),
          focusedBorder: const OutlineInputBorder(
            borderRadius:
                BorderRadius.all(Radius.circular(BcTheme.radius)),
            borderSide: BorderSide(color: BcTheme.primary),
          ),
          disabledBorder: OutlineInputBorder(
            borderRadius:
                const BorderRadius.all(Radius.circular(BcTheme.radius)),
            borderSide: BorderSide(
                color: BcTheme.border.withValues(alpha: 0.5)),
          ),
        ),
        style: const TextStyle(fontSize: 13),
      ),
      cupertino: (_, __) => CupertinoTextFieldData(
        placeholder: widget.placeholder ?? 'Type a message...',
        padding: const EdgeInsets.symmetric(
            horizontal: 10, vertical: 8),
        decoration: BoxDecoration(
          color: widget.disabled ? BcTheme.muted : BcTheme.background,
          border: Border.all(color: BcTheme.border),
          borderRadius:
              const BorderRadius.all(Radius.circular(BcTheme.radius)),
        ),
        style: const TextStyle(fontSize: 13),
      ),
    );

    return Padding(
      key: const ValueKey<String>('bc_chatter_compose'),
      padding: const EdgeInsets.fromLTRB(8, 4, 8, 8),
      child: Row(
        children: [
          Expanded(child: input),
          const SizedBox(width: 6),
          PlatformIconButton(
            key: const ValueKey<String>('bc_chatter_send'),
            onPressed: widget.disabled ? null : () => _send(source: BcEventSource.pointer),
            icon: Icon(
              isIOS ? CupertinoIcons.arrow_up_circle_fill : Icons.send,
              size: 24,
              color: widget.disabled ? BcTheme.mutedForeground : BcTheme.primary,
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      key: const ValueKey<String>('bc_chatter_root'),
      decoration: BoxDecoration(
        color: BcTheme.background,
        border: Border.all(color: BcTheme.border),
        borderRadius: const BorderRadius.all(Radius.circular(BcTheme.radius)),
      ),
      child: Column(
        children: [
          Expanded(
            child: widget.messages.isEmpty
                ? (widget.emptyState ??
                    Center(
                      key: const ValueKey<String>('bc_chatter_empty'),
                      child: Text(
                        widget.emptyText,
                        style: const TextStyle(
                          fontSize: 14,
                          color: BcTheme.mutedForeground,
                        ),
                      ),
                    ))
                : ListView.builder(
                    key: const ValueKey<String>('bc_chatter_list'),
                    controller: _scrollController,
                    padding: const EdgeInsets.symmetric(vertical: 8),
                    itemCount: widget.messages.length,
                    itemBuilder: (context, index) =>
                        _buildMessage(widget.messages[index], index),
                  ),
          ),
          _buildInput(),
        ],
      ),
    );
  }
}
