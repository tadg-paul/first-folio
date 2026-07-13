// ABOUTME: Escapes plain content and string literals for Typst templates.
// ABOUTME: Centralizes the base contexts shared by script and letter renderers.
package typst

import "strings"

func EscapeContent(value string) string {
	return strings.NewReplacer(
		`\`, `\\`,
		`$`, `\$`,
		`#`, `\#`,
		`@`, `\@`,
		`[`, `\[`,
		`]`, `\]`,
		`*`, `\*`,
		`_`, `\_`,
	).Replace(value)
}

func EscapeString(value string) string {
	return strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"\n", `\n`,
		"\r", `\r`,
	).Replace(value)
}
