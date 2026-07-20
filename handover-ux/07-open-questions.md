# 07 · Open questions

Every entry here needs Taḋg's answer before the affected work can start. Do not resolve them yourself; do not implement around them.

When Taḋg answers, move the entry from this file to `02-decisions.md` under **Decided**, and record the answer verbatim so future agents can trace it.

## Q1 · New repo location and name

Taḋg said 2026-07-20: "i take your rec on separate project and i've created one." Path and name are not recorded.

**Needed:** absolute path on this Mac, repo name, GitHub owner/name, whether it has any initial contents.

## Q2 · Confirm the R-series recommendations from `02-decisions.md`

Each is currently only **Recommended**:

- **R1** — MVP scope: all three screens, no Preferences.
- **R2** — Distribution outside Mac App Store, sandbox off.
- **R3** — Discover `folio` on PATH; do not bundle.
- **R4** — Add `folio letter --list` to the CLI; use it for recipient discovery in the Letter screen.
- **R5** — Sidebar as primary layout; no landing page.

**Needed:** yes/no on each, or a counter-proposal.

## Q3 · Deployment target (O4 in `02-decisions.md`)

Suggested: macOS 14 (Sonoma).

**Needed:** minimum macOS version to support.

## Q4 · Bundle identifier and app name (O3)

Suggested: `ie.tigger.first-folio-mac`, product name `First Folio`.

**Needed:** confirmed reverse-DNS identifier and product name to write into Info.plist.

## Q5 · Tool-not-found guidance (O2)

If `folio` exits with a "typst not found" or "pandoc not found" style error, does the app:

- (a) relay the CLI's stderr verbatim — recommended.
- (b) recognise the error and offer install guidance (e.g. `brew install typst`) inline in the app.

**Needed:** (a) or (b), with a note if guidance should be a separate follow-up rather than never.

## Q6 · Fixture policy

Copy fixtures from `/Users/tigger/code/tigoss/first-folio/examples/` into the new repo, or reference by absolute path in tests?

Recommendation: **copy** (avoids the new repo depending on the old repo's on-disk state), with a `NOTE.md` recording source commit and copy date. Do not vendor any file that may contain personal correspondence — `about-time.org` should be checked before vendoring.

**Needed:** copy vs reference, and permission to copy `about-time.org` (or a smaller synthetic fixture instead).

## Q7 · CLI feature request cadence

If the app needs a new CLI capability (e.g. R4's `folio letter --list`, or the post-MVP schema/validate flags described in `docs/ux_proposal.md` §Open questions 6):

- Does Taḋg want the app-side agent to file the issue in `tigger-developer/first-folio`, wait for the CLI change to ship, then continue?
- Or does the app-side agent implement the CLI change (Go) itself, submit a PR, then continue?

Recommendation: **file the issue, wait, continue** — keeps the app agent's toolchain scope small (Swift only). CLI changes are governed by the first-folio SDLC and would still need Gate 1/2 there.

**Needed:** confirmation of the cadence.

## Q8 · Icon and branding

The ux_proposal explicitly does not commit to a colour scheme or iconography. MVP presumably needs at least a placeholder app icon.

**Needed:** placeholder icon acceptable (e.g. a text-only "FF" glyph), or does Taḋg have an intended asset?

## Q9 · Localization scope

`~/.claude/CLAUDE.md` requires Hiberno-English with OED spellings for Taḋg's prose. The app's user-visible strings should follow the same rule.

**Needed:** confirm the app ships English-only (Hiberno-English / OED) for MVP, with no other locales.

## Q10 · Test data privacy

If Taḋg's `about-time.org` (in `/Users/tigger/code/tigoss/first-folio/`) is the intended manual-test fixture, it should not be committed to a public repo. If the new repo is public on GitHub, this becomes an active concern.

**Needed:** is the new repo public? If yes, confirm what may and may not be vendored.
