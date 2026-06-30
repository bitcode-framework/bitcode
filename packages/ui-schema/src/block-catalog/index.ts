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

// ============================================================================
// Policy validators (ADR-017/018, Agent L)
//
// Compose over the shape validator above. Same diagnostic model, separate
// BLOCK_POLICY_* code space. Never flip readiness (no Ready promotion).
// ============================================================================
export {
  DEFAULT_POLICY_CONFIG,
  gateSeverity,
  hasAnyCapability,
  hasCapability,
  isDeprecatedEntry,
  isReadyEntry,
  mergeResults,
  POLICY_DIAGNOSTIC_CODES,
  policyDiagnostic,
  resolvePolicyConfig,
  resultFromDiagnostics,
  validateAccessibilityPolicy,
  validateAdaptivePolicy,
  validateBlockCatalogPolicies,
  validateOfflinePolicy,
  validatePerformancePolicy,
  validateSecurityPolicy,
} from './policies/index.ts';
export type {
  PolicyConfig,
  PolicyResult,
  ResolvedPolicyConfig,
} from './policies/index.ts';

// ============================================================================
// Per-category Candidate metadata seed + platform audit + pilot analysis
// (ADR-017, Agent J). Canonical Candidate metadata for all 44 matrix rows.
// Every entry is `readiness: 'candidate'` — NO Ready promotion.
// ============================================================================

// Static platform availability audit (@internal — P4 live scanner is the
// authoritative Ready gate).
export {
  PLATFORM_AUDIT,
  PLATFORM_AUDIT_META,
  componentPlatforms,
  getPlatformAuditEntry,
  isComponentOnBothPlatforms,
} from './platform-audit.ts';
export type {
  PlatformAuditEntry,
  PlatformAuditMeta,
} from './platform-audit.ts';

// Per-category candidate metadata + helpers.
export {
  ALL_CANDIDATE_ENTRIES,
  CANDIDATE_COUNT,
  assertUniqueIds,
  getCandidateById,
  getCandidatesByBatch,
  getCandidatesByCategory,
} from './categories/index.ts';
export {
  analyticsCandidates,
  builderCandidates,
  commerceCandidates,
  contentCandidates,
  dataCandidates,
  excludedCandidates,
  interactionCandidates,
  navigationCandidates,
  pluginCandidates,
  workflowCandidates,
} from './categories/index.ts';

// Pilot candidate analysis (NOT promotion).
export {
  CONDITIONAL_CANDIDATES,
  PILOT_CANDIDATES,
  PILOT_SUMMARY,
  allComponentsOnBothPlatforms,
  validatePilotConsistency,
} from './pilot-candidates.ts';
export type {
  ConditionalCandidate,
  PilotCandidate,
} from './pilot-candidates.ts';
