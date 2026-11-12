package nvtval

import (
	"reflect"
	"testing"
)

func TestMessage(t *testing.T) {
	tests := []struct {
		name              string
		tag, field, param string
		kind              reflect.Kind
		want              string
	}{
		// Messages that existed before keep their wording.
		{"required", "required", "Email", "", reflect.String, "Email is required and must not be empty"},
		{"email", "email", "Email", "", reflect.String, "Email must be a valid email address"},
		{"uuid4", "uuid4", "ID", "", reflect.String, "ID must be a valid UUID"},
		{"min string", "min", "Name", "3", reflect.String, "Name must be at least 3 characters"},
		{"max string", "max", "Name", "20", reflect.String, "Name must not exceed 20 characters"},
		{"len string", "len", "Code", "6", reflect.String, "Code must be exactly 6 characters long"},
		{"oneof", "oneof", "Role", "admin user", reflect.String, "Role must be one of the following values: [admin user]"},
		{"eqfield", "eqfield", "Confirm", "Password", reflect.String, "Confirm didn't match"},
		{"excludes space", "excludes", "Username", " ", reflect.String, "Username cannot contain space in between characters"},
		{"excludes char", "excludes", "Username", "@", reflect.String, "Username cannot contain character '@'"},

		// Size rules read correctly for numbers and collections.
		{"min number", "min", "Age", "18", reflect.Int, "Age must be at least 18"},
		{"max number", "max", "Age", "99", reflect.Float64, "Age must not exceed 99"},
		{"len number", "len", "Pin", "4", reflect.Int, "Pin must be exactly 4"},
		{"min slice", "min", "Tags", "1", reflect.Slice, "Tags must contain at least 1 items"},
		{"max map", "max", "Meta", "5", reflect.Map, "Meta must not contain more than 5 items"},
		{"gt number", "gt", "Price", "0", reflect.Float64, "Price must be greater than 0"},
		{"gte number", "gte", "Price", "1", reflect.Int, "Price must be at least 1"},
		{"lt string", "lt", "Bio", "100", reflect.String, "Bio must be shorter than 100 characters"},
		{"lte number", "lte", "Qty", "10", reflect.Int, "Qty must not exceed 10"},

		// Fallbacks always name the field.
		{"unknown tag", "isbn", "Book", "", reflect.String, "Book failed on validation rule 'isbn'"},
		{"unknown tag with param", "weird", "Book", "x", reflect.String, "Book failed on validation rule 'weird' with value 'x'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := message(tt.tag, tt.field, tt.param, tt.kind); got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}
