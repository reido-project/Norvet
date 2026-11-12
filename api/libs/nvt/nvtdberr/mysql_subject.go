package nvtdberr

import (
	"strings"

	"github.com/go-sql-driver/mysql"
)

// subjectFunc extracts the database objects an error message names.
type subjectFunc func(msg string) Info

// Server error numbers whose message names an index instead of a column.
const (
	errDupEntry            = 1062 // ER_DUP_ENTRY
	errDupEntryWithKeyName = 1586 // ER_DUP_ENTRY_WITH_KEY_NAME
	errDupUnknownInIndex   = 1859 // ER_DUP_UNKNOWN_IN_INDEX
)

// quoted returns the text between the first pair of quote characters q in s
// and the remainder of s after the closing quote.
func quoted(s string, q byte) (val, rest string, ok bool) {
	i := strings.IndexByte(s, q)
	if i < 0 {
		return "", s, false
	}
	j := strings.IndexByte(s[i+1:], q)
	if j < 0 {
		return "", s, false
	}
	return s[i+1 : i+1+j], s[i+j+2:], true
}

// splitDot splits "a.b" into ("a","b"); with no dot returns ("", "a").
func splitDot(s string) (left, right string) {
	if i := strings.LastIndexByte(s, '.'); i >= 0 {
		return s[:i], s[i+1:]
	}
	return "", s
}

// subjDupEntry parses duplicate-key messages:
//
//	1062 and 1586: Duplicate entry '<value>' for key 'users.username'
//	1859:          Duplicate entry for key 'users.username'
//
// The key name is always the LAST "for key '" part. It is searched from the
// end because the value is user-controlled and may itself contain quotes or
// that very phrase, so parsing from the start could return a fragment of the
// value as the "constraint".
func subjDupEntry(msg string) Info {
	const marker = "for key '"
	i := strings.LastIndex(msg, marker)
	if i < 0 {
		return Info{}
	}
	key := strings.TrimSuffix(msg[i+len(marker):], "'")
	tbl, idx := splitDot(key) // MySQL 8.0+: 'table.index'; older MySQL/MariaDB: 'index'
	return Info{Table: tbl, Constraint: idx}
}

func subjColumn(msg string) Info { // 1048, 1364, 1406 ...
	col, _, ok := quoted(msg, '\'')
	if !ok {
		return Info{}
	}
	return Info{Column: col}
}

func subjUnknownColumn(msg string) Info { // 1054: Unknown column 't.c' in 'field list'
	v, _, ok := quoted(msg, '\'')
	if !ok {
		return Info{}
	}
	tbl, col := splitDot(v)
	return Info{Table: tbl, Column: col}
}

func subjDatabase(msg string) Info { // 1049, 1007, 1008
	v, _, ok := quoted(msg, '\'')
	if !ok {
		return Info{}
	}
	return Info{Database: v}
}

func subjTable(msg string) Info { // 1050, 1051, 1146: Table 'db.tbl' doesn't exist
	v, _, ok := quoted(msg, '\'')
	if !ok {
		return Info{}
	}
	db, tbl := splitDot(v)
	return Info{Database: db, Table: tbl}
}

func subjForeignKey(msg string) Info { // 1451/1452
	// ... fails (`db`.`child`, CONSTRAINT `fk` FOREIGN KEY (`col`) REFERENCES ...
	var info Info
	rest := msg
	var ok bool
	if info.Database, rest, ok = quoted(rest, '`'); !ok {
		return Info{}
	}
	if info.Table, rest, ok = quoted(rest, '`'); !ok {
		return info
	}
	if info.Constraint, rest, ok = quoted(rest, '`'); !ok {
		return info
	}
	info.Column, _, _ = quoted(rest, '`')
	return info
}

// mysqlSubjects is sparse: only the codes whose message names a subject.
// Shared by MySQL and MariaDB.
var mysqlSubjects = map[uint16]subjectFunc{
	1007:                   subjDatabase,
	1008:                   subjDatabase,
	1048:                   subjColumn,
	1049:                   subjDatabase,
	1050:                   subjTable,
	1051:                   subjTable,
	1054:                   subjUnknownColumn,
	errDupEntry:            subjDupEntry,
	errDupEntryWithKeyName: subjDupEntry,
	errDupUnknownInIndex:   subjDupEntry,
	1146:                   subjTable,
	1364:                   subjColumn,
	1406:                   subjColumn,
	1451:                   subjForeignKey,
	1452:                   subjForeignKey,
}

// subjectOf extracts the Info for a MySQL/MariaDB server error. Duplicate-key
// errors name an index, not a column, so the column is resolved from the
// loaded index metadata (see LoadIndexes).
func subjectOf(me *mysql.MySQLError) Info {
	f, ok := mysqlSubjects[me.Number]
	if !ok {
		return Info{}
	}
	info := f(me.Message)
	if info.Column == "" && isDupEntry(me.Number) {
		info.Column = resolveColumn(info.Table, info.Constraint)
	}
	return info
}

func isDupEntry(number uint16) bool {
	return number == errDupEntry || number == errDupEntryWithKeyName || number == errDupUnknownInIndex
}
