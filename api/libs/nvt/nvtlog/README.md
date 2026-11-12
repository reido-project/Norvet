# nvtlog

Request-error logging and responses for Fiber v3, built on [`nvterr`](../nvterr).

```go
nvtlog.Error(c, err)                      // log only
return nvtlog.Response(c, err, nil)       // log + respond (level from status)
return nvtlog.WarnResponse(c, err, body)  // forced level, custom JSON body
return nvtlog.Response(c, err, nil, 409)  // override the status
```

| Function | Level |
|---|---|
| `Info` / `Warn` / `Error` | fixed, log only |
| `Response` | from the status sent: 5xx `Error`, 4xx `Warn`, else not logged |
| `InfoResponse` / `WarnResponse` / `ErrorResponse` | fixed |

`Response*` send `body` if given, else the `nvterr` JSON of the error, and set
`Retry-After` when the error has one. The internal cause appears in the log,
never in the response. `c` may be `nil` for logging (background jobs); the
route then shows as `UNKNOWN`.

Console output is set in `nvt.Config`:

| Option | Effect |
|---|---|
| `NvtLogColor` | colors the line; `NO_COLOR=1` still disables it |
| `NvtLogTime` | prefixes the time in the app timezone: `16:45 -  [Warn] Status 422 in POST /auth/register endpoint: validation failed` |

Colors are not tied to a terminal check, so they also work in IDE run windows;
turn `NvtLogColor` off when output goes to a file or log collector. In
diagnostic mode entries are also written to `<LogOutput>/nvtlog.sqlite` (table
`logs`; timestamps are fixed-width RFC 3339 in the app timezone, so they sort as
text), without colors or the time prefix.

Needs the revised `nvterr` (`IsKnown`, `RetryAfterSeconds`).
