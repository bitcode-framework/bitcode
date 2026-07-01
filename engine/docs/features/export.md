# Export Data

> ADR-014 Phase 1 — Feature #7 Export Data (BcExport wire)

Model-scoped data export endpoint that streams records to CSV, XLSX, or JSON
using the existing query builder (`persistence.Query`) for filter/sort,
column sanitization for field-level permissions, and mask enforcement for
sensitive fields. Consumed by the `bc-export` Stencil component, the Mitosis
`BcExport.lite.tsx` component, and the Flutter `bc_export.dart` widget.

## Architecture

```
POST /api/v1/{module}/{model}/export  ──► ExportHandler
                                              │ 1. resolve model (registry)
                                              │ 2. check export permission (group-based, 12-perm system)
                                              │ 3. resolve field-level permission (groups + can_mask)
                                              │ 4. sanitize columns (drop group-restricted fields)
                                              │ 5. build query (filter + sort + RLS)
                                              │ 6. stream rows in batches (1000/batch)
                                              │      └─ apply field mask per batch
                                              ▼
                                    ExportService (runtime/export)
                                      ├─ GenerateCSV  (encoding/csv, UTF-8 BOM)
                                      ├─ GenerateJSON (json.Encoder, streaming)
                                      └─ GenerateXLSX (excelize streaming writer)
```

The export is **transport-agnostic**: the `ExportService` formats records
supplied by a `FetchFunc`, while the `ExportHandler` wires the HTTP layer
(permission, query building, streaming response). This separation lets the
service be reused by the CLI and tests without an HTTP server.

## ExportService (`internal/runtime/export/service.go`)

The `Service` renders model data to three formats via batched fetching:

| Format | Generator           | Content-Type                                                        | Streaming |
|--------|---------------------|---------------------------------------------------------------------|-----------|
| CSV    | `GenerateCSV`       | `text/csv` (UTF-8 BOM for Excel)                                    | Yes       |
| XLSX   | `GenerateXLSX`      | `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` | Buffer*   |
| JSON   | `GenerateJSON`      | `application/json`                                                  | Yes       |

\* XLSX uses excelize's streaming writer (`NewStreamWriter`) but the final
ZIP must be serialized to a buffer before sending — excelize cannot stream
directly to an HTTP response writer. CSV and JSON stream row-by-row via
fasthttp's `SetBodyStreamWriter` so memory stays bounded for large datasets.

### Column resolution

- **Explicit columns**: caller supplies `[{field, label}]` — labels default
  to the model field label, then the field name.
- **Default columns**: `DefaultColumnsFromModel` derives columns from the
  model definition, excluding relational collections (one2many, many2many,
  morph_one, morph_many, morph_to_many, morph_by_many) since they represent
  collections rather than scalar export values. many2one fields are included
  (they store a foreign-key scalar).

### Cell formatting

`FormatCellValue` renders values to strings for CSV; `normalizeXLSXValue`
preserves native types (string, int, float, bool, time) for XLSX so
spreadsheet apps treat numbers/dates correctly.

## API Endpoint

### `POST /api/v1/{module}/{model}/export`

Streams the exported file as an attachment.

**Request body** (JSON):

```json
{
  "format": "csv",
  "columns": [{"field": "name", "label": "Full Name"}],
  "filter": {"wheres": [{"field": "status", "op": "=", "value": "active"}]},
  "sort": {"field": "name", "order": "asc"}
}
```

| Field     | Type     | Required | Description |
|-----------|----------|----------|-------------|
| `format`  | string   | no (default `csv`) | `csv` \| `xlsx` \| `json` |
| `columns` | array   | no (defaults from model) | `[{field, label?}]` — sanitized per field-level permission |
| `filter`  | object  | no | Query builder DSL (`persistence.Query` map form) |
| `sort`    | object  | no | `{field, order}` where order is `asc` \| `desc` |

**Response**: streaming file with `Content-Type` + `Content-Disposition:
attachment; filename="<model>_<timestamp>.<ext>"`.

### `GET /api/v1/{module}/{model}/export`

Browser-initiated download variant. Query params: `format`, `columns`
(comma-separated field list), `oql` (filter), `sort_field`, `sort_order`.

## Permissions

Export enforces the **12-permission group-based system** at two levels:

### Model-level (export permission)

- **Authentication required** (401 if no `user_id`).
- The `export` permission is checked via `PermissionService.GetModelPermissions`.
  Users without `can_export` on the target model receive **403 Forbidden**.
- Superusers bypass permission checks.
- Record-level security (RLS) filters from middleware (`rls_filters`,
  `rls_expr_filters`) are merged into the query so only permitted records
  are exported.

### Field-level (column sanitization + mask)

- **Column sanitization**: fields declaring field-level `groups` are dropped
  from the export when the user belongs to none of those groups. If all
  requested columns are dropped, the handler returns **403**.
- **Field masking**: fields flagged `mask: true` are masked (last
  `mask_length` chars visible, rest `*`) when the user lacks the `can_mask`
  permission. Masking is applied per-batch inside the fetcher wrapper before
  any value reaches the format generator.
- Superusers see all fields unmasked (via the superuser bypass).

