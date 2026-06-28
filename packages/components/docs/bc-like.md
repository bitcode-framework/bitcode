# BcLike

Like/unlike toggle with optional count.

## Props

| Prop            | Type                                                          | Default  | Description                          |
|-----------------|---------------------------------------------------------------|----------|--------------------------------------|
| `liked`         | `boolean`                                                     | —        | Controlled liked state               |
| `defaultLiked`  | `boolean`                                                     | `false`  | Initial liked state (uncontrolled)   |
| `count`         | `number`                                                      | —        | Like count to display (omit to hide) |
| `disabled`      | `boolean`                                                     | `false`  | Disable interaction                  |
| `ariaLabel`     | `string`                                                      | `'Like'` | Button aria-label                    |
| `onLikedChange` | `(value: boolean, meta: BcChangeMeta<BcLikeChangeReason>) => void` | — | Fired on toggle (reason: toggle-on/toggle-off) |

## Example

```tsx
<BcLike liked={liked} count={count} onLikedChange={setLiked} />
```

## Flutter

```dart
BcLike(liked: liked, count: count, onLikedChange: (v, _) => setLiked(v))
```
