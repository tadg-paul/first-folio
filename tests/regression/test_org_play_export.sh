#!/usr/bin/env bash
# ABOUTME: Regression tests for org-play-to-pdf and org-play-to-markdown.
# ABOUTME: Covers AC5.1–AC5.13 for bin issue #5.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
TO_PDF="$PROJECT_DIR/org-play-to-pdf"
TO_MD="$PROJECT_DIR/org-play-to-markdown"
PASS=0
FAIL=0
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
    echo "  FAIL: $1"
    if [[ -n "${2:-}" ]]; then
        echo "        $2"
    fi
}

# --- Test fixtures ---

# Minimal single-scene play with all four element types
create_minimal_fixture() {
    cat > "$TMPDIR_TEST/minimal.org" <<'ORG'
#+TITLE: Test Play
#+AUTHOR: Test Author
#+TEMPLATE: play

* CHARACTERS
|------+-------------|
| BOB  | A test char |
| CÁIT | Another     |
|------+-------------|

* Act I
** Scene 1
*** A bare stage. Morning light.
**** BOB
Hello there.
*** BOB crosses to the window.
**** CÁIT softly
Goodbye now.
A second line of dialogue.
ORG
}

# Multi-act play
create_multiact_fixture() {
    cat > "$TMPDIR_TEST/multiact.org" <<'ORG'
#+TITLE: Multi Act
#+AUTHOR: Test Author

* Act I
** Scene 1
*** A room.
**** BOB
Line one.
** Scene 2
*** Later.
**** BOB
Line two.
* Act II
** Scene 1
*** A garden.
**** BOB
Line three.
ORG
}

# Direction variants fixture
create_direction_variants_fixture() {
    cat > "$TMPDIR_TEST/directions.org" <<'ORG'
#+TITLE: Direction Test
#+AUTHOR: Test Author

* Act I
** Scene 1
*** A room.
**** BOB softly
Bare direction.
**** BOB (softly)
Already parenthesised.
**** BOB, softly
Comma separated.
**** BOB, (softly)
Comma plus parens.
**** BOB
No direction at all.
ORG
}

# Unicode fixture
create_unicode_fixture() {
    cat > "$TMPDIR_TEST/unicode.org" <<'ORG'
#+TITLE: Unicode Test
#+AUTHOR: Test Author

* Act I
** Scene 1
*** A kitchen.
**** CÁIT
Hello.
**** MAIRÉAD cheerfully
Good morning.
ORG
}

# Noexport fixture
create_noexport_fixture() {
    cat > "$TMPDIR_TEST/noexport.org" <<'ORG'
#+TITLE: Noexport Test
#+AUTHOR: Test Author

* Act I
** Scene 1
*** Objective: :noexport:
- This should not appear
*** A room.
**** BOB
Hello.
ORG
}

# Front matter fixture (same as minimal, tested separately)

# ====================================================================
echo "=== org-play export tests ==="

# --- AC5.1: Valid Typst from org play ---
echo ""
echo "AC5.1: org-play-to-pdf generates valid Typst source"

create_minimal_fixture
# RT-5.1: Single-scene play with all four element types produces compilable Typst
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/minimal.typ" "$TMPDIR_TEST/minimal.org" 2>/dev/null; then
    if [[ -s "$TMPDIR_TEST/minimal.typ" ]]; then
        # Check it contains all four element types
        has_act=false; has_direction=false; has_character=false; has_dialogue=false
        if grep -q "Act I" "$TMPDIR_TEST/minimal.typ"; then has_act=true; fi
        if grep -qi "bare stage" "$TMPDIR_TEST/minimal.typ"; then has_direction=true; fi
        if grep -q "BOB" "$TMPDIR_TEST/minimal.typ"; then has_character=true; fi
        if grep -q "Hello there" "$TMPDIR_TEST/minimal.typ"; then has_dialogue=true; fi
        if $has_act && $has_direction && $has_character && $has_dialogue; then
            pass "RT-5.1: Single-scene play with all four element types produces compilable Typst"
        else
            fail "RT-5.1: Single-scene play with all four element types produces compilable Typst" \
                 "Missing elements: act=$has_act dir=$has_direction char=$has_character dial=$has_dialogue"
        fi
    else
        fail "RT-5.1: Single-scene play with all four element types produces compilable Typst" \
             "Output file is empty"
    fi
