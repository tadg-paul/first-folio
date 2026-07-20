// ABOUTME: Regression tests for issue #23 -- title/TOC pages must honour gutter.
// ABOUTME: The document-top `#set page(...)` previously used a flat margin scalar
// ABOUTME: which meant title page and TOC ignored `folio.manuscript.gutter`.
package manuscript

import (
	"strings"
	"testing"
)

// RT-23.1: with a non-zero gutter, the document-top `#set page(...)` emits the
// parity-aware inside/outside/rest form so title and TOC pages honour the gutter.
func TestRT_23_1_DocumentTopPageHonoursGutter(t *testing.T) {
	typst := renderIssue15Manuscript(t, strings.Join([]string{
		"folio:",
		"  manuscript:",
		"    margin: 20mm",
		"    gutter: 15mm",
		"",
	}, "\n"))
	// The document-top #set page(...) (which title and TOC inherit from) must use
	// the parity-aware margin form when gutter is active.
	preTOC := extractDocumentTopPageBlock(t, typst)
	assertContains(t, preTOC, `inside: 20mm + 15mm`)
	assertContains(t, preTOC, `outside: 20mm`)
	assertContains(t, preTOC, `rest: 20mm`)
	assertNotContains(t, preTOC, `margin: 20mm,`)
}

// RT-23.2: with gutter unset (default 0mm), the document-top `#set page(...)`
// emits the flat scalar margin form (backwards-compatible byte-for-byte).
func TestRT_23_2_DocumentTopPageFlatMarginWithoutGutter(t *testing.T) {
	typst := renderIssue15Manuscript(t, strings.Join([]string{
		"folio:",
		"  manuscript:",
		"    margin: 20mm",
		"",
	}, "\n"))
	preTOC := extractDocumentTopPageBlock(t, typst)
	assertContains(t, preTOC, `margin: 20mm,`)
	assertNotContains(t, preTOC, `inside:`)
	assertNotContains(t, preTOC, `outside:`)
}

// RT-23.3: with gutter 0mm explicitly configured, same as unset -- flat margin.
func TestRT_23_3_DocumentTopPageExplicitZeroGutter(t *testing.T) {
	typst := renderIssue15Manuscript(t, strings.Join([]string{
		"folio:",
		"  manuscript:",
		"    margin: 20mm",
		"    gutter: 0mm",
		"",
	}, "\n"))
	preTOC := extractDocumentTopPageBlock(t, typst)
	assertContains(t, preTOC, `margin: 20mm,`)
	assertNotContains(t, preTOC, `inside:`)
}

// extractDocumentTopPageBlock returns the substring of the generated Typst covering
// the first #set page(...) call (which title page and TOC inherit from), stopping
// before the body-page setup marker.
func extractDocumentTopPageBlock(t *testing.T, typst string) string {
	t.Helper()
	bodyMarker := `// counter(page) is intentionally NOT reset here.`
	end := strings.Index(typst, bodyMarker)
	if end == -1 {
		t.Fatalf("body-page marker not found in generated Typst")
	}
	return typst[:end]
}
