package nvtlog

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// recordingDB is a minimal database/sql driver that records every statement.
type recordingDB struct {
	mu       sync.Mutex
	execs    []execRecord
	failOn   string // fail statements containing this text
	failWith error
}

type execRecord struct {
	query string
	args  []driver.Value
}

func (d *recordingDB) Open(string) (driver.Conn, error) { return &recConn{d}, nil }

type recConn struct{ d *recordingDB }

func (c *recConn) Prepare(q string) (driver.Stmt, error) {
	if c.d.failOn != "" && strings.Contains(q, c.d.failOn) {
		return nil, c.d.failWith
	}
	return &recStmt{c.d, q}, nil
}
func (c *recConn) Close() error              { return nil }
func (c *recConn) Begin() (driver.Tx, error) { return nil, errors.New("not supported") }

type recStmt struct {
	d *recordingDB
	q string
}

func (s *recStmt) Close() error  { return nil }
func (s *recStmt) NumInput() int { return -1 }
func (s *recStmt) Exec(args []driver.Value) (driver.Result, error) {
	s.d.mu.Lock()
	defer s.d.mu.Unlock()
	s.d.execs = append(s.d.execs, execRecord{s.q, append([]driver.Value(nil), args...)})
	return driver.RowsAffected(1), nil
}
func (s *recStmt) Query([]driver.Value) (driver.Rows, error) { return nil, errors.New("not supported") }

var (
	driverSeq   int
	driverSeqMu sync.Mutex
)

func openRecording(t *testing.T, d *recordingDB) *sql.DB {
	t.Helper()
	driverSeqMu.Lock()
	driverSeq++
	name := fmt.Sprintf("nvtlog-recording-%d", driverSeq)
	driverSeqMu.Unlock()

	sql.Register(name, d)
	db, err := sql.Open(name, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestNewStore_PreparesSchemaInOrder(t *testing.T) {
	d := &recordingDB{}
	if _, err := newStore(openRecording(t, d)); err != nil {
		t.Fatal(err)
	}
	if len(d.execs) != 3 {
		t.Fatalf("executed %d statements, want 3 (pragma, table, index): %+v", len(d.execs), d.execs)
	}
	for i, want := range []string{"PRAGMA journal_mode=WAL", "CREATE TABLE IF NOT EXISTS logs", "CREATE INDEX IF NOT EXISTS idx_logs_timestamp"} {
		if !strings.Contains(d.execs[i].query, want) {
			t.Errorf("statement %d = %q, want it to contain %q", i, d.execs[i].query, want)
		}
	}
	for _, col := range []string{"timestamp TEXT NOT NULL", "level TEXT NOT NULL", "endpoint TEXT NOT NULL DEFAULT ''", "message TEXT NOT NULL"} {
		if !strings.Contains(d.execs[1].query, col) {
			t.Errorf("table definition lost column %q", col)
		}
	}
}

func TestNewStore_FailuresAreReturned(t *testing.T) {
	for _, stage := range []string{"PRAGMA", "CREATE TABLE", "CREATE INDEX", "INSERT INTO logs"} {
		d := &recordingDB{failOn: stage, failWith: errors.New("disk full")}
		if _, err := newStore(openRecording(t, d)); err == nil || !strings.Contains(err.Error(), "disk full") {
			t.Errorf("failure at %q: err = %v", stage, err)
		}
	}
}

func TestStoreAdd_ColumnsAndTimestamp(t *testing.T) {
	d := &recordingDB{}
	st, err := newStore(openRecording(t, d))
	if err != nil {
		t.Fatal(err)
	}

	loc := time.FixedZone("UTC+7", 7*3600)
	ts := time.Date(2026, 11, 10, 8, 30, 5, 100_000_000, loc)
	if err := st.add(ts, levelWarn, "GET /users/:id", "Status 404 in GET /users/:id endpoint: x"); err != nil {
		t.Fatal(err)
	}

	got := d.execs[len(d.execs)-1]
	if !strings.HasPrefix(got.query, "INSERT INTO logs") {
		t.Fatalf("query = %q", got.query)
	}
	want := []driver.Value{"2026-11-10T08:30:05.100000000+07:00", "Warn", "GET /users/:id", "Status 404 in GET /users/:id endpoint: x"}
	for i := range want {
		if got.args[i] != want[i] {
			t.Errorf("arg %d = %#v, want %#v", i, got.args[i], want[i])
		}
	}
}

func TestTimestampLayoutSortsChronologically(t *testing.T) {
	base := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	times := []time.Time{base, base.Add(time.Nanosecond), base.Add(100 * time.Millisecond), base.Add(123 * time.Millisecond), base.Add(time.Second)}

	for i := 1; i < len(times); i++ {
		a, b := times[i-1].Format(timestampLayout), times[i].Format(timestampLayout)
		if len(a) != len(b) || a >= b {
			t.Errorf("%q must sort before %q as text", a, b)
		}
	}

	// Why the width is fixed: the previous time.RFC3339Nano trims trailing zeros,
	// so 100ms renders as ".1Z" and sorts AFTER 123ms (".123Z").
	if x, y := times[2].Format(time.RFC3339Nano), times[3].Format(time.RFC3339Nano); x <= y {
		t.Errorf("expected RFC3339Nano to mis-sort %q and %q; the premise of timestampLayout no longer holds", x, y)
	}
}

func TestStoreAdd_ConcurrentWritersAreSafe(t *testing.T) {
	d := &recordingDB{}
	st, err := newStore(openRecording(t, d))
	if err != nil {
		t.Fatal(err)
	}
	before := len(d.execs)

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := st.add(time.Now(), levelInfo, "GET /x", fmt.Sprint("m", i)); err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
	if got := len(d.execs) - before; got != 50 {
		t.Fatalf("recorded %d inserts, want 50", got)
	}
}