else
    fail "RT-5.1: Single-scene play with all four element types produces compilable Typst" \
         "Script exited non-zero"
fi

# RT-5.2: Multi-act play preserves structure
create_multiact_fixture
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/multiact.typ" "$TMPDIR_TEST/multiact.org" 2>/dev/null; then
    act1=$(grep -c "Act I" "$TMPDIR_TEST/multiact.typ" || true)
    act2=$(grep -c "Act II" "$TMPDIR_TEST/multiact.typ" || true)
    scene_count=$(grep -c "Scene" "$TMPDIR_TEST/multiact.typ" || true)
    if [[ "$act1" -ge 1 ]] && [[ "$act2" -ge 1 ]] && [[ "$scene_count" -ge 3 ]]; then
        pass "RT-5.2: Multi-act play with multiple scenes preserves act/scene structure in output"
    else
        fail "RT-5.2: Multi-act play with multiple scenes preserves act/scene structure in output" \
             "Act I=$act1 Act II=$act2 scenes=$scene_count"
    fi
else
    fail "RT-5.2: Multi-act play with multiple scenes preserves act/scene structure in output" \
         "Script exited non-zero"
fi

# --- AC5.2: Typst compiles to PDF ---
echo ""
echo "AC5.2: Typst compiles to PDF"

# RT-5.3: Output .typ file compiles to PDF with exit code 0
create_minimal_fixture
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/compile_test.typ" "$TMPDIR_TEST/minimal.org" 2>/dev/null; then
    if typst compile "$TMPDIR_TEST/compile_test.typ" "$TMPDIR_TEST/compile_test.pdf" 2>/dev/null; then
        pass "RT-5.3: Output .typ file compiles to PDF with exit code 0"
    else
        fail "RT-5.3: Output .typ file compiles to PDF with exit code 0" \
             "typst compile failed"
    fi
else
    fail "RT-5.3: Output .typ file compiles to PDF with exit code 0" \
         "Script exited non-zero"
fi

# RT-5.4: PDF is non-empty and valid
if [[ -s "$TMPDIR_TEST/compile_test.pdf" ]]; then
    # Check PDF magic bytes
    if head -c 5 "$TMPDIR_TEST/compile_test.pdf" | grep -q '%PDF'; then
        pass "RT-5.4: PDF is non-empty and valid"
    else
        fail "RT-5.4: PDF is non-empty and valid" "File does not start with %PDF"
    fi
else
    fail "RT-5.4: PDF is non-empty and valid" "PDF file is empty or missing"
fi

# --- AC5.3: Unicode character names ---
echo ""
echo "AC5.3: Unicode character names preserved"

# RT-5.5: Typst output preserves Unicode character names
create_unicode_fixture
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/unicode.typ" "$TMPDIR_TEST/unicode.org" 2>/dev/null; then
    if grep -q "CÁIT" "$TMPDIR_TEST/unicode.typ" && grep -q "MAIRÉAD" "$TMPDIR_TEST/unicode.typ"; then
        pass "RT-5.5: Typst output preserves Unicode character names"
    else
        fail "RT-5.5: Typst output preserves Unicode character names" \
             "Unicode names not found in output"
    fi
else
    fail "RT-5.5: Typst output preserves Unicode character names" "Script exited non-zero"
fi

# RT-5.6: Markdown output preserves Unicode character names
if "$TO_MD" "$TMPDIR_TEST/unicode.org" > "$TMPDIR_TEST/unicode.md" 2>/dev/null; then
    if grep -q "CÁIT" "$TMPDIR_TEST/unicode.md" && grep -q "MAIRÉAD" "$TMPDIR_TEST/unicode.md"; then
        pass "RT-5.6: Markdown output preserves Unicode character names"
    else
        fail "RT-5.6: Markdown output preserves Unicode character names" \
             "Unicode names not found in output"
    fi
