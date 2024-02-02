package internal

import (
	"bufio"
	"io"
	"strings"
)

type SchemaCrawler interface {
	Crawl(schemaName string) (*DatabaseSchema, error)
}

// yes, it is a verb here :3
// we're always using \n as a line separator because we can (and because go standard library does the same)
func Mermaid(crawler SchemaCrawler, schemaName string, w io.Writer) error {
	schema, err := crawler.Crawl(schemaName)
	if err != nil {
		return err
	}

	prefix := strings.Repeat(" ", 4)

	b := bufio.NewWriter(w)
	defer b.Flush()

	b.WriteString("erDiagram\n")

	for _, t := range schema.Tables {

		b.WriteString(prefix + t.Name + " {\n")

		for _, c := range t.Columns {
			str := prefix + prefix + c.Name + " " + strings.ReplaceAll(c.Datatype, " ", "_") + "\n"
			b.WriteString(str)
		}

		b.WriteString(prefix + "}\n")

		b.Flush()
	}

	return nil
}
