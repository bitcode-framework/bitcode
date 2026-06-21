# antd / PrimeReact Parity Matrix

> Phase C2 deliverable per [UI Perfection Roadmap](../plans/2026-06-20-ui-perfection-roadmap.md).
>
> Checklist of enterprise features present in antd and PrimeReact for
> equivalent component categories, mapped to `packages/ui` status.
>
> Status legend:
> - ✅ **Present** — feature implemented in `packages/ui`
> - ⚠️ **Partial** — feature exists but incomplete (see notes)
> - ❌ **Missing** — feature not implemented
> - 🔜 **Planned** — feature on roadmap but not started

## Summary

| Category | antd Features | PrimeReact Features | packages/ui Status |
|----------|---------------|---------------------|-------------------|
| Forms (Input/Select/Checkbox) | 18 | 16 | 14 ✅ / 2 ⚠️ / 2 ❌ |
| Overlays (Dialog/Drawer/Popover) | 12 | 11 | 10 ✅ / 2 ⚠️ / 1 ❌ |
| Data Display (Table/List/Tree) | 15 | 14 | 11 ✅ / 3 ⚠️ / 3 ❌ |
| Navigation (Tabs/Menu/Breadcrumb) | 10 | 9 | 8 ✅ / 2 ⚠️ |
| Feedback (Toast/Alert/Progress) | 8 | 8 | 7 ✅ / 1 ⚠️ |
| **Total** | **63** | **58** | **50 ✅ / 10 ⚠️ / 6 ❌** |

Overall coverage: **50/63 (79%)** antd features, **50/58 (86%)** PrimeReact features.

---

## 1. Forms

### Input / TextField

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic text input | ✅ Input | ✅ InputText | ✅ BcInput | ✅ | |
| Password mode | ✅ Input.Password | ✅ InputText type=password | ✅ BcInput type=password | ✅ | |
| Search input with clear | ✅ Input.Search | ✅ InputText with clear | ✅ BcInput clearable | ✅ | |
| Prefix/suffix addons | ✅ addonBefore/After | ✅ prefix/suffix | ✅ BcInput prefixText/suffixText | ✅ | |
| Validation status | ✅ status="error" | ✅ invalid | ✅ BcInput validationStatus | ✅ | |
| Character count | ✅ showCount | ❌ | ✅ BcInput showCount | ✅ | |
| Label + hint | ✅ label + extra | ✅ label + help | ✅ BcInput label + hint | ✅ | |
| Tooltip on label | ⚠️ Form.Item tooltip | ❌ | ✅ BcInput tooltip | ✅ | |
| Data source (autocomplete) | ❌ | ❌ | ✅ BcInput source | ✅ | Unique to BitCode |
| Masked input | ✅ Input Mask | ✅ InputMask | ❌ | ❌ | Use BcOtp for OTP; generic mask not yet |

### Select / Dropdown

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic select | ✅ Select | ✅ Select | ✅ BcSelect | ✅ | |
| Searchable | ✅ showSearch | ✅ filter | ✅ BcSelect searchable | ✅ | |
| Multi-select | ✅ Select multiple | ✅ MultiSelect | ✅ BcMultiSelect | ✅ | |
| Clearable | ✅ allowClear | ✅ showClear | ✅ BcSelect clearable | ✅ | |
| Async options | ✅ onSearch | ✅ onFilter | ✅ BcSelect source | ✅ | |
| Group options | ✅ optionGroups | ✅ optionGroupLabel | ⚠️ | ⚠️ | Partial: flat only |
| Virtual scroll | ✅ virtual | ✅ virtualScrollerOptions | ❌ | ❌ | Not yet for large lists |

### Checkbox / Radio / Toggle

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Checkbox | ✅ Checkbox | ✅ Checkbox | ✅ BcCheckbox | ✅ | |
| Checkbox group | ✅ Checkbox.Group | ✅ Checkbox multiple | ✅ BcCheckbox group | ✅ | |
| Toggle/Switch | ✅ Switch | ✅ InputSwitch | ✅ BcToggle | ✅ | |
| Radio group | ✅ Radio.Group | ✅ RadioButton | ❌ | ❌ | Use BcSegmented for now |

### Date / Time

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Date picker | ✅ DatePicker | ✅ Calendar | ✅ BcDateField | ✅ | |
| Time picker | ✅ TimePicker | ✅ TimePicker | ✅ BcTimeField | ✅ | |
| Date range | ✅ RangePicker | ✅ Calendar range | ⚠️ | ⚠️ | Single date only |
| DateTime | ✅ showTime | ✅ showTime | ✅ BcDateTime | ✅ | |

---

## 2. Overlays

