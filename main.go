package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/microsoft/go-mssqldb/msdsn"
	"github.com/thecoldwine/mermaider/internal"
)

var connString string
var schemaName string
var dbType string

func guessDatabase(cs string) string {
	_, err := pgx.ParseConfig(cs)
	if err == nil {
		return "postgres"
	}

	_, err = msdsn.Parse(cs)
	if err == nil {
		return "sqlserver"
	}

	return "unknown"
}

func main() {
	flag.StringVar(&connString, "connection-string", "", "Connection string to a target database, driver will be inferred automatically")
	flag.StringVar(&schemaName, "schema", "", "Schema name for the database, defaults to dbo for MSSQL, public for postgres")
	flag.StringVar(&dbType, "db-type", "postgres", "Type of the database if it cannot be guessed right from the connection string, defaults to postgres")

	flag.Parse()

	guessedDb := guessDatabase(connString)
	if guessedDb == "unknown" {
		guessedDb = dbType
	}

	var crawler internal.SchemaCrawler
	switch guessedDb {
	case "postgres":
		db, err := sql.Open("pgx", connString)
		if err != nil {
			log.Fatalln(err)
		}

		if schemaName == "" {
			schemaName = "public"
		}

		crawler = internal.NewPostgresCrawler(db)
	case "sqlserver":
		if schemaName == "" {
			schemaName = "dbo"
		}
		log.Fatalln("Not implemented yet")
	default:
		log.Fatalf("Unknown database type: %s\n", guessedDb)
	}

	err := internal.Mermaid(crawler, schemaName, os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
}
