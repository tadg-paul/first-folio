# 03 · Architecture

**Read first:** `/Users/tigger/code/tigoss/first-folio/docs/ux_proposal.md` §Architecture. This file distils and augments it.

## Runtime shape

Thin SwiftUI shell over the `folio` subprocess. No conversion or rendering logic in Swift.

```
+-----------------------------------------------------------------+
|  SwiftUI App (macOS)                                            |
|                                                                 |
|  NavigationSplitView                                            |
|  ├── Sidebar (Convert | Letter | Manuscript)                    |
|  └── Detail pane (per-command view)                             |
|                                                                 |
|      Views call ArgBuilder → FolioRunner → Process(folio)       |
|                                                                 |
+-----------------------------------------------------------------+
                              │
                              ▼
+-----------------------------------------------------------------+
|  folio (external binary on PATH)                                |
|      spawns typst / pandoc as needed                            |
|      reads ~/.config/first-folio/script.yaml                    |
|      writes output files the user chose                         |
+-----------------------------------------------------------------+
```

## Wrapper components (Swift-side)

All in a `FolioCore` module; views depend on it, it does not depend on views. Names below match `docs/ux_proposal.md` §Architecture.

### FolioRunner

Only place that spawns processes.

- Input: `argv: [String]`, working directory, environment.
- Behaviour: spawns `folio` via `Process`, streams stdout and stderr through `AsyncThrowingStream<String, Error>`, returns exit status.
- Cancellation: exposes a cancel primitive that sends `SIGTERM`. `Run` becomes `Cancel` in the UI while a subprocess is active.
- Environment: inherits the user's `PATH` (see `06-first-steps.md` §PATH capture). Sets `LC_ALL=en_IE.UTF-8` (or `en_GB.UTF-8`; not decided — probably safest to inherit).

### ArgBuilder

Per-command state → validated `[String]` argv. **Blank fields are omitted rather than passed as empty strings** — this is load-bearing. It lets `folio`'s config-file values and built-in defaults apply exactly as they would on the terminal.

Example (Convert): if `--font` is left blank in the UI, do not pass `--font ""`. Do not pass `--font` at all.

### OrgLetterScanner

**Only needed if R4 (adding `folio letter --list`) is rejected** (see `02-decisions.md`). Otherwise delete this component from the plan.

If retained: scans `**** Name :tag:` H4 headings under `:letter:` sections. Deliberately shallow — enough to list recipients, not to interpret bodies. Document it as a mirror of `docs/format-org.md` and note that the CLI is authoritative for interpretation.

### ConfigReader

Read-only. Parses `~/.config/first-folio/script.yaml` (and a local `script.yaml` next to a loaded source file) to prefill the Convert screen's Style / Font / Page-size defaults.

Never writes. Never mutates. If the file is missing, return an empty struct — that is the normal case for a fresh install.

Use `Yams` (or another well-maintained SwiftPM YAML library). Do not roll your own parser.

### OutputParser

Watches `FolioRunner`'s streams for the "output produced" signals per subcommand. **The ux_proposal is inaccurate here — see `04-cli-contract.md` §Output patterns for the actual behaviour.** Summary:

- `convert`: silent on success. The app already knows the target path from its own argv; use that to add a Reveal-in-Finder / Open-in-Preview row when exit is 0.
- `letter`: parse lines matching `^Generated: (.+)$` — one PDF per line, followed by a summary `^N cover letter\(s\) generated\.$`.
- `manuscript`: silent on success. Use the app's known target path.

## Views

Three views, one per command. Match the specs in `docs/ux_proposal.md` §Screen: Convert / §Screen: Letter / §Screen: Manuscript exactly.

Shared bottom strip on every detail pane:

- Primary action (`Run` / `Convert` / `Generate Letters` / `Render Manuscript`) — becomes `Cancel` when running.
- Monospace output log streaming stdout/stderr.
- Reveal-in-Finder / Open-in-Preview affordances derived from `OutputParser` output.
- Non-zero exit surfaces an alert: last stderr line as title, full stderr in a "Show details" disclosure.

## Storage

- `UserDefaults` for the recent-files list per command (no cross-command sharing) and the discovered/overridden `folio` binary path.
- No CoreData, no SwiftData, no on-disk databases. If you feel one is needed, you have gone off-scope — pause and check.

## Suggested project layout (Xcode / SwiftPM)

Not prescriptive; adjust to fit Xcode's own conventions if there's a clash.

```
FirstFolio.xcodeproj                (or Package.swift for a SwiftPM-only build)
Sources/
  FirstFolioApp/                    (SwiftUI app target)
    App.swift
    RootView.swift                  (NavigationSplitView)
    Convert/
      ConvertView.swift
      ConvertViewModel.swift
    Letter/
      LetterView.swift
      LetterViewModel.swift
    Manuscript/
      ManuscriptView.swift
      ManuscriptViewModel.swift
    Shared/
      OutputLogView.swift
      RunButton.swift
  FolioCore/                        (SwiftPM library target — unit-testable)
    FolioRunner.swift
    ArgBuilder+Convert.swift
    ArgBuilder+Letter.swift
    ArgBuilder+Manuscript.swift
    ConfigReader.swift
    OutputParser.swift
    OrgLetterScanner.swift          (only if R4 is rejected)
Tests/
  FolioCoreTests/
    ArgBuilderTests.swift
    ConfigReaderTests.swift
    OutputParserTests.swift
    FolioRunnerTests.swift          (integration; requires folio on PATH)
```

## Testing strategy summary

Detail in `05-process.md`. The short version:

- **Unit tests** in `FolioCoreTests` for `ArgBuilder`, `ConfigReader`, `OutputParser`. These do not spawn processes.
- **Integration tests** for `FolioRunner` — spawn a real `folio` against fixtures copied from `/Users/tigger/code/tigoss/first-folio/examples/`. Requires `folio` on `PATH`. Skip with a clear message if not present.
- **UT (user tests)** — every UI screen is a UT. Only Taḋg marks these ✅/❌.

## Deltas from the ux_proposal

Two you must not repeat:

1. **`convert`'s "wrote X" line does not exist.** The proposal describes an `OutputParser` that "watches stderr/stdout for the 'wrote X' lines emitted by `folio`". That is only accurate for `letter` (which emits `Generated: <path>`). `convert` and `manuscript` are silent on success. Use the app's own known target path.
2. **The `--to` regex form.** The proposal's Letter screen says "builds a `--to` regex of the form `^(A|B|C)$` (escaped)". The CLI's `--to` is documented in `docs/folio-letter-help.md` as "matching recipients, using a regex substring". Anchored alternation still works fine — it just is not required. Proceed with the proposal's approach but do not treat it as the only valid one.