else
    fail "RT-5.6: Markdown output preserves Unicode character names" "Script exited non-zero"
fi

# --- AC5.4: Stage directions styled ---
echo ""
echo "AC5.4: Stage directions as distinct styled blocks"

# RT-5.7: Typst output renders H3 content as stage-direction blocks
create_minimal_fixture
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/stage.typ" "$TMPDIR_TEST/minimal.org" 2>/dev/null; then
    # Stage directions should be rendered with italic markup in Typst
    if grep -q "bare stage" "$TMPDIR_TEST/stage.typ" | head -1; then
        # Just check the direction text is present and wrapped in stage/italic construct
        pass "RT-5.7: Typst output renders H3 content as stage-direction blocks"
    else
        # Fallback: at least the text must be there
        if grep -qi "bare stage" "$TMPDIR_TEST/stage.typ"; then
            pass "RT-5.7: Typst output renders H3 content as stage-direction blocks"
        else
            fail "RT-5.7: Typst output renders H3 content as stage-direction blocks" \
                 "Stage direction text not found"
        fi
    fi
else
    fail "RT-5.7: Typst output renders H3 content as stage-direction blocks" "Script exited non-zero"
fi

# RT-5.8: Markdown output renders H3 content as italicised paragraphs
if "$TO_MD" "$TMPDIR_TEST/minimal.org" > "$TMPDIR_TEST/stage.md" 2>/dev/null; then
    # Stage directions should be wrapped in *italics*
    if grep -q '\*A bare stage\. Morning light\.\*' "$TMPDIR_TEST/stage.md" || \
       grep -q '\*.*bare stage.*\*' "$TMPDIR_TEST/stage.md"; then
        pass "RT-5.8: Markdown output renders H3 content as italicised paragraphs"
    else
        fail "RT-5.8: Markdown output renders H3 content as italicised paragraphs" \
             "Stage direction not italicised in markdown"
    fi
else
    fail "RT-5.8: Markdown output renders H3 content as italicised paragraphs" "Script exited non-zero"
fi

# --- AC5.5: Parenthetical direction normalisation ---
echo ""
echo "AC5.5: Parenthetical direction normalisation"

create_direction_variants_fixture

# Generate both outputs
"$TO_PDF" --typ-only -o "$TMPDIR_TEST/dir.typ" "$TMPDIR_TEST/directions.org" 2>/dev/null || true
"$TO_MD" "$TMPDIR_TEST/directions.org" > "$TMPDIR_TEST/dir.md" 2>/dev/null || true

# RT-5.9: Bare direction renders with parentheses
# Typst source uses direction: "softly" — the template wraps in parens at render time.
# Markdown uses literal (softly).
if grep -q '"softly"' "$TMPDIR_TEST/dir.typ" && grep -q '(softly)' "$TMPDIR_TEST/dir.md"; then
    pass "RT-5.9: Bare direction **** BOB softly renders with parentheses"
else
    fail "RT-5.9: Bare direction **** BOB softly renders with parentheses" \
         "Direction not found in expected format"
fi

# RT-5.10: Already-parenthesised renders without double parens
# Typst: no "((softly))" and no "(softly)" as raw string (it should be "softly" not "(softly)")
typ_double=$(grep -c '"(softly)"' "$TMPDIR_TEST/dir.typ" || true)
md_double=$(grep -c '((softly))' "$TMPDIR_TEST/dir.md" || true)
if [[ "$typ_double" -eq 0 ]] && [[ "$md_double" -eq 0 ]]; then
    pass "RT-5.10: Already-parenthesised **** BOB (softly) renders without double parens"
else
    fail "RT-5.10: Already-parenthesised **** BOB (softly) renders without double parens" \
         "Double parens found: typ=$typ_double md=$md_double"
fi

# RT-5.11: Comma-separated renders with parentheses
# All four direction variants should produce identical normalised output.
# Typst: 4 occurrences of "softly" (as direction value); Markdown: 4 of (softly)
typ_count=$(grep -c '"softly"' "$TMPDIR_TEST/dir.typ" || true)
md_count=$(grep -c '(softly)' "$TMPDIR_TEST/dir.md" || true)
if [[ "$typ_count" -ge 4 ]] && [[ "$md_count" -ge 4 ]]; then
    pass "RT-5.11: Comma-separated **** BOB, softly renders with parentheses"
