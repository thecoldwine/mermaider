package internal

import (
	"database/sql"
	"errors"
)

const tablesQuery = `
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
		t.table_schema = 'public' and t.table_type  = 'BASE TABLE'
	order by t.table_name, c.ordinal_position
`

type PostgresCrawler struct {
}

func (p *PostgresCrawler) Crawl(db *sql.DB) (*DatabaseSchema, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	rows, err := db.Query(tablesQuery)
	if err != nil {
		return nil, err
	}

	// we will have a capacity of 20 tables by default
	// and each table will have 20 columns capacity
	tables := make([]Table, 0, 20)

	t := ""
	var table Table

	for rows.Next() {
		var (
			tableName  string
			columnName string
			dataType   string
			nullable   bool
			relType    sql.NullString
		)

		err = rows.Scan(&tableName, &columnName, &dataType, &nullable, &relType)
		if err != nil {
			continue
		}

		if t != tableName {
			// create new Table entity and add current column
			table = Table{
				Columns: make([]Column, 0, 20),
			}

			tables = append(tables, table)
		}

		table.Columns = append(table.Columns, Column{
			Name:       columnName,
			Datatype:   dataType,
			Nullable:   nullable,
			PrimaryKey: relType.Valid && relType.String == "p",
			FK:         relType.Valid && relType.String == "f",
		})
	}

	return &DatabaseSchema{
		Tables:    tables,
		Relations: make([]Relation, 0),
	}, nil
}
