package nvtdberr

import (
	"api/libs/nvt/nvterr"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net"

	"github.com/go-sql-driver/mysql"
)

// identifyDriver maps database/sql and driver failures that carry no server
// error code. It reports false for errors it does not recognize.
func identifyDriver(err error) (*nvterr.Error, bool) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// Only the caller knows the subject; use nvterr.NewNotFound there for a better message.
		return nvterr.ErrNotFound, true
	case errors.Is(err, driver.ErrBadConn),
		errors.Is(err, mysql.ErrInvalidConn),
		errors.Is(err, sql.ErrConnDone):
		return nvterr.ErrDatabaseUnavailable, true
	case errors.Is(err, sql.ErrTxDone):
		// A programming error (use of a finished transaction): technical, not client-facing.
		return nvterr.ErrInternal, true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return nvterr.ErrDatabaseTimeout, true
		}
		return nvterr.ErrDatabaseUnavailable, true
	}
	return nil, false
}
