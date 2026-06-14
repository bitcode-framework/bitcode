# ADR-001: Mitosis as Component Compiler for BitCode UI

## Status: **Proposed**

## Date: 2026-06-05

## Context

BitCode Platform saat ini menggunakan **Stencil Web Components** (`@bitcode/components`) untuk UI layer. Ada 128 komponen dalam 12 grup yang berjalan di browser dan Tauri WebView.

### Masalah yang Ada

1. **React integration buruk** — Stencil Web Components tidak seamless dengan React. Event handling (`CustomEvent`), type definitions, dan HMR sering bermasalah di Vite + React setup.
2. **Tidak ada output ke framework lain** — React Native, Flutter, Vue, Qwik, Solid tidak didukung.
3. **Styling tidak konsisten** — CSS per-component files, bukan Tailwind/utility-first.
4. **Naming inconsistency** — `bc-field-string` (kebab-case) tidak natural di React (`<BcFieldString>`).

### Kebutuhan Masa Depan

- Support untuk **web** (React, Vue, Stencil, Qwik, Solid, Svelte)
- Support untuk **mobile** (React Native, Flutter, SwiftUI)
- **Tailwind CSS** sebagai styling system
- **Theme system** yang compatible dengan shadcn/ui ecosystem (tweakcn.com, dll)
- **RTL/LTR/BiDi** support
- **Enterprise-grade** events, methods, dan properties
- **Standalone** (bisa dipakai tanpa BitCode) tapi mudah integrate dengan BitCode framework

## Decision

### 1. Gunakan Mitosis sebagai Component Compiler

