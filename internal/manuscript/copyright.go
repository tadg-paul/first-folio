// ABOUTME: Renders the copyright frontmatter page (issue #21) from declarative
// ABOUTME: config into a single Typst content block for template insertion.
package manuscript

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// renderCopyrightPage composes the entire copyright-page emit as a single Typst
// content block: blank-page-before directive, skip-header/skip-footer state
// updates, alignment wrapper, credit blocks, body paragraphs, separator,
// publication lines, publisher, ISBN, and (optionally) barcode. Returns the
// empty string when the copyright block is not enabled.
func renderCopyrightPage(meta Metadata, cfg Config) string {
	c := cfg.Folio.Manuscript.Copyright
	if !c.Enabled {
		return ""
	}
	var b strings.Builder

	// Blank-page-before directive (enforce-left / enforce-right / true / false).
	if d := c.BlankPageBefore.TypstDirective(); d != "" {
		b.WriteString(d)
		b.WriteString("\n")
	}

	// Skip-header / skip-footer: record this page in the appropriate skip lists so
	// the running header/footer context suppresses them here. Runs inside a context
	// block because we need the current physical page number.
	if (c.SkipHeader != nil && *c.SkipHeader) || (c.SkipFooter != nil && *c.SkipFooter) {
		b.WriteString("#context {\n")
		b.WriteString("  let pg = counter(page).at(here()).first()\n")
		if c.SkipHeader != nil && *c.SkipHeader {
			b.WriteString(`  state("folio-skip-header-pages", ()).update(pages => pages + (pg,))` + "\n")
		}
		if c.SkipFooter != nil && *c.SkipFooter {
			b.WriteString(`  state("folio-skip-footer-pages", ()).update(pages => pages + (pg,))` + "\n")
		}
		b.WriteString("}\n")
	}

	// Content wrapper: font + alignment + block layout.
	b.WriteString(fmt.Sprintf("#text(font: %q, size: %s)[\n", c.Font, orFallback(c.FontSize, "10pt")))
	b.WriteString(fmt.Sprintf("#set par(leading: %s)\n", copyrightLeading(c.LineSpacing)))
	b.WriteString(fmt.Sprintf("#align(%s)[\n", copyrightAlignExpr(c.Align)))

	blocks := composeCopyrightBlocks(meta, cfg, c)
	for i, block := range blocks {
		if i > 0 {
			b.WriteString(fmt.Sprintf("#v(%s)\n", c.BlockSpacing))
		}
		b.WriteString(block)
		b.WriteString("\n")
	}

	b.WriteString("]\n") // close align
	b.WriteString("]\n") // close text

	// Blank-page-after directive.
	if d := c.BlankPageAfter.TypstDirective(); d != "" {
		b.WriteString(d)
		b.WriteString("\n")
	}

	// End of copyright page: hard pagebreak so the following block starts fresh.
	b.WriteString("#pagebreak()\n")

	return b.String()
}

