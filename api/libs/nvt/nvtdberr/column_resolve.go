package nvtdberr

import "strings"

// resolveColumn turns a duplicate-entry (table, index) pair into a column name:
// the loaded metadata first, then inference from the index name.
func resolveColumn(table, constraint string) string {
	if constraint == "" {
		return ""
	}
	m := indexes.Load()
	if m != nil {
		if col, ok := m.lookup(table, constraint); ok {
			return col
		}
	}

	col := inferColumn(table, constraint)
	if col != "" {
		key := table + "." + constraint
		if m == nil {
			warnOnce("infer-unloaded:"+key, "nvtdberr: unique-index metadata is not loaded (nvtdberr.LoadIndexes was not called or failed); "+
				"inferred column %q from index %q. Call nvtdberr.LoadIndexes at startup for reliable results.", col, key)
		} else {
			warnOnce("infer-missing:"+key, "nvtdberr: index %q is not in the loaded metadata; inferred column %q from its name. "+
				"If the index was added by a recent migration, call nvtdberr.LoadIndexes again.", key, col)
		}
	}
	return col
}

// inferColumn guesses the column from the index name. It handles MySQL's
// default naming (index named after its column) and conventions that embed the
// table name ("{type}_{table}_{column}", "{table}_{column}"). It is anchored on
// the table name instead of splitting on the last "_" because column names
// contain underscores themselves ("phone_number").
func inferColumn(table, constraint string) string {
	if strings.EqualFold(constraint, "PRIMARY") {
		return ""
	}
	if table == "" {
		return constraint
	}

	if n := len(table); len(constraint) > n+1 &&
		strings.EqualFold(constraint[:n], table) && constraint[n] == '_' {
		return constraint[n+1:]
	}

	// Infix search only when lower-casing preserves byte offsets.
	if lc, lt := strings.ToLower(constraint), strings.ToLower(table); len(lc) == len(constraint) && len(lt) == len(table) {
		if i := strings.Index(lc, "_"+lt+"_"); i >= 0 && len(constraint) > i+len(table)+2 {
			return constraint[i+len(table)+2:]
		}
	}
	return constraint
}
