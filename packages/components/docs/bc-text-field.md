# BcTextField

Rich text input with multiline auto-grow, character counter, and
prefix/suffix slots.

## Contract scope

`BcInput` is the basic variant (single-line, simpler API). `BcTextField` is
the rich variant: multiline auto-grow, character counter, and affixed content
slots.

## Props

| Prop                | Type                                                | Default      | Description                              |
|---------------------|-----------------------------------------------------|--------------|------------------------------------------|
| `name`              | `string`                                            | (required)   | Field identifier                         |
| `label`             | `string`                                            | —            | Field label                              |
| `value`             | `string`                                            | —            | Controlled value                         |
| `defaultValue`      | `string`                                            | —            | Uncontrolled initial value               |
| `placeholder`       | `string`                                            | —            | Placeholder text                         |
| `multiline`         | `boolean`                                           | `false`      | Render as `<textarea>` instead of `<input>` |
| `autoGrow`          | `boolean`                                           | `false`      | Auto-grow textarea (requires `multiline`) |
| `rows`              | `number`                                            | `4`          | Visible rows when multiline              |
| `minRows`           | `number`                                            | `1`          | Min rows for autoGrow                    |
| `maxRows`           | `number`                                            | `8`          | Max rows for autoGrow (overflow scrolls) |
| `showCount`         | `boolean`                                           | `false`      | Show character counter                   |
| `maxLength`         | `number`                                            | —            | Max characters (enforced)                |
| `prefixText`        | `string`                                            | —            | Prefix slot text                         |
| `suffixText`        | `string`                                            | —            | Suffix slot text                         |
| `clearable`         | `boolean`                                           | `false`      | Show clear button when non-empty         |
| `disabled`          | `boolean`                                           | `false`      | Disable the field                        |
| `readonly`          | `boolean`                                           | `false`      | Read-only mode                           |
| `required`          | `boolean`                                           | `false`      | Show required marker in label            |
| `loading`           | `boolean`                                           | `false`      | Show loading spinner                     |
| `autofocus`         | `boolean`                                           | `false`      | Autofocus on mount                       |
| `description`       | `string`                                            | —            | Helper text                              |
| `errorMessage`      | `string`                                            | —            | Error text (shown when invalid)          |
| `validationStatus`  | `'none' \| 'valid' \| 'invalid'`                    | `'none'`     | Border/error token driver                |
| `autoComplete`      | `string`                                            | —            | HTML autocomplete attribute              |
| `onValueChange`     | `(value: string, meta: BcChangeMeta<BcTextFieldChangeReason>) => void` | — | Value change handler |
| `onClear`           | `() => void`                                        | —            | Fired on clear                           |
| `onEnter`           | `() => void`                                        | —            | Fired on Enter (single-line only)        |
| `onFocus` / `onBlur` / `onKeyDown` | `() => void`                          | —            | Standard field events                    |

## A11y

- `role=textbox` with `aria-multiline` set per `multiline`
- `aria-invalid` reflects `validationStatus === 'invalid'`
- `aria-describedby` concatenates description + error + counter ids
- Clear button has `aria-label="Clear"` and >=44px touch target

## Example

```tsx
<BcTextField
  name="bio"
  label="Bio"
  multiline
  autoGrow
  minRows={2}
  maxRows={8}
  showCount
  maxLength={500}
  value={bio}
  onValueChange={setBio}
/>
```

## Flutter

```dart
BcTextField(
  name: 'bio',
  label: 'Bio',
  multiline: true,
  autoGrow: true,
  minRows: 2,
  showCount: true,
  maxLength: 500,
  onValueChange: (v, meta) => setBio(v),
)
```

Platform parity: Material `TextField` with `maxLines: null` for autoGrow,
Cupertino-equivalent styling via the surrounding theme.
