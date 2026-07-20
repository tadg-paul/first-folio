<!-- Prepared: 2026-07-20 · Prepared by: Claude Opus 4.7 · For: agent picking up the SwiftUI companion app -->

# Handover: First Folio SwiftUI companion app

You are picking up a new project — a native macOS SwiftUI companion app that wraps the existing `folio` CLI. The CLI is authoritative; the app is a window onto it.

**Nothing in this project has been built yet.** Your job is to bootstrap the Xcode project, implement the MVP, and document as you go.

## Read in this order

1. `01-project.md` — what you are building and why. Do not skip. It also names the paths you must reference outside this repo.
2. `02-decisions.md` — what is decided and what is still open. Note the "Pending decisions" section: those are Taḋg-only calls.
3. `03-architecture.md` — the runtime shape (SwiftUI shell + Swift wrapper + `folio` subprocess).
4. `04-cli-contract.md` — every `folio` subcommand's argv shape, output pattern, and the two known deltas from the UX proposal you must not repeat.
5. `05-process.md` — the SDLC that governs this work, gates, canary, and how to hand back to Taḋg.
6. `06-first-steps.md` — concrete boot: create the Xcode project, discover `folio`, first end-to-end smoke.
7. `07-open-questions.md` — decisions you must escalate before implementing certain parts.

## Reference material outside this repo

Per `~/.claude/CLAUDE.md` §1.1 you may read outside cwd only when directed. **You are directed to read the following** as needed:

- `/Users/tigger/code/tigoss/first-folio/docs/ux_proposal.md` — the authoritative UX spec. Read in full before touching any view code.
- `/Users/tigger/code/tigoss/first-folio/docs/vision.md` — the product vision for First Folio itself.
- `/Users/tigger/code/tigoss/first-folio/docs/config.md` — full config schema (needed for the `ConfigReader` and the post-MVP Preferences pane).
- `/Users/tigger/code/tigoss/first-folio/docs/folio-*-help.md` — subcommand help texts (cross-check against live `folio <cmd> --help`).
- `/Users/tigger/code/tigoss/first-folio/ARCHITECTURE.md` — the CLI-side architecture.
- `/Users/tigger/code/tigoss/first-folio/README.md` — the CLI-side quickstart.
- `~/code/agents/AGENTS.md` and `~/code/agents/SDLC.md` — the governing process; you must follow these.

Do not read the first-folio Go source unless you are chasing a specific behaviour question. The CLI is a black box to the app.

## What you should not do

- Do not reimplement any parsing, rendering, or conversion logic in Swift. Every conversion goes through `folio`.
- Do not bundle a copy of `script.yaml`, `folio`, `typst`, or `pandoc` in the app in MVP — see `02-decisions.md`.
- Do not write to `script.yaml` in MVP. The app reads config; it does not create or modify it. This mirrors the CLI's contract in `docs/config.md`.
- Do not implement the Preferences pane in MVP. The ux_proposal describes it in §Preferences (post-MVP) and it stays out of scope.
- Do not begin implementation until Taḋg has resolved the "Pending decisions" in `02-decisions.md` and issued a Gate 1 authorization per `05-process.md`.

## What "done" for MVP looks like

Working app with three functional screens (Convert, Letter, Manuscript) that shell out to `folio`, surface progress/output/errors, and open produced PDFs in Preview. No signing, no notarization, no App Store. See `02-decisions.md` §MVP scope.
