# Multi-Brain Master Index

> **What this is:** Lightweight, per-project cross-session memory for AI agents.
> Each bucket = one active workstream. Agents write short entries after meaningful
> work so the next session picks up where the last one left off — without re-reading
> every handoff.
>
> **Relationship to other memory systems in this repo:**
> - `handoffs/` — detailed per-session handoff docs (long form, human + AI context). Multi-brain buckets *reference* these, they don't replace them.
> - `.opencode-jce/context/` — heavier JCE context-keeper with evidence/telemetry JSON. Multi-brain is the lighter, faster sibling.
> - `.opencode-context.md` — top-level always-loaded context. Keep it ≤40 lines.
>
> **Rule:** If a bucket file grows past ~150 lines, archive old entries to `archive/<bucket>-YYYY-MM.md` and keep only the last ~10 entries in the active bucket file.

---

## Active Buckets

| Bucket | Workstream | Last updated | Index |
|--------|-----------|--------------|-------|
| `ui-library` | `packages/ui` Mitosis components, ADR-008 contract freeze, ADR-013 stub completion | 2026-10-29 | [buckets/ui-library.md](buckets/ui-library.md) |
| `engine` | Go engine runtime, JSON definitions, modules, query builder, security | 2026-07-13 | [buckets/engine.md](buckets/engine.md) |
| `flutter` | Flutter adapter (`packages/ui-flutter`), token bridge, demo parity, ADR-011/015 | 2026-10-29 | [buckets/flutter.md](buckets/flutter.md) |
| `offline-mode` | Tauri shell, sync engine, bc-native bridge, Phases 3–5 | 2026-06-30 | [buckets/offline-mode.md](buckets/offline-mode.md) |
| `adr` | Architecture Decision Records (active frontier ADR-021 modular definition surface draft) | 2026-08-02 | [buckets/adr.md](buckets/adr.md) |
| `demo-app` | React demo (`packages/ui/demo-app`) + Flutter demo parity, QA bug sweeps | 2026-10-29 | [buckets/demo-app.md](buckets/demo-app.md) |

## How to Use (Agents)

1. **Session start:** Read `.multibrain/session.md` first. Pick the bucket(s) matching
   the current task. Open only those bucket files — don't read all of them.
2. **Before working:** Scan the bucket's "Current State" + "Open Threads" + "Gotchas"
   sections so you don't repeat finished work or break a known constraint.
3. **After meaningful work** (feature done, bug fixed, decision made, blocker hit):
   append ONE short entry to the matching bucket. Update "Current State" if it changed.
4. **If a bucket doesn't exist** for your workstream, create it via the pattern in
   [buckets/README.md](buckets/README.md) and add a row to this index.
5. **Archiving:** When a bucket exceeds ~150 lines, move the oldest entries to
   `archive/<bucket>-YYYY-MM.md` (create the file if needed) and leave the recent ones.

## Entry Format (append to top of bucket, newest first)

```
## YYYY-MM-DD — <short title>

**Status:** done | in-progress | blocked | decision
**Files touched:** <paths or "see handoff YYYY-MM-DD-*.md">
**Evidence:** <command run + result, or commit SHA>

<2–8 lines: what changed, why, what's next. Link to handoff if long.>
```

Keep entries short. Detail belongs in `handoffs/`. Multi-brain is the *index* to find the right handoff fast.
