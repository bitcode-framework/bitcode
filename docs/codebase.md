# BitCode Platform — Codebase Map

## Repository Structure

```
bitcode/
├── docs/                                       # Project-level documentation
│   ├── architecture.md                         # System design, data flow, concepts
│   ├── codebase.md                             # This file — full file map
│   └── plans/                                  # Design documents & implementation plans
│       ├── DESIGN.md
│       ├── 2026-04-18-component-system-design.md
│       ├── 2026-04-18-component-system-plan.md
│       ├── 2026-04-22-fix-process-engine-design.md
│       ├── 2026-04-22-i18n-implementation-plan.md
│       ├── 2026-04-22-i18n-stencil-components-design.md
        │       ├── 2026-04-25-model-options-design.md
        │       └── 2026-04-26-media-viewers-design.md
│
├── engine/                                     # Go runtime (the core)
├── packages/                                   # Shared libraries
│   ├── components/                             # Stencil Web Components
│   ├── ui/                                     # Mitosis cross-framework UI library (React/Vue/Stencil; React Native via ADR-006 Phase 0)
│   ├── ui-schema/                              # @bitcode/ui-schema — standalone UI schema types + ReactiveStore (ADR-012, zero-dep, pure TS)
│   ├── go-json/                                # go-json — JSON/JSONC programming language engine (Go)
│   ├── go-json-runtimes/                       # Script runtime engines for go-json (separate go.mod)
│   └── tauri/                                  # Tauri native shell (desktop + mobile)
├── samples/                                    # Example applications
│   └── erp/                                    # Full ERP sample
├── sprints/                                    # Sprint tracking
└── .gitignore
```

---

## Engine (`engine/`)

The Go backend that reads JSON definitions and runs the application.

