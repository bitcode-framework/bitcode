# Stencil-to-Mitosis Parity Matrix

> **Phase A2 deliverable** per the UI Perfection Roadmap.
>
> This document maps every `@Prop` / `@Event` / `@Method` in the Stencil Web Components
> (`packages/components/src/components/`) to its Mitosis equivalent in
> (`packages/ui/src/components/`), and identifies gaps that need to be closed.

## Summary

| Metric | Count |
|--------|-------|
| Stencil components scanned | 127 |
| Mitosis components scanned | 176 |
| Stencil components with Mitosis equivalent | 103 |
| Stencil components with NO Mitosis equivalent | 24 |
| Unique Mitosis components mapped from Stencil | 91 |
| New Mitosis components (no Stencil equivalent) | 85 |
| Stencil coverage (has Mitosis equiv) | 81.1% |

### Group Summary

| Stencil Group | Stencil Count | Mapped | Missing |
|---------------|--------------|--------|---------|
| charts | 26 | 25 | 1 |
| datatable | 5 | 1 | 4 |
| dialogs | 7 | 6 | 1 |
| fields | 36 | 36 | 0 |
| kanban | 5 | 1 | 4 |
| layout | 10 | 6 | 4 |
| media | 8 | 8 | 0 |
| print | 3 | 3 | 0 |
| search | 2 | 2 | 0 |
| social | 3 | 3 | 0 |
| views | 10 | 10 | 0 |
| widgets | 12 | 2 | 10 |

| Mitosis Group | Mitosis Count | Mapped from Stencil | New in Mitosis |
|---------------|--------------|---------------------|----------------|
| charts | 29 | 25 | 4 |
| data | 5 | 5 | 0 |
| data-display | 12 | 0 | 12 |
| feedback | 7 | 0 | 7 |
| forms | 41 | 25 | 16 |
| input | 1 | 1 | 0 |
| layout | 10 | 5 | 5 |
| media | 9 | 8 | 1 |
| navigation | 16 | 2 | 14 |
| overlays | 14 | 5 | 9 |
| primitives | 18 | 2 | 16 |
| print | 3 | 3 | 0 |
| social | 3 | 3 | 0 |
| views | 8 | 7 | 1 |

## Stencil Group: charts

### bc-chart-area

> **Mitosis equivalent:** `BcChartArea`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `smooth` | boolean | `smooth` | boolean | mapped |
| `stacked` | boolean | `stacked` | boolean | mapped |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `showSymbols` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-bar

> **Mitosis equivalent:** `BcChartBar`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `stacked` | boolean | `stacked` | boolean | mapped |
| `horizontal` | boolean | `horizontal` | boolean | mapped |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-boxplot

> **Mitosis equivalent:** `BcChartBoxplot`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-candlestick

> **Mitosis equivalent:** `BcChartCandlestick`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `showMA` | boolean | - | - | missing |
| `maPeriod` | number | - | - | missing |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-custom

> **Mitosis equivalent:** `BcChartCustom`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `option` | string | - | - | missing |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| - | - | `data` | string | ChartData | new in Mitosis |
| - | - | `colors` | string | string[] | new in Mitosis |
| - | - | `chartTitle` | string | new in Mitosis |
| - | - | `legend` | boolean | new in Mitosis |
| - | - | `tooltipEnabled` | boolean | new in Mitosis |
| - | - | `animate` | boolean | new in Mitosis |
| - | - | `toolbox` | boolean | new in Mitosis |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setOption` | - | missing (Mitosis uses callback props, not imperative methods) |
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-funnel

> **Mitosis equivalent:** `BcChartFunnel`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| `sortOrder` | string | - | - | missing |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `sort` | 'ascending' | 'descending' | 'none' | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-gauge

> **Mitosis equivalent:** `BcChartGauge`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `value` | string | `value` | number | mapped |
| `max` | string | `max` | number | mapped |
| `min` | string | `min` | number | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `locale` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `segments` | string | - | - | missing |
| - | - | `data` | string | ChartData | new in Mitosis |
| - | - | `legend` | boolean | new in Mitosis |
| - | - | `tooltipEnabled` | boolean | new in Mitosis |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `thresholds` | Array<{ value: number | new in Mitosis |
| - | - | `color` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | - | missing |
| `lcChartHover` | - | missing |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-graph

> **Mitosis equivalent:** `BcChartGraph`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| `layout` | string | `layout` | 'force' | 'circular' | 'none' | mapped |
| `roam` | boolean | - | - | missing |
| `draggable` | boolean | - | - | missing |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-heatmap

> **Mitosis equivalent:** `BcChartHeatmap`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `visualMapMin` | number | - | - | missing |
| `visualMapMax` | number | - | - | missing |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `xAxis` | string[] | new in Mitosis |
| - | - | `yAxis` | string[] | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-kpi

> **Mitosis equivalent:** `BcKpi`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `value` | string | `value` | string | number | mapped |
| `label` | string | - | - | missing |
| `trend` | string | `trend` | 'up' | 'down' | 'neutral' | mapped |
| `valuePrefix` | string | - | - | missing |
| `valueSuffix` | string | - | - | missing |
| `color` | string | `color` | string | mapped |
| `loading` | boolean | - | - | missing |
| - | - | `title` | string | new in Mitosis |
| - | - | `subtitle` | string | new in Mitosis |
| - | - | `trendValue` | string | new in Mitosis |
| - | - | `icon` | string | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-line

> **Mitosis equivalent:** `BcChartLine`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `smooth` | boolean | `smooth` | boolean | mapped |
| `stacked` | boolean | - | - | missing |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `showArea` | boolean | new in Mitosis |
| - | - | `showSymbols` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-mixed

> **Mitosis equivalent:** `BcChartMixed`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `dualAxis` | boolean | - | - | missing |
| `stacked` | string | `stacked` | boolean | mapped |
| - | - | `horizontal` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-parallel

> **Mitosis equivalent:** `BcChartParallel`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-pictorialbar

> **Mitosis equivalent:** `BcChartPictorialBar`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `symbol` | string | `symbol` | string | mapped |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-pie

> **Mitosis equivalent:** `BcChartPie`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| `donut` | boolean | `donut` | boolean | mapped |
| `roseType` | string | `roseType` | 'radius' | 'area' | false | mapped |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `showLabels` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-pivot

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `data` | string | - | no Mitosis equivalent |
| `rows` | string | - | no Mitosis equivalent |
| `cols` | string | - | no Mitosis equivalent |
| `valueField` | string | - | no Mitosis equivalent |
| `aggFunc` | string | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | no Mitosis equivalent |
| `refresh` | - | no Mitosis equivalent |

### bc-chart-polar

> **Mitosis equivalent:** `BcChartPolar`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| `polarType` | string | - | - | missing |
| `stacked` | boolean | - | - | missing |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-progress

> **Mitosis equivalent:** `BcProgress` (shared with: bc-progress, bc-chart-progress)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `value` | string | `value` | number | mapped |
| `max` | string | `max` | number | mapped |
| `label` | string | `label` | string | mapped |
| `color` | string | - | - | missing |
| `showPercent` | boolean | - | - | missing |
| - | - | `size` | 'sm' | 'md' | 'lg' | new in Mitosis |
| - | - | `variant` | 'default' | 'success' | 'warning' | 'destructive' | new in Mitosis |
| - | - | `showLabel` | boolean | new in Mitosis |
| - | - | `striped` | boolean | new in Mitosis |
| - | - | `animated` | boolean | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-radar

> **Mitosis equivalent:** `BcChartRadar`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `shape` | string | - | - | missing |
| `filled` | boolean | - | - | missing |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `indicators` | Array<{ name: string | new in Mitosis |
| - | - | `max` | number | new in Mitosis |
| - | - | `min` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | - | missing |
| `lcChartHover` | - | missing |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | - | missing |
| `lcError` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-sankey

> **Mitosis equivalent:** `BcChartSankey`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| `orient` | string | `orient` | 'horizontal' | 'vertical' | mapped |
| `draggable` | boolean | - | - | missing |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `nodeAlign` | 'left' | 'justify' | 'right' | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-scatter

> **Mitosis equivalent:** `BcChartScatter`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `bubble` | boolean | - | - | missing |
| `effectScatter` | boolean | - | - | missing |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `showRegression` | boolean | new in Mitosis |
| - | - | `bubbleSize` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-scorecard

> **Mitosis equivalent:** `BcChartScorecard`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `value` | string | - | - | missing |
| `target` | string | - | - | missing |
| `label` | string | - | - | missing |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `locale` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| - | - | `data` | string | ChartData | new in Mitosis |
| - | - | `chartTitle` | string | new in Mitosis |
| - | - | `legend` | boolean | new in Mitosis |
| - | - | `tooltipEnabled` | boolean | new in Mitosis |
| - | - | `loading` | boolean | new in Mitosis |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| - | `onError` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-sunburst

> **Mitosis equivalent:** `BcChartSunburst`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| `innerRadius` | string | - | - | missing |
| `outerRadius` | string | - | - | missing |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-themeriver

> **Mitosis equivalent:** `BcChartThemeRiver`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-tree

> **Mitosis equivalent:** `BcChartTree`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `locale` | string | - | - | missing |
| `layout` | string | - | - | missing |
| `orient` | string | `orient` | 'LR' | 'RL' | 'TB' | 'BT' | mapped |
| - | - | `dataZoom` | boolean | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |
| - | - | `expandAndCollapse` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-chart-treemap

> **Mitosis equivalent:** `BcChartTreemap`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `data` | string | `data` | string | ChartData | mapped |
| `chartTitle` | string | `chartTitle` | string | mapped |
| `colors` | string | `colors` | string | string[] | mapped |
| `legend` | boolean | `legend` | boolean | mapped |
| `tooltipEnabled` | boolean | `tooltipEnabled` | boolean | mapped |
| `animate` | boolean | `animate` | boolean | mapped |
| `height` | string | `height` | string | mapped |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `model` | string | - | - | missing |
| `theme` | string | `theme` | string | mapped |
| `renderer` | string | `renderer` | string | mapped |
| `toolbox` | boolean | `toolbox` | boolean | mapped |
| `dataZoom` | boolean | `dataZoom` | boolean | mapped |
| `locale` | string | - | - | missing |
| `leafDepth` | number | `leafDepth` | number | mapped |
| `roam` | boolean | - | - | missing |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChartClick` | `onChartClick` | mapped |
| `lcChartHover` | `onChartHover` | mapped |
| `lcChartLegendSelect` | - | missing |
| `lcChartDataZoom` | - | missing |
| `lcChartReady` | `onChartReady` | mapped |
| `lcError` | `onError` | mapped |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `updateData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `resize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportImage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChartInstance` | - | missing (Mitosis uses callback props, not imperative methods) |

---

## Stencil Group: datatable

### bc-child-table

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `field` | string | - | no Mitosis equivalent |
| `columns` | string | - | no Mitosis equivalent |
| `data` | string | - | no Mitosis equivalent |
| `summary` | string | - | no Mitosis equivalent |
| `readonly` | boolean | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | no Mitosis equivalent |

### bc-datatable

> **Mitosis equivalent:** `BcDataTable`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `columns` | string | `columns` | ColumnDef[] | mapped |
| `apiUrl` | string | - | - | missing |
| `pageSize` | number | `pageSize` | number | mapped |
| `selectable` | boolean | - | - | missing |
| `draggableColumns` | boolean | - | - | missing |
| `exportXls` | boolean | `exportXls` | boolean | mapped |
| `showFilterBuilder` | boolean | - | - | missing |
| `showJsonFilter` | boolean | - | - | missing |
| `actions` | string | - | - | missing |
| `serverSide` | boolean | - | - | missing |
| `savedPresets` | string | - | - | missing |
| `permissions` | string | `permissions` | BcDataTablePermissions | mapped |
| `createUrl` | string | `createUrl` | string | mapped |
| `editUrl` | string | `editUrl` | string | mapped |
| `detailUrl` | string | `detailUrl` | string | mapped |
| `moduleName` | string | - | - | missing |
| `modalMode` | boolean | `modalMode` | boolean | mapped |
| `formFields` | string | `formFields` | CrudFormField[] | mapped |
| `viewTitle` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `localData` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `emptyText` | string | `emptyText` | string | mapped |
| - | - | `source` | DataSource | new in Mitosis |
| - | - | `data` | unknown[] | new in Mitosis |
| - | - | `striped` | boolean | new in Mitosis |
| - | - | `bordered` | boolean | new in Mitosis |
| - | - | `loading` | boolean | new in Mitosis |
| - | - | `errorText` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcDataTableParts, Record<string, unknown>>> | new in Mitosis |
| - | - | `searchable` | boolean | new in Mitosis |
| - | - | `searchFields` | string[] | new in Mitosis |
| - | - | `multiSort` | boolean | new in Mitosis |
| - | - | `sortModel` | SortRule[] | new in Mitosis |
| - | - | `defaultSortModel` | SortRule[] | new in Mitosis |
| - | - | `filterModel` | FilterRule[] | new in Mitosis |
| - | - | `defaultFilterModel` | FilterRule[] | new in Mitosis |
| - | - | `initialSort` | SortRule[] | new in Mitosis |
| - | - | `initialFilters` | FilterRule[] | new in Mitosis |
| - | - | `selectionMode` | SelectionMode | new in Mitosis |
| - | - | `selectedKeys` | Array<string | number> | new in Mitosis |
| - | - | `defaultSelectedKeys` | Array<string | number> | new in Mitosis |
| - | - | `initialSelected` | Array<string | number> | new in Mitosis |
| - | - | `getRowId` | RowIdAccessor | new in Mitosis |
| - | - | `reorderable` | boolean | new in Mitosis |
| - | - | `resizable` | boolean | new in Mitosis |
| - | - | `showColumnToggle` | boolean | new in Mitosis |
| - | - | `virtualized` | boolean | new in Mitosis |
| - | - | `rowHeight` | number | new in Mitosis |
| - | - | `viewportHeight` | number | new in Mitosis |
| - | - | `overscan` | number | new in Mitosis |
| - | - | `table` | accessor returning child rows. */ | new in Mitosis |
| - | - | `getChildren` | GetChildren | new in Mitosis |
| - | - | `table` | stable row key (required when `getChildren` is set). */ | new in Mitosis |
| - | - | `getTreeKey` | GetRowKey | new in Mitosis |
| - | - | `expandedKeys` | string[] | new in Mitosis |
| - | - | `defaultExpandedKeys` | string[] | new in Mitosis |
| - | - | `initialExpanded` | string[] | new in Mitosis |
| - | - | `groupBy` | string | new in Mitosis |
| - | - | `aggregates` | AggregateSpec[] | new in Mitosis |
| - | - | `showGroupFooter` | boolean | new in Mitosis |
| - | - | `editMode` | EditMode | new in Mitosis |
| - | - | `exportable` | boolean | new in Mitosis |
| - | - | `exportFilename` | string | new in Mitosis |
| - | - | `exportFormat` | 'csv' | 'xlsx' | new in Mitosis |
| - | - | `persistenceStorage` | PersistenceStorage | new in Mitosis |
| - | - | `persistenceKey` | string | new in Mitosis |
| - | - | `showRowContextMenu` | boolean | new in Mitosis |
| - | - | `createButton` | boolean | new in Mitosis |
| - | - | `Default` | "New". */ | new in Mitosis |
| - | - | `createButtonText` | string | new in Mitosis |
| - | - | `use` | id placeholder). */ | new in Mitosis |
| - | - | `formValues` | Record<string, unknown> | new in Mitosis |
| - | - | `showRowActions` | boolean | new in Mitosis |
| - | - | `Default` | ['edit', 'delete']. */ | new in Mitosis |
| - | - | `rowActions` | CrudRowAction[] | new in Mitosis |
| - | - | `total` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcRowClick` | `onRowClick` | mapped |
| `lcSelectionChange` | - | missing |
| `lcBulkAction` | - | missing |
| `lcPageChange` | - | missing |
| `lcSortChange` | - | missing |
| `lcFilterChange` | - | missing |
| - | `onExportData` | new in Mitosis |
| - | `onError` | new in Mitosis |
| - | `onLoaded` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getSelected` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearSelection` | - | missing (Mitosis uses callback props, not imperative methods) |
| `selectAll` | - | missing (Mitosis uses callback props, not imperative methods) |
| `goToPage` | - | missing (Mitosis uses callback props, not imperative methods) |
| `sortBy` | - | missing (Mitosis uses callback props, not imperative methods) |
| `exportCSV` | - | missing (Mitosis uses callback props, not imperative methods) |
| `scrollToRow` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-filter-bar

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `value` | string | - | no Mitosis equivalent |
| `presets` | string | - | no Mitosis equivalent |
| `placeholder` | string | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFilterChange` | - | no Mitosis equivalent |
| `lcSearch` | - | no Mitosis equivalent |

### bc-filter-builder

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `fields` | string | - | no Mitosis equivalent |
| `operators` | string | - | no Mitosis equivalent |
| `value` | string | - | no Mitosis equivalent |
| `showJsonToggle` | boolean | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFilterChange` | - | no Mitosis equivalent |

### bc-filter-panel

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `value` | string | - | no Mitosis equivalent |
| `categories` | string | - | no Mitosis equivalent |
| `placeholder` | string | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFilterChange` | - | no Mitosis equivalent |
| `lcSearch` | - | no Mitosis equivalent |

## Stencil Group: dialogs

### bc-dialog-alert

> **Mitosis equivalent:** `BcAlertDialog` (shared with: bc-dialog-confirm, bc-dialog-alert)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `open` | boolean | `open` | boolean | mapped |
| `dialogTitle` | string | - | - | missing |
| `message` | string | - | - | missing |
| `okText` | string | - | - | missing |
| `variant` | 'default' | 'success' | 'danger' | 'warning' | 'info' | `variant` | 'default' | 'destructive' | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| - | - | `title` | string | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `confirmLabel` | string | new in Mitosis |
| - | - | `cancelLabel` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcAlertDialogParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcDialogConfirm` | - | missing |
| `lcDialogClose` | - | missing |
| - | `onConfirm` | new in Mitosis |
| - | `onCancel` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `openDialog` | - | missing (Mitosis uses callback props, not imperative methods) |
| `closeDialog` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-dialog-confirm

> **Mitosis equivalent:** `BcAlertDialog` (shared with: bc-dialog-confirm, bc-dialog-alert)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `open` | boolean | `open` | boolean | mapped |
| `dialogTitle` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `message` | string | - | - | missing |
| `confirmText` | string | - | - | missing |
| `cancelText` | string | - | - | missing |
| `variant` | 'default' | 'danger' | 'warning' | 'info' | `variant` | 'default' | 'destructive' | mapped |
| - | - | `title` | string | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `confirmLabel` | string | new in Mitosis |
| - | - | `cancelLabel` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcAlertDialogParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcDialogClose` | - | missing |
| `lcDialogConfirm` | - | missing |
| `lcDialogCancel` | - | missing |
| - | `onConfirm` | new in Mitosis |
| - | `onCancel` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `openDialog` | - | missing (Mitosis uses callback props, not imperative methods) |
| `closeDialog` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-dialog-modal

> **Mitosis equivalent:** `BcDialog`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `open` | boolean | `open` | boolean | mapped |
| `dialogTitle` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | 'xl' | `size` | 'sm' | 'md' | 'lg' | 'xl' | 'full' | mapped |
| `loading` | boolean | - | - | missing |
| - | - | `defaultOpen` | boolean | new in Mitosis |
| - | - | `title` | string | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `closable` | boolean | new in Mitosis |
| - | - | `closeOnOverlay` | boolean | new in Mitosis |
| - | - | `closeOnEscape` | boolean | new in Mitosis |
| - | - | `modal` | boolean | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcDialogParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcDialogClose` | - | missing |
| - | `onOpenChange` | new in Mitosis |
| - | `onClose` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `openDialog` | - | missing (Mitosis uses callback props, not imperative methods) |
| `closeDialog` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isOpen` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-dialog-prompt

> **Mitosis equivalent:** `BcPrompt`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `open` | boolean | `open` | boolean | mapped |
| `dialogTitle` | string | - | - | missing |
| `message` | string | - | - | missing |
| `placeholder` | string | `placeholder` | string | mapped |
| `inputValue` | string | - | - | missing |
| `okText` | string | - | - | missing |
| `cancelText` | string | - | - | missing |
| `inputType` | 'text' | 'number' | 'password' | - | - | missing |
| `required` | boolean | `required` | boolean | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| - | - | `title` | string | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `defaultValue` | string | new in Mitosis |
| - | - | `confirmLabel` | string | new in Mitosis |
| - | - | `cancelLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcDialogConfirm` | - | missing |
| `lcDialogCancel` | - | missing |
| `lcDialogClose` | - | missing |
| - | `onConfirm` | new in Mitosis |
| - | `onCancel` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `openDialog` | - | missing (Mitosis uses callback props, not imperative methods) |
| `closeDialog` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-dialog-quickentry

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `open` | boolean | - | no Mitosis equivalent |
| `dialogTitle` | string | - | no Mitosis equivalent |
| `model` | string | - | no Mitosis equivalent |
| `fields` | string | - | no Mitosis equivalent |
| `dataSource` | string | - | no Mitosis equivalent |
| `fetchHeaders` | string | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcDialogClose` | - | no Mitosis equivalent |
| `lcError` | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `openDialog` | - | no Mitosis equivalent |
| `closeDialog` | - | no Mitosis equivalent |

### bc-dialog-wizard

> **Mitosis equivalent:** `BcWizard`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `open` | boolean | `open` | boolean | mapped |
| `dialogTitle` | string | - | - | missing |
| `steps` | string | `steps` | WizardStep[] | mapped |
| - | - | `title` | string | new in Mitosis |
| - | - | `step` | number | new in Mitosis |
| - | - | `defaultStep` | number | new in Mitosis |
| - | - | `canAdvance` | (step: number) => boolean | new in Mitosis |
| - | - | `finishLabel` | string | new in Mitosis |
| - | - | `nextLabel` | string | new in Mitosis |
| - | - | `prevLabel` | string | new in Mitosis |
| - | - | `cancelLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcDialogClose` | - | missing |
| `lcWizardComplete` | - | missing |
| - | `onStepChange` | new in Mitosis |
| - | `onFinish` | new in Mitosis |
| - | `onCancel` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `openDialog` | - | missing (Mitosis uses callback props, not imperative methods) |
| `closeDialog` | - | missing (Mitosis uses callback props, not imperative methods) |
| `goToStep` | - | missing (Mitosis uses callback props, not imperative methods) |
| `nextStep` | - | missing (Mitosis uses callback props, not imperative methods) |
| `prevStep` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getCurrentStep` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-toast

> **Mitosis equivalent:** `BcToast`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `open` | boolean | - | - | missing |
| `dialogTitle` | string | - | - | missing |
| `message` | string | `message` | string | mapped |
| `variant` | string | - | - | missing |
| `duration` | number | `duration` | number | mapped |
| `position` | string | - | - | missing |
| - | - | `type` | 'success' | 'error' | 'warning' | 'info' | new in Mitosis |
| - | - | `visible` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcDialogClose` | - | missing |
| - | `onClose` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `show` | - | missing (Mitosis uses callback props, not imperative methods) |
| `dismiss` | - | missing (Mitosis uses callback props, not imperative methods) |

---

## Stencil Group: fields

### bc-field-barcode

> **Mitosis equivalent:** `BcBarcode`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `format` | string | `format` | string | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `placeholder` | string | `placeholder` | string | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-checkbox

> **Mitosis equivalent:** `BcCheckbox` (shared with: bc-field-checkbox, bc-field-multicheck)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `tooltip` | string | - | - | missing |
| `defaultValue` | boolean | - | - | missing |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `checked` | boolean | new in Mitosis |
| - | - | `defaultChecked` | boolean | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcCheckboxParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onCheckedChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-code

> **Mitosis equivalent:** `BcCode`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | - | - | missing |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `language` | string | `language` | string | mapped |
| `toolbar` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `minLength` | number | - | - | missing |
| `maxLength` | number | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `showCount` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `rows` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-color

