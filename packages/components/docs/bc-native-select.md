# BcNativeSelect

Native picker from a list of options. iOS renders a CupertinoPicker wheel; Android renders a Material dropdown. Flutter typedef alias to `BcSelect` (zero duplication).

## Props

| Prop            | Type                                              | Default  | Description                              |
|-----------------|---------------------------------------------------|----------|------------------------------------------|
| `value`         | `string \| number`                                | —        | Controlled selected value                |
| `defaultValue`  | `string \| number`                                | —        | Initial value (uncontrolled)             |
| `options`       | `{ label: string; value: string \| number }[]`    | `[]`     | Available choices                        |
| `placeholder`   | `string`                                          | —        | Placeholder when no selection            |
| `label`         | `string`                                          | —        | Leading label                            |
| `disabled`      | `boolean`                                         | `false`  | Disable interaction                      |
| `onValueChange` | `(value: string \| number, meta: BcChangeMeta) => void` | —  | Fired on selection change                |

## A11y

- `label` renders as the accessible name
- iOS CupertinoPicker wheel / Android Material dropdown via platform defaults

## Example

```tsx
<BcNativeSelect
  label="Country"
  placeholder="Select country"
  options={[
    { label: 'Indonesia', value: 'id' },
    { label: 'Singapore', value: 'sg' },
  ]}
  value={country}
  onValueChange={setCountry}
/>
```

## Flutter

```dart
BcSelect( // typedef alias to BcSelect
  label: 'Country',
  placeholder: 'Select country',
  options: [
    BcSelectOption(label: 'Indonesia', value: 'id'),
    BcSelectOption(label: 'Singapore', value: 'sg'),
  ],
  value: country,
  onValueChange: (v, _) => setCountry(v),
)
```
