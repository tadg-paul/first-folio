<!-- Version: 0.1 | Last updated: 2026-07-06 -->

# Markdown Manuscript Format

Markdown manuscript input is a prose contract, separate from the Markdown stage-play contract.

## Element Schema

| Markdown syntax | Manuscript meaning |
|---|---|
| `# Title` | Manuscript title |
| `**Subtitle**` after the title | Subtitle |
| `*by Author*` | Author name |
| `--- Draft | Date ---` | Version and date |
| Metadata table with `Metadata`/`Value` headings | Optional manuscript metadata: word count, address, phone, email, website |
| `## PART ONE` | Part divider page |
| `### Chapter 1` | Chapter start page |
| `#### Section` and deeper | Local section heading |
| Plain paragraphs | Body prose |
| `***` on its own line | Scene break |
| Backticks and fenced code blocks | Monospace text |
| `[^name]` and `[^name]: text` | Footnote reference and definition |
| HTML comments | Private notes, excluded |
| Heading ending `<!-- noexport -->` | Private section excluded until the next same-or-higher heading |

Fountain is not accepted by manuscript mode.

## Example

```markdown
# The Glass Orchard

**A Novel**

*by Example Author*

--- Draft 2 | July 2026 ---

| Metadata | Value |
|---|---|
| Wordcount | 90000 |
| Address | 100 Example Street / Sample City / Exampleland |
| Phone | +353 1 000 0000 |
| Email | author@example.invalid |
| Website | https://example.invalid |

## PART ONE

### Chapter 1

The rain had been falling since Tuesday.

***

By noon, the hands had moved backwards twice.

### Notes <!-- noexport -->

This planning note is excluded.
```
