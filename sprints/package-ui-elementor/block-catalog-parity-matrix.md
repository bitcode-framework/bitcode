# Block Catalog Candidate Parity Matrix

Status: Draft 0.1

This matrix maps Elementor, Website/CMS, EasyERP, portal, commerce, workflow, analytics, and builder references into a BitCode standalone Block Catalog plan. It is intentionally matrix-first: ADRs should use this as evidence, not duplicate every candidate row.

## Non-negotiable gates

- Ready Catalog entries are production-grade across web and Flutter from day one.
- No stub or placeholder implementations are allowed for Ready entries.
- Flutter is first-class, especially for ERP use cases.
- Existing ADR-003 component taxonomy remains unchanged. Block categories are metadata.
- Integration is based on general-purpose `DataSource` and future `ActionSource` adapters, not BitCode-specific bindings.
- Required Ready Catalog reactivity uses Portable Expression AST, not string expressions as the foundation.
- Duplicate analysis is based on contract equivalence, not label similarity.

## Categories

1. `content`
2. `data`
3. `interaction`
4. `commerce`
5. `workflow`
6. `builder`
7. `analytics`
8. `navigation`

`account` is a subcategory/tag, not a top-level category.

## Status legend

| Status | Meaning |
|---|---|
| `existing-component` | Existing web and Flutter components cover the concept; usually needs only catalog entry or schema preset. |
| `existing-component-needs-binding` | Existing components cover UI, but the block needs DataSource/ActionAdapter wiring or catalog presets. |
| `schema-preset-only` | No new component should be created; compose existing primitives/components via schema preset. |
| `component-gap` | New `Bc*` component required in `packages/ui` and `packages/ui-flutter`. |
| `contract-gap` | Component exists, but its contract is insufficient for the target semantics. |
| `adapter-gap` | UI exists or can be composed, but general-purpose adapter/capability is missing. |
| `duplicate-alias` | Reference name is alias of a canonical entry and must not become a separate Ready entry. |
| `excluded` | Raw reference intentionally excluded from v1 core Ready Catalog. |

## Candidate matrix