```
engine/
├── cmd/                                        # Entry points
│   └── bitcode/
│       ├── main.go                             # CLI — serve, dev, init, validate, module, user, db, seed, version, publish, security
│       ├── publish.go                          # Module publish command (extract embedded modules)
│       ├── publish_crud.go                     # publish:crud command — generate API/page override files from auto-generated CRUD
│       ├── security.go                         # security CLI — load/export/diff/validate/history (JSON↔DB sync)
│       ├── backup.go                           # db backup/restore commands (SQLite/Postgres/MySQL)
│       └── seed.go                             # Data migration CLI — seed run/rollback/status/fresh/create
│
├── internal/                                   # Private application code
│   ├── app.go                                  # Central wiring — NewApp(), LoadModules(), Start(), Shutdown()
│   │                                           #   Registers step handlers, middleware, routes
│   │                                           #   Module install: parse → register → migrate → seed
│   ├── config.go                               # Viper-based config — env vars + TOML/YAML file
│   │
│   ├── compiler/parser/                        # JSON → Go struct parsers
│   │   ├── model.go                            # ModelDefinition, FieldDefinition, field types, validation rules, APIConfig, ProtocolConfig, mask/mask_length/groups on fields
│   │   ├── model_test.go                       # 13 tests (valid, inheritance, missing fields, relationships, API config, field mask/groups)
│   │   ├── api.go                              # APIDefinition, ExpandAutoCRUD(), GetBasePath()
│   │   ├── api_test.go                         # 8 tests (auto_crud, workflow, custom, RLS)
│   │   ├── security.go                         # SecurityDefinition, SecurityACL (with "all" shorthand), SecurityRuleDefinition, ParseSecurity/ParseSecurityFile
│   │   ├── security_test.go                    # 4 tests (basic group, all shorthand, validation errors, rule defaults)
│   │   ├── process.go                          # ProcessDefinition, StepDefinition, 14 step type constants
│   │   ├── view.go                             # ViewDefinition, 6 view types (list/form/kanban/calendar/chart/custom)
│   │   ├── view_test.go                        # 6 tests
│   │   ├── agent.go                            # AgentDefinition, triggers, cron expressions, retry config
│   │   ├── migration.go                        # MigrationDefinition, source types (JSON/CSV/XLSX/XML), processors, conflict modes
│   │   ├── module.go                           # ModuleDefinition, permissions, groups, securities, pages, menu (with groups), settings, migrations, i18n patterns
│   │   ├── workflow.go                         # WorkflowDefinition, states, transitions, CanTransition()
│   │   └── workflow_test.go                    # 3 tests (parse, transitions, multi-from)
│   │
│   │   ├── domain/                                 # Business logic (no DB imports)
│   │   ├── model/
│   │   │   ├── registry.go                     # Register/Get/List/Has models, TableName()
│   │   │   └── registry_test.go                # 6 tests
│   │   ├── security/
│   │   │   ├── user.go                         # User aggregate — NewUser, CheckPassword, IsSuperuser, AllGroupNames
│   │   │   ├── role.go                         # Role aggregate — HasPermission (with inheritance) [DEPRECATED — being replaced by Group]
│   │   │   ├── group.go                        # Group aggregate — AllGroupNames (with implied groups), share, comment, module, modified_source
│   │   │   ├── permission.go                   # Permission value object [DEPRECATED — replaced by ModelAccess]
│   │   │   ├── model_access.go                 # ModelAccess entity — 12 ERPNext-style permissions per model per group (select/read/write/create/delete/print/email/report/export/import/mask/clone)
│   │   │   ├── model_access_test.go            # 4 tests (HasPermission, AllPermissions, SetFromList, IsGlobal)
│   │   │   ├── security_history.go             # SecurityHistory entity — audit trail for group/ACL/rule changes with snapshot + rollback
│   │   │   ├── security_history_test.go        # 2 tests (create, update)
│   │   │   ├── record_rule.go                  # RecordRule — m2m Groups, AppliesToGroupNames, IsGlobal, module, modified_source
│   │   │   └── security_test.go                # 14 tests (user, role inheritance, groups, record rules, superuser, AllGroupNames, share, m2m)
│   │   ├── event/
│   │   │   ├── bus.go                          # In-process event bus — Subscribe, SubscribeAll, Publish
│   │   │   └── bus_test.go                     # 4 tests
│   │   ├── setting/
│   │   │   ├── setting.go                      # Key-value store — Get, Set, GetWithDefault, LoadDefaults
│   │   │   └── setting_test.go                 # 5 tests
│   │   └── storage/
│   │       ├── storage.go                      # StorageDriver interface, PutOptions, URLOptions, ScanHook interface
│   │       └── attachment.go                   # Attachment entity — GORM model for attachments table
│   │
│   ├── runtime/                                # Execution engines
│   │   ├── executor/
│   │   │   ├── executor.go                     # ProcessExecutor — RegisterHandler, Execute, step dispatch
│   │   │   ├── executor_test.go                # 3 tests
│   │   │   ├── process_loader.go               # ProcessRegistry — Register, Get, List processes by name
│   │   │   ├── dag.go                          # DAG executor for parallel step execution
│   │   │   ├── dag_test.go                     # DAG tests
│   │   │   └── steps/                          # Step handler implementations
│   │   │       ├── validate.go                 # Validate step — eq, neq, required rules
│   │   │       ├── data.go                     # Query/Create/Update/Delete steps via GenericRepository
│   │   │       ├── control.go                  # If/Switch/Loop steps — condition evaluation, variable resolution
│   │   │       ├── emit.go                     # Emit step — add domain event to context
│   │   │       ├── call.go                     # Call step — invoke sub-process by name
│   │   │       ├── script.go                   # Script step — invoke TS/Python plugin via ScriptRunner
│   │   │       ├── http.go                     # HTTP step — external API calls with method/url/headers/body
│   │   │       ├── util.go                     # Assign + Log step handlers
│   │   │       └── steps_test.go               # 9 tests (validate, emit, assign, if, parse)
│   │   ├── expression/
│   │   │   ├── evaluator.go                    # Expression evaluator — lexer, parser, AST, arithmetic/comparison/boolean/functions
│   │   │   ├── evaluator_test.go               # 17 tests (arithmetic, fields, aggregates, comparisons, functions)
│   │   │   └── hydrator.go                     # Computed field hydrator — loads one2many children, evaluates computed/formula fields
│   │   ├── validation/
│   │   │   ├── validator.go                    # Field validation engine — ValidateCreate, ValidateUpdate, auto-map, short-circuit
│   │   │   ├── rules.go                        # Built-in validation rules — email, phone, regex, date, etc.
│   │   │   ├── conditional.go                  # Conditional validators — required_if, when, expression evaluator
│   │   │   ├── sanitizer.go                    # Sanitization engine — trim, lowercase, slugify, strip_tags, etc.
│   │   │   ├── errors.go                       # ValidationErrors — per-field error accumulation
│   │   │   ├── adapter.go                      # ValidatorAdapter — implements persistence.FieldValidator interface
│   │   │   ├── validator_test.go               # 18 tests
│   │   │   └── sanitizer_test.go               # 10 tests
│   │   ├── hook/
│   │   │   ├── dispatcher.go                   # Event dispatcher — priority sort, condition eval, sync/async, retry, timeout
│   │   │   ├── context.go                      # EventContext — model, data, old_data, changes, session, bulk info
│   │   │   ├── model_hooks.go                  # ModelHookDispatcher — bridges repository interface with dispatcher
│   │   │   ├── expr.go                         # Expression evaluator for handler conditions
│   │   │   └── dispatcher_test.go              # 10 tests
│   │   ├── agent/
│   │   │   ├── worker.go                       # Agent worker — subscribe to events, execute with retry
│   │   │   └── cron.go                         # Cron scheduler — periodic job execution
│   │   ├── workflow/
│   │   │   └── engine.go                       # Workflow engine — Register, ExecuteTransition, GetInitialState
│   │   ├── model_process.go                    # ModelProcessRegistry — built-in model operations (CRUD, aggregates, dynamic finders)
│   │   ├── dynamic_finder.go                   # Dynamic finder parser — FindBy{Field}, CountBy, SumBy, etc. with And/Or/operator suffixes
│   │   ├── dynamic_finder_test.go              # 24 tests for dynamic finder parsing
│   │   └── plugin/
│   │       └── manager.go                      # Plugin manager — spawn TS/Python processes, JSON-RPC over stdin/stdout
│   │   ├── comment/
│   │   │   ├── service.go                      # Comment service — Create (mention parse + resolve + notify), ListByRecord (paginated), Delete (soft), EnsureTables, ResolveMentions, recordFromMap
│   │   │   ├── mention.go                      # Mention parser — @username extraction (email-aware, dedup, boundary-safe), ResolveMentions (batch user lookup)
│   │   │   └── service_test.go                 # 13 tests — parse mentions (edge cases), resolve, create+notify, list by record
│   │
│   ├── infrastructure/                         # External concerns
│   │   ├── persistence/
│   │   │   ├── database.go                     # NewDatabase() — SQLite/Postgres/MySQL connection via GORM
│   │   │   ├── repository_interface.go         # Repository interface, SystemRepository, SequenceEngine, MigrationEngine interfaces
│   │   │   ├── query.go                        # Comprehensive Query builder — OR/AND/NOT groups, JOINs, HAVING, DISTINCT,
│   │   │   │                                   #   aggregates (COUNT/SUM/AVG/MIN/MAX), subqueries, UNION, raw expressions,
│   │   │   │                                   #   scopes, eager loading (WITH/preload), locking, soft delete scopes,
│   │   │   │                                   #   field sanitization, JSON/Map/Domain parsers
│   │   │   ├── query_test.go                   # Tests for query builder (all features)
│   │   │   ├── oql.go                          # OQL (Object Query Language) parser — 3 syntax styles:
│   │   │   │                                   #   Style A: SQL-like (JPQL/HQL), Style B: Simplified DSL,
│   │   │   │                                   #   Style C: Dot-notation. Auto-detect via ParseOQL()
│   │   │   ├── oql_test.go                     # Tests for OQL parser (all 3 styles)
│   │   │   ├── mongo_connection.go             # OpenMongoDB() — MongoDB connection via official driver
│   │   │   ├── mongo_repository.go             # MongoRepository — MongoDB implementation of Repository interface
│   │   │   ├── mongo_migration.go              # MongoMigrationEngine — index creation, system collection setup
│   │   │   ├── mongo_sequence.go               # MongoSequenceEngine — counter collection pattern for sequences
│   │   │   ├── mongo_system.go                 # MongoSystemRepository + MongoAuditLogRepository
│   │   │   ├── dynamic_model.go                # MigrateModel() — CREATE TABLE from ModelDefinition, dialect-aware DDL
│   │   │   │                                   #   MergeInheritedFields() — model inheritance field merging
│   │   │   │                                   #   Auto-creates junction tables for many2many
│   │   │   ├── repository.go                   # GenericRepository (SQL) — implements Repository interface with GORM
│   │   │   │                                   #   Full query translation: JOINs, OR/AND/NOT, HAVING, DISTINCT,
│   │   │   │                                   #   subqueries, locking, soft delete scopes, eager loading
│   │   │   │                                   #   Avg/Min/Max/Pluck/Exists/Aggregate/Chunk/Increment/Decrement
│   │   │   │                                   #   Transaction support, RawQuery/RawExec, relation preloading
│   │   │   │                                   #   Computed field hydration via expression.Hydrator
│   │   │   │                                   #   Data revision snapshots on write operations
│   │   │   ├── data_revision.go                # DataRevision — full record snapshots for rollback/restore
│   │   │   │                                   #   Monotonic versioning per (model, record_id), change diff
│   │   │   ├── data_revision_test.go           # 7 tests (create, version increment, list, get, cleanup, changes, latest)
│   │   │   ├── view_revision.go                # ViewRevision — view JSON snapshots for editor versioning
│   │   │   ├── view_revision_test.go           # 6 tests
│   │   │   ├── audit_log.go                    # AuditLogRepository — persistent audit log writer with async support
│   │   │   │                                   #   FindByRecord, FindByUser, FindLoginHistory, FindRequests
│   │   │   │                                   #   ImpersonatedBy field, AutoMigrateAuditLog()
│   │   │   ├── audit_log_test.go               # 5 tests (write, find by record, requests, user, login history)
│   │   │   ├── permission_checker.go            # PermissionService — resolves 12 permissions per model per user via ModelAccess + Group chain (additive, default-deny, superuser bypass)
│   │   │   ├── permission_checker_test.go      # 14 tests (superuser, default-deny, single group, additive, global ACL, implied groups)
│   │   │   ├── record_rule_service.go          # RecordRuleService — Odoo-compatible rule composition: global INTERSECT, group UNION. Domain interpolation.
│   │   │   ├── record_rule_service_test.go     # 11 tests (superuser, no rules, global, group in/out, operation filter, implied, legacy, inactive, interpolation)
│   │   │   ├── migration_tracker.go            # MigrationTracker — ir_migration table, batch tracking, status
│   │   │   └── backup.go                       # Backup/Restore — driver-aware (SQLite copy, pg_dump, mysqldump)
│   │   ├── cache/
│   │   │   ├── cache.go                        # Cache interface + NewCache() factory
│   │   │   ├── memory.go                       # MemoryCache — in-process with TTL
│   │   │   ├── memory_test.go                  # 5 tests (set/get, TTL expiry, delete, clear)
│   │   │   └── redis.go                        # RedisCache — Redis-backed implementation
│   │   ├── module/
│   │   │   ├── registry.go                     # Module registry — Register/Get/IsInstalled/InstalledNames
│   │   │   ├── dependency.go                   # ResolveDependencies() — topological sort, circular detection
│   │   │   ├── loader.go                       # LoadModule() — parse module dir, collect models + APIs + securities
│   │   │   ├── auto_api.go                     # GenerateAPIFromModel() — auto-creates APIDefinition from model "api" config. MergeAPIs() — merge auto-generated + override APIs. pluralize()
│   │   │   ├── auto_api_test.go                # 6 tests (basic, no-api, override endpoints, custom API, workflow override, pluralize)
│   │   │   ├── security_loader.go              # SecurityLoader — loads securities/*.json, syncs groups/ACL/rules/menus/pages to DB. Respects modified_source="ui" (noupdate). Idempotent.
│   │   │   ├── security_loader_test.go         # 6 tests (basic sync, implies, record rules, UI-modified preservation, all 12 permissions, idempotent)
│   │   │   ├── fs.go                           # DiskFS, EmbedFS, LayeredFS — module filesystem abstraction
│   │   │   ├── fs_test.go                      # FS tests
│   │   │   ├── reader.go                       # Multi-format data readers (JSON, CSV, XLSX, XML)
│   │   │   ├── migration.go                    # MigrationEngine, DataInserter (GORM+Mongo), RunUp/RunDown, processors, depends_on topo sort
│   │   │   ├── migration_test.go               # 26 tests (readers, upsert, field mapping, defaults, tracker, topo sort, circular detection)
│   │   │   ├── module_test.go                  # 7 tests (registry, dependencies, parse)
│   │   │   └── integration_test.go             # Integration tests
│   │   ├── i18n/
│   │   │   ├── loader.go                       # Translator — LoadFile/LoadJSON, Translate with locale fallback
│   │   │   └── i18n_test.go                    # 4 tests
│   │   ├── storage/
│   │   │   ├── config.go                       # StorageConfig, LocalStorageConfig, S3StorageConfig, ThumbnailConfig
│   │   │   ├── local.go                        # LocalStorage — filesystem StorageDriver implementation
│   │   │   ├── s3.go                           # S3Storage — AWS S3 StorageDriver implementation
│   │   │   ├── formatter.go                    # FormatPath/FormatName — template variable resolution
│   │   │   │                                   #   NewStorageDriver() — factory for local/S3
│   │   │   ├── repository.go                   # AttachmentRepository — GORM CRUD for attachments table
│   │   │   │                                   #   AutoMigrateAttachments(), FindByHash, FindVersions, CleanupVersions
│   │   │   └── thumbnail.go                    # ThumbnailService — generate thumbnails, on-demand resize
│   │   └── watcher/
│   │       └── watcher.go                      # FileWatcher — poll for .json/.html changes, trigger reload
│   │
│   └── presentation/                           # HTTP layer
│       ├── api/
│       │   ├── router.go                       # Dynamic route registration — auto-CRUD from model + override merge, permission + record rule middleware wiring
│       │   ├── crud_handler.go                 # Auto-CRUD handler — List/Read/Create/Update/Delete with pagination, field masking, field groups filtering, permission injection in response
│       │   ├── field_filter.go                 # Server-side field filtering — field groups (hide), field masking (****1234), per-response
│       │   ├── field_filter_test.go            # 11 tests (masking, field groups, nil model, unknown fields, list filtering)
│       │   ├── swagger.go                      # SwaggerGenerator — auto-generates OpenAPI 3.0 spec from model + API definitions. Swagger UI at /api/v1/docs
│       │   ├── auth_handler.go                 # POST /auth/login, /register, /logout, /2fa/enable, /2fa/disable, /2fa/validate
│       │   ├── upload_handler.go               # Legacy upload handler (replaced by file_handler)
│       │   ├── file_handler.go                 # FileHandler — upload, download, list, delete, versions, resize, thumbnail
│       │   ├── comment_handler.go              # CommentHandler — polymorphic comments per record, threaded replies, @mention parsing, permission + record rule enforcement
│       │   ├── comment_handler_test.go         # 5 tests (create+list, threading, missing record_id, empty body, permission denied)
│       │   ├── activity_handler.go             # ActivityHandler — GET /api/v1/activity aggregated feed, permission + record rule enforcement
│       │   ├── activity_handler_test.go        # 10 tests (happy path, filters, pagination, permission denial, auth)
│       │   ├── export_handler.go               # ExportHandler — POST/GET /api/v1/{module}/{model}/export (CSV/XLSX/JSON streaming), export permission + field-level mask + column sanitization
│       │   └── export_handler_test.go          # 12 tests (csv, xlsx, json, streaming, permission denial, filter, columns, field mask, column sanitization, superuser bypass)
│       ├── middleware/
│       │   ├── auth.go                         # JWT validation, user context injection, impersonated_by extraction
│       │   ├── permission.go                   # Permission checking via PermissionChecker interface (ModelAccess-based)
│       │   ├── record_rule.go                  # RLS filter injection with {{user.id}} interpolation
│       │   ├── audit.go                        # Audit logging for write operations (includes impersonated_by)
│       │   ├── ratelimit.go                    # Rate limiting middleware (Fiber limiter, tiered: global + auth)
│       │   ├── ipwhitelist.go                  # IP whitelist middleware (exact IP + CIDR, admin-only or global)
│       │   └── tenant.go                       # Multi-tenancy middleware (header/subdomain/path)
│       ├── graphql/
│       │   ├── schema.go                       # SchemaBuilder — auto-generates GraphQL schema from model definitions (types, queries, mutations)
│       │   ├── resolver.go                     # Resolver — CRUD resolvers with permission + record rule enforcement, context-based user ID
│       │   ├── handler.go                      # Fiber HTTP handler for GraphQL at POST /api/v1/graphql
│       │   └── schema_test.go                  # 5 tests (empty schema, model fields, skip non-graphql, mutations, field type mapping)
│       ├── template/
│       │   ├── engine.go                       # Go html/template engine — LoadDirectory, Render, RenderWithLayout
│       │   └── engine_test.go                  # 5 tests
│       ├── view/
│       │   ├── renderer.go                     # View renderer — list, form, kanban, calendar, chart, custom (SSR)
│       │   ├── spa.go                          # SPA envelope renderer — IsSPAMode, RenderSPAEnvelope, SPAEnvelope, SPAMeta, per-view data loaders
│       │   ├── spa_test.go                     # 11 SPA tests — parser, envelope, data parity, SSR backwards compat
│       │   ├── auto_page_generator.go          # GenerateListView/GenerateFormView — auto-generates pages from model fields
│       │   ├── auto_page_generator_test.go     # 6 tests (list generation, form with tabs, auto_pages detection)
│       │   ├── component_compiler.go           # Compiles view JSON into Stencil Web Component HTML. CompileListDatatable() emits bc-datatable with permissions
│       │   └── component_compiler_test.go      # Component compiler tests
│       ├── admin/                              # Split into 7 files for maintainability
│       │   ├── admin.go                        # Core: types, constructor, RegisterRoutes, dashboard, sidebar (with Security section), CSS, helpers
│       │   ├── admin_models.go                 # Model pages: list, view (5 tabs: Form/Fields/Connections/Schema/API), data. Fields tab with MASK/GROUPS columns. API tab with config + generated endpoints preview
│       │   ├── admin_modules.go                # Module pages: list, view (3 tabs: Overview/Permissions/Menu)
│       │   ├── admin_views.go                  # View pages: list, detail (4 tabs: Info/Preview/Editor/Revisions)
│       │   ├── admin_audit.go                  # Health, login history, request log, impersonate/stop-impersonate
│       │   ├── admin_groups.go                 # Group pages: list + detail with 7 Odoo-style tabs (Users/Inherited/Menus/Pages/Access Rights/Record Rules/Notes). Group CRUD API
│       │   ├── admin_security.go               # Security sync page: Load/Export/Upload/Download buttons, history table with rollback. API handlers
│       │   └── admin_api.go                    # View API handlers, data revision handlers
│       ├── assets/
│       │   └── handler.go                      # Static asset serving
│       └── websocket/
│           ├── hub.go                          # WebSocket hub — connect to EventBus, broadcast domain events, route CRUD messages
│           ├── crud.go                         # CRUDHandler — CRUD over WebSocket with request/reply protocol, permission + record rule enforcement
│           └── crud_test.go                    # 7 tests (create, list, read, update, delete, model enable, permission denied)
│
├── pkg/                                        # Public packages (reusable outside engine)
│   ├── ddd/                                    # Domain-Driven Design building blocks
│   │   ├── entity.go                           # Entity interface + BaseEntity
│   │   ├── aggregate.go                        # Aggregate interface + BaseAggregate (with domain events)
│   │   ├── domain_event.go                     # DomainEvent interface + BaseDomainEvent
│   │   ├── repository.go                       # Repository[T] generic interface
│   │   ├── value_object.go                     # ValueObject interface
│   │   └── ddd_test.go                         # 3 tests
│   ├── security/                               # Auth & crypto utilities
│   │   ├── password.go                         # HashPassword (bcrypt), CheckPassword
│   │   ├── jwt.go                              # GenerateToken (with options), ValidateToken, Claims (ImpersonatedBy, Purpose)
│   │   ├── otp.go                              # GenerateOTP — crypto-secure 6-digit code
│   │   ├── encryption.go                       # FieldEncryptor — AES-256-GCM encrypt/decrypt with key versioning
│   │   └── security_test.go                    # 5 tests
│   ├── email/                                  # Email sending
│   │   ├── sender.go                           # SMTPSender — SMTP with TLS, NoopSender fallback
│   │   └── templates.go                        # HTML email templates (OTP code)
│   └── plugin/                                 # Plugin SDK
│       └── proto/
│           └── plugin.proto                    # gRPC service definition for plugins
│
├── embedded/                                   # Compiled-in modules
│   ├── embed.go                                # go:embed directive for modules/
│   ├── embed_test.go                           # Embed tests
│   └── modules/                                # Embedded module files (base, auth)
│       ├── base/                               # Core module (always available)
│       └── auth/                               # Auth module — login, register, forgot, reset, 2FA (i18n x11)
│
├── modules/                                    # Built-in modules (on disk)
│   ├── base/                                   # Core module — users, groups, model_access, security_history, record_rules, settings
│   │   ├── module.json                         # 11 permissions, 2 groups (user, manager), menu
│   │   ├── models/                             # user (is_superuser), role, group (share, comment, module), permission, model_access (12 ERPNext permissions), security_history (ir_security_histories), record_rule (module, modified_source), audit_log, setting
│   │   ├── apis/                               # auth_api, user_api, group_api, role_api, permission_api, etc.
│   │   ├── views/                              # CRUD views for all base models
│   │   ├── data/                               # default_roles, default_groups, default_users
│   │   └── templates/                          # Default UI templates
│   │       ├── layout.html                     # Base layout (full page)
│   │       ├── layout-app.html                 # App layout (sidebar + navbar + content)
│   │       ├── partials/                       # Reusable: sidebar, navbar, pagination, status_badge, actions
│   │       └── views/                          # View templates: list, form, kanban, calendar, chart, login, home
│   ├── crm/                                    # CRM module — contacts, leads
│   │   ├── module.json                         # securities glob, menu with groups
│   │   ├── models/                             # contact (api config, mask on phone), lead (api config)
│   │   ├── securities/                         # crm_user.json (ACL + rules), crm_manager.json (full access)
│   │   ├── apis/                               # contact_api, lead_api (override auto-generated)
│   │   └── views/                              # contact_list, lead_list
│   └── sales/                                  # Sales module — orders
│       ├── module.json
│       ├── models/                             # order, order_line
│       ├── apis/                               # order_api
│       ├── processes/                          # confirm_order
│       ├── views/                              # order_form, order_list
│       └── i18n/                               # Indonesian translations
│
├── plugins/                                    # Plugin runtimes
│   ├── typescript/
│   │   └── index.js                            # Node.js JSON-RPC server (stdin/stdout)
│   └── python/
│       └── runtime.py                          # Python JSON-RPC server (stdin/stdout)
│
├── Dockerfile                                  # Multi-stage build (Go build → minimal runtime)
├── docker-compose.yml                          # Engine + PostgreSQL + Redis
├── Makefile                                    # build, cli, dev, test, lint, clean, tidy
├── go.mod                                      # Go 1.24+, Fiber v2, GORM, Viper, Cobra, JWT, Redis
└── go.sum
```

