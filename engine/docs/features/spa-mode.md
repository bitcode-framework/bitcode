# SPA Mode (Render Envelope)

> **ADR-014 scope note:** This feature is part of the "Mitosis compiler, cross-platform
> architecture" track. ADR-014 D3 specifies SSR-as-default with opt-in SPA — this
> feature implements the opt-in path while keeping SSR backwards compatible.

## Overview

SPA Mode adds a `render: "spa"` flag to `view.json` that swaps the engine's
default HTML response (SSR) for a JSON envelope consumed by a client runtime.
The client runtime then renders the view using the existing BcRenderer +
ReactiveStore stack, reusing the portable BlockRenderer work.

**Backwards compatible:** Views that omit `render` (or set `render: "ssr"`)
continue to produce HTML. No existing module, view, or template changes are
required to keep working.

## Where it lives

| Layer | File | Role |
|------|------|------|
| Parser | `engine/internal/compiler/parser/view.go` | `RenderMode` enum, `Render` field, validation |
| Renderer (SPA) | `engine/internal/presentation/view/spa.go` | `SPAEnvelope`, `RenderSPAEnvelope`, per-view data loaders |
| Renderer (SSR) | `engine/internal/presentation/view/renderer.go` | Unchanged; still produces HTML by default |
| HTTP route | `engine/internal/app.go` (`handleViewGet`) | Branches on `view.IsSPAMode(def)` — SPA returns `application/json`, SSR returns HTML |
| Tests | `engine/internal/presentation/view/spa_test.go` | 11 tests: parser, envelope structure, data parity, SSR compat |
| Web runtime | `packages/ui/src/spa-runtime.ts` | Fetch envelope → map view JSON → build `BcRendererProps` + `ReactiveStore` |
| Web tests | `packages/ui/src/spa-runtime.test.ts` | 13 source-level tests (node:test): envelope parsing, mapping, ReactiveStore seeding, fetcher |
| Flutter runtime | `packages/ui-flutter/lib/src/spa_runtime.dart` | Dart parity: `SPAEnvelope`, `SPARuntime`, `viewDefinitionToUiSchema` |
| Flutter tests | `packages/ui-flutter/test/spa_runtime_test.dart` | 28 tests: envelope parsing, mapping, descriptor, URL resolution, errors |

## JSON contract

A view that opts in:

```json
{
  "name": "contact_list",
  "type": "list",
  "render": "spa",
  "title": "Contacts",
  "fields": ["id", "name", "email"]
}
```

The engine responds with a JSON envelope:

```json
{
  "view": { "name": "contact_list", "type": "list", "render": "spa", ... },
  "data": {
    "records": [ { "id": "1", "name": "Alice", "email": "a@x" } ],
    "apiSource": "/api/v1/crm/contacts"
  },
  "meta": {
    "module": "crm",
    "viewPath": "/app/crm/contact_list",
    "viewType": "list",
    "apiSource": "/api/v1/crm/contacts",
    "page": 1,
    "locale": "en",
    "username": "admin"
  }
}
```

`Content-Type: application/json` is set explicitly on the response.

## Why JSON, not HTML

The view.json is already the source of truth for view structure (layout, fields,
actions, etc.). Pushing the *same JSON* to the client lets the existing BcRenderer
re-render the same view deterministically without re-parsing the HTML DOM. This
keeps the SSR and SPA data paths parallel — the engine loads the data once, then
either renders HTML or serializes the envelope. See `TestDataParity_SPAEqualsSSR_FormView`.

## Parser contract

```go
type RenderMode string

const (
    RenderSSR RenderMode = "ssr"
    RenderSPA RenderMode = "spa"
)

type ViewDefinition struct {
    // ...
    Render RenderMode
}
```

Validation in `ParseView`:

| Input | Result |
|------|--------|
| `"render": "spa"` | accepted |
| `"render": "ssr"` | accepted (default) |
| `"render": ""` (absent) | accepted, defaults to SSR |
| `"render": "foo"` | rejected (parser error) |

## Envelope construction

`view.RenderSPAEnvelope(def, opts)` reuses the same data-loading helpers as the
SSR path (`loadFormData`, `loadListData`, `loadKanbanData`, etc.) so the data in
`data` is byte-for-byte the data the SSR path would have inlined into the HTML
template. This is the property `TestDataParity_SPAEqualsSSR_FormView` enforces.

Per-view-type behaviour:

| View type | `data` shape | `meta` extras |
|----------|------------|---------------|
| `list`   | `{ records: [...], apiSource, page, pageSize, totalPages }` | `apiSource` |
| `form`   | `{ record: {...} }` | `recordId` |
| `detail` | `{ record: {...} }` | `recordId` |
| `kanban` | `{ stages: [...], records: [...] }` | `apiSource` |
| `calendar` / `chart` / `custom` / `activity` | pass-through `data` | as-is |

## Route wiring

```go
// engine/internal/app.go
func (a *App) handleViewGet(ctx ...) error {
    def := resolveViewDef(ctx)
    if view.IsSPAMode(def) {
        env, err := view.RenderSPAEnvelope(def, a.renderOptsFor(ctx))
        if err != nil { return err }
        return ctx.JSON(env)        // Content-Type: application/json
    }
    return a.serveViewHTML(ctx)     // existing SSR path
}
```

