# 01 · The project

## What First Folio is

A CLI-first Go application that converts stage plays between formats and renders prose manuscripts. Three subcommands:

- `folio convert` — Org / Markdown / Fountain / Typst / PDF for stage plays.
- `folio letter` — recipient-specific PDF cover letters from `:letter:` tagged Org sections.
- `folio manuscript` — prose Markdown / Org → Typst or PDF.

Full quickstart: `/Users/tigger/code/tigoss/first-folio/README.md`. Vision: `/Users/tigger/code/tigoss/first-folio/docs/vision.md`. Currently at v0.4.10.

## What you are building

A macOS SwiftUI companion app that wraps `folio`. Target user: playwrights and dramaturgs who do not want to use a terminal, YAML, or Fountain-style markup on the command line.

Core principle (from `docs/ux_proposal.md` §Core principle):

> **The app is an abstraction layer over the CLI. It offers no capability the `folio` CLI cannot perform.** If a workflow requires new behaviour, the CLI grows first and the app catches up.

Non-goals for the app: in-app editing, live preview (out of MVP), config-file writing (out of MVP), any re-implementation of parsing or rendering in Swift.

## Relationship to the first-folio repo

- **Separate repo.** Confirmed by Taḋg 2026-07-20. Rationale: different toolchain, different release cadence, keeps the Go repo Go-only per its `ARCHITECTURE.md`.
- **CLI is authoritative.** If the app needs behaviour the CLI does not have, the CLI grows first (issue against `tigger-developer/first-folio`) and the app follows.
- **No cross-repo code sharing.** The app does not vendor Go source, does not link against Go, does not parse `folio`'s internal formats beyond what CLI stdout tells it.

## Where this project lives

Taḋg has created the new repo. **Location and name are not yet recorded in this handover.**

Your first action after reading these notes: ask Taḋg to confirm the repo path and name, and record it here (edit this file, add a "Repository" section beneath this line). Do not assume a path.

## Who this handover is for

An AI coding agent (Claude Code, Codex, or equivalent) working in the new repo's cwd. The agent inherits the universal prohibitions in `~/.claude/CLAUDE.md` (`~/code/agents/AGENTS.md`) and the code SDLC in `~/code/agents/SDLC.md`. See `05-process.md`.
