# BcNativeToast

Native toast notification with auto-close timer and optional action button. Flutter renders a controlled-open widget with a SnackBar-like overlay.

## Props

| Prop            | Type                                              | Default        | Description                              |
|-----------------|---------------------------------------------------|----------------|------------------------------------------|
| `open`          | `boolean`                                         | —              | Controlled open state                    |
| `defaultOpen`   | `boolean`                                         | `false`        | Initial open state (uncontrolled)        |
| `message`       | `string`                                          | —              | Toast message text                       |
| `actionLabel`   | `string`                                          | —              | Optional action button label             |
| `duration`      | `number`                                          | `3000`         | Auto-close delay in ms (0 = manual)      |
| `onAction`      | `() => void`                                      | —              | Fired on action button press             |
| `onOpenChange`  | `(open: boolean, meta: BcChangeMeta) => void`     | —              | Fired on open/close (including auto)     |

## A11y

- `role=status`, `aria-live=polite`
- Action button is keyboard-focusable

## Example

```tsx
<BcNativeToast
  open={open}
  message="Saved to draft"
  actionLabel="Undo"
  duration={4000}
  onAction={undoSave}
  onOpenChange={setOpen}
/>
```

## Flutter

```dart
BcNativeToast(
  open: open,
  message: 'Saved to draft',
  actionLabel: 'Undo',
  duration: Duration(milliseconds: 4000),
  onAction: undoSave,
  onOpenChange: (v, _) => setOpen(v),
)
```

Controlled-open widget rendered as a SnackBar-like overlay.
