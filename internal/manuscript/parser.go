// ABOUTME: Manuscript Markdown and org-mode parsers.
// ABOUTME: Produces prose manuscript blocks rather than stage-play events.
package manuscript

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func Parse(format string, text string) (Document, error) {
	switch format {
	case "markdown":
		return parseMarkdown(text)
	case "org":
		return parseOrg(text), nil
	default:
		return Document{}, fmt.Errorf("unsupported manuscript format: %s", format)
	}
}

func parseMarkdown(text string) (Document, error) {
	var doc Document
	var paragraph []string
	inNoExport := false
	noExportLevel := 0
	inCode := false
	var codeLines []string
	inMetadataTable := false
	metadataOpen := true

	var err error
	text, err = parseMarkdownFrontmatter(&doc.Metadata, text)
	if err != nil {
		return Document{}, err
	}
	if text != "" && doc.Metadata != (Metadata{}) {
		metadataOpen = false
	}

	flushParagraph := func() {
		if len(paragraph) > 0 && !inNoExport {
			doc.Blocks = append(doc.Blocks, Block{Kind: "paragraph", Text: strings.Join(paragraph, " ")})
		}
		paragraph = nil
	}
	flushCode := func() {
		if len(codeLines) > 0 && !inNoExport {
			doc.Blocks = append(doc.Blocks, Block{Kind: "code", Text: strings.Join(codeLines, "\n")})
		}
		codeLines = nil
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")

		if strings.HasPrefix(line, "```") {
			flushParagraph()
			if inCode {
				flushCode()
				inCode = false
				continue
			}
			inCode = true
			continue
		}
		if inCode {
			codeLines = append(codeLines, line)
			continue
		}

		level, heading := markdownHeading(line)
		if level > 0 {
			flushParagraph()
			metadataOpen = false
			if strings.Contains(heading, "<!-- noexport -->") {
				inNoExport = true
				noExportLevel = level
				continue
			}
			if inNoExport && level <= noExportLevel {
				inNoExport = false
			}
			if inNoExport {
				continue
			}
			heading = strings.TrimSpace(strings.ReplaceAll(heading, "<!-- noexport -->", ""))
			addMarkdownHeading(&doc, level, heading)
			continue
		}

		if inNoExport {
			continue
		}
		if strings.TrimSpace(line) == "" {
			flushParagraph()
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "<!--") {
			continue
		}
		if metadataOpen && parseMarkdownMetadataTable(&doc.Metadata, line, &inMetadataTable) {
			flushParagraph()
			continue
		}
		if !strings.HasPrefix(strings.TrimSpace(line), "|") {
			inMetadataTable = false
		}
		if strings.TrimSpace(line) == "***" {
			flushParagraph()
			metadataOpen = false
			doc.Blocks = append(doc.Blocks, Block{Kind: "scene-break"})
			continue
		}
		if name, value, ok := markdownFootnote(line); ok {
			flushParagraph()
			metadataOpen = false
			doc.Blocks = append(doc.Blocks, Block{Kind: "footnote", Text: name + "\t" + value})
			continue
		}
		if metadataOpen && parseMarkdownMetadata(&doc.Metadata, line) {
			continue
		}
		metadataOpen = false
		paragraph = append(paragraph, strings.TrimSpace(line))
	}
	flushParagraph()
	flushCode()
	return doc, nil
}

func parseOrg(text string) Document {
	var doc Document
	var paragraph []string
	inNoExport := false
	noExportLevel := 0
	inCode := false
	var codeLines []string

	flushParagraph := func() {
		if len(paragraph) > 0 && !inNoExport {
			doc.Blocks = append(doc.Blocks, Block{Kind: "paragraph", Text: strings.Join(paragraph, " ")})
		}
		paragraph = nil
	}
	flushCode := func() {
		if len(codeLines) > 0 && !inNoExport {
			doc.Blocks = append(doc.Blocks, Block{Kind: "code", Text: strings.Join(codeLines, "\n")})
		}
		codeLines = nil
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		upper := strings.ToUpper(strings.TrimSpace(line))

		if upper == "#+BEGIN_SRC" || strings.HasPrefix(upper, "#+BEGIN_SRC ") {
			flushParagraph()
			inCode = true
			continue
		}
		if upper == "#+END_SRC" {
			flushCode()
			inCode = false
			continue
		}
		if inCode {
			codeLines = append(codeLines, line)
			continue
		}
		if parseOrgMetadata(&doc.Metadata, line) {
			continue
		}

		level, heading := orgHeading(line)
		if level > 0 {
			flushParagraph()
			if strings.Contains(strings.ToLower(heading), ":noexport:") {
				inNoExport = true
				noExportLevel = level
				continue
			}
			if inNoExport && level <= noExportLevel {
				inNoExport = false
			}
			if inNoExport {
				continue
			}
			addOrgHeading(&doc, level, heading)
			continue
		}
		if inNoExport {
			continue
		}
		if strings.TrimSpace(line) == "" {
			flushParagraph()
			continue
		}
		if strings.TrimSpace(line) == "-----" || strings.TrimSpace(line) == "_____" {
			flushParagraph()
			doc.Blocks = append(doc.Blocks, Block{Kind: "scene-break"})
			continue
		}
		if name, value, ok := orgFootnote(line); ok {
			flushParagraph()
			doc.Blocks = append(doc.Blocks, Block{Kind: "footnote", Text: name + "\t" + value})
			continue
		}
		paragraph = append(paragraph, canonicalizeOrgInline(strings.TrimSpace(line)))
	}
	flushParagraph()
	flushCode()
	return doc
}

