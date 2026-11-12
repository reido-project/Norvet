package nvtlog

import (
	"api/libs/nvt"
	"time"
)

// settings is the part of the nvt context that nvtlog depends on. It is the
// only place in the package that reads the nvt context, so changing how nvt
// is configured touches this file alone.
type settings struct {
	diagnostic bool           // persist entries to SQLite
	logDir     string         // directory of the SQLite file
	location   *time.Location // timezone for timestamps; nil = local
	color      bool           // color the console output
	showTime   bool           // prefix console lines with the time
}

func currentSettings() settings {
	c := nvt.Ctx()
	if c == nil {
		return settings{}
	}
	return settings{
		diagnostic: c.Diagnostic,
		logDir:     c.LogOutput,
		location:   c.AppTimezone,
		color:      c.NvtLogColor,
		showTime:   c.NvtLogTime,
	}
}

// now returns the current time in the application timezone.
func (s settings) now() time.Time {
	if s.location != nil {
		return time.Now().In(s.location)
	}
	return time.Now()
}