// composeCopyrightBlocks returns the ordered non-empty content blocks that make up
// the copyright page body: credits, body paragraphs, separator, publication,
// publisher, ISBN, barcode. Empty blocks are omitted so block-spacing collapses.
func composeCopyrightBlocks(meta Metadata, cfg Config, c CopyrightConfig) []string {
	var out []string

	// Credit blocks.
	credits := effectiveCredits(meta, c)
	for _, credit := range credits {
		if len(credit.Holders) == 0 {
			continue
		}
		out = append(out, renderCreditBlock(credit, c))
	}

	// Body paragraphs.
	for _, para := range c.Body {
		trimmed := strings.TrimSpace(para)
		if trimmed == "" {
			continue
		}
		out = append(out, renderBodyParagraph(trimmed))
	}

	// Separator.
	if c.Separator != "" {
		out = append(out, fmt.Sprintf("#v(%s)\n%s\n#v(%s)",
			c.SeparatorSpaceBefore,
			escapeTypst(c.Separator),
			c.SeparatorSpaceAfter,
		))
	}

	// Publication lines.
	if len(c.Publication) > 0 {
		var pub strings.Builder
		for i, line := range c.Publication {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			if i > 0 {
				pub.WriteString(" \\\n")
			}
			pub.WriteString(escapeTypst(trimmed))
		}
		if pub.Len() > 0 {
			out = append(out, pub.String())
		}
	}

	// Publisher.
	if c.Publisher != "" {
		out = append(out, fmt.Sprintf("%s #text(weight: %q)[%s]",
			escapeTypst(c.PublisherPreposition),
			c.HeadingFontWeight,
			escapeTypst(c.Publisher),
		))
	}

	// ISBN (label bold, value regular).
	if c.ISBN != "" {
		out = append(out, fmt.Sprintf("#text(weight: %q)[%s]: %s",
			c.HeadingFontWeight,
			escapeTypst(c.ISBNLabel),
			escapeTypst(c.ISBN),
		))
	}

	// Barcode (render or render-and-file). SVG is embedded as inline Typst image.
	if c.ISBNBarcode == "render" || c.ISBNBarcode == "render-and-file" {
		if svg, err := renderEAN13SVG(c.ISBN); err == nil {
			// Typst image() accepts SVG data via bytes; embed as a UTF-8 string.
			out = append(out, fmt.Sprintf("#image(bytes(%q), format: \"svg\", width: 40mm)", svg))
		}
	}

	return out
}

// effectiveCredits returns the user-configured credits list, or a single default
// entry ({Copyright, folio.date year, [folio.author]}) if the list is empty.
func effectiveCredits(meta Metadata, c CopyrightConfig) []CopyrightCredit {
	if len(c.Credits) > 0 {
		// Fill in defaults on each entry.
		primaryYear := ""
		for i := range c.Credits {
			if c.Credits[i].Year == "" && primaryYear != "" {
				c.Credits[i].Year = primaryYear
			}
			if c.Credits[i].Year == "" {
				c.Credits[i].Year = deriveYear(meta.Date)
				primaryYear = c.Credits[i].Year
			} else if primaryYear == "" {
				primaryYear = c.Credits[i].Year
			}
		}
		return c.Credits
	}
	// Default single-author copyright.
	year := deriveYear(meta.Date)
	holders := []string{}
	if meta.Author != "" {
		holders = append(holders, meta.Author)
	}
	if len(holders) == 0 {
		return nil
	}
	return []CopyrightCredit{{
		Heading: "Copyright",
		Year:    year,
		Holders: holders,
	}}
}

// deriveYear returns the 4-digit year parsed from meta.Date (YYYY-MM-DD form).
// meta.Date is guaranteed non-empty by applyMetadataOverrides (defaults to today).
func deriveYear(date string) string {
	if len(date) >= 4 {
		return date[:4]
	}
	return strconv.Itoa(time.Now().Year())
}

// renderCreditBlock returns "<heading> © YEAR" (bold heading) followed by the
// comma-joined holders with a trailing full stop.
func renderCreditBlock(credit CopyrightCredit, c CopyrightConfig) string {
	var b strings.Builder
	headingLine := credit.Heading
	if credit.Year != "" {
		headingLine = fmt.Sprintf("%s © %s", credit.Heading, credit.Year)
	}
	b.WriteString(fmt.Sprintf("#text(weight: %q)[%s]", c.HeadingFontWeight, escapeTypst(headingLine)))
	b.WriteString(" \\\n")
	names := make([]string, 0, len(credit.Holders))
	for _, h := range credit.Holders {
		trimmed := strings.TrimSpace(h)
		if trimmed == "" {
			continue
		}
		names = append(names, escapeTypst(trimmed))
	}
	if len(names) > 0 {
		b.WriteString(strings.Join(names, ", "))
		b.WriteString(".")
	}
	return b.String()
}

