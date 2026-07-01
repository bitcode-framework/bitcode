# Schema Builder Engine

> ADR-014 Phase 4 — Visual Schema Builder Validation Endpoint

Admin-only endpoint that validates `model.json` definitions before they are
saved by the visual schema builder. Checks structural integrity, field type
validity, relation integrity (including circular relation detection), and
reserved namespace enforcement. Returns a structured
`{valid, errors[], warnings[]}` response consumed by the Mitosis
`BcSchemaBuilder.lite.tsx` component and the Flutter
`bc_schema_builder.dart` widget.

## Architecture

```
POST /api/v1/admin/models/validate  ──► BuilderHandler.ValidateModel
                                              │ 1. auth middleware (JWT)
                                              │ 2. read raw body
                                              │ 3. delegate to ModelValidator.Validate()
                                              ▼
                                    ModelValidator (runtime/builder)
                                      ├─ parse JSON → ModelDefinition
                                      ├─ validateStructure (name, fields present)
                                      ├─ validateFields (type validity, selection options, vector dims)
                                      ├─ validateReservedNamespaces (ctx, input, old, session, self)
                                      ├─ validateRelations (config, targets, circular detection)
                                      └─ return ValidationResult {valid, errors, warnings}
```

## Endpoint

### POST /api/v1/admin/models/validate

Validates a `model.json` definition. Requires JWT authentication (admin-only).

**Request body:** Raw model JSON.

**Response (200):**

```json
{
  "valid": true,
  "errors": [],
  "warnings": []
}
```

**Response with errors:**

```json
{
  "valid": false,
  "errors": [
    {"field": "fields[\"qty\"].type", "message": "invalid field type \"nonexistent\""},
    {"field": "fields[\"ctx\"]", "message": "field name \"ctx\" is a reserved expression namespace"}
  ],
  "warnings": []
}
```

## Validation Rules

### Structure
- `name` is required and must not contain spaces
- `fields` must have at least one entry

### Field Types
- Every field must have a `type` from the supported set (50+ types including
  string, integer, float, boolean, datetime, json, selection, many2one,
  one2many, many2many, morph variants, computed, etc.)
- `selection` and `radio` fields must have `options`
- `vector` fields must specify `dimensions`

### Reserved Namespaces
- Field names cannot use reserved expression namespaces: `ctx`, `input`,
  `old`, `session`, `self` (these conflict with the expression engine)

### Relations
- `many2one`, `many2many`: must specify `model`
- `one2many`: must specify both `model` and `inverse`
- `morph_one`, `morph_many`, `morph_to_many`, `morph_by_many`: must specify
  both `model` and `morph`
- Relational target models are checked against the registry (if available)
- Circular relations are detected: inverse chains that form closed loops on
  the same model are rejected

## Files

| File | Purpose |
|------|---------|
| `engine/internal/runtime/builder/model_validator.go` | ModelValidator service |
| `engine/internal/runtime/builder/model_validator_test.go` | Unit tests (21 cases) |
| `engine/internal/presentation/api/builder_handler.go` | Handler + route registration |
| `engine/internal/presentation/api/builder_handler_test.go` | Handler tests (4 cases) |
| `engine/internal/app.go` | Wiring (NewModelValidator + BuilderHandlerConfig) |
| `engine/embedded/modules/base/i18n/*.json` | i18n keys (11 locales) |

## i18n Keys

All user-facing strings use `schema_builder.*` prefix. Available in 11 locales:
en, id, ar, de, es, fr, ja, ko, pt-BR, ru, zh-CN.
