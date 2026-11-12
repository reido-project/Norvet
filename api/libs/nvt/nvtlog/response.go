package nvtlog

import (
	"net/http"
	"strconv"

	"api/libs/nvt/nvterr"

	"github.com/gofiber/fiber/v3"
)

// Response logs err and answers the request with it, choosing the log level
// from the response status: 5xx logs at error level, 4xx at warn level, and
// anything else is not logged.
//
//   - err == nil: responds with body as JSON, or with just the status when body
//     is nil. The status is the optional code, or 200.
//   - err != nil: responds with the optional code, else the nvterr status, else
//     500. The JSON body is body when given, otherwise the nvterr
//     representation of err (code, message, violations, meta), which never
//     contains the internal cause. A Retry-After header is set when the error
//     carries one.
//
// The internal cause appears in the log line but never in the response.
func Response(ctx fiber.Ctx, err error, body map[string]any, code ...int) error {
	return response(ctx, err, levelAuto, body, code)
}

// InfoResponse is Response, but always logs at info level.
func InfoResponse(ctx fiber.Ctx, err error, body map[string]any, code ...int) error {
	return response(ctx, err, levelInfo, body, code)
}

// WarnResponse is Response, but always logs at warn level.
func WarnResponse(ctx fiber.Ctx, err error, body map[string]any, code ...int) error {
	return response(ctx, err, levelWarn, body, code)
}

// ErrorResponse is Response, but always logs at error level.
func ErrorResponse(ctx fiber.Ctx, err error, body map[string]any, code ...int) error {
	return response(ctx, err, levelError, body, code)
}

// reply is what to send and log for one request. plan decides it without
// touching the framework, so the rules can be tested on their own.
type reply struct {
	status     int
	payload    any   // JSON body; nil means send the status alone
	retryAfter int   // Retry-After seconds; 0 means no header
	logAs      level // "" means do not log
}

func plan(err error, lvl level, body map[string]any, code []int) reply {
	if err == nil {
		r := reply{status: resolveStatus(nil, code, http.StatusOK)}
		if body != nil {
			r.payload = body
		}
		return r
	}

	r := reply{status: resolveStatus(err, code, http.StatusInternalServerError), logAs: lvl}
	if lvl == levelAuto {
		r.logAs = levelFor(r.status) // based on the status actually sent, including an override
	}

	e := nvterr.From(err)
	r.retryAfter = e.RetryAfterSeconds()
	if body != nil {
		r.payload = body
	} else {
		r.payload = e
	}
	return r
}

func response(ctx fiber.Ctx, err error, lvl level, body map[string]any, code []int) error {
	r := plan(err, lvl, body, code)

	if r.logAs != "" {
		write(ctx, err, r.logAs, r.status)
	}
	if r.retryAfter > 0 {
		ctx.Set(fiber.HeaderRetryAfter, strconv.Itoa(r.retryAfter))
	}
	if r.payload == nil {
		return ctx.SendStatus(r.status)
	}
	return ctx.Status(r.status).JSON(r.payload)
}
