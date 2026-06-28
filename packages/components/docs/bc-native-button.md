# BcNativeButton

Native-feeling button with 5 variants, haptic feedback (React Native), and loading state. Adapt component platform-aware via `Pressable` (RN) and `PlatformElevatedButton` (Flutter).

## Props

| Prop        | Type                                                                      | Default       | Description                              |
|-------------|---------------------------------------------------------------------------|---------------|------------------------------------------|
| `variant`   | `'primary' \| 'secondary' \| 'destructive' \| 'outline' \| 'plain'`      | `'primary'`   | Visual variant (token-driven)            |
| `loading`   | `boolean`                                                                 | `false`       | Show spinner, disable interaction        |
| `block`     | `boolean`                                                                 | `false`       | Stretch to full container width          |
| `disabled`  | `boolean`                                                                 | `false`       | Disable interaction                      |
| `haptic`    | `boolean`                                                                 | `true`        | Haptic feedback on press (RN only; web no-op) |
| `label`     | `string`                                                                  | —             | Button label (if no slot children)       |
| `onPress`   | `() => void`                                                              | —             | Fired on press                           |

## A11y

- `role=button` with the label as the accessible name
- Min 44px touch target (mobile)
- `loading`/`disabled` expose `aria-disabled`

## Example

```tsx
<BcNativeButton variant="primary" loading={submitting} block onPress={submit}>
  Save
</BcNativeButton>
```

## Flutter

```dart
BcNativeButton(
  variant: BcNativeButtonVariant.primary,
  loading: submitting,
  block: true,
  haptic: true,
  child: Text('Save'),
  onPress: submit,
)
```

Platform-aware via `PlatformElevatedButton` (Material) / `CupertinoButton` (Cupertino).
