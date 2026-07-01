# Activity Feed

> ADR-014 Phase 1 — Feature #6 Activity Feed (BcActivity)

Aggregated activity feed endpoint that transforms `audit_log` entries into
`BcActivityEntry[]` items consumable by the `bc-activity` Stencil component,
the Mitosis `BcActivity.lite.tsx` component, and the Flutter `bc_activity.dart`
widget.

## Architecture

```
audit_logs (table) ──► AuditLogRepository.FindActivity ──► activity.Service
                                                              │ transform (verb/variant/icon/actor)
                                                              ▼
                                              ActivityEntry[] (icon, verb, variant, actor, timestamp)
                                                              │
GET /api/v1/activity ◄──────────────────────────────── ActivityHandler
                                                              │
                                              permission + record-rule check (record-scoped)
```

The activity feed is **read-only** — it aggregates the existing `audit_log`
infrastructure table (written by the audit middleware + `AuditLogRepository`).
No new persistence model is required; the service layer transforms raw audit
rows into the UI-facing `BcActivityEntry` shape.

## ActivityService (`internal/runtime/activity/service.go`)

The `Service` queries `audit_logs` via `AuditLogRepository.FindActivity` and
transforms each row into an `ActivityEntry`:

| Field      | Source                              | Notes |
|------------|-------------------------------------|-------|
| `id`       | `audit_logs.id`                     | UUID |
| `actor`    | resolved from `res_user`            | name → username → id fallback; `"system"` for null user |
| `actor_id` | `audit_logs.user_id`                | raw UUID |
| `action`   | `audit_logs.action`                 | create / update / delete / login / logout / register / request |
| `verb`     | derived from action                 | "created", "updated", "deleted", "logged in", … |
| `variant`  | derived from action                 | success / default / destructive / info |
| `icon`     | derived from variant                | "+", "×", "!", "i", "•" |
| `target`   | model + record_id OR request_path   | "contact #rec-1" or "/api/v1/contacts" |
| `model`    | `audit_logs.model_name`             | optional |
| `record_id`| `audit_logs.record_id`              | optional |
| `timestamp`| `audit_logs.created_at` (RFC3339)   | UTC |
| `ip_address`| `audit_logs.ip_address`            | optional |
| `changes`  | `audit_logs.changes`                | optional JSON |

### Actor resolution

The service batch-resolves actor display names in a single query against the
user table (`res_user` by default; configurable via `SetUserTable`). Resolution
chain: `name` → `username` → `id` → `"system"` (for null user_id).

### Verb / variant mapping

| Action     | Verb          | Variant      |
|------------|---------------|--------------|
| `create`   | created       | success      |
| `update`   | updated       | default      |
| `delete`   | deleted       | destructive  |
| `login`    | logged in     | info         |
| `logout`   | logged out    | info         |
| `register` | registered    | success      |
| `request`  | requested     | default      |
| other      | <action>      | default      |

Custom mapping is supported via `SetVerbResolver` and `SetActorRenderer`.

## API Endpoint

### `GET /api/v1/activity`

Returns aggregated activity entries (newest-first).

**Query parameters** (all optional):

| Param       | Type      | Description |
|-------------|-----------|-------------|
| `model`     | string    | Filter by model name (e.g. `contact`) |
| `record_id` | string    | Filter by record ID (requires `model`) |
| `user`      | string    | Filter by actor user ID |
| `since`     | ISO8601   | Only entries after this timestamp |
| `limit`     | int       | Page size (default 50, max 200) |
| `page`      | int       | Page number (1-based, default 1) |

**Response:**

```json
{
  "data": [
    {
      "id": "uuid",
      "actor": "Alice Wong",
      "actor_id": "u1",
      "action": "create",
      "verb": "created",
      "target": "contact #rec-1",
      "model": "contact",
      "record_id": "rec-1",
      "timestamp": "2026-07-01T10:00:00Z",
      "variant": "success",
      "icon": "+",
      "ip_address": "127.0.0.1"
    }
  ],
  "total": 42,
  "page": 1,
  "page_size": 50
}
```

### Permissions

