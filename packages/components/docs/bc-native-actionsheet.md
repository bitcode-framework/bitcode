# BcNativeActionSheet

Native bottom action sheet with a list of items (each optionally destructive or disabled). Flutter exposes a `show()` static method via `showCupertinoModalPopup` (CupertinoActionSheet) / `showModalBottomSheet`.

## Props

| Prop            | Type                                                                | Default     | Description                              |
|-----------------|---------------------------------------------------------------------|-------------|------------------------------------------|
| `open`          | `boolean`                                                           | —           | Controlled open state                    |
| `defaultOpen`   | `boolean`                                                           | `false`     | Initial open state (uncontrolled)        |
| `items`         | `{ label: string; value: string \| number; destructive?: boolean; disabled?: boolean }[]` | `[]` | Action items |
| `title`         | `string`                                                            | —           | Optional sheet title                     |
| `cancelLabel`   | `string`                                                            | `'Cancel'`  | Cancel button label (omit to hide)       |
| `onSelect`      | `(value: string \| number, meta: BcChangeMeta) => void`             | —           | Fired on item select                     |
| `onCancel`      | `() => void`                                                        | —           | Fired on cancel                          |
| `onOpenChange`  | `(open: boolean, meta: BcChangeMeta) => void`                       | —           | Fired on open/close                      |

## A11y

- `role=dialog` with focus trap
- Destructive items use warning tokens; disabled items are non-interactive

## Example

```tsx
<BcNativeActionSheet
  open={open}
  title="Photo"
  cancelLabel="Cancel"
  items={[
    { label: 'Camera', value: 'cam' },
    { label: 'Gallery', value: 'gal' },
    { label: 'Remove', value: 'del', destructive: true },
  ]}
  onSelect={(v) => pickPhoto(v)}
  onCancel={() => setOpen(false)}
  onOpenChange={setOpen}
/>
```

## Flutter

```dart
final result = await BcNativeActionSheet.show(
  context: context,
  title: 'Photo',
  cancelLabel: 'Cancel',
  items: [
    BcActionItem(label: 'Camera', value: 'cam'),
    BcActionItem(label: 'Gallery', value: 'gal'),
    BcActionItem(label: 'Remove', value: 'del', destructive: true),
  ],
);
if (result != null) pickPhoto(result);
```

`show()` wraps `showCupertinoModalPopup` (CupertinoActionSheet) / `showModalBottomSheet` (Material).
