# 2026-06-30 — Block Catalog Grilling + Matrix + ADR Skeletons

## Apa yang sudah dikerjakan

Sesi grilling-with-docs untuk merancang **Block Catalog Architecture** (website + ERP blocks di atas `packages/ui` dan `packages/ui-flutter`). ~58 pertanyaan grill dijawab dengan keputusan tegas. Matrix formal dibuat, ADR skeleton dibuat, CONTEXT.md + components-progress.md + ADR README diupdate. Sesi lanjutan 2026-06-30 melakukan review evidence-first, membuat `REVIEW-FINDINGS.md`, dan mendraft full body ADR-017/018.

## Keputusan inti (non-negotiable)

1. **2 istilah resmi**: Website Block + ERP Interaction Component (bukan satu istilah "widget"). Sudah di CONTEXT.md.
2. **Block Catalog** = metadata layer di atas ADR-003 taxonomy. Source taxonomy TIDAK diubah. Kategori = capability (`content`, `data`, `interaction`, `commerce`, `workflow`, `builder`, `analytics`, `navigation`); `account` = subcategory. Channel tags (`website`, `portal`, `internal`, `saas`) bukan kategori.
3. **Sub-entrypoint `@bitcode/ui/website`, `/erp`, dll** = curated barrels only, bukan taxonomy struktural baru. Single Caller API ADR-006 tetap.
4. **Default artifact = schema preset**; naik jadi `Bc* component` jika ada behavior kompleks / data lifecycle / mutation / security / native / schema duplikasi berat.
5. **Component exports pakai `Bc*`; block preset IDs pakai namespaced slug** (`content.hero`, `commerce.product-grid`).
6. **Block Catalog shape/types di `@bitcode/ui-schema`** (ADR-012 canonical). Preset bawaan di `packages/ui` per-category files + generated index. Engine/module JSON hanya reference/override.
7. **Built-in preset TIDAK hardcode model/API**; abstract DataSource contract. Engine = salah satu provider adapter saja.
8. **Flutter-first, no stub**: setiap Ready Block wajib full render web + Flutter sejak awal. Entry yang belum Flutter-ready hanya boleh Candidate, bukan Ready.
9. **Candidate Catalog = staging/backlog**; target semua dipromosikan Ready; runtime hanya baca Ready Catalog.
10. **Prioritas eksekusi ERP interaction dulu**, tapi scope akhir tetap seluruh catalog v1.
11. **Curated parity** (bukan blind 1:1 Elementor/EasyERP); dedupe by contract equivalence; ambiguous wajib `keep/merge/split/alias` decision + reason.
12. **ADR = aturan dedupe**; parity matrix = keputusan merge/split per item; contract file = bukti teknis.
13. **Ready Block wajib demo cross-platform** di `packages/ui/demo-app` + `packages/ui-flutter/example` (section Block Catalog).
14. **Data-bound Ready Block**: mock demo wajib; engine integration demo wajib jika engine-bound.
15. **Generic reusable block dulu**; module-specific binding preset TIDAK masuk core (prinsip standalone). BitCode engine = optional adapter via DataSource/ActionAdapter.
16. **DataSource general-purpose 4-layer** (`local/fetcher/url/model` via `ModelAdapter`) — BUKAN BitCode-specific.
17. **ActionSource perlu** (mirror DataSource: `local/handler/url/model`); Ready JSON hanya serializable (`local/url/model`), `handler` hanya runtime injection.
18. **ActionResult normalized** `{ok,status,data,message,errors,errorList,effects}` + ActionError structured.
19. **ActionLifecycle events**: `bc-action-start/success/error/validation-error/settled` + Dart notifiers.
20. **ActionEffect typed & safe** (`setValue,setMeta,setError,refreshDataSource,invalidateKey,toast,navigate,closeOverlay,resetForm,emit`); no arbitrary string command.
21. **ADR Block Catalog wajib bagian "Data, Action, and Reactivity Contract"** menyatukan DataSource + ActionSource + ADR-004 reactions.
22. **Existing component adjustment terarah** (bukan rewrite): `ui-schema`, `BcRenderer`, `BcForm`, action-runner, komponen interaction berat. Primitives/display pasif tidak diubah.
23. **Ready Catalog JSON hanya serializable reactions** (Portable Expression AST); function reactions boleh untuk consumer runtime code.
24. **String expr = advanced/server-evaluated escape hatch**; fondasi Ready = Portable Expression AST (karena expr-js/go-json belum ada parity Dart/Swift/KMP).
25. **Portable Expression AST scope kecil & composable**: `literal/path/eq/ne/gt/gte/lt/lte/and/or/not/add/sub/mul/div/sum/count/min/max/avg/concat/contains/startsWith/endsWith/isEmpty/isNotEmpty/coalesce/if`. Function call boleh via whitelisted built-in registry (`math.*,string.*,array.*,date.*,format.*,validate.*`), pure & deterministic, versioned.
26. **Host/project function register manual/auto via capability manifest**; Ready universal hanya boleh depends built-in portable; host function = explicit manifest + `requiredFunctions` + fallback. Runtime fail-loud jika function belum registered.
27. **RuntimeContext** = explicit named registry (`values/services/functions/adapters/permissions/i18n`), BUKAN positional `[x,z]`. Extends AppProvider/BcAppConfig (ADR-013), bukan global baru.
28. **BcBlockRenderer** = thin layer di atas `BcRenderer`; dua strategy `schema` | `component` (eksplisit `renderStrategy` field).
29. **Component strategy**: catalog JSON pakai string component key; runtime pakai typed registry per platform (`registerBlockComponent("BcCheckout", BcCheckout)`).
30. **Ready Block wajib deklarasi** `requiredComponents`, `requiredFunctions`, `requiredAdapters`.
31. **Permissions/security** = typed adapter capability (`permission`/`authz`), bukan hardcoded RBAC/BitCode.
32. **Validation build-time + runtime** (build-time cegah catalog rusak; runtime cegah host salah config).
33. **Ready Catalog source** per-category files + generated index (bukan satu file besar).
34. **Candidate Catalog/matrix** = machine-readable JSON source + generated markdown report.
35. **8 top-level categories final**; `account` subcategory/tag.
36. **Full scope v1** = semua candidate dari referensi awal setelah dedupe/curate; bukan infinite universe.
37. **Raw references lengkap** (semua Elementor/EasyERP/ChatGPT); curated candidates hanya yang relevan; excluded wajib reason.
38. **Third-party-specific** (Facebook, PayPal, Stripe, Google Maps) → adapter/plugin catalog, BUKAN core Ready, kecuali generic equivalent.
39. **Plugin/adapter catalog** masuk ADR sebagai extensibility model; implementasi v1 fokus generic core blocks.
40. **Ready Block wajib adaptive policy** (minimal `inherits child components`; eksplisit untuk block kompleks: breakpoints, iOS Cupertino vs Android Material, overlay morph, touch target, reduced motion, offline/loading/error).
41. **Performance budget** per block type (list/table/data = pagination/virtualization/debounce/abort; Flutter no unbounded scroll; content = no layout shift). Aturan teknis + test ringan, bukan benchmark berat untuk semua.
42. **Security/privacy gate** wajib untuk data/action/commerce/workflow/builder: no secrets in schema, adapters own auth, permission adapter, safe URL allowlist, no arbitrary code execution, sanitize rich text, file validation, payment via adapter. Sanitize default ON; opt-out explicit + policy-gated + auditable (mirip Laravel).
43. **i18n**: Ready Block strings wajib key + English fallback (`labelKey`/`fallbackLabel`).
44. **A11y**: Ready Block wajib block-level a11y contract (landmark, heading hierarchy, labels/errors, keyboard order, focus mgmt, ARIA live, touch target, reduced motion).
45. **Customization**: typed simple props (`variant/density/layout/themeRole/mediaPosition`) + advanced `slotProps`/parts.
46. **Offline/cache**: declarative policy `offlinePolicy`/`cachePolicy`/`conflictPolicy`; tidak semua block wajib full offline.
47. **Versioning/migration**: BlockCatalogEntry wajib `version` + `schemaVersion` + `migration` metadata.
48. **Testing gate**: per-block tests (schema/web/Flutter/action mock/a11y) + category smoke.
49. **Matrix-first sebelum ADR**; ADR setelah matrix stabil.
50. **2 ADR utama + 1 parity matrix/plan appendix**: ADR-017 Block Catalog Architecture, ADR-018 Data/Action/Reactivity Contract, parity matrix appendix.
51. **ADR = arsitektur + full target/gate; plan appendix = batch implementasi sampai semua Ready.**

