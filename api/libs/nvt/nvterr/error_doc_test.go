package nvterr

import (
	"cmp"
	"flag"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

//go:generate go test -run TestErrorsDocUpToDate -update .

var update = flag.Bool("update", false, "rewrite ERRORS.md from the registered errors")

// TestErrorsDocUpToDate keeps ERRORS.md in sync with the definitions.
// After adding or changing an error, run `go generate ./...` in this package.
func TestErrorsDocUpToDate(t *testing.T) {
	want := renderErrorsDoc(registry())

	if *update {
		if err := os.WriteFile("ERRORS.md", []byte(want), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}

	got, err := os.ReadFile("ERRORS.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatal("ERRORS.md is out of date; run `go generate ./...` in the nvterr package")
	}
}

// renderErrorsDoc renders the error tree as a Markdown table: roots ordered
// by status then code, children by code under their parent, indented by depth.
func renderErrorsDoc(all []*Error) string {
	children := make(map[*Error][]*Error)
	var roots []*Error
	for _, e := range all {
		if e.parent == nil {
			roots = append(roots, e)
		} else {
			children[e.parent] = append(children[e.parent], e)
		}
	}
	slices.SortFunc(roots, func(a, b *Error) int {
		return cmp.Or(cmp.Compare(a.status, b.status), cmp.Compare(a.code, b.code))
	})

	var sb strings.Builder
	sb.WriteString("# nvterr error reference\n\n")
	sb.WriteString(strconv.Itoa(len(all)))
	sb.WriteString(" errors. Indentation shows inheritance: `errors.Is(child, parent)` is true, never the reverse.\n\n")
	sb.WriteString("| Code | Status | Message |\n|---|---|---|\n")

	var walk func(e *Error, depth int)
	walk = func(e *Error, depth int) {
		sb.WriteString("| ")
		sb.WriteString(strings.Repeat("&nbsp;", depth*4))
		sb.WriteString("`" + e.code + "` | " + strconv.Itoa(e.status) + " | " + e.message + " |\n")

		kids := children[e]
		slices.SortFunc(kids, func(a, b *Error) int { return cmp.Compare(a.code, b.code) })
		for _, k := range kids {
			walk(k, depth+1)
		}
	}
	for _, r := range roots {
		walk(r, 0)
	}
	return sb.String()
}