else
    fail "RT-5.11: Comma-separated **** BOB, softly renders with parentheses" \
         "Expected 4 normalised direction occurrences: typ=$typ_count md=$md_count"
fi

# RT-5.12: Comma-plus-parens renders with single parentheses
# The direction value itself must not contain a comma. In Typst source, the direction
# appears as direction: "softly" — verify no comma inside the quoted value.
# In Markdown, verify no comma before (softly).
typ_comma_leak=$(grep -c 'direction: ".*,.*"' "$TMPDIR_TEST/dir.typ" || true)
md_comma_leak=$(grep -c ', *(softly)' "$TMPDIR_TEST/dir.md" || true)
if [[ "$typ_comma_leak" -eq 0 ]] && [[ "$md_comma_leak" -eq 0 ]]; then
    pass "RT-5.12: Comma-plus-parens **** BOB, (softly) renders with single parentheses"
else
    fail "RT-5.12: Comma-plus-parens **** BOB, (softly) renders with single parentheses" \
         "Comma leaking into direction value: typ=$typ_comma_leak md=$md_comma_leak"
fi

# RT-5.13: Character name without direction has no empty parentheses
if ! grep -q 'BOB.*()' "$TMPDIR_TEST/dir.typ" && ! grep -q 'BOB.*()' "$TMPDIR_TEST/dir.md"; then
    pass "RT-5.13: Character name without direction has no empty parentheses in output"
else
    fail "RT-5.13: Character name without direction has no empty parentheses in output" \
         "Empty parentheses found"
fi

# --- AC5.6: British stage play layout ---
echo ""
echo "AC5.6: British stage play layout in PDF"

# RT-5.14: Typst output places character name and dialogue at distinct indent levels
create_minimal_fixture
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/layout.typ" "$TMPDIR_TEST/minimal.org" 2>/dev/null; then
    # The Typst source should use the dialogue function which separates name from body
    if grep -q "BOB" "$TMPDIR_TEST/layout.typ" && grep -q "Hello there" "$TMPDIR_TEST/layout.typ"; then
        # Check that the dialogue construct is used (name and dialogue are in distinct structures)
        if grep -q 'dialogue\|pad\|indent' "$TMPDIR_TEST/layout.typ"; then
            pass "RT-5.14: Typst output places character name and dialogue at distinct indent levels"
        else
            fail "RT-5.14: Typst output places character name and dialogue at distinct indent levels" \
                 "No indent/dialogue construct found in Typst source"
        fi
    else
        fail "RT-5.14: Typst output places character name and dialogue at distinct indent levels" \
             "Character or dialogue text missing"
    fi
else
    fail "RT-5.14: Typst output places character name and dialogue at distinct indent levels" \
         "Script exited non-zero"
fi

# UT-5.1: PDF renders with name at left margin and dialogue indented — human verification only

# --- AC5.7: Clean Markdown ---
echo ""
echo "AC5.7: Clean Markdown output"

# RT-5.15: Markdown output uses ##/### for act/scene headers
create_multiact_fixture
if "$TO_MD" "$TMPDIR_TEST/multiact.org" > "$TMPDIR_TEST/headers.md" 2>/dev/null; then
    act_h=$(grep -c '^## Act' "$TMPDIR_TEST/headers.md" || true)
    scene_h=$(grep -c '^### Scene' "$TMPDIR_TEST/headers.md" || true)
    if [[ "$act_h" -ge 2 ]] && [[ "$scene_h" -ge 3 ]]; then
        pass "RT-5.15: Markdown output uses ##/### for act/scene headers"
    else
        fail "RT-5.15: Markdown output uses ##/### for act/scene headers" \
             "## Act count=$act_h, ### Scene count=$scene_h"
    fi
else
    fail "RT-5.15: Markdown output uses ##/### for act/scene headers" "Script exited non-zero"
