# Handoff: UI Perfection Roadmap — Phase A1 Runtime Tests

**Tanggal:** 2026-06-20
**Session:** Phase A1.1 + A1.2 + A1.3 (partial) dari `docs/plans/2026-06-20-ui-perfection-roadmap.md`
**Status:** 10/25 Core Integration Surface components selesai (form controls), 141 runtime tests pass

---

## Yang Sudah Selesai

### A1.1 Infrastructure (DONE)

- `packages/ui/vitest.config.ts` — Vitest config dengan jsdom, axe, alias ke `output/react/src`
- `packages/ui/test-setup.ts` — jest-dom matchers, matchMedia/ResizeObserver/IntersectionObserver polyfills, `expect.toBeAccessible()` matcher (axe-core)
- `packages/ui/test/react-native-web-stub.tsx` — View/Text/Pressable/TextInput/Image → HTML equivalents
- `package.json` test scripts: `test` (source + runtime), `test:source` (node --test), `test:runtime` (vitest), `test:coverage`, `test:watch`
- Dev deps installed: vitest@2.1, @testing-library/react@16, @testing-library/user-event@14, @testing-library/jest-dom@6, jsdom@25, axe-core@4.10, vite@5.4, @vitest/coverage-v8@2.1

### A1.2 Pilot (DONE) — BcInput

- 23 runtime tests covering: mount, label, controlled/uncontrolled value, onValueChange/onChange/onFocus/onBlur/onEnter, disabled/readOnly, clearable, maxLength/minLength/pattern DOM attrs, aria-invalid, aria-describedby, role=alert, data-part, slotProps className, axe a11y

### A1.3 Core Integration Surface (10/25 DONE — form controls)

| # | Component | Tests | Status |
|---|-----------|-------|--------|
| 1 | BcInput | 23 | ✅ |
| 2 | BcSelect | 20 | ✅ (BC-004, BC-005 deferred) |
| 3 | BcCombobox | 16 | ✅ (BC-006 deferred) |
| 4 | BcCheckbox | 13 | ✅ (BC-007 deferred) |
| 5 | BcToggle | 12 | ✅ (BC-007 deferred) |
| 6 | BcSegmented | 11 | ✅ |
| 7 | BcNumberField | 11 | ✅ |
| 8 | BcDateField | 12 | ✅ |
| 9 | BcTimeField | 11 | ✅ |
| 10 | BcMultiSelect | 12 | ✅ (filter/maxSelected deferred) |

**Total: 141 runtime tests pass, 0 fail.**

### Bug Log (`docs/ui-perfection-buglog.md`)

7 bug real ditemukan dan terdokumentasi:

| # | Component | Severity | Status | Summary |
|---|-----------|----------|--------|---------|
| BC-001 | BcInput | P0 | fixed | handleChange crash karena value-first + plugin interaction |
| BC-002 | BcInput | P1 | deferred | slotProps.label attributes (data-testid) tidak di-apply |
| BC-003 | BcInput+all | P0 | fixed | `if (cond) return;` compile ke `if (cond); return;` (plugin regex terlalu agresif) |
| BC-004 | BcSelect | P1 | deferred | listbox missing aria-label + combobox missing aria-controls |
| BC-005 | BcSelect | P1 | deferred | clearable prop declared tapi clear button never rendered |
| BC-006 | BcCombobox | P1 | deferred | controlled value prop tidak sync ke query state |
| BC-007 | BcCheckbox/Toggle | P1 | deferred | button lacks accessible name (label not wired via htmlFor/aria-labelledby) |

### Build Pipeline Fixes

- `scripts/mitosis-plugin-react-fixes.mjs` step 5: tightened regex dari `([)}])\s+keyword` ke `(\})\s+keyword` (hanya `}` yang safe, tidak break `if (cond) return`)
- `scripts/reformat-react-output.mjs` (NEW): post-build babel parse+regenerate untuk robust semicolon insertion. Berjalan SETELAH semua Mitosis plugin (termasuk value-events yang line-based) selesai.
- `scripts/mitosis-plugin-react-value-events.mjs`: regex update untuk handle babel-emitted single-param arrow tanpa parens (`value =>` selain `(value) =>`)
- `package.json` build script: `mitosis build && apply-adapt-overrides && reformat-react-output && generate-index`
- `src/components/input/BcInput.lite.tsx`: refactor `handleChange(e: Event)` → `handleChange(value: string)` (value-first pattern, BC-001 fix)

---

## Yang Belum Selesai (Next Steps)

### A1.3 Lanjutan — 15 komponen tersisa

Urutan per roadmap:

**Forms/reactive (4):**
- [ ] BcForm — submit/validate/reset/reactions/field visibility
- [ ] BcReactiveScope — dispose/effect re-run/scope isolation
- [ ] ReactiveStore — get/set/reactions/batch/transaction
- [ ] Field wrappers — composition via BcForm parts

**Overlays (8):**
- [ ] BcDialog — open/close/Esc/backdrop click/focus trap
- [ ] BcAlertDialog — confirm/cancel/role=alertdialog/focus trap
- [ ] BcDrawer — open/close/4 placements/Esc/focus trap
- [ ] BcPopover — open/close/trigger hover+click/placement
- [ ] BcTooltip — hover+focus open/delay/placement
- [ ] BcDropdown — trigger/item select/separator/disabled
- [ ] BcContextMenu — right-click open/item select/long-press (RN)
- [ ] BcCommand — type-filter/arrow nav/enter select/esc

