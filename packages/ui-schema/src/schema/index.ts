// @bitcode/ui-schema — Schema module barrel (ADR-005)
export type {
  BcRendererProps,
  FieldNode,
  FieldType,
  UiAction,
  UiSchema,
  UiSchemaKind,
} from './types.ts';

export {
  defaultFieldResolver,
  resolveField,
  resolveFieldType,
  resolveFields,
} from './resolver.ts';

export type {
  FieldResolver,
  ResolvedFormField,
} from './resolver.ts';
