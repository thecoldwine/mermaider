package internal

import (
	"os"
	"testing"
)

type mockCrawler struct {
}

func (m *mockCrawler) Crawl(schemaName string) (*DatabaseSchema, error) {
	return &DatabaseSchema{
		Tables: []Table{
			Table{
				Name: "table1",
				Columns: []Column{
					Column{
						Name:       "pk",
						PrimaryKey: true,
						Datatype:   "character varying",
						Nullable:   false,
						FK:         false,
					},
				},
			},
		},
		Relations: []Relation{},
	}, nil
}

func TestMermaider(t *testing.T) {
	m := &mockCrawler{}
	err := Mermaid(m, "public", os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
