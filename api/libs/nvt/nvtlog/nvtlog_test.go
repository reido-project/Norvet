package nvtlog

import (
	"api/libs/nvt/nvterr"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

// capture replaces the console, color and persistence seams for one test.
type capture struct {
	console bytes.Buffer
	saved   []entry
	msgs    []string
}

func newCapture(t *testing.T, diagnostic bool, c colors) *capture {
	t.Helper()
	cp := &capture{}
	oldOut, oldColors, oldSettings, oldPersist := out, consoleColors, loadSettings, persist
	out = &cp.console
	consoleColors = func(settings) colors { return c }
	loadSettings = func() settings { return settings{diagnostic: diagnostic} }
	persist = func(_ settings, e entry) {
		cp.saved = append(cp.saved, e)
		cp.msgs = append(cp.msgs, plainMessage(e))
	}
	t.Cleanup(func() { out, consoleColors, loadSettings, persist = oldOut, oldColors, oldSettings, oldPersist })
	return cp
}

// Golden strings: the console format and color scheme as specified by the owner.
func TestConsoleLine_Golden(t *testing.T) {
	tests := []struct {
		name    string
		e       entry
		console string
		saved   string
	}{
		{"info, no status",
			entry{levelInfo, "GET", "/users/:id", 0, errors.New("boom")},
			"\x1b[97m[\x1b[0m\x1b[94mInfo\x1b[0m\x1b[97m]\x1b[0m \x1b[92mGET\x1b[0m \x1b[1m\x1b[97m/users/:id\x1b[0m \x1b[97mendpoint: \x1b[97mboom\x1b[0m",
			"GET /users/:id endpoint: boom"},
		{"warn 404 POST",
			entry{levelWarn, "POST", "/a", 404, nvterr.ErrNotFound},
			"\x1b[97m[\x1b[0m\x1b[93mWarn\x1b[0m\x1b[97m]\x1b[0m \x1b[97mStatus\x1b[0m \x1b[93m404\x1b[0m \x1b[97min\x1b[0m \x1b[95mPOST\x1b[0m \x1b[1m\x1b[97m/a\x1b[0m \x1b[97mendpoint: \x1b[97mresource not found\x1b[0m",
			"Status 404 in POST /a endpoint: resource not found"},
		{"error 500 DELETE includes the cause",
			entry{levelError, "DELETE", "/a", 500, nvterr.ErrInternal.WithCause(errors.New("db down"))},
			"\x1b[97m[\x1b[0m\x1b[91mError\x1b[0m\x1b[97m]\x1b[0m \x1b[97mStatus\x1b[0m \x1b[95m500\x1b[0m \x1b[97min\x1b[0m \x1b[91mDELETE\x1b[0m \x1b[1m\x1b[97m/a\x1b[0m \x1b[97mendpoint: \x1b[97minternal server error: db down\x1b[0m",
			"Status 500 in DELETE /a endpoint: internal server error: db down"},
		{"info 204 PUT",
			entry{levelInfo, "PUT", "/a", 204, errors.New("ok")},
			"\x1b[97m[\x1b[0m\x1b[94mInfo\x1b[0m\x1b[97m]\x1b[0m \x1b[97mStatus\x1b[0m \x1b[92m204\x1b[0m \x1b[97min\x1b[0m \x1b[93mPUT\x1b[0m \x1b[1m\x1b[97m/a\x1b[0m \x1b[97mendpoint: \x1b[97mok\x1b[0m",
			"Status 204 in PUT /a endpoint: ok"},
		{"warn 302 unknown method is white",
			entry{levelWarn, "OPTIONS", "/a", 302, errors.New("x")},
			"\x1b[97m[\x1b[0m\x1b[93mWarn\x1b[0m\x1b[97m]\x1b[0m \x1b[97mStatus\x1b[0m \x1b[94m302\x1b[0m \x1b[97min\x1b[0m \x1b[97mOPTIONS\x1b[0m \x1b[1m\x1b[97m/a\x1b[0m \x1b[97mendpoint: \x1b[97mx\x1b[0m",
			"Status 302 in OPTIONS /a endpoint: x"},
		{"no request context",
			entry{levelError, "UNKNOWN", "UNKNOWN", 0, errors.New("job failed")},
			"\x1b[97m[\x1b[0m\x1b[91mError\x1b[0m\x1b[97m]\x1b[0m \x1b[97mUNKNOWN\x1b[0m \x1b[1m\x1b[97mUNKNOWN\x1b[0m \x1b[97mendpoint: \x1b[97mjob failed\x1b[0m",
			"UNKNOWN UNKNOWN endpoint: job failed"},
		{"the documented example",
			entry{levelError, "POST", "/auth/register", 500, nvterr.ErrInternal},
			"\x1b[97m[\x1b[0m\x1b[91mError\x1b[0m\x1b[97m]\x1b[0m \x1b[97mStatus\x1b[0m \x1b[95m500\x1b[0m \x1b[97min\x1b[0m \x1b[95mPOST\x1b[0m \x1b[1m\x1b[97m/auth/register\x1b[0m \x1b[97mendpoint: \x1b[97minternal server error\x1b[0m",
			"Status 500 in POST /auth/register endpoint: internal server error"},
	}
	ansiRe := regexp.MustCompile("\x1b\\[[0-9;]*m")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ansiColors.line(tt.e); got != tt.console {
				t.Errorf("colored line\n got  %q\n want %q", got, tt.console)
			}
			if got := plainMessage(tt.e); got != tt.saved {
				t.Errorf("persisted message\n got  %q\n want %q", got, tt.saved)
			}
			// The uncolored line is the colored one without its escape codes.
			if got, want := (colors{}).line(tt.e), ansiRe.ReplaceAllString(tt.console, ""); got != want {
				t.Errorf("plain line\n got  %q\n want %q", got, want)
			}
			if want := "[" + string(tt.e.level) + "] " + tt.saved; (colors{}).line(tt.e) != want {
				t.Errorf("plain line must be level + persisted message:\n got  %q\n want %q", (colors{}).line(tt.e), want)
			}
		})
	}
}

