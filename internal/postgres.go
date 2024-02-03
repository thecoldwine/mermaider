package internal

import (
	"database/sql"
	"errors"
)

const pgTablesQuery = `
	with col_constraints as (select
		k.table_schema,
		k.table_name,
		k.column_name,
		pgc.contype
	from
		information_schema.key_column_usage k
		inner join  pg_constraint pgc on k.constraint_name = pgc.conname)
	select
		t.table_name,
		c.column_name,
		c.data_type,
		c.is_nullable,
		cc.contype
	from
		information_schema.tables t
		inner join information_schema.columns c on t.table_name = c.table_name and t.table_schema = c.table_schema
		left join col_constraints cc on t.table_schema = cc.table_schema and c.column_name = cc.column_name and c.table_name = cc.table_name
	where
		t.table_schema = $1 and t.table_type = 'BASE TABLE'
	order by t.table_name, c.ordinal_position
`

// Mermaid isn't very nice when it comes to support actual
// links between columns in ER, but we still like to know
// the column names
// This query _always_ interprets connection type as
// many-to-one
const pgRelationsQuery = `
	SELECT
		tc.table_name,
		kcu.column_name,
		ccu.table_name AS foreign_table_name,
		ccu.column_name AS foreign_column_name
	FROM information_schema.table_constraints AS tc
	JOIN information_schema.key_column_usage AS kcu
		ON tc.constraint_name = kcu.constraint_name
		AND tc.table_schema = kcu.table_schema
	JOIN information_schema.constraint_column_usage AS ccu
		ON ccu.constraint_name = tc.constraint_name
	WHERE tc.constraint_type = 'FOREIGN KEY'
		and tc.table_schema = $1
`

type PostgresCrawler struct {
	db *sql.DB
}

func NewPostgresCrawler(db *sql.DB) *PostgresCrawler {
	return &PostgresCrawler{
		db: db,
	}
}

func pgCrawlTables(db *sql.DB, schemaName string) ([]Table, error) {
	rows, err := db.Query(pgTablesQuery, schemaName)
	if err != nil {
		return nil, err
	}

	// we will have a capacity of 20 tables by default
	// and each table will have 20 columns capacity
	tables := make([]Table, 0, 20)

	var table *Table = nil

	for rows.Next() {
		var (
			tableName  string
			columnName string
			dataType   string
			nullable   string
			conType    sql.NullString
		)

		err = rows.Scan(&tableName, &columnName, &dataType, &nullable, &conType)
		if err != nil {
			continue
		}

		if table == nil {
			table = &Table{
				Name:    tableName,
				Columns: make([]Column, 0, 20),
			}
		}

		if table.Name != tableName {
			tables = append(tables, *table)

			table = &Table{
				Name:    tableName,
				Columns: make([]Column, 0, 20),
			}
		}

		table.Columns = append(table.Columns, Column{
			Name:       columnName,
			Datatype:   dataType,
			Nullable:   nullable == "YES",
			PrimaryKey: conType.Valid && conType.String == "p",
			ForeignKey: conType.Valid && conType.String == "f",
		})
	}

	if table != nil {
		tables = append(tables, *table)
	}

	return tables, nil
}

func pgCrawlRelations(db *sql.DB, schemaName string) ([]Relation, error) {
	rows, err := db.Query(pgRelationsQuery, schemaName)
	if err != nil {
		return nil, err
	}

	relations := make([]Relation, 0, 20)

	for rows.Next() {
		var (
			srcTable  string
			srcColumn string
			dstTable  string
			dstColumn string
		)

		err = rows.Scan(&srcTable, &srcColumn, &dstTable, &dstColumn)
		if err != nil {
			continue
		}

		relations = append(relations, Relation{
			SourceTable:       srcTable,
			SourceColumn:      srcColumn,
			DestinationTable:  dstTable,
			DestinationColumn: dstColumn,
			RelType:           ManyToOne,
		})
	}

	return relations, nil
}

func (p *PostgresCrawler) Crawl(schemaName string) (*DatabaseSchema, error) {
	if p.db == nil {
		return nil, errors.New("database is nil")
	}

	if err := p.db.Ping(); err != nil {
		return nil, err
	}

	tables, err := pgCrawlTables(p.db, schemaName)
	if err != nil {
		return nil, err
	}

	relations, err := pgCrawlRelations(p.db, schemaName)
	if err != nil {
		return nil, err
	}

	return &DatabaseSchema{
		Tables:    tables,
		Relations: relations,
	}, nil
}
