# BcNativePicker

Native scroll-wheel picker with controlled open + value. iOS renders a CupertinoPicker wheel; Android renders a Material dialog. Flutter exposes a `show()` static method via `showCupertinoModalPopup`/`showDialog`.

## Props

| Prop            | Type                                                  | Default  | Description                              |
|-----------------|-------------------------------------------------------|----------|------------------------------------------|
| `open`          | `boolean`                                             | —        | Controlled open state                    |
| `defaultOpen`   | `boolean`                                             | `false`  | Initial open state (uncontrolled)        |
| `value`         | `string \| number`                                    | —        | Controlled selected value                |
| `options`       | `{ label: string; value: string \| number }[]`        | `[]`     | Available choices                        |
| `title`         | `string`                                              | —        | Picker title / header label              |
| `onValueChange` | `(value: string \| number, meta: BcChangeMeta) => void` | —      | Fired on selection                       |
| `onOpenChange`  | `(open: boolean, meta: BcChangeMeta) => void`         | —        | Fired on open/close                      |

## A11y

- Open state renders `role=dialog`
- iOS CupertinoPicker wheel / Android Material dialog via platform defaults

## Example

```tsx
<BcNativePicker
  title="Select hour"
  options={[{ label: '09:00', value: 9 }, { label: '10:00', value: 10 }]}
  open={open}
  value={hour}
  onValueChange={setHour}
  onOpenChange={setOpen}
/>
```

## Flutter

```dart
final result = await BcNativePicker.show(
  context: context,
  title: 'Select hour',
  options: [
    BcPickerOption(label: '09:00', value: 9),
    BcPickerOption(label: '10:00', value: 10),
  ],
  value: hour,
);
if (result != null) setHour(result);
```

`show()` wraps `showCupertinoModalPopup` (CupertinoPicker wheel) / `showDialog` (Material).
