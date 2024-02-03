package internal

import (
	"database/sql"
	"errors"
)

const mssqlTablesQuery = `
	with keys as (
		select
			kcu.TABLE_NAME,
			kcu.COLUMN_NAME,
			IIF(tc_p.CONSTRAINT_TYPE = 'PRIMARY KEY', 'YES', NULL) pk,
			IIF(tc_f.CONSTRAINT_TYPE = 'FOREIGN KEY', 'YES', NULL) fk
		from
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			left join INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc_p on
				tc_p.TABLE_SCHEMA = kcu.TABLE_SCHEMA
				and tc_p.TABLE_NAME = kcu.TABLE_NAME
				and tc_p.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
				and tc_p.CONSTRAINT_TYPE = 'PRIMARY KEY'
			left join INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc_f on
				tc_f.TABLE_SCHEMA = kcu.TABLE_SCHEMA
				and tc_f.TABLE_NAME = tc_f.TABLE_NAME
				and tc_f.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
				and tc_f.CONSTRAINT_TYPE = 'FOREIGN KEY'
		where kcu.TABLE_SCHEMA = @schema
	)
	select
		c.TABLE_NAME,
		c.COLUMN_NAME,
		c.DATA_TYPE,
		c.IS_NULLABLE,
		k.pk,
		k.fk
	from
		INFORMATION_SCHEMA.COLUMNS c
		left join keys k on k.TABLE_NAME = c.TABLE_NAME and k.COLUMN_NAME = c.COLUMN_NAME
	where
		c.TABLE_SCHEMA = @schema
`

const mssqlRelationsQuery = `
	select
		ftc.TABLE_NAME src_table,
		kcu1.COLUMN_NAME src_col,
		ptc.TABLE_NAME dst_table,
		kcu2.COLUMN_NAME dst_col
	from
		INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
		inner join INFORMATION_SCHEMA.TABLE_CONSTRAINTS ftc
			on rc.CONSTRAINT_NAME = ftc.CONSTRAINT_NAME and rc.CONSTRAINT_SCHEMA = ftc.CONSTRAINT_SCHEMA and ftc.CONSTRAINT_TYPE = 'FOREIGN KEY'
		inner join INFORMATION_SCHEMA.TABLE_CONSTRAINTS ptc
			on rc.UNIQUE_CONSTRAINT_NAME = ptc.CONSTRAINT_NAME and rc.UNIQUE_CONSTRAINT_SCHEMA = ptc.CONSTRAINT_SCHEMA and ptc.CONSTRAINT_TYPE = 'PRIMARY KEY'
		inner join INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu1
			on kcu1.CONSTRAINT_SCHEMA = ftc.CONSTRAINT_SCHEMA and kcu1.CONSTRAINT_NAME = ftc.CONSTRAINT_NAME
		inner join INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu2
			on kcu2.CONSTRAINT_SCHEMA = ptc.CONSTRAINT_SCHEMA and kcu2.CONSTRAINT_NAME = ptc.CONSTRAINT_NAME
	where
		rc.CONSTRAINT_SCHEMA = @schema
`

type MssqlCrawler struct {
	db *sql.DB
}

func NewMssqlCrawler(db *sql.DB) *MssqlCrawler {
	return &MssqlCrawler{
		db: db,
	}
}

func mssqlCrawlTables(db *sql.DB, schemaName string) ([]Table, error) {
	rows, err := db.Query(mssqlTablesQuery, sql.Named("schema", schemaName))
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
			pk         sql.NullString
			fk         sql.NullString
		)

		err = rows.Scan(&tableName, &columnName, &dataType, &nullable, &pk, &fk)
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
			PrimaryKey: pk.Valid && pk.String == "YES",
			ForeignKey: fk.Valid && fk.String == "YES",
		})
	}

	if table != nil {
		tables = append(tables, *table)
	}

	return tables, nil
}

func mssqlCrawlRelations(db *sql.DB, schemaName string) ([]Relation, error) {
	rows, err := db.Query(mssqlRelationsQuery, sql.Named("schema", schemaName))
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

func (p *MssqlCrawler) Crawl(schemaName string) (*DatabaseSchema, error) {
	if p.db == nil {
		return nil, errors.New("database is nil")
	}

	if err := p.db.Ping(); err != nil {
		return nil, err
	}

	tables, err := mssqlCrawlTables(p.db, schemaName)
	if err != nil {
		return nil, err
	}

	relations, err := mssqlCrawlRelations(p.db, schemaName)
	if err != nil {
		return nil, err
	}

	return &DatabaseSchema{
		Tables:    tables,
		Relations: relations,
	}, nil
}
