# 02 · Decisions

Every decision below has a status: **Decided** (act on it), **Recommended** (proposed by Claude 2026-07-20, awaiting Taḋg's confirmation), or **Open** (needs Taḋg before you touch the affected code).

## Decided

### D1 · Separate repository

The app lives in its own repository, not in the `first-folio` Go repo. Confirmed by Taḋg 2026-07-20.

### D2 · The CLI is authoritative

The app never adds capability the `folio` CLI cannot perform. When the app appears to add something (recipient list surfaced from an Org file, PDF auto-open, YAML editing when Preferences ships), it is either presenting an existing CLI capability, or performing an OS-level action (Finder reveal, Preview open) that is not `folio`'s job.

Source: `docs/ux_proposal.md` §Core principle.

### D3 · Preferences pane is post-MVP

MVP treats `~/.config/first-folio/script.yaml` as read-only. The app reads it at launch to prefill defaults; it never writes.

Source: `docs/ux_proposal.md` §Scope and non-goals, §Preferences (post-MVP).

## Recommended (awaiting Taḋg's confirmation — do not implement affected behaviour until confirmed)

### R1 · MVP scope: all three screens, no Preferences

Ship a working macOS app with Convert, Letter, and Manuscript screens as specified in `docs/ux_proposal.md` §Screen: Convert / Letter / Manuscript. No live PDF preview surface (auto-open in macOS Preview covers "see the result immediately"). No Preferences pane.

Rationale: matches the Vision doc's stated bar ("one click"), and gives a shippable app without waiting on the CLI additions the Preferences pane would require.

### R2 · Distribution: direct download, notarized, sandbox off

Do not target the Mac App Store for MVP. Sandbox is disabled. The app runs with the user's normal filesystem privileges, spawns `folio` as a child process, reads `~/.config/first-folio/script.yaml` without prompts, and discovers `folio` on the user's `PATH`.

Rationale: the App Sandbox conflicts with the core architecture in three concrete ways — subprocess execution of arbitrary binaries, reads outside the container, and `PATH` lookups across the sandbox boundary — each documented in `docs/ux_proposal.md` §Open questions 5.

Consequence: for MVP you produce an unsigned `.app` runnable locally by the developer. Codesigning, notarization, and distribution via Homebrew Cask or direct download are separate follow-up work, tracked as their own issues later.

### R3 · `folio` binary discovery: PATH only

The app does not bundle the `folio` binary. It discovers `folio` in this order:

1. `~/.local/bin/folio`
2. Homebrew paths (`/opt/homebrew/bin/folio`, `/usr/local/bin/folio`)
3. `PATH`
4. If none found, show a first-run sheet with install instructions and a `Choose folio binary...` picker that persists the chosen path in `UserDefaults`.

Rationale: keeps releases simple, avoids version-drift between bundled `folio` and the user's shell `folio`.

Consequence: same rule for `typst` and `pandoc`. The app does not discover them itself; if `folio` fails because `typst` is missing, relay the CLI's stderr as-is. Do not offer install guidance for third-party tools — it goes stale.

### R4 · UX Open Q1 — Recipient listing via `folio letter --list`

Add a new subcommand mode to the CLI: `folio letter <source.org> --list` emits the recipient names (one per line, or JSON if we want structured tags) and exits without generating any PDFs. The app uses this to populate the recipient table.

Rationale: single source of truth for what a "recipient" is; keeps the app's parser surface at zero; matches the D2 principle.

Consequence: this is a **new CLI feature that must ship before the app's Letter screen can be built**. Open an issue against `tigger-developer/first-folio` describing the flag, its output format (recommend: one recipient name per line, UTF-8, no ANSI, exit 0 on success, exit non-zero on parse failure with a diagnostic on stderr), and any tag surfacing rules. See `07-open-questions.md`.

### R5 · Landing page vs sidebar — sidebar as proposed

No landing page. Fixed three-item sidebar (Convert / Letter / Manuscript) as the primary layout. First-run guidance, if needed, is a sheet inside the currently-active view.

Rationale: `docs/ux_proposal.md` §Layout already argues this convincingly and matches macOS Mail / Notes / Music conventions.

## Open (Taḋg must decide before implementation begins)

### O1 · New repo location and name

Taḋg has created the repo; the path is not yet recorded here. See `01-project.md`.

### O2 · UX Open Q4 — Tool-not-found guidance

If `folio` exits with `typst not found` or `pandoc not found`, does the app:

- (a) relay the CLI's stderr verbatim; or
- (b) recognize the error and offer install guidance (e.g. `brew install typst`)?

Recommendation: **(a)**. Guidance goes stale, and the CLI's own error is already actionable.

### O3 · Bundle identifier and app name

Suggested: `ie.tigger.first-folio-mac` or similar; product name `First Folio`. Not decided.

### O4 · Deployment target

Minimum macOS version to support. Suggested: **macOS 14 (Sonoma)** — SwiftUI `NavigationSplitView` is well-behaved from 13 and settled by 14; anything older adds workarounds.

## Explicitly out of scope for MVP

Do not build these unless Taḋg re-scopes:

- In-app text editing of source files.
- Live PDF preview surface (a Preview.app auto-open is fine; embedding `PDFKit` for a preview pane is out).
- Writing to `script.yaml` (the whole Preferences pane in `docs/ux_proposal.md` §Preferences).
- Codesigning, notarization, Sparkle-style auto-update, Homebrew Cask packaging.
- iOS/iPadOS/visionOS variants.
- Cloud sync, collaboration, telemetry.
- Any Swift-side parser more than the shallowest sniff (and even that only if R4 is rejected).