## Approval terakhir (5 pertanyaan)

1. Action-trigger block (CTA/pricing/product-detail) → `existing-component-needs-binding` sampai ActionSource ada. ✅ applied
2. Workflow diagram → split viewer (schema-preset-only) + editor (component-gap). ✅ applied
3. Notification Center → composed block dari existing primitives (`BcNotification`+`BcToastQueue`+`BcDropdown`); missing capability = realtime/action via ActionSource + DataSource realtime policy. ✅ applied
4. Commerce product-grid/product-detail → preset generic `data.record-list/detail` dengan commerce schema variant (bukan primitive terpisah). ✅ applied
5. Provider embeds/payments/maps → plugin/provider catalog, bukan core Ready v1 (maps sudah ada `BcMap`/`BcViewMap`). ✅ applied

## File yang diubah/dibuat

- `CONTEXT.md` — tambah 4 glossary: Website Block, ERP Interaction Component, Block Catalog, Ready Block.
- `sprints/package-ui-elementor/block-catalog-candidate-matrix.json` — machine-readable matrix (44 candidates, 12 raw groups, 5 foundational work). Version 0.1 accepted baseline for ADR-017/018.
- `sprints/package-ui-elementor/block-catalog-parity-matrix.md` — generated/readable report.
- `docs/components-progress.md` — tambah section "Block Catalog (Baseline 0.1)" dengan counts + foundational work + categories + key findings.
- `docs/adr/017-block-catalog-architecture.md` — Accepted (catalog surface, categories/channels/capabilities, render strategies, Ready gate, RuntimeContext, plugin model, Flutter parity, validation, consequences, phases, AC).
- `docs/adr/018-block-catalog-data-action-reactivity.md` — Accepted (ActionSource, ActionResult/ActionError, lifecycle, effects, Portable Expression AST, function registry, RuntimeContext, runAction, security, phases, AC).
- `docs/adr/README.md` — row ADR-017/018 updated to full body status and next ADR number corrected to ADR-019.
- `sprints/package-ui-elementor/REVIEW-FINDINGS.md` — review findings for count drift, Flutter parity gaps, web export gaps, BcMarkdown gap, and ADR-014 overlap.

