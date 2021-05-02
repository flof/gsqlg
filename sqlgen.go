package gsqlg

import "strings"

type StatementGenerator struct {
	SelectColumnAliaser SelectColumnAliaser
}

type SelectColumnAliaser interface {
	Alias(s SelectColumn) string
}

type DefaultSelectColumnAliaser struct {
}

func (d DefaultSelectColumnAliaser) Alias(s SelectColumn) string {
	return s.SourceRef().Ref() + "_" + s.ColumnName()
}

type Statement struct {
	Text string
	Args []interface{}
}

func (s StatementGenerator) Select(sel Select) Statement {
	var stmt strings.Builder
	var args []interface{}

	stmt.WriteString("select ")

	for idx, col := range sel.Columns {
		if idx > 0 {
			stmt.WriteString(", ")
		}

		stmt.WriteString(col.SourceRef().Ref())
		stmt.WriteString(".")
		stmt.WriteString(col.ColumnName())

		alias, ok := col.(Aliaser)
		if ok {
			stmt.WriteString(" as ")
			stmt.WriteString(alias.Alias())
		} else {
			if s.SelectColumnAliaser != nil {
				stmt.WriteString(" as ")
				stmt.WriteString(s.SelectColumnAliaser.Alias(col))
			}
		}
	}

	stmt.WriteString(" from ")
	stmt.WriteString(sel.BaseSource.Def())
	alias, ok := sel.BaseSource.(Aliaser)
	if ok {
		stmt.WriteString(" as ")
		stmt.WriteString(alias.Alias())
	}

	for _, join := range sel.Joins {
		stmt.WriteString(" ")
		switch join.JoinType {
		case JoinTypeInner:
			stmt.WriteString("inner")
		case JoinTypeLeft:
			stmt.WriteString("left")
		case JoinTypeRight:
			stmt.WriteString("right")
		case JoinTypeCross:
			stmt.WriteString("cross")
		}

		stmt.WriteString(" join ")
		stmt.WriteString(join.Source.Def())
		alias, ok := join.Source.(Aliaser)
		if ok {
			stmt.WriteString(" as ")
			stmt.WriteString(alias.Alias())
		}

		for idx, joinCon := range join.Constraints {
			if idx == 0 {
				stmt.WriteString(" on ")
			} else {
				stmt.WriteString(" and ")
			}

			stmt.WriteString(joinCon.RefColumn.SourceRef().Ref())
			stmt.WriteString(".")
			stmt.WriteString(joinCon.RefColumn.Name)
			stmt.WriteString(" = ")
			stmt.WriteString(joinCon.TargetColumn.SourceRef().Ref())
			stmt.WriteString(".")
			stmt.WriteString(joinCon.TargetColumn.Name)
		}
	}

	return Statement{
		Text: stmt.String(),
		Args: args,
	}
}
