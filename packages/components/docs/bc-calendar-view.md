# BcCalendarView

Month/week/day calendar with event overlay and navigation.

## Contract scope

Adapt tier. Web: month grid (7x6), week grid (7x1), day view (hourly column)
with event badges. Native: vertical day list with event counts.

## Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `view` | `'month' \| 'week' \| 'day'` | `'month'` | Active view mode |
| `date` | `Date` | today | Controlled center date |
| `events` | `BcCalendarViewEvent[]` | `[]` | Events to overlay |
| `className` | `string` | — | Extra class for root |
| `slotProps` | `Partial<Record<BcCalendarViewParts, ...>>` | — | Per-part overrides |
| `onDateSelect` | `(date, meta) => void` | — | Fired on day cell click |
| `onEventClick` | `(event) => void` | — | Fired on event badge click |
| `onViewChange` | `(view, meta) => void` | — | Fired on view mode change |

## A11y

- Web: `role=application` + `aria-label="calendar"`, cells `role=gridcell` + `tabIndex=0`
- Prev/next/today nav buttons with aria-labels

## Example

```tsx
<BcCalendarView
  view="month"
  date={new Date(2026, 5, 15)}
  events={[
    { id: 'e1', title: 'Meeting', start: new Date(2026, 5, 15, 10), end: new Date(2026, 5, 15, 11) },
  ]}
  onDateSelect={(date, meta) => console.log(date)}
  onEventClick={(event) => openEvent(event)}
/>
```

Tier: Adapt (`.web.tsx` + `.native.tsx`). Flutter parity: month/week/day grid with BcTheme tokens.
