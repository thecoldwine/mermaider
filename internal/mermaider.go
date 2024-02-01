package internal

import (
	"bufio"
	"database/sql"
	"io"
	"strings"
)

type SchemaCrawler interface {
	Crawl(db *sql.DB) (*DatabaseSchema, error)
}

// yes, it is a verb here :3
// we're always using \n as a line separator because we can (and because go standard library does the same)
func Mermaid(db *sql.DB, crawler SchemaCrawler, w io.Writer) error {
	schema, err := crawler.Crawl(db)

	if err != nil {
		return err
	}

	prefix := strings.Repeat(" ", 4)

	b := bufio.NewWriter(w)
	b.WriteString("erDiagram\n")

	for _, t := range schema.Tables {

		b.WriteString(prefix + t.Name + "{\n")

		for _, c := range t.Columns {
			str := prefix + prefix + c.Name + " " + c.Datatype + "\n"
			b.WriteString(str)
		}

		b.WriteString(prefix + "}\n")
	}

	return nil
}
