// Package nvtlog logs request errors to the console and answers Fiber requests
// with them.
//
// Log only:
//
//	nvtlog.Warn(c, err)        // [Warn] Status 404 in GET /users/:id endpoint: user not found
//
// Log and respond, in one call:
//
//	return nvtlog.Response(c, err, nil)
//
// Response picks the level from the status (5xx error, 4xx warn, otherwise not
// logged) and sends the nvterr representation of err, so internal causes show
// in the log but never in the response. InfoResponse, WarnResponse and
// ErrorResponse force a level. An optional trailing code overrides the status.
//
// Console output is configured through nvt.Config:
//
//	NvtLogColor   colors the line (NO_COLOR=1 still disables it)
//	NvtLogTime    prefixes the line with the time: "16:45 -  [Warn] Status 422 ..."
//
// In diagnostic mode (nvt context) every entry is also stored in
// <LogOutput>/nvtlog.sqlite, without colors or the time prefix.
//
// # Layout
//
//	log.go        Info, Warn, Error
//	response.go   Response and friends; plan decides, response executes
//	status.go     status resolution
//	level.go      levels
//	writer.go     entry, emit, request resolution
//	console.go    line format, colors
//	store.go      SQLite persistence
//	settings.go   the only place that reads the nvt context
package nvtlog
