// ABOUTME: Typst rendering for manuscript documents.
// ABOUTME: Uses a file-backed template and generated block markup.
package manuscript

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type templateData struct {
	Config           Config
	Meta             Metadata
	Header           string
	Body             string
	IsUS             bool
	Leading          string
	Spacing          string
	PartVertical     string
	ChapterPosition  string
	SceneBreakMarker string
	HasContact       bool
}

func RenderTypst(doc Document, cfg Config) (string, error) {
	body, err := renderBlocks(doc.Blocks, cfg)
	if err != nil {
		return "", err
	}
	safeMeta := escapedMetadata(doc.Metadata)
	safeMeta.Date = escapeTypst(renderDate(doc.Metadata.Date, cfg.Folio.Manuscript.DateFormat))
	data := templateData{
		Config:           cfg,
		Meta:             safeMeta,
		Header:           renderHeader(doc.Metadata, cfg),
		Body:             body,
		IsUS:             cfg.Folio.Manuscript.Style == "us",
		Leading:          lineSpacingLeading(cfg.Folio.Manuscript.LineSpacing),
		Spacing:          paragraphSpacing(cfg.Folio.Manuscript.ParagraphSpacing, cfg.Folio.Manuscript.LineSpacing),
		PartVertical:     typstVerticalAlign(cfg.Folio.Manuscript.Part.VerticalAlign),
		ChapterPosition:  chapterPosition(cfg.Folio.Manuscript.Chapter.Position),
		SceneBreakMarker: escapeTypst(cfg.Folio.Manuscript.SceneBreak.Marker),
		HasContact:       hasContactBlock(doc.Metadata, cfg),
	}
	root, err := projectRoot()
	if err != nil {
		return "", err
	}
	tmplPath := filepath.Join(root, "templates", "manuscript.typ")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("parsing Typst template: %w", err)
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return "", fmt.Errorf("executing Typst template: %w", err)
	}
	return out.String(), nil
}

func renderDate(value string, layout string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	parsed, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return value
	}
	if strings.TrimSpace(layout) == "" {
		layout = "2 January 2006"
	}
	return parsed.Format(layout)
}

func lineSpacingLeading(lineSpacing string) string {
	multiplier, err := strconv.ParseFloat(strings.TrimSpace(lineSpacing), 64)
	if err != nil {
		return lineSpacing + "em"
	}
	if multiplier < 1 {
		multiplier = 1
	}
	return strconv.FormatFloat(multiplier, 'f', -1, 64) + "em"
}

func paragraphSpacing(spacing string, lineSpacing string) string {
	trimmed := strings.TrimSpace(spacing)
	if trimmed == "" || trimmed == "0" || trimmed == "0pt" {
		multiplier, err := strconv.ParseFloat(strings.TrimSpace(lineSpacing), 64)
		if err != nil {
			return lineSpacing + "em"
		}
		if multiplier < 1 {
			multiplier = 1
		}
		return strconv.FormatFloat(multiplier, 'f', -1, 64) + "em"
	}
	return trimmed
}

