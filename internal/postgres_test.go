package internal

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresCrawl(t *testing.T) {
	ctx := context.Background()

	dbName := "sakila"
	dbUser := "postgres"
	dbPassword := "password"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithInitScripts(filepath.Join("../examples", "pg-sakila.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	cs, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("can't obtain a connection string: %s", err)
	}

	db, err := sql.Open("pgx", cs)
	if err != nil {
		t.Fatalf("error while creating a db %s", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		t.Fatalf("error while pinging the context %s", err)
	}

	pgCrawler := PostgresCrawler{db: db}
	schema, err := pgCrawler.Crawl("public")
	if err != nil {
		t.Fatalf("Cannot crawl a database schema: %s", err)
	}

	if len(schema.Tables) != 21 {
		t.Fatalf("Expected 21 table, but got %d", len(schema.Tables))
	}

	log.Printf("%v", schema.Tables)

	for _, table := range schema.Tables {
		t.Logf("Table name: %s, data: %v\n", table.Name, table.Columns)
	}

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

}
