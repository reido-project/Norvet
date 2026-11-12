package nvterr

import "net/http"

// StatusClientClosedRequest is the de-facto (nginx) status for requests the
// client abandoned before the server responded. It is not in net/http.
const StatusClientClosedRequest = 499

// FromStatus maps a bare HTTP status code to the matching top-level sentinel.
// It is meant for errors a framework produced itself and that carry only a
// status, such as Fiber's *fiber.Error for an unmatched route (404), a wrong
// method (405) or an oversized body (413):
//
//	var fe *fiber.Error
//	if errors.As(err, &fe) {
//		e = nvterr.FromStatus(fe.Code)
//	}
//
// Statuses without a dedicated sentinel yield a generic error that keeps the
// status and uses its standard text as the message ("client_error" for 4xx,
// "server_error" for 5xx). Anything outside 400-599 maps to ErrInternal.
func FromStatus(status int) *Error {
	if e, ok := sentinelByStatus[status]; ok {
		return e
	}
	switch {
	case status >= 400 && status < 500:
		return New("client_error", status, statusText(status))
	case status >= 500 && status < 600:
		return New("server_error", status, statusText(status))
	}
	return ErrInternal
}

// sentinelByStatus lists, for each status that has one, the broadest sentinel.
var sentinelByStatus = map[int]*Error{
	http.StatusBadRequest:            ErrBadRequest,
	http.StatusUnauthorized:          ErrUnauthenticated,
	http.StatusPaymentRequired:       ErrPaymentRequired,
	http.StatusForbidden:             ErrForbidden,
	http.StatusNotFound:              ErrNotFound,
	http.StatusMethodNotAllowed:      ErrMethodNotAllowed,
	http.StatusNotAcceptable:         ErrNotAcceptable,
	http.StatusRequestTimeout:        ErrRequestTimeout,
	http.StatusConflict:              ErrConflict,
	http.StatusGone:                  ErrGone,
	http.StatusPreconditionFailed:    ErrPreconditionFailed,
	http.StatusRequestEntityTooLarge: ErrPayloadTooLarge,
	http.StatusUnsupportedMediaType:  ErrUnsupportedMediaType,
	http.StatusUnprocessableEntity:   ErrUnprocessable,
	http.StatusLocked:                ErrLocked,
	http.StatusPreconditionRequired:  ErrPreconditionRequired,
	http.StatusTooManyRequests:       ErrRateLimited,
	StatusClientClosedRequest:        ErrClientClosedRequest,
	http.StatusInternalServerError:   ErrInternal,
	http.StatusNotImplemented:        ErrNotImplemented,
	http.StatusBadGateway:            ErrDependencyFailed,
	http.StatusServiceUnavailable:    ErrServiceUnavailable,
	http.StatusGatewayTimeout:        ErrGatewayTimeout,
}

// statusText is http.StatusText with a fallback for unregistered codes.
func statusText(status int) string {
	if t := http.StatusText(status); t != "" {
		return t
	}
	return "http error"
}