func TestUseColor(t *testing.T) {
	env := func(kv ...string) func(string) string {
		m := map[string]string{}
		for i := 0; i < len(kv); i += 2 {
			m[kv[i]] = kv[i+1]
		}
		return func(k string) string { return m[k] }
	}
	tests := []struct {
		name    string
		enabled bool
		env     func(string) string
		want    bool
	}{
		{"enabled", true, env(), true},
		{"disabled by config", false, env(), false},
		{"NO_COLOR beats config", true, env("NO_COLOR", "1"), false},
		{"empty NO_COLOR is ignored", true, env("NO_COLOR", ""), true},
		{"config off stays off", false, env("FORCE_COLOR", "1"), false},
	}
	for _, tt := range tests {
		if got := useColor(tt.enabled, tt.env); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestClock(t *testing.T) {
	at := time.Date(2026, 11, 12, 16, 45, 9, 0, time.UTC)
	if got, want := (colors{}).clock(at), "16:45 -  "; got != want {
		t.Errorf("plain clock = %q, want %q", got, want)
	}
	if got, want := ansiColors.clock(at), "\x1b[97m16:45\x1b[0m -  "; got != want {
		t.Errorf("colored clock = %q, want %q", got, want)
	}
}

func TestEmit_TimePrefixIsConsoleOnly(t *testing.T) {
	cp := newCapture(t, true, colors{})
	loadSettings = func() settings { return settings{diagnostic: true, showTime: true} }

	Warn(nil, nvterr.ErrNotFound)

	got := strings.TrimSuffix(cp.console.String(), "\n")
	re := regexp.MustCompile(`^\d\d:\d\d -  \[Warn\] Status 404 in UNKNOWN UNKNOWN endpoint: resource not found$`)
	if !re.MatchString(got) {
		t.Errorf("console line = %q", got)
	}
	if len(cp.msgs) != 1 || cp.msgs[0] != "Status 404 in UNKNOWN UNKNOWN endpoint: resource not found" {
		t.Errorf("persisted message must have no time prefix or colors: %q", cp.msgs)
	}
}

func TestEmit_NoTimePrefixByDefault(t *testing.T) {
	cp := newCapture(t, false, colors{})
	Warn(nil, nvterr.ErrNotFound)
	if got, want := cp.console.String(), "[Warn] Status 404 in UNKNOWN UNKNOWN endpoint: resource not found\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPublicLoggers(t *testing.T) {
	cp := newCapture(t, false, colors{})

	Info(nil, errors.New("a"))
	Warn(nil, nvterr.ErrNotFound)
	Error(nil, errors.New("c"), 503)
	Info(nil, nil) // nothing to log
	Warn(nil, nil, 400)

	want := "[Info] UNKNOWN UNKNOWN endpoint: a\n" +
		"[Warn] Status 404 in UNKNOWN UNKNOWN endpoint: resource not found\n" +
		"[Error] Status 503 in UNKNOWN UNKNOWN endpoint: c\n"
	if got := cp.console.String(); got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
	if len(cp.saved) != 0 {
		t.Errorf("nothing may be persisted outside diagnostic mode, got %d entries", len(cp.saved))
	}
}

func TestEmit_PersistsOnlyInDiagnosticMode(t *testing.T) {
	cp := newCapture(t, true, colors{})
	Error(nil, errors.New("x"), 500)

	if len(cp.saved) != 1 {
		t.Fatalf("saved %d entries, want 1", len(cp.saved))
	}
	e := cp.saved[0]
	if e.level != levelError || e.status != 500 || e.method != "UNKNOWN" || cp.msgs[0] != "Status 500 in UNKNOWN UNKNOWN endpoint: x" {
		t.Fatalf("saved %+v / %q", e, cp.msgs[0])
	}
}

func TestResolveRequest_NilContext(t *testing.T) {
	if m, e := resolveRequest(nil); m != "UNKNOWN" || e != "UNKNOWN" {
		t.Fatalf("got %q %q", m, e)
	}
}

func TestResolveStatus(t *testing.T) {
	typedNil := (*nvterr.Error)(nil)
	tests := []struct {
		name     string
		err      error
		code     []int
		fallback int
		want     int
	}{
		{"explicit code wins", nvterr.ErrNotFound, []int{418}, 0, 418},
		{"first code only", nvterr.ErrNotFound, []int{418, 503}, 0, 418},
		{"zero code is ignored", nvterr.ErrNotFound, []int{0}, 0, 404},
		{"negative code is ignored", nvterr.ErrNotFound, []int{-1}, 0, 404},
		{"nvterr status", nvterr.ErrNotFound, nil, 0, 404},
		{"wrapped nvterr", fmt.Errorf("w: %w", nvterr.ErrConflict), nil, 0, 409},
		{"canceled", context.Canceled, nil, 0, 499},
		{"deadline", context.DeadlineExceeded, nil, 0, 504},
		{"unknown error logs no status", errors.New("x"), nil, 0, 0},
		{"unknown error responds 500", errors.New("x"), nil, 500, 500},
		{"nil error uses the fallback", nil, nil, 200, 200},
		{"typed nil does not panic", typedNil, nil, 500, 500},
	}
	for _, tt := range tests {
		if got := resolveStatus(tt.err, tt.code, tt.fallback); got != tt.want {
			t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestLevelFor(t *testing.T) {
	for status, want := range map[int]level{
		200: "", 204: "", 302: "", 399: "",
		400: levelWarn, 404: levelWarn, 499: levelWarn,
		500: levelError, 503: levelError, 599: levelError,
	} {
		if got := levelFor(status); got != want {
			t.Errorf("levelFor(%d) = %q, want %q", status, got, want)
		}
	}
}

func TestPlan(t *testing.T) {
	boom := errors.New("boom")
	body := map[string]any{"ok": false}

	tests := []struct {
		name string
		err  error
		lvl  level
		body map[string]any
		code []int
		want reply
	}{
		{"no error, no body: bare 200", nil, levelAuto, nil, nil, reply{status: 200}},
		{"no error, body", nil, levelAuto, body, nil, reply{status: 200, payload: body}},
		{"no error, code override", nil, levelAuto, nil, []int{201}, reply{status: 201}},
		{"no error is never logged, whatever the level", nil, levelError, nil, nil, reply{status: 200}},
		{"empty non-nil body is still sent", nil, levelAuto, map[string]any{}, nil, reply{status: 200, payload: map[string]any{}}},

		{"nvterr 404 auto-logs as warn", nvterr.ErrNotFound, levelAuto, nil, nil, reply{status: 404, payload: nvterr.ErrNotFound, logAs: levelWarn}},
		{"nvterr 500 auto-logs as error", nvterr.ErrInternal, levelAuto, nil, nil, reply{status: 500, payload: nvterr.ErrInternal, logAs: levelError}},
		{"explicit level beats the status", nvterr.ErrInternal, levelInfo, nil, nil, reply{status: 500, payload: nvterr.ErrInternal, logAs: levelInfo}},
		{"override picks the level of the status sent", nvterr.ErrNotFound, levelAuto, nil, []int{500}, reply{status: 500, payload: nvterr.ErrNotFound, logAs: levelError}},
		{"override to a success status logs nothing", nvterr.ErrNotFound, levelAuto, nil, []int{200}, reply{status: 200, payload: nvterr.ErrNotFound}},
		{"custom body replaces the error body", nvterr.ErrNotFound, levelAuto, body, nil, reply{status: 404, payload: body, logAs: levelWarn}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plan(tt.err, tt.lvl, tt.body, tt.code)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got  %+v\nwant %+v", got, tt.want)
			}
		})
	}

	// A plain error answers 500 with the generic nvterr body, never its text.
	r := plan(boom, levelAuto, nil, nil)
	e, ok := r.payload.(*nvterr.Error)
	if !ok || r.status != http.StatusInternalServerError || r.logAs != levelError {
		t.Fatalf("plain error: %+v", r)
	}
	if e.Message() != "internal server error" {
		t.Errorf("client message leaked the internal error: %q", e.Message())
	}
	if raw, _ := e.MarshalJSON(); strings.Contains(string(raw), "boom") {
		t.Errorf("client JSON leaked the internal error: %s", raw)
	}
}

func TestPlan_RetryAfter(t *testing.T) {
	for _, tt := range []struct {
		err  error
		want int
	}{
		{nvterr.NewRateLimited(1500 * time.Millisecond), 2},
		{nvterr.NewTooManyAttempts(400 * time.Millisecond), 1},
		{nvterr.NewRateLimited(30 * time.Second), 30},
		{nvterr.ErrRateLimited, 0},
		{nvterr.ErrNotFound, 0},
		{errors.New("x"), 0},
	} {
		if got := plan(tt.err, levelAuto, nil, nil).retryAfter; got != tt.want {
			t.Errorf("%v: retryAfter = %d, want %d", tt.err, got, tt.want)
		}
	}
}

func TestPlainLine_MatchesTheDocumentedFormat(t *testing.T) {
	e := entry{levelError, "POST", "/auth/register", 500, nvterr.ErrInternal}
	if got, want := (colors{}).line(e), "[Error] Status 500 in POST /auth/register endpoint: internal server error"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStatusColors(t *testing.T) {
	// Only the number is colored; the word "Status" stays white.
	for status, color := range map[int]string{
		100: "\x1b[94m", 302: "\x1b[94m", // light blue
		200: "\x1b[92m", 204: "\x1b[92m", // green
		404: "\x1b[93m", // yellow
		500: "\x1b[95m", // magenta
	} {
		want := "\x1b[97mStatus\x1b[0m " + color + fmt.Sprint(status) + "\x1b[0m"
		if got := ansiColors.status(status); got != want {
			t.Errorf("status %d: got %q, want %q", status, got, want)
		}
	}
}
