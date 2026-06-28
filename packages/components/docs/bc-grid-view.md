# BcGridView

Spreadsheet-like editable grid with cell navigation and inline editing.

## Contract scope

Adapt tier. Web: HTML table with `role=grid` and arrow-key cell navigation,
click-to-edit cells. Native: RN FlatList with display-only rows (no inline
edit on mobile, rows are display + `onRowSelect`).

## Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `columns` | `BcGridViewColumn[]` | — (required) | Column definitions |
| `rows` | `Record<string, any>[]` | — (required) | Row data |
| `cellEditable` | `boolean` | `false` | Enable inline cell editing (web only) |
| `className` | `string` | — | Extra class for root |
| `slotProps` | `Partial<Record<BcGridViewParts, ...>>` | — | Per-part overrides |
| `onCellChange` | `(row, col, value, meta) => void` | — | Fired on cell edit commit |
| `onRowSelect` | `(row, index) => void` | — | Fired on row click |

## A11y

- Web: `role=grid` + `role=row` + `role=gridcell`, `tabIndex` for focus, `ring-ring` focus indicator
- Arrow Up/Down/Left/Right navigate cells; Enter/F2 to edit; Escape cancels

## Example

```tsx
<BcGridView
  columns={[
    { key: 'name', label: 'Name', editable: true },
    { key: 'qty', label: 'Qty', editable: true },
  ]}
  rows={[{ name: 'Widget', qty: 10 }]}
  cellEditable
  onCellChange={(row, col, value, meta) => console.log(meta.rowIndex, meta.colKey, value)}
  onRowSelect={(row) => console.log(row)}
/>
```

Tier: Adapt (`.web.tsx` + `.native.tsx`). Flutter parity: Material-style display-only rows with 44px touch target.
