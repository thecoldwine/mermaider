package internal

import (
	"io"
	"strings"
	"text/template"
)

const mermaidTemplate = `
erDiagram
{{ range .Tables }}
  {{ .Name }} {
   {{ range .Columns }} {{ .Name }} {{ escape .Datatype }} {{ renderKeys . }} {{ renderNullability . }}
   {{end}}
  }
{{ end }}{{ range .Relations }}
  {{ .DestinationTable }} ||--o{ {{ .SourceTable }} : ""{{ end }}
`

func escapeSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}

func renderKeys(c Column) string {
	if c.PrimaryKey && c.ForeignKey {
		return "PK,FK"
	}

	if c.PrimaryKey {
		return "PK"
	}

	if c.ForeignKey {
		return "FK"
	}

	return ""
}

func renderNullability(c Column) string {
	if c.Nullable {
		return ""
	}

	return "\"not null\""
}

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

	t, err := template.New("m1").Funcs(template.FuncMap{
		"escape":            escapeSpaces,
		"renderKeys":        renderKeys,
		"renderNullability": renderNullability,
	}).Parse(mermaidTemplate)
	if err != nil {
		return err
	}

	err = t.Execute(w, schema)
	if err != nil {
		return err
	}

	return nil
}
