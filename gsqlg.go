package gsqlg

type Source interface {
	Aliaser
}

type Projection interface{}

type Aliaser interface {
	GetAlias() string
}

var (
	_ Source     = &Table{}
	_ Source     = &SubQuery{}
	_ Projection = SourceColumn{}
	_ Projection = RawSelectColumn{}
)

type Select struct {
	BaseSource  Source
	Joins       []Join
	Projections []Projection
}

type Join struct {
	JoinType    JoinType
	Source      Source
	Constraints []JoinConstraint
}

type JoinConstraint struct {
	RefColumn    SourceColumn
	TargetColumn SourceColumn
}

type RawSelectColumn struct {
	Raw   string
	Alias string
}

type SourceColumn struct {
	Source Source
	Name   string
}

type JoinType int

const (
	JoinTypeInner JoinType = iota
	JoinTypeLeft
	JoinTypeRight
	JoinTypeCross
)

type Table struct {
	Name  string
	Alias string
}

func (t *Table) Column(name string) SourceColumn {
	return SourceColumn{
		Source: t,
		Name:   name,
	}
}

func (t *Table) GetAlias() string {
	return t.Alias
}

type SubQuery struct {
	Select Select
	Alias  string
}

func (s *SubQuery) Column(name string) SourceColumn {
	return SourceColumn{
		Source: s,
		Name:   name,
	}
}

func (s *SubQuery) GetAlias() string {
	return s.Alias
}
