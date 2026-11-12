package nvtdberr

import (
	"log"
	"sync"
)

// warned remembers which warnings were already logged. Its keys come from a
// bounded set (error types, index names from the schema, dialect names), so it
// cannot grow without limit.
var warned sync.Map

// warnOnce logs a warning the first time key is seen, so a hot path cannot
// flood the log.
func warnOnce(key, format string, args ...any) {
	if _, loaded := warned.LoadOrStore(key, struct{}{}); !loaded {
		log.Printf(format, args...)
	}
}
