package mermaider

import (
	"flag"
)

var connString string
var schemaName string

func main() {
	flag.StringVar(&connString, "connection-string", "", "Connection string to a target database, driver will be inferred automatically")
	flag.StringVar(&schemaName, "schema", "", "Schema name for the database, defaults to dbo for MSSQL, public for postgres")

	flag.Parse()

}
