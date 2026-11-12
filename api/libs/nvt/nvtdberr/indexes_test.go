package nvtdberr

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
)

func TestBuildIndexMap(t *testing.T) {
	rows := []indexRow{
		{"Users", "PRIMARY", "id"},
		{"users", "username", "username"},
		{"users", "uq_name_org", "name"}, {"users", "uq_name_org", "org_id"},
		{"orders", "PRIMARY", "order_no"},                        // PRIMARY differs between tables
		{"a", "idx_email", "email"}, {"b", "idx_email", "email"}, // same name, same columns
		{"a", "idx_x", "x"}, {"b", "idx_x", "y"}, // same name, different columns
	}
	// Map iteration order is random; the result must not depend on it.
	for i := 0; i < 50; i++ {
		m := buildIndexMap(rows)

		for _, tt := range []struct {
			table, index, want string
			ok                 bool
		}{
			{"users", "primary", "id", true}, // case-insensitive on both sides
			{"USERS", "USERNAME", "username", true},
			{"users", "uq_name_org", "name, org_id", true},
			{"orders", "PRIMARY", "order_no", true},
			{"", "username", "username", true},
			{"", "uq_name_org", "name, org_id", true},
			{"", "PRIMARY", "", false}, // ambiguous without a table
			{"", "idx_email", "email", true},
			{"", "idx_x", "", false}, // ambiguous without a table
			{"a", "idx_x", "x", true},
			{"b", "idx_x", "y", true},
			{"users", "nope", "", false},
		} {
			got, ok := m.lookup(tt.table, tt.index)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("lookup(%q, %q) = %q, %v; want %q, %v", tt.table, tt.index, got, ok, tt.want, tt.ok)
			}
		}
	}
}

func TestInferColumn(t *testing.T) {
	tests := []struct{ table, index, want string }{
		{"users", "PRIMARY", ""},
		{"users", "primary", ""},
		{"", "username", "username"},                       // MySQL default: index named after the column
		{"users", "username", "username"},                  // no table in the name: assume column
		{"users", "users_email", "email"},                  // {table}_{column}
		{"users", "Users_Email", "Email"},                  // case-insensitive table match
		{"users", "uq_users_phone_number", "phone_number"}, // underscores in column names
		{"users", "uq_users_", "uq_users_"},                // nothing after the table: not a guess
		{"users", "users_", "users_"},
	}
	for _, tt := range tests {
		if got := inferColumn(tt.table, tt.index); got != tt.want {
			t.Errorf("inferColumn(%q, %q) = %q, want %q", tt.table, tt.index, got, tt.want)
		}
	}
}

func TestResolveColumn_WarnsOncePerIndexAndSaysWhy(t *testing.T) {
	logs := isolate(t)

	for i := 0; i < 4; i++ {
		if got := resolveColumn("users", "users_email"); got != "email" {
			t.Fatalf("got %q", got)
		}
	}
	if n := strings.Count(logs.String(), "metadata is not loaded"); n != 1 {
		t.Fatalf("not-loaded warning logged %d times, want 1:\n%s", n, logs)
	}

	// Once metadata is loaded, an unknown index warns with a different reason.
	indexes.Store(buildIndexMap([]indexRow{{"users", "other", "other"}}))
	resolveColumn("users", "users_email")
	resolveColumn("users", "users_email")
	if n := strings.Count(logs.String(), "is not in the loaded metadata"); n != 1 {
		t.Fatalf("missing-index warning logged %d times, want 1:\n%s", n, logs)
	}

	// A hit is silent and exact.
	before := logs.Len()
	if got := resolveColumn("users", "other"); got != "other" || logs.Len() != before {
		t.Fatalf("hit: got %q, log grew by %d", got, logs.Len()-before)
	}
	if got := resolveColumn("users", ""); got != "" {
		t.Fatalf("empty constraint: got %q", got)
	}
}

func TestLoadIndexes(t *testing.T) {
	logs := isolate(t)
	f := &fakeDB{rows: [][]driver.Value{
		{"users", "PRIMARY", "id"},
		{"users", "uq_name_org", "name"},
		{"users", "uq_name_org", "org_id"},
		{"users", "uq_lower_email", nil}, // functional index: NULL column is skipped
	}}
	db := openFake(t, f)

	if err := loadIndexes(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	m := indexes.Load()
	if m == nil {
		t.Fatal("snapshot not stored")
	}
	if got, _ := m.lookup("users", "uq_name_org"); got != "name, org_id" {
		t.Errorf("composite: %q", got)
	}
	if _, ok := m.lookup("users", "uq_lower_email"); ok {
		t.Error("functional index must be skipped")
	}
	if len(f.queries) != 1 || !strings.Contains(f.queries[0], "information_schema.STATISTICS") {
		t.Errorf("queries: %v", f.queries)
	}
	if logs.Len() != 0 {
		t.Errorf("unexpected log: %s", logs)
	}
}

func TestLoadIndexes_FailureKeepsPreviousSnapshot(t *testing.T) {
	logs := isolate(t)
	old := buildIndexMap([]indexRow{{"t", "i", "c"}})
	indexes.Store(old)

	f := &fakeDB{err: errors.New("SELECT command denied")}
	err := loadIndexes(context.Background(), openFake(t, f))
	if err == nil || !strings.Contains(err.Error(), "SELECT command denied") {
		t.Fatalf("err = %v", err)
	}
	if indexes.Load() != old {
		t.Error("a failed reload must keep the previous snapshot")
	}
	if !strings.Contains(logs.String(), "information_schema") {
		t.Errorf("the failure must be logged with guidance, got: %s", logs)
	}
}

func TestLoadIndexes_EmptyResultWarns(t *testing.T) {
	logs := isolate(t)
	if err := loadIndexes(context.Background(), openFake(t, &fakeDB{})); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(logs.String(), "returned no unique indexes") {
		t.Errorf("got: %s", logs)
	}
}
