# BcAlert

Persistent inline alert banner with 4 severity levels (info, success, warning, error).

## Contract scope

Distinct from `BcNotification` by being **persistent** (no auto-dismiss, no action
button). Intended for in-flow form/page validation messages and static notices.
Visibility is presence-driven (no `open` prop).

## Props

| Prop          | Type                                                            | Default     | Description                          |
|---------------|----------------------------------------------------------------|-------------|--------------------------------------|
| `severity`    | `'info' \| 'success' \| 'warning' \| 'error'`                  | `'info'`    | Severity variant + color/icon token |
| `title`       | `string`                                                       | —           | Bold title line                      |
| `description` | `string`                                                       | —           | Body text below the title           |
| `dismissible` | `boolean`                                                      | `false`     | Show the close (dismiss) control    |
| `hideIcon`    | `boolean`                                                      | `false`     | Hide the leading severity icon      |
| `className`   | `string`                                                       | —           | Extra class for the root element    |
| `slotProps`   | `Partial<Record<BcAlertParts, Record<string, unknown>>>`       | —           | Per-part className/data/aria overrides (ADR-009) |
| `onDismiss`   | `(value: true, meta: BcChangeMeta<BcAlertDismissReason>) => void` | —        | Fired when user dismisses            |

## A11y

- `error` severity: `role=alert`, `aria-live=assertive`
- Other severities: `role=status`, `aria-live=polite`
- Dismiss button has `aria-label="Dismiss alert"`

## Example

```tsx
<BcAlert
  severity="warning"
  title="Storage almost full"
  description="You have used 95% of your quota."
  dismissible
  onDismiss={() => setHidden(true)}
/>
```

## Flutter

```dart
BcAlert(
  severity: BcAlertSeverity.warning,
  title: 'Storage almost full',
  description: '95% used.',
  dismissible: true,
  onDismiss: (v, meta) => setHidden(),
)
```

Platform parity: Material + Cupertino (severity-driven icon + token colors).
