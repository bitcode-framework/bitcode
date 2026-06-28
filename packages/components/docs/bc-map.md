# BcMap

Interactive map wrapper with lazy-loaded Leaflet (web) and dependency-free
delegation (native/Flutter).

## Contract scope

Renders an interactive map with markers. Web lazy-loads Leaflet via a
registered loader (`setLeafletLoader(() => import('leaflet'))`). Native and
Flutter are dependency-free: the actual map module is provided by the caller
via a builder, keeping the package free of platform permission manifests.

## Props

| Prop           | Type                                                            | Default   | Description                          |
|----------------|-----------------------------------------------------------------|-----------|--------------------------------------|
| `center`       | `{ lat: number; lng: number }`                                  | `{0,0}`   | Map center coordinate                |
| `zoom`         | `number`                                                        | `13`      | Zoom level                           |
| `markers`      | `{ id, lat, lng, label? }[]`                                    | `[]`      | Marker pins                          |
| `height`       | `string`                                                        | `'400px'` | Container height (CSS)               |
| `adapt`        | `AdaptProps`                                                    | —         | Responsive/platform morph (ADR-013)  |
| `className`    | `string`                                                        | —         | Extra class for the root element     |
| `slotProps`    | `Partial<Record<BcMapParts, Record<string, unknown>>>`         | —         | Per-part className/data/aria (ADR-009) |
| `onMarkerClick`| `(marker: BcMapMarker, meta: BcChangeMeta<'marker-click'\|'programmatic'>) => void` | — | Fired when a marker is clicked |

## A11y

- `role=application`, `aria-label="interactive map"`

## Example

```tsx
import { setLeafletLoader } from '@bitcode/ui/src/components/media/map/BcMap.web';
setLeafletLoader(() => import('leaflet'));

<BcMap
  center={{ lat: 40.7128, lng: -74.0060 }}
  zoom={12}
  markers={[{ id: 'nyc', lat: 40.7128, lng: -74.0060, label: 'NYC' }]}
  onMarkerClick={(m) => console.info(m)}
/>
```

## Flutter

```dart
BcMap(
  center: BcMapLatLng(lat: 40.7128, lng: -74.0060),
  zoom: 12,
  markers: [BcMapMarker(id: 'nyc', lat: 40.7128, lng: -74.0060, label: 'NYC')],
  mapBuilder: (ctx) => GoogleMap(...),  // caller wires google_maps_flutter
)
```

Platform parity: web Leaflet lazy-load; Flutter dependency-free (caller wires
`google_maps_flutter` via `mapBuilder`).
