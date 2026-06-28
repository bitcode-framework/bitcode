# BcNativeInput

Native text input with controlled/uncontrolled value, secure entry, and keyboard-type variants. Flutter typedef alias to `BcInput` (zero duplication since `BcInput` is already platform-aware via `flutter_platform_widgets`).

## Props

| Prop               | Type                                                | Default       | Description                              |
|--------------------|-----------------------------------------------------|---------------|------------------------------------------|
| `value`            | `string`                                            | —             | Controlled value                         |
| `defaultValue`     | `string`                                            | `''`          | Initial value (uncontrolled)             |
| `label`            | `string`                                            | —             | Floating/leading label                   |
| `placeholder`      | `string`                                            | —             | Placeholder text                         |
| `secureTextEntry`  | `boolean`                                           | `false`       | Mask input (passwords)                   |
| `keyboardType`     | `'default' \| 'numeric' \| 'email' \| 'phone' \| 'url' \| 'decimal' \| 'ascii'` | `'default'` | RN/web keyboard hint |
| `disabled`         | `boolean`                                           | `false`       | Disable interaction                      |
| `onValueChange`    | `(value: string, meta: BcChangeMeta) => void`       | —             | Fired on input change                    |

## A11y

- `label` renders as the accessible name (`aria-label` / Flutter `Semantics`)
- `disabled` exposes `aria-disabled`

## Example

```tsx
<BcNativeInput
  label="Email"
  keyboardType="email"
  placeholder="you@example.com"
  value={email}
  onValueChange={setEmail}
/>
```

## Flutter

```dart
BcInput( // typedef alias to BcInput — zero duplication
  label: 'Email',
  keyboardType: BcKeyboardType.email,
  placeholder: 'you@example.com',
  value: email,
  onValueChange: (v, _) => setEmail(v),
)
```
