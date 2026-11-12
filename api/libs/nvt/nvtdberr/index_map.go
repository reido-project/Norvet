package nvtdberr

import "strings"

type tableIndex struct{ table, index string }

type indexRow struct{ table, index, column string }

// indexMap is immutable once built, so readers need no locking.
type indexMap struct {
	byTableIndex map[tableIndex]string // (table, index) -> "col" or "col1, col2"
	byIndex      map[string]string     // index -> columns; only names unambiguous across tables
}

// buildIndexMap expects rows ordered by column position within each index.
func buildIndexMap(rows []indexRow) *indexMap {
	cols := make(map[tableIndex][]string)
	for _, r := range rows {
		k := tableIndex{strings.ToLower(r.table), strings.ToLower(r.index)}
		cols[k] = append(cols[k], r.column)
	}

	m := &indexMap{
		byTableIndex: make(map[tableIndex]string, len(cols)),
		byIndex:      make(map[string]string, len(cols)),
	}
	ambiguous := make(map[string]struct{})
	for k, c := range cols {
		joined := strings.Join(c, ", ")
		m.byTableIndex[k] = joined
		if prev, seen := m.byIndex[k.index]; seen && prev != joined {
			ambiguous[k.index] = struct{}{}
		}
		m.byIndex[k.index] = joined
	}
	for name := range ambiguous { // e.g. "PRIMARY" when tables have different primary keys
		delete(m.byIndex, name)
	}
	return m
}

// lookup finds the columns of an index. Index names are case-insensitive in
// MySQL.
func (m *indexMap) lookup(table, index string) (string, bool) {
	idx := strings.ToLower(index)
	if table != "" { // MySQL 8+ messages name the table: 'users.username'
		v, ok := m.byTableIndex[tableIndex{strings.ToLower(table), idx}]
		return v, ok
	}
	v, ok := m.byIndex[idx] // older servers: 'username'
	return v, ok
}
