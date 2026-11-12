package nvtdberr

import (
	"api/libs/nvt"
	"api/libs/nvt/nvterr"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
)

func myErr(n uint16, msg string) error { return &mysql.MySQLError{Number: n, Message: msg} }

func TestIdentify_Nil(t *testing.T) {
	isolate(t)
	if e, info := Identify(nil); e != nil || info != (Info{}) {
		t.Fatalf("got %v, %+v", e, info)
	}
}

func TestIdentify_PassesThroughNvterrAndContextErrors(t *testing.T) {
	isolate(t)
	tests := []struct {
		err  error
		want string
	}{
		{nvterr.ErrNotFound, "not_found"},
		{fmt.Errorf("repo: %w", nvterr.ErrConflict), "conflict"},
		{context.Canceled, "client_closed_request"},
		{fmt.Errorf("query: %w", context.DeadlineExceeded), "gateway_timeout"},
	}
	for _, tt := range tests {
		e, info := identify(nvt.MariaDB, tt.err)
		if e.Code() != tt.want || info != (Info{}) {
			t.Errorf("identify(%v) = %s, %+v; want %s", tt.err, e.Code(), info, tt.want)
		}
	}
}

func TestIdentify_MariaDB_DuplicateEntry(t *testing.T) {
	logs := isolate(t)
	indexes.Store(buildIndexMap([]indexRow{{"users", "username", "username"}, {"users", "uq_users_email", "email"}}))

	raw := myErr(1062, "Duplicate entry 'bob' for key 'users.username'")
	e, info := identify(nvt.MariaDB, raw)

	if e.Code() != "duplicate_value" || e.Status() != 422 {
		t.Fatalf("got %s/%d", e.Code(), e.Status())
	}
	if want := (Info{Table: "users", Column: "username", Constraint: "username"}); info != want {
		t.Fatalf("info = %+v, want %+v", info, want)
	}
	if !errors.Is(e, nvterr.ErrDuplicateValue) || !errors.Is(e, nvterr.ErrValidationFailed) {
		t.Error("must match ErrDuplicateValue and its ancestors")
	}
	if !errors.Is(e, raw) {
		t.Error("original driver error must stay reachable as the cause")
	}
	if strings.Contains(e.Message(), "bob") {
		t.Errorf("row value leaked into the client message: %q", e.Message())
	}
	if logs.Len() != 0 {
		t.Errorf("unexpected log output: %s", logs)
	}
}

func TestIdentify_DuplicateEntryVariants(t *testing.T) {
	isolate(t)
	indexes.Store(buildIndexMap([]indexRow{{"users", "uq_email", "email"}}))

	for _, tt := range []struct {
		name string
		err  error
	}{
		{"1062", myErr(1062, "Duplicate entry 'a@b.c' for key 'users.uq_email'")},
		{"1586 names the key", myErr(1586, "Duplicate entry 'a@b.c' for key 'users.uq_email'")},
		{"1859 has no value", myErr(1859, "Duplicate entry for key 'users.uq_email'")},
	} {
		e, info := identify(nvt.MariaDB, tt.err)
		if e.Code() != "duplicate_value" || info.Column != "email" || info.Table != "users" {
			t.Errorf("%s: got %s %+v", tt.name, e.Code(), info)
		}
	}
}

func TestIdentify_MySQLUsesItsOwnTable(t *testing.T) {
	isolate(t)
	e, _ := identify(nvt.MySQL, myErr(1062, "Duplicate entry 'x' for key 'k'"))
	if e.Code() != "internal" { // mysqlErrors is still empty (TODO in mysql.go)
		t.Fatalf("got %s", e.Code())
	}
}

func TestIdentify_PostgresIsRecognizedByType(t *testing.T) {
	isolate(t)
	pe := &pgconn.PgError{Code: "23505", TableName: "users", ConstraintName: "users_username_key",
		Detail: "Key (username)=(bob) already exists."}

	// The configured dialect does not matter for a pgconn error.
	for _, d := range []nvt.Dialect{nvt.Postgres, nvt.MariaDB, ""} {
		e, info := identify(d, pe)
		want := Info{Table: "users", Column: "username", Constraint: "users_username_key"}
		if info != want {
			t.Errorf("dialect %q: info = %+v, want %+v", d, info, want)
		}
		if e.Code() != "internal" { // postgresErrors is still empty (TODO in postgres.go)
			t.Errorf("dialect %q: code = %s", d, e.Code())
		}
	}
}

