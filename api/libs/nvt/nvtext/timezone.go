package nvtext

import (
	"strings"
	"time"
)

// upperNames are the zone names, or parts of a zone name between slashes,
// that the IANA database writes entirely in capitals ("UTC", "EST5EDT",
// "Etc/UCT", "US/Pacific", "Australia/NSW").
var upperNames = map[string]bool{
	"utc": true, "gmt": true, "uct": true,
	"est": true, "mst": true, "hst": true,
	"cet": true, "eet": true, "met": true, "wet": true,
	"est5edt": true, "cst6cdt": true, "mst7mdt": true, "pst8pdt": true,
	"gb": true, "nz": true, "prc": true, "roc": true, "rok": true,
	"us": true, "act": true, "nsw": true, "lhi": true, "nz-chat": true, "w-su": true,
}

// lowerWords are the words that stay lower-case inside a zone name
// ("America/Port_of_Spain", "Africa/Dar_es_Salaam", "America/Port-au-Prince").
var lowerWords = map[string]bool{"of": true, "es": true, "au": true}

// NormalizeTimezone repairs the capitalization of an IANA timezone name so it
// can be passed to time.LoadLocation, which is case-sensitive:
//
//	NormalizeTimezone("asia/jakarta")        // "Asia/Jakarta"
//	NormalizeTimezone("AMERICA/NEW_YORK")    // "America/New_York"
//	NormalizeTimezone("america/port_of_spain") // "America/Port_of_Spain"
//	NormalizeTimezone("etc/gmt+5")           // "Etc/GMT+5"
//	NormalizeTimezone("utc")                 // "UTC"
//
// A name that time.LoadLocation already accepts is returned unchanged, which
// keeps unusual spellings such as "Antarctica/DumontDUrville" intact. Any
// other input is title-cased by rule, so mixed-case names written in the wrong
// case ("antarctica/mcmurdo" gives "Antarctica/Mcmurdo") are not recoverable;
// the caller finds out when LoadLocation fails. Surrounding spaces are
// trimmed and an empty input gives "".
func NormalizeTimezone(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	if _, err := time.LoadLocation(input); err == nil {
		return input
	}

	parts := strings.Split(input, "/")
	for i, part := range parts {
		parts[i] = titleZonePart(part)
	}
	return strings.Join(parts, "/")
}

// titleZonePart capitalizes one "/"-separated part of a zone name: each word
// separated by "_" or "-" starts with a capital, except the lowerWords after
// the first. The upperNames and offset names ("GMT+5", "GMT0") are upper-cased
// as a whole.
func titleZonePart(part string) string {
	lower := strings.ToLower(part)

	if upperNames[lower] {
		return strings.ToUpper(lower)
	}
	if strings.HasPrefix(lower, "gmt") && (len(lower) == 3 || strings.ContainsRune("+-0123456789", rune(lower[3]))) {
		return strings.ToUpper(lower)
	}

	var sb strings.Builder
	sb.Grow(len(lower))

	start := 0
	for i := 0; i <= len(lower); i++ {
		if i < len(lower) && lower[i] != '_' && lower[i] != '-' {
			continue
		}
		word := lower[start:i]
		if start == 0 || !lowerWords[word] {
			word = capitalize(word)
		}
		sb.WriteString(word)
		if i < len(lower) {
			sb.WriteByte(lower[i])
		}
		start = i + 1
	}
	return sb.String()
}