**Pilihan:** [Builder.io Mitosis](https://github.com/BuilderIO/mitosis) (13.8k stars, 253 releases, MIT license)

**Alasan:**
- Write once → compile ke 21+ output targets
- Production-proven: DB UX Design System (Deutsche Bahn) menggunakan Mitosis + Tailwind
- Stencil output target tersedia (backward compatibility dengan BitCode engine)
- React Native output target tersedia
- SwiftUI output target tersedia (bisa jadi base untuk Flutter generator)
- Plugin system untuk custom generators (Flutter)

**Output targets yang akan di-support:**

| Priority | Target | Kegunaan |
|----------|--------|----------|
| P0 | Stencil | Backward compatibility dengan Go engine (bc-* tags) |
| P0 | React | Solve pain point utama (Vite React integration) |
| P1 | Vue | Secondary web framework |
| P1 | React Native | Mobile option A |
| P2 | Qwik | Performance-first web |
| P2 | Solid | Reactive web |
| P3 | SwiftUI | iOS native |
| P3 | Flutter | Mobile option B (custom generator) |
| P4 | Svelte, Angular, Lit | Bonus targets |

### 2. Naming Convention: PascalCase

```
Source:    BcInput, BcChartLine, BcViewForm
React:     <BcInput />
Vue:       <BcInput />
Stencil:   <bc-input />         (auto kebab-case)
Flutter:   BcInput()
SwiftUI:   BcInput()
```

**Prefix:** `Bc` (bukan `BC`, bukan `bc-`)
**Pattern:** `Bc{Category}{Component}` — contoh: `BcInput`, `BcChartLine`, `BcViewForm`

### 3. Component Taxonomy (128 → ~85 components)

| Group | Mitosis Components | Covers |
|-------|-------------------|--------|
| Fields | BcInput, BcTextarea, BcRichText, BcCode, BcNumber, BcCheckbox, BcToggle, BcSelect, BcDateTime, BcFile, BcLookup, BcBarcode, BcColor, BcGeo, BcRating, BcJson, BcSignature | 36 Stencil → 17 |
| Layout | BcRow, BcColumn, BcSection, BcTabs, BcCard, BcHeader, BcSeparator, BcButtonGroup, BcHtmlBlock, BcSidebar | 10 → 10 |
| Views | BcList, BcForm, BcKanban, BcCalendar, BcGantt, BcMap, BcTree, BcReport, BcActivity, BcEditor | 10 → 10 |
| Charts | BcChart (type prop) | 26 → 1 |
| DataTable | BcDataTable, BcFilter, BcLookupModal | 5 → 3 |
| Dialogs | BcDialog, BcAlertDialog, BcPrompt, BcWizard, BcToast | 7 → 5 |
| Widgets | BcBadge, BcCopy, BcPhone, BcEmail, BcUrl, BcProgress, BcStatusBar, BcPriority, BcHandle, BcDomain, BcPlaceholder, BcSyncStatus | 12 → 12 |
| Media | BcViewerPdf, BcViewerImage, BcViewerDocument, BcViewerVideo, BcViewerAudio, BcViewerYoutube, BcViewerInstagram, BcViewerTiktok | 8 → 8 |
| Search | BcSearch, BcFavorites | 2 → 2 |
| Social | BcActivity, BcChatter, BcTimeline | 3 → 3 |
| Print | BcExport, BcPrint, BcReportLink | 3 → 3 |
| Navigation | BcBreadcrumb, BcPagination, BcMenu, BcDropdown | (new) |

**Total: ~85 components, 100% coverage dari 128 Stencil components.**

### 4. Data Source: 4-Layer Unified Type

```typescript
export type DataSource =
  | { type: 'local'; data: unknown[] }                              // Layer 1
  | { type: 'fetcher'; fetcher: (params: FetchParams) => Promise<FetchResult> }  // Layer 2
  | { type: 'url'; url: string; headers?: Record<string, string> }  // Layer 3
  | { type: 'model'; name: string; adapter?: ModelAdapter }         // Layer 4
```

- **Standalone** (Layer 1-3): Bisa dipakai tanpa BitCode
- **BitCode integration** (Layer 4): Via `ModelAdapter` yang di-register

### 5. Events: Framework-Native Callbacks

```typescript
// BUKAN: CustomEvent + dispatchEvent (DOM-only)
// TAPI: Framework-native callbacks

export interface BcInputProps {
  onChange?: (event: { name: string; value: string; oldValue: string }) => void;
  onFocus?: (event: { name: string; value: string }) => void;
  onBlur?: (event: { name: string; value: string; dirty: boolean; touched: boolean }) => void;
}
```

Mitosis compile ke:
- React: `onChange={(e) => ...}` (function prop)
- Vue: `@change="(e) => ..."` (emit)
- Stencil: `onLcFieldChange` (CustomEvent, backward compat)
- Flutter: `onChange: (e) => ...` (callback)

### 6. Methods: useTarget + Ref Pattern

```typescript
// Mitosis useTarget — per-framework method exposure
// React: useImperativeHandle + ref
// Vue: defineExpose
// Stencil: @Method decorator
// Flutter: direct method
```

### 7. BcSetup: Hybrid (Singleton + Context)

```typescript
// Mode 1: Singleton (simple standalone)
BcSetup.configure({ baseUrl: '/api', theme: 'dark' });

// Mode 2: Context (per-tree override)
<BcProvider context={{ theme: 'dark', locale: 'id' }}>
  <BcInput name="email" />
</BcProvider>

// Mode 3: Props (per-component override)
<BcInput name="email" theme="custom" />

// Resolution: Props → Context → Singleton
```

### 8. Theme System: shadcn/ui Standard (No Prefix)

```css
/* Compatible dengan tweakcn.com, shadcn themes, dll */
:root {
  --background: oklch(1 0 0);
  --primary: oklch(0.205 0 0);
  --radius: 0.625rem;
  --font-sans: "Inter", sans-serif;
}
.dark { ... }

@theme inline {
  --color-background: var(--background);
  --color-primary: var(--primary);
  /* ... */
}
```

**Theme installation:**
```bash
pnpm dlx shadcn@latest add https://tweakcn.com/r/themes/xxx
```

**Dark mode:** `.dark` class selector (shadcn standard)
**System detect:** `prefers-color-scheme` + JS sync

### 9. RTL/LTR/BiDi

- **Web:** CSS logical properties (`margin-inline-start`, `text-align: start`) + Tailwind `rtl:` variant
- **React Native:** `I18nManager.isRTL`
- **Flutter:** `Directionality` widget
- **SwiftUI:** `environment(\.layoutDirection, .rightToLeft)`

### 10. Tailwind CSS Integration

- **Web:** Tailwind v4 `@theme inline` (CSS variables → Tailwind utilities)
- **React Native:** NativeWind (Tailwind for RN)
- **Flutter:** Can't use Tailwind directly — use Material 3 / Cupertino themes generated from CSS tokens
- **SwiftUI:** Color assets generated from CSS tokens

### 11. Font Strategy (Inspired by bud.app)

bud.app menggunakan multiple custom fonts premium dengan fallback system. Kita adopt pattern yang sama:

```css
/* @bitcode/ui/themes/fonts.css */

/* Primary: Inter (free, universal) */
@font-face {
  font-family: 'BcSans';
  src: url('./fonts/Inter-Regular.woff2') format('woff2');
  font-weight: 400;
  font-style: normal;
  font-display: swap;
}
@font-face {
  font-family: 'BcSans';
  src: url('./fonts/Inter-Medium.woff2') format('woff2');
  font-weight: 500;
  font-style: normal;
  font-display: swap;
}
@font-face {
  font-family: 'BcSans';
  src: url('./fonts/Inter-SemiBold.woff2') format('woff2');
  font-weight: 600;
  font-style: normal;
  font-display: swap;
}
@font-face {
  font-family: 'BcSans';
  src: url('./fonts/Inter-Bold.woff2') format('woff2');
  font-weight: 700;
  font-style: normal;
  font-display: swap;
}

/* Fallback (system fonts) */
@font-face {
  font-family: 'BcSans Fallback';
  src: local('Arial');
  ascent-override: 90.44%;
  descent-override: 22.52%;
  line-gap-override: 0.0%;
  size-adjust: 107.12%;
}

/* Monospace: JetBrains Mono */
@font-face {
  font-family: 'BcMono';
  src: url('./fonts/JetBrainsMono-Regular.woff2') format('woff2');
  font-weight: 400;
  font-style: normal;
  font-display: swap;
}
```

**CSS Custom Properties:**
```css
:root {
  --font-sans: 'BcSans', 'BcSans Fallback', ui-sans-serif, system-ui, sans-serif;
  --font-mono: 'BcMono', ui-monospace, monospace;
  --font-serif: ui-serif, Georgia, serif;
}
```

**Key decisions:**
- **`font-display: swap`** — text langsung visible dengan fallback, font custom load di background
- **`size-adjust`** pada fallback — mencegah layout shift saat font load
- **Naming: `BcSans` / `BcMono`** — brand-agnostic, user bisa ganti font via CSS variable
- **Weights: 400, 500, 600, 700** — cukup untuk UI (tidak perlu 300/800/900)
- **Format: woff2** — best compression, universal browser support

**User bisa override font:**
```css
/* User theme dari tweakcn.com */
:root {
  --font-sans: 'Outfit', sans-serif;  /* Override default */
}
```

### 12. CSS Modules for Component-Specific Styles

bud.app menggunakan **CSS Modules** untuk component-specific styles alongside Tailwind utilities. Kita adopt pattern yang sama:

**Kapan pakai Tailwind vs CSS Modules:**

| Use Case | Approach | Contoh |
|----------|----------|--------|
| Layout, spacing, positioning | Tailwind utility | `class="flex items-center gap-2 p-4"` |
| Component-specific animations | CSS Module | `@keyframes bc-slide-down { ... }` |
| Complex pseudo-selectors | CSS Module | `.bc-input:focus-within .bc-input-icon { ... }` |
| Theme-aware custom styles | CSS Module + CSS vars | `background: var(--primary);` |
| Responsive breakpoints | Tailwind | `class="text-sm md:text-base lg:text-lg"` |

**Implementation:**

```tsx
// BcInput component — hybrid approach
// Tailwind untuk layout/utilities
// CSS Module untuk component-specific patterns

import styles from './BcInput.module.css';

export default function BcInput(props) {
  return (
    <div className="flex flex-col gap-1.5">  {/* Tailwind */}
      <input
        className={[styles.input, sizeClass, statusClass].join(' ')}  {/* CSS Module + Tailwind */}
      />
      <div className={styles.optionsDropdown}>  {/* CSS Module */}
        {/* ... */}
      </div>
    </div>
  );
}
```

```css
/* BcInput.module.css */
.input {
  /* Component-specific styles that are hard to express in Tailwind */
  transition: box-shadow 0.15s ease, border-color 0.15s ease;
}

.input:focus-within {
  box-shadow: 0 0 0 2px var(--ring);
}

.optionsDropdown {
  animation: bc-slide-down 0.15s ease-out;
}

@keyframes bc-slide-down {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}
```

**Key decisions:**
- **Tailwind sebagai default** — 80% styling pakai Tailwind utilities
- **CSS Module sebagai fallback** — untuk animations, complex selectors, component-specific patterns
- **Naming: `ComponentName.module.css`** — standard CSS Modules convention
- **Mitosis compatibility** — CSS Modules bisa di-import di semua output targets (React, Vue, Stencil)
- **No CSS-in-JS** — tidak pakai styled-components/emotion (bundle size, runtime overhead)

### 13. Build Tool: Vite (Universal) + Turbopack (Next.js)

bud.app menggunakan **Turbopack** (Next.js build tool). Untuk BitCode UI library:

**Decision: Vite sebagai primary build tool**

| Aspect | Vite | Turbopack |
|--------|------|-----------|
| Framework support | Universal (React, Vue, Svelte, dll) | Next.js only |
| Build speed | Very fast (esbuild + Rollup) | Very fast (Rust-based) |
| Library mode | ✅ Built-in (`build.lib`) | ❌ App-only |
| CSS Modules | ✅ Native | ✅ Native |
| Tailwind v4 | ✅ via `@tailwindcss/vite` | ✅ via PostCSS |
| HMR | ✅ Fast | ✅ Fast |
| Community | 70k+ stars | Next.js ecosystem |

**Why Vite:**
- Library mode untuk build `@bitcode/ui` package
- Universal — works dengan semua demo apps (React, Vue, Stencil)
- Tailwind v4 native plugin (`@tailwindcss/vite`)
- Tidak lock-in ke Next.js

**Demo apps bisa pakai framework apapun:**
```
demos/
├── react-app/    → Vite + React
├── vue-app/      → Vite + Vue
├── next-app/     → Next.js + Turbopack (optional)
└── stencil-app/  → Stencil CLI
```

**Key decisions:**
- **Vite untuk library build** — `vite build --lib`
- **Vite untuk demo apps** — React, Vue demos
- **Turbopack optional** — kalau user mau Next.js app, bisa pakai Turbopack
- **`@tailwindcss/vite`** — Tailwind v4 native Vite plugin (bukan PostCSS)

### 14. React Attribute Conversion: Mitosis Plugin (Option 2)

**Masalah:** Mitosis React generator tidak selalu convert HTML attributes ke React camelCase:
- `autofocus` → harusnya `autoFocus` (React)
- `readonly` → harusnya `readOnly` (React)
- `class` → harusnya `className` (React)
- `stroke-width` → harusnya `strokeWidth` (React)

**Opsi yang dipertimbangkan:**

| Option | Approach | Verdict |
|--------|----------|---------|
| 1. Post-build script | `sedreplace` setelah `mitosis build` | ✅ Work, tapi fragile |
| 2. Mitosis plugin | Custom plugin yang auto-convert | ✅ **DIPILIH** — cleanest |
| 3. React-style di source | Pakai `autoFocus` di `.lite.tsx` | ❌ Break Vue/Stencil output |

**Decision: Option 2 — Mitosis Plugin**

```js
// mitosis.config.js
import { reactAttributePlugin } from './plugins/react-attributes';

export default {
  // ...
  options: {
    react: {
      plugins: [reactAttributePlugin()],
    },
  },
};
```

```js
// plugins/react-attributes.ts
const HTML_TO_REACT: Record<string, string> = {
  'autofocus': 'autoFocus',
  'readonly': 'readOnly',
  'classname': 'className',
  'stroke-width': 'strokeWidth',
  'stroke-linecap': 'strokeLinecap',
  'stroke-linejoin': 'strokeLinejoin',
  'fill-rule': 'fillRule',
  'clip-rule': 'clipRule',
  'font-size': 'fontSize',
  'font-weight': 'fontWeight',
  'text-anchor': 'textAnchor',
  'dominant-baseline': 'dominantBaseline',
  'alignment-baseline': 'alignmentBaseline',
};

export function reactAttributePlugin() {
  return {
    code: {
      // Pre-process: convert HTML attributes to React camelCase
      pre: (code: string) => {
        let result = code;
        for (const [html, react] of Object.entries(HTML_TO_REACT)) {
          // Match attribute in JSX: autofocus= or autofocus={
          const regex = new RegExp(`\\b${html}(?==)`, 'g');
          result = result.replace(regex, react);
        }
        return result;
      },
    },
  };
}
```

**Key decisions:**
- **Plugin approach** — source tetap HTML standard, output React camelCase
- **Whitelist-based** — explicit mapping, tidak auto-convert (safety)
- **Extensible** — mudah tambah attribute baru kalau Mitosis miss
- **No impact ke Vue/Stencil** — plugin hanya jalan di React target
- **Fallback** — kalau plugin gagal, post-build script sebagai backup

**Milestone:**
- Phase 0 (sekarang): Post-build manual fix (sudah dilakukan)
- Phase 1 (saat 10+ components): Tulis plugin, automated
- Phase 2 (stable): Plugin jadi bagian dari build pipeline

### 15. Lynx (ReactLynx) — Native Mobile Target Only

**Status: ⛔ SUPERSEDED by [ADR-006](./006-universal-component-strategy.md) (2026-06-14)**

> **Lynx target dropped.** Mobile sekarang ditangani oleh React Native via NativeWind
> sebagai inti Tier A (lihat ADR-006 §6). Artifacts Lynx (`output/lynx/`,
> `scripts/lynx-generator.mjs`, `scripts/generate-lynx.mjs`, `demos/lynx-app/`) sudah
> dibersihkan. Section di bawah dipertahankan sebagai **historical record** untuk
> menjelaskan kenapa Lynx dipertimbangkan dan kenapa akhirnya tidak dipakai.

[Lynx](https://lynxjs.org/) adalah cross-platform native rendering engine dari ByteDance (TikTok). Production-proven dengan 2B+ users.

**Lynx = native mobile target ONLY. Web renderer TIDAK layak untuk interactive UI.**

| Platform | Target | Technology | Status |
|----------|--------|-----------|--------|
| **Web** | Browser | React/Vue/Stencil dari Mitosis | ✅ Sudah jalan |
| **Mobile** | iOS/Android | Lynx native dari Mitosis (custom generator) | 🧪 Phase candidate |
| **Desktop** | Windows/macOS/Linux | Tauri (existing) | ✅ Sudah ada |

**Lynx Web Renderer Limitations (verified 2026-06-05):**
- `bindtap` events trigger NYI (Not Yet Implemented) in WASM module
- `<input>` is XElement (requires extra deps, limited web support)
- `global-bind` explicitly "no Web" per docs
- Interactive elements (combobox, dropdown) not functional on web
- `NYI: profileStart/profileEnd` noise on every interaction

**Kenapa Lynx untuk mobile (bukan React Native/Flutter):**

| Aspect | Lynx | React Native | Flutter |
|--------|------|--------------|---------|
| Performance | ✅ C++ engine, dual-thread | ⚠️ JS bridge | ✅ Dart VM |
| Startup | ✅ Instant first-frame | ❌ Slow | ⚠️ Medium |
| React syntax | ✅ ReactLynx | ✅ React | ❌ Dart |
| Tailwind | ✅ @lynx-js/tailwind-preset | ✅ NativeWind | ❌ |
| Production | TikTok (2B+ users) | Meta | Google |

**Approach:**
1. **Phase 0-4:** Focus web outputs (React, Vue, Stencil) — sudah jalan
2. **Phase candidate:** Custom Mitosis → Lynx generator
   - Direct Mitosis AST/block renderer, not React output post-processing
   - Map HTML elements → Lynx elements (`<div>`→`<view>`, `<span>`/`<label>`/`<button>`→`<text>`, `<img>`→`<image>`)
   - Map events (`onClick`→`bindtap`, `onChange`/`onInput`→`bindinput`, `onFocus`→`bindfocus`, `onBlur`→`bindblur`)
   - Normalize Lynx input events (`e.detail.value` instead of DOM `e.target.value`)
   - Lower Mitosis state assignments to ReactLynx `useState` setters
   - Tailwind via `@lynx-js/tailwind-preset`

**Proof-of-concept evidence (2026-06-05):**
- `packages/ui/scripts/lynx-generator.mjs` implements standalone `componentToLynx()` / `blockToLynx()` rendering.
- `packages/ui/scripts/generate-lynx.mjs` generates both library output and the Lynx demo component from `BcInput.lite.tsx`.
- Verified commands in `packages/ui/demos/lynx-app`: `npm run generate:lynx`, `npm exec -- tsc --noEmit`, `npm run build`.
- Generated Lynx outputs were scanned for DOM/React leftovers (`<div`, `<span`, `<label`, `<button`, `<svg`, `onClick`, `onChange=`, `onFocus=`, `onBlur=`, `e.target`, `resolveDataSource`, `componentToReact`) with no matches in Lynx output.

**Remaining risk:** native-device validation is still required; Lynx Web remains smoke/reference only.

**Key insight:** Lynx web rendering (@lynx-js/web-core) masih immature — banyak fitur broken. Gunakan Lynx hanya untuk native mobile, bukan web.

**Key decisions:**
- Lynx target **bukan priority** — focus web (React, Vue, Stencil) dulu
- Kalau butuh mobile, React Native target sudah tersedia di Mitosis
- Lynx target sebagai **Phase 8** (future) — setelah Flutter generator

**Update 2026-06-14 (ADR-006):** Lynx track ditutup. Visi "satu library, jalan di
web + desktop + mobile" dipenuhi dengan **React Native + NativeWind** sebagai inti
Tier A (bersama React + Vue + Stencil), bukan via target tambahan Lynx. Alasan:

- React Native + NativeWind sudah memberi Tailwind-driven mobile UI tanpa generator custom.
- Maintenance dua mobile-stack (Lynx + RN) tidak proporsional untuk tim ukuran sekarang.
- Lynx web renderer masih tidak layak (lihat di atas) → tidak menambah nilai untuk Tauri/web.
- ADR-006 fokus ke "Single Caller API" cross-platform; Lynx menambah complexity tanpa
  benefit yang setimpal di luar use-case TikTok-style native shell.

Lihat **[ADR-006 §6 (Cleanup Tasks)](./006-universal-component-strategy.md#section-6--cleanup-tasks-phase-0)** untuk daftar artifact yang dihapus.

## Consequences

### Positive

- **Solve React integration problem** — native React output, bukan wrapper
- **Future-proof** — bisa output ke framework baru tanpa rewrite
- **Consistent naming** — satu source, nama sama di semua target
- **Theme ecosystem** — compatible dengan shadcn/ui, tweakcn.com
- **Standalone** — components bisa dipakai tanpa BitCode
- **Enterprise-grade** — events, methods, properties yang proper

### Negative

- **Migration effort besar** — 128 components perlu di-rewrite dari Stencil ke Mitosis
- **Go engine perlu adaptasi** — `component_compiler.go` perlu update (minimal untuk Stencil output)
- **Custom Flutter generator** — perlu 4-6 minggu development
- **Learning curve** — team perlu belajar Mitosis syntax dan patterns
- **BcSetup refactor** — dari singleton DOM-coupled ke pluggable context

### Risks

| Risk | Level | Mitigation |
|------|-------|-----------|
| Mitosis breaking changes | Medium | Pin version, community 13.8k stars |
| Flutter generator tidak perfect | High | SwiftUI as fallback, React Native as alternative |
| Tailwind di Flutter tidak ada | High | Material 3 / Cupertino theme adapter |
| BcSetup refactor break existing code | Medium | Backward compatible — singleton tetap work |
| 128 components rewrite lambat | Medium | Phase-based migration, mulai dari primitives |

## Migration Strategy

### Phase 0: Proof of Concept (1 minggu)
- Setup Mitosis project + Tailwind v4
- Buat 1 component (BcInput) dengan full features
- Output ke Stencil + React
- Verify backward compatibility (Stencil output = `<bc-input>`)
- Verify React output natural (bukan wrapper)

### Phase 1: Foundation (2-3 minggu)
- Setup monorepo structure (`packages/ui-core`, `packages/ui-react`, `packages/ui-vue`, dll)
- Implement BcSetup hybrid (singleton + context)
- Implement theme system (shadcn standard, dark/light/system)
- Build primitives: BcButton, BcBadge, BcLabel, BcSpinner, BcSkeleton, BcAvatar, BcSeparator, BcCard

### Phase 2: Forms (3-4 minggu)
- BcInput, BcTextarea, BcRichText, BcCode, BcNumber
- BcCheckbox, BcToggle, BcSelect, BcDateTime
- BcFile, BcLookup, BcBarcode, BcColor, BcGeo, BcRating, BcJson, BcSignature
- 4-layer data source implementation
- Field events + validation

### Phase 3: Layout + Overlays (3-4 minggu)
- BcRow, BcColumn, BcSection, BcTabs, BcHeader, BcSeparator, BcButtonGroup
- BcDialog, BcAlertDialog, BcPrompt, BcWizard, BcToast
- BcBreadcrumb, BcPagination, BcMenu, BcDropdown

### Phase 4: Data + Charts (6-8 minggu)
- BcDataTable (full enterprise features)
- BcFilter, BcLookupModal
- BcKanban, BcCalendar, BcGantt, BcMap, BcTree
- BcChart (26 chart types via type prop)
- BcExport, BcPrint, BcReportLink

### Phase 5: Media + Specialized (2-3 minggu)
- BcViewerPdf, BcViewerImage, BcViewerDocument, BcViewerVideo, BcViewerAudio
- BcSearch, BcFavorites
- BcActivity, BcChatter, BcTimeline
- All widgets (BcBadge, BcCopy, BcPhone, dll)

### Phase 6: Flutter Generator (4-6 minggu)
- Write custom Mitosis → Flutter/Dart generator
- Material 3 theme adapter
- Cupertino theme adapter
- Test dengan subset components

### Phase 7: Go Engine Adaptation (2-3 minggu)
- Update `component_compiler.go` untuk generate Stencil tags (backward compat)
- Minimal changes — Stencil output preserves `bc-*` tags

## References

- [Mitosis GitHub](https://github.com/BuilderIO/mitosis) — 13.8k stars, MIT license
- [Mitosis Documentation](https://mitosis.builder.io/docs/)
- [DB UX Design System](https://github.com/db-ux-design-system/core-web) — Production Mitosis example
- [shadcn/ui](https://ui.shadcn.com/) — Theme format standard
- [tweakcn.com](https://tweakcn.com/) — Custom theme generator
- [Tailwind CSS v4](https://tailwindcss.com/) — `@theme inline` syntax
- [NativeWind](https://www.nativewind.dev/) — Tailwind for React Native
- [bud.app](https://bud.app/) — Reference for Tailwind v4 + CSS Modules + custom font strategy
