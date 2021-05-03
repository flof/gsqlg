package gsqlg

import "strings"

type StatementGenerator struct {
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

		stmt.WriteString(col.Definition())
		stmt.WriteString(" as ")
		stmt.WriteString(col.Alias())
	}

	stmt.WriteString(" from ")
	stmt.WriteString(sel.BaseSource.Definition())
	stmt.WriteString(" ")
	stmt.WriteString(sel.BaseSource.Alias())

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

		stmt.WriteString(join.Source.Definition())
		stmt.WriteString(" ")
		stmt.WriteString(join.Source.Alias())

		for idx, joinCon := range join.Constraints {
			if idx == 0 {
				stmt.WriteString(" on ")
			} else {
				stmt.WriteString(" and ")
			}

			stmt.WriteString(joinCon.RefColumn.SourceRef.Alias())
			stmt.WriteString(".")
			stmt.WriteString(joinCon.RefColumn.Name)
			stmt.WriteString(" = ")
			stmt.WriteString(joinCon.TargetColumn.SourceRef.Alias())
			stmt.WriteString(".")
			stmt.WriteString(joinCon.TargetColumn.Name)
		}
	}

	stmt.WriteString(";")

	return Statement{
		Text: stmt.String(),
		Args: args,
	}
}
