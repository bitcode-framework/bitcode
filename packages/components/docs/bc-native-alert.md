# BcNativeAlert

Native modal alert dialog with controlled open, title/message, and configurable confirm/cancel labels. Flutter exposes a `show()` static method via `showPlatformDialog` (CupertinoAlertDialog / AlertDialog) with focus trap.

## Props

| Prop            | Type                                              | Default       | Description                              |
|-----------------|---------------------------------------------------|---------------|------------------------------------------|
| `open`          | `boolean`                                         | —             | Controlled open state                    |
| `defaultOpen`   | `boolean`                                         | `false`       | Initial open state (uncontrolled)        |
| `title`         | `string`                                          | —             | Alert title                              |
| `message`       | `string`                                          | —             | Alert body text                          |
| `confirmLabel`  | `string`                                          | `'OK'`        | Confirm button label                     |
| `cancelLabel`   | `string`                                          | `'Cancel'`    | Cancel button label (omit to hide)       |
| `onConfirm`     | `() => void`                                      | —             | Fired on confirm                         |
| `onCancel`      | `() => void`                                      | —             | Fired on cancel                          |
| `onOpenChange`  | `(open: boolean, meta: BcChangeMeta) => void`     | —             | Fired on open/close                      |

## A11y

- `role=alertdialog` with focus trap
- Confirm is the primary action; cancel dismisses

## Example

```tsx
<BcNativeAlert
  open={show}
  title="Delete account"
  message="This action cannot be undone."
  confirmLabel="Delete"
  cancelLabel="Keep"
  onConfirm={deleteAccount}
  onCancel={() => setShow(false)}
  onOpenChange={setShow}
/>
```

## Flutter

```dart
final result = await BcNativeAlert.show(
  context: context,
  title: 'Delete account',
  message: 'This action cannot be undone.',
  confirmLabel: 'Delete',
  cancelLabel: 'Keep',
);
if (result == true) deleteAccount();
```

`show()` wraps `showPlatformDialog` (CupertinoAlertDialog / AlertDialog) with focus trap.
