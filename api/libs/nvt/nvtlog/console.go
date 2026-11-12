package nvtlog

import (
	"fmt"
	"os"
	"time"
)

// colors holds the ANSI escape sequences for console output. The zero value
// is the plain set: every sequence is empty, so the same formatting code
// produces uncolored text.
type colors struct {
	reset, red, green, yellow, magenta, lightBlue, white, bold string
}

// ansiColors uses the bright (high-intensity) ANSI colors so the output stays
// readable on dark terminal themes.
var ansiColors = colors{
	reset:     "[0m",
	red:       "[91m",
	green:     "[92m",
	yellow:    "[93m",
	magenta:   "[95m",
	lightBlue: "[94m",
	white:     "[97m",
	bold:      "[1m",
}

// colorsFor returns the color set for the console: ANSI when the application
// enabled colors (nvt.Config.NvtLogColor), plain otherwise. The NO_COLOR
// environment variable (any non-empty value) always disables colors.
func colorsFor(s settings) colors {
	if useColor(s.color, os.Getenv) {
		return ansiColors
	}
	return colors{}
}

// useColor reports whether console output should be colored. It does not
// check for a terminal: IDE run windows are not terminals but do render
// colors, so the application decides through its configuration.
func useColor(enabled bool, getenv func(string) string) bool {
	return enabled && getenv("NO_COLOR") == ""
}

// line renders the console line for e:
//
//	[Warn] Status 404 in GET /users/:id endpoint: user not found
func (p colors) line(e entry) string {
	endpointPart := fmt.Sprintf("%s %s %sendpoint: %s%v%s",
		p.method(e.method), p.bold+p.white+e.endpoint+p.reset, p.white, p.white, e.err, p.reset,
	)

	message := endpointPart
	if e.status != 0 {
		message = fmt.Sprintf("%s %s %s", p.status(e.status), p.white+"in"+p.reset, endpointPart)
	}

	return fmt.Sprintf("%s%s%s %s", p.white+"["+p.reset, p.level(e.level), p.white+"]"+p.reset, message)
}

// clock renders the time prefix of a console line: "16:45 -  ". The time is
// printed in t's own timezone.
func (p colors) clock(t time.Time) string {
	return p.white + t.Format("15:04") + p.reset + " - "
}

// plainMessage renders the message that is persisted: the console line without
// the level (stored in its own column) and without colors.
func plainMessage(e entry) string {
	message := fmt.Sprintf("%s %s endpoint: %v", e.method, e.endpoint, e.err)
	if e.status != 0 {
		message = fmt.Sprintf("Status %d in %s", e.status, message)
	}
	return message
}

func (p colors) level(lvl level) string {
	color := p.white
	switch lvl {
	case levelInfo:
		color = p.lightBlue
	case levelWarn:
		color = p.yellow
	case levelError:
		color = p.red
	}
	return color + string(lvl) + p.reset
}

func (p colors) method(method string) string {
	color := p.white
	switch method {
	case "GET":
		color = p.green
	case "POST":
		color = p.magenta
	case "DELETE":
		color = p.red
	case "PUT", "PATCH":
		color = p.yellow
	}
	return color + method + p.reset
}

// status renders "Status <code>": the word stays white, the number is colored
// by class (1xx/3xx light blue, 2xx green, 4xx yellow, 5xx magenta).
func (p colors) status(status int) string {
	color := p.white
	switch {
	case status >= 500:
		color = p.magenta
	case status >= 400:
		color = p.yellow
	case status >= 300, status < 200:
		color = p.lightBlue
	default:
		color = p.green
	}
	return fmt.Sprintf("%sStatus%s %s%d%s", p.white, p.reset, color, status, p.reset)
}
