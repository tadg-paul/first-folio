// ABOUTME: Pure-Go EAN-13 barcode SVG generator for ISBN rendering (issue #21).
// ABOUTME: No external dependency; produces a compact SVG that Typst can embed.
package manuscript

import (
	"fmt"
	"strings"
)

// EAN-13 encoding tables. Each digit encodes to 7 modules (bits: 1 = bar,
// 0 = space). Three sets: A (odd parity), B (even parity), C (right group).
// The first digit selects which A/B parity pattern is used across positions
// 2-7 of the left group; positions 8-13 always use set C.
var (
	ean13SetA = [10]string{
		"0001101", "0011001", "0010011", "0111101", "0100011",
		"0110001", "0101111", "0111011", "0110111", "0001011",
	}
	ean13SetB = [10]string{
		"0100111", "0110011", "0011011", "0100001", "0011101",
		"0111001", "0000101", "0010001", "0001001", "0010111",
	}
	ean13SetC = [10]string{
		"1110010", "1100110", "1101100", "1000010", "1011100",
		"1001110", "1010000", "1000100", "1001000", "1110100",
	}
	// Parity pattern for the left group, indexed by the leading (first) digit.
	// "A" = set A, "B" = set B.
	ean13LeftParity = [10]string{
		"AAAAAA", "AABABB", "AABBAB", "AABBBA", "ABAABB",
		"ABBAAB", "ABBBAA", "ABABAB", "ABABBA", "ABBABA",
	}
)

// renderEAN13SVG returns an EAN-13 barcode as an SVG string for the given ISBN.
// The ISBN may contain hyphens or spaces; they are stripped. The input must have
// 13 digits with a valid check digit (validated separately in validateISBN13).
func renderEAN13SVG(isbn string) (string, error) {
	digits := strings.ReplaceAll(isbn, "-", "")
	digits = strings.ReplaceAll(digits, " ", "")
	if len(digits) != 13 {
		return "", fmt.Errorf("EAN-13 requires 13 digits, got %d in %q", len(digits), isbn)
	}
	nums := make([]int, 13)
	for i, r := range digits {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("EAN-13 requires numeric digits, got %q in %q", r, isbn)
		}
		nums[i] = int(r - '0')
	}
	// Build the module string: guard + left group (6 digits, A/B parity per first digit) + middle guard + right group (6 digits, C set) + guard.
	var modules strings.Builder
	modules.WriteString("101") // left guard
	parity := ean13LeftParity[nums[0]]
	for i := 0; i < 6; i++ {
		digit := nums[i+1]
		if parity[i] == 'A' {
			modules.WriteString(ean13SetA[digit])
		} else {
			modules.WriteString(ean13SetB[digit])
		}
	}
	modules.WriteString("01010") // middle guard
	for i := 0; i < 6; i++ {
		modules.WriteString(ean13SetC[nums[i+7]])
	}
	modules.WriteString("101") // right guard

	// Render as SVG. Standard EAN-13 dimensions: 95 modules wide + margins; each
	// module is 0.33mm at 100% magnification. We render at a normalized viewBox
	// with 1-unit modules and let the consumer scale via image() width.
	const (
		moduleWidth = 1
		height      = 60
		textHeight  = 8
	)
	viewWidth := len(modules.String())*moduleWidth + 20 // margins
	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d">`,
		viewWidth, height+textHeight+4,
	))
	b.WriteString(`<rect width="100%" height="100%" fill="white"/>`)
	// Draw bars for each module set to '1'.
	x := 10 // left margin
	mods := modules.String()
	for i := 0; i < len(mods); i++ {
		if mods[i] == '1' {
			// Extend guard bars slightly below the numeric row for authenticity.
			barHeight := height
			if isGuardModule(i, len(mods)) {
				barHeight = height + 4
			}
			b.WriteString(fmt.Sprintf(
				`<rect x="%d" y="0" width="%d" height="%d" fill="black"/>`,
				x, moduleWidth, barHeight,
			))
		}
		x += moduleWidth
	}
	// Human-readable digit row below the bars.
	// First digit sits to the left of the left guard; digits 2-7 under the left group;
	// digits 8-13 under the right group.
	textY := height + textHeight
	b.WriteString(fmt.Sprintf(
		`<text x="4" y="%d" font-family="monospace" font-size="8">%d</text>`,
		textY, nums[0],
	))
	leftGroupStart := 13
	for i := 0; i < 6; i++ {
		cx := leftGroupStart + i*7 + 3
		b.WriteString(fmt.Sprintf(
			`<text x="%d" y="%d" font-family="monospace" font-size="8" text-anchor="middle">%d</text>`,
			cx, textY, nums[i+1],
		))
	}
	rightGroupStart := 13 + 42 + 5 // left guard + 6*7 modules + middle guard
	for i := 0; i < 6; i++ {
		cx := rightGroupStart + i*7 + 3
		b.WriteString(fmt.Sprintf(
			`<text x="%d" y="%d" font-family="monospace" font-size="8" text-anchor="middle">%d</text>`,
			cx, textY, nums[i+7],
		))
	}
	b.WriteString(`</svg>`)
	return b.String(), nil
}

// isGuardModule returns true if the module at index i falls within one of the
// three EAN-13 guard patterns (left/middle/right), so guard bars can be drawn
// slightly taller than digit bars for visual authenticity.
func isGuardModule(i, total int) bool {
	// Left guard: 3 modules at positions 0-2
	if i < 3 {
		return true
	}
	// Right guard: last 3 modules
	if i >= total-3 {
		return true
	}
	// Middle guard: 5 modules at positions 45-49 (after left guard + 6*7 = 3 + 42 = 45)
	if i >= 45 && i < 50 {
		return true
	}
	return false
}
