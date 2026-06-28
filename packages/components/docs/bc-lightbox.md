# BcLightbox

Full-screen image/video gallery with keyboard navigation and prev/next controls.

## Contract scope

An open-state overlay (like `BcDialog`) specialized for media. Controlled open
(`open`) and active index (`index`) both use resolveControllable so they work
in controlled or uncontrolled mode. Items are a typed list of image/video
sources.

## Props

| Prop           | Type                                                                       | Default     | Description                              |
|----------------|----------------------------------------------------------------------------|-------------|------------------------------------------|
| `open`         | `boolean`                                                                  | —           | Controlled open state                    |
| `defaultOpen`  | `boolean`                                                                  | `false`     | Initial open state (uncontrolled)        |
| `items`        | `{ src, type: 'image'\|'video', alt?, caption? }[]`                       | `[]`        | Media items to display                   |
| `index`        | `number`                                                                   | —           | Controlled active index                  |
| `defaultIndex` | `number`                                                                   | `0`         | Initial active index (uncontrolled)      |
| `closable`     | `boolean`                                                                  | `true`      | Show the close button                    |
| `adapt`        | `AdaptProps`                                                               | —           | Responsive/platform morph (ADR-013)      |
| `className`    | `string`                                                                   | —           | Extra class for the root element         |
| `slotProps`    | `Partial<Record<BcLightboxParts, Record<string, unknown>>>`               | —           | Per-part className/data/aria overrides (ADR-009) |
| `onIndexChange`| `(index: number, meta: BcChangeMeta<'prev'\|'next'\|'programmatic'>) => void` | —    | Fired when active index changes          |
| `onOpenChange` | `(open: boolean, meta: BcChangeMeta<'close-button'\|'overlay-click'\|'escape'\|'programmatic'>) => void` | — | Fired when open state changes |

## A11y

- `role=dialog`, `aria-modal=true`, `aria-label="image gallery"`
- Keyboard: `Escape` closes, `ArrowLeft`/`ArrowRight` navigate
- Prev/next/close buttons have `min 44px` touch targets

## Example

```tsx
<BcLightbox
  open={isOpen}
  onOpenChange={setIsOpen}
  items={gallery}
  index={currentIdx}
  onIndexChange={(i) => setCurrentIdx(i)}
/>
```

## Flutter

```dart
BcLightbox(
  open: true,
  items: items,
  onOpenChange: (v, meta) => setOpen(v),
  onIndexChange: (i, meta) => setIndex(i),
)
```

Platform parity: web overlay + keyboard nav; Flutter PageView swipe + semantics.
