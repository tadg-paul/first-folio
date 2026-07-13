// ABOUTME: Characterizes template-backed script Typst and PDF output through the Go CLI.
// ABOUTME: Covers configured layout, Unicode, hostile characters, and real compilation.
package app

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertTypstUsesConfiguredTemplate(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	source := filepath.Join(dir, "play.org")
	target := filepath.Join(dir, "play.typ")
	writeAppFile(t, source, `#+TITLE: Syntax # [Test]
#+AUTHOR: Taḋg Paul
* ACT ONE
** Scene One
*** A sign reads #OPEN [NOW].
**** CÁIT quietly
Price is $5 and path is C:\tmp.[fn:cost]
[fn:cost] A **bold** note.
`)
	writeAppFile(t, filepath.Join(dir, "script.yaml"), `folio:
  font: Libertinus Serif
  font-size: 11pt
  page: a4
  margin: 20mm
  positioning:
    speech:
      space-before: 2em
      speaker:
        bold: false
      dialogue:
        wrap-indent: 8em
`)

	status, stdout, stderr := runApp(t, "convert", source, target)
	if status != 0 {
		t.Fatalf("status %d\nstdout:%s\nstderr:%s", status, stdout, stderr)
	}
	typst := readAppFile(t, target)
	for _, fragment := range []string{
		`#set page(paper: "a4", margin: 20mm`,
		`#set text(font: "Libertinus Serif", size: 11pt`,
		`columns: (8em, 1fr)`,
		`weight: "regular"`,
		`Syntax \# \[Test\]`,
		`Price is \$5 and path is C:\\tmp.`,
		`#footnote[A *bold* note.]`,
	} {
		if !strings.Contains(typst, fragment) {
			t.Errorf("Typst missing %q:\n%s", fragment, typst)
		}
	}
}

func TestConvertPDFCompiles(t *testing.T) {
	if _, err := exec.LookPath("typst"); err != nil {
		t.Skip("typst is not installed")
	}
	dir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	source := filepath.Join(dir, "play.md")
	target := filepath.Join(dir, "play.pdf")
	writeAppFile(t, source, "# Samhain\n\n*by Taḋg Paul*\n\n## ACT ONE\n\n**CÁIT:**\nHello.\n")

	status, stdout, stderr := runApp(t, "convert", source, target)
	if status != 0 {
		t.Fatalf("status %d\nstdout:%s\nstderr:%s", status, stdout, stderr)
	}
	raw, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(raw, []byte("%PDF")) || len(raw) < 1000 {
		t.Fatalf("invalid PDF output: %d bytes", len(raw))
	}
}
