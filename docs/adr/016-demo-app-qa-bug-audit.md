# ADR-016: Demo App QA Bug Audit & Improvement Plan

| Field | Value |
|-------|-------|
| **Status** | Proposed (2026-06-28) |
| **Depends on** | 001, 002, 006, 007, 011 |
| **Impact** | Fixes Mitosis compiler ref bug (charts), overlay portal bug, 25+ stub implementations, Flutter layout crashes |

## Context

A full browser QA pass of both demo apps (`packages/ui/demo-app` for React,
`packages/ui-flutter/example` for Flutter) revealed a large number of broken,
visually broken, or stub-only components. This ADR categorizes every finding
into severity tiers with root cause and fix plan.

## Audit Method

- **React demo**: Vite preview production build at `localhost:3006`, inspected
  via Playwright MCP (DOM snapshots, computed styles, console capture,
  interactive button clicks for overlay testing).
- **Flutter demo**: `flutter run -d chrome` at `localhost:8099`, inspected
  via Playwright MCP + widget test suite (`section_error_test.dart`) that
  taps through every sidebar section and captures `ErrorWidget` exceptions.

---

## React Demo Findings

### CRITICAL — Charts do not render (0 canvas, 0 SVG, 0 echarts instances)

**Root cause**: Mitosis compiler bug — `ref={state.containerRef}` in `.lite.tsx`
compiles to `ref={containerRef}` in React output, but `containerRef` is a
`useState` value (`const [containerRef, setContainerRef] = useState(() => null)`),
not a `useRef`. React's `ref` attribute needs either a `useRef` object (with
`.current`) or a callback ref. The state value is always `null`, so echarts
never initializes.

**Affected**: `BcChart.lite.tsx` and all 23 chart-type wrapper components
(`BcChartBar`, `BcChartLine`, `BcChartPie`, etc.).

**Fix options**:
1. (Preferred) Mitosis compiler fix: when compiling `ref={state.xxx}` to React,
   emit `ref={setXxx}` (callback ref) instead of `ref={xxx}`.
2. (Workaround) Post-build script: regex-replace `ref={containerRef}` →
   `ref={setContainerRef}` in chart output files.
3. (Manual) Rewrite BcChart as a non-Mitosis `.tsx` file using `useRef`
   directly, bypassing the compiler.

### CRITICAL — Overlays render without backdrop/portal

**Root cause**: Overlay components (`BcDialog`, `BcDrawer`, `BcPopover`,
`BcTooltip`) render their `fixed inset-0 z-50` wrapper **inside** the demo
component inline rather than via `ReactDOM.createPortal`. The wrapper inherits
the parent's overflow/bounds, resulting in:
- `offsetHeight: 118px` instead of full viewport height.
- No dark backdrop (`rgba(0,0,0,0)` background).
- Dialog content is visible but appears inline, not centered over the page.

**Fix options**:
1. (Preferred) Add `createPortal` to overlay components in Mitosis source
   (wrap the overlay content in a portal to `document.body`).
2. (Alternative) Fix the demo to wrap overlay children in a portal host.

### HIGH — BcLookupModal: Modal does not open on click

**Symptom**: "Open Contact Lookup" button clicked, but no modal appears.
**Root cause**: Same overlay portal issue as above (dialog renders inline
inside the demo card, invisible to user).

### HIGH — BcCheckbox: Checked state visually invisible

**Symptom**: When unchecked, the checkbox appears empty/invisible. The checked
state shows an SVG checkmark icon but unchecked state has no visible border.
**Root cause**: BcCheckbox uses custom SVG rendering (`4 SVGs` found, `0
input[type=checkbox]`). The unchecked state likely has transparent or
mismatched border color from theme tokens.

### HIGH — BcProgress: 27 zero-height elements

**Symptom**: Progress bars show only text labels, no visible bar/track.
**Root cause**: Progress bar inner elements have `height: 0` due to CSS
(token variable not resolving or flexbox collapse).

