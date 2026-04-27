#!/usr/bin/env bash
# ABOUTME: Format conversion matrix tests — every source→target pair.
# ABOUTME: Verifies structural elements survive each conversion path.
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

# --- Source fixtures in each format ---

cat > "$TMPDIR_TEST/source.org" <<'ORG'
#+TITLE: Matrix Test
#+AUTHOR: Test Author
#+SUBTITLE: A Subtitle

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
*** BOB crosses to the window.
**** CÁIT softly
Goodbye now.
***** FADE OUT.
ORG

cat > "$TMPDIR_TEST/source.md" <<'MD'
# Matrix Test

**A Subtitle**

*by Test Author*

| Character | Description |
|-----------|-------------|
| BOB       | A test char |
| CÁIT      | Another     |

## Act I

### Scene 1

*A bare stage.*

**BOB:**
Hello there.

*BOB crosses to the window.*

**CÁIT:** *(softly)*
Goodbye now.

> FADE OUT.
MD

cat > "$TMPDIR_TEST/source.fountain" <<'FTN'
Title: Matrix Test
Author: Test Author

> **ACT I** <

.Scene 1

A bare stage.

BOB
Hello there.

BOB crosses to the window.

CÁIT
(softly)
Goodbye now.

FADE OUT.
FTN

# --- Verification helpers ---

check_org() {
    local file="$1" desc="$2"
    local has_title=false has_act=false has_scene=false has_dir=false has_char=false has_dial=false
    if grep -q '#+TITLE' "$file"; then has_title=true; fi
    if grep -qi '^\* Act I' "$file"; then has_act=true; fi
    if grep -qi '^\*\* Scene 1' "$file"; then has_scene=true; fi
    if grep -q '^\*\*\* ' "$file"; then has_dir=true; fi
    if grep -q '^\*\*\*\* BOB' "$file"; then has_char=true; fi
    if grep -q 'Hello there' "$file"; then has_dial=true; fi
    if $has_title && $has_act && $has_scene && $has_dir && $has_char && $has_dial; then
        pass "$desc"
    else
        fail "$desc" "title=$has_title act=$has_act scene=$has_scene dir=$has_dir char=$has_char dial=$has_dial"
    fi
}

check_md() {
    local file="$1" desc="$2"
    local has_title=false has_act=false has_scene=false has_dir=false has_char=false has_dial=false
    if grep -q '^# ' "$file"; then has_title=true; fi
    if grep -qi '^## Act I' "$file"; then has_act=true; fi
    if grep -q '^### Scene 1' "$file"; then has_scene=true; fi
    if grep -q '^\*.*bare stage.*\*$' "$file"; then has_dir=true; fi
    if grep -q '^\*\*BOB:\*\*' "$file"; then has_char=true; fi
    if grep -q 'Hello there' "$file"; then has_dial=true; fi
    if $has_title && $has_act && $has_scene && $has_dir && $has_char && $has_dial; then
        pass "$desc"
    else
        fail "$desc" "title=$has_title act=$has_act scene=$has_scene dir=$has_dir char=$has_char dial=$has_dial"
    fi
}

check_fountain() {
    local file="$1" desc="$2"
    local has_title=false has_act=false has_scene=false has_char=false has_dial=false
    if grep -q '^Title:' "$file"; then has_title=true; fi
    if grep -qi 'ACT I' "$file"; then has_act=true; fi
    if grep -q '^\.' "$file"; then has_scene=true; fi
    if grep -q '^BOB$' "$file"; then has_char=true; fi
    if grep -q 'Hello there' "$file"; then has_dial=true; fi
    if $has_title && $has_act && $has_scene && $has_char && $has_dial; then
        pass "$desc"
    else
        fail "$desc" "title=$has_title act=$has_act scene=$has_scene char=$has_char dial=$has_dial"
    fi
}

check_pdf() {
    local file="$1" desc="$2"
    if [[ -s "$file" ]] && head -c 5 "$file" | grep -q '%PDF'; then
        pass "$desc"
    else
        fail "$desc" "not a valid PDF"
    fi
}

# ====================================================================
echo "=== Format conversion matrix tests ==="

# --- Org source ---
echo ""
echo "Org → all targets"

"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/org2md.md" 2>/dev/null || true
check_md "$TMPDIR_TEST/org2md.md" "RT-3.60: org→md preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/org2ftn.fountain" 2>/dev/null || true
check_fountain "$TMPDIR_TEST/org2ftn.fountain" "RT-3.61: org→fountain preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/org2org.org" 2>/dev/null || true
check_org "$TMPDIR_TEST/org2org.org" "RT-3.62: org→org round-trip preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/org2pdf.pdf" 2>/dev/null || true
check_pdf "$TMPDIR_TEST/org2pdf.pdf" "RT-3.63: org→pdf produces valid PDF"

# --- Markdown source ---
echo ""
echo "Markdown → all targets"

