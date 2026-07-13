// ABOUTME: Verifies manuscript dispatch through the single in-process Go application.
// ABOUTME: Covers embedded help and dry-run output without the legacy helper boundary.
package app

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestManuscriptDispatchInProcess(t *testing.T) {
	status, stdout, stderr := runApp(t, "manuscript", "--help")
	if status != 0 || stderr != "" || !strings.Contains(stdout, "Usage: folio manuscript") {
		t.Fatalf("help status %d\nstdout:%s\nstderr:%s", status, stdout, stderr)
	}

	dir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	source := filepath.Join(dir, "chapter.md")
	writeAppFile(t, source, "---\ntitle: Test Manuscript\nauthor: Example Author\n---\n\n## Chapter 1\n\nBody.\n")
	status, stdout, stderr = runApp(t, "manuscript", "--dry-run", source, filepath.Join(dir, "out.pdf"))
	if status != 0 || stderr != "" {
		t.Fatalf("dry-run status %d\nstdout:%s\nstderr:%s", status, stdout, stderr)
	}
	for _, fragment := range []string{"format: markdown", "style: british", "page: a4"} {
		if !strings.Contains(stdout, fragment) {
			t.Errorf("dry-run missing %q:\n%s", fragment, stdout)
		}
	}
}
