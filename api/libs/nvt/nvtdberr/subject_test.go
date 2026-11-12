package nvtdberr

import (
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestSubjectOf(t *testing.T) {
	isolate(t)
	tests := []struct {
		name string
		num  uint16
		msg  string
		want Info
	}{
		{"1048 null column", 1048, "Column 'name' cannot be null", Info{Column: "name"}},
		{"1364 no default", 1364, "Field 'name' doesn't have a default value", Info{Column: "name"}},
		{"1406 too long", 1406, "Data too long for column 'bio' at row 1", Info{Column: "bio"}},
		{"1049 unknown db", 1049, "Unknown database 'shop'", Info{Database: "shop"}},
		{"1146 table", 1146, "Table 'shop.users' doesn't exist", Info{Database: "shop", Table: "users"}},
		{"1146 no db", 1146, "Table 'users' doesn't exist", Info{Table: "users"}},
		{"1054 unknown column", 1054, "Unknown column 'u.age' in 'field list'", Info{Table: "u", Column: "age"}},
		{"1054 no table", 1054, "Unknown column 'age' in 'where clause'", Info{Column: "age"}},
		{"1451 fk parent", 1451, "Cannot delete or update a parent row: a foreign key constraint fails (`shop`.`orders`, CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`))",
			Info{Database: "shop", Table: "orders", Constraint: "fk_orders_user", Column: "user_id"}},
		{"1452 fk child", 1452, "Cannot add or update a child row: a foreign key constraint fails (`shop`.`orders`, CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`))",
			Info{Database: "shop", Table: "orders", Constraint: "fk_orders_user", Column: "user_id"}},
		{"1451 truncated message", 1451, "fails (`shop`.`orders`", Info{Database: "shop", Table: "orders"}},
		{"unlisted code", 1045, "Access denied for user 'x'", Info{}},
		{"malformed", 1048, "no quotes here", Info{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subjectOf(&mysql.MySQLError{Number: tt.num, Message: tt.msg})
			if got != tt.want {
				t.Fatalf("got  %+v\nwant %+v", got, tt.want)
			}
		})
	}
}

func TestSubjDupEntry(t *testing.T) {
	tests := []struct {
		name, msg string
		want      Info
	}{
		{"mysql 8 table.index", "Duplicate entry 'bob' for key 'users.username'", Info{Table: "users", Constraint: "username"}},
		{"mariadb index only", "Duplicate entry 'bob' for key 'username'", Info{Constraint: "username"}},
		{"primary", "Duplicate entry '7' for key 'PRIMARY'", Info{Constraint: "PRIMARY"}},
		{"no value (1859)", "Duplicate entry for key 'users.uq_email'", Info{Table: "users", Constraint: "uq_email"}},
		// The value is user-controlled: it must never be mistaken for the key.
		{"value with quotes", "Duplicate entry 'a' for key 'x' for key 'users.username'", Info{Table: "users", Constraint: "username"}},
		{"value is the marker", "Duplicate entry 'for key 'evil'' for key 'users.username'", Info{Table: "users", Constraint: "username"}},
		{"no marker", "Duplicate entry", Info{}},
	}
	for _, tt := range tests {
		if got := subjDupEntry(tt.msg); got != tt.want {
			t.Errorf("%s: got %+v, want %+v", tt.name, got, tt.want)
		}
	}
}

func TestSubjectOf_ResolvesColumnFromLoadedIndexes(t *testing.T) {
	isolate(t)
	indexes.Store(buildIndexMap([]indexRow{
		{"users", "uq_name_org", "name"}, {"users", "uq_name_org", "org_id"},
	}))
	got := subjectOf(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'a-1' for key 'users.uq_name_org'"})
	if got.Column != "name, org_id" {
		t.Fatalf("composite key: got %q", got.Column)
	}
}

func TestInfoSubject(t *testing.T) {
	for _, tt := range []struct {
		in   Info
		want string
	}{
		{Info{Database: "d", Table: "t", Column: "c"}, "c"},
		{Info{Database: "d", Table: "t"}, "t"},
		{Info{Database: "d"}, "d"},
		{Info{}, ""},
	} {
		if got := tt.in.Subject(); got != tt.want {
			t.Errorf("%+v.Subject() = %q, want %q", tt.in, got, tt.want)
		}
	}
}
