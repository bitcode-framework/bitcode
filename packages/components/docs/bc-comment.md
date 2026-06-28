# BcComment

Threaded comment with author, body, optional timestamp, and inline reply form.

## Props

| Prop              | Type                                                          | Default                | Description                              |
|-------------------|---------------------------------------------------------------|------------------------|------------------------------------------|
| `author`          | `string`                                                      | (required)             | Author display name                      |
| `avatar`          | `string`                                                      | —                      | Avatar URL/data-URI or 2-char initials   |
| `body`            | `string`                                                      | (required)             | Comment body text                        |
| `timestamp`       | `string`                                                      | —                      | Timestamp (formatted by caller)          |
| `replyable`       | `boolean`                                                     | `false`                | Show inline reply form                   |
| `replyPlaceholder`| `string`                                                      | `'Write a reply...'`   | Reply textarea placeholder               |
| `disabled`        | `boolean`                                                     | `false`                | Disable the reply form                   |
| `onReply`         | `(value: string, meta: BcChangeMeta<BcCommentSubmitReason>) => void` | —        | Fired when user submits a reply          |

## Behavior

- Reply form: Enter submits, Shift+Enter inserts newline.
- Submit trims + guards empty/disabled drafts.
- Avatar renders `<img>` for URL, initials for short strings.

## Example

```tsx
<BcComment
  author="Jane Doe"
  avatar="JD"
  body="What do you think?"
  timestamp="5 min ago"
  replyable
  onReply={(text) => submitReply(text)}
/>
```

## Flutter

```dart
BcComment(
  author: 'Jane Doe',
  avatar: 'JD',
  body: 'What do you think?',
  replyable: true,
  onReply: (text, _) => submitReply(text),
)
```