This mirrors the `field_filter.go` helpers used by the CRUD handler so
export respects the same field-level rules as the regular API responses.

## View Renderer

The `<bc-export>` Stencil component is emitted by the component compiler
inside list-view toolbars when `viewDef.Exportable` is true:

```go
exportAttrs := fmt.Sprintf(`model="%s" formats='["csv","xlsx","json"]'`, esc(viewDef.Model))
if opts.ModuleName != "" {
    exportAttrs += fmt.Sprintf(` endpoint="/api/v1/%s/%s/export"`, esc(opts.ModuleName), esc(viewDef.Model))
}
// → <div class="bc-list-toolbar"><bc-export model="contact" formats='["csv","xlsx","json"]' endpoint="/api/v1/crm/contact/export"></bc-export></div><bc-datatable ...></bc-datatable>
```

This is the SSR-side `part:"export"` wiring — the export button appears in
list views and POSTs to the model's export endpoint with the configured
format.

## UI Components

### Web (`packages/ui/src/components/print/export/`)

- `BcExport.contract.ts` — ADR-008 contract: `formats`, `disabled`, `label`,
  `model`, `endpoint`, `columns` (permission-gated `BcExportColumn` with
  `permitted?`), `filter`, `sort`, `loading`, `progress`, `filename`,
  `onExport`, `onExportComplete`, `onExportError`, `onProgress`.
- `BcExport.lite.tsx` — Mitosis Pure component. Declarative export via
  `endpoint` (POST → blob download with progress), fallback to `onExport`
  callback when no endpoint. Progress indicator (`aria-busy`), token-driven
  palette, a11y (role=group, per-button aria-label).
- `BcExport.test.ts` — 13 source tests (contract surface, a11y, formats,
  disabled gate, progress, declarative fetch, columns, permission-gated
  column).
- `BcExport.rtl.test.tsx` — vitest a11y smoke test (axe-core zero violations).

### Flutter (`packages/ui-flutter/lib/src/components/bc_export.dart`)

- Parity with web contract: `formats`, `disabled`, `label`, `endpoint`,
  `columns` (`BcExportColumn`), `loading`, `progress` (`BcExportProgress`),
  `onExport`, `onExportComplete`, `onProgress`.
- Platform-aware: Material popup overlay / Cupertino action sheet (iOS/macOS).
- Loading indicator (`CircularProgressIndicator`) on the trigger when
  `loading` is true.
- Adaptive theming via `BcTheme` tokens.
- 8 widget tests (`test/bc_export_test.dart`).

## i18n Keys

All 11 locales (`engine/embedded/modules/base/i18n/*.json`):

| Key | Description |
|-----|-------------|
| `export.title` | Export button label |
| `export.empty` | No data to export |
| `export.no_columns` | No exportable columns permitted |
| `export.permission_denied` | Export permission denied message |
| `export.authentication_required` | Auth required message |
| `export.unsupported_format` | Unsupported format error |
| `export.model_required` | Model parameter required error |
| `export.exporting` | Export in progress label |
| `export.complete` | Export complete label |
| `export.failed` | Export failed label |

## Tests

- **Service** (`runtime/export/service_test.go`): 21 tests — CSV (BOM,
  quoting, empty, multi-page, nil values), XLSX (header, empty, content),
  JSON (keys, empty, valid array), filename + timestamp, column derivation
  from model, `QueryFetcher` adapter, cell formatting, `Export` service
  orchestration (CSV/JSON/XLSX, unsupported format).
- **Handler** (`presentation/api/export_handler_test.go`): 12 tests — CSV,
  XLSX, JSON output + content-types, permission denial (403), column
  selection, filter, default columns, GET variant, unsupported format,
  **field mask** (non-mask user sees masked phone), **column sanitization**
  (group-restricted salary dropped), **superuser mask bypass**.

## Key Files

| File | Purpose |
|------|---------|
| `engine/internal/runtime/export/service.go` | ExportService (CSV/XLSX/JSON generators, batched fetch) |
| `engine/internal/runtime/export/service_test.go` | Service tests (21) |
| `engine/internal/presentation/api/export_handler.go` | HTTP handler + permission + field-level sanitization/mask |
| `engine/internal/presentation/api/export_handler_test.go` | Handler tests (12) |
| `engine/internal/presentation/api/field_filter.go` | Shared field-level filter + mask helpers |
| `engine/internal/presentation/view/component_compiler.go` | `<bc-export>` emit in list toolbars |
| `engine/internal/app.go` | Wiring (handler + routes per module/model) |
| `packages/ui/src/components/print/export/BcExport.contract.ts` | Web contract (ADR-008) |
| `packages/ui/src/components/print/export/BcExport.lite.tsx` | Web component (declarative export + progress) |
| `packages/ui/src/components/print/export/BcExport.test.ts` | Web source tests (13) |
| `packages/ui-flutter/lib/src/components/bc_export.dart` | Flutter widget (adaptive) |
| `packages/ui-flutter/test/bc_export_test.dart` | Flutter widget tests (8) |
| `engine/embedded/modules/base/i18n/*.json` | 11 locale files (export.* keys) |
