package nvt

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"api/libs/nvt/nvtext"
)

// ErrNotInitialized is returned by packages that need the Norvet context when
// Init has not been called yet.
var ErrNotInitialized = errors.New("nvt: Init has not been called")

// Config is the application-level input to Init. Everything else (database
// driver, timezones) is read from the environment.
type Config struct {
	// Diagnostic turns on diagnostic features, such as persisting nvtlog
	// entries to SQLite.
	Diagnostic bool

	// LogOutput is the directory where diagnostic files are written.
	LogOutput string

	// NvtLogColor colors nvtlog's console output. The NO_COLOR environment
	// variable still disables it. Turn it off when output goes to a file or a
	// log collector.
	NvtLogColor bool

	// NvtLogTime prefixes nvtlog's console lines with the time (15:04) in the
	// application timezone.
	NvtLogTime bool
}

// Context is the immutable, process-wide Norvet configuration built by Init.
// Treat it as read-only.
type Context struct {
	Diagnostic  bool
	LogOutput   string
	NvtLogColor bool
	NvtLogTime  bool

	// AppTimezone is the timezone of the application (APP_TIMEZONE), with its
	// name in AppTimezoneStr. It is the local timezone when unset or invalid.
	AppTimezone    *time.Location
	AppTimezoneStr string

	// Database is the SQL dialect selected by DB_DRIVER; empty when unset or
	// unsupported.
	Database Dialect

	// DBTimezone is the timezone used for database timestamps (DB_TIMEZONE),
	// with its name in DBTimezoneStr. It is the application timezone when
	// unset or invalid.
	DBTimezone    *time.Location
	DBTimezoneStr string
}

var (
	once     sync.Once
	current  atomic.Pointer[Context]
	warnings []error
)

// Init builds the Norvet context from config and the environment
// (DB_DRIVER, APP_TIMEZONE, DB_TIMEZONE). Call it once at startup, before
// anything that uses Ctx.
//
// Init never fails; problems with the environment fall back to defaults and
// are returned as warnings so the caller can report them. Only the first call
// has any effect: later calls ignore config and return the same warnings.
func Init(config Config) []error {
	once.Do(func() {
		var w []error

		dialect, ok := ParseDialect(os.Getenv("DB_DRIVER"))
		if raw := os.Getenv("DB_DRIVER"); !ok && raw != "" {
			w = append(w, fmt.Errorf("nvt: unsupported DB_DRIVER %q", raw))
		}

		appTz := nvtext.NormalizeTimezone(os.Getenv("APP_TIMEZONE"))
		appLoc, err := time.LoadLocation(appTz)
		if err != nil {
			w = append(w, fmt.Errorf("nvt: cannot load APP_TIMEZONE, using local timezone: %w", err))
			appLoc, appTz = time.Local, "Local"
		}

		dbTz := nvtext.NormalizeTimezone(os.Getenv("DB_TIMEZONE"))
		dbLoc, err := time.LoadLocation(dbTz)
		if err != nil {
			w = append(w, fmt.Errorf("nvt: cannot load DB_TIMEZONE, using the application timezone: %w", err))
			dbLoc, dbTz = appLoc, appTz
		}

		warnings = w
		current.Store(&Context{
			Diagnostic:  config.Diagnostic,
			LogOutput:   config.LogOutput,
			NvtLogColor: config.NvtLogColor,
			NvtLogTime:  config.NvtLogTime,

			AppTimezone:    appLoc,
			AppTimezoneStr: appTz,

			Database:      dialect,
			DBTimezone:    dbLoc,
			DBTimezoneStr: dbTz,
		})
	})

	return slices.Clone(warnings)
}

// Ctx returns the context built by Init, or nil if Init has not been called.
// It is safe for concurrent use.
func Ctx() *Context {
	return current.Load()
}
