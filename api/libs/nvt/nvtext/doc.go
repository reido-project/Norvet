// Package nvtext has small text helpers: identifier case conversion, sentence
// capitalization and timezone-name repair. It depends only on the standard
// library.
//
// # Case conversion
//
// The case functions split their input into words first, so they accept any
// style as input (snake_case, kebab-case, camelCase, Title Case, "HTTPServer",
// "user2name") and give the same result:
//
//	nvtext.PascalCase("user_id")   // UserId
//	nvtext.CamelCase("user_id")    // userId
//	nvtext.SnakeCase("UserID")     // user_id
//	nvtext.KebabCase("UserID")     // user-id
//	nvtext.TitleCase("user_id")    // User Id
//
// Words break at non-alphanumeric characters, letter/digit boundaries,
// lower-to-upper transitions and at the end of an acronym.
//
// # Sentences
//
//	nvtext.SentenceCase("hello world. visit example.com!") // Hello world. Visit example.com!
//
// # Timezones
//
//	nvtext.NormalizeTimezone("asia/jakarta") // Asia/Jakarta
//
// NormalizeTimezone fixes the capitalization of a timezone name taken from
// configuration so that time.LoadLocation accepts it.
//
// # Layout
//
//	case.go       tokenize and the *Case functions
//	sentence.go   SentenceCase
//	timezone.go   NormalizeTimezone
package nvtext
