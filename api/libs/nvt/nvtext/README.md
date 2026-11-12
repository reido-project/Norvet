# nvtext

Small text helpers: identifier case conversion, sentence capitalization and timezone-name repair. Standard library only.

## Case conversion

Input can be in any style; the functions split it into words first.

```go
nvtext.PascalCase("user_id") // UserId
nvtext.CamelCase("user_id")  // userId
nvtext.SnakeCase("UserID")   // user_id
nvtext.KebabCase("UserID")   // user-id
nvtext.TitleCase("user_id")  // User Id
```

Words break at non-alphanumeric characters, between letters and digits (`user2name` gives `user`, `2`, `name`), at lower-to-upper transitions and at the end of an acronym (`HTTPServer` gives `http`, `server`). Apostrophes split too: `don't` gives `don`, `t`.

## Sentences

```go
nvtext.SentenceCase("hello world. visit example.com! it costs 3.14")
// Hello world. Visit example.com! It costs 3.14
```

Capitalizes the first letter after `.`, `?`, `!`, `;` or `:` followed by whitespace or the end of the text, and leaves everything else alone. A sentence starting with a digit is unchanged. Punctuation inside a token (`example.com`, `3.14`) is not a boundary.

## Timezones

```go
nvtext.NormalizeTimezone("asia/jakarta")          // Asia/Jakarta
nvtext.NormalizeTimezone("america/port_of_spain") // America/Port_of_Spain
nvtext.NormalizeTimezone("etc/gmt+5")             // Etc/GMT+5
nvtext.NormalizeTimezone("utc")                   // UTC
```

`time.LoadLocation` is case-sensitive, so a timezone typed into `.env` in the wrong case fails. `NormalizeTimezone` fixes the case. A name that `time.LoadLocation` already accepts is returned unchanged. Mixed-case names written in the wrong case (`antarctica/mcmurdo`) cannot be recovered by rule and still fail at `LoadLocation`.

## Renamed from the previous version

| Before | Now |
|---|---|
| `Title` | `TitleCase` |
| `Sentence` | `SentenceCase` |
| `Timezone` | `NormalizeTimezone` |
