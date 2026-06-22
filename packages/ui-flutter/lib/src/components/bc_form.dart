library bc_form;

import 'package:flutter/material.dart';
import '../core/reactive_store.dart';

/// Material-only form container with reactive store binding (ADR-011 Section 5).
/// Pure Material 3 — no Cupertino counterpart.
class BcForm extends StatefulWidget {
  final Map<String, dynamic>? initialValues;
  final List<Reaction>? reactions;
  final void Function(ReactiveStore store)? onCreateStore;
  final Widget child;
  final void Function(Map<String, dynamic> values)? onSubmit;
  const BcForm({
    super.key,
    this.initialValues,
    this.reactions,
    this.onCreateStore,
    required this.child,
    this.onSubmit,
  });

  @override
  State<BcForm> createState() => _BcFormState();

  /// Access the nearest BcForm's store via context.
  static ReactiveStore? storeOf(BuildContext context) {
    final scope = context
        .dependOnInheritedWidgetOfExactType<_BcFormScope>();
    return scope?.store;
  }
}

class _BcFormState extends State<BcForm> {
  late ReactiveStore _store;

  @override
  void initState() {
    super.initState();
    _store = ReactiveStore(ReactiveStoreOptions(
      initialValues: widget.initialValues,
      reactions: widget.reactions,
    ));
    _store.runAll();
    widget.onCreateStore?.call(_store);
  }

  @override
  void dispose() {
    _store.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return _BcFormScope(
      store: _store,
      child: Form(
        child: widget.child,
      ),
    );
  }
}

class _BcFormScope extends InheritedWidget {
  final ReactiveStore store;

  const _BcFormScope({required this.store, required super.child});

  @override
  bool updateShouldNotify(_BcFormScope oldWidget) =>
      store != oldWidget.store;
}
