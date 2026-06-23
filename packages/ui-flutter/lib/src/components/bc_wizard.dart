library bc_wizard;

import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import '../core/types.dart';
import '../theme/bc_theme.dart';

/// Single wizard step descriptor.
///
/// Mirrors TS `BcWizardStep` ({ title, description?, blocking? }).
class BcWizardStep {
  final String title;
  final String? description;
  final bool blocking;

  const BcWizardStep({
    required this.title,
    this.description,
    this.blocking = false,
  });
}

/// Controller returned by [BcWizard.show] for programmatic step control.
///
/// [setStep] clamps the target index to `[0, steps.length - 1]` and emits
/// [onStepChange] through the original callback. When the per-step
/// [BcWizard.canAdvance] predicate returns `false` for an advance attempt,
/// the call is a no-op.
class BcWizardController {
  final int Function() _getCurrentStep;
  final void Function(int, BcChangeMeta) _doSetStep;
  final int _stepCount;
  final bool Function(int) _canAdvance;

  BcWizardController._(
    int Function() getCurrentStep,
    void Function(int, BcChangeMeta) doSetStep,
    bool Function(int) canAdvance,
    int stepCount,
  )  : _getCurrentStep = getCurrentStep,
        _doSetStep = doSetStep,
        _stepCount = stepCount,
        _canAdvance = canAdvance;
  int get currentStep => _getCurrentStep();

  void setStep(int target, BcChangeMeta meta) {
    final current = _getCurrentStep();
    if (target > current && !_canAdvance(current)) {
      return;
    }
    final clamped = _stepCount <= 0
        ? 0
        : target.clamp(0, _stepCount - 1);
    if (clamped == current) {
      return;
    }
    _doSetStep(clamped, meta);
  }

  void next(BcChangeMeta meta) => setStep(_getCurrentStep() + 1, meta);

  void previous(BcChangeMeta meta) => setStep(_getCurrentStep() - 1, meta);
}

/// Multi-step modal wizard with stepper.
///
/// Mirrors the TS `BcWizardContract`: a modal overlay rendering a stepper
/// (one entry per step) plus the current step's content area. Navigation
/// controls (Previous/Next/Finish/Cancel) drive step transitions.
///
/// The per-step [canAdvance] predicate gates forward movement: when it
/// returns `false` for the current step, the Next control is disabled and
/// a programmatic [BcWizardController.setStep] targeting a higher index is
/// a no-op. [BcWizardController.setStep] always clamps to `[0,
/// steps.length - 1]`, so `setStep(-1)` resolves to step 0 and
/// `setStep(steps.length)` resolves to the last index without throwing a
/// `RangeError`.
///
/// Any transition (pointer-driven Next/Previous, or programmatic via the
/// controller) emits [onStepChange] carrying the destination step index and
/// a [BcChangeMeta] whose source distinguishes the origin. The Finish
/// control is enabled only on the last step; activating it invokes
/// [onFinish] and dismisses the wizard. Cancel/Close disposes all
/// step-local controllers.
///
/// Controlled mode: the parent provides [step]; internal state is not
/// self-mutated on Next/Previous until the parent updates the prop.
/// Uncontrolled mode: [step] is null and [defaultStep] seeds the initial
/// state; Next/Previous update internally.
///
/// Platform-aware: renders a `CupertinoAlertDialog`-style layout on
/// iOS/macOS, `AlertDialog`-style elsewhere.
class BcWizard {
  BcWizard._();

  static BcWizardController show({
    required BuildContext context,
    required List<BcWizardStep> steps,
    int? step,
    int defaultStep = 0,
    bool Function(int)? canAdvance,
    String? title,
    String nextLabel = 'Next',
    String prevLabel = 'Previous',
    String finishLabel = 'Finish',
    String cancelLabel = 'Cancel',
    void Function(int index, BcChangeMeta meta)? onStepChange,
    void Function(BcChangeMeta meta)? onFinish,
    void Function(BcChangeMeta meta)? onCancel,
    List<Widget> Function(int index)? stepContentBuilder,
  }) {
    final gate = canAdvance ?? (_) => true;

    final route = _BcWizardRoute(
      context: context,
      steps: steps,
      step: step,
      defaultStep: defaultStep,
      canAdvance: gate,
      title: title,
      nextLabel: nextLabel,
      prevLabel: prevLabel,
      finishLabel: finishLabel,
      cancelLabel: cancelLabel,
      onStepChange: onStepChange,
      onFinish: onFinish,
      onCancel: onCancel,
      stepContentBuilder: stepContentBuilder,
    );

    Navigator.of(context, rootNavigator: true).push<dynamic>(route);

    return BcWizardController._(
      () => route.currentStepNotifier.value,
      (target, meta) => route.programmaticSetStep(target, meta),
      gate,
      steps.length,
    );
  }
}

