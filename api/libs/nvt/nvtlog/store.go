package nvtlog

import (
	"api/libs/nvt/nvtsql/sqlite"
	"database/sql"
	stdlog "log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// timestampLayout is RFC 3339 with a fixed-width fraction. Unlike
// time.RFC3339Nano it does not trim trailing zeros, so the strings sort
// chronologically as text and the timestamp index can serve range queries.
const timestampLayout = "2006-01-02T15:04:05.000000000Z07:00"

const (
	createTableSQL = `
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			level TEXT NOT NULL,
			endpoint TEXT NOT NULL DEFAULT '',
			message TEXT NOT NULL
		)`
	createIndexSQL = `CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp)`
	insertSQL      = `INSERT INTO logs (timestamp, level, endpoint, message) VALUES (?, ?, ?, ?)`
)

// store writes log entries to the SQLite "logs" table.
type store struct {
	mu     sync.Mutex // SQLite allows one writer at a time; serializing avoids "database is locked"
	insert *sql.Stmt
}

// newStore prepares db for logging: WAL mode, the table, its index and the
// insert statement.
func newStore(db *sql.DB) (*store, error) {
	for _, q := range []string{"PRAGMA journal_mode=WAL;", createTableSQL, createIndexSQL} {
		if _, err := db.Exec(q); err != nil {
			return nil, err
		}
	}
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return nil, err
	}
	return &store{insert: stmt}, nil
}

func (s *store) add(ts time.Time, lvl level, endpoint, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.insert.Exec(ts.Format(timestampLayout), string(lvl), endpoint, message)
	return err
}

var (
	defaultMu    sync.Mutex
	defaultStore *store
)

// openDefaultStore opens <dir>/nvtlog.sqlite on first use and reuses it after.
// A failure is not cached: the next call tries again.
func openDefaultStore(dir string) (*store, error) {
	defaultMu.Lock()
	defer defaultMu.Unlock()

	if defaultStore != nil {
		return defaultStore, nil
	}

	dirPath := filepath.Join(".", dir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}

	db, err := sqlite.Open(filepath.Join(dirPath, "nvtlog.sqlite"))
	if err != nil {
		return nil, err
	}

	st, err := newStore(db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	defaultStore = st
	return st, nil
}

// saveToSQLite persists e. Failures are reported on the standard logger and
// never returned: logging must not break a request.
func saveToSQLite(s settings, e entry) {
	st, err := openDefaultStore(s.logDir)
	if err != nil {
		stdlog.Printf("Unable to prepare nvtlog database: %v", err)
		return
	}
	if err := st.add(s.now(), e.level, e.method+" "+e.endpoint, plainMessage(e)); err != nil {
		stdlog.Printf("Unable to insert nvtlog into sqlite: %v", err)
	}
}