func markdownHeading(line string) (int, string) {
	hashes := 0
	for hashes < len(line) && line[hashes] == '#' {
		hashes++
	}
	if hashes == 0 || hashes > 6 || len(line) <= hashes || line[hashes] != ' ' {
		return 0, ""
	}
	return hashes, strings.TrimSpace(line[hashes+1:])
}

func addMarkdownHeading(doc *Document, level int, heading string) {
	switch level {
	case 1:
		doc.Blocks = append(doc.Blocks, Block{Kind: "part", Level: level, Text: heading})
	case 2:
		doc.Blocks = append(doc.Blocks, Block{Kind: "chapter", Level: level, Text: heading})
	default:
		doc.Blocks = append(doc.Blocks, Block{Kind: "section", Level: level, Text: heading})
	}
}

func addOrgHeading(doc *Document, level int, heading string) {
	switch level {
	case 1:
		doc.Blocks = append(doc.Blocks, Block{Kind: "part", Level: level, Text: heading})
	case 2:
		doc.Blocks = append(doc.Blocks, Block{Kind: "chapter", Level: level, Text: heading})
	default:
		doc.Blocks = append(doc.Blocks, Block{Kind: "section", Level: level, Text: heading})
	}
}

func orgHeading(line string) (int, string) {
	stars := 0
	for stars < len(line) && line[stars] == '*' {
		stars++
	}
	if stars == 0 || len(line) <= stars || line[stars] != ' ' {
		return 0, ""
	}
	return stars, strings.TrimSpace(line[stars+1:])
}

func parseMarkdownMetadata(meta *Metadata, line string) bool {
	trimmed := strings.TrimSpace(line)
	if meta.Subtitle == "" && strings.HasPrefix(trimmed, "**") && strings.HasSuffix(trimmed, "**") {
		meta.Subtitle = strings.TrimSuffix(strings.TrimPrefix(trimmed, "**"), "**")
		return true
	}
	if meta.Author == "" && strings.HasPrefix(trimmed, "*by ") && strings.HasSuffix(trimmed, "*") {
		meta.Author = strings.TrimSuffix(strings.TrimPrefix(trimmed, "*by "), "*")
		meta.AuthorAttribution = "by"
		return true
	}
	if strings.HasPrefix(trimmed, "--- ") && strings.HasSuffix(trimmed, " ---") {
		content := strings.TrimSuffix(strings.TrimPrefix(trimmed, "--- "), " ---")
		parts := strings.Split(content, "|")
		if len(parts) == 2 {
			meta.Version = strings.TrimSpace(parts[0])
			meta.Date = strings.TrimSpace(parts[1])
			return true
		}
		meta.Version = strings.TrimSpace(content)
		return true
	}
	return false
}

func parseMarkdownFrontmatter(meta *Metadata, text string) (string, error) {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	if !strings.HasPrefix(normalized, "---\n") {
		return text, nil
	}
	end := strings.Index(normalized[len("---\n"):], "\n---\n")
	closingLen := len("\n---\n")
	if end < 0 && strings.HasSuffix(normalized, "\n---") {
		end = len(normalized) - len("---\n") - len("\n---")
		closingLen = len("\n---")
	}
	if end < 0 {
		return "", fmt.Errorf("markdown frontmatter starts with --- but has no closing ---")
	}
	content := normalized[len("---\n") : len("---\n")+end]
	remaining := normalized[len("---\n")+end+closingLen:]
	values := map[string]any{}
	if err := yaml.Unmarshal([]byte(content), &values); err != nil {
		return "", fmt.Errorf("parsing markdown frontmatter: %w", err)
	}
	applyMarkdownFrontmatter(meta, values)
	return remaining, nil
}