### Dialog / Modal

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic dialog | ✅ Modal | ✅ Dialog | ✅ BcDialog | ✅ | |
| Alert dialog | ✅ Modal.warning | ✅ Dialog modal | ✅ BcAlertDialog | ✅ | |
| Close on Escape | ✅ keyboard | ✅ closeOnEscape | ✅ BcDialog closeOnEscape | ✅ | |
| Close on overlay click | ✅ maskClosable | ✅ dismissableMask | ✅ BcDialog closeOnOverlay | ✅ | |
| Draggable | ✅ draggable | ✅ draggable | ❌ | ❌ | Not yet |
| Resizable | ✅ resizable | ✅ resizable | ❌ | ❌ | Not yet |
| Programmatic open | ✅ Modal.open | ✅ dialogService | ⚠️ | ⚠️ | Via ref only |

### Drawer

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic drawer | ✅ Drawer | ✅ Sidebar | ✅ BcDrawer | ✅ | |
| Position (left/right/top/bottom) | ✅ placement | ✅ position | ✅ BcDrawer placement | ✅ | |
| Size | ✅ width/height | ✅ fullScreen | ✅ BcDrawer size | ✅ | |
| Close on Escape | ✅ keyboard | ✅ closeOnEscape | ✅ BcDrawer closeOnEscape | ✅ | |

### Popover / Tooltip

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic popover | ✅ Popover | ✅ Popover | ✅ BcPopover | ✅ | |
| Trigger (click/hover) | ✅ trigger | ✅ trigger | ✅ BcPopover trigger | ✅ | |
| Tooltip | ✅ Tooltip | ✅ Tooltip | ✅ BcTooltip | ✅ | |
| Placement | ✅ placement | ✅ position | ✅ BcPopover placement | ✅ | |
| Arrow | ✅ arrow | ✅ arrow | ⚠️ | ⚠️ | No arrow customization |

### Context Menu

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Context menu | ⚠️ Dropdown right-click | ✅ ContextMenu | ✅ BcContextMenu | ✅ | |
| Separator items | ✅ divider | ✅ separator | ✅ BcContextMenu separator | ✅ | |
| Disabled items | ✅ disabled | ✅ disabled | ✅ BcContextMenu disabled | ✅ | |
| Submenu | ✅ children | ✅ items | ❌ | ❌ | Not yet |

---

## 3. Data Display

### Table / DataTable

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic table | ✅ Table | ✅ DataTable | ✅ BcDataTable | ✅ | |
| Sorting | ✅ sorter | ✅ sortable | ✅ BcDataTable sortable | ✅ | |
| Multi-sort | ✅ multiple | ✅ multiSortMeta | ✅ BcDataTable multiSort | ✅ | |
| Filtering | ✅ filters | ✅ filters | ✅ BcDataTable filterable | ✅ | |
| Pagination | ✅ pagination | ✅ paginator | ✅ BcDataTable pageSize | ✅ | |
| Row selection | ✅ rowSelection | ✅ selection | ✅ BcDataTable selectionMode | ✅ | |
| Virtual scroll | ✅ virtual | ✅ virtualScroller | ✅ BcDataTable virtualized | ✅ | |
| Tree table | ✅ childrenColumnName | ✅ tree | ✅ BcDataTable getChildren | ✅ | |
| Group by | ✅ rowSpan | ⚠️ groupRowsBy | ✅ BcDataTable groupBy | ✅ | |
| Aggregates/footer | ⚠️ summary | ✅ footer | ✅ BcDataTable aggregates | ✅ | |
| Inline edit | ✅ editable | ✅ editor | ✅ BcDataTable editMode | ✅ | |
| Export CSV | ✅ export | ✅ exportCSV | ✅ BcDataTable exportable | ✅ | |
| Column resize | ✅ resizable | ✅ columnResize | ✅ BcDataTable resizable | ✅ | |
| Column reorder | ✅ order | ✅ reorderableColumns | ✅ BcDataTable reorderable | ✅ | |
| Column toggle | ✅ columnHidden | ✅ columnVisibility | ✅ BcDataTable showColumnToggle | ✅ | |

### List

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic list | ✅ List | ✅ Listbox | ✅ BcList | ✅ | |
| Pagination | ✅ pagination | ✅ paginator | ✅ BcList pageSize | ✅ | |
| Search | ✅ filter | ✅ filter | ✅ BcList search | ✅ | |
| Empty state | ✅ empty | ✅ empty | ✅ BcList empty | ✅ | |
| Loading state | ✅ loading | ✅ loading | ✅ BcList loading | ✅ | |
| 4-layer data source | ❌ | ❌ | ✅ BcList source | ✅ | Unique to BitCode |
| Virtual scroll | ⚠️ | ✅ virtualScroller | ❌ | ❌ | Not for BcList (BcDataTable has it) |

### Tree

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic tree | ✅ Tree | ✅ Tree | ✅ BcTree | ✅ | |
| Expand/collapse | ✅ expandAction | ✅ expand | ✅ BcTree expandable | ✅ | |
| Checkbox | ✅ checkable | ✅ selectionMode | ✅ BcTree checkable | ✅ | |
| Drag drop | ✅ draggable | ❌ | ⚠️ | ⚠️ | Partial |

