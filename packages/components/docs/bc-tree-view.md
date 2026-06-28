# BcTreeView

Hierarchical tree view with expand/collapse and keyboard navigation.

## Contract scope

Adapt tier. Web: nested divs with `role=tree`/`role=treeitem`, `aria-expanded`,
`aria-selected`. Native: recursive `View` with `Pressable` rows (tap to select,
tap chevron to expand).

## Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `data` | `BcTreeViewNode[]` | — (required) | Tree node data |
| `expandedKeys` | `string[]` | — | Controlled expanded node IDs |
| `defaultExpandedKeys` | `string[]` | — | Initial expanded node IDs |
| `selectedKey` | `string` | — | Controlled selected node ID |
| `className` | `string` | — | Extra class for root |
| `slotProps` | `Partial<Record<BcTreeViewParts, ...>>` | — | Per-part overrides |
| `onSelect` | `(node, meta) => void` | — | Fired on node selection |
| `onExpand` | `(key, expanded, meta) => void` | — | Fired on expand/collapse |

## A11y

- Web: `role=tree` + `role=treeitem` + `aria-expanded` + `aria-selected`
- Arrow Up/Down navigate; ArrowRight/Left expand/collapse; Enter/Space selects

## Example

```tsx
<BcTreeView
  data={[
    { id: '1', label: 'Root', children: [
      { id: '1-1', label: 'Child' },
    ]},
  ]}
  defaultExpandedKeys={['1']}
  onSelect={(node, meta) => console.log(node, meta)}
/>
```

Tier: Adapt (`.web.tsx` + `.native.tsx`). Flutter parity: recursive widget with 44px touch target.
