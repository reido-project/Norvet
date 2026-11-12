package nvtlog

import "github.com/gofiber/fiber/v3"

// Info logs err at info level. The optional code overrides the status shown
// in the log line; without it the status comes from the error when it is an
// nvterr error (or a context cancellation/deadline), and is omitted otherwise.
//
// ctx supplies the method and route; it may be nil (for example in a
// background job), in which case both show as UNKNOWN. A nil err logs nothing.
func Info(ctx fiber.Ctx, err error, code ...int) {
	write(ctx, err, levelInfo, resolveStatus(err, code, 0))
}

// Warn logs err at warn level. See Info.
func Warn(ctx fiber.Ctx, err error, code ...int) {
	write(ctx, err, levelWarn, resolveStatus(err, code, 0))
}

// Error logs err at error level. See Info.
func Error(ctx fiber.Ctx, err error, code ...int) {
	write(ctx, err, levelError, resolveStatus(err, code, 0))
}
