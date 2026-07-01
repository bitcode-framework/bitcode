# Workflow Builder Engine

> ADR-014 Phase 4 — Visual Workflow Builder Validation Endpoint

Admin-only endpoint that validates `process.json` definitions before they are
saved by the visual workflow builder. Checks structural integrity, step type
validity, step configuration completeness, transition target references, and
cross-process call references. Supports both sequential (linear `steps[]`)
and DAG (`nodes` + `edges`) process formats. Returns a structured
`{valid, errors[], warnings[]}` response consumed by the Mitosis
`BcWorkflowBuilder.lite.tsx` component and the Flutter
`bc_workflow_builder.dart` widget.

## Architecture

```
POST /api/v1/admin/processes/validate  ──► BuilderHandler.ValidateProcess
                                                  │ 1. auth middleware (JWT)
                                                  │ 2. read raw body
                                                  │ 3. delegate to ProcessValidator.Validate()
                                                  ▼
                                        ProcessValidator (runtime/builder)
                                          ├─ parse JSON → ProcessDefinition
                                          ├─ validateStructure (name, steps/nodes present)
                                          ├─ if DAG:
                                          │    ├─ validateDAG: node types, edge refs
                                          │    └─ validateDAGReachability (orphan warnings)
                                          ├─ if sequential:
                                          │    ├─ validateStepType (type exists)
                                          │    ├─ validateStepConfig (model, condition, etc.)
                                          │    ├─ validateStepTransitions (then/else/case refs)
                                          │    └─ validateCallReference (process registered)
                                          └─ return ValidationResult {valid, errors, warnings}
```

## Endpoint

### `POST /api/v1/admin/processes/validate`

Admin-only endpoint (requires JWT authentication). Accepts a raw JSON body
containing a `process.json` definition.

**Request body** — a sequential process JSON:

```json
{
  "name": "create_lead",
  "steps": [
    {"name": "validate_input", "type": "validate", "rules": {}},
    {"name": "create_record", "type": "create", "model": "lead"},
    {"name": "emit_created", "type": "emit", "event": "lead_created"}
  ]
}
```

**Request body** — a DAG process JSON:

```json
{
  "name": "approval_flow",
  "nodes": {
    "start":    {"type": "query", "model": "invoice"},
    "validate": {"type": "validate", "rules": {}},
    "approve":  {"type": "update", "model": "invoice"},
    "notify":   {"type": "emit", "event": "approved"}
  },
  "edges": [
    {"from": "start", "to": "validate"},
    {"from": "validate", "to": "approve"},
    {"from": "approve", "to": "notify"}
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

**Response with errors (invalid step ref):**

```json
{
  "valid": false,
  "errors": [
    {
      "field": "steps[\"check\"].then",
      "message": "transition target \"nonexistent_step\" does not match any step name"
    }
  ],
  "warnings": []
}
```

**Response with warnings (orphan DAG node):**

```json
{
  "valid": true,
  "errors": [],
  "warnings": [
    {
      "field": "nodes[\"orphan\"]",
      "message": "node \"orphan\" is isolated: no edges target or originate from it"
    }
  ]
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
| `name` required | error | Missing or empty process name |
| `name` no spaces | error | Process name contains spaces |
| At least one step | error | Sequential processes need `steps[]` |
| At least one node | error | DAG processes need `nodes` map |
| Valid step type | error | Step type must be one of: validate, query, create, update, delete, if, switch, loop, emit, call, script, http, assign, log, upsert, count, sum |
| Step model required | error | query/create/update/delete/upsert/count/sum steps require `model` |
| If condition required | error | `if` steps require `condition` |
| Switch field required | error | `switch` steps require `field` |
| Loop over required | error | `loop` steps require `over` |
| Emit event required | error | `emit` steps require `event` |
| HTTP url required | error | `http` steps require `url` |
| Assign variable required | error | `assign` steps require `variable` |
| Transition target exists | error | `then`, `else`, `cases` targets must match a step name |
| Call process registered | error | `call` step process must exist in registry |
| DAG edge from/to exist | error | Edge endpoints must reference existing nodes |
| DAG edge completeness | error | Edges must have both `from` and `to` |
| DAG node type valid | error | Node type must be a valid step type |
| Isolated DAG node | warning | Node with no incoming or outgoing edges |

## ProcessValidator Service

Location: `engine/internal/runtime/builder/process_validator.go`

The `ProcessValidator` struct wraps a `ProcessResolver` interface (implemented
by `runtime/executor.ProcessRegistry`) and performs structural and reference
validation without needing a database connection.

```go
type ProcessValidator struct {
    processRegistry ProcessResolver
}

func NewProcessValidator(processRegistry ProcessResolver) *ProcessValidator

func (v *ProcessValidator) Validate(data []byte) *ValidationResult
```

The `ProcessResolver` interface:

```go
type ProcessResolver interface {
    LoadProcess(name string) (*parser.ProcessDefinition, error)
}
```

This interface is satisfied by `executor.ProcessRegistry`, allowing the
validator to look up process definitions and verify that `call` step
references point to registered processes.

## Wiring

The handler is wired in `engine/internal/app.go`:

```go
builderViewValidator := builder.NewViewValidator(a.ModelRegistry)
builderProcessValidator := builder.NewProcessValidator(a.ProcessRegistry)
builderHandler := api.NewBuilderHandler(api.BuilderHandlerConfig{
    ViewValidator:    builderViewValidator,
    ProcessValidator: builderProcessValidator,
})
api.RegisterBuilderRoutes(a.Fiber, "/api/v1/admin", builderHandler, a.JWTConfig)
```

Routes are registered under `/api/v1/admin` with JWT auth middleware.

## i18n

Workflow builder-related i18n keys are provided in all 11 required locales
(`en`, `id`, `ar`, `de`, `es`, `fr`, `ja`, `ko`, `pt-BR`, `ru`, `zh-CN`)
at `engine/embedded/modules/base/i18n/*.json`.

Keys cover:
- Builder labels: `workflow_builder.title`, `workflow_builder.canvas`, `workflow_builder.property_panel`, etc.
- Mode labels: `workflow_builder.dag_mode`, `workflow_builder.sequential_mode`
- Validation messages: `workflow_builder.error_invalid_json`, `workflow_builder.error_name_required`, etc.

## Tests

- `engine/internal/runtime/builder/process_validator_test.go` — 23 unit tests
  covering valid sequential/DAG processes, missing name, name with spaces, no
  steps, invalid/missing step types, step config validation (model, condition,
  event, url), transition target validation (then/else/case refs), DAG edge
  ref validation, DAG orphan node warnings, call step cross-process ref
  validation, invalid JSON, and nil registry.

- `engine/internal/presentation/api/builder_handler_test.go` — 4 handler tests
  covering valid process, invalid step ref, orphan transition, and empty body.
