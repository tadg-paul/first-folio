// ABOUTME: Verifies deterministic Homebrew formula updates without network or publication.
// ABOUTME: Protects URL and digest replacement used by the release target.
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateFormula(t *testing.T) {
	path := filepath.Join(t.TempDir(), "formula.rb")
	input := "class FirstFolio < Formula\n  url \"old\"\n  sha256 \"oldsum\"\nend\n"
	if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := updateFormula(path, "https://example.invalid/v1.tar.gz", "abc123"); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	output := string(raw)
	if !strings.Contains(output, `url "https://example.invalid/v1.tar.gz"`) || !strings.Contains(output, `sha256 "abc123"`) {
		t.Fatalf("formula not updated:\n%s", output)
	}
}
