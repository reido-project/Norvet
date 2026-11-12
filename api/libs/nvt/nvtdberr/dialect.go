package nvtdberr

import "api/libs/nvt"

// activeDialect returns the database dialect configured in the nvt context. ok
// is false when nvt has not been initialized. It is the only place in the
// package that reads the nvt context.
func activeDialect() (dialect nvt.Dialect, ok bool) {
	c := nvt.Ctx()
	if c == nil {
		return "", false
	}
	return c.Database, true
}
