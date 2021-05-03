package gsqlg

import "fmt"

type Select struct {
	BaseSource Source
	Joins      []Join
	Columns    []SelectColumn
}

type Source interface {
	Definition() string
	Alias() string
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

type SelectColumn interface {
	Definition() string
	Alias() string
}

type RawSelectColumn struct {
	Raw      string
	RawAlias string
}

func (r RawSelectColumn) Definition() string {
	return r.Raw
}

func (r RawSelectColumn) Alias() string {
	return r.RawAlias
}

var (
	_ SelectColumn = SourceColumn{}
)

type SourceColumn struct {
	SourceRef Source
	Name      string
}

func (s SourceColumn) Definition() string {
	return fmt.Sprintf("%s.%s", s.SourceRef.Alias(), s.Name)
}

func (s SourceColumn) Alias() string {
	return fmt.Sprintf("%s_%s", s.SourceRef.Alias(), s.Name)
}

type JoinType int

const (
	JoinTypeInner JoinType = iota
	JoinTypeLeft
	JoinTypeRight
	JoinTypeCross
)

var (
	_ Source = &Table{}
)

type Table struct {
	TableName  string
	TableAlias string
}

func (t *Table) Definition() string {
	return t.TableName
}

func (t *Table) Alias() string {
	return t.TableAlias
}

func (t *Table) Column(name string) SourceColumn {
	return SourceColumn{
		SourceRef: t,
		Name:      name,
	}
}
