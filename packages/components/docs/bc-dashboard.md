# BcDashboard

Grid of widgets with configurable layout and custom rendering.

## Contract scope

Pure tier. Renders a CSS grid of widget cards. The actual widget content is
delegated to a `renderWidget` prop (defaults to a placeholder showing the
kind + title), keeping the component Pure without coupling to chart/table deps.

## Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `widgets` | `BcDashboardWidget[]` | — (required) | Widget configurations |
| `columns` | `number` | `3` | Grid column count (clamped 1-12) |
| `gap` | `'none' \| 'sm' \| 'md' \| 'lg'` | `'md'` | Gap between widgets |
| `loading` | `boolean` | `false` | Show loading skeleton |
| `emptyText` | `string` | `'No widgets'` | Empty-state message |
| `className` | `string` | — | Extra class for root |
| `slotProps` | `Partial<Record<BcDashboardParts, ...>>` | — | Per-part overrides |
| `renderWidget` | `(widget) => unknown` | — | Custom widget renderer |
| `onWidgetClick` | `(widget, meta) => void` | — | Fired on widget card click |

## A11y

- `role=region` + `aria-label="dashboard"`
- Widget cards clickable with pointer cursor

## Example

```tsx
<BcDashboard
  columns={4}
  widgets={[
    { id: 'rev', title: 'Revenue', kind: 'stat', span: 1 },
    { id: 'chart', title: 'Sales Chart', kind: 'chart', span: 2 },
  ]}
  renderWidget={(w) => w.kind === 'stat' ? <StatCard {...w.props} /> : <Chart {...w.props} />}
  onWidgetClick={(w, meta) => console.log(w.id)}
/>
```

## Flutter

```dart
BcDashboard(
  columns: 4,
  widgets: [
    BcDashboardWidget(id: 'rev', title: 'Revenue', kind: BcDashboardWidgetKind.stat),
  ],
  renderWidget: (w) => StatCard(),
)
```

Tier: Pure (`.lite.tsx`). Flutter parity: GridView.builder with BcTheme tokens.