---

## Packages / Components (`packages/components/`)

Stencil.js Web Component library (`@bitcode/components`).

---

## Packages / UI (`packages/ui/`)

Cross-framework BitCode UI library using Mitosis source components compiled to React, Vue, Stencil, and React Native. Contract-first API per ADR-007/008.

```
packages/ui/
├── package.json                                # @bitcode/ui (Mitosis-driven cross-framework UI library, 177 components)
├── mitosis.config.js                           # Targets: stencil, react, vue, react-native (per ADR-006)
├── tokens/                                     # Design tokens (theme.css = source of truth)
│   ├── theme.css                               #   --bc-* HSL channel variables (light + .dark)
│   ├── tailwind.preset.cjs                     #   Tailwind v3 preset (auto-generated)
│   ├── tailwind.theme.css                      #   Tailwind v4 @theme inline (auto-generated)
│   └── nativewind.preset.cjs                   #   NativeWind preset for React Native
├── src/
│   ├── core/                                   # Reactive store, registry, context, data-source, types, wc-event-bridge, part-props, renderer, providers, collection contract, haptics, animation runtime, adapt runtime
│   │   ├── hooks/                              # Portable headless SSR-safe hooks (ADR-002 §9.2, ADR-013). P0: useId, useDisclosure, useControllableState, useDebouncedValue, useClipboard, useElementSize, useFocusTrap, useMediaQuery, useAppContext/useAppConfigStore. P1 (ADR-013 §3): usePrevious, useIntersection, useInterval, useTimeout, useClickOutside, useKeyboard, useScrollLock, useLocalStorage — all gated behind typeof window/document, barrel-exported via core/hooks/index.ts.
│   │   ├── animation/                          # Animation runtime (ADR-002 §9.3, ADR-013 §4). types.ts = AnimationProps/AnimationKind/AnimationEasing/ANIMATION_DEFAULT contract slot. runtime.ts = resolveAnimation() merges props with default + collapses to 'none' on prefers-reduced-motion (SSR-safe), returns ResolvedAnimation {kind,duration,easing,cssTransition}. getAnimationClasses() returns enter/exit Tailwind utility strings. prefersReducedMotion()/detectPrefersReducedMotion() exported detector. Barrel core/animation/index.ts.
│   │   ├── block-renderer/                     # Block Catalog renderer web binding (ADR-017 section 3, Agent H). web-registry.ts = WebBlockComponentRegistry (componentKey -> Mitosis component + metadata, the ONLY place holding framework components). BcBlockRenderer.lite.tsx = thin Adapt-tier component translating portable render plans into BcRenderer/registry output. index.ts = barrel re-exporting web + portable surface.
│   │   └── adapt/                              # Adapt runtime (ADR-002 §9.4, ADR-013 §4). types.ts = AdaptProps/AdaptCondition/AdaptTarget contract slot. runtime.ts = resolveAdapt(props?, ctx?) evaluates when-condition (breakpoint via createMediaQueryWatcher OR platform: tag) and returns ResolvedAdapt {applies,target}. AdaptContext injectable for test-friendly/SSR-safe platform+breakpoint detection. Barrel core/adapt/index.ts.
│   ├── spa-runtime.ts                          # ADR-014 Phase 2 SPA bootstrap — fetch JSON envelope from engine (render:'spa'), map ViewDefinition → UiSchema, build BcRendererProps + ReactiveStore. SPARuntime.load() / buildDescriptor() / viewDefinitionToUiSchema(). Reuses block-renderer + core/registry.
│   ├── spa-runtime.test.ts                     # 13 source-level SPA runtime tests (node:test) — envelope parsing, mapping for list/form/custom, ReactiveStore seeding, fetcher errors, humanize helper
│   ├── components/                             # 189 Mitosis components across 14 groups (charts, data, data-display, feedback, forms, input, layout, media, navigation, overlays, primitives, print, social, views, adapt-poc). 177 .lite.tsx Pure + 12 .web.tsx/.native.tsx Adapt (ADR-006).
│   └── themes/default.css                      # Default theme bundle
├── scripts/                                    # Build & tooling helpers
│   ├── contract-audit.mjs                      # ADR-007/008/009 contract compliance auditor (blocking mode, Props-interface-aware parsing)
│   ├── audit-touch-targets.mjs                 # Touch target size auditor for .native.tsx files (44pt minimum, brace-depth-aware tag extraction)
│   ├── generate-component-docs.mjs             # ADR-008 Phase 4 per-component API doc generator (parses .lite.tsx + .contract.ts → markdown)
│   ├── generate-a11y-tests.mjs                 # A11y runtime test generator (axe-core smoke tests for all components)
│   ├── generate-tokens.mjs                     # Design token sync (theme.css ↔ tailwind/nativewind presets)
│   ├── generate-index.mjs                      # Component index generator (src/index.ts + per-target index)
│   ├── generate-registry.mjs                   # Component registry generator
│   ├── sync-regular-ts.mjs                     # Mirrors .ts files across all 4 output targets
│   ├── apply-adapt-overrides.mjs               # Copies hand-tuned native impls (7 adapt components) into react-native output
│   ├── reformat-react-output.mjs               # Post-processes React output for hooks compatibility
│   ├── run-source-tests.mjs                    # Source-level test runner (node:test)
│   └── dev.mjs                                 # Dev server
├── docs/                                       # ADR-008 Phase 4 documentation
│   ├── api/                                    # 177 per-component API docs + categorized sitemap README
│   └── guides/                                 # 4 framework integration guides (React, Vue, RN, Stencil) + index README
├── output/                                     # Generated targets (gitignored — produced by mitosis build)
├── demo-app/                                   # Production Vite + React demo app (all 177 components, 604 sections)
│   ├── package.json                            #   @bitcode/ui-demo — Vite + React + Tailwind v4
│   ├── vite.config.ts                          #   Core-import redirect plugin + Mitosis runtime stub alias
│   ├── src/
│   │   ├── App.tsx                             #   Main demo shell — sidebar nav + section renderer
│   │   ├── main.tsx                            #   React entry — AppProvider + App
│   │   ├── index.css                           #   Tailwind v4 import + demo-specific styles
│   │   ├── mitosis-browser-stub.ts             #   Browser-safe stub for non-Mitosis .ts files importing @builder.io/mitosis
│   │   ├── react-native-web-stub.tsx           #   react-native → react-dom shim for web demo
│   │   ├── bitcode-ui.d.ts                     #   Ambient module declarations
│   │   ├── demos/                              #   Per-component demo files (one per component group)
│   │   ├── mocks/                              #   Mock data for demos
│   │   └── utils/                              #   Demo utilities
│   └── tsconfig.json                           #   Strict TS config with path aliases
└── demos/
    └── react-app/                              # Legacy Vite + React playground (Playwright e2e)
```

