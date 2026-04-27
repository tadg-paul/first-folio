#!/usr/bin/env bash
# ABOUTME: Extended config parameter tests covering title page, spacing, transitions, and more.
# ABOUTME: Complements test_folio_config_params.sh with deeper coverage of the YAML schema.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
FOLIO="$PROJECT_DIR/bin/folio"
PASS=0
FAIL=0
FAILURES=()
TMPDIR_TEST="$(mktemp -d)"

cleanup() {
    rm -rf -- "$TMPDIR_TEST"
}
trap cleanup EXIT

pass() {
    PASS=$((PASS + 1))
    echo "  PASS: $1"
}

fail() {
    FAIL=$((FAIL + 1))
    FAILURES+=("$1")
    echo "  FAIL: $1"
    if [[ -n "${2:-}" ]]; then
        echo "        $2"
    fi
}

# Fixture with all element types
cat > "$TMPDIR_TEST/play.org" <<'ORG'
#+TITLE: Test Play
#+AUTHOR: Test Author
#+SUBTITLE: A Test Subtitle
#+DATE: 2026-01-01
#+VERSION: Draft v1

* Introduction
Some introductory text.

* CHARACTERS
|------+-------------|
| BOB  | A test char |
| CÁIT | Another     |
|------+-------------|

* Act I
** Scene 1
*** A bare stage.
**** BOB
Hello there.
A second line.
*** BOB crosses to the window.
**** CÁIT softly
Goodbye now.
***** CUT TO
ORG

# Helper: write YAML config, convert to .typ, grep for pattern
set_config() {
    printf '%s\n' "$1" > "$TMPDIR_TEST/script.yaml"
}

clear_config() {
    rm -f "$TMPDIR_TEST/script.yaml"
}

typ_has() {
    local desc="$1" pattern="$2"
    if "$FOLIO" convert "$TMPDIR_TEST/play.org" "$TMPDIR_TEST/out.typ" 2>/dev/null; then
        if grep -q "$pattern" "$TMPDIR_TEST/out.typ"; then
            pass "$desc"
        else
            fail "$desc" "pattern not found: $pattern"
        fi
    else
        fail "$desc" "folio exited non-zero"
    fi
    rm -f "$TMPDIR_TEST/out.typ"
}

typ_absent() {
    local desc="$1" pattern="$2"
    if "$FOLIO" convert "$TMPDIR_TEST/play.org" "$TMPDIR_TEST/out.typ" 2>/dev/null; then
        if grep -q "$pattern" "$TMPDIR_TEST/out.typ"; then
            fail "$desc" "pattern should be absent: $pattern"
        else
            pass "$desc"
        fi
    else
        fail "$desc" "folio exited non-zero"
    fi
    rm -f "$TMPDIR_TEST/out.typ"
}

md_has() {
    local desc="$1" pattern="$2"
    local out
    out=$("$FOLIO" convert "$TMPDIR_TEST/play.org" --to md 2>/dev/null || true)
    if echo "$out" | grep -qi "$pattern"; then
        pass "$desc"
    else
        fail "$desc" "pattern not found: $pattern"
    fi
}

md_absent() {
    local desc="$1" pattern="$2"
    local out
    out=$("$FOLIO" convert "$TMPDIR_TEST/play.org" --to md 2>/dev/null || true)
    if echo "$out" | grep -qi "$pattern"; then
        fail "$desc" "pattern should be absent: $pattern"
    else
        pass "$desc"
    fi
}

# ====================================================================
echo "=== Extended config parameter tests (issue #3) ==="

# --- Title page ---
echo ""
echo "Title page parameters"

clear_config
typ_has "RT-3.21: title.font-size defaults to 24pt" "24pt"
typ_has "RT-3.22: title.bold defaults to true" '"bold"'
typ_has "RT-3.23: subtitle appears in Typst output" "Test Subtitle"
typ_has "RT-3.24: author prefix 'by ' appears" "by Test Author"
typ_has "RT-3.25: date appears in footer" "2026-01-01"
typ_has "RT-3.26: version appears in footer" "Draft v1"

set_config 'folio:
  title-page:
    title:
      font-size: 18pt'
typ_has "RT-3.27: title.font-size overridden to 18pt" "18pt"

set_config 'folio:
  title-page:
    title:
      bold: false'
typ_absent "RT-3.28: title.bold: false removes bold" '"bold".*Test Play'

set_config 'folio:
  title-page:
    title:
      italic: true'
typ_has "RT-3.29: title.italic: true adds italic" '"italic"'

set_config 'folio:
  title-page:
    subtitle:
      italic: false'
# The subtitle should not have italic style when overridden
typ_has "RT-3.30: subtitle.italic: false respected" "Test Subtitle"

set_config 'folio:
  title-page:
    author:
      prefix: ""'
typ_absent "RT-3.31: author.prefix: empty removes 'by'" "by Test Author"

set_config 'folio:
  title-page:
    date:
      position: bottom-right'
typ_has "RT-3.32: date.position: bottom-right" "2026-01-01"

# --- Speech spacing ---
echo ""
echo "Speech spacing and dialogue"