func applyMarkdownFrontmatter(meta *Metadata, values map[string]any) {
	for key, value := range values {
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "title":
			meta.Title = frontmatterString(value)
		case "subtitle":
			meta.Subtitle = frontmatterString(value)
		case "author":
			meta.Author = frontmatterString(value)
		case "author-attribution", "author_attribution":
			meta.AuthorAttribution = frontmatterString(value)
		case "date":
			meta.Date = frontmatterString(value)
		case "version", "draft":
			meta.Version = frontmatterString(value)
		case "wordcount", "word count", "word-count":
			meta.WordCount = frontmatterString(value)
		case "contact-name", "contact_name", "contact":
			meta.ContactName = frontmatterString(value)
		case "address":
			meta.Address = frontmatterString(value)
		case "phone", "telephone":
			meta.Phone = frontmatterString(value)
		case "email":
			meta.Email = frontmatterString(value)
		case "website":
			meta.Website = frontmatterString(value)
		}
	}
}

func frontmatterString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case int:
		return fmt.Sprintf("%d", typed)
	case int64:
		return fmt.Sprintf("%d", typed)
	case float64:
		return fmt.Sprintf("%.0f", typed)
	case time.Time:
		return typed.Format("2006-01-02")
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			if rendered := frontmatterString(item); rendered != "" {
				parts = append(parts, rendered)
			}
		}
		return strings.Join(parts, " / ")
	default:
		return fmt.Sprint(typed)
	}
}

func parseMarkdownMetadataTable(meta *Metadata, line string, inMetadataTable *bool) bool {
	if !strings.HasPrefix(strings.TrimSpace(line), "|") {
		return false
	}
	cells := markdownTableCells(line)
	if len(cells) < 2 {
		return false
	}
	key := strings.ToLower(strings.TrimSpace(cells[0]))
	value := strings.TrimSpace(cells[1])
	if key == "metadata" || key == "field" {
		*inMetadataTable = true
		return true
	}
	if strings.Trim(key, "-: ") == "" {
		return *inMetadataTable
	}
	if !*inMetadataTable && !isMarkdownMetadataKey(key) {
		return false
	}
	switch key {
	case "title":
		meta.Title = value
	case "subtitle":
		meta.Subtitle = value
	case "author":
		meta.Author = value
	case "date":
		meta.Date = value
	case "version", "draft":
		meta.Version = value
	case "wordcount", "word count", "word-count":
		meta.WordCount = value
	case "contact-name", "contact name", "contact":
		meta.ContactName = value
	case "address":
		meta.Address = value
	case "phone", "telephone":
		meta.Phone = value
	case "email":
		meta.Email = value
	case "website":
		meta.Website = value
	default:
		return *inMetadataTable
	}
	*inMetadataTable = true
	return true
}

func markdownTableCells(line string) []string {
	trimmed := strings.Trim(strings.TrimSpace(line), "|")
	rawCells := strings.Split(trimmed, "|")
	cells := make([]string, 0, len(rawCells))
	for _, cell := range rawCells {
		cells = append(cells, strings.TrimSpace(cell))
	}
	return cells
}

func isMarkdownMetadataKey(key string) bool {
	switch key {
	case "title", "subtitle", "author", "date", "version", "draft", "wordcount", "word count", "word-count", "contact-name", "contact name", "contact", "address", "phone", "telephone", "email", "website":
		return true
	default:
		return false
	}
}

func parseOrgMetadata(meta *Metadata, line string) bool {
	if !strings.HasPrefix(line, "#+") {
		return false
	}
	parts := strings.SplitN(line[2:], ":", 2)
	if len(parts) != 2 {
		return false
	}
	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])
	switch key {
	case "title":
		meta.Title = value
	case "subtitle":
		meta.Subtitle = value
	case "author":
		meta.Author = value
	case "date":
		meta.Date = value
	case "version":
		meta.Version = value
	case "wordcount":
		meta.WordCount = value
	case "contact-name", "contact_name", "contact":
		meta.ContactName = value
	case "address":
		meta.Address = value
	case "phone":
		meta.Phone = value
	case "email":
		meta.Email = value
	case "website":
		meta.Website = value
	}
	return true
}

func markdownFootnote(line string) (string, string, bool) {
	re := regexp.MustCompile(`^\[\^([^\]]+)\]:\s+(.+)$`)
	match := re.FindStringSubmatch(line)
	if match == nil {
		return "", "", false
	}
	return match[1], match[2], true
}

func orgFootnote(line string) (string, string, bool) {
	re := regexp.MustCompile(`^\[fn:([^\]]+)\]\s+(.+)$`)
	match := re.FindStringSubmatch(line)
	if match == nil {
		return "", "", false
	}
	return match[1], match[2], true
}

func canonicalizeOrgInline(text string) string {
	text = regexp.MustCompile(`\*([^*\n]+)\*`).ReplaceAllString(text, `**$1**`)
	text = regexp.MustCompile(`/([^/\n]+?)/`).ReplaceAllString(text, `*$1*`)
	text = regexp.MustCompile(`~([^~\n]+?)~`).ReplaceAllString(text, "`$1`")
	text = regexp.MustCompile(`=([^=\n]+?)=`).ReplaceAllString(text, "`$1`")
	return text
}
