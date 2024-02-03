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
			{
				Name: "table1",
				Columns: []Column{
					{
						Name:       "pkey",
						PrimaryKey: true,
						Datatype:   "character varying",
						Nullable:   false,
						ForeignKey: false,
					},
				},
			},
			{
				Name: "table2",
				Columns: []Column{
					{
						Name:       "pkey",
						PrimaryKey: true,
						Datatype:   "character varying",
						Nullable:   false,
						ForeignKey: false,
					},
					{
						Name:       "pkey2",
						PrimaryKey: true,
						Datatype:   "character varying",
						Nullable:   false,
						ForeignKey: true,
					},
					{
						Name:       "attr",
						PrimaryKey: false,
						Datatype:   "character varying",
						Nullable:   true,
						ForeignKey: false,
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
