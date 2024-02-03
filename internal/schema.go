package internal

type DatabaseSchema struct {
	Tables    []Table
	Relations []Relation
}

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name       string
	Datatype   string
	Nullable   bool
	PrimaryKey bool
	ForeignKey bool
}

type RelationType int8

// it is enum-like, NOT BIT FLAGS
const (
	OneToOne   RelationType = 0
	OneToMany  RelationType = 1
	ManyToOne  RelationType = 2
	ManyToMany RelationType = 3
)

type Relation struct {
	SourceTable       string
	SourceColumn      string
	DestinationTable  string
	DestinationColumn string
	RelType           RelationType
}
