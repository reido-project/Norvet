package nvtdberr

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

// isolate resets the package's global state and captures log output for the
// duration of the test. It returns the captured log.
func isolate(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	indexes.Store(nil)
	warned.Clear()
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
		indexes.Store(nil)
		warned.Clear()
	})
	return &buf
}

// fakeDB is a minimal database/sql driver that answers every query with a
// canned result, so LoadIndexes can run against real *sql.Rows.
type fakeDB struct {
	mu      sync.Mutex
	rows    [][]driver.Value
	err     error
	queries []string
}

func (f *fakeDB) Open(string) (driver.Conn, error) { return &fakeConn{f}, nil }

type fakeConn struct{ f *fakeDB }

func (c *fakeConn) Prepare(string) (driver.Stmt, error) { return nil, fmt.Errorf("not supported") }
func (c *fakeConn) Close() error                        { return nil }
func (c *fakeConn) Begin() (driver.Tx, error)           { return nil, fmt.Errorf("not supported") }
func (c *fakeConn) QueryContext(_ context.Context, q string, _ []driver.NamedValue) (driver.Rows, error) {
	c.f.mu.Lock()
	defer c.f.mu.Unlock()
	c.f.queries = append(c.f.queries, q)
	if c.f.err != nil {
		return nil, c.f.err
	}
	return &fakeRows{rows: c.f.rows}, nil
}

type fakeRows struct {
	rows [][]driver.Value
	i    int
}

func (r *fakeRows) Columns() []string { return []string{"TABLE_NAME", "INDEX_NAME", "COLUMN_NAME"} }
func (r *fakeRows) Close() error      { return nil }
func (r *fakeRows) Next(dest []driver.Value) error {
	if r.i >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.i])
	r.i++
	return nil
}

var fakeSeq int
var fakeSeqMu sync.Mutex

// openFake registers a fresh fake driver and returns a *sql.DB on it.
func openFake(t *testing.T, f *fakeDB) *sql.DB {
	t.Helper()
	fakeSeqMu.Lock()
	fakeSeq++
	name := fmt.Sprintf("nvtdberr-fake-%d", fakeSeq)
	fakeSeqMu.Unlock()

	sql.Register(name, f)
	db, err := sql.Open(name, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}
