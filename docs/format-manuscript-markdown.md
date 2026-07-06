<!-- Version: 0.1 | Last updated: 2026-07-06 -->

# Markdown Manuscript Format

Markdown manuscript input is a prose contract, separate from the Markdown stage-play contract.

## Metadata Contract

All YAML frontmatter values are treated as manuscript strings. Quote values when that keeps the intent clearest, but the parser also accepts YAML scalars and converts them to strings, so `wordcount: about 90,000 words` and `wordcount: 90000` are both valid. Dates should be written as ISO strings such as `2026-07-06`.

Supported frontmatter fields are `title`, `subtitle`, `author`, `author-attribution`, `date`, `version`, `wordcount`, `contact-name`, `address`, `phone`, `email`, and `website`. `contact-name` is optional and is used only for the title-page contact block; it does not default to the manuscript author.

## Element Schema

| Markdown syntax | Manuscript meaning |
|---|---|
| YAML frontmatter bounded by `---` | Manuscript metadata |
| `# PART ONE` | Part divider page |
| `## Chapter 1` | Chapter start page |
| `### Section` and deeper | Local section heading |
| Plain paragraphs | Body prose |
| `***` on its own line | Scene break |
| `**bold**` | Bold text |
| `*italic*` | Italic text |
| `` `code` `` and fenced code blocks | Monospace text |
| `--` and `---` | En dash and em dash |
| `[^name]` and `[^name]: text` | Footnote reference and definition |
| HTML comments | Private notes, excluded |
| Heading ending `<!-- noexport -->` | Private section excluded until the next same-or-higher heading |

Fountain is not accepted by manuscript mode.

## Example

```markdown
---
title: The Glass Orchard
subtitle: A Novel
author: Example Author
author-attribution: by
date: 2026-07-06
version: Draft 2
wordcount: about 90,000 words
contact-name: Example Agent
address: 100 Example Street / Sample City / Exampleland
phone: +353 1 000 0000
email: author@example.invalid
website: https://example.invalid
---

# PART ONE

## Chapter 1

The rain had been falling since Tuesday. The ledger flashed **WAIT** -- then the latch answered --- and Mira typed `nine-bell`.

***

By noon, the hands had moved backwards twice.

### Notes <!-- noexport -->

This planning note is excluded.
```