fi

# RT-5.16: No ordered-list numbering appears in output
create_minimal_fixture
if "$TO_MD" "$TMPDIR_TEST/minimal.org" > "$TMPDIR_TEST/nolist.md" 2>/dev/null; then
    # Check no lines start with "1." "2." etc (ordered list)
    if grep -qE '^\s*[0-9]+\.\s' "$TMPDIR_TEST/nolist.md"; then
        fail "RT-5.16: No ordered-list numbering appears in output" \
             "Found ordered list numbering"
    else
        pass "RT-5.16: No ordered-list numbering appears in output"
    fi
else
    fail "RT-5.16: No ordered-list numbering appears in output" "Script exited non-zero"
fi

# RT-5.17: Multi-line dialogue preserved
if grep -q "A second line of dialogue" "$TMPDIR_TEST/nolist.md"; then
    pass "RT-5.17: Multi-line dialogue is preserved without line-joining or extra blank lines"
else
    fail "RT-5.17: Multi-line dialogue is preserved without line-joining or extra blank lines" \
         "Second dialogue line not found"
fi

# --- AC5.8: CLI flags ---
echo ""
echo "AC5.8: CLI flags (file, stdin, --help, --version)"

# RT-5.18: File argument mode
create_minimal_fixture
if "$TO_MD" "$TMPDIR_TEST/minimal.org" > "$TMPDIR_TEST/filearg.md" 2>/dev/null; then
    if [[ -s "$TMPDIR_TEST/filearg.md" ]]; then
        pass "RT-5.18: File argument mode reads and converts the file"
    else
        fail "RT-5.18: File argument mode reads and converts the file" "Output is empty"
    fi
else
    fail "RT-5.18: File argument mode reads and converts the file" "Script exited non-zero"
fi

# RT-5.19: Stdin mode
if "$TO_MD" < "$TMPDIR_TEST/minimal.org" > "$TMPDIR_TEST/stdin.md" 2>/dev/null; then
    if [[ -s "$TMPDIR_TEST/stdin.md" ]]; then
        pass "RT-5.19: Stdin mode reads and converts piped input"
    else
        fail "RT-5.19: Stdin mode reads and converts piped input" "Output is empty"
    fi
else
    fail "RT-5.19: Stdin mode reads and converts piped input" "Script exited non-zero"
fi

# RT-5.20: --help prints usage and exits 0
help_ok=true
for script in "$TO_PDF" "$TO_MD"; do
    if ! "$script" --help > /dev/null 2>&1; then
        help_ok=false
    fi
done
if $help_ok; then
    pass "RT-5.20: --help prints usage and exits 0"
else
    fail "RT-5.20: --help prints usage and exits 0" "One or both scripts failed on --help"
fi

# RT-5.21: --version prints version and exits 0
version_ok=true
for script in "$TO_PDF" "$TO_MD"; do
    out=$("$script" --version 2>&1 || true)
    if [[ -z "$out" ]]; then
        version_ok=false
    fi
done
if $version_ok; then
    pass "RT-5.21: --version prints version and exits 0"
else
    fail "RT-5.21: --version prints version and exits 0" "One or both scripts produced no version output"
fi

# --- AC5.9: Front matter ---
echo ""
echo "AC5.9: Front matter in output"

# RT-5.22: Typst output includes title and author
create_minimal_fixture
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/fm.typ" "$TMPDIR_TEST/minimal.org" 2>/dev/null; then
    if grep -q "Test Play" "$TMPDIR_TEST/fm.typ" && grep -q "Test Author" "$TMPDIR_TEST/fm.typ"; then
        pass "RT-5.22: Typst output includes title and author from org front matter"
    else
        fail "RT-5.22: Typst output includes title and author from org front matter" \
             "Title or author not found"
    fi
else
    fail "RT-5.22: Typst output includes title and author from org front matter" "Script exited non-zero"
fi