> Lynx target (formerly ADR-001 §15) was dropped per ADR-006. Mobile is now React
> Native + NativeWind, sharing the same Mitosis pipeline as React/Vue.

---

## UI Schema (`packages/ui-schema/`)
Standalone UI schema definitions — domain-agnostic, framework-neutral types for schema-driven UI. Zero dependencies, zero platform globals. Extracted from `@bitcode/ui` per ADR-012 to provide a canonical source for all renderers (React/Vue/RN via Mitosis, Flutter via Dart port).
```
packages/ui-schema/
├── package.json                                # @bitcode/ui-schema (private, zero-dep, file: dependency of @bitcode/ui)
├── README.md
└── src/
    ├── index.ts                                # Root barrel — re-exports all types + ReactiveStore + resolver
    ├── types.ts                                # Core types: BcChangeMeta, BcEventSource, DataSource, BcKey, BcSelection, field events
    ├── reactive/
    │   ├── index.ts                            # Reactive barrel
    │   ├── types.ts                            # Reaction, Effect, Evaluator, Snapshot, FieldMeta, ValidationRule, ReactiveScopeLike
    │   └── store.ts                            # ReactiveStore class (headless reactive engine, pure TS)
    ├── expressions/                             # Portable Expression AST (ADR-018 S5-S6) - canonical types, validator, serializer, fixtures
    │   ├── index.ts                            # Expressions barrel - re-exports all expression types/functions
    │   ├── ast.ts                              # 5-node-kind discriminated union (Literal/Path/Op/Call/If), type guards, constructors, versioned envelope
    │   ├── operators.ts                        # PortableOperator (24 ops, 6 families), OPERATOR_SPECS table with arity bounds, family helpers
    │   ├── scope.ts                            # ExpressionScope, PermissionSnapshot, parsePath (reserved-roots + forbidden-segments enforcement)
    │   ├── errors.ts                           # ExpressionErrorCode (14 codes), ExpressionError with nodePath/suggestions, depth/step caps
    │   ├── permission-scope.ts                 # buildPermissionFilteredScope (drops unauthorized fields), pathAuthorization
    │   ├── evaluator.ts                        # ExpressionEvaluator interface (pure contract for host adapters)
    │   ├── function-registry.ts                # FunctionManifestEntry (pure, deterministic, timezone, locale), FunctionRegistry interface
    │   ├── builtin-manifest.ts                 # 70+ built-in functions (math/string/array/date/format/validate), all platforms
    │   ├── golden-types.ts                     # GoldenFixture, GoldenCase, GoldenExpect types for cross-platform parity fixtures
    │   ├── validate.ts                         # validateExpression, validatePath, FORBIDDEN_FUNCTION_NAMES denylist, suggestFunctions (Levenshtein)
    │   └── serialize.ts                        # serializeExpression (versioned envelope), deserializeExpression (bare+envelope), roundTripExpression
    └── schema/
        ├── index.ts                            # Schema barrel
        ├── types.ts                            # UiSchema, FieldNode, FieldType, UiAction, BcRendererProps
        └── resolver.ts                         # FieldResolver, resolveField, resolveFieldType, resolveFields
    └── block-renderer/                         # Block Renderer portable engine (ADR-017 section 3, Agent H). Emits RENDER PLANS only.
        ├── index.ts                            # Block-renderer barrel (BlockRenderer, registry contract, types, validator)
        ├── types.ts                            # BlockNode, BlockBindings, RenderPlan union (Schema/Component/Hidden/Diagnostic), ResolvedBlockTree, ActionHandlerPlan (deferred), diagnostic codes
        ├── registry-contract.ts                # BlockComponentRegistry interface (METADATA ONLY), InMemoryBlockComponentRegistry
        ├── validate-runtime.ts                 # Runtime host validation (required components/functions/adapters) -> BLOCK_RENDER_* diagnostics
        └── renderer.ts                         # BlockRenderer engine: catalog resolution, readiness gate, F evaluator bindings, permission-filtered scope, strategy switch -> render plan, deferred action binding
```

---

