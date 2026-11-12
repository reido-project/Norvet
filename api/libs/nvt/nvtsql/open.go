package nvtsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"api/libs/nvt"
)

// Errors returned by Open.
var (
	// ErrUnsupportedDialect means Config.Dialect is empty or not a known
	// dialect, usually because DB_DRIVER is unset or misspelled.
	ErrUnsupportedDialect = errors.New("nvtsql: unsupported or missing database dialect")

	// ErrDriverNotLinked means the dialect is valid but its driver package was
	// not imported into the program.
	ErrDriverNotLinked = errors.New("nvtsql: driver for the dialect is not linked into the program")
)

// pingTimeout bounds the connectivity check done by Open.
const pingTimeout = 10 * time.Second

// Opener creates the *sql.DB for a dialect. It only has to open the handle;
// Open applies the pool settings and verifies connectivity afterwards.
type Opener func(cfg Config) (*sql.DB, error)

var (
	mu      sync.RWMutex
	openers = map[nvt.Dialect]Opener{}
)

// Register makes opener the connector for dialect. Driver packages call it
// from init. It panics on a nil opener or a dialect registered twice, like
// database/sql.Register.
func Register(dialect nvt.Dialect, opener Opener) {
	if opener == nil {
		panic("nvtsql: Register opener is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, dup := openers[dialect]; dup {
		panic("nvtsql: Register called twice for dialect " + string(dialect))
	}
	openers[dialect] = opener
}

// Open connects to the database described by cfg: it opens the handle with
// the registered driver, applies cfg.Pool and pings the server (10 seconds at
// most). The caller owns the returned *sql.DB and must Close it.
func Open(cfg Config) (*sql.DB, error) {
	if _, ok := nvt.ParseDialect(string(cfg.Dialect)); !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedDialect, cfg.Dialect)
	}

	mu.RLock()
	opener := openers[cfg.Dialect]
	mu.RUnlock()
	if opener == nil {
		return nil, fmt.Errorf("%w: %q (import _ \"api/libs/nvt/nvtsql/drivers\")", ErrDriverNotLinked, cfg.Dialect)
	}

	db, err := opener(cfg)
	if err != nil {
		return nil, err
	}

	cfg.Pool.apply(db)

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("nvtsql: ping %s database: %w", cfg.Dialect, err)
	}
	return db, nil
}

// OpenFromEnv is LoadConfig followed by Open: the one-call way to get a
// configured *sql.DB.
func OpenFromEnv() (*sql.DB, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return Open(cfg)
}

// apply sets every non-zero pool field on db.
func (p Pool) apply(db *sql.DB) {
	if p.MaxOpenConns > 0 {
		db.SetMaxOpenConns(p.MaxOpenConns)
	}
	if p.MaxIdleConns > 0 {
		db.SetMaxIdleConns(p.MaxIdleConns)
	}
	if p.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(p.ConnMaxLifetime)
	}
	if p.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(p.ConnMaxIdleTime)
	}
}