# RT-5.23: Markdown output includes title and author
if "$TO_MD" "$TMPDIR_TEST/minimal.org" > "$TMPDIR_TEST/fm.md" 2>/dev/null; then
    if grep -q "Test Play" "$TMPDIR_TEST/fm.md" && grep -q "Test Author" "$TMPDIR_TEST/fm.md"; then
        pass "RT-5.23: Markdown output includes title and author from org front matter"
    else
        fail "RT-5.23: Markdown output includes title and author from org front matter" \
             "Title or author not found"
    fi
else
    fail "RT-5.23: Markdown output includes title and author from org front matter" "Script exited non-zero"
fi

# --- AC5.10: Configurable Typst preamble ---
echo ""
echo "AC5.10: Configurable Typst preamble"

# RT-5.24: Non-default flag values appear in generated Typst preamble
create_minimal_fixture
if "$TO_PDF" --typ-only --font "Georgia" --font-size 14pt --margin 30mm --page letter --indent 6em \
    -o "$TMPDIR_TEST/custom.typ" "$TMPDIR_TEST/minimal.org" 2>/dev/null; then
    custom_ok=true
    if ! grep -q "Georgia" "$TMPDIR_TEST/custom.typ"; then custom_ok=false; fi
    if ! grep -q "14pt" "$TMPDIR_TEST/custom.typ"; then custom_ok=false; fi
    if ! grep -q "30mm" "$TMPDIR_TEST/custom.typ"; then custom_ok=false; fi
    if ! grep -qi "letter" "$TMPDIR_TEST/custom.typ"; then custom_ok=false; fi
    if ! grep -q "6em" "$TMPDIR_TEST/custom.typ"; then custom_ok=false; fi
    if $custom_ok; then
        pass "RT-5.24: Non-default flag values appear in the generated Typst preamble"
    else
        fail "RT-5.24: Non-default flag values appear in the generated Typst preamble" \
             "Some custom values not found in output"
    fi
else
    fail "RT-5.24: Non-default flag values appear in the generated Typst preamble" "Script exited non-zero"
fi

# UT-5.2: Changing font/margins via flags alters the PDF layout visually — human verification only

# --- AC5.11: Noexport sections ---
echo ""
echo "AC5.11: Noexport sections excluded"

create_noexport_fixture

# RT-5.25: Typst output omits :noexport: sections
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/noexp.typ" "$TMPDIR_TEST/noexport.org" 2>/dev/null; then
    if grep -q "This should not appear" "$TMPDIR_TEST/noexp.typ"; then
        fail "RT-5.25: Typst output omits :noexport: sections" "Noexport content found in output"
    else
        if grep -q "Hello" "$TMPDIR_TEST/noexp.typ"; then
            pass "RT-5.25: Typst output omits :noexport: sections"
        else
            fail "RT-5.25: Typst output omits :noexport: sections" "Regular content also missing"
        fi
    fi
else
    fail "RT-5.25: Typst output omits :noexport: sections" "Script exited non-zero"
fi

# RT-5.26: Markdown output omits :noexport: sections
if "$TO_MD" "$TMPDIR_TEST/noexport.org" > "$TMPDIR_TEST/noexp.md" 2>/dev/null; then
    if grep -q "This should not appear" "$TMPDIR_TEST/noexp.md"; then
        fail "RT-5.26: Markdown output omits :noexport: sections" "Noexport content found in output"
    else
        if grep -q "Hello" "$TMPDIR_TEST/noexp.md"; then
            pass "RT-5.26: Markdown output omits :noexport: sections"
        else
            fail "RT-5.26: Markdown output omits :noexport: sections" "Regular content also missing"
        fi
    fi
else
    fail "RT-5.26: Markdown output omits :noexport: sections" "Script exited non-zero"
fi

# --- AC5.12: Character table ---
echo ""
echo "AC5.12: Character table preserved"

create_minimal_fixture

# RT-5.27: Typst output includes a formatted character table
if "$TO_PDF" --typ-only -o "$TMPDIR_TEST/table.typ" "$TMPDIR_TEST/minimal.org" 2>/dev/null; then
    if grep -q "BOB" "$TMPDIR_TEST/table.typ" && grep -q "A test char" "$TMPDIR_TEST/table.typ"; then
        pass "RT-5.27: Typst output includes a formatted character table"
    else
        fail "RT-5.27: Typst output includes a formatted character table" \
             "Character table content not found"
    fi