class _BcWizardRoute extends PopupRoute<dynamic> {
  final BuildContext context;
  final List<BcWizardStep> steps;
  final int? step;
  final int defaultStep;
  final bool Function(int) canAdvance;
  final String? title;
  final String nextLabel;
  final String prevLabel;
  final String finishLabel;
  final String cancelLabel;
  final void Function(int index, BcChangeMeta meta)? onStepChange;
  final void Function(BcChangeMeta meta)? onFinish;
  final void Function(BcChangeMeta meta)? onCancel;
  final List<Widget> Function(int index)? stepContentBuilder;

  final ValueNotifier<int> currentStepNotifier;
  void Function(int target, BcChangeMeta meta)? _programmaticSetStepHandler;

  _BcWizardRoute({
    required this.context,
    required this.steps,
    required this.step,
    required this.defaultStep,
    required this.canAdvance,
    required this.title,
    required this.nextLabel,
    required this.prevLabel,
    required this.finishLabel,
    required this.cancelLabel,
    required this.onStepChange,
    required this.onFinish,
    required this.onCancel,
    required this.stepContentBuilder,
  }) : currentStepNotifier = ValueNotifier<int>(
          steps.isEmpty
              ? 0
              : defaultStep.clamp(0, steps.length - 1),
        );

  void registerProgrammaticSetStepHandler(
      void Function(int, BcChangeMeta) handler) {
    _programmaticSetStepHandler = handler;
  }

  void programmaticSetStep(int target, BcChangeMeta meta) {
    _programmaticSetStepHandler?.call(target, meta);
  }

  @override
  Color? get barrierColor => Colors.black54;

  @override
  bool get barrierDismissible => false;

  @override
  String? get barrierLabel => null;

  @override
  Widget buildPage(
    BuildContext context,
    Animation<double> animation,
    Animation<double> secondaryAnimation,
  ) {
    return _BcWizardPage(route: this);
  }

  @override
  Widget buildTransitions(
    BuildContext context,
    Animation<double> animation,
    Animation<double> secondaryAnimation,
    Widget child,
  ) {
    return FadeTransition(opacity: animation, child: child);
  }

  @override
  Duration get transitionDuration => const Duration(milliseconds: 200);

  void dismiss([dynamic result]) {
    if (isCurrent) {
      Navigator.of(context).pop(result);
    }
  }
}

class _BcWizardPage extends StatefulWidget {
  final _BcWizardRoute route;

  const _BcWizardPage({required this.route});

  @override
  State<_BcWizardPage> createState() => _BcWizardPageState();
}

class _BcWizardPageState extends State<_BcWizardPage> {
  late int _internalStep;

  int get currentStep => _internalStep;

  @override
  void initState() {
    super.initState();
    final seed = widget.route.step ?? widget.route.defaultStep;
    if (widget.route.steps.isNotEmpty) {
      _internalStep = seed.clamp(0, widget.route.steps.length - 1);
    } else {
      _internalStep = 0;
    }
    widget.route.currentStepNotifier.value = currentStep;
    widget.route.registerProgrammaticSetStepHandler(programmaticSetStep);
  }

  @override
  void didUpdateWidget(covariant _BcWizardPage oldWidget) {
    super.didUpdateWidget(oldWidget);
  }

  void _notifyState() {
    widget.route.currentStepNotifier.value = currentStep;
  }

  void _transitionTo(int target, BcChangeMeta meta) {
    final clamped = widget.route.steps.isEmpty
        ? 0
        : target.clamp(0, widget.route.steps.length - 1);
    if (clamped == currentStep) return;

    setState(() {
      _internalStep = clamped;
    });
    widget.route.onStepChange?.call(clamped, meta);
    _notifyState();
  }

  void programmaticSetStep(int target, BcChangeMeta meta) {
    final current = currentStep;
    if (target > current && !widget.route.canAdvance(current)) {
      return;
    }
    _transitionTo(target, meta);
  }

  void _handleNext() {
    _transitionTo(
      currentStep + 1,
      const BcChangeMeta(
        reason: BcChangeReason.increment,
        source: BcEventSource.pointer,
      ),
    );
  }

  void _handlePrev() {
    _transitionTo(
      currentStep - 1,
      const BcChangeMeta(
        reason: BcChangeReason.decrement,
        source: BcEventSource.pointer,
      ),
    );
  }

  void _handleFinish() {
    widget.route.onFinish?.call(
      const BcChangeMeta(
        reason: BcChangeReason.submit,
        source: BcEventSource.pointer,
      ),
    );
    widget.route.dismiss(true);
  }

  void _handleCancel() {
    widget.route.onCancel?.call(
      const BcChangeMeta(
        reason: BcChangeReason.cancel,
        source: BcEventSource.pointer,
      ),
    );
    widget.route.dismiss(false);
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final platform = Theme.of(context).platform;
    final isIOS = platform == TargetPlatform.iOS ||
        platform == TargetPlatform.macOS;
    return isIOS ? _buildCupertino() : _buildMaterial();
  }

