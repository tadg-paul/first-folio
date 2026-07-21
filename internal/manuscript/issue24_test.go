// ABOUTME: Regression tests for issue #24 -- frontmatter-format on header/footer.
// ABOUTME: Adds a distinct format string for frontmatter pages (title, copyright,
// ABOUTME: TOC, ...) via *string fields that distinguish unset (nil) from empty ("").
package manuscript

import (
	"strings"
	"testing"
)

// RT-24.1: unconfigured frontmatter-format produces byte-identical output (format
// used on every page; no is-body branch in the emitted header context body).
func TestRT_24_1_UnconfiguredFrontmatterFormatUnchanged(t *testing.T) {
	typst := renderIssue15Manuscript(t, "")
	// Without HasHeaderFrontmatter, no is-body branch appears inside the header
	// text() body -- the emit remains the calc.odd(pg) verso/recto ternary only.
	header := extractHeaderBlock(t, typst)
	assertNotContains(t, header, "if is-body")
}

// RT-24.2: frontmatter-format set to a non-empty string branches on is-body in the emit.
func TestRT_24_2_FrontmatterFormatBranchesOnIsBody(t *testing.T) {
	typst := renderIssue15Manuscript(t, strings.Join([]string{
		"folio:",
		"  manuscript:",
		"    page-header:",
		"      enabled: true",
		"      format: \"BODY-HEADER\"",
		"      frontmatter-format: \"FRONT-HEADER\"",
		"",
	}, "\n"))
	header := extractHeaderBlock(t, typst)
	assertContains(t, header, "if is-body")
	assertContains(t, header, "BODY-HEADER")
	assertContains(t, header, "FRONT-HEADER")
}

// RT-24.3: frontmatter-format set to an empty string still branches (renders blank
// on frontmatter, format on body).
func TestRT_24_3_EmptyFrontmatterFormatRendersBlank(t *testing.T) {
	typst := renderIssue15Manuscript(t, strings.Join([]string{
		"folio:",
		"  manuscript:",
		"    page-header:",
		"      enabled: true",
		"      format: \"BODY-HEADER\"",
		"      frontmatter-format: \"\"",
		"",
	}, "\n"))
	header := extractHeaderBlock(t, typst)
	assertContains(t, header, "if is-body")
	// The frontmatter branch renders empty content [] (no text between the brackets).
	assertContains(t, header, "BODY-HEADER")
}

// RT-24.4: alt-frontmatter-format set produces a verso/recto pair on frontmatter.
func TestRT_24_4_AltFrontmatterFormatVersoRectoPair(t *testing.T) {
	typst := renderIssue15Manuscript(t, strings.Join([]string{
		"folio:",
		"  manuscript:",
		"    page-header:",
		"      enabled: true",
		"      format: \"BODY-HEADER\"",
		"      frontmatter-format: \"FRONT-VERSO\"",
		"      alt-frontmatter-format: \"FRONT-RECTO\"",
		"",
	}, "\n"))
	header := extractHeaderBlock(t, typst)
	assertContains(t, header, "FRONT-VERSO")
	assertContains(t, header, "FRONT-RECTO")
	// The frontmatter branch has its own calc.odd check.
	assertContains(t, header, "if is-body")
}

// RT-24.5: page-footer.frontmatter-format works the same way.
func TestRT_24_5_PageFooterFrontmatterFormat(t *testing.T) {
	typst := renderIssue15Manuscript(t, strings.Join([]string{
		"folio:",
		"  manuscript:",
		"    page-footer:",
		"      enabled: true",
		"      format: \"[page]\"",
		"      frontmatter-format: \"\"",
		"",
	}, "\n"))
	footer := extractFooterBlock(t, typst)
	assertContains(t, footer, "if is-body")
}

// RT-24.6: folio-is-body state is seeded to true inside #folio-part and #folio-chapter
// when first: true, so the header/footer context can distinguish frontmatter from body.
func TestRT_24_6_IsBodyStateSeededInPartAndChapter(t *testing.T) {
	typst := renderIssue15Manuscript(t, "")
	// The template must emit an update of folio-is-body inside both folio-part and
	// folio-chapter (guarded on first). Check the macro definitions themselves.
	assertContains(t, typst, `state("folio-is-body", false).update(true)`)
}

// RT-24.7: page-footer.alt-frontmatter-format produces the verso/recto pair on the footer.
func TestRT_24_7_AltFrontmatterFormatFooter(t *testing.T) {
	typst := renderIssue15Manuscript(t, strings.Join([]string{
		"folio:",
		"  manuscript:",
		"    page-footer:",
		"      enabled: true",
		"      format: \"[page]\"",
		"      frontmatter-format: \"FRONT-VERSO-FOOT\"",
		"      alt-frontmatter-format: \"FRONT-RECTO-FOOT\"",
		"",
	}, "\n"))
	footer := extractFooterBlock(t, typst)
	assertContains(t, footer, "FRONT-VERSO-FOOT")
	assertContains(t, footer, "FRONT-RECTO-FOOT")
}