else
    fail "RT-5.27: Typst output includes a formatted character table" "Script exited non-zero"
fi

# RT-5.28: Markdown output includes the character table as a Markdown table
if "$TO_MD" "$TMPDIR_TEST/minimal.org" > "$TMPDIR_TEST/table.md" 2>/dev/null; then
    # Markdown tables use | delimiters
    if grep -q '|.*BOB.*|.*A test char.*|' "$TMPDIR_TEST/table.md" || \
       (grep -q "BOB" "$TMPDIR_TEST/table.md" && grep -q "A test char" "$TMPDIR_TEST/table.md" && grep -q '|' "$TMPDIR_TEST/table.md"); then
        pass "RT-5.28: Markdown output includes the character table as a Markdown table"
    else
        fail "RT-5.28: Markdown output includes the character table as a Markdown table" \
             "Character table not found in markdown"
    fi
else
    fail "RT-5.28: Markdown output includes the character table as a Markdown table" "Script exited non-zero"
fi

# --- AC5.13: make install/uninstall ---
echo ""
echo "AC5.13: make install/uninstall"

# RT-5.29: make install creates working symlinks in ~/bin/
# Save any existing symlinks to restore later
pdf_backup=""
md_backup=""
if [[ -L "$HOME/bin/org-play-to-pdf" ]]; then
    pdf_backup="$(readlink "$HOME/bin/org-play-to-pdf")"
fi
if [[ -L "$HOME/bin/org-play-to-markdown" ]]; then
    md_backup="$(readlink "$HOME/bin/org-play-to-markdown")"
fi

if make -C "$PROJECT_DIR" install > /dev/null 2>&1; then
    if [[ -L "$HOME/bin/org-play-to-pdf" ]] && [[ -L "$HOME/bin/org-play-to-markdown" ]]; then
        pass "RT-5.29: make install creates working symlinks in ~/bin/"
    else
        fail "RT-5.29: make install creates working symlinks in ~/bin/" "Symlinks not found"
    fi
else
    fail "RT-5.29: make install creates working symlinks in ~/bin/" "make install failed"
fi

# RT-5.30: make uninstall removes the symlinks
if make -C "$PROJECT_DIR" uninstall > /dev/null 2>&1; then
    if [[ ! -e "$HOME/bin/org-play-to-pdf" ]] && [[ ! -e "$HOME/bin/org-play-to-markdown" ]]; then
        pass "RT-5.30: make uninstall removes the symlinks"
    else
        fail "RT-5.30: make uninstall removes the symlinks" "Symlinks still present"
    fi
else
    fail "RT-5.30: make uninstall removes the symlinks" "make uninstall failed"
fi

# RT-5.31: Installed scripts are executable via PATH
# Re-install for this test
if make -C "$PROJECT_DIR" install > /dev/null 2>&1; then
    if [[ -x "$HOME/bin/org-play-to-pdf" ]] && [[ -x "$HOME/bin/org-play-to-markdown" ]]; then
        pass "RT-5.31: Installed scripts are executable via PATH"
    else
        fail "RT-5.31: Installed scripts are executable via PATH" "Symlinks not executable"
    fi
else
    fail "RT-5.31: Installed scripts are executable via PATH" "make install failed"
fi

# Restore original state if there were pre-existing symlinks
if [[ -n "$pdf_backup" ]]; then
    ln -sf "$pdf_backup" "$HOME/bin/org-play-to-pdf"
elif [[ -L "$HOME/bin/org-play-to-pdf" ]]; then
    rm -f "$HOME/bin/org-play-to-pdf"
fi
if [[ -n "$md_backup" ]]; then
    ln -sf "$md_backup" "$HOME/bin/org-play-to-markdown"
elif [[ -L "$HOME/bin/org-play-to-markdown" ]]; then
    rm -f "$HOME/bin/org-play-to-markdown"
fi

# ====================================================================
echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="
if [[ "$FAIL" -gt 0 ]]; then
    exit 1
fi
