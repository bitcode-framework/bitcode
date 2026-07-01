# Wizard (Multi-Step Form)

> **ADR-014 scope note:** The wizard block catalog entry (`interaction.multi-step-wizard`)
> defines the schema convention, presets, and Ready gates. Engine-backed wiring
> (server validation endpoint auto-generation, finish action dispatch via
> ActionSource) remains ADR-014 scope.

## Overview

The wizard component (`BcWizard`) provides a multi-step form experience with:

- **Per-step icon** — rendered in the step indicator
- **Schema-driven content** — each step can declare a `schema` (UiSchema fragment)
  rendered via the `renderStepContent` callback (ADR-005 §2 compositor pattern)
- **Per-step validation** — `validationMode: 'client' | 'server' | 'both'`
  - Client: ADR-004 ValidationRule (min, max, length, pattern, custom)
  - Server: POST `{ step, values }` → `{ ok, errors }` to a consumer-set
    `validationEndpoint` (standalone-compatible, no BitCode coupling)
- **Skippable steps** — `skippable: true` bypasses validation for that step
- **Merged values** — `onFinish(values)` receives the merged values from all steps
- **Controlled/uncontrolled** — `values` + `onValuesChange` for controlled mode,
  or internal state for uncontrolled mode
- **Legacy mode** — steps without `schema` fall back to `<Slot/>` (backward compat)

## Schema Convention

```typescript
// packages/ui-schema/src/schema/types.ts
type WizardValidationMode = 'client' | 'server' | 'both';

interface WizardStepSchema {
  title: string;
  description?: string;
  icon?: string;
  schema?: UiSchema;           // Form fragment for this step
  validationMode?: WizardValidationMode;  // Default: 'client'
  validationEndpoint?: string;  // Generic endpoint (consumer-set)
  skippable?: boolean;
  blocking?: boolean;
}

// Added to UiSchema:
// steps?: WizardStepSchema[]
```

## Contract

The web contract lives in `packages/ui/src/components/overlays/wizard/BcWizard.contract.ts`.
The Flutter contract is in `packages/ui-flutter/lib/src/components/bc_wizard.dart`.

### Key Props

| Prop | Type | Description |
|------|------|-------------|
| `steps` | `WizardStep[]` | Step definitions with optional schema/icon/validation |
| `step` / `defaultStep` | `number` | Controlled / uncontrolled current step |
| `canAdvance` | `(step: number) => boolean` | Gate before next/finish (legacy, preserved) |
| `onFinish` | `(values: Record<string, unknown>) => void` | Called with merged values |
| `finishAction` | `ActionSource?` | Optional declarative finish action (ADR-018) |
| `values` / `onValuesChange` | `Record<string, unknown>` / `(values) => void` | Controlled mode |
| `renderStepContent` | `(step, values, onChange) => ReactNode` | Schema-driven content renderer |

## Cross-Platform Parity

| Feature | Web (Mitosis) | Flutter |
|---------|---------------|---------|
| Per-step icon | ✅ `data-bc-wizard-icon` | ✅ `Icons` mapping (~45 icons) |
| Schema-driven content | ✅ `renderStepContent` callback | ✅ `stepContentBuilder(index, step, values)` |
| Client validation | ✅ ADR-004 ValidationRule | ✅ Mirror implementation |
| Server validation | ✅ `fetch()` POST | ✅ `http.post()` |
| Merged values | ✅ `onFinish(values)` | ✅ `onFinish(values, meta)` |
| Skippable | ✅ | ✅ |
| Legacy Slot mode | ✅ `<Slot/>` | ✅ Default title+description |

## ADR Alignment

- **ADR-004** — ValidationRule (client-side field validation)
- **ADR-005 §2** — Compositor pattern (`renderStepContent` callback, no direct BcRenderer import)
- **ADR-005 §3** — Generic server validation endpoint (standalone-compatible)
- **ADR-012** — Canonical types in `ui-schema` package
- **ADR-018** — Optional `finishAction: ActionSource` for declarative mode
