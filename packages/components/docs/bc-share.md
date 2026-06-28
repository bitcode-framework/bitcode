# BcShare

Share intent with platform-native fallback chain (Adapt component).

## Behavior

- **Web**: prefers `navigator.share()` (Web Share API). When unavailable
  or when the user dismisses/cancels, falls back to an inline dialog with
  copy-link + social text.
- **React Native**: calls `Share.share()` from the react-native module.
- **Flutter**: emits the trigger via `onShare` and delegates the actual
  platform share to the caller (keeps `bc_ui_flutter` dependency-free;
  use `share_plus`'s `Share.share(payload.message ?? payload.url)` at the
  call site).

## Props

| Prop        | Type                                                              | Default  | Description                          |
|-------------|-------------------------------------------------------------------|----------|--------------------------------------|
| `title`     | `string`                                                          | —        | Title to share                       |
| `text`      | `string`                                                          | —        | Body text to share                   |
| `url`       | `string`                                                          | —        | URL to share                         |
| `label`     | `string`                                                          | `'Share'`| Trigger button label                 |
| `disabled`  | `boolean`                                                         | `false`  | Disable the trigger                  |
| `adapt`     | `AdaptProps`                                                      | —        | Mobile morph (ADR-013)               |
| `onShare`   | `(value: BcShareResult, meta: BcChangeMeta<BcShareTriggerReason>) => void` | — | Share outcome callback |
| `onError`   | `(error: unknown) => void`                                        | —        | Fired on share error                 |

## BcShareResult

- `'shared'` — share succeeded
- `'dismissed'` — user dismissed the share sheet
- `'unavailable'` — platform share unavailable
- `'error'` — an error occurred

## Example

```tsx
<BcShare
  title="BitCode"
  text="Check this out"
  url="https://example.com"
  onShare={(result) => analytics.track('share', { result })}
/>
```

## Flutter

```dart
BcShare(
  title: 'BitCode',
  text: 'Check this out',
  url: 'https://example.com',
  onShare: (result, _) async {
    final payload = BcSharePayload(text: '...', url: '...');
    await Share.share(payload.message ?? '');
  },
)
```
