# BcNativeSwitch

Native on/off toggle with controlled/uncontrolled checked state. Flutter typedef alias to `BcToggle` (zero duplication).

## Props

| Prop              | Type                                              | Default  | Description                              |
|-------------------|---------------------------------------------------|----------|------------------------------------------|
| `checked`         | `boolean`                                         | —        | Controlled checked state                 |
| `defaultChecked`  | `boolean`                                         | `false`  | Initial checked state (uncontrolled)     |
| `label`           | `string`                                          | —        | Accessible label / leading text          |
| `disabled`        | `boolean`                                         | `false`  | Disable interaction                      |
| `onCheckedChange` | `(value: boolean, meta: BcChangeMeta<BcSwitchToggleReason>) => void` | — | Fired on toggle (reason: toggle) |

## A11y

- `role=switch`, `aria-checked` reflects checked state
- `disabled` exposes `aria-disabled`

## Example

```tsx
<BcNativeSwitch
  label="Airplane mode"
  checked={airplane}
  onCheckedChange={setAirplane}
/>
```

## Flutter

```dart
BcToggle( // typedef alias to BcToggle
  label: 'Airplane mode',
  checked: airplane,
  onCheckedChange: (v, _) => setAirplane(v),
)
```