| Canonical ID | Category | Decision | Status | Render | Key aliases | Existing mapping | Main gaps | Batch |
|---|---|---:|---|---|---|---|---|---|
| `content.hero` | content | keep | schema-preset-only | schema | Hero Banner, Hero Section, Slides hero | BcSection, BcStack, BcImage, BcButton | BlockCatalog schema preset and Flutter demo example | content-website |
| `content.cta` | content | merge | existing-component-needs-binding | schema | CTA, Call to Action, Button Banner | BcCard, BcButton, BcImage | ActionSource trigger support for non-navigation actions | content-website |
| `content.rich-text` | content | merge | existing-component-needs-binding | schema | Text Editor, Rich Editor, HTML Block, Post Content | BcRichText, BcMarkdown | Sanitize policy in BlockCatalog security metadata | content-website |
| `content.plain-text` | content | keep | schema-preset-only | schema | Text, Heading, Page Title, Post Title | BcLabel, BcSection | Text node support or canonical text primitive in BlockCatalog schema | content-website |
| `content.media` | content | merge | existing-component-needs-binding | schema | Image, Video, Audio, Featured Image, Video Playlist | BcImage, BcViewerImage, BcViewerVideo, BcViewerAudio, BcViewerYoutube, BcGalleria, BcCarousel | Provider adapter metadata and media URL safety policy | content-website |
| `content.gallery` | content | merge | existing-component-needs-binding | schema | Gallery, Basic Gallery, Image Carousel, Media Carousel | BcGalleria, BcCarousel, BcViewerImage | Gallery preset variants and datasource shape | content-website |
| `content.faq` | content | keep | schema-preset-only | schema | FAQ, Toggle, Accordion FAQ | BcAccordion, BcDisclosure | FAQ schema preset, heading hierarchy policy, i18n keys | content-website |
| `content.testimonial` | content | merge | schema-preset-only | schema | Testimonial, Testimonial Carousel, Reviews | BcCard, BcCarousel, BcRate, BcAvatar | Testimonial schema shape and avatar fallback policy | content-website |
| `content.team-member` | content | merge | schema-preset-only | schema | Team Member, Author Box, Profile Card | BcCard, BcAvatar, BcDescriptions, BcFollow | BcAvatar Flutter parity check | content-website |
| `content.pricing-table` | content | keep | existing-component-needs-binding | schema | Pricing Table, Price Table, Price List | BcCard, BcDataTable, BcButton, BcStatistic | ActionSource for plan selection actions | content-website |
| `content.countdown` | content | keep | component-gap | component | Countdown, Timer | none | Need BcCountdown or portable time/reactive block | content-website |
| `content.social-links` | content | merge | existing-component-needs-binding | schema | Social Links, Social Icons, Share Buttons, Link in Bio | BcLike, BcFollow, BcButton (+ BcShare source export gap) | External URL safety policy, provider icon strategy, and BcShare export/registry fix | content-website |
| `data.record-list` | data | merge | existing-component-needs-binding | schema | Table, List View, Order List, Invoice List, Ticket List, Posts, Loop Grid | BcDataTable, BcList, BcVirtualScroller (+ BcGridView source export gap) | Datasource conventions, performance budget, mock demo data, and BcGridView export/registry fix | erp-data-interaction |
| `data.record-detail` | data | merge | existing-component-needs-binding | schema | Detail View, Product Detail, Blog Post, Job Detail | BcDetailView, BcDescriptions, BcForm | Detail schema preset and permission adapter integration | erp-data-interaction |
| `data.account-summary` | data | keep | schema-preset-only | schema | AccountWidget, My Account, Profile Summary, Billing Summary | BcCard, BcDescriptions, BcStatistic, BcAvatar | Account summary preset and avatar Flutter parity check | erp-data-interaction |
| `data.download-center` | data | keep | existing-component-needs-binding | schema | Download Center, File List, Document Center | BcFile, BcViewerDocument, BcViewerPdf, BcDataTable, BcReportLink | Safe download action contract and file validation policy | erp-data-interaction |
| `data.map-view` | data | merge | existing-component-needs-binding | schema | Map, Google Maps, Location Map | BcMap, BcViewMap, BcGeo | Generic map provider adapter, not Google-specific core block | content-website |
| `data.audit-log` | data | merge | existing-component-needs-binding | schema | Audit Log, Activity Feed, Activity Timeline | BcActivity, BcTimelineSocial, BcViewActivity, BcTimeline | Generic activity item contract and datasource examples | erp-workflow |
| `data.pivot-table` | data | keep | component-gap | component | Pivot, Pivot Table, Spreadsheet | BcDataTable | BcPivot or pivot mode contract plus aggregation adapter | analytics |
| `interaction.form` | interaction | merge | existing-component-needs-binding | schema | Form, Contact Form, Login, Apply Jobs, Dynamic Form | BcForm, BcRenderer, BcWizard | ActionSource/ActionRunner, Portable Expression AST, and Flutter renderer export validation | foundation |
| `interaction.multi-step-wizard` | interaction | keep | existing-component-needs-binding | schema | Wizard, Apply Jobs, Import Wizard | BcWizard, BcForm, BcSteps | Wizard schema convention and action lifecycle integration | erp-data-interaction |
| `interaction.record-picker` | interaction | merge | contract-gap | component | Relation Picker, OdooSelect, Search More, Lookup Modal | BcLookupModal, BcFieldTableSelect, BcFieldDynLink, BcFieldMorph, BcSelect, BcDataTable | Unified record picker block contract, quick-create action support, JSON snapshot convention, shared UI surface | erp-data-interaction |
| `interaction.import-wizard` | interaction | keep | component-gap | component | ImportExcelAction, Excel Import, CSV Import | BcWizard, BcFile, BcDataTable, BcDropZone | BcImportWizard/schema-driven import block, parser adapter, preview grid | erp-data-interaction |
| `interaction.notification-center` | interaction | keep | existing-component-needs-binding | component | Notifications, Notification Center, WooCommerce Notices | BcNotification, BcToastQueue, BcDropdown | Compose from existing primitives; missing realtime/action capability solved by ActionSource and DataSource realtime policy | erp-data-interaction |
| `workflow.approval-actions` | workflow | keep | component-gap | component | Approval Flow, Approval Actions, Change stage/status, Assignment | BcSteps, BcWizard, BcPrompt, BcActivity | Generic workflow action component and ActionSource lifecycle | erp-workflow |
| `workflow.state-indicator` | workflow | keep | contract-gap | schema | State Badge, Status Bar, Progress Tracker, Priority | BcBadge, BcSteps, BcProgress, BcTagGroup | Confirm BcBadge/BcProgress Flutter parity or add BcStateBadge/BcStatusBar | erp-workflow |
| `workflow.diagram-viewer` | workflow | keep | schema-preset-only | schema | Workflow Diagram, Process Diagram, Workflow Chart | BcChartGraph, BcSteps | Read-only process graph preset and example data shape | erp-workflow |
| `workflow.diagram-editor` | workflow | keep | component-gap | component | Workflow Diagram Editor, Process Designer, BcWorkflowBuilder | BcViewEditor, BcChartGraph, BcTree | Dedicated workflow authoring component (drag, edges, properties, persistence, undo/redo) | builder-surface |
| `workflow.chatter` | workflow | merge | existing-component-needs-binding | component | Chatter, Chat, Comments, Post Comments, Facebook Comments | BcChatter, BcComment, BcMentions, BcActivity | Generic comment/message datasource and action adapter examples | erp-workflow |
| `builder.query-builder` | builder | keep | component-gap | component | Query Builder, Filter Builder, Rule Builder | BcDataTable, BcForm, BcTree, BcCode | Dedicated builder surface, Portable Expression AST editor, validation, undo/versioning | builder-surface |
| `builder.expression-builder` | builder | keep | component-gap | component | Expression Builder, Formula Builder | BcCode, BcJson, BcTree, BcForm | Portable Expression AST types/evaluator | foundation |
| `builder.page-builder` | builder | keep | component-gap | component | Page Builder, Template, Dynamic Form, Dashboard Builder | BcViewEditor, BcRenderer, BcDashboard | Application Builder surface separate from Ready Block runtime and Flutter renderer export validation | builder-surface |
| `builder.report-builder` | builder | keep | component-gap | component | Report Builder, Dashboard Builder | BcViewReport, BcDashboard, BcChart, BcDataTable, BcExport | Report/dashboard builder components and persistence/versioning | builder-surface |
| `analytics.kpi-grid` | analytics | merge | schema-preset-only | schema | KPI Card, Stat Card, Counter, Statistic | BcKpi, BcStatistic, BcDashboard | Verify BcKpi Flutter parity or standardize on BcStatistic | analytics |
| `analytics.chart` | analytics | merge | existing-component-needs-binding | schema | Chart Line, Chart Bar, Chart Pie, Gauge, Heatmap, Funnel, Sankey, Treemap | BcChart, BcSparkline, BcProgress | Use ADR-015 shared ECharts option contract | analytics |
| `commerce.product-grid` | commerce | merge | existing-component-needs-binding | schema | Product List, Product Grid, Products, Product Categories | BcDataTable, BcCard, BcImage, BcRate (+ BcGridView source export gap) | Reuse data.record-list commerce schema variant, commerce item contract, add-to-cart ActionSource, and BcGridView export/registry fix | commerce |
| `commerce.product-detail` | commerce | merge | existing-component-needs-binding | schema | Product Detail, Product Title, Product Images, Product Price | BcDetailView, BcGalleria, BcRate, BcTabs, BcButton | Reuse data.record-detail commerce schema variant and ActionSource for cart actions | commerce |
| `commerce.cart-checkout` | commerce | merge | component-gap | component | Cart, Checkout, Menu Cart, Add To Cart, Purchase Summary | BcDataTable, BcForm, BcWizard, BcButton | Generic cart/checkout components and payment ActionAdapter | commerce |
| `commerce.payment-action` | commerce | merge | adapter-gap | schema | PayPal Button, Stripe Button, Payment Button | BcButton, BcDialog | Generic payment adapter, provider buttons as plugin candidates | commerce |
| `commerce.review` | commerce | merge | existing-component-needs-binding | schema | Review, Reviews, Product Rating | BcRate, BcComment, BcChatter | Generic review datasource/action shape | commerce |
| `navigation.menu` | navigation | merge | existing-component-needs-binding | schema | Nav Menu, Mega Menu, Menu Anchor, Sidebar, Off-Canvas | BcNavigationMenu, BcMegaMenu, BcMenubar, BcPanelMenu, BcSlideMenu, BcSidebar, BcBreadcrumb, BcAnchor | Responsive/adapt policies per menu variant | navigation-cms |
| `navigation.search-toc-sitemap` | navigation | merge | existing-component-needs-binding | schema | Search, Search Bar, TOC, Table Of Content, Sitemap, Category List | BcSearch, BcAnchor, BcTree, BcList, BcBreadcrumb | CMS navigation schemas and search datasource conventions | navigation-cms |
| `plugin.embed-social` | content | alias | duplicate-alias | component | Facebook Page, Facebook Embed, SoundCloud, Shortcode | BcViewerYoutube, BcViewerInstagram, BcViewerTiktok, BcComment, BcChatter | Provider-specific embeds are plugin/provider catalog, not core Ready blocks | plugin-extensibility |
| `excluded.wordpress-specific` | content | alias | excluded | schema | Template, WooCommerce Pages, WooCommerce Notices | generic equivalents | Excluded: third-party/WordPress-specific implementation detail | none |

