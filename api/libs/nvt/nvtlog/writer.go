package nvtlog

import (
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v3"
)

// entry is one thing to log.
type entry struct {
	level            level
	method, endpoint string
	status           int // 0 = unknown, omitted from the line
	err              error
}

// Seams that tests replace.
var (
	out           io.Writer = os.Stdout
	consoleColors           = colorsFor
	loadSettings            = currentSettings
	persist                 = saveToSQLite
)

// write logs err for the request ctx.
func write(ctx fiber.Ctx, err error, lvl level, status int) {
	if err == nil {
		return
	}
	method, endpoint := resolveRequest(ctx)
	emit(entry{level: lvl, method: method, endpoint: endpoint, status: status, err: err})
}

// emit prints the entry to the console and, in diagnostic mode, persists it.
// Persisted messages never contain colors or the console time prefix.
func emit(e entry) {
	s := loadSettings()
	p := consoleColors(s)

	line := p.line(e)
	if s.showTime {
		line = p.clock(s.now()) + line
	}
	_, _ = fmt.Fprintln(out, line)

	if s.diagnostic {
		persist(s, e)
	}
}

// resolveRequest extracts the method and route pattern from Fiber. The route
// pattern ("/users/:id") is used rather than the raw path so no identifiers
// from the URL end up in the log.
func resolveRequest(ctx fiber.Ctx) (method, endpoint string) {
	method, endpoint = "UNKNOWN", "UNKNOWN"

	if ctx == nil {
		return
	}
	if value := ctx.Method(); value != "" {
		method = value
	}
	if route := ctx.Route(); route != nil && route.Path != "" {
		endpoint = route.Path
	}
	return
}
