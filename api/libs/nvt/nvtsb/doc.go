// Package sqlbuilder connects a Norvet application to its database and hands
// back a query builder, configured from the environment.
//
//	nvt.Init(nvt.Config{})
//
//	db, err := sqlbuilder.ConnectGoqu()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer sqlbuilder.CloseGoqu(db)
//
//	var users []User
//	err = db.From("users").Where(goqu.C("active").IsTrue()).ScanStructs(&users)
//
// The connection itself (driver, DSN, pool, ping) is handled by package
// nvtsql, so every builder gets the same behavior. The environment variables
// are documented on nvtsql.LoadConfig.
//
// # Supported builders
//
//	goqu   ConnectGoqu, WrapGoqu, CloseGoqu
//
// To use plain database/sql without a builder, call nvtsql.OpenFromEnv.
//
// # Adding a builder
//
// Create one file per builder with its own Connect function; no existing code
// changes:
//
//	func ConnectSquirrel() (squirrel.StatementBuilderType, error) {
//		return connect(WrapSquirrel)
//	}
//
//	func WrapSquirrel(db *sql.DB, dialect nvt.Dialect) (squirrel.StatementBuilderType, error) {
//		...
//	}
//
// # Layout
//
//	builder.go   connect: the shared open-then-wrap step
//	goqu.go      goqu builder
package sqlbuilder
