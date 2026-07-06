<!-- Version: 0.1 | Last updated: 2026-07-06 -->

# Org-mode Manuscript Format

Org-mode manuscript input uses org front matter and headings for prose manuscript structure. It is separate from the org-mode stage-play contract.

## Element Schema

| Org syntax | Manuscript meaning |
|---|---|
| `#+TITLE: The Glass Orchard` | Manuscript title |
| `#+SUBTITLE: A Novel` | Subtitle |
| `#+AUTHOR: Example Author` | Author name |
| `#+DATE: July 2026` | Manuscript date |
| `#+VERSION: Draft 4` | Draft/version marker |
| `#+WORDCOUNT: 80000` | Approximate word count |
| `#+ADDRESS: ...` | Postal address |
| `#+PHONE: ...` | Phone number |
| `#+EMAIL: ...` | Email address |
| `#+WEBSITE: ...` | Website |
| `* PART ONE` | Part divider page |
| `** Chapter 1` | Chapter start page |
| `*** Section` and deeper | Local section heading |
| Plain paragraphs | Body prose |
| `-----` on its own line | Scene break |
| `~code~`, `=verbatim=`, and source blocks | Monospace text |
| `[fn:name]` and `[fn:name] Text` | Footnote reference and definition |
| Heading tagged `:noexport:` | Private section excluded with children |

Fountain is not accepted by manuscript mode.

## Example

```org
#+TITLE: The Glass Orchard
#+SUBTITLE: A Novel
#+AUTHOR: Example Author
#+DATE: July 2026
#+VERSION: Draft 4
#+WORDCOUNT: 80000
#+ADDRESS: 100 Example Street / Sample City / Exampleland
#+PHONE: +353 1 000 0000
#+EMAIL: author@example.invalid
#+WEBSITE: https://example.invalid

* PART ONE
** Chapter 1
The rain had been falling since Tuesday.

-----

By noon, the hands had moved backwards twice.

*** Notes :noexport:
This planning note is excluded.
```
