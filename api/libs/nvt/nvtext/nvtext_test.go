package nvtext

import (
	"archive/zip"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestCaseConversion(t *testing.T) {
	tests := []struct {
		in                          string
		pascal, camel, snake, kebab string
		title                       string
	}{
		{"user_id", "UserId", "userId", "user_id", "user-id", "User Id"},
		{"UserID", "UserId", "userId", "user_id", "user-id", "User Id"},
		{"HTTPServer", "HttpServer", "httpServer", "http_server", "http-server", "Http Server"},
		{"user2name", "User2Name", "user2Name", "user_2_name", "user-2-name", "User 2 Name"},
		{"  hello   world  ", "HelloWorld", "helloWorld", "hello_world", "hello-world", "Hello World"},
		{"çok güzel", "ÇokGüzel", "çokGüzel", "çok_güzel", "çok-güzel", "Çok Güzel"},
		{"", "", "", "", "", ""},
		{"___", "", "", "", "", ""},
	}
	for _, tt := range tests {
		if got := PascalCase(tt.in); got != tt.pascal {
			t.Errorf("PascalCase(%q) = %q, want %q", tt.in, got, tt.pascal)
		}
		if got := CamelCase(tt.in); got != tt.camel {
			t.Errorf("CamelCase(%q) = %q, want %q", tt.in, got, tt.camel)
		}
		if got := SnakeCase(tt.in); got != tt.snake {
			t.Errorf("SnakeCase(%q) = %q, want %q", tt.in, got, tt.snake)
		}
		if got := KebabCase(tt.in); got != tt.kebab {
			t.Errorf("KebabCase(%q) = %q, want %q", tt.in, got, tt.kebab)
		}
		if got := TitleCase(tt.in); got != tt.title {
			t.Errorf("TitleCase(%q) = %q, want %q", tt.in, got, tt.title)
		}
	}
}

func TestSentenceCase(t *testing.T) {
	tests := []struct{ in, want string }{
		{"hello world. this is go! really? yes", "Hello world. This is go! Really? Yes"},
		{"note: this matters; so does this", "Note: This matters; So does this"},
		{"visit example.com now", "Visit example.com now"},
		{"it costs 3.14 euros", "It costs 3.14 euros"},
		{"3 apples are here. 4 pears too", "3 apples are here. 4 pears too"},
		{`"quoted" start. (bracket) next`, `"Quoted" start. (Bracket) next`},
		{"wait... what", "Wait... What"},
		{"ends with a dot.", "Ends with a dot."},
		{"  leading space", "  Leading space"},
		{"éclair. ünder", "Éclair. Ünder"},
		{"   ", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := SentenceCase(tt.in); got != tt.want {
			t.Errorf("SentenceCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestNormalizeTimezone(t *testing.T) {
	tests := []struct{ in, want string }{
		{"asia/jakarta", "Asia/Jakarta"},
		{"  AMERICA/NEW_YORK ", "America/New_York"},
		{"america/argentina/buenos_aires", "America/Argentina/Buenos_Aires"},
		{"america/port_of_spain", "America/Port_of_Spain"},
		{"africa/dar_es_salaam", "Africa/Dar_es_Salaam"},
		{"america/port-au-prince", "America/Port-au-Prince"},
		{"asia/ust-nera", "Asia/Ust-Nera"},
		{"etc/gmt+5", "Etc/GMT+5"},
		{"etc/utc", "Etc/UTC"},
		{"us/pacific", "US/Pacific"},
		{"australia/nsw", "Australia/NSW"},
		{"utc", "UTC"},
		{"gmt", "GMT"},
		{"est5edt", "EST5EDT"},
		{"", ""},
		{"   ", ""},
	}
	for _, tt := range tests {
		if got := NormalizeTimezone(tt.in); got != tt.want {
			t.Errorf("NormalizeTimezone(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// TestNormalizeTimezone_RoundTripsKnownZones lower-cases every zone name known
// to Go's timezone database and checks it is repaired, except the names whose
// capitalization is not derivable by rule (mixed-case ones).
func TestNormalizeTimezone_RoundTripsKnownZones(t *testing.T) {
	archive, err := zip.OpenReader(filepath.Join(runtime.GOROOT(), "lib", "time", "zoneinfo.zip"))
	if err != nil {
		t.Skip("zoneinfo.zip not available:", err)
	}
	defer archive.Close()

	var failed []string
	checked := 0
	for _, f := range archive.File {
		name := f.Name
		if f.FileInfo().IsDir() {
			continue
		}
		checked++

		// A correctly written name is always returned unchanged.
		if got := NormalizeTimezone(name); got != name {
			t.Errorf("NormalizeTimezone(%q) changed a valid name to %q", name, got)
		}

		if got := NormalizeTimezone(strings.ToLower(name)); got != name {
			failed = append(failed, name+" -> "+got)
		}
	}

	// Mixed-case zone names (McMurdo, DumontDUrville, ComodRivadavia, ...) and
	// a few legacy aliases cannot be recovered from lower case.
	if limit := 12; len(failed) > limit {
		t.Errorf("%d of %d zones not repaired from lower case (limit %d):\n%s",
			len(failed), checked, limit, strings.Join(failed, "\n"))
	}
	t.Logf("repaired %d of %d zones from lower case; unrecoverable: %v", checked-len(failed), checked, failed)
}

func TestNormalizeTimezone_ResultLoads(t *testing.T) {
	for _, in := range []string{"asia/jakarta", "europe/isle_of_man", "etc/gmt-14"} {
		if _, err := time.LoadLocation(NormalizeTimezone(in)); err != nil {
			t.Errorf("%q normalizes to a zone that does not load: %v", in, err)
		}
	}
}