## Status counts final (matrix 0.1, 44 candidates)

- schema-preset-only: 8
- existing-component-needs-binding: 21
- component-gap: 10
- contract-gap: 2
- adapter-gap: 1
- duplicate-alias: 1
- excluded: 1

## Konteks penting / gotcha

- ADR-017 dan ADR-018 sudah dipromote ke **Accepted** pada 2026-06-30 setelah final read-through; tidak ada architectural blocker tersisa sebelum implementation planning.
- Matrix 0.1 sudah reclassify dan dilabel accepted baseline untuk ADR-017/018. Later changes require version bump atau documented change record.
- `workflow.diagram-editor` = component-gap (dedicated authoring component, drag/edges/properties/persistence/undo-redo). `workflow.diagram-viewer` = schema-preset-only (BcChartGraph + BcSteps preset read-only).
- Flutter parity decision: add canonical Flutter `BcAvatar`, `BcBadge`, `BcProgress`, and `BcKpi` with source, tests, and barrel exports before Ready; do not remap to generic primitives just to pass parity.
- Export gap decision: `BcShare` and `BcGridView` source exists, so treat missing root export/registry as export bug, not component-gap. Dependents stay non-Ready until public exports are fixed.
- `BcMarkdown` decision: implement canonical web `BcMarkdown` with tests and public export; do not permanently remap markdown-backed rich text to `BcRichText` only.
- `BcRenderer` di Flutter ada sebagai renderer infrastructure (`packages/ui-flutter/lib/src/renderer/bc_renderer.dart`), bukan component file. Ready validation harus cek renderer export terpisah.
- ADR-017/018 depend on ADR-013 (cross-tier infra) + ADR-014 (engine UI wiring, masih Proposed). Boundary resolved: Block Catalog owns catalog/preset/generic contracts; ADR-014 owns endpoint, model, permission, realtime, persistence, and engine wiring for overlap features.
- Minor review edits applied: inventory counts clarified as background only, offline action Ready proof tightened, permission/context validation moved before sensitive expression evaluation, `handler` portability manifest rules tightened, and a11y evidence location required.
- Provider adapter decision: core ships generic adapter contracts plus mock/local adapters for tests only; provider SDK integrations stay in plugin packages.
- Matrix JSON validated: parse OK, no missing mapped components (setelah BcShare/BcGridView dihapus dari webComponents mapping dan dipindah ke gaps).

## Next steps (konkret)

1. Create implementation plan appendix (batch sampai semua Ready) in `docs/plans/` before code changes.
2. Implement resolved gap fixes in a separate code phase: canonical Flutter widgets (`BcAvatar`, `BcBadge`, `BcProgress`, `BcKpi`), canonical web `BcMarkdown`, and public export/registry fixes for `BcShare`/`BcGridView`.
3. Keep overlap rows non-Ready until ADR-014 surface exists for endpoint/model/permission/realtime/persistence/engine wiring.
4. Keep provider-specific SDK integrations out of core; implement only generic adapter contracts and mock/local test adapters in core.
5. If ADRs are approved, implement foundational work in order: block-catalog types, ActionSource/ActionRunner, Portable Expression AST, BcBlockRenderer, validation tooling.
6. Create implementation plan appendix (batch sampai semua Ready) in `docs/plans/` after ADR approval.
7. Commit current decision-resolution docs/matrix changes with a scoped message after verification.

## Prinsip yang harus diingat (dari user)

> Pilih yang paling robust, reliable, future-proof, easy maintain, production grade, efficient, flexible (mudah custom/extend), simple, powerful (handle berbagai edge case). Cukup 1 code menjadi platform-aware, mendukung web/desktop/mobile. Mobile platform-aware iOS→Cupertino, Android→Material, sesuai ADR. Bukan MVP, full feature production grade, no stub (kecuali phase berikutnya implement stub). Kerjakan dengan ketelitian tinggi, disiplin, cermat, matang, detail, baik benar bagus sempurna lengkap dan terbaik.
