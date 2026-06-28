# BcPrintPreview

Preview pane (iframe on web) with a print trigger button.

## Props

| Prop           | Type                                                  | Default            | Description                          |
|----------------|-------------------------------------------------------|--------------------|--------------------------------------|
| `content`      | `string`                                              | —                  | HTML string to render in the preview |
| `title`        | `string`                                              | `'Print Preview'`  | Toolbar title                        |
| `triggerLabel` | `string`                                              | `'Print'`          | Print trigger button label           |
| `disabled`     | `boolean`                                             | `false`            | Disable the trigger button           |
| `hideToolbar`  | `boolean`                                             | `false`            | Hide the toolbar (preview-only)      |
| `height`       | `string`                                              | `'600px'`          | Height of the preview pane (CSS)     |
| `className`    | `string`                                              | —                  | Extra class for the root element     |
| `slotProps`    | `Partial<Record<BcPrintPreviewParts, ...>>`           | —                  | Per-part overrides (ADR-009)         |
| `onPrint`      | `(value: void, meta: BcChangeMeta<BcPrintTriggerReason>) => void` | —        | Fired when user clicks Print         |
| `onError`      | `(value: string, meta: BcChangeMeta<'iframe-load'>) => void` | —          | Fired if the iframe fails to load    |

## A11y

- `role=region` with `aria-label="Print preview"`
- Trigger button has >=44px touch target
- Iframe has `title` + `sandbox="allow-same-origin"`

## Behavior

On web, the component renders an iframe whose content is written on mount and
when `content` changes. The trigger button first attempts `iframe.contentWindow.print()`,
falling back to `window.print()`.

## Example

```tsx
<BcPrintPreview
  title="Invoice Preview"
  content={invoiceHtml}
  height="400px"
  onPrint={() => analytics.track('print')}
/>
```

## Flutter

```dart
BcPrintPreview(
  title: 'Invoice Preview',
  content: invoiceHtml,
  height: 400,
  onPrint: (_, meta) => doPlatformPrint(),
)
```

Note: Flutter preview surfaces stripped text (HTML rendering requires a
package the consumer chooses). The parent wires `onPrint` to the platform
print API of their choice (e.g. the `printing` package).
