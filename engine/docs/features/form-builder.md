# Form Builder Engine

> ADR-014 Phase 4 — Visual Form Builder Validation Endpoint

Admin-only endpoint that validates `view.json` definitions before they are
saved by the visual form builder. Checks structural integrity, view type
validity, and field reference correctness against the model registry.
Returns a structured `{valid, errors[], warnings[]}` response consumed by
the Mitosis `BcFormBuilder.lite.tsx` component and the Flutter
`bc_form_builder.dart` widget.

## Architecture

```
POST /api/v1/admin/views/validate  ──► BuilderHandler.ValidateView
                                              │ 1. auth middleware (JWT)
                                              │ 2. read raw body
                                              │ 3. delegate to ViewValidator.Validate()
                                              ▼
                                    ViewValidator (runtime/builder)
                                      ├─ parse JSON → ViewDefinition
                                      ├─ validateStructure (name, type, render, template)
                                      ├─ validateDataSources (model vs process mutual exclusion)
                                      ├─ validateModelReference (model exists in registry)
                                      ├─ validateFieldReferences (fields, filters, sort, layout)
                                      └─ return ValidationResult {valid, errors, warnings}
```

## Endpoint

### `POST /api/v1/admin/views/validate`

Admin-only endpoint (requires JWT authentication). Accepts a raw JSON body
containing a `view.json` definition.

**Request body** — a view JSON definition:

```json
{
  "name": "contact_form",
  "type": "form",
  "model": "contact",
  "title": "Contact Form",
  "fields": ["name", "email"],
  "layout": [
    {
      "row": [
        {"field": "name", "width": 6},
        {"field": "email", "width": 6}
      ]
    }
  ]
}
```

**Response** (200 OK):

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
    {
      "field": "fields[\"nonexistent_field\"]",
      "message": "field \"nonexistent_field\" does not exist on model \"contact\""
    }
  ],
  "warnings": []
}
```

**Response with empty body** (400):

```json
{
  "valid": false,
  "errors": [
    {
      "field": "root",
      "message": "request body is required"
    }
  ]
}
```

### Validation Rules

| Rule | Severity | Condition |
|------|----------|-----------|
| `name` required | error | Missing or empty view name |
| `name` no spaces | error | View name contains spaces |
| `type` required | error | Missing view type |
| Valid view type | error | Type must be one of: list, form, kanban, calendar, chart, custom, gantt, map, tree, activity, report, pivot |
| Custom view template | error | Custom views require a `template` field |
| Valid render mode | error | If `render` is set, must be `ssr` or `spa` |
| Model exists | error | Referenced model must exist in model registry |
| Data source mutual exclusion | error | A data source cannot have both `model` and `process` |
| Field references | error | Fields in `fields[]`, `filters[]`, `sort.field`, `date_field`, `start_field`, `end_field`, `parent_field`, `group_by` must exist on the model |
| Layout field references | error | Fields referenced in layout rows, tabs, and child tables must exist on the model |
| Form without layout | warning | Form views without a layout will auto-generate all fields |

## ViewValidator Service

Location: `engine/internal/runtime/builder/view_validator.go`

The `ViewValidator` struct wraps a `ModelResolver` interface (implemented by
`domain/model.Registry`) and performs structural and reference validation
without needing a database connection.

```go
type ViewValidator struct {
    modelRegistry ModelResolver
}

func NewViewValidator(modelRegistry ModelResolver) *ViewValidator

func (v *ViewValidator) Validate(data []byte) *ValidationResult
```

The `ModelResolver` interface:

```go
type ModelResolver interface {
    Get(name string) (*parser.ModelDefinition, error)
}
```

This interface is satisfied by `domain/model.Registry`, allowing the
validator to look up model definitions and verify that referenced fields
actually exist on the model.

## Wiring

The handler is wired in `engine/internal/app.go`:

```go
builderViewValidator := builder.NewViewValidator(a.ModelRegistry)
builderHandler := api.NewBuilderHandler(api.BuilderHandlerConfig{
    ViewValidator: builderViewValidator,
})
api.RegisterBuilderRoutes(a.Fiber, "/api/v1/admin", builderHandler, a.JWTConfig)
```

Routes are registered under `/api/v1/admin` with JWT auth middleware.

## i18n

Builder-related i18n keys are provided in all 11 required locales
(`en`, `id`, `ar`, `de`, `es`, `fr`, `ja`, `ko`, `pt-BR`, `ru`, `zh-CN`)
at `engine/embedded/modules/base/i18n/*.json`.

Keys cover:
- Builder labels: `builder.title`, `builder.canvas`, `builder.property_panel`, etc.
- Field type labels: `builder.type_string`, `builder.type_integer`, `builder.type_many2one`, etc.
- Validation messages: `builder.error_invalid_json`, `builder.error_name_required`, etc.

## Tests

- `engine/internal/runtime/builder/view_validator_test.go` — 16 unit tests
  covering valid views, invalid field refs, missing required fields, invalid
  view types, model not found, invalid JSON, nil registry, layout field refs,
  sort/filter field refs, data source mutual exclusion, render mode validation,
  custom view template requirement, and JSON serialization.

- `engine/internal/presentation/api/builder_handler_test.go` — 5 handler tests
  covering valid view, invalid field ref, missing required fields, empty body,
  and invalid JSON.
