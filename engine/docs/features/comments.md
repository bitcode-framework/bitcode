# Comments & Mentions

Record-level commenting system with `@mention` parsing, threaded replies, real-time broadcast, and permission-aware access control.

## Overview

Comments are **polymorphic** — they attach to any record of any model via `commentable_type` (model name) and `commentable_id` (record UUID). This means a CRM contact, a sales order, and an HR leave request can all have independent comment threads without per-model comment tables.

### Key Features

| Feature | Description |
|---------|-------------|
| `@mention` parsing | `@username` syntax in comment body auto-parsed and resolved to user IDs |
| Threaded replies | Comments support `parent_id` for nested reply threads |
| Real-time broadcast | New comments + mentions pushed via WebSocket (`comment:{type}:{id}` channel) |
| Mention notifications | Mentioned users get a row in the `notifications` table + `notification:user:{id}` WS event |
| Permission-aware | Create/List/Delete checked against group-based permissions + record rules |
| Polymorphic | Any model can receive comments — no per-model setup needed |
| Soft delete | Comments are deactivated (`active = false`), not hard-deleted |
| Author-scoped delete | Only the comment author (or users with delete permission) can delete |

## Architecture

```
internal/runtime/comment/
├── service.go          Service: Create, ListByRecord, Delete, ResolveMentions, ParseMentions
├── types.go            CommentRecord, Mention, CreateCommentInput, ServiceConfig
└── service_test.go     15 tests (mention parsing, resolution, create+notify, list)

internal/presentation/api/
├── comment_handler.go          HTTP handler: Create, List (threaded), Delete + checkPermission
├── comment_handler_test.go     5 handler tests (create+list, threading, validation, permission denial)
└── routes.go                   RegisterCommentRoutes mounts /api/v1/:module/comments
```

### Data Flow

```
Client POST /api/v1/crm/comments  {body: "Hello @bob", record_id: "rec-1"}
  │
  ▼
CommentHandler.Create
  ├─ Parse body + record_id from request
  ├─ Resolve morph type (model → commentable_type)
  ├─ checkPermission(userID, model, recordID, "write")
  │    ├─ PermissionService: group-based can_write on model
  │    └─ RecordRuleService: domain filter on record
  ├─ service.Create(input)
  │    ├─ ParseMentions(body) → [{username:"bob", ...}]
  │    ├─ ResolveMentions(mentions) → lookup user table, set UserID + Resolved
  │    ├─ INSERT INTO comments (...)
  │    ├─ For each resolved mention: INSERT INTO notifications (...)
  │    └─ broadcaster.Broadcast("comment:crm:rec-1", {type:"comment.created", data})
  └─ Return 201 {data: comment}

Client GET /api/v1/crm/comments?record_id=rec-1
  │
  ▼
CommentHandler.List
  ├─ checkPermission(userID, model, recordID, "read")
  ├─ service.ListByRecord(type, id, page, pageSize) → flat list
  ├─ buildThread(records) → nested tree (parent → replies[])
  └─ Return 200 {data: [threaded_comments], total, page}
```

## API Endpoints

