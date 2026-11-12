package nvtdberr

import (
	"api/libs/nvt/nvterr"

	"github.com/go-sql-driver/mysql"
)

// mysqlErrors maps MySQL server error numbers to client-safe nvterr errors.
// TODO: populate once the MySQL error data is mapped (it may share entries with mariadbErrors).
var mysqlErrors = map[uint16]*nvterr.Error{}

// identifyMySQL translates a MySQL error number into an nvterr error.
// Unmapped numbers return nvterr.ErrInternal.
func identifyMySQL(me *mysql.MySQLError) (*nvterr.Error, Info) {
	return identifyMySQLFamily(mysqlErrors, me)
}

// identifyMySQLFamily is the lookup shared by MySQL and MariaDB, which speak
// the same protocol and report errors as *mysql.MySQLError.
func identifyMySQLFamily(table map[uint16]*nvterr.Error, me *mysql.MySQLError) (*nvterr.Error, Info) {
	e, ok := table[me.Number]
	if !ok {
		e = nvterr.ErrInternal
	}
	return e, subjectOf(me)
}