set_config 'folio:
  positioning:
    speech:
      space-before: 2.5em'
typ_has "RT-3.33: speech.space-before: 2.5em" "2.5em"

set_config 'folio:
  positioning:
    speech:
      dialogue:
        wrap-indent: 10em'
typ_has "RT-3.34: dialogue.wrap-indent: 10em" "10em"

# --- Speech instruction ---
echo ""
echo "Speech instruction (parenthetical direction)"

set_config 'folio:
  positioning:
    speech:
      speech-instruction:
        prefix: "["
        suffix: "]"'
typ_has "RT-3.35: speech-instruction prefix/suffix override in template" '\[#direction\]'
clear_config
typ_has "RT-3.36: default speech-instruction uses parentheses" "(.*softly.*)"

# --- Stage direction ---
echo ""
echo "Stage direction positioning"

set_config 'folio:
  positioning:
    stage-direction:
      space-before: 3em'
typ_has "RT-3.37: stage-direction.space-before: 3em" "3em"

set_config 'folio:
  positioning:
    stage-direction:
      align: center'
typ_has "RT-3.38: stage-direction.align: center" "align(center)"

# --- Transition ---
echo ""
echo "Transition positioning"

clear_config
typ_has "RT-3.39: transition renders right-aligned by default" "align(right)"
typ_has "RT-3.40: transition text present" "CUT TO"

set_config 'render:
  transitions: false'
md_absent "RT-3.41: render.transitions: false suppresses transitions" "CUT TO"

# --- Act header ---
echo ""
echo "Act header extended"

set_config 'folio:
  positioning:
    act-header:
      font-size: 18pt'
typ_has "RT-3.42: act-header.font-size: 18pt" "18pt"

set_config 'folio:
  positioning:
    act-header:
      bold: false'
typ_has "RT-3.43: act-header.bold: false" '"regular"'

clear_config
typ_has "RT-3.44: act-header.page-break-before: true (default)" "pagebreak"

# --- Scene header ---
echo ""
echo "Scene header extended"

set_config 'folio:
  positioning:
    scene-header:
      font-size: 14pt'
typ_has "RT-3.45: scene-header.font-size: 14pt" "14pt"

set_config 'folio:
  positioning:
    scene-header:
      bold: false'
typ_has "RT-3.46: scene-header.bold: false" '"regular"'

set_config 'folio:
  positioning:
    scene-header:
      space-before: 3em'
typ_has "RT-3.47: scene-header.space-before: 3em" "3em"

# --- Render toggles (remaining) ---
echo ""
echo "Remaining render toggles"

set_config 'render:
  footnotes: false'
md_absent "RT-3.48: render.footnotes: false suppresses footnotes" "footnote"

# --- Frontmatter header ---
echo ""
echo "Frontmatter header styling"

clear_config
typ_has "RT-3.49: frontmatter header renders Introduction" "Introduction"

# --- Metadata ---
echo ""
echo "Metadata: subtitle, date, version"

set_config 'subtitle: "Override Sub"'
md_has "RT-3.50: config subtitle overrides source" "Override Sub"

set_config 'date: "2099-12-31"'
typ_has "RT-3.51: config date overrides source in PDF" "2099-12-31"

set_config 'version: "Final"'
typ_has "RT-3.52: config version overrides source in PDF" "Final"

clear_config

# --- Style-specific config files ---
echo ""
echo "Style-specific config files (script-{style}.yaml)"

# Global script-us.yaml should apply when style=us
mkdir -p "$HOME/.config/first-folio"
printf 'folio:\n  font: Courier Prime\n' > "$HOME/.config/first-folio/script-us.yaml"
set_config 'folio:
  style: us'
typ_has "RT-3.53: global script-us.yaml applies when style=us" "Courier Prime"
trash "$HOME/.config/first-folio/script-us.yaml" 2>/dev/null

# Local script-british.yaml should apply when style=british
printf 'folio:\n  margin: 40mm\n' > "$TMPDIR_TEST/script-british.yaml"
clear_config
typ_has "RT-3.54: local script-british.yaml applies for british style" "40mm"
rm -f "$TMPDIR_TEST/script-british.yaml"

# Style-specific overrides style-agnostic: local script-us.yaml overrides global script.yaml
mkdir -p "$HOME/.config/first-folio"
printf 'folio:\n  font: Georgia\n' > "$HOME/.config/first-folio/script.yaml"
printf 'folio:\n  font: Courier Prime\n' > "$TMPDIR_TEST/script-us.yaml"
set_config 'folio:
  style: us'
typ_has "RT-3.55: local script-us.yaml overrides global script.yaml font" "Courier Prime"
trash "$HOME/.config/first-folio/script.yaml" 2>/dev/null
rm -f "$TMPDIR_TEST/script-us.yaml"

clear_config

# ====================================================================
echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="
if [[ ${#FAILURES[@]} -gt 0 ]]; then
    echo ""
    echo "Failures:"
    for f in "${FAILURES[@]}"; do
        echo "  - $f"
    done
    exit 1
fi
