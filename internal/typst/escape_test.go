// ABOUTME: Verifies shared Typst content and string-literal escaping.
// ABOUTME: Covers syntax-significant ASCII and Unicode preservation.
package typst

import "testing"

func TestEscapeContexts(t *testing.T) {
	if got := EscapeContent(`Taḋg [#1] costs $5 @ home_*`); got != `Taḋg \[\#1\] costs \$5 \@ home\_\*` {
		t.Fatalf("content escape = %q", got)
	}
	if got := EscapeString("Taḋg \\\"\n"); got != `Taḋg \\\"\n` {
		t.Fatalf("string escape = %q", got)
	}
}
