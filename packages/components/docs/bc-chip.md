# BcChip

Compact tag/chip with optional avatar and remove control.

## Contract scope

Distinct from `BcBadge` (display-only) by supporting interactive removal,
and from `BcTag` (slightly larger, often static) by being compact and
removal-oriented.

## Props

| Prop             | Type                                            | Default       | Description                              |
|------------------|-------------------------------------------------|---------------|------------------------------------------|
| `label`          | `string`                                        | (required)    | Chip label text                          |
| `variant`        | `'default' \| 'primary' \| 'success' \| 'warning' \| 'destructive' \| 'info'` | `'default'` | Visual variant (token-driven) |
| `avatar`         | `string`                                        | —             | Image URL/data-URI or 2-char initials    |
| `removable`      | `boolean`                                       | `false`       | Show the remove (X) control              |
| `disabled`       | `boolean`                                       | `false`       | Disable interaction (dim + non-interactive) |
| `removed`        | `boolean`                                       | —             | Controlled removed state                 |
| `defaultRemoved` | `boolean`                                       | `false`       | Initial removed state (uncontrolled)     |
| `onRemove`       | `(value: true, meta: BcChangeMeta) => void`     | —             | Fired when user removes                  |
| `onClick`        | `() => void`                                    | —             | Whole-chip click handler                 |

## A11y

- `role=group` with the chip label as the semantic label
- Remove button `aria-label="Remove {label}"`
- Disabled chip exposes `aria-disabled="true"`

## Example

```tsx
<BcChip
  label="TypeScript"
  avatar="TS"
  variant="primary"
  removable
  onRemove={() => remove('ts')}
/>
```

## Flutter

```dart
BcChip(
  label: 'TypeScript',
  avatar: 'TS',
  variant: BcChipVariant.primary,
  removable: true,
  onRemove: (v, meta) => remove('ts'),
)
```