All routes are under `/api/v1/:module/comments` and require JWT authentication.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/:module/comments?record_id={id}` | List comments for a record (threaded) |
| `POST` | `/api/v1/:module/comments` | Create a comment or reply |
| `DELETE` | `/api/v1/:module/comments/:id` | Soft-delete a comment (author or delete-perm only) |

### Create Comment

```json
POST /api/v1/crm/comments
{
  "body": "Hey @alice, can you review this?",
  "record_id": "rec-1",
  "parent_id": "optional-parent-uuid"
}
```

Response `201 Created`:
```json
{
  "data": {
    "id": "uuid",
    "body": "Hey @alice, can you review this?",
    "commentable_type": "crm",
    "commentable_id": "rec-1",
    "parent_id": null,
    "author_id": "user-1",
    "mentions": [{"username": "alice", "user_id": "user-2", "resolved": true}],
    "created_at": "2026-07-01T12:00:00Z",
    "updated_at": "2026-07-01T12:00:00Z"
  }
}
```

### List Comments (Threaded)

```
GET /api/v1/crm/comments?record_id=rec-1
```

Response `200 OK`:
```json
{
  "data": [
    {
      "id": "parent-uuid",
      "body": "Parent comment",
      "replies": [
        {
          "id": "reply-uuid",
          "body": "Reply to parent",
          "replies": []
        }
      ]
    }
  ],
  "total": 2,
  "page": 1,
  "page_size": 50
}
```

## Database Tables

### `comments` (infrastructure table, Go-generated)

| Column | Type | Description |
|--------|------|-------------|
| `id` | TEXT PK | UUID |
| `body` | TEXT | Comment text (may contain `@mentions`) |
| `commentable_type` | TEXT | Model name (e.g., `crm`, `sales_order`) |
| `commentable_id` | TEXT | Record UUID |
| `parent_id` | TEXT NULL | Parent comment UUID (for threaded replies) |
| `author_id` | TEXT | User UUID of comment author |
| `mentions` | TEXT | JSON array of resolved mentions |
| `created_at` | DATETIME | |
| `updated_at` | DATETIME | |
| `active` | INTEGER | Soft delete flag (1=active, 0=deleted) |

### `notifications` (infrastructure table, Go-generated)

| Column | Type | Description |
|--------|------|-------------|
| `id` | TEXT PK | UUID |
| `user_id` | TEXT | Mentioned user UUID |
| `type` | VARCHAR(50) | Always `mention` |
| `title` | TEXT | "You were mentioned" |
| `body` | TEXT | Comment excerpt (max 200 chars) |
| `model_name` | TEXT | commentable_type |
| `record_id` | TEXT | commentable_id |
| `comment_id` | TEXT | Comment UUID |
| `read` | INTEGER | 0=unread, 1=read |
| `active` | INTEGER | |
| `created_at` | DATETIME | |
| `updated_at` | DATETIME | |

> **Note:** These tables are Go-generated infrastructure tables (like `audit_log`, `ir_migration`), not JSON model definitions. They do not appear in the admin panel and do not get auto-CRUD.

## Mention Parsing

Mentions are parsed from comment body text using `@username` syntax.

### Parsing Rules

- `@alice` → mention `{username: "alice"}`
- `@alice.bob` → mention `{username: "alice.bob"}` (dots and dashes allowed)
- `@alice!` → mention `{username: "alice"}` (punctuation stripped)
- `alice@example.com` → NOT a mention (email excluded)
- `@alice @alice` → deduplicated to one mention

### Resolution

After parsing, `ResolveMentions` queries the user table to find matching usernames. If found, the mention's `UserID` and `Resolved=true` are set. Only resolved mentions trigger notifications.

## Permission Model

Comment access is controlled at two levels:

### 1. Group-Based Permissions (PermissionService)

Checked via `PermissionService.HasPermission(userID, modelName, operation)`:
- **create/write**: `can_write` on the model's group permissions
- **read/list**: `can_read` on the model's group permissions
- **delete**: `can_delete` on the model's group permissions

Superusers bypass all permission checks.

### 2. Record Rules (RecordRuleService)

Checked via `RecordRuleService.CheckRecordAccess(userID, modelName, recordID, operation)`:
- Domain filters restrict which records a user can interact with
- Global rules intersect (AND), group rules union (OR)
- For comments, the record rule is evaluated against the `commentable_id` record

> **Important:** The handler uses `fiber.NewError(403, ...)` for permission denials, not `c.Status(403).JSON(...)`. This is because `c.Status().JSON()` returns `nil` on success, which would cause the caller to incorrectly proceed. `fiber.NewError` returns a non-nil error that properly halts handler execution.

## WebSocket Events

| Event | Channel | Payload |
|-------|---------|---------|
| Comment created | `comment:{commentable_type}:{commentable_id}` | `{type: "comment.created", data: CommentRecord}` |
| Mention notification | `notification:user:{user_id}` | `{type: "notification.created", data: notification}` |

## Wiring (app.go)

```go
commentSvc := comment.NewService(a.DB)
commentSvc.SetBroadcaster(a.WSHub)
commentSvc.SetUserTable(a.ModelRegistry.TableName("user"))

commentPermSvc := persistence.NewPermissionService(a.DB)
commentPermSvc.SetTableNameResolver(func(name string) string { return a.ModelRegistry.TableName(name) })

commentRRSvc := persistence.NewRecordRuleService(a.DB)
commentRRSvc.SetTableNameResolver(func(name string) string { return a.ModelRegistry.TableName(name) })

commentHandler := api.NewCommentHandler(api.CommentHandlerConfig{
    Service:           commentSvc,
    PermissionService: commentPermSvc,
    RecordRuleService: commentRRSvc,
    ModelRegistry:     a.ModelRegistry,
})

api.RegisterCommentRoutes(a.Fiber, "/api/v1", a.ModuleRegistry.InstalledNames(), commentHandler, a.JWTConfig)
```

## Tests

| Package | Tests | Coverage |
|---------|-------|----------|
| `internal/runtime/comment` | 15 | Mention parsing (9 cases), mention resolution, create+notify, list by record |
| `internal/presentation/api` | 5 | Create+list, threading, missing record_id, empty body, permission denied |

```bash
cd engine
go test ./internal/runtime/comment/ -v -count=1
go test ./internal/presentation/api/ -run TestCommentHandler -v -count=1
```

## Future Enhancements

- [ ] Comment editing (with edit history)
- [ ] Reactions/emoji on comments
- [ ] File attachments in comments
- [ ] `@group` mentions (notify all group members)
- [ ] Email notifications for mentions (in addition to in-app)
- [ ] Comment search within a record thread
- [ ] Pagination within deep threads
