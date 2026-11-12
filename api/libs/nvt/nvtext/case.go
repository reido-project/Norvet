package nvtext

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// tokenize splits input into lower-case words. A new word starts at every
// non-alphanumeric character, at a letter/digit boundary, at a lower-to-upper
// transition ("fooBar") and at the end of an acronym ("HTTPServer" gives
// "http", "server"). Apostrophes split too, so "don't" gives "don", "t".
func tokenize(input string) []string {
	runes := []rune(input)
	var words []string
	var current []rune

	flush := func() {
		if len(current) > 0 {
			words = append(words, strings.ToLower(string(current)))
			current = current[:0]
		}
	}

	for i, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			flush()
			continue
		}

		if len(current) > 0 {
			prev := current[len(current)-1]
			nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])

			if unicode.IsDigit(prev) != unicode.IsDigit(r) ||
				(unicode.IsUpper(r) && (unicode.IsLower(prev) || unicode.IsDigit(prev))) ||
				(unicode.IsUpper(prev) && unicode.IsUpper(r) && nextIsLower) {
				flush()
			}
		}

		current = append(current, r)
	}

	flush()
	return words
}

// capitalize upper-cases the first letter of word.
func capitalize(word string) string {
	r, size := utf8.DecodeRuneInString(word)
	if size == 0 {
		return word
	}
	return string(unicode.ToUpper(r)) + word[size:]
}

// joinCapitalized capitalizes words[from:] in place and joins all words with sep.
func joinCapitalized(words []string, from int, sep string) string {
	for i := from; i < len(words); i++ {
		words[i] = capitalize(words[i])
	}
	return strings.Join(words, sep)
}

// PascalCase converts input to PascalCase: "user_id" gives "UserId".
func PascalCase(input string) string {
	return joinCapitalized(tokenize(input), 0, "")
}

// CamelCase converts input to camelCase: "user_id" gives "userId".
func CamelCase(input string) string {
	return joinCapitalized(tokenize(input), 1, "")
}

// SnakeCase converts input to snake_case: "UserID" gives "user_id".
func SnakeCase(input string) string {
	return strings.Join(tokenize(input), "_")
}

// KebabCase converts input to kebab-case: "UserID" gives "user-id".
func KebabCase(input string) string {
	return strings.Join(tokenize(input), "-")
}

// TitleCase converts input to Title Case: "user_id" gives "User Id".
func TitleCase(input string) string {
	return joinCapitalized(tokenize(input), 0, " ")
}
