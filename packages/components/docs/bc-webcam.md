# BcWebcam

Camera capture with getUserMedia (web) and dependency-free delegation
(native/Flutter).

## Contract scope

An open-state overlay that captures a photo from the device camera. Web uses
`getUserMedia` streamed into a `<video>` element; the Capture button draws
the video frame to a canvas and calls `canvas.toDataURL`. Tracks are stopped
on unmount. Permission denial surfaces an error state. Native and Flutter are
dependency-free: the actual camera module is provided by the caller via a
builder.

## Props

| Prop            | Type                                                                                | Default     | Description                          |
|-----------------|-------------------------------------------------------------------------------------|-------------|--------------------------------------|
| `open`          | `boolean`                                                                           | —           | Controlled open state                |
| `defaultOpen`   | `boolean`                                                                           | `false`     | Initial open state (uncontrolled)    |
| `facingMode`    | `'user' \| 'environment'`                                                           | `'user'`    | Camera facing mode                   |
| `captureLabel`  | `string`                                                                            | `'Capture'` | Capture button label                 |
| `adapt`         | `AdaptProps`                                                                        | —           | Responsive/platform morph (ADR-013)  |
| `className`     | `string`                                                                            | —           | Extra class for the root element     |
| `slotProps`     | `Partial<Record<BcWebcamParts, Record<string, unknown>>>`                          | —           | Per-part className/data/aria (ADR-009) |
| `onError`       | `(error: Error) => void`                                                            | —           | Fired on permission/device errors    |
| `onCapture`     | `(dataURL: string, meta: BcChangeMeta<'capture-button'\|'programmatic'>) => void`  | —           | Fired with the captured image dataURL |
| `onOpenChange`  | `(open: boolean, meta: BcChangeMeta<'close-button'\|'overlay-click'\|'escape'\|'programmatic'>) => void` | — | Fired when open state changes |

## A11y

- `role=application`, `aria-label="camera capture"`
- Keyboard: `Escape` closes
- Capture/close buttons have `min 44px` touch targets

## Example

```tsx
<BcWebcam
  open={isOpen}
  onOpenChange={setIsOpen}
  facingMode="user"
  onCapture={(dataURL) => savePhoto(dataURL)}
  onError={(err) => showError(err.message)}
/>
```

## Flutter

```dart
BcWebcam(
  facingMode: BcWebcamFacingMode.user,
  captureLabel: 'Capture',
  cameraBuilder: (ctx) => CameraPreview(controller),  // caller wires camera plugin
  onCapture: (dataURL, meta) => savePhoto(),
)
```

Platform parity: web getUserMedia + canvas; Flutter dependency-free (caller
wires `camera` plugin via `cameraBuilder`).
