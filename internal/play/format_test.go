// ABOUTME: Characterizes supported stage-play format names and file extensions.
// ABOUTME: Protects aliases, readable formats, and diagnostics during the Go migration.
package play

import "testing"

func TestFormatDetection(t *testing.T) {
	tests := []struct {
		path string
		want Format
	}{
		{"play.org", FormatOrg},
		{"play.md", FormatMarkdown},
		{"play.markdown", FormatMarkdown},
		{"play.fountain", FormatFountain},
		{"play.ftn", FormatFountain},
		{"play.pdf", FormatPDF},
	}
	for _, tt := range tests {
		got, err := FormatFromPath(tt.path)
		if err != nil {
			t.Fatalf("FormatFromPath(%q): %v", tt.path, err)
		}
		if got != tt.want {
			t.Errorf("FormatFromPath(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestFormatValidation(t *testing.T) {
	if _, err := FormatFromPath("play"); err == nil {
		t.Fatal("extensionless input should fail")
	}
	if _, err := FormatFromPath("play.xyz"); err == nil {
		t.Fatal("unknown extension should fail")
	}
	if FormatPDF.Readable() {
		t.Fatal("PDF must remain write-only")
	}
	for _, name := range []string{"org", "md", "markdown", "fountain", "ftn", "pdf"} {
		if _, err := ParseFormat(name); err != nil {
			t.Errorf("ParseFormat(%q): %v", name, err)
		}
	}
}
