# BcKanbanView

Kanban board with drag-and-drop cards between columns.

## Contract scope

Adapt tier. Web: HTML5 DnD API (draggable cards, droppable columns).
Native: long-press to select a card, then tap the destination column to move it.

## Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `columns` | `BcKanbanViewColumn[]` | — (required) | Column definitions |
| `cards` | `BcKanbanViewCard[]` | — (required) | Card data |
| `className` | `string` | — | Extra class for root |
| `slotProps` | `Partial<Record<BcKanbanViewParts, ...>>` | — | Per-part overrides |
| `onCardMove` | `(cardId, fromColumn, toColumn, meta) => void` | — | Fired on card move |
| `onCardClick` | `(card) => void` | — | Fired on card click |

## A11y

- Web: `role=region` + `aria-label="kanban board"`, cards `role=listitem` + `tabIndex=0`
- Cards draggable via keyboard (Enter/Space activates `onCardClick`)

## Example

```tsx
<BcKanbanView
  columns={[
    { id: 'todo', title: 'To Do' },
    { id: 'done', title: 'Done', color: '#22c55e' },
  ]}
  cards={[
    { id: 'c1', title: 'Task 1', columnId: 'todo', priority: 'high' },
  ]}
  onCardMove={(cardId, from, to, meta) => moveCard(cardId, to)}
  onCardClick={(card) => openCard(card)}
/>
```

Tier: Adapt (`.web.tsx` + `.native.tsx`). Flutter parity: Draggable + DragTarget with BcTheme tokens.
