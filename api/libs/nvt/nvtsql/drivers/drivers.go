// Package drivers links every database connector Norvet supports (MySQL,
// MariaDB, PostgreSQL and SQLite) into the program. Import it for its side
// effect:
//
//	import _ "api/libs/nvt/nvtsql/drivers"
//
// To keep a binary smaller, import only the connectors you need instead:
// nvtsql/mysql, nvtsql/postgres or nvtsql/sqlite.
package drivers

import (
	_ "api/libs/nvt/nvtsql/mysql"
	_ "api/libs/nvt/nvtsql/postgres"
	_ "api/libs/nvt/nvtsql/sqlite"
)
