# Pivot Table

> ADR-014 Phase 3 — Feature #42 Pivot Table

OLAP pivot aggregation endpoint that groups model data across row and column
dimensions and computes aggregate measures (sum, count, avg, min, max).
Returns a 2D matrix with row/column headers, cell values, and totals. Consumed
by the Mitosis `BcPivot.lite.tsx` component and the Flutter `bc_pivot.dart`
widget.

## Architecture

```
POST /api/v1/{module}/{model}/pivot  ──► PivotHandler
                                              │ 1. resolve model (registry)
                                              │ 2. check read permission (group-based, 12-perm system)
                                              │ 3. parse pivot request (rows, cols, values, filter)
                                              │ 4. validate aggregations + field existence
                                              │ 5. fetch all matching records (RLS filter applied)
                                              │ 6. sanitize field-level permissions
                                              │ 7. delegate to PivotService.Aggregate()
                                              ▼
                                    PivotService (runtime/pivot)
                                      ├─ Group records by row+col dimensions
                                      ├─ Aggregate value fields (sum/count/avg/min/max)
                                      ├─ Build 2D matrix (rowHeaders × colHeaders)
                                      ├─ Compute row totals, column totals, grand total
                                      └─ Return PivotResult JSON
```

## Endpoint

### `POST /api/v1/{module}/{model}/pivot`

**Request body:**

```json
{
  "rows": ["region"],
  "cols": ["product"],
  "values": [
    { "field": "amount", "agg": "sum" },
    { "field": "quantity", "agg": "count" }
  ],
  "filter": {
    "conditions": [
      { "field": "status", "op": "=", "value": "confirmed" }
    ]
  }
}
```

| Field     | Type             | Required | Description                        |
|-----------|------------------|----------|------------------------------------|
| `rows`    | `[]string`       | Yes      | Row dimension fields               |
| `cols`    | `[]string`       | No       | Column dimension fields            |
| `values`  | `[]ValueSpec`    | Yes      | Value measures with aggregation    |
| `filter`  | `object`         | No       | Query filter (same format as list) |

**ValueSpec:**

| Field  | Type   | Required | Description |
|--------|--------|----------|-------------|
| `field`| string | Yes      | Field name to aggregate |
| `agg`  | string | Yes      | Aggregation: `sum`, `count`, `avg`, `min`, `max` |

**Response (200):**

```json
{
  "rowHeaders": [["North"], ["South"]],
  "colHeaders": [["Widget"], ["Gadget"]],
  "cells": [[[100], [200]], [[150], [300]]],
  "rowTotals": [[300], [450]],
  "colTotals": [[250], [500]],
  "grandTotal": [750],
  "valueFields": [
    { "field": "amount", "agg": "sum" }
  ]
}
```

**Error responses:**

| Status | Condition |
|--------|-----------|
| 400    | Missing `rows` or `values` |
| 400    | Invalid aggregation function |
| 401    | Not authenticated |
| 403    | No read permission on model |

## PivotService (`internal/runtime/pivot/service.go`)

Pure Go OLAP engine with no database dependency. Works across SQL and MongoDB
backends because it operates on pre-fetched `[]map[string]interface{}` records.

### Key methods

| Method | Description |
|--------|-------------|
| `Validate(req *Request, fields map[string]bool) error` | Validates request fields exist and aggregations are supported |
| `Aggregate(req *Request, records []map[string]interface{}) (*Result, error)` | Computes the full pivot result |

### Aggregation functions

| Function | Description |
|----------|-------------|
| `sum`    | Sum of numeric values |
| `count`  | Count of non-null values |
| `avg`    | Average of numeric values |
| `min`    | Minimum numeric value |
| `max`    | Maximum numeric value |

### Null handling

- Records missing a row or col dimension value are grouped under `"null"` bucket
- Null/missing value fields are excluded from aggregation (count counts non-null only)
- `avg` divides by count of non-null values

### Totals

- **Row totals**: Sum across columns for each row
- **Column totals**: Sum across rows for each column
- **Grand total**: Sum across all rows and columns

## Permission Integration

The `PivotHandler` enforces:

1. **Model read permission** — user's groups must have `read` access on the model
2. **Field-level permission** — fields restricted by `groups` are excluded from rows/cols/values
3. **Record-level security (RLS)** — `RecordRuleMiddleware` filter is applied to the data fetch

## Component Wiring

### View renderer

`ViewPivot` in `compiler/parser/view.go` — view type `"pivot"`.

### Component compiler

`CompilePivot` in `presentation/view/component_compiler.go` — emits
`<bc-pivot>` custom element with data-source attributes.

### BcPivot component (Mitosis)

`packages/ui/src/components/data/pivot/BcPivot.lite.tsx` — renders the pivot
matrix with expandable rows, totals, grand total, and no-data state.

### BcPivot component (Flutter)

`packages/ui-flutter/lib/src/components/bc_pivot.dart` — Material `DataTable`
on Android, Cupertino-styled container on iOS.

### Block Catalog

`analytics.pivot` entry in
`packages/ui-schema/src/block-catalog/categories/analytics.ts`.

## i18n Keys

| Key | English |
|-----|---------|
| `pivot.total` | Total |
| `pivot.grandTotal` | Grand Total |
| `pivot.noData` | No data available |

Available in all 11 locales: en, id, ar, de, es, fr, ja, ko, pt-BR, ru, zh-CN.

## Tests

### Engine tests (`internal/runtime/pivot/service_test.go`)

20 tests covering: single aggregation, multiple aggregations, multi-dimension
rows, multi-dimension columns, null handling (null row/col/value), totals
correctness, avg correctness, min/max, request validation.

### Handler tests (`internal/presentation/api/pivot_handler_test.go`)

9 tests covering: single agg, multi agg, missing rows, missing values, invalid
agg, permission denied, field-level permission, auth required, no cols.
