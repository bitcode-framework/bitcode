# Notification System

> Feature #67 — In-App Notification

## Overview

Persistent in-app notification system with user-scoped delivery, WebSocket push, @mention integration, and cross-platform UI (web + Flutter).

## Model

**File:** `embedded/modules/base/models/notification.json`

| Field | Type | Description |
|-------|------|-------------|
| `user_id` | many2one (user) | Recipient user |
| `type` | selection | `info`, `warning`, `error`, `success`, `mention`, `message` |
| `title` | text | Notification title |
| `body` | text | Notification body |
| `data` | json | Arbitrary structured data (metadata, context) |
| `notifiable_type` | text (morph_to) | Polymorphic ref type (e.g., `comment`, `order`) |
| `notifiable_id` | text (morph_to) | Polymorphic ref UUID |
| `read_at` | datetime | When the user read the notification (NULL = unread) |
| `created_at` | datetime | Auto-timestamp |

**Indexes:** `idx_notifications_user_read` on (`user_id`, `read_at`) for efficient unread-first queries.

**API:** `"api": true` — convention-driven CRUD auto-generated. Custom endpoints override via the dedicated NotificationHandler.

## NotificationService

**File:** `engine/internal/runtime/notification/service.go`

| Method | Description |
|--------|-------------|
| `Create(input NotificationCreateInput)` | Creates a notification and pushes via WebSocket |
| `List(userID, opts)` | User-scoped list, unread-first ordering, pagination (`page`, `page_size`), optional `unread_only` filter |
| `UnreadCount(userID)` | Returns count of unread notifications for user |
| `MarkRead(userID, notificationID)` | Marks a single notification as read |
| `MarkAllRead(userID)` | Marks all unread notifications for user as read |

WebSocket push: On `Create`, the service calls `hub.BroadcastToUser(userID, "notification", payload)` to deliver the notification to the connected user in real time.

## REST Endpoints

**File:** `engine/internal/presentation/api/notification_handler.go`

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/notifications` | List notifications for authenticated user (unread-first, pagination) |
| POST | `/api/v1/notifications/:id/read` | Mark a single notification as read |
| POST | `/api/v1/notifications/read-all` | Mark all notifications as read |
| GET | `/api/v1/notifications/unread-count` | Get unread notification count |

All endpoints require authentication. Responses are user-scoped — a user can only see and interact with their own notifications.

**Query parameters (GET /notifications):**
- `page` — page number (default: 1)
- `page_size` — items per page (default: 20)
- `unread_only` — if `true`, only return unread notifications

## WebSocket

**File:** `engine/internal/presentation/websocket/hub.go`

Added `BroadcastToUser(userID, channel, data)` — iterates all connected clients, filters by `Client.UserID`, and sends the payload on the specified channel. Used by `NotificationService.Create` to push new notifications to the user's browser/app in real time.

## @Mention Integration

The comment service (`engine/internal/runtime/comment/service.go`) automatically creates a notification when a user is @mentioned in a comment. The notification includes:
- `type`: `mention`
- `notifiable_type`: `comment`, `notifiable_id`: the comment UUID
- `data`: comment context (record type, record ID, mentioned by)

## UI Components

### Web (packages/ui)

| Component | Description |
|-----------|-------------|
| `BcNotificationCenter` | Bell button + dropdown panel, unread badge, mark-read, mark-all-read, infinite scroll. Declarative `data-source` attribute for API binding. Pure/Adapt pattern (ADR-014). |
| `BcNotificationItem` | Individual notification row with 6 variant icons (info/warning/error/success/mention/message), timestamp, unread indicator, mark-read button. |

**Renderer:** `<bc-notification-center data-source="/api/v1/notifications" data-part="notification-center">` in navbar partial template.

### Flutter (packages/ui-flutter)

`BcNotificationCenter` widget with adaptive rendering:
- **iOS:** Cupertino dropdown with iOS-style bell icon and badge
- **Android:** Material dropdown with Material bell icon and badge

## i18n

Notification keys added to all 11 locale files (`embedded/modules/base/i18n/*.json`):
- `notification.title`, `notification.empty`, `notification.mark_all_read`, `notification.mark_read`
- `notification.load_more`, `notification.loading`, `notification.unread`, `notification.mentioned_you`
- `notification.authentication_required`, `notification.not_found`
- `notification.type_info/warning/error/success/mention/message`

## Tests

| Suite | Count | Coverage |
|-------|-------|----------|
| `runtime/notification` | 16 | Create, List (unread-first, user-scope, pagination, unread_only), MarkRead, MarkAllRead, UnreadCount, WS push |
| `presentation/api` (notification) | 7 | List unread-first, mark single read, mark all read, unread count, pagination, auth required, user scope isolation |
| `packages/ui` (notification) | 21 | Contract, bell button, badge, dropdown panel, mark-all-read, mark-read, infinite scroll, loading/error/empty states, declarative data source, token palette, 6 variants, BcNotificationItem subcomponent |
| `packages/ui-flutter` (notification) | 10 | fromJson, isUnread, copyWith, equality, bell renders with badge, no badge when 0, 99+ overflow, onMarkRead callback, disabled state |
