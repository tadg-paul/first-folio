# First Folio

A format converter for stage plays. Reads plays in structured source formats and produces output in multiple target formats, preserving the semantic structure: acts, scenes, stage directions, characters, dialogue, and front matter.

## Quickstart

```bash
# Convert an org-mode play to Markdown
org-play-to-markdown play.org > play.md

# Convert to PDF (requires Typst)
org-play-to-pdf play.org

# Read from stdin
cat play.org | org-play-to-markdown > play.md

# Customise PDF output
org-play-to-pdf --font "Georgia" --font-size 14pt --page letter play.org
```

## Installation

```bash
make install    # symlinks scripts to ~/.local/bin/
make uninstall  # removes the symlinks
```

### Dependencies

- **Perl 5** (core modules only - no CPAN dependencies)
- **Typst** (required only for PDF output)

## Supported Formats

| Format | Read | Write |
|--------|------|-------|
| Org-mode play | Yes | - |
| Markdown | - | Yes |
| PDF (via Typst) | - | Yes |
| Fountain | Planned | Planned |

## Org-mode Play Format

The input format uses org-mode heading levels to encode play structure:

| Element | Org syntax |
|---------|-----------|
| Front matter | `#+TITLE:`, `#+AUTHOR:`, etc. |
| Act | `* Act I` (H1) |
| Scene | `** Scene 1` (H2) |
| Stage direction | `*** A bare stage.` (H3) |
| Character + direction | `**** BOB softly` (H4) |
| Dialogue | Plain text after H4 |
| Character table | `* CHARACTERS` followed by an org table |
| Excluded sections | `:noexport:` tag on any heading |
| Footnotes | `[fn:name] text` |
| Prop text | `- *"TEXT"*` |

## Configuration

First Folio reads configuration from `script.yaml` files. It never creates or modifies config files.

```yaml
# ~/.config/first-folio/script.yaml (global defaults)
# or alongside your source file (per-project overrides)

title: "About Time"
author: "Tadhg Paul"

render-stage-directions: true
render-intro: true
render-footnotes: true
render-character-table: true

folio:
  font: EB Garamond
  font-size: 11pt
  margin: 25mm
  page: a4
  indent: 5em
  dialogue-spacing: 1.1em
  direction-italic: true
  default-format: pdf
```

All config sources are merged in precedence order: CLI flags > local `script.yaml` > global `script.yaml` > built-in defaults. Each layer overrides individual keys, not the entire config.

The config file is shared with [yapper](https://github.com/tigger04/yapper) (TTS rendering). Each tool reads its own namespace and ignores the rest.

See [docs/config.md](docs/config.md) for the full schema and [examples/script.yaml](examples/script.yaml) for a complete annotated example.

## Project Structure

| Path | Purpose |
|------|---------|
| `org-play-to-markdown` | CLI: org-mode play to Markdown |
| `org-play-to-pdf` | CLI: org-mode play to PDF via Typst |
| `lib/OrgPlay/Parser.pm` | Shared parser - line-by-line state machine emitting typed events |
| `lib/OrgPlay/TypstTemplate.pm` | Typst preamble and title page template |
| `tests/regression/` | Regression test suite (run via `make test`) |
| `tests/one_off/` | One-off tests for specific issues |
| `docs/vision.md` | Project vision and goals |
| `docs/config.md` | Configuration system - schema, precedence, migration |
| `examples/script.yaml` | Complete annotated config example |
| `Makefile` | Build, install, test targets |

## Running Tests

```bash
make test           # regression suite
make test-one-off   # all one-off tests
make test-one-off ISSUE=5  # one-off tests for a specific issue
```

## Documentation

- [Vision](docs/vision.md) - project goals, supported formats, and direction of travel
- [Configuration](docs/config.md) - config schema, precedence, shared keys, migration
- [Formats](docs/formats.md) - format overview, event stream, and fidelity matrix
  - [Org-mode](docs/format-org.md) - org-mode play format schema
  - [Markdown](docs/format-markdown.md) - Markdown play format schema
  - [Fountain](docs/format-fountain.md) - Fountain format schema and fidelity analysis

## Licence

MIT - Copyright Taḋg Paul