> **Mitosis equivalent:** `BcColorPicker`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `required` | boolean | - | - | missing |
| `readonly` | boolean | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `presets` | string[] | new in Mitosis |
| - | - | `showValue` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-currency

> **Mitosis equivalent:** `BcNumber` (shared with: bc-field-integer, bc-field-float, bc-field-currency, bc-field-percent, bc-field-decimal)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | number | `value` | number | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `currency` | string | `currency` | string | mapped |
| `precision` | number | - | - | missing |
| `min` | number | `min` | number | mapped |
| `max` | number | `max` | number | mapped |
| `step` | number | `step` | number | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `prefixText` | string | - | - | missing |
| `suffixText` | string | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | number | `defaultValue` | number | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `format` | 'integer' | 'decimal' | 'currency' | 'percent' | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-date

> **Mitosis equivalent:** `BcDateField`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `min` | string | `min` | string | mapped |
| `max` | string | `max` | string | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| `dependOn` | string | - | - | missing |
| - | - | `locale` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcDateFieldParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-datetime

> **Mitosis equivalent:** `BcDateTime`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `min` | string | `min` | string | mapped |
| `max` | string | `max` | string | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `mode` | 'date' | 'time' | 'datetime' | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-decimal

> **Mitosis equivalent:** `BcNumber` (shared with: bc-field-integer, bc-field-float, bc-field-currency, bc-field-percent, bc-field-decimal)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | number | `value` | number | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `precision` | number | - | - | missing |
| `min` | number | `min` | number | mapped |
| `max` | number | `max` | number | mapped |
| `step` | number | `step` | number | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `prefixText` | string | - | - | missing |
| `suffixText` | string | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | number | `defaultValue` | number | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `format` | 'integer' | 'decimal' | 'currency' | 'percent' | new in Mitosis |
| - | - | `currency` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-duration

> **Mitosis equivalent:** `BcFieldDuration`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | number | `value` | number | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | `size` | 'sm' | 'md' | 'lg' | mapped |
| `clearable` | boolean | `clearable` | boolean | mapped |
| `tooltip` | string | - | - | missing |
| `defaultValue` | number | `defaultValue` | number | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `showSeconds` | boolean | new in Mitosis |
| - | - | `maxHours` | number | new in Mitosis |
| - | - | `step` | number | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcFieldDurationParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-dynlink

