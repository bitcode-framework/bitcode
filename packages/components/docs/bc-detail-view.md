# BcDetailView

Master-detail record renderer with field grid, optional tabs, and action buttons.

## Contract scope

Pure tier. Reads already-resolved field data (no fetching). Composes either
a descriptions-style field grid or a tabbed body with a header (title +
subtitle + action buttons).

## Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `title` | `string` | — | Header title |
| `subtitle` | `string` | — | Header subtitle (muted) |
| `fields` | `BcDetailViewField[]` | — | Field rows for descriptions grid |
| `tabs` | `BcDetailViewTab[]` | — | Tabbed body (takes precedence over fields) |
| `actions` | `BcDetailViewAction[]` | — | Header action buttons |
| `emptyText` | `string` | `'No data'` | Empty-state message |
| `loading` | `boolean` | `false` | Show loading skeleton |
| `className` | `string` | — | Extra class for root |
| `slotProps` | `Partial<Record<BcDetailViewParts, ...>>` | — | Per-part overrides (ADR-009) |
| `onAction` | `(action, meta) => void` | — | Fired when an action button is clicked |

## A11y

- `role=region` + `aria-label`
- Tabs: `role=tab` + `aria-selected`
- Action buttons: `min-h-[36px]` touch target

## Example

```tsx
<BcDetailView
  title="Contact"
  subtitle="John Doe"
  fields={[
    { label: 'Email', value: 'john@test.com' },
    { label: 'Phone', value: '+1-555-0100' },
  ]}
  actions={[{ label: 'Edit', variant: 'primary', onClick: handleEdit }]}
  onAction={(action, meta) => console.log(action, meta)}
/>
```

## Flutter

```dart
BcDetailView(
  title: 'Contact',
  subtitle: 'John Doe',
  fields: [
    BcDetailViewField(label: 'Email', value: 'john@test.com'),
  ],
  actions: [
    BcDetailViewAction(label: 'Edit', variant: BcDetailViewActionVariant.primary),
  ],
)
```

Tier: Pure (`.lite.tsx`). Flutter parity: dependency-free StatelessWidget with BcTheme tokens.