// renderBodyParagraph converts a markdown-mini body string to a Typst content
// paragraph, translating **bold** -> *bold*, *italic* -> _italic_, --- -> em-dash,
// -- -> en-dash. Emitted as a single paragraph -- callers apply block-spacing.
func renderBodyParagraph(md string) string {
	return markdownMiniToTypst(md)
}

var (
	markdownBoldRE = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	markdownItalicRE = regexp.MustCompile(`(?:^|[^*])(\*[^*]+\*)(?:[^*]|$)`)
)

// markdownMiniToTypst converts a small subset of Markdown inline markup to Typst
// syntax. This intentionally keeps the surface small and predictable: bold,
// italic, en-dash, em-dash. Existing Typst inline syntax in the input passes
// through unchanged (double-conversion is idempotent for this subset).
func markdownMiniToTypst(md string) string {
	s := md
	// Em-dash / en-dash first so they don't get eaten by italic detection.
	s = strings.ReplaceAll(s, "---", "—")
	s = strings.ReplaceAll(s, "--", "–")
	// **bold** -> sentinel form so the italic pass below doesn't re-interpret the
	// single asterisks we just produced. We restore to Typst *bold* at the end.
	const (
		boldOpen  = "\x00BOLD_OPEN\x00"
		boldClose = "\x00BOLD_CLOSE\x00"
	)
	s = markdownBoldRE.ReplaceAllString(s, boldOpen+`$1`+boldClose)
	// Remaining single-* pairs are Markdown italic. Convert to Typst _italic_.
	s = convertMarkdownItalic(s)
	// Restore bold sentinels to Typst *bold*.
	s = strings.ReplaceAll(s, boldOpen, "*")
	s = strings.ReplaceAll(s, boldClose, "*")
	return s
}

// convertMarkdownItalic pairs up single-asterisks left after **bold** conversion
// and rewrites them as _italic_ Typst emphasis. Single-asterisk sequences at word
// boundaries only (avoids catching e.g. arithmetic 2*3 in copyright body).
func convertMarkdownItalic(s string) string {
	var out strings.Builder
	inItalic := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '*' {
			// Word-boundary check: opening italic requires a non-word char before
			// (or start), and non-space after; closing italic requires non-space
			// before and non-word after (or end).
			if !inItalic {
				// Opening: prev is start / non-word, next is non-space.
				prevOK := i == 0 || !isWordChar(s[i-1])
				nextOK := i+1 < len(s) && s[i+1] != ' '
				if prevOK && nextOK {
					out.WriteByte('_')
					inItalic = true
					continue
				}
			} else {
				// Closing: prev is non-space, next is end / non-word.
				prevOK := i > 0 && s[i-1] != ' '
				nextOK := i+1 == len(s) || !isWordChar(s[i+1])
				if prevOK && nextOK {
					out.WriteByte('_')
					inItalic = false
					continue
				}
			}
		}
		out.WriteByte(ch)
	}
	return out.String()
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// copyrightAlignExpr converts a compass keyword to a Typst align expression.
func copyrightAlignExpr(align string) string {
	switch strings.ToLower(strings.TrimSpace(align)) {
	case "left":
		return "left"
	case "right":
		return "right"
	case "top-center", "top-centre":
		return "top + center"
	case "bottom-center", "bottom-centre":
		return "bottom + center"
	default:
		return "center"
	}
}

func copyrightLeading(spacing string) string {
	if strings.TrimSpace(spacing) == "" {
		return "0.65em"
	}
	// spacing here is a Typst-compatible length (e.g. "1.15em"); convert to a
	// leading value (multiplier - 1)em when we detect a bare multiplier.
	if v, err := strconv.ParseFloat(strings.TrimSpace(spacing), 64); err == nil {
		return fmt.Sprintf("%gem", v-1)
	}
	// Otherwise pass through (already an em/pt length).
	return spacing
}

func orFallback(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}