> **Mitosis equivalent:** `BcFieldDynLink`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | BcDynLinkValue | mapped |
| `placeholder` | string | - | - | missing |
| `modelfield` | string | - | - | missing |
| `model` | string | - | - | missing |
| `required` | boolean | - | - | missing |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | - | - | missing |
| `options` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | BcDynLinkValue | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `models` | Array<{ model: string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-file

> **Mitosis equivalent:** `BcFile` (shared with: bc-field-image, bc-field-file)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | - | - | missing |
| `accept` | string | `accept` | string | mapped |
| `maxSize` | string | `maxSize` | number | mapped |
| `multiple` | boolean | `multiple` | boolean | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `pathFormat` | string | - | - | missing |
| `nameFormat` | string | - | - | missing |
| `apiBase` | string | - | - | missing |
| `model` | string | - | - | missing |
| `recordId` | string | - | - | missing |
| `fieldName` | string | - | - | missing |
| `preview` | boolean | - | - | missing |
| `showDownload` | boolean | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | - | - | missing |
| `validateOn` | ValidateOn | '' | - | - | missing |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-float

> **Mitosis equivalent:** `BcNumber` (shared with: bc-field-integer, bc-field-float, bc-field-currency, bc-field-percent, bc-field-decimal)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | number | `value` | number | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `precision` | number | - | - | missing |
| `min` | number | `min` | number | mapped |
| `max` | number | `max` | number | mapped |
| `step` | number | `step` | number | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `prefixText` | string | - | - | missing |
| `suffixText` | string | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | number | `defaultValue` | number | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `format` | 'integer' | 'decimal' | 'currency' | 'percent' | new in Mitosis |
| - | - | `currency` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-geo

> **Mitosis equivalent:** `BcGeo`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | - | - | missing |
| `disabled` | boolean | - | - | missing |
| `required` | boolean | - | - | missing |
| `readonly` | boolean | - | - | missing |
| `placeholder` | string | - | - | missing |
| `drawMode` | string | - | - | missing |
| `zoom` | number | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `defaultValue` | string | - | - | missing |
| `validateOn` | ValidateOn | '' | - | - | missing |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-html

> **Mitosis equivalent:** `BcRichText` (shared with: bc-field-markdown, bc-field-html, bc-field-richtext)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `language` | string | - | - | missing |
| `toolbar` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `minLength` | number | - | - | missing |
| `maxLength` | number | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `showCount` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `rows` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-image

> **Mitosis equivalent:** `BcFile` (shared with: bc-field-image, bc-field-file)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | - | - | missing |
| `accept` | string | `accept` | string | mapped |
| `maxSize` | string | `maxSize` | number | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `preview` | boolean | - | - | missing |
| `showDownload` | boolean | - | - | missing |
| `model` | string | - | - | missing |
| `recordId` | string | - | - | missing |
| `fieldName` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | - | - | missing |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `multiple` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-integer

> **Mitosis equivalent:** `BcNumber` (shared with: bc-field-integer, bc-field-float, bc-field-currency, bc-field-percent, bc-field-decimal)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | number | `value` | number | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `min` | number | `min` | number | mapped |
| `max` | number | `max` | number | mapped |
| `step` | number | `step` | number | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `prefixText` | string | - | - | missing |
| `suffixText` | string | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | number | `defaultValue` | number | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `format` | 'integer' | 'decimal' | 'currency' | 'percent' | new in Mitosis |
| - | - | `currency` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-json

> **Mitosis equivalent:** `BcJson`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | - | - | missing |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `language` | string | - | - | missing |
| `toolbar` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `rows` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-link

> **Mitosis equivalent:** `BcCombobox`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `model` | string | - | - | missing |
| `displayField` | string | - | - | missing |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `options` | string | `options` | ComboboxOption[] | mapped |
| `lookupColumns` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| `dependOn` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `open` | boolean | new in Mitosis |
| - | - | `defaultOpen` | boolean | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcComboboxParts, Record<string, unknown>>> | new in Mitosis |
| - | - | `emptyText` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onOpenChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-markdown

> **Mitosis equivalent:** `BcRichText` (shared with: bc-field-markdown, bc-field-html, bc-field-richtext)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `language` | string | - | - | missing |
| `toolbar` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `minLength` | number | - | - | missing |
| `maxLength` | number | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `showCount` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `rows` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-morph

> **Mitosis equivalent:** `BcFieldMorph`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | BcFieldMorphValue | mapped |
| `morphType` | string | - | - | missing |
| `placeholder` | string | `placeholder` | string | mapped |
| `models` | string | - | - | missing |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | `clearable` | boolean | mapped |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | BcFieldMorphValue | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `types` | BcFieldMorphType[] | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcFieldMorphParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-multicheck

> **Mitosis equivalent:** `BcCheckbox` (shared with: bc-field-checkbox, bc-field-multicheck)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | - | - | missing |
| `options` | string | - | - | missing |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `defaultValue` | string | - | - | missing |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `checked` | boolean | new in Mitosis |
| - | - | `defaultChecked` | boolean | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcCheckboxParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onCheckedChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-password

> **Mitosis equivalent:** `BcInput` (shared with: bc-field-password, bc-field-string)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | `validationStatus` | ValidationStatus | mapped |
| `validationMessage` | string | `validationMessage` | string | mapped |
| `hint` | string | `hint` | string | mapped |
| `minLength` | number | `minLength` | number | mapped |
| `maxLength` | number | `maxLength` | number | mapped |
| `pattern` | string | `pattern` | string | mapped |
| `patternMessage` | string | `patternMessage` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | `size` | InputSize | mapped |
| `clearable` | boolean | `clearable` | boolean | mapped |
| `prefixText` | string | `prefixText` | string | mapped |
| `suffixText` | string | `suffixText` | string | mapped |
| `tooltip` | string | `tooltip` | string | mapped |
| `showCount` | boolean | `showCount` | boolean | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `autofocus` | boolean | `autofocus` | boolean | mapped |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | `validateOn` | 'blur' | 'change' | 'submit' | 'manual' | mapped |
| `showReveal` | boolean | - | - | missing |
| - | - | `type` | InputType | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `errorMessage` | string | new in Mitosis |
| - | - | `autoComplete` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `style` | unknown | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcInputParts, Record<string, unknown>>> | new in Mitosis |
| - | - | `source` | DataSource | new in Mitosis |
| - | - | `sourceType` | DataSource['type'] | new in Mitosis |
| - | - | `total` | number | new in Mitosis |
| - | - | `labels` | string[] | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |
| - | `onInput` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |
| - | `onClear` | new in Mitosis |
| - | `onKeyDown` | new in Mitosis |
| - | `onEnter` | new in Mitosis |
| - | `onOptionsResolved` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-percent

> **Mitosis equivalent:** `BcNumber` (shared with: bc-field-integer, bc-field-float, bc-field-currency, bc-field-percent, bc-field-decimal)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | number | `value` | number | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `min` | number | `min` | number | mapped |
| `max` | number | `max` | number | mapped |
| `step` | number | `step` | number | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `prefixText` | string | - | - | missing |
| `suffixText` | string | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | number | `defaultValue` | number | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `format` | 'integer' | 'decimal' | 'currency' | 'percent' | new in Mitosis |
| - | - | `currency` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-radio

> **Mitosis equivalent:** `BcSegmented`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | - | - | missing |
| `value` | string | `value` | string | mapped |
| `options` | string | `options` | BcSegmentedOption[] | mapped |
| `direction` | string | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `required` | boolean | - | - | missing |
| `readonly` | boolean | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | `size` | BcSegmentedSize | mapped |
| `tooltip` | string | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcSegmentedParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-rating

> **Mitosis equivalent:** `BcRate`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | - | - | missing |
| `value` | number | `value` | number | mapped |
| `maxStars` | number | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `required` | boolean | - | - | missing |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | `size` | 'sm' | 'md' | 'lg' | mapped |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `defaultValue` | number | `defaultValue` | number | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `count` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-richtext

> **Mitosis equivalent:** `BcRichText` (shared with: bc-field-markdown, bc-field-html, bc-field-richtext)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `toolbar` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `minLength` | number | - | - | missing |
| `maxLength` | number | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `showCount` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `rows` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-select

> **Mitosis equivalent:** `BcSelect`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `options` | string | `options` | BcSelectOption[] | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | `clearable` | boolean | mapped |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| `dependOn` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `searchable` | boolean | - | - | missing |
| `serverSide` | boolean | - | - | missing |
| `multiple` | boolean | - | - | missing |
| `displayField` | string | - | - | missing |
| `valueField` | string | - | - | missing |
| `groupBy` | string | - | - | missing |
| `creatable` | boolean | - | - | missing |
| `pageSize` | number | - | - | missing |
| `debounceMs` | number | - | - | missing |
| `minSearchLength` | number | - | - | missing |
| `noResultsText` | string | - | - | missing |
| `loadingText` | string | - | - | missing |
| `model` | string | - | - | missing |
| - | - | `open` | boolean | new in Mitosis |
| - | - | `defaultOpen` | boolean | new in Mitosis |
| - | - | `source` | DataSource | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcSelectParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| `lcOptionsLoad` | - | missing |
| `lcOptionsError` | - | missing |
| `lcOptionCreate` | - | missing |
| `lcDropdownOpen` | - | missing |
| `lcDropdownClose` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onOpenChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setOptions` | - | missing (Mitosis uses callback props, not imperative methods) |
| `loadOptions` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reloadOptions` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getOptions` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getSelectedOptions` | - | missing (Mitosis uses callback props, not imperative methods) |
| `open` | - | missing (Mitosis uses callback props, not imperative methods) |
| `close` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-signature

> **Mitosis equivalent:** `BcSignature`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `width` | number | `width` | number | mapped |
| `height` | number | `height` | number | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `tooltip` | string | - | - | missing |
| `defaultValue` | string | - | - | missing |
| `validateOn` | ValidateOn | '' | - | - | missing |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-smalltext

> **Mitosis equivalent:** `BcTextarea` (shared with: bc-field-smalltext, bc-field-text)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `rows` | number | `rows` | number | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `minLength` | number | - | - | missing |
| `maxLength` | number | `maxLength` | number | mapped |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `showCount` | boolean | `showCount` | boolean | mapped |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-string

> **Mitosis equivalent:** `BcInput` (shared with: bc-field-password, bc-field-string)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `max` | number | - | - | missing |
| `widget` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | `validationStatus` | ValidationStatus | mapped |
| `validationMessage` | string | `validationMessage` | string | mapped |
| `hint` | string | `hint` | string | mapped |
| `minLength` | number | `minLength` | number | mapped |
| `maxLength` | number | `maxLength` | number | mapped |
| `pattern` | string | `pattern` | string | mapped |
| `patternMessage` | string | `patternMessage` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | `size` | InputSize | mapped |
| `clearable` | boolean | `clearable` | boolean | mapped |
| `prefixText` | string | `prefixText` | string | mapped |
| `suffixText` | string | `suffixText` | string | mapped |
| `tooltip` | string | `tooltip` | string | mapped |
| `showCount` | boolean | `showCount` | boolean | mapped |
| `loading` | boolean | `loading` | boolean | mapped |
| `autofocus` | boolean | `autofocus` | boolean | mapped |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | `validateOn` | 'blur' | 'change' | 'submit' | 'manual' | mapped |
| `dependOn` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| - | - | `type` | InputType | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `errorMessage` | string | new in Mitosis |
| - | - | `autoComplete` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `style` | unknown | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcInputParts, Record<string, unknown>>> | new in Mitosis |
| - | - | `source` | DataSource | new in Mitosis |
| - | - | `sourceType` | DataSource['type'] | new in Mitosis |
| - | - | `total` | number | new in Mitosis |
| - | - | `labels` | string[] | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |
| - | `onInput` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |
| - | `onClear` | new in Mitosis |
| - | `onKeyDown` | new in Mitosis |
| - | `onEnter` | new in Mitosis |
| - | `onOptionsResolved` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-tableselect

> **Mitosis equivalent:** `BcFieldTableSelect`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | string[] | mapped |
| `placeholder` | string | - | - | missing |
| `model` | string | - | - | missing |
| `required` | boolean | - | - | missing |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | - | - | missing |
| `options` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | string[] | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `columns` | Array<{ key: string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-tags

> **Mitosis equivalent:** `BcTagsInput`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string[] | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `model` | string | - | - | missing |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `options` | string | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string[] | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `maxTags` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-text

> **Mitosis equivalent:** `BcTextarea` (shared with: bc-field-smalltext, bc-field-text)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | `readonly` | boolean | mapped |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `rows` | number | `rows` | number | mapped |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `minLength` | number | - | - | missing |
| `maxLength` | number | `maxLength` | number | mapped |
| `pattern` | string | - | - | missing |
| `patternMessage` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `showCount` | boolean | `showCount` | boolean | mapped |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onFocus` | new in Mitosis |
| - | `onBlur` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-time

> **Mitosis equivalent:** `BcTimeField`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `required` | boolean | `required` | boolean | mapped |
| `readonly` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `min` | string | - | - | missing |
| `max` | string | - | - | missing |
| `step` | number | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | - | - | missing |
| `size` | 'sm' | 'md' | 'lg' | - | - | missing |
| `clearable` | boolean | - | - | missing |
| `tooltip` | string | - | - | missing |
| `loading` | boolean | - | - | missing |
| `autofocus` | boolean | - | - | missing |
| `defaultValue` | string | `defaultValue` | string | mapped |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `HH` | mm` or `HH:mm:ss`. Controlled. */ | new in Mitosis |
| - | - | `withSeconds` | boolean | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcTimeFieldParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-field-toggle

> **Mitosis equivalent:** `BcToggle`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `name` | string | `name` | string | mapped |
| `label` | string | `label` | string | mapped |
| `value` | boolean | - | - | missing |
| `disabled` | boolean | `disabled` | boolean | mapped |
| `required` | boolean | - | - | missing |
| `readonly` | boolean | - | - | missing |
| `validationStatus` | 'none' | 'validating' | 'valid' | 'invalid' | - | - | missing |
| `validationMessage` | string | - | - | missing |
| `hint` | string | `hint` | string | mapped |
| `size` | 'sm' | 'md' | 'lg' | `size` | 'sm' | 'md' | 'lg' | mapped |
| `tooltip` | string | - | - | missing |
| `defaultValue` | boolean | - | - | missing |
| `validateOn` | ValidateOn | '' | - | - | missing |
| - | - | `checked` | boolean | new in Mitosis |
| - | - | `defaultChecked` | boolean | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcToggleParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | missing |
| `lcFieldFocus` | - | missing |
| `lcFieldBlur` | - | missing |
| `lcFieldClear` | - | missing |
| `lcFieldInvalid` | - | missing |
| `lcFieldValid` | - | missing |
| - | `onCheckedChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `validate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `reset` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clear` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `focusField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `blurField` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isDirty` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTouched` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setError` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearError` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-lookup-modal

> **Mitosis equivalent:** `BcLookupModal`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `open` | boolean | `open` | boolean | mapped |
| `model` | string | `model` | string | mapped |
| `displayField` | string | - | - | missing |
| `columns` | string | `columns` | BcLookupModalColumn[] | mapped |
| `multiple` | boolean | `multiple` | boolean | mapped |
| `apiUrl` | string | - | - | missing |
| `modalTitle` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `defaultOpen` | boolean | new in Mitosis |
| - | - | `title` | string | new in Mitosis |
| - | - | `data` | unknown[] | new in Mitosis |
| - | - | `searchable` | boolean | new in Mitosis |
| - | - | `searchFields` | string[] | new in Mitosis |
| - | - | `pageSize` | number | new in Mitosis |
| - | - | `value` | string | string[] | new in Mitosis |
| - | - | `defaultValue` | string | string[] | new in Mitosis |
| - | - | `closeable` | boolean | new in Mitosis |
| - | - | `closeOnOverlay` | boolean | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcLookupModalParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcLookupSelect` | - | missing |
| `lcLookupClose` | - | missing |
| `lcError` | - | missing |
| - | `onValueChange` | new in Mitosis |
| - | `onOpenChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |

---

## Stencil Group: kanban

### bc-kanban-board

> **Mitosis equivalent:** `BcKanbanBoard`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `boardTitle` | string | - | - | missing |
| `groupBy` | string | - | - | missing |
| `cardTitleField` | string | - | - | missing |
| `cardDescriptionField` | string | - | - | missing |
| `cardCoverField` | string | - | - | missing |
| `cardAssigneesField` | string | - | - | missing |
| `cardDueDateField` | string | - | - | missing |
| `cardStartDateField` | string | - | - | missing |
| `cardPriorityField` | string | - | - | missing |
| `cardLabelsField` | string | - | - | missing |
| `cardPositionField` | string | - | - | missing |
| `columnModel` | string | - | - | missing |
| `columnDataSource` | string | - | - | missing |
| `allowAddColumn` | boolean | - | - | missing |
| `allowAddCard` | boolean | - | - | missing |
| `allowRenameColumn` | boolean | - | - | missing |
| `allowDeleteColumn` | boolean | - | - | missing |
| `allowReorderColumns` | boolean | - | - | missing |
| `cardDetailModel` | string | - | - | missing |
| `cardDetailDataSource` | string | - | - | missing |
| `commentModel` | string | - | - | missing |
| `commentDataSource` | string | - | - | missing |
| `subtaskModel` | string | - | - | missing |
| `subtaskDataSource` | string | - | - | missing |
| `attachmentModel` | string | - | - | missing |
| `attachmentDataSource` | string | - | - | missing |
| `activityModel` | string | - | - | missing |
| `activityDataSource` | string | - | - | missing |
| `mentionModel` | string | - | - | missing |
| `mentionDataSource` | string | - | - | missing |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `kanbanCardMove` | - | missing |
| `kanbanColumnReorder` | - | missing |
| `kanbanCardCreate` | - | missing |
| `kanbanCardUpdate` | - | missing |
| `kanbanCardDelete` | - | missing |
| `kanbanColumnAdd` | - | missing |
| `kanbanColumnRename` | - | missing |
| `kanbanColumnDelete` | - | missing |
| `kanbanError` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-kanban-card-activity

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `cardId` | string | - | no Mitosis equivalent |
| `dataSource` | string | - | no Mitosis equivalent |
| `model` | string | - | no Mitosis equivalent |
| `fetchHeaders` | string | - | no Mitosis equivalent |
| `filterBy` | string | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | no Mitosis equivalent |

### bc-kanban-card-attachments

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `cardId` | string | - | no Mitosis equivalent |
| `dataSource` | string | - | no Mitosis equivalent |
| `model` | string | - | no Mitosis equivalent |
| `fetchHeaders` | string | - | no Mitosis equivalent |
| `filterBy` | string | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `kanbanAttachmentUpload` | - | no Mitosis equivalent |
| `kanbanAttachmentDelete` | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | no Mitosis equivalent |

### bc-kanban-card-comments

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `cardId` | string | - | no Mitosis equivalent |
| `dataSource` | string | - | no Mitosis equivalent |
| `model` | string | - | no Mitosis equivalent |
| `fetchHeaders` | string | - | no Mitosis equivalent |
| `filterBy` | string | - | no Mitosis equivalent |
| `mentionModel` | string | - | no Mitosis equivalent |
| `mentionDataSource` | string | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `kanbanCommentCreate` | - | no Mitosis equivalent |
| `kanbanCommentDelete` | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | no Mitosis equivalent |

### bc-kanban-card-subtasks

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `cardId` | string | - | no Mitosis equivalent |
| `dataSource` | string | - | no Mitosis equivalent |
| `model` | string | - | no Mitosis equivalent |
| `fetchHeaders` | string | - | no Mitosis equivalent |
| `filterBy` | string | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `kanbanSubtaskToggle` | - | no Mitosis equivalent |
| `kanbanSubtaskCreate` | - | no Mitosis equivalent |
| `kanbanSubtaskDelete` | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | no Mitosis equivalent |

## Stencil Group: layout

### bc-button-box

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `buttons` | string | - | no Mitosis equivalent |

### bc-column

> **Mitosis equivalent:** `BcColumn`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `width` | number | - | - | missing |
| - | - | `gap` | ComponentSize | 'none' | new in Mitosis |
| - | - | `align` | 'start' | 'center' | 'end' | 'stretch' | 'baseline' | new in Mitosis |
| - | - | `justify` | 'start' | 'center' | 'end' | 'between' | new in Mitosis |
| - | - | `flex` | number | new in Mitosis |

---

### bc-header

> **Mitosis equivalent:** `BcHeader`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `buttons` | string | - | - | missing |
| `statusField` | string | - | - | missing |
| `statusValue` | string | - | - | missing |
| `states` | string | - | - | missing |
| - | - | `title` | string | new in Mitosis |
| - | - | `subtitle` | string | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `size` | 'sm' | 'md' | 'lg' | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcActionClick` | - | missing |

---

### bc-html-block

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `content` | string | - | no Mitosis equivalent |

### bc-row

> **Mitosis equivalent:** `BcRow`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `gap` | string | `gap` | ComponentSize | 'none' | mapped |
| - | - | `align` | 'start' | 'center' | 'end' | 'stretch' | 'baseline' | new in Mitosis |
| - | - | `justify` | 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly' | new in Mitosis |
| - | - | `wrap` | boolean | new in Mitosis |

---

### bc-section

> **Mitosis equivalent:** `BcSection`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `sectionTitle` | string | - | - | missing |
| `description` | string | - | - | missing |
| `collapsible` | boolean | `collapsible` | boolean | mapped |
| `collapsed` | boolean | - | - | missing |
| - | - | `title` | string | new in Mitosis |
| - | - | `subtitle` | string | new in Mitosis |
| - | - | `defaultCollapsed` | boolean | new in Mitosis |
| - | - | `headerAction` | any | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `toggle` | - | missing (Mitosis uses callback props, not imperative methods) |
| `expand` | - | missing (Mitosis uses callback props, not imperative methods) |
| `collapse` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-separator

> **Mitosis equivalent:** `BcSeparator`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `label` | string | - | - | missing |
| - | - | `orientation` | 'horizontal' | 'vertical' | new in Mitosis |
| - | - | `decorative` | boolean | new in Mitosis |

---

### bc-sheet

> **Status:** no Mitosis equivalent

### bc-tab

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `label` | string | - | no Mitosis equivalent |

### bc-tabs

> **Mitosis equivalent:** `BcTabs`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| - | - | `tabs` | TabItem[] | new in Mitosis |
| - | - | `activeIndex` | number | new in Mitosis |
| - | - | `defaultIndex` | number | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcTabsParts, Record<string, unknown>>> | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| - | `onChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `selectTab` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getActiveIndex` | - | missing (Mitosis uses callback props, not imperative methods) |

---

## Stencil Group: media

### bc-viewer-audio

> **Mitosis equivalent:** `BcViewerAudio`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `src` | string | `src` | string | mapped |
| `type` | string | `type` | string | mapped |
| `controls` | boolean | `controls` | boolean | mapped |
| `autoplay` | boolean | `autoplay` | boolean | mapped |
| `loop` | boolean | `loop` | boolean | mapped |
| `download` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `dataSource` | string | - | - | missing |
| `srcField` | string | - | - | missing |
| - | - | `muted` | boolean | new in Mitosis |
| - | - | `title` | string | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-viewer-document

> **Mitosis equivalent:** `BcViewerDocument`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `src` | string | `src` | string | mapped |
| `height` | string | `height` | string | mapped |
| `provider` | 'microsoft' | 'google' | - | - | missing |
| `download` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `dataSource` | string | - | - | missing |
| `srcField` | string | - | - | missing |
| - | - | `engine` | 'office' | 'google' | 'native' | new in Mitosis |
| - | - | `width` | string | new in Mitosis |
| - | - | `title` | string | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-viewer-image

> **Mitosis equivalent:** `BcViewerImage`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `src` | string | `src` | string | mapped |
| `alt` | string | `alt` | string | mapped |
| `width` | string | `width` | string | mapped |
| `height` | string | `height` | string | mapped |
| `zoomable` | boolean | `zoomable` | boolean | mapped |
| `lightbox` | boolean | - | - | missing |
| `download` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `dataSource` | string | - | - | missing |
| `srcField` | string | - | - | missing |
| - | - | `fit` | 'cover' | 'contain' | 'fill' | 'none' | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-viewer-instagram

> **Mitosis equivalent:** `BcViewerInstagram`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `src` | string | - | - | missing |
| `width` | string | `width` | string | mapped |
| `captioned` | boolean | `captioned` | boolean | mapped |
| `loading` | boolean | - | - | missing |
| `dataSource` | string | - | - | missing |
| `srcField` | string | - | - | missing |
| - | - | `url` | string | new in Mitosis |
| - | - | `height` | string | new in Mitosis |
| - | - | `title` | string | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-viewer-pdf

> **Mitosis equivalent:** `BcViewerPdf`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `src` | string | `src` | string | mapped |
| `height` | string | `height` | string | mapped |
| `toolbar` | boolean | `toolbar` | boolean | mapped |
| `download` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `dataSource` | string | - | - | missing |
| `srcField` | string | - | - | missing |
| - | - | `width` | string | new in Mitosis |
| - | - | `title` | string | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-viewer-tiktok

> **Mitosis equivalent:** `BcViewerTiktok`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `src` | string | - | - | missing |
| `width` | string | `width` | string | mapped |
| `loading` | boolean | - | - | missing |
| `dataSource` | string | - | - | missing |
| `srcField` | string | - | - | missing |
| - | - | `videoId` | string | new in Mitosis |
| - | - | `height` | string | new in Mitosis |
| - | - | `title` | string | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-viewer-video

> **Mitosis equivalent:** `BcViewerVideo`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `src` | string | `src` | string | mapped |
| `type` | string | `type` | string | mapped |
| `poster` | string | `poster` | string | mapped |
| `controls` | boolean | `controls` | boolean | mapped |
| `autoplay` | boolean | `autoplay` | boolean | mapped |
| `loop` | boolean | `loop` | boolean | mapped |
| `muted` | boolean | `muted` | boolean | mapped |
| `width` | string | `width` | string | mapped |
| `height` | string | `height` | string | mapped |
| `download` | boolean | - | - | missing |
| `loading` | boolean | - | - | missing |
| `dataSource` | string | - | - | missing |
| `srcField` | string | - | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-viewer-youtube

> **Mitosis equivalent:** `BcViewerYoutube`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `src` | string | - | - | missing |
| `width` | string | - | - | missing |
| `height` | string | - | - | missing |
| `autoplay` | boolean | `autoplay` | boolean | mapped |
| `controls` | boolean | `controls` | boolean | mapped |
| `start` | number | `start` | number | mapped |
| `loading` | boolean | - | - | missing |
| `dataSource` | string | - | - | missing |
| `srcField` | string | - | - | missing |
| - | - | `videoId` | string | new in Mitosis |
| - | - | `ratio` | string | new in Mitosis |
| - | - | `title` | string | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

## Stencil Group: print

### bc-export

> **Mitosis equivalent:** `BcExport`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `label` | string | `label` | string | mapped |
| `href` | string | - | - | missing |
| - | - | `formats` | BcExportFormat[] | new in Mitosis |
| - | - | `disabled` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| - | `onExport` | new in Mitosis |

---

### bc-print

> **Mitosis equivalent:** `BcPrint`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `label` | string | `label` | string | mapped |
| `href` | string | - | - | missing |
| - | - | `target` | string | new in Mitosis |
| - | - | `disabled` | boolean | new in Mitosis |
| - | - | `variant` | 'default' | 'ghost' | 'outline' | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| - | `onBeforePrint` | new in Mitosis |

---

### bc-report-link

> **Mitosis equivalent:** `BcReportLink`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `label` | string | `label` | string | mapped |
| `href` | string | `href` | string | mapped |
| - | - | `format` | string | new in Mitosis |
| - | - | `fileSize` | string | new in Mitosis |
| - | - | `description` | string | new in Mitosis |
| - | - | `external` | boolean | new in Mitosis |

---

## Stencil Group: search

### bc-favorites

> **Mitosis equivalent:** `BcFavorites`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `value` | string | - | - | missing |
| `placeholder` | string | - | - | missing |
| - | - | `items` | FavoriteItem[] | new in Mitosis |
| - | - | `emptyText` | string | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFavoriteSelect` | - | missing |
| `lcSearch` | - | missing |
| - | `onSelect` | new in Mitosis |
| - | `onRemove` | new in Mitosis |

---

### bc-search

> **Mitosis equivalent:** `BcSearch`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `value` | string | `value` | string | mapped |
| `placeholder` | string | `placeholder` | string | mapped |
| `model` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `results` | SearchResult[] | new in Mitosis |
| - | - | `loading` | boolean | new in Mitosis |
| - | - | `emptyText` | string | new in Mitosis |
| - | - | `ariaLabel` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcSearch` | - | missing |
| `lcError` | - | missing |
| - | `onChange` | new in Mitosis |
| - | `onSelect` | new in Mitosis |

---

## Stencil Group: social

### bc-activity

> **Mitosis equivalent:** `BcActivity`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `recordId` | string | - | - | missing |
| `model` | string | - | - | missing |
| - | - | `entries` | ActivityEntry[] | new in Mitosis |
| - | - | `emptyText` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcActivitySchedule` | - | missing |

---

### bc-chatter

> **Mitosis equivalent:** `BcChatter`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `recordId` | string | - | - | missing |
| `model` | string | - | - | missing |
| - | - | `messages` | ChatterMessage[] | new in Mitosis |
| - | - | `placeholder` | string | new in Mitosis |
| - | - | `disabled` | boolean | new in Mitosis |
| - | - | `currentUser` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcChatterSend` | - | missing |
| - | `onSend` | new in Mitosis |

---

### bc-timeline

> **Mitosis equivalent:** `BcTimelineSocial`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `recordId` | string | - | - | missing |
| `model` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `posts` | SocialPost[] | new in Mitosis |
| - | - | `emptyText` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcError` | - | missing |
| - | `onLike` | new in Mitosis |
| - | `onComment` | new in Mitosis |
| - | `onShare` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

## Stencil Group: views

### bc-view-activity

> **Mitosis equivalent:** `BcViewActivity`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `fields` | string | - | - | missing |
| `config` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `entries` | ViewActivityEntry[] | new in Mitosis |
| - | - | `groupByDay` | boolean | new in Mitosis |
| - | - | `emptyText` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcError` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-calendar

> **Mitosis equivalent:** `BcCalendar`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `fields` | string | - | - | missing |
| `config` | string | - | - | missing |
| `dateField` | string | - | - | missing |
| `titleField` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `events` | CalendarEvent[] | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcError` | - | missing |
| - | `onDateClick` | new in Mitosis |
| - | `onEventClick` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-editor

> **Mitosis equivalent:** `BcViewEditor`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `viewJson` | string | - | - | missing |
| `modelFields` | string | - | - | missing |
| `readonly` | boolean | - | - | missing |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `viewChanged` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-form

> **Mitosis equivalent:** `BcForm`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `recordId` | string | - | - | missing |
| `fields` | string | `fields` | FormFieldDef[] | mapped |
| `config` | string | - | - | missing |
| `permissions` | string | - | - | missing |
| `moduleName` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `values` | Record<string, unknown> | new in Mitosis |
| - | - | `layout` | 'vertical' | 'horizontal' | new in Mitosis |
| - | - | `submitLabel` | string | new in Mitosis |
| - | - | `disabled` | boolean | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcFormParts, Record<string, unknown>>> | new in Mitosis |
| - | - | `hideSubmit` | boolean | new in Mitosis |
| - | - | `it` | cascading re-queries, conditional visibility / | new in Mitosis |
| - | - | `reactions` | Reaction[] | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFormSubmit` | - | missing |
| `lcError` | - | missing |
| - | `onReactiveChange` | new in Mitosis |
| - | `onChange` | new in Mitosis |
| - | `onSubmit` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-gantt

> **Mitosis equivalent:** `BcViewGantt`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `tasks` | string | `tasks` | BcGanttTask[] | mapped |
| `links` | string | `links` | BcGanttLink[] | mapped |
| `columns` | string | - | - | missing |
| `scales` | string | - | - | missing |
| `markers` | string | - | - | missing |
| `config` | string | - | - | missing |
| `templates` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `height` | string | - | - | missing |
| `model` | string | `model` | string | mapped |
| `fields` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| `loading` | boolean | - | - | missing |
| `readonly` | boolean | - | - | missing |
| `fieldMapping` | string | - | - | missing |
| - | - | `timelineStart` | string | new in Mitosis |
| - | - | `timelineEnd` | string | new in Mitosis |
| - | - | `zoomLevel` | BcGanttZoomLevel | new in Mitosis |
| - | - | `unitWidth` | number | new in Mitosis |
| - | - | `rowHeight` | number | new in Mitosis |
| - | - | `showCriticalPath` | boolean | new in Mitosis |
| - | - | `showTodayLine` | boolean | new in Mitosis |
| - | - | `showDependencies` | boolean | new in Mitosis |
| - | - | `showProgress` | boolean | new in Mitosis |
| - | - | `showAssignees` | boolean | new in Mitosis |
| - | - | `allowTaskDrag` | boolean | new in Mitosis |
| - | - | `allowTaskResize` | boolean | new in Mitosis |
| - | - | `allowDependencyCreate` | boolean | new in Mitosis |
| - | - | `taskNameField` | string | new in Mitosis |
| - | - | `taskStartField` | string | new in Mitosis |
| - | - | `taskEndField` | string | new in Mitosis |
| - | - | `taskDurationField` | string | new in Mitosis |
| - | - | `taskProgressField` | string | new in Mitosis |
| - | - | `taskDependencyField` | string | new in Mitosis |
| - | - | `taskParentField` | string | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcGanttTaskClick` | - | missing |
| `lcGanttTaskDblClick` | - | missing |
| `lcGanttTaskDrag` | - | missing |
| `lcGanttTaskAdd` | - | missing |
| `lcGanttTaskUpdate` | - | missing |
| `lcGanttTaskDelete` | - | missing |
| `lcGanttTaskSelect` | - | missing |
| `lcGanttTaskUnselect` | - | missing |
| `lcGanttTaskMove` | - | missing |
| `lcGanttTaskRowClick` | - | missing |
| `lcGanttLinkClick` | - | missing |
| `lcGanttLinkDblClick` | - | missing |
| `lcGanttLinkAdd` | - | missing |
| `lcGanttLinkUpdate` | - | missing |
| `lcGanttLinkDelete` | - | missing |
| `lcGanttLinkDrag` | - | missing |
| `lcGanttScaleClick` | - | missing |
| `lcGanttScroll` | - | missing |
| `lcGanttColumnResize` | - | missing |
| `lcGanttRowResize` | - | missing |
| `lcGanttContextMenu` | - | missing |
| `lcGanttEmptyClick` | - | missing |
| `lcGanttReady` | - | missing |
| `lcGanttRendered` | - | missing |
| `lcGanttParsed` | - | missing |
| `lcGanttError` | - | missing |
| `lcGanttTaskOpened` | - | missing |
| `lcGanttTaskClosed` | - | missing |
| `lcGanttBeforeTaskDrag` | - | missing |
| `lcGanttAfterTaskDrag` | - | missing |
| `lcGanttBeforeTaskAdd` | - | missing |
| `lcGanttBeforeTaskDelete` | - | missing |
| `lcGanttBeforeTaskUpdate` | - | missing |
| `lcGanttBeforeTaskMove` | - | missing |
| `lcGanttBeforeLinkAdd` | - | missing |
| `lcGanttBeforeLinkDelete` | - | missing |
| `lcGanttBeforeLinkUpdate` | - | missing |
| - | `onTaskClick` | new in Mitosis |
| - | `onTaskMove` | new in Mitosis |
| - | `onTaskResize` | new in Mitosis |
| - | `onTaskProgressChange` | new in Mitosis |
| - | `onDependencyCreate` | new in Mitosis |
| - | `onZoomChange` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `parseData` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getTasks` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getTaskCount` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getVisibleTaskCount` | - | missing (Mitosis uses callback props, not imperative methods) |
| `addTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `updateTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `deleteTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `moveTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getLink` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getLinks` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getLinkCount` | - | missing (Mitosis uses callback props, not imperative methods) |
| `addLink` | - | missing (Mitosis uses callback props, not imperative methods) |
| `updateLink` | - | missing (Mitosis uses callback props, not imperative methods) |
| `deleteLink` | - | missing (Mitosis uses callback props, not imperative methods) |
| `selectTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `unselectTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getSelectedId` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isSelected` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTaskExists` | - | missing (Mitosis uses callback props, not imperative methods) |
| `isTaskVisible` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getChildren` | - | missing (Mitosis uses callback props, not imperative methods) |
| `hasChild` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getParent` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getSiblings` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getNext` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getPrev` | - | missing (Mitosis uses callback props, not imperative methods) |
| `open` | - | missing (Mitosis uses callback props, not imperative methods) |
| `close` | - | missing (Mitosis uses callback props, not imperative methods) |
| `openAll` | - | missing (Mitosis uses callback props, not imperative methods) |
| `closeAll` | - | missing (Mitosis uses callback props, not imperative methods) |
| `scrollTo` | - | missing (Mitosis uses callback props, not imperative methods) |
| `showTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `showDate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getScrollState` | - | missing (Mitosis uses callback props, not imperative methods) |
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |
| `rerender` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearAll` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getState` | - | missing (Mitosis uses callback props, not imperative methods) |
| `addMarker` | - | missing (Mitosis uses callback props, not imperative methods) |
| `deleteMarker` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getMarker` | - | missing (Mitosis uses callback props, not imperative methods) |
| `undo` | - | missing (Mitosis uses callback props, not imperative methods) |
| `redo` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getUndoStack` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getRedoStack` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearUndoStack` | - | missing (Mitosis uses callback props, not imperative methods) |
| `clearRedoStack` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getTaskIndex` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getGlobalTaskIndex` | - | missing (Mitosis uses callback props, not imperative methods) |
| `calculateDuration` | - | missing (Mitosis uses callback props, not imperative methods) |
| `calculateEndDate` | - | missing (Mitosis uses callback props, not imperative methods) |
| `sortBy` | - | missing (Mitosis uses callback props, not imperative methods) |
| `eachTask` | - | missing (Mitosis uses callback props, not imperative methods) |
| `serialize` | - | missing (Mitosis uses callback props, not imperative methods) |
| `setSizes` | - | missing (Mitosis uses callback props, not imperative methods) |
| `expand` | - | missing (Mitosis uses callback props, not imperative methods) |
| `collapse` | - | missing (Mitosis uses callback props, not imperative methods) |
| `dateFromPos` | - | missing (Mitosis uses callback props, not imperative methods) |
| `posFromDate` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-kanban

> **Mitosis equivalent:** `BcKanban`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `fields` | string | - | - | missing |
| `config` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `columns` | KanbanColumn[] | new in Mitosis |
| - | - | `cards` | KanbanCard[] | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcKanbanMove` | - | missing |
| `lcKanbanColumnReorder` | - | missing |
| `lcError` | - | missing |
| - | `onCardMove` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-list

> **Mitosis equivalent:** `BcList`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `fields` | string | - | - | missing |
| `config` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| `refreshInterval` | number | - | - | missing |
| - | - | `source` | DataSource | new in Mitosis |
| - | - | `data` | unknown[] | new in Mitosis |
| - | - | `titleField` | string | new in Mitosis |
| - | - | `subtitleField` | string | new in Mitosis |
| - | - | `metaField` | string | new in Mitosis |
| - | - | `pageSize` | number | new in Mitosis |
| - | - | `bordered` | boolean | new in Mitosis |
| - | - | `emptyText` | string | new in Mitosis |
| - | - | `loadingText` | string | new in Mitosis |
| - | - | `className` | string | new in Mitosis |
| - | - | `slotProps` | Partial<Record<keyof BcListParts, Record<string, unknown>>> | new in Mitosis |
| - | - | `total` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcRowSelect` | - | missing |
| `lcError` | `onError` | mapped |
| - | `onItemClick` | new in Mitosis |
| - | `onLoaded` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-map

> **Mitosis equivalent:** `BcViewMap`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `fields` | string | - | - | missing |
| `config` | string | - | - | missing |
| `geoField` | string | - | - | missing |
| `nameField` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `lat` | number | new in Mitosis |
| - | - | `lng` | number | new in Mitosis |
| - | - | `zoom` | number | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcError` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-report

> **Mitosis equivalent:** `BcViewReport`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `fields` | string | - | - | missing |
| `config` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `title` | string | new in Mitosis |
| - | - | `subtitle` | string | new in Mitosis |
| - | - | `sections` | ReportSection[] | new in Mitosis |
| - | - | `generatedAt` | string | new in Mitosis |
| - | - | `printable` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcError` | - | missing |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-view-tree

> **Mitosis equivalent:** `BcTree`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `model` | string | - | - | missing |
| `viewTitle` | string | - | - | missing |
| `fields` | string | - | - | missing |
| `config` | string | - | - | missing |
| `parentField` | string | - | - | missing |
| `dataSource` | string | - | - | missing |
| `fetchHeaders` | string | - | - | missing |
| - | - | `data` | TreeNode[] | new in Mitosis |
| - | - | `selectable` | boolean | new in Mitosis |

#### Events

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcError` | - | missing |
| - | `onSelect` | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refresh` | - | missing (Mitosis uses callback props, not imperative methods) |

---

## Stencil Group: widgets

### bc-badge

> **Mitosis equivalent:** `BcBadge`

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `value` | string | - | - | missing |
| `variant` | string | `variant` | ComponentVariant | mapped |
| - | - | `children` | string | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-copy

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `value` | string | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setValue` | - | no Mitosis equivalent |
| `getValue` | - | no Mitosis equivalent |

### bc-domain

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `value` | string | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setValue` | - | no Mitosis equivalent |
| `getValue` | - | no Mitosis equivalent |

### bc-email

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `value` | string | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setValue` | - | no Mitosis equivalent |
| `getValue` | - | no Mitosis equivalent |

### bc-handle

> **Status:** no Mitosis equivalent

### bc-phone

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `value` | string | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setValue` | - | no Mitosis equivalent |
| `getValue` | - | no Mitosis equivalent |

### bc-placeholder

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `text` | string | - | no Mitosis equivalent |

### bc-priority

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `value` | number | - | no Mitosis equivalent |
| `max` | number | - | no Mitosis equivalent |
| `disabled` | boolean | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `lcFieldChange` | - | no Mitosis equivalent |

### bc-progress

> **Mitosis equivalent:** `BcProgress` (shared with: bc-progress, bc-chart-progress)

#### Props

| Stencil Prop | Stencil Type | Mitosis Prop | Mitosis Type | Status |
|-------------|-------------|--------------|--------------|--------|
| `value` | number | `value` | number | mapped |
| `max` | number | `max` | number | mapped |
| `variant` | string | `variant` | 'default' | 'success' | 'warning' | 'destructive' | mapped |
| - | - | `size` | 'sm' | 'md' | 'lg' | new in Mitosis |
| - | - | `showLabel` | boolean | new in Mitosis |
| - | - | `label` | string | new in Mitosis |
| - | - | `striped` | boolean | new in Mitosis |
| - | - | `animated` | boolean | new in Mitosis |

#### Methods

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setValue` | - | missing (Mitosis uses callback props, not imperative methods) |
| `getValue` | - | missing (Mitosis uses callback props, not imperative methods) |

---

### bc-statusbar

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `states` | string | - | no Mitosis equivalent |
| `value` | string | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setValue` | - | no Mitosis equivalent |
| `getValue` | - | no Mitosis equivalent |

### bc-sync-status

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `pollInterval` | number | - | no Mitosis equivalent |
| `compact` | boolean | - | no Mitosis equivalent |
| `showSyncButton` | boolean | - | no Mitosis equivalent |

| Stencil Event | Mitosis Equivalent | Status |
|---------------|---------------------|--------|
| `bcSyncTriggered` | - | no Mitosis equivalent |
| `bcSyncCompleted` | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `refreshStatus` | - | no Mitosis equivalent |
| `triggerSync` | - | no Mitosis equivalent |

### bc-url

> **Status:** no Mitosis equivalent

| Stencil Prop | Type | Mitosis Equivalent | Status |
|-------------|------|---------------------|--------|
| `value` | string | - | no Mitosis equivalent |

| Stencil Method | Mitosis Equivalent | Status |
|----------------|---------------------|--------|
| `setValue` | - | no Mitosis equivalent |
| `getValue` | - | no Mitosis equivalent |

## New in Mitosis (No Stencil Equivalent)

These Mitosis components have no direct Stencil counterpart. They are new additions
introduced by the Mitosis-based UI library.

| Mitosis Component | Group | Props | Events |
|-------------------|-------|-------|--------|
| `BcChart` | charts | 19 | 4 |
| `BcChartLegend` | charts | 3 | 1 |
| `BcChartTooltip` | charts | 5 | 0 |
| `BcSparkline` | charts | 7 | 0 |
| `BcCarousel` | data-display | 5 | 1 |
| `BcDescriptions` | data-display | 5 | 0 |
| `BcDisclosure` | data-display | 6 | 1 |
| `BcDisclosureGroup` | data-display | 1 | 0 |
| `BcImage` | data-display | 7 | 0 |
| `BcInplace` | data-display | 8 | 1 |
| `BcOrgChart` | data-display | 2 | 1 |
| `BcQRCode` | data-display | 3 | 0 |
| `BcStatistic` | data-display | 6 | 0 |
| `BcTagGroup` | data-display | 3 | 1 |
| `BcTimeline` | data-display | 1 | 0 |
| `BcWatermark` | data-display | 5 | 0 |
| `BcBackdrop` | feedback | 7 | 1 |
| `BcEmpty` | feedback | 3 | 0 |
| `BcLoadingOverlay` | feedback | 7 | 0 |
| `BcMeterGroup` | feedback | 3 | 0 |
| `BcNotification` | feedback | 6 | 2 |
| `BcResult` | feedback | 4 | 0 |
| `BcTour` | feedback | 7 | 3 |
| `BcCascader` | forms | 8 | 1 |
| `BcColorArea` | forms | 5 | 1 |
| `BcColorSlider` | forms | 6 | 1 |
| `BcColorSwatch` | forms | 4 | 1 |
| `BcColorWheel` | forms | 4 | 1 |
| `BcDropZone` | forms | 8 | 1 |
| `BcKnob` | forms | 6 | 1 |
| `BcMentions` | forms | 6 | 2 |
| `BcMultiSelect` | forms | 14 | 3 |
| `BcNumberField` | forms | 14 | 2 |
| `BcOtp` | forms | 6 | 2 |
| `BcRangeSlider` | forms | 10 | 1 |
| `BcSlider` | forms | 10 | 1 |
| `BcToggleGroup` | forms | 7 | 1 |
| `BcTransfer` | forms | 5 | 1 |
| `BcTreeSelect` | forms | 8 | 1 |
| `BcAccordion` | layout | 5 | 1 |
| `BcFlex` | layout | 7 | 0 |
| `BcSidebar` | layout | 5 | 1 |
| `BcSplitter` | layout | 9 | 1 |
| `BcStack` | layout | 1 | 0 |
| `BcGalleria` | media | 3 | 1 |
| `BcAffix` | navigation | 2 | 1 |
| `BcAnchor` | navigation | 3 | 1 |
| `BcBackTop` | navigation | 2 | 1 |
| `BcBreadcrumb` | navigation | 2 | 1 |
| `BcDock` | navigation | 3 | 1 |
| `BcDropdown` | navigation | 6 | 1 |
| `BcMegaMenu` | navigation | 2 | 1 |
| `BcMenubar` | navigation | 2 | 1 |
| `BcNavigationMenu` | navigation | 2 | 1 |
| `BcPagination` | navigation | 6 | 1 |
| `BcPanelMenu` | navigation | 2 | 2 |
| `BcSlideMenu` | navigation | 2 | 3 |
| `BcSteps` | navigation | 4 | 1 |
| `BcTieredMenu` | navigation | 2 | 2 |
| `BcActionsheet` | overlays | 5 | 2 |
| `BcCommand` | overlays | 6 | 2 |
| `BcContextMenu` | overlays | 4 | 1 |
| `BcDrawer` | overlays | 9 | 1 |
| `BcHoverCard` | overlays | 6 | 1 |
| `BcPopover` | overlays | 7 | 2 |
| `BcSpeedDial` | overlays | 5 | 2 |
| `BcToastQueue` | overlays | 5 | 1 |
| `BcTooltip` | overlays | 5 | 0 |
| `BcAspectRatio` | primitives | 1 | 0 |
| `BcAvatar` | primitives | 4 | 0 |
| `BcBox` | primitives | 4 | 0 |
| `BcButton` | primitives | 5 | 1 |
| `BcCard` | primitives | 0 | 0 |
| `BcCenter` | primitives | 2 | 0 |
| `BcCircle` | primitives | 2 | 0 |
| `BcKbd` | primitives | 3 | 0 |
| `BcLabel` | primitives | 4 | 0 |
| `BcSkeleton` | primitives | 3 | 0 |
| `BcSlot` | primitives | 3 | 0 |
| `BcSpinner` | primitives | 1 | 0 |
| `BcSquare` | primitives | 2 | 0 |
| `BcVisuallyHidden` | primitives | 4 | 0 |
| `BcXStack` | primitives | 6 | 0 |
| `BcYStack` | primitives | 6 | 0 |
| `BcVirtualScroller` | views | 4 | 0 |

## Gap Summary

### Missing Stencil Components (no Mitosis equivalent)

These Stencil components have no corresponding Mitosis component:

| Stencil Component | Group | Props | Events | Methods |
|-------------------|-------|-------|--------|---------|
| `bc-chart-pivot` | charts | 5 | 0 | 2 |
| `bc-child-table` | datatable | 5 | 1 | 0 |
| `bc-filter-bar` | datatable | 3 | 2 | 0 |
| `bc-filter-builder` | datatable | 4 | 1 | 0 |
| `bc-filter-panel` | datatable | 3 | 2 | 0 |
| `bc-dialog-quickentry` | dialogs | 6 | 2 | 2 |
| `bc-kanban-card-activity` | kanban | 5 | 0 | 1 |
| `bc-kanban-card-attachments` | kanban | 5 | 2 | 1 |
| `bc-kanban-card-comments` | kanban | 7 | 2 | 1 |
| `bc-kanban-card-subtasks` | kanban | 5 | 3 | 1 |
| `bc-button-box` | layout | 1 | 0 | 0 |
| `bc-html-block` | layout | 1 | 0 | 0 |
| `bc-sheet` | layout | 0 | 0 | 0 |
| `bc-tab` | layout | 1 | 0 | 0 |
| `bc-copy` | widgets | 1 | 0 | 2 |
| `bc-domain` | widgets | 1 | 0 | 2 |
| `bc-email` | widgets | 1 | 0 | 2 |
| `bc-handle` | widgets | 0 | 0 | 0 |
| `bc-phone` | widgets | 1 | 0 | 2 |
| `bc-placeholder` | widgets | 1 | 0 | 0 |
| `bc-priority` | widgets | 3 | 1 | 0 |
| `bc-statusbar` | widgets | 2 | 0 | 2 |
| `bc-sync-status` | widgets | 3 | 2 | 2 |
| `bc-url` | widgets | 1 | 0 | 2 |

### Missing Props in Mitosis (mapped components)

| Stencil Component | Mitosis Component | Missing Prop |
|-------------------|-------------------|--------------|
| `bc-activity` | `BcActivity` | `model` |
| `bc-activity` | `BcActivity` | `recordId` |
| `bc-badge` | `BcBadge` | `value` |
| `bc-chart-area` | `BcChartArea` | `dataSource` |
| `bc-chart-area` | `BcChartArea` | `fetchHeaders` |
| `bc-chart-area` | `BcChartArea` | `locale` |
| `bc-chart-area` | `BcChartArea` | `model` |
| `bc-chart-area` | `BcChartArea` | `refreshInterval` |
| `bc-chart-bar` | `BcChartBar` | `dataSource` |
| `bc-chart-bar` | `BcChartBar` | `fetchHeaders` |
| `bc-chart-bar` | `BcChartBar` | `locale` |
| `bc-chart-bar` | `BcChartBar` | `model` |
| `bc-chart-bar` | `BcChartBar` | `refreshInterval` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `dataSource` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `fetchHeaders` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `locale` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `model` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `refreshInterval` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `dataSource` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `fetchHeaders` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `locale` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `maPeriod` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `model` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `refreshInterval` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `showMA` |
| `bc-chart-custom` | `BcChartCustom` | `dataSource` |
| `bc-chart-custom` | `BcChartCustom` | `fetchHeaders` |
| `bc-chart-custom` | `BcChartCustom` | `model` |
| `bc-chart-custom` | `BcChartCustom` | `option` |
| `bc-chart-custom` | `BcChartCustom` | `refreshInterval` |
| `bc-chart-funnel` | `BcChartFunnel` | `dataSource` |
| `bc-chart-funnel` | `BcChartFunnel` | `fetchHeaders` |
| `bc-chart-funnel` | `BcChartFunnel` | `locale` |
| `bc-chart-funnel` | `BcChartFunnel` | `model` |
| `bc-chart-funnel` | `BcChartFunnel` | `refreshInterval` |
| `bc-chart-funnel` | `BcChartFunnel` | `sortOrder` |
| `bc-chart-gauge` | `BcChartGauge` | `locale` |
| `bc-chart-gauge` | `BcChartGauge` | `segments` |
| `bc-chart-graph` | `BcChartGraph` | `dataSource` |
| `bc-chart-graph` | `BcChartGraph` | `draggable` |
| `bc-chart-graph` | `BcChartGraph` | `fetchHeaders` |
| `bc-chart-graph` | `BcChartGraph` | `locale` |
| `bc-chart-graph` | `BcChartGraph` | `model` |
| `bc-chart-graph` | `BcChartGraph` | `refreshInterval` |
| `bc-chart-graph` | `BcChartGraph` | `roam` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `dataSource` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `fetchHeaders` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `locale` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `model` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `refreshInterval` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `visualMapMax` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `visualMapMin` |
| `bc-chart-kpi` | `BcKpi` | `label` |
| `bc-chart-kpi` | `BcKpi` | `loading` |
| `bc-chart-kpi` | `BcKpi` | `valuePrefix` |
| `bc-chart-kpi` | `BcKpi` | `valueSuffix` |
| `bc-chart-line` | `BcChartLine` | `dataSource` |
| `bc-chart-line` | `BcChartLine` | `fetchHeaders` |
| `bc-chart-line` | `BcChartLine` | `locale` |
| `bc-chart-line` | `BcChartLine` | `model` |
| `bc-chart-line` | `BcChartLine` | `refreshInterval` |
| `bc-chart-line` | `BcChartLine` | `stacked` |
| `bc-chart-mixed` | `BcChartMixed` | `dataSource` |
| `bc-chart-mixed` | `BcChartMixed` | `dualAxis` |
| `bc-chart-mixed` | `BcChartMixed` | `fetchHeaders` |
| `bc-chart-mixed` | `BcChartMixed` | `locale` |
| `bc-chart-mixed` | `BcChartMixed` | `model` |
| `bc-chart-mixed` | `BcChartMixed` | `refreshInterval` |
| `bc-chart-parallel` | `BcChartParallel` | `dataSource` |
| `bc-chart-parallel` | `BcChartParallel` | `fetchHeaders` |
| `bc-chart-parallel` | `BcChartParallel` | `locale` |
| `bc-chart-parallel` | `BcChartParallel` | `model` |
| `bc-chart-parallel` | `BcChartParallel` | `refreshInterval` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `dataSource` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `fetchHeaders` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `locale` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `model` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `refreshInterval` |
| `bc-chart-pie` | `BcChartPie` | `dataSource` |
| `bc-chart-pie` | `BcChartPie` | `fetchHeaders` |
| `bc-chart-pie` | `BcChartPie` | `locale` |
| `bc-chart-pie` | `BcChartPie` | `model` |
| `bc-chart-pie` | `BcChartPie` | `refreshInterval` |
| `bc-chart-polar` | `BcChartPolar` | `dataSource` |
| `bc-chart-polar` | `BcChartPolar` | `fetchHeaders` |
| `bc-chart-polar` | `BcChartPolar` | `locale` |
| `bc-chart-polar` | `BcChartPolar` | `model` |
| `bc-chart-polar` | `BcChartPolar` | `polarType` |
| `bc-chart-polar` | `BcChartPolar` | `refreshInterval` |
| `bc-chart-polar` | `BcChartPolar` | `stacked` |
| `bc-chart-progress` | `BcProgress` | `color` |
| `bc-chart-progress` | `BcProgress` | `showPercent` |
| `bc-chart-radar` | `BcChartRadar` | `dataSource` |
| `bc-chart-radar` | `BcChartRadar` | `fetchHeaders` |
| `bc-chart-radar` | `BcChartRadar` | `filled` |
| `bc-chart-radar` | `BcChartRadar` | `locale` |
| `bc-chart-radar` | `BcChartRadar` | `model` |
| `bc-chart-radar` | `BcChartRadar` | `refreshInterval` |
| `bc-chart-radar` | `BcChartRadar` | `shape` |
| `bc-chart-sankey` | `BcChartSankey` | `dataSource` |
| `bc-chart-sankey` | `BcChartSankey` | `draggable` |
| `bc-chart-sankey` | `BcChartSankey` | `fetchHeaders` |
| `bc-chart-sankey` | `BcChartSankey` | `locale` |
| `bc-chart-sankey` | `BcChartSankey` | `model` |
| `bc-chart-sankey` | `BcChartSankey` | `refreshInterval` |
| `bc-chart-scatter` | `BcChartScatter` | `bubble` |
| `bc-chart-scatter` | `BcChartScatter` | `dataSource` |
| `bc-chart-scatter` | `BcChartScatter` | `effectScatter` |
| `bc-chart-scatter` | `BcChartScatter` | `fetchHeaders` |
| `bc-chart-scatter` | `BcChartScatter` | `locale` |
| `bc-chart-scatter` | `BcChartScatter` | `model` |
| `bc-chart-scatter` | `BcChartScatter` | `refreshInterval` |
| `bc-chart-scorecard` | `BcChartScorecard` | `label` |
| `bc-chart-scorecard` | `BcChartScorecard` | `locale` |
| `bc-chart-scorecard` | `BcChartScorecard` | `target` |
| `bc-chart-scorecard` | `BcChartScorecard` | `value` |
| `bc-chart-sunburst` | `BcChartSunburst` | `dataSource` |
| `bc-chart-sunburst` | `BcChartSunburst` | `fetchHeaders` |
| `bc-chart-sunburst` | `BcChartSunburst` | `innerRadius` |
| `bc-chart-sunburst` | `BcChartSunburst` | `locale` |
| `bc-chart-sunburst` | `BcChartSunburst` | `model` |
| `bc-chart-sunburst` | `BcChartSunburst` | `outerRadius` |
| `bc-chart-sunburst` | `BcChartSunburst` | `refreshInterval` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `dataSource` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `fetchHeaders` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `locale` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `model` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `refreshInterval` |
| `bc-chart-tree` | `BcChartTree` | `dataSource` |
| `bc-chart-tree` | `BcChartTree` | `fetchHeaders` |
| `bc-chart-tree` | `BcChartTree` | `layout` |
| `bc-chart-tree` | `BcChartTree` | `locale` |
| `bc-chart-tree` | `BcChartTree` | `model` |
| `bc-chart-tree` | `BcChartTree` | `refreshInterval` |
| `bc-chart-treemap` | `BcChartTreemap` | `dataSource` |
| `bc-chart-treemap` | `BcChartTreemap` | `fetchHeaders` |
| `bc-chart-treemap` | `BcChartTreemap` | `locale` |
| `bc-chart-treemap` | `BcChartTreemap` | `model` |
| `bc-chart-treemap` | `BcChartTreemap` | `refreshInterval` |
| `bc-chart-treemap` | `BcChartTreemap` | `roam` |
| `bc-chatter` | `BcChatter` | `model` |
| `bc-chatter` | `BcChatter` | `recordId` |
| `bc-column` | `BcColumn` | `width` |
| `bc-datatable` | `BcDataTable` | `actions` |
| `bc-datatable` | `BcDataTable` | `apiUrl` |
| `bc-datatable` | `BcDataTable` | `dataSource` |
| `bc-datatable` | `BcDataTable` | `draggableColumns` |
| `bc-datatable` | `BcDataTable` | `fetchHeaders` |
| `bc-datatable` | `BcDataTable` | `localData` |
| `bc-datatable` | `BcDataTable` | `model` |
| `bc-datatable` | `BcDataTable` | `moduleName` |
| `bc-datatable` | `BcDataTable` | `savedPresets` |
| `bc-datatable` | `BcDataTable` | `selectable` |
| `bc-datatable` | `BcDataTable` | `serverSide` |
| `bc-datatable` | `BcDataTable` | `showFilterBuilder` |
| `bc-datatable` | `BcDataTable` | `showJsonFilter` |
| `bc-datatable` | `BcDataTable` | `viewTitle` |
| `bc-dialog-alert` | `BcAlertDialog` | `dialogTitle` |
| `bc-dialog-alert` | `BcAlertDialog` | `message` |
| `bc-dialog-alert` | `BcAlertDialog` | `okText` |
| `bc-dialog-alert` | `BcAlertDialog` | `size` |
| `bc-dialog-confirm` | `BcAlertDialog` | `cancelText` |
| `bc-dialog-confirm` | `BcAlertDialog` | `confirmText` |
| `bc-dialog-confirm` | `BcAlertDialog` | `dialogTitle` |
| `bc-dialog-confirm` | `BcAlertDialog` | `message` |
| `bc-dialog-confirm` | `BcAlertDialog` | `size` |
| `bc-dialog-modal` | `BcDialog` | `dialogTitle` |
| `bc-dialog-modal` | `BcDialog` | `loading` |
| `bc-dialog-prompt` | `BcPrompt` | `cancelText` |
| `bc-dialog-prompt` | `BcPrompt` | `dialogTitle` |
| `bc-dialog-prompt` | `BcPrompt` | `inputType` |
| `bc-dialog-prompt` | `BcPrompt` | `inputValue` |
| `bc-dialog-prompt` | `BcPrompt` | `message` |
| `bc-dialog-prompt` | `BcPrompt` | `okText` |
| `bc-dialog-prompt` | `BcPrompt` | `size` |
| `bc-dialog-wizard` | `BcWizard` | `dialogTitle` |
| `bc-export` | `BcExport` | `href` |
| `bc-favorites` | `BcFavorites` | `placeholder` |
| `bc-favorites` | `BcFavorites` | `value` |
| `bc-field-barcode` | `BcBarcode` | `autofocus` |
| `bc-field-barcode` | `BcBarcode` | `clearable` |
| `bc-field-barcode` | `BcBarcode` | `hint` |
| `bc-field-barcode` | `BcBarcode` | `loading` |
| `bc-field-barcode` | `BcBarcode` | `readonly` |
| `bc-field-barcode` | `BcBarcode` | `size` |
| `bc-field-barcode` | `BcBarcode` | `tooltip` |
| `bc-field-barcode` | `BcBarcode` | `validateOn` |
| `bc-field-barcode` | `BcBarcode` | `validationMessage` |
| `bc-field-barcode` | `BcBarcode` | `validationStatus` |
| `bc-field-checkbox` | `BcCheckbox` | `defaultValue` |
| `bc-field-checkbox` | `BcCheckbox` | `readonly` |
| `bc-field-checkbox` | `BcCheckbox` | `size` |
| `bc-field-checkbox` | `BcCheckbox` | `tooltip` |
| `bc-field-checkbox` | `BcCheckbox` | `validateOn` |
| `bc-field-checkbox` | `BcCheckbox` | `validationMessage` |
| `bc-field-checkbox` | `BcCheckbox` | `validationStatus` |
| `bc-field-checkbox` | `BcCheckbox` | `value` |
| `bc-field-code` | `BcCode` | `clearable` |
| `bc-field-code` | `BcCode` | `hint` |
| `bc-field-code` | `BcCode` | `loading` |
| `bc-field-code` | `BcCode` | `maxLength` |
| `bc-field-code` | `BcCode` | `minLength` |
| `bc-field-code` | `BcCode` | `placeholder` |
| `bc-field-code` | `BcCode` | `readonly` |
| `bc-field-code` | `BcCode` | `showCount` |
| `bc-field-code` | `BcCode` | `size` |
| `bc-field-code` | `BcCode` | `toolbar` |
| `bc-field-code` | `BcCode` | `tooltip` |
| `bc-field-code` | `BcCode` | `validateOn` |
| `bc-field-code` | `BcCode` | `validationMessage` |
| `bc-field-code` | `BcCode` | `validationStatus` |
| `bc-field-color` | `BcColorPicker` | `clearable` |
| `bc-field-color` | `BcColorPicker` | `hint` |
| `bc-field-color` | `BcColorPicker` | `readonly` |
| `bc-field-color` | `BcColorPicker` | `required` |
| `bc-field-color` | `BcColorPicker` | `size` |
| `bc-field-color` | `BcColorPicker` | `tooltip` |
| `bc-field-color` | `BcColorPicker` | `validateOn` |
| `bc-field-color` | `BcColorPicker` | `validationMessage` |
| `bc-field-color` | `BcColorPicker` | `validationStatus` |
| `bc-field-currency` | `BcNumber` | `autofocus` |
| `bc-field-currency` | `BcNumber` | `clearable` |
| `bc-field-currency` | `BcNumber` | `loading` |
| `bc-field-currency` | `BcNumber` | `precision` |
| `bc-field-currency` | `BcNumber` | `prefixText` |
| `bc-field-currency` | `BcNumber` | `size` |
| `bc-field-currency` | `BcNumber` | `suffixText` |
| `bc-field-currency` | `BcNumber` | `tooltip` |
| `bc-field-currency` | `BcNumber` | `validateOn` |
| `bc-field-currency` | `BcNumber` | `validationMessage` |
| `bc-field-currency` | `BcNumber` | `validationStatus` |
| `bc-field-date` | `BcDateField` | `autofocus` |
| `bc-field-date` | `BcDateField` | `clearable` |
| `bc-field-date` | `BcDateField` | `dependOn` |
| `bc-field-date` | `BcDateField` | `hint` |
| `bc-field-date` | `BcDateField` | `loading` |
| `bc-field-date` | `BcDateField` | `readonly` |
| `bc-field-date` | `BcDateField` | `size` |
| `bc-field-date` | `BcDateField` | `tooltip` |
| `bc-field-date` | `BcDateField` | `validateOn` |
| `bc-field-date` | `BcDateField` | `validationMessage` |
| `bc-field-date` | `BcDateField` | `validationStatus` |
| `bc-field-datetime` | `BcDateTime` | `autofocus` |
| `bc-field-datetime` | `BcDateTime` | `clearable` |
| `bc-field-datetime` | `BcDateTime` | `loading` |
| `bc-field-datetime` | `BcDateTime` | `size` |
| `bc-field-datetime` | `BcDateTime` | `tooltip` |
| `bc-field-datetime` | `BcDateTime` | `validateOn` |
| `bc-field-datetime` | `BcDateTime` | `validationMessage` |
| `bc-field-datetime` | `BcDateTime` | `validationStatus` |
| `bc-field-decimal` | `BcNumber` | `autofocus` |
| `bc-field-decimal` | `BcNumber` | `clearable` |
| `bc-field-decimal` | `BcNumber` | `loading` |
| `bc-field-decimal` | `BcNumber` | `precision` |
| `bc-field-decimal` | `BcNumber` | `prefixText` |
| `bc-field-decimal` | `BcNumber` | `size` |
| `bc-field-decimal` | `BcNumber` | `suffixText` |
| `bc-field-decimal` | `BcNumber` | `tooltip` |
| `bc-field-decimal` | `BcNumber` | `validateOn` |
| `bc-field-decimal` | `BcNumber` | `validationMessage` |
| `bc-field-decimal` | `BcNumber` | `validationStatus` |
| `bc-field-duration` | `BcFieldDuration` | `tooltip` |
| `bc-field-duration` | `BcFieldDuration` | `validateOn` |
| `bc-field-duration` | `BcFieldDuration` | `validationMessage` |
| `bc-field-duration` | `BcFieldDuration` | `validationStatus` |
| `bc-field-dynlink` | `BcFieldDynLink` | `clearable` |
| `bc-field-dynlink` | `BcFieldDynLink` | `disabled` |
| `bc-field-dynlink` | `BcFieldDynLink` | `hint` |
| `bc-field-dynlink` | `BcFieldDynLink` | `loading` |
| `bc-field-dynlink` | `BcFieldDynLink` | `model` |
| `bc-field-dynlink` | `BcFieldDynLink` | `modelfield` |
| `bc-field-dynlink` | `BcFieldDynLink` | `options` |
| `bc-field-dynlink` | `BcFieldDynLink` | `placeholder` |
| `bc-field-dynlink` | `BcFieldDynLink` | `readonly` |
| `bc-field-dynlink` | `BcFieldDynLink` | `required` |
| `bc-field-dynlink` | `BcFieldDynLink` | `size` |
| `bc-field-dynlink` | `BcFieldDynLink` | `tooltip` |
| `bc-field-dynlink` | `BcFieldDynLink` | `validateOn` |
| `bc-field-dynlink` | `BcFieldDynLink` | `validationMessage` |
| `bc-field-dynlink` | `BcFieldDynLink` | `validationStatus` |
| `bc-field-file` | `BcFile` | `apiBase` |
| `bc-field-file` | `BcFile` | `defaultValue` |
| `bc-field-file` | `BcFile` | `fieldName` |
| `bc-field-file` | `BcFile` | `hint` |
| `bc-field-file` | `BcFile` | `loading` |
| `bc-field-file` | `BcFile` | `model` |
| `bc-field-file` | `BcFile` | `nameFormat` |
| `bc-field-file` | `BcFile` | `pathFormat` |
| `bc-field-file` | `BcFile` | `preview` |
| `bc-field-file` | `BcFile` | `recordId` |
| `bc-field-file` | `BcFile` | `showDownload` |
| `bc-field-file` | `BcFile` | `size` |
| `bc-field-file` | `BcFile` | `tooltip` |
| `bc-field-file` | `BcFile` | `validateOn` |
| `bc-field-file` | `BcFile` | `validationMessage` |
| `bc-field-file` | `BcFile` | `validationStatus` |
| `bc-field-file` | `BcFile` | `value` |
| `bc-field-float` | `BcNumber` | `autofocus` |
| `bc-field-float` | `BcNumber` | `clearable` |
| `bc-field-float` | `BcNumber` | `loading` |
| `bc-field-float` | `BcNumber` | `precision` |
| `bc-field-float` | `BcNumber` | `prefixText` |
| `bc-field-float` | `BcNumber` | `size` |
| `bc-field-float` | `BcNumber` | `suffixText` |
| `bc-field-float` | `BcNumber` | `tooltip` |
| `bc-field-float` | `BcNumber` | `validateOn` |
| `bc-field-float` | `BcNumber` | `validationMessage` |
| `bc-field-float` | `BcNumber` | `validationStatus` |
| `bc-field-geo` | `BcGeo` | `clearable` |
| `bc-field-geo` | `BcGeo` | `defaultValue` |
| `bc-field-geo` | `BcGeo` | `disabled` |
| `bc-field-geo` | `BcGeo` | `drawMode` |
| `bc-field-geo` | `BcGeo` | `hint` |
| `bc-field-geo` | `BcGeo` | `placeholder` |
| `bc-field-geo` | `BcGeo` | `readonly` |
| `bc-field-geo` | `BcGeo` | `required` |
| `bc-field-geo` | `BcGeo` | `size` |
| `bc-field-geo` | `BcGeo` | `tooltip` |
| `bc-field-geo` | `BcGeo` | `validateOn` |
| `bc-field-geo` | `BcGeo` | `validationMessage` |
| `bc-field-geo` | `BcGeo` | `validationStatus` |
| `bc-field-geo` | `BcGeo` | `value` |
| `bc-field-geo` | `BcGeo` | `zoom` |
| `bc-field-html` | `BcRichText` | `clearable` |
| `bc-field-html` | `BcRichText` | `hint` |
| `bc-field-html` | `BcRichText` | `language` |
| `bc-field-html` | `BcRichText` | `loading` |
| `bc-field-html` | `BcRichText` | `maxLength` |
| `bc-field-html` | `BcRichText` | `minLength` |
| `bc-field-html` | `BcRichText` | `readonly` |
| `bc-field-html` | `BcRichText` | `showCount` |
| `bc-field-html` | `BcRichText` | `size` |
| `bc-field-html` | `BcRichText` | `toolbar` |
| `bc-field-html` | `BcRichText` | `tooltip` |
| `bc-field-html` | `BcRichText` | `validateOn` |
| `bc-field-html` | `BcRichText` | `validationMessage` |
| `bc-field-html` | `BcRichText` | `validationStatus` |
| `bc-field-image` | `BcFile` | `defaultValue` |
| `bc-field-image` | `BcFile` | `fieldName` |
| `bc-field-image` | `BcFile` | `hint` |
| `bc-field-image` | `BcFile` | `loading` |
| `bc-field-image` | `BcFile` | `model` |
| `bc-field-image` | `BcFile` | `preview` |
| `bc-field-image` | `BcFile` | `readonly` |
| `bc-field-image` | `BcFile` | `recordId` |
| `bc-field-image` | `BcFile` | `showDownload` |
| `bc-field-image` | `BcFile` | `size` |
| `bc-field-image` | `BcFile` | `tooltip` |
| `bc-field-image` | `BcFile` | `validateOn` |
| `bc-field-image` | `BcFile` | `validationMessage` |
| `bc-field-image` | `BcFile` | `validationStatus` |
| `bc-field-image` | `BcFile` | `value` |
| `bc-field-integer` | `BcNumber` | `autofocus` |
| `bc-field-integer` | `BcNumber` | `clearable` |
| `bc-field-integer` | `BcNumber` | `loading` |
| `bc-field-integer` | `BcNumber` | `prefixText` |
| `bc-field-integer` | `BcNumber` | `size` |
| `bc-field-integer` | `BcNumber` | `suffixText` |
| `bc-field-integer` | `BcNumber` | `tooltip` |
| `bc-field-integer` | `BcNumber` | `validateOn` |
| `bc-field-integer` | `BcNumber` | `validationMessage` |
| `bc-field-integer` | `BcNumber` | `validationStatus` |
| `bc-field-json` | `BcJson` | `clearable` |
| `bc-field-json` | `BcJson` | `hint` |
| `bc-field-json` | `BcJson` | `language` |
| `bc-field-json` | `BcJson` | `loading` |
| `bc-field-json` | `BcJson` | `placeholder` |
| `bc-field-json` | `BcJson` | `readonly` |
| `bc-field-json` | `BcJson` | `size` |
| `bc-field-json` | `BcJson` | `toolbar` |
| `bc-field-json` | `BcJson` | `tooltip` |
| `bc-field-json` | `BcJson` | `validateOn` |
| `bc-field-json` | `BcJson` | `validationMessage` |
| `bc-field-json` | `BcJson` | `validationStatus` |
| `bc-field-link` | `BcCombobox` | `dataSource` |
| `bc-field-link` | `BcCombobox` | `dependOn` |
| `bc-field-link` | `BcCombobox` | `displayField` |
| `bc-field-link` | `BcCombobox` | `fetchHeaders` |
| `bc-field-link` | `BcCombobox` | `hint` |
| `bc-field-link` | `BcCombobox` | `loading` |
| `bc-field-link` | `BcCombobox` | `lookupColumns` |
| `bc-field-link` | `BcCombobox` | `model` |
| `bc-field-link` | `BcCombobox` | `readonly` |
| `bc-field-link` | `BcCombobox` | `size` |
| `bc-field-link` | `BcCombobox` | `tooltip` |
| `bc-field-link` | `BcCombobox` | `validateOn` |
| `bc-field-link` | `BcCombobox` | `validationMessage` |
| `bc-field-link` | `BcCombobox` | `validationStatus` |
| `bc-field-markdown` | `BcRichText` | `clearable` |
| `bc-field-markdown` | `BcRichText` | `hint` |
| `bc-field-markdown` | `BcRichText` | `language` |
| `bc-field-markdown` | `BcRichText` | `loading` |
| `bc-field-markdown` | `BcRichText` | `maxLength` |
| `bc-field-markdown` | `BcRichText` | `minLength` |
| `bc-field-markdown` | `BcRichText` | `readonly` |
| `bc-field-markdown` | `BcRichText` | `showCount` |
| `bc-field-markdown` | `BcRichText` | `size` |
| `bc-field-markdown` | `BcRichText` | `toolbar` |
| `bc-field-markdown` | `BcRichText` | `tooltip` |
| `bc-field-markdown` | `BcRichText` | `validateOn` |
| `bc-field-markdown` | `BcRichText` | `validationMessage` |
| `bc-field-markdown` | `BcRichText` | `validationStatus` |
| `bc-field-morph` | `BcFieldMorph` | `loading` |
| `bc-field-morph` | `BcFieldMorph` | `models` |
| `bc-field-morph` | `BcFieldMorph` | `morphType` |
| `bc-field-morph` | `BcFieldMorph` | `size` |
| `bc-field-morph` | `BcFieldMorph` | `tooltip` |
| `bc-field-morph` | `BcFieldMorph` | `validateOn` |
| `bc-field-morph` | `BcFieldMorph` | `validationMessage` |
| `bc-field-morph` | `BcFieldMorph` | `validationStatus` |
| `bc-field-multicheck` | `BcCheckbox` | `clearable` |
| `bc-field-multicheck` | `BcCheckbox` | `defaultValue` |
| `bc-field-multicheck` | `BcCheckbox` | `options` |
| `bc-field-multicheck` | `BcCheckbox` | `readonly` |
| `bc-field-multicheck` | `BcCheckbox` | `size` |
| `bc-field-multicheck` | `BcCheckbox` | `tooltip` |
| `bc-field-multicheck` | `BcCheckbox` | `validateOn` |
| `bc-field-multicheck` | `BcCheckbox` | `validationMessage` |
| `bc-field-multicheck` | `BcCheckbox` | `validationStatus` |
| `bc-field-multicheck` | `BcCheckbox` | `value` |
| `bc-field-password` | `BcInput` | `showReveal` |
| `bc-field-percent` | `BcNumber` | `autofocus` |
| `bc-field-percent` | `BcNumber` | `clearable` |
| `bc-field-percent` | `BcNumber` | `loading` |
| `bc-field-percent` | `BcNumber` | `prefixText` |
| `bc-field-percent` | `BcNumber` | `size` |
| `bc-field-percent` | `BcNumber` | `suffixText` |
| `bc-field-percent` | `BcNumber` | `tooltip` |
| `bc-field-percent` | `BcNumber` | `validateOn` |
| `bc-field-percent` | `BcNumber` | `validationMessage` |
| `bc-field-percent` | `BcNumber` | `validationStatus` |
| `bc-field-radio` | `BcSegmented` | `direction` |
| `bc-field-radio` | `BcSegmented` | `hint` |
| `bc-field-radio` | `BcSegmented` | `label` |
| `bc-field-radio` | `BcSegmented` | `readonly` |
| `bc-field-radio` | `BcSegmented` | `required` |
| `bc-field-radio` | `BcSegmented` | `tooltip` |
| `bc-field-radio` | `BcSegmented` | `validateOn` |
| `bc-field-radio` | `BcSegmented` | `validationMessage` |
| `bc-field-radio` | `BcSegmented` | `validationStatus` |
| `bc-field-rating` | `BcRate` | `clearable` |
| `bc-field-rating` | `BcRate` | `hint` |
| `bc-field-rating` | `BcRate` | `label` |
| `bc-field-rating` | `BcRate` | `maxStars` |
| `bc-field-rating` | `BcRate` | `required` |
| `bc-field-rating` | `BcRate` | `tooltip` |
| `bc-field-rating` | `BcRate` | `validateOn` |
| `bc-field-rating` | `BcRate` | `validationMessage` |
| `bc-field-rating` | `BcRate` | `validationStatus` |
| `bc-field-richtext` | `BcRichText` | `clearable` |
| `bc-field-richtext` | `BcRichText` | `hint` |
| `bc-field-richtext` | `BcRichText` | `loading` |
| `bc-field-richtext` | `BcRichText` | `maxLength` |
| `bc-field-richtext` | `BcRichText` | `minLength` |
| `bc-field-richtext` | `BcRichText` | `readonly` |
| `bc-field-richtext` | `BcRichText` | `showCount` |
| `bc-field-richtext` | `BcRichText` | `size` |
| `bc-field-richtext` | `BcRichText` | `toolbar` |
| `bc-field-richtext` | `BcRichText` | `tooltip` |
| `bc-field-richtext` | `BcRichText` | `validateOn` |
| `bc-field-richtext` | `BcRichText` | `validationMessage` |
| `bc-field-richtext` | `BcRichText` | `validationStatus` |
| `bc-field-select` | `BcSelect` | `creatable` |
| `bc-field-select` | `BcSelect` | `dataSource` |
| `bc-field-select` | `BcSelect` | `debounceMs` |
| `bc-field-select` | `BcSelect` | `dependOn` |
| `bc-field-select` | `BcSelect` | `displayField` |
| `bc-field-select` | `BcSelect` | `fetchHeaders` |
| `bc-field-select` | `BcSelect` | `groupBy` |
| `bc-field-select` | `BcSelect` | `loading` |
| `bc-field-select` | `BcSelect` | `loadingText` |
| `bc-field-select` | `BcSelect` | `minSearchLength` |
| `bc-field-select` | `BcSelect` | `model` |
| `bc-field-select` | `BcSelect` | `multiple` |
| `bc-field-select` | `BcSelect` | `noResultsText` |
| `bc-field-select` | `BcSelect` | `pageSize` |
| `bc-field-select` | `BcSelect` | `readonly` |
| `bc-field-select` | `BcSelect` | `searchable` |
| `bc-field-select` | `BcSelect` | `serverSide` |
| `bc-field-select` | `BcSelect` | `size` |
| `bc-field-select` | `BcSelect` | `tooltip` |
| `bc-field-select` | `BcSelect` | `validateOn` |
| `bc-field-select` | `BcSelect` | `validationMessage` |
| `bc-field-select` | `BcSelect` | `validationStatus` |
| `bc-field-select` | `BcSelect` | `valueField` |
| `bc-field-signature` | `BcSignature` | `defaultValue` |
| `bc-field-signature` | `BcSignature` | `hint` |
| `bc-field-signature` | `BcSignature` | `readonly` |
| `bc-field-signature` | `BcSignature` | `size` |
| `bc-field-signature` | `BcSignature` | `tooltip` |
| `bc-field-signature` | `BcSignature` | `validateOn` |
| `bc-field-signature` | `BcSignature` | `validationMessage` |
| `bc-field-signature` | `BcSignature` | `validationStatus` |
| `bc-field-signature` | `BcSignature` | `value` |
| `bc-field-smalltext` | `BcTextarea` | `autofocus` |
| `bc-field-smalltext` | `BcTextarea` | `clearable` |
| `bc-field-smalltext` | `BcTextarea` | `loading` |
| `bc-field-smalltext` | `BcTextarea` | `minLength` |
| `bc-field-smalltext` | `BcTextarea` | `size` |
| `bc-field-smalltext` | `BcTextarea` | `tooltip` |
| `bc-field-smalltext` | `BcTextarea` | `validateOn` |
| `bc-field-smalltext` | `BcTextarea` | `validationMessage` |
| `bc-field-smalltext` | `BcTextarea` | `validationStatus` |
| `bc-field-string` | `BcInput` | `dataSource` |
| `bc-field-string` | `BcInput` | `dependOn` |
| `bc-field-string` | `BcInput` | `max` |
| `bc-field-string` | `BcInput` | `widget` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `clearable` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `disabled` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `hint` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `loading` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `model` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `options` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `placeholder` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `readonly` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `required` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `size` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `tooltip` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `validateOn` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `validationMessage` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `validationStatus` |
| `bc-field-tags` | `BcTagsInput` | `clearable` |
| `bc-field-tags` | `BcTagsInput` | `hint` |
| `bc-field-tags` | `BcTagsInput` | `loading` |
| `bc-field-tags` | `BcTagsInput` | `model` |
| `bc-field-tags` | `BcTagsInput` | `options` |
| `bc-field-tags` | `BcTagsInput` | `readonly` |
| `bc-field-tags` | `BcTagsInput` | `size` |
| `bc-field-tags` | `BcTagsInput` | `tooltip` |
| `bc-field-tags` | `BcTagsInput` | `validateOn` |
| `bc-field-tags` | `BcTagsInput` | `validationMessage` |
| `bc-field-tags` | `BcTagsInput` | `validationStatus` |
| `bc-field-text` | `BcTextarea` | `autofocus` |
| `bc-field-text` | `BcTextarea` | `clearable` |
| `bc-field-text` | `BcTextarea` | `loading` |
| `bc-field-text` | `BcTextarea` | `minLength` |
| `bc-field-text` | `BcTextarea` | `pattern` |
| `bc-field-text` | `BcTextarea` | `patternMessage` |
| `bc-field-text` | `BcTextarea` | `size` |
| `bc-field-text` | `BcTextarea` | `tooltip` |
| `bc-field-text` | `BcTextarea` | `validateOn` |
| `bc-field-text` | `BcTextarea` | `validationMessage` |
| `bc-field-text` | `BcTextarea` | `validationStatus` |
| `bc-field-time` | `BcTimeField` | `autofocus` |
| `bc-field-time` | `BcTimeField` | `clearable` |
| `bc-field-time` | `BcTimeField` | `hint` |
| `bc-field-time` | `BcTimeField` | `loading` |
| `bc-field-time` | `BcTimeField` | `max` |
| `bc-field-time` | `BcTimeField` | `min` |
| `bc-field-time` | `BcTimeField` | `readonly` |
| `bc-field-time` | `BcTimeField` | `size` |
| `bc-field-time` | `BcTimeField` | `step` |
| `bc-field-time` | `BcTimeField` | `tooltip` |
| `bc-field-time` | `BcTimeField` | `validateOn` |
| `bc-field-time` | `BcTimeField` | `validationMessage` |
| `bc-field-time` | `BcTimeField` | `validationStatus` |
| `bc-field-toggle` | `BcToggle` | `defaultValue` |
| `bc-field-toggle` | `BcToggle` | `readonly` |
| `bc-field-toggle` | `BcToggle` | `required` |
| `bc-field-toggle` | `BcToggle` | `tooltip` |
| `bc-field-toggle` | `BcToggle` | `validateOn` |
| `bc-field-toggle` | `BcToggle` | `validationMessage` |
| `bc-field-toggle` | `BcToggle` | `validationStatus` |
| `bc-field-toggle` | `BcToggle` | `value` |
| `bc-header` | `BcHeader` | `buttons` |
| `bc-header` | `BcHeader` | `states` |
| `bc-header` | `BcHeader` | `statusField` |
| `bc-header` | `BcHeader` | `statusValue` |
| `bc-kanban-board` | `BcKanbanBoard` | `activityDataSource` |
| `bc-kanban-board` | `BcKanbanBoard` | `activityModel` |
| `bc-kanban-board` | `BcKanbanBoard` | `allowAddCard` |
| `bc-kanban-board` | `BcKanbanBoard` | `allowAddColumn` |
| `bc-kanban-board` | `BcKanbanBoard` | `allowDeleteColumn` |
| `bc-kanban-board` | `BcKanbanBoard` | `allowRenameColumn` |
| `bc-kanban-board` | `BcKanbanBoard` | `allowReorderColumns` |
| `bc-kanban-board` | `BcKanbanBoard` | `attachmentDataSource` |
| `bc-kanban-board` | `BcKanbanBoard` | `attachmentModel` |
| `bc-kanban-board` | `BcKanbanBoard` | `boardTitle` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardAssigneesField` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardCoverField` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardDescriptionField` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardDetailDataSource` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardDetailModel` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardDueDateField` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardLabelsField` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardPositionField` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardPriorityField` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardStartDateField` |
| `bc-kanban-board` | `BcKanbanBoard` | `cardTitleField` |
| `bc-kanban-board` | `BcKanbanBoard` | `columnDataSource` |
| `bc-kanban-board` | `BcKanbanBoard` | `columnModel` |
| `bc-kanban-board` | `BcKanbanBoard` | `commentDataSource` |
| `bc-kanban-board` | `BcKanbanBoard` | `commentModel` |
| `bc-kanban-board` | `BcKanbanBoard` | `dataSource` |
| `bc-kanban-board` | `BcKanbanBoard` | `fetchHeaders` |
| `bc-kanban-board` | `BcKanbanBoard` | `groupBy` |
| `bc-kanban-board` | `BcKanbanBoard` | `mentionDataSource` |
| `bc-kanban-board` | `BcKanbanBoard` | `mentionModel` |
| `bc-kanban-board` | `BcKanbanBoard` | `model` |
| `bc-kanban-board` | `BcKanbanBoard` | `subtaskDataSource` |
| `bc-kanban-board` | `BcKanbanBoard` | `subtaskModel` |
| `bc-lookup-modal` | `BcLookupModal` | `apiUrl` |
| `bc-lookup-modal` | `BcLookupModal` | `dataSource` |
| `bc-lookup-modal` | `BcLookupModal` | `displayField` |
| `bc-lookup-modal` | `BcLookupModal` | `fetchHeaders` |
| `bc-lookup-modal` | `BcLookupModal` | `modalTitle` |
| `bc-print` | `BcPrint` | `href` |
| `bc-search` | `BcSearch` | `dataSource` |
| `bc-search` | `BcSearch` | `fetchHeaders` |
| `bc-search` | `BcSearch` | `model` |
| `bc-section` | `BcSection` | `collapsed` |
| `bc-section` | `BcSection` | `description` |
| `bc-section` | `BcSection` | `sectionTitle` |
| `bc-separator` | `BcSeparator` | `label` |
| `bc-timeline` | `BcTimelineSocial` | `dataSource` |
| `bc-timeline` | `BcTimelineSocial` | `fetchHeaders` |
| `bc-timeline` | `BcTimelineSocial` | `model` |
| `bc-timeline` | `BcTimelineSocial` | `recordId` |
| `bc-toast` | `BcToast` | `dialogTitle` |
| `bc-toast` | `BcToast` | `open` |
| `bc-toast` | `BcToast` | `position` |
| `bc-toast` | `BcToast` | `variant` |
| `bc-view-activity` | `BcViewActivity` | `config` |
| `bc-view-activity` | `BcViewActivity` | `dataSource` |
| `bc-view-activity` | `BcViewActivity` | `fetchHeaders` |
| `bc-view-activity` | `BcViewActivity` | `fields` |
| `bc-view-activity` | `BcViewActivity` | `model` |
| `bc-view-activity` | `BcViewActivity` | `viewTitle` |
| `bc-view-calendar` | `BcCalendar` | `config` |
| `bc-view-calendar` | `BcCalendar` | `dataSource` |
| `bc-view-calendar` | `BcCalendar` | `dateField` |
| `bc-view-calendar` | `BcCalendar` | `fetchHeaders` |
| `bc-view-calendar` | `BcCalendar` | `fields` |
| `bc-view-calendar` | `BcCalendar` | `model` |
| `bc-view-calendar` | `BcCalendar` | `titleField` |
| `bc-view-calendar` | `BcCalendar` | `viewTitle` |
| `bc-view-editor` | `BcViewEditor` | `modelFields` |
| `bc-view-editor` | `BcViewEditor` | `readonly` |
| `bc-view-editor` | `BcViewEditor` | `viewJson` |
| `bc-viewer-audio` | `BcViewerAudio` | `dataSource` |
| `bc-viewer-audio` | `BcViewerAudio` | `download` |
| `bc-viewer-audio` | `BcViewerAudio` | `loading` |
| `bc-viewer-audio` | `BcViewerAudio` | `srcField` |
| `bc-viewer-document` | `BcViewerDocument` | `dataSource` |
| `bc-viewer-document` | `BcViewerDocument` | `download` |
| `bc-viewer-document` | `BcViewerDocument` | `loading` |
| `bc-viewer-document` | `BcViewerDocument` | `provider` |
| `bc-viewer-document` | `BcViewerDocument` | `srcField` |
| `bc-viewer-image` | `BcViewerImage` | `dataSource` |
| `bc-viewer-image` | `BcViewerImage` | `download` |
| `bc-viewer-image` | `BcViewerImage` | `lightbox` |
| `bc-viewer-image` | `BcViewerImage` | `loading` |
| `bc-viewer-image` | `BcViewerImage` | `srcField` |
| `bc-viewer-instagram` | `BcViewerInstagram` | `dataSource` |
| `bc-viewer-instagram` | `BcViewerInstagram` | `loading` |
| `bc-viewer-instagram` | `BcViewerInstagram` | `src` |
| `bc-viewer-instagram` | `BcViewerInstagram` | `srcField` |
| `bc-viewer-pdf` | `BcViewerPdf` | `dataSource` |
| `bc-viewer-pdf` | `BcViewerPdf` | `download` |
| `bc-viewer-pdf` | `BcViewerPdf` | `loading` |
| `bc-viewer-pdf` | `BcViewerPdf` | `srcField` |
| `bc-viewer-tiktok` | `BcViewerTiktok` | `dataSource` |
| `bc-viewer-tiktok` | `BcViewerTiktok` | `loading` |
| `bc-viewer-tiktok` | `BcViewerTiktok` | `src` |
| `bc-viewer-tiktok` | `BcViewerTiktok` | `srcField` |
| `bc-viewer-video` | `BcViewerVideo` | `dataSource` |
| `bc-viewer-video` | `BcViewerVideo` | `download` |
| `bc-viewer-video` | `BcViewerVideo` | `loading` |
| `bc-viewer-video` | `BcViewerVideo` | `srcField` |
| `bc-viewer-youtube` | `BcViewerYoutube` | `dataSource` |
| `bc-viewer-youtube` | `BcViewerYoutube` | `height` |
| `bc-viewer-youtube` | `BcViewerYoutube` | `loading` |
| `bc-viewer-youtube` | `BcViewerYoutube` | `src` |
| `bc-viewer-youtube` | `BcViewerYoutube` | `srcField` |
| `bc-viewer-youtube` | `BcViewerYoutube` | `width` |
| `bc-view-form` | `BcForm` | `config` |
| `bc-view-form` | `BcForm` | `dataSource` |
| `bc-view-form` | `BcForm` | `fetchHeaders` |
| `bc-view-form` | `BcForm` | `model` |
| `bc-view-form` | `BcForm` | `moduleName` |
| `bc-view-form` | `BcForm` | `permissions` |
| `bc-view-form` | `BcForm` | `recordId` |
| `bc-view-form` | `BcForm` | `viewTitle` |
| `bc-view-gantt` | `BcViewGantt` | `columns` |
| `bc-view-gantt` | `BcViewGantt` | `config` |
| `bc-view-gantt` | `BcViewGantt` | `dataSource` |
| `bc-view-gantt` | `BcViewGantt` | `fetchHeaders` |
| `bc-view-gantt` | `BcViewGantt` | `fieldMapping` |
| `bc-view-gantt` | `BcViewGantt` | `fields` |
| `bc-view-gantt` | `BcViewGantt` | `height` |
| `bc-view-gantt` | `BcViewGantt` | `loading` |
| `bc-view-gantt` | `BcViewGantt` | `markers` |
| `bc-view-gantt` | `BcViewGantt` | `readonly` |
| `bc-view-gantt` | `BcViewGantt` | `refreshInterval` |
| `bc-view-gantt` | `BcViewGantt` | `scales` |
| `bc-view-gantt` | `BcViewGantt` | `templates` |
| `bc-view-gantt` | `BcViewGantt` | `viewTitle` |
| `bc-view-kanban` | `BcKanban` | `config` |
| `bc-view-kanban` | `BcKanban` | `dataSource` |
| `bc-view-kanban` | `BcKanban` | `fetchHeaders` |
| `bc-view-kanban` | `BcKanban` | `fields` |
| `bc-view-kanban` | `BcKanban` | `model` |
| `bc-view-kanban` | `BcKanban` | `viewTitle` |
| `bc-view-list` | `BcList` | `config` |
| `bc-view-list` | `BcList` | `dataSource` |
| `bc-view-list` | `BcList` | `fetchHeaders` |
| `bc-view-list` | `BcList` | `fields` |
| `bc-view-list` | `BcList` | `model` |
| `bc-view-list` | `BcList` | `refreshInterval` |
| `bc-view-list` | `BcList` | `viewTitle` |
| `bc-view-map` | `BcViewMap` | `config` |
| `bc-view-map` | `BcViewMap` | `dataSource` |
| `bc-view-map` | `BcViewMap` | `fetchHeaders` |
| `bc-view-map` | `BcViewMap` | `fields` |
| `bc-view-map` | `BcViewMap` | `geoField` |
| `bc-view-map` | `BcViewMap` | `model` |
| `bc-view-map` | `BcViewMap` | `nameField` |
| `bc-view-map` | `BcViewMap` | `viewTitle` |
| `bc-view-report` | `BcViewReport` | `config` |
| `bc-view-report` | `BcViewReport` | `dataSource` |
| `bc-view-report` | `BcViewReport` | `fetchHeaders` |
| `bc-view-report` | `BcViewReport` | `fields` |
| `bc-view-report` | `BcViewReport` | `model` |
| `bc-view-report` | `BcViewReport` | `viewTitle` |
| `bc-view-tree` | `BcTree` | `config` |
| `bc-view-tree` | `BcTree` | `dataSource` |
| `bc-view-tree` | `BcTree` | `fetchHeaders` |
| `bc-view-tree` | `BcTree` | `fields` |
| `bc-view-tree` | `BcTree` | `model` |
| `bc-view-tree` | `BcTree` | `parentField` |
| `bc-view-tree` | `BcTree` | `viewTitle` |

### Missing Events in Mitosis (mapped components)

| Stencil Component | Mitosis Component | Missing Event |
|-------------------|-------------------|---------------|
| `bc-activity` | `BcActivity` | `lcActivitySchedule` |
| `bc-chart-area` | `BcChartArea` | `lcChartDataZoom` |
| `bc-chart-area` | `BcChartArea` | `lcChartLegendSelect` |
| `bc-chart-bar` | `BcChartBar` | `lcChartDataZoom` |
| `bc-chart-bar` | `BcChartBar` | `lcChartLegendSelect` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `lcChartDataZoom` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `lcChartLegendSelect` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `lcChartDataZoom` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `lcChartLegendSelect` |
| `bc-chart-funnel` | `BcChartFunnel` | `lcChartDataZoom` |
| `bc-chart-funnel` | `BcChartFunnel` | `lcChartLegendSelect` |
| `bc-chart-gauge` | `BcChartGauge` | `lcChartClick` |
| `bc-chart-gauge` | `BcChartGauge` | `lcChartDataZoom` |
| `bc-chart-gauge` | `BcChartGauge` | `lcChartHover` |
| `bc-chart-gauge` | `BcChartGauge` | `lcChartLegendSelect` |
| `bc-chart-gauge` | `BcChartGauge` | `lcChartReady` |
| `bc-chart-graph` | `BcChartGraph` | `lcChartDataZoom` |
| `bc-chart-graph` | `BcChartGraph` | `lcChartLegendSelect` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `lcChartDataZoom` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `lcChartLegendSelect` |
| `bc-chart-line` | `BcChartLine` | `lcChartDataZoom` |
| `bc-chart-line` | `BcChartLine` | `lcChartLegendSelect` |
| `bc-chart-mixed` | `BcChartMixed` | `lcChartDataZoom` |
| `bc-chart-mixed` | `BcChartMixed` | `lcChartLegendSelect` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `lcChartDataZoom` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `lcChartLegendSelect` |
| `bc-chart-pie` | `BcChartPie` | `lcChartDataZoom` |
| `bc-chart-pie` | `BcChartPie` | `lcChartLegendSelect` |
| `bc-chart-polar` | `BcChartPolar` | `lcChartDataZoom` |
| `bc-chart-polar` | `BcChartPolar` | `lcChartLegendSelect` |
| `bc-chart-radar` | `BcChartRadar` | `lcChartClick` |
| `bc-chart-radar` | `BcChartRadar` | `lcChartDataZoom` |
| `bc-chart-radar` | `BcChartRadar` | `lcChartHover` |
| `bc-chart-radar` | `BcChartRadar` | `lcChartLegendSelect` |
| `bc-chart-radar` | `BcChartRadar` | `lcChartReady` |
| `bc-chart-radar` | `BcChartRadar` | `lcError` |
| `bc-chart-sankey` | `BcChartSankey` | `lcChartDataZoom` |
| `bc-chart-sankey` | `BcChartSankey` | `lcChartLegendSelect` |
| `bc-chart-scatter` | `BcChartScatter` | `lcChartDataZoom` |
| `bc-chart-scatter` | `BcChartScatter` | `lcChartLegendSelect` |
| `bc-chart-scorecard` | `BcChartScorecard` | `lcChartDataZoom` |
| `bc-chart-scorecard` | `BcChartScorecard` | `lcChartLegendSelect` |
| `bc-chart-sunburst` | `BcChartSunburst` | `lcChartDataZoom` |
| `bc-chart-sunburst` | `BcChartSunburst` | `lcChartLegendSelect` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `lcChartLegendSelect` |
| `bc-chart-tree` | `BcChartTree` | `lcChartDataZoom` |
| `bc-chart-tree` | `BcChartTree` | `lcChartLegendSelect` |
| `bc-chart-treemap` | `BcChartTreemap` | `lcChartDataZoom` |
| `bc-chart-treemap` | `BcChartTreemap` | `lcChartLegendSelect` |
| `bc-chatter` | `BcChatter` | `lcChatterSend` |
| `bc-datatable` | `BcDataTable` | `lcBulkAction` |
| `bc-datatable` | `BcDataTable` | `lcFilterChange` |
| `bc-datatable` | `BcDataTable` | `lcPageChange` |
| `bc-datatable` | `BcDataTable` | `lcSelectionChange` |
| `bc-datatable` | `BcDataTable` | `lcSortChange` |
| `bc-dialog-alert` | `BcAlertDialog` | `lcDialogClose` |
| `bc-dialog-alert` | `BcAlertDialog` | `lcDialogConfirm` |
| `bc-dialog-confirm` | `BcAlertDialog` | `lcDialogCancel` |
| `bc-dialog-confirm` | `BcAlertDialog` | `lcDialogClose` |
| `bc-dialog-confirm` | `BcAlertDialog` | `lcDialogConfirm` |
| `bc-dialog-modal` | `BcDialog` | `lcDialogClose` |
| `bc-dialog-prompt` | `BcPrompt` | `lcDialogCancel` |
| `bc-dialog-prompt` | `BcPrompt` | `lcDialogClose` |
| `bc-dialog-prompt` | `BcPrompt` | `lcDialogConfirm` |
| `bc-dialog-wizard` | `BcWizard` | `lcDialogClose` |
| `bc-dialog-wizard` | `BcWizard` | `lcWizardComplete` |
| `bc-favorites` | `BcFavorites` | `lcFavoriteSelect` |
| `bc-favorites` | `BcFavorites` | `lcSearch` |
| `bc-field-barcode` | `BcBarcode` | `lcFieldBlur` |
| `bc-field-barcode` | `BcBarcode` | `lcFieldChange` |
| `bc-field-barcode` | `BcBarcode` | `lcFieldClear` |
| `bc-field-barcode` | `BcBarcode` | `lcFieldFocus` |
| `bc-field-barcode` | `BcBarcode` | `lcFieldInvalid` |
| `bc-field-barcode` | `BcBarcode` | `lcFieldValid` |
| `bc-field-checkbox` | `BcCheckbox` | `lcFieldBlur` |
| `bc-field-checkbox` | `BcCheckbox` | `lcFieldChange` |
| `bc-field-checkbox` | `BcCheckbox` | `lcFieldClear` |
| `bc-field-checkbox` | `BcCheckbox` | `lcFieldFocus` |
| `bc-field-checkbox` | `BcCheckbox` | `lcFieldInvalid` |
| `bc-field-checkbox` | `BcCheckbox` | `lcFieldValid` |
| `bc-field-code` | `BcCode` | `lcFieldBlur` |
| `bc-field-code` | `BcCode` | `lcFieldChange` |
| `bc-field-code` | `BcCode` | `lcFieldClear` |
| `bc-field-code` | `BcCode` | `lcFieldFocus` |
| `bc-field-code` | `BcCode` | `lcFieldInvalid` |
| `bc-field-code` | `BcCode` | `lcFieldValid` |
| `bc-field-color` | `BcColorPicker` | `lcFieldBlur` |
| `bc-field-color` | `BcColorPicker` | `lcFieldChange` |
| `bc-field-color` | `BcColorPicker` | `lcFieldClear` |
| `bc-field-color` | `BcColorPicker` | `lcFieldFocus` |
| `bc-field-color` | `BcColorPicker` | `lcFieldInvalid` |
| `bc-field-color` | `BcColorPicker` | `lcFieldValid` |
| `bc-field-currency` | `BcNumber` | `lcFieldBlur` |
| `bc-field-currency` | `BcNumber` | `lcFieldChange` |
| `bc-field-currency` | `BcNumber` | `lcFieldClear` |
| `bc-field-currency` | `BcNumber` | `lcFieldFocus` |
| `bc-field-currency` | `BcNumber` | `lcFieldInvalid` |
| `bc-field-currency` | `BcNumber` | `lcFieldValid` |
| `bc-field-date` | `BcDateField` | `lcFieldBlur` |
| `bc-field-date` | `BcDateField` | `lcFieldChange` |
| `bc-field-date` | `BcDateField` | `lcFieldClear` |
| `bc-field-date` | `BcDateField` | `lcFieldFocus` |
| `bc-field-date` | `BcDateField` | `lcFieldInvalid` |
| `bc-field-date` | `BcDateField` | `lcFieldValid` |
| `bc-field-datetime` | `BcDateTime` | `lcFieldBlur` |
| `bc-field-datetime` | `BcDateTime` | `lcFieldChange` |
| `bc-field-datetime` | `BcDateTime` | `lcFieldClear` |
| `bc-field-datetime` | `BcDateTime` | `lcFieldFocus` |
| `bc-field-datetime` | `BcDateTime` | `lcFieldInvalid` |
| `bc-field-datetime` | `BcDateTime` | `lcFieldValid` |
| `bc-field-decimal` | `BcNumber` | `lcFieldBlur` |
| `bc-field-decimal` | `BcNumber` | `lcFieldChange` |
| `bc-field-decimal` | `BcNumber` | `lcFieldClear` |
| `bc-field-decimal` | `BcNumber` | `lcFieldFocus` |
| `bc-field-decimal` | `BcNumber` | `lcFieldInvalid` |
| `bc-field-decimal` | `BcNumber` | `lcFieldValid` |
| `bc-field-duration` | `BcFieldDuration` | `lcFieldBlur` |
| `bc-field-duration` | `BcFieldDuration` | `lcFieldChange` |
| `bc-field-duration` | `BcFieldDuration` | `lcFieldClear` |
| `bc-field-duration` | `BcFieldDuration` | `lcFieldFocus` |
| `bc-field-duration` | `BcFieldDuration` | `lcFieldInvalid` |
| `bc-field-duration` | `BcFieldDuration` | `lcFieldValid` |
| `bc-field-dynlink` | `BcFieldDynLink` | `lcFieldBlur` |
| `bc-field-dynlink` | `BcFieldDynLink` | `lcFieldChange` |
| `bc-field-dynlink` | `BcFieldDynLink` | `lcFieldClear` |
| `bc-field-dynlink` | `BcFieldDynLink` | `lcFieldFocus` |
| `bc-field-dynlink` | `BcFieldDynLink` | `lcFieldInvalid` |
| `bc-field-dynlink` | `BcFieldDynLink` | `lcFieldValid` |
| `bc-field-file` | `BcFile` | `lcFieldBlur` |
| `bc-field-file` | `BcFile` | `lcFieldChange` |
| `bc-field-file` | `BcFile` | `lcFieldClear` |
| `bc-field-file` | `BcFile` | `lcFieldFocus` |
| `bc-field-file` | `BcFile` | `lcFieldInvalid` |
| `bc-field-file` | `BcFile` | `lcFieldValid` |
| `bc-field-float` | `BcNumber` | `lcFieldBlur` |
| `bc-field-float` | `BcNumber` | `lcFieldChange` |
| `bc-field-float` | `BcNumber` | `lcFieldClear` |
| `bc-field-float` | `BcNumber` | `lcFieldFocus` |
| `bc-field-float` | `BcNumber` | `lcFieldInvalid` |
| `bc-field-float` | `BcNumber` | `lcFieldValid` |
| `bc-field-geo` | `BcGeo` | `lcFieldBlur` |
| `bc-field-geo` | `BcGeo` | `lcFieldChange` |
| `bc-field-geo` | `BcGeo` | `lcFieldClear` |
| `bc-field-geo` | `BcGeo` | `lcFieldFocus` |
| `bc-field-geo` | `BcGeo` | `lcFieldInvalid` |
| `bc-field-geo` | `BcGeo` | `lcFieldValid` |
| `bc-field-html` | `BcRichText` | `lcFieldBlur` |
| `bc-field-html` | `BcRichText` | `lcFieldChange` |
| `bc-field-html` | `BcRichText` | `lcFieldClear` |
| `bc-field-html` | `BcRichText` | `lcFieldFocus` |
| `bc-field-html` | `BcRichText` | `lcFieldInvalid` |
| `bc-field-html` | `BcRichText` | `lcFieldValid` |
| `bc-field-image` | `BcFile` | `lcFieldBlur` |
| `bc-field-image` | `BcFile` | `lcFieldChange` |
| `bc-field-image` | `BcFile` | `lcFieldClear` |
| `bc-field-image` | `BcFile` | `lcFieldFocus` |
| `bc-field-image` | `BcFile` | `lcFieldInvalid` |
| `bc-field-image` | `BcFile` | `lcFieldValid` |
| `bc-field-integer` | `BcNumber` | `lcFieldBlur` |
| `bc-field-integer` | `BcNumber` | `lcFieldChange` |
| `bc-field-integer` | `BcNumber` | `lcFieldClear` |
| `bc-field-integer` | `BcNumber` | `lcFieldFocus` |
| `bc-field-integer` | `BcNumber` | `lcFieldInvalid` |
| `bc-field-integer` | `BcNumber` | `lcFieldValid` |
| `bc-field-json` | `BcJson` | `lcFieldBlur` |
| `bc-field-json` | `BcJson` | `lcFieldChange` |
| `bc-field-json` | `BcJson` | `lcFieldClear` |
| `bc-field-json` | `BcJson` | `lcFieldFocus` |
| `bc-field-json` | `BcJson` | `lcFieldInvalid` |
| `bc-field-json` | `BcJson` | `lcFieldValid` |
| `bc-field-link` | `BcCombobox` | `lcFieldBlur` |
| `bc-field-link` | `BcCombobox` | `lcFieldChange` |
| `bc-field-link` | `BcCombobox` | `lcFieldClear` |
| `bc-field-link` | `BcCombobox` | `lcFieldFocus` |
| `bc-field-link` | `BcCombobox` | `lcFieldInvalid` |
| `bc-field-link` | `BcCombobox` | `lcFieldValid` |
| `bc-field-markdown` | `BcRichText` | `lcFieldBlur` |
| `bc-field-markdown` | `BcRichText` | `lcFieldChange` |
| `bc-field-markdown` | `BcRichText` | `lcFieldClear` |
| `bc-field-markdown` | `BcRichText` | `lcFieldFocus` |
| `bc-field-markdown` | `BcRichText` | `lcFieldInvalid` |
| `bc-field-markdown` | `BcRichText` | `lcFieldValid` |
| `bc-field-morph` | `BcFieldMorph` | `lcFieldBlur` |
| `bc-field-morph` | `BcFieldMorph` | `lcFieldChange` |
| `bc-field-morph` | `BcFieldMorph` | `lcFieldClear` |
| `bc-field-morph` | `BcFieldMorph` | `lcFieldFocus` |
| `bc-field-morph` | `BcFieldMorph` | `lcFieldInvalid` |
| `bc-field-morph` | `BcFieldMorph` | `lcFieldValid` |
| `bc-field-multicheck` | `BcCheckbox` | `lcFieldBlur` |
| `bc-field-multicheck` | `BcCheckbox` | `lcFieldChange` |
| `bc-field-multicheck` | `BcCheckbox` | `lcFieldClear` |
| `bc-field-multicheck` | `BcCheckbox` | `lcFieldFocus` |
| `bc-field-multicheck` | `BcCheckbox` | `lcFieldInvalid` |
| `bc-field-multicheck` | `BcCheckbox` | `lcFieldValid` |
| `bc-field-password` | `BcInput` | `lcFieldBlur` |
| `bc-field-password` | `BcInput` | `lcFieldChange` |
| `bc-field-password` | `BcInput` | `lcFieldClear` |
| `bc-field-password` | `BcInput` | `lcFieldFocus` |
| `bc-field-password` | `BcInput` | `lcFieldInvalid` |
| `bc-field-password` | `BcInput` | `lcFieldValid` |
| `bc-field-percent` | `BcNumber` | `lcFieldBlur` |
| `bc-field-percent` | `BcNumber` | `lcFieldChange` |
| `bc-field-percent` | `BcNumber` | `lcFieldClear` |
| `bc-field-percent` | `BcNumber` | `lcFieldFocus` |
| `bc-field-percent` | `BcNumber` | `lcFieldInvalid` |
| `bc-field-percent` | `BcNumber` | `lcFieldValid` |
| `bc-field-radio` | `BcSegmented` | `lcFieldBlur` |
| `bc-field-radio` | `BcSegmented` | `lcFieldChange` |
| `bc-field-radio` | `BcSegmented` | `lcFieldClear` |
| `bc-field-radio` | `BcSegmented` | `lcFieldFocus` |
| `bc-field-radio` | `BcSegmented` | `lcFieldInvalid` |
| `bc-field-radio` | `BcSegmented` | `lcFieldValid` |
| `bc-field-rating` | `BcRate` | `lcFieldBlur` |
| `bc-field-rating` | `BcRate` | `lcFieldChange` |
| `bc-field-rating` | `BcRate` | `lcFieldClear` |
| `bc-field-rating` | `BcRate` | `lcFieldFocus` |
| `bc-field-rating` | `BcRate` | `lcFieldInvalid` |
| `bc-field-rating` | `BcRate` | `lcFieldValid` |
| `bc-field-richtext` | `BcRichText` | `lcFieldBlur` |
| `bc-field-richtext` | `BcRichText` | `lcFieldChange` |
| `bc-field-richtext` | `BcRichText` | `lcFieldClear` |
| `bc-field-richtext` | `BcRichText` | `lcFieldFocus` |
| `bc-field-richtext` | `BcRichText` | `lcFieldInvalid` |
| `bc-field-richtext` | `BcRichText` | `lcFieldValid` |
| `bc-field-select` | `BcSelect` | `lcDropdownClose` |
| `bc-field-select` | `BcSelect` | `lcDropdownOpen` |
| `bc-field-select` | `BcSelect` | `lcFieldBlur` |
| `bc-field-select` | `BcSelect` | `lcFieldChange` |
| `bc-field-select` | `BcSelect` | `lcFieldClear` |
| `bc-field-select` | `BcSelect` | `lcFieldFocus` |
| `bc-field-select` | `BcSelect` | `lcFieldInvalid` |
| `bc-field-select` | `BcSelect` | `lcFieldValid` |
| `bc-field-select` | `BcSelect` | `lcOptionCreate` |
| `bc-field-select` | `BcSelect` | `lcOptionsError` |
| `bc-field-select` | `BcSelect` | `lcOptionsLoad` |
| `bc-field-signature` | `BcSignature` | `lcFieldBlur` |
| `bc-field-signature` | `BcSignature` | `lcFieldChange` |
| `bc-field-signature` | `BcSignature` | `lcFieldClear` |
| `bc-field-signature` | `BcSignature` | `lcFieldFocus` |
| `bc-field-signature` | `BcSignature` | `lcFieldInvalid` |
| `bc-field-signature` | `BcSignature` | `lcFieldValid` |
| `bc-field-smalltext` | `BcTextarea` | `lcFieldBlur` |
| `bc-field-smalltext` | `BcTextarea` | `lcFieldChange` |
| `bc-field-smalltext` | `BcTextarea` | `lcFieldClear` |
| `bc-field-smalltext` | `BcTextarea` | `lcFieldFocus` |
| `bc-field-smalltext` | `BcTextarea` | `lcFieldInvalid` |
| `bc-field-smalltext` | `BcTextarea` | `lcFieldValid` |
| `bc-field-string` | `BcInput` | `lcFieldBlur` |
| `bc-field-string` | `BcInput` | `lcFieldChange` |
| `bc-field-string` | `BcInput` | `lcFieldClear` |
| `bc-field-string` | `BcInput` | `lcFieldFocus` |
| `bc-field-string` | `BcInput` | `lcFieldInvalid` |
| `bc-field-string` | `BcInput` | `lcFieldValid` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `lcFieldBlur` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `lcFieldChange` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `lcFieldClear` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `lcFieldFocus` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `lcFieldInvalid` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `lcFieldValid` |
| `bc-field-tags` | `BcTagsInput` | `lcFieldBlur` |
| `bc-field-tags` | `BcTagsInput` | `lcFieldChange` |
| `bc-field-tags` | `BcTagsInput` | `lcFieldClear` |
| `bc-field-tags` | `BcTagsInput` | `lcFieldFocus` |
| `bc-field-tags` | `BcTagsInput` | `lcFieldInvalid` |
| `bc-field-tags` | `BcTagsInput` | `lcFieldValid` |
| `bc-field-text` | `BcTextarea` | `lcFieldBlur` |
| `bc-field-text` | `BcTextarea` | `lcFieldChange` |
| `bc-field-text` | `BcTextarea` | `lcFieldClear` |
| `bc-field-text` | `BcTextarea` | `lcFieldFocus` |
| `bc-field-text` | `BcTextarea` | `lcFieldInvalid` |
| `bc-field-text` | `BcTextarea` | `lcFieldValid` |
| `bc-field-time` | `BcTimeField` | `lcFieldBlur` |
| `bc-field-time` | `BcTimeField` | `lcFieldChange` |
| `bc-field-time` | `BcTimeField` | `lcFieldClear` |
| `bc-field-time` | `BcTimeField` | `lcFieldFocus` |
| `bc-field-time` | `BcTimeField` | `lcFieldInvalid` |
| `bc-field-time` | `BcTimeField` | `lcFieldValid` |
| `bc-field-toggle` | `BcToggle` | `lcFieldBlur` |
| `bc-field-toggle` | `BcToggle` | `lcFieldChange` |
| `bc-field-toggle` | `BcToggle` | `lcFieldClear` |
| `bc-field-toggle` | `BcToggle` | `lcFieldFocus` |
| `bc-field-toggle` | `BcToggle` | `lcFieldInvalid` |
| `bc-field-toggle` | `BcToggle` | `lcFieldValid` |
| `bc-header` | `BcHeader` | `lcActionClick` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanCardCreate` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanCardDelete` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanCardMove` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanCardUpdate` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanColumnAdd` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanColumnDelete` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanColumnRename` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanColumnReorder` |
| `bc-kanban-board` | `BcKanbanBoard` | `kanbanError` |
| `bc-lookup-modal` | `BcLookupModal` | `lcError` |
| `bc-lookup-modal` | `BcLookupModal` | `lcLookupClose` |
| `bc-lookup-modal` | `BcLookupModal` | `lcLookupSelect` |
| `bc-search` | `BcSearch` | `lcError` |
| `bc-search` | `BcSearch` | `lcSearch` |
| `bc-timeline` | `BcTimelineSocial` | `lcError` |
| `bc-toast` | `BcToast` | `lcDialogClose` |
| `bc-view-activity` | `BcViewActivity` | `lcError` |
| `bc-view-calendar` | `BcCalendar` | `lcError` |
| `bc-view-editor` | `BcViewEditor` | `viewChanged` |
| `bc-view-form` | `BcForm` | `lcError` |
| `bc-view-form` | `BcForm` | `lcFormSubmit` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttAfterTaskDrag` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttBeforeLinkAdd` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttBeforeLinkDelete` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttBeforeLinkUpdate` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttBeforeTaskAdd` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttBeforeTaskDelete` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttBeforeTaskDrag` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttBeforeTaskMove` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttBeforeTaskUpdate` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttColumnResize` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttContextMenu` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttEmptyClick` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttError` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttLinkAdd` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttLinkClick` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttLinkDblClick` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttLinkDelete` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttLinkDrag` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttLinkUpdate` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttParsed` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttReady` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttRendered` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttRowResize` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttScaleClick` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttScroll` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskAdd` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskClick` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskClosed` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskDblClick` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskDelete` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskDrag` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskMove` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskOpened` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskRowClick` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskSelect` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskUnselect` |
| `bc-view-gantt` | `BcViewGantt` | `lcGanttTaskUpdate` |
| `bc-view-kanban` | `BcKanban` | `lcError` |
| `bc-view-kanban` | `BcKanban` | `lcKanbanColumnReorder` |
| `bc-view-kanban` | `BcKanban` | `lcKanbanMove` |
| `bc-view-list` | `BcList` | `lcRowSelect` |
| `bc-view-map` | `BcViewMap` | `lcError` |
| `bc-view-report` | `BcViewReport` | `lcError` |
| `bc-view-tree` | `BcTree` | `lcError` |

### Missing Methods in Mitosis (mapped components)

> Note: Mitosis components use callback props (e.g., `onChange`, `onFocus`) instead of
> imperative methods. Stencil `@Method()` declarations represent imperative APIs that
> need to be re-implemented as either callback props or ref-based method exposure in Mitosis.

| Stencil Component | Mitosis Component | Missing Method |
|-------------------|-------------------|----------------|
| `bc-badge` | `BcBadge` | `getValue` |
| `bc-badge` | `BcBadge` | `setValue` |
| `bc-chart-area` | `BcChartArea` | `exportImage` |
| `bc-chart-area` | `BcChartArea` | `getChartInstance` |
| `bc-chart-area` | `BcChartArea` | `refresh` |
| `bc-chart-area` | `BcChartArea` | `resize` |
| `bc-chart-area` | `BcChartArea` | `setData` |
| `bc-chart-area` | `BcChartArea` | `updateData` |
| `bc-chart-bar` | `BcChartBar` | `exportImage` |
| `bc-chart-bar` | `BcChartBar` | `getChartInstance` |
| `bc-chart-bar` | `BcChartBar` | `refresh` |
| `bc-chart-bar` | `BcChartBar` | `resize` |
| `bc-chart-bar` | `BcChartBar` | `setData` |
| `bc-chart-bar` | `BcChartBar` | `updateData` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `exportImage` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `getChartInstance` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `refresh` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `resize` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `setData` |
| `bc-chart-boxplot` | `BcChartBoxplot` | `updateData` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `exportImage` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `getChartInstance` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `refresh` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `resize` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `setData` |
| `bc-chart-candlestick` | `BcChartCandlestick` | `updateData` |
| `bc-chart-custom` | `BcChartCustom` | `exportImage` |
| `bc-chart-custom` | `BcChartCustom` | `getChartInstance` |
| `bc-chart-custom` | `BcChartCustom` | `refresh` |
| `bc-chart-custom` | `BcChartCustom` | `resize` |
| `bc-chart-custom` | `BcChartCustom` | `setData` |
| `bc-chart-custom` | `BcChartCustom` | `setOption` |
| `bc-chart-custom` | `BcChartCustom` | `updateData` |
| `bc-chart-funnel` | `BcChartFunnel` | `exportImage` |
| `bc-chart-funnel` | `BcChartFunnel` | `getChartInstance` |
| `bc-chart-funnel` | `BcChartFunnel` | `refresh` |
| `bc-chart-funnel` | `BcChartFunnel` | `resize` |
| `bc-chart-funnel` | `BcChartFunnel` | `setData` |
| `bc-chart-funnel` | `BcChartFunnel` | `updateData` |
| `bc-chart-gauge` | `BcChartGauge` | `exportImage` |
| `bc-chart-gauge` | `BcChartGauge` | `getChartInstance` |
| `bc-chart-gauge` | `BcChartGauge` | `refresh` |
| `bc-chart-gauge` | `BcChartGauge` | `resize` |
| `bc-chart-gauge` | `BcChartGauge` | `setData` |
| `bc-chart-gauge` | `BcChartGauge` | `updateData` |
| `bc-chart-graph` | `BcChartGraph` | `exportImage` |
| `bc-chart-graph` | `BcChartGraph` | `getChartInstance` |
| `bc-chart-graph` | `BcChartGraph` | `refresh` |
| `bc-chart-graph` | `BcChartGraph` | `resize` |
| `bc-chart-graph` | `BcChartGraph` | `setData` |
| `bc-chart-graph` | `BcChartGraph` | `updateData` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `exportImage` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `getChartInstance` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `refresh` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `resize` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `setData` |
| `bc-chart-heatmap` | `BcChartHeatmap` | `updateData` |
| `bc-chart-kpi` | `BcKpi` | `refresh` |
| `bc-chart-kpi` | `BcKpi` | `updateData` |
| `bc-chart-line` | `BcChartLine` | `exportImage` |
| `bc-chart-line` | `BcChartLine` | `getChartInstance` |
| `bc-chart-line` | `BcChartLine` | `refresh` |
| `bc-chart-line` | `BcChartLine` | `resize` |
| `bc-chart-line` | `BcChartLine` | `setData` |
| `bc-chart-line` | `BcChartLine` | `updateData` |
| `bc-chart-mixed` | `BcChartMixed` | `exportImage` |
| `bc-chart-mixed` | `BcChartMixed` | `getChartInstance` |
| `bc-chart-mixed` | `BcChartMixed` | `refresh` |
| `bc-chart-mixed` | `BcChartMixed` | `resize` |
| `bc-chart-mixed` | `BcChartMixed` | `setData` |
| `bc-chart-mixed` | `BcChartMixed` | `updateData` |
| `bc-chart-parallel` | `BcChartParallel` | `exportImage` |
| `bc-chart-parallel` | `BcChartParallel` | `getChartInstance` |
| `bc-chart-parallel` | `BcChartParallel` | `refresh` |
| `bc-chart-parallel` | `BcChartParallel` | `resize` |
| `bc-chart-parallel` | `BcChartParallel` | `setData` |
| `bc-chart-parallel` | `BcChartParallel` | `updateData` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `exportImage` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `getChartInstance` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `refresh` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `resize` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `setData` |
| `bc-chart-pictorialbar` | `BcChartPictorialBar` | `updateData` |
| `bc-chart-pie` | `BcChartPie` | `exportImage` |
| `bc-chart-pie` | `BcChartPie` | `getChartInstance` |
| `bc-chart-pie` | `BcChartPie` | `refresh` |
| `bc-chart-pie` | `BcChartPie` | `resize` |
| `bc-chart-pie` | `BcChartPie` | `setData` |
| `bc-chart-pie` | `BcChartPie` | `updateData` |
| `bc-chart-polar` | `BcChartPolar` | `exportImage` |
| `bc-chart-polar` | `BcChartPolar` | `getChartInstance` |
| `bc-chart-polar` | `BcChartPolar` | `refresh` |
| `bc-chart-polar` | `BcChartPolar` | `resize` |
| `bc-chart-polar` | `BcChartPolar` | `setData` |
| `bc-chart-polar` | `BcChartPolar` | `updateData` |
| `bc-chart-progress` | `BcProgress` | `refresh` |
| `bc-chart-progress` | `BcProgress` | `updateData` |
| `bc-chart-radar` | `BcChartRadar` | `exportImage` |
| `bc-chart-radar` | `BcChartRadar` | `getChartInstance` |
| `bc-chart-radar` | `BcChartRadar` | `refresh` |
| `bc-chart-radar` | `BcChartRadar` | `resize` |
| `bc-chart-radar` | `BcChartRadar` | `setData` |
| `bc-chart-radar` | `BcChartRadar` | `updateData` |
| `bc-chart-sankey` | `BcChartSankey` | `exportImage` |
| `bc-chart-sankey` | `BcChartSankey` | `getChartInstance` |
| `bc-chart-sankey` | `BcChartSankey` | `refresh` |
| `bc-chart-sankey` | `BcChartSankey` | `resize` |
| `bc-chart-sankey` | `BcChartSankey` | `setData` |
| `bc-chart-sankey` | `BcChartSankey` | `updateData` |
| `bc-chart-scatter` | `BcChartScatter` | `exportImage` |
| `bc-chart-scatter` | `BcChartScatter` | `getChartInstance` |
| `bc-chart-scatter` | `BcChartScatter` | `refresh` |
| `bc-chart-scatter` | `BcChartScatter` | `resize` |
| `bc-chart-scatter` | `BcChartScatter` | `setData` |
| `bc-chart-scatter` | `BcChartScatter` | `updateData` |
| `bc-chart-scorecard` | `BcChartScorecard` | `exportImage` |
| `bc-chart-scorecard` | `BcChartScorecard` | `getChartInstance` |
| `bc-chart-scorecard` | `BcChartScorecard` | `refresh` |
| `bc-chart-scorecard` | `BcChartScorecard` | `resize` |
| `bc-chart-scorecard` | `BcChartScorecard` | `setData` |
| `bc-chart-scorecard` | `BcChartScorecard` | `updateData` |
| `bc-chart-sunburst` | `BcChartSunburst` | `exportImage` |
| `bc-chart-sunburst` | `BcChartSunburst` | `getChartInstance` |
| `bc-chart-sunburst` | `BcChartSunburst` | `refresh` |
| `bc-chart-sunburst` | `BcChartSunburst` | `resize` |
| `bc-chart-sunburst` | `BcChartSunburst` | `setData` |
| `bc-chart-sunburst` | `BcChartSunburst` | `updateData` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `exportImage` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `getChartInstance` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `refresh` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `resize` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `setData` |
| `bc-chart-themeriver` | `BcChartThemeRiver` | `updateData` |
| `bc-chart-tree` | `BcChartTree` | `exportImage` |
| `bc-chart-tree` | `BcChartTree` | `getChartInstance` |
| `bc-chart-tree` | `BcChartTree` | `refresh` |
| `bc-chart-tree` | `BcChartTree` | `resize` |
| `bc-chart-tree` | `BcChartTree` | `setData` |
| `bc-chart-tree` | `BcChartTree` | `updateData` |
| `bc-chart-treemap` | `BcChartTreemap` | `exportImage` |
| `bc-chart-treemap` | `BcChartTreemap` | `getChartInstance` |
| `bc-chart-treemap` | `BcChartTreemap` | `refresh` |
| `bc-chart-treemap` | `BcChartTreemap` | `resize` |
| `bc-chart-treemap` | `BcChartTreemap` | `setData` |
| `bc-chart-treemap` | `BcChartTreemap` | `updateData` |
| `bc-datatable` | `BcDataTable` | `clearSelection` |
| `bc-datatable` | `BcDataTable` | `exportCSV` |
| `bc-datatable` | `BcDataTable` | `getData` |
| `bc-datatable` | `BcDataTable` | `getSelected` |
| `bc-datatable` | `BcDataTable` | `goToPage` |
| `bc-datatable` | `BcDataTable` | `refresh` |
| `bc-datatable` | `BcDataTable` | `scrollToRow` |
| `bc-datatable` | `BcDataTable` | `selectAll` |
| `bc-datatable` | `BcDataTable` | `setData` |
| `bc-datatable` | `BcDataTable` | `sortBy` |
| `bc-dialog-alert` | `BcAlertDialog` | `closeDialog` |
| `bc-dialog-alert` | `BcAlertDialog` | `openDialog` |
| `bc-dialog-confirm` | `BcAlertDialog` | `closeDialog` |
| `bc-dialog-confirm` | `BcAlertDialog` | `openDialog` |
| `bc-dialog-modal` | `BcDialog` | `closeDialog` |
| `bc-dialog-modal` | `BcDialog` | `isOpen` |
| `bc-dialog-modal` | `BcDialog` | `openDialog` |
| `bc-dialog-prompt` | `BcPrompt` | `closeDialog` |
| `bc-dialog-prompt` | `BcPrompt` | `openDialog` |
| `bc-dialog-wizard` | `BcWizard` | `closeDialog` |
| `bc-dialog-wizard` | `BcWizard` | `getCurrentStep` |
| `bc-dialog-wizard` | `BcWizard` | `goToStep` |
| `bc-dialog-wizard` | `BcWizard` | `nextStep` |
| `bc-dialog-wizard` | `BcWizard` | `openDialog` |
| `bc-dialog-wizard` | `BcWizard` | `prevStep` |
| `bc-field-barcode` | `BcBarcode` | `blurField` |
| `bc-field-barcode` | `BcBarcode` | `clear` |
| `bc-field-barcode` | `BcBarcode` | `clearError` |
| `bc-field-barcode` | `BcBarcode` | `focusField` |
| `bc-field-barcode` | `BcBarcode` | `getValue` |
| `bc-field-barcode` | `BcBarcode` | `isDirty` |
| `bc-field-barcode` | `BcBarcode` | `isTouched` |
| `bc-field-barcode` | `BcBarcode` | `reset` |
| `bc-field-barcode` | `BcBarcode` | `setError` |
| `bc-field-barcode` | `BcBarcode` | `setValue` |
| `bc-field-barcode` | `BcBarcode` | `validate` |
| `bc-field-checkbox` | `BcCheckbox` | `blurField` |
| `bc-field-checkbox` | `BcCheckbox` | `clear` |
| `bc-field-checkbox` | `BcCheckbox` | `clearError` |
| `bc-field-checkbox` | `BcCheckbox` | `focusField` |
| `bc-field-checkbox` | `BcCheckbox` | `getValue` |
| `bc-field-checkbox` | `BcCheckbox` | `isDirty` |
| `bc-field-checkbox` | `BcCheckbox` | `isTouched` |
| `bc-field-checkbox` | `BcCheckbox` | `reset` |
| `bc-field-checkbox` | `BcCheckbox` | `setError` |
| `bc-field-checkbox` | `BcCheckbox` | `setValue` |
| `bc-field-checkbox` | `BcCheckbox` | `validate` |
| `bc-field-code` | `BcCode` | `blurField` |
| `bc-field-code` | `BcCode` | `clear` |
| `bc-field-code` | `BcCode` | `clearError` |
| `bc-field-code` | `BcCode` | `focusField` |
| `bc-field-code` | `BcCode` | `getValue` |
| `bc-field-code` | `BcCode` | `isDirty` |
| `bc-field-code` | `BcCode` | `isTouched` |
| `bc-field-code` | `BcCode` | `reset` |
| `bc-field-code` | `BcCode` | `setError` |
| `bc-field-code` | `BcCode` | `setValue` |
| `bc-field-code` | `BcCode` | `validate` |
| `bc-field-color` | `BcColorPicker` | `blurField` |
| `bc-field-color` | `BcColorPicker` | `clear` |
| `bc-field-color` | `BcColorPicker` | `clearError` |
| `bc-field-color` | `BcColorPicker` | `focusField` |
| `bc-field-color` | `BcColorPicker` | `getValue` |
| `bc-field-color` | `BcColorPicker` | `isDirty` |
| `bc-field-color` | `BcColorPicker` | `isTouched` |
| `bc-field-color` | `BcColorPicker` | `reset` |
| `bc-field-color` | `BcColorPicker` | `setError` |
| `bc-field-color` | `BcColorPicker` | `setValue` |
| `bc-field-color` | `BcColorPicker` | `validate` |
| `bc-field-currency` | `BcNumber` | `blurField` |
| `bc-field-currency` | `BcNumber` | `clear` |
| `bc-field-currency` | `BcNumber` | `clearError` |
| `bc-field-currency` | `BcNumber` | `focusField` |
| `bc-field-currency` | `BcNumber` | `getValue` |
| `bc-field-currency` | `BcNumber` | `isDirty` |
| `bc-field-currency` | `BcNumber` | `isTouched` |
| `bc-field-currency` | `BcNumber` | `reset` |
| `bc-field-currency` | `BcNumber` | `setError` |
| `bc-field-currency` | `BcNumber` | `setValue` |
| `bc-field-currency` | `BcNumber` | `validate` |
| `bc-field-date` | `BcDateField` | `blurField` |
| `bc-field-date` | `BcDateField` | `clear` |
| `bc-field-date` | `BcDateField` | `clearError` |
| `bc-field-date` | `BcDateField` | `focusField` |
| `bc-field-date` | `BcDateField` | `getValue` |
| `bc-field-date` | `BcDateField` | `isDirty` |
| `bc-field-date` | `BcDateField` | `isTouched` |
| `bc-field-date` | `BcDateField` | `reset` |
| `bc-field-date` | `BcDateField` | `setError` |
| `bc-field-date` | `BcDateField` | `setValue` |
| `bc-field-date` | `BcDateField` | `validate` |
| `bc-field-datetime` | `BcDateTime` | `blurField` |
| `bc-field-datetime` | `BcDateTime` | `clear` |
| `bc-field-datetime` | `BcDateTime` | `clearError` |
| `bc-field-datetime` | `BcDateTime` | `focusField` |
| `bc-field-datetime` | `BcDateTime` | `getValue` |
| `bc-field-datetime` | `BcDateTime` | `isDirty` |
| `bc-field-datetime` | `BcDateTime` | `isTouched` |
| `bc-field-datetime` | `BcDateTime` | `reset` |
| `bc-field-datetime` | `BcDateTime` | `setError` |
| `bc-field-datetime` | `BcDateTime` | `setValue` |
| `bc-field-datetime` | `BcDateTime` | `validate` |
| `bc-field-decimal` | `BcNumber` | `blurField` |
| `bc-field-decimal` | `BcNumber` | `clear` |
| `bc-field-decimal` | `BcNumber` | `clearError` |
| `bc-field-decimal` | `BcNumber` | `focusField` |
| `bc-field-decimal` | `BcNumber` | `getValue` |
| `bc-field-decimal` | `BcNumber` | `isDirty` |
| `bc-field-decimal` | `BcNumber` | `isTouched` |
| `bc-field-decimal` | `BcNumber` | `reset` |
| `bc-field-decimal` | `BcNumber` | `setError` |
| `bc-field-decimal` | `BcNumber` | `setValue` |
| `bc-field-decimal` | `BcNumber` | `validate` |
| `bc-field-duration` | `BcFieldDuration` | `blurField` |
| `bc-field-duration` | `BcFieldDuration` | `clear` |
| `bc-field-duration` | `BcFieldDuration` | `clearError` |
| `bc-field-duration` | `BcFieldDuration` | `focusField` |
| `bc-field-duration` | `BcFieldDuration` | `getValue` |
| `bc-field-duration` | `BcFieldDuration` | `isDirty` |
| `bc-field-duration` | `BcFieldDuration` | `isTouched` |
| `bc-field-duration` | `BcFieldDuration` | `reset` |
| `bc-field-duration` | `BcFieldDuration` | `setError` |
| `bc-field-duration` | `BcFieldDuration` | `setValue` |
| `bc-field-duration` | `BcFieldDuration` | `validate` |
| `bc-field-dynlink` | `BcFieldDynLink` | `blurField` |
| `bc-field-dynlink` | `BcFieldDynLink` | `clear` |
| `bc-field-dynlink` | `BcFieldDynLink` | `clearError` |
| `bc-field-dynlink` | `BcFieldDynLink` | `focusField` |
| `bc-field-dynlink` | `BcFieldDynLink` | `getValue` |
| `bc-field-dynlink` | `BcFieldDynLink` | `isDirty` |
| `bc-field-dynlink` | `BcFieldDynLink` | `isTouched` |
| `bc-field-dynlink` | `BcFieldDynLink` | `reset` |
| `bc-field-dynlink` | `BcFieldDynLink` | `setError` |
| `bc-field-dynlink` | `BcFieldDynLink` | `setValue` |
| `bc-field-dynlink` | `BcFieldDynLink` | `validate` |
| `bc-field-file` | `BcFile` | `blurField` |
| `bc-field-file` | `BcFile` | `clear` |
| `bc-field-file` | `BcFile` | `clearError` |
| `bc-field-file` | `BcFile` | `focusField` |
| `bc-field-file` | `BcFile` | `getValue` |
| `bc-field-file` | `BcFile` | `isDirty` |
| `bc-field-file` | `BcFile` | `isTouched` |
| `bc-field-file` | `BcFile` | `reset` |
| `bc-field-file` | `BcFile` | `setError` |
| `bc-field-file` | `BcFile` | `setValue` |
| `bc-field-file` | `BcFile` | `validate` |
| `bc-field-float` | `BcNumber` | `blurField` |
| `bc-field-float` | `BcNumber` | `clear` |
| `bc-field-float` | `BcNumber` | `clearError` |
| `bc-field-float` | `BcNumber` | `focusField` |
| `bc-field-float` | `BcNumber` | `getValue` |
| `bc-field-float` | `BcNumber` | `isDirty` |
| `bc-field-float` | `BcNumber` | `isTouched` |
| `bc-field-float` | `BcNumber` | `reset` |
| `bc-field-float` | `BcNumber` | `setError` |
| `bc-field-float` | `BcNumber` | `setValue` |
| `bc-field-float` | `BcNumber` | `validate` |
| `bc-field-geo` | `BcGeo` | `blurField` |
| `bc-field-geo` | `BcGeo` | `clear` |
| `bc-field-geo` | `BcGeo` | `clearError` |
| `bc-field-geo` | `BcGeo` | `focusField` |
| `bc-field-geo` | `BcGeo` | `getValue` |
| `bc-field-geo` | `BcGeo` | `isDirty` |
| `bc-field-geo` | `BcGeo` | `isTouched` |
| `bc-field-geo` | `BcGeo` | `reset` |
| `bc-field-geo` | `BcGeo` | `setError` |
| `bc-field-geo` | `BcGeo` | `setValue` |
| `bc-field-geo` | `BcGeo` | `validate` |
| `bc-field-html` | `BcRichText` | `blurField` |
| `bc-field-html` | `BcRichText` | `clear` |
| `bc-field-html` | `BcRichText` | `clearError` |
| `bc-field-html` | `BcRichText` | `focusField` |
| `bc-field-html` | `BcRichText` | `getValue` |
| `bc-field-html` | `BcRichText` | `isDirty` |
| `bc-field-html` | `BcRichText` | `isTouched` |
| `bc-field-html` | `BcRichText` | `reset` |
| `bc-field-html` | `BcRichText` | `setError` |
| `bc-field-html` | `BcRichText` | `setValue` |
| `bc-field-html` | `BcRichText` | `validate` |
| `bc-field-image` | `BcFile` | `blurField` |
| `bc-field-image` | `BcFile` | `clear` |
| `bc-field-image` | `BcFile` | `clearError` |
| `bc-field-image` | `BcFile` | `focusField` |
| `bc-field-image` | `BcFile` | `getValue` |
| `bc-field-image` | `BcFile` | `isDirty` |
| `bc-field-image` | `BcFile` | `isTouched` |
| `bc-field-image` | `BcFile` | `reset` |
| `bc-field-image` | `BcFile` | `setError` |
| `bc-field-image` | `BcFile` | `setValue` |
| `bc-field-image` | `BcFile` | `validate` |
| `bc-field-integer` | `BcNumber` | `blurField` |
| `bc-field-integer` | `BcNumber` | `clear` |
| `bc-field-integer` | `BcNumber` | `clearError` |
| `bc-field-integer` | `BcNumber` | `focusField` |
| `bc-field-integer` | `BcNumber` | `getValue` |
| `bc-field-integer` | `BcNumber` | `isDirty` |
| `bc-field-integer` | `BcNumber` | `isTouched` |
| `bc-field-integer` | `BcNumber` | `reset` |
| `bc-field-integer` | `BcNumber` | `setError` |
| `bc-field-integer` | `BcNumber` | `setValue` |
| `bc-field-integer` | `BcNumber` | `validate` |
| `bc-field-json` | `BcJson` | `blurField` |
| `bc-field-json` | `BcJson` | `clear` |
| `bc-field-json` | `BcJson` | `clearError` |
| `bc-field-json` | `BcJson` | `focusField` |
| `bc-field-json` | `BcJson` | `getValue` |
| `bc-field-json` | `BcJson` | `isDirty` |
| `bc-field-json` | `BcJson` | `isTouched` |
| `bc-field-json` | `BcJson` | `reset` |
| `bc-field-json` | `BcJson` | `setError` |
| `bc-field-json` | `BcJson` | `setValue` |
| `bc-field-json` | `BcJson` | `validate` |
| `bc-field-link` | `BcCombobox` | `blurField` |
| `bc-field-link` | `BcCombobox` | `clear` |
| `bc-field-link` | `BcCombobox` | `clearError` |
| `bc-field-link` | `BcCombobox` | `focusField` |
| `bc-field-link` | `BcCombobox` | `getValue` |
| `bc-field-link` | `BcCombobox` | `isDirty` |
| `bc-field-link` | `BcCombobox` | `isTouched` |
| `bc-field-link` | `BcCombobox` | `reset` |
| `bc-field-link` | `BcCombobox` | `setError` |
| `bc-field-link` | `BcCombobox` | `setValue` |
| `bc-field-link` | `BcCombobox` | `validate` |
| `bc-field-markdown` | `BcRichText` | `blurField` |
| `bc-field-markdown` | `BcRichText` | `clear` |
| `bc-field-markdown` | `BcRichText` | `clearError` |
| `bc-field-markdown` | `BcRichText` | `focusField` |
| `bc-field-markdown` | `BcRichText` | `getValue` |
| `bc-field-markdown` | `BcRichText` | `isDirty` |
| `bc-field-markdown` | `BcRichText` | `isTouched` |
| `bc-field-markdown` | `BcRichText` | `reset` |
| `bc-field-markdown` | `BcRichText` | `setError` |
| `bc-field-markdown` | `BcRichText` | `setValue` |
| `bc-field-markdown` | `BcRichText` | `validate` |
| `bc-field-morph` | `BcFieldMorph` | `blurField` |
| `bc-field-morph` | `BcFieldMorph` | `clear` |
| `bc-field-morph` | `BcFieldMorph` | `clearError` |
| `bc-field-morph` | `BcFieldMorph` | `focusField` |
| `bc-field-morph` | `BcFieldMorph` | `getValue` |
| `bc-field-morph` | `BcFieldMorph` | `isDirty` |
| `bc-field-morph` | `BcFieldMorph` | `isTouched` |
| `bc-field-morph` | `BcFieldMorph` | `reset` |
| `bc-field-morph` | `BcFieldMorph` | `setError` |
| `bc-field-morph` | `BcFieldMorph` | `setValue` |
| `bc-field-morph` | `BcFieldMorph` | `validate` |
| `bc-field-multicheck` | `BcCheckbox` | `blurField` |
| `bc-field-multicheck` | `BcCheckbox` | `clear` |
| `bc-field-multicheck` | `BcCheckbox` | `clearError` |
| `bc-field-multicheck` | `BcCheckbox` | `focusField` |
| `bc-field-multicheck` | `BcCheckbox` | `getValue` |
| `bc-field-multicheck` | `BcCheckbox` | `isDirty` |
| `bc-field-multicheck` | `BcCheckbox` | `isTouched` |
| `bc-field-multicheck` | `BcCheckbox` | `reset` |
| `bc-field-multicheck` | `BcCheckbox` | `setError` |
| `bc-field-multicheck` | `BcCheckbox` | `setValue` |
| `bc-field-multicheck` | `BcCheckbox` | `validate` |
| `bc-field-password` | `BcInput` | `blurField` |
| `bc-field-password` | `BcInput` | `clear` |
| `bc-field-password` | `BcInput` | `clearError` |
| `bc-field-password` | `BcInput` | `focusField` |
| `bc-field-password` | `BcInput` | `getValue` |
| `bc-field-password` | `BcInput` | `isDirty` |
| `bc-field-password` | `BcInput` | `isTouched` |
| `bc-field-password` | `BcInput` | `reset` |
| `bc-field-password` | `BcInput` | `setError` |
| `bc-field-password` | `BcInput` | `setValue` |
| `bc-field-password` | `BcInput` | `validate` |
| `bc-field-percent` | `BcNumber` | `blurField` |
| `bc-field-percent` | `BcNumber` | `clear` |
| `bc-field-percent` | `BcNumber` | `clearError` |
| `bc-field-percent` | `BcNumber` | `focusField` |
| `bc-field-percent` | `BcNumber` | `getValue` |
| `bc-field-percent` | `BcNumber` | `isDirty` |
| `bc-field-percent` | `BcNumber` | `isTouched` |
| `bc-field-percent` | `BcNumber` | `reset` |
| `bc-field-percent` | `BcNumber` | `setError` |
| `bc-field-percent` | `BcNumber` | `setValue` |
| `bc-field-percent` | `BcNumber` | `validate` |
| `bc-field-radio` | `BcSegmented` | `blurField` |
| `bc-field-radio` | `BcSegmented` | `clear` |
| `bc-field-radio` | `BcSegmented` | `clearError` |
| `bc-field-radio` | `BcSegmented` | `focusField` |
| `bc-field-radio` | `BcSegmented` | `getValue` |
| `bc-field-radio` | `BcSegmented` | `isDirty` |
| `bc-field-radio` | `BcSegmented` | `isTouched` |
| `bc-field-radio` | `BcSegmented` | `reset` |
| `bc-field-radio` | `BcSegmented` | `setError` |
| `bc-field-radio` | `BcSegmented` | `setValue` |
| `bc-field-radio` | `BcSegmented` | `validate` |
| `bc-field-rating` | `BcRate` | `blurField` |
| `bc-field-rating` | `BcRate` | `clear` |
| `bc-field-rating` | `BcRate` | `clearError` |
| `bc-field-rating` | `BcRate` | `focusField` |
| `bc-field-rating` | `BcRate` | `getValue` |
| `bc-field-rating` | `BcRate` | `isDirty` |
| `bc-field-rating` | `BcRate` | `isTouched` |
| `bc-field-rating` | `BcRate` | `reset` |
| `bc-field-rating` | `BcRate` | `setError` |
| `bc-field-rating` | `BcRate` | `setValue` |
| `bc-field-rating` | `BcRate` | `validate` |
| `bc-field-richtext` | `BcRichText` | `blurField` |
| `bc-field-richtext` | `BcRichText` | `clear` |
| `bc-field-richtext` | `BcRichText` | `clearError` |
| `bc-field-richtext` | `BcRichText` | `focusField` |
| `bc-field-richtext` | `BcRichText` | `getValue` |
| `bc-field-richtext` | `BcRichText` | `isDirty` |
| `bc-field-richtext` | `BcRichText` | `isTouched` |
| `bc-field-richtext` | `BcRichText` | `reset` |
| `bc-field-richtext` | `BcRichText` | `setError` |
| `bc-field-richtext` | `BcRichText` | `setValue` |
| `bc-field-richtext` | `BcRichText` | `validate` |
| `bc-field-select` | `BcSelect` | `blurField` |
| `bc-field-select` | `BcSelect` | `clear` |
| `bc-field-select` | `BcSelect` | `clearError` |
| `bc-field-select` | `BcSelect` | `close` |
| `bc-field-select` | `BcSelect` | `focusField` |
| `bc-field-select` | `BcSelect` | `getOptions` |
| `bc-field-select` | `BcSelect` | `getSelectedOptions` |
| `bc-field-select` | `BcSelect` | `getValue` |
| `bc-field-select` | `BcSelect` | `isDirty` |
| `bc-field-select` | `BcSelect` | `isTouched` |
| `bc-field-select` | `BcSelect` | `loadOptions` |
| `bc-field-select` | `BcSelect` | `open` |
| `bc-field-select` | `BcSelect` | `reloadOptions` |
| `bc-field-select` | `BcSelect` | `reset` |
| `bc-field-select` | `BcSelect` | `setError` |
| `bc-field-select` | `BcSelect` | `setOptions` |
| `bc-field-select` | `BcSelect` | `setValue` |
| `bc-field-select` | `BcSelect` | `validate` |
| `bc-field-signature` | `BcSignature` | `blurField` |
| `bc-field-signature` | `BcSignature` | `clear` |
| `bc-field-signature` | `BcSignature` | `clearError` |
| `bc-field-signature` | `BcSignature` | `focusField` |
| `bc-field-signature` | `BcSignature` | `getValue` |
| `bc-field-signature` | `BcSignature` | `isDirty` |
| `bc-field-signature` | `BcSignature` | `isTouched` |
| `bc-field-signature` | `BcSignature` | `reset` |
| `bc-field-signature` | `BcSignature` | `setError` |
| `bc-field-signature` | `BcSignature` | `setValue` |
| `bc-field-signature` | `BcSignature` | `validate` |
| `bc-field-smalltext` | `BcTextarea` | `blurField` |
| `bc-field-smalltext` | `BcTextarea` | `clear` |
| `bc-field-smalltext` | `BcTextarea` | `clearError` |
| `bc-field-smalltext` | `BcTextarea` | `focusField` |
| `bc-field-smalltext` | `BcTextarea` | `getValue` |
| `bc-field-smalltext` | `BcTextarea` | `isDirty` |
| `bc-field-smalltext` | `BcTextarea` | `isTouched` |
| `bc-field-smalltext` | `BcTextarea` | `reset` |
| `bc-field-smalltext` | `BcTextarea` | `setError` |
| `bc-field-smalltext` | `BcTextarea` | `setValue` |
| `bc-field-smalltext` | `BcTextarea` | `validate` |
| `bc-field-string` | `BcInput` | `blurField` |
| `bc-field-string` | `BcInput` | `clear` |
| `bc-field-string` | `BcInput` | `clearError` |
| `bc-field-string` | `BcInput` | `focusField` |
| `bc-field-string` | `BcInput` | `getValue` |
| `bc-field-string` | `BcInput` | `isDirty` |
| `bc-field-string` | `BcInput` | `isTouched` |
| `bc-field-string` | `BcInput` | `reset` |
| `bc-field-string` | `BcInput` | `setError` |
| `bc-field-string` | `BcInput` | `setValue` |
| `bc-field-string` | `BcInput` | `validate` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `blurField` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `clear` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `clearError` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `focusField` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `getValue` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `isDirty` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `isTouched` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `reset` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `setError` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `setValue` |
| `bc-field-tableselect` | `BcFieldTableSelect` | `validate` |
| `bc-field-tags` | `BcTagsInput` | `blurField` |
| `bc-field-tags` | `BcTagsInput` | `clear` |
| `bc-field-tags` | `BcTagsInput` | `clearError` |
| `bc-field-tags` | `BcTagsInput` | `focusField` |
| `bc-field-tags` | `BcTagsInput` | `getValue` |
| `bc-field-tags` | `BcTagsInput` | `isDirty` |
| `bc-field-tags` | `BcTagsInput` | `isTouched` |
| `bc-field-tags` | `BcTagsInput` | `reset` |
| `bc-field-tags` | `BcTagsInput` | `setError` |
| `bc-field-tags` | `BcTagsInput` | `setValue` |
| `bc-field-tags` | `BcTagsInput` | `validate` |
| `bc-field-text` | `BcTextarea` | `blurField` |
| `bc-field-text` | `BcTextarea` | `clear` |
| `bc-field-text` | `BcTextarea` | `clearError` |
| `bc-field-text` | `BcTextarea` | `focusField` |
| `bc-field-text` | `BcTextarea` | `getValue` |
| `bc-field-text` | `BcTextarea` | `isDirty` |
| `bc-field-text` | `BcTextarea` | `isTouched` |
| `bc-field-text` | `BcTextarea` | `reset` |
| `bc-field-text` | `BcTextarea` | `setError` |
| `bc-field-text` | `BcTextarea` | `setValue` |
| `bc-field-text` | `BcTextarea` | `validate` |
| `bc-field-time` | `BcTimeField` | `blurField` |
| `bc-field-time` | `BcTimeField` | `clear` |
| `bc-field-time` | `BcTimeField` | `clearError` |
| `bc-field-time` | `BcTimeField` | `focusField` |
| `bc-field-time` | `BcTimeField` | `getValue` |
| `bc-field-time` | `BcTimeField` | `isDirty` |
| `bc-field-time` | `BcTimeField` | `isTouched` |
| `bc-field-time` | `BcTimeField` | `reset` |
| `bc-field-time` | `BcTimeField` | `setError` |
| `bc-field-time` | `BcTimeField` | `setValue` |
| `bc-field-time` | `BcTimeField` | `validate` |
| `bc-field-toggle` | `BcToggle` | `blurField` |
| `bc-field-toggle` | `BcToggle` | `clear` |
| `bc-field-toggle` | `BcToggle` | `clearError` |
| `bc-field-toggle` | `BcToggle` | `focusField` |
| `bc-field-toggle` | `BcToggle` | `getValue` |
| `bc-field-toggle` | `BcToggle` | `isDirty` |
| `bc-field-toggle` | `BcToggle` | `isTouched` |
| `bc-field-toggle` | `BcToggle` | `reset` |
| `bc-field-toggle` | `BcToggle` | `setError` |
| `bc-field-toggle` | `BcToggle` | `setValue` |
| `bc-field-toggle` | `BcToggle` | `validate` |
| `bc-kanban-board` | `BcKanbanBoard` | `refresh` |
| `bc-progress` | `BcProgress` | `getValue` |
| `bc-progress` | `BcProgress` | `setValue` |
| `bc-section` | `BcSection` | `collapse` |
| `bc-section` | `BcSection` | `expand` |
| `bc-section` | `BcSection` | `toggle` |
| `bc-tabs` | `BcTabs` | `getActiveIndex` |
| `bc-tabs` | `BcTabs` | `selectTab` |
| `bc-timeline` | `BcTimelineSocial` | `refresh` |
| `bc-toast` | `BcToast` | `dismiss` |
| `bc-toast` | `BcToast` | `show` |
| `bc-view-activity` | `BcViewActivity` | `refresh` |
| `bc-view-calendar` | `BcCalendar` | `refresh` |
| `bc-view-editor` | `BcViewEditor` | `refresh` |
| `bc-viewer-audio` | `BcViewerAudio` | `refresh` |
| `bc-viewer-document` | `BcViewerDocument` | `refresh` |
| `bc-viewer-image` | `BcViewerImage` | `refresh` |
| `bc-viewer-instagram` | `BcViewerInstagram` | `refresh` |
| `bc-viewer-pdf` | `BcViewerPdf` | `refresh` |
| `bc-viewer-tiktok` | `BcViewerTiktok` | `refresh` |
| `bc-viewer-video` | `BcViewerVideo` | `refresh` |
| `bc-viewer-youtube` | `BcViewerYoutube` | `refresh` |
| `bc-view-form` | `BcForm` | `refresh` |
| `bc-view-gantt` | `BcViewGantt` | `addLink` |
| `bc-view-gantt` | `BcViewGantt` | `addMarker` |
| `bc-view-gantt` | `BcViewGantt` | `addTask` |
| `bc-view-gantt` | `BcViewGantt` | `calculateDuration` |
| `bc-view-gantt` | `BcViewGantt` | `calculateEndDate` |
| `bc-view-gantt` | `BcViewGantt` | `clearAll` |
| `bc-view-gantt` | `BcViewGantt` | `clearRedoStack` |
| `bc-view-gantt` | `BcViewGantt` | `clearUndoStack` |
| `bc-view-gantt` | `BcViewGantt` | `close` |
| `bc-view-gantt` | `BcViewGantt` | `closeAll` |
| `bc-view-gantt` | `BcViewGantt` | `collapse` |
| `bc-view-gantt` | `BcViewGantt` | `dateFromPos` |
| `bc-view-gantt` | `BcViewGantt` | `deleteLink` |
| `bc-view-gantt` | `BcViewGantt` | `deleteMarker` |
| `bc-view-gantt` | `BcViewGantt` | `deleteTask` |
| `bc-view-gantt` | `BcViewGantt` | `eachTask` |
| `bc-view-gantt` | `BcViewGantt` | `expand` |
| `bc-view-gantt` | `BcViewGantt` | `getChildren` |
| `bc-view-gantt` | `BcViewGantt` | `getGlobalTaskIndex` |
| `bc-view-gantt` | `BcViewGantt` | `getLink` |
| `bc-view-gantt` | `BcViewGantt` | `getLinkCount` |
| `bc-view-gantt` | `BcViewGantt` | `getLinks` |
| `bc-view-gantt` | `BcViewGantt` | `getMarker` |
| `bc-view-gantt` | `BcViewGantt` | `getNext` |
| `bc-view-gantt` | `BcViewGantt` | `getParent` |
| `bc-view-gantt` | `BcViewGantt` | `getPrev` |
| `bc-view-gantt` | `BcViewGantt` | `getRedoStack` |
| `bc-view-gantt` | `BcViewGantt` | `getScrollState` |
| `bc-view-gantt` | `BcViewGantt` | `getSelectedId` |
| `bc-view-gantt` | `BcViewGantt` | `getSiblings` |
| `bc-view-gantt` | `BcViewGantt` | `getState` |
| `bc-view-gantt` | `BcViewGantt` | `getTask` |
| `bc-view-gantt` | `BcViewGantt` | `getTaskCount` |
| `bc-view-gantt` | `BcViewGantt` | `getTaskIndex` |
| `bc-view-gantt` | `BcViewGantt` | `getTasks` |
| `bc-view-gantt` | `BcViewGantt` | `getUndoStack` |
| `bc-view-gantt` | `BcViewGantt` | `getVisibleTaskCount` |
| `bc-view-gantt` | `BcViewGantt` | `hasChild` |
| `bc-view-gantt` | `BcViewGantt` | `isSelected` |
| `bc-view-gantt` | `BcViewGantt` | `isTaskExists` |
| `bc-view-gantt` | `BcViewGantt` | `isTaskVisible` |
| `bc-view-gantt` | `BcViewGantt` | `moveTask` |
| `bc-view-gantt` | `BcViewGantt` | `open` |
| `bc-view-gantt` | `BcViewGantt` | `openAll` |
| `bc-view-gantt` | `BcViewGantt` | `parseData` |
| `bc-view-gantt` | `BcViewGantt` | `posFromDate` |
| `bc-view-gantt` | `BcViewGantt` | `redo` |
| `bc-view-gantt` | `BcViewGantt` | `refresh` |
| `bc-view-gantt` | `BcViewGantt` | `rerender` |
| `bc-view-gantt` | `BcViewGantt` | `scrollTo` |
| `bc-view-gantt` | `BcViewGantt` | `selectTask` |
| `bc-view-gantt` | `BcViewGantt` | `serialize` |
| `bc-view-gantt` | `BcViewGantt` | `setData` |
| `bc-view-gantt` | `BcViewGantt` | `setSizes` |
| `bc-view-gantt` | `BcViewGantt` | `showDate` |
| `bc-view-gantt` | `BcViewGantt` | `showTask` |
| `bc-view-gantt` | `BcViewGantt` | `sortBy` |
| `bc-view-gantt` | `BcViewGantt` | `undo` |
| `bc-view-gantt` | `BcViewGantt` | `unselectTask` |
| `bc-view-gantt` | `BcViewGantt` | `updateLink` |
| `bc-view-gantt` | `BcViewGantt` | `updateTask` |
| `bc-view-kanban` | `BcKanban` | `refresh` |
| `bc-view-list` | `BcList` | `refresh` |
| `bc-view-map` | `BcViewMap` | `refresh` |
| `bc-view-report` | `BcViewReport` | `refresh` |
| `bc-view-tree` | `BcTree` | `refresh` |

### Top 5 Components with Most Missing Props/Events/Methods

| Component | Missing Props | Missing Events | Missing Methods | Total |
|-----------|--------------|----------------|-----------------|-------|
| `bc-view-gantt` | 14 | 37 | 61 | 112 |
| `bc-field-select` | 23 | 11 | 18 | 52 |
| `bc-kanban-board` | 33 | 9 | 1 | 43 |
| `bc-field-file` | 17 | 6 | 11 | 34 |
| `bc-field-geo` | 15 | 6 | 11 | 32 |

## Statistics

| Stat | Value |
|------|-------|
| Total Stencil components | 127 |
| Total Mitosis components | 176 |
| Mapped Stencil components (have Mitosis equiv) | 103 |
| Missing Stencil components (no Mitosis equiv) | 24 |
| New in Mitosis components | 85 |
| Stencil coverage | 81.1% |
| Total Stencil props | 1472 |
| Total Stencil events | 452 |
| Total Stencil methods | 673 |
| Total Mitosis props | 1343 |
| Total Mitosis events | 247 |
| Missing props (in mapped components) | 731 |
| Missing events (in mapped components) | 349 |
| Missing methods (all Stencil methods, pattern gap) | 673 |

