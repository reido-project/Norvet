package nvtext

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// sentenceBoundary matches the punctuation that ends a sentence together with
// the whitespace after it. Punctuation inside a token ("3.14", "example.com")
// is not a boundary because no whitespace follows it.
var sentenceBoundary = regexp.MustCompile(`[.?!;:]+(?:\s+|$)`)

// SentenceCase capitalizes the first letter of every sentence in input and
// leaves the rest untouched. Sentences end at ".", "?", "!", ";" or ":"
// followed by whitespace or the end of the text.
//
//	SentenceCase("hello world. visit example.com! it costs 3.14")
//	// "Hello world. Visit example.com! It costs 3.14"
//
// A sentence that starts with a digit is left as is, and leading quotes,
// brackets and spaces are skipped. A blank input gives "".
func SentenceCase(input string) string {
	if strings.TrimSpace(input) == "" {
		return ""
	}

	var sb strings.Builder
	sb.Grow(len(input))

	last := 0
	for _, m := range sentenceBoundary.FindAllStringIndex(input, -1) {
		sb.WriteString(capitalizeFirstLetter(input[last:m[0]]))
		sb.WriteString(input[m[0]:m[1]])
		last = m[1]
	}
	sb.WriteString(capitalizeFirstLetter(input[last:]))

	return sb.String()
}

// capitalizeFirstLetter upper-cases the first letter of s, unless a digit
// comes before it. Anything else before the letter (spaces, quotes, brackets)
// is skipped.
func capitalizeFirstLetter(s string) string {
	for i, r := range s {
		switch {
		case unicode.IsLetter(r):
			_, size := utf8.DecodeRuneInString(s[i:])
			return s[:i] + string(unicode.ToUpper(r)) + s[i+size:]
		case unicode.IsDigit(r):
			return s
		}
	}
	return s
}
