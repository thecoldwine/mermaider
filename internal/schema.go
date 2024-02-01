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
	FK         bool
}

type RelationType int8

// it is enum-like, NOT BIT FLAGS
const (
	OneToOne   RelationType = 0
	OneToMany  RelationType = 1
	ManyToMany RelationType = 2
)

type Relation struct {
	Source      string
	Destination string
	RelType     RelationType
}
