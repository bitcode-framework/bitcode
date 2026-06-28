# BcNativeModal

Native modal dialog with 3 presentation styles and a default Slot for content. Flutter renders a controlled-open widget via `showGeneralDialog`.

## Props

| Prop             | Type                                              | Default         | Description                              |
|------------------|---------------------------------------------------|-----------------|------------------------------------------|
| `open`           | `boolean`                                         | —               | Controlled open state                    |
| `defaultOpen`    | `boolean`                                         | `false`         | Initial open state (uncontrolled)        |
| `presentation`   | `'fullscreen' \| 'formSheet' \| 'pageSheet'`      | `'formSheet'`   | Presentation style                       |
| `closable`       | `boolean`                                         | `true`          | Allow close via backdrop / Escape        |
| `title`          | `string`                                          | —               | Optional modal title                     |
| `onOpenChange`   | `(open: boolean, meta: BcChangeMeta) => void`     | —               | Fired on open/close                      |

## Slots

- **default** — modal body content

## A11y

- `role=dialog` with focus trap
- Backdrop dismiss disabled when `closable=false`

## Example

```tsx
<BcNativeModal
  open={open}
  presentation="formSheet"
  title="Edit profile"
  closable
  onOpenChange={setOpen}
>
  <BcNativeInput label="Name" />
  <BcNativeButton onPress={save}>Save</BcNativeButton>
</BcNativeModal>
```

## Flutter

```dart
BcNativeModal(
  open: open,
  presentation: BcModalPresentation.formSheet,
  title: 'Edit profile',
  closable: true,
  onOpenChange: (v, _) => setOpen(v),
  child: Column(children: [
    BcInput(label: 'Name'),
    BcNativeButton(onPress: save, child: Text('Save')),
  ]),
)
```

Controlled-open widget rendered via `showGeneralDialog`.
