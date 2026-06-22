// ============================================================================
// @bitcode/ui-schema — Root barrel
//
// Standalone UI schema definitions — domain-agnostic, framework-neutral types
// for schema-driven UI. Consumed by:
//   - @bitcode/ui (React/Vue/RN via Mitosis)
//   - @bitcode/ui-flutter (Dart port, Phase 5a)
//   - BitCode engine (optional mapping from domain models to UiSchema)
//   - Any standalone consumer (JsonForms-style usage)
//
// Zero framework dependency. Zero platform globals. Safe in Node/SSR/RN/Dart.
// ============================================================================

// Core types
export type {
  BcChangeHandler,
  BcChangeMeta,
  BcEventBase,
  BcEventSource,
  BcKey,
  BcSelection,
  ComponentSize,
  ComponentVariant,
  DataSource,
  FetchParams,
  FetchResult,
  FieldBlurEvent,
  FieldChangeEvent,
  FieldClearEvent,
  FieldInputEvent,
  FieldKeyEvent,
  FieldValidateEvent,
  InputSize,
  InputType,
  ModelAdapter,
  ResolvedTheme,
  ThemeMode,
  ValidationStatus,
} from './types.ts';

// Reactive engine
export { ReactiveStore } from './reactive/store.ts';
export type { Listener } from './reactive/store.ts';
export type {
  CompiledExpr,
  DataEffect,
  Effect,
  Evaluator,
  FieldMeta,
  FieldPath,
  Reaction,
  ReactiveScopeLike,
  ReactiveStoreOptions,
  Snapshot,
  UiEffect,
  ValidationEffect,
  ValidationRule,
  ValueEffect,
} from './reactive/types.ts';

// Schema + resolver
export type {
  BcRendererProps,
  FieldNode,
  FieldType,
  UiAction,
  UiSchema,
  UiSchemaKind,
} from './schema/types.ts';

export {
  defaultFieldResolver,
  resolveField,
  resolveFieldType,
  resolveFields,
} from './schema/resolver.ts';

export type {
  FieldResolver,
  ResolvedFormField,
} from './schema/resolver.ts';