  Widget _buildMaterial() {
    final isLast = currentStep == widget.route.steps.length - 1;
    final isFirst = currentStep == 0;
    final canAdvanceNow = widget.route.canAdvance(currentStep);

    return Center(
      child: AlertDialog(
        key: const Key('bc-wizard-dialog'),
        title: widget.route.title != null ? Text(widget.route.title!) : null,
        content: SizedBox(
          width: 400,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildStepper(),
              const Divider(height: 24),
              Flexible(child: _buildContent()),
            ],
          ),
        ),
        actions: [
          TextButton(
            key: const Key('bc-wizard-cancel'),
            onPressed: _handleCancel,
            child: Text(widget.route.cancelLabel),
          ),
          if (!isFirst)
            TextButton(
              key: const Key('bc-wizard-prev'),
              onPressed: _handlePrev,
              child: Text(widget.route.prevLabel),
            ),
          if (isLast)
            FilledButton(
              key: const Key('bc-wizard-finish'),
              onPressed: _handleFinish,
              child: Text(widget.route.finishLabel),
            )
          else
            FilledButton(
              key: const Key('bc-wizard-next'),
              onPressed: canAdvanceNow ? _handleNext : null,
              child: Text(widget.route.nextLabel),
            ),
        ],
      ),
    );
  }

  Widget _buildCupertino() {
    final isLast = currentStep == widget.route.steps.length - 1;
    final isFirst = currentStep == 0;
    final canAdvanceNow = widget.route.canAdvance(currentStep);

    return Center(
      child: CupertinoAlertDialog(
        key: const Key('bc-wizard-dialog'),
        title: widget.route.title != null ? Text(widget.route.title!) : null,
        content: Padding(
          padding: const EdgeInsets.only(top: 8),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildStepper(),
              const SizedBox(height: 12),
              _buildContent(),
            ],
          ),
        ),
        actions: [
          CupertinoDialogAction(
            key: const Key('bc-wizard-cancel'),
            onPressed: _handleCancel,
            child: Text(widget.route.cancelLabel),
          ),
          if (!isFirst)
            CupertinoDialogAction(
              key: const Key('bc-wizard-prev'),
              onPressed: _handlePrev,
              child: Text(widget.route.prevLabel),
            ),
          if (isLast)
            CupertinoDialogAction(
              key: const Key('bc-wizard-finish'),
              isDefaultAction: true,
              onPressed: _handleFinish,
              child: Text(widget.route.finishLabel),
            )
          else
            CupertinoDialogAction(
              key: const Key('bc-wizard-next'),
              isDefaultAction: true,
              onPressed: canAdvanceNow ? _handleNext : null,
              child: Text(widget.route.nextLabel),
            ),
        ],
      ),
    );
  }

  Widget _buildStepper() {
    return Container(
      key: const Key('bc-wizard-stepper'),
      child: Row(
        children: [
          for (var i = 0; i < widget.route.steps.length; i++) ...[
            _buildStepNode(i),
            if (i < widget.route.steps.length - 1)
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 4),
                  child: Container(
                    height: 2,
                    color: i < currentStep
                        ? BcTheme.primary
                        : BcTheme.border,
                  ),
                ),
              ),
          ],
        ],
      ),
    );
  }

  Widget _buildStepNode(int index) {
    final isCurrent = index == currentStep;
    final isCompleted = index < currentStep;
    final step = widget.route.steps[index];

    final circle = Container(
      width: 24,
      height: 24,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: isCurrent
            ? BcTheme.primary
            : isCompleted
                ? BcTheme.primary
                : BcTheme.muted,
        border: isCurrent
            ? Border.all(color: BcTheme.ring, width: 2)
            : null,
      ),
      child: Center(
        child: isCompleted
            ? const Icon(Icons.check, size: 14, color: BcTheme.primaryForeground)
            : Text(
                '${index + 1}',
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: isCurrent
                      ? BcTheme.primaryForeground
                      : BcTheme.mutedForeground,
                ),
              ),
      ),
    );

    return Column(
      key: isCurrent
          ? const Key('bc-wizard-marker-current')
          : ValueKey<String>('bc-wizard-marker-$index'),
      mainAxisSize: MainAxisSize.min,
      children: [
        circle,
        const SizedBox(height: 4),
        ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 80),
          child: Text(
            step.title,
            textAlign: TextAlign.center,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
            style: TextStyle(
              fontSize: 11,
              fontWeight: isCurrent ? FontWeight.w600 : FontWeight.normal,
              color: isCurrent ? BcTheme.foreground : BcTheme.mutedForeground,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildContent() {
    if (widget.route.stepContentBuilder != null) {
      return Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: widget.route.stepContentBuilder!(currentStep),
      );
    }
    if (widget.route.steps.isEmpty) return const SizedBox.shrink();
    final step = widget.route.steps[currentStep];
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          step.title,
          style: const TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
            color: BcTheme.foreground,
          ),
        ),
        if (step.description != null) ...[
          const SizedBox(height: 8),
          Text(
            step.description!,
            style: const TextStyle(
              fontSize: 14,
              color: BcTheme.mutedForeground,
            ),
          ),
        ],
      ],
    );
  }
}
