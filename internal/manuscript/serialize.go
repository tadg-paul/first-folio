// ABOUTME: Manuscript document serializers for canonical Markdown and org-mode.
// ABOUTME: Keeps Markdown and org contracts in sync through the shared manuscript AST.
package manuscript

import (
	"strings"
)

func RenderMarkdown(doc Document) string {
	var lines []string
	meta := doc.Metadata
	if meta.Title != "" {
		lines = append(lines, "# "+meta.Title, "")
	}
	if meta.Subtitle != "" {
		lines = append(lines, "**"+meta.Subtitle+"**", "")
	}
	if meta.Author != "" {
		attribution := meta.AuthorAttribution
		if attribution == "" {
			attribution = "by"
		}
		lines = append(lines, "*"+strings.TrimSpace(attribution+" "+meta.Author)+"*", "")
	}
	if meta.Version != "" || meta.Date != "" {
		dateLine := strings.TrimSpace(meta.Version)
		if meta.Date != "" {
			if dateLine != "" {
				dateLine += " | "
			}
			dateLine += meta.Date
		}
		lines = append(lines, "--- "+dateLine+" ---", "")
	}
	appendMarkdownMetadataTable(&lines, meta)
	appendMarkdownBlocks(&lines, doc.Blocks)
	return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
}

func appendMarkdownMetadataTable(lines *[]string, meta Metadata) {
	if meta.WordCount == "" && meta.Address == "" && meta.Phone == "" && meta.Email == "" && meta.Website == "" {
		return
	}
	*lines = append(*lines, "| Metadata | Value |", "|---|---|")
	appendMarkdownMetadataRow(lines, "Wordcount", meta.WordCount)
	appendMarkdownMetadataRow(lines, "Address", meta.Address)
	appendMarkdownMetadataRow(lines, "Phone", meta.Phone)
	appendMarkdownMetadataRow(lines, "Email", meta.Email)
	appendMarkdownMetadataRow(lines, "Website", meta.Website)
	*lines = append(*lines, "")
}

func appendMarkdownMetadataRow(lines *[]string, key string, value string) {
	if value != "" {
		*lines = append(*lines, "| "+key+" | "+value+" |")
	}
}

func appendMarkdownBlocks(lines *[]string, blocks []Block) {
	for _, block := range blocks {
		switch block.Kind {
		case "part":
			*lines = append(*lines, "## "+block.Text, "")
		case "chapter":
			*lines = append(*lines, "### "+block.Text, "")
		case "section":
			level := block.Level
			if level < 4 {
				level = 4
			}
			*lines = append(*lines, strings.Repeat("#", level)+" "+block.Text, "")
		case "paragraph":
			*lines = append(*lines, block.Text, "")
		case "scene-break":
			*lines = append(*lines, "***", "")
		case "code":
			*lines = append(*lines, "```"+block.Lang, block.Text, "```", "")
		}
	}
}