- **Authentication required** (401 if no `user_id`).
- When `model` + `record_id` are supplied, the handler enforces:
  - Group-based **read permission** on the target model (403 if denied).
  - **Record rule** access check on the specific record (403 if denied).
- Superusers bypass permission checks.
- Global timeline queries (no model/record_id) are available to any
  authenticated user.

## View Renderer

`ViewType "activity"` is wired in `presentation/view/renderer.go`. A standalone
activity view renders a `<bc-activity>` Stencil component bound to the activity
API:

```go
case parser.ViewActivity:
    contentHTML, err = r.renderActivity(ctx, viewDef, opt)
```

The `ComponentCompiler.CompileActivity` method emits:

```html
<bc-activity data-model="contact" data-record-id="rec-1"
             data-source="/api/v1/activity?model=contact&amp;record_id=rec-1"
             data-variant="timeline"></bc-activity>
```

The `activity: true` flag on a `LayoutItem` (form views) also emits a
`<bc-activity>` part scoped to the current record.

## UI Components

### Web (`packages/ui/src/components/social/activity/`)

- `BcActivity.contract.ts` — ADR-008 contract: `entries`, `emptyText`,
  `loading`, `variant` (timeline/feed/compact), `dataSource` (endpoint/model/
  recordId/params), `onLoadMore`, `onItemAction`.
- `BcActivity.lite.tsx` — Mitosis Pure component. Declarative data-source fetch
  on mount, skeleton loading state, 5 semantic variants, a11y (role=region,
  aria-busy, aria-hidden decorative icons), token-driven palette.
- `BcActivity.test.ts` — 15 source tests (contract surface, a11y, variants,
  data-source wiring, load-more, item action).
- `BcActivity.rtl.test.tsx` — vitest a11y smoke test.

### Flutter (`packages/ui-flutter/lib/src/components/bc_activity.dart`)

- Parity with web contract: `items`, `variant` (timeline/feed/compact),
  `loading`, `dataSource`, `onLoadMore`, `onItemAction`, `onItemClick`.
- Platform-aware: `PlatformWidget` (Material `InkWell` / Cupertino
  `CupertinoButton`).
- Adaptive theming via `BcTheme` tokens.
- 14 widget tests (`test/bc_activity_test.dart`).

## i18n Keys

All 11 locales (`engine/embedded/modules/base/i18n/*.json`):

| Key | Description |
|-----|-------------|
| `activity.title` | Activity feed title |
| `activity.empty` | Empty state text |
| `activity.load_more` | Load more button label |
| `activity.forbidden_read` | Permission denied message |
| `activity.authentication_required` | Auth required message |

## Tests

- **Service** (`runtime/activity/service_test.go`): 12 tests — action mapping,
  batch user resolution, filters (model, record, user, since), pagination,
  limit clamping, empty result, global timeline, target building, system actor,
  verb resolver coverage.
- **Handler** (`presentation/api/activity_handler_test.go`): 10 tests — happy
  path, global timeline, user filter, since filter, invalid since, pagination,
  empty result, permission denied, unauthenticated, limit clamping.
- **View** (`presentation/view/component_compiler_test.go`): activity part on
  form + standalone activity view compilation.

## Key Files

| File | Purpose |
|------|---------|
| `engine/internal/runtime/activity/service.go` | ActivityService (transform + aggregation) |
| `engine/internal/runtime/activity/service_test.go` | Service tests |
| `engine/internal/presentation/api/activity_handler.go` | HTTP handler + route registration |
| `engine/internal/presentation/api/activity_handler_test.go` | Handler tests |
| `engine/internal/presentation/view/renderer.go` | `renderActivity` (ViewActivity case) |
| `engine/internal/presentation/view/component_compiler.go` | `CompileActivity` (bc-activity emit) |
| `engine/internal/infrastructure/persistence/audit_log.go` | `FindActivity` query |
| `engine/internal/app.go` | Wiring (service + handler + routes) |
| `packages/ui/src/components/social/activity/` | Web component + contract + tests |
| `packages/ui-flutter/lib/src/components/bc_activity.dart` | Flutter widget |
| `packages/ui-flutter/test/bc_activity_test.dart` | Flutter widget tests |
| `engine/embedded/modules/base/i18n/*.json` | 11 locale files (activity.* keys) |
