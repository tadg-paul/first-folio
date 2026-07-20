# 05 · Process (SDLC, gates, canary)

## The governing documents

You are bound by:

- `~/.claude/CLAUDE.md` (or `~/code/agents/AGENTS.md` — same content) — universal prohibitions and relationship principles.
- `~/code/agents/SDLC.md` — code-specific rules, gates, and modes.

Read both before doing any work. This file summarizes the parts most likely to bite you on this project; it is not a substitute.

## Canary

Begin every message to Taḋg with the mode-prefixed canary plus SDLC suffix, e.g. `[MODE PAIR] EHLO SDLC`. If you have not read and agreed with `AGENTS.md` and `SDLC.md`, do not utter it.

## Modes

- **MODE PAIR** is the default at the start of every conversation. Only Taḋg can change mode.
- **MODE DELIVER** is possible for this project but not automatic. If Taḋg requests it, produce a compliant goal statement (goal, definition of done, scope, exclusions, human-only checkpoints) and wait for `CONFIRM DELIVER`.

## Gates

- **Gate 1 — PROCEED n [n ...]** authorises implementation for the listed issue(s).
- **Gate 2 — APPROVED n [n ...]** authorises closure of the listed issue(s).

You may not write these keywords yourself. Only Taḋg does. `/proceed N` (a skill) counts as `PROCEED N`; same for `/approved N`.

**Do not write or modify source code until you hold a valid PROCEED for the relevant issue, or the issue is an in-scope child of a confirmed MODE DELIVER master.**

Documentation-only work (like this handover) is exempt from the gate rule.

## Issues

Every unit of code work has a GitHub issue in the new repo. Draft each with the appropriate skill:

- `/draft-design-issue` — combined draft + solution design (this project's default).
- `/draft-issue` — issue + ACs only, with `/design-solution` following.
- `/draft-bug-fix` — bug fix referencing existing ACs, no new AC table.

The skills produce the issue via `gh` in the current repo. Do not create issues on `tigger-developer/first-folio` for app work — CLI-side changes (e.g. `folio letter --list`, per `02-decisions.md` R4) do go there.

## AC discipline (project-specific reminders)

Read `~/code/agents/docs/ISSUES.md` for the full standard. Project-specific traps:

- **The app's ACs are user-facing behaviour, not implementation detail.** "The app launches an Xcode-generated shell" is not an AC. "Selecting an Org file with `:letter:` sections and clicking Generate Letters produces one PDF per recipient in the chosen directory" is.
- **UI ACs will need UT (user test) rows.** UTs are human-only judgement (visual layout, feel, subjective quality). You may never mark a UT ✅ or ❌. Present the app in the state the UT requires — start it up, load fixtures, place it at the exact screen — and hand back to Taḋg for the call.
- **UT presentation obligation** (SDLC §Hard blocks): "the agent is responsible for starting the application or tool, preparing representative data and configuration, and placing it at the point the human must inspect." Do not hand Taḋg a click-sequence or setup checklist for a UT.

## Testing

- `make test` (or Xcode's test action, wired through a Makefile) must run all `FolioCoreTests` and pass with zero errors and no new warnings.
- Unit tests do not spawn `folio`. Integration tests that do must skip cleanly when `folio` is not on `PATH` (print a clear reason, exit skipped, not failed).
- Do not use UI test frameworks unless Taḋg asks. The macOS UI test tooling is flaky and the ux_proposal never asked for it. UI behaviour is validated via UT.

## Commit cadence (SDLC §Commit cadence)

Commit at useful checkpoints: after red tests are written, after implementation passes, after doc updates. Do not accumulate a large dirty tree while continuing to change code.

Never use auto-close keywords in commit messages. Never sign commits with AI attribution.

## Hand-back rules (SDLC §Problems, obstacles, and blockers)

An obstacle is not a blocker just because your first attempt failed. Diagnose, propose a solution, try it, verify, and continue. Hand back only when you genuinely need Taḋg's authorization / access / product decision / hardware / a security-widening call.

Every hand-back must include: intended outcome, current state, evidence, root-cause diagnosis, remedies attempted with results, recommended resolution, and the one precise thing you need from Taḋg.

## When to hand back on this project specifically

- Any decision in `07-open-questions.md` — those are Taḋg-only.
- Any change to the CLI (`folio` behaviour, flags, output format). File a CLI-side issue and hand back with the link and impact.
- Any UT — you never mark it passed.
- Any request that would widen access (e.g. sandbox exception, network permission entitlement, requesting user data outside `~/.config/first-folio/`).

## Language and voice

`~/.claude/CLAUDE.md` requires Hiberno-English with OED spellings (`-ize`, `-our`). In code (identifiers, comments, docstrings) match the surrounding style; in prose (issues, ACs, this handover) use OED.

## What the harness will not do for you

- The harness will not remember decisions between conversations. If Taḋg makes a call ("we're going with R2"), record it in `handover-ux/02-decisions.md` under Decided **and** on the relevant issue, so the next agent sees it.
- The harness will not create the repo, the Xcode project, the `.gitignore`, or the initial commit. `06-first-steps.md` walks you through those.
