// ============================================================================
// @bitcode/ui-schema - Block Catalog module barrel (ADR-017)
//
// Canonical Block Catalog contracts: entry types, candidate metadata,
// category/channel/capability/readiness/render-strategy enums, policies,
// migration/a11y/test metadata, capability normalization, and a minimal
// shape validator. Pure types + pure functions. No runtime, no renderer.
// ============================================================================

// Canonical types
export type {
  BlockAdaptivePolicy,
  BlockA11yEvidence,
  BlockCachePolicy,
  BlockCapability,
  BlockCatalogEntry,
  BlockCategory,
  BlockChannel,
  BlockMigrationMetadata,
  BlockOfflinePolicy,
  BlockPlatform,
  BlockPlatformFamily,
  BlockReadiness,
  BlockRenderStrategy,
  BlockSecurityPolicy,
  BlockTestReference,
  CandidateCatalogEntry,
  RequiredAdapter,
  RequiredComponent,
  RequiredFunction,
} from './types.ts';

// Capability normalization helpers
export {
  BLOCK_CAPABILITY_ALIASES,
  CANONICAL_BLOCK_CAPABILITIES,
  findCapabilityAliases,
  isCanonicalCapability,
  isKnownCapabilityAlias,
  normalizeCapability,
  normalizeCapabilities,
} from './capability.ts';

// Validator
export {
  BLOCK_DIAGNOSTIC_CODES,
  validateBlockCatalogEntry,
  validateCandidateCatalogEntry,
} from './validate.ts';
export type {
  BlockCatalogDiagnostic,
  BlockCatalogValidationResult,
  BlockDiagnosticSeverity,
} from './validate.ts';

// Candidate fixtures (@internal - not a Ready catalog)
export {
  candidateContentHero,
  candidateFixtures,
  candidateInteractionNotificationCenter,
} from './fixtures.ts';