## Required new foundational work

1. `foundation.block-catalog-types`: add BlockCatalogEntry, Ready/Candidate catalog types, category/channel/capability enums, readiness metadata, examples, adaptive policy, security/i18n/offline/versioning/test gates to `@bitcode/ui-schema`.
2. `foundation.action-source`: add ActionSource, ActionAdapter, ActionResult, ActionError, ActionEffect, ActionRunner lifecycle, and Web Component/Flutter event parity.
3. `foundation.portable-expression-ast`: add Portable Expression AST plus built-in function registry manifest. Ready Catalog cannot depend on string expressions for required behavior.
4. `foundation.block-renderer`: add `BcBlockRenderer` as a thin layer over `BcRenderer` with schema/component strategies and platform component registries.
5. `foundation.validation-tooling`: add build-time and runtime validators for Ready Catalog: no stubs, web + Flutter availability, schema/action/expression/security/i18n/adaptive/test gates.

## Key findings

- Most public/portal/content blocks are schema presets over existing components, not component gaps.
- The largest true gaps are not Elementor-style content. They are ActionSource, Portable Expression AST, builder surfaces, import wizard, checkout/payment, workflow diagram editor, pivot, timer/countdown.
- Several Elementor third-party widgets should be generic adapters or plugin catalog entries, not core blocks.
- Some existing mappings require Flutter parity verification: `BcAvatar`, `BcBadge`, `BcProgress`, `BcKpi` are used in candidate mappings but are not clearly present in the Flutter component file list from this session.
- Existing `DataSource` is general-purpose and should remain the read contract. A matching `ActionSource` is the main missing mutation contract.
- Reclassification rule from review session: blocks with `action-trigger` semantics must not be marked `schema-preset-only`; they belong in `existing-component-needs-binding` until ActionSource/ActionRunner are available.
- Workflow Diagram is split into `workflow.diagram-viewer` (static, BcChartGraph/BcSteps preset) and `workflow.diagram-editor` (component-gap, requires dedicated authoring component).
- Commerce blocks (`commerce.product-grid`, `commerce.product-detail`) are documented as Ready Catalog entries that reuse `data.record-list`/`data.record-detail` with commerce schema variants, rather than as separate primitives.
- Notification Center is treated as a composed block from existing primitives (`BcNotification` + `BcToastQueue` + `BcDropdown`); the missing capability is realtime/action, which is solved by ActionSource plus DataSource realtime policy.
