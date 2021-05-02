package gsqlg

type Select struct {
	BaseSource SourceDef
	Joins      []Join
	Columns    []SelectColumn
}

type Join struct {
	JoinType    JoinType
	Source      SourceDef
	Constraints []JoinConstraint
}

type JoinConstraint struct {
	RefColumn    SourceColumn
	TargetColumn SourceColumn
}

type SelectColumn interface {
	SourceRef() SourceRef
	ColumnName() string
}

var (
	_ SelectColumn = SourceColumn{}
)

type SourceColumn struct {
	Source SourceRef
	Name   string
}

func (s SourceColumn) SourceRef() SourceRef {
	return s.Source
}

func (s SourceColumn) ColumnName() string {
	return s.Name
}

func (s SourceColumn) WithAlias(alias string) SourceColumnWithAlias {
	return SourceColumnWithAlias{
		SourceColumn: s,
		ColumnAlias: alias,
	}
}

var (
	_ Aliaser = SourceColumnWithAlias{}
)

type SourceColumnWithAlias struct {
	SourceColumn
	ColumnAlias string
}

func (s SourceColumnWithAlias) Alias() string {
	return s.ColumnAlias
}

type JoinType int

const (
	JoinTypeInner JoinType = iota
	JoinTypeLeft
	JoinTypeRight
	JoinTypeCross
)

type SourceRef interface {
	Ref() string
}

type SourceDef interface {
	Def() string
}

type Aliaser interface {
	Alias() string
}

var (
	_ SourceRef = &Table{}
	_ SourceDef = &Table{}
)

type Table struct {
	TableName string
}

func (t *Table) Ref() string {
	return t.TableName
}

func (t *Table) Def() string {
	return t.TableName
}

func (t *Table) Column(name string) SourceColumn {
	return SourceColumn{
		Source: t,
		Name:   name,
	}
}

var (
	_ Aliaser = &TableWithAlias{}
)

type TableWithAlias struct {
	Table      Table
	TableAlias string
}

func (t *TableWithAlias) Ref() string {
	return t.TableAlias
}

func (t *TableWithAlias) Def() string {
	return t.Table.TableName
}

func (t *TableWithAlias) Alias() string {
	return t.TableAlias
}

func (t *TableWithAlias) Column(name string) SourceColumn {
	return SourceColumn{
		Source: t,
		Name:   name,
	}
}