```
packages/components/
├── package.json                                # @bitcode/components v0.1.0
├── stencil.config.ts                           # Namespace: bc-components, output: dist + dist-custom-elements + www
├── tsconfig.json
│
├── scripts/
│   └── generate-exports.mjs                    # Post-build: generates group barrels in dist/, full.js, theme copy, package.json exports (151 paths)
│
└── src/
    ├── components.d.ts                         # Auto-generated type declarations
    ├── declarations.d.ts                       # Module declarations
    │
    ├── core/                                   # Shared infrastructure
    │   ├── types.ts                            # FieldType (30+ types), WidgetType, event interfaces, FetchParams/FetchResult, ValidationResult, BcConfig, BcTheme, BcRtl, ThemeChangeCallback
    │   ├── bc-setup.ts                         # BcSetup singleton — global config (auth, headers, theme, RTL, validators, reactivity rules). setTheme/getTheme/onThemeChange/setRtl/isRtl
    │   ├── bc-native.ts                        # BcNative — bridge abstraction for native capabilities (camera, GPS, SQLite, barcode, biometrics). Detects Tauri vs browser.
    │   ├── bc-native.spec.ts                   # BcNative tests (10 tests — environment detection, browser fallbacks)
    │   ├── code-editor-theme.ts                # Shared CodeMirror theme extensions — IDE-style syntax highlighting, gutters, caret, selection, and editor chrome tokens.
    │   ├── chart-theme.ts                      # Pure chart theme patching helpers — resolves active theme, reads CSS vars, and applies text/axis/legend color defaults.
    │   ├── chart-utils.ts                      # Shared ECharts bootstrap, option builders, and runtime theme sync for canvas chart text, axis, legend, and toolbox colors.
    │   ├── chart-utils.spec.ts                 # Chart theming tests for title, axis, parallel-axis, and series label color patching.
    │   ├── offline-store.ts                    # OfflineStore — CRUD routing layer. SQLite for offline models, fetch() for online. Outbox recording.
    │   ├── offline-store.spec.ts               # OfflineStore tests (11 tests — routing, CRUD, outbox, table mapping)
    │   ├── data-fetcher.ts                     # 4-level data fetching (local, URL, event intercept, custom fetcher). Standalone — uses native fetch()
    │   ├── validation-engine.ts                # 3-level validation pipeline (built-in, custom JS, server-side). Uses validators.ts
    │   ├── field-utils.ts                      # Shared field utilities (dirty/touched tracking, ARIA attrs, CSS classes, FormProxy, debounce)
    │   ├── field-base.ts                      # FieldBase<V> helper class (composition pattern for field components)
    │   ├── dialog-service.ts                  # Programmatic BcDialog service (SweetAlert-style: alert, confirm, prompt)
    │   ├── api-client.ts                       # HTTP client for engine REST APIs (BitCode-specific, optional fallback)
    │   ├── event-bus.ts                        # Cross-component event bus
    │   ├── form-engine.ts                      # Form state management, validation, submission (BitCode-specific)
    │   └── i18n.ts                             # Client-side i18n utilities
    │
    ├── registry.ts                             # Component registry — 122 components, 11 groups, dependency graph, include/exclude selector
    ├── full.ts                                 # Re-exports Registry + types. defineAll() generated post-build into dist/full.js
    │
    ├── global/
    │   ├── global.css                          # CSS custom properties (80+ tokens), reset styles, --bc-direction var, auto dark mode
    │   └── themes/
    │       ├── light.css                       # Explicit light theme (data-bc-theme="light")
    │       └── dark.css                        # Dark theme overrides (data-bc-theme="dark")
    │
    ├── i18n/                                   # Translation files
    │   ├── index.ts                            # i18n initialization (global script)
    │   ├── en.json                             # English
    │   ├── id.json                             # Indonesian
    │   ├── ar.json                             # Arabic
    │   ├── de.json                             # German
    │   ├── es.json                             # Spanish
    │   ├── fr.json                             # French
    │   ├── ja.json                             # Japanese
    │   ├── ko.json                             # Korean
    │   ├── pt-BR.json                          # Portuguese (Brazil)
    │   ├── ru.json                             # Russian
    │   └── zh-CN.json                          # Chinese (Simplified)
    │
    ├── utils/                                  # Shared utilities
    │   ├── expression-eval.ts                  # Expression evaluator for computed fields
    │   ├── format.ts                           # Number, date, currency formatting
    │   └── validators.ts                       # Field validation rules
    │
    └── components/                             # Web Components (each has .tsx + .css)
        ├── fields/                             # 30+ field components
        │   ├── bc-field-string/                # Text input
        │   ├── bc-field-smalltext/             # Small textarea
        │   ├── bc-field-text/                  # Large textarea
        │   ├── bc-field-richtext/              # Tiptap rich text editor
        │   ├── bc-field-markdown/              # Markdown editor with preview
        │   ├── bc-field-html/                  # HTML editor
        │   ├── bc-field-code/                  # CodeMirror code editor
        │   ├── bc-field-password/              # Password input
        │   ├── bc-field-integer/               # Integer input
        │   ├── bc-field-float/                 # Float input
        │   ├── bc-field-decimal/               # Decimal input
        │   ├── bc-field-currency/              # Currency input with formatting
        │   ├── bc-field-percent/               # Percentage input
        │   ├── bc-field-checkbox/              # Checkbox
        │   ├── bc-field-toggle/                # Toggle switch
        │   ├── bc-field-select/                # Dropdown select
        │   ├── bc-field-radio/                 # Radio buttons
        │   ├── bc-field-multicheck/            # Multi-checkbox
        │   ├── bc-field-tags/                  # Tag input
        │   ├── bc-field-date/                  # Date picker
        │   ├── bc-field-time/                  # Time picker
        │   ├── bc-field-datetime/              # DateTime picker
        │   ├── bc-field-duration/              # Duration input
        │   ├── bc-field-file/                  # File upload
        │   ├── bc-field-image/                 # Image upload with preview
        │   ├── bc-field-signature/             # Signature pad (signature_pad)
        │   ├── bc-field-barcode/               # Barcode/QR generator (JsBarcode + QRCode)
        │   ├── bc-field-color/                 # Color picker
        │   ├── bc-field-geo/                   # Geolocation (Leaflet map)
        │   ├── bc-field-rating/                # Star rating
        │   ├── bc-field-json/                  # JSON editor
        │   ├── bc-field-link/                  # Many2one link field
        │   ├── bc-field-dynlink/               # Dynamic link field
        │   ├── bc-field-tableselect/           # Table multi-select
        │   └── field-base.css                  # Shared field styles
        │
        ├── layout/                             # Layout components
        │   ├── bc-row/                         # Flex row
        │   ├── bc-column/                      # Flex column
        │   ├── bc-section/                     # Collapsible section
        │   ├── bc-tabs/ + bc-tab/              # Tab container + tab panels
        │   ├── bc-sheet/                       # Card/sheet container
        │   ├── bc-header/                      # Section header
        │   ├── bc-separator/                   # Visual separator
        │   ├── bc-button-box/                  # Button group
        │   └── bc-html-block/                  # Raw HTML block
        │
        ├── views/                              # View components
        │   ├── bc-view-list/                   # Data table view
        │   ├── bc-view-form/                   # Record form view
        │   ├── bc-view-kanban/                 # Kanban board
        │   ├── bc-view-calendar/               # Calendar (FullCalendar)
        │   ├── bc-view-gantt/                  # Gantt chart (frappe-gantt)
        │   ├── bc-view-map/                    # Map view (Leaflet)
        │   ├── bc-view-tree/                   # Tree/hierarchy view
        │   ├── bc-view-report/                 # Report view
        │   └── bc-view-activity/               # Activity stream view
        │
        ├── kanban/                             # Enterprise kanban components (Trello/Jira-level)
        │   ├── core/                           # Shared utilities
        │   │   ├── kanban-types.ts             # 25+ interfaces (cards, columns, events, configs)
        │   │   ├── kanban-data-fetcher.ts      # Shared 4-layer CRUD (kanbanFetch/Create/Update/Remove/Upload)
        │   │   └── mention-parser.ts           # @mention parse/highlight/extract utilities
        │   ├── bc-kanban-board/                # Main board — drag-drop, column CRUD, card detail panel, inline dialogs
        │   ├── bc-kanban-card-subtasks/        # Checklist with progress bar, toggle/create/delete
        │   ├── bc-kanban-card-comments/        # Comments with @mention dropdown, inline file attachments
        │   ├── bc-kanban-card-attachments/     # File list with upload, thumbnails, delete
        │   └── bc-kanban-card-activity/        # Timeline with relative time (read-only)
        │
        ├── charts/                             # Chart components (23 ECharts + 3 pure CSS)
        │   ├── bc-chart-line/                  # Line chart
        │   ├── bc-chart-bar/                   # Bar chart
        │   ├── bc-chart-pie/                   # Pie/donut chart
        │   ├── bc-chart-area/                  # Area chart
        │   ├── bc-chart-scatter/               # Scatter/bubble chart
        │   ├── bc-chart-radar/                 # Radar chart
        │   ├── bc-chart-gauge/                 # Gauge chart
        │   ├── bc-chart-funnel/                # Funnel chart
        │   ├── bc-chart-heatmap/               # Heatmap
        │   ├── bc-chart-treemap/               # Treemap (hierarchical)
        │   ├── bc-chart-sunburst/              # Sunburst chart (hierarchical)
        │   ├── bc-chart-candlestick/           # Candlestick/OHLC chart
        │   ├── bc-chart-boxplot/               # Box plot (statistics)
        │   ├── bc-chart-mixed/                 # Mixed/combo chart (bar+line+scatter)
        │   ├── bc-chart-sankey/                # Sankey flow diagram
        │   ├── bc-chart-graph/                 # Network/relationship graph
        │   ├── bc-chart-tree/                  # Tree diagram (org chart)
        │   ├── bc-chart-polar/                 # Polar coordinate chart
        │   ├── bc-chart-parallel/              # Parallel coordinates
        │   ├── bc-chart-themeriver/            # Theme river/stream graph
        │   ├── bc-chart-pictorialbar/          # Pictorial bar chart
        │   ├── bc-chart-custom/                # Custom chart (raw ECharts option)
        │   ├── bc-chart-pivot/                 # Pivot table (pure CSS)
        │   ├── bc-chart-kpi/                   # KPI card (pure CSS)
        │   ├── bc-chart-scorecard/             # Scorecard
        │   └── bc-chart-progress/              # Progress indicator (pure CSS)
        │
        ├── datatable/                          # Data table components
        │   ├── bc-datatable/                   # Full-featured data table
        │   ├── bc-child-table/                 # Inline child table (one2many)
        │   ├── bc-filter-bar/                  # Filter bar
        │   ├── bc-filter-panel/                # Filter panel
        │   ├── bc-filter-builder/              # Advanced filter builder
        │   └── bc-lookup-modal/                # Record lookup modal
        │
        ├── dialogs/                            # Dialog components
        │   ├── bc-dialog-modal/                # Generic modal
        │   ├── bc-dialog-confirm/              # Confirmation dialog
        │   ├── bc-dialog-alert/                # Alert dialog (OK only)
        │   ├── bc-dialog-prompt/               # Prompt dialog (with input)
        │   ├── bc-dialog-quickentry/           # Quick record entry
        │   ├── bc-dialog-wizard/               # Multi-step wizard
        │   └── bc-toast/                       # Toast notifications
        │
        ├── search/                             # Search & filter components
        │   ├── bc-search/                      # Search bar
        │   └── bc-favorites/                   # Saved filters
        │
        ├── widgets/                            # Widget components
        │   ├── bc-badge/                       # Status badge
        │   ├── bc-copy/                        # Copy to clipboard
        │   ├── bc-phone/                       # Phone link
        │   ├── bc-email/                       # Email link
        │   ├── bc-url/                         # URL link
        │   ├── bc-progress/                    # Progress bar
        │   ├── bc-statusbar/                   # Status bar (workflow states)
        │   ├── bc-priority/                    # Priority indicator
        │   ├── bc-handle/                      # Drag handle
        │   ├── bc-domain/                      # Domain filter widget
        │   ├── bc-placeholder/                 # Placeholder component
        │   └── bc-sync-status/                 # Sync status indicator
        │
        ├── media/                              # Media viewer/player components
        │   ├── bc-viewer-pdf/                  # PDF viewer (iframe + object fallback)
        │   ├── bc-viewer-image/                # Image viewer (zoom, lightbox)
        │   ├── bc-viewer-document/             # Office doc viewer (Microsoft/Google iframe)
        │   ├── bc-viewer-youtube/              # YouTube embed (full video + shorts)
        │   ├── bc-viewer-instagram/            # Instagram post embed (oEmbed)
        │   ├── bc-viewer-tiktok/               # TikTok video embed (oEmbed)
        │   ├── bc-viewer-video/                # HTML5 video player (mp4, webm, ogg)
         │   └── bc-viewer-audio/                # HTML5 audio player (mp3, m4a, aac, ogg, webm)
         │
         ├── print/                              # Print & export
        │   ├── bc-export/                      # Data export (XLSX)
        │   ├── bc-print/                       # Print view
        │   └── bc-report-link/                 # Report link
        │
        └── social/                             # Social/collaboration
            ├── bc-activity/                    # Activity log
            ├── bc-chatter/                     # Comment thread
            └── bc-timeline/                    # Timeline view
```

---

