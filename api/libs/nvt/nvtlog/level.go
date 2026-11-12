package nvtlog

// level is the severity of a log entry. The string is what appears in the
// console line and in the persisted log.
type level string

const (
	levelInfo  level = "Info"
	levelWarn  level = "Warn"
	levelError level = "Error"

	// levelAuto asks response to pick the level from the status being sent. It
	// is never printed or stored.
	levelAuto level = "auto"
)

// levelFor picks the level for a response status: 5xx is an error, 4xx a
// warning. Anything else is not worth logging and gets "".
func levelFor(status int) level {
	switch {
	case status >= 500:
		return levelError
	case status >= 400:
		return levelWarn
	}
	return ""
}