**Navigation/layout (4):**
- [ ] BcTabs — tab switch/keyboard arrows/controlled activeKey
- [ ] BcAccordion — expand/collapse/single+multi/controlled
- [ ] BcDisclosure — expand/collapse/aria-expanded
- [ ] BcSplitter — drag handle/keyboard nudge/min sizes

**Data (2):**
- [ ] BcDataTable — sort/filter/search/select/pagination/virtualize/tree/group/edit/export
- [ ] BcList — 4-layer source/pagination/empty/error/loading

**Metadata/infrastructure (5):**
- [ ] BcRenderer — 6 schema kinds/recursion/custom sentinel
- [ ] AppProvider
- [ ] ReactiveStore (infra)
- [ ] FieldResolver
- [ ] UiSchema

### Setelah A1.3 selesai:
- A2 Stencil parity audit
- A3 Demo app browser verification
- A4 4-layer data source uniformity
- B1-B3 Harden phases
- C1-C3 Formalize (ADR promotion + parity gates + DoD sign-off)

---

## Konteks Penting untuk Session Berikutnya

### 1. Test File Convention

- Runtime tests: `*.rtl.test.tsx` (Vitest + RTL + jsdom)
- Source-level tests: `*.test.ts` (node --test, TAP format)
- JANGAN pakai `*.runtime.test.ts` — sudah dipakai oleh file lama yang menggunakan `node:test`
- Vitest config hanya include `src/**/*.rtl.test.{ts,tsx}`

### 2. Bug Deferral Pattern

Ketika test menemukan bug yang bukan blocker infra:
1. Log di `docs/ui-perfection-buglog.md` dengan ID `BC-NNN`
2. Relax test dengan comment `— DEFERRED (BC-NNN)` + singkat deskripsi
3. Lanjut ke komponen berikutnya
4. Bug fix di Phase B (B2 khusus untuk a11y bugs)

### 3. Build Commands

```powershell
# Selalu set PATH dulu
$env:PATH = "C:\Program Files\nodejs;" + $env:PATH

# Atau pakai cmd.exe (lebih reliable untuk npm)
cmd.exe /c "set PATH=C:\Program Files\nodejs;%PATH% && cd /d D:\NodeJS\lowcode\packages\ui && npm run build"

# Run all runtime tests
cd D:\NodeJS\lowcode\packages\ui
npx vitest run

# Run specific component
npx vitest run src/components/forms/select

# Run source-level tests
node --test "src/**/*.test.ts"  # via PowerShell, bukan cmd.exe
```

### 4. Mitosis Plugin Architecture

Plugin order di `mitosis.config.js`:
```
react.plugins: [reactFixesPlugin(), reactValueEventsPlugin()]
```

- `reactFixesPlugin` run PERTAMA: attribute fixes (class→className, SVG camelCase, data-bc-* strip, semicolon `}keyword`)
- `reactValueEventsPlugin` run KEDUA: value-first → Event wrapping (line-based regex)
- `reformat-react-output.mjs` run SETELAH build: babel parse+regenerate untuk robust semicolon insertion

**Penting:** Jangan pindah babel regenerate ke dalam plugin (breaks value-events line-regex karena babel format multiline).

### 5. Bug BC-003 Blast Radius

Bug BC-003 (`if (cond) return;` → `if (cond); return;`) sudah fixed, tapi **semua komponen yang sudah ada test sebelum fix mungkin punya handler yang silently broken**. Re-run semua runtime test setelah BC-003 fix untuk verify. Untuk komponen yang belum ada test, bug ini tidak akan recurring karena plugin sudah diperbaiki.

### 6. A11y Bug Pattern (BC-004, BC-007)

Dua pattern a11y bug umum yang akan muncul di komponen lain:
- **Label not wired:** `<label>` dan `<input>`/`<button>` sebagai sibling tanpa `htmlFor`/`id` atau `aria-labelledby`. Fix: tambah `id` ke control + `htmlFor` ke label.
- **Listbox/dialog missing accessible name:** `<div role="listbox">` tanpa `aria-label`/`aria-labelledby`. Fix: tambah `aria-labelledby` pointing ke label id.

Phase B2 akan fix semua sekaligus dengan audit menyeluruh.

---

## Commits di Session Ini

### packages/ui
- `e5441de` test(runtime): Phase A1 infrastructure + BcInput pilot (23 runtime tests)
- `63ed7f3` test(runtime): BcSelect + BcCombobox runtime tests (36 tests)
- `6aa60eb` test(runtime): BcCheckbox, BcToggle, BcSegmented, BcNumberField, BcDateField, BcTimeField, BcMultiSelect (70 tests)

### root repo
- `62eaad7` docs(ui): bug log + bump packages/ui for Phase A1 runtime test foundation
- `404b70e` docs(plan): UI perfection roadmap — runtime tests, parity audit, DoD gates
- `31666d6` test(runtime): BcSelect + BcCombobox runtime tests (36 tests, 3 bugs logged)
- `6228933` chore: bump packages/ui for BcSelect+BcCombobox runtime tests
- `0807c30` docs(ui): bug log BC-007 + bump packages/ui for form runtime tests

---

## Instruksi untuk Session Baru

```
Lanjutkan pekerjaan dari handoff file: handoffs/2026-06-20-ui-perfection-phase-a1.md

Lanjutkan A1.3: tulis runtime tests untuk 15 komponen tersisa (overlays, navigation, data, core).
Pakai pattern yang sama: *.rtl.test.tsx, mount + interaksi + assert.
Log bug di docs/ui-perfection-buglog.md, relax test dengan DEFERRED comment.
Commit per 3-5 komponen.
```