"$FOLIO" convert "$TMPDIR_TEST/source.md" "$TMPDIR_TEST/md2org.org" 2>/dev/null || true
check_org "$TMPDIR_TEST/md2org.org" "RT-3.64: md→org preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.md" "$TMPDIR_TEST/md2ftn.fountain" 2>/dev/null || true
check_fountain "$TMPDIR_TEST/md2ftn.fountain" "RT-3.65: md→fountain preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.md" "$TMPDIR_TEST/md2md.md" 2>/dev/null || true
check_md "$TMPDIR_TEST/md2md.md" "RT-3.66: md→md round-trip preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.md" "$TMPDIR_TEST/md2pdf.pdf" 2>/dev/null || true
check_pdf "$TMPDIR_TEST/md2pdf.pdf" "RT-3.67: md→pdf produces valid PDF"

# --- Fountain source ---
echo ""
echo "Fountain → all targets"

"$FOLIO" convert "$TMPDIR_TEST/source.fountain" "$TMPDIR_TEST/ftn2org.org" 2>/dev/null || true
check_org "$TMPDIR_TEST/ftn2org.org" "RT-3.68: fountain→org preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.fountain" "$TMPDIR_TEST/ftn2md.md" 2>/dev/null || true
check_md "$TMPDIR_TEST/ftn2md.md" "RT-3.69: fountain→md preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.fountain" "$TMPDIR_TEST/ftn2ftn.fountain" 2>/dev/null || true
check_fountain "$TMPDIR_TEST/ftn2ftn.fountain" "RT-3.70: fountain→fountain round-trip preserves structure"

"$FOLIO" convert "$TMPDIR_TEST/source.fountain" "$TMPDIR_TEST/ftn2pdf.pdf" 2>/dev/null || true
check_pdf "$TMPDIR_TEST/ftn2pdf.pdf" "RT-3.71: fountain→pdf produces valid PDF"

# --- Unicode preservation ---
echo ""
echo "Unicode preservation across formats"

"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/unicode.md" 2>/dev/null || true
if grep -q 'CÁIT' "$TMPDIR_TEST/unicode.md"; then
    pass "RT-3.72: org→md preserves Unicode (CÁIT)"
else
    fail "RT-3.72: org→md preserves Unicode (CÁIT)" "CÁIT not found"
fi

"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/unicode.fountain" 2>/dev/null || true
if grep -q 'CÁIT' "$TMPDIR_TEST/unicode.fountain"; then
    pass "RT-3.73: org→fountain preserves Unicode (CÁIT)"
else
    fail "RT-3.73: org→fountain preserves Unicode (CÁIT)" "CÁIT not found"
fi

"$FOLIO" convert "$TMPDIR_TEST/source.fountain" "$TMPDIR_TEST/unicode2.org" 2>/dev/null || true
if grep -q 'CÁIT' "$TMPDIR_TEST/unicode2.org"; then
    pass "RT-3.74: fountain→org preserves Unicode (CÁIT)"
else
    fail "RT-3.74: fountain→org preserves Unicode (CÁIT)" "CÁIT not found"
fi

# --- Round-trip fidelity ---
echo ""
echo "Round-trip fidelity"

# org → fountain → org
"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/rt_step1.fountain" 2>/dev/null || true
"$FOLIO" convert "$TMPDIR_TEST/rt_step1.fountain" "$TMPDIR_TEST/rt_step2.org" 2>/dev/null || true
check_org "$TMPDIR_TEST/rt_step2.org" "RT-3.75: org→fountain→org round-trip"

# org → md → org
"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/rt_md1.md" 2>/dev/null || true
"$FOLIO" convert "$TMPDIR_TEST/rt_md1.md" "$TMPDIR_TEST/rt_md2.org" 2>/dev/null || true
check_org "$TMPDIR_TEST/rt_md2.org" "RT-3.76: org→md→org round-trip"

# fountain → md → fountain
"$FOLIO" convert "$TMPDIR_TEST/source.fountain" "$TMPDIR_TEST/rt_ftn1.md" 2>/dev/null || true
"$FOLIO" convert "$TMPDIR_TEST/rt_ftn1.md" "$TMPDIR_TEST/rt_ftn2.fountain" 2>/dev/null || true
check_fountain "$TMPDIR_TEST/rt_ftn2.fountain" "RT-3.77: fountain→md→fountain round-trip"

# --- Transition preservation ---
echo ""
echo "Transition across formats"

"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/trans.md" 2>/dev/null || true
if grep -q '> FADE OUT' "$TMPDIR_TEST/trans.md"; then
    pass "RT-3.78: org→md preserves transition as blockquote"
else
    fail "RT-3.78: org→md preserves transition" "not found"
fi

"$FOLIO" convert "$TMPDIR_TEST/source.org" "$TMPDIR_TEST/trans.fountain" 2>/dev/null || true
if grep -qi 'FADE OUT' "$TMPDIR_TEST/trans.fountain"; then
    pass "RT-3.79: org→fountain preserves transition"
else
    fail "RT-3.79: org→fountain preserves transition" "not found"
fi

"$FOLIO" convert "$TMPDIR_TEST/source.md" "$TMPDIR_TEST/trans2.org" 2>/dev/null || true
if grep -q '^\*\*\*\*\* FADE OUT' "$TMPDIR_TEST/trans2.org"; then
    pass "RT-3.80: md→org preserves transition as H5"
else
    fail "RT-3.80: md→org preserves transition as H5" "not found"
fi

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
