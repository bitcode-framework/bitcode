# BcPrintArea

Marks a printable region of the page. With `auto=true` (default), only this
region prints; everything else is hidden via `@media print` isolation.

## Props

| Prop        | Type                                              | Default              | Description                              |
|-------------|---------------------------------------------------|----------------------|------------------------------------------|
| `auto`      | `boolean`                                         | `true`               | Inject global print-isolation CSS         |
| `ariaLabel` | `string`                                          | `'Printable region'` | Accessibility label for the region        |
| `className` | `string`                                          | —                    | Extra class for the root element         |
| `slotProps` | `Partial<Record<BcPrintAreaParts, ...>>`          | —                    | Per-part overrides (ADR-009)             |
| `onReady`   | `(value: void, meta: BcChangeMeta<'mount'>) => void` | —                | Fired on mount (telemetry)               |

## A11y

- `role=region` with `aria-label`
- No special keyboard handling (printable region is a landmark, not interactive)

## Behavior

When `auto=true`, the component injects an idempotent `@media print` style on
mount that hides everything outside any `.bc-print-area-isolate` element. The
style is only added once per page.

## Example

```tsx
<BcPrintArea ariaLabel="Invoice" auto>
  <InvoiceContent />
</BcPrintArea>
```

## Flutter

```dart
BcPrintArea(
  ariaLabel: 'Invoice',
  child: InvoiceContent(),
  onReady: () => print('ready'),
)
```

Note: Flutter printing routes through the platform print dialog. `auto` is
accepted for contract parity but has no effect on Flutter (the consumer
always manages their own print flow).
