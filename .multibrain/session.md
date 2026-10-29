# Multi-Brain Session Memory

> Per-project lightweight cross-session memory for AI agents.
> Active whenever this file exists. Read this FIRST at session start.

**Last session:** 2026-10-29 (Droid — go-json examples + cross-platform installers)
**Last bucket touched:** engine

---

## How to use (read this before working)

1. **Session start:** Read this file, then open `INDEX.md`. Pick the 1-2 buckets
   relevant to the current task. Do NOT open all buckets â€” only what matches.
2. **During work:** When you make a meaningful decision, finish a phase, hit a
   blocker, or learn a gotcha, write ONE short entry to the matching bucket.
   Prepend it to the top (newest first). Keep entries under ~15 lines.
3. **After entry:** Update `INDEX.md` "Last updated" for that bucket + this
   file's "Last session" / "Last bucket touched" lines.
4. **Archive, don't delete:** When a bucket grows past ~50 entries or a milestone
   fully closes, move old entries to `archive/<bucket>-YYYY-MM.md`.
5. **Non-destructive:** Never overwrite another agent's entry without reason.
   Append/prepend your own. If you disagree, add a correction note with your
   name/date.

## Entry format

```
### YYYY-MM-DD â€” <short title>
- **What:** <one line>
- **Decision / outcome:** <one or two lines>
- **Files touched:** <paths>
- **Gotcha / note:** <optional, the valuable part>
```

Sign entries with your agent identity if ambiguous (e.g. `- @worker`).

---

## Current active frontiers (2026-07-13)

- **adr** â€” ADR-021 Modular Definition Surface Architecture is the active draft. Locked: components (not blocks); JSON UI canonical / HTML secondary; process vs workflow; BitCode-native process JSON; local `name` / cross-module `module::name` logical refs; optional bindings; shell module shared defs; scripts as executable assets; local `scripts/`-relative paths; cross-module scripts via explicit exports; three separate resolvers; `includes`/`excludes` constrained POSIX globs with exclusion precedence; `data/` only for current static seed/demo/reference data (`fixtures/` deferred); repository-backed resources allowed in Phase A with full CRUD behind the `resource_provider` contract (typed schema, stable identity, explicit ops, deterministic collections, idempotent create, concurrency-aware update, declared delete semantics, write audit, normalized errors, no raw responses, no offline sync). Component packaging locked: module-owned definitions via `module::name`; packages/ui* are renderer primitive libraries only (no package registry, no `core::`). Websites/pages split locked: `websites/` public-by-default UI folder, `pages/` internal-by-default, per-file auth overrides, API auth always in `apis/`, public static files in `public/`, mount config deferred. HTML renderer files locked under `templates/` with explicit JSON references. Asset policy locked: private `assets/`, public `assets/public/**`, runtime `storage/uploads/`. Automation boundary locked: `automations/` trigger adapter, `processes/` business execution, `agents/` AI-only. ADR-022 created as living design discussion log (master `docs/adr/022-modular-surface-design-log.md` + detail files `docs/adr/022/001-006`); locked decisions mirror into ADR-021.
- **ui-library** â€” ADR-013 cross-tier infra complete (2026-06-29). 28 stub
  components ported. `packages/ui` 1538 tests pass, 4-target build green.
  Remaining: Phase 3-5 of ADR-008 (a11y/perf testing, docs, Flutter integration).
- **flutter** â€” `packages/ui-flutter` has 120+ components ported. Demo parity
  gap: React demo has 127 components / 14 groups vs Flutter ~25 / 6 sections.
- **demo-app** â€” Bug-fix sweeps ongoing. Recent: Mitosis ref compiler bug,
  overlay portal bug, BcTabs/BcSplitter runtime bugs fixed.
- **engine** â€” Stable. 570 Go tests pass. Open: Redis cache wiring, NATS bus,
  Marketplace.
- **offline-mode** â€” Phase 1-2.5 complete. Phase 3-5 pending (sync engine,
  conflict resolution, polish).

## Buckets

| Bucket | Covers | Last updated |
|--------|--------|--------------|
| [ui-library](buckets/ui-library.md) | `packages/ui`, Mitosis, contracts, ADR-008 freeze | 2026-06-30 |
| [engine](buckets/engine.md) | Go engine, runtime, JSON defs, modules | 2026-07-13 |
| [flutter](buckets/flutter.md) | `packages/ui-flutter`, token bridge, FPW | 2026-06-30 |
| [offline-mode](buckets/offline-mode.md) | Tauri, bc-native, sync engine | 2026-06-30 |
| [adr](buckets/adr.md) | Architecture Decision Records | 2026-07-13 |
| [demo-app](buckets/demo-app.md) | Demo app parity + bug-fix sweeps | 2026-06-30 |

See [INDEX.md](INDEX.md) for the full dependency map and bucket-to-file ownership.
