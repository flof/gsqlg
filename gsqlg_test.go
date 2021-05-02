package gsqlg

import (
	"log"
	"testing"
)

func TestModel(t *testing.T) {
	var (
		userTable           = Table{TableName: "users"}
		groupTable          = Table{TableName: "groups"}
		primaryGroupTable   = TableWithAlias{Table: groupTable, TableAlias: "primary_group"}
		secondaryGroupTable = TableWithAlias{Table: groupTable, TableAlias: "secondary_group"}
	)

	s := Select{
		BaseSource: &userTable,
		Joins: []Join{
			{
				JoinType: JoinTypeLeft,
				Source:   &primaryGroupTable,
				Constraints: []JoinConstraint{
					{
						RefColumn:    SourceColumn{Source: &primaryGroupTable, Name: "id"},
						TargetColumn: SourceColumn{Source: &userTable, Name: "primary_group_id"},
					},
				},
			},
			{
				JoinType: JoinTypeLeft,
				Source:   &secondaryGroupTable,
				Constraints: []JoinConstraint{
					{
						RefColumn:    SourceColumn{Source: &secondaryGroupTable, Name: "id"},
						TargetColumn: SourceColumn{Source: &userTable, Name: "secondary_group_id"},
					},
				},
			},
		},
		Columns: []SelectColumn{
			userTable.Column("id"),
			userTable.Column("first_name"),
			userTable.Column("last_name"),
			primaryGroupTable.Column("name").WithAlias("primary_group_name"),
			secondaryGroupTable.Column("name").WithAlias("secondary_group_name"),
		},
	}

	gen := StatementGenerator{SelectColumnAliaser: DefaultSelectColumnAliaser{}}

	stmt := gen.Select(s)
	log.Print(stmt.Text)
}
