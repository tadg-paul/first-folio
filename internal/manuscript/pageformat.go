// ABOUTME: Manuscript page-format parser for named presets and custom WxH dimensions.
// ABOUTME: Rejects malformed values at config load time so Typst gets a well-formed argument.
package manuscript

import (
	"fmt"
	"regexp"
)

// PageSpec is the parsed shape of a manuscript page value.
type PageSpec struct {
	Custom bool
	Named  string
	Width  string
	Height string
}

var (
	pageCustomRE = regexp.MustCompile(`^(\d+(?:\.\d+)?)x(\d+(?:\.\d+)?)(mm|in)$`)
	pageNamedRE  = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
)

// knownPapers is the allowlist of named page presets. It combines the Typst standard papers
// with a small set of UK book sizes referenced by first-folio users.
var knownPapers = map[string]struct{}{
	"a0": {}, "a1": {}, "a2": {}, "a3": {}, "a4": {}, "a5": {}, "a6": {}, "a7": {}, "a8": {}, "a9": {}, "a10": {}, "a11": {},
	"iso-b1": {}, "iso-b2": {}, "iso-b3": {}, "iso-b4": {}, "iso-b5": {}, "iso-b6": {}, "iso-b7": {}, "iso-b8": {},
	"iso-c3": {}, "iso-c4": {}, "iso-c5": {}, "iso-c6": {}, "iso-c7": {}, "iso-c8": {},
	"din-d3": {}, "din-d4": {}, "din-d5": {}, "din-d6": {}, "din-d7": {}, "din-d8": {},
	"iso-id-1": {}, "iso-id-2": {}, "iso-id-3": {},
	"ansi-a": {}, "ansi-b": {}, "ansi-c": {}, "ansi-d": {}, "ansi-e": {},
	"arch-a": {}, "arch-b": {}, "arch-c": {}, "arch-d": {}, "arch-e": {}, "arch-e1": {}, "arch-e2": {}, "arch-e3": {},
	"jis-b0": {}, "jis-b1": {}, "jis-b2": {}, "jis-b3": {}, "jis-b4": {}, "jis-b5": {}, "jis-b6": {}, "jis-b7": {}, "jis-b8": {}, "jis-b9": {}, "jis-b10": {}, "jis-b11": {},
	"jp-shiroku-ban-4": {}, "jp-shiroku-ban-5": {}, "jp-shiroku-ban-6": {},
	"jp-kiku-4": {}, "jp-kiku-5": {},
	"us-letter": {}, "us-legal": {}, "us-tabloid": {}, "us-executive": {}, "us-statement": {}, "us-financial": {},
	"us-trade": {}, "us-digest": {}, "us-mass-market": {},
	"presentation-16-9": {}, "presentation-4-3": {},
	"uk-royal": {}, "uk-crown": {}, "uk-large-crown": {}, "uk-demy": {}, "uk-medium": {},
	"uk-a-format-paperback": {}, "uk-b-format-paperback": {},
	"uk-pinched-crown": {}, "uk-pinched-post": {}, "uk-post": {}, "uk-pinched-demy": {},
	"uk-foolscap": {}, "uk-imperial": {},
	"uk-book-a": {}, "uk-book-b": {}, "uk-book-c": {},
	"letter": {}, "legal": {},
}

// ParsePageSpec parses a manuscript page value into a PageSpec, returning an error if the value
// is neither a known named preset nor a valid WxH custom dimension.
func ParsePageSpec(value string) (PageSpec, error) {
	if m := pageCustomRE.FindStringSubmatch(value); m != nil {
		return PageSpec{
			Custom: true,
			Width:  m[1] + m[3],
			Height: m[2] + m[3],
		}, nil
	}
	if pageNamedRE.MatchString(value) {
		if _, ok := knownPapers[value]; ok {
			return PageSpec{Named: value}, nil
		}
	}
	return PageSpec{}, fmt.Errorf("invalid manuscript page value %q: expected a known preset (e.g. a4, us-letter, uk-book-b) or custom WxH in mm or in (e.g. 200x300mm, 5.5x8.5in)", value)
}
