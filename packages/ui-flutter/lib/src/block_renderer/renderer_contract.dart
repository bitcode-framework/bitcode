/// BlockRenderer Contract - Dart mirror of @bitcode/ui-schema block-renderer/
/// renderer.ts (ADR-017 §3, Agent H).
///
/// The portable block composition engine CONTRACT. This file defines the
/// interfaces the Flutter widget (follow-up) will implement. The full engine
/// port + Flutter BcBlockRenderer widget is a dedicated follow-up; this
/// contract mirror keeps Dart parity with the portable engine so the
/// follow-up can land against a stable surface.
///
/// DESIGN (matches the TypeScript engine):
///   - Emits render plans only — never widgets.
///   - Registry stores metadata only.
///   - Actions bound as deferred plans; executed on events, not during render.
///   - Dependencies injected (no globals, ADR-018 S7 discipline).
///
/// DEFERRED (amendment 1 - G not yet ported to Dart): the [BlockRendererDeps]
/// `actionRunner` field references ActionRunner, part of the ADR-018 action
/// contract (Agent G). Agent G landed the TypeScript ActionRunner (commit
/// a8a28a8) but the Dart port does not yet exist. The actionRunner dependency
/// is typed as `Object?` here so the contract mirror compiles today; replace
/// with the typed `ActionRunner?` when Agent G's Dart port lands.
library bitcode.ui_schema.block_renderer.renderer_contract;

import 'block_node.dart'
    show
        BlockNode,
        RenderOptions,
        ResolvedBlockTree;
import 'registry.dart' show BlockComponentRegistry;
import '../expressions/evaluator.dart' show ExpressionEvaluator;

/// The catalog resolver the engine consults.
abstract class CatalogResolver {
  /// Resolve a catalog entry by id (+ optional version). Null when unknown.
  Object? resolve(String id, String? version);
}

/// The runtime scope surface the engine reads bindings from + evaluates ASTs
/// against. Built from the host's current state.
class BlockRenderScope {
  const BlockRenderScope({
    required this.values,
    required this.meta,
    required this.data,
    required this.services,
    this.permissions,
    this.args,
  });

  final Map<String, dynamic> values;
  final Map<String, dynamic> meta;
  final Map<String, dynamic> data;
  final Map<String, dynamic> services;
  /// Permission snapshot. DEFERRED typing until the full permission surface is
  /// finalized; passed through to the F permission-filter builder verbatim.
  final Object? permissions;
  final Map<String, dynamic>? args;
}

/// Runtime availability surface (registry + functions + adapters + platform).
/// The validator checks declared requirements against this without coupling to
/// the full ActionRuntimeContext.
class RuntimeAvailability {
  const RuntimeAvailability({
    required this.componentRegistry,
    required this.registeredFunctions,
    required this.registeredAdapters,
    this.platform,
  });

  final BlockComponentRegistry componentRegistry;
  final Set<String> registeredFunctions;
  final Set<String> registeredAdapters;
  final String? platform;
}

/// Dependencies the BlockRenderer is constructed with (injected, no globals).
class BlockRendererDeps {
  const BlockRendererDeps({
    required this.catalog,
    required this.componentRegistry,
    this.evaluator,
    this.actionRunner,
  });

  final CatalogResolver catalog;
  final BlockComponentRegistry componentRegistry;
  final ExpressionEvaluator? evaluator;
  /// DEFERRED: typed as `Object?` until the ADR-018 action contract (Agent G)
  /// Dart port lands. The TS engine consumes an `ActionRunner` instance.
  final Object? actionRunner;
}

/// BlockRenderer engine contract. Implementations resolve a BlockNode tree
/// into a ResolvedBlockTree of render plans. Never emits widgets; never throws
/// (failures captured as DiagnosticPlan entries).
abstract class BlockRenderer {
  /// Resolve a block node tree into render plans.
  ///
  /// [node] root block node (null -> error plan).
  /// [scope] runtime scope (values/data/meta/permissions).
  /// [availability] runtime availability (registry/functions/adapters).
  /// [options] render options (allowCandidate, platform).
  ResolvedBlockTree render(
    BlockNode? node,
    BlockRenderScope scope,
    RuntimeAvailability availability, [
    RenderOptions options = const RenderOptions(),
  ]);
}