func TestIdentifyPostgres_CompositeKeyAndMissingDetail(t *testing.T) {
	_, info := identifyPostgres(&pgconn.PgError{Code: "23505", Detail: "Key (org_id, email)=(1, a) already exists."})
	if info.Column != "org_id, email" {
		t.Errorf("composite: %q", info.Column)
	}
	_, info = identifyPostgres(&pgconn.PgError{Code: "23505", Detail: "garbage"})
	if info.Column != "" {
		t.Errorf("no parens must give no column, got %q", info.Column)
	}
	_, info = identifyPostgres(&pgconn.PgError{Code: "23502", ColumnName: "name", TableName: "t"})
	if info.Column != "name" || info.Table != "t" {
		t.Errorf("server-provided column must win: %+v", info)
	}
}

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "i/o timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

var _ net.Error = timeoutErr{}

func TestIdentify_DriverErrors(t *testing.T) {
	isolate(t)
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"no rows", sql.ErrNoRows, "not_found"},
		{"wrapped no rows", fmt.Errorf("find: %w", sql.ErrNoRows), "not_found"},
		{"bad conn", driver.ErrBadConn, "database_unavailable"},
		{"invalid conn", mysql.ErrInvalidConn, "database_unavailable"},
		{"conn done", sql.ErrConnDone, "database_unavailable"},
		{"tx done", sql.ErrTxDone, "internal"},
		{"net timeout", &net.OpError{Op: "read", Err: timeoutErr{}}, "database_timeout"},
		{"net refused", &net.OpError{Op: "dial", Err: errors.New("connection refused")}, "database_unavailable"},
	}
	for _, tt := range tests {
		// Driver errors are dialect independent, including when the dialect is unknown.
		for _, d := range []nvt.Dialect{nvt.MariaDB, nvt.Postgres, ""} {
			e, info := identify(d, tt.err)
			if e.Code() != tt.want || info != (Info{}) {
				t.Errorf("%s (dialect %q): got %s, want %s", tt.name, d, e.Code(), tt.want)
			}
			if !errors.Is(e, tt.err) {
				t.Errorf("%s: cause lost", tt.name)
			}
		}
	}
}

func TestIdentify_UnknownDialectWarnsOnceAndStillMapsDriverErrors(t *testing.T) {
	logs := isolate(t)
	for i := 0; i < 5; i++ {
		e, _ := identify("", myErr(1062, "Duplicate entry 'x' for key 'k'"))
		if e.Code() != "internal" {
			t.Fatalf("got %s", e.Code())
		}
	}
	if n := strings.Count(logs.String(), "active database dialect"); n != 1 {
		t.Fatalf("warned %d times, want exactly once:\n%s", n, logs)
	}
	if e, _ := identify("", sql.ErrNoRows); e.Code() != "not_found" {
		t.Errorf("sql.ErrNoRows must not depend on the dialect, got %s", e.Code())
	}
}

func TestIdentify_NotADatabaseErrorWarnsOncePerType(t *testing.T) {
	logs := isolate(t)
	for i := 0; i < 3; i++ {
		e, _ := identify(nvt.MariaDB, errors.New("boom"))
		if e.Code() != "internal" {
			t.Fatalf("got %s", e.Code())
		}
	}
	if n := strings.Count(logs.String(), "not a database error"); n != 1 {
		t.Fatalf("warned %d times, want 1:\n%s", n, logs)
	}
	identify(nvt.MariaDB, fmt.Errorf("wrapped: %w", errors.New("x"))) // *fmt.wrapError, a different type
	if n := strings.Count(logs.String(), "not a database error"); n != 2 {
		t.Fatalf("a different error type must warn again, got %d warnings", n)
	}
}

func TestIdentify_UnmappedServerCodeIsInternal(t *testing.T) {
	isolate(t)
	e, _ := identify(nvt.MariaDB, myErr(1045, "Access denied for user 'root'@'localhost'"))
	if e.Code() != "internal" || e.Message() != "internal server error" {
		t.Fatalf("got %s %q", e.Code(), e.Message())
	}
	if !strings.Contains(e.Error(), "Access denied") {
		t.Error("the cause must be available for logs")
	}
}

func TestMariaDBTableInvariants(t *testing.T) {
	if len(mariadbErrors) < 600 {
		t.Fatalf("table shrank to %d entries", len(mariadbErrors))
	}
	for n, e := range mariadbErrors {
		switch {
		case e == nil:
			t.Errorf("%d: nil entry", n)
		case e.Status() == 500:
			t.Errorf("%d: %s is a 500; the table must list only non-500 errors (unlisted numbers already fall back to ErrInternal)", n, e.Code())
		}
	}
	// Spot-check the semantics the rest of the package relies on.
	for n, want := range map[uint16]*nvterr.Error{
		1048: nvterr.ErrRequired,
		1062: nvterr.ErrDuplicateValue,
		1146: nvterr.ErrNotFound,
		1451: nvterr.ErrConflict,
	} {
		if mariadbErrors[n] != want {
			t.Errorf("%d: got %v, want %v", n, mariadbErrors[n], want)
		}
	}
}
