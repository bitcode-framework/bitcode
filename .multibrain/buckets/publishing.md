# Publishing & Org Profile Bucket

Publishing pipeline for the `bitcode-framework` GitHub org: the org profile
README, per-repo metadata, and the live demos linked from both.

## Current State

Every public repo except `ui-tauri` has a working browser demo on GitHub Pages.
`ui-tauri` is a native desktop/mobile shell with no honest web equivalent.

| Repo | Demo | Status |
|---|---|---|
| `ui` | https://bitcode-framework.github.io/ui/ | 200, renders (canvas + DOM) |
| `ui-flutter` | https://bitcode-framework.github.io/ui-flutter/ | 200, canvas only |
| `ui-schema` | https://bitcode-framework.github.io/ui-schema/ | 200, 4 scenarios verified |
| `go-json` | https://bitcode-framework.github.io/go-json/ | 200, 8 samples verified |
| `go-json-runtimes` | https://bitcode-framework.github.io/go-json-runtimes/ | 200, 4 samples verified |
| `expr-js` | https://bitcode-framework.github.io/expr-js/ | 200, 12 samples verified |
| `ui-web-components` | https://bitcode-framework.github.io/ui-web-components/demo/ | 200, deprecated but working; **note the `/demo/` subpath** |

## How publishing works

There is no publish script. Every repo is a manual subtree split or a manual
clone-and-copy. Always `git push --force <remote> "${sha}:refs/heads/main"`
when pushing a split.

Pages must be enabled per repo before the first deploy, or `deploy-pages`
fails with a bare `404 Not Found` that says nothing about the cause:

```bash
gh api -X POST repos/bitcode-framework/<repo>/pages -f build_type=workflow
```

## Gotchas

- **Never re-point a pushed tag.** `sum.golang.org` is append-only. Learned the
  hard way in the go-json task.
- **`go-json` and `go-json-runtimes` tag independently.** Latest go-json is
  v0.1.2; latest runtimes is v0.1.4. An install command naming the wrong tag
  fails, so check `gh api repos/.../tags` before writing one into docs.
- **Component counts drift.** The profile README said "119 Stencil components"
  and the `ui` repo description said 204; the actual Mitosis source count is
  194. Prefer no number over a stale number.
- **`ui-web-components` is deprecated** and superseded by `ui`, `ui-schema`,
  and `ui-flutter`. Its description is prefixed `DEPRECATED` but it is not
  archived, because archiving hides the migration note.
- **`ui-web-components` serves its gallery from `/demo/`, not the site root.**
  The root 404s even though the deploy succeeds and `www/demo/` exists. Probing
  subpaths found it. Do not "fix" this by linking the root.
- **A successful `deploy-pages` run does not prove the URL works.** The
  ui-web-components deploy reported success repeatedly while the root URL kept
  404ing; the artifact was fine and the path was wrong. Always fetch the actual
  URL, don't trust the green check.
- **`lynx-stack` is an upstream fork** (lynx-family/lynx-stack) and must not be
  presented as ours.
- **The `erp` repo is private** and deliberately absent from the profile README.
  The word "ERP" still appears there as an application type, which is intended.
- **Do not state repo visibility in the profile README.** Visibility is a repo
  setting the owner controls; asserting "(private)" goes stale the moment it
  changes. Removed 2026-10-30 at the owner's request.
- **`bitcode` is private.** Any `go install` or clone command pointing at it is
  a dead link for visitors. Install commands must target `go-json` or
  `go-json-runtimes`.
- **Do not state a license on the org profile.** Licensing is per repo, and
  `bitcode` has no LICENSE yet. The org page used to assert MIT for everything;
  that was wrong and was removed 2026-10-30.
- **`packages/components` is NOT dead code.** It looks deprecated, and its own
  README advertises 127 Stencil components, but it is what the engine's
  `bc-*` markup and `tauri.conf.json` (`beforeDevCommand`) actually build
  against. `packages/ui` emits a *different* tag set (`bc-input` vs
  `bc-field-string`), so pointing the engine or Tauri at `packages/ui` is a
  migration, not a rename. Verify before "fixing" any doc that says Stencil.
- **The `bitcode` repo on GitHub is empty.** `bitcode-framework/bitcode` has no
  commits, so it has no LICENSE and no README. Pushing the monorepo there is
  still outstanding.

## 2026-10-30 — org profile corrected, three demos built, one credential leak found

- Rewrote `bitcode-framework/.github` → `profile/README.md`. Removed the
  non-existent `pg_fast_data_transfer` / `pg_data_transfer`, removed `erp`,
  marked `ui-web-components` deprecated, added `ui`, `ui-schema`, `ui-flutter`,
  `expr-js`, and pointed Quick Start at public packages instead of the private
  monorepo.
- Built and published three demos: `ui-schema` (reactive engine, 4 scenarios),
  `go-json-runtimes` (`script:` bridge, 4 samples), plus the profile links.
- Updated descriptions and homepages for `ui`, `ui-schema`, `ui-flutter`,
  `expr-js`, `go-json`, `go-json-runtimes`, and `bitcode`.
- **Found a real Stripe live key** (`sk_live_…`, 44 chars) committed in
  `packages/ui-schema/test/block-catalog/security-policy.test.ts:65`. GitHub
  push protection caught it. Replaced with a placeholder in commit `a899e196`.
  The key is still in monorepo history and **must be rotated** — the monorepo is
  private, and the public `ui-schema` repo never contained it, so exposure was
  limited to org members.

## Open Threads

- **Push the monorepo to `bitcode-framework/bitcode`.** That repo is empty. Two
  blockers found 2026-10-30 and not yet resolved:
  1. The Stripe live key is still reachable in history (`465826e1`, `71885aa3`).
     Pushing publishes it. Rotate the key and rewrite those commits first, or
     accept that the org repo carries a dead credential.
  2. `packages/expr-js` is an uninitialized submodule whose gitlink
     (`257c4fab`) does not exist in the public `expr-js` remote, so a fresh
     clone of the org repo cannot check it out. Re-point the gitlink at a real
     commit first.
- Rotate the leaked Stripe key. Removing it from the working tree is not
  enough; history still carries it.
- Decide whether to archive `ui-web-components` or leave it as a redirect
  surface. It is currently live and linked, so archiving would break the link.
- `ui-tauri` has no demo. Options: omit, link release binaries, or add a
  screenshot gallery.
- The profile README links `github.com/bitcode-framework/bitcode`, which 404s
  for anonymous visitors while that repo is private. The link is correct and
  will work once visibility changes; do not remove it.
- Migrating the engine and Tauri shell from `packages/components` to
  `packages/ui` is open work. The tag names differ, so it is not a drop-in swap.
