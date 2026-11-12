package postgres

import (
	"net/url"
	"testing"
	"time"

	"api/libs/nvt/nvtsql"
)

func TestDSN(t *testing.T) {
	jakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Skip("tzdata not available")
	}

	cfg := nvtsql.Config{
		Host: "db.local", Port: "5432", User: "app", Password: "p@ss:w/rd%", Name: "shop",
		Location: jakarta,
		Postgres: nvtsql.PostgresOptions{SSLMode: "require"},
	}
	u, err := url.Parse(dsn(cfg))
	if err != nil {
		t.Fatal(err)
	}
	if pw, _ := u.User.Password(); u.User.Username() != "app" || pw != "p@ss:w/rd%" {
		t.Errorf("credentials did not round-trip: %v", u.User)
	}
	if u.Host != "db.local:5432" || u.Path != "/shop" {
		t.Errorf("host/path = %q %q", u.Host, u.Path)
	}
	if q := u.Query(); q.Get("sslmode") != "require" || q.Get("TimeZone") != "Asia/Jakarta" {
		t.Errorf("query = %v", q)
	}
}

func TestDSNOmitsLocalTimezone(t *testing.T) {
	u, err := url.Parse(dsn(nvtsql.Config{Host: "h", Port: "1", Location: time.Local}))
	if err != nil {
		t.Fatal(err)
	}
	if u.Query().Has("TimeZone") {
		t.Errorf("TimeZone must not be set for the Local location: %v", u.Query())
	}
}