func renderBlocks(blocks []Block, cfg Config) (string, error) {
	var lines []string
	firstPageBlock := true
	for _, block := range blocks {
		switch block.Kind {
		case "part":
			lines = append(lines, fmt.Sprintf("#folio-part(first: %t)[%s]",
				firstPageBlock,
				caseTransform(block.Text, cfg.Folio.Manuscript.Part.CaseTransform)))
			firstPageBlock = false
		case "chapter":
			lines = append(lines, fmt.Sprintf("#folio-chapter(first: %t)[%s]",
				firstPageBlock,
				caseTransform(block.Text, cfg.Folio.Manuscript.Chapter.CaseTransform)))
			firstPageBlock = false
		case "section":
			lines = append(lines, "#folio-section["+typstInline(block.Text, cfg)+"]")
		case "paragraph":
			lines = append(lines, typstInline(block.Text, cfg))
		case "scene-break":
			lines = append(lines, "#folio-scene-break()")
		case "code":
			lines = append(lines, "#folio-code["+escapeTypst(block.Text)+"]")
		case "raw-typst":
			lines = append(lines, block.Text)
		case "footnote":
			continue
		default:
			return "", fmt.Errorf("unknown manuscript block kind: %s", block.Kind)
		}
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n"), nil
}

func renderHeader(meta Metadata, cfg Config) string {
	header := cfg.Folio.Manuscript.PageHeader.Format
	header = strings.ReplaceAll(header, "[author]", escapeTypst(meta.Author))
	header = strings.ReplaceAll(header, "[title]", escapeTypst(meta.Title))
	header = strings.ReplaceAll(header, "[page]", "#context counter(page).display()")
	return header
}

func hasContactBlock(meta Metadata, cfg Config) bool {
	titlePage := cfg.Folio.Manuscript.TitlePage
	return titlePage.IncludeContactName && meta.ContactName != "" ||
		titlePage.IncludeAddress && meta.Address != "" ||
		titlePage.IncludePhone && meta.Phone != "" ||
		titlePage.IncludeEmail && meta.Email != "" ||
		titlePage.IncludeWebsite && meta.Website != ""
}

func typstInline(text string, cfg Config) string {
	return renderInlineMarkup(text, cfg.Folio.Manuscript.MonoFont, cfg.Folio.Manuscript.MonoFontSize, cfg.Folio.Manuscript.MonoFontWeight)
}

func renderInlineMarkup(text string, monoFont string, monoSize string, monoWeight string) string {
	var out strings.Builder
	for i := 0; i < len(text); {
		switch {
		case strings.HasPrefix(text[i:], "`"):
			if end := strings.Index(text[i+1:], "`"); end >= 0 {
				content := text[i+1 : i+1+end]
				out.WriteString(fmt.Sprintf(`#text(font: "%s", size: %s, weight: "%s")[%s]`, escapeTypst(monoFont), monoSize, escapeTypst(monoWeight), escapeTypst(content)))
				i += end + 2
				continue
			}
		case strings.HasPrefix(text[i:], "**"):
			if end := strings.Index(text[i+2:], "**"); end >= 0 {
				content := text[i+2 : i+2+end]
				out.WriteString("*")
				out.WriteString(renderInlineMarkup(content, monoFont, monoSize, monoWeight))
				out.WriteString("*")
				i += end + 4
				continue
			}
		case strings.HasPrefix(text[i:], "*"):
			if end := strings.Index(text[i+1:], "*"); end >= 0 {
				content := text[i+1 : i+1+end]
				out.WriteString("_")
				out.WriteString(renderInlineMarkup(content, monoFont, monoSize, monoWeight))
				out.WriteString("_")
				i += end + 2
				continue
			}
		case strings.HasPrefix(text[i:], "[fn:"):
			if end := strings.Index(text[i:], "]"); end >= 0 {
				out.WriteString("#footnote[")
				out.WriteString(escapeTypst(text[i+4 : i+end]))
				out.WriteString("]")
				i += end + 1
				continue
			}
		case strings.HasPrefix(text[i:], "[^"):
			if end := strings.Index(text[i:], "]"); end >= 0 {
				out.WriteString("#footnote[")
				out.WriteString(escapeTypst(text[i+2 : i+end]))
				out.WriteString("]")
				i += end + 1
				continue
			}
		case strings.HasPrefix(text[i:], "---"):
			out.WriteString("—")
			i += 3
			continue
		case strings.HasPrefix(text[i:], "--"):
			out.WriteString("–")
			i += 2
			continue
		}

		next := nextInlineMarker(text[i+1:])
		if next < 0 {
			out.WriteString(escapeTypst(applyMarkdownDashes(text[i:])))
			break
		}
		next += i + 1
		out.WriteString(escapeTypst(applyMarkdownDashes(text[i:next])))
		i = next
	}
	return out.String()
}

func nextInlineMarker(text string) int {
	index := -1
	for _, marker := range []string{"`", "**", "*", "[fn:", "[^", "---", "--"} {
		if found := strings.Index(text, marker); found >= 0 && (index < 0 || found < index) {
			index = found
		}
	}
	return index
}

func applyMarkdownDashes(text string) string {
	text = strings.ReplaceAll(text, "---", "—")
	return strings.ReplaceAll(text, "--", "–")
}

func typstVerticalAlign(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "center", "middle", "horizon", "":
		return "horizon"
	case "top":
		return "top"
	case "bottom":
		return "bottom"
	default:
		return value
	}
}

func chapterPosition(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "one-third", "third", "":
		return "30%"
	case "center", "middle":
		return "50%"
	case "top":
		return "0%"
	default:
		return value
	}
}

func escapeTypst(text string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`[`, `\[`,
		`]`, `\]`,
		`#`, `\#`,
		`$`, `\$`,
		`@`, `\@`,
		`_`, `\_`,
		`*`, `\*`,
		`/`, `\/`,
	)
	return replacer.Replace(text)
}

func caseTransform(text string, mode string) string {
	if mode == "upper" {
		return escapeTypst(strings.ToUpper(text))
	}
	return typstInline(text, Config{})
}

func escapedMetadata(meta Metadata) Metadata {
	return Metadata{
		Title:             escapeTypst(meta.Title),
		Subtitle:          escapeTypst(meta.Subtitle),
		Author:            escapeTypst(meta.Author),
		AuthorAttribution: escapeTypst(meta.AuthorAttribution),
		Date:              escapeTypst(meta.Date),
		Version:           escapeTypst(meta.Version),
		WordCount:         escapeTypst(meta.WordCount),
		ContactName:       escapeTypst(meta.ContactName),
		Address:           escapeTypst(meta.Address),
		Phone:             escapeTypst(meta.Phone),
		Email:             escapeTypst(meta.Email),
		Website:           escapeTypst(meta.Website),
	}
}

func templateExistsForTests() bool {
	root, err := projectRoot()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(root, "templates", "manuscript.typ"))
	return err == nil
}