### HIGH — BcTabs: 8 empty divs (tab panels not rendering)

**Symptom**: Tab buttons may be present but tab panel content is empty.
**Root cause**: Tab content rendering likely depends on active state that
is not properly initialized in the demo.

### HIGH — BcSplitter: 16 empty divs (panes not rendering)

**Symptom**: Splitter shows labels but no resizable panes or dividers.
**Root cause**: Splitter panes are likely collapsed to zero height/width.

### MEDIUM — BcBadge: Font overlap / layout issues

**Symptom**: Badge text overlaps with container.
**Root cause**: Insufficient padding or line-height in badge styling.

### MEDIUM — BcStepper: Misaligned steps, unclear connector lines

**Symptom**: Step bullets and active/inactive indicators are not vertically
aligned. Connector lines are invisible.
**Root cause**: CSS layout issue in stepper connector element.

### MEDIUM — BcMenu: When expanding a submenu, sibling items shift position

**Symptom**: Clicking "File" shows its children, but "Edit" gets pushed down
instead of the submenu floating over.
**Root cause**: Submenu renders inline (pushes siblings) instead of as a
floating dropdown.

### MEDIUM — BcCalendar: Day cells render vertically or in wrong layout

**Symptom**: Calendar days are not in a 7-column grid; they appear stacked.
**Root cause**: CSS grid/flex not applied correctly to calendar day cells.

### MEDIUM — BcTimeline: Last item position misaligned

**Symptom**: The last timeline item has different positioning/alignment from
previous items. 37 empty divs indicate rendering issues.
**Root cause**: Timeline connector line or dot positioning for last item.

### MEDIUM — BcCarousel: Navigation arrows and slide controls broken

**Symptom**: Carousel slides and navigation arrows are visually broken.
28 overflow-hidden elements suggest clipping issues.

### MEDIUM — BcToast: Click does nothing

**Symptom**: "Show Info Toast" button clicked, nothing appears.
**Root cause**: Toast portal/render issue (toast container not mounted).

### MEDIUM — BcFieldDuration: Layout broken

**Symptom**: Duration input fields are visually misaligned/berantakan.
12 SVGs + 15 overflow-hidden suggest excessive icon clipping.

### MEDIUM — BcCommand / BcDropdown / BcContextMenu: Not working correctly

**Symptom**: These menu-like components don't open or display correctly.
**Root cause**: Same overlay/portal issue as dialogs.

### LOW — BcTopbar: Unclear/vague

**Symptom**: Header renders but layout is unclear about what it should look like.

### STUB — "Coming Soon" placeholders (not yet implemented in Mitosis)

These components show a "Coming Soon" placeholder in the demo. The Stencil
Web Component version exists in `@bitcode/components` but has not been ported
to Mitosis (`.lite.tsx`):

| Component | Section |
|-----------|---------|
| BcTextField | Input (Extended) — "Not Yet Implemented" (BcInput is the equivalent) |
| BcChip | Data Display |
| BcAlert | Feedback |
| BcPrintArea | Print |
| BcPrintPreview | Print |
| BcDetailView | Views |
| BcGridView | Views |
| BcTreeView | Views |
| BcKanbanView | Views |
| BcCalendarView | Views |
| BcDashboard | Views |
| BcLike | Social |
| BcShare | Social |
| BcComment | Social |
| BcFollow | Social |
| BcLightbox | Media |
| BcMap | Media |
| BcWebcam | Media |
| BcNativeButton | Adapt POC |
| BcNativeInput | Adapt POC |
| BcNativeSelect | Adapt POC |
| BcNativeSwitch | Adapt POC |
| BcNativeSlider | Adapt POC |
| BcNativePicker | Adapt POC |
| BcNativeAlert | Adapt POC |
| BcNativeActionSheet | Adapt POC |
| BcNativeModal | Adapt POC |
| BcNativeToast | Adapt POC |

**Total stubs**: 29 components.

### WORKING but could be improved

These components render and function but have minor visual polish issues:

| Component | Note |
|-----------|------|
| BcFileInput | Functional but not styled like typical file inputs (drag-drop zone) |
| BcRating | 87 zero-height elements (star SVGs may be collapsed) |
| BcGantt | 2 zero-height elements, otherwise rendering well |

---

## Flutter Demo Findings

### CRITICAL — Data/Display section crashes: "Vertical viewport was given unbounded height"

**Root cause**: `data_section.dart` line 17 uses `ListView` as its root widget.
When this section is rendered inside `SingleChildScrollView` (added in
`main.dart` to fix the previous RenderFlex overflow), the `ListView` gets
unbounded vertical constraints, causing a crash.

**Fix**: Change `data_section.dart` root from `ListView` to `Column` (since
the parent `SingleChildScrollView` already provides scrolling). Same fix likely
needed for `display_section.dart` if it also uses `ListView`.

### MEDIUM — Overlays section: "ListTile background color or ink splashes may be invisible"

**Root cause**: `ListTile` widgets inside `BcBox` (which has a `DecoratedBox`
with background color) cause Flutter's assertion: ListTile's ink splash is
hidden by the parent's background. ListTile needs to be wrapped in its own
`Material` widget.

**Fix**: Wrap `ListTile` in `Material(color: Colors.transparent, child: ListTile(...))`
or restructure to avoid nesting ListTile inside BcBox with background.

### NOT FULLY REVIEWED

The Flutter demo's Forms, Overlays, Layout, Navigation, and Misc sections
were not individually reviewed due to Flutter web's canvas-based rendering
being difficult to inspect via DOM. The widget test confirms sections 0 (Forms)
and 1 (Overlays) render without `ErrorWidget`, but section 2 (Data/Display)
crashes.

---

## Summary Statistics

### React Demo

| Category | Count |
|----------|-------|
| Critical (charts/overlays) | 2 bug clusters affecting ~30 components |
| High (functional broken) | 6 components |
| Medium (visual bugs) | 10 components |
| Low (minor polish) | 1 component |
| Stubs (Coming Soon) | 29 components |
| **Total issues** | **48** |

### Flutter Demo

| Category | Count |
|----------|-------|
| Critical (crash) | 1 section (Data/Display) |
| Medium (assertion warning) | 1 section (Overlays) |
| Not reviewed | 4 sections (Flutter web canvas inspection limitation) |
| **Total issues** | **2+ (confirmed), up to 6 sections need review** |

---

## Recommended Fix Priority

### Phase 1: Critical (blocks entire feature categories)

1. **Mitosis ref compiler bug** — Fix `ref={state.xxx}` compilation for React
   target. This unblocks ALL chart components (24 components).
2. **Overlay portal** — Add `createPortal` to BcDialog/BcDrawer/BcPopover/
   BcTooltip/BcLookupModal/BcCommand/BcDropdown/BcContextMenu/BcToast.
   This unblocks ALL overlay components (~10 components).
3. **Flutter Data section crash** — Change `ListView` → `Column` in
   `data_section.dart` and `display_section.dart`.

### Phase 2: High (individual component bugs)

4. BcCheckbox unchecked state visibility
5. BcProgress bar rendering (zero-height elements)
6. BcTabs panel content rendering
7. BcSplitter pane rendering

### Phase 3: Medium (visual polish)

8. BcBadge font/layout
9. BcStepper alignment
10. BcMenu submenu floating
11. BcCalendar grid layout
12. BcTimeline last-item alignment
13. BcCarousel navigation
14. BcFieldDuration layout
15. Flutter ListTile ink splash warning

### Phase 4: Stub implementation (longest tail)

16-44. Implement the 29 "Coming Soon" stubs (each requires porting from
   Stencil Web Component to Mitosis `.lite.tsx`).

---

## Decision

Document and categorize all bugs. Phase 1 fixes (Mitosis ref bug + overlay
portal + Flutter crash) should be prioritized as they unblock the largest
number of components.
