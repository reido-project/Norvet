package nvtsql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"os"
	"testing"
	"time"

	"api/libs/nvt"
)

func TestMain(m *testing.M) {
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("APP_TIMEZONE", "UTC")
	os.Setenv("DB_TIMEZONE", "UTC")
	nvt.Init(nvt.Config{})
	os.Exit(m.Run())
}

// fakeDriver is a database/sql driver whose connections always succeed.
type fakeDriver struct{}

func (fakeDriver) Open(string) (driver.Conn, error) { return fakeConn{}, nil }

type fakeConn struct{}

func (fakeConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("not supported") }
func (fakeConn) Close() error                        { return nil }
func (fakeConn) Begin() (driver.Tx, error)           { return nil, errors.New("not supported") }

func init() { sql.Register("nvtsql-fake", fakeDriver{}) }

func TestLoadConfigDefaultsPerDialect(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Dialect != nvt.Postgres || cfg.Port != "5432" || cfg.User != "postgres" {
		t.Errorf("postgres defaults wrong: %+v", cfg)
	}
	if cfg.Pool.MaxOpenConns != 25 || cfg.Pool.ConnMaxLifetime != 30*time.Minute {
		t.Errorf("pool defaults wrong: %+v", cfg.Pool)
	}
	if cfg.Postgres.SSLMode != "disable" {
		t.Errorf("SSLMode = %q", cfg.Postgres.SSLMode)
	}
}

func TestEnvInt(t *testing.T) {
	t.Setenv("NVTSQL_TEST", "abc")
	if got := envInt("NVTSQL_TEST", 7); got != 7 {
		t.Errorf("invalid value: got %d, want 7", got)
	}
	t.Setenv("NVTSQL_TEST", "-3")
	if got := envInt("NVTSQL_TEST", 7); got != 7 {
		t.Errorf("non-positive value: got %d, want 7", got)
	}
	t.Setenv("NVTSQL_TEST", "12")
	if got := envInt("NVTSQL_TEST", 7); got != 12 {
		t.Errorf("valid value: got %d, want 12", got)
	}
}

func TestOpenErrors(t *testing.T) {
	if _, err := Open(Config{}); !errors.Is(err, ErrUnsupportedDialect) {
		t.Errorf("empty dialect: got %v, want ErrUnsupportedDialect", err)
	}
	if _, err := Open(Config{Dialect: "oracle"}); !errors.Is(err, ErrUnsupportedDialect) {
		t.Errorf("unknown dialect: got %v, want ErrUnsupportedDialect", err)
	}
	if _, err := Open(Config{Dialect: nvt.MariaDB}); !errors.Is(err, ErrDriverNotLinked) {
		t.Errorf("unregistered dialect: got %v, want ErrDriverNotLinked", err)
	}
}

func TestOpenUsesRegisteredOpenerAndAppliesPool(t *testing.T) {
	Register(nvt.MySQL, func(Config) (*sql.DB, error) { return sql.Open("nvtsql-fake", "") })

	db, err := Open(Config{Dialect: nvt.MySQL, Pool: Pool{MaxOpenConns: 3}})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if got := db.Stats().MaxOpenConnections; got != 3 {
		t.Errorf("MaxOpenConnections = %d, want 3", got)
	}
}

func TestRegisterPanicsOnDuplicate(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic on duplicate Register")
		}
	}()
	Register(nvt.MySQL, func(Config) (*sql.DB, error) { return nil, nil })
}