## Client runtime contract

Both web (`packages/ui/src/spa-runtime.ts`) and Flutter
(`packages/ui-flutter/lib/src/spa_runtime.dart`) expose the same shape:

```ts
class SPARuntime {
  config: { baseUrl: string; headers: Record<string,string> };
  fetcher: (url, headers) => Promise<string>;

  load(viewUrl: string): Promise<SPADescriptor>;
  buildDescriptor(envelope: SPAEnvelope): SPADescriptor;
}

interface SPADescriptor {
  envelope: SPAEnvelope;
  schema: UiSchema;          // engine view JSON → Mitosis/Flutter UiSchema
  values: Record<string, unknown>;  // form initial values
  store: ReactiveStore;      // pre-seeded
  props: BcRendererProps;    // ready to feed BcRenderer
}
```

The mapping function `viewDefinitionToUiSchema(view, data, meta)` is exposed
publicly on both runtimes so hosts can inspect the mapping without a network
fetch. The mapping rule is:

- `type: "list"` → `UiSchemaKind.list` with `columnsDef` from `view.fields` and
  `DataSource.url` from `meta.apiSource` (falls back to `data.apiSource`).
- `type: "form"` → `UiSchemaKind.form` with `FieldNode`s from `view.layout[*].row[*].field`
  and `UiAction`s from `view.actions`. Falls back to flat `view.fields` when no
  layout is present.
- `type: "detail"` → `UiSchemaKind.detail` (form, but `disabled: true`).
- anything else → `UiSchemaKind.custom` with `payload = {view, data, meta}` so the
  host's `customBuilder` can render it (kanban, calendar, chart, etc.).

## Tests

- **Engine (11 tests, `spa_test.go`):**
  - `TestParseView_RenderModeValidation` — parser accepts spa/ssr/empty, rejects
    unknown values.
  - `TestIsSPAMode` — predicate returns true only for `"spa"`.
  - `TestRenderSPAEnvelope_ListView` / `_FormView_DataParity` /
    `_FormView_NewRecord` / `_KanbanView` / `_NilView` / `_KanbanMissingGroupBy` —
    envelope structure for each view type.
  - `TestDataParity_SPAEqualsSSR_FormView` — explicitly asserts the SSR template
    and SPA envelope carry the same record fields.
  - `TestSSR_RenderProducesHTML` / `TestSSR_ExplicitRenderMode_ProducesHTML` —
    backwards compat: views with `render: "ssr"` or no `render` produce HTML.

- **Web (`spa-runtime.test.ts`, 13 tests):** envelope parsing, mapping for
  list/form/detail/custom, descriptor construction with `ReactiveStore` seeding,
  fetcher error propagation, URL resolution (base + relative, no double slash).

- **Flutter (`spa_runtime_test.dart`, 28 tests):** same surface as web, using a
  fake `SPAFetcher` (no real HTTP) and Dart 3 pattern matching for the freezed
  `DataSource` sealed union.

## Cross-platform parity

| Behaviour | Web | Flutter |
|-----------|-----|---------|
| `fetch JSON envelope` | `fetch()` | `HttpClient` (injectable) |
| `parse envelope` | `JSON.parse` | `jsonDecode` |
| `map view JSON → UiSchema` | `viewDefinitionToUiSchema` | `viewDefinitionToUiSchema` |
| `seed ReactiveStore` | `ReactiveStore` ctor | `ReactiveStore(ReactiveStoreOptions(initialValues: ...))` |
| `build BcRendererProps` | `BcRendererProps` ctor | `BcRendererProps` ctor |
| `error: SPA envelope` | `SPARuntimeError` (custom) | `SPARuntimeException` |
| `injectable fetcher` | yes | yes |

## Backwards compatibility

- No existing view.json file is broken: omitting `render` defaults to SSR.
- No existing API surface is renamed or removed.
- No existing test is affected (test count went from 941 → 952 in
  `engine/internal/presentation/view/` after adding the 11 new tests; all old
  tests still pass).
- The new public surface is additive: `SPAEnvelope`, `SPAMeta`, `SPARuntime`,
  `viewDefinitionToUiSchema`, `IsSPAMode`, `RenderSPAEnvelope`. None collide
  with existing exported names.

## i18n

No new user-facing strings were introduced. The error messages raised by
`SPARuntimeError` / `SPARuntimeException` are developer/runtime diagnostics
(network failure, malformed envelope) and are not shown to end users.

## ADR Alignment

- **ADR-014 D3** — SSR default, opt-in SPA. This feature implements the opt-in.
- **ADR-005** — Client rendering goes through the BcRenderer + UiSchema contract
  (no second renderer is introduced; the SPA runtime reuses the existing
  portable block-renderer work).
- **ADR-004 / ADR-011** — `ReactiveStore` is pre-seeded so existing form
  machinery (reactions, computed values, effects) works in SPA mode without
  changes.