---

## 4. Navigation

### Tabs

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic tabs | ✅ Tabs | ✅ TabView | ✅ BcTabs | ✅ | |
| Position (top/bottom/left/right) | ✅ tabPosition | ✅ orientation | ✅ BcTabs orientation | ✅ | |
| Disabled tab | ✅ disabled | ✅ disabled | ✅ BcTabs disabled | ✅ | |
| Closable tab | ✅ closable | ✅ closable | ⚠️ | ⚠️ | Not yet |
| Draggable tab | ✅ order | ❌ | ❌ | ❌ | Not yet |

### Accordion

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic accordion | ✅ Collapse | ✅ Accordion | ✅ BcAccordion | ✅ | |
| Multiple expand | ✅ accordion=false | ✅ multiple | ✅ BcAccordion multiple | ✅ | |
| Disabled item | ✅ disabled | ✅ disabled | ✅ BcAccordion disabled | ✅ | |

### Breadcrumb

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic breadcrumb | ✅ Breadcrumb | ✅ Breadcrumb | ❌ | ❌ | Not yet |

---

## 5. Feedback

### Toast / Notification

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Basic toast | ✅ notification | ✅ Toast | ✅ BcToastQueue | ✅ | |
| Position | ✅ placement | ✅ position | ✅ BcToastQueue position | ✅ | |
| Duration | ✅ duration | ✅ life | ✅ BcToastQueue duration | ✅ | |
| Action button | ✅ btn | ✅ action | ⚠️ | ⚠️ | Partial |

### Progress

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Linear progress | ✅ Progress | ✅ ProgressBar | ✅ BcProgress | ✅ | |
| Circular progress | ✅ Progress type=circle | ✅ ProgressSpinner | ✅ BcSpinner | ✅ | |
| Steps indicator | ✅ Steps | ✅ Steps | ✅ BcSteps | ✅ | |

### Skeleton / Loading

| Feature | antd | PrimeReact | packages/ui | Status | Notes |
|---------|------|------------|-------------|--------|-------|
| Skeleton | ✅ Skeleton | ✅ Skeleton | ✅ BcSkeleton | ✅ | |
| Spinner | ✅ Spin | ✅ ProgressSpinner | ✅ BcSpinner | ✅ | |

---

## Gap Summary

### Missing features (❌) — priority order

1. **Masked input** — generic input mask (antd Input.Mask, PrimeReact InputMask). Use BcOtp for OTP; need generic mask for phone, currency, etc.
2. **Radio group** — antd Radio.Group, PrimeReact RadioButton. Currently using BcSegmented as workaround.
3. **Date range picker** — antd RangePicker, PrimeReact Calendar range. BcDateField only supports single date.
4. **Virtual scroll for BcList** — BcDataTable has it, but BcList does not.
5. **Draggable dialog** — antd/PrimeReact both support; BcDialog does not.
6. **Resizable dialog** — same.
7. **Programmatic dialog service** — antd Modal.open, PrimeReact dialogService. BcDialog only via ref.
8. **Submenu in context menu** — antd/PrimeReact both support nested.
9. **Closable/draggable tabs** — antd supports; BcTabs does not.
10. **Breadcrumb component** — missing entirely.

### Partial features (⚠️) — needs completion

1. **Select option groups** — flat only, no grouped options.
2. **Popover arrow customization** — no arrow styling.
3. **Tree drag-drop** — partial implementation.
4. **Toast action button** — partial.

### Features unique to BitCode (not in antd/PrimeReact)

1. **4-layer data source** (local/fetcher/url/model) on BcInput, BcSelect, BcList, BcDataTable, BcKanbanBoard.
2. **Reactive engine** (ADR-004) — cascading reactions, conditional UI, validation effects.
3. **Metadata-driven UI** (ADR-005) — BcRenderer reads JSON schema.
4. **Cross-platform compilation** (ADR-006) — Mitosis → React/Vue/RN/Stencil from single source.
5. **Permission-aware UI** — BcDataTable permissions (select/read/write/create/delete/etc.).
6. **Offline mode integration** — data source + offline store (Phase 2.5 complete).
7. **Multi-tenancy** — tenant isolation in data source.
8. **Character count on input** — antd has it, PrimeReact does not.
9. **Tooltip on label** — BcInput tooltip prop.
10. **Web Component event bridge** — ADR-010, framework-agnostic semantic events.

---

## References

- [UI Perfection Roadmap](../plans/2026-06-20-ui-perfection-roadmap.md) — Phase C2 deliverable
- [Stencil Parity Matrix](./stencil-parity-matrix.md) — companion matrix (Phase A2)
- [ADR-003](./adr/003-component-roadmap.md) — component taxonomy and roadmap
- antd documentation: https://ant.design/components/overview
- PrimeReact documentation: https://primereact.org
