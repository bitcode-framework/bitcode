# BcFollow

Follow/unfollow toggle with optional follower count.

## Props

| Prop                  | Type                                                              | Default       | Description                          |
|-----------------------|-------------------------------------------------------------------|---------------|--------------------------------------|
| `following`           | `boolean`                                                         | —             | Controlled following state           |
| `defaultFollowing`    | `boolean`                                                         | `false`       | Initial following state (uncontrolled) |
| `count`               | `number`                                                          | —             | Follower count to display            |
| `disabled`            | `boolean`                                                         | `false`       | Disable interaction                  |
| `followLabel`         | `string`                                                          | `'Follow'`    | Label when not following             |
| `followingLabel`      | `string`                                                          | `'Following'` | Label when following                 |
| `onFollowingChange`   | `(value: boolean, meta: BcChangeMeta<BcFollowChangeReason>) => void` | —          | Fired on toggle (reason: follow/unfollow) |

## Example

```tsx
<BcFollow following={following} count={followers} onFollowingChange={setFollowing} />
```

## Flutter

```dart
BcFollow(following: following, count: followers, onFollowingChange: (v, _) => setFollowing(v))
```
