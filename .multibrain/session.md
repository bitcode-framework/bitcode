# Multi-Brain Session Memory

> Per-project lightweight cross-session memory for AI agents.
> Active whenever this file exists. Read this FIRST at session start.

**Last session:** 2026-07-01 (Agent cleanup and verification sweep)
**Last bucket touched:** ui-library (ui-schema strict tsc cleanup after Agent I)

---

## How to use (read this before working)

1. **Session start:** Read this file, then open `INDEX.md`. Pick the 1-2 buckets
   relevant to the current task. Do NOT open all buckets — only what matches.
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
### YYYY-MM-DD — <short title>
- **What:** <one line>
- **Decision / outcome:** <one or two lines>
- **Files touched:** <paths>
- **Gotcha / note:** <optional, the valuable part>
```

Sign entries with your agent identity if ambiguous (e.g. `- @worker`).

---

## Current active frontiers (2026-06-30)

- **adr** — ADR-017 (Block Catalog Architecture) + ADR-018 (Data/Action/Reactivity
  Contract) are **Accepted** with matrix 0.1 baseline. Next: implementation plan
  appendix, then code phase for Flutter/web/export gaps and foundational runtime work.
- **ui-library** — ADR-013 cross-tier infra complete (2026-06-29). 28 stub
  components ported. `packages/ui` 1538 tests pass, 4-target build green.
  Remaining: Phase 3-5 of ADR-008 (a11y/perf testing, docs, Flutter integration).
- **flutter** — `packages/ui-flutter` has 120+ components ported. Demo parity
  gap: React demo has 127 components / 14 groups vs Flutter ~25 / 6 sections.
- **demo-app** — Bug-fix sweeps ongoing. Recent: Mitosis ref compiler bug,
  overlay portal bug, BcTabs/BcSplitter runtime bugs fixed.
- **engine** — Stable. 570 Go tests pass. Open: Redis cache wiring, NATS bus,
  Marketplace.
- **offline-mode** — Phase 1-2.5 complete. Phase 3-5 pending (sync engine,
  conflict resolution, polish).

## Buckets

| Bucket | Covers | Last updated |
|--------|--------|--------------|
| [ui-library](buckets/ui-library.md) | `packages/ui`, Mitosis, contracts, ADR-008 freeze | 2026-06-30 |
| [engine](buckets/engine.md) | Go engine, runtime, JSON defs, modules | 2026-06-30 |
| [flutter](buckets/flutter.md) | `packages/ui-flutter`, token bridge, FPW | 2026-06-30 |
| [offline-mode](buckets/offline-mode.md) | Tauri, bc-native, sync engine | 2026-06-30 |
| [adr](buckets/adr.md) | Architecture Decision Records | 2026-06-30 |
| [demo-app](buckets/demo-app.md) | Demo app parity + bug-fix sweeps | 2026-06-30 |

See [INDEX.md](INDEX.md) for the full dependency map and bucket-to-file ownership.