## Flutter Adapter (`packages/ui-flutter/`)
Dart port of the BitCode UI layer per ADR-011 Phase 5a/5b/5c + ADR-016 Phase 4. Pure-Dart reactive engine (no Flutter imports) with thin Flutter `ChangeNotifier` adapter, auto-generated theme tokens, UiSchema types port (freezed + json_serializable), contract→Dart stub generator, and 160+ cross-platform widgets (130 Phase 5c + 28 ADR-016 Phase 4 stub-completion + BcChart). Uses `flutter_platform_widgets` (FPW) for Material/Cupertino adaptation where a Cupertino counterpart exists; Material-only for widgets without one.
```
packages/ui-flutter/
├── pubspec.yaml                                # bc_ui_flutter — depends on flutter, flutter_platform_widgets 0xfw, flutter_echarts 2.5 (ADR-015), flutter_test, freezed, json_serializable, build_runner, collection, webview_flutter_platform_interface (dev, test stub)
├── analysis_options.yaml                       # flutter_lints, prefer_const, excludes generated *.freezed.dart/*.g.dart
├── lib/
│   ├── bc_ui_flutter.dart                      # Public barrel — exports 160+ components (24 P0 + 8 P1 Phase 5b.2 + ~98 Phase 5c M1-M7 + 28 ADR-016 Phase 4 + BcChart) + schema types + core + theme + BcApp + spa_runtime
│   └── src/
│       ├── bc_app.dart                         # PlatformApp wrapper with BcTheme defaults + variant override
│       ├── spa_runtime.dart                    # ADR-014 Phase 2 SPA bootstrap — SPARuntime / SPAEnvelope / SPAMeta / viewDefinitionToUiSchema / SPADescriptor. Fetch JSON envelope via dart:io HttpClient (injectable), map ViewDefinition → UiSchema, build BcRendererProps + ReactiveStore. Mirrors packages/ui/src/spa-runtime.ts.
│       ├── theme/
│       │   └── bc_theme.dart                   # AUTO-GENERATED by token bridge (25 light + 25 dark const Color, 4 ColorScheme factories, radius)
│       ├── core/
│       │   ├── types.dart                      # BcEventSource, BcChangeReason (18 values), BcChangeMeta, BcChangeHandler, BcComponentSize, BcComponentVariant; re-exports chart/ chart_types.dart (ADR-015)
│       │   ├── reactive_store.dart             # Pure-Dart reactive engine: VALUE/UI/VALIDATION/DATA effects, cascade propagation, snapshot with $self/$root, async DATA fetcher, dispose
│       │   ├── reactive_store_notifier.dart    # Thin Flutter adapter: ChangeNotifier bridging ReactiveStore subscriptions to notifyListeners()
│       │   ├── app_config_snapshot.dart        # ADR-013 Phase 1 — AppConfigSnapshot + enums (AppThemeMode, ResolvedAppTheme, AppDirection, AppSize, AppMessageQueue) — pure data, parity with TS AppConfigSnapshot
│       │   ├── app_config_store.dart           # ADR-013 Phase 1 — headless AppConfigStore (getSnapshot/setTheme/setLocale/setDirection/setSize/setDateLocale/setMessageQueue/subscribe/dispose/isDisposed)
│       │   ├── bc_app_config.dart              # ADR-013 Phase 1 — BcAppConfig InheritedWidget exposing BcAppConfigData (snapshot + store) to descendants; BcAppConfig.of(context) parity with TS useAppContext
│       │   ├── bc_breakpoints.dart             # ADR-013 Phase 5b — bcBreakpointQueries (sm/md/lg/xl/2xl Tailwind-aliased min-width queries), parity with TS BC_BREAKPOINTS
│       │   ├── hooks/                          # ADR-013 Phase 5b — Dart P1 hook parity (hooks.dart barrel)
│       │   │   ├── hooks.dart                  # Barrel — re-exports use_previous, use_interval, use_timeout, use_scroll_lock, use_local_storage
│       │   │   ├── use_previous.dart           # PreviousValue<T> + usePrevious() — pure Dart port of TS use-previous (ref cache, no DOM)
│       │   │   ├── use_interval.dart           # BcIntervalController + createInterval() — Timer.periodic parity with TS createInterval (injectable createTimer for FakeAsync tests)
│       │   │   ├── use_timeout.dart            # BcTimeoutController + createTimeout() — Timer parity with TS createTimeout (start/clear/dispose)
│       │   │   ├── use_scroll_lock.dart        # BcScrollLock + createScrollLock() — ref-counted lock exposing ValueNotifier<bool> (Flutter has no document.body; host swaps ScrollPhysics)
│       │   │   └── use_local_storage.dart      # BcLocalStorageStore<T> + createLocalStorage() — injectable BcStorageBackend (default in-memory, parity with TS SSR no-op); JSON round-trip, malformed-as-default
│       │   ├── animation/                      # ADR-013 Phase 5b — Dart animation runtime parity
│       │   │   └── animation.dart              # BcAnimationKind enum, BcAnimationEasing, BcAnimationProps, ANIMATION_DEFAULT, BcResolvedAnimation, resolveAnimation() (reduced-motion via injectable flag → host passes MediaQuery.disableAnimationsOf), getAnimationCurves()
│       │   ├── adapt/                          # ADR-013 Phase 5b — Dart adapt runtime parity
│       │   │   └── adapt.dart                  # BcAdaptBreakpoint/Platform/Target, BcAdaptCondition, BcAdaptProps, BcAdaptContext, BcResolvedAdapt, BcAdaptResolver, resolveAdapt() (breakpoint via injectable matcher → MediaQuery.size, platform via ctx.platform → Platform.isIOS/Android)
│       │   └── chart/                          # Dart port of TS chart-option contract (ADR-015 F2) — pure functions, no Flutter deps
│       │       ├── chart_types.dart            # BcChartType enum (26 variants + BC_CHART_TYPES), BcChartItem, BcChartSeries, BuildOptionContext, ChartData, EChartsOption typedefs
│       │       ├── chart_helpers.dart          # parseColors(), parseChartData(), isSingleSeries(), isMultiSeries(), buildTitleOption/Tooltip/Legend/Toolbox/DataZoom/GridOption(), DEFAULT_PALETTE (7 token colors)
│       │       ├── chart_factory.dart          # buildChartOption(ctx) — 26-type ECharts option factory, pure function, _stripNulls for TS undefined parity
│       │       └── chart_contract.dart         # Barrel export (mirrors TS contracts/chart-option/index.ts)
│       ├── schema/                             # UiSchema types port (ADR-011 Phase 5b.1) — freezed + json_serializable
│       │   ├── ui_schema_types.dart            # FetchParams, FetchResult, SortRule, BcSelection (sealed), DataSourceFetcher, enums (ThemeMode, ResolvedTheme, ComponentSize, ComponentVariant, InputType, InputSize, ValidationStatus), BcKey, DataSource
│       │   ├── ui_schema_types.freezed.dart    # AUTO-GENERATED by build_runner
│       │   ├── ui_schema_types.g.dart          # AUTO-GENERATED by build_runner
│       │   ├── ui_schema_events.dart           # FieldChangeEvent, FieldInputEvent, FieldFocusEvent, FieldBlurEvent, FieldKeyEvent, FieldClearEvent, BcCustomEvent
│       │   ├── ui_schema_events.freezed.dart   # AUTO-GENERATED by build_runner
│       │   ├── ui_schema_events.g.dart         # AUTO-GENERATED by build_runner
│       │   ├── ui_schema.dart                  # UiSchema, FieldNode, LayoutNode, PageNode, BindingNode, ConditionNode, ValidationError, ViewNode — composite tree with toJson/fromJson
│       │   ├── ui_schema.freezed.dart          # AUTO-GENERATED by build_runner
│       │   ├── ui_schema.g.dart                # AUTO-GENERATED by build_runner
│       │   ├── field_resolver.dart             # FieldResolver — resolves field paths, bindings, conditions against a ReactiveStore snapshot
│       │   ├── field_resolver.freezed.dart     # AUTO-GENERATED by build_runner
│       │   └── field_resolver.g.dart           # AUTO-GENERATED by build_runner
│       ├── expressions/                        # Portable Expression AST (ADR-018 S5-S6) - Dart mirror of TS package. Pure types + validator + scope + manifest. No widgets, no I/O.
│       │   ├── expressions.dart                # Barrel — re-exports all expression types
│       │   ├── ast.dart                        # Manual sealed class (5 kinds: Literal/Path/Op/Call/If) with kind discriminator for TS wire-shape parity, JSON roundtrip, ExpressionParseException
│       │   ├── operators.dart                  # PortableOperator enum (24 ops), operatorSpecs map with arity bounds, family helpers
│       │   ├── errors.dart                     # ExpressionErrorCode enum (14 codes), ExpressionError, depth/step caps
│       │   ├── scope.dart                      # ExpressionScope, PermissionSnapshot, parsePath with reserved-roots + forbidden-segments enforcement
│       │   ├── permission_scope.dart           # buildPermissionFilteredScope (drops unauthorized fields), pathAuthorization
│       │   ├── function_registry.dart          # FunctionManifestEntry, FunctionRegistry abstract, InMemoryFunctionRegistry
│       │   ├── builtin_manifest.dart           # 70+ built-in functions manifest (mirrors TS builtin-manifest.ts)
│       │   ├── validate.dart                   # validateExpression with FORBIDDEN_FUNCTION_NAMES denylist, suggestFunctions (Levenshtein)
│       │   └── evaluator.dart                  # ExpressionEvaluator abstract + EvalResult (pure contract for host adapters)
│       ├── components/
│           ├── bc_input.dart                   # FPW adaptive — text input (P0 hand-written)
│           ├── bc_select.dart                  # FPW adaptive — dropdown select (P0 hand-written)
│           ├── bc_checkbox.dart                # FPW adaptive — checkbox (P0 hand-written)
│           ├── bc_toggle.dart                  # FPW adaptive — switch/toggle (P0 hand-written)
│           ├── bc_segmented.dart               # FPW adaptive — segmented control (P0 hand-written)
│           ├── bc_number_field.dart            # FPW adaptive — numeric input (P0 hand-written)
│           ├── bc_date_field.dart              # FPW adaptive — date picker (P0 hand-written)
│           ├── bc_time_field.dart              # FPW adaptive — time picker (P0 hand-written)
│           ├── bc_dialog.dart                  # FPW adaptive — platform dialog (P0 hand-written)
│           ├── bc_alert_dialog.dart            # FPW adaptive — alert dialog (P0 hand-written)
│           ├── bc_tooltip.dart                 # FPW adaptive — tooltip (P0 hand-written)
│           ├── bc_dropdown.dart                # FPW adaptive — dropdown menu (P0 hand-written)
│           ├── bc_combobox.dart                # Material-only — filterable combobox (P0 hand-written)
│           ├── bc_multi_select.dart            # Material-only — multi-select chips (P0 hand-written)
│           ├── bc_form.dart                    # Material-only — form with ReactiveStore scope (P0 hand-written)
│           ├── bc_drawer.dart                  # Material-only — navigation drawer (P0 hand-written)
│           ├── bc_popover.dart                 # Material-only — popover (P0 hand-written)
│           ├── bc_context_menu.dart            # Material-only — context menu (P0 hand-written)
│           ├── bc_command.dart                 # Material-only — command palette (Cmd+K) (P0 hand-written)
│           ├── bc_tabs.dart                    # Material-only — tab bar + view (P0 hand-written)
│           ├── bc_accordion.dart               # Material-only — expandable accordion (P0 hand-written)
│           ├── bc_splitter.dart                # Material-only — resizable split pane (P0 hand-written)
│           ├── bc_data_table.dart              # Material-only — data table with columns (P0 hand-written)
│           ├── bc_list.dart                    # Material-only — list view (P0 hand-written)
│           ├── bc_disclosure.dart              # P1 Phase 5b.2 — disclosure / disclosure group
│           ├── bc_field_duration.dart          # P1 Phase 5b.2 — duration field
│           ├── bc_field_dyn_link.dart          # P1 Phase 5b.2 — dynamic link field
│           ├── bc_field_morph.dart             # P1 Phase 5b.2 — morph field
│           ├── bc_field_table_select.dart      # P1 Phase 5b.2 — table select field
│           ├── bc_lookup_modal.dart            # P1 Phase 5b.2 — lookup modal
│           ├── bc_portal.dart                  # P1 Phase 5b.2 — portal overlay host
│           ├── bc_view_gantt.dart              # P1 Phase 5b.2 — Gantt view
│           └── bc_chart.dart                   # ADR-015 — ECharts-backed chart widget (26 types via flutter_echarts WebView, replaces CustomPaint)
│           # ADR-016 Phase 4 stub-completion (28 new components):
│           ├── bc_alert.dart                   # Inline alert banner (4 severity, persistent, dismissible)
│           ├── bc_chip.dart                    # Compact tag/chip (6 variants, avatar, removable)
│           ├── bc_print_area.dart              # Printable region marker (semantic landmark)
│           ├── bc_print_preview.dart           # Preview pane + print trigger (iframe on web)
│           ├── bc_text_field.dart              # Rich text input (multiline auto-grow, counter, prefix/suffix)
│           ├── bc_like.dart                    # Like/unlike toggle + count
│           ├── bc_follow.dart                  # Follow/unfollow toggle + count
│           ├── bc_comment.dart                 # Threaded comment with reply form
│           ├── bc_share.dart                   # Share intent (dependency-free delegation on Flutter)
│           ├── bc_native_button.dart           # Platform-idiomatic button (PlatformElevatedButton)
│           ├── bc_native_input.dart            # Typedef → BcInput (already platform-aware)
│           ├── bc_native_select.dart           # Typedef → BcSelect
│           ├── bc_native_switch.dart           # Typedef → BcToggle
│           ├── bc_native_slider.dart           # Typedef → BcSlider
│           ├── bc_native_picker.dart           # FPW picker (CupertinoPicker wheel / Material dialog)
│           ├── bc_native_alert.dart            # FPW alert (CupertinoAlertDialog / Material AlertDialog)
│           ├── bc_native_actionsheet.dart      # FPW action sheet (CupertinoActionSheet / bottom sheet)
│           ├── bc_native_modal.dart            # FPW modal (fullscreen/formSheet/pageSheet)
│           ├── bc_native_toast.dart            # FPW toast (overlay / SnackBar)
│           ├── bc_lightbox.dart                # Full-screen gallery (PageView swipe nav)
│           ├── bc_map.dart                     # Map (dependency-free, caller wires google_maps_flutter)
│           ├── bc_webcam.dart                  # Camera capture (dependency-free, caller wires camera plugin)
│           ├── bc_detail_view.dart             # Master-detail renderer (fields + tabs)
│           ├── bc_grid_view.dart               # Spreadsheet grid (header + rows + onRowSelect)
│           ├── bc_tree_view.dart               # Hierarchical tree (expand/collapse + onSelect)
│           ├── bc_kanban_view.dart             # Kanban board (columns + cards + onCardMove)
│           ├── bc_calendar_view.dart           # Calendar (month/week/day + events + onDateSelect)
│           └── bc_dashboard.dart               # Widget grid (configurable columns + renderWidget)
├── test/
│   ├── types_test.dart                         # 7 tests — BcChangeMeta, BcEventSource, BcChangeReason, BcComponentSize, BcComponentVariant
│   ├── reactive_store_test.dart                # 12 tests — get/set, setMany, metadata, VALUE effects, cascade, conditions, snapshot, subscription, dispose
│   ├── spa_runtime_test.dart                   # 28 tests — SPAEnvelope/SPAMeta parsing, viewDefinitionToUiSchema list/form/detail/custom mapping, SPARuntime.buildDescriptor with ReactiveStore seeding, fetcher error propagation, URL resolution (base + relative, no double slash)
│   ├── bc_input_test.dart                      # 3 tests — text input, disabled, label+error
│   ├── bc_checkbox_test.dart                   # 3 tests — toggle, disabled, label
│   ├── bc_toggle_test.dart                     # 2 tests — label, toggle
│   ├── bc_select_test.dart                     # 2 tests — selected value, placeholder
│   ├── bc_tabs_test.dart                       # 2 tests — labels, content switch
│   ├── bc_accordion_test.dart                  # 2 tests — items, expand
│   ├── bc_data_table_test.dart                 # 2 tests — headers+rows, empty state
│   ├── bc_list_test.dart                       # 3 tests — items, empty state, onTap
│   ├── bc_form_test.dart                       # 2 tests — store provision, reactions
│   ├── bc_disclosure_test.dart                 # Phase 5b.2 tests — disclosure and disclosure group behavior
│   ├── bc_portal_test.dart                     # Phase 5b.2 tests — portal host, overlay, lifecycle behavior
│   ├── bc_field_duration_test.dart             # Phase 5b.2 tests — duration parsing, formatting, interaction
│   ├── bc_field_table_select_test.dart         # Phase 5b.2 tests — table selection behavior
│   ├── bc_field_morph_test.dart                # Phase 5b.2 tests — polymorphic field selection behavior
│   ├── bc_field_dyn_link_test.dart             # Phase 5b.2 tests — dynamic link model and value behavior
│   ├── bc_lookup_modal_test.dart               # Phase 5b.2 tests — lookup modal search, selection, and close behavior
│   ├── bc_view_gantt_test.dart                 # Phase 5b.2 tests — Gantt rendering, tasks, scales, and interactions
│   ├── bc_chart_test.dart                      # ADR-015 F3 — 25 widget-contract tests (WC-001..012): container, loading, empty, semantics, Echarts presence, data normalization, backward-compat; uses fake WebViewPlatform
│   ├── chart_factory_test.dart                 # ADR-015 F2 — 48 tests: 36 golden-fixture parity (CP-001..017, 0 diff vs TS) + 3 BC_CHART_TYPES enum + 9 structural invariants (CP-015, XS-003)
│   ├── chart_helpers_test.dart                 # ADR-015 F2 — 34 unit tests for parseColors, parseChartData, buildTitleOption/Tooltip/Legend/Toolbox/DataZoom/Grid
│   ├── helpers/
│   │   └── fake_webview_platform.dart          # ADR-015 F3 — no-op WebViewPlatform stub for widget tests that instantiate flutter_echarts Echarts widget
│   ├── fixtures/
│   │   └── chart-options/
│   │       └── golden-fixtures.json            # ADR-015 F2 — 36 golden JSON fixtures generated from TS chart-option factory
│   ├── bc_integration_gate_test.dart           # Integration gate for Phase 5b.2 widgets and exports
│   ├── hooks_parity_test.dart                  # ADR-013 Phase 5b — 38 tests: PreviousValue, BcIntervalController, BcTimeoutController, BcScrollLock, BcLocalStorageStore (FakeAsync-driven timers)
│   ├── animation_adapt_parity_test.dart        # ADR-013 Phase 5b — 16 tests: resolveAnimation (default merge, kind→css, reduced-motion), resolveAdapt (breakpoint/platform/SSR), bcBreakpointQueries parity
│   ├── schema/
│   │   ├── ui_schema_test.dart                 # 55 tests — FetchParams, FetchResult, SortRule, BcSelection (sealed), enums, UiSchema/FieldNode/LayoutNode JSON round-trip, FieldResolver path resolution
│   │   └── field_resolver_test.dart            # 21 tests — field path resolution, binding evaluation, condition matching
│   ├── app_config_store_test.dart              # ADR-013 Phase 1/5a — AppConfigStore parity (defaults, setTheme resolves, subscribe, dispose, direction auto-resolve, dateLocale fallback)
│   ├── hooks_parity_test.dart                  # ADR-013 Phase 5b — 38 tests: PreviousValue, BcIntervalController, BcTimeoutController, BcScrollLock, BcLocalStorageStore (timer paths via FakeAsync)
│   └── animation_adapt_parity_test.dart        # ADR-013 Phase 5b — 16 tests: resolveAnimation (default merge, kind→css, reduced-motion collapse), resolveAdapt (breakpoint/platform/SSR), breakpoint-query parity
├── example/                                    # Demo Flutter app (BcUiDemoApp) — single-page showcase of all 32 widgets
│   ├── lib/
│   │   ├── main.dart                           # App entry, theme, drawer + section list, debugPrint for widget events
│   │   ├── sections/                           # 6 demo section files (fields, overlays, layout, data_display, views, misc)
│   │   └── shared/                             # Shared demo primitives (section card, heading, demo state)
│   ├── test/
│   │   └── widget_test.dart                    # Smoke test — BcUiDemoApp renders, drawer + 6 sections present
│   ├── web/                                    # Flutter web bootstrap (index.html, manifest.json)
│   └── pubspec.yaml                            # Path-dep on bc_ui_flutter, flutter_platform_widgets
└── tool/
    └── run-flutter.ps1                         # PATH-isolation wrapper for flutter/dart on machines with bloated PATH
```
**1581 Flutter tests** across 27 test files (incl. 107 chart tests: 48 factory golden-fixture parity + 34 helpers + 25 bc_chart widget contract; + 54 ADR-013 Phase 5b parity tests). `flutter analyze` returns 0 errors (300 info-level lints; invoke via `tool/run-flutter.ps1` on machines with a bloated PATH). `build_runner` succeeds. **Demo app** (`example/`) — `flutter analyze` 0 issues, `flutter build web` succeeds; 6 section files render all 32 widgets in a single-page showcase with Material/Cupertino adaptation via FPW.

