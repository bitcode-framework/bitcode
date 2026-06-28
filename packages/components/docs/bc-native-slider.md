# BcNativeSlider

Native range slider with min/max/step and drag/commit change reasons. Flutter typedef alias to `BcSlider` (zero duplication).

## Props

| Prop            | Type                                                | Default | Description                              |
|-----------------|-----------------------------------------------------|---------|------------------------------------------|
| `value`         | `number`                                            | `0`     | Controlled value                         |
| `defaultValue`  | `number`                                            | `0`     | Initial value (uncontrolled)             |
| `min`           | `number`                                            | `0`     | Minimum value                            |
| `max`           | `number`                                            | `100`   | Maximum value                            |
| `step`          | `number`                                            | `1`     | Value step increment                     |
| `disabled`      | `boolean`                                           | `false` | Disable interaction                      |
| `onValueChange` | `(value: number, meta: BcChangeMeta<BcSliderDragReason>) => void` | — | Fired on change (`drag` while dragging, `commit` on release) |

## A11y

- `role=slider`, `aria-valuemin`/`aria-valuemax`/`aria-valuenow` reflect bounds + value
- Keyboard: arrow keys to adjust

## Example

```tsx
<BcNativeSlider
  min={0}
  max={100}
  step={5}
  value={volume}
  onValueChange={(v, meta) => {
    setVolume(v);
    if (meta.reason === 'commit') saveVolume(v);
  }}
/>
```

## Flutter

```dart
BcSlider( // typedef alias to BcSlider
  min: 0,
  max: 100,
  step: 5,
  value: volume,
  onValueChange: (v, meta) {
    setVolume(v);
    if (meta.reason == BcSliderDragReason.commit) saveVolume(v);
  },
)
```