**Token bridge** (`packages/ui/scripts/generate-flutter-tokens.mjs`): Parses `theme.css` HSL triplets, converts to hex at build time, emits `bc_theme.dart` with 25 light + 25 dark `const Color`, 4 factory methods (materialLight/Dark, cupertinoLight/Dark), and radius. `--check` mode for CI.

**Stub generator** (`packages/ui/scripts/generate-flutter-stubs.mjs`): Parses `.contract.ts` files for flutter adapter hints, extracts props/methods via brace-depth-aware TS interface parsing, maps TS types to Dart (primitives, arrays, unions, function types, `Partial<T>`, `BcChangeHandler<T>`), and emits StatefulWidget stubs with auto-generated header + TODO comments. Protects P0 hand-written files (no auto-generated header). `--check` mode for CI. Phase 5b.2 fixed extends clause walking and `Promise<T>` to `Future<T>` mapping, then promoted the 8 generated P1 stubs to completed widgets. 32 contracts scanned, 8 P1 widgets completed, 24 P0 protected, 98 skipped (no flutter hint). 80 Node.js tests across parser, type_mapping, widget_map suites.

---

## Tauri Native Shell (`packages/tauri/`)

Tauri 2.0 project that wraps Stencil components in a native WebView for desktop and mobile.

```
packages/tauri/
├── package.json                                # @bitcode/tauri — npm scripts for dev/build (desktop, android, ios)
├── .gitignore                                  # Excludes target/, gen/, node_modules/
│
└── src-tauri/
    ├── Cargo.toml                              # Tauri 2.10 + plugins (sql, fs, notification, barcode, biometric)
    ├── build.rs                                # Tauri build script
    ├── tauri.conf.json                         # frontendDist → ../../components/www, withGlobalTauri, CSP, window config
    ├── capabilities/
    │   └── default.json                        # Permissions: core, sql, fs, notification
    ├── icons/                                  # App icons for all platforms (generated via cargo tauri icon)
    └── src/
        └── main.rs                             # Entry point — plugin registration, SQLite migrations for _off_* tables
```

---

## go-json Language Engine (`packages/go-json/`)

Standalone JSON/JSONC programming language engine. Embeddable in Go applications. Powered by expr-lang/expr for expression evaluation.

**Design doc:** `docs/plans/2026-07-14-runtime-engine-phase-4.5a-go-json-core-language.md`
**Implementation plan:** `docs/plans/2026-07-14-runtime-engine-phase-4.5a-go-json-core-language-plan.md`

```
packages/go-json/
├── go.mod                                      # module github.com/bitcode-framework/go-json (Go 1.24+, expr-lang/expr)
├── go.sum
│
├── lang/                                       # Core language engine
│   ├── ast.go                                  # AST node types — 15 step types (let, set, if, switch, for, while, break, continue, return, call, try, error, log, comment) + Program, FuncDef
│   ├── parser.go                               # JSONC/JSON → AST. Handles all step types, overloaded nodes (return/error/log), ordered function params
│   ├── compiler.go                             # AST → CompiledProgram. Structural validation (break/continue outside loop), limit resolution
│   ├── vm.go                                   # Tree-walk interpreter. All step types, scope isolation, resource limits, debug hooks, trace
│   ├── scope.go                                # Variable scoping — block scope (if/for/while), chain lookup, isolated child (functions), ToMap() for expr-lang
│   ├── types.go                                # Gradual type system — InferType (JSON float64→int detection), TypesCompatible, nullable (?T)
│   ├── errors.go                               # GoJSONError with enrichment, levenshtein "did you mean?", JSON/Short/Error output, fluent builder
│   ├── expr_engine.go                          # ExprEngine interface + ExprLangEngine (expr-lang/expr). Compiled expression cache, error enrichment
│   ├── debugger.go                             # Debugger interface (OnStep/OnVariable/OnError/OnFunctionCall/OnFunctionReturn) + ExecutionTrace
│   ├── program.go                              # Immutable CompiledProgram, CompiledFunc, ParamDef, ResolvedLimits, DefaultLimits(), HardLimits()
│   ├── preprocess.go                           # JSONC → JSON: strip // and /* */ comments, trailing commas. String-aware (won't strip inside strings)
│   ├── preprocess_test.go                      # 18 tests — comments, strings, escaped quotes, trailing commas, edge cases
│   ├── integration_test.go                     # 35 tests — full pipeline (parse → compile → execute) for all features
│   └── edge_cases_test.go                      # 22 tests — nil, types, overflow, unicode, malformed JSON, error formatting
│
├── stdlib/                                     # Layer 2 stdlib (Layer 1 = expr-lang built-ins, ~68 functions)
│   ├── registry.go                             # Function registry + DefaultRegistry() with all Layer 2 functions
│   ├── math.go                                 # 7 functions: clamp, sign, randomInt, randomFloat, pow, sqrt, mod
│   ├── strings.go                              # 5 functions: padLeft, padRight, substring, format, matches (regex)
│   ├── arrays.go                               # 5 functions: append, prepend, slice, chunk, zip
│   └── types.go                                # 2 functions: bool (truthy conversion), isNil
│
├── runtime/                                    # Runtime API
│   ├── runtime.go                              # NewRuntime(opts...), Compile(), Execute(), ExecuteJSON(). Program cache (SHA256 keyed). Concurrent-safe
│   ├── limits.go                               # Limits struct, DefaultLimits(), HardLimits(), ToResolved()
│   ├── context.go                              # Session (UserID, Locale, TenantID, Groups) + ExecutionMeta (ID, Program, StartedAt)
│   └── logger.go                               # Logger interface (Log(level, message, data)) + DefaultLogger (stdout)
│
├── cmd/go-json/
│   └── main.go                                 # CLI placeholder
│
├── codegen/                                    # Reserved for Phase 4.5b (struct codegen)
├── io/                                         # Reserved for Phase 4.5c (I/O modules)
│
└── testdata/                                   # Test fixture programs
    ├── hello.json                              # Minimal program
    ├── hello.jsonc                              # Same with JSONC comments
    ├── variables.json                           # let, set, value/expr/with modes
    ├── variables.jsonc                          # Same with JSONC comments + trailing commas
    ├── control_flow.json                        # if/elif/else, switch
    ├── loops.json                               # for, while, range, break, continue
    ├── functions.json                           # Function definition, call, recursion (factorial)
    ├── recursion.json                           # Factorial + fibonacci recursive functions
    ├── error_handling.json                      # try/catch/finally, error throw
    ├── stdlib_math.json                         # Layer 2 math: clamp, sign, pow, sqrt, mod
    ├── stdlib_strings.json                      # Layer 2 strings: padLeft, substring, format, matches
    ├── stdlib_arrays.json                       # Layer 2 arrays: append, prepend, slice, chunk, zip
    ├── limits.json                              # Program-level resource limits
    ├── comments.json                            # _c inline, standalone, multi-line
    └── types.json                               # Type inference, nullable, any, strict assignment
```

**82 tests** across 4 packages. `go build ./...` + `go vet ./...` clean.

---

## Samples (`samples/`)

### ERP Sample (`samples/erp/`)

Full ERP application demonstrating all engine features.

```
samples/erp/
├── bitcode.yaml                                # Project config (port 8989, SQLite)
├── bitcode.toml                                # Alternative TOML config
├── .env.example                                # Environment variable template
├── README.md                                   # Comprehensive feature documentation
├── run.bat / run.ps1 / run.sh                  # Cross-platform run scripts (go install + bitcode dev)
│
└── modules/
    ├── base/                                   # Core module (users, roles, groups, permissions)
    │   ├── module.json
    │   ├── models/                             # 7 models: user, role, group, permission, record_rule, audit_log, setting
    │   ├── apis/                               # 8 APIs: auth, user, group, role, permission, record_rule, setting, audit_log
    │   ├── views/                              # CRUD views for all models
    │   ├── data/                               # Default roles, groups, users
    │   └── templates/                          # Layout, partials, view templates
    │
    ├── crm/                                    # CRM Module (depends: base)
    │   ├── module.json                         # 9 permissions, 2 groups (user, manager), menu, settings
    │   ├── models/                             # contact, lead, activity, tag, vip_contact (inherits contact)
    │   ├── apis/                               # contact_api, lead_api, activity_api, tag_api
    │   ├── processes/                          # lead_workflow, qualify_lead, convert_lead, win_lead, lose_lead,
    │   │                                       #   log_activity, enrich_leads (DAG), enrich_lead_dag
    │   ├── agents/                             # lead_agent (3 event triggers + 2 cron jobs)
    │   ├── views/                              # list, form, kanban, dashboard, pipeline_chart
    │   ├── scripts/                            # TS: on_deal_won, on_deal_lost, notify_manager, weekly_report, stale_leads
    │   │                                       # Python: analyze_pipeline.py
    │   ├── templates/                          # dashboard.html, partials/lead_card.html
    │   ├── data/demo.json                      # Demo contacts, leads, tags
    │   └── i18n/id.json                        # Indonesian translations
    │
    └── hrm/                                    # HRM Module (depends: base)
        ├── module.json                         # 13 permissions, 3 groups (user, officer, manager), menu
        ├── models/                             # department, job_position, employee, leave_request
        ├── apis/                               # department_api, position_api, employee_api, leave_api
        ├── processes/                          # employee_workflow, leave_workflow, submit_leave, approve_leave,
        │                                       #   reject_leave, promote_employee, onboard_employee
        ├── agents/                             # hr_agent (3 event triggers + 1 cron)
        ├── views/                              # list, form, calendar, dashboard
        ├── scripts/                            # TS: on_leave_approved, on_leave_submitted, on_promotion, weekly_attendance
        │                                       # Python: calculate_leave_balance.py, generate_onboard_checklist.py
        ├── templates/                          # dashboard.html, partials/employee_card.html
        ├── data/demo.json                      # Demo departments, employees, positions
        └── i18n/id.json                        # Indonesian translations
```

---

## Test Coverage

| Package | Tests | What's Tested |
|---------|-------|---------------|
| `pkg/ddd` | 3 | Entity, Aggregate events, DomainEvent |
| `pkg/security` | 5 | Password hash/check, JWT generate/validate/wrong-secret |
| `compiler/parser` | 27 | Model (10), API (8), View (6), Workflow (3) |
| `domain/model` | 6 | Registry CRUD, module prefix, table name |
| `domain/security` | 9 | User lifecycle, role inheritance, groups, record rules |
| `domain/event` | 4 | Pub/sub, subscribe all, no subscribers |
| `domain/setting` | 5 | Get/set, defaults, all |
| `infrastructure/cache` | 5 | Memory cache set/get, TTL, delete, clear |
| `infrastructure/module` | 7+ | Registry, dependency resolution, circular detection, FS |
| `infrastructure/i18n` | 4 | Translate, fallback, locale detection |
| `presentation/template` | 5 | Load/render, helpers (truncate, dict, eq) |
| `presentation/view` | — | Component compiler tests |
| `runtime/executor` | 3+ | Step dispatch, unknown type, step error, DAG |
| `runtime/executor/steps` | 9 | Validate (eq/fail/required), emit, assign, if, process parse |
| `runtime/validation` | 28 | Field validators (required, email, regex, conditional, etc.), sanitizers, model-level validators, update merge |
| `runtime/hook` | 10 | Event dispatcher (priority, condition, sync/async, retry, on_change cascade), context, expression eval |
| `embedded` | — | Embed FS tests |
| **Total** | **131+** | |

Run all tests: `cd engine && go test ./... -v`

---

## Key Interfaces

```go
// Cache — pluggable caching (memory or redis)
type Cache interface {
    Get(key string) (any, bool)
    Set(key string, value any, ttl time.Duration)
    Delete(key string)
    Clear()
}

// StepHandler — process step execution
type StepHandler interface {
    Execute(ctx context.Context, execCtx *Context, step StepDefinition) error
}

// ScriptRunner — plugin execution (TS/Python)
type ScriptRunner interface {
    Run(ctx context.Context, script string, params map[string]any) (any, error)
}

// PermissionChecker — RBAC authorization
type PermissionChecker interface {
    UserHasPermission(userID string, permission string) (bool, error)
}

// RecordRuleEngine — row-level security
type RecordRuleEngine interface {
    GetFilters(userID string, modelName string, operation string) ([][]any, error)
}

// Repository[T] — generic data access (DDD)
type Repository[T any] interface {
    Save(ctx context.Context, entity *T) error
    FindByID(ctx context.Context, id string) (*T, error)
    FindAll(ctx context.Context, filters map[string]interface{}, page int, pageSize int) ([]T, int64, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id string) error
}

// ModuleFS — filesystem abstraction for module loading
type ModuleFS interface {
    ReadFile(name string) ([]byte, error)
    ReadDir(name string) ([]os.DirEntry, error)
    Glob(pattern string) ([]string, error)
}
```

---

## Build & Run

```bash
# Engine
cd engine
make build          # Build bitcode binary → bin/bitcode
make install        # Install bitcode to $GOPATH/bin
make serve          # Start production server
make dev            # Start dev server with hot reload
make test           # Run all tests

# Components
cd packages/components
npm install
npm run build       # Build Stencil components → dist/
npm run start       # Dev server with watch

# Sample ERP
cd samples/erp
./run.sh            # Linux/macOS
.\run.bat           # Windows CMD
.\run.ps1           # Windows PowerShell
```

---

## Dependencies

### Engine (Go)

| Dependency | Purpose |
|------------|---------|
| `gofiber/fiber/v2` | HTTP framework |
| `gorm.io/gorm` | ORM (SQLite, Postgres, MySQL) |
| `glebarez/sqlite` | Pure-Go SQLite driver |
| `gorm.io/driver/postgres` | PostgreSQL driver |
| `gorm.io/driver/mysql` | MySQL driver |
| `spf13/viper` | Configuration (env + file) |
| `spf13/cobra` | CLI framework |
| `golang-jwt/jwt/v5` | JWT tokens |
| `google/uuid` | UUID generation |
| `redis/go-redis/v9` | Redis client |
| `golang.org/x/crypto` | bcrypt password hashing |
| `gofiber/contrib/websocket` | WebSocket support |
| `fsnotify/fsnotify` | File watching (dev mode) |

### Components (TypeScript)

| Dependency | Purpose |
|------------|---------|
| `@stencil/core` | Web Component compiler |
| `echarts` | Charts |
| `@tiptap/*` | Rich text editor |
| `@codemirror/*` | Code editor |
| `@fullcalendar/*` | Calendar view |
| `frappe-gantt` | Gantt chart |
| `leaflet` | Maps |
| `sortablejs` | Drag & drop |
| `signature_pad` | Signature capture |
| `jsbarcode` + `qrcode` | Barcode/QR generation |
| `markdown-it` | Markdown rendering |
| `xlsx` | Excel export |
